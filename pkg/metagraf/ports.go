/*
Copyright 2021 The metaGraf Authors

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

package metagraf

import (
	"strconv"
	"strings"

	"github.com/laetho/metagraf/internal/pkg/helpers"
	"github.com/laetho/metagraf/internal/pkg/imageurl"
	"github.com/laetho/metagraf/internal/pkg/k8sclient"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	log "k8s.io/klog"
)

func (mg MetaGraf) GetServicePorts() []corev1.ServicePort {
	if mg.HasServicePorts() {
		return mg.ServicePorts()
	} else {
		return mg.DefaultServicePorts()
	}
}

func (mg MetaGraf) HasServicePorts() bool {
	if len(mg.Spec.Ports) > 0 {
		return true
	}
	return false
}

func (mg MetaGraf) ServicePorts() []corev1.ServicePort {

	serviceports := mg.ServicePortsBySpec()

	// Find all ports from the container image not defined in
	//mg.Spec.Ports section and add them to the list of ServicePorts
	for _, ip := range mg.ImagePorts() {
		portMatched := false
		for _, sp := range serviceports {
			if ip.Port == sp.TargetPort.IntVal {
				portMatched = true
			}
		}

		if !portMatched {
			serviceports = append(serviceports, ip)
		}
	}

	return serviceports
}


func (mg MetaGraf) ImagePorts() []corev1.ServicePort {

	var imgurl imageurl.ImageURL
	_ = imgurl.Parse(mg.GetDockerImageURL())

	ImageInfo, err := helpers.ImageInfo(mg.GetDockerImageURL())
	if err != nil {
		return []corev1.ServicePort{}
	}

	client := k8sclient.GetImageClient()
	ist := helpers.GetImageStreamTags(
		client,
		imgurl.Namespace,
		imgurl.Image+":"+imgurl.Tag)
	ImageInfo = helpers.GetDockerImageFromIST(ist)

	return helpers.ImageExposedPortsToServicePorts(ImageInfo.Config)
}

// Returns the mg tool's opinionated default if ports in a metagraf spec
// or the container image has no default exposed ports.
func (mg MetaGraf) DefaultServicePorts() []corev1.ServicePort {
	var serviceports []corev1.ServicePort

	serviceports = append(serviceports, corev1.ServicePort{
		Name:     "http",
		Port:     int32(80),
		Protocol: "TCP",
		TargetPort: intstr.IntOrString{
			Type:   0,
			IntVal: int32(8080),
			StrVal: "8080",
		},
	})
	return serviceports
}

// Checks the metagraf specification for k8s.io namespaced port information.
// Format:
// <protocol>.service.k8s.io/port : value
func (mg MetaGraf) ServicePortsByAnnotation() []corev1.ServicePort {
	var ports []corev1.ServicePort
	for k, v := range mg.Metadata.Annotations {
		if strings.Contains(k, ".service.k8s.io/port") {
			protocol := strings.Split(k, ".")[0]
			switch protocol {
			case "http":
				intport, err := strconv.Atoi(v)
				if err != nil {
					log.Warningf("Unable to convert port to numeric value for annotation: %v", k)
					continue
				}
				ports = append(ports, corev1.ServicePort{
					Name:     "http",
					Port:     int32(80),
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: int32(intport),
						StrVal: v,
					},
				})
			case "https":
				intport, err := strconv.Atoi(v)
				if err != nil {
					log.Warningf("Unable to convert port to numeric value for annotation: %v", k)
					continue
				}
				ports = append(ports, corev1.ServicePort{
					Name:     "https",
					Port:     int32(443),
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: int32(intport),
						StrVal: v,
					},
				})
			}
		}
	}
	return ports
}

func (mg MetaGraf) ServicePortsBySpec() []corev1.ServicePort {
	var ports []corev1.ServicePort
	for protocol, port := range mg.Spec.Ports {
		switch protocol {
		case "http":
			ports = append(ports, corev1.ServicePort{
				Name:     "http",
				Port:     int32(80),
				Protocol: "TCP",
				TargetPort: intstr.IntOrString{
					Type:   0,
					IntVal: port,
					StrVal: protocol,
				},
			})
		case "https":
			ports = append(ports, corev1.ServicePort{
				Name:     "https",
				Port:     int32(443),
				Protocol: "TCP",
				TargetPort: intstr.IntOrString{
					Type:   0,
					IntVal: port,
					StrVal: protocol,
				},
			})
		default:
			ports = append(ports, corev1.ServicePort{
				Name:     protocol,
				Protocol: "TCP",
				Port:     port,
				TargetPort: intstr.IntOrString{
					Type:   0,
					IntVal: port,
					StrVal: protocol,
				},
			})
		}
	}
	return ports
}