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
	// Documentation specific information for the class.
	Documentation ClassDocumentation
	// List of all identifying properties of the class.
	// These are properties that are part if the relative name (RN) format, ex 'tn-{name}' where name is the identifying property.
	IdentifiedBy []string
	// Indicates that the class is migrated from previous version of the provider.
	// This is used to determine if legacy attributes to be exposed in the resource.
	// This is used to determine if the plan modification logic.
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
	// The raw data from the meta file.
	// Storing the raw data proactively in case we need to access the data at a later stage.
	RawMetaData map[string]interface{}
	// Indicates if the class is required when defined as a child in a parent resource.
	Required bool
	// The resource name is the name of the resource in the provider, ex "aci_tenant".
	ResourceName string
	// The nested attribute name when part of a parent resource.
	// ex. "aci_relation_from_bridge_domain_to_netflow_monitor_policy" would translate to "aci_relation_to_netflow_monitor_policy".
	ResourceNameNested string
	// The relative name (RN) format of the class, ex "tn-{name}".
	RnFormat string
	// The supported APIC versions for the class.
	// Each version range is separated by a comma, ex "4.2(7f)-4.2(7w),5.2(1g)-".
	// The first version is the minimum version and the second version is the maximum version.
	// A dash at the end of a range (ex. 4.2(7f)-) indicates that the class is supported from the first version to the latest version.
	Versions string
}

type Child struct {
	// The name of the child class, ex "fvTenant".
	// This is used as the key for the map of all classes.
	ClassName string
	// When it is a relationship class, this is the class to which the relationship points.
	PointsToClass string
}

type ClassDocumentation struct {
	Description string
	// List of all child classes which are not included inside the resource but have a separate resource
	// Used to reference child resource in documentation
	Children []string
	// List of notes to be added to the top of the documentation
	Notes []string
	// List of warnings to be added to the top of the documentation
	Warnings []string
	// List of DN formats
	DnFormats []string
}

func NewClass(className string) *Class {
	genLogger.Trace(fmt.Sprintf("Creating new class struct for class: %s.", className))

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
	if shortName == "" || packageName == "" {
		genLogger.Fatal(fmt.Sprintf("Failed to split class name '%s' for name space separation.", className))
	}

	return &Class{
		ClassName:             className,
		ClassNameShort:        shortName,
		ClassNameForFunctions: cases.Title(language.Und, cases.NoLower).String(className),
		ClassNamePackage:      packageName,
		Properties:            make(map[string]*Property),
	}
}

func (c *Class) LoadMetaFile() {
	genLogger.Trace(fmt.Sprintf("Loading meta file for class '%s'.", c.ClassName))

	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s.json", metaPath, c.ClassName))
	if err != nil {
		genLogger.Fatal(fmt.Sprintf("Error during reading of file: %s.", err.Error()))
	}

	genLogger.Trace(fmt.Sprintf("Parsing meta file for class '%s'.", c.ClassName))
	// For now, the raw data is unmarshalled into a map[string]interface{} and then set the class data.
	// This is done because we add logic on top of the raw data to set the class data.
	// ENHANCEMENT: investigate if we can unmarshal the raw data directly into a class struct specific for meta.
	var classDetails map[string]interface{}
	err = json.Unmarshal(fileContent, &classDetails)
	if err != nil {
		genLogger.Fatal(fmt.Sprintf("Error during Unmarshal(): %s.", err.Error()))
	}

	c.RawMetaData = classDetails[fmt.Sprintf("%s:%s", c.ClassNamePackage, c.ClassNameShort)].(map[string]interface{})
	c.setClassData()
}

func (c *Class) setClassData() {
	genLogger.Trace(fmt.Sprintf("Setting class data for class '%s'.", c.ClassName))

	c.setResourceName()

	// TODO: add functions to set the other class data

	if properties, ok := c.RawMetaData["properties"]; ok {
		c.setProperties(properties.(map[string]interface{}))
	}

	genLogger.Trace(fmt.Sprintf("Succesfully set class data for class '%s'.", c.ClassName))
}

func (c *Class) setResourceName() {
	genLogger.Trace(fmt.Sprintf("Implement setting resourceName for class '%s'.", c.ClassName))
}

func (c *Class) setProperties(properties map[string]interface{}) {
	genLogger.Trace(fmt.Sprintf("Setting properties for class '%s'.", c.ClassName))
	for name, propertyDetails := range properties {
		details := propertyDetails.(map[string]interface{})
		// TODO: add logic to set the property data based on ignore/include/exclude overwrites (read-only) from definition files.
		if details["isConfigurable"] == true {
			c.Properties[name] = NewProperty(name, details)
			c.PropertiesAll = append(c.PropertiesAll, name)
			// TODO: add logic to set the required/optional/read-only list logic
		}
	}

	genLogger.Trace(fmt.Sprintf("Succesfully set properties for class '%s'.", c.ClassName))
	// TODO: add sorting logic for the properties
}
