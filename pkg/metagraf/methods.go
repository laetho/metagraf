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
	"metagraf/mg/params"
	"strconv"
	"strings"
)

// Returns a metagraf adressable key for a property.
func (mgp *MGProperty) MGKey() string {
	return mgp.Source+"|"+mgp.Key
}

// Returns a list of GroupKind's described by the parsed metaGraf
// specification
func (mg *MetaGraf) GroupKinds() []metav1.GroupKind {
	var gks []metav1.GroupKind

	sgk := metav1.GroupKind{Group: "core", Kind: "Service"}
	dgk := metav1.GroupKind{Group: "apps", Kind: "Deployment"}
	rgk := metav1.GroupKind{Group: "core", Kind: "Route"}

	gks = append(gks, sgk, dgk, rgk)

	return gks
}

// Returns a struct (MGProperties) of all MGProperty addressable
// in the metaGraf specification.
func (mg *MetaGraf) GetProperties() MGProperties {
	props := MGProperties{}

	// Config section, find parameters from
	for _,conf := range mg.Spec.Config {
		if len(conf.Options) == 0 {
			continue
		}

		switch conf.Type {
		case "parameters":
			for _, opts := range conf.Options {
				p := MGProperty{
					Source:   conf.Name,
					Key:      opts.Name,
					Value:    "",
					Required: opts.Required,
					Default: opts.Default,
				}
				props[p.MGKey()] = p
			}
		case "JVM_SYS_PROP":
			for _, opts := range conf.Options {
				p := MGProperty{
					Source:   "JVM_SYS_PROP",
					Key:      opts.Name,
					Value:    "",
					Required: opts.Required,
					Default: opts.Default,
				}
				props[p.MGKey()] = p
			}
		}
	}

	// Environment Section
	for _,env := range mg.Spec.Environment.Local {
		if len(env.SecretFrom) > 0 {continue}
		if len(env.EnvFrom) > 0 {continue}
		p := MGProperty{
			Source:   "local",
			Key:      env.Name,
			Value:    "",
			Required: env.Required,
			Default: env.Default,
		}

		// Environment variables of type JVM_SYS_PROP will
		// be implicitly populated by values from config
		// named JVM_SYS_PROP
		if env.Type == "JVM_SYS_PROP" {
			continue
		}
		props[p.MGKey()] = p
	}
	for _,env := range mg.Spec.Environment.External.Introduces {
		p := MGProperty{
			Source:   "external",
			Key:      env.Name,
			Value:    "",
			Required: env.Required,
			Default: env.Default,
		}
		props[p.MGKey()] = p
	}
	for _,env := range mg.Spec.Environment.External.Consumes {
		p := MGProperty{
			Source:   "external",
			Key:      env.Name,
			Value:    "",
			Required: env.Required,
			Default: env.Default,
		}
		props[p.MGKey()] = p
	}


	return props
}

// Returns the MGProperty.Required = true
func (mgp MGProperties) GetRequired() MGProperties {
	props := MGProperties{}

	for _, prop := range mgp {
		if prop.Required && prop.Source != "external" {
			props[prop.MGKey()] = prop
		}
	}
	return props
}
// Returns a slice of Keys
func (mgp MGProperties) Keys() []string {
	var keys []string
	for _, prop := range mgp {
		keys = append(keys, prop.Key)
	}
	return keys
}

func (mgp MGProperties) GetByKey(key string) (MGProperty, error){
	for _, p := range mgp {
		if p.Key == key {
			return p, nil
		}
	}
	return MGProperty{}, errors.Errorf("Key not found!")
}

// Return a slice of property keys. If required == true only return required keys.
func (mgp MGProperties) SourceKeys(required bool) []string {
	var keys []string
	for _, prop := range mgp {
		if required {
			if prop.Required == required {
				keys = append(keys, prop.MGKey())
			}
		} else {
			keys = append(keys, prop.MGKey())
		}
	}
	return keys
}

// Returns a map of key,values
func (mgp MGProperties) KeyMap() map[string]string {
	keys := make(map[string]string)
	for _, prop := range mgp {
		keys[prop.Key] = prop.Value
	}
	return keys
}

// Returns a map of MGProperty.Source+Key, MGProperty.Value
// It takes a boolean as argument. Return only required or all keys?
func (mgp MGProperties) SourceKeyMap(required bool) map[string]string {
	keys := make(map[string]string)
	for _, prop := range mgp {
		if prop.Required == required {
			keys[prop.MGKey()] = prop.Value
			continue
		}
		keys[prop.MGKey()] = prop.Value
	}
	return keys
}

func (mg *MetaGraf) GetResourceByName(name string) (Resource, error) {
	for _,r := range mg.Spec.Resources{
		if r.Name == name {
			return r, nil
		}
	}
	return Resource{}, errors.New("Resource{} not found, name: "+name)
}

//
func (mg *MetaGraf) GetSecretByName(name string) (Secret, error) {
	for _,s := range mg.Spec.Secret{
		if s.Name == name {
			return s, nil
		}
	}
	return Secret{}, errors.New("Secret{} not found, name: "+name)
}

//
func (mg *MetaGraf) GetConfigByName(name string) (Config, error) {
	for _,c := range mg.Spec.Config{
		if c.Name == name {
			return c, nil
		}
	}
	return Config{}, errors.New("Config{} not found, name: "+name)
}

func (mg *MetaGraf) Labels(name string) map[string]string {
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
func (mg *MetaGraf) AnnotationServicePorts() ([]corev1.ServicePort, error) {
	var ports []corev1.ServicePort
	for k,v := range mg.Metadata.Annotations {
		if strings.Contains(k, ".service.k8s.io/port") {
			protocol := strings.Split(k, ".")[0]
			switch protocol {
			case "http":
				intport, err := strconv.Atoi(v)
				if err != nil {
					return ports, errors.Errorf("Unable to convert port to numeric value for annotation: %v", k)
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
					return ports, errors.Errorf("Unable to convert port to numeric value for annotation: %v", k)
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
	return ports, nil
}

func sanitizeLabelValue(val string) string {
	ret := strings.Replace(val, " ", "_", -1)
	ret = strings.Replace(ret, ",", "-", -1 )
	return ret
}

func sanitizeKey(key string) string {
	if params.NameSpacingStripHost {
		parts := strings.Split(key,"/")
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
