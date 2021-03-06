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
	RootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createDeploymentCmd)
	createDeploymentCmd.Flags().StringVarP(&Namespace, "namespace", "n","", "namespace to work on, if not supplied it will use current working namespace")
	createDeploymentCmd.Flags().StringVar(&OName, "name", "", "Overrides name of deployment.")
	createDeploymentCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	createDeploymentCmd.Flags().StringVar(&params.PropertiesFile, "cvfile","", "File with component configuration values. (key=value pairs)")
	createDeploymentCmd.Flags().BoolVar(&BaseEnvs, "baseenv", false, "Hydrate deploymentconfig with baseimage environment variables")
	createDeploymentCmd.Flags().BoolVar(&Defaults, "defaults", false, "Populate Environment variables with default values from metaGraf")
	createDeploymentCmd.Flags().StringVarP(&ImageNS,"imagens", "i", "", "Image Namespace, used to override default namespace")
	createDeploymentCmd.Flags().StringVarP(&Registry,"registry", "r",viper.GetString("registry"), "Specify container registry host")
	createDeploymentCmd.Flags().StringVarP(&Tag,"tag", "t", "latest", "specify custom tag")
	createDeploymentCmd.Flags().Int32Var(&params.Replicas,"replicas", params.DefaultReplicas, "Number of replicas.")
	createDeploymentCmd.Flags().BoolVar(&params.DisableDeploymentImageAliasing, "disable-aliasing", false, "Only applies to .spec.image references. Aliasing will use mg conventions for image references. Setting this to true will disable that behavior.")
}

var createDeploymentCmd = &cobra.Command{
	Use:   "deployment <metagraf>",
	Short: "create Deployment from metaGraf file",
	Long:  MGBanner + `create Deployment`,
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

		if len(modules.NameSpace) == 0 {
			modules.NameSpace = Namespace
		}
		modules.GenDeployment(&mg, Namespace)
	},
}