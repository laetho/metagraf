/*
Copyright 2018-2019 The MetaGraph Authors

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
	"github.com/blang/semver"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"metagraf/pkg/metagraf"
	"os"
)

func init() {
	RootCmd.AddCommand(injectCmd)
	injectCmd.AddCommand(injectAnnotationsCmd)
	injectCmd.AddCommand(injectVersionCmd)
	injectCmd.AddCommand(injectSemVerCmd)
	//injectAnnotationsCmd.Flags().StringSliceVar(&CVars, "values", []string{}, "Slice of key=value pairs, seperated by ,")
}

var injectCmd = &cobra.Command{
	Use:   "inject",
	Short: "Inject operations",
	Long:  Banner + ` inject `,
}

var injectAnnotationsCmd = &cobra.Command{
	Use: "annotations <metagraf> <arg>",
	Short: "Injects annotations",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args[0]) < 1 {
			glog.Error(StrMissingMetaGraf)
			os.Exit(1)
		}

		if len(args[1]) < 1 {
			glog.Error("Missing key")
			os.Exit(1)
		}

		if len(args[2]) < 1 {
			glog.Error("Missing value for ", args[1] )
			os.Exit(1)
		}

		mg := metagraf.Parse(args[0])
		mg.Metadata.Annotations[args[1]] = args[2]

		metagraf.Store(args[0], &mg)
	},
}

var injectVersionCmd = &cobra.Command{
	Use: "version <metagraf> <version>",
	Short: "Injects a custom version for the component.",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args[0]) < 1 {
			glog.Error(StrMissingMetaGraf)
			os.Exit(1)
		}

		if len(args[1]) < 1 {
			glog.Error("You have to specify a version.")
			os.Exit(1)
		}

		mg := metagraf.Parse(args[0])
		mg.Spec.Version = args[1]

		metagraf.Store(args[0], &mg)

	},
}

var injectSemVerCmd = &cobra.Command{
	Use: "semver <metagraf> <version>",
	Short: "Injects a SemVer 2.0 version for the component.",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args[0]) < 1 {
			glog.Error(StrMissingMetaGraf)
			os.Exit(1)
		}

		if len(args[1]) < 1 {
			glog.Error("You have to specify a version.")
			os.Exit(1)
		}

		sv, err := semver.Parse(args[1])
		if err != nil {
			glog.Error(err)
			os.Exit(1)
		}

		err = sv.Validate()
		if err != nil {
			glog.Error("err")
			os.Exit(1)
		}

		mg := metagraf.Parse(args[0])
		mg.Spec.Version = args[1]

		metagraf.Store(args[0], &mg)

	},
}