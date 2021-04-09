/*
Copyright 2020 The metaGraf Authors

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
	"github.com/laetho/metagraf/internal/pkg/params/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/oam"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	RootCmd.AddCommand(oamCmd)
	oamCmd.AddCommand(oamCreateCmd)
	oamCreateCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "verbose output")
	oamCreateCmd.PersistentFlags().BoolVar(&Output, "output", false, "also output objects")
	oamCreateCmd.PersistentFlags().StringVarP(&Format, "format", "o", "json", "specify json or yaml, json id default")
	oamCreateCmd.PersistentFlags().BoolVar(&Dryrun, "dryrun", false, "do not create objects, only output")
	oamCreateCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "", "namespace to work on")
	oamCreateCmd.AddCommand(oamCreateComponentCmd)
	oamCreateCmd.AddCommand(oamCreateApplicationConfigurationCmd)
}

var oamCmd = &cobra.Command{
	Use:   "oam",
	Short: "oam subcommands",
	Long:  `Subcommands for oam`,
}

var oamCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "oam create subcommands",
	Long:  `Create Subcommands for oam`,
}

var oamCreateComponentCmd = &cobra.Command{
	TraverseChildren: true,
	Use:              "component <metagraf>",
	Short:            "oam create component",
	Long:             `Creates an oam component from a metagraf specification`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Insufficient arguments")
			os.Exit(-1)
		}
		mg := metagraf.Parse(args[0])

		if len(params.NameSpace) == 0 {
			params.NameSpace = Namespace
		}
		oam.GenOAMComponent(&mg)
	},
}

var oamCreateApplicationConfigurationCmd = &cobra.Command{
	TraverseChildren: true,
	Use:              "configuration <metagraf>",
	Short:            "oam create configuration",
	Long:             `Creates an oam application configuration from a metagraf specification and properties`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Insufficient arguments")
			os.Exit(-1)
		}
		mg := metagraf.Parse(args[0])

		if len(params.NameSpace) == 0 {
			params.NameSpace = Namespace
		}

		oam.GenOAMApplicationConfiguration(&mg)
	},
}
