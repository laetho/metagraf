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
	"github.com/golang/glog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"metagraf/pkg/modules"
	"metagraf/pkg/metagraf"
)



func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createConfigMapCmd)
	createCmd.AddCommand(createDeploymentConfigCmd)
	createCmd.AddCommand(createBuildConfigCmd)
	createCmd.AddCommand(createImageStreamCmd)
	createCmd.AddCommand(createServiceCmd)
	createCmd.AddCommand(createDotCmd)
	createCmd.AddCommand(createRefCmd)
	createCmd.AddCommand(createSecretCmd)
	createCmd.Flags().StringArray("cvars", CVars, "String array of KEY=VALUE variables.", )
	createDeploymentConfigCmd.Flags().StringVar(&Namespace, "namespace", "", "namespace to work on, if not supplied it will use current working namespace")
	createBuildConfigCmd.Flags().StringVar(&Namespace, "namespace", "", "namespace to work on, if not supplied it will use current working namespace")
	createSecretCmd.Flags().StringVar(&Namespace, "namespace", "", "namespace to work on, if not supplied it will use current working namespace")
	createConfigMapCmd.Flags().StringVar(&Namespace, "namespace", "", "namespace to work on, if not supplied it will use current working namespace")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create operations",
	Long:  Banner + ` create `,
}

var createBuildConfigCmd = &cobra.Command{
	Use:   "buildconfig <metagraf>",
	Short: "create BuildConfig from metaGraf file",
	Long:  Banner + `create BuildConfig`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
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
		modules.GenBuildConfig(&mg)
	},
}

var createConfigMapCmd = &cobra.Command{
	Use:   "configmap <metagraf>",
	Short: "create ConfigMaps from metaGraf file",
	Long:  Banner + `create ConfigMap`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
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
		modules.GenConfigMaps(&mg)
	},
}

var createDeploymentConfigCmd = &cobra.Command{
	Use:   "deploymentconfig <metagraf>",
	Short: "create DeploymentConfig from metaGraf file",
	Long:  Banner + `create DeploymentConfig`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
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
		modules.GenDeploymentConfig(&mg, Namespace)
	},
}

var createImageStreamCmd = &cobra.Command{
	Use:   "imagestream <metagraf>",
	Short: "create ImageStream from metaGraf file",
	Long:  Banner + `create ImageStream`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
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
		modules.GenImageStream(&mg, Namespace)
	},
}

var createServiceCmd = &cobra.Command{
	Use:   "service <metagraf>",
	Short: "create Service from metaGraf file",
	Long:  Banner + `create Service`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
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
		modules.GenService(&mg)
	},
}

var createDotCmd = &cobra.Command{
	Use:   "dot <collection directory>",
	Short: "create Graphviz service graph from collectio of metaGraf's",
	Long:  Banner + `create dot`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing path to collection of metaGraf specifications")
			return
		}
		modules.GenDotFromPath(args[0])
	},
}

var createRefCmd = &cobra.Command{
	Use:   "ref <metaGraf>",
	Short: "create ref document from metaGraf specification",
	Long:  Banner + `create ref`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing path to metaGraf specifications")
			return
		}
		mg := metagraf.Parse(args[0])
		modules.GenRef(&mg)
	},
}

var createSecretCmd = &cobra.Command{
	Use:   "secret <metaGraf>",
	Short: "create Secrets from metaGraf specification",
	Long:  Banner + `create Secret`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			glog.Info("Active project is:", viper.Get("namespace"))
			glog.Error("Missing path to metaGraf specification")
			os.Exit(1)
		}

		if len(Namespace) == 0 {
			Namespace = viper.GetString("namespace")
			if len(Namespace) == 0 {
				glog.Error("Namespace must be supplied")
				os.Exit(1)
			}
		}

		modules.NameSpace = Namespace
		mg := metagraf.Parse(args[0])
		modules.GenSecrets(&mg)
	},
}