package cmd

import (
	"fmt"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	"testing"
)

func TestInspectProperties(t *testing.T) {
	// Mocking a metagraf.json
	mg := metagraf.MetaGraf{}
	envvars := []metagraf.EnvironmentVar{}
	envvars = append(envvars, metagraf.EnvironmentVar{
		Name:     "REQUIRED",
		Required: true,
		Type:     "string",
		Default:  "defaultvalue",
	},
		metagraf.EnvironmentVar{
			Name:     "OPTIONAL",
			Required: false,
			Type:     "string",
			Default:  "defaultvalue",
		},
	)
	mg.Spec.Environment.Local = envvars

	t.Run("FailWithoutFileProperties", func(t *testing.T) {
		modules.Variables = GetCmdProperties(mg.GetProperties())

		fmt.Printf("%+v\n", modules.Variables)
		if ValidateProperties(modules.Variables) {
			t.Errorf("ValidateProperties should not have returned true.")
		}
	})

	t.Run("WithDefaultsNoFileProperties", func(t *testing.T) {
		Defaults = true
		modules.Variables = GetCmdProperties(mg.GetProperties())

		mapkey := "local|REQUIRED"
		expected := "defaultvalue"

		fmt.Printf("%+v\n", modules.Variables)
		if !ValidateProperties(modules.Variables) {
			t.Errorf("ValidateProperties should have returned true. Returned false.")
		}

		if _, ok := modules.Variables[mapkey]; !ok {
			t.Errorf("Expected MGProperty with key (%v) not fond.", mapkey)
		}

		property := modules.Variables[mapkey]
		if property.Value != expected {
			t.Errorf("Default value not as expected: %v, got: %v", expected, property.Value)
		}

	})

	t.Run("WithServiceDiscoveryConfiguration", func(t *testing.T) {

		expected := "someservice.example.com"
		mapkey := "local|SOMESERVICEV1_SERVICE_HOST"
		key := "SOMESERVICEV1_SERVICE_HOST"

		fileprops := make(metagraf.MGProperties)
		fileprops["local|REQUIRED"] = metagraf.MGProperty{Source: "local", Key: "REQUIRED", Value: "MyValue"}
		fileprops[mapkey] = metagraf.MGProperty{Source: "local", Key: key, Value: expected}
		modules.Variables = MergeAndValidateProperties(GetCmdProperties(mg.GetProperties()), fileprops, true)

		fmt.Printf("FileProps: %+v\n", fileprops)
		fmt.Printf("Props: %+v\n", modules.Variables)
		if !ValidateProperties(modules.Variables) {
			t.Errorf("InspectFileProperties should have returned true.")
		}

		if _, ok := modules.Variables[mapkey]; !ok {
			t.Errorf("Expected MGProperty with key (%v) not fond.", mapkey)
		}
		property := modules.Variables[mapkey]
		if property.Value != expected {
			t.Errorf("MGProperty for key: %v, had the incorrect value: %v, expected, %v", mapkey, property.Value, expected)
		}
	})
}
