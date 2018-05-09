/*
Copyright 2018 The MetaGraf Authors

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

package metagraf

import (
	"fmt"
	"text/template"
	"os"
)

var TmplBasePath = "/home/a01595/go/src/metagraf/templates"

// Takes a pointer to a metagraf.MetaGraf struct
func Refgen(mg *MetaGraf) {
	tmpl := template.Must(template.ParseFiles(TmplBasePath + "/refdoc.html"))
	f, err := os.OpenFile("/home/a01595/go/src/metagraf/docs/refdoc/"+mg.Metadata.Name+"-"+mg.Metadata.Version+".html", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = tmpl.Execute(f, mg)
	if err != nil {
		fmt.Println(err)
	}
}

func ResourceDotGen(mgs *[]MetaGraf, col string) {
	tmpl := template.Must(template.ParseFiles(TmplBasePath + "/digraph.dot"))
	f, err := os.OpenFile("/home/a01595/go/src/metagraf/docs/dots/"+col+".dot", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = tmpl.Execute(f, mgs)
	if err != nil {
		fmt.Println(err)
	}
}
