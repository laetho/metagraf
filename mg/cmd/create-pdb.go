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

package cmd

import (
	"github.com/laetho/metagraf/internal/pkg/params/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/pdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "k8s.io/klog"
	"os"
)

func init() {
	createCmd.AddCommand(createPodDisruptionBudget)
	createPodDisruptionBudget.Flags().Int32Var(&params.Replicas, "replicas", params.DefaultReplicas, "Number of replicas.")
	createPodDisruptionBudget.Flags().StringVarP(&params.NameSpace, "namespace", "n", "", "Set namespace for generated resource.")
}

var createPodDisruptionBudget = &cobra.Command{
	Use:   "poddisruptionbudget <metagraf>",
	Short: "create PodDisruptionBudget from metaGraf file",
	Long:  MGBanner + `create PodDisruptionBudget`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.V(2).Info(StrActiveProject, viper.Get("namespace"))
			log.Error(StrMissingMetaGraf)
			os.Exit(1)
		}

		if len(Namespace) == 0 {
			Namespace = viper.GetString("namespace")
			if len(Namespace) == 0 {
				log.Error(StrMissingNamespace)
				os.Exit(1)
			}
		}
		// Migration to params not complete.
		params.Dryrun = Dryrun
		params.Output = Output
		mg := metagraf.Parse(args[0])

		pdb.GenPodDisruptionBudget(&mg, params.Replicas)
	},
}
