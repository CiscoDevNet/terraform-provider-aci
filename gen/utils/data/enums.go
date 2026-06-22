package data

import "fmt"

// This file defines every typed enum used by the gen/utils/data package and the
// definition YAML schema. Two rules apply uniformly to all of them:
//
//  1. UnmarshalText is strict — a YAML field that is *present but empty* (e.g.
//     `reference_type: ""`) is treated as a typo and errors out. The zero value
//     is only reached by *omitting* the field from the YAML entirely.
//  2. String() round-trips the canonical YAML form of every valid constant.
//
// Beyond those, every enum falls into exactly one of three categories. Each
// category implies a different convention for the iota zero, whether an
// `Undefined…` sentinel exists, and what a consumer does with the zero value.
// The file is physically grouped into three sections below, separated by
// section banners. To add a new enum, pick its category and place it in the
// matching section so the convention stays obvious at a glance.
//
// ─── Category 1: Default IS a real value ─────────────────────────────────────
//   The default is declared as the iota zero of the enum itself. NO Undefined
//   sentinel exists. Consumers use the zero value directly with no fixup.
//   Members: PlatformTypeEnum, LegacyStatusEnum, ReferenceTypeEnum,
//   ValueRenderTypeEnum.
//
// ─── Category 2: No default; absence is a real runtime state ─────────────────
//   The iota zero IS an `UndefinedXxx` sentinel. Consumers either propagate it
//   ("the field really wasn't set, and that's meaningful") or use it as the
//   trigger for a validation error. Members: RelationshipTypeEnum,
//   TestDependencyRoleEnum, RegexStatementTypeEnum, MigrationSourceEnum.
//
// ─── Category 3: Default is a lookup, not a static value ─────────────────────
//   The iota zero is an `UndefinedXxx` sentinel meaning "no override; fall back
//   to another source" — typically the meta file or the current schema. The
//   default cannot be expressed as a static constant because it varies per
//   class/property. Members: ValueTypeEnum, RestrictionEnum,
//   LegacyAttributeTypeEnum.

// ════════════════════════════════════════════════════════════════════════════
// Category 1: Default IS a real value
// ════════════════════════════════════════════════════════════════════════════
// These enums have NO `UndefinedXxx` sentinel. The iota zero is a legitimate,
// documented default applied automatically when the YAML field is omitted.
// Consumers can use the zero value directly. UnmarshalText is strict — `""`
// errors as a typo guard, identical to every other enum in this file.

// PlatformTypeEnum represents the APIC platform type. The iota zero (Apic) is
// the documented default applied when the platform_type YAML field is omitted.
type PlatformTypeEnum int

const (
	// Apic indicates that the class is available on the on-premises version of APIC. Iota zero / default.
	Apic PlatformTypeEnum = iota
	// Both indicates that the class is available on both the on-premises and cloud versions of APIC.
	Both
	// Cloud indicates that the class is available on the cloud version of APIC.
	Cloud
)

func (p PlatformTypeEnum) String() string {
	switch p {
	case Both:
		return "both"
	case Cloud:
		return "cloud"
	default:
		return "apic"
	}
}

func (p *PlatformTypeEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: the zero value (Apic) is reached by
	// omitting the YAML field entirely. A field present-but-empty is treated as a typo
	// and errors out.
	switch string(text) {
	case "apic":
		*p = Apic
	case "cloud":
		*p = Cloud
	case "both":
		*p = Both
	default:
		return fmt.Errorf("unknown platform_type %q (expected one of: apic, cloud, both)", string(text))
	}
	return nil
}

// LegacyStatusEnum describes the current-schema exposure of a legacy attribute
// declared in state_upgrades. The iota zero (Functioning) is the natural state
// of a freshly renamed attribute where both legacy and replacement names work,
// and is the documented default applied when legacy_status is omitted.
type LegacyStatusEnum int

const (
	// Functioning is the iota zero / default: legacy name remains in the current schema
	// alongside the replacement attribute, with full device round-trip via the
	// replacement. The typical "we just renamed, both work" state.
	Functioning LegacyStatusEnum = iota
	// Frozen exposes the legacy name in the current schema with Deprecated but
	// without device round-trip. A state-preserving plan modifier suppresses the
	// diff so existing configs keep working without touching the device.
	Frozen
	// Removed drops the legacy name from the current schema entirely. The
	// state_upgrades entry is retained only to drive UpgradeResourceState.
	Removed
)

func (l LegacyStatusEnum) String() string {
	switch l {
	case Frozen:
		return "frozen"
	case Removed:
		return "removed"
	default:
		return "functioning"
	}
}

func (l *LegacyStatusEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: the zero value (Functioning) is
	// reached by omitting the YAML field entirely. A field present-but-empty is
	// treated as a typo and errors out.
	switch string(text) {
	case "functioning":
		*l = Functioning
	case "frozen":
		*l = Frozen
	case "removed":
		*l = Removed
	default:
		return fmt.Errorf("unknown legacy_status %q (expected one of: functioning, frozen, removed)", string(text))
	}
	return nil
}

// ReferenceTypeEnum indicates how to interpret a TestDependency.Reference value.
// The iota zero (ResourceReference) is the documented default applied when
// reference_type is omitted from YAML — by far the most common reference shape
// in test configs.
type ReferenceTypeEnum int

const (
	// ResourceReference is a Terraform resource attribute path (e.g. "aci_tenant.test.id").
	// Iota zero / default.
	ResourceReference ReferenceTypeEnum = iota
	// StaticReference is a hardcoded DN string (e.g. "uni/vmmp-VMware/dom-domain_1").
	StaticReference
	// DataSourceReference is a Terraform data source attribute path (e.g. "data.aci_tenant.test.id").
	DataSourceReference
)

func (r ReferenceTypeEnum) String() string {
	switch r {
	case StaticReference:
		return "static"
	case DataSourceReference:
		return "data_source"
	default:
		return "resource"
	}
}

func (r *ReferenceTypeEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: the zero value (ResourceReference) is
	// reached by omitting the YAML field entirely. A field present-but-empty is treated
	// as a typo and errors out.
	switch string(text) {
	case "resource":
		*r = ResourceReference
	case "static":
		*r = StaticReference
	case "data_source":
		*r = DataSourceReference
	default:
		return fmt.Errorf("unknown reference_type %q (expected one of: resource, static, data_source)", string(text))
	}
	return nil
}

// ValueRenderTypeEnum controls how a TestValueEntry is rendered in HCL
// configuration. The iota zero (StringValue) is the documented default applied
// when value_type is omitted from YAML — the overwhelmingly common case for
// property test values.
type ValueRenderTypeEnum int

const (
	// StringValue renders as a quoted string: attribute = "value". Iota zero / default.
	StringValue ValueRenderTypeEnum = iota
	// ReferenceValue renders as an unquoted reference expression: attribute = aci_tenant.test.id
	ReferenceValue
)

func (v ValueRenderTypeEnum) String() string {
	switch v {
	case ReferenceValue:
		return "reference"
	default:
		return "string"
	}
}

func (v *ValueRenderTypeEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: the zero value (StringValue) is
	// reached by omitting the YAML field entirely. A field present-but-empty is treated
	// as a typo and errors out.
	switch string(text) {
	case "string":
		*v = StringValue
	case "reference":
		*v = ReferenceValue
	default:
		return fmt.Errorf("unknown value_type %q (expected one of: string, reference)", string(text))
	}
	return nil
}

// referenceValueRenderType maps a dependency's ReferenceType to the HCL render
// type its value should use: StaticReference DNs are quoted strings, Terraform
// resource and data-source paths are unquoted expressions.
func referenceValueRenderType(referenceType ReferenceTypeEnum) ValueRenderTypeEnum {
	if referenceType == StaticReference {
		return StringValue
	}
	return ReferenceValue
}

// ════════════════════════════════════════════════════════════════════════════
// Category 2: No default; absence is a real runtime state
// ════════════════════════════════════════════════════════════════════════════
// These enums HAVE an `UndefinedXxx` iota zero that consumers check explicitly.
// The zero value either propagates as a meaningful "field not set" signal or
// triggers a validation error (e.g. depth-0 dependencies must declare a role).
// There is no static default — absence is the runtime state.

// RelationshipTypeEnum identifies how a relational class connects from one MO to another.
type RelationshipTypeEnum int

const (
	// UndefinedRelationshipType is the zero value: no relationship type assigned.
	UndefinedRelationshipType RelationshipTypeEnum = iota
	// Named indicates that the relationship is a named relation.
	Named
	// Explicit indicates that the relationship is an explicit relation.
	Explicit
)

func (r RelationshipTypeEnum) String() string {
	switch r {
	case Named:
		return "named"
	case Explicit:
		return "explicit"
	default:
		return ""
	}
}

func (r *RelationshipTypeEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: the zero value (UndefinedRelationshipType)
	// is reached by omitting the YAML field entirely. A field present-but-empty is
	// treated as a typo and errors out.
	switch string(text) {
	case "named":
		*r = Named
	case "explicit":
		*r = Explicit
	default:
		return fmt.Errorf("unknown relationship type %q (expected one of: named, explicit)", string(text))
	}
	return nil
}

// TestDependencyRoleEnum indicates how the dependency is consumed in HCL.
type TestDependencyRoleEnum int

const (
	// UndefinedRole is the zero value: no role assigned. Valid for nested
	// dependencies (depth > 0) which are pure prerequisites; invalid at depth 0
	// where Parent or Target is required.
	UndefinedRole TestDependencyRoleEnum = iota
	// Parent means the dependency provides the parent_dn attribute.
	Parent
	// Target means the dependency provides the target_dn attribute (relation classes).
	Target
)

func (r TestDependencyRoleEnum) String() string {
	switch r {
	case Parent:
		return "parent"
	case Target:
		return "target"
	default:
		return ""
	}
}

func (r *TestDependencyRoleEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: the zero value (UndefinedRole) is
	// reached by omitting the YAML field entirely. A field present-but-empty is
	// treated as a typo and errors out.
	switch string(text) {
	case "parent":
		*r = Parent
	case "target":
		*r = Target
	default:
		return fmt.Errorf("unknown role %q (expected one of: parent, target)", string(text))
	}
	return nil
}

// RegexStatementTypeEnum identifies the kind of regex match constraint in a Validator.
type RegexStatementTypeEnum int

const (
	// UndefinedRegexStatementType is the zero value: no regex statement type assigned.
	UndefinedRegexStatementType RegexStatementTypeEnum = iota
	// Include indicates that the value must match the regex statement.
	Include
)

func (r RegexStatementTypeEnum) String() string {
	switch r {
	case Include:
		return "include"
	default:
		return ""
	}
}

func (r *RegexStatementTypeEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: the zero value (UndefinedRegexStatementType)
	// is reached by omitting the field entirely. A field present-but-empty is treated
	// as a typo and errors out.
	switch string(text) {
	case "include":
		*r = Include
	default:
		return fmt.Errorf("unknown regex statement type %q (expected one of: include)", string(text))
	}
	return nil
}

// MigrationSourceEnum records the lineage of a resource — the prior provider or
// generator the resource was migrated from. Drives the documentation migration
// warning and any future migration-source-specific codegen. Extensible: additional
// sources can be added as plain new iota constants + String()/UnmarshalText cases
// without touching the field type or consumers that check != UndefinedMigrationSource.
type MigrationSourceEnum int

const (
	// UndefinedMigrationSource is the zero value: no migration history. The
	// resource was born in the current framework provider; no migration warning
	// is rendered in the docs.
	UndefinedMigrationSource MigrationSourceEnum = iota
	// FromSDKv2 indicates the resource was migrated from the SDKv2 provider
	// implementation. Renders the SDKv2-specific migration warning in the docs.
	FromSDKv2
)

func (m MigrationSourceEnum) String() string {
	switch m {
	case FromSDKv2:
		return "from_sdkv2"
	default:
		return ""
	}
}

func (m *MigrationSourceEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: the zero value
	// (UndefinedMigrationSource) is reached by omitting the YAML field entirely.
	// A field present-but-empty is treated as a typo and errors out.
	switch string(text) {
	case "from_sdkv2":
		*m = FromSDKv2
	default:
		return fmt.Errorf("unknown migration_source %q (expected one of: from_sdkv2)", string(text))
	}
	return nil
}
// ArtifactEnum identifies a generated artifact kind that the renderer can
// produce for a class. Used inside `ClassDefinition.Artifacts` to control
// which artifacts are emitted. A nil slice (the YAML field omitted entirely)
// is the signal to auto-derive the default set from IdentifiedBy; an empty
// slice (`artifacts: []`) excludes the class from both `provider.Resources()`
// and `provider.DataSources()`.
type ArtifactEnum int

const (
	// UndefinedArtifact is the zero value: no artifact kind assigned. Reaching
	// this inside a slice element means the YAML carried a present-but-empty
	// entry (treated as a typo). Field-level absence yields a nil slice instead.
	UndefinedArtifact ArtifactEnum = iota
	// ResourceArtifact identifies the Terraform resource artifact.
	ResourceArtifact
	// DatasourceArtifact identifies the Terraform datasource artifact.
	DatasourceArtifact
)

func (a ArtifactEnum) String() string {
	switch a {
	case ResourceArtifact:
		return "resource"
	case DatasourceArtifact:
		return "datasource"
	default:
		return ""
	}
}

func (a *ArtifactEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: the zero value (UndefinedArtifact)
	// is reached only by omitting the slice element entirely. A present-but-empty
	// entry is treated as a typo and errors out.
	switch string(text) {
	case "resource":
		*a = ResourceArtifact
	case "datasource":
		*a = DatasourceArtifact
	default:
		return fmt.Errorf("unknown artifact %q (expected one of: resource, datasource)", string(text))
	}
	return nil
}

// IgnoreTestEnum identifies a test artifact kind that should be skipped by
// the test renderer. Used inside `ClassTestConfigDefinition.IgnoreTests` to
// suppress generation of specific test buckets. The legacy
// `exclude_from_testing: true` flag migrates to `ignore_tests: [child]`; the
// `resource` and `datasource` values have no legacy driver and exist as
// future opt-ins for skipping resource / datasource acceptance tests.
type IgnoreTestEnum int

const (
	// UndefinedIgnoreTest is the zero value: no ignore-test target assigned.
	// Reaching this inside a slice element means the YAML carried a
	// present-but-empty entry (treated as a typo). Field-level absence yields
	// a nil slice instead.
	UndefinedIgnoreTest IgnoreTestEnum = iota
	// ChildIgnoreTest skips the child-as-part-of-parent test bucket
	// (the bucket exercised when the class is rendered as a child block of
	// another resource's HCL config).
	ChildIgnoreTest
	// ResourceIgnoreTest skips the resource acceptance test bucket.
	ResourceIgnoreTest
	// DatasourceIgnoreTest skips the datasource acceptance test bucket.
	DatasourceIgnoreTest
)

func (i IgnoreTestEnum) String() string {
	switch i {
	case ChildIgnoreTest:
		return "child"
	case ResourceIgnoreTest:
		return "resource"
	case DatasourceIgnoreTest:
		return "datasource"
	default:
		return ""
	}
}

func (i *IgnoreTestEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: the zero value (UndefinedIgnoreTest)
	// is reached only by omitting the slice element entirely. A present-but-empty
	// entry is treated as a typo and errors out.
	switch string(text) {
	case "child":
		*i = ChildIgnoreTest
	case "resource":
		*i = ResourceIgnoreTest
	case "datasource":
		*i = DatasourceIgnoreTest
	default:
		return fmt.Errorf("unknown ignore_tests entry %q (expected one of: child, resource, datasource)", string(text))
	}
	return nil
}
// ════════════════════════════════════════════════════════════════════════════
// Category 3: Default is a lookup, not a static value
// ════════════════════════════════════════════════════════════════════════════
// These enums HAVE an `UndefinedXxx` iota zero that means "no override; consult
// another source for the actual value." The fallback source varies per consumer
// — typically the meta file or the current-schema attribute type — and cannot
// be expressed as a static constant on the enum itself.

// ValueTypeEnum identifies the data shape of a property and is the single dispatch point
// for type-specific schema and template behavior. New named custom types should be added
// here as constants and registered in UnmarshalText so the same definition `value_type`
// override key can select them.
type ValueTypeEnum int

const (
	// UndefinedValueType is the zero value: no explicit override. Consumers (e.g.
	// setValueType) treat this as "fall back to meta-derived type".
	UndefinedValueType ValueTypeEnum = iota
	// String indicates that the property is a plain string value.
	String
	// Set indicates that the property is a set value (driven by meta uitype "bitmask").
	Set
	// Object indicates that the property is a structured object — a single nested
	// block or a map of attributes — as opposed to a flat scalar or a homogeneous
	// collection. Templates render this with the object-shaped path.
	Object
	// IpAddress indicates that the property is an IP address (IPv4 or IPv6); driven by
	// meta `validateAsIPv4OrIPv6`. Renders with the IP-address custom type for parsing,
	// validation, and semantic-equality (e.g. zero-padding normalization).
	IpAddress
	// SemanticEquality indicates that the property has both ValidValues and Validators,
	// meaning the wire form (e.g. "22") and the human form (e.g. "ssh") must compare equal.
	// Templates render this with the semantic-equality custom type.
	SemanticEquality
)

func (v ValueTypeEnum) String() string {
	switch v {
	case String:
		return "string"
	case Set:
		return "set"
	case Object:
		return "object"
	case IpAddress:
		return "ip_address"
	case SemanticEquality:
		return "semantic_equality"
	default:
		return ""
	}
}

func (v *ValueTypeEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: the zero value (UndefinedValueType) is
	// reached by omitting the YAML field entirely, and setValueType then falls back to
	// the meta-derived type. A field present-but-empty is treated as a typo and errors out.
	switch string(text) {
	case "string":
		*v = String
	case "set":
		*v = Set
	case "object":
		*v = Object
	case "ip_address":
		*v = IpAddress
	case "semantic_equality":
		*v = SemanticEquality
	default:
		return fmt.Errorf("unknown value_type %q (expected one of: string, set, ip_address, semantic_equality, object)", string(text))
	}
	return nil
}

// RestrictionEnum controls the schema behavior of a property as declared by the class definition.
type RestrictionEnum int

const (
	// UndefinedRestriction is the zero value: derive behavior from the meta file.
	UndefinedRestriction RestrictionEnum = iota
	// Required marks the property as required in the resource and datasource schemas.
	Required
	// Optional marks the property as optional in the resource and datasource schemas.
	Optional
	// ReadOnly includes a non-configurable property as a computed-only attribute.
	ReadOnly
	// Exclude omits the property entirely from generated schemas.
	Exclude
)

func (r RestrictionEnum) String() string {
	switch r {
	case Required:
		return "required"
	case Optional:
		return "optional"
	case ReadOnly:
		return "read_only"
	case Exclude:
		return "exclude"
	default:
		return ""
	}
}

func (r *RestrictionEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: the zero value (UndefinedRestriction)
	// is reached by omitting the YAML field entirely. A field present-but-empty is
	// treated as a typo and errors out.
	switch string(text) {
	case "required":
		*r = Required
	case "optional":
		*r = Optional
	case "read_only":
		*r = ReadOnly
	case "exclude":
		*r = Exclude
	default:
		return fmt.Errorf("unknown restriction %q (expected one of: required, optional, read_only, exclude)", string(text))
	}
	return nil
}

// LegacyAttributeTypeEnum identifies the Terraform plugin framework attribute type of a
// prior-schema attribute. Used inside state_upgrades to describe how the saved state
// should be deserialized when the type differs from the current schema.
type LegacyAttributeTypeEnum int

const (
	// UndefinedLegacyAttributeType is the zero value: the prior-schema type is the
	// same as the current-schema attribute type (no override needed).
	UndefinedLegacyAttributeType LegacyAttributeTypeEnum = iota
	StringAttribute
	BoolAttribute
	Int64Attribute
	Float64Attribute
	ListAttribute
	SetAttribute
	MapAttribute
	SingleNestedAttribute
	ListNestedAttribute
	SetNestedAttribute
	MapNestedAttribute
)

func (l LegacyAttributeTypeEnum) String() string {
	switch l {
	case StringAttribute:
		return "string_attribute"
	case BoolAttribute:
		return "bool_attribute"
	case Int64Attribute:
		return "int64_attribute"
	case Float64Attribute:
		return "float64_attribute"
	case ListAttribute:
		return "list_attribute"
	case SetAttribute:
		return "set_attribute"
	case MapAttribute:
		return "map_attribute"
	case SingleNestedAttribute:
		return "single_nested_attribute"
	case ListNestedAttribute:
		return "list_nested_attribute"
	case SetNestedAttribute:
		return "set_nested_attribute"
	case MapNestedAttribute:
		return "map_nested_attribute"
	default:
		return ""
	}
}

func (l *LegacyAttributeTypeEnum) UnmarshalText(text []byte) error {
	// Empty/missing is intentionally not handled: callers reach the zero value
	// (UndefinedLegacyAttributeType) by omitting the YAML field entirely. A field
	// present-but-empty is treated as a typo and errors out.
	switch string(text) {
	case "string_attribute":
		*l = StringAttribute
	case "bool_attribute":
		*l = BoolAttribute
	case "int64_attribute":
		*l = Int64Attribute
	case "float64_attribute":
		*l = Float64Attribute
	case "list_attribute":
		*l = ListAttribute
	case "set_attribute":
		*l = SetAttribute
	case "map_attribute":
		*l = MapAttribute
	case "single_nested_attribute":
		*l = SingleNestedAttribute
	case "list_nested_attribute":
		*l = ListNestedAttribute
	case "set_nested_attribute":
		*l = SetNestedAttribute
	case "map_nested_attribute":
		*l = MapNestedAttribute
	default:
		return fmt.Errorf("unknown legacy_type %q (expected one of: string_attribute, bool_attribute, int64_attribute, float64_attribute, list_attribute, set_attribute, map_attribute, single_nested_attribute, list_nested_attribute, set_nested_attribute, map_nested_attribute)", string(text))
	}
	return nil
}

func legacyTypeToValueType(t LegacyAttributeTypeEnum) ValueTypeEnum {
	// Map the framework-attribute-typed enum used in state_upgrades to the
	// renderer's ValueTypeEnum vocabulary. Collection-shaped types collapse to
	// Set, single-nested and map-shaped types collapse to Object, and every
	// scalar (and the zero value) collapses to String. The renderer reads
	// LegacyAttributeType directly from the StateUpgrades tree when it needs
	// framework-level fidelity (e.g. SingleNestedAttribute vs MapAttribute).
	// TODO: extend the switch with explicit cases for IpAddress, SemanticEquality
	switch t {
	case ListAttribute, ListNestedAttribute, SetAttribute, SetNestedAttribute:
		return Set
	case MapAttribute, MapNestedAttribute, SingleNestedAttribute:
		return Object
	default:
		return String
	}
}
