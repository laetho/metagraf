/*
Copyright 2018 The MetaGraph Authors

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
	"github.com/golang/glog"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	"os"
	"strings"
)

// Returns a slice of strings of potential parameterized variables in a
// metaGraf specification that can be found in the execution environment.
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

func VarsFromCmd(mgv metagraf.MGVars, cvars CmdVars) map[string]string {
	vars := make(map[string]string)

	for k, v := range cvars {
		if _, ok := mgv[k]; ok {
			vars[k] = v
		}
	}
	return vars
}

// Reads input from properties file.
func VarsFromFile(mgv metagraf.MGVars) map[string]string {
	vars := make(map[string]string)
	if len(CVfile) == 0 {
		return vars
	}

	file, err := os.Open( CVfile )
	if err != nil {
		glog.Error(err)
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
		vl := strings.Split(line,"=")
		if len(vl) < 2 {
			glog.Errorf("Properties are formatted improperly in: %v", CVfile)
			break
		}
		if len(vl) == 2 {
			vars[vl[1]] = ""
		} else if len(vl) >= 4 {
			vars[vl[1]] = strings.ReplaceAll(strings.Join(vl[2:], "="), "\n", "")
		} else {
			vars[vl[1]] = strings.ReplaceAll(vl[2], "\n", "")
		}
	}
	return vars
}


func keyValueFromEnv(s string) (string, string) {
	return strings.Split(s, "=")[0], strings.Split(s, "=")[1]
}

// Returns a list of variables from command line or environment where
// command line is the most significant.
// Precedence is Environment, File and Command
func OverrideVars(mgv metagraf.MGVars, cvars CmdVars) map[string]string {
	ovars := make(map[string]string)

	// Fetch possible variables form metaGraf specification
	for k, v := range VarsFromEnv(mgv) {
		ovars[k] = v
	}

	// Fetch variable overrides from file if specified with --cvfile
	for k,v := range VarsFromFile(mgv) {
		ovars[k] = v
	}

	// Fetch from commandline
	for k, v := range VarsFromCmd(mgv, cvars) {
		ovars[k] = v
	}

	return ovars
}

func MergeVars(base metagraf.MGVars, override map[string]string) metagraf.MGVars {
	for k, v := range override {
		base[k] = v
	}
	glog.Info("Calling MergeVars: ", base)
	return base
}

func MergeSourceVars(base metagraf.MGVars, override map[string]string) metagraf.MGVars {
	glog.Info("MergeSourcevars(): base mg vars", base)
	glog.Info("MergeSourceVars(): with override values: ", override)

	// Translation map, key = untyped key, value typed key name
	keys := make(map[string]string)

	// Strip Source Label from base
	for k,_ := range base {
		key := strings.Split(k, "=")[1]
		keys[key] = k
	}

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


