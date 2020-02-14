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
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"metagraf/pkg/metagraf"
)

func GenKnativeService(mg *metagraf.MetaGraf) {

	obj := knservingv1.Service{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: knservingv1.ServiceSpec{
			ConfigurationSpec: knservingv1.ConfigurationSpec{
				Template: knservingv1.RevisionTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name:                       "",
						GenerateName:               "",
						Namespace:                  "",
						SelfLink:                   "",
						UID:                        "",
						ResourceVersion:            "",
						Generation:                 0,
						CreationTimestamp:          metav1.Time{},
						DeletionTimestamp:          nil,
						DeletionGracePeriodSeconds: nil,
						Labels:                     nil,
						Annotations:                nil,
						OwnerReferences:            nil,
						Finalizers:                 nil,
						ClusterName:                "",
						ManagedFields:              nil,
					},
					Spec: knservingv1.RevisionSpec{
						PodSpec: corev1.PodSpec{
							Volumes:                       nil,
							InitContainers:                nil,
							Containers:                    nil,
							RestartPolicy:                 "",
							TerminationGracePeriodSeconds: nil,
							ActiveDeadlineSeconds:         nil,
							DNSPolicy:                     "",
							NodeSelector:                  nil,
							ServiceAccountName:            "",
							DeprecatedServiceAccount:      "",
							AutomountServiceAccountToken:  nil,
							NodeName:                      "",
							HostNetwork:                   false,
							HostPID:                       false,
							HostIPC:                       false,
							ShareProcessNamespace:         nil,
							SecurityContext:               nil,
							ImagePullSecrets:              nil,
							Hostname:                      "",
							Subdomain:                     "",
							Affinity:                      nil,
							SchedulerName:                 "",
							Tolerations:                   nil,
							HostAliases:                   nil,
							PriorityClassName:             "",
							Priority:                      nil,
							DNSConfig:                     nil,
							ReadinessGates:                nil,
							RuntimeClassName:              nil,
							EnableServiceLinks:            nil,
						},
						ContainerConcurrency: nil,
						TimeoutSeconds:       nil,
					},
				},
			},
			RouteSpec: knservingv1.RouteSpec{
				Traffic: []knservingv1.TrafficTarget{},
			},
		},
		Status: knservingv1.ServiceStatus{
			Status:                    duckv1.Status{},
			ConfigurationStatusFields: knservingv1.ConfigurationStatusFields{},
			RouteStatusFields:         knservingv1.RouteStatusFields{},
		},
	}
	fmt.Println(obj)

}
