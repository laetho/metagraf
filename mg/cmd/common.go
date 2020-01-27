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


// Builds and returns a EnvVars{} map of shell environment
// variables that matches addressable fields in a metaGraf
// specification.
func VarsFromEnv(mgv metagraf.MGVars) EnvVars {
	envs := EnvVars{}
	for _, v := range os.Environ() {
		key, val := keyValueFromEnv(v)
		if _, ok := mgv[key]; ok {
			envs[key] = val
		}
	}
	return envs
}

func PropertiesFromEnv(mgp metagraf.MGProperties) {
	for _, v := range os.Environ() {
		key, val := keyValueFromEnv(v)
		if p, ok := mgp["local:"+key]; ok {
			p.Value = val
			mgp[p.MGKey()] = p
		}
	}
}

// Returns a map of key, value pairs that matched addressable fields
// in a metaGraf specification from --cvars argument.
func VarsFromCmd(mgv metagraf.MGVars) map[string]string {
	vars := make(map[string]string)

	// Parse and get values from --cvars
	cvars := CmdCVars(CVars).Parse()

	for k, v := range cvars {
		if _, ok := mgv[k]; ok {
			vars[k] = v
		}
	}
	return vars
}

// Modifies a MGProperties map with values from --cvars
// argument. Only supports local environment variables
// for now.
func PropertiesFromCmd(mgp metagraf.MGProperties) {
	// Parse and get values from --cvars
	cvars := CmdCVars(CVars).Parse()

	for k, v := range cvars {
		if p, ok := mgp["local:"+k]; ok {
			p.Value = v
			mgp[p.MGKey()] = p
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

		// Skip empty lines
		if len(line) == 1 {
			if strings.Contains(line, "\n") { continue }
		}
		a := strings.FieldsFunc(line, MgPropertyLineSplit)
		if len(a) != 3 {
			fmt.Println("Configuration format error! Is it in: \"souce|key=value\"")
			os.Exit(1)
		}
		t := metagraf.MGProperty{
			Source:   a[0],
			Key:      a[1],
			Value:    strings.TrimRight(a[2], "\n"),
			Required: true,
			Default:  "",
		}
		t.Default = mgp[t.MGKey()].Default
		t.Required = mgp[t.MGKey()].Required

		if len(t.Value) == 0 {
			fail = true
			fmt.Printf("Configured property %v must have a value in %v\n", t.MGKey(),CVfile )
		}
		// Only set in mgp MGProperties if the key is valid.
		if _, ok := mgp[t.MGKey()]; ok {
			log.V(1).Infof("Found invalid key: %s while reading configuration file.\n", t.MGKey())
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


