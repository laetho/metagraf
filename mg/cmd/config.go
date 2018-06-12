package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config operations",
	Long:  `set, get, list, delete configuration parameters`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mg config operations")
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set configuration",
	Long:  `set`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mg config set")
	},
}