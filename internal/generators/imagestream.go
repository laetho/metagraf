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

	"metagraf/internal/metagraf"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	imagev1 "github.com/openshift/api/image/v1"
)

func GenImageStream( mg *metagraf.MetaGraf, namespace string) {
	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		fmt.Println(err)
	}

	isname := strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10))

	objref := corev1.ObjectReference{}
	objref.Kind = ""

	is := imagev1.ImageStream{
		TypeMeta: metav1.TypeMeta{
			Kind: "ImageStream",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: isname,
		},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				{
					From: &corev1.ObjectReference{
						Kind: "DockerImage",
						Name: "docker-registry.default.svc:5000/"+namespace+"/"+isname+":latest",
					},
					Name: "latest",
				},
			},
		},
	}

	ba, err := json.Marshal(is)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(ba))

}