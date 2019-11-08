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

package modules

import (
	"github.com/blang/semver"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"metagraf/pkg/metagraf"
	"os"
	"strconv"
	"strings"


	appsv1 "github.com/openshift/api/apps/v1"
	authorizationv1 "github.com/openshift/api/authorization/v1"
	buildv1 "github.com/openshift/api/build/v1"
	imagev1 "github.com/openshift/api/image/v1"
	networkv1 "github.com/openshift/api/network/v1"
	oauthv1 "github.com/openshift/api/oauth/v1"
	projectv1 "github.com/openshift/api/project/v1"
	quotav1 "github.com/openshift/api/quota/v1"
	routev1 "github.com/openshift/api/route/v1"
	securityv1 "github.com/openshift/api/security/v1"
	templatev1 "github.com/openshift/api/template/v1"
	userv1 "github.com/openshift/api/user/v1"

	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

// This is a complete hack. todo: fix this shit, restructure packages
var (
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
)

var Variables map[string]string

// Returns a name for a resource based on convention as follows.
func Name(mg *metagraf.MetaGraf) string {
	var objname string

	if len(OName) > 0 {
		glog.Infof("ObjectName overridden with: %v",OName)
		return OName
	}

	if len(Version) > 0 {
		sv, err := semver.Parse(Version)
		if err != nil {
			return strings.ToLower(mg.Metadata.Name)+"-"+Version
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
		if v, t := Variables[e.Name]; t {
			if len(v) > 0 {
				value = v
				glog.Info("Found override value for: ", name, " Override value: ", Variables[name])
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
		if v, t := Variables[e.Name]; t {
			value = v
		}

		return corev1.EnvVar{
			Name:  name,
			Value: value,
		}
	}
}

func ValueFromEnv(key string) bool {
	if _, t := Variables[key]; t {
		return true
	}
	return false
}

// Marshal kubernetes resource to json
func MarshalObject(obj runtime.Object) {
	_ = appsv1.AddToScheme(scheme.Scheme)
	_ = authorizationv1.AddToScheme(scheme.Scheme)
	_ = buildv1.AddToScheme(scheme.Scheme)
	_ = imagev1.AddToScheme(scheme.Scheme)
	_ = networkv1.AddToScheme(scheme.Scheme)
	_ = oauthv1.AddToScheme(scheme.Scheme)
	_ = projectv1.AddToScheme(scheme.Scheme)
	_ = quotav1.AddToScheme(scheme.Scheme)
	_ = routev1.AddToScheme(scheme.Scheme)
	_ = securityv1.AddToScheme(scheme.Scheme)
	_ = templatev1.AddToScheme(scheme.Scheme)
	_ = userv1.AddToScheme(scheme.Scheme)

	switch Format {
		case "json":
			s := json.NewSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)
			err := s.Encode(obj,os.Stdout)
			if err != nil {
				glog.Error(err)
				os.Exit(1)
			}
		case "yaml":
			s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
			err := s.Encode(obj, os.Stdout)
			if err != nil {
				glog.Error(err)
				os.Exit(1)
			}
	}
}
