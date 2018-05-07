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

package main

import (
	"fmt"
	"metagraf/mgraf/metagraf"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	// Calling poc() function
	poc()
}

func poc() {
	var files []string
	basepath := "/home/a01595/go/src/metagraf/collections/example"

	// Wal the directory provided in basepath
	err := filepath.Walk(basepath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Loop through all files in the directory given in basepath,
	// ignore basepath itself, and ignore file names not containing "json"
	// Parse each json file
	fmt.Println("digraph {")
	for _, file := range files {
		if file == basepath {
			continue
		}
		if !strings.Contains(file, "json") {
			continue
		}
		mg := metagraf.Parse(file)
		// wtf... mg.Metadata.Name == "stuff"
		metagraf.Refgen(&mg)
		metagraf.ResourceDotGen(&mg)
	}
	fmt.Println("}")
}
