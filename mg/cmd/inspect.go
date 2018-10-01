package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"metagraf/pkg/modules"
	"metagraf/pkg/metagraf"
	"os"
)

var InspectCmd = &cobra.Command{
	Use:   "inspect <metaGraf>",
	Short: "inspect a metaGraf specification",
	Long:  `inspect a metaGraf specification and list objects that will be created or patched.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Active project is:", viper.Get("namespace"))
			fmt.Println("Missing path to metaGraf specification")
			return
		}

		if len(Namespace) == 0 {
			Namespace = viper.GetString("namespace")
			if len(Namespace) == 0 {
				fmt.Println("Namespace must be supplied")
				os.Exit(1)
			}
		}

		mg := metagraf.Parse(args[0])
		if modules.Variables == nil {
			vars := MergeVars(
				mg.GetVars(),
				OverrideVars(mg.GetVars(), CmdCVars(CVars).Parse()))
			modules.Variables = vars
		}
		name := modules.Name(&mg)
		for k,v := range modules.Variables {
			fmt.Println(name,"Variable:",k,v)
		}

		modules.InspectSecrets(&mg)
		modules.InspectConfigMaps(&mg)
	},
}

func init() {
	RootCmd.AddCommand(InspectCmd)
}