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

package metagraf

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	log "k8s.io/klog"
	"metagraf/mg/params"
	"strconv"
	"strings"
)

// Returns a list of GroupKind's described by the parsed metaGraf
// specification
func (mg MetaGraf) GroupKinds() []metav1.GroupKind {
	var gks []metav1.GroupKind

	sgk := metav1.GroupKind{Group: "core", Kind: "Service"}
	dgk := metav1.GroupKind{Group: "apps", Kind: "Deployment"}
	rgk := metav1.GroupKind{Group: "core", Kind: "Route"}

	gks = append(gks, sgk, dgk, rgk)

	return gks
}

func (mg MetaGraf) GetResourceByName(name string) (Resource, error) {
	for _, r := range mg.Spec.Resources {
		if r.Name == name {
			return r, nil
		}
	}
	return Resource{}, errors.New("Resource{} not found, name: " + name)
}

//
func (mg MetaGraf) GetSecretByName(name string) (Secret, error) {
	for _, s := range mg.Spec.Secret {
		if s.Name == name {
			return s, nil
		}
	}
	return Secret{}, errors.New("Secret{} not found, name: " + name)
}

//
func (mg MetaGraf) GetConfigByName(name string) (Config, error) {
	for _, c := range mg.Spec.Config {
		if c.Name == name {
			return c, nil
		}
	}
	return Config{}, errors.New("Config{} not found, name: " + name)
}

func (mg MetaGraf) Labels(name string) map[string]string {
	l := make(map[string]string)
	l["app"] = name

	for k, v := range mg.Metadata.Annotations {
		if !validLabelValue(sanitizeLabelValue(v)) {
			continue
		}
		if strings.Contains(k, params.NameSpacingFilter) {
			l[sanitizeKey(k)] = sanitizeLabelValue(v)
		}
	}
	for k, v := range mg.Metadata.Labels {
		l[sanitizeKey(k)] = v
	}
	return l
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
		}
	}
	return ports
}

func sanitizeLabelValue(val string) string {
	ret := strings.Replace(val, " ", "_", -1)
	ret = strings.Replace(ret, ",", "-", -1)
	return ret
}

func sanitizeKey(key string) string {
	if params.NameSpacingStripHost {
		parts := strings.Split(key, "/")
		if len(parts) > 1 {
			return strings.Join(parts[1:], "")
		}
	}
	return key
}

func validLabelValue(val string) bool {
	if len(val) > 64 {
		return false
	}
	if strings.Contains(val, "/") {
		return false
	}
	return true
}
