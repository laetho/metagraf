package modules

import (
	"testing"
)

func TestDefaultServicePorts(t *testing.T) {
	input := defaultServicePorts()

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