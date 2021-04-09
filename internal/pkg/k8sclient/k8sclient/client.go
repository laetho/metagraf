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

package k8sclient

import (
	"flag"
	argocdv1alpha1client "github.com/argoproj/argo-cd/pkg/client/clientset/versioned/typed/application/v1alpha1"
	monitoringv1client "github.com/coreos/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	appsv1client "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	buildv1client "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	imagev1client "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	routev1client "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	log "k8s.io/klog"
	"os"
	"path/filepath"
	"runtime"
)

var RestConfig *rest.Config

// Returns the current users home directory for both windows, mac and linux.
func getHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")

}

func getKubeConfig() string {
	var kubeconfig string
	if home := getHomeDir(); home != "" {
		kubeconfig = *flag.String(
			"kubeconfig", filepath.Join(home, ".kube", "config"),
			"(optional) path to a kubeconfig file")
	} else {
		kubeconfig = *flag.String("kubeconfig", "", "(optional) path to a kubeconfig file")
	}
	flag.Parse()
	return kubeconfig
}

// Get rest.Config from outside or inside cluster
func getRestConfig(kc string) *rest.Config {
	log.Infof("kubeconfig: %v", kc)
	var config *rest.Config

	if _, err := os.Stat(kc); os.IsNotExist(err) {
		var err error
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Infof("Did not find InClusterConfig: %v", err)
		}
	} else {
		var err error
		config, err = clientcmd.BuildConfigFromFlags("", kc)
		if err != nil {
			log.Errorf("Unable to build config from flags: %v", err)
			os.Exit(1)
		}
	}
	return config
}

// todo handle error
func GetCoreClient() *corev1client.CoreV1Client {
	if RestConfig == nil {
		RestConfig = getRestConfig(getKubeConfig())
	}

	client, err := corev1client.NewForConfig(RestConfig)
	if err != nil {
		panic(err)
	}

	return client
}

// Returns a ImageV1Client
func GetImageClient() *imagev1client.ImageV1Client {

	if RestConfig == nil {
		RestConfig = getRestConfig(getKubeConfig())
	}

	client, err := imagev1client.NewForConfig(RestConfig)
	if err != nil {
		panic(err)
	}

	return client
}

// Returns a Apps client
func GetAppsClient() *appsv1client.AppsV1Client {
	if RestConfig == nil {
		RestConfig = getRestConfig(getKubeConfig())
	}

	client, err := appsv1client.NewForConfig(RestConfig)
	if err != nil {
		panic(err)
	}

	return client
}

// Returns a K8S Apps client
func GetKubernetesClient() *kubernetes.Clientset {
	if RestConfig == nil {
		RestConfig = getRestConfig(getKubeConfig())
	}

	client, err := kubernetes.NewForConfig(RestConfig)
	if err != nil {
		panic(err)
	}

	return client
}

// Returns a build client
func GetBuildClient() *buildv1client.BuildV1Client {
	if RestConfig == nil {
		RestConfig = getRestConfig(getKubeConfig())
	}

	client, err := buildv1client.NewForConfig(RestConfig)
	if err != nil {
		panic(err)
	}

	return client
}

// Returns a build client
func GetRouteClient() *routev1client.RouteV1Client {
	if RestConfig == nil {
		RestConfig = getRestConfig(getKubeConfig())
	}

	client, err := routev1client.NewForConfig(RestConfig)
	if err != nil {
		panic(err)
	}

	return client
}

// Returns a build client
func GetMonitoringV1Client() *monitoringv1client.MonitoringV1Client {
	if RestConfig == nil {
		RestConfig = getRestConfig(getKubeConfig())
	}

	client, err := monitoringv1client.NewForConfig(RestConfig)
	if err != nil {
		panic(err)
	}

	return client
}

func GetArgoCDClient() *argocdv1alpha1client.ArgoprojV1alpha1Client {
	if RestConfig == nil {
		RestConfig = getRestConfig(getKubeConfig())
	}

	client, err := argocdv1alpha1client.NewForConfig(RestConfig)
	if err != nil {
		panic(err)
	}
	return client
}
