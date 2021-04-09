/*
Copyright 2020 The metaGraf Authors

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
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	k8sclient2 "github.com/laetho/metagraf/internal/pkg/k8sclient"
	params2 "github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	log "k8s.io/klog"
	"os"
	"strconv"
)

func GenServiceMonitorAndService(mg *metagraf.MetaGraf) {

}

func GenServiceMonitor(mg *metagraf.MetaGraf) {
	objname := Name(mg)
	// Resource labels
	l := make(map[string]string)
	l["app"] = objname
	l["app.kubernetes.io/instance"] = objname
	l["prometheus"] = params2.ServiceMonitorOperatorName

	// Selector
	s := make(map[string]string)
	s["app"] = objname

	eps := []monitoringv1.Endpoint{}
	ep := monitoringv1.Endpoint{
		TargetPort: &intstr.IntOrString{
			IntVal: FindServiceMonitorPort(mg),
		},
		Path:     FindServiceMonitorPath(mg),
		Scheme:   params2.ServiceMonitorScheme,
		Interval: params2.ServiceMonitorInterval,
	}
	eps = append(eps, ep)

	var obj = monitoringv1.ServiceMonitor{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceMonitor",
			APIVersion: "monitoring.coreos.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      objname,
			Namespace: NameSpace,
			Labels:    l,
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			Endpoints: eps,
			Selector: metav1.LabelSelector{
				MatchLabels: s,
			},
		},
	}

	if !Dryrun {
		StoreServiceMonitor(obj)
	}
	if Output {
		MarshalObject(obj.DeepCopyObject())
	}
}

// Parses metaGraf specification to look for annotation to
// control scrape path for ServiceMonitor resource.
func FindServiceMonitorPath(mg *metagraf.MetaGraf) string {
	// mg cli value, return provided if not default.
	if len(params2.ServiceMonitorPath) > 0 && params2.ServiceMonitorPath != params2.ServiceMonitorPathDefault {
		return params2.ServiceMonitorPath
	}
	// Annotation, if not provided return annotation value.
	if len(mg.Metadata.Annotations["servicemonitor.monitoring.coreos.com/path"]) > 0 {
		return mg.Metadata.Annotations["servicemonitor.monitoring.coreos.com/path"]
	}
	// Default, return default value
	return params2.ServiceMonitorPathDefault
}

// Parses metaGraf specification to look for annotation to
// control scrape port when generating ServiceMonitor resource.
func FindServiceMonitorPort(mg *metagraf.MetaGraf) int32 {
	// mg cli value, return provided if not default.
	if params2.ServiceMonitorPort > 1024 && params2.ServiceMonitorPort != params2.ServiceMonitorPortDefault {
		return params2.ServiceMonitorPort
	}
	// Annotation, if not provided return annotation value.
	if _, ok := mg.Metadata.Annotations["servicemonitor.monitoring.coreos.com/port"]; ok {
		port, err := strconv.Atoi(mg.Metadata.Annotations["servicemonitor.monitoring.coreos.com/port"])
		if err != nil {
			panic(err)
		}
		if int32(port) > 1024 {
			return int32(port)
		}
	}
	// Default, return default value
	return params2.ServiceMonitorPortDefault
}

func StoreServiceMonitor(obj monitoringv1.ServiceMonitor) {
	client := k8sclient2.GetMonitoringV1Client().ServiceMonitors(NameSpace)
	res, _ := client.Get(context.TODO(), obj.Name, metav1.GetOptions{})

	if len(res.ResourceVersion) > 0 {
		obj.ResourceVersion = res.ResourceVersion
		_, err := client.Update(context.TODO(), &obj, metav1.UpdateOptions{})
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Updated ServiceMonitor: ", obj.Name, " in Namespace: ", NameSpace)
	} else {
		_, err := client.Create(context.TODO(), &obj, metav1.CreateOptions{})
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created ServiceMonitor: ", obj.Name, " in Namespace: ", NameSpace)
	}
}

func DeleteServiceMonitor(name string) {
	client := k8sclient2.GetMonitoringV1Client().ServiceMonitors(NameSpace)

	_, err := client.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("The ServiceMonitor: ", name, "does not exist in namespace: ", NameSpace, ", skipping...")
		return
	}

	err = client.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("Unable to delete ServiceMonitor: ", name, " in namespace: ", NameSpace)
		log.Error(err)
		return
	}
	fmt.Println("Deleted ServiceMonitor: ", name, ", in namespace: ", NameSpace)
}
