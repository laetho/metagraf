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
	"metagraf/mg/ocpclient"
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

		glog.Error("Malformed BuildImage url provided in metaGraf file; %v", mg.Spec.BuildImage)
		os.Exit(1)
	}

	objname := Name(mg)

	if len(mg.Spec.BaseRunImage) > 0 && len(mg.Spec.Repository) > 0 {
		buildsource = genBinaryBuildSource()
	} else if len(mg.Spec.BuildImage) > 0 && len(mg.Spec.BaseRunImage) < 1 {
		buildsource = genGitBuildSource(mg)
	}

	client := ocpclient.GetImageClient()

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

	// Resource labels
	l := make(map[string]string)
	l["app"] = objname

	for _, e := range mg.Spec.Environment.Build {
		if e.Required == true {
			EnvVars = append(EnvVars, corev1.EnvVar{Name: e.Name, Value: e.Default})
		} else if e.Required == false {
			EnvVars = append(EnvVars, corev1.EnvVar{Name: e.Name, Value: "null"})
		}
	}

	bc := buildv1.BuildConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "BuildConfig",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: objname,
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
					To: &corev1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: objname + ":latest",
					},
				},
			},
		},
	}

	fmt.Println(Dryrun,Output)

	if !Dryrun { StoreBuildConfig(bc) }
	if Output { MarshalObject(bc) }
}

func genBinaryBuildSource() buildv1.BuildSource {
	return buildv1.BuildSource{
		Type:   "Source",
		Binary: &buildv1.BinaryBuildSource{},
	}
}

func genGitBuildSource(mg *metagraf.MetaGraf) buildv1.BuildSource {
	return buildv1.BuildSource{
		Type: "Git",
		Git: &buildv1.GitBuildSource{
			URI: mg.Spec.Repository,
			Ref: mg.Spec.Branch,
		},
		SourceSecret: &corev1.LocalObjectReference{
			Name: mg.Spec.RepSecRef,
		},
	}
}

func StoreBuildConfig(obj buildv1.BuildConfig) {

	glog.Infof("ResourceVersion: %v Length: %v", obj.ResourceVersion, len(obj.ResourceVersion))
	glog.Infof("Namespace: %v", NameSpace)

	client := ocpclient.GetBuildClient().BuildConfigs(NameSpace)

	if len(obj.ResourceVersion) > 0 {
		// update
		result, err := client.Update(&obj)
		if err != nil {
			glog.Info(err)
		}
		glog.Infof("Updated BuildConfig: %v(%v)", result.Name, obj.Name)
	} else {
		result, err := client.Create(&obj)
		if err != nil {
			glog.Info(err)
		}
		glog.Infof("Created BuildConfig: %v(%v)", result.Name, obj.Name)
	}
}