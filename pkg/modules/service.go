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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"metagraf/mg/ocpclient"
	"metagraf/pkg/helpers"
	"metagraf/pkg/imageurl"
	"metagraf/pkg/metagraf"
	"os"
	"strconv"
	"strings"
)

func GenService(mg *metagraf.MetaGraf) {
	objname := Name(mg)

	var serviceports []corev1.ServicePort

	// todo dockerimage should inspected once and return a pointer to the first instance inspection data
	var DockerImage string
	if len(mg.Spec.BaseRunImage) > 0 {
		DockerImage = mg.Spec.BaseRunImage
	} else if len(mg.Spec.BuildImage) > 0 {
		DockerImage = mg.Spec.BuildImage
	} else {
		DockerImage = ""
	}

	var imgurl imageurl.ImageURL
	imgurl.Parse(DockerImage)

	client := ocpclient.GetImageClient()

	ist := helpers.GetImageStreamTags(
		client,
		imgurl.Namespace,
		imgurl.Image+":"+imgurl.Tag)

	ImageInfo := helpers.GetDockerImageFromIST(ist)

	for k := range ImageInfo.Config.ExposedPorts {
		ss := strings.Split(k, "/")
		port, _ := strconv.Atoi(ss[0])
		ContainerPort := corev1.ServicePort{
			Name:     strings.ToLower(ss[0]) + "-" + ss[1],
			Port:     int32(port),
			Protocol: corev1.Protocol(strings.ToUpper(ss[1])),
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: int32(port),
				StrVal: ss[1],
			},
		}
		serviceports = append(serviceports, ContainerPort)
	}

	selectors := make(map[string]string)
	selectors["deploymentconfig"] = objname

	labels := make(map[string]string)
	labels["app"] = objname

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
			Ports:           serviceports,
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

}

func StoreService(obj corev1.Service) {

	glog.Infof("ResourceVersion: %v Length: %v", obj.ResourceVersion, len(obj.ResourceVersion))
	glog.Infof("Namespace: %v", NameSpace)

	client := ocpclient.GetCoreClient().Services(NameSpace)

	if len(obj.ResourceVersion) > 0 {
		// update
		result, err := client.Update(&obj)
		if err != nil {
			glog.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		glog.Infof("Updated Service: %v(%v)", result.Name, obj.Name)
	} else {
		result, err := client.Create(&obj)
		if err != nil {
			glog.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		glog.Infof("Created Service: %v(%v)", result.Name, obj.Name)
	}
}
