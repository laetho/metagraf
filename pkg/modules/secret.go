/*
Copyright 2019 The metaGraf Authors

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
	log "k8s.io/klog"
	"metagraf/pkg/metagraf"
	"os"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"metagraf/mg/k8sclient"
)

func FindSecrets(mg *metagraf.MetaGraf) map[string]string {
	maps := make(map[string]string)

	for _, r := range mg.Spec.Resources {
		if len(r.User) > 0 {
			maps[strings.ToLower(strings.Replace(r.Name,"_","-", -1))+"-"+strings.ToLower(r.User)] = "password"
		}
	}

	for _, s := range mg.Spec.Secret {
		maps[strings.ToLower(s.Name)] = strconv.FormatBool(s.Global)
	}

	return maps
}

// todo: flag creation of resource secrets?
func GenSecrets(mg *metagraf.MetaGraf) {
	for _, r := range mg.Spec.Resources {

		// Is secret generation necessary?
		if len(r.Secret) == 0 && len(r.User) == 0 {
			continue
		}

		// Do not create secret if it already exist!
		if secretExists(ResourceSecretName(&r)) {
			fmt.Println("Skipping resource: ", r.Name, "already exist in Namespace: ", NameSpace)
			continue
		}

		obj := genResourceSecret(&r, mg)
		if !Dryrun {
			StoreSecret(*obj)
		}
		if Output {
			MarshalObject(obj.DeepCopyObject())
		}
	}

	for _, s := range mg.Spec.Secret{
		// Do not create global secrets unless CreateGlobals is true.
		if s.Global == true && !CreateGlobals {
			log.Info("Skipping creation of global secret named: "+ strings.ToLower(s.Name))
			continue
		}

		if secretExists(strings.ToLower(s.Name)) {
			log.Info("Skipping secret: ", Name(mg)+"-"+strings.ToLower(s.Name))
			continue
		}

		obj := genSecret(&s, mg)
		if !Dryrun {
			StoreSecret(*obj)
		}
		if Output{
			MarshalObject(obj.DeepCopyObject())
		}
	}
}

// Check if a named secret exsist in the current namespace.
func secretExists(name string) bool {
	cli := k8sclient.GetCoreClient()
	obj, err := cli.Secrets(NameSpace).Get(name,metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		return false
	}
	log.Info("Secret ", obj.Name, " exists in namespace: ", NameSpace)
	return true
}

//
func GetSecret(name string) (*corev1.Secret, error) {
	cli := k8sclient.GetCoreClient()
	sec, err := cli.Secrets(NameSpace).Get(name, metav1.GetOptions{})
	if err != nil {
		return sec, err
	}
	return sec, nil
}

func genSecret(s *metagraf.Secret, mg *metagraf.MetaGraf) *corev1.Secret {
	objname := Name(mg)

	// Resource labels
	l := make(map[string]string)

	if s.Global == true {
		l["name"] = strings.ToLower(s.Name)
	} else {
		l["name"] = objname+"-"+strings.ToLower(s.Name)
	}

	// Populate v1.Secret StringData and Data
	stringdata := make(map[string]string)
	data := make(map[string][]byte)
		name := strings.Replace(s.Name, ".", "_", -1)
		data[name] = []byte("data")
		stringdata[name] = "data"

	sec := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   l["name"],
			Labels: l,
		},
		Type:       "opaque",
		StringData: stringdata,
		Data:       data,
	}

	return &sec
}

/*
 *	Generates implicit secrets based on fields in the resource section.
 */
func genResourceSecret(res *metagraf.Resource, mg *metagraf.MetaGraf) *corev1.Secret {

	objname := Name(mg)

	// Resource labels
	l := make(map[string]string)
	l["name"] = ResourceSecretName(res)
	l["app"] = objname

	// Populate v1.Secret StringData and Data
	stringdata := make(map[string]string)
	data := make(map[string][]byte)

	if len(res.User) > 0 {
		stringdata["type"] = res.Type
		stringdata["templateref"] = res.TemplateRef
		stringdata["user"] = res.User
		stringdata["password"] = "replaceme"
	}

	//if len(res.Secret) > 0 && res.SecretType == "cert" {
	// 	data[res.Secret] = []byte("Replace this")
	//}

	sec := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   l["name"],
			Labels: l,
		},
		Type:       "opaque",
		StringData: stringdata,
		Data:       data,
	}

	return &sec
}

func StoreSecret(obj corev1.Secret) {
	client := k8sclient.GetCoreClient().Secrets(NameSpace)
	sec, err := client.Get(obj.Name, metav1.GetOptions{})
	if err != nil {
		log.Infof("Could not fetch Secret: %v", err)
	}
	if len(sec.ResourceVersion) > 0 {
		result, err := client.Update(&obj)
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Updated Secret: %v(%v)", result.Name, obj.Name)
	} else {
		result, err := client.Create(&obj)
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		log.Infof("Created Secret: %v(%v)", result.Name, obj.Name)
	}
}

func DeleteSecrets(mg *metagraf.MetaGraf) {
	obname := Name(mg)

	for _, s := range mg.Spec.Secret {
		// Do not delete global configuration maps.
		if s.Global == true {
			continue
		}

		name := strings.ToLower(obname+"-"+s.Name)
		name = strings.Replace(name, "_", "-", -1)
		name = strings.Replace(name, ".", "-", -1)
		DeleteSecret(name)
	}
}

func DeleteSecret(name string) {
	client := k8sclient.GetCoreClient().Secrets(NameSpace)

	_, err := client.Get(name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("Secret: ", name, "does not exist in namespace: ", NameSpace,", skipping...")
		return
	}

	err = client.Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("Unable to delete Secret: ", name, " in namespace: ", NameSpace)
		log.Error(err)
		return
	}
	fmt.Println("Deleted Secret: ", name, ", in namespace: ", NameSpace)
}
