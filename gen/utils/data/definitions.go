package data

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type GlobalMetaDefinition struct {
	// A list of class names that should always be included as children regardless of standard inclusion logic.
	AlwaysIncludeAsChild []string `yaml:"always_include_as_child"`
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
	// Overrides the default deletion behavior from meta file. Set to "never" to prevent deletion of the class.
	// The value "never" is used to keep the input consistent with the meta data file.
	AllowDelete string `yaml:"allow_delete"`
	// Indicates that the resource and datasource are deprecated. A deprecation warning will be included in the schemas.
	Deprecated bool `yaml:"deprecated"`
	// The deprecated APIC versions for the class. Format: "1.0(1e)-" or "4.2(7f)-4.2(7w),5.2(1g)-".
	// Used to indicate versions where the class is deprecated but still functional.
	DeprecatedVersions string `yaml:"deprecated_versions"`
	// A list of child class names to exclude from the Children list.
	ExcludeChildren []string `yaml:"exclude_children"`
	// A list of parent class names to exclude from the Parents list.
	ExcludeParents []string `yaml:"exclude_parents"`
	// A list of identifier attributes for the class.
	IdentifiedBy []string `yaml:"identified_by"`
	// A list of child class names to include in the Children list outside of the standard inclusion logic.
	IncludeChildren []string `yaml:"include_children"`
	// A list of parent class names to include in the Parents list outside of the standard inclusion logic.
	IncludeParents []string `yaml:"include_parents"`
	// Overrides the default single nested behavior. When true, the class is treated as a single nested attribute
	// when defined as a child in a parent resource, regardless of whether it has identifying properties.
	IsSingleNestedWhenDefinedAsChild bool `yaml:"is_single_nested_when_defined_as_child"`
	// Overrides the versions from the meta file. Format: "1.0(1e)-" or "4.2(7f)-4.2(7w),5.2(1g)-".
	SupportedVersions string `yaml:"supported_versions"`
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
