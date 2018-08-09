package cmd

import (
	"github.com/spf13/cobra"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/generators"
	)

func init() {
	createMiddlearthCmd.Flags().StringVar(&Metagraf, "metagraf", "","path to metaGraf file")
	createCmd.AddCommand(createMiddlearthCmd)
}

var createMiddlearthCmd = &cobra.Command{
	Use:   "middlearth",
	Short: "Generate middlearth application json",
	Long:  `Outputs a middlearth application json from a metaGraf definition`,
	Run: func(cmd *cobra.Command, args []string) {
		createMiddlearth(Metagraf)
	},
}

func createMiddlearth( mgf string) {
	mg := metagraf.Parse(mgf)
	generators.MiddlearthApp(&mg)
}