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

func keyValueFromEnv(s string) (string, string) {
	return strings.Split(s, "=")[0], strings.Split(s, "=")[1]
}

// Returns a list of variables from command line or environment where
// command line is the most significant.
func OverrideVars(mgv metagraf.MGVars, cvars CmdVars) map[string]string {
	ovars := make(map[string]string)

	// Fetch possible variables form metaGraf specification
	for k, v := range VarsFromEnv(mgv) {
		ovars[k] = v
	}
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
}
