package data

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	// The human-readable documentation label derived from class.ResourceName, with
	// DocumentationLabelOverrides applied (e.g., "Application EPG", "BGP Timers").
	// Resolved as ClassDocumentationDefinition.Label when set, otherwise humanized from class.ResourceName.
	Label string
	// The description line for the resource documentation (e.g., "Manages ACI Application EPG").
	ResourceDescription string
	// The description line for the datasource documentation (e.g., "Data source for ACI Application EPG").
	DatasourceDescription string
	// The description used when this class appears as a nested child in a parent resource.
	DescriptionWhenDefinedAsChild string
	// DN format strings from meta file (e.g., "uni/tn-{name}").
	DnFormats []string
	// Parent DN references rendered in the documentation. Built from class.Parents.
	// Resources (parent classes that have a Terraform resource) are listed first in
	// alphabetical order, followed by parent classes without a resource under an
	// explanatory note. Each section is independently capped by constMaxParentDnsToDisplay;
	// a notice line replaces that section's list when exceeded.
	ParentDns []string
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

	c.Documentation.setLabel(c, ds)

	c.Documentation.setDescription(c)

	c.Documentation.setDescriptionWhenDefinedAsChild(c)

	c.Documentation.setDnFormats(c)

	c.Documentation.setParentDns(c, ds)

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

func (d *ClassDocumentation) setLabel(class *Class, ds *DataStore) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation Label for class '%s'.", class.Name.full))

	if class.ClassDefinition.Documentation.Label != "" {
		d.Label = class.ClassDefinition.Documentation.Label
	} else {
		d.Label = cases.Title(language.English).String(strings.ReplaceAll(class.ResourceName, "_", " "))
		// Apply word substitutions: multi-word keys are matched as substrings;
		// single-word keys are only replaced on whole-word matches to avoid partial-word collisions.
		for key, replacement := range ds.GlobalMetaDefinition.DocumentationLabelOverrides {
			if strings.Contains(key, " ") {
				d.Label = strings.ReplaceAll(d.Label, key, replacement)
			} else if slices.Contains(strings.Split(d.Label, " "), key) {
				d.Label = strings.ReplaceAll(d.Label, key, replacement)
			}
		}
	}

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation Label for class '%s'. Label: %s", class.Name.full, d.Label))
}

func (d *ClassDocumentation) setDescription(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation Description for class '%s'.", class.Name.full))

	docDef := class.ClassDefinition.Documentation

	// Build the resource and datasource description lines. The standard prefix sentence is always
	// applied; shared and artifact-specific text is appended afterwards, separated by a single space.
	resourceParts := []string{fmt.Sprintf("Manages ACI %s.", d.Label)}
	if docDef.Description != "" {
		resourceParts = append(resourceParts, docDef.Description)
	}
	if docDef.Resource.Description != "" {
		resourceParts = append(resourceParts, docDef.Resource.Description)
	}
	d.ResourceDescription = strings.Join(resourceParts, " ")

	datasourceParts := []string{fmt.Sprintf("Data source for ACI %s.", d.Label)}
	if docDef.Description != "" {
		datasourceParts = append(datasourceParts, docDef.Description)
	}
	if docDef.Datasource.Description != "" {
		datasourceParts = append(datasourceParts, docDef.Datasource.Description)
	}
	d.DatasourceDescription = strings.Join(datasourceParts, " ")

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation Description for class '%s'. ResourceDescription: %s, DatasourceDescription: %s", class.Name.full, d.ResourceDescription, d.DatasourceDescription))
}

func (d *ClassDocumentation) setDescriptionWhenDefinedAsChild(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation DescriptionWhenDefinedAsChild for class '%s'.", class.Name.full))
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation DescriptionWhenDefinedAsChild for class '%s'.", class.Name.full))
}

func (d *ClassDocumentation) setDnFormats(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation DnFormats for class '%s'.", class.Name.full))

	if override := class.ClassDefinition.Documentation.DnFormats; len(override) > 0 {
		d.DnFormats = override
	} else if rawFormats, ok := class.MetaFileContent["dnFormats"].([]interface{}); ok {
		for _, raw := range rawFormats {
			if s, ok := raw.(string); ok {
				d.DnFormats = append(d.DnFormats, s)
			}
		}
	} else {
		genLogger.Debug(fmt.Sprintf("No dnFormats available for class '%s'.", class.Name.full))
	}

	// Sort to guarantee stable, deterministic output across regenerations.
	slices.Sort(d.DnFormats)

	if len(d.DnFormats) > constMaxDnFormatsToDisplay {
		notice := fmt.Sprintf("Too many DN formats to display, see model documentation for all possible parents of %s.", d.ClassName)
		d.DnFormats = append([]string{notice}, d.DnFormats[:constMaxDnFormatsToDisplay]...)
	}

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation DnFormats for class '%s'. DnFormats: %v", class.Name.full, d.DnFormats))
}

func (d *ClassDocumentation) setMigrationWarning(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation MigrationWarning for class '%s'.", class.Name.full))

	if class.IsMigration {
		d.MigrationWarning = "This resource has been migrated to the terraform plugin protocol version 6, refer to the [migration guide](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/guides/migration) for more details and implications for already managed resources."
	}

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation MigrationWarning for class '%s'. MigrationWarning: %s", class.Name.full, d.MigrationWarning))
}

func (d *ClassDocumentation) setParentDns(class *Class, ds *DataStore) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation ParentDns for class '%s'.", class.Name.full))

	var resourceEntries, classOnlyEntries []string
	for _, parent := range class.Parents {
		parentLink := fmt.Sprintf("[%s](https://%s/app/index.html#/objects/%s/overview)", parent.full, constPubhubDevnetHost, parent.full)
		parentClass, knownInStore := ds.Classes[parent.full]
		if knownInStore && parentClass.ResourceName != "" {
			resourceEntries = append(resourceEntries, fmt.Sprintf("[%s_%s](%s/%s) (%s)", constProviderName, parentClass.ResourceName, constRegistryResourceBaseUrl, parentClass.ResourceName, parentLink))
		} else {
			classOnlyEntries = append(classOnlyEntries, parentLink)
		}
	}

	// class.Parents is already sorted by sortAndConvertToClassNames, so the partitioned
	// slices preserve that order; no extra sort needed.

	if len(resourceEntries) > constMaxParentDnsToDisplay {
		d.ParentDns = []string{fmt.Sprintf("Too many parent DNs to display, see model documentation for all possible parents of %s.", d.ClassName)}
	} else {
		d.ParentDns = resourceEntries
	}

	if len(classOnlyEntries) > 0 {
		d.ParentDns = append(d.ParentDns, "The distinguished name (DN) of classes below can be used but currently there is no available resource for it:")
		if len(classOnlyEntries) > constMaxParentDnsToDisplay {
			d.ParentDns = append(d.ParentDns, fmt.Sprintf("Too many classes to display, see model documentation for all possible classes of %s.", d.ClassName))
		} else {
			d.ParentDns = append(d.ParentDns, classOnlyEntries...)
		}
	}

	genLogger.Debug(fmt.Sprintf("Successfully set Documentation ParentDns for class '%s'. ParentDns: %v", class.Name.full, d.ParentDns))
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
