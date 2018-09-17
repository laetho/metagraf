/*
Copyright 2018 The MetaGraph Authors

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

package ocpclient

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	imagev1client "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
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
func getRestConfig(kc string) *rest.Config{
	var config *rest.Config

	if kc == "" {
		var err error
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error %v", err)
		}
	} else {
		var err error
		config, err = clientcmd.BuildConfigFromFlags("", kc)
		if err != nil {
			fmt.Fprintf(os.Stderr,"error %v", err)
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

