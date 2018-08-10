package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/spf13/viper"
)

var ProjectCmd = &cobra.Command{
	Use:   "project <name>",
	Short: "set active project / namespace",
	Long:  `sets the `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1  {
			fmt.Println("Active project is:", viper.Get("namespace"))
			return
		}
		name := args[0]
		viper.Set("namespace", name)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		fmt.Printf("Active namespace is now %v\n", name)
	},
}

func init() {
	RootCmd.AddCommand(ProjectCmd)
}
