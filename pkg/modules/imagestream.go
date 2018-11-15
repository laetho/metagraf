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
	"github.com/blang/semver"
	"github.com/golang/glog"
	"metagraf/mg/ocpclient"
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

	if !Dryrun { StoreImageStream(is) }
	if Output { MarshalObject(is) }

}

func StoreImageStream(obj imagev1.ImageStream) {

	glog.Infof("ResourceVersion: %v Length: %v", obj.ResourceVersion, len(obj.ResourceVersion))
	glog.Infof("Namespace: %v", NameSpace)

	client := ocpclient.GetImageClient().ImageStreams(NameSpace)

	if len(obj.ResourceVersion) > 0 {
		// update
		result, err := client.Update(&obj)
		if err != nil {
			glog.Info(err)
		}
		glog.Infof("Updated ImageStream: %v(%v)", result.Name, obj.Name)
	} else {
		result, err := client.Create(&obj)
		if err != nil {
			glog.Info(err)
		}
		glog.Infof("Created ImageStream: %v(%v)", result.Name, obj.Name)
	}
}