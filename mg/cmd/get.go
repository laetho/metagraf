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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"metagraf/pkg/metagraf"
	"os"
	"strings"
	"github.com/tidwall/gjson"
)

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getCmdJSONPatch)
	getCmd.AddCommand(getCmdJSONPath)
	getCmdJSONPatch.AddCommand(getCmdJSONPatchLabels)


}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get subcommands",
	Long:  `get subcommands`,
}

var getCmdJSONPath = &cobra.Command{
	Use:   "jsonpath <metagraf> <query>",
	Short: "jsonpath ",
	Long:  `jsonpath `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(StrMissingMetaGraf)
			os.Exit(1)
		}
		data,_ := json.Marshal(metagraf.Parse(args[0]))
		value := gjson.Get(string(data), args[1])
		fmt.Println(value.String())
	},
}

var getCmdJSONPatch = &cobra.Command{
	Use:   "jsonpatch",
	Short: "patch subcommands",
	Long:  `patch subcommands`,
}

var getCmdJSONPatchLabels = &cobra.Command{
	Use:   "labels <metagraf>",
	Short: "fetch labels as a json patch structure for kubernetes resource model",
	Long:  `fetch labels as a json patch structure for kubernetes resource model`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(StrMissingMetaGraf)
			os.Exit(1)
		}
		mg := metagraf.Parse(args[0])

		// Anonymous struct and initialization
		data := struct {
			Metadata struct {
				Labels            map[string]string	`json:"labels,omitempty"`
			} `json:"metadata"`
		}{}
		data.Metadata.Labels = make(map[string]string)

		for k,v := range mg.Metadata.Annotations {
			if !validLabel(sanitizeLabel(v)) {
				continue
			}
			if strings.Contains(k,"norsk-tipping.no") {
				data.Metadata.Labels[k] = sanitizeLabel(v)
			}
		}
		for k,v := range mg.Metadata.Labels {
			data.Metadata.Labels[k] = v
		}
		labelString, _ := json.Marshal(data)
		fmt.Println(string(labelString))
	},
}

func sanitizeLabel(val string) string {
	ret := strings.Replace(val, " ", "_", -1)
	return ret
}

func validLabel(val string) bool {
	if len(val) > 64 {
		return false
	}
	return true
}