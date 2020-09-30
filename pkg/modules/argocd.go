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
	"context"
	argoapp "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"metagraf/mg/k8sclient"
	"metagraf/mg/params"
	"metagraf/pkg/metagraf"
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
	}
	return &asd
}

func GenArgoApplication(mg *metagraf.MetaGraf) {

	meta := []argoapp.Info{}

	obj := argoapp.Application{
		TypeMeta:   metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: 	Name(mg),
			Namespace: GetArgoCDApplicationNamespace(),
		},
		Spec: argoapp.ApplicationSpec{
			Destination: argoapp.ApplicationDestination{
				Server: "https://kubernetes.default.svc",
				Namespace: NameSpace,
			},
			Source: argoapp.ApplicationSource{
				RepoURL: 	params.ArgoCDApplicationRepoURL,
				Path: 		params.ArgoCDApplicationRepoPath,
				Directory:	GetArgoCDSourceDirectory(),
			},

			Project:              params.ArgoCDApplicationProject,
			SyncPolicy:           GetArgoCDApplicationSyncPolicy(),
			Info:                 meta,
		},
		Operation:  nil,
	}

	if !Dryrun {
		StoreArgoCDApplication(obj)
	}
	if Output {
		MarshalObject(obj.DeepCopyObject())
	}
}

func StoreArgoCDApplication(obj argoapp.Application) {

	glog.Infof("ResourceVersion: %v Length: %v", obj.ResourceVersion, len(obj.ResourceVersion))
	glog.Infof("Namespace: %v", NameSpace)

	client := k8sclient.GetArgoCDClient().Applications(NameSpace)
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