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
	"context"
	"encoding/base64"
	"fmt"
	"github.com/laetho/metagraf/internal/pkg/k8sclient/k8sclient"
	"github.com/laetho/metagraf/pkg/metagraf"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	log "k8s.io/klog"
	"os"
	"strings"
)

/*
	This function will inspect the metaGraf specification
	for which configmaps it will need to create and their
	type in a map.
*/
func FindMetagrafConfigMaps(mg *metagraf.MetaGraf) map[string]string {
	maps := make(map[string]string)

	for _, c := range mg.Spec.Config {
		// Skip envRef configmaps
		if strings.ToLower(c.Type) == "envref" {
			continue
		}
		if strings.ToUpper(c.Type) == "JVM_SYS_PROP" {
			continue
		}

		if strings.ToUpper(c.Type) == "TRUSTED-CA" {
			continue
		}

		if strings.ToLower(c.Type) == "cert" {
			fmt.Println("The Config type \"cert\" is deprecated!")
			os.Exit(1)
		}
		maps[strings.ToLower(c.Name)] = "config"
	}

	for _, r := range mg.Spec.Resources {
		if r.Type == "jdbc:oracle:thin" {
			maps[strings.ToLower(r.User)] = "resource"
		}

		if len(r.TemplateRef) > 0 {
			cm, err := GetConfigMap(r.TemplateRef)
			if err != nil {
				log.Error(err)
				os.Exit(-1)
			}
			maps[cm.Name] = "template"
		}
	}

	log.Info("FindMetagrafConfigMaps(): Found", len(maps), " ConfigMaps to mount...")

	return maps
}

/*
	Fetch a ConfigMap resource the connected kubernetes cluster.
*/
func GetConfigMap(name string) (*corev1.ConfigMap, error) {
	cli := k8sclient.GetCoreClient()
	cm, err := cli.ConfigMaps(NameSpace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return cm, err
	}
	return cm, nil
}

/*
	Returns a slice of metagraf Config structs that match specific ctype string
*/
func GetMetagrafConfigsByType(mg *metagraf.MetaGraf, ctype string) []metagraf.Config {
	configs := []metagraf.Config{}

	for _, c := range mg.Spec.Config {
		if strings.ToLower(c.Type) == strings.ToLower(ctype) {
			configs = append(configs, c)
		}
	}
	return configs
}

/*
	Entry function for creating a slew of configmaps, this will be somewhat
	specific to NT internal workings for now.
*/
func GenConfigMaps(mg *metagraf.MetaGraf) {
	log.Info("GenConfigMaps: Handle", len(mg.Spec.Config), " configs...")
	for _, c := range mg.Spec.Config {
		if c.Type != "parameters" {
			continue
		}
		genConfigMapsFromConfig(&c, mg)
	}
	//genConfigMapsFromResources(mg)
}

/*
	Generates a configmap for jvm.params file for Liberty java apps
*/
func genConfigMapsFromConfig(conf *metagraf.Config, mg *metagraf.MetaGraf) {

	if conf.Global {
		return
	}

	objname := Name(mg)

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
	cm.ObjectMeta.Labels = l

	// Check if configmap should be treated as a global name.
	if conf.Global {
		cm.Name = conf.Name
	} else {
		cm.Name = objname + "-" + strings.Replace(strings.ToLower(conf.Name), ".", "-", -1)
	}

	for _, o := range conf.Options {
		prop := Variables[conf.Name+"|"+o.Name]
		if len(o.SecretFrom) > 0 {
			sec, err := GetSecret(o.SecretFrom)
			if err != nil {
				log.Error(err)
			}
			cm.Data[o.Name] = base64.StdEncoding.EncodeToString(sec.Data[o.SecretFrom])
		} else {
			if len(prop.Value) > 0 {
				cm.Data[prop.Key] = prop.Value
			} else {
				cm.Data[prop.Key] = prop.Default
			}
		}
	}

	if !Dryrun {
		StoreConfigMap(cm)
	}
	if Output {
		MarshalObject(cm.DeepCopyObject())
	}

}

/*
	todo: redo this to use template in spec to do this.
*/
func genConfigMapsFromResources(mg *metagraf.MetaGraf) {
	return
	// objname := Name(mg)

	/*
		for _, r := range mg.Spec.Resources {
			if strings.Contains(r.Type, "oracle") {
				cm := genJDBCOracle(objname, &r)
				if !Dryrun {
					StoreConfigMap(cm)
				}
				if Output {
					MarshalObject(cm.DeepCopyObject())
				}
			}
		}
	*/
}

func StoreConfigMap(m corev1.ConfigMap) {
	cmclient := k8sclient.GetCoreClient().ConfigMaps(NameSpace)
	cm, _ := cmclient.Get(context.TODO(), m.Name, metav1.GetOptions{})

	if len(cm.ResourceVersion) > 0 {
		m.ResourceVersion = cm.ResourceVersion
		_, err := cmclient.Update(context.TODO(), &m, metav1.UpdateOptions{})
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Updated ConfigMap: ", m.Name, " in Namespace: ", NameSpace)
	} else {
		_, err := cmclient.Create(context.TODO(), &m, metav1.CreateOptions{})
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created ConfigMap: ", m.Name, " in Namespace: ", NameSpace)
	}
}

func DeleteConfigMaps(mg *metagraf.MetaGraf) {
	obname := Name(mg)

	for _, c := range mg.Spec.Config {
		// Do not delete global configuration maps.
		if c.Global == true {
			continue
		}

		name := strings.ToLower(obname + "-" + c.Name)
		name = strings.Replace(name, "_", "-", -1)
		name = strings.Replace(name, ".", "-", -1)
		DeleteConfigMap(name)
	}
}

func DeleteConfigMap(name string) {
	client := k8sclient.GetCoreClient().ConfigMaps(NameSpace)

	_, err := client.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("ConfigMap: ", name, "does not exist in namespace: ", NameSpace, ", skipping...")
		return
	}

	err = client.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("Unable to delete ConfigMap: ", name, " in namespace: ", NameSpace)
		log.Error(err)
		return
	}
	fmt.Println("Deleted Configmap: ", name, ", in namespace: ", NameSpace)
}
