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
	"metagraf/internal/metagraf"

	buildv1 "github.com/openshift/api/build/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenBuildConfig(mg *metagraf.MetaGraf) {
	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		fmt.Println(err)
	}

	bc := buildv1.BuildConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10)),
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
				},
			},
		},
	}

	/*
		bc := buildv1.BuildConfig{}
		bc.SetName( strings.ToLower(mg.Metadata.Name +"v"+strconv.FormatUint(sv.Major, 10)) )
		btp := buildv1.BuildTriggerPolicy{}
		btp.Type = "GitHub"
		bc.Spec.Triggers = append(bc.Spec.Triggers, btp)
	*/

	ba, err := json.Marshal(bc)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(ba))
}
