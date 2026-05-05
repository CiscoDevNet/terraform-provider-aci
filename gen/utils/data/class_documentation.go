package data

import (
	"fmt"
	"slices"
)

// validSubCategories contains the allowed sub-category values for Terraform Registry sidebar grouping.
var validSubCategories = []string{
	"AAA",
	"Access Policies",
	"Application Management",
	"Cloud",
	"Contract",
	"Fabric Access Policies",
	"Fabric Inventory",
	"Fabric Policies",
	"Firmware",
	"Generic",
	"Import/Export",
	"L2Out",
	"L3Out",
	"L4-L7",
	"Monitoring",
	"Multi-Site",
	"Networking",
	"Node Management",
	"Scheduler",
	"System Settings",
	"Tenant Infra Policies",
	"Tenant Policies",
	"Virtual Networking",
}

// validateSubCategory validates that the given sub-category value is one of the allowed sub-categories.
func validateSubCategory(class *Class) error {
	if class.ClassDefinition.Documentation.SubCategory == "" {
		if class.IsSingleNestedWhenDefinedAsChild {
			return nil
		}
		return fmt.Errorf("sub_category not specified for class '%s': add documentation.sub_category to the class definition file", class.Name.full)
	}
	if !slices.Contains(validSubCategories, class.ClassDefinition.Documentation.SubCategory) {
		return fmt.Errorf("invalid sub_category '%s'", class.ClassDefinition.Documentation.SubCategory)
	}
	return nil
}

type ClassDocumentation struct {
	// A markdown link to the DevNet documentation page for the class, used to reference the class in the documentation.
	// Format: "[<className>](https://<host>/app/index.html#/objects/<className>/overview)".
	ClassName string
	// Child classes that have their own standalone resource but are not inline children.
	Children []string
	// A deprecation warning message for the documentation when the class is deprecated.
	DeprecationWarning string
	// The description of the class, used at the top of the documentation.
	// From meta comment[] or definition override.
	Description string
	// The description used when this class appears as a nested child in a parent resource.
	DescriptionWhenDefinedAsChild string
	// DN format strings from meta file (e.g., "uni/tn-{name}").
	DnFormats []string
	// A migration warning message for the documentation when the class has been migrated from a previous provider version.
	MigrationWarning string
	// The supported APIC versions string for display in documentation.
	SupportedVersions string
	// Notes rendered with -> prefix in the resource documentation.
	// Resolved as the shared ClassDocumentationDefinition.Notes followed by ClassDocumentationDefinition.Resource.Notes.
	ResourceNotes []string
	// Warnings rendered with !> prefix in the resource documentation.
	// Resolved as the shared ClassDocumentationDefinition.Warnings followed by ClassDocumentationDefinition.Resource.Warnings.
	ResourceWarnings []string
	// Notes rendered with -> prefix in the datasource documentation.
	// Resolved as the shared ClassDocumentationDefinition.Notes followed by ClassDocumentationDefinition.Datasource.Notes.
	DatasourceNotes []string
	// Warnings rendered with !> prefix in the datasource documentation.
	// Resolved as the shared ClassDocumentationDefinition.Warnings followed by ClassDocumentationDefinition.Datasource.Warnings.
	DatasourceWarnings []string
	// Sub-category for Terraform Registry sidebar grouping.
	SubCategory string
	// GUI locations in APIC (e.g., "Tenants -> Networking -> VRFs").
	UiLocations []string
}

func (c *Class) setDocumentation(ds *DataStore) error {
	genLogger.Debug(fmt.Sprintf("Setting Documentation for class '%s'.", c.Name))

	c.Documentation.setClassName(c)

	c.Documentation.setChildren(c, ds)

	c.Documentation.setDeprecationWarning(c)

	c.Documentation.setDescription(c)

	c.Documentation.setDescriptionWhenDefinedAsChild(c)

	c.Documentation.setDnFormats(c)

	c.Documentation.setMigrationWarning(c)

	c.Documentation.setSupportedVersions(c)

	c.Documentation.setNotes(c)

	err := c.Documentation.setSubCategory(c)
	if err != nil {
		return err
	}

	err = c.Documentation.setUiLocations(c)
	if err != nil {
		return err
	}

	c.Documentation.setWarnings(c)

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation for class '%s'.", c.Name))
	return nil
}

func (d *ClassDocumentation) setClassName(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation ClassName for class '%s'.", class.Name.full))

	d.ClassName = fmt.Sprintf("[%s](https://%s/app/index.html#/objects/%s/overview)", class.Name.full, constPubhubDevnetHost, class.Name.full)

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation ClassName for class '%s'. ClassName: %s", class.Name.full, d.ClassName))
}

func (d *ClassDocumentation) setChildren(class *Class, ds *DataStore) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation Children for class '%s'.", class.Name.full))

	rnMap, ok := class.MetaFileContent["rnMap"].(map[string]interface{})
	if !ok {
		genLogger.Debug(fmt.Sprintf("No rnMap available for class '%s'; skipping documentation children.", class.Name.full))
		return
	}

	// Build a set of children already embedded as nested attributes in this resource
	// so they can be excluded from the documentation children list.
	childrenIncludedInResource := make(map[string]struct{}, len(class.Children))
	for _, child := range class.Children {
		childrenIncludedInResource[child.full] = struct{}{}
	}

	links := make([]string, 0)
	for _, classNameInterface := range rnMap {
		childName, err := sanitizeClassName(classNameInterface.(string))
		if err != nil {
			genLogger.Warn(fmt.Sprintf("Skipping invalid child class name in rnMap for class '%s': %s", class.Name.full, err))
			continue
		}
		if _, isIncluded := childrenIncludedInResource[childName]; isIncluded {
			continue
		}
		childClass, ok := ds.Classes[childName]
		if !ok || childClass.ResourceName == "" {
			continue
		}
		links = append(links, fmt.Sprintf("[%s_%s](%s/%s)", constProviderName, childClass.ResourceName, constRegistryResourceBaseUrl, childClass.ResourceName))
	}

	slices.Sort(links)
	d.Children = slices.Compact(links)

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation Children for class '%s'. Children: %v", class.Name.full, d.Children))
}

func (d *ClassDocumentation) setDeprecationWarning(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation DeprecationWarning for class '%s'.", class.Name.full))

	if class.Deprecated {
		d.DeprecationWarning = fmt.Sprintf("The %s class is deprecated and will be removed in a future release.", d.ClassName)
	}

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation DeprecationWarning for class '%s'. DeprecationWarning: %s", class.Name.full, d.DeprecationWarning))
}

func (d *ClassDocumentation) setDescription(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation Description for class '%s'.", class.Name.full))
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation Description for class '%s'.", class.Name.full))
}

func (d *ClassDocumentation) setDescriptionWhenDefinedAsChild(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation DescriptionWhenDefinedAsChild for class '%s'.", class.Name.full))
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation DescriptionWhenDefinedAsChild for class '%s'.", class.Name.full))
}

func (d *ClassDocumentation) setDnFormats(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation DnFormats for class '%s'.", class.Name.full))
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation DnFormats for class '%s'.", class.Name.full))
}

func (d *ClassDocumentation) setMigrationWarning(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation MigrationWarning for class '%s'.", class.Name.full))
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation MigrationWarning for class '%s'.", class.Name.full))
}

func (d *ClassDocumentation) setSupportedVersions(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation SupportedVersions for class '%s'.", class.Name.full))
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation SupportedVersions for class '%s'.", class.Name.full))
}

func (d *ClassDocumentation) setNotes(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation Notes for class '%s'.", class.Name.full))

	docDef := class.ClassDefinition.Documentation
	d.ResourceNotes = slices.Concat(docDef.Notes, docDef.Resource.Notes)
	d.DatasourceNotes = slices.Concat(docDef.Notes, docDef.Datasource.Notes)

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation Notes for class '%s'. ResourceNotes: %v, DatasourceNotes: %v", class.Name.full, d.ResourceNotes, d.DatasourceNotes))
}

func (d *ClassDocumentation) setUiLocations(class *Class) error {
	genLogger.Debug(fmt.Sprintf("Setting Documentation UiLocations for class '%s'.", class.Name.full))

	if len(class.ClassDefinition.Documentation.UiLocations) == 0 {
		if class.IsSingleNestedWhenDefinedAsChild {
			return nil
		}
		return fmt.Errorf("class '%s': ui_locations not specified: add documentation.ui_locations to the class definition file", class.Name.full)
	}

	d.UiLocations = class.ClassDefinition.Documentation.UiLocations

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation UiLocations for class '%s'. UiLocations: %v", class.Name.full, d.UiLocations))
	return nil
}

func (d *ClassDocumentation) setWarnings(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation Warnings for class '%s'.", class.Name.full))

	docDef := class.ClassDefinition.Documentation
	d.ResourceWarnings = slices.Concat(docDef.Warnings, docDef.Resource.Warnings)
	d.DatasourceWarnings = slices.Concat(docDef.Warnings, docDef.Datasource.Warnings)

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation Warnings for class '%s'. ResourceWarnings: %v, DatasourceWarnings: %v", class.Name.full, d.ResourceWarnings, d.DatasourceWarnings))
}

func (d *ClassDocumentation) setSubCategory(class *Class) error {
	genLogger.Debug(fmt.Sprintf("Setting Documentation SubCategory for class '%s'.", class.Name.full))

	if err := validateSubCategory(class); err != nil {
		return fmt.Errorf("class '%s': %w", class.Name.full, err)
	}

	d.SubCategory = class.ClassDefinition.Documentation.SubCategory

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation SubCategory '%s' for class '%s'.", d.SubCategory, class.Name.full))
	return nil
}
