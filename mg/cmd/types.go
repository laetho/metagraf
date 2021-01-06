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
	log "k8s.io/klog"
	"strings"
)

type CmdMessage string

var StrActiveProject CmdMessage = "Active project is:"
var StrMissingMetaGraf CmdMessage = "Missing path to metaGraf specification."
var StrMissingCollection CmdMessage = "Missing path to collection of metaGraf specifications."
var StrMissingNamespace CmdMessage = "Namespace must be supplied or configured."
var StrMalformedVar CmdMessage = "Malformed key=value pair supplied through --cvars :"

// Type for mg custom variables
type EnvVars map[string]string		// Map for holding addressable key, value pairs from os.Environ().
type CmdCVars []string				// Map
type CmdVars map[string]string


// Returns a map (CmdVars) parsed from --cvars flag
// todo: fix parsing of , seperated values for a key
func (v CmdCVars) Parse() CmdVars {
	cm := make(CmdVars)
	for _, str := range v {
		log.Info("Cmd Var string:", str)
		split := strings.Split(str, "=")
		if len(split) <= 1 {
			log.Info("Split:", split)
			log.Warning(StrMalformedVar)
			continue
		}
		cm[split[0]] = split[1]
	}
	log.V(2).Info("CmdCVars: ", cm)
	return cm
}
