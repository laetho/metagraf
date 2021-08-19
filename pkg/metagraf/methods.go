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
	params "github.com/laetho/metagraf/internal/pkg/params"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (mg MetaGraf) GetEnvVarByType(envtype string) []EnvironmentVar {
	var envs []EnvironmentVar
	for _, env := range mg.Spec.Environment.Local {
		if env.Type == envtype {
			envs = append(envs, env)
		}
	}
	return envs
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
		l[sanitizeKey(k)] = sanitizeLabelValue(v)
	}
	return l
}



func sanitizeLabelValue(val string) string {
	// If a value includes a (, split the string and only return the part
	// leading up to the first (.
	ret := strings.Split(val, "(")[0]
	ret = strings.Replace(ret, " ", "_", -1)
	ret = strings.Replace(ret, ",", "-", -1)
	ret = strings.TrimSuffix(ret, "_")
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
