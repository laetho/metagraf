/*
Copyright 2018-2020 The metaGraf Authors

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

// Structure to hold specification section sourced parameters from input. Should
// solve key collisions and generally be a more workable solution.
type MGProperty struct {
	Source   string `json:"source"`
	Key      string `json:"key"`
	Value    string `json:"value,omitempty"`
	Required bool   `json:"required,omitempty"`
	Default  string `json:"default,omitempty"`
}

// map for holding MGProperty structs, should be keyed by
// MGProperty.Source + ":" + MGProperty.Key
type MGProperties map[string]MGProperty

type MetagrafType string

const (
	Application   MetagrafType = "application"
	Configuration MetagrafType = "config"
)

type ResourceType string

const (
	ClusterService ResourceType = "clusterservice"
	ExternalName   ResourceType = "externalname"
)

// JSON structure for a MetaGraf entity
type MetaGraf struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Name              string            `json:"name"`
		ResourceVersion   string            `json:"resourceversion"`
		Namespace         string            `json:"namespace"`
		CreationTimestamp string            `json:"creationtimestamp,omitempty"`
		Labels            map[string]string `json:"labels,omitempty"`
		Annotations       map[string]string `json:"annotations,omitempty"`
	} `json:"metadata"`
	Spec struct {
		// Describes the type of metaGraf specification. What types we have are not formalized yet.
		Type        string `json:"type"`
		Version     string `json:"version"`
		Description string `json:"description"`
		// Ports is a map keyed on protocol name (string) with port as it's value.
		Ports map[string]int32 `json:"ports,omitempty"`
		// Git repository URL for the source code of the described software component.
		Repository string `json:"repository,omitempty"`
		// Repository Secret Reference, git pull secret
		RepSecRef string `json:"repsecref,omitempty"`

		// Check out and build code from another branch than master. Defaults to master if
		// not provided.
		Branch string `json:"branch,omitempty"`

		// When a spec.image is specified, we want to deliver an existing image with
		// manifest generation provided with tools like mg.
		Image string `json:"image,omitempty"`

		// When spec.dockerfile is provided we will attempt to build the container image
		// with local tools, if present. Image and Dockerfile are mutually exclusive
		Dockerfile string `json:"dockerfile,omitempty"`

		// Image used to build the software referenced in Repository.
		BuildImage string `json:"buildimage,omitempty"`

		// Image to inject artifacts from above build.
		BaseRunImage string `json:"baserunimage,omitempty"`

		// StartupProbe, a v1.Probe{} definition from upstream Kubernetes.
		StartupProbe v1.Probe `json:"startupProbe,omitempty"`
		// LivenessProbe, a v1.Probe{} definition from upstream Kubernetes.
		LivenessProbe v1.Probe `json:"livenessProbe,omitempty"`
		// ReadinessProbe, a v1.Probe{} definition from upstream Kubernetes.
		ReadinessProbe v1.Probe `json:"readinessProbe,omitempty"`
		// Slice of Resource structs for holding information about attached resources.
		Resources []Resource `json:"resources,omitempty"`
		// Slice of strings to reference kubernetes resources manually maintained within the
		// repository in Spec.Resource. Downstream tooling may care about these.
		LocalManifests []string `json:"localManifests,omitempty"`

		// Structure for holding diffrent kind of environment variables.
		Environment struct {
			// Slice for holding environmentvariables for the build process.
			Build []EnvironmentVar `json:"build,omitempty"`
			// Environment variables that should be set on the Deployment resource.
			Local []EnvironmentVar `json:"local,omitempty"`
			// Environment variables or configuration keys that comes from some kind of
			// central configuration mechanism. Redis, etcd, your solution.
			External struct {
				// Slice for holding environment variable or configuration keys  that
				// are introduces to a central configuration solution.
				Introduces []EnvironmentVar `json:"introduces,omitempty"`
				// Slice of environment variable or configuration keys that this
				// compoent consumes from the central configuration solution.
				Consumes []EnvironmentVar `json:"consumes,omitempty"`
			} `json:"external,omitempty"`
		} `json:"environment,omitempty"`

		Config []Config `json:"config,omitempty"`

		// Slice of metagraf.Secret's for describing secrets needed in execution context.
		Secret []Secret `json:"secret,omitempty"`

		// Slice of metagraf.Secret's needed in build context.
		BuildSecret []Secret `json:"buildsecret,omitempty"`
	} `json:"spec"`
}

// Describes attached resources for a component. Ref, 12 factor app.
// This section is currently a mess because of "lift and shift" approach
// we  should never have done. Going forward all attached resources
// should become a Kubernets Service of some kind.
type Resource struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	External    bool   `json:"external"`
	Semop       string `json:"semop,omitempty"` // Semantic operator, how to evaluate version match/requirements.

	Semver string `json:"semver,omitempty"` // Semantic version to evaluate for attached resource
	EnvRef string `json:"envref,omitempty"` // Reference an Environment variable

	// Used when we need to generate configuration for connection to the described attached resource.
	Template    string `json:"template,omitempty"`    // Go txt template string for generating resource configuration.
	TemplateRef string `json:"templateref,omitempty"` // ConfigMap Reference, OUTDATED Use ConfigRef
	ConfigRef   string `json:"configref,omitempty"`   // ConfigMap Reference, Replaces TemplateRef which was not a good name.
	User        string `json:"user,omitempty"`
	Secret      string `json:"secret,omitempty"` // k8s Secret reference
}

type Config struct {
	Name string `json:"name"`
	Type string `json:"type"`
	// If this is set to true, this will just be a refernce to a existing ConfigMap
	Global bool `json:"global,omitempty"`
	// Controls the mount point for the ConfigMap
	MountPath   string        `json:"mountpath,omitempty"`
	Description string        `json:"description,omitempty"`
	Options     []ConfigParam `json:"options,omitempty"`
}

type ConfigParam struct {
	Name        string `json:"name"`
	Required    bool   `json:"required"`
	Dynamic     bool   `json:"dynamic,omitempty"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	SecretFrom  string `json:"secretfrom,omitempty"` // References a value from a k8s secret resource
}

type Secret struct {
	Name        string `json:"name"`
	Global      bool   `json:"global,omitempty"`
	Description string `json:"description,omitempty"`

	// Indicated where in the filesystem we would like to mount the Secret.
	MountPath string `json:"mountpath,omitempty"`

	// Using Kubernetes v1.KeyToPath struct for mapping individual secret key's to filenames.
	Items []v1.KeyToPath `json:"items,omitempty"`
}

type EnvironmentVar struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Type     string `json:"type,omitempty"`
	// Expose environment variables from ConfigMap resources.
	// All keys, value pairs in secret will be exported from
	// the ConfigMap into a running Pod.  The Environment.Name
	// will just be a placeholder value.
	EnvFrom string `json:"envfrom,omitempty"`
	// Expose  contents of a kubernets Secret as environment variables
	// exported into a running container. The values are only available
	// inside a running Pod or if you have access to view secrets in the
	// namespace. Exposes all key, values from the Secret. The
	// EnvironmentVar.Name will just be a placeholder.
	SecretFrom string `json:"secretfrom,omitempty"`
	// When exporting environment variables from a Secret or Configmap resource, you
	// have the option to specify the name of a key to export. If provided
	// the value from the referenced key will appear as EnvironmentVar.Name
	// inside the running Pod.
	Key         string `json:"key,omitempty"`
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
	Example     string `json:"example,omitempty"`
}

type Generator interface {
}
