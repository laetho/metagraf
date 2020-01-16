/*
Copyright 2019 The MetaGraph Authors

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
	"github.com/spf13/cobra/doc"
	log "k8s.io/klog"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	"os"
)

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(generatePropertiesCmd)
	generatePropertiesCmd.Flags().BoolVar(&Defaults, "defaults", false, "Populate Environment variables with default values from metaGraf")
	generatePropertiesCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	generateCmd.AddCommand(generateMarkdownCmd)
	generateCmd.AddCommand(generateManPagesCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate operations",
	Long:  Banner + ` create `,
}

var generateMarkdownCmd = &cobra.Command{
	Use:   "markdown",
	Short: "create markdown documentation for the mg command.",
	Long:  Banner + `generate markdown`,
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree( RootCmd, "./")
		if err != nil {
			log.Fatalf("%v", err)
			os.Exit(1)
		}
	},
}

var generateManPagesCmd = &cobra.Command{
	Use:   "manpages",
	Short: "create manpages for the mg command",
	Long:  Banner + `generate markdown`,
	Run: func(cmd *cobra.Command, args []string) {

		header := &doc.GenManHeader{
			Title: "mg",
			Section: "3",
		}
		err := doc.GenManTree( RootCmd, header, "./")
		if err != nil {
			log.Fatalf("%v", err)
			os.Exit(1)
		}
	},
}
var generatePropertiesCmd = &cobra.Command{
	Use:   "properties <metagraf>",
	Short: "create configuration properties from metaGraf file",
	Long:  Banner + `generate properties`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error(StrMissingMetaGraf)
			os.Exit(1)
		}

		FlagPassingHack()

		mg := metagraf.Parse(args[0])
		if modules.Variables == nil {
			vars := MergeSourceVars(
				mg.GetVarsFromSource(Defaults),
				OverrideVars(mg.GetVars()))
			modules.Variables = vars
		}
		for k,v := range modules.Variables {
			fmt.Println(k+"="+v)
		}
	},
}