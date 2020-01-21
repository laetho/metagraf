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

// Returns a map of key, value pairs that matched addressable fields
// in a metaGraf specification from a properties file generated by
// "mg generate properties" and provided through
// the --cvfile argument.
func VarsFromFile(mgv metagraf.MGVars) map[string]string {
	vars := make(map[string]string)
	if len(CVfile) == 0 {
		return vars
	}

	file, err := os.Open( CVfile )
	if err != nil {
		log.Error(err)
		return vars
	}
	defer file.Close()
	reader := bufio.NewReader(file)

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

		if strings.ContainsRune( line, 58 ) {
			fmt.Println("new handling")
		} else {
			// Old --cvfile format
			vl := strings.Split(line, "=")
			if len(vl) < 2 {
				log.Errorf("Properties are formatted improperly in: %s", CVfile)
				os.Exit(1)
			}
			if len(vl) == 2 {
				vars[vl[1]] = ""
			} else if len(vl) >= 4 {
				vars[vl[1]] = strings.ReplaceAll(strings.Join(vl[2:], "="), "\n", "")
			} else {
				vars[vl[1]] = strings.ReplaceAll(vl[2], "\n", "")
			}
		}
	}
	return vars
}

// Splits a shell environment variable into key and value parts
// and returns them seperatly.
func keyValueFromEnv(s string) (string, string) {
	return strings.Split(s, "=")[0], strings.Split(s, "=")[1]
}

// Returns a map of key, value pairs of addressable fields in a
// metaGraf specification from Environment, --cvfile argument and
// --cvars argument.
//
// Precedence is:
//   1. --cvars argument
//   2. --cvfile argument
//   3. Environment
func OverrideVars(keys metagraf.MGVars) map[string]string {
	ovars := make(map[string]string)

	// Fetch possible variables form metaGraf specification
	for k, v := range VarsFromEnv(keys) {
		ovars[k] = v
	}

	// Fetch variable overrides from file if specified with --cvfile
	for k,v := range VarsFromFile(keys) {
		ovars[k] = v
	}

	// Fetch from commandline
	for k, v := range VarsFromCmd(keys) {
		ovars[k] = v
	}

	return ovars
}

func MergeVars(base metagraf.MGVars, override map[string]string) metagraf.MGVars {
	log.Info("Calling MergeVars: ", base)
	for k, v := range override {
		base[k] = v
	}
	return base
}


// Used when parsing --cvfile
func MergeSourceKeyedVars(base metagraf.MGVars, override map[string]string) metagraf.MGVars {
	log.Info("MergeSourcevars(): base mg vars", base)
	log.Info("MergeSourceVars(): with override values: ", override)

	// Translation map, key = untyped key, value typed key name
	keys := make(map[string]string)

	// @todo Also support : as seperator to make the format of the .properties more readable.
	// Strip Source Label from base
	for k,_ := range base {
		key := strings.Split(k, "=")[1] // Stripping source reference.
		keys[key] = k
	}

	// Apply only valid variables from basevars.
	for k, v := range override {
		base[keys[k]] = v
	}
	return base
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


