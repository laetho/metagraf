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
	"encoding/json"
	"fmt"
	"github.com/blang/semver"
	"metagraf/pkg/metagraf"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
Entry function for creating a slew of configmaps, this will be somewhat
specific to NT internal workings for now.
*/
func GenConfigMaps(mg *metagraf.MetaGraf) {

	/*
		We need to create the following ConfigMaps:
			* INTPL.Config.properties
			* jvm.params
			* server.xml
	*/

	for _, c := range mg.Spec.Config {
		if c.Type != "parameters" {
			continue
		}
		genConfigMapFromConfig(&c, mg)
	}

}

/*
Generates a configmap for jvm.params file for Liberty java apps
*/
func genConfigMapFromConfig(conf *metagraf.Config, mg *metagraf.MetaGraf) {
	var objname string
	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		objname = strings.ToLower(mg.Metadata.Name)
	} else {
		objname = strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10))
	}

	l := make(map[string]string)
	l["app"] = objname

	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   objname,
			Labels: l,
		},
	}

	// Need to initialize the map defined in struct
	// Should be done in a factory maybe.
	cm.Data = make(map[string]string)
	//cm.TypeMeta.Kind = "ConfigMap"
	//cm.TypeMeta.APIVersion = "v1"
	cm.ObjectMeta.Labels = l

	cm.Name = strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10) + "-" + strings.Replace(conf.Name, ".", "-", -1))
	for _, o := range conf.Options {
		if ValueFromEnv(o.Name) {
			cm.Data[o.Name] = Variables[o.Name]
		} else {
			cm.Data[o.Name] = o.Default
		}
	}

	b, err := json.Marshal(cm)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

