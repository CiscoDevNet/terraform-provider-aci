package data

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type GlobalMetaDefinition struct {
	// A list of class names that should always be included as children regardless of standard inclusion logic.
	AlwaysIncludeAsChild []string `yaml:"always_include_as_child"`
	// A map of meta property names to their desired attribute names.
	// Used to globally override the default snake_case derivation for specific properties (e.g., "descr" → "description").
	// Per-class PropertyDefinition.AttributeName takes precedence over this.
	AttributeNameOverrides map[string]string `yaml:"attribute_name_overrides"`
	// A list of property names to exclude from all classes.
	// A class-level PropertyDefinition entry for the same property takes precedence over this global exclude.
	ExcludeProperties []string `yaml:"exclude_properties"`
	// A map containing class names as keys and their corresponding resource names as values.
	// This is used to search for the resource name of a class when it is not defined in meta directory.
	NoMetaFile map[string]string `yaml:"no_meta_file"`
	// A map of word substitutions applied when humanizing a snake_case resource name into a documentation label.
	// e.g., "Bgp" → "BGP", "External Network Instance Profile" → "External EPG".
	// Multi-word keys are matched as substrings; single-word keys are only replaced on whole-word matches.
	DocumentationLabelOverrides map[string]string `yaml:"documentation_label_overrides"`
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

type ClassDocumentationDefinition struct {
	// Overrides the humanized documentation label for this class (e.g., "Application EPG").
	// Used to build "Manages ACI <Label>" / "Data source for ACI <Label>" and any other
	// documentation text that references the class by its label.
	// When empty, the label is derived from class.ResourceName via humanization + DocumentationLabelOverrides.
	Label string `yaml:"label"`
	// Additional sentence(s) appended after the standard description prefix.
	// Shared and applied to both the resource and datasource documentation.
	// Per-artifact entries in Resource.Description / Datasource.Description are appended after this shared text.
	Description string `yaml:"description"`
	// Overrides the DN format strings sourced from the meta file. When set, these values are
	// used verbatim. Sorting and the constMaxDnFormatsToDisplay cap still apply.
	DnFormats []string `yaml:"dn_formats"`
	// A list of child class names to exclude from the documentation children list.
	ExcludeChildren []string `yaml:"exclude_children"`
	// A list of child class names to force-include in the documentation children list.
	IncludeChildren []string `yaml:"include_children"`
	// Notes rendered with -> prefix in the documentation.
	// These are shared and applied to both the resource and datasource documentation.
	// Per-artifact entries in Resource.Notes / Datasource.Notes are appended after this shared list.
	Notes []string `yaml:"notes"`
	// Sub-category for Terraform Registry sidebar grouping.
	SubCategory string `yaml:"sub_category"`
	// GUI locations in APIC (e.g., "Tenants -> Networking -> VRFs").
	UiLocations []string `yaml:"ui_locations"`
	// Warnings rendered with !> prefix in the documentation.
	// These are shared and applied to both the resource and datasource documentation.
	// Per-artifact entries in Resource.Warnings / Datasource.Warnings are appended after this shared list.
	Warnings []string `yaml:"warnings"`
	// Resource-only notes/warnings appended to the shared lists for the resource documentation.
	Resource ArtifactDocumentationDefinition `yaml:"resource"`
	// Datasource-only notes/warnings appended to the shared lists for the datasource documentation.
	Datasource ArtifactDocumentationDefinition `yaml:"datasource"`
}

// ArtifactDocumentationDefinition holds documentation overrides specific to a single
// generated artifact (resource or datasource). Description/notes/warnings entries here
// are appended to the shared ClassDocumentationDefinition lists.
type ArtifactDocumentationDefinition struct {
	// Additional sentence(s) appended after the shared description text for this artifact.
	Description string `yaml:"description"`
	// Notes rendered with -> prefix in the documentation, appended after the shared notes.
	Notes []string `yaml:"notes"`
	// Warnings rendered with !> prefix in the documentation, appended after the shared warnings.
	Warnings []string `yaml:"warnings"`
}

type PropertyDefinition struct {
	// Overrides the default attribute name derived from the meta property name in snake_case notation.
	AttributeName string `yaml:"attribute_name"`
	// Indicates that the property is deprecated. A deprecation warning will be included in the schemas.
	// When true, this overrides the meta `isDeprecated` flag with logical OR semantics: definition can flip true on top of meta but cannot force-off.
	Deprecated bool `yaml:"deprecated"`
	// The deprecated APIC versions for the property. Format: "1.0(1e)-" or "4.2(7f)-4.2(7w),5.2(1g)-".
	// Used to indicate versions where the property is deprecated but still functional.
	// When non-empty, this overrides the meta `deprecatedSince` value.
	DeprecatedVersions string `yaml:"deprecated_versions"`
	// Indicates that the property is hidden by the APIC API (no longer accepted).
	// When true, this overrides the meta `isHidden` flag with logical OR semantics.
	Hidden bool `yaml:"hidden"`
	// The hidden APIC versions for the property. Format: "1.0(1e)-" or "4.2(7f)-4.2(7w),5.2(1g)-".
	// When non-empty, this overrides the meta `hiddenSince` value.
	HiddenVersions string `yaml:"hidden_versions"`
	// Controls the schema behavior of the property.
	// Valid values: "required", "optional", "read_only", "exclude".
	// When empty, the default behavior is derived from the meta file (isConfigurable+isNaming → required, isConfigurable → optional).
	Restriction string `yaml:"restriction"`
	// Overrides the sensitive flag for the property.
	// When true, the property is marked as sensitive in the Terraform schema regardless of the meta file.
	Sensitive bool `yaml:"sensitive"`
	// Overrides the versions from the meta file. Format: "1.0(1e)-" or "4.2(7f)-4.2(7w),5.2(1g)-".
	SupportedVersions string `yaml:"supported_versions"`
	// Overrides the meta `validators` array. When non-empty, replaces the meta validators entirely.
	Validators []ValidatorDefinition `yaml:"validators"`
	// Adds extra valid values to the property on top of the meta `validValues` array.
	// Each entry is treated as both the wire value and the localName.
	// A warning is logged when an entry is already present in the meta valid values.
	AddValidValues []string `yaml:"add_valid_values"`
	// Removes valid values from the meta `validValues` array by localName.
	// A warning is logged when an entry is not present in the meta valid values.
	RemoveValidValues []string `yaml:"remove_valid_values"`
	// Overrides the value type derived from the meta `uitype`.
	// Accepted values mirror the ValueTypeEnum vocabulary: "string", "set", "ip_address", "semantic_equality".
	// An error is returned during generation when set to an unrecognized value.
	ValueType string `yaml:"value_type"`
}

// ValidatorDefinition mirrors the YAML/JSON shape of a single validator entry (min/max plus optional regex statements).
type ValidatorDefinition struct {
	Min       int64                      `yaml:"min"`
	Max       int64                      `yaml:"max"`
	RegexList []RegexStatementDefinition `yaml:"regexs"`
}

// RegexStatementDefinition mirrors the YAML/JSON shape of a single regex entry inside a validator.
type RegexStatementDefinition struct {
	Regex string `yaml:"regex"`
	Type  string `yaml:"type"`
}

// RelationInfoDefinition mirrors the meta `relationInfo` block and allows per-field overrides
// from the class definition file. Empty fields fall back to the corresponding meta value.
type RelationInfoDefinition struct {
	// When true, the class is treated as non-relational regardless of the meta `relationInfo`
	// block. Used to opt out of relational handling for classes whose meta declares a relation
	// that should not be exposed by the provider. Mutually exclusive with the override fields
	// below; an error is returned during generation when `Disabled` is combined with any of
	// `Type`, `FromClass`, or `ToClasses`.
	Disabled bool `yaml:"disabled"`
	// Relationship type. Valid values: "named", "explicit".
	Type string `yaml:"type"`
	// Source class of the relation in `pkg:Class` form (e.g., "fv:EPg").
	FromClass string `yaml:"from_class"`
	// Target classes of the relation in `pkg:Class` form (e.g., ["vz:BrCP"]).
	// Replaces the meta `toMo` entirely when set; a single-element list is the common case,
	// while a multi-element list is required when the meta `toMo` is an abstract superclass
	// (e.g., "infra:DomP") that maps to multiple concrete target classes. When more than one
	// class is listed, the class definition must also provide an explicit `resource_name`
	// since auto-naming from a single target is no longer meaningful.
	ToClasses []string `yaml:"to_classes"`
}

type ClassDefinition struct {
	// Overrides the default deletion behavior from meta file. Set to "never" to prevent deletion of the class.
	// The value "never" is used to keep the input consistent with the meta data file.
	AllowDelete string `yaml:"allow_delete"`
	// Indicates that the resource and datasource are deprecated. A deprecation warning will be included in the schemas.
	// When true, this overrides the meta `isDeprecated` flag with logical OR semantics: definition can flip true on top of meta but cannot force-off.
	Deprecated bool `yaml:"deprecated"`
	// The deprecated APIC versions for the class. Format: "1.0(1e)-" or "4.2(7f)-4.2(7w),5.2(1g)-".
	// Used to indicate versions where the class is deprecated but still functional.
	// When non-empty, this overrides the meta `deprecatedSince` value.
	DeprecatedVersions string `yaml:"deprecated_versions"`
	// Indicates that the class is hidden by the APIC API (no longer accepted).
	// When true, this overrides the meta `isHidden` flag with logical OR semantics.
	Hidden bool `yaml:"hidden"`
	// The hidden APIC versions for the class. Format: "1.0(1e)-" or "4.2(7f)-4.2(7w),5.2(1g)-".
	// When non-empty, this overrides the meta `hiddenSince` value.
	HiddenVersions string `yaml:"hidden_versions"`
	// Documentation specific overrides for the class.
	Documentation ClassDocumentationDefinition `yaml:"documentation"`
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
	// Overrides the platform type from the meta file. Valid values: "apic", "cloud", "both".
	PlatformType string `yaml:"platform_type"`
	// Property-level overrides keyed by the meta property name (e.g., "pcTag", "name").
	// Used to override the attribute name, or control the schema restriction (required, optional, read_only, exclude).
	Properties map[string]PropertyDefinition `yaml:"properties"`
	// Overrides (or supplies) the meta `relationInfo` block on a per-field basis.
	// Any non-empty field replaces the matching field from the meta file; empty fields fall back to meta.
	// When the meta file has no `relationInfo` and this definition supplies at least one field,
	// the class is treated as a relational class driven entirely by the definition.
	// `to_classes` directly maps to `Relation.ToClasses`; when more than one class is listed,
	// `resource_name` must also be set explicitly since auto-naming from a single target is no longer meaningful.
	// Set `disabled: true` to opt out of relational handling entirely, even when the meta declares a relation.
	RelationInfo RelationInfoDefinition `yaml:"relation_info"`
	// Indicates that the class is required when defined as a child in a parent resource.
	RequiredAsChild bool `yaml:"required_as_child"`
	// Overrides the resource name derived from the meta file label (e.g., "vrf" instead of "context").
	ResourceName string `yaml:"resource_name"`
	// Overrides the rnFormat from the meta file. The full RN format string (e.g., "custom-{name}").
	RnFormat string `yaml:"rn_format"`
	// Prepends a path prefix to the resolved RN format (e.g., "infra" results in "infra/{rnFormat}").
	RnPrepend string `yaml:"rn_prepend"`
	// Overrides the versions from the meta file. Format: "1.0(1e)-" or "4.2(7f)-4.2(7w),5.2(1g)-".
	SupportedVersions string `yaml:"supported_versions"`
}

func loadClassDefinition(className string) ClassDefinition {
	classDefinitionPath := fmt.Sprintf("%s/%s.yaml", constDefinitionsPath, className)
	var classDefinitionData ClassDefinition

	classDefinitionBytes, err := os.ReadFile(classDefinitionPath)
	if err != nil {
		genLogger.Debugf("The file '%s' was not found in the definitions folder.", classDefinitionPath)
		return classDefinitionData
	}

	err = yaml.Unmarshal(classDefinitionBytes, &classDefinitionData)
	if err != nil {
		genLogger.Fatal(err.Error())
	}

	return classDefinitionData
}
