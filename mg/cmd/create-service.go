package cmd

import (
	params2 "github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "k8s.io/klog"
	"os"
)

func init() {
	createCmd.AddCommand(createServiceCmd)
	createServiceCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "namespace to work on, if not supplied it will use current working namespace")
	createServiceCmd.Flags().StringVar(&OName, "name", "", "Overrides name of application used to prefix configmaps.")
	createServiceCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	createServiceCmd.Flags().BoolVar(&params2.ServiceMonitor, "service-monitor", false, "Set flag to also create a ServiceMonitor resource. Requires a cluster with the prometheus-operator.")
	createServiceCmd.Flags().StringVar(&params2.ServiceMonitorPath, "service-monitor-path", params2.ServiceMonitorPathDefault, "Path to scrape metrics from.")
	createServiceCmd.Flags().Int32Var(&params2.ServiceMonitorPort, "service-monitor-port", params2.ServiceMonitorPortDefault, "Set Service port to scrape by a ServiceMonitor.")
	createServiceCmd.Flags().StringVar(&params2.ServiceMonitorOperatorName, "service-monitor-operator-name", params2.ServiceMonitorOperatorName, "Name of prometheus-operator instance to create ServiceMonitor for.")
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
