/*
Copyright 2019 The metaGraf Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ProjectCmd = &cobra.Command{
	Use:   "project <name>",
	Short: "set active project / namespace",
	Long:  `sets the `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
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
