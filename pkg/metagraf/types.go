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

import (
	v1 "k8s.io/api/core/v1"
)

// Map to hold all variables from a specification
type MGVars			map[string]string

// Structure to hold specification secrtion sourced parameters. Should
// solve key collisions and generally be a more workable solution.
type MGProperty struct {
	Source		string	`json:source`
	Key			string	`json:key`
	Value		string	`json:value, omitempty`
	//Required	bool	`json:required, omitempty`	// Add this when refactoring is done.
}


type MGProperties []MGProperty

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
		Type			string		`json:"type"`
		Version			string		`json:"version"`
		Description		string		`json:"description"`
		Repository		string  	`json:"repository,omitempty"`
		RepSecRef		string		`json:"repsecref,omitempty"`
		Branch			string		`json:"branch,omitempty"`
		BuildImage		string		`json:"buildimage,omitempty"`	// Builder Image for s2i builds.
		BaseRunImage 	string		`json:"baserunimage,omitempty"`	// Runtime Container Image for binary build.
		Image		 	string		`json:"image,omitempty"`		// Container Image URL, for wrapping upstream images.
		LivenessProbe	v1.Probe	`json:"livenessProbe,omitempty"`	// Using the k8s Probe type
		ReadinessProbe	v1.Probe	`json:"readinessProbe,omitempty"`	// Using the k8s Probe type

		Resources	 []Resource	`json:"resources,omitempty"`
		Environment struct {
			Build []EnvironmentVar	`json:"build,omitempty"`
			Local []EnvironmentVar	`json:"local,omitempty"`
			External struct {
				Introduces []EnvironmentVar `json:"introduces,omitempty"`
				Consumes   []EnvironmentVar `json:"consumes,omitempty"`
			} `json:"external,omitempty"`
		} `json:"environment,omitempty"`
		Config []Config `json:"config,omitempty"`
		Secret []Secret `json:"secret,omitempty"`
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
	EnvRef		string			`json:"envref,omitempty"`		// Reference an Environment variable

	// Used when we need to generate configuration for connection to the described attached resource.
	Template	string  		`json:"template,omitempty"`		// Go txt template string for generating resource configuration.
	TemplateRef string			`json:"templateref,omitempty"`	// ConfigMap Reference, OUTDATED Use ConfigRef
	ConfigRef	string			`json:"configref,omitempty"`	// ConfigMap Reference, Replaces TemplateRef which was not a good name.
	User 		string			`json:"user,omitempty"`
	Secret	    string			`json:"secret,omitempty"`		// k8s Secret reference
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
	SecretFrom  string			`json:"secretfrom,omitempty"`	// References a value from a k8s secret resource
}

type Secret struct {
	Name    	string			`json:"name"`
	Global		bool			`json:"global,omitempty"`
	Description string			`json:"description,omitempty"`
	Value		string			`json:"value,omitempty"`		// Never use this!
}

type EnvironmentVar struct {
	Name        string			`json:"name"`
	Required    bool			`json:"required"`
	Type        string			`json:"type,omitempty"`
	EnvFrom		string			`json:"envfrom,omitempty"`		// Looks for a globally named configmap.
	SecretFrom	string			`json:"secretfrom,omitempty"`	// References a k8s Secret resource.
	Description string			`json:"description"`
	Default		string			`json:"default,omitempty"`
	Example		string			`json:"example,omitempty"`
}
