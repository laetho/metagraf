package cmd

import (
	"fmt"
	params2 "github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "k8s.io/klog"
	"os"
)

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createDeploymentCmd)
	createDeploymentCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "namespace to work on, if not supplied it will use current working namespace")
	createDeploymentCmd.Flags().StringVar(&OName, "name", "", "Overrides name of deployment.")
	createDeploymentCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	createDeploymentCmd.Flags().StringVar(&params2.PropertiesFile, "cvfile", "", "File with component configuration values. (key=value pairs)")
	createDeploymentCmd.Flags().BoolVar(&BaseEnvs, "baseenv", false, "Hydrate deploymentconfig with baseimage environment variables")
	createDeploymentCmd.Flags().BoolVar(&Defaults, "defaults", false, "Populate Environment variables with default values from metaGraf")
	createDeploymentCmd.Flags().StringVarP(&ImageNS, "imagens", "i", "", "Image Namespace, used to override default namespace")
	createDeploymentCmd.Flags().StringVarP(&Registry, "registry", "r", viper.GetString("registry"), "Specify container registry host")
	createDeploymentCmd.Flags().StringVarP(&Tag, "tag", "t", "latest", "specify custom tag")
	createDeploymentCmd.Flags().Int32Var(&params2.Replicas, "replicas", params2.DefaultReplicas, "Number of replicas.")
	createDeploymentCmd.Flags().BoolVar(&params2.DisableDeploymentImageAliasing, "disable-aliasing", false, "Only applies to .spec.image references. Aliasing will use mg conventions for image references. Setting this to true will disable that behavior.")
	createDeploymentCmd.Flags().BoolVar(&params2.WithAffinityRules, "with-affinity-rules", params2.WithPodAffinityRulesDefault, "Enable generation of pod affinity or anti-affinity rules.")
	createDeploymentCmd.Flags().StringVar(&params2.PodAntiAffinityTopologyKey, "anti-affinity-topology-key", "", "Define which node label to use as a topologyKey (describing a datacenter, zone or a rack as an example)")
	createDeploymentCmd.Flags().Int32Var(&params2.PodAntiAffinityWeight, "pod-anti-affinity-weight", params2.PodAntiAffinityWeightDefault, "Provide weight for WeightedPodAffinityTerm.")
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
		params2.NameSpace = Namespace

		if params2.WithAffinityRules && len(params2.PodAntiAffinityTopologyKey) == 0 {
			fmt.Println("ERROR: --affinity-topology-key cannot be empty when --with-affinity-rules is active!")
			os.Exit(1)
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
