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
	"github.com/blang/semver"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	log "k8s.io/klog"
	"metagraf/pkg/metagraf"
	"os"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

// This is a complete hack. todo: fix this shit, restructure packages
var (
	All		  bool
	NameSpace string // Used to pass namespace from cmd to module to avoid import cycle.
	Output    bool   // Flag passing hack
	Version   string // Flag passing hack
	Verbose   bool   // Flag passing hack
	Dryrun    bool   // Flag passing hack
	Branch	  string // Flag passing hack
	BaseEnvs  bool 		//Flag passing hack
	CVfile	  string	//Flag passing hack
	Defaults  bool 		//Flag passing hack
	Format	  string	// Flag passing hack
	Template  string	// Flag passing hack
	Suffix	  string	// Flag passing hack
	Enforce	  bool
	ImageNS	  string
	Registry  string
	Tag		  string
	OName	  string
	Context   string	// Application context root from FlagPassingHack.
	CreateGlobals bool
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

func GetEnvVars( vars []metagraf.EnvironmentVar, mgp metagraf.MGProperties) []corev1.EnvVar {
	var envs []corev1.EnvVar

	for _, e := range vars {
		env := e.ToEnvVar()
		prop,err := mgp.GetByKey(e.Name)
		if err != nil {
			log.V(2).Infof("Did not find input value for %s", e.Name)
		} else {
			env.Value = prop.Value
		}
		envs = append(envs, env)
	}
	return envs
}

func GetBuildEnvVars(mg *metagraf.MetaGraf, mgp metagraf.MGProperties) []corev1.EnvVar {
	var envs []corev1.EnvVar

	for _, e := range mg.Spec.Environment.Build {
		envs = append(envs, e.ToEnvVar())
	}

	return envs
}



func parseEnvVars(mg *metagraf.MetaGraf) []corev1.EnvVar {
	var EnvVars []corev1.EnvVar

	// Adding name and version of component as en environment variable
	EnvVars = append(EnvVars, corev1.EnvVar{
		Name:  "MG_APP_NAME",
		Value: MGAppName(mg),
	})

	var oversion string
	if len(Version) > 0 {
		oversion = Version
	} else {
		oversion = mg.Spec.Version
	}

	EnvVars = append(EnvVars, corev1.EnvVar{
		Name:  "MG_APP_VERSION",
		Value: oversion,
	})

	// Local variables from metagraf as deployment envvars
	for _, e := range mg.Spec.Environment.Local {

		if len(e.SecretFrom) > 0 && len(e.Key) == 0 {
			continue
		}
		if len(e.EnvFrom) > 0 && len(e.Key) == 0 {
			continue
		}

		// Inject JVM_SYS_PROP as an EnvVar. Content comes from a Config
		// section of type JVM_SYS_PROPS. But requires a Environment variable
		// of type JVM_SYS_PROP as well.
		if strings.ToUpper(e.Type) == "JVM_SYS_PROP" {
			if HasJVM_SYS_PROP(mg) {
				EnvVars = append(EnvVars, GenEnvVar_JVM_SYS_PROP(Variables, e.Name))
			}
			continue
		}

		if len(e.Key) > 0 {
			EnvVars = append(EnvVars, genValueFrom(&e))
			continue
		}

		// Omit optional EnvVar's that has no value provided through explicit config.
		if ev, t := Variables["local|"+e.Name]; t {
			if len(ev.Value) == 0 && ev.Required == false {
				log.V(3).Infof("Omitting optional environment variable %s without explicit value", e.Name)
				continue
			}
		}

		// Default behaviour
		EnvVars = append(EnvVars, EnvToEnvVar(&e, false))
	}

	// External variables from metagraf as deployment envvars
	for _, e := range mg.Spec.Environment.External.Consumes {
		EnvVars = append(EnvVars, EnvToEnvVar(&e, true))
	}
	for _, e := range mg.Spec.Environment.External.Introduces {
		EnvVars = append(EnvVars, EnvToEnvVar(&e, true))
	}

	return EnvVars
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
			custver = "-"+Version
		} else {
			custver = "v" + strings.ToLower(strconv.FormatUint(sv.Major, 10))
		}
	} else {
		sv, err := semver.Parse(mg.Spec.Version)
		if err != nil {
			custver = "-"+mg.Spec.Version
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
			return mg.Metadata.Name+"-"+ Version
		} else {
			objname = mg.Metadata.Name+"v"+strconv.FormatUint(sv.Major, 10)
			return objname + "-" + Version
		}
	}

	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		objname = mg.Metadata.Name
	} else {
		objname = mg.Metadata.Name+"v"+strconv.FormatUint(sv.Major, 10)
	}
	return objname
}

// Returns a name for a secret for a resource based on convention as follows.
func ResourceSecretName(r *metagraf.Resource) string {
	if len(r.User) > 0 && len(r.Secret) == 0 {
		// When an implicit secret is created it's resource name will
		// prepended to the user. They resourcename + user will get treated as a global secret.
		return strings.ToLower(strings.Replace(r.Name,"_","-", -1))+"-"+strings.ToLower(r.User)
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

// Prepends _ to indicate environment variable name to indicate a variable that
// comes from some external configuration repository. How you use this is an
// implementation detail in the runtime container image.
func ExternalEnvToEnvVar(e *metagraf.EnvironmentVar ) corev1.EnvVar {
	v:= EnvToEnvVar(e, true)
	v.Name = "_"+v.Name
	return v
}

/*
	Applies conventions and override logic to an environment variable and returns
	a corev1.EnvVar{}.
*/
func EnvToEnvVar(e *metagraf.EnvironmentVar, ext bool) corev1.EnvVar {
	name := ""		// Var for holding potentially modified name.
	value := ""		// Var for holding potentially modified or overridden value.

	// Handle external vs local variable notation
	if ext {
		name = "_"+e.Name
	} else {
		name = e.Name
	}

	if e.Required {

		// Set default value first if the Environment variable is not external
		// If Defaults flag is set, always populate default values, might get overridden by eksplicitly set values.
		if len(e.Default)> 0 && !ext {
			value = e.Default
		} else if Defaults {
			value = e.Default
		}

		// Handle possible override value for non required fields
		if p, t := Variables["local|"+e.Name]; t {
			if len(p.Value) > 0 {
				value = p.Value
				log.V(1).Info("Found override value for: ", p.Key, " Override value: ", p.Value)
			}
		}

		return corev1.EnvVar{
			Name:  name,
			Value: value,
		}
	} else {
		// Optional EnvironmentVariables should be populated but empty. Unless we choose to populate defaults.
		if Defaults {
			value = e.Default
		}
		// Handle override values for optional fields
		if p, t := Variables["local|"+e.Name]; t {
			value = p.Value
		}

		return corev1.EnvVar{
			Name:  name,
			Value: value,
		}
	}
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
			err := s.Encode(obj,os.Stdout)
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

func GetGlobalConfigMapVolumes(mg *metagraf.MetaGraf, Volumes *[]corev1.Volume, VolumeMounts *[]corev1.VolumeMount ) {
	// Only handle global ConfigMaps
	for _, v := range mg.Spec.Config {
		if !v.Global {
			continue
		}

		paths := []string{"/etc/pki/ca-trust/extracted/pem","/etc/pki/ca-trust/extracted/java"}
		items := []string{"tls-ca-bundle.pem","cacerts"}

		if strings.ToUpper(v.Type) == "TRUSTED-CA" {

			itms := []corev1.KeyToPath{}

			for _,key := range items {
				itm := corev1.KeyToPath{
					Key:  "ca-bundle.crt",
					Path: key,
				}
				itms = append(itms, itm)
			}

			vol := corev1.Volume{
				Name:         "trusted-ca",
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

