/*
Copyright 2020 The MetaGraph Authors

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
	"metagraf/pkg/metagraf"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	istionetv1alpha3 "istio.io/api/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
)

func GenIstioVirtualService(mg *metagraf.MetaGraf, namespace string) {

	vs := istionetv1alpha3.VirtualService{
		v1.TypeMeta{},
		v1.ObjectMeta{}
		Hosts:                nil,
		Gateways:             nil,
		Http:                 nil,
		Tls:                  nil,
		Tcp:                  nil,
		ExportTo:             nil,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}

}