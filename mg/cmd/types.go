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
	"strings"
)

type CmdMessage string

var StrActiveProject		CmdMessage = "Active project is:"
var StrMissingMetaGraf		CmdMessage = "Missing path to metaGraf specification."
var StrMissingCollection	CmdMessage = "Missing path to collection of metaGraf specifications."
var StrMissingNamespace 	CmdMessage = "Namespace must be supplied or configured."
var StrMalformedVar			CmdMessage = "Malformed key=value pair supplied through --cvars :"

// Type for mg custom variables
type EnvVars		map[string]string
type MGVars			map[string]string
type CmdCVars		[]string
type CmdVars		map[string]string

var CVars 	[]string

// Returns a map (CmdVars) parsed from --cvars flag
func (v CmdCVars) Parse() CmdVars {
	cm := make(CmdVars)
	for _,str := range v {
		split := strings.Split(str,"=")
		if len(split) < 1 {
			glog.Info(StrMalformedVar)
			continue
		}
		cm[split[0]]=split[1]
	}
	return cm
}

