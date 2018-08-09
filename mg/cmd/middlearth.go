package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

func init() {
	createCmd.AddCommand(createMiddlearthCmd)
}

var createMiddlearthCmd = &cobra.Command{
	Use:   "middlearth",
	Short: "Generate middlearth application json",
	Long:  `tja`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mg create operations")
	},
}