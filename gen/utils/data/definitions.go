package data

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type GlobalMetaDefinition struct {
	// A map containing class names as keys and their corresponding resource names as values.
	// This is used to search for the resource name of a class when it is not defined in meta directory.
	NoMetaFile map[string]string `yaml:"no_meta_file"`
}

func loadGlobalMetaDefinition() GlobalMetaDefinition {
	definition, err := os.ReadFile(constGlobalDefinitionFilePath)
	if err != nil {
		genLogger.Fatal("A file 'global.yaml' is required to be defined in the definitions folder.")
	}

	var definitionGlobalMetaData GlobalMetaDefinition
	err = yaml.Unmarshal(definition, &definitionGlobalMetaData)
	if err != nil {
		genLogger.Fatal(err.Error())
	}

	return definitionGlobalMetaData
}

type ClassDefinition struct {
	AllowDelete string `yaml:"allow_delete"`
}

func loadClassDefinition(className string) ClassDefinition {
	classDefinitionPath := fmt.Sprintf("%s/%s.yaml", constDefinitionsPath, className)
	var classDefinitionData ClassDefinition

	classDefinitionBytes, err := os.ReadFile(classDefinitionPath)
	if err != nil {
		genLogger.Debug(fmt.Sprintf("The file '%s' was not found in the definitions folder.", classDefinitionPath))
		return classDefinitionData
	}

	err = yaml.Unmarshal(classDefinitionBytes, &classDefinitionData)
	if err != nil {
		genLogger.Fatal(err.Error())
	}

	return classDefinitionData
}
