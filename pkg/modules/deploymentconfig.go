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
	"github.com/openshift/api/image/docker10"
	"github.com/spf13/viper"
	"metagraf/mg/ocpclient"
	"os"
	"strconv"
	"strings"

	"metagraf/pkg/helpers"
	"metagraf/pkg/imageurl"
	"metagraf/pkg/metagraf"

	appsv1 "github.com/openshift/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)


// Todo: Still needs to be split up, but some refactoring has been done.
func GenDeploymentConfig(mg *metagraf.MetaGraf, namespace string) {
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
		EnvFrom:		 EnvFrom,
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
			Replicas:             1,
			RevisionHistoryLimit: &RevisionHistoryLimit,
			Selector:             s,
			Strategy: appsv1.DeploymentStrategy{
				ActiveDeadlineSeconds: &ActiveDeadlineSeconds,
				Type:                  appsv1.DeploymentStrategyTypeRolling,
				RollingParams:         &rollingParams,
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

	if !Dryrun {
		StoreDeploymentConfig(obj)
	}
	if Output {
		MarshalObject(obj.DeepCopyObject())
	}

}

func envs(mg *metagraf.MetaGraf)  {

}

/**
	Builds up slices of corev1.Volume and corev1.VolumeMount structs and returns them.
	Should maybe consider splitting this up even further.
 */
func volumes(mg *metagraf.MetaGraf, ImageInfo *docker10.DockerImage ) ([]corev1.Volume, []corev1.VolumeMount) {
	objname := Name(mg)
	var Volumes []corev1.Volume
	var VolumeMounts []corev1.VolumeMount

	// Volumes & VolumeMounts from base image into podspec
	glog.Info("ImageInfo: Got ", len(ImageInfo.Config.Volumes), " volumes from base image...")
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

	// Put ConfigMap volumes and mounts into PodSpec
	for n, t := range FindMetagrafConfigMaps(mg) {
		var mode int32 = 420
		var vname string
		var oname string

		vname = "cm-"+strings.Replace(n,".","-", -1)

		glog.V(2).Infof("Name,Type: %v,%v", n,t)

		if t == "template" {
			oname =  strings.Replace(n,".","-", -1)
		} else {
			oname = objname+"-"+strings.Replace(n,".","-", -1)
		}

		vol := corev1.Volume{
			Name: vname,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: oname,
					},
					DefaultMode: &mode,
				},
			},
		}

		volm := corev1.VolumeMount{}
		volm.Name = vname

		// Special handling of resource because of old hackish handling of oracle jdbc connections
		if t == "resource" {
			volm.MountPath = "/mg/" + n
		} else {
			volm.MountPath = "/mg/"+t+"/"+n
		}

		Volumes = append(Volumes, vol)
		VolumeMounts = append(VolumeMounts, volm)
	}

	for n, t := range FindSecrets(mg) {
		glog.V(2).Infof("Secret: %v,%t", n, t)
		voln := strings.Replace(n,".", "-", -1)
		var mode int32 = 420
		vol := corev1.Volume{
			Name: voln,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: n,
					DefaultMode: &mode,
				},
			},
		}
		volm := corev1.VolumeMount{
			Name: voln,
			MountPath: "/mg/secret/"+n,
		}
		Volumes = append(Volumes, vol)
		VolumeMounts = append(VolumeMounts, volm)
	}

	return Volumes, VolumeMounts
}


func StoreDeploymentConfig(obj appsv1.DeploymentConfig) {
	client := ocpclient.GetAppsClient().DeploymentConfigs(NameSpace)
	dc, _ := client.Get(obj.Name, metav1.GetOptions{})

	if len(dc.ResourceVersion) > 0 {
		obj.ResourceVersion = dc.ResourceVersion
		_, err := client.Update(&obj)
		if err != nil {
			glog.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Updated DeploymentConfig: ", obj.Name," in Namespace: ", obj.Name)
	} else {
		result, err := client.Create(&obj)
		if err != nil {
			glog.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created DeploymentConfig: ", result.Name," in Namespace: ", obj.Name)
	}
}

func DeleteDeploymentConfig(name string) {
	client := ocpclient.GetAppsClient().DeploymentConfigs(NameSpace)

	_ , err := client.Get(name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("DeploymentConfig: ", name, "does not exist in namespace: ", NameSpace,", skipping...")
		return
	}

	err = client.Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		fmt.Println( "Service to delete DeploymentConfig: ", name, " in namespace: ", NameSpace)
		glog.Error(err)
		return
	}
	fmt.Println("Deleted DeploymentConfig: ", name, ", in namespace: ", NameSpace)
}