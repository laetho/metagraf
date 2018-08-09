/*
Copyright 2018 The MetaGraph Authors

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