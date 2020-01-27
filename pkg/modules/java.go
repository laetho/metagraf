/*
Copyright 2020 The MetaGraph Authors

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
	"metagraf/pkg/metagraf"
	"strings"
	corev1 "k8s.io/api/core/v1"
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
func GenEnvVar_JVM_SYS_PROP(mg *metagraf.MetaGraf, name string) corev1.EnvVar {
	var props []string

	for _,c := range mg.Spec.Config {
		if strings.ToUpper(c.Type) != "JVM_SYS_PROP" {
			continue
		}
		for _,o := range c.Options {
			prop := Variables[c.Name+"|"+o.Name]
			str := ""
			if o.Required {
				// Set default value if --defaults arg is given.
				if Defaults {
					str = "-D" + o.Name + "=" + o.Default
				}
				// Set value of key from either --cvars of --cvfile or Environment
				if len(prop.Value) > 0  {
					str = "-D" + o.Name + "=" + prop.Value
				}
				// Append generated string to props, could be empty if non of the above
				// logic kicks in.
				props = append(props, str)
			}
			// If we find an optional value in --cvars og --cvfile or Environment
			if !o.Required {
				if len(prop.Value) > 0 {
					str = "-D" + o.Name + "=" + prop.Value
				}
				props = append(props, str)
			}
		}
	}
	return corev1.EnvVar{
		Name: name,
		Value: strings.Join(props, " "),
	}
}