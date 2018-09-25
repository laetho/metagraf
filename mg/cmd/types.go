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

type CmdMessage string

var StrActiveProject		CmdMessage = "Active project is:"
var StrMissingMetaGraf		CmdMessage = "Missing path to metaGraf specification."
var StrMissingCollection	CmdMessage = "Missing path to collection of metaGraf specifications."
var StrMissingNamespace 	CmdMessage = "Namespace must be supplied or configured."
