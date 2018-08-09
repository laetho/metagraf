package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create operations",
	Long:  `create `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mg create operations")
	},
}
