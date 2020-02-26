/*
Copyright 2019 The MetaGraph Authors

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
	log "k8s.io/klog"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	"os"
	"strings"
)

func PropertiesFromEnv(mgp metagraf.MGProperties) {
	for _, v := range os.Environ() {
		key, val := keyValueFromEnv(v)
		if p, ok := mgp["local:"+key]; ok {
			p.Value = val
			mgp[p.MGKey()] = p
		}
	}
}

// Modifies a MGProperties map with values from --cvars
// argument. Only supports local environment variables
// for now.
func PropertiesFromCmd(mgp metagraf.MGProperties) {
	// Parse and get values from --cvars
	cvars := CmdCVars(CVars).Parse()

	for k, v := range cvars {
		if p, ok := mgp["local|"+k]; ok {
			p.Value = v
			mgp[p.MGKey()] = p
			continue
		}
		if p, ok := mgp["JVM_SYS_PROP|"+k]; ok {
			p.Value = v
			mgp[p.MGKey()] = p
			continue
		}
	}
}

// Used for splitting --cvfile .properties files with strings.FieldsFunc()
func MgPropertyLineSplit(r rune) bool {
	return r == '|' || r == '='
}

// Modifies MGProfperties map with information from files
// and also return a map of only the properties on file..
func PropertiesFromFile(mgp metagraf.MGProperties) metagraf.MGProperties {
	props := metagraf.MGProperties{}

	if len(CVfile) == 0 {
		return props
	}

	file, err := os.Open( CVfile )
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	fail := false
	var line string
	for {
		line,err = reader.ReadString('\n')
		if err != nil {
			break
		}

		// Skip emptyish lines, there will never be a line shorter than 3 characters.
		if len(line) <= 3 {
			continue
		}

		a := strings.FieldsFunc(line, MgPropertyLineSplit)
		var v []string	// to hold soruce, key, value
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
			v = append(v, a[0]) 	// Set Source
			v = append(v, a[1])		// Set Key
			v = append(v, val[1])	// Set value
		} else {
			v =  a 	// Set Source, Key, Value from simple line parse
		}

		if len(v) != 3 {
			fmt.Println("Configuration format error! Is it in: \"<source>|key=value\"")
			os.Exit(1)
		}

		t := metagraf.MGProperty{
			Source:   v[0],
			Key:      v[1],
			Value:    strings.TrimRight(v[2], "\n"),
			Required: true,
			Default:  "",
		}
		t.Default = mgp[t.MGKey()].Default
		t.Required = mgp[t.MGKey()].Required

		// Only set in mgp MGProperties if the key is valid.
		if _, ok := mgp[t.MGKey()]; !ok {
			log.V(1).Infof("Found invalid key: %s while reading configuration file.\n", t.MGKey())
		} else {
			if len(t.Value) == 0 {
				fail = true
				fmt.Printf("Configured property %v must have a value in %v\n", t.MGKey(),CVfile )
			}
			mgp[t.MGKey()] = t
		}

		props[t.MGKey()] = t
	}
	if fail {
		os.Exit(1)
	}
	return props
}


// Splits a shell environment variable into key and value parts
// and returns them seperatly.
func keyValueFromEnv(s string) (string, string) {
	return strings.Split(s, "=")[0], strings.Split(s, "=")[1]
}

func OverrideProperties(mgp metagraf.MGProperties) {
	// Fetch possible variables form metaGraf specification
	PropertiesFromEnv(mgp)
	// Fetch variable overrides from file if specified with --cvfile
	PropertiesFromFile(mgp)
	// Fetch from commandline
	PropertiesFromCmd(mgp)
}

func FlagPassingHack() {
	if Dryrun {
		Output = true
	}
	// Push flags to modules (hack)
	modules.Version = Version
	modules.Output = Output
	modules.Dryrun = Dryrun
	modules.NameSpace = Namespace
	modules.Verbose = Verbose
	modules.CVfile = CVfile
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


