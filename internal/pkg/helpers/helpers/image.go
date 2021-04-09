/*
Copyright 2018-2020 The metaGrafAuthors

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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/laetho/metagraf/internal/pkg/imageurl/imageurl"
	k8sclient2 "github.com/laetho/metagraf/internal/pkg/k8sclient"
	"github.com/laetho/metagraf/pkg/metagraf"
	dockerv10 "github.com/openshift/api/image/docker10"
	imagev1 "github.com/openshift/api/image/v1"
	imagev1client "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	log "k8s.io/klog"
	"os"
	"os/exec"
)

func GetImageStreamTags(c *imagev1client.ImageV1Client, ns string, n string) *imagev1.ImageStreamTag {
	ist, err := c.ImageStreamTags(ns).Get(context.TODO(), n, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if len(ist.Image.DockerImageMetadata.Raw) != 0 {
		di := &dockerv10.DockerImage{}
		err = json.Unmarshal(ist.Image.DockerImageMetadata.Raw, di)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		ist.Image.DockerImageMetadata.Object = di
	}
	return ist
}

// Unmarshals json data into a DockerImage type
func GetDockerImageFromIST(i *imagev1.ImageStreamTag) *dockerv10.DockerImage {
	di := dockerv10.DockerImage{}
	if len(i.Image.DockerImageMetadata.Raw) != 0 {
		err := json.Unmarshal(i.Image.DockerImageMetadata.Raw, &di)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		i.Image.DockerImageMetadata.Object = &di
	}
	return &di
}

// Inspect a docker image with skope and unmarshal that json output to
// a dockerv10.DockerImage{} struct.
func SkopeoImageInfo(image string) *dockerv10.DockerImage {
	image = "docker://" + image
	cmd := exec.Command("skopeo", "inspect", "--config", image)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	list := dockerv10.DockerImage{}
	err = json.Unmarshal(stdout, &list)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return &list
}

// Unmarshalljson data into a DockerImage type
func GetDockerImageFromImage(i *imagev1.Image) *dockerv10.DockerImage {
	di := dockerv10.DockerImage{}
	if len(i.DockerImageMetadata.Raw) != 0 {
		err := json.Unmarshal(i.DockerImageMetadata.Raw, &di)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		i.DockerImageMetadata.Object = &di
	}
	return &di
}

func GetImage(c *imagev1client.ImageV1Client, i string) *imagev1.Image {
	img, err := c.Images().Get(context.TODO(), i, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return img
}

// Returns docker information based on images referenced in spec.
func ImageInfo(mg *metagraf.MetaGraf) (*dockerv10.DockerImage, error) {
	var DockerImage string
	if len(mg.Spec.BaseRunImage) > 0 {
		DockerImage = mg.Spec.BaseRunImage
	} else if len(mg.Spec.BuildImage) > 0 {
		DockerImage = mg.Spec.BuildImage
	} else if len(mg.Spec.Image) > 0 {
		DockerImage = mg.Spec.Image
		// @todo, we need a way to natively inspect a upstream image when we're not running in openshift, or make the new way default.
		return nil, errors.New("blah")
	}

	var imgurl imageurl.ImageURL
	imgurl.Parse(DockerImage)

	client := k8sclient2.GetImageClient()

	ist := GetImageStreamTags(
		client,
		imgurl.Namespace,
		imgurl.Image+":"+imgurl.Tag)

	return GetDockerImageFromIST(ist), nil
}
