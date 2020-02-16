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
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/util/intstr"
	"metagraf/mg/ocpclient"
	"metagraf/pkg/helpers"
	"metagraf/pkg/imageurl"
	"metagraf/pkg/metagraf"
	"strconv"
	"strings"

	//corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenDeployment(mg *metagraf.MetaGraf, namespace string) {
	objname := Name(mg)
	registry := viper.GetString("registry")

	if len(ImageNS) > 0 {
		namespace = ImageNS
	}

	// If container registry host is set and it differs from default, use value from -r (--registry) flag.
	if len(Registry) > 0 && registry != Registry {
		registry = Registry
	}

	// Resource labels
	l := make(map[string]string)
	l["app"] = objname
	l["deployment"] = objname
	// Propagate labels from metagraf to Deployment and Pod resources
	for k,v := range mg.Metadata.Labels {
		l[k]=v
	}

	// Selector
	sm := make(map[string]string)
	sm["app"] = objname
	sm["deployment"] = objname

	s := metav1.LabelSelector{
		MatchLabels: sm,
	}

	var RevisionHistoryLimit int32 = 5

	var MaxSurge intstr.IntOrString
	MaxSurge.StrVal = "25%"
	MaxSurge.Type = 1
	var MaxUnavailable intstr.IntOrString
	MaxUnavailable.StrVal = "25%"
	MaxUnavailable.Type = 1

	// Instance of RollingDeploymentStrategyParams
	rollingParams := appsv1.RollingUpdateDeployment{
		MaxSurge:            &MaxSurge,
		MaxUnavailable:      &MaxUnavailable,
	}

	var replicas int32 = 1

	// Containers
	var Containers []corev1.Container
	var ContainerPorts []corev1.ContainerPort
	//var ContainerVolumes []string
	var Volumes []corev1.Volume
	var VolumeMounts []corev1.VolumeMount
	// Environment
	var EnvVars []corev1.EnvVar

	var DockerImage string
	if len(mg.Spec.BaseRunImage) > 0 {
		DockerImage = mg.Spec.BaseRunImage
	} else if len(mg.Spec.BuildImage) > 0 {
		DockerImage = mg.Spec.BuildImage
	} else {
		DockerImage = mg.Spec.Image
	}

	var imgurl imageurl.ImageURL
	imgurl.Parse(DockerImage)

	/*
	client := ocpclient.GetImageClient()

	ist := helpers.GetImageStreamTags(
		client,
		imgurl.Namespace,
		imgurl.Image+":"+imgurl.Tag)
*/
	ImageInfo := helpers.SkopeoImageInfo(DockerImage)

	// Adding name and version of component as en environment variable
	EnvVars = append(EnvVars, corev1.EnvVar{
		Name:  "MG_APP_NAME",
		Value: MGAppName(mg),
	})

	// todo: can be removed?
	var oversion string
	if len(Version) > 0 {
		oversion = Version
	} else {
		oversion = mg.Spec.Version
	}

	EnvVars = append(EnvVars, corev1.EnvVar{
		Name:  "MG_APP_VERSION",
		Value: oversion,
	})

	// Environment Variables from baserunimage
	if BaseEnvs {
		for _, e := range ImageInfo.Config.Env {
			es := strings.Split(e, "=")
			if helpers.SliceInString(EnvBlacklistFilter, strings.ToLower(es[0])) {
				continue
			}
			EnvVars = append(EnvVars, corev1.EnvVar{Name: es[0], Value: es[1]})
		}
	}

	// Handle EnvFrom
	var EnvFrom []corev1.EnvFromSource

	// Local variables from metagraf as deployment envvars
	for _, e := range mg.Spec.Environment.Local {
		// Skip environment variable if SecretFrom
		if len(e.SecretFrom) > 0 {
			continue
		}
		if len(e.EnvFrom) > 0 {
			continue
		}
		// Use EnvToEnvVar to potentially use override values.
		EnvVars = append(EnvVars, EnvToEnvVar(&e, false))
	}

	// External variables from metagraf as deployment envvars
	for _, e := range mg.Spec.Environment.External.Consumes {
		EnvVars = append(EnvVars, EnvToEnvVar(&e, true))
	}
	for _, e := range mg.Spec.Environment.External.Introduces {
		EnvVars = append(EnvVars, EnvToEnvVar(&e, true))
	}

	// EnvVars from ConfigMaps, fetch Metagraf config resources that is of
	for _, e := range mg.Spec.Environment.Local {
		if len(e.EnvFrom) == 0 {
			continue
		}

		EnvFrom = append(EnvFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: e.EnvFrom,
				},
			},
		})
	}

	/*
		EnvVars from Secrets. Find all environment variables
		that containers the SecretFrom field and append to the
		EnvFrom as EnvFromSource->SecretRef.
	*/
	for _, e := range mg.Spec.Environment.Local {
		if len(e.SecretFrom) == 0 {
			continue
		}
		cmref := corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: e.SecretFrom,
				},
			},
		}
		EnvFrom = append(EnvFrom, cmref)
	}

	/* Norsk Tipping Specific Logic regarding
	   WLP / OpenLiberty Features. Should maybe
	   look at some plugin approach to this later.
	   todo: Add annotations from metagraf to deployment and expose them to pod using downward api.
	   info: https://kubernetes.io/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/#the-downward-api
	*/
	if len(mg.Metadata.Annotations["norsk-tipping.no/libertyfeatures"]) > 0 {
		EnvVars = append(EnvVars, corev1.EnvVar{
			Name:  "LIBERTY_FEATURES",
			Value: mg.Metadata.Annotations["norsk-tipping.no/libertyfeatures"],
		})
	}

	// Labels from baserunimage
	/*
		for k, v := range ImageInfo.Config.Labels {
			if helpers.SliceInString(LabelBlacklistFilter, strings.ToLower(k)) {
				continue
			}
			l[k] = helpers.LabelString(v)
		}
	*/

	// ContainerPorts
	for k := range ImageInfo.Config.ExposedPorts {
		ss := strings.Split(k, "/")
		port, _ := strconv.Atoi(ss[0])
		ContainerPort := corev1.ContainerPort{
			ContainerPort: int32(port),
			Protocol:      corev1.Protocol(strings.ToUpper(ss[1])),
		}
		ContainerPorts = append(ContainerPorts, ContainerPort)
	}

	Volumes, VolumeMounts = volumes(mg, ImageInfo)

	// Tying Container PodSpec together
	Container := corev1.Container{
		Name:            objname,
		Image:           registry + "/" + namespace + "/" + SpecName(mg) + ":" + Tag,
		ImagePullPolicy: corev1.PullAlways,
		Ports:           ContainerPorts,
		VolumeMounts:    VolumeMounts,
		Env:             EnvVars,
		EnvFrom:         EnvFrom,
	}
	Containers = append(Containers, Container)

	// Tying the DeploymentObject together, literally :)
	obj := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   objname,
			Labels: l,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:             &replicas,
			RevisionHistoryLimit: &RevisionHistoryLimit,
			Selector:             &s,
			Strategy: appsv1.DeploymentStrategy{
				Type:                  appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate:         &rollingParams,
			},
			Template: corev1.PodTemplateSpec{
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
		Status: appsv1.DeploymentStatus{},
	}

	if !Dryrun {
		StoreDeployment(obj)
	}
	if Output {
		MarshalObject(obj.DeepCopyObject())
	}
}


func StoreDeployment(obj appsv1.Deployment) {

	glog.Infof("ResourceVersion: %v Length: %v", obj.ResourceVersion, len(obj.ResourceVersion))
	glog.Infof("Namespace: %v", NameSpace)

	client := ocpclient.GetKubernetesClient().AppsV1().Deployments(NameSpace)
	if len(obj.ResourceVersion) > 0 {
		// update
		result, err := client.Update(&obj)
		if err != nil {
			glog.Info(err)
		}
		glog.Infof("Updated Deployment: %v(%v)", result.Name, obj.Name)
	} else {
		result, err := client.Create(&obj)
		if err != nil {
			glog.Info(err)
		}
		glog.Infof("Created Deployment: %v(%v)", result.Name, obj.Name)
	}
}
