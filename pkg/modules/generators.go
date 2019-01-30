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

package modules

// todo this should be made a configuration option, distribution of templates are also "unsolved"
var TmplBasePath = "/usr/share/metagraf/templates"

// Slice of strings that we wash labels from baseimage with when
// creating a deploymentconfig
var LabelBlacklistFilter []string = []string{
	"openshift",
	"s2i",
	"license",
	"k8s",
}

// Slice of strings that we wash environment variables with from
// the base image when creating a deploymentconfig
var EnvBlacklistFilter []string = []string{
	"path",
	"home",
	"bash",
	"env",
	"sti",
	"openshift",
}
