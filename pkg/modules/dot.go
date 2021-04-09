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

package modules

import (
	"fmt"
	"html/template"
	"github.com/laetho/metagraf/pkg/metagraf"
	"os"
	"path/filepath"
	"strings"
)

func GenDotFromPath(cpath string) {

	var files []string
	var mgs []metagraf.MetaGraf

	// Walk the directory passed with cpath
	err := filepath.Walk(cpath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file == cpath {
			continue
		}
		if !strings.Contains(file, "json") {
			continue
		}

		mg := metagraf.Parse(file)
		mgs = append(mgs, mg)
	}

	genDot(&mgs)
}

//
func genDot(mgs *[]metagraf.MetaGraf) {
	tmpl := template.Must(template.ParseFiles(TmplBasePath + "/digraph.dot"))
	f, err := os.OpenFile("/tmp/output.dot", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = tmpl.Execute(f, mgs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("File written to: /tmp/output.dot")
}
