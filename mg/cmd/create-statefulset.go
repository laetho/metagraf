package cmd

import (
	"fmt"
	"os"

	"github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "k8s.io/klog"
)

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createStatefulsetCmd)
	createStatefulsetCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "namespace to work on, if not supplied it will use current working namespace")
	createStatefulsetCmd.Flags().StringVar(&OName, "name", "", "Overrides name of StatefulSet.")
	createStatefulsetCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	createStatefulsetCmd.Flags().StringVar(&params.PropertiesFile, "cvfile", "", "File with component configuration values. (key=value pairs)")
	createStatefulsetCmd.Flags().BoolVar(&BaseEnvs, "baseenv", false, "Hydrate deploymentconfig with baseimage environment variables")
	createStatefulsetCmd.Flags().BoolVar(&Defaults, "defaults", false, "Populate Environment variables with default values from metaGraf")
	createStatefulsetCmd.Flags().StringVarP(&ImageNS, "imagens", "i", "", "Image Namespace, used to override default namespace")
	createStatefulsetCmd.Flags().StringVar(&params.ImageName, "imagename", "", "Set image artifact name. Overrides imagename from metaGraf spec parsing behaviour.")
	createStatefulsetCmd.Flags().StringVarP(&Registry, "registry", "r", viper.GetString("registry"), "Specify container registry host")
	createStatefulsetCmd.Flags().StringVarP(&Tag, "tag", "t", "latest", "specify custom tag")
	createStatefulsetCmd.Flags().Int32Var(&params.Replicas, "replicas", params.DefaultReplicas, "Number of replicas.")
	createStatefulsetCmd.Flags().BoolVar(&params.DisableDeploymentImageAliasing, "disable-aliasing", false, "Only applies to .spec.image references. Aliasing will use mg conventions for image references. Setting this to true will disable that behavior.")
	createStatefulsetCmd.Flags().BoolVar(&params.WithAffinityRules, "with-affinity-rules", params.WithPodAffinityRulesDefault, "Enable generation of pod affinity or anti-affinity rules.")
	createStatefulsetCmd.Flags().StringVar(&params.PodAntiAffinityTopologyKey, "anti-affinity-topology-key", "", "Define which node label to use as a topologyKey (describing a datacenter, zone or a rack as an example)")
	createStatefulsetCmd.Flags().Int32Var(&params.PodAntiAffinityWeight, "pod-anti-affinity-weight", params.PodAntiAffinityWeightDefault, "Weight for WeightedPodAffinityTerm.")
	createStatefulsetCmd.Flags().BoolVar(&params.DownwardAPIEnvVars, "downward-api-envvars", false, "Enables generation of environment variables from Downward API. An opinionated selection.")
	createStatefulsetCmd.Flags().BoolVar(&params.CreateStatefulSetPersistentVolumeClaim, "create-persistent-volumeclaim", false, "Enables the creation of a persistent volume claims template.")
	createStatefulsetCmd.Flags().StringVar(&params.StatefulSetPersistentVolumeClaimStorageClass, "statefulset-pv-storageclass", "", "StatefulSet persistent volume claim storage class")
	createStatefulsetCmd.Flags().StringVar(&params.StatefulSetPersistentVolumeClaimSize, "statefulset-pv-size", "", "StatefulSet persistent volume claim size")
}

var createStatefulsetCmd = &cobra.Command{
	Use:   "statefulset <metagraf>",
	Short: "create StatefulSet from metaGraf file",
	Long:  MGBanner + `create StatefulSet`,
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
		params.NameSpace = Namespace

		if params.WithAffinityRules && len(params.PodAntiAffinityTopologyKey) == 0 {
			fmt.Println("ERROR: --affinity-topology-key cannot be empty when --with-affinity-rules is active!")
			os.Exit(1)
		}

		mg := metagraf.Parse(args[0])
		FlagPassingHack()

		modules.Variables = GetCmdProperties(mg.GetProperties())

		if len(modules.NameSpace) == 0 {
			modules.NameSpace = Namespace
		}
		modules.GenStatefulSet(&mg, Namespace)
	},
}
