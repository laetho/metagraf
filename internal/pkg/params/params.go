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

package params

var (

	// Everything bool is a flag for indicating if we want to delete or operate on all resources.
	// As an example, mg dev down --everything. Which deletes all created resources and artifacts from
	// mg dev up command.
	Everything bool = false

	// DefaultReplicas
	DefaultReplicas int32 = 1

	// Boolean flag indicating if we want to enable population of
	// Pod specification fields as environment variables.
	DownwardAPIEnvVars = false

	// Boolean flag indicating fi we want to mount a volume
	// containing annotations and labels
	DownwardAPIVolume = false

	// Replicas, indicate how many of a thing we want.
	Replicas int32

	// PropertiesFile, assigned with --cvfile.
	PropertiesFile string
	// BuildParams Slice of strings to hold overridden values when using dev-build command.
	BuildParams    []string

	// Potentially used by BuildConfig creation to override output imagestream
	OutputImagestream string
	// Override BuildSourceRef with somthing other than provided in specification.
	SourceRef string

	// Label and annotation namespacing filter
	NameSpacingFilter string
	// If set to true, will strip hostname namespacing from annotatons and labels
	// when generating a jsonpatch.
	NameSpacingStripHost bool

	// Namespace
	NameSpace string

	// Deployment image aliasing will expect mg convention tagging of upstream images.
	// mysql:1.2.3 becomes mg-mysqlv1 if your metagraf name is mg-mysqlv1 and version is in 1.x.x range.
	DisableDeploymentImageAliasing bool = false

	// Set to true for generating ServiceMonitor objects when creating services.
	ServiceMonitor bool = false
	// ServiceMonitor definition of which port to scrape.
	ServiceMonitorPort        int32
	ServiceMonitorPortDefault int32 = 8080
	// ServiceMonitor definition of scraping interval.
	ServiceMonitorInterval string = "10s"
	// ServiceMonitor definition of scraping scheme.
	ServiceMonitorScheme string = "http"
	// ServiceMonitor definition of scrape path.
	ServiceMonitorPath        string
	ServiceMonitorPathDefault string = "/prometheus"

	// Name of prometheus-operator instance that should discover the generated ServiceMonitor or PodMonitor resources.
	ServiceMonitorOperatorName string = "prometheus"

	// Relative path to template file for use when creating a software component reference document.
	RefTemplateFile       string = ""
	RefTemplateOutputFile string = "REF.md"

	// RegistryUser stores the explicitly defined username for a private registry. Usually passed to mg with --reguser.
	RegistryUser string
	// RegistryPassword stores the explicitly defined password for a private registry. Usually passed to mg with --regpass.
	RegistryPassword string

	// Toggle for generating corev1.Affinity{} in Pod Templates in Deployment or DeploymentConfig.
	WithAffinityRules bool
	// Default value for WithAffinityRules if it's not set.
	WithPodAffinityRulesDefault bool = false
	// Name of node label to use as a topologyKey when generating pod affinity rules.
	PodAntiAffinityTopologyKey string

	// Flag for defining weight in a WeightedPodAffinityTerm.
	PodAntiAffinityWeight int32
	// The default value for PodAntiAffinityWeight
	PodAntiAffinityWeightDefault int32 = 100

	// String to hold a container image name override
	ImageName string

	CreateStatefulSetPersistentVolumeClaim       bool = false
	StatefulSetPersistentVolumeClaimStorageClass string
	StatefulSetPersistentVolumeClaimSize         string
)
