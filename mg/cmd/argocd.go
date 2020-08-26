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

package cmd

import "github.com/spf13/cobra"

func init() {
	RootCmd.AddCommand(argocdCmd)
	argocdCmd.AddCommand(argocdCreateCmd)
	argocdCreateCmd.AddCommand(argocdCreateApplicationCmd)
}

var argocdCmd = &cobra.Command{
	Use:   "argocd",
	Short: "argocd subcommands",
	Long:  `Subcommands for ArgoCD`,
}

var argocdCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "argocd create subcommands",
	Long:  `Create Subcommands for ArgoCD`,
}

var argocdCreateApplicationCmd = &cobra.Command{
	Use:   "application",
	Short: "argocd create application",
	Long:  `Creates an ArgoCD Application from a metagraf specification`,
}
