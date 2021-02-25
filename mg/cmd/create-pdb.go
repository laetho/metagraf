package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "k8s.io/klog"
	"metagraf/internal/pkg/params/params"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/pdb"
	"os"
)

func init() {
	createCmd.AddCommand(createPodDisruptionBudget)
	createPodDisruptionBudget.Flags().IntVarP(&params.Replicas,"replicas", "r", 2, "Number of replicas.")
	createPodDisruptionBudget.Flags().StringVarP(&params.NameSpace,"namespace", "n", "", "Set namespace for generated resource.")
}


var createPodDisruptionBudget = &cobra.Command{
	Use:   "poddisruptionbudget <metagraf>",
	Short: "create PodDisruptionBudget from metaGraf file",
	Long:  MGBanner + `create PodDisruptionBudget`,
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
		params.Dryrun = Dryrun
		params.Output = Output
		mg := metagraf.Parse(args[0])

		pdb.GenPodDisruptionBudget(&mg, params.Replicas)
	},
}
