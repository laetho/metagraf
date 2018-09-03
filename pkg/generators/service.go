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

package generators

import (
	"encoding/json"
	"fmt"
	"github.com/blang/semver"
	"github.com/fsouza/go-dockerclient"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"metagraf/pkg/helpers"
	"metagraf/pkg/metagraf"
	"strconv"
	"strings"
)

func GenService(mg *metagraf.MetaGraf) {
	var objname string
	var serviceports []corev1.ServicePort

	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		objname = strings.ToLower(mg.Metadata.Name)
	} else {
		objname = strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10))
	}

	// todo dockerimage should inspected once and return a pointer to the first instance inspection data
	var DockerImage string
	if len(mg.Spec.BaseRunImage) > 0 {
		DockerImage = mg.Spec.BaseRunImage
	} else if len(mg.Spec.BuildImage) > 0 {
		DockerImage = mg.Spec.BuildImage
	} else {
		DockerImage = ""
	}

	auth := docker.AuthConfiguration{
		Username: viper.GetString("user"),
		Password: viper.GetString("password"),
	}

	ImageInfo, err := helpers.DockerInspectImage(DockerImage, "latest", auth)
	if err != nil {
		// @todo should we fail here?
		// Create an empty Image struct so the following logic can run
		ImageInfo = &docker.Image{
			Config: &docker.Config{},
		}
	}

	for k := range ImageInfo.Config.ExposedPorts {
		port, _ := strconv.Atoi(k.Port())
		ContainerPort := corev1.ServicePort{
			Name:     strings.ToLower(k.Proto()) + "-" + k.Port(),
			Port:     int32(port),
			Protocol: corev1.Protocol(strings.ToUpper(k.Proto())),
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: int32(port),
				StrVal: k.Port(),
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

	ba, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(ba))

}
