/*
Copyright 2018-2020 The metaGraf Authors

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
	"strings"

	"github.com/laetho/metagraf/internal/pkg/helpers"
	"github.com/laetho/metagraf/internal/pkg/imageurl"
	"github.com/laetho/metagraf/internal/pkg/k8sclient"
	"github.com/laetho/metagraf/internal/pkg/params"
	log "k8s.io/klog"

	"github.com/laetho/metagraf/pkg/metagraf"

	buildv1 "github.com/openshift/api/build/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"github.com/openshift/oc/pkg/helpers/source-to-image/tar"
)

func GenBuildConfig(mg *metagraf.MetaGraf) {

	imgurl := imageurl.NewImageUrl(mg.Spec.BuildImage)
	imageStreamTagObjRef := constructImageStreamTagObjRef(imgurl)

	containerEnvVarsMap := make(map[string]corev1.EnvVar)
	if BaseEnvs {
		imageEnvVars := fetchEnvironmentVariablesFromBuildImage(imgurl)
		envVarsFromBaseImage := mapEnvVarsFromBaseImage(imageEnvVars)

		for _, envVar := range envVarsFromBaseImage {
			containerEnvVarsMap[envVar.Name] = envVar
		}

	}

	envVarsFromMetaGraf := mapBuildParamsFromMetaGraf(mg)
	for _, envVar := range envVarsFromMetaGraf {
		containerEnvVarsMap[envVar.Name] = envVar
	}

	envVarsFromBuildParams := mapBuildParamsFromCommandline(params.BuildParams)
	for _, envVar := range envVarsFromBuildParams {
		containerEnvVarsMap[envVar.Name] = envVar
	}


	containerEnvVars := make([]corev1.EnvVar, 0, len(containerEnvVarsMap))

	for  _, value := range containerEnvVarsMap {
		containerEnvVars = append(containerEnvVars, value)
	}

	// Construct toObjRef for BuildConfig output overrides
	// Resource labels
	objname := Name(mg)
	labels := Labels(objname, labelsFromParams(params.Labels))
	var toObjRefName = objname
	var toObjRefTag = "latest"
	if len(params.OutputImagestream) > 0 {
		toObjRefName = params.OutputImagestream
	}
	if len(Tag) > 0 {
		toObjRefTag = Tag
	}
	var toObjRef = &corev1.ObjectReference{
		Kind: "ImageStreamTag",
		Name: toObjRefName + ":" + toObjRefTag,
	}

	bc := buildv1.BuildConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "BuildConfig",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   objname,
			Labels: labels,
		},
		Spec: buildv1.BuildConfigSpec{
			RunPolicy: buildv1.BuildRunPolicySerial,
			CommonSpec: buildv1.CommonSpec{
				Source: ConstructBuildSource(mg),
				Strategy: buildv1.BuildStrategy{
					Type: buildv1.SourceBuildStrategyType,
					SourceStrategy: &buildv1.SourceBuildStrategy{
						Env:  containerEnvVars,
						From: imageStreamTagObjRef,
					},
				},
				Output: buildv1.BuildOutput{
					To: toObjRef,
				},
			},
		},
	}

	if !Dryrun {
		StoreBuildConfig(bc)
	}
	if Output {
		MarshalObject(bc.DeepCopyObject())
	}
}

func mapBuildParamsFromCommandline(buildParams []string) []corev1.EnvVar {
	var containerEnvVars []corev1.EnvVar
	for _, envVarRaw := range buildParams {
		if strings.Contains(envVarRaw, "=") {
			keyVal := strings.SplitN(envVarRaw, "=", 2) // Use strings.SplitN with n=2 to limit the result to two substrings.
			key := keyVal[0]
			val := keyVal[1]
			envVar := corev1.EnvVar{Name: key, Value: val}
			containerEnvVars = append(containerEnvVars, envVar)
			log.Infof("Added build param %s", envVar)
		} else {
			log.Warningf("BuildParam [%s] doesn't fit expected pattern of KEY=VAL!", envVarRaw)
		}
	}

	return containerEnvVars
}

func mapBuildParamsFromMetaGraf(mg *metagraf.MetaGraf) []corev1.EnvVar {
	var containerEnvVars []corev1.EnvVar
	configPropertiesFromMetaGraf := Variables.KeyMap()
	for _, buildEnvVar := range mg.Spec.Environment.Build {
		if buildEnvVar.Required == true {
			if len(configPropertiesFromMetaGraf[buildEnvVar.Name]) > 0 {
				containerEnvVars = append(containerEnvVars, corev1.EnvVar{Name: buildEnvVar.Name, Value: configPropertiesFromMetaGraf[buildEnvVar.Name]})
			} else {
				containerEnvVars = append(containerEnvVars, corev1.EnvVar{Name: buildEnvVar.Name, Value: buildEnvVar.Default})
			}
		} else if buildEnvVar.Required == false {
			containerEnvVars = append(containerEnvVars, corev1.EnvVar{Name: buildEnvVar.Name, Value: "null"})
		}
	}
	return containerEnvVars
}

func fetchEnvironmentVariablesFromBuildImage(imgurl imageurl.ImageURL) []string {
	client := k8sclient.GetImageClient()
	ist := helpers.GetImageStreamTags(
		client,
		imgurl.Namespace,
		imgurl.Image+":"+imgurl.Tag)
	dockerImage := helpers.GetDockerImageFromIST(ist)
	// Environment Variables from buildimage
	imageEnvs := dockerImage.Config.Env
	return imageEnvs
}

func constructImageStreamTagObjRef(imgurl imageurl.ImageURL) corev1.ObjectReference {
	imageStreamTag := corev1.ObjectReference{
		Kind:      "ImageStreamTag",
		Namespace: imgurl.Namespace,
		Name:      imgurl.Image + ":" + imgurl.Tag,
	}
	return imageStreamTag
}

func mapEnvVarsFromBaseImage(imageEnvs []string) []corev1.EnvVar {
	log.V(2).Info("Populate environment variables form base image.")

	var containerEnvVars []corev1.EnvVar
	for _, e := range imageEnvs {
		es := strings.Split(e, "=")
		if helpers.SliceInString(EnvBlacklistFilter, strings.ToLower(es[0])) {
			continue
		}
		containerEnvVars = append(containerEnvVars, corev1.EnvVar{Name: es[0], Value: es[1]})
	}
	return containerEnvVars
}

func StoreBuildConfig(obj buildv1.BuildConfig) {
	client := k8sclient.GetBuildClient().BuildConfigs(params.NameSpace)
	bc, _ := client.Get(context.TODO(), obj.Name, metav1.GetOptions{})

	if len(bc.ResourceVersion) > 0 {
		obj.ResourceVersion = bc.ResourceVersion
		_, err := client.Update(context.TODO(), &obj, metav1.UpdateOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Updated BuildConfig: ", obj.Name, " in Namespace: ", params.NameSpace)
	} else {
		_, err := client.Create(context.TODO(), &obj, metav1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created BuildConfig: ", obj.Name, " in Namespace: ", params.NameSpace)
	}
}

func DeleteBuildConfig(name string) {
	client := k8sclient.GetBuildClient().BuildConfigs(NameSpace)

	_, err := client.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("The BuildConfig: ", name, "does not exist in namespace: ", NameSpace, ", skipping...")
		return
	}

	err = client.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("Unable to delete BuildConfig: ", name, " in namespace: ", NameSpace)
		log.Error(err)
		return
	}
	fmt.Println("Deleted BuildConfig: ", name, ", in namespace: ", NameSpace)
}
