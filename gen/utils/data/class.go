package data

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"unicode"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils"
	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
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
	Children []string
	// List of all possible parent classes.
	ContainedBy []string
	// Deprecated resources include a warning the resource and datasource schemas.
	Deprecated bool
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
	// Indicates that when the class is included in a resource as a child it can only be configured once.
	// This is used to determine the type of the nested attribute to be a map or list.
	IsSingleNested bool
	// The platform type is used to indicate on which APIC platform the class is available.
	PlatformType PlatformTypeEnum
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
	// The relationship information of the class.
	// If the class is a relationship class, it will contain the information about the relationship.
	Relation Relation
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

type PlatformTypeEnum int

// The enumeration options of the Platform type.
const (
	// Apic indicates that the class is available on the on-premises version of APIC.
	Apic PlatformTypeEnum = iota + 1
	// Both indicates that the class is available on both the on-premises and cloud versions of APIC.
	Both
	// Cloud indicates that the class is available on the cloud version of APIC.
	Cloud
)

type Relation struct {
	// The class from which the relationship is defined.
	FromClass string
	// Indicates if _from_ should be included in the resource name.
	//  ex. "fvRsBdToOut" would be "data_source_aci_relation_from_bridge_domain_to_l3_outside".
	IncludeFrom bool
	// Indicates if the class is a relational class.
	RelationalClass bool
	// The class to which the relationship points.
	ToClass string
	// The type of the relationship.
	Type RelationshipTypeEnum
}

// The enumeration options of the relationship type.
type RelationshipTypeEnum int

const (
	// Named indicates that the relationship is a named relation.
	Named RelationshipTypeEnum = iota + 1
	// Explicit indicates that the relationship is an explicit relation.
	Explicit
	// Undefined indicates that the relationship type is unknown.
	Undefined RelationshipTypeEnum = iota
)

func setRelationshipTypeEnum(relationType string) (RelationshipTypeEnum, error) {
	switch relationType {
	case "named":
		return Named, nil
	case "explicit":
		return Explicit, nil
	default:
		return Undefined, fmt.Errorf("undefined relationship type '%s'", relationType)
	}
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
	Max provider.Version
	// The minimum version of the range.
	// This is the first version of the range.
	// The version is in the format "4.2(7f)".
	Min provider.Version
}

func NewClass(className string, ds *DataStore) (*Class, error) {
	genLogger.Trace(fmt.Sprintf("Creating new class struct with class name: %s.", className))
	// Splitting the class name into the package and short name.
	packageName, shortName, err := splitClassNameToPackageNameAndShortName(className)
	if err != nil {
		return nil, err
	}

	class := Class{
		ClassName:             className,
		ClassNameShort:        shortName,
		ClassNameForFunctions: cases.Title(language.Und, cases.NoLower).String(className),
		ClassNamePackage:      packageName,
		Properties:            make(map[string]*Property),
	}

	genLogger.Trace(fmt.Sprintf("Successfully created new class struct with class name: %s.", className))

	err = class.loadMetaFile()
	if err != nil {
		return nil, err
	}

	err = class.setClassData(ds)
	if err != nil {
		return nil, err
	}

	return &class, nil
}

func (c *Class) loadMetaFile() error {
	genLogger.Debug(fmt.Sprintf("Loading meta file for class '%s'.", c.ClassName))

	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s.json", constMetaPath, c.ClassName))
	if err != nil {
		return fmt.Errorf("failed to load meta file for class '%s': %s", c.ClassName, err.Error())
	}

	genLogger.Trace(fmt.Sprintf("Parsing meta file for class '%s'.", c.ClassName))
	// For now, the file content is unmarshalled into a map[string]interface{} and then set the class data.
	// This is done because we add logic on top of the file content to set the class data.
	// ENHANCEMENT: investigate if we can unmarshal the file content directly into a class struct specific for meta.
	var metaFileContent map[string]interface{}
	err = json.Unmarshal(fileContent, &metaFileContent)
	if err != nil {
		return fmt.Errorf("failed to parse meta file for class '%s': %s", c.ClassName, err.Error())
	}

	c.MetaFileContent = metaFileContent[fmt.Sprintf("%s:%s", c.ClassNamePackage, c.ClassNameShort)].(map[string]interface{})

	genLogger.Debug(fmt.Sprintf("Successfully loaded meta file for class '%s'.", c.ClassName))

	return nil
}

func (c *Class) setClassData(ds *DataStore) error {
	genLogger.Debug(fmt.Sprintf("Setting class data for class '%s'.", c.ClassName))

	var err error

	// TODO: add function to set AllowDelete
	c.setAllowDelete()

	// TODO: add function to set Children
	c.setChildren()

	// TODO: add function to set ContainedBy
	c.setContainedBy()

	// TODO: add placeholder function for Deprecated
	c.setDeprecated()

	// TODO: add placeholder function for DeprecatedVersions
	c.setDeprecatedVersions()

	// TODO: add function to set Documentation
	c.setDocumentation()

	// TODO: add function to set IdentifiedBy
	c.setIdentifiedBy()

	// TODO: add function to set IsMigration
	c.setIsMigration()

	err = c.setRelation()
	if err != nil {
		return err
	}

	// TODO: add function to set IsSingleNested
	c.setIsSingleNested()

	// TODO: add function to set PlatformType
	c.setPlatformType()

	// TODO: add function to set Properties
	c.setProperties()

	// TODO: add function to set RequiredAsChild
	c.setRequiredAsChild()

	err = c.setResourceName(ds)
	if err != nil {
		return err
	}

	// TODO: add function to set RnFormat
	c.setRnFormat()

	// TODO: add function to set Versions
	c.setVersions()

	genLogger.Debug(fmt.Sprintf("Successfully set class data for class '%s'.", c.ClassName))
	return nil
}

func (c *Class) setAllowDelete() {
	// Determine if the class can be deleted.
	genLogger.Debug(fmt.Sprintf("Setting AllowDelete for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set AllowDelete for class '%s'.", c.ClassName))
}

func (c *Class) setChildren() {
	// Determine the child classes for the class.
	genLogger.Debug(fmt.Sprintf("Setting Children for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set Children for class '%s'.", c.ClassName))
}

func (c *Class) setContainedBy() {
	// Determine the parent classes for the class.
	genLogger.Debug(fmt.Sprintf("Setting ContainedBy for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set ContainedBy for class '%s'.", c.ClassName))
}

func (c *Class) setDeprecated() {
	// Determine if the class is deprecated.
	genLogger.Debug(fmt.Sprintf("Setting Deprecated for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set Deprecated for class '%s'.", c.ClassName))
}

func (c *Class) setDeprecatedVersions() {
	// Determine the APIC versions in which the class is deprecated.
	genLogger.Debug(fmt.Sprintf("Setting DeprecatedVersions for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set DeprecatedVersions for class '%s'.", c.ClassName))
}

func (c *Class) setDocumentation() {
	// Determine the documentation specific information for the class.
	genLogger.Debug(fmt.Sprintf("Setting Documentation for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation for class '%s'.", c.ClassName))
}

func (c *Class) setIdentifiedBy() {
	// Determine the identifying properties of the class.
	genLogger.Debug(fmt.Sprintf("Setting IdentifiedBy for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set IdentifiedBy for class '%s'.", c.ClassName))
}

func (c *Class) setIsMigration() {
	// Determine if the class is migrated from previous version of the provider.
	genLogger.Debug(fmt.Sprintf("Setting IsMigration for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set IsMigration for class '%s'.", c.ClassName))
}

func (c *Class) setIsSingleNested() {
	// Determine if the class can only be configured once when used as a nested attribute in a parent resource.
	genLogger.Debug(fmt.Sprintf("Setting IsSingleNested for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set IsSingleNested for class '%s'.", c.ClassName))
}

func (c *Class) setPlatformType() {
	// Determine the platform type of the class.
	genLogger.Debug(fmt.Sprintf("Setting PlatformType for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set PlatformType for class '%s'.", c.ClassName))
}

func (c *Class) setProperties() {
	genLogger.Debug(fmt.Sprintf("Setting properties for class '%s'.", c.ClassName))

	if properties, ok := c.MetaFileContent["properties"]; ok {
		for name, propertyDetails := range properties.(map[string]interface{}) {
			details := propertyDetails.(map[string]interface{})
			// TODO: add logic to set the property data based on ignore/include/exclude overwrites (read-only) from definition files.
			if details["isConfigurable"] == true {
				c.Properties[name] = NewProperty(name, details)
				c.PropertiesAll = append(c.PropertiesAll, name)
				// TODO: add logic to set the required/optional/read-only list logic
			}
		}
	}

	genLogger.Debug(fmt.Sprintf("Successfully set properties for class '%s'.", c.ClassName))
	// TODO: add sorting logic for the properties
}

func (c *Class) setRelation() error {
	// Determine if the class is a relational class.
	genLogger.Debug(fmt.Sprintf("Setting Relation details for class '%s'.", c.ClassName))

	// TODO: add logic to override the relational status from a definition file.
	if relationInfo, ok := c.MetaFileContent["relationInfo"]; ok {
		relationType, err := setRelationshipTypeEnum(relationInfo.(map[string]interface{})["type"].(string))
		if err != nil {
			return err
		}

		c.Relation.FromClass = strings.Replace(relationInfo.(map[string]interface{})["fromMo"].(string), ":", "", -1)
		if strings.Contains(c.ClassName, "To") {
			c.Relation.IncludeFrom = true
		}
		c.Relation.RelationalClass = true
		c.Relation.ToClass = strings.Replace(relationInfo.(map[string]interface{})["toMo"].(string), ":", "", -1)
		c.Relation.Type = RelationshipTypeEnum(relationType)

	}
	genLogger.Debug(fmt.Sprintf("Successfully set Relation details for class '%s'.", c.ClassName))

	return nil
}

func (c *Class) setRequiredAsChild() {
	// Determine if the class is required when defined as a child in a parent resource.
	genLogger.Debug(fmt.Sprintf("Setting RequiredAsChild for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set RequiredAsChild for class '%s'.", c.ClassName))
}

func (c *Class) setResourceName(ds *DataStore) error {
	genLogger.Debug(fmt.Sprintf("Setting resource name for class '%s'.", c.ClassName))

	// TODO: add logic to override the resource name from a definition file.
	// TODO: add logic to override the label from a definition file.
	// TODO: add logic to override the class the nested resource name from the parent
	// 	ex. is fvRsSecInherited with fvESg and fvAEPg.
	// 		fvESg: relation_to_end_point_security_groups
	// 		fvAEPg: relation_to_application_epgs
	if label, ok := c.MetaFileContent["label"]; ok && label != "" {
		if c.Relation.RelationalClass {
			// If the relation includes 'To' in the classname, the resource name will be in the format 'relation_from_{from_class}_to_{to_class}'.
			// If the relation does not include 'To' in the classname, the resource name will be in the format 'relation_to_{to_class}'.
			toClass := getRelationshipResourceName(ds, c.Relation.ToClass)
			if c.Relation.IncludeFrom {
				// If the class is a relational class and the relation includes 'from', the resource name will be in the format 'relation_from_{from_class}_to_{to_class}'.
				c.ResourceName = fmt.Sprintf("relation_from_%s_to_%s", getRelationshipResourceName(ds, c.Relation.FromClass), toClass)
				c.ResourceNameNested = fmt.Sprintf("relation_to_%s", toClass)
			} else {
				// If the class is a relational class and the relation does not include 'from', the resource name will be in the format 'relation_to_{to_class}'.
				c.ResourceName = fmt.Sprintf("relation_to_%s", toClass)
				c.ResourceNameNested = fmt.Sprintf("relation_to_%s", toClass)
			}
		} else {
			c.ResourceName = utils.Underscore(label.(string))
			c.ResourceNameNested = c.ResourceName
		}

		// If the class is a relational class and the relation has identifiers, the plural form of the resource name will be set.
		if len(c.IdentifiedBy) != 0 {
			pluralForm, err := utils.Plural(c.ResourceName)
			if err != nil {
				return err
			}
			c.ResourceNameNested = pluralForm
		}
	} else {
		return fmt.Errorf("failed to set resource name for class '%s': label not found", c.ClassName)
	}

	genLogger.Debug(fmt.Sprintf("Successfully set resource name '%s' for class '%s'.", c.ResourceName, c.ClassName))
	return nil
}

func (c *Class) setRnFormat() {
	// Determine the relative name (RN) format of the class.
	genLogger.Debug(fmt.Sprintf("Setting RnFormat for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set RnFormat for class '%s'.", c.ClassName))
}

func (c *Class) setVersions() {
	// Determine the supported APIC versions for the class.
	genLogger.Debug(fmt.Sprintf("Setting Versions for class '%s'.", c.ClassName))
	genLogger.Debug(fmt.Sprintf("Successfully set Versions for class '%s'.", c.ClassName))
}

func splitClassNameToPackageNameAndShortName(className string) (string, string, error) {
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
		genLogger.Error(fmt.Sprintf("Failed to split class name '%s' for name space separation.", className))
		return "", "", fmt.Errorf("failed to split class name '%s' for name space separation", className)
	}

	return packageName, shortName, nil

}

func getRelationshipResourceName(ds *DataStore, toClass string) string {
	err := ds.loadClass(toClass)
	if err != nil {
		// If the class is not found, try returning the resource name of the class from global definition file.
		if resourceName, ok := ds.GlobalMetaDefinition.NoMetaFile[toClass]; ok {
			genLogger.Debug(fmt.Sprintf("Failed to load class '%s'. Using resource name from global definition file: %s.", toClass, resourceName))
			return resourceName
		}
		if !slices.Contains(failedToLoadClasses, toClass) {
			// If the class is not found and it is not already in the failed to load classes, add it to the list to be errored.
			failedToLoadClasses = append(failedToLoadClasses, toClass)
		}
		return toClass
	}
	// If the class is found, return the resource name of the class.
	return ds.Classes[toClass].ResourceName
}
