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
	"github.com/tidwall/gjson"
	"metagraf/mg/params"
	"metagraf/pkg/metagraf"
	"os"
	"strings"
)

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getCmdJSONPatch)
	getCmd.AddCommand(getCmdGJSONPath)
	getCmdJSONPatch.AddCommand(getCmdJSONPatchLabels)
	getCmdJSONPatchLabels.Flags().StringVarP(
		&params.NameSpacingFilter,
		"filter",
		"f",
		"",
		"Filter to use while fetching annotations and turning them into labels. Example: \"kubernetes.io\" or \"productowner\"")
	getCmdJSONPatchLabels.Flags().BoolVarP(
		&params.NameSpacingStripHost,
		"striphost",
		"s",
		false,
		"Strips hostname from annotations and labels when creating a jsonpatch.")

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get subcommands",
	Long:  `get subcommands`,
}

// Check the following link for gjson query syntax:
// https://github.com/tidwall/gjson/blob/master/SYNTAX.md
var getCmdGJSONPath = &cobra.Command{
	Use:   "gjson <metagraf> <query>",
	Short: "get json or string with gjson query",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(StrMissingMetaGraf)
			os.Exit(1)
		}
		data, _ := json.Marshal(metagraf.Parse(args[0]))
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
				Labels map[string]string `json:"labels,omitempty"`
			} `json:"metadata"`
		}{}
		data.Metadata.Labels = make(map[string]string)

		for k, v := range mg.Metadata.Annotations {
			if !validLabel(sanitizeLabelValue(v)) {
				continue
			}
			if strings.Contains(k, params.NameSpacingFilter) {
				data.Metadata.Labels[sanitizeKey(k)] = sanitizeLabelValue(v)
			}
		}
		for k, v := range mg.Metadata.Labels {
			data.Metadata.Labels[sanitizeKey(k)] = v
		}
		labelString, _ := json.Marshal(data)
		fmt.Println(string(labelString))
	},
}

func sanitizeLabelValue(val string) string {
	ret := strings.Replace(val, " ", "_", -1)
	ret = strings.Replace(ret, ",", "-", -1 )
	return ret
}

func sanitizeKey(key string) string {
	if params.NameSpacingStripHost {
		parts := strings.Split(key,"/")
		if len(parts) > 1 {
			return strings.Join(parts[1:], "")
		}
	}
	return key
}

func validLabel(val string) bool {
	if len(val) > 64 {
		return false
	}
	return true
}
