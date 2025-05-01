package data

import (
	"encoding/json"
	"fmt"
	"os"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Class struct {
	// This is used to prevent the deletion of the class if it is not allowed on APIC.
	AllowDelete bool
	// Full name of the class, ex "fvTenant".
	ClassName string
	// Capitalized name of the class, ex "FvTenant".
	ClassNameForFunctions string
	// Package part of the class, ex "fv".
	ClassNamePackage string
	// Name part of the class, ex "Tenant".
	ClassNameShort string
	// List of all child classes which are included inside the resource.
	// When looping over maps in golang the order of the returned elements is random, thus list is used for order consistency.
	Children []Child
	// List of all possible parent classes.
	ContainedBy []string
	// The APIC versions in which the class is deprecated.
	DeprecatedVersions []VersionRange
	// Documentation specific information for the class.
	Documentation ClassDocumentation
	// List of all identifying properties of the class.
	// These are properties that are part of the relative name (RN) format, ex 'tn-{name}' where name is the identifying property.
	IdentifiedBy []string
	// Indicates that the class is migrated from previous version of the provider.
	// This is used to determine if legacy attributes have to be exposed in the resource.
	IsMigration bool
	// Indicates that the class is a relationship class.
	IsRelational bool
	// Indicates that when the class is included in a resource as a child it can only be configured once.
	// This is used to determine the type of the nested attribute to be a map or list.
	IsSingleNested bool
	// Indicates for which APIC platform the class is available, cloud vs onprem.
	// TODO: investigate enum type for this.
	Platform string
	// A map containing all the properties of the class.
	Properties map[string]*Property
	// Each property is categorised in all, required, optional, or read-only.
	// When looping over maps in golang the order of the returned elements is random, thus lists are used for order consistency.
	PropertiesAll      []string
	PropertiesRequired []string
	PropertiesOptional []string
	PropertiesReadOnly []string
	// The full content from the meta file.
	// Storing the content proactively in case we need to access the data at a later stage.
	MetaFileContent map[string]interface{}
	// Indicates if the class is required when defined as a child in a parent resource.
	RequiredAsChild bool
	// The resource name is the name of the resource in the provider, ex "aci_tenant".
	ResourceName string
	// The nested attribute name when part of a parent resource.
	// ex. "relation_from_bridge_domain_to_netflow_monitor_policy" would translate to "relation_to_netflow_monitor_policy".
	ResourceNameNested string
	// The relative name (RN) format of the class, ex "tn-{name}".
	RnFormat string
	// The supported APIC versions for the class.
	// Each version range is separated by a comma, ex "4.2(7f)-4.2(7w),5.2(1g)-".
	// The first version is the minimum version and the second version is the maximum version.
	// A dash at the end of a range (ex. 4.2(7f)-) indicates that the class is supported from the first version to the latest version.
	Versions []VersionRange
}

type Child struct {
	// The name of the child class, ex "fvTenant".
	// This is used as the key for the map of all classes.
	ClassName string
	// When it is a relationship class, this is the class to which the relationship points.
	PointsToClass string
}

type ClassDocumentation struct {
	// List of all child classes which are not included inside the resource but have a separate resource
	// Used to reference child resource in documentation
	Children []string
	// The description of the class, which is used at the top of the documentation.
	Description string
	// List of DN formats
	DnFormats []string
	// List of notes to be added to the top of the documentation
	Notes []string
	// List of warnings to be added to the top of the documentation
	Warnings []string
}

type VersionRange struct {
	// The maximum version of the range.
	// This is the second version of the range.
	// The version is in the format "4.2(7w)".
	// A dash at the end of a range (ex. 4.2(7f)-) indicates that the class is supported from the first version to the latest version.
	Max string
	// The minimum version of the range.
	// This is the first version of the range.
	// The version is in the format "4.2(7f)".
	Min string
}

func NewClass(className string) *Class {
	genLogger.Trace(fmt.Sprintf("Creating new class struct with class name: %s.", className))
	// Splitting the class name into the package and short name.
	packageName, shortName := splitClassNameToPackageNameAndShortName(className)
	class := Class{
		ClassName:             className,
		ClassNameShort:        shortName,
		ClassNameForFunctions: cases.Title(language.Und, cases.NoLower).String(className),
		ClassNamePackage:      packageName,
		Properties:            make(map[string]*Property),
	}
	genLogger.Trace(fmt.Sprintf("Succesfully created new class struct with class name: %s.", className))
	class.loadMetaFile()
	class.setClassData()
	return &class
}

func (c *Class) loadMetaFile() {
	genLogger.Debug(fmt.Sprintf("Loading meta file for class '%s'.", c.ClassName))

	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s.json", constMetaPath, c.ClassName))
	if err != nil {
		genLogger.Fatal(fmt.Sprintf("Error during reading of file: %s.", err.Error()))
	}

	genLogger.Trace(fmt.Sprintf("Parsing meta file for class '%s'.", c.ClassName))
	// For now, the file content is unmarshalled into a map[string]interface{} and then set the class data.
	// This is done because we add logic on top of the file content to set the class data.
	// ENHANCEMENT: investigate if we can unmarshal the file content directly into a class struct specific for meta.
	var metaFileContent map[string]interface{}
	err = json.Unmarshal(fileContent, &metaFileContent)
	if err != nil {
		genLogger.Fatal(fmt.Sprintf("Error during Unmarshal(): %s.", err.Error()))
	}

	c.MetaFileContent = metaFileContent[fmt.Sprintf("%s:%s", c.ClassNamePackage, c.ClassNameShort)].(map[string]interface{})

	genLogger.Debug(fmt.Sprintf("Succesfully loaded meta file for class '%s'.", c.ClassName))
}

func (c *Class) setClassData() {
	genLogger.Debug(fmt.Sprintf("Setting class data for class '%s'.", c.ClassName))

	c.setResourceName()

	// TODO: add functions to set the other class data

	if properties, ok := c.MetaFileContent["properties"]; ok {
		c.setProperties(properties.(map[string]interface{}))
	}

	genLogger.Debug(fmt.Sprintf("Succesfully set class data for class '%s'.", c.ClassName))
}

func (c *Class) setResourceName() {
	genLogger.Trace(fmt.Sprintf("Implement setting resourceName for class '%s'.", c.ClassName))
}

func (c *Class) setProperties(properties map[string]interface{}) {
	genLogger.Debug(fmt.Sprintf("Setting properties for class '%s'.", c.ClassName))
	for name, propertyDetails := range properties {
		details := propertyDetails.(map[string]interface{})
		// TODO: add logic to set the property data based on ignore/include/exclude overwrites (read-only) from definition files.
		if details["isConfigurable"] == true {
			c.Properties[name] = NewProperty(name, details)
			c.PropertiesAll = append(c.PropertiesAll, name)
			// TODO: add logic to set the required/optional/read-only list logic
		}
	}

	genLogger.Debug(fmt.Sprintf("Succesfully set properties for class '%s'.", c.ClassName))
	// TODO: add sorting logic for the properties
}

func splitClassNameToPackageNameAndShortName(className string) (string, string) {
	// Splitting the class name into the package and short name.
	// The package and short names are used for the meta file download, documentation links and lookup in the raw data.
	var shortName, packageName string
	genLogger.Trace(fmt.Sprintf("Splitting class name '%s' for name space separation.", className))
	for index, character := range className {
		if unicode.IsUpper(character) {
			shortName = className[index:]
			packageName = className[:index]
			break
		}
	}

	genLogger.Debug(fmt.Sprintf("Class name '%s' got split into package name '%s' and short name '%s'.", className, packageName, shortName))

	if packageName == "" || shortName == "" {
		genLogger.Fatal(fmt.Sprintf("Failed to split class name '%s' for name space separation.", className))
	}

	return packageName, shortName

}
