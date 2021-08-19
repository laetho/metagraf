/*
Copyright 2021 The metaGraf Authors

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

package metagraf

import (
	"fmt"
	"testing"
)

func TestDefaultServicePorts(t *testing.T) {
	mg := MetaGraf{}
	input := mg.DefaultServicePorts()

	expectedName := "http"
	expectePort := int32(80)
	expectedTargetPort := "8080"

	t.Run("Got Default ServicePorts", func(t *testing.T) {
		if len(input) == 0 {
			t.Error("Found no default ServicePort!")
		}
	})

	t.Run("Check default ServicePort values.", func(t *testing.T) {

		var foundName = false
		var foundPort = false
		var foundTargetPort = false

		for _, port := range input {
			if port.Name == expectedName {
				foundName = true
			}
			if port.Port == expectePort {
				foundPort = true
			}
			if port.TargetPort.StrVal == expectedTargetPort {
				foundTargetPort = true
			}
		}

		if !foundName {
			t.Error("Did not find a default ServicePort named:", expectedName)
		}
		if !foundPort {
			t.Error("Did not find a ServicePort with Port:", expectePort)
		}
		if !foundTargetPort {
			t.Error("Did not find a ServicePort TargetPort:", expectedTargetPort)
		}
	})
}

func TestServicePortsBySpec(t *testing.T) {
	mg := MetaGraf{}
	mg.Spec.Ports = make(map[string]int32)
	mg.Spec.Ports["http"]=8080
	mg.Spec.Ports["metrics"]=9090

	ports := mg.ServicePortsBySpec()

	fmt.Println(ports)

	t.Run("Found named http port", func(t *testing.T) {
		found := false
		for _, sp := range ports {
			if sp.Name == "http" {
				found = true
			}
		}
		if !found {
			t.Error("Failed to find ServicePort with name \"http\"")
		}
	})

	t.Run("Found named metrics port", func(t *testing.T) {
		found := false
		for _, sp := range ports {
			if sp.Name == "metrics" {
				found = true
			}
		}
		if !found {
			t.Error("Failed to find ServicePort with name \"http\"")
		}
	})
}