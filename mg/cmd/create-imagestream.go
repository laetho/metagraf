/*
Copyright 2018-2020 The metaGraf Authors

Licensed under the Apache oc License, Version 2.0 (the "License");
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
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	log "k8s.io/klog"
	"os"
)

func init() {
	createCmd.AddCommand(createImageStreamCmd)
	createImageStreamCmd.Flags().StringVarP(&Namespace, "namespace", "n","", "namespace to work on, if not supplied it will use current working namespace")
	createImageStreamCmd.Flags().StringVar(&OName, "name", "", "Overrides name of application basename to generate imagestream for.")
	createImageStreamCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")

}

var createImageStreamCmd = &cobra.Command{
	Use:   "imagestream <metagraf>",
	Short: "create ImageStream from metaGraf file",
	Long:  MGBanner + `create ImageStream`,
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

		mg := metagraf.Parse(args[0])
		FlagPassingHack()

		if len(modules.NameSpace) == 0 {
			modules.NameSpace = Namespace
		}
		modules.GenImageStream(&mg, Namespace)
	},
}