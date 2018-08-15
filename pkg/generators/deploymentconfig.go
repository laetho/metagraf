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
	"strconv"
	"strings"

	"github.com/blang/semver"
	"github.com/fsouza/go-dockerclient"
	"github.com/spf13/viper"

	"metagraf/pkg/helpers"
	"metagraf/pkg/metagraf"

	appsv1 "github.com/openshift/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
	l["deploymentconfig"] = objname

	// Selector
	s := make(map[string]string)
	s["app"] = objname
	s["deploymentconfig"] = objname

	var RevisionHistoryLimit int32 = 5
	var ActiveDeadlineSeconds int64 = 21600
	var TimeoutSeconds int64 = 600
	var UpdatePeriodSeconds int64 = 1
	var IntervalSeconds int64 = 1

	var MaxSurge intstr.IntOrString
	MaxSurge.StrVal = "25%"
	MaxSurge.Type = 1
	var MaxUnavailable intstr.IntOrString
	MaxUnavailable.StrVal = "25%"
	MaxUnavailable.Type = 1

	// Instance of RollingDeploymentStrategyParams
	rollingParams := appsv1.RollingDeploymentStrategyParams{
		MaxSurge:            &MaxSurge,
		MaxUnavailable:      &MaxUnavailable,
		TimeoutSeconds:      &TimeoutSeconds,
		IntervalSeconds:     &IntervalSeconds,
		UpdatePeriodSeconds: &UpdatePeriodSeconds,
	}

	// Containers
	var Containers []corev1.Container
	var ContainerPorts []corev1.ContainerPort
	//var ContainerVolumes []string
	var Volumes []corev1.Volume
	var VolumeMounts []corev1.VolumeMount
	// Environment
	var EnvVars []corev1.EnvVar

	auth := docker.AuthConfiguration{
		Username: viper.GetString("user"),
		Password: viper.GetString("password"),
	}
	ImageInfo := helpers.DockerInspectImage(mg.Spec.BaseRunImage, "latest", auth)

	for _, e := range ImageInfo.Config.Env {
		es := strings.Split(e, "=")
		if helpers.SliceInString(EnvBlacklistFilter, strings.ToLower(es[0])) {
			continue
		}
		EnvVars = append(EnvVars, corev1.EnvVar{Name: es[0], Value: es[1]})
	}
	for k, v := range ImageInfo.Config.Labels {
		if helpers.SliceInString(LabelBlacklistFilter, strings.ToLower(k)) {
			continue
		}
		l[k] = v
	}

	// ContainerPorts
	for k := range ImageInfo.Config.ExposedPorts {
		port, _ := strconv.Atoi(k.Port())
		ContainerPort := corev1.ContainerPort{
			ContainerPort: int32(port),
			Protocol:      corev1.Protocol(strings.ToUpper(k.Proto())),
		}
		ContainerPorts = append(ContainerPorts, ContainerPort)
	}

	// Volumes & VolumeMounts in podspec
	for k := range ImageInfo.Config.Volumes {
		// Volume Definitions
		Volume := corev1.Volume{
			Name: objname + helpers.PathToIdentifier(k),
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}
		Volumes = append(Volumes, Volume)

		VolumeMount := corev1.VolumeMount{
			MountPath: k,
			Name:      objname + helpers.PathToIdentifier(k),
		}
		VolumeMounts = append(VolumeMounts, VolumeMount)
	}

	// Tying Container PodSpec together
	Container := corev1.Container{
		Name:            objname,
		Image:           "registry-default.ocp.norsk-tipping.no:443/devpipeline/" + objname + ":latest",
		ImagePullPolicy: corev1.PullAlways,
		Ports:           ContainerPorts,
		VolumeMounts:    VolumeMounts,
		Env:             EnvVars,
	}
	Containers = append(Containers, Container)

	// Tying the DeploymentObject together, literally :)
	obj := appsv1.DeploymentConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DeploymentConfig",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   objname,
			Labels: l,
		},
		Spec: appsv1.DeploymentConfigSpec{
			Replicas:             0,
			RevisionHistoryLimit: &RevisionHistoryLimit,
			Selector:             s,
			Strategy: appsv1.DeploymentStrategy{
				ActiveDeadlineSeconds: &ActiveDeadlineSeconds,
				Type:          appsv1.DeploymentStrategyTypeRolling,
				RollingParams: &rollingParams,
			},
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   objname,
					Labels: l,
				},
				Spec: corev1.PodSpec{
					Containers: Containers,
					Volumes:    Volumes,
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
