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
)
