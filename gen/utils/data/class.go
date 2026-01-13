package data

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils"
	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Class struct {
	// This is used to prevent the deletion of the class if it is not allowed on APIC.
	AllowDelete bool
	// Custom class definition to override class meta properties
	ClassDefinition ClassDefinition
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
	Parents []string
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

	class := Class{
		ClassDefinition:       loadClassDefinition(className),
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

	c.setAllowDelete()

	err = c.setChildren(ds)
	if err != nil {
		return err
	}

	err = c.setParents()
	if err != nil {
		return err
	}

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
	if c.ClassDefinition.AllowDelete == "" {
		isCreatableDeletable, ok := c.MetaFileContent["isCreatableDeletable"]
		if ok && isCreatableDeletable.(string) != "never" {
			c.AllowDelete = true
		}
	} else if c.ClassDefinition.AllowDelete != "never" {
		c.AllowDelete = true
	}
	genLogger.Debug(fmt.Sprintf("The AllowDelete property was successfully set to '%t' for the class '%s'.", c.AllowDelete, c.ClassName))
}

func (c *Class) setChildren(ds *DataStore) error {
	// Determine the child classes for the class.
	genLogger.Debug(fmt.Sprintf("Setting Children for class '%s'.", c.ClassName))

	// Warn if a class is defined in both IncludeChildren and ExcludeChildren.
	for _, includedChild := range c.ClassDefinition.IncludeChildren {
		if slices.Contains(c.ClassDefinition.ExcludeChildren, includedChild) {
			genLogger.Warn(fmt.Sprintf("Child class '%s' is defined in both IncludeChildren and ExcludeChildren for class '%s'. IncludeChildren takes precedence.", includedChild, c.ClassName))
		}
	}

	// Initialize with children from ClassDefinition.IncludeChildren.
	childClasses := c.ClassDefinition.IncludeChildren

	// Retrieve the rnMap from MetaFileContent which holds the rnFormat as key and class name as value.
	// The class name is in the format "{ClassNamePackage}:{ClassName}".
	if rnMap, ok := c.MetaFileContent["rnMap"].(map[string]interface{}); ok {
		for rn, classNameInterface := range rnMap {
			className := classNameInterface.(string)
			// Remove the colon separator from the class name (e.g., "fv:Tenant" -> "fvTenant").
			className, err := sanitizeClassName(className)
			if err != nil {
				return err
			}

			if shouldIncludeChild(rn, className, c.ClassDefinition.ExcludeChildren, ds.GlobalMetaDefinition.AlwaysIncludeAsChild) {
				childClasses = append(childClasses, className)
			}
		}
	}

	// Sort the children for consistent ordering and remove duplicates.
	slices.Sort(childClasses)
	c.Children = slices.Compact(childClasses)

	genLogger.Debug(fmt.Sprintf("Successfully set Children for class '%s'. Found %d children.", c.ClassName, len(c.Children)))
	return nil
}

func (c *Class) setParents() error {
	// Determine the parent classes for the class.
	genLogger.Debug(fmt.Sprintf("Setting Parents for class '%s'.", c.ClassName))

	// Warn if a class is defined in both IncludeParents and ExcludeParents.
	for _, included := range c.ClassDefinition.IncludeParents {
		if slices.Contains(c.ClassDefinition.ExcludeParents, included) {
			genLogger.Warn(fmt.Sprintf("Parent class '%s' is defined in both IncludeParents and ExcludeParents for class '%s'. IncludeParents takes precedence.", included, c.ClassName))
		}
	}

	// Initialize with parents from ClassDefinition.IncludeParents.
	parentClasses := c.ClassDefinition.IncludeParents

	// Retrieve the containedBy from MetaFileContent which holds the parent class names as keys.
	if containedBy, ok := c.MetaFileContent["containedBy"].(map[string]interface{}); ok {
		for classNameWithColon := range containedBy {
			// Remove the colon separator from the class name (e.g., "fv:AEPg" -> "fvAEPg").
			className, err := sanitizeClassName(classNameWithColon)
			if err != nil {
				return err
			}

			// Exclude if the class is explicitly excluded via ClassDefinition.ExcludeParents.
			if !slices.Contains(c.ClassDefinition.ExcludeParents, className) {
				parentClasses = append(parentClasses, className)
			}
		}
	}

	// Sort the parents for consistent ordering and remove duplicates.
	slices.Sort(parentClasses)
	c.Parents = slices.Compact(parentClasses)

	genLogger.Debug(fmt.Sprintf("Successfully set Parents for class '%s'. Found %d parents.", c.ClassName, len(c.Parents)))
	return nil
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

		fromClass, err := sanitizeClassName(relationInfo.(map[string]interface{})["fromMo"].(string))
		if err != nil {
			return err
		}
		c.Relation.FromClass = fromClass

		if strings.Contains(c.ClassName, "To") {
			c.Relation.IncludeFrom = true
		}
		c.Relation.RelationalClass = true

		toClass, err := sanitizeClassName(relationInfo.(map[string]interface{})["toMo"].(string))
		if err != nil {
			return err
		}
		c.Relation.ToClass = toClass
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

// shouldIncludeChild determines if a child class should be included.

func shouldIncludeChild(rn, className string, excludeChildrenFromClassDef, alwaysIncludeFromGlobalDef []string) bool {
	// Determines if a child class should be included, with a default behavior of excluding a child class.
	// A child class is included if any of the following conditions are met (in order of precedence):
	//  1. The className is NOT in the excludeChildrenFromClassDef list (if it is, the child is excluded regardless of other rules)
	//  2. The className is in the alwaysIncludeFromGlobalDef list
	//  3. The RN format starts with "rs" (relation source classes)
	//  4. The RN format does NOT end with "-" (non-named classes without a specific identifier)
	// Exclude if the class is explicitly excluded via ClassDefinition.ExcludeChildren.
	if slices.Contains(excludeChildrenFromClassDef, className) {
		genLogger.Trace(fmt.Sprintf("Child class '%s' excluded via excludeChildrenFromClassDef.", className))
		return false
	}

	// Include if the class is in GlobalMetaDefinition.AlwaysIncludeAsChild.
	if slices.Contains(alwaysIncludeFromGlobalDef, className) {
		genLogger.Trace(fmt.Sprintf("Child class '%s' included via alwaysIncludeFromGlobalDef.", className))
		return true
	}

	// Include classes where the RN starts with "rs" (relation source).
	if strings.HasPrefix(rn, "rs") {
		genLogger.Trace(fmt.Sprintf("Child class '%s' included because RN starts with 'rs'.", className))
		return true
	}

	// Include classes where the RN does NOT end with "-" (non-named classes without a specific identifier).
	if !strings.HasSuffix(rn, "-") {
		genLogger.Trace(fmt.Sprintf("Child class '%s' included because RN '%s' does not end with '-'.", className, rn))
		return true
	}

	genLogger.Trace(fmt.Sprintf("Child class '%s' excluded by default (RN ends with '-').", className))
	return false
}
