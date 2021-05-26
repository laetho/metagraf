package cmd

import (
	"github.com/laetho/metagraf/pkg/generators/kaniko"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	"github.com/spf13/cobra"
	log "k8s.io/klog"
)

func init() {
	RootCmd.AddCommand(buildCmd)
	buildCmd.AddCommand(buildKanikoCmd)
	buildKanikoCmd.Flags().StringVarP(&kaniko.KanikoPodOpts.Namespace,"namespace", "n", "", "namespace to operate on" )

	buildKanikoCmd.Flags().StringVar(&kaniko.KanikoPodOpts.DockerfileArg, "dockerfile", "Dockerfile", "Specify Kaniko --dockerfile argument")
	buildKanikoCmd.Flags().StringVar(&kaniko.KanikoPodOpts.ContextArg, "context","", "Specify Kaniko --context argument. Overrides Git ref from metaGraf specification.")
	buildKanikoCmd.Flags().StringVar(&kaniko.KanikoPodOpts.DestinationArg, "destination", "", "Specify Kaniko --destination argument. Registry reference.")

	//buildKanikoCmd.Flags().BoolVarP()

}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build operations",
	Long:  MGBanner + ` build `,
}

var buildKanikoCmd = &cobra.Command{
	Use:   "kaniko <metagraf>",
	Short: "create a kaniko build pod from metaGraf specification",
	Long:  MGBanner + `build kaniko <metagraf.json>`,
	Run: func(cmd *cobra.Command, args []string) {
		requireMetagraf(args)

		mg := metagraf.Parse(args[0])
		modules.Variables = GetCmdProperties(mg.GetProperties())
		log.V(2).Info("Current MGProperties: ", modules.Variables)

		generator := kaniko.NewKanikoPodGenerator(mg, metagraf.MGProperties{}, kaniko.KanikoPodOpts)
		obj := generator.GenerateKanikoPod("blah")
		kaniko.MarshalToYaml(obj)
		kaniko.MarshalToJson(obj)
	},
}