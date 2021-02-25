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

import (
	"fmt"
	"github.com/spf13/cobra"
	"metagraf/internal/pkg/params/params"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	"os"
)

func init() {
	RootCmd.AddCommand(argocdCmd)
	argocdCmd.AddCommand(argocdCreateCmd)
	argocdCreateCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "verbose output")
	argocdCreateCmd.PersistentFlags().BoolVar(&Output, "output", false, "also output objects")
	argocdCreateCmd.PersistentFlags().StringVarP(&Format, "format","o","json", "specify json or yaml, json id default")
	argocdCreateCmd.PersistentFlags().BoolVar(&Dryrun, "dryrun", false, "do not create objects, only output")
	argocdCreateCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n","", "namespace to work on")
	argocdCreateCmd.AddCommand(argocdCreateApplicationCmd)
	argocdCreateApplicationCmd.Flags().StringVar(&params.ArgoCDApplicationProject, "project","", "Project reference")
	argocdCreateApplicationCmd.Flags().StringVar(&params.ArgoCDApplicationNamespace, "argo-namespace", "", "Namespace for the ArgoCD Application resource.")
	argocdCreateApplicationCmd.Flags().StringVarP(&params.ArgoCDApplicationRepoURL, "repo", "r", "", "Repository URL")
	argocdCreateApplicationCmd.Flags().StringVarP(&params.ArgoCDApplicationRepoPath, "path", "p", "", "Path to manifests inside the repository")
	argocdCreateApplicationCmd.Flags().StringVar(&params.ArgoCDApplicationTargetRevision, "target-revision", params.ArgoCDApplicationTargetRevision, "Git ref for commit to synchronise.")
	argocdCreateApplicationCmd.Flags().BoolVar(&params.ArgoCDApplicationSourceDirectoryRecurse, "recurse", false, "Recursively traverse basepath looking for manifests")
	argocdCreateApplicationCmd.Flags().BoolVar(&params.ArgoCDSyncPolicyRetry, "retry", false, "Retry failed synchronizations?")
	argocdCreateApplicationCmd.Flags().Int64Var(&params.ArgoCDSyncPolicyRetryLimit, "retry-limit", 2, "Retry limit")
	argocdCreateApplicationCmd.Flags().BoolVarP(&params.ArgoCDAutomatedSyncPolicy, "auto", "a", false, "Generate an automated sync policy?")
	argocdCreateApplicationCmd.Flags().BoolVar(&params.ArgoCDAutomatedSyncPolicyPrune, "auto-prune", false, "Automatically delete removed items")
	argocdCreateApplicationCmd.Flags().BoolVar(&params.ArgoCDAutomatedSyncPolicySelfHeal, "auto-heal", false, "Try to self heal?")

	_ = argocdCreateCmd.MarkPersistentFlagRequired("namespace")
	_ = argocdCreateApplicationCmd.MarkFlagRequired("project")
	_ = argocdCreateApplicationCmd.MarkFlagRequired("repo")
	_ = argocdCreateApplicationCmd.MarkFlagRequired("path")
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
	TraverseChildren: true,
	Use:   "application <metagraf>",
	Short: "argocd create application",
	Long:  `Creates an ArgoCD Application from a metagraf specification`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Insufficient arguments")
			os.Exit(-1)
		}
		mg := metagraf.Parse(args[0])
		FlagPassingHack()

		if len(modules.NameSpace) == 0 {
			modules.NameSpace = Namespace
		}
		modules.GenArgoApplication(&mg)
	},
}
