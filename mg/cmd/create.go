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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"metagraf/pkg/metagraf"
	"metagraf/pkg/generators"
)

var (
	Namespace string
)

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createDeploymentConfigCmd)
	createDeploymentConfigCmd.Flags().StringVar(&Namespace, "namespace", "", "namespace to work on, if not supplied it will use current working namespace")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create operations",
	Long:  Banner + ` create `,
}

var createDeploymentConfigCmd = &cobra.Command{
	Use:   "deploymentconfig <metagraf>",
	Short: "create DeploymentConfig from metaGraf file",
	Long:  Banner + `create DeploymentConfig`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1  {
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
		generators.GenDeploymentConfig(&mg)
	},
}
