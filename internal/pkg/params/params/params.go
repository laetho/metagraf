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

	// Dryrun bool indicated if mg should do operations against a Kubernetes API.
	Dryrun bool

	// Output bool indicates if mg should output the generated objects.
	Output bool

	// Format value can be either json or yaml. Controls output format.
	Format string = "json"

	// DefaultReplicas
	DefaultReplicas int32 = 1

	// Replicas, indicate how many of a thing we want.
	Replicas int32

	// PropertiesFile, assigned with --cvfile.
	PropertiesFile string

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

	ArgoCDApplicationProject string
	// In which namespace do we want the ArgoCD Application CR to be created
	ArgoCDApplicationNamespace string
	ArgoCDApplicationRepoURL   string
	ArgoCDApplicationRepoPath  string
	// Git Reference (tag/commit)
	ArgoCDApplicationTargetRevision         string = "HEAD"
	ArgoCDApplicationSourceDirectoryRecurse bool
	ArgoCDSyncPolicyRetry                   bool
	ArgoCDSyncPolicyRetryLimit              int64
	ArgoCDAutomatedSyncPolicy               bool
	ArgoCDAutomatedSyncPolicyPrune          bool
	ArgoCDAutomatedSyncPolicySelfHeal       bool

	// Deployment image aliasing will expect mg convention tagging of upstream images.
	// mysql:1.2.3 becomes mg-mysqlv1 if your metagraf name is mg-mysqlv1 and version is in 1.x.x range.
	DisableDeploymentImageAliasing bool = false

	// Set to true for generating ServiceMonitor objects when creating services.
	ServiceMonitor bool = false
	// ServiceMonitor definition of which port to scrape.
	ServiceMonitorPort int32
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
	WithAffinityRules			bool
	// Default value for WithAffinityRules if it's not set.
	WithPodAffinityRulesDefault bool = false
	// Name of node label to use as a topologyKey when generating pod affinity rules.
	PodAffinityTopologyKey 		string

	// Flag for defining weight in a WeightedPodAffinityTerm.
	PodAffinityWeight			int32
	// The default value for PodAffinityWeight
	PodAffinityWeightDefault	int32 = 100

)
