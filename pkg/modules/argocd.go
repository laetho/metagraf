/*
Copyright 2020 The metaGraf Authors

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
	"context"
	gojson "encoding/json"
	"fmt"
	argoapp "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/golang/glog"
	"github.com/laetho/metagraf/internal/pkg/k8sclient"
	"github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	log "k8s.io/klog"
	"os"
)

func GetArgoCDApplicationSyncPolicy() *argoapp.SyncPolicy {
	sp := &argoapp.SyncPolicy{}

	if params.ArgoCDAutomatedSyncPolicy {
		sp.Automated = &argoapp.SyncPolicyAutomated{
			Prune:    params.ArgoCDAutomatedSyncPolicyPrune,
			SelfHeal: params.ArgoCDAutomatedSyncPolicySelfHeal,
		}
	}

	if params.ArgoCDSyncPolicyRetry {

	}
	return sp
}

func GetArgoCDApplicationNamespace() string {
	if len(params.ArgoCDApplicationNamespace) > 0 {
		return params.ArgoCDApplicationNamespace
	}
	return NameSpace
}

func GetArgoCDSourceDirectory() *argoapp.ApplicationSourceDirectory {
	asd := argoapp.ApplicationSourceDirectory{
		Recurse: params.ArgoCDApplicationSourceDirectoryRecurse,
		Jsonnet: argoapp.ApplicationSourceJsonnet{},
	}
	return &asd
}

func GenArgoApplication(mg *metagraf.MetaGraf) argoapp.Application {

	var meta []argoapp.Info

	obj := argoapp.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name(mg),
			Namespace: GetArgoCDApplicationNamespace(),
			Labels:    mg.Metadata.Labels,
		},
		Spec: argoapp.ApplicationSpec{
			Destination: argoapp.ApplicationDestination{
				Server:    "https://kubernetes.default.svc",
				Namespace: params.NameSpace,
			},
			Source: argoapp.ApplicationSource{
				RepoURL:        params.ArgoCDApplicationRepoURL,
				Path:           params.ArgoCDApplicationRepoPath,
				TargetRevision: params.ArgoCDApplicationTargetRevision,
			},

			Project:    params.ArgoCDApplicationProject,
			SyncPolicy: GetArgoCDApplicationSyncPolicy(),
			Info:       meta,
		},
		Operation: nil,
	}
	return obj
}

func OutputArgoCDApplication(obj argoapp.Application) {
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

func StoreArgoCDApplication(obj argoapp.Application) {

	glog.V(2).Infof("ResourceVersion: %v Length: %v", obj.ResourceVersion, len(obj.ResourceVersion))
	glog.V(2).Infof("Namespace: %v", params.NameSpace)

	client := k8sclient.GetArgoCDClient().Applications(params.NameSpace)
	if len(obj.ResourceVersion) > 0 {
		// update
		result, err := client.Update(context.TODO(), &obj, metav1.UpdateOptions{})
		if err != nil {
			glog.Info(err)
		}
		glog.Infof("Updated ArgoCD Application: %v(%v)", result.Name, obj.Name)
	} else {
		result, err := client.Create(context.TODO(), &obj, metav1.CreateOptions{})
		if err != nil {
			glog.Info(err)
		}
		glog.Infof("Created ArgoCD Application: %v(%v)", result.Name, obj.Name)
	}
}
