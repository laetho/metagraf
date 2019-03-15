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

// Map to hold all variables from a specification
type MGVars			map[string]string

// JSON structure for a MetaGraf entity
type MetaGraf struct {
	Kind     string		`json:"kind"`
	Metadata struct {
		Name              string	`json:"name"`
		ResourceVersion   string	`json:"resourceversion"`
		Namespace         string	`json:"namespace"`
		CreationTimestamp string	`json:"creationtimestamp,omitempty"`
		Labels            map[string]string	`json:"labels,omitempty"`
		Annotations       map[string]string	`json:"annotations,omitempty"`
	} `json:"metadata"`
	Spec struct {
		Type         string		`json:"type"`
		Version      string		`json:"version"`
		Description  string		`json:"description"`
		Repository   string  	`json:"repository,omitempty"`
		RepSecRef	 string		`json:"repsecref,omitempty"`
		Branch 		 string		`json:"branch,omitempty"`
		BuildImage   string		`json:"buildimage,omitempty"`
		BaseRunImage string		`json:"baserunimage,omitempty"`
		Image		 string		`json:"image,omitempty"`
		Resources   []Resource	`json:"resources,omitempty"`
		Environment struct {
			Build []EnvironmentVar	`json:"build,omitempty"`
			Local []EnvironmentVar	`json:"local,omitempty"`
			External struct {
				Introduces []EnvironmentVar `json:"introduces,omitempty"`
				Consumes   []EnvironmentVar `json:"consumes,omitempty"`
			} `json:"external,omitempty"`
		} `json:"environment,omitempty"`
		Config []Config `json:"config,omitempty"`
	} `json:"spec"`
}


/*
 * Describes attached resources for a component. Ref, 12 factor app.
 */
type Resource struct {
	Name     	string			`json:"name"`
	Description string 			`json:"description,omitempty"`
	Type     	string			`json:"type"`
	Required 	bool			`json:"required"`
	External 	bool    		`json:"external"`
	Semop		string			`json:"semop,omitempty"`		// Semantic operator, how to evaluate version match/requirements.
	Semver  	string			`json:"semver,omitempty"`		// Semantic version to evaluate for attached resource

	Template	string  		`json:"template,omitempty"`		// Go txt template string for generating resource configuration.
	TemplateFrom string			`json:"envref,omitempty"`		//

	User 		string			`json:"user,omitempty"`
	UserFrom	string			`json:"userfrom,omitempty"`		// UserEnvRef overrides User, UserEnvRef must be in Environment->Local section.

	Secret		string			`json:"secret,omitempty"`		// Not Used?
	SecretType  string			`json:"secrettype,omitempty"`	// Not Used?
	SecretFrom	string			`json:"secretfrom,omitempty"`	// Secret reference for

}

type Config struct {
	Name    	string			`json:"name"`
	Type        string			`json:"type"`
	Global		bool			`json:"global,omitempty"`
	Description string			`json:"description,omitempty"`
	Options     []ConfigParam	`json:"options,omitempty"`
}

type ConfigParam struct {
	Name        string			`json:"name"`
	Required    bool			`json:"required"`
	Dynamic 	bool			`json:"dynamic,omitempty"`
	Description string			`json:"description"`
	Type        string			`json:"type"`
	Default     string			`json:"default"`
}

type EnvironmentVar struct {
	Name        string			`json:"name"`
	Required    bool			`json:"required"`
	Type        string			`json:"type,omitempty"`
	EnvFrom		string			`json:"envfrom,omitempty"`	// Looks for a globally named configmap.
	Description string			`json:"description"`
	Default		string			`json:"default,omitempty"`
	Example		string			`json:"example,omitempty"`
}
