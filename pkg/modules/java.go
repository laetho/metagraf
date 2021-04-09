/*
Copyright 2020 The metaGraf Authors

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
	corev1 "k8s.io/api/core/v1"
	"github.com/laetho/metagraf/pkg/metagraf"
	"strings"
)

// Identify if a metagraf specification has a JVM_SYS_PROP configuration type.
// This is a special type for handling system properties in the java space that
// your application might not know about.
func HasJVM_SYS_PROP(mg *metagraf.MetaGraf) bool {
	for _,c := range mg.Spec.Config {
		if strings.ToUpper(c.Type) == "JVM_SYS_PROP" {
			return true
		}
	}
	return false
}

// Generate an EnvVar for a config section.
// SecretFrom or EnvFrom will not be processed.
func GenEnvVar_JVM_SYS_PROP(mgp metagraf.MGProperties, name string) corev1.EnvVar {
	var props []string
	for _, p := range mgp {
		// Skip processing if not JVM_SYS_PROP
		if p.Source != "JVM_SYS_PROP" {
			continue
		}

		if Defaults {
			props = append(props, "-D"+p.Key+"="+p.Default)
		} else {
			props = append(props, "-D"+p.Key+"="+p.Value)
		}
	}
	return corev1.EnvVar{
		Name: name,
		Value: strings.Join(props, " "),
	}
}