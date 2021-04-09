/*
Copyright 2018-2020 The metaGraf Authors

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
	"github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	"github.com/spf13/cobra"
	log "k8s.io/klog"
	"os"
)

func init() {
	createCmd.AddCommand(createApplicationCmd)
	createApplicationCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "namespace to work on, if not supplied it will use current working namespace")
	createApplicationCmd.Flags().StringSliceVar(&CVars, "cvars", []string{}, "Slice of key=value pairs, seperated by ,")
	createApplicationCmd.Flags().StringVar(&params.PropertiesFile, "cvfile", "", "File with component configuration values. (source|key=value) pairs")
}

var createApplicationCmd = &cobra.Command{
	Use:   "application <metaGraf>",
	Short: "create a Kubernetes Application SIG resource from metaGraf specification.",
	Long:  `reate a Kubernetes Application SIG resource from metaGraf specification.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			log.Error(StrMissingMetaGraf)
			fmt.Println(StrMissingMetaGraf)
			os.Exit(1)
		}
		if len(Namespace) == 0 {
			log.Error(StrMissingNamespace)
			fmt.Println(StrMissingNamespace)
			os.Exit(1)
		}

		FlagPassingHack()
		mg := metagraf.Parse(args[0])
		modules.GenApplication(&mg)
	},
}
