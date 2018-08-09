package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createPipelineCmd)
	createCmd.AddCommand(createMiddlearthCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create operations",
	Long:  `create `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mg create operations")
	},
}

var createPipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Generate kubernetes primitives",
	Long:  `create`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mg config set")
	},
}

var createMiddlearthCmd = &cobra.Command{
	Use:   "middlearth",
	Short: "Generate middlearth application json",
	Long:  `Parse a metaGraf and produce a Middlearth Application JSON`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mg config set")
	},
}