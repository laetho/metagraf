package cmd

import (
	"metagraf/pkg/metagraf"
	"metagraf/pkg/modules"
	"testing"
)

func TestInspectProperties(t *testing.T) {
	mg := metagraf.MetaGraf{}
	mg.

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
}