package cmd

import (
	"github.com/laetho/metagraf/internal/pkg/params/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "k8s.io/klog"
	"os"
)

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createDeploymentConfigCmd)
	createDeploymentConfigCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "namespace to work on, if not supplied it will use current working namespace")
	createDeploymentConfigCmd.Flags().StringVar(&OName, "name", "", "Overrides name of deployment.")
	createDeploymentConfigCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	createDeploymentConfigCmd.Flags().StringVar(&params.PropertiesFile, "cvfile", "", "File with component configuration values. (key=value pairs)")
	createDeploymentConfigCmd.Flags().BoolVar(&BaseEnvs, "baseenv", false, "Hydrate deploymentconfig with baseimage environment variables")
	createDeploymentConfigCmd.Flags().BoolVar(&Defaults, "defaults", false, "Populate Environment variables with default values from metaGraf")
	createDeploymentConfigCmd.Flags().StringVarP(&ImageNS, "imagens", "i", "", "Image Namespace, used to override default namespace")
	createDeploymentConfigCmd.Flags().StringVarP(&Registry, "registry", "r", viper.GetString("registry"), "Specify container registry host")
	createDeploymentConfigCmd.Flags().StringVarP(&Tag, "tag", "t", "latest", "specify custom tag")
	createDeploymentConfigCmd.Flags().Int32Var(&params.Replicas, "replicas", params.DefaultReplicas, "Number of replicas.")
	createDeploymentConfigCmd.Flags().BoolVar(&params.DisableDeploymentImageAliasing, "disable-aliasing", false, "Only applies to .spec.image references. Aliasing will use mg conventions for image references. Setting this to true will disable that behavior.")
}

var createDeploymentConfigCmd = &cobra.Command{
	Use:   "deploymentconfig <metagraf>",
	Short: "create DeploymentConfig from metaGraf file",
	Long:  MGBanner + `create DeploymentConfig`,
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

		modules.Variables = GetCmdProperties(mg.GetProperties())
		log.V(2).Info("Current MGProperties: ", modules.Variables)

		if len(modules.NameSpace) == 0 {
			modules.NameSpace = Namespace
		}

		modules.GenDeploymentConfig(&mg)
	},
}
