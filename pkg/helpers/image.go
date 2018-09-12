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

package helpers

import (
	"fmt"
	imagev1client "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	imagev1 "github.com/openshift/api/image/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetImageStreamTags(c *imagev1client.ImageV1Client, ns string, n string) *imagev1.ImageStreamTag {
	//fmt.Println("Image:",ns,n)
	ist, err := c.ImageStreamTags(ns).Get(n, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println(len(ist.Image.DockerImageConfig))
	fmt.Printf("Type: %t",ist.Image.DockerImageConfig)

	return ist
}

func GetImage(c *imagev1client.ImageV1Client, i string) *imagev1.Image {
	img, err := c.Images().Get(i, metav1.GetOptions{})
	if err != nil {
		fmt.Println(err.Error())
		panic(1)
	}
	return img
}
