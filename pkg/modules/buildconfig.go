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
	"fmt"
	log "k8s.io/klog"
	"metagraf/mg/k8sclient"
	"metagraf/mg/params"
	"metagraf/pkg/helpers"
	"os"
	"strings"

	"metagraf/pkg/imageurl"
	"metagraf/pkg/metagraf"

	buildv1 "github.com/openshift/api/build/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenBuildConfig(mg *metagraf.MetaGraf) {
	var buildsource buildv1.BuildSource
	var imgurl imageurl.ImageURL
	var EnvVars []corev1.EnvVar

	err := imgurl.Parse(mg.Spec.BuildImage)
	if err != nil {

		log.Errorf("Malformed BuildImage url provided in metaGraf file; %v", mg.Spec.BuildImage)
		os.Exit(1)
	}

	objname := Name(mg)

	if len(mg.Spec.BaseRunImage) > 0 && len(mg.Spec.Repository) > 0 {
		buildsource = genBinaryBuildSource()
	} else if len(mg.Spec.BuildImage) > 0 && len(mg.Spec.BaseRunImage) < 1 {
		buildsource = genGitBuildSource(mg)
	}

	if BaseEnvs {
		log.V(2).Info("Populate environment variables form base image.")
		client := k8sclient.GetImageClient()

		ist := helpers.GetImageStreamTags(
			client,
			imgurl.Namespace,
			imgurl.Image+":"+imgurl.Tag)

		ImageInfo := helpers.GetDockerImageFromIST(ist)

		// Environment Variables from buildimage
		for _, e := range ImageInfo.Config.Env {
			es := strings.Split(e, "=")
			if helpers.SliceInString(EnvBlacklistFilter, strings.ToLower(es[0])) {
				continue
			}
			EnvVars = append(EnvVars, corev1.EnvVar{Name: es[0], Value: es[1]})
		}
	}

	// Resource labels
	l := make(map[string]string)
	l["app"] = objname
	l["deploymentconfig"] = objname


	km := Variables.KeyMap()
	for _, e := range mg.Spec.Environment.Build {
		if e.Required == true {
			if len(km[e.Name]) > 0 {
				EnvVars = append(EnvVars, corev1.EnvVar{Name: e.Name, Value: km[e.Name]})
			} else {
				EnvVars = append(EnvVars, corev1.EnvVar{Name: e.Name, Value: e.Default})
			}
		} else if e.Required == false {
			EnvVars = append(EnvVars, corev1.EnvVar{Name: e.Name, Value: "null"})
		}
	}

	// Construct toObjRef for BuildConfig output overrides
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
		Name: toObjRefName +":"+ toObjRefTag,
	}

	bc := buildv1.BuildConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "BuildConfig",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: objname,
			Labels: l,
		},
		Spec: buildv1.BuildConfigSpec{
			Triggers: []buildv1.BuildTriggerPolicy{
				{
					Type: "generic",
					GenericWebHook: &buildv1.WebHookTrigger{
						Secret: "wakkawakkawakka",
					},
				},
			},
			RunPolicy: buildv1.BuildRunPolicySerial,
			CommonSpec: buildv1.CommonSpec{
				Source: buildsource,
				Strategy: buildv1.BuildStrategy{
					Type: buildv1.SourceBuildStrategyType,
					SourceStrategy: &buildv1.SourceBuildStrategy{
						Env: EnvVars,
						From: corev1.ObjectReference{
							Kind:      "ImageStreamTag",
							Namespace: imgurl.Namespace,
							Name:      imgurl.Image + ":" + imgurl.Tag,
						},
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

func genBinaryBuildSource() buildv1.BuildSource {
	return buildv1.BuildSource{
		Type:   "Source",
		Binary: &buildv1.BinaryBuildSource{},
	}
}

func genGitBuildSource(mg *metagraf.MetaGraf) buildv1.BuildSource {
	var branch string
	if len(params.SourceRef) > 0 {
		branch = params.SourceRef
	} else {
		branch = mg.Spec.Branch
	}

	return buildv1.BuildSource{
		Type: "Git",
		Git: &buildv1.GitBuildSource{
			URI: mg.Spec.Repository,
			Ref: branch,
		},
		SourceSecret: &corev1.LocalObjectReference{
			Name: mg.Spec.RepSecRef,
		},
	}
}

func StoreBuildConfig(obj buildv1.BuildConfig) {
	client := k8sclient.GetBuildClient().BuildConfigs(NameSpace)
	bc, _ := client.Get(obj.Name, metav1.GetOptions{})

	if len(bc.ResourceVersion) > 0 {
		obj.ResourceVersion = bc.ResourceVersion
		_, err := client.Update(&obj)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Updated BuildConfig: ", obj.Name, " in Namespace: ", NameSpace)
	} else {
		_, err := client.Create(&obj)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created BuildConfig: ", obj.Name, " in Namespace: ", NameSpace)
	}
}

func DeleteBuildConfig(name string) {
	client := k8sclient.GetBuildClient().BuildConfigs(NameSpace)

	_, err := client.Get(name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("The BuildConfig: ", name, "does not exist in namespace: ", NameSpace,", skipping...")
		return
	}

	err = client.Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		fmt.Println( "Unable to delete BuildConfig: ", name, " in namespace: ", NameSpace)
		log.Error(err)
		return
	}
	fmt.Println("Deleted BuildConfig: ", name, ", in namespace: ", NameSpace)
}