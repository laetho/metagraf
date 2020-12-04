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
	createServiceCmd.Flags().BoolVar(&params.ServiceMonitor, "service-monitor",false, "Set flag to also create a ServiceMonitor resource. Requires a cluster with the prometheus-operator.")
	createServiceCmd.Flags().StringVar(&params.ServiceMonitorPath, "service-monitor-path", params.ServiceMonitorPathDefault, "Path to scrape metrics from.")
	createServiceCmd.Flags().Int32Var(&params.ServiceMonitorPort, "service-monitor-port", params.ServiceMonitorPortDefault, "Set Service port to scrape by a ServiceMonitor.")
	createServiceCmd.Flags().StringVar(&params.ServiceMonitorOperatorName, "service-monitor-operator-name", params.ServiceMonitorOperatorName,"Name of prometheus-operator instance to create ServiceMonitor for.")
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
	},
}

