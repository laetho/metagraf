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
	"fmt"
	"os"

	"github.com/laetho/metagraf/internal/pkg/k8sclient"
	"github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	log "k8s.io/klog"
)

func GenService(mg *metagraf.MetaGraf) {
	objname := Name(mg)

	selectors := make(map[string]string)
	selectors["app"] = objname

	labels := Labels(objname, labelsFromParams(params.Labels))

	obj := corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   objname,
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports:           mg.GetServicePorts(),
			Selector:        selectors,
			Type:            "ClusterIP",
			SessionAffinity: "None",
		},
	}

	if !Dryrun {
		StoreService(obj)
	}
	if Output {
		MarshalObject(obj.DeepCopyObject())
	}

	// Optinonally also create a ServiceMonitor resource.
	if params.ServiceMonitor {
		if Output && Format == "yaml" {
			fmt.Println("---")
		}
		GenServiceMonitor(mg)
	}
}


func StoreService(obj corev1.Service) {
	client := k8sclient.GetCoreClient().Services(NameSpace)
	svc, _ := client.Get(context.TODO(), obj.Name, metav1.GetOptions{})

	if len(svc.ResourceVersion) > 0 {
		obj.ResourceVersion = svc.ResourceVersion
		obj.Spec.ClusterIP = svc.Spec.ClusterIP
		_, err := client.Update(context.TODO(), &obj, metav1.UpdateOptions{})
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Updated Service: ", obj.Name, " in Namespace: ", NameSpace)
	} else {
		_, err := client.Create(context.TODO(), &obj, metav1.CreateOptions{})
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created Service: ", obj.Name, " in Namespace: ", NameSpace)
	}
}

func DeleteService(name string) {
	client := k8sclient.GetCoreClient().Services(NameSpace)

	_, err := client.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("Service: ", name, "does not exist in namespace: ", NameSpace, ", skipping...")
		return
	}

	err = client.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("Unable to delete Service: ", name, " in namespace: ", NameSpace)
		log.Error(err)
		return
	}
	fmt.Println("Deleted Service: ", name, ", in namespace: ", NameSpace)
}
