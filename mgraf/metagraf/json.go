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

package metagraf

// JSON structure for a MetaGraf entity
type MetaGraf struct {
	Kind     string
	Version  string
	Metadata struct {
		Name              string
		Version           string
		ResourceVersion   string
		Namespace         string
		CreationTimestamp string
		Description       string
		Labels            map[string]string
		Annotations       map[string]string
	}
	Spec struct {
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
		FileName string
		Options []ConfigParam
	}
}

type Resource struct {
	Name     string
	Type     string
	Version  string
	Match    string
	Required bool
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
}
