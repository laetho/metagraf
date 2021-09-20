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
	"bytes"
	gojson "encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/blang/semver"
	"github.com/laetho/metagraf/pkg/metagraf"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	log "k8s.io/klog"
	"k8s.io/klog/v2"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

// This is a complete hack. todo: fix this shit, restructure packages
var (
	NameSpace     string // Used to pass namespace from cmd to module to avoid import cycle.
	Output        bool   // Flag passing hack
	Version       string // Flag passing hack
	Dryrun        bool   // Flag passing hack
	BaseEnvs      bool   //Flag passing hack
	Defaults      bool   //Flag passing hack
	Format        string // Flag passing hack
	Template      string // Flag passing hack
	Suffix        string // Flag passing hack
	ImageNS       string
	Registry      string
	Tag           string
	OName         string
	Context       string // Application context root from FlagPassingHack.
	CreateGlobals bool
	// Sets the default pull policy for all metagraf modules
	PullPolicy corev1.PullPolicy = corev1.PullIfNotPresent
)

var Variables metagraf.MGProperties

// Returns a corev1.EnvVar{} with a valueFrom construct if the metagraf.EnvironmentVar
// has a SecretFrom or EnvFrom and the Key reference is set.
func genValueFrom(e *metagraf.EnvironmentVar) corev1.EnvVar {
	var EnvVar corev1.EnvVar

	if len(e.SecretFrom) > 0 {
		EnvVar = corev1.EnvVar{
			Name: e.Name,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: e.SecretFrom,
					},
					Key: e.Key,
				},
			},
		}
	}

	if len(e.EnvFrom) > 0 {
		EnvVar = corev1.EnvVar{
			Name: e.Name,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: e.EnvFrom,
					},
					Key: e.Key,
				},
			},
		}
	}

	return EnvVar
}

// Returns a slice of metagraf convention informative k8s EnvVar{}'s
func GetMGEnvVars(mg *metagraf.MetaGraf) []corev1.EnvVar {
	var envs []corev1.EnvVar
	// Adding name and version of component as en environment variable
	envs = append(envs, corev1.EnvVar{
		Name:  "MG_APP_NAME",
		Value: MGAppName(mg),
	})

	var oversion string
	if len(Version) > 0 {
		oversion = Version
	} else {
		oversion = mg.Spec.Version
	}

	envs = append(envs, corev1.EnvVar{
		Name:  "MG_API_VERSION",
		Value: oversion,
	})
	return envs
}

func GetEnvVars(mg *metagraf.MetaGraf, inputprops metagraf.MGProperties) []corev1.EnvVar {
	specenvs := mg.Spec.Environment.Local
	var outputevars []corev1.EnvVar

	outputevars = append(outputevars, GetMGEnvVars(mg)...)

	// If the spec has any JVM_SYS_PROP variable generate and inject values from
	// input properties.
	if HasJVM_SYS_PROP(mg) {
		envvars := mg.GetEnvVarByType("JVM_SYS_PROP")
		if len(envvars) > 1 {
			klog.Warning("Only one variable of JVM_SYS_PROP is supported. If you have multiple all JVM_SYS_PROP definitions will get duplicated.")
		}
		for _, env := range envvars {
			outputevars = append(outputevars, GenEnvVar_JVM_SYS_PROP(inputprops, env.Name))
		}
	}

	// Create corev1.EnvVar for all .spec.environment.local of with type secretfrom or envfrom.
	for _, e := range specenvs {
		if e.GetType() == "default" {
			continue
		}
		env := e.ToEnvVar()
		outputevars = append(outputevars, env)
	}

	specenvsmap := make(map[string]metagraf.EnvironmentVar)
	// Generate EnvironmentVar from MGProperty that is not defined in metagraf
	// provided by vars argument here.
	for _, i := range specenvs {
		specenvsmap["local|"+i.Name] = i
	}

	// Find and set input values for all variables in .spec.environment.local
	for _, p := range inputprops {
		// Skip input props if they are not from source "local".
		if p.Source != "local" {
			continue
		}

		// If input property is found in spec environment vars. Transform to a corev1.EnvVar
		// and set it's value from the inputproperty. If a local input property is not from spec,
		// transform MGProperty into a corev1.EnvVar containing the input value.
		if _, ok := specenvsmap[p.MGKey()]; ok {
			localenv := specenvsmap[p.MGKey()]
			outvar := localenv.ToEnvVar()
			outvar.Value = p.Value
			outputevars = append(outputevars, outvar)
		} else {
			environmentVar := p.ToEnvironmentVar()
			envvar := environmentVar.ToEnvVar()
			envvar.Value = p.Value
			outputevars = append(outputevars, envvar)
		}
	}
	return outputevars
}

func GetBuildEnvVars(mg *metagraf.MetaGraf, mgp metagraf.MGProperties) []corev1.EnvVar {
	var envs []corev1.EnvVar

	for _, e := range mg.Spec.Environment.Build {
		envs = append(envs, e.ToEnvVar())
	}

	return envs
}

// Parses the metagraf specification for instances of reading environment variables
// from either a ConfigMap or a Secret. If the metagraf.Environment.Key is provided
func parseEnvFrom(mg *metagraf.MetaGraf) []corev1.EnvFromSource {
	var EnvFrom []corev1.EnvFromSource

	for _, e := range mg.Spec.Environment.Local {
		if len(e.EnvFrom) == 0 || len(e.Key) > 0 {
			continue
		}
		EnvFrom = append(EnvFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: e.EnvFrom,
				},
			},
		})
	}

	for _, e := range mg.Spec.Environment.Local {
		if len(e.SecretFrom) == 0 || len(e.Key) > 0 {
			continue
		}
		cmref := corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: e.SecretFrom,
				},
			},
		}
		EnvFrom = append(EnvFrom, cmref)
	}

	return EnvFrom
}

// Returns a name for a resource based on convention as follows.
func Name(mg *metagraf.MetaGraf) string {
	var custver string
	var custname string

	if len(OName) > 0 {
		custname = OName
	} else {
		custname = strings.ToLower(mg.Metadata.Name)
	}

	if len(Version) > 0 {
		sv, err := semver.Parse(Version)
		if err != nil {
			custver = "-" + Version
		} else {
			custver = "v" + strings.ToLower(strconv.FormatUint(sv.Major, 10))
		}
	} else {
		sv, err := semver.Parse(mg.Spec.Version)
		if err != nil {
			custver = "-" + mg.Spec.Version
		} else {
			custver = "v" + strings.ToLower(strconv.FormatUint(sv.Major, 10))
		}
	}

	return strings.ToLower(custname + custver)
}

// Return a specification name for a resource base on convention. Does not adhere to override flags.
func SpecName(mg *metagraf.MetaGraf) string {
	var objname string

	if len(Version) > 0 {
		sv, err := semver.Parse(Version)
		if err != nil {
			return strings.ToLower(mg.Metadata.Name+"-") + Version
		} else {
			objname = strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10))
			return objname + "-" + Version
		}
	}

	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		objname = strings.ToLower(mg.Metadata.Name)
	} else {
		objname = strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10))
	}
	return objname
}

func MGAppName(mg *metagraf.MetaGraf) string {
	var objname string

	if len(Version) > 0 {
		sv, err := semver.Parse(Version)
		if err != nil {
			return mg.Metadata.Name + "-" + Version
		} else {
			objname = mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10)
			return objname + "-" + Version
		}
	}

	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		objname = mg.Metadata.Name
	} else {
		objname = mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10)
	}
	return objname
}

// Returns a name for a secret for a resource based on convention as follows.
func ResourceSecretName(r *metagraf.Resource) string {
	if len(r.User) > 0 && len(r.Secret) == 0 {
		// When an implicit secret is created it's resource name will
		// prepended to the user. They resourcename + user will get treated as a global secret.
		return strings.ToLower(strings.Replace(r.Name, "_", "-", -1)) + "-" + strings.ToLower(r.User)
	} else if len(r.User) == 0 && len(r.Secret) > 0 {
		// Explicit secret name generation
		return strings.ToLower(r.Secret)
	} else {
		return strings.ToLower(r.Name)
	}
}

func ConfigSecretName(c *metagraf.Config) string {
	return strings.ToLower(c.Name)
}

func ValueFromEnv(key string) bool {
	if _, t := Variables["local|"+key]; t {
		return true
	}
	return false
}

// Marshal kubernetes resource to json
func MarshalObject(obj runtime.Object) {
	switch Format {
	case "json":
		opt := json.SerializerOptions{
			Yaml:   false,
			Pretty: true,
			Strict: true,
		}
		s := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, opt)
		err := s.Encode(obj, os.Stdout)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	case "yaml":
		opt := json.SerializerOptions{
			Yaml:   true,
			Pretty: true,
			Strict: true,
		}
		s := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, opt)
		err := s.Encode(obj, os.Stdout)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
}

func MarshalObjectWithoutStatus(obj runtime.Object) {
	opt := json.SerializerOptions{
		Yaml:   false,
		Pretty: true,
		Strict: true,
	}
	s := json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, opt)

	var buff bytes.Buffer
	err := s.Encode(obj.DeepCopyObject(), &buff)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	jsonMap := make(map[string]interface{})
	err = gojson.Unmarshal(buff.Bytes(), &jsonMap)
	if err != nil {
		panic(err)
	}

	delete(jsonMap, "status")

	if Format == "json" {
		oj, err := gojson.MarshalIndent(jsonMap, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(oj))
		return
	} else {
		oy, err := yaml.Marshal(jsonMap)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(oy))
	}
}

func GetGlobalConfigMapVolumes(mg *metagraf.MetaGraf, Volumes *[]corev1.Volume, VolumeMounts *[]corev1.VolumeMount) {
	// Only handle global ConfigMaps
	for _, v := range mg.Spec.Config {
		if !v.Global {
			continue
		}

		paths := []string{"/etc/pki/ca-trust/extracted/pem", "/etc/pki/ca-trust/extracted/java"}
		items := []string{"tls-ca-bundle.pem", "cacerts"}

		if strings.ToUpper(v.Type) == "TRUSTED-CA" {

			itms := []corev1.KeyToPath{}

			for _, key := range items {
				itm := corev1.KeyToPath{
					Key:  "ca-bundle.crt",
					Path: key,
				}
				itms = append(itms, itm)
			}

			vol := corev1.Volume{
				Name: "trusted-ca",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "trusted-ca",
						},
						Items: itms,
					},
				},
			}
			*Volumes = append(*Volumes, vol)

			for _, path := range paths {
				volm := corev1.VolumeMount{
					Name:      "trusted-ca",
					ReadOnly:  true,
					MountPath: path,
				}
				*VolumeMounts = append(*VolumeMounts, volm)
			}
		}
	}
}

func labelsFromParams(labels []string) map[string]string {
	ret := make(map[string]string)
	for _, s := range labels {
		split := strings.Split(s, "=")
		if len(split) != 2 {
			continue
		}
		ret[split[0]] = split[1]
	}
	return ret
}

// Generate standardised labels map
func Labels(name string, input map[string]string) map[string]string {
	// Resource labels
	l := make(map[string]string)
	l["app"] = name
	for k, v := range input {
		l[k] = v
	}
	return l
}

// Builds and returns slice of Kubernetes EnvVars for common values
// extracted from DownwardAPI.
func DownwardAPIEnvVars() []corev1.EnvVar {
	vars := []corev1.EnvVar{}

	// Map for fieldRefs: Name, fieldRef
	fieldRefs := map[string]string{
		"POD_NAME":      "metadata.name",
		"POD_NAMESPACE": "metadata.namespace",
		"NODE_NAME":     "spec.nodeName",
	}

	for k, v := range fieldRefs {
		ev := corev1.EnvVar{
			Name: k,
			ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: v,
			}},
		}
		vars = append(vars, ev)
	}
	return vars
}
