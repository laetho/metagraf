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

package modules

import (
	"fmt"
	"github.com/golang/glog"
	"html/template"
	"metagraf/pkg/metagraf"
	"os"
)

func GenRef(mg *metagraf.MetaGraf) {
	glog.Info("Fetching template: %v", Template)
	cm, err  := GetConfigMap(Template)
	if err != nil {
		glog.Error(err)
		os.Exit(-1)
	}
	tmpl, _ := template.New("refdoc").Parse(cm.Data["template"])

	filename := "/tmp/"+Name(mg)+Suffix

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = tmpl.Execute(f, mg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote ref file to: ", filename)
}
