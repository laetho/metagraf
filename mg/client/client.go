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

package client

import (
	"flag"
	"fmt"
    "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"runtime"
	buildv1 "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
)


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

func GetClient() buildv1.BuildV1Client {
	var config *rest.Config

	kubeconfig := getKubeConfig()

	// Determine if we are running inside the cluster or not
	if kubeconfig == "" {
		var err error
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error %v", err)
		}
	} else {
		var err error
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Fprintf(os.Stderr,"error %v", err)
			os.Exit(1)
		}
	}

	client, err := buildv1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return *client
}