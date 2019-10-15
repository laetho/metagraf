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
)

func init() {
	initCmd.AddCommand(configCmdList)
	initCmd.AddCommand(configCmdSet)
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a metagraf.json file",
	Long:  `initialize a metagraf.json file`,
	Run: func(cmd *cobra.Command, args []string) {

		newmg := metagraf.MetaGraf{}

		// Populate default annotations and labels
		var annotations = map[string]string{
			"myannotation": "myvalue",
		}
		var labels = map[string]string{
			"mylabel": "labelvalue",
		}
		newmg.Metadata.Annotations = annotations
		newmg.Metadata.Labels = labels




		metagraf.Store("./metagraf.json", &newmg)

	},
}


