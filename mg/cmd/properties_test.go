package cmd

import (
	"bytes"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestSimpleLocalSourceProperty(t *testing.T) {
	var propertyFileContent bytes.Buffer
	propertyFileContent.WriteString("local|Nisse=Nasse")

	actualResult := ParseProps(&propertyFileContent)
	actualProperty := actualResult["local|Nisse"]

	expectedProperty := metagraf.MGProperty{
		Source: "local",
		Key: "Nisse",
		Value: "Nasse",
	}

	assert.Len(t, actualResult, 1, "Expected 1 MGProperty, but got %d", len(actualResult))
	assert.Equal(t, expectedProperty, actualProperty)
}

func TestPropertyWithMultipleEqualsSigns(t *testing.T) {
	var propertyFileContent bytes.Buffer
	propertyFileContent.WriteString("local|JAVA_OPTIONS=-Dconfig.file=/config/application.properties")

	actualResult := ParseProps(&propertyFileContent)
	actualProperty := actualResult["local|JAVA_OPTIONS"]

	expectedProperty := metagraf.MGProperty{
		Source: "local",
		Key: "JAVA_OPTIONS",
		Value: "-Dconfig.file=/config/application.properties",
	}

	assert.Len(t, actualResult, 1, "Expected 1 MGProperty, but got %d", len(actualResult))
	assert.Equal(t, expectedProperty, actualProperty)
}


func TestMissingSourceHintPropertyFunctionShouldExit(t *testing.T) {
	var propertyFileContent bytes.Buffer
	propertyFileContent.WriteString("Nisse=Nasse")

	if os.Getenv("BE_CRASHER") == "1" {
		ParseProps(&propertyFileContent)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMissingSourceHintProperty")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestMultipleProperties(t *testing.T) {
	var propertyFileContent bytes.Buffer
	propertyFileContent.WriteString(
		"local|JAVA_OPTIONS=-Dconfig.file=/config/application.properties\n" +
			"build|MAVEN_OPS=-xms123m,-aaa4311G\n")

	actualResult := ParseProps(&propertyFileContent)
	actualProperty1 := actualResult["local|JAVA_OPTIONS"]
	actualProperty2 := actualResult["build|MAVEN_OPS"]

	expectedProperty1 := metagraf.MGProperty{
		Source: "local",
		Key: "JAVA_OPTIONS",
		Value: "-Dconfig.file=/config/application.properties",
	}

	expectedProperty2 := metagraf.MGProperty{
		Source: "build",
		Key: "MAVEN_OPS",
		Value: "-xms123m,-aaa4311G",
	}

	assert.Len(t, actualResult, 2, "Expected 2 MGProperties, but got %d", len(actualResult))
	assert.Equal(t, expectedProperty1, actualProperty1)
	assert.Equal(t, expectedProperty2, actualProperty2)
}
