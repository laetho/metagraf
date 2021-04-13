/*
Copyright 2021 The metaGraf Authors

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
)
