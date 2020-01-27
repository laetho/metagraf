/*
Copyright 2018 The MetaGraph Authors

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
	log "k8s.io/klog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	"os"
)

func init() {
	RootCmd.AddCommand(devCmd)
	devCmd.AddCommand(devCmdUp)
	devCmd.AddCommand(devCmdDown)
	devCmdUp.Flags().StringVarP(&Namespace, "namespace", "n","","namespace to work on, if not supplied it will use current active namespace.")
	devCmdUp.Flags().StringVar(&Branch, "branch","", "Override branch to build from. Used when generating BuildConfig object.")
	devCmdUp.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	devCmdUp.Flags().StringVar(&CVfile, "cvfile","", "Property file with component configuration values. Can be generated with \"mg generate properties\" command.)")
	devCmdUp.Flags().StringVar(&OName, "name", "", "Overrides name of application.")
	devCmdUp.Flags().StringVarP(&Context,"context", "c","/","Application contextroot. (\"/<context>\"). Used when creating Route object.")
	devCmdUp.Flags().BoolVarP(&CreateGlobals, "globals", "g", false, "Override default behavior and force creation of global secrets. Will not overwrite existing ones.")
	devCmdDown.Flags().StringVarP(&Namespace, "namespace", "n","","namespace to work on, if not supplied it will use current active namespace.")
	devCmdDown.Flags().BoolVarP(&All, "all", "a", false,"Delete all component resources including images.")
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "dev subcommands",
	Long:  `dev subcommands`,
}

var devCmdUp = &cobra.Command{
	Use:   "up <metagraf.json>",
	Short: "creates the required component resources.",
	Long:  `sets up the required resources to test the component in the platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Info(StrActiveProject, viper.Get("namespace"))
			fmt.Println(StrMissingMetaGraf)
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

		devUp(args[0])
	},
}

var devCmdDown = &cobra.Command{
	Use:   "down <metagraf.json>",
	Short: "deletes component resources",
	Long:  `dev subcommands`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Info(StrActiveProject, viper.Get("namespace"))
			log.Error(StrMissingMetaGraf)
			return
		}

		if len(Namespace) == 0 {
			Namespace = viper.GetString("namespace")
			if len(Namespace) == 0 {
				log.Error(StrMissingNamespace)
				os.Exit(1)
			}
		}
		FlagPassingHack()

		devDown(args[0])
	},
}

func devUp(mgf string) {
	mg := metagraf.Parse(mgf)
	modules.Variables = mg.GetProperties()
	OverrideProperties(modules.Variables)
	log.V(2).Info("Current MGProperties: ", modules.Variables)

	modules.GenSecrets(&mg)
	modules.GenConfigMaps(&mg)
	modules.GenImageStream(&mg, Namespace)
	modules.GenBuildConfig(&mg)
	modules.GenDeploymentConfig(&mg, Namespace)
	modules.GenService(&mg)
	modules.GenRoute(&mg)

}

func devDown(mgf string) {
	mg := metagraf.Parse(mgf)
	basename := modules.Name(&mg)

	modules.DeleteRoute(basename)
	modules.DeleteService(basename)
	modules.DeleteDeploymentConfig(basename)
	modules.DeleteBuildConfig(basename)
	modules.DeleteConfigMaps(&mg)

	if All {
		modules.DeleteSecrets(&mg)
		modules.DeleteImageStream(basename)
	}
/*
	client := ocpclient.GetCoreClient()



	modules.DeleteBuildConfig()
*/
}