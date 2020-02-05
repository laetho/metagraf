/*
Copyright 2019 The MetaGraph Authors

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
	log "k8s.io/klog"
	"metagraf/pkg/metagraf"
	"os"
	"strings"
)

// Returns a metagraf.MGProperties map of parameters mainly for use in
// template processing. It filters out JVM_SYS_PROP type.
func getParamMap(mg *metagraf.MetaGraf) metagraf.MGProperties {
	params := metagraf.MGProperties{}

	for _, e := range mg.Spec.Environment.Local {
		if e.Type == "JVM_SYS_PROP" { continue }
		p := metagraf.MGProperty{
			Source:   "local",
			Key:      e.Name,
			Value:    "",
			Required: e.Required,
			Default:  e.Default,
		}
		params[p.Key] = p
	}

	for _, c := range mg.Spec.Config {
		if c.Name != "JVM_SYS_PROP" || c.Name != "jvm.options"  { continue }
		for _, o := range c.Options {
			p := metagraf.MGProperty{
				Source:   c.Name,
				Key:      o.Name,
				Value:    "",
				Required: o.Required,
				Default:  o.Default,
			}
			params[p.Key] = p
		}
	}

	return params
}
// Filters out only required property elements.
func getReqParamMap(p metagraf.MGProperties) metagraf.MGProperties {
	params := metagraf.MGProperties{}
	for _, p := range p {
		if p.Required {
			params[p.Key] = p
		}
	}
	return params
}

func GenRef(mg *metagraf.MetaGraf) {
	log.Info("Fetching template: %v", Template)
	cm, err  := GetConfigMap(Template)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	tmpl, err := template.New("refdoc").Funcs(
		template.FuncMap{
			"getParamMap": getParamMap(mg),
			"getReqParamMap": getReqParamMap(getParamMap(mg)),
			"split": func(s string, d string) []string {
				return strings.Split(s, d)
			},
			"numOfLocal": func(l []metagraf.EnvironmentVar) int {
				return len(l)
			},
			"numOfOptions": func(l []metagraf.ConfigParam) int {
				return len(l)
			},
			"isLast": func(a []string, k int) bool {
				if (len(a)-1 == k) {
					return true
				}
				return false
			},
			"last": func(t int, c int) bool {
				if (t-1 == c) {
					return true
				}
				return false
			},
		}).Parse(cm.Data["template"])
	if (err != nil) {
		log.Error(err)
		os.Exit(1)
	}

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
