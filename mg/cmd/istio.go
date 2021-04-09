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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/laetho/metagraf/pkg/metagraf"
	"os"
	log "k8s.io/klog"
	"github.com/laetho/metagraf/pkg/modules"
)

func init() {
	RootCmd.AddCommand(istioCmd)
	istioCmd.AddCommand(istioCreateCmd)
	istioCreateCmd.AddCommand(istioCreateServiceEntryCmd)
	istioCreateCmd.AddCommand(istioCreateVirtualServiceCmd)
}

var istioCmd = &cobra.Command{
	Use:   "istio",
	Short: "istio subcommands",
	Long:  `Subcommands for Istio ServiceMesh`,
}

var istioCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "istio create subcommands",
	Long:  `Create Subcommands for Istio ServiceMesh`,
}

//
var istioCreateServiceEntryCmd = &cobra.Command{
	Use:   "serviceentry <metagraf>",
	Short: "create istio service resource",
	Long:  `
ServiceEntry enables adding additional entries into Istio’s internal service registry, 
so that auto-discovered services in the mesh can access/route to these manually specified services.
A service entry describes the properties of a service (DNS name, VIPs, ports, protocols, endpoints). 
These services could be external to the mesh (e.g., web APIs) or mesh-internal services that are 
not part of the platform’s service registry (e.g., a set of VMs talking to services in Kubernetes).
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error(StrMissingMetaGraf)
			fmt.Println(StrMissingMetaGraf)
			os.Exit(1)
		}
		FlagPassingHack()
		mg := metagraf.Parse(args[0])
		modules.GenIstioServiceEntry(&mg)

	},
}

//
var istioCreateVirtualServiceCmd = &cobra.Command{
	Use:   "virtualservice <metagraf>",
	Short: "create istio VirtualService resource",
	Long:  `
A VirtualService defines a set of traffic routing rules to apply when a host is addressed. 
Each routing rule defines matching criteria for traffic of a specific protocol. If the 
traffic is matched, then it is sent to a named destination service (or subset/version of 
it) defined in the registry.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Error(StrMissingMetaGraf)
			fmt.Println(StrMissingMetaGraf)
			os.Exit(1)
		}
		FlagPassingHack()
		mg := metagraf.Parse(args[0])
		modules.GenIstioVirtualService(&mg)

	},
}