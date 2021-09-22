/*
Copyright 2019 The metaGraf Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	log "k8s.io/klog"
)

func PropertiesFromEnv(mgp metagraf.MGProperties) metagraf.MGProperties {
	for _, v := range os.Environ() {
		key, val := keyValueFromEnv(v)
		if p, ok := mgp["local:"+key]; ok {
			p.Value = val
			mgp[p.MGKey()] = p
		}
	}
	return mgp
}

// Modifies a MGProperties map with values from --cvars
// argument. Only supports local environment variables
// for now.
func PropertiesFromCmd(mgp metagraf.MGProperties) metagraf.MGProperties {
	// Parse and get values from --cvars
	cvars := CmdCVars(CVars).Parse()
	for k, v := range cvars {
		if p, ok := mgp[k]; ok {
			p.Value = v
			mgp[p.MGKey()] = p
		} else {
			property := metagraf.MGProperty{
				Source:   "local",
				Key:      k,
				Value:    v,
				Required: false,
				Default:  "",
			}
			mgp[property.MGKey()] = property
		}
	}
	return mgp
}

// MgPropertyLineSplit Used for splitting --cvfile .properties files with strings.FieldsFunc()
func MgPropertyLineSplit(r rune) bool {
	return r == '|' || r == '='
}

// ReadPropertiesFromFile Reads a properties file and returns a metagraf.MGProperties structure
func ReadPropertiesFromFile(propfile string) metagraf.MGProperties {
	if len(propfile) > 0 {
		file, err := os.Open(propfile)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Warningf("Unable to close file: %v", err)
			}
		}()
		return ParseProps(file)
	} else {
		return metagraf.MGProperties{}
	}

}


func ParseProps(reader io.Reader) metagraf.MGProperties {
	mgProps := metagraf.MGProperties{}
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	var line string
	for scanner.Scan() {
		line = scanner.Text()
		// Skip emptyish lines, there will never be a line shorter than 3 characters.
		if len(line) <= 3 {
			continue
		}

		propertySlice := parsePropertyLine(line)

		mgProperty := metagraf.MGProperty{
			Source: propertySlice[0],
			Key:    propertySlice[1],
			Value:  strings.TrimRight(propertySlice[2], "\n"),
		}
		mgProps[mgProperty.MGKey()] = mgProperty
	}
	return mgProps
}

func parsePropertyLine(line string) []string {
	lineSplitSlice := strings.FieldsFunc(line, MgPropertyLineSplit)
	var propertySlice []string
	if splitByMultipleEqualSigns(lineSplitSlice) {
		propertySlice = concatenateValueSplitByEqualsSign(line, lineSplitSlice)
	} else {
		propertySlice = lineSplitSlice // Set Source, Key, Value from simple line parse
	}

	if len(propertySlice) != 3 {
		fmt.Println("Configuration format error! Is it in: \"<source>|key=value\"")
		os.Exit(1)
	}

	return propertySlice
}

// Occasionally properties will contain multiple equals sign,
// then we must concatenate the parts from the initial split
func concatenateValueSplitByEqualsSign(line string, lineSplitSlice []string) []string {
	val := strings.SplitAfter(line, lineSplitSlice[0]+"|"+lineSplitSlice[1]+"=")
	if len(val) == 1 {
		val = strings.SplitAfter(line, lineSplitSlice[0]+"="+lineSplitSlice[1]+"=")
		if len(val) == 1 {
			fmt.Println("Unable to parse config format!")
			os.Exit(1)
		}
	}
	var propertySlice []string
	propertySlice = append(propertySlice, lineSplitSlice[0]) // Set Source
	propertySlice = append(propertySlice, lineSplitSlice[1]) // Set Key
	propertySlice = append(propertySlice, val[1])            // Set value
	return propertySlice
}

func splitByMultipleEqualSigns(lineSplitSlice []string) bool {
	return len(lineSplitSlice) >= 3
}

//
func MergeAndValidateProperties(base metagraf.MGProperties, merge metagraf.MGProperties, novalidate bool) metagraf.MGProperties {
	for _, p := range merge {

		// Do not allow setting values on a sticky key.
		// Sticky keys are values fetched from secrets or configmaps.
		if p.Source == "local" && novalidate {
			if _, ok := base["sticky|"+p.Key]; ok {
				log.Fatalf("You tried to set a custom value on a secretfrom og valuefrom property. This is not allowed! Check you properties file on key: %v=%v", p.MGKey(), p.Value)
			}
		}

		if novalidate {
			if _, ok := base[p.MGKey()]; !ok {
				base[p.MGKey()] = p
			} else {
				property := base[p.MGKey()]
				property.Value = p.Value
				base[p.MGKey()] = property
			}
			continue
		}
		// Allow for setting Kubernetes service discovery environment variables
		// even though they are not part of of the metagraf specification.
		if strings.Contains(p.Key, "_SERVICE_") && !novalidate {
			base[p.MGKey()] = p
		}

		// Only set in base MGProperties if the key is valid.
		if _, ok := base[p.MGKey()]; !ok {
			if len(p.Value) == 0 && p.Required {
				fmt.Printf("Configured property %v must have a value in %v\n", p.MGKey(), params.PropertiesFile)
			} else {
				target := base[p.MGKey()]
				target.Value = p.Value
				base[target.MGKey()] = target
			}
		}
	}
	return base
}

// Splits a shell environment variable into key and value parts
// and returns them seperatly.
func keyValueFromEnv(s string) (string, string) {
	return strings.Split(s, "=")[0], strings.Split(s, "=")[1]
}

// Process cmd parameters based on metagraf defined properties.
func GetCmdProperties(mgp metagraf.MGProperties) metagraf.MGProperties {
	if Defaults {
		for _, property := range mgp {
			property.DefaultValueAsValue()
			mgp[property.MGKey()] = property
		}
	}
	fileprops := ReadPropertiesFromFile(params.PropertiesFile)
	// Fetch possible variables from metaGraf specification
	mgp = MergeAndValidateProperties(mgp, PropertiesFromEnv(mgp), false)
	// Fetch variable overrides from file if specified with --cvfile
	mgp = MergeAndValidateProperties(mgp, fileprops, true)
	// Fetch from commandline with --cvars
	mgp = MergeAndValidateProperties(mgp, PropertiesFromCmd(mgp), false)
	return mgp
}

func FlagPassingHack() {
	if Dryrun {
		Output = true
	}
	// Push flags to modules (hack) All of these should be replaced by using params. package
	modules.Version = Version
	modules.Output = Output
	modules.Dryrun = Dryrun
	modules.NameSpace = Namespace
	modules.Defaults = Defaults
	modules.Format = Format
	modules.Suffix = Suffix
	modules.Template = Template
	modules.ImageNS = ImageNS
	modules.Registry = Registry
	modules.Tag = Tag
	modules.OName = OName
	modules.Context = Context
	modules.CreateGlobals = CreateGlobals
}
