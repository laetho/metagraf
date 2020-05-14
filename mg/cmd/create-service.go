package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"metagraf/mg/params"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	log "k8s.io/klog"
	"os"
)

func init() {
	createCmd.AddCommand(createServiceCmd)
	createServiceCmd.Flags().StringVarP(&Namespace, "namespace", "n","", "namespace to work on, if not supplied it will use current working namespace")
	createServiceCmd.Flags().StringVar(&OName, "name", "", "Overrides name of application used to prefix configmaps.")
	createServiceCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	createServiceCmd.Flags().BoolVarP(&params.ServiceMonitor, "monitor","m",false, "Set flag to also create a ServiceMonitor resource. Requires a cluster with the prometheus-operator.")
}

var createServiceCmd = &cobra.Command{
	Use:   "service <metagraf>",
	Short: "create Service from metaGraf file",
	Long:  MGBanner + `create Service`,
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
		modules.GenService(&mg)
		if params.ServiceMonitor {
			modules.GenServiceMonitor(&mg)
		}

	},
}

