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
	"metagraf/pkg/metagraf"
	"metagraf/pkg/generators"
	)

func init() {
	createMiddlearthCmd.Flags().StringVar(&Metagraf, "metagraf", "","path to metaGraf file")
	createCmd.AddCommand(createMiddlearthCmd)
}

var createMiddlearthCmd = &cobra.Command{
	Use:   "middlearth",
	Short: "Generate middlearth application json",
	Long:  `Outputs a middlearth application json from a metaGraf definition`,
	Run: func(cmd *cobra.Command, args []string) {
		createMiddlearth(Metagraf)
	},
}

func createMiddlearth( mgf string) {
	mg := metagraf.Parse(mgf)
	generators.MiddlearthApp(&mg)
}