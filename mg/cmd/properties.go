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
	"github.com/laetho/metagraf/internal/pkg/params/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	log "k8s.io/klog"
	"os"
	"strings"
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
			continue
		}
	}
	return mgp
}

// Used for splitting --cvfile .properties files with strings.FieldsFunc()
func MgPropertyLineSplit(r rune) bool {
	return r == '|' || r == '='
}

// Reads a properties file and returns a metagraf.MGProperties structure
func ReadPropertiesFile(propfile string) metagraf.MGProperties {
	props := metagraf.MGProperties{}

	if len(propfile) == 0 {
		return props
	}

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
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var line string
	for scanner.Scan() {
		line = scanner.Text()
		// Skip emptyish lines, there will never be a line shorter than 3 characters.
		if len(line) <= 3 {
			continue
		}

		a := strings.FieldsFunc(line, MgPropertyLineSplit)
		var v []string // to hold soruce, key, value
		if len(a) >= 3 {
			// Handle multiple = in value.
			val := strings.SplitAfter(line, a[0]+"|"+a[1]+"=")
			if len(val) == 1 {
				val = strings.SplitAfter(line, a[0]+"="+a[1]+"=")
				if len(val) == 1 {
					fmt.Println("Unable to parse config format!")
					os.Exit(1)
				}
			}
			v = append(v, a[0])   // Set Source
			v = append(v, a[1])   // Set Key
			v = append(v, val[1]) // Set value
		} else {
			v = a // Set Source, Key, Value from simple line parse
		}

		if len(v) != 3 {
			fmt.Println("Configuration format error! Is it in: \"<source>|key=value\"")
			os.Exit(1)
		}

		t := metagraf.MGProperty{
			Source: v[0],
			Key:    v[1],
			Value:  strings.TrimRight(v[2], "\n"),
		}
		props[t.MGKey()] = t
	}
	return props
}

//
func MergeAndValidateProperties(base metagraf.MGProperties, merge metagraf.MGProperties, novalidate bool) metagraf.MGProperties {
	for _, p := range merge {
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
	fileprops := ReadPropertiesFile(params.PropertiesFile)
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
	modules.Verbose = Verbose
	modules.Defaults = Defaults
	modules.Format = Format
	modules.Suffix = Suffix
	modules.Template = Template
	modules.ImageNS = ImageNS
	modules.Enforce = Enforce
	modules.Registry = Registry
	modules.Tag = Tag
	modules.OName = OName
	modules.Context = Context
	modules.CreateGlobals = CreateGlobals
}
