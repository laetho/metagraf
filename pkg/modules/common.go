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
	corev1 "k8s.io/api/core/v1"
	"metagraf/pkg/metagraf"
	"strconv"
	"strings"
)

var (
	NameSpace string	// Used to pass namespace from cmd to module to avoid import cycle.
)

var Variables map[string]string

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
	if len(r.User) > 0 && len(r.Secret) == 0 {
		// Implicit secret name generation
		return strings.ToLower(r.Name+"-"+r.User)
	} else if len(r.User) == 0 && len(r.Secret) > 0 {
		// Explicit secret name generation
		return strings.ToLower(r.Name+"-"+r.Secret)
	} else {
		return strings.ToLower(r.Name)
	}
}

// Applies conventions and overridden logic to an environment variable and returns a corev1.EnvVar{}
func EnvToEnvVar(e *metagraf.EnvironmentVar) corev1.EnvVar {
	if e.Required == false {
		value := ""

		// Handle possible override value for non required fields
		if v, t := Variables[e.Name]; t {
			value = v
		}

		return corev1.EnvVar{
			Name: e.Name,
			Value: value,
		}
	}

	if len(e.Default) == 0 {e.Default = "$"+e.Name}
	// Handle possible overridden values for required fields
	if v, t := Variables[e.Name]; t {
		e.Default = v
	}

	return corev1.EnvVar{
		Name:  e.Name,
		Value: e.Default,
	}
}


func ValueFromEnv(key string) bool {
	if _, t := Variables[key]; t {
		return true
	}
	return false
}