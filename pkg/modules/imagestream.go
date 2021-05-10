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
	"context"
	"fmt"
	"os"

	k8sclient "github.com/laetho/metagraf/internal/pkg/k8sclient"
	"github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	log "k8s.io/klog"

	imagev1 "github.com/openshift/api/image/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenImageStream(mg *metagraf.MetaGraf, namespace string) {

	objname := Name(mg)
	log.V(2).Infof("Generated ImageStream name: %v", objname)

	// Resource labels
	l := Labels(objname, labelsFromParams(params.Labels))

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

	if !Dryrun {
		StoreImageStream(is)
	}
	if Output {
		MarshalObject(is.DeepCopyObject())
	}

}

func StoreImageStream(obj imagev1.ImageStream) {

	log.Infof("ResourceVersion: %v Length: %v", obj.ResourceVersion, len(obj.ResourceVersion))
	log.Infof("Namespace: %v", NameSpace)

	client := k8sclient.GetImageClient().ImageStreams(NameSpace)

	im, err := client.Get(context.TODO(), obj.Name, metav1.GetOptions{})
	if len(im.ResourceVersion) > 0 {
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		log.Infof("ImageStream: %v exists, skipping...", obj.Name)
	} else {
		result, err := client.Create(context.TODO(), &obj, metav1.CreateOptions{})
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		log.Infof("Created ImageStream: %v(%v)", result.Name, obj.Name)
	}
}

func DeleteImageStream(name string) {
	client := k8sclient.GetImageClient().ImageStreams(NameSpace)

	_, err := client.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("ImageStream: ", name, "does not exist in namespace: ", NameSpace, ", skipping...")
		return
	}

	err = client.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("Unable to delete ImageStream: ", name, " in namespace: ", NameSpace)
		log.Error(err)
		return
	}
	fmt.Println("Deleted ImageStream: ", name, ", in namespace: ", NameSpace)
}
