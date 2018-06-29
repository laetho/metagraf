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
	"metagraf/pkg/metagraf"

	corev1 "k8s.io/api/core/v1"
	buildv1 "github.com/openshift/api/build/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
 @todo: Figure out how to best inject artifact
 */
func GenBuildConfig(mg *metagraf.MetaGraf) {
	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		fmt.Println(err)
	}

	objname := strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10))

	// Resource labels
	l := make(map[string]string)
	l["app"] = objname

	bc := buildv1.BuildConfig{
		TypeMeta: metav1.TypeMeta{
			Kind: "BuildConfig",
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
						Secret: "moregibberish",
					},
				},
			},
			RunPolicy:  buildv1.BuildRunPolicySerial,
			CommonSpec: buildv1.CommonSpec{
				Source: buildv1.BuildSource{
					Type: "Source",
					Binary: &buildv1.BinaryBuildSource{},
				},
				Strategy: buildv1.BuildStrategy{
					Type: buildv1.SourceBuildStrategyType,
					SourceStrategy: &buildv1.SourceBuildStrategy{
						From: corev1.ObjectReference{
							Kind: "ImageStreamTag",
							Namespace: "openshift",
							Name: "nt-wlp-pipeline:latest", // @todo: parameterize buildimage
						},
					},
				},
				Output: buildv1.BuildOutput{
					To: &corev1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: objname+":latest",
					},
				},
			},
		},
	}

	ba, err := json.Marshal(bc)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(ba))
}
