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
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"metagraf/pkg/metagraf"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"metagraf/mg/ocpclient"
)

func GenSecrets(mg *metagraf.MetaGraf) {
	for _,r := range mg.Spec.Resources {
		// Is secret generation necessary?
		if len(r.SecretRef) == 0 && len(r.User) == 0 {
			continue
		}

		// Do not create secret if it already exist!
		if secretExists(ResourceSecretName(&r)) {
			continue
		}

		obj := genResourceSecret(&r, mg)
		ba, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(ba))
	}
}

// Check if a named secret exsist in the current namespace.
func secretExists(name string) bool {
	cli := ocpclient.GetCoreClient()
	l, err := cli.Secrets(NameSpace).List(metav1.ListOptions{LabelSelector:"name = "+name})

	if err != nil{
		glog.Error(err)
		os.Exit(1)
	}

	if len(l.Items) > 0 {
		return true
	}
	return false
}

func genResourceSecret(res *metagraf.Resource, mg *metagraf.MetaGraf) *corev1.Secret {

	objname := Name(mg)

	// Resource labels
	l := make(map[string]string)
	l["name"] = ResourceSecretName(res)
	l["app"] = objname


	sec := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: objname,
			Labels: l,
		},
	}

	return &sec
}



func InspectSecrets(mg *metagraf.MetaGraf) {

	for _,r := range mg.Spec.Resources {
		fmt.Println("Resource type:", r.Type)
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