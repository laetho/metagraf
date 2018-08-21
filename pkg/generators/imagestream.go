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
	"github.com/blang/semver"
	"strconv"
	"strings"

	"metagraf/pkg/metagraf"

	imagev1 "github.com/openshift/api/image/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenImageStream(mg *metagraf.MetaGraf, namespace string) {
	var objname string
	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		objname = strings.ToLower(mg.Metadata.Name)
	} else {
		objname = strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10))
	}

	// Resource labels
	l := make(map[string]string)
	l["app"] = objname

	objref := corev1.ObjectReference{}
	objref.Kind = ""

	is := imagev1.ImageStream{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ImageStream",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   objname,
			Labels: l,
		},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				{
					From: &corev1.ObjectReference{
						Kind: "DockerImage",
						Name: "docker-registry.default.svc:5000/" + namespace + "/" + objname + ":latest",
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
