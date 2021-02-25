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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "k8s.io/klog"
	"metagraf/internal/pkg/params/params"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	"os"
)

func init() {
	createCmd.AddCommand(createBuildConfigCmd)
	createBuildConfigCmd.Flags().StringVar(&OName, "name", "", "Overrides name of BuildConfig.")
	createBuildConfigCmd.Flags().StringVarP(&Tag,"tag", "t", "latest", "specifies custom output tag")
	createBuildConfigCmd.Flags().StringVarP(&params.OutputImagestream,"istream", "i", "", "specify if you want to output to another imagestream than the component name")
	createBuildConfigCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "namespace to work on, if not supplied it will use current working namespace")
	createBuildConfigCmd.Flags().StringVar(&params.SourceRef, "ref", "", "specify source ref or branch name.")
	createBuildConfigCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
}

var createBuildConfigCmd = &cobra.Command{
	Use:   "buildconfig <metagraf>",
	Short: "create BuildConfig from metaGraf file",
	Long:  MGBanner + `create BuildConfig`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Info(StrActiveProject, viper.Get("namespace"))
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
		FlagPassingHack()
		mg := metagraf.Parse(args[0])

		modules.Variables = GetCmdProperties(mg.GetProperties())
		log.V(2).Info("Current MGProperties: ", modules.Variables)

		if len(modules.NameSpace) == 0 {
			modules.NameSpace = Namespace
		}

		modules.GenBuildConfig(&mg)
	},
}