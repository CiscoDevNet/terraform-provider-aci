package data

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils"
)

type Class struct {
	// This is used to prevent the deletion of the class if it is not allowed on APIC.
	AllowDelete bool
	// List of all child classes which are included inside the resource.
	// When looping over maps in golang the order of the returned elements is random, thus list is used for order consistency.
	Children []*ClassName
	// Custom class definition to override class meta properties
	ClassDefinition ClassDefinition
	// Deprecated resources include a warning the resource and datasource schemas.
	// Driven by the meta `isDeprecated` flag with an optional definition override (logical OR).
	Deprecated bool
	// The deprecated APIC versions for the class.
	// Used to indicate versions where the class is deprecated but still functional.
	// Driven by the meta `deprecatedSince` value with an optional definition override.
	DeprecatedVersions *Versions
	// Documentation specific information for the class.
	Documentation ClassDocumentation
	// Hidden indicates that the class is no longer accepted by the APIC API.
	// Driven by the meta `isHidden` flag with an optional definition override (logical OR).
	Hidden bool
	// The hidden APIC versions for the class.
	// Driven by the meta `hiddenSince` value with an optional definition override.
	HiddenVersions *Versions
	// List of all identifying properties of the class.
	// These are properties that are part of the relative name (RN) format, ex 'tn-{name}' where name is the identifying property.
	IdentifiedBy []string
	// Records the lineage of this resource: the prior provider or generator it was
	// migrated from. Driven by ClassDefinition.MigrationSource. Drives the
	// documentation migration warning when non-zero. Orthogonal to StateUpgrades:
	// a resource born in the framework can still have StateUpgrades (v0→v1 framework
	// schema bumps) without a MigrationSource set.
	MigrationSource MigrationSourceEnum
	// Resolved per-version state upgrade definitions, copied verbatim from
	// ClassDefinition.StateUpgrades after validation against the resolved
	// Properties / Children sets. Each entry is direct-to-current per the
	// Terraform plugin framework contract.
	StateUpgrades []StateUpgradeDefinition
	// Indicates that when the class is included in a resource as a child it can only be configured once.
	// This is used to determine the type of the nested attribute to be a map or list.
	IsSingleNestedWhenDefinedAsChild bool
	// The full content from the meta file.
	// Storing the content proactively in case we need to access the data at a later stage.
	MetaFileContent map[string]any
	// The class name with all its representations.
	Name *ClassName
	// List of all possible parent classes.
	Parents []*ClassName
	// The platform type is used to indicate on which APIC platform the class is available.
	PlatformType PlatformTypeEnum
	// A map containing all the properties of the class.
	Properties map[string]*Property
	// Each property is categorised in all, required, optional, or read-only.
	// When looping over maps in golang the order of the returned elements is random, thus lists are used for order consistency.
	PropertiesAll      []string
	PropertiesOptional []string
	PropertiesReadOnly []string
	PropertiesRequired []string
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
	// Parsed from the "versions" field in the meta file (e.g., "1.0(1e)-", "4.2(7f)-4.2(7w),5.2(1g)-").
	SupportedVersions *Versions
	// Test dependencies resolved from explicit definitions first (allowing overrides),
	// then auto-filled from Parents and Relation for any remaining references.
	// A slice is used instead of a map to preserve declaration order, which is required for deterministic HCL output
	// (dependencies must appear before the resources that reference them).
	TestDependencies []*TestDependency
	// Test children resolved from child classes' TestValues with optional manual override.
	// A slice is used instead of a map to preserve declaration order for deterministic HCL output.
	TestChildren []*TestChild
}

type Relation struct {
	// The class from which the relationship is defined.
	FromClass *ClassName
	// Indicates if _from_ should be included in the resource name.
	//  ex. "fvRsBdToOut" would be "data_source_aci_relation_from_bridge_domain_to_l3_outside".
	IncludeFrom bool
	// Indicates if the class is a relational class.
	RelationalClass bool
	// The list of concrete target classes the relationship points to.
	// Seeded from the meta `relationInfo.toMo` (single element) and replaced by the
	// `relation_to_classes` definition override when present. A length greater than 1
	// indicates a multi-target (plural) relation, in which case the class definition
	// must provide an explicit `resource_name`.
	ToClasses []*ClassName
	// The type of the relationship.
	Type RelationshipTypeEnum
}

// TestDependency describes a prerequisite resource needed for tests (parent or target).
// Dependencies are recursive: a dependency can itself require other resources to exist.
type TestDependency struct {
	// The class of the dependency resource.
	Class *ClassName
	// The reference value: either a static DN string or a Terraform resource/datasource attribute path.
	Reference string
	// How to interpret the Reference field.
	ReferenceType ReferenceTypeEnum
	// The role of this dependency relative to the resource-under-test.
	// Only meaningful on TOP-LEVEL entries in Class.TestDependencies.
	// Nested entries (inside Dependencies) are pure prerequisites — Role is UndefinedRole.
	Role TestDependencyRoleEnum
	// Recursive dependencies: resources that THIS dependency needs to exist first.
	// These are order-only prerequisites — they have no Role (always UndefinedRole).
	Dependencies []*TestDependency
	// Optional property overrides for the dependency resource's HCL configuration.
	// When empty, the dependency is rendered using its class's own TestValues.Create.
	// When populated, the specified properties override the auto-derived values.
	ConfigOverrides map[string]string
	// Optional child block overrides for the dependency resource's HCL nested blocks.
	// When empty, the dependency is rendered using its class's own TestChildren.
	// When populated, the specified children override the auto-derived nested blocks.
	Children map[string]*TestChild
}

// TestChildInstance holds the property values and nested children for a single child instance.
type TestChildInstance struct {
	// The property values for this child instance, keyed by attribute name.
	Properties map[string]TestValueEntry
	// Children nested inside this specific instance block.
	// Per-instance because in HCL, each block owns its own nested blocks.
	Children []*TestChild
}

// TestChild represents a child class's test data with one or more instances.
type TestChild struct {
	// The child class.
	Class *ClassName
	// The instances to render in test HCL.
	Instances []TestChildInstance
}

func NewClass(className string, ds *DataStore) (*Class, error) {
	genLogger.Tracef("Creating new class struct with class name: %s.", className)

	name, err := NewClassName(className)
	if err != nil {
		return nil, err
	}

	class := Class{
		ClassDefinition: loadClassDefinition(className),
		Name:            name,
		Properties:      make(map[string]*Property),
	}

	genLogger.Tracef("Successfully created new class struct with class name: %s.", className)

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
	genLogger.Debugf("Loading meta file for class '%s'.", c.Name)

	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s.json", constMetaPath, c.Name))
	if err != nil {
		return fmt.Errorf("failed to load meta file for class '%s': %s", c.Name, err.Error())
	}

	genLogger.Tracef("Parsing meta file for class '%s'.", c.Name)
	// For now, the file content is unmarshalled into a map[string]any and then set the class data.
	// This is done because we add logic on top of the file content to set the class data.
	// ENHANCEMENT: investigate if we can unmarshal the file content directly into a class struct specific for meta.
	var metaFileContent map[string]any
	err = json.Unmarshal(fileContent, &metaFileContent)
	if err != nil {
		return fmt.Errorf("failed to parse meta file for class '%s': %s", c.Name, err.Error())
	}

	c.MetaFileContent = metaFileContent[c.Name.MetaStyle()].(map[string]any)

	genLogger.Debugf("Successfully loaded meta file for class '%s'.", c.Name)

	return nil
}

func (c *Class) setClassData(ds *DataStore) error {
	genLogger.Debugf("Setting class data for class '%s'.", c.Name)

	var err error

	c.setAllowDelete()

	err = c.setChildren(ds)
	if err != nil {
		return err
	}

	c.setDeprecated()

	err = c.setDeprecatedVersions()
	if err != nil {
		return err
	}

	c.setHidden()

	err = c.setHiddenVersions()
	if err != nil {
		return err
	}

	c.setIdentifiedBy()

	c.setIsSingleNestedWhenDefinedAsChild()

	err = c.setParents(ds)
	if err != nil {
		return err
	}

	c.setPlatformType()

	err = c.setProperties(ds)
	if err != nil {
		return err
	}

	c.setStateUpgrades(ds)

	c.setPropertyStateUpgradeValues()

	err = c.setRelation()
	if err != nil {
		return err
	}

	c.setRequiredAsChild()

	err = c.setResourceName(ds)
	if err != nil {
		return err
	}

	err = c.setRnFormat()
	if err != nil {
		return err
	}

	err = c.setSupportedVersions()
	if err != nil {
		return err
	}

	genLogger.Debugf("Successfully set class data for class '%s'.", c.Name)
	return nil
}

func (c *Class) setAllowDelete() {
	// Determine if the class can be deleted.
	genLogger.Debugf("Setting AllowDelete for class '%s'.", c.Name)
	if c.ClassDefinition.AllowDelete == "" {
		isCreatableDeletable, ok := c.MetaFileContent["isCreatableDeletable"]
		if ok && isCreatableDeletable.(string) != "never" {
			c.AllowDelete = true
		}
	} else if c.ClassDefinition.AllowDelete != "never" {
		c.AllowDelete = true
	}
	genLogger.Debugf("The AllowDelete property was successfully set to '%t' for the class '%s'.", c.AllowDelete, c.Name)
}

func (c *Class) setChildren(ds *DataStore) error {
	// Determine the child classes for the class.
	genLogger.Debugf("Setting Children for class '%s'.", c.Name)

	// Warn if a class is defined in both IncludeChildren and ExcludeChildren.
	for _, includedChild := range c.ClassDefinition.IncludeChildren {
		if slices.Contains(c.ClassDefinition.ExcludeChildren, includedChild) {
			genLogger.Warnf("Child class '%s' is defined in both IncludeChildren and ExcludeChildren for class '%s'. IncludeChildren takes precedence.", includedChild, c.Name)
		}
	}

	// Initialize with children from ClassDefinition.IncludeChildren (as strings first).
	childClassNames := c.ClassDefinition.IncludeChildren

	// Retrieve the rnMap from MetaFileContent which holds the rnFormat as key and class name as value.
	// The class name is in the format "{ClassNamePackage}:{ClassName}".
	if rnMap, ok := c.MetaFileContent["rnMap"].(map[string]any); ok {
		for rn, classNameInterface := range rnMap {
			// Remove the colon separator from the class name (e.g., "fv:Tenant" -> "fvTenant").
			classNameStr, err := sanitizeClassName(classNameInterface.(string))
			if err != nil {
				return err
			}

			if shouldIncludeChild(rn, classNameStr, c.ClassDefinition.ExcludeChildren, ds.GlobalMetaDefinition.AlwaysIncludeAsChild) {
				childClassNames = append(childClassNames, classNameStr)
			}
		}
	}

	// Sort, deduplicate, and convert to ClassName pointers.
	children, err := sortAndConvertToClassNames(childClassNames)
	if err != nil {
		return err
	}
	c.Children = children

	genLogger.Debugf("Successfully set Children for class '%s'. Found %d children.", c.Name, len(c.Children))
	return nil
}

func (c *Class) setDeprecated() {
	genLogger.Debugf("Setting Deprecated for class '%s'.", c.Name)

	if c.ClassDefinition.Deprecated {
		c.Deprecated = true
	} else if metaDeprecated, ok := c.MetaFileContent["isDeprecated"].(bool); ok {
		c.Deprecated = metaDeprecated
	}

	genLogger.Debugf("Successfully set Deprecated for class '%s'. Deprecated: %t", c.Name, c.Deprecated)
}

func (c *Class) setDeprecatedVersions() error {
	// Determine the deprecated APIC versions for the class from the definition file (override) or meta `deprecatedSince`.
	genLogger.Debugf("Setting DeprecatedVersions for class '%s'.", c.Name)

	deprecatedVersions := c.ClassDefinition.DeprecatedVersions
	if deprecatedVersions == "" {
		deprecatedVersions, _ = c.MetaFileContent["deprecatedSince"].(string)
	}
	if deprecatedVersions == "" {
		genLogger.Debugf("No DeprecatedVersions specified for class '%s'.", c.Name)
		return nil
	}

	parsedVersions, err := NewVersions(deprecatedVersions)
	if err != nil {
		return fmt.Errorf("failed to parse deprecated versions for class '%s': %w", c.Name, err)
	}
	c.DeprecatedVersions = parsedVersions

	genLogger.Debugf("Successfully set DeprecatedVersions for class '%s'. Versions: '%s'", c.Name, c.DeprecatedVersions)
	return nil
}

func (c *Class) setHidden() {
	genLogger.Debugf("Setting Hidden for class '%s'.", c.Name)

	if c.ClassDefinition.Hidden {
		c.Hidden = true
	} else if metaHidden, ok := c.MetaFileContent["isHidden"].(bool); ok {
		c.Hidden = metaHidden
	}

	genLogger.Debugf("Successfully set Hidden for class '%s'. Hidden: %t", c.Name, c.Hidden)
}

func (c *Class) setHiddenVersions() error {
	// Determine the hidden APIC versions for the class from the definition file (override) or meta `hiddenSince`.
	genLogger.Debugf("Setting HiddenVersions for class '%s'.", c.Name)

	hiddenVersions := c.ClassDefinition.HiddenVersions
	if hiddenVersions == "" {
		hiddenVersions, _ = c.MetaFileContent["hiddenSince"].(string)
	}
	if hiddenVersions == "" {
		genLogger.Debugf("No HiddenVersions specified for class '%s'.", c.Name)
		return nil
	}

	parsedVersions, err := NewVersions(hiddenVersions)
	if err != nil {
		return fmt.Errorf("failed to parse hidden versions for class '%s': %w", c.Name, err)
	}
	c.HiddenVersions = parsedVersions

	genLogger.Debugf("Successfully set HiddenVersions for class '%s'. Versions: '%s'", c.Name, c.HiddenVersions)
	return nil
}

func (c *Class) setIdentifiedBy() {
	// Determine the identifying properties of the class.
	genLogger.Debugf("Setting IdentifiedBy for class '%s'.", c.Name)

	// Use ClassDefinition override if specified, otherwise read from meta file.
	if len(c.ClassDefinition.IdentifiedBy) > 0 {
		c.IdentifiedBy = c.ClassDefinition.IdentifiedBy
	} else if identifiedBy, ok := c.MetaFileContent["identifiedBy"].([]any); ok {
		for _, identifier := range identifiedBy {
			if identifierString, ok := identifier.(string); ok && !slices.Contains(c.IdentifiedBy, identifierString) {
				c.IdentifiedBy = append(c.IdentifiedBy, identifierString)
			}
		}
	}

	slices.Sort(c.IdentifiedBy)

	genLogger.Debugf("Successfully set IdentifiedBy for class '%s'. IdentifiedBy: %v", c.Name, c.IdentifiedBy)
}

func (c *Class) setIsSingleNestedWhenDefinedAsChild() {
	// Determine if the class can only be configured once when used as a nested attribute in a parent resource.
	genLogger.Debugf("Setting IsSingleNestedWhenDefinedAsChild for class '%s'.", c.Name)

	c.IsSingleNestedWhenDefinedAsChild = c.ClassDefinition.IsSingleNestedWhenDefinedAsChild || len(c.IdentifiedBy) == 0

	genLogger.Debugf("Successfully set IsSingleNestedWhenDefinedAsChild for class '%s'. IsSingleNestedWhenDefinedAsChild: %t", c.Name, c.IsSingleNestedWhenDefinedAsChild)
}

func (c *Class) setParents(ds *DataStore) error {
	// Determine the parent classes for the class.
	genLogger.Debugf("Setting Parents for class '%s'.", c.Name)

	// Warn if a class is defined in both IncludeParents and ExcludeParents.
	for _, included := range c.ClassDefinition.IncludeParents {
		if slices.Contains(c.ClassDefinition.ExcludeParents, included) {
			genLogger.Warnf("Parent class '%s' is defined in both IncludeParents and ExcludeParents for class '%s'. IncludeParents takes precedence.", included, c.Name)
		}
	}

	// Initialize with parents from ClassDefinition.IncludeParents (as strings first).
	parentClassNames := c.ClassDefinition.IncludeParents

	// Retrieve the containedBy from MetaFileContent which holds the parent class names as keys.
	if containedBy, ok := c.MetaFileContent["containedBy"].(map[string]any); ok {
		for classNameWithColon := range containedBy {
			// Remove the colon separator from the class name (e.g., "fv:AEPg" -> "fvAEPg").
			classNameStr, err := sanitizeClassName(classNameWithColon)
			if err != nil {
				return err
			}

			// Exclude if the class is explicitly excluded via ClassDefinition.ExcludeParents or global ExcludeParents.
			if !slices.Contains(c.ClassDefinition.ExcludeParents, classNameStr) &&
				!slices.Contains(ds.GlobalMetaDefinition.ExcludeParents, classNameStr) {
				parentClassNames = append(parentClassNames, classNameStr)
			}
		}
	}

	// Sort, deduplicate, and convert to ClassName pointers.
	parents, err := sortAndConvertToClassNames(parentClassNames)
	if err != nil {
		return err
	}
	c.Parents = parents

	genLogger.Debugf("Successfully set Parents for class '%s'. Found %d parents.", c.Name, len(c.Parents))
	return nil
}

func (c *Class) setPlatformType() {
	// Determine the platform type of the class.
	genLogger.Debugf("Setting PlatformType for class '%s'.", c.Name)

	// Use ClassDefinition override if specified.
	if c.ClassDefinition.PlatformType != "" {
		switch c.ClassDefinition.PlatformType {
		case "apic":
			c.PlatformType = Apic
		case "cloud":
			c.PlatformType = Cloud
		case "both":
			c.PlatformType = Both
		}
		genLogger.Debugf("Successfully set PlatformType for class '%s' from definition: %s.", c.Name, c.PlatformType)
		return
	}

	// Fall back to meta file platformFlavors.
	if platformFlavors, ok := c.MetaFileContent["platformFlavors"].([]any); ok {
		hasApic := false
		hasCloud := false
		for _, flavor := range platformFlavors {
			if flavorStr, ok := flavor.(string); ok {
				switch flavorStr {
				case "apic":
					hasApic = true
				case "capic":
					hasCloud = true
				default:
					genLogger.Warnf("Unknown platform flavor '%s' found for class '%s'.", flavorStr, c.Name)
				}
			}
		}
		if hasApic && hasCloud {
			c.PlatformType = Both
		} else if hasCloud {
			c.PlatformType = Cloud
		} else if hasApic {
			c.PlatformType = Apic
		}
	}

	genLogger.Debugf("Successfully set PlatformType for class '%s': %s.", c.Name, c.PlatformType)
}

func (c *Class) setProperties(ds *DataStore) error {
	genLogger.Debugf("Setting properties for class '%s'.", c.Name)

	if properties, ok := c.MetaFileContent["properties"]; ok {
		for name, propertyDetails := range properties.(map[string]any) {
			details := propertyDetails.(map[string]any)

			// Look up the property definition override from the class definition file.
			propertyDefinition, hasClassDefinition := c.ClassDefinition.Properties[name]

			// Skip the property entirely when the restriction is Exclude.
			if propertyDefinition.Restriction == Exclude {
				genLogger.Debugf("Property '%s' excluded via definition restriction for class '%s'.", name, c.Name)
				continue
			}

			// Skip globally excluded properties unless the class definition explicitly defines the property.
			if !hasClassDefinition && slices.Contains(ds.GlobalMetaDefinition.ExcludeProperties, name) {
				genLogger.Debugf("Property '%s' excluded via global exclude for class '%s'.", name, c.Name)
				continue
			}

			// Include configurable properties (default behavior from meta file).
			// Include non-configurable properties only when the restriction is ReadOnly.
			if details["isConfigurable"] == true || propertyDefinition.Restriction == ReadOnly {
				property, err := NewProperty(name, details, propertyDefinition, ds.GlobalMetaDefinition)
				if err != nil {
					return err
				}
				c.Properties[name] = property
				c.PropertiesAll = append(c.PropertiesAll, name)

				if property.Required {
					c.PropertiesRequired = append(c.PropertiesRequired, name)
				} else if property.ReadOnly {
					c.PropertiesReadOnly = append(c.PropertiesReadOnly, name)
				} else if property.Optional {
					c.PropertiesOptional = append(c.PropertiesOptional, name)
				}
			}
		}
	}

	// Add synthetic parentDn property when the class has parents.
	// parentDn is not in the meta file; it represents the parent DN in Terraform.
	if len(c.Parents) > 0 {
		parentDnDefinition := c.ClassDefinition.Properties["parentDn"]
		property, err := NewProperty("parentDn", map[string]any{
			"isConfigurable": true,
			"isNaming":       true,
		}, parentDnDefinition, ds.GlobalMetaDefinition)
		if err != nil {
			return err
		}
		property.AttributeName = "parent_dn"
		property.Documentation.Description = "The distinguished name (DN) of the parent object."
		c.Properties["parentDn"] = property
		c.PropertiesAll = append(c.PropertiesAll, "parentDn")
		c.PropertiesRequired = append(c.PropertiesRequired, "parentDn")
	}

	// Sort all property lists alphabetically to ensure deterministic output when iterating over properties in templates.
	slices.Sort(c.PropertiesAll)
	slices.Sort(c.PropertiesOptional)
	slices.Sort(c.PropertiesReadOnly)
	slices.Sort(c.PropertiesRequired)

	genLogger.Debugf("Successfully set properties for class '%s'.", c.Name)
	return nil
}

func (c *Class) setRelation() error {
	// Determine if the class is a relational class.
	genLogger.Debugf("Setting Relation details for class '%s'.", c.Name)

	metaRelationInfo, _ := c.MetaFileContent["relationInfo"].(map[string]any)
	defRelationInfo := c.ClassDefinition.RelationInfo

	// Allow the definition file to opt out of relational handling entirely. This is mutually
	// exclusive with supplying any other `relation_info` override; mixing them is almost
	// always a YAML authoring mistake, so fail fast rather than silently ignoring fields.
	if defRelationInfo.Disabled {
		if defRelationInfo.Type != UndefinedRelationshipType || defRelationInfo.FromClass != "" || len(defRelationInfo.ToClasses) > 0 {
			return fmt.Errorf("failed to set relation for class '%s': relation_info.disabled is mutually exclusive with type, from_class, and to_classes", c.Name)
		}
		genLogger.Debugf("Class '%s' has relation_info.disabled=true; skipping relational handling.", c.Name)
		return nil
	}

	// Skip non-relational classes entirely. A class is relational when either the meta file
	// declares a `relationInfo` block or the definition file supplies a `relation_info` override.
	if metaRelationInfo == nil && defRelationInfo.Type == UndefinedRelationshipType && defRelationInfo.FromClass == "" && len(defRelationInfo.ToClasses) == 0 {
		genLogger.Debugf("Class '%s' is not a relational class; no Relation details to set.", c.Name)
		return nil
	}
	c.Relation.RelationalClass = true

	// Resolve each relationInfo field with the definition taking precedence over meta.
	// defRelationInfo.Type.String() returns "" for the UndefinedRelationshipType zero value,
	// which lets the override-merge fall back to the meta value.
	typeStr := utils.GetValueFromMapWithOverride(metaRelationInfo, "type", defRelationInfo.Type.String())
	fromMo := utils.GetValueFromMapWithOverride(metaRelationInfo, "fromMo", defRelationInfo.FromClass)
	metaToMo := utils.GetValueFromMapWithOverride(metaRelationInfo, "toMo", "")

	// Validate all required relation_info fields together so the user sees every missing
	// field in a single error rather than fix-rebuild-fix-rebuild.
	var missing []string
	if typeStr == "" {
		missing = append(missing, "type")
	}
	if fromMo == "" {
		missing = append(missing, "from_class")
	}
	if metaToMo == "" && len(defRelationInfo.ToClasses) == 0 {
		missing = append(missing, "to_classes")
	}
	if len(missing) > 0 {
		return fmt.Errorf("failed to set relation for class '%s': missing required relation_info fields: %s", c.Name, strings.Join(missing, ", "))
	}

	var relationType RelationshipTypeEnum
	if err := relationType.UnmarshalText([]byte(typeStr)); err != nil {
		return fmt.Errorf("failed to set relation for class '%s': %w", c.Name, err)
	}
	c.Relation.Type = relationType

	fromClass, err := NewClassName(fromMo)
	if err != nil {
		return err
	}
	c.Relation.FromClass = fromClass

	if strings.Contains(c.Name.String(), "To") {
		c.Relation.IncludeFrom = true
	}

	// Resolve the target classes. The `relation_info.to_classes` definition override replaces
	// the meta `toMo` entirely when present, supporting both a single concrete target and a list
	// of concrete targets when the meta `toMo` is abstract.
	toMos := defRelationInfo.ToClasses
	if len(toMos) == 0 {
		toMos = []string{metaToMo}
	}
	toClasses := make([]*ClassName, 0, len(toMos))
	for _, target := range toMos {
		toClass, err := NewClassName(target)
		if err != nil {
			return err
		}
		toClasses = append(toClasses, toClass)
	}
	// Sort and deduplicate by full class name. Pointer comparison is not meaningful here,
	// so SortFunc/CompactFunc are used instead of Sort/Compact.
	slices.SortFunc(toClasses, func(a, b *ClassName) int { return strings.Compare(a.String(), b.String()) })
	c.Relation.ToClasses = slices.CompactFunc(toClasses, func(a, b *ClassName) bool { return a.String() == b.String() })

	genLogger.Debugf("Successfully set Relation details for class '%s'.", c.Name)

	return nil
}

func (c *Class) setRequiredAsChild() {
	// Determine if the class is required when defined as a child in a parent resource.
	genLogger.Debugf("Setting RequiredAsChild for class '%s'.", c.Name)

	c.RequiredAsChild = c.ClassDefinition.RequiredAsChild

	genLogger.Debugf("Successfully set RequiredAsChild for class '%s'. RequiredAsChild: %t", c.Name, c.RequiredAsChild)
}

func (c *Class) setResourceName(ds *DataStore) error {
	genLogger.Debugf("Setting resource name for class '%s'.", c.Name)

	// Get the resource name from definition, else construct the resource name from label, error when both not found.
	if c.ClassDefinition.ResourceName != "" {
		c.ResourceName = c.ClassDefinition.ResourceName
	} else if label, ok := c.MetaFileContent["label"]; ok && label != "" {
		c.ResourceName = utils.Underscore(label.(string))
	} else {
		return fmt.Errorf("failed to set resource name for class '%s': resource_name not defined and label not found", c.Name)
	}
	c.ResourceNameNested = c.ResourceName

	// Determine if the class is relational and set the ResourceName and ResourceNameNested based on the relation.
	if c.Relation.RelationalClass {
		// When the relation has more than one target class, the meta `toMo` is abstract or
		// has been overridden with a list of concrete targets via `relation_info.to_classes`.
		// In that case no single target can drive auto-naming, so the class definition
		// must provide an explicit `resource_name`.
		if len(c.Relation.ToClasses) > 1 {
			if c.ClassDefinition.ResourceName == "" {
				return fmt.Errorf("failed to set resource name for class '%s': resource_name is required when relation_info.to_classes has more than one entry", c.Name)
			}
			// Keep the definition-provided ResourceName and ResourceNameNested as-is.
		} else {
			// Single-target relation: auto-generate `relation_to_<x>` (or
			// `relation_from_<from>_to_<x>`) from the only target class.
			toClass := getRelationshipResourceName(ds, c.Relation.ToClasses[0].String())
			if c.Relation.IncludeFrom {
				c.ResourceName = fmt.Sprintf("relation_from_%s_to_%s", getRelationshipResourceName(ds, c.Relation.FromClass.String()), toClass)
			} else {
				c.ResourceName = fmt.Sprintf("relation_to_%s", toClass)
			}
			c.ResourceNameNested = fmt.Sprintf("relation_to_%s", toClass)
		}
	}

	// Set the plural form for the nested resource name when the class has identifiers.
	if len(c.IdentifiedBy) != 0 {
		c.ResourceNameNested = utils.Plural(c.ResourceNameNested)
	}

	genLogger.Debugf("Successfully set resource name '%s' for class '%s'.", c.ResourceName, c.Name)
	return nil
}

func (c *Class) setRnFormat() error {
	// Determine the relative name (RN) format of the class.
	genLogger.Debugf("Setting RnFormat for class '%s'.", c.Name)

	// Use ClassDefinition override if specified, otherwise read from meta file.
	if c.ClassDefinition.RnFormat != "" {
		c.RnFormat = c.ClassDefinition.RnFormat
	} else if rnFormat, ok := c.MetaFileContent["rnFormat"].(string); ok {
		c.RnFormat = rnFormat
	}

	// When rnFormat is not specified error to force users to add rnFormat.
	if c.RnFormat == "" {
		return fmt.Errorf("rnFormat not specified for class '%s': add rn_format to the class definition file", c.Name)
	}

	// Prepend the RN path prefix if specified in the definition.
	if c.ClassDefinition.RnPrepend != "" {
		c.RnFormat = fmt.Sprintf("%s/%s", c.ClassDefinition.RnPrepend, c.RnFormat)
	}

	genLogger.Debugf("Successfully set RnFormat '%s' for class '%s'.", c.RnFormat, c.Name)
	return nil
}

func (c *Class) setSupportedVersions() error {
	// Determine the supported APIC versions for the class.
	genLogger.Debugf("Setting SupportedVersions for class '%s'.", c.Name)

	// Initialize with versions from ClassDefinition, if not defined set the versions from meta file.
	metaVersions := c.ClassDefinition.SupportedVersions
	if metaVersions == "" {
		metaVersions, _ = c.MetaFileContent["versions"].(string)
	}

	// When versions are not specified error to force users to add versions.
	if metaVersions == "" {
		return fmt.Errorf("versions not specified for class '%s': add versions to the class definition file", c.Name)
	}

	versions, err := NewVersions(metaVersions)
	if err != nil {
		return fmt.Errorf("failed to parse versions for class '%s': %w", c.Name, err)
	}
	c.SupportedVersions = versions

	genLogger.Debugf("Successfully set SupportedVersions for class '%s'. Versions: '%s'", c.Name, c.SupportedVersions)
	return nil
}

func getRelationshipResourceName(ds *DataStore, toClass string) string {
	err := ds.loadClass(toClass)
	if err != nil {
		// If the class is not found, try returning the resource name of the class from global definition file.
		if resourceName, ok := ds.GlobalMetaDefinition.NoMetaFile[toClass]; ok {
			genLogger.Debugf("Failed to load class '%s'. Using resource name from global definition file: %s.", toClass, resourceName)
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

func shouldIncludeChild(rn, className string, excludeChildrenFromClassDef, alwaysIncludeFromGlobalDef []string) bool {
	// Determines if a child class should be included, with a default behavior of excluding a child class.
	// A child class is included if any of the following conditions are met (in order of precedence):
	//  1. The className is NOT in the excludeChildrenFromClassDef list (if it is, the child is excluded regardless of other rules)
	//  2. The className is in the alwaysIncludeFromGlobalDef list
	//  3. The RN format starts with "rs" (relation source classes)
	//  4. The RN format does NOT end with "-" (non-named classes without a specific identifier)
	// Exclude if the class is explicitly excluded via ClassDefinition.ExcludeChildren.
	if slices.Contains(excludeChildrenFromClassDef, className) {
		genLogger.Tracef("Child class '%s' excluded via excludeChildrenFromClassDef.", className)
		return false
	}

	// Include if the class is in GlobalMetaDefinition.AlwaysIncludeAsChild.
	if slices.Contains(alwaysIncludeFromGlobalDef, className) {
		genLogger.Tracef("Child class '%s' included via alwaysIncludeFromGlobalDef.", className)
		return true
	}

	// Include classes where the RN starts with "rs" (relation source).
	if strings.HasPrefix(rn, "rs") {
		genLogger.Tracef("Child class '%s' included because RN starts with 'rs'.", className)
		return true
	}

	// Include classes where the RN does NOT end with "-" (non-named classes without a specific identifier).
	if !strings.HasSuffix(rn, "-") {
		genLogger.Tracef("Child class '%s' included because RN '%s' does not end with '-'.", className, rn)
		return true
	}

	genLogger.Tracef("Child class '%s' excluded by default (RN ends with '-').", className)
	return false
}

// sortAndConvertToClassNames sorts, deduplicates, and converts a slice of class name strings to ClassName pointers.
func sortAndConvertToClassNames(classNameStrings []string) ([]*ClassName, error) {
	slices.Sort(classNameStrings)
	classNameStrings = slices.Compact(classNameStrings)

	classNames := make([]*ClassName, 0, len(classNameStrings))
	for _, classNameStr := range classNameStrings {
		name, err := NewClassName(classNameStr)
		if err != nil {
			return nil, err
		}
		classNames = append(classNames, name)
	}
	return classNames, nil
}

func (c *Class) setTestDependencies(ds *DataStore) {
	// Resolve the test dependencies for the class.
	// By default, explicit definitions are processed first (allowing overrides with ConfigOverrides),
	// then auto-resolution fills in the remainder from Parents and Relation.
	// When ReplaceAutoResolved is true, only explicit dependencies are used.
	genLogger.Debugf("Setting TestDependencies for class '%s'.", c.Name)

	testDependencies := make(map[string]*TestDependency)

	if c.ClassDefinition.TestConfig.ReplaceAutoResolved {
		// Full replacement: skip auto-resolution entirely.
		c.TestDependencies = c.getTestDependenciesFromDefinitions(c.ClassDefinition.TestConfig.Dependencies, ds, testDependencies, 0)
	} else {
		// Process explicit definitions first (allows overriding auto-resolved defaults).
		if len(c.ClassDefinition.TestConfig.Dependencies) > 0 {
			c.TestDependencies = c.getTestDependenciesFromDefinitions(c.ClassDefinition.TestConfig.Dependencies, ds, testDependencies, 0)
		}

		// Auto-resolve remainder from Parents and Relation (skips already-defined references).
		for _, resolvedTestDependency := range slices.Concat(c.resolveParentDependencies(ds, testDependencies), c.resolveTargetDependencies(ds, testDependencies)) {
			if !slices.ContainsFunc(c.TestDependencies, func(td *TestDependency) bool {
				return td.Reference == resolvedTestDependency.Reference
			}) {
				c.TestDependencies = append(c.TestDependencies, resolvedTestDependency)
			}
		}
	}

	// Resolve placeholders in ConfigOverrides against the full DAG.
	c.resolveConfigOverridePlaceholders(testDependencies)
	genLogger.Debugf("Successfully set TestDependencies for class '%s'. Top-level: %d, total in DAG: %d", c.Name, len(c.TestDependencies), len(testDependencies))
}

func (c *Class) getTestDependenciesFromDefinitions(testDependencyDefinitions []TestDependencyDefinition, ds *DataStore, testDependencies map[string]*TestDependency, depth int) []*TestDependency {
	// Convert YAML-driven TestDependencyDefinition entries to TestDependency structs.
	// Validates Role at the given depth (required at depth 0, ignored at depth > 0).
	// Returns a slice rather than assigning directly because this function is called recursively:
	// at depth 0 the result is assigned to c.TestDependencies, at depth > 0 to td.Dependencies.
	genLogger.Tracef("Converting %d test dependency definitions for class '%s' at depth %d.", len(testDependencyDefinitions), c.Name, depth)
	result := make([]*TestDependency, 0, len(testDependencyDefinitions))

	for _, testDependencyDefinition := range testDependencyDefinitions {
		// Validate role based on depth.
		if depth == 0 && testDependencyDefinition.Role == UndefinedRole {
			ds.ctx.Diagnostics.AddError("Class '%s': top-level test dependency for '%s' is missing required 'role' field.", c.Name, testDependencyDefinition.ClassName)
			continue
		}
		if depth > 0 && testDependencyDefinition.Role != UndefinedRole {
			genLogger.Tracef("Class '%s': nested test dependency for '%s' has 'role' field set (role is ignored for nested dependencies).", c.Name, testDependencyDefinition.ClassName)
		}

		// Dedup by reference: reuse existing node when the same reference appears more than once.
		// At depth 0 we still validate the duplicate's Role to allow promoting an existing dep that
		// was first introduced nested (Role=UndefinedRole) into a Parent/Target slot.
		if existing, ok := testDependencies[testDependencyDefinition.Reference]; ok {
			if depth == 0 {
				newRole := testDependencyDefinition.Role
				switch {
				case existing.Role == UndefinedRole && newRole != UndefinedRole:
					existing.Role = newRole
				case existing.Role != UndefinedRole && newRole != UndefinedRole && existing.Role != newRole:
					ds.ctx.Diagnostics.AddError("Class '%s': duplicate test dependency for '%s' declares conflicting roles ('%s' vs '%s'); first declaration wins.", c.Name, testDependencyDefinition.ClassName, existing.Role, newRole)
				}
				if len(testDependencyDefinition.ConfigOverrides) > 0 {
					ds.ctx.Diagnostics.AddError("Class '%s': duplicate test dependency for '%s' carries config_overrides; merge them into the first declaration.", c.Name, testDependencyDefinition.ClassName)
				}
				if len(testDependencyDefinition.Dependencies) > 0 {
					ds.ctx.Diagnostics.AddError("Class '%s': duplicate test dependency for '%s' carries dependencies; merge them into the first declaration.", c.Name, testDependencyDefinition.ClassName)
				}
			} else {
				genLogger.Tracef("Class '%s': nested duplicate reference '%s' reuses existing DAG node.", c.Name, testDependencyDefinition.Reference)
			}
			result = append(result, existing)
			continue
		}

		className, err := NewClassName(testDependencyDefinition.ClassName)
		if err != nil {
			ds.ctx.Diagnostics.AddError("Class '%s': failed to parse dependency class name '%s': %v", c.Name, testDependencyDefinition.ClassName, err)
			continue
		}

		testDependency := &TestDependency{
			Class:           className,
			Reference:       testDependencyDefinition.Reference,
			ReferenceType:   testDependencyDefinition.ReferenceType,
			Role:            testDependencyDefinition.Role,
			ConfigOverrides: testDependencyDefinition.ConfigOverrides,
		}

		testDependencies[testDependencyDefinition.Reference] = testDependency

		// Recursively resolve nested dependencies.
		if len(testDependencyDefinition.Dependencies) > 0 {
			testDependency.Dependencies = c.getTestDependenciesFromDefinitions(testDependencyDefinition.Dependencies, ds, testDependencies, depth+1)
		}

		// Populate per-dependency child overrides (used when this dependency's resource block
		// is rendered in tests). Keys are child class names; values are TestChild instances
		// built from the override, with grandchildren preserved via merge.
		if len(testDependencyDefinition.Children) > 0 {
			testDependency.Children = make(map[string]*TestChild, len(testDependencyDefinition.Children))
			for childClassStr, childOverride := range testDependencyDefinition.Children {
				if len(childOverride.Instances) == 0 {
					continue
				}
				childClassName, err := NewClassName(childClassStr)
				if err != nil {
					ds.ctx.Diagnostics.AddError("Class '%s': dependency '%s' has invalid child class name '%s': %v", c.Name, testDependencyDefinition.Reference, childClassStr, err)
					continue
				}
				var grandChildClass *Class
				if gc, ok := ds.Classes[childClassStr]; ok {
					grandChildClass = &gc
				}
				testDependency.Children[childClassStr] = &TestChild{
					Class:     childClassName,
					Instances: buildOverrideInstances(ds, grandChildClass, childOverride.Instances),
				}
			}
		}

		result = append(result, testDependency)
	}

	genLogger.Tracef("Successfully converted %d test dependency definitions for class '%s' at depth %d.", len(result), c.Name, depth)
	return result
}

func (c *Class) resolveParentDependencies(ds *DataStore, testDependencies map[string]*TestDependency) []*TestDependency {
	// Resolve parent test dependencies.
	// First parent class gets 2 instances (for ForceNew testing), second parent class gets 1.
	genLogger.Tracef("Resolving parent dependencies for class '%s'.", c.Name)
	var result []*TestDependency

	for i, parent := range c.Parents {
		if i >= 2 {
			// Only auto-resolve first 2 parent classes; remaining parents are available via explicit test_config.dependencies if needed.
			genLogger.Tracef("Class '%s': has %d parents, auto-resolving first 2 only.", c.Name, len(c.Parents))
			break
		}

		resourceName := c.getResourceNameForClass(parent.String(), ds)
		if resourceName == "" {
			genLogger.Tracef("Class '%s': parent '%s' not found in DataStore or NoMetaFile, skipping.", c.Name, parent)
			continue
		}

		if i == 0 {
			// First parent: 2 instances for ForceNew testing.
			result = append(result,
				c.buildDependency(parent, fmt.Sprintf("aci_%s.test.id", resourceName), ResourceReference, Parent, ds, testDependencies),
				c.buildDependency(parent, fmt.Sprintf("aci_%s.test_2.id", resourceName), ResourceReference, Parent, ds, testDependencies),
			)
		} else {
			// Additional parent: 1 instance for compatibility testing.
			result = append(result, c.buildDependency(parent, fmt.Sprintf("aci_%s.test.id", resourceName), ResourceReference, Parent, ds, testDependencies))
		}
	}

	genLogger.Tracef("Successfully resolved %d parent dependencies for class '%s'.", len(result), c.Name)
	return result
}

func (c *Class) resolveTargetDependencies(ds *DataStore, testDependencies map[string]*TestDependency) []*TestDependency {
	// Resolve target test dependencies from Relation.ToClasses.
	genLogger.Tracef("Resolving target dependencies for class '%s'.", c.Name)
	if !c.Relation.RelationalClass || len(c.Relation.ToClasses) == 0 {
		genLogger.Tracef("No target dependencies to resolve for class '%s'.", c.Name)
		return nil
	}

	// Multi-target: require explicit YAML.
	if len(c.Relation.ToClasses) > 1 {
		if !slices.ContainsFunc(c.TestDependencies, func(td *TestDependency) bool {
			return td.Role == Target
		}) {
			ds.ctx.Diagnostics.AddError("Class '%s': multi-target relation (%d targets) requires explicit test_config.dependencies with role 'target'.", c.Name, len(c.Relation.ToClasses))
		}
		return nil
	}

	// Single-target: 2 instances for toggling in update tests.
	target := c.Relation.ToClasses[0]
	resourceName := c.getResourceNameForClass(target.String(), ds)
	if resourceName == "" {
		ds.ctx.Diagnostics.AddError("Class '%s': target '%s' has no resource (not in DataStore or NoMetaFile). Provide explicit test_config.dependencies.", c.Name, target)
		return nil
	}

	return []*TestDependency{
		c.buildDependency(target, fmt.Sprintf("aci_%s.test.id", resourceName), ResourceReference, Target, ds, testDependencies),
		c.buildDependency(target, fmt.Sprintf("aci_%s.test_2.id", resourceName), ResourceReference, Target, ds, testDependencies),
	}
}

func (c *Class) buildDependency(className *ClassName, reference string, refType ReferenceTypeEnum, role TestDependencyRoleEnum, ds *DataStore, testDependencies map[string]*TestDependency) *TestDependency {
	// Create a TestDependency node and recursively resolve its parent chain.
	genLogger.Tracef("Building dependency '%s' for class '%s'.", reference, c.Name)
	if existing, ok := testDependencies[reference]; ok {
		genLogger.Tracef("Reusing existing dependency '%s' for class '%s'.", reference, c.Name)
		return existing
	}

	testDependency := &TestDependency{
		Class:         className,
		Reference:     reference,
		ReferenceType: refType,
		Role:          role,
	}
	testDependencies[reference] = testDependency

	// Recursively resolve the dependency's own parents as prerequisites.
	depClass, exists := ds.Classes[className.String()]
	if exists {
		for _, depParent := range depClass.Parents {
			dependencyResourceName := c.getResourceNameForClass(depParent.String(), ds)
			if dependencyResourceName == "" {
				genLogger.Tracef("Class '%s': dependency parent '%s' not found in DataStore or NoMetaFile, skipping.", c.Name, depParent)
				continue
			}
			testDependency.Dependencies = append(testDependency.Dependencies, c.buildDependency(depParent, fmt.Sprintf("aci_%s.test.id", dependencyResourceName), ResourceReference, UndefinedRole, ds, testDependencies))
		}
	}

	genLogger.Tracef("Successfully built dependency '%s' for class '%s'.", reference, c.Name)
	return testDependency
}

func (c *Class) getResourceNameForClass(className string, ds *DataStore) string {
	// Return the resource name for a class, or empty string if not found in DataStore or NoMetaFile.
	if class, ok := ds.Classes[className]; ok {
		return class.ResourceName
	}
	if resourceName, ok := ds.GlobalMetaDefinition.NoMetaFile[className]; ok {
		return resourceName
	}
	return ""
}

func (c *Class) resolveConfigOverridePlaceholders(testDependencies map[string]*TestDependency) {
	// Resolve {{<reference>}} placeholders in ConfigOverrides values for THIS class's
	// TestDependencies only. The testDependencies map is the per-class DAG keyed by
	// class name; it is not shared across classes. Cross-class resolution is not
	// supported here.
	// Unresolved placeholders are left as-is: resolution happens in multiple passes
	// (dependencies → properties → children), so erroring immediately on an unresolved
	// placeholder mid-resolution would be premature. validateTestCompleteness catches
	// everything that is still unresolved after all passes complete.
	genLogger.Tracef("Resolving ConfigOverride placeholders for class '%s'.", c.Name)
	for _, testDependency := range testDependencies {
		if len(testDependency.ConfigOverrides) == 0 {
			continue
		}
		for key, value := range testDependency.ConfigOverrides {
			reference, ok := parsePlaceholder(value)
			if !ok {
				continue
			}
			if resolved, ok := testDependencies[reference]; ok {
				testDependency.ConfigOverrides[key] = resolved.Reference
			}
		}
	}
	genLogger.Tracef("Successfully resolved ConfigOverride placeholders for class '%s'.", c.Name)
}

func (c *Class) setPropertyTestValues(ds *DataStore) {
	// Wire dependency-derived property values (tDn, parent_dn, tn<Cap>Name) and resolve placeholders in TestValues.
	genLogger.Debugf("Resolving property test values for class '%s'.", c.Name)

	// Auto-wire parent_dn from Parent dependencies.
	c.setParentDn()

	// Auto-wire target property for relational classes.
	// Explicit relations expose `tDn` (full DN). Named relations expose `tn<TargetCap>Name`
	// (target's name attribute). Other RelationalClass types are unsupported.
	if c.Relation.RelationalClass {
		switch c.Relation.Type {
		case Explicit:
			c.setTargetDn()
		case Named:
			c.setTargetNameProperty(ds)
		default:
			ds.ctx.Diagnostics.AddError("Class '%s': relational class has unsupported relationship type (expected named or explicit).", c.Name)
		}
	}

	// Resolve placeholders in any explicitly-defined property TestValues.
	c.resolvePlaceholdersInProperties()

	genLogger.Debugf("Successfully resolved property test values for class '%s'.", c.Name)
}

func (c *Class) setParentDn() {
	// Wire the parent_dn property's TestValues from Parent dependencies.
	// parent_dn is a synthetic property that exists only when the class has resolvable parents.
	// Create, Update, and Default use the SAME first parent reference (parent doesn't change).
	// ForceNew uses the second parent reference to trigger destroy+recreate.
	genLogger.Tracef("Setting parentDn test values for class '%s'.", c.Name)

	// Find the parentDn property.
	parentDn, exists := c.Properties["parentDn"]
	if !exists {
		return
	}

	// Only auto-wire if no explicit TestConfig was provided for this property.
	if parentDn.hasTestConfigDefinition() {
		return
	}

	// Collect up to the first two Parent dependencies (mirrors resolveTargetDependencies).
	var parents []*TestDependency
	for _, testDependency := range c.TestDependencies {
		if testDependency.Role != Parent {
			continue
		}
		parents = append(parents, testDependency)
		if len(parents) == 2 {
			break
		}
	}
	if len(parents) == 0 {
		return
	}

	// Wire Create, Update, and Default from the first parent reference.
	parentDn.TestValues = &TestValues{
		Create: []TestValueEntry{{
			ConfigValue:   parents[0].Reference,
			ConfigInclude: true,
			AssertValue:   parents[0].Reference,
			ValueType:     ReferenceValue,
		}},
		Update: []TestValueEntry{{
			ConfigValue:   parents[0].Reference,
			ConfigInclude: true,
			AssertValue:   parents[0].Reference,
			ValueType:     ReferenceValue,
		}},
		Default: []TestValueEntry{{
			ConfigValue:   parents[0].Reference,
			ConfigInclude: true,
			AssertValue:   parents[0].Reference,
			ValueType:     ReferenceValue,
		}},
	}

	// Wire ForceNew from the second parent reference (triggers destroy+recreate).
	if len(parents) > 1 {
		parentDn.TestValues.ForceNew = []TestValueEntry{{
			ConfigValue:   parents[1].Reference,
			ConfigInclude: true,
			AssertValue:   parents[1].Reference,
			ValueType:     ReferenceValue,
		}}
	}
	genLogger.Tracef("Successfully set parentDn test values for class '%s'.", c.Name)
}

func (c *Class) setTargetDn() {
	// Wire the tDn property's TestValues from Target dependencies (Explicit relations).
	// Caller (setPropertyTestValues) already gates on c.Relation.RelationalClass && Explicit.
	genLogger.Tracef("Setting targetDn test values for class '%s'.", c.Name)

	// Find the target dependencies. At least one is required to wire the tDn property.
	var targets []*TestDependency
	for _, testDependency := range c.TestDependencies {
		if testDependency.Role == Target {
			targets = append(targets, testDependency)
		}
	}
	if len(targets) == 0 {
		return
	}

	// Find the tDn property.
	tDn, exists := c.Properties["tDn"]
	if !exists {
		return
	}

	// Only auto-wire if no explicit TestConfig was provided for this property.
	if tDn.hasTestConfigDefinition() {
		return
	}

	// With a single target, every step (Create/Update/Default/ForceNew) uses the same reference.
	// With two or more targets, Create uses the first and Update uses the second so plan diffs
	// exercise both target references.
	createTarget := targets[0]
	updateTarget := targets[0]
	if len(targets) > 1 {
		updateTarget = targets[1]
	} else {
		genLogger.Tracef("Class '%s': only one Target dependency available; Update reuses Create target.", c.Name)
	}

	tDn.TestValues = &TestValues{
		Create: []TestValueEntry{{
			ConfigValue:   createTarget.Reference,
			ConfigInclude: true,
			AssertValue:   createTarget.Reference,
			ValueType:     ReferenceValue,
		}},
		Update: []TestValueEntry{{
			ConfigValue:   updateTarget.Reference,
			ConfigInclude: true,
			AssertValue:   updateTarget.Reference,
			ValueType:     ReferenceValue,
		}},
		Default: []TestValueEntry{{
			ConfigValue:   createTarget.Reference,
			ConfigInclude: true,
			AssertValue:   createTarget.Reference,
			ValueType:     ReferenceValue,
		}},
		ForceNew: []TestValueEntry{{
			ConfigValue:   createTarget.Reference,
			ConfigInclude: true,
			AssertValue:   createTarget.Reference,
			ValueType:     ReferenceValue,
		}},
	}
	genLogger.Tracef("Successfully set targetDn test values for class '%s'.", c.Name)
}

func (c *Class) setTargetNameProperty(ds *DataStore) {
	// Wire the tn<TargetCap>Name property's TestValues from Target dependencies (Named relations).
	// Named relations always have exactly one target class (ToClasses[0]); the property
	// holds the target resource's name attribute rather than its full DN.
	genLogger.Tracef("Setting target-name test values for class '%s'.", c.Name)

	if len(c.Relation.ToClasses) == 0 {
		ds.ctx.Diagnostics.AddError("Class '%s': named relation has no target classes.", c.Name)
		return
	}
	targetClass := c.Relation.ToClasses[0]
	propertyName := "tn" + targetClass.Capitalized() + "Name"

	property, exists := c.Properties[propertyName]
	if !exists {
		// Not every named-relation class exposes the canonical tn<Cap>Name property
		// (some use property overrides). Caller should set TestValues explicitly via test_config.
		genLogger.Tracef("Class '%s': named target property '%s' not found; skipping auto-wire.", c.Name, propertyName)
		return
	}

	// Skip when an explicit TestConfig already drives the property.
	if property.hasTestConfigDefinition() {
		return
	}

	// Find Target dependencies.
	var targets []*TestDependency
	for _, testDependency := range c.TestDependencies {
		if testDependency.Role == Target {
			targets = append(targets, testDependency)
		}
	}
	if len(targets) == 0 {
		return
	}

	// Static references resolve to literal DNs; they have no `.name` attribute to project.
	// Diagnose and bail rather than emit a broken Terraform expression.
	for _, target := range targets {
		if target.ReferenceType == StaticReference {
			ds.ctx.Diagnostics.AddError("Class '%s': named relation target '%s' uses a static reference; property '%s' needs a resource or data_source target (or an explicit test_config).", c.Name, target.Reference, propertyName)
			return
		}
	}

	createTarget := targets[0]
	updateTarget := targets[0]
	if len(targets) > 1 {
		updateTarget = targets[1]
	} else {
		genLogger.Tracef("Class '%s': only one Target dependency available for named relation; Update reuses Create target.", c.Name)
	}

	createRef := targetReferenceToName(createTarget.Reference)
	updateRef := targetReferenceToName(updateTarget.Reference)

	property.TestValues = &TestValues{
		Create: []TestValueEntry{{
			ConfigValue:   createRef,
			ConfigInclude: true,
			AssertValue:   createRef,
			ValueType:     ReferenceValue,
		}},
		Update: []TestValueEntry{{
			ConfigValue:   updateRef,
			ConfigInclude: true,
			AssertValue:   updateRef,
			ValueType:     ReferenceValue,
		}},
		Default: []TestValueEntry{{
			ConfigValue:   createRef,
			ConfigInclude: true,
			AssertValue:   createRef,
			ValueType:     ReferenceValue,
		}},
		ForceNew: []TestValueEntry{{
			ConfigValue:   createRef,
			ConfigInclude: true,
			AssertValue:   createRef,
			ValueType:     ReferenceValue,
		}},
	}
	genLogger.Tracef("Successfully set target-name test values for class '%s'.", c.Name)
}

func targetReferenceToName(ref string) string {
	// Rewrite a target-resource reference from its `.id` form to its `.name` form so the
	// named-relation property points at the target's name attribute rather than its DN.
	if strings.HasSuffix(ref, ".id") {
		return strings.TrimSuffix(ref, ".id") + ".name"
	}
	return ref
}

func (c *Class) resolvePlaceholdersInProperties() {
	// Resolve {{<reference>}} placeholders in property TestValues against the class's TestDependencies.
	genLogger.Tracef("Resolving placeholders in property TestValues for class '%s'.", c.Name)
	for _, property := range c.Properties {
		if property.TestValues == nil {
			continue
		}
		c.resolvePlaceholdersInEntries(property.TestValues.Create)
		c.resolvePlaceholdersInEntries(property.TestValues.Update)
		c.resolvePlaceholdersInEntries(property.TestValues.Default)
		c.resolvePlaceholdersInEntries(property.TestValues.ForceNew)
	}
	genLogger.Tracef("Successfully resolved placeholders in property TestValues for class '%s'.", c.Name)
}

func (c *Class) resolvePlaceholdersInEntries(entries []TestValueEntry) {
	// Resolve placeholders in a slice of TestValueEntry.
	// Unresolved placeholders are left as-is (see resolveConfigOverridePlaceholders for rationale).
	genLogger.Tracef("Resolving placeholders in %d test value entries.", len(entries))
	for i := range entries {
		entry := &entries[i]
		reference, ok := parsePlaceholder(entry.ConfigValue)
		if !ok {
			continue
		}
		resolved := c.findDependencyByRefRecursive(c.TestDependencies, reference)
		if resolved == nil {
			continue
		}
		entry.ConfigValue = resolved.Reference
		entry.AssertValue = resolved.Reference
		entry.ValueType = ReferenceValue
	}
}

func isPlaceholder(value string) bool {
	// Check if a value is a placeholder (e.g., "{{aci_bd.test.id}}").
	return strings.HasPrefix(value, "{{") && strings.HasSuffix(value, "}}")
}

func parsePlaceholder(value string) (string, bool) {
	// Extract the reference from a placeholder string.
	// Returns the trimmed reference and true if the value is a placeholder, or empty and false otherwise.
	if !isPlaceholder(value) {
		return "", false
	}
	return strings.TrimSpace(value[2 : len(value)-2]), true
}

func (c *Class) findDependencyByRefRecursive(testDependencies []*TestDependency, reference string) *TestDependency {
	// Search the entire DAG for a dependency with the given reference.
	for _, testDependency := range testDependencies {
		if testDependency.Reference == reference {
			return testDependency
		}
		if found := c.findDependencyByRefRecursive(testDependency.Dependencies, reference); found != nil {
			return found
		}
	}
	return nil
}

func (c *Class) setChildTestValues(ds *DataStore) {
	// Build TestChildren from child classes' TestValues.
	// Also auto-collects child-driven dependencies into the parent's TestDependencies.
	genLogger.Debugf("Resolving child test values for class '%s'.", c.Name)

	// Track visited classes to prevent infinite recursion from circular child relationships.
	visited := map[string]bool{c.Name.String(): true}

	// Build TestChildren from the class's direct children.
	c.TestChildren = c.buildTestChildren(ds, c.Children, visited)

	// Apply overrides from ClassDefinition.TestConfig.Children.
	if len(c.ClassDefinition.TestConfig.Children) > 0 {
		c.applyChildOverrides(ds)
	}

	// Collect child-driven dependencies.
	c.collectChildDrivenDependencies(ds)

	// Resolve placeholders in child instance properties against parent's TestDependencies.
	c.resolvePlaceholdersInChildren()

	genLogger.Debugf("Successfully resolved child test values for class '%s'. TestChildren count: %d", c.Name, len(c.TestChildren))
}

func (c *Class) buildTestChildren(ds *DataStore, children []*ClassName, visited map[string]bool) []*TestChild {
	// Create TestChild entries for a list of child classes.
	// The visited map tracks the current ancestor chain to prevent infinite
	// recursion from circular child relationships in the ACI meta data. It is
	// NOT a "seen anywhere" set — entries are removed after each child's subtree
	// is built so that siblings (and cousins) can independently include the same
	// child class (e.g. tagAnnotation / tagTag, which apply to almost every class).
	genLogger.Tracef("Building test children for class '%s'. Child count: %d.", c.Name, len(children))
	var result []*TestChild

	for _, childClassName := range children {
		if visited[childClassName.String()] {
			genLogger.Tracef("Class '%s': child class '%s' is an ancestor in the current branch, skipping to prevent cycle.", c.Name, childClassName)
			continue
		}

		childClass, exists := ds.Classes[childClassName.String()]
		if !exists {
			genLogger.Tracef("Class '%s': child class '%s' not found in DataStore, skipping.", c.Name, childClassName)
			continue
		}

		visited[childClassName.String()] = true
		testChild := c.buildTestChild(ds, &childClass, childClassName, visited)
		delete(visited, childClassName.String())
		if testChild != nil {
			result = append(result, testChild)
		}
	}

	genLogger.Tracef("Successfully built %d test children for class '%s'.", len(result), c.Name)
	return result
}

func (c *Class) buildTestChild(ds *DataStore, childClass *Class, childClassName *ClassName, visited map[string]bool) *TestChild {
	// Create a TestChild for a single child class with appropriate instances.
	genLogger.Tracef("Building test child '%s' for class '%s'.", childClassName, c.Name)
	instanceCount := 2
	if childClass.IsSingleNestedWhenDefinedAsChild {
		instanceCount = 1
	}

	testChild := &TestChild{
		Class:     childClassName,
		Instances: make([]TestChildInstance, 0, instanceCount),
	}

	// Instance 0: uses child's TestValues.Create.
	instance0 := buildChildInstance(childClass, true)
	instance0.Children = c.buildTestChildren(ds, childClass.Children, visited)
	testChild.Instances = append(testChild.Instances, instance0)

	// Instance 1 (if list-type): uses child's TestValues.Update.
	if instanceCount > 1 {
		instance1 := buildChildInstance(childClass, false)
		instance1.Children = c.buildTestChildren(ds, childClass.Children, visited)
		testChild.Instances = append(testChild.Instances, instance1)
	}

	genLogger.Tracef("Successfully built test child '%s' for class '%s'.", childClassName, c.Name)
	return testChild
}

func buildChildInstance(childClass *Class, useCreate bool) TestChildInstance {
	// Create a TestChildInstance from a child class's properties.
	// If useCreate is true, uses Create values; otherwise uses Update values (falling back to Create).
	instance := TestChildInstance{
		Properties: make(map[string]TestValueEntry),
	}

	for _, property := range childClass.Properties {
		if property.TestValues == nil || property.IgnoreInTest || property.ReadOnly {
			continue
		}

		if useCreate {
			if len(property.TestValues.Create) > 0 {
				instance.Properties[property.AttributeName] = property.TestValues.Create[0]
			}
			continue
		}

		if len(property.TestValues.Update) > 0 {
			instance.Properties[property.AttributeName] = property.TestValues.Update[0]
		} else if len(property.TestValues.Create) > 0 {
			instance.Properties[property.AttributeName] = property.TestValues.Create[0]
		}
	}

	return instance
}

func (c *Class) applyChildOverrides(ds *DataStore) {
	// Apply explicit overrides from ClassDefinition.TestConfig.Children.
	genLogger.Tracef("Applying child overrides for class '%s'.", c.Name)
	for childClassStr, override := range c.ClassDefinition.TestConfig.Children {
		// Find the matching TestChild.
		var targetChild *TestChild
		for _, testChild := range c.TestChildren {
			if testChild.Class.String() == childClassStr {
				targetChild = testChild
				break
			}
		}
		if targetChild == nil {
			genLogger.Warnf("Class '%s': child override for '%s' does not match any resolved child.", c.Name, childClassStr)
			continue
		}

		// Full replacement semantics: if instances are specified, replace all.
		if len(override.Instances) > 0 {
			var childClass *Class
			if cc, ok := ds.Classes[childClassStr]; ok {
				childClass = &cc
			}
			targetChild.Instances = buildOverrideInstances(ds, childClass, override.Instances)
		}
	}
	genLogger.Tracef("Successfully applied child overrides for class '%s'.", c.Name)
}

func buildOverrideInstances(ds *DataStore, childClass *Class, defs []ChildTestInstanceOverrideDefinition) []TestChildInstance {
	// Convert override instance definitions into TestChildInstance values.
	// For each instance, its Children default to the underlying child class's auto-derived
	// TestChildren (when available) and are then overlaid per-key from the override.
	// Result: grandchildren of classes not mentioned in the override are preserved.
	instances := make([]TestChildInstance, 0, len(defs))
	for _, def := range defs {
		instance := TestChildInstance{
			Properties: make(map[string]TestValueEntry),
		}
		for key, value := range def.Properties {
			entry := TestValueEntry{
				ConfigValue:   value,
				ConfigInclude: true,
				AssertValue:   value,
				ValueType:     StringValue,
			}
			if isPlaceholder(value) {
				entry.ValueType = ReferenceValue
			}
			instance.Properties[key] = entry
		}

		// Start each instance's Children from the underlying child class's auto-derived
		// grandchildren so unrelated child types survive the override.
		var baseChildren []*TestChild
		if childClass != nil {
			baseChildren = childClass.TestChildren
		}
		instance.Children = mergeOverrideChildren(ds, baseChildren, def.Children)
		instances = append(instances, instance)
	}
	return instances
}

func mergeOverrideChildren(ds *DataStore, base []*TestChild, overlay map[string]ChildTestOverrideDefinition) []*TestChild {
	// Per-key merge of an instance's grandchildren overlay onto its auto-derived base.
	// - Base entries whose class is NOT in overlay are kept as-is.
	// - Base entries whose class IS in overlay have their Instances rebuilt from the overlay.
	// - Overlay entries with no matching base entry are appended as override-only TestChildren.
	if len(base) == 0 && len(overlay) == 0 {
		return nil
	}
	result := make([]*TestChild, 0, len(base)+len(overlay))
	matched := make(map[string]bool, len(overlay))
	for _, baseChild := range base {
		classStr := baseChild.Class.String()
		if override, ok := overlay[classStr]; ok && len(override.Instances) > 0 {
			matched[classStr] = true
			var grandChildClass *Class
			if gc, ok := ds.Classes[classStr]; ok {
				grandChildClass = &gc
			}
			result = append(result, &TestChild{
				Class:     baseChild.Class,
				Instances: buildOverrideInstances(ds, grandChildClass, override.Instances),
			})
			continue
		}
		result = append(result, baseChild)
	}
	for classStr, override := range overlay {
		if matched[classStr] {
			continue
		}
		if len(override.Instances) == 0 {
			continue
		}
		className, err := NewClassName(classStr)
		if err != nil {
			genLogger.Warnf("Override-only child class '%s' failed to parse: %v.", classStr, err)
			continue
		}
		grandChildClass := ds.Classes[classStr]
		result = append(result, &TestChild{
			Class:     className,
			Instances: buildOverrideInstances(ds, &grandChildClass, override.Instances),
		})
	}
	return result
}

func (c *Class) collectChildDrivenDependencies(ds *DataStore) {
	// Auto-collect dependencies from child instances that have ReferenceValue properties
	// pointing to resources not already in TestDependencies.
	genLogger.Tracef("Collecting child-driven dependencies for class '%s'.", c.Name)
	c.collectFromTestChildren(ds, c.TestChildren)
	genLogger.Tracef("Successfully collected child-driven dependencies for class '%s'.", c.Name)
}

func (c *Class) collectFromTestChildren(ds *DataStore, testChildren []*TestChild) {
	// Walks a TestChild slice (and recursively each instance's nested Children) collecting
	// child-driven dependencies. Pulled out so grandchild instance properties also contribute.
	for _, testChild := range testChildren {
		childClass, exists := ds.Classes[testChild.Class.String()]
		if !exists {
			continue
		}
		for _, instance := range testChild.Instances {
			for _, entry := range instance.Properties {
				if entry.ValueType != ReferenceValue {
					continue
				}
				// Skip when this reference is already present anywhere in our dependency DAG.
				if c.findDependencyByRefRecursive(c.TestDependencies, entry.ConfigValue) != nil {
					continue
				}
				// Search the child class's full TestDependencies DAG (not just top-level) for the reference.
				if found := c.findDependencyByRefRecursive(childClass.TestDependencies, entry.ConfigValue); found != nil {
					c.TestDependencies = append(c.TestDependencies, found)
				}
			}
			// Recurse into per-instance nested children so grandchild references are collected too.
			if len(instance.Children) > 0 {
				c.collectFromTestChildren(ds, instance.Children)
			}
		}
	}
}

func (c *Class) resolvePlaceholdersInChildren() {
	// Resolve {{<reference>}} placeholders in child instance properties against the parent
	// class's TestDependencies.
	// Called after collectChildDrivenDependencies so the parent has both own and child-collected dependencies.
	genLogger.Tracef("Resolving placeholders in children for class '%s'.", c.Name)
	c.resolvePlaceholdersInTestChildren(c.TestChildren)
	// Also resolve placeholders inside per-dependency child overrides.
	visited := make(map[*TestDependency]bool)
	c.resolvePlaceholdersInDependencyChildren(c.TestDependencies, visited)
	genLogger.Tracef("Successfully resolved placeholders in children for class '%s'.", c.Name)
}

func (c *Class) resolvePlaceholdersInDependencyChildren(deps []*TestDependency, visited map[*TestDependency]bool) {
	// Walks the TestDependency DAG and resolves placeholders inside each dependency's Children
	// overrides. The visited map handles DAG-sharing without re-resolving the same node.
	for _, testDependency := range deps {
		if visited[testDependency] {
			continue
		}
		visited[testDependency] = true
		if len(testDependency.Children) > 0 {
			children := make([]*TestChild, 0, len(testDependency.Children))
			for _, testChild := range testDependency.Children {
				children = append(children, testChild)
			}
			c.resolvePlaceholdersInTestChildren(children)
		}
		c.resolvePlaceholdersInDependencyChildren(testDependency.Dependencies, visited)
	}
}

func (c *Class) resolvePlaceholdersInTestChildren(testChildren []*TestChild) {
	// Recursively resolve placeholders in a slice of TestChild.
	// Unresolved placeholders are left as-is (see resolveConfigOverridePlaceholders for rationale).
	genLogger.Tracef("Resolving placeholders in %d test children.", len(testChildren))
	for _, testChild := range testChildren {
		for i := range testChild.Instances {
			instance := &testChild.Instances[i]
			for key, entry := range instance.Properties {
				reference, ok := parsePlaceholder(entry.ConfigValue)
				if !ok {
					continue
				}
				resolved := c.findDependencyByRefRecursive(c.TestDependencies, reference)
				if resolved == nil {
					continue
				}
				entry.ConfigValue = resolved.Reference
				entry.AssertValue = resolved.Reference
				entry.ValueType = ReferenceValue
				instance.Properties[key] = entry
			}
			// Recurse into grandchildren.
			c.resolvePlaceholdersInTestChildren(instance.Children)
		}
	}
}

func (c *Class) validateTestCompleteness(ctx *Context) {
	// Validate the resolved test data for unresolved placeholders and missing values.
	// Runs after all resolution steps are complete so all errors are reported in a single pass.
	// TODO: Add validation for missing test values, empty TestDependencies when parents exist,
	// properties without any TestValues entries, and other completeness checks.

	// Check ConfigOverrides for unresolved placeholders.
	// validateTestDependencyPlaceholders walks both top-level and nested deps using a
	// shared visited map, so DAG-shared dependencies are diagnosed at most once.
	visited := make(map[*TestDependency]bool)
	c.validateTestDependencyPlaceholders(ctx, c.TestDependencies, visited)

	// Check property TestValues for unresolved placeholders.
	for _, property := range c.Properties {
		if property.TestValues == nil {
			continue
		}
		c.validateEntriesPlaceholders(ctx, property.AttributeName, "Create", property.TestValues.Create)
		c.validateEntriesPlaceholders(ctx, property.AttributeName, "Update", property.TestValues.Update)
		c.validateEntriesPlaceholders(ctx, property.AttributeName, "Default", property.TestValues.Default)
		c.validateEntriesPlaceholders(ctx, property.AttributeName, "ForceNew", property.TestValues.ForceNew)
	}

	// Check child instance properties for unresolved placeholders.
	c.validateChildrenPlaceholders(ctx, c.TestChildren)
}

func (c *Class) validateTestDependencyPlaceholders(ctx *Context, testDependencies []*TestDependency, visited map[*TestDependency]bool) {
	// Recursively check nested dependency ConfigOverrides for unresolved placeholders.
	// The visited map prevents infinite recursion from circular dependency references.
	for _, testDependency := range testDependencies {
		if visited[testDependency] {
			continue
		}
		visited[testDependency] = true
		for key, value := range testDependency.ConfigOverrides {
			if isPlaceholder(value) {
				ctx.Diagnostics.AddError("Class '%s': ConfigOverrides placeholder '%s' (key '%s') on dependency '%s' could not be resolved.", c.Name, value, key, testDependency.Reference)
			}
		}
		// Also validate per-dependency child overrides for unresolved placeholders.
		if len(testDependency.Children) > 0 {
			children := make([]*TestChild, 0, len(testDependency.Children))
			for _, testChild := range testDependency.Children {
				children = append(children, testChild)
			}
			c.validateChildrenPlaceholders(ctx, children)
		}
		c.validateTestDependencyPlaceholders(ctx, testDependency.Dependencies, visited)
	}
}

func (c *Class) validateEntriesPlaceholders(ctx *Context, propName, step string, entries []TestValueEntry) {
	// Check a slice of TestValueEntry for unresolved placeholders.
	for _, entry := range entries {
		if isPlaceholder(entry.ConfigValue) {
			ctx.Diagnostics.AddError("Class '%s': property '%s' %s placeholder '%s' could not be resolved.", c.Name, propName, step, entry.ConfigValue)
		}
	}
}

func (c *Class) validateChildrenPlaceholders(ctx *Context, children []*TestChild) {
	// Recursively check child instance properties for unresolved placeholders.
	for _, testChild := range children {
		for _, instance := range testChild.Instances {
			for key, entry := range instance.Properties {
				if isPlaceholder(entry.ConfigValue) {
					ctx.Diagnostics.AddError("Class '%s': child '%s' property '%s' placeholder '%s' could not be resolved.", c.Name, testChild.Class, key, entry.ConfigValue)
				}
			}
			c.validateChildrenPlaceholders(ctx, instance.Children)
		}
	}
}

func (c *Class) setStateUpgrades(ds *DataStore) {
	// Copy MigrationSource and StateUpgrades from the class definition onto the
	// resolved Class struct and validate the upgrade tree against the resolved
	// properties / children sets. Runs after setProperties so Property name
	// lookups can validate keys. Authoring mistakes accumulate on
	// ds.ctx.Diagnostics for a single-pass summary; this matches the policy used
	// by setTestDependencies and validateTestCompleteness.
	genLogger.Debugf("Setting StateUpgrades for class '%s'.", c.Name)

	c.MigrationSource = c.ClassDefinition.MigrationSource
	c.StateUpgrades = c.ClassDefinition.StateUpgrades

	c.validateStateUpgrades(ds.ctx)

	genLogger.Debugf("Successfully set StateUpgrades for class '%s'. Entries: %d.", c.Name, len(c.StateUpgrades))
}

func (c *Class) validateStateUpgrades(ctx *Context) {
	// Run cross-field and cross-entry validation on the resolved StateUpgrades tree.
	// Strict-decode of the typed enums in definitions.go already rejects unknown
	// enum values at parse time; this method enforces the rules that cannot be
	// expressed via the YAML schema alone.
	//
	// All issues are written to ctx.Diagnostics rather than returned, and no check
	// short-circuits, so a single pass surfaces every authoring mistake.

	// Migration-source coherence: a non-zero MigrationSource requires at least
	// one state_upgrades entry. The specific prior_schema_version is whatever
	// the prior provider's resource was at when migrated (commonly 0 for a
	// resource that never had a schema bump in SDKv2, but can be any value if
	// the SDKv2 resource itself went through one or more upgrades).
	if c.MigrationSource != UndefinedMigrationSource && len(c.StateUpgrades) == 0 {
		ctx.Diagnostics.AddError("Class '%s': migration_source %q requires at least one state_upgrades entry describing the migration hop from the prior provider", c.Name, c.MigrationSource)
	}

	// Unique prior_schema_version values.
	seenVersion := make(map[int]struct{}, len(c.StateUpgrades))
	for _, entry := range c.StateUpgrades {
		if _, dup := seenVersion[entry.PriorSchemaVersion]; dup {
			ctx.Diagnostics.AddError("Class '%s': duplicate prior_schema_version %d", c.Name, entry.PriorSchemaVersion)
		}
		seenVersion[entry.PriorSchemaVersion] = struct{}{}

		c.validateStateUpgradeEntry(ctx, entry, fmt.Sprintf("Class '%s': prior_schema_version %d", c.Name, entry.PriorSchemaVersion))
	}

	// Exhaustiveness on top-level Attributes / Children: keys must resolve to a
	// current Property / child class on this class (or be marked "removed", in
	// which case the key is allowed to be a meta-only / excluded name).
	for _, entry := range c.StateUpgrades {
		for propertyName, attributeUpgrade := range entry.Attributes {
			if _, ok := c.Properties[propertyName]; ok {
				continue
			}
			if attributeUpgrade.LegacyStatus == Removed {
				continue
			}
			ctx.Diagnostics.AddError("Class '%s': prior_schema_version %d: attribute %q not found in resolved properties", c.Name, entry.PriorSchemaVersion, propertyName)
		}
		for childClassName, childUpgrade := range entry.Children {
			if c.hasChild(childClassName) {
				continue
			}
			if childUpgrade.LegacyStatus == Removed {
				continue
			}
			ctx.Diagnostics.AddError("Class '%s': prior_schema_version %d: child %q not found in resolved children", c.Name, entry.PriorSchemaVersion, childClassName)
		}
	}
}

func (c *Class) hasChild(childClassName string) bool {
	for _, child := range c.Children {
		if child.String() == childClassName {
			return true
		}
	}
	return false
}

func (c *Class) validateStateUpgradeEntry(ctx *Context, entry StateUpgradeDefinition, prefix string) {
	// Per-entry checks: collision detection across legacy_attribute values within
	// the same prior-version entry, plus structural shape rules per node.
	// The prefix already contains the class name and prior_schema_version, so
	// per-node diagnostics only need to append the bucket and key.
	legacyNames := make(map[string]struct{})
	for propertyName, attributeUpgrade := range entry.Attributes {
		attributeUpgrade.validate(ctx, fmt.Sprintf("%s: attributes[%q]", prefix, propertyName))
		if attributeUpgrade.LegacyAttribute != "" {
			if _, dup := legacyNames[attributeUpgrade.LegacyAttribute]; dup {
				ctx.Diagnostics.AddError("%s: duplicate legacy_attribute %q within the same prior_schema_version entry", prefix, attributeUpgrade.LegacyAttribute)
			}
			legacyNames[attributeUpgrade.LegacyAttribute] = struct{}{}
		}
	}
	for childClassName, childUpgrade := range entry.Children {
		childUpgrade.validateChild(ctx, fmt.Sprintf("%s: children[%q]", prefix, childClassName))
	}
}

func (c *Class) setPropertyStateUpgradeValues() {
	// Distribute the class-level StateUpgrades tree into per-Property maps
	// so renderer templates iterating one property at a time can look up
	// that property's prior-schema shape per version without a back-pointer
	// to the owning class.
	//
	// Walk: the outer loop is one entry per prior_schema_version; the inner
	// Attributes map's keys are CURRENT property names, which is the join
	// key against c.Properties. Each (property, version) pair lands at
	// property.StateUpgradeValues[version], so a property referenced in N
	// version entries accumulates N entries. The map is lazy-initialised,
	// so properties no upgrade ever touches keep StateUpgradeValues == nil.
	//
	// Attributes whose key does not resolve to a current Property are
	// silently skipped — by design, validateStateUpgrades permits this
	// when LegacyStatus == Removed. The underlying meta property may
	// still be alive but intentionally excluded from the current schema
	// (a future dual-expose + plan-modify target), or it may be truly
	// gone; either way there is no schema attribute on c.Properties to
	// wire StateUpgradeValues onto here.
	//
	// Class-side ownership is required because Property has no parentClass
	// reference and the per-property setter has no way to locate the
	// class's upgrade tree.
	//
	// Inner Children entries are NOT distributed here: child Property maps
	// would need cross-class lookups before all classes are loaded into the
	// DataStore. Templates that need inner-child state-upgrade data walk
	// the full StateUpgrades tree on the owning class directly.
	genLogger.Debugf("Distributing StateUpgradeValues for class '%s'.", c.Name)
	for _, stateUpgradeEntry := range c.StateUpgrades {
		for propertyName, attributeUpgradeDefinition := range stateUpgradeEntry.Attributes {
			property, ok := c.Properties[propertyName]
			if !ok {
				// No current schema attribute under this key — by design
				// when LegacyStatus == Removed. The underlying meta
				// property may still be alive but intentionally excluded
				// from the current schema (a future dual-expose +
				// plan-modify target), or it may be truly gone. Nothing
				// to distribute either way.
				continue
			}
			if property.StateUpgradeValues == nil {
				property.StateUpgradeValues = make(map[int]StateUpgradeValue)
			}
			property.StateUpgradeValues[stateUpgradeEntry.PriorSchemaVersion] = buildStateUpgradeValue(property, attributeUpgradeDefinition)
		}
	}
	genLogger.Debugf("Successfully distributed StateUpgradeValues for class '%s'.", c.Name)
}

func buildStateUpgradeValue(property *Property, attributeUpgradeDefinition AttributeUpgradeDefinition) StateUpgradeValue {
	// Build the prior-schema view of one Property at one version: seed from
	// the current property's resolved attributes (AttributeName, Type, and
	// the Required/Optional/Computed triplet), then overlay any explicit
	// legacy_* overrides from the upgrade node. Unset overrides (empty
	// string for LegacyAttribute, zero enum for LegacyType and
	// LegacyRestriction) leave the seeded value in place.
	//
	// LegacyType is mapped through legacyTypeToValueType to collapse the
	// framework-attribute vocabulary into the renderer's ValueTypeEnum.
	// LegacyRestriction rewrites the entire Required/Optional/Computed
	// triplet rather than merging field by field, because the three flags
	// are mutually constrained.
	genLogger.Tracef("Building state upgrade value for property '%s'.", property.PropertyName)
	stateUpgradeValue := StateUpgradeValue{
		AttributeName: property.AttributeName,
		Required:      property.Required,
		Optional:      property.Optional,
		Computed:      property.Computed,
		Type:          property.ValueType,
		Status:        attributeUpgradeDefinition.LegacyStatus,
	}
	if attributeUpgradeDefinition.LegacyAttribute != "" {
		stateUpgradeValue.AttributeName = attributeUpgradeDefinition.LegacyAttribute
	}
	if attributeUpgradeDefinition.LegacyType != UndefinedLegacyAttributeType {
		stateUpgradeValue.Type = legacyTypeToValueType(attributeUpgradeDefinition.LegacyType)
	}
	switch attributeUpgradeDefinition.LegacyRestriction {
	case Required:
		stateUpgradeValue.Required, stateUpgradeValue.Optional, stateUpgradeValue.Computed = true, false, false
	case Optional:
		stateUpgradeValue.Required, stateUpgradeValue.Optional, stateUpgradeValue.Computed = false, true, true
	case ReadOnly:
		stateUpgradeValue.Required, stateUpgradeValue.Optional, stateUpgradeValue.Computed = false, false, true
	}
	genLogger.Tracef("Successfully built state upgrade value for property '%s'.", property.PropertyName)
	return stateUpgradeValue
}
