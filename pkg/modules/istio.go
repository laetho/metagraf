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
	"fmt"
	"metagraf/pkg/metagraf"
	istionetv1alpha3 "istio.io/api/networking/v1alpha3"
)

func GenIstioServiceEntry(mg *metagraf.MetaGraf) {
	svc := istionetv1alpha3.ServiceEntry{
		Hosts:                nil,
		Addresses:            nil,
		Ports:                nil,
		Location:             0,
		Resolution:           0,
		Endpoints:            nil,
		ExportTo:             nil,
		SubjectAltNames:      nil,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	fmt.Println(svc)
}

func GenIstioVirtualService(mg *metagraf.MetaGraf) {
	vs := istionetv1alpha3.VirtualService{
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
	fmt.Println(vs)
}
