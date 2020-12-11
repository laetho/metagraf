/*
Copyright 2018 The metaGraf Authors

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
	"github.com/openshift/api/image/docker10"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"strings"
)

// Returns a slice of k8s.io v1 ServicePort{} types from a Docker image config.
func ImageExposedPortsToServicePorts(config *docker10.DockerConfig ) []corev1.ServicePort {
	var ports []corev1.ServicePort
	for k := range config.ExposedPorts {
		ss := strings.Split(k, "/")
		port, _ := strconv.Atoi(ss[0])
		ContainerPort := corev1.ServicePort{
			Name:     strings.ToLower(ss[0]) + "-" + ss[1],
			Port:     int32(port),
			Protocol: corev1.Protocol(strings.ToUpper(ss[1])),
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: int32(port),
				StrVal: ss[1],
			},
		}
		ports = append(ports, ContainerPort)
	}
	return ports
}

/*
// todo: verify image string

func DockerInspectImage(image string, tag string, auth docker.AuthConfiguration) (*docker.Image, error) {

	endpoint := "unix://var/run/docker.sock"

	cli, err := docker.NewClient(endpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to authenticate to registry")
		return nil, err
	}

	pull := docker.PullImageOptions{
		Tag: tag,
		Repository: image,
	}

	err = cli.PullImage( pull, auth )
	if err != nil {
		return nil, err
	}

	i, err := cli.InspectImage(image)
	if err != nil {
		return nil, err
	}
	return i, nil
}
*/