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
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	"github.com/blang/semver"

	"metagraf/pkg/metagraf"
	"k8s.io/apimachinery/pkg/util/intstr"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "github.com/openshift/api/apps/v1"

)


func GenDeploymentConfig(mg *metagraf.MetaGraf) {
	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		fmt.Println(err)
	}

	objname := strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10))

	// Resource labels
	l := make(map[string]string)
	l["app"] = objname

	s := make(map[string]string)
	s["app"] = objname
	s["deploymentconfig"] = objname

	var RevisionHistoryLimit int32 = 5
	var ActiveDeadlineSeconds int64 = 21600
	var TimeoutSeconds int64 = 600
	var UpdatePeriodSeconds int64 = 1
	var IntervalSeconds	int64 = 1

	var MaxSurge intstr.IntOrString
	MaxSurge.StrVal = "25%"
	MaxSurge.Type = 1
	var MaxUnavailable intstr.IntOrString
	MaxUnavailable.StrVal = "25%"
	MaxUnavailable.Type = 1

	// Instance of RollingDeploymentStrategyParams
	rollingParams := appsv1.RollingDeploymentStrategyParams{
		MaxSurge: &MaxSurge,
		MaxUnavailable: &MaxUnavailable,
		TimeoutSeconds: &TimeoutSeconds,
		IntervalSeconds: &IntervalSeconds,
		UpdatePeriodSeconds: &UpdatePeriodSeconds,
	}




	// Containers
	var Containers []corev1.Container
	var ContainerPorts []corev1.ContainerPort

	// Build Container ports, @todo this should be done by inspecting the image (docker inspect)
	ContainerPort := corev1.ContainerPort{
		Name: objname,
		ContainerPort: 8080,
		Protocol: "TCP",
	}
	ContainerPorts = append(ContainerPorts, ContainerPort)

	Container := corev1.Container{
		Name: objname,
		Image: "docker-registry.default.svc:5000/devpipeline/customeridentity:latest",
		ImagePullPolicy: corev1.PullAlways,
		Ports: ContainerPorts,
	}
	Containers = append( Containers, Container)


	obj := appsv1.DeploymentConfig{
		TypeMeta: metav1.TypeMeta{
			Kind: "DeploymentConfig",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: objname,
			Labels: l,
		},
		Spec: appsv1.DeploymentConfigSpec{
			Replicas: 0,
			RevisionHistoryLimit: &RevisionHistoryLimit,
			Selector: s,
			Strategy: appsv1.DeploymentStrategy{
				ActiveDeadlineSeconds: &ActiveDeadlineSeconds,
				Type: appsv1.DeploymentStrategyTypeRolling,
				RollingParams: &rollingParams,
			},
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: objname,
					Labels: l,
				},
				Spec: corev1.PodSpec{
					Containers: Containers,
				},
			},



		},
		Status: appsv1.DeploymentConfigStatus{},
	}

	ba, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(ba))

}