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
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"metagraf/pkg/metagraf"
)

func GenSecrets(mg *metagraf.MetaGraf) {
	for _,r := range mg.Spec.Resources {
		if len(r.SecretRef) == 0 && len(r.User) == 0 {
			continue
		}
		obj := genResourceSecret(r)
		ba, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(ba))

	}
}

func genResourceSecret(res metagraf.Resource) *corev1.Secret {
	sec := corev1.Secret{}

	return &sec
}



func InspectSecrets(mg *metagraf.MetaGraf) {

	for _,r := range mg.Spec.Resources {
		if len(r.SecretRef) == 0 && len(r.User) > 0 {
			fmt.Println(mg.Metadata.Name, "needs implicit secret for", r.User)
		}
	}

	for _,r := range mg.Spec.Config {
		if r.Type != "secret" {
			continue
		}
		fmt.Println(mg.Metadata.Name, "needs the following secret", r.Name)
	}

}