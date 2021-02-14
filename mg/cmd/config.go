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
	"metagraf/internal/pkg/helpers/helpers"
)

func init() {
	configCmd.AddCommand(configCmdList)
	configCmd.AddCommand(configCmdSet)
	RootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config operations",
	Long:  `set, get, list, delete configuration parameters`,
}

var configCmdSet = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "<key> <value>",
	Long:  `set`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Insufficient arguments")
			return
		}
		if helpers.StringInSlice(args[0], configkeys) {
			viper.Set(args[0], args[1])
			err := viper.WriteConfig()
			if err != nil {
				fmt.Println("ERROR:", err)
				return
			}
		}
	},
}

var configCmdList = &cobra.Command{
	Use:   "list",
	Short: "list configuration",
	Long:  `list current configuration settings`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, ck := range configkeys {
			fmt.Println(ck, ":", viper.GetString(ck))
		}
	},
}
