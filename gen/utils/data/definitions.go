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
	// A list of class names to globally exclude from parent resolution (root-level singletons like polUni, fabricInst).
	ExcludeParents []string `yaml:"exclude_parents"`
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
	// Per-meta-property documentation overrides applied as a global layer between
	// any per-class PropertyDefinition.Documentation.Description and the meta
	// `comment` / `label` fallbacks (consumed by setDescription).
	// Keys are meta property names (e.g., "descr", "nameAlias"); values are the
	// English sentence rendered as the property description. When the value
	// contains `%s`, the renderer interpolates the class's humanised resource
	// name (same substitution as GetResourceNameAsDescription).
	// Normalises inconsistent meta `comment` wording across the ~180 classes
	// without duplicating the same text into every per-class file.
	PropertyDocumentationOverrides map[string]string `yaml:"property_documentation_overrides"`
}

func loadGlobalMetaDefinition() GlobalMetaDefinition {
	definition, err := os.ReadFile(constGlobalDefinitionFilePath)
	if err != nil {
		genLogger.Fatal("A file 'global.yaml' is required to be defined in the definitions folder.")
	}

	definitionGlobalMetaData, err := parseGlobalMetaDefinition(definition)
	if err != nil {
		genLogger.Fatal(err.Error())
	}

	return definitionGlobalMetaData
}

// parseGlobalMetaDefinition decodes raw YAML bytes into a GlobalMetaDefinition.
// UnmarshalStrict rejects unknown YAML keys so renamed/typo'd fields surface
// as a generator error instead of being silently ignored.
func parseGlobalMetaDefinition(data []byte) (GlobalMetaDefinition, error) {
	var definitionGlobalMetaData GlobalMetaDefinition
	err := yaml.UnmarshalStrict(data, &definitionGlobalMetaData)
	return definitionGlobalMetaData, err
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
	// Curated subset of `containedBy` parent class names used to render example
	// HCL snippets (one block per entry) in the resource/datasource docs,
	// `resource_example.tf.tmpl`, `datasource_example.tf.tmpl`,
	// `resource_example_all_attributes.tf.tmpl`, `resource.md.tmpl`, and
	// `testvars.yaml.tmpl`. Used when meta `containedBy` is too large (e.g.
	// relation/tag classes) to render every parent without producing dozens of
	// near-identical snippets. Each entry must be a meta class name; the
	// renderer resolves it to the generated resource name via `getResourceName`.
	ExampleParentClasses []string `yaml:"example_parent_classes"`
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
	// When the field is omitted or empty, the default behavior is derived from the meta file
	// (isConfigurable+isNaming → required, isConfigurable → optional).
	Restriction RestrictionEnum `yaml:"restriction"`
	// Overrides the RequiresReplace behavior for the property.
	// When non-nil, takes precedence over the meta-derived `isNaming` logic.
	// Use true to force replacement on change; false to suppress even for naming properties.
	RequiresReplace *bool `yaml:"requires_replace"`
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
	// Accepted values mirror the ValueTypeEnum vocabulary: "string", "set", "object", "ip_address", "semantic_equality".
	// An error is returned by YAML decoding when set to an unrecognized value.
	ValueType ValueTypeEnum `yaml:"value_type"`
	// Per-property documentation overrides (description, notes, warnings).
	// Reuses the shared ArtifactDocumentationDefinition struct.
	Documentation ArtifactDocumentationDefinition `yaml:"documentation"`
	// Overrides the default values derived from the meta `default` field.
	// Each key is the default value and the map value is an optional version range string (e.g., "5.0(1a)-").
	// An empty version string means the default applies to all versions.
	// When non-empty, completely replaces the meta-derived defaults.
	DefaultValues map[string]string `yaml:"default_values"`
	// Test configuration for the property. Controls test value generation and inclusion/exclusion.
	TestConfig TestConfigDefinition `yaml:"test_config"`
}

// TestValueEntryDefinition is the YAML representation of a single test value entry.
type TestValueEntryDefinition struct {
	// The value to write into HCL configuration.
	ConfigValue string `yaml:"config_value"`
	// Whether to include this value in the test config. Pointer: nil = default true, explicit false = omit.
	ConfigInclude *bool `yaml:"config_include"`
	// Expected value in state after apply. Empty = same as ConfigValue.
	AssertValue string `yaml:"assert_value"`
	// Controls HCL rendering: "string" (default, quoted), "reference" (unquoted expression).
	ValueType ValueRenderTypeEnum `yaml:"value_type"`
}

// TestConfigDefinition groups test-related overrides for a property definition.
type TestConfigDefinition struct {
	// Values for the "all attributes" create step. Each entry becomes a TestValueEntry.
	Create []TestValueEntryDefinition `yaml:"create"`
	// Values for the "required-only" step. Typically auto-derived; explicit overrides here.
	Default []TestValueEntryDefinition `yaml:"default"`
	// Values for the update step. Each entry becomes a TestValueEntry.
	Update []TestValueEntryDefinition `yaml:"update"`
	// Values for the ForceNew step. Typically auto-derived (same as Create for non-parent_dn).
	ForceNew []TestValueEntryDefinition `yaml:"force_new"`
	// Values for the Legacy step, exercising state_upgrades legacy attribute
	// aliases (Functioning / Frozen). Independent of the standard buckets:
	// supplying legacy alone is permitted, and supplying legacy is required
	// when the legacy alias has a different Terraform type than the current
	// attribute (auto-derivation skips that case with a warning).
	Legacy []TestValueEntryDefinition `yaml:"legacy"`
	// When true, the property is excluded from generated tests entirely.
	IgnoreInTest bool `yaml:"ignore_in_test"`
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
	// Relationship type. Valid values: "named", "explicit". Zero value (UndefinedRelationshipType)
	// means "no override" — fall back to meta `relationInfo.type`.
	Type RelationshipTypeEnum `yaml:"type"`
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
	// Selects which generated artifacts the renderer emits for this class.
	// A nil slice (the YAML field omitted entirely) signals the resolver to
	// auto-derive: classes with non-empty `IdentifiedBy` get [resource, datasource],
	// classes with empty `IdentifiedBy` get nothing (the legacy default).
	// A non-nil but empty slice (`artifacts: []`) is an explicit opt-out and
	// removes the class from both `provider.Resources()` and
	// `provider.DataSources()`. A non-empty list overrides the auto-derivation:
	// `[resource, datasource]` opts an empty-`IdentifiedBy` class in as both;
	// `[datasource]` or `[resource]` selects a single artifact (e.g. `topSystem`
	// renders as a datasource only).
	Artifacts []ArtifactEnum `yaml:"artifacts"`
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
	// Records the lineage of this resource: the prior provider or generator it was
	// migrated from. Today only "from_sdkv2" is recognized. Drives the documentation
	// migration warning. A non-zero value also requires at least one state_upgrades
	// entry describing the migration hop from the prior provider; the specific
	// prior_schema_version depends on what the prior resource was at when migrated
	// (validated in class.setStateUpgrades).
	MigrationSource MigrationSourceEnum `yaml:"migration_source"`
	// Per-version state upgrade definitions. Each entry describes the prior-schema
	// shape used by Terraform's UpgradeResourceState RPC and the current-schema
	// exposure (or non-exposure) of any legacy attribute alias. Independent
	// direct-to-current upgraders per the framework contract: the framework selects
	// the upgrader matching the saved version and does not chain entries.
	StateUpgrades []StateUpgradeDefinition `yaml:"state_upgrades"`
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
	// Alternate parent-DN placements for classes that legitimately resolve under
	// more than one user-facing parent and reach APIC via different request paths
	// per placement (e.g. pkiKeyRing as system-scoped vs tenant-scoped). The
	// generated resource branches on the user's `parent_dn` and selects the
	// matching variant's (api_endpoint, json_envelope) pair. Each entry pins one
	// (parent_class, rn_prepend, wrapper_class, test_platform) tuple; not
	// representable in meta `containedBy`, which lists direct parents without the
	// implicit-wrapper or user-facing-parent annotations.
	ParentDnVariants []ParentDnVariantDefinition `yaml:"parent_dn_variants"`
	// Overrides the versions from the meta file. Format: "1.0(1e)-" or "4.2(7f)-4.2(7w),5.2(1g)-".
	SupportedVersions string `yaml:"supported_versions"`
	// Test configuration for the class. Controls dependency resolution and child test value overrides.
	TestConfig ClassTestConfigDefinition `yaml:"test_config"`
}

// TestDependencyDefinition is used at ALL levels of test_config.dependencies (top-level and nested).
// Role is required at the top level (tells resource-under-test how to consume this dep).
// Role must be empty at nested levels (pure prerequisites). Validated based on depth.
type TestDependencyDefinition struct {
	// The class name of the dependency resource. Always required.
	ClassName string `yaml:"class_name"`
	// The reference value: either a static DN string or a Terraform resource/datasource attribute path.
	Reference string `yaml:"reference"`
	// How to interpret the Reference field. Valid values: "static", "resource", "data_source".
	// When omitted, defaults to ResourceReference (the iota zero of ReferenceTypeEnum).
	ReferenceType ReferenceTypeEnum `yaml:"reference_type"`
	// Role of this dependency. Valid values: "parent", "target". Required at top level, empty for nested.
	Role TestDependencyRoleEnum `yaml:"role"`
	// Recursive dependencies: resources that THIS dependency needs to exist first.
	Dependencies []TestDependencyDefinition `yaml:"dependencies"`
	// Optional property overrides for the dependency resource's HCL configuration.
	ConfigOverrides map[string]string `yaml:"config_overrides"`
	// Children of THIS dependency resource (i.e. nested blocks within the dependency's HCL),
	// keyed by child class name. Overrides auto-derived child test values for the dependency.
	Children map[string]ChildTestOverrideDefinition `yaml:"children"`
}

// ChildTestOverrideDefinition is the YAML representation of child test value overrides.
// Keyed by child class name in the parent map. No instance_count — count is determined
// by the child class's own IsSingleNestedWhenDefinedAsChild setting.
type ChildTestOverrideDefinition struct {
	// Full replacement: when present, ALL auto-derived instances are discarded and replaced by these.
	Instances []ChildTestInstanceOverrideDefinition `yaml:"instances"`
}

// ChildTestInstanceOverrideDefinition represents a single instance override with its properties and nested children.
type ChildTestInstanceOverrideDefinition struct {
	// Property overrides for this instance, keyed by attribute name.
	Properties map[string]string `yaml:"properties"`
	// Grandchildren of THIS override instance (i.e. nested blocks within this child instance),
	// keyed by grandchild class name. Recursively overrides auto-derived nested children.
	Children map[string]ChildTestOverrideDefinition `yaml:"children"`
}

// ClassTestConfigDefinition groups all test-related configuration for a class.
// Consistent with PropertyDefinition.TestConfig — both use the `test_config` YAML key.
type ClassTestConfigDefinition struct {
	// When true, Dependencies fully replaces all auto-resolved dependencies.
	// When false (default), Dependencies are merged on top of auto-resolved dependencies.
	ReplaceAutoResolved bool `yaml:"replace_auto_resolved"`
	// Test dependencies for the class. By default additive on top of auto-resolved dependencies.
	// Set replace_auto_resolved to true to skip auto-resolution entirely.
	Dependencies []TestDependencyDefinition `yaml:"dependencies"`
	// Children of THIS class (the one being generated), keyed by child class name.
	// When an entry's `instances` is set, it fully replaces the auto-derived instances
	// for that child class; unspecified child classes keep their auto-derived values.
	Children map[string]ChildTestOverrideDefinition `yaml:"children"`
	// Suppresses generation of specific test buckets. `child` skips the entry in
	// every parent's testvars.yaml iteration (the class still emits its own
	// resource / datasource tests). `resource` skips this class's own
	// resource_aci_<x>_test.go. `datasource` skips this class's own
	// data_source_aci_<x>_test.go. Combinations are valid (e.g. [child, resource]).
	// A nil or empty slice means no skips.
	IgnoreTests []IgnoreTestEnum `yaml:"ignore_tests"`
	// Suppresses just the ImportStateVerify assertion inside the import test (the
	// import smoke test still runs). Required for classes whose APIC response
	// carries non-roundtrip state that would fail attribute-equality verification.
	IgnoreImportStateVerify bool `yaml:"ignore_import_state_verify"`
}

// StateUpgradeDefinition describes a single prior-schema version's state upgrade.
// Each entry is direct-to-current: the Terraform plugin framework selects the
// upgrader matching the saved state's version and does not chain intermediate
// upgraders. Keys under Attributes are meta PropertyName values (or known
// synthetic names like "parentDn"/"tDn"); keys under Children are meta child
// class names. The two splits keep scalar/leaf attributes separate from nested
// blocks at every level.
type StateUpgradeDefinition struct {
	// The prior schema version this upgrader handles (the value saved in state).
	// Must be unique across all StateUpgrades entries for a class.
	PriorSchemaVersion int `yaml:"prior_schema_version"`
	// Top-level scalar attributes that changed between the prior schema and the
	// current schema. Map keys are meta PropertyName values; map values describe
	// how the prior attribute differed from the current attribute and what to do
	// with its legacy name in the current schema.
	Attributes map[string]AttributeUpgradeDefinition `yaml:"attributes"`
	// Top-level nested blocks (child classes) that changed between the prior
	// schema and the current schema. Map keys are meta child class names; values
	// recursively describe the block-level diff and any inner attribute changes.
	Children map[string]AttributeUpgradeDefinition `yaml:"children"`
}

// AttributeUpgradeDefinition is the unified node type used for both scalar
// attributes and nested blocks inside a state_upgrades entry. Recursive: an
// attributes/children entry can itself contain attributes and children to
// describe inner shape changes at any depth.
type AttributeUpgradeDefinition struct {
	// The Terraform attribute name as it appeared in the prior schema. Omit when
	// the prior name matches the current name. On a children entry, presence
	// records a block rename or a scalar-wrap (when the new block didn't exist
	// in the prior schema, an inner attribute carries the prior flat scalar
	// name instead).
	LegacyAttribute string `yaml:"legacy_attribute"`
	// The Terraform plugin framework attribute type in the prior schema. Omit
	// when the prior type matches the current type. Required on a "removed"
	// node since there is no current attribute to inherit from.
	LegacyType LegacyAttributeTypeEnum `yaml:"legacy_type"`
	// The schema restriction in the prior schema (required/optional/read_only).
	// Omit to inherit from the current property (UndefinedRestriction). Required
	// on a "removed" node since there is no current attribute to inherit from.
	LegacyRestriction RestrictionEnum `yaml:"legacy_restriction"`
	// Controls how this legacy attribute is exposed in the CURRENT schema:
	//   - functioning (default): legacy name still exposed, full device round-trip.
	//   - frozen: legacy name still exposed but no device round-trip.
	//   - removed: legacy name dropped from current schema; entry retained for
	//     state migration only.
	LegacyStatus LegacyStatusEnum `yaml:"legacy_status"`
	// Inner scalar attribute changes for a nested block. Same key/value
	// semantics as StateUpgradeDefinition.Attributes.
	Attributes map[string]AttributeUpgradeDefinition `yaml:"attributes"`
	// Inner child class changes for a nested block. Same key/value semantics
	// as StateUpgradeDefinition.Children.
	Children map[string]AttributeUpgradeDefinition `yaml:"children"`
}

// validate enforces the shape rules for an attributes-bucket node. Errors are
// accumulated on ctx.Diagnostics rather than returned so a single pass surfaces
// every authoring mistake in one summary. The path argument carries the fully
// qualified location of this node (e.g. "Class 'fvCtx': prior_schema_version 0:
// attributes[\"name\"]") so each diagnostic is self-describing.
func (n AttributeUpgradeDefinition) validate(ctx *Context, path string) {
	if n.LegacyStatus == Removed {
		if n.LegacyAttribute == "" {
			ctx.Diagnostics.AddError("%s: legacy_status: removed requires legacy_attribute to be set", path)
		}
		if n.LegacyType == UndefinedLegacyAttributeType {
			ctx.Diagnostics.AddError("%s: legacy_status: removed requires legacy_type to be set", path)
		}
		if n.LegacyRestriction == UndefinedRestriction {
			ctx.Diagnostics.AddError("%s: legacy_status: removed requires legacy_restriction to be set", path)
		}
	}
}

// validateChild enforces the shape rules for a children-bucket node, including
// the scalar-wrap shape: a children entry without legacy_attribute is only valid
// when at least one inner attributes entry carries legacy_attribute (the inner
// attribute(s) inherit the prior flat scalar value(s)). Errors are accumulated
// on ctx.Diagnostics, see validate for the path-prefix convention.
func (n AttributeUpgradeDefinition) validateChild(ctx *Context, path string) {
	if n.LegacyStatus == Removed && n.LegacyAttribute == "" {
		ctx.Diagnostics.AddError("%s: legacy_status: removed requires legacy_attribute to be set", path)
	}
	if n.LegacyAttribute == "" && n.LegacyType == UndefinedLegacyAttributeType {
		if !n.hasInnerLegacyAttribute() && len(n.Attributes) == 0 && len(n.Children) == 0 {
			ctx.Diagnostics.AddError("%s: entry has neither legacy_attribute / legacy_type on the block nor any inner attributes/children annotations", path)
		}
	}
	for innerProp, innerNode := range n.Attributes {
		innerNode.validate(ctx, fmt.Sprintf("%s: attributes[%q]", path, innerProp))
	}
	for innerChild, innerNode := range n.Children {
		innerNode.validateChild(ctx, fmt.Sprintf("%s: children[%q]", path, innerChild))
	}
}

func (n AttributeUpgradeDefinition) hasInnerLegacyAttribute() bool {
	for _, inner := range n.Attributes {
		if inner.LegacyAttribute != "" {
			return true
		}
	}
	for _, inner := range n.Children {
		if inner.LegacyAttribute != "" || inner.hasInnerLegacyAttribute() {
			return true
		}
	}
	return false
}

// ParentDnVariantDefinition pins a single (parent_class, rn_prepend,
// wrapper_class, test_platform) tuple for a ClassDefinition.ParentDnVariants
// entry. Each variant describes one valid placement of the class under a
// distinct user-facing parent that reaches APIC via its own request path.
type ParentDnVariantDefinition struct {
	// User-facing parent class for this variant (e.g. "fvTenant"). Distinguishes
	// system-scoped placements from wrapped tenant-scoped placements that meta
	// `containedBy` lists as a single direct parent.
	ParentClass string `yaml:"parent_class"`
	// Intermediate RN segment that selects this variant. The generated resource
	// matches the user's `parent_dn` against this segment to route the API call.
	RnPrepend string `yaml:"rn_prepend"`
	// Implicit container the request nests the resource inside (e.g. "cloudCertStore"
	// for a tenant-scoped pkiKeyRing). Empty for variants that POST against a real,
	// user-addressable parent.
	WrapperClass string `yaml:"wrapper_class"`
	// Platform profile that exercises this variant in tests (apic / cloud / both).
	// Lets testvars.yaml.tmpl gate variant-specific test cases without splitting the
	// test file.
	TestPlatform PlatformTypeEnum `yaml:"test_platform"`
}

func loadClassDefinition(className string) ClassDefinition {
	classDefinitionPath := fmt.Sprintf("%s/%s.yaml", constDefinitionsPath, className)
	var classDefinitionData ClassDefinition

	classDefinitionBytes, err := os.ReadFile(classDefinitionPath)
	if err != nil {
		genLogger.Debugf("The file '%s' was not found in the definitions folder.", classDefinitionPath)
		return classDefinitionData
	}

	classDefinitionData, err = parseClassDefinition(classDefinitionBytes)
	if err != nil {
		genLogger.Fatal(err.Error())
	}

	return classDefinitionData
}

// parseClassDefinition decodes raw YAML bytes into a ClassDefinition.
// UnmarshalStrict rejects unknown YAML keys so renamed/typo'd fields surface
// as a generator error instead of being silently ignored.
func parseClassDefinition(data []byte) (ClassDefinition, error) {
	var classDefinitionData ClassDefinition
	err := yaml.UnmarshalStrict(data, &classDefinitionData)
	return classDefinitionData, err
}
