/*
Copyright 2018 The MetaGraf Authors

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

package metagraf

// JSON structure for a MetaGraf entity
type MetaGraf struct {
	Kind     string
	Metadata struct {
		Name              string
		ResourceVersion   string
		Namespace         string
		CreationTimestamp string
		Labels            map[string]string
		Annotations       map[string]string
	}
	Spec struct {
		Type 		string
		Version		string
		Description	string
		Resources   []Resource
		Environment struct {
			Local    []EnvironmentVar
			External struct {
				Introduces []EnvironmentVar
				Consumes   []EnvironmentVar
			}
		}
	}
	Config []struct{
		FileName 	string
		Type 		string
		Description string
		Options 	[]ConfigParam
	}
}

type Resource struct {
	Name     	string
	Type     	string
	Semop		string
	Semver  	string
	Required 	bool
}
type ConfigParam struct {
	Name        string
	Required    bool
	Description string
	Type        string
	Default     string
}
type EnvironmentVar struct {
	Name        string
	Required    bool
	Type        string
	Description string
	Default		string
}
