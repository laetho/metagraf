package test

import (
	"metagraf/mg/cmd"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	"testing"
)

func TestGenerateJvmSysPropValue(t *testing.T) {
	mg := metagraf.MetaGraf{}
	defaultValue := "default_value"
	mg.Spec.Config = append(mg.Spec.Config, metagraf.Config{
		Name: "example.jvm.prop",
		Type: "JVM_SYS_PROP",
		Options: append([]metagraf.ConfigParam{}, metagraf.ConfigParam{
			Name:     "my.test.prop",
			Required: true,
			Type:     "string",
			Default:  defaultValue,
		}),
	})

	t.Run("FromDefaults", func(t *testing.T) {
		modules.Defaults = true

		generatedEnvVar := modules.GenEnvVar_JVM_SYS_PROP(&mg, "JAVA_OPTIONS")

		expected := "-Dmy.test.prop=" + defaultValue
		if generatedEnvVar.Value != expected {
			t.Errorf("Wrong JVM SYS PROP env var generated, expected %v got %v", expected, generatedEnvVar.Value)
		}
	})

	t.Run("FromCmdVars", func(t *testing.T) {
		cvarsValue := "cvars_value"
		cmd.CVars = []string{"my.test.prop=" + cvarsValue}
		modules.Variables = cmd.MergeVars(
			mg.GetVars(),
			cmd.OverrideVars(mg.GetVars()))

		generatedEnvVar := modules.GenEnvVar_JVM_SYS_PROP(&mg, "JAVA_OPTIONS")

		expected := "-Dmy.test.prop=" + cvarsValue
		if generatedEnvVar.Value != expected {
			t.Errorf("Wrong JVM SYS PROP env var generated, expected %v got %v", expected, generatedEnvVar.Value)
		}
	})
}
