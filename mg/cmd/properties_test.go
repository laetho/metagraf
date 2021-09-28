package cmd

import (
	"bytes"
	"fmt"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestSimpleLocalSourceProperty(t *testing.T) {
	propertyFileContent := fileContent("local|Nisse=Nasse")

	actualResult := parseProps(&propertyFileContent)
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
	propertyFileContent := fileContent(
		"local|JAVA_OPTIONS=-Dconfig.file=/config/application.properties")

	actualResult := parseProps(&propertyFileContent)
	actualProperty := actualResult["local|JAVA_OPTIONS"]

	expectedProperty := metagraf.MGProperty{
		Source: "local",
		Key: "JAVA_OPTIONS",
		Value: "-Dconfig.file=/config/application.properties",
	}

	assert.Len(t, actualResult, 1, "Expected 1 MGProperty, but got %d", len(actualResult))
	assert.Equal(t, expectedProperty, actualProperty)
}


func TestMissingSourceHintFunctionShouldExit(t *testing.T) {
	propertyFileContent := fileContent("Nisse=Nasse")

	// Pretty weird way to handle a test case where the function under test does an os.exit(1) during execution
	// Solution is proposed by golang authors, though :-D
	// A better solution is to refactor the function under test to return error as opposed to do a hard exit of the program
	if os.Getenv("BE_CRASHER") == "1" {
		parseProps(&propertyFileContent)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMissingSourceHintFunctionShouldExit")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestMultipleProperties(t *testing.T) {
	propertyFileContent := fileContent(
		"local|JAVA_OPTIONS=-Dconfig.file=/config/application.properties\n" +
				"build|MAVEN_OPS=-xms123m,-aaa4311G\n")

	actualResult := parseProps(&propertyFileContent)
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

func TestPropertiesFile(t *testing.T) {

	tempFile, err := os.CreateTemp("", "sample")
	check(err)
	defer os.Remove(tempFile.Name())
	fmt.Println("Temp file name:", tempFile.Name())

	fileContent := "local|JAVA_OPTIONS=-Dconfig.file=/config/application.properties\n" +
				   "build|MAVEN_OPS=-xms123m,-aaa4311G\n"

	_, err = tempFile.WriteString(fileContent)
	check(err)

	actualResult := propertiesFromFile(tempFile.Name())

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

func fileContent(content string) bytes.Buffer {
	var propertyFileContent bytes.Buffer
	propertyFileContent.WriteString(content)
	return propertyFileContent
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}