package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	"metagraf/mg/params"
	"os"
	log "k8s.io/klog"
)

func init() {
	createBuildConfigCmd.Flags().StringVar(&OName, "name", "", "Overrides name of BuildConfig.")
	createBuildConfigCmd.Flags().StringVarP(&Tag,"tag", "t", "latest", "specifies custom output tag")
	createBuildConfigCmd.Flags().StringVarP(&params.OutputImagestream,"ostream", "i", "", "specify if you want to output to another imagestream than the component name")
	createBuildConfigCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "namespace to work on, if not supplied it will use current working namespace")
	createBuildConfigCmd.Flags().StringVar(&Branch, "branch", "", "Override branch to build from.")
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

		fmt.Println(mg)

		modules.Variables = mg.GetProperties()
		OverrideProperties(modules.Variables)
		log.V(2).Info("Current MGProperties: ", modules.Variables)

		if len(modules.NameSpace) == 0 {
			modules.NameSpace = Namespace
		}

		modules.GenBuildConfig(&mg)
	},
}