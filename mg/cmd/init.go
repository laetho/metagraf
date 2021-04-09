/*
Copyright 2018 The metaGraf Authors

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
	"bufio"
	"fmt"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a metagraf.json file",
	Long:  `initialize a metagraf.json file`,
	Run: func(cmd *cobra.Command, args []string) {

		newmg := metagraf.MetaGraf{}
		newmg.Kind = "metagraf"
		initInput(&newmg)

		// Populate default annotations and labels
		var annotations = map[string]string{
			"app": newmg.Metadata.Name,
		}
		var labels = map[string]string{
			"metagraf": "true",
		}
		newmg.Metadata.Annotations = annotations
		newmg.Metadata.Labels = labels
		newmg.Spec.Environment.Local = append(newmg.Spec.Environment.Local, metagraf.EnvironmentVar{
			Name:        "LOCAL_EXAMPLE_ENV",
			Required:    false,
			Type:        "string",
			Description: "An example local environment variable.",
			Default:     "default",
		})
		newmg.Spec.Environment.Build = append(newmg.Spec.Environment.Build, metagraf.EnvironmentVar{
			Name:        "BUILD_EXAMPLE_ENV",
			Required:    false,
			Type:        "string",
			Description: "An example build environment variable.",
			Default:     "default",
		})
		newmg.Spec.Config = append(newmg.Spec.Config, metagraf.Config{
			Name:        "example.properties",
			Type:        "parameters",
			Global:      false,
			Description: "An example file",
			Options: append([]metagraf.ConfigParam{}, metagraf.ConfigParam{
				Name:        "key",
				Required:    true,
				Dynamic:     false,
				Description: "A key",
				Type:        "string",
				Default:     "value",
			}),
		})
		newmg.Spec.Secret = append(newmg.Spec.Secret, metagraf.Secret{
			Name:        "ca.cert",
			Global:      true,
			Description: "Certificate Authority",
		})

		metagraf.Store("./metagraf.json", &newmg)

	},
}

func initInput(mg *metagraf.MetaGraf) {

	reader := bufio.NewReader(os.Stdin)
	text := ""

	fmt.Println("Initialize a metaGraf specification:")

	fmt.Print(" Name of component -> ")
	text, _ = reader.ReadString('\n')
	text = strings.Trim(text, "\n\r")
	if len(text) > 0 {
		mg.Metadata.Name = text
	}

	fmt.Print(" Description of component -> ")
	text, _ = reader.ReadString('\n')
	text = strings.Trim(text, "\n\r")
	if len(text) > 0 {
		mg.Spec.Description = text
	}

	fmt.Print("GIT Repository url -> ")
	text, _ = reader.ReadString('\n')
	text = strings.Trim(text, "\n\r")
	if len(text) > 0 {
		mg.Spec.Repository = text
	}

	fmt.Print("GIT Repository branch -> ")
	text, _ = reader.ReadString('\n')
	text = strings.Trim(text, "\n\r")
	if len(text) > 0 {
		mg.Spec.Branch = text
	}

}
