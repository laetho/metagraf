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
	"path/filepath"
	"os"
	"strings"
	"github.com/spf13/cobra"
	"fmt"
	"metagraf/internal/metagraf"
	"metagraf/internal/generators"
)

var pmgf string


func init() {
	pipelineCreateCmd.Flags().StringVar(&pmgf, "metagraf", "","path to metaGraf file")
	pipelineCmd.AddCommand(pipelineCreateCmd)
	rootCmd.AddCommand(pipelineCmd)
}

var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "pipeline operations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mg pipeline operations")

	},
}

var pipelineCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create pipeline",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("mg pipeline create")
		pipelineCreate( pmgf )
	},
}

func pipelineCreate(mgf string) {
	mg := metagraf.Parse(mgf)
	generators.GenConfigMaps(&mg)
	generators.GenBuildConfig(&mg)

}


func poc(cpath string) {
	var files []string

	// Walk the directory passed with cpath
	err := filepath.Walk(cpath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Loop through all files in the directory given in basepath,
	// ignore basepath itself, and ignore file names not containing "json"
	// Parse each json file

	//var mgs []metagraf.MetaGraf

	for _, file := range files {
		if file == cpath {
			continue
		}
		if !strings.Contains(file, "json") {
			continue
		}

		//mg := metagraf.Parse(file)

		//metagraf.Refgen(&mg)
		//generators.MiddlearthApp(&mg)
		//mgs = append(mgs, mg)
	}
	//sp := strings.Split(strings.TrimRight(cpath, "/"),"/")
	//metagraf.ResourceDotGen(&mgs, sp[len(sp)-1])
}
