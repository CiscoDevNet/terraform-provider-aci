package data

import "fmt"

// SubCategoryEnum represents the sub-category for Terraform Registry sidebar grouping.
type SubCategoryEnum string

const (
	SubCategoryAAA                 SubCategoryEnum = "AAA"
	SubCategoryAccessPolicies      SubCategoryEnum = "Access Policies"
	SubCategoryApplicationMgmt     SubCategoryEnum = "Application Management"
	SubCategoryCloud               SubCategoryEnum = "Cloud"
	SubCategoryContract            SubCategoryEnum = "Contract"
	SubCategoryFabricAccessPol     SubCategoryEnum = "Fabric Access Policies"
	SubCategoryFabricInventory     SubCategoryEnum = "Fabric Inventory"
	SubCategoryFabricPolicies      SubCategoryEnum = "Fabric Policies"
	SubCategoryFirmware            SubCategoryEnum = "Firmware"
	SubCategoryGeneric             SubCategoryEnum = "Generic"
	SubCategoryImportExport        SubCategoryEnum = "Import/Export"
	SubCategoryL2Out               SubCategoryEnum = "L2Out"
	SubCategoryL3Out               SubCategoryEnum = "L3Out"
	SubCategoryL4L7Services        SubCategoryEnum = "L4-L7"
	SubCategoryMonitoring          SubCategoryEnum = "Monitoring"
	SubCategoryMultiSite           SubCategoryEnum = "Multi-Site"
	SubCategoryNetworking          SubCategoryEnum = "Networking"
	SubCategoryNodeMgmt            SubCategoryEnum = "Node Management"
	SubCategoryScheduler           SubCategoryEnum = "Scheduler"
	SubCategorySystemSettings      SubCategoryEnum = "System Settings"
	SubCategoryTenantInfraPolicies SubCategoryEnum = "Tenant Infra Policies"
	SubCategoryTenantPolicies      SubCategoryEnum = "Tenant Policies"
	SubCategoryVirtualNetworking   SubCategoryEnum = "Virtual Networking"
)

// validateSubCategory validates that the given SubCategoryEnum value is one of the allowed sub-categories.
func validateSubCategory(class *Class) error {
	switch class.ClassDefinition.Documentation.SubCategory {
	case "":
		if class.IsSingleNestedWhenDefinedAsChild {
			return nil
		}
		return fmt.Errorf("sub_category not specified for class '%s': add documentation.sub_category to the class definition file", class.Name.full)
	case SubCategoryAAA, SubCategoryAccessPolicies, SubCategoryApplicationMgmt,
		SubCategoryCloud, SubCategoryContract, SubCategoryFabricAccessPol,
		SubCategoryFabricInventory, SubCategoryFabricPolicies, SubCategoryFirmware,
		SubCategoryGeneric, SubCategoryImportExport, SubCategoryL2Out,
		SubCategoryL3Out, SubCategoryL4L7Services, SubCategoryMonitoring,
		SubCategoryMultiSite, SubCategoryNetworking, SubCategoryNodeMgmt,
		SubCategoryScheduler, SubCategorySystemSettings, SubCategoryTenantInfraPolicies,
		SubCategoryTenantPolicies, SubCategoryVirtualNetworking:
		return nil
	default:
		return fmt.Errorf("invalid sub_category '%s'", class.ClassDefinition.Documentation.SubCategory)
	}
}

type ClassDocumentation struct {
	// The class name for referencing in documentation.
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
	// Notes rendered with -> prefix in docs.
	Notes []string
	// Sub-category for Terraform Registry sidebar grouping.
	SubCategory SubCategoryEnum
	// GUI locations in APIC (e.g., "Tenants -> Networking -> VRFs").
	UiLocations []string
	// Warnings rendered with !> prefix in docs.
	Warnings []string
}

func (c *Class) setDocumentation() error {
	genLogger.Debug(fmt.Sprintf("Setting Documentation for class '%s'.", c.Name))

	c.Documentation.setClassName(c)

	c.Documentation.setChildren(c)

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
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation ClassName for class '%s'.", class.Name.full))
}

func (d *ClassDocumentation) setChildren(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation Children for class '%s'.", class.Name.full))
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation Children for class '%s'.", class.Name.full))
}

func (d *ClassDocumentation) setDeprecationWarning(class *Class) {
	genLogger.Debug(fmt.Sprintf("Setting Documentation DeprecationWarning for class '%s'.", class.Name.full))
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation DeprecationWarning for class '%s'.", class.Name.full))
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
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation Notes for class '%s'.", class.Name.full))
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
	genLogger.Debug(fmt.Sprintf("Successfully set Documentation Warnings for class '%s'.", class.Name.full))
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
