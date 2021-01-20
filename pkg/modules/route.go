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
	"k8s.io/apimachinery/pkg/util/intstr"
	log "k8s.io/klog"
	"metagraf/pkg/helpers"
	"os"
	"sort"
	"strings"

	"metagraf/mg/k8sclient"
	"metagraf/pkg/imageurl"
	"metagraf/pkg/metagraf"

	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenRoute(mg *metagraf.MetaGraf) {
	var weight int32 = 100

	objname := Name(mg)

	var DockerImage string
	if len(mg.Spec.BaseRunImage) > 0 {
		DockerImage = mg.Spec.BaseRunImage
	} else if len(mg.Spec.BuildImage) > 0 {
		DockerImage = mg.Spec.BuildImage
	} else if len(mg.Spec.Image) > 0 {
		DockerImage = mg.Spec.Image
	} else {
		DockerImage = ""
	}

	client := k8sclient.GetImageClient()
	var imgurl imageurl.ImageURL
	err := imgurl.Parse(DockerImage)
	ist := helpers.GetImageStreamTags(
		client,
		imgurl.Namespace,
		imgurl.Image+":"+imgurl.Tag)
	if err != nil {
		log.Errorf("%v", err)
	}

	ImageInfo := helpers.GetDockerImageFromIST(ist)
	log.V(2).Infof("Docker image ports: %v", ImageInfo.Config.ExposedPorts)

	serviceports := GetServicePorts(mg, helpers.ImageExposedPortsToServicePorts(ImageInfo.Config))
	// Find http port
	// todo: This needs to be more solid
	var ports []string
	for _, port := range serviceports {
		if port.Name == "http" {
			ports = append(ports, "http")
			break
		} else {
			ports = append(ports, port.Name)
		}
	}
	sort.Strings(ports)

	l := make(map[string]string)
	l["app"] = objname

	obj := routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: "route.openshift.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   objname,
			Labels: l,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: objname,
				Weight: &weight,
			},
			Path: Context,
			Port: &routev1.RoutePort{
				TargetPort: intstr.IntOrString{
					Type: 1,
					StrVal: strings.Replace(ports[0], "/", "-", -1),

				},
			},
		},
	}

	if !Dryrun {
		StoreRoute(obj)
	}
	if Output {
		MarshalObject(obj.DeepCopyObject())
	}
}

func StoreRoute(obj routev1.Route) {
	client := k8sclient.GetRouteClient().Routes(NameSpace)
    route, _ := client.Get(context.TODO(), obj.Name, metav1.GetOptions{} )
	if len(route.ResourceVersion) > 0 {
		obj.ResourceVersion = route.ResourceVersion
		_, err := client.Update(context.TODO(), &obj, metav1.UpdateOptions{})
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Updated Route: ", obj.Name, " in Namespace: ", NameSpace)
	} else {
		_, err := client.Create(context.TODO(), &obj, metav1.CreateOptions{})
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created Route: ", obj.Name, " in Namespace: ", NameSpace)
	}
}

func DeleteRoute(name string) {
	client := k8sclient.GetRouteClient().Routes(NameSpace)

	_, err := client.Get(context.TODO(),name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("Route: ", name, "does not exist in namespace: ", NameSpace,", skipping...")
		return
	}

	err = client.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println( "Unable to delete Route: ", name, " in namespace: ", NameSpace)
		return
	}
	fmt.Println("Deleted Route: ", name, ", in namespace: ", NameSpace)
}