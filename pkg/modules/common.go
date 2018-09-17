/*
Copyright 2018 The metaGraf Authors

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

import (
	"github.com/blang/semver"
	"metagraf/pkg/metagraf"
	"strconv"
	"strings"
)

var (
	NameSpace string	// Used to pass namespace from cmd to module to avoid import cycle.
)

// Returns a name for a resource based on convention as follows.
func Name(mg *metagraf.MetaGraf) string {
	var objname string
	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		objname = strings.ToLower(mg.Metadata.Name)
	} else {
		objname = strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10))
	}
	return objname
}

// Returns a name for a secret for a resource based on convention as follows.
func ResourceSecretName(r *metagraf.Resource) string {
	if len(r.User) > 0 && len(r.SecretRef) == 0 {
		// Implicit secret name generation
		return r.Name+"-"+r.User
	} else if len(r.User) == 0 && len(r.SecretRef) > 0 {
		// Explicit secret name generation
		return r.Name+"-"+r.SecretRef
	} else {
		return r.Name
	}
}