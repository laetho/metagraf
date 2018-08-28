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
	"github.com/blang/semver"
	"metagraf/pkg/metagraf"
	"strconv"
	"strings"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenService(mg *metagraf.MetaGraf) {
	var objname string
	var serviceports []corev1.ServicePort


	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		objname = strings.ToLower(mg.Metadata.Name)
	} else {
		objname = strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10))
	}

	selectors := make(map[string]string)
	selectors["deploymentconfig"] = objname

	labels := make(map[string]string)
	labels["app"] = objname

	obj := corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: objname,
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: serviceports,
			Selector: selectors,
		},
	}

	ba, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(ba))

}