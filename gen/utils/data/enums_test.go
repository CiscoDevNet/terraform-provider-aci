package data

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestPlatformTypeEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "apic", Apic.String())
	assert.Equal(t, "cloud", Cloud.String())
	assert.Equal(t, "both", Both.String())

	assert.Equal(t, Apic, test.MustUnmarshalText[PlatformTypeEnum](t, "apic"))
	assert.Equal(t, Cloud, test.MustUnmarshalText[PlatformTypeEnum](t, "cloud"))
	assert.Equal(t, Both, test.MustUnmarshalText[PlatformTypeEnum](t, "both"))

	// Category 1: iota zero IS the default (Apic). Omitting the YAML field yields
	// this automatically; a field present-but-empty is treated as a typo.
	var zero PlatformTypeEnum
	assert.Equal(t, Apic, zero)
	assert.Error(t, test.UnmarshalTextErr[PlatformTypeEnum](""))
	assert.Error(t, test.UnmarshalTextErr[PlatformTypeEnum]("garbage"))
}

func TestRelationshipTypeEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "", UndefinedRelationshipType.String())
	assert.Equal(t, "named", Named.String())
	assert.Equal(t, "explicit", Explicit.String())

	assert.Equal(t, Named, test.MustUnmarshalText[RelationshipTypeEnum](t, "named"))
	assert.Equal(t, Explicit, test.MustUnmarshalText[RelationshipTypeEnum](t, "explicit"))

	// Empty and unknown both error: the zero value (UndefinedRelationshipType) is
	// reached by omitting the YAML field entirely, so a field present-but-empty is
	// treated as a typo.
	assert.Error(t, test.UnmarshalTextErr[RelationshipTypeEnum](""))
	assert.Error(t, test.UnmarshalTextErr[RelationshipTypeEnum]("garbage"))
}

func TestTestDependencyRoleEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "", UndefinedRole.String())
	assert.Equal(t, "parent", Parent.String())
	assert.Equal(t, "target", Target.String())

	assert.Equal(t, Parent, test.MustUnmarshalText[TestDependencyRoleEnum](t, "parent"))
	assert.Equal(t, Target, test.MustUnmarshalText[TestDependencyRoleEnum](t, "target"))

	// Empty and unknown both error: the zero value (UndefinedRole) is reached by
	// omitting the YAML field entirely, so a field present-but-empty is treated as
	// a typo.
	assert.Error(t, test.UnmarshalTextErr[TestDependencyRoleEnum](""))
	assert.Error(t, test.UnmarshalTextErr[TestDependencyRoleEnum]("garbage"))
}

func TestReferenceTypeEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "resource", ResourceReference.String())
	assert.Equal(t, "static", StaticReference.String())
	assert.Equal(t, "data_source", DataSourceReference.String())

	assert.Equal(t, ResourceReference, test.MustUnmarshalText[ReferenceTypeEnum](t, "resource"))
	assert.Equal(t, StaticReference, test.MustUnmarshalText[ReferenceTypeEnum](t, "static"))
	assert.Equal(t, DataSourceReference, test.MustUnmarshalText[ReferenceTypeEnum](t, "data_source"))

	// Category 1: iota zero IS the default (ResourceReference). Omitting reference_type
	// in YAML yields this automatically; a field present-but-empty is treated as a typo.
	var zero ReferenceTypeEnum
	assert.Equal(t, ResourceReference, zero)
	assert.Error(t, test.UnmarshalTextErr[ReferenceTypeEnum](""))
	assert.Error(t, test.UnmarshalTextErr[ReferenceTypeEnum]("garbage"))
}

func TestRegexStatementTypeEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "", UndefinedRegexStatementType.String())
	assert.Equal(t, "include", Include.String())

	assert.Equal(t, Include, test.MustUnmarshalText[RegexStatementTypeEnum](t, "include"))

	// Empty and unknown both error: the zero value (UndefinedRegexStatementType) is
	// reached by omitting the field entirely, so a field present-but-empty is treated
	// as a typo.
	assert.Error(t, test.UnmarshalTextErr[RegexStatementTypeEnum](""))
	assert.Error(t, test.UnmarshalTextErr[RegexStatementTypeEnum]("garbage"))
}

func TestValueRenderTypeEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "string", StringValue.String())
	assert.Equal(t, "reference", ReferenceValue.String())

	assert.Equal(t, StringValue, test.MustUnmarshalText[ValueRenderTypeEnum](t, "string"))
	assert.Equal(t, ReferenceValue, test.MustUnmarshalText[ValueRenderTypeEnum](t, "reference"))

	// Category 1: iota zero IS the default (StringValue). Omitting value_type in YAML
	// yields this automatically; a field present-but-empty is treated as a typo.
	var zero ValueRenderTypeEnum
	assert.Equal(t, StringValue, zero)
	assert.Error(t, test.UnmarshalTextErr[ValueRenderTypeEnum](""))
	assert.Error(t, test.UnmarshalTextErr[ValueRenderTypeEnum]("garbage"))
}

func TestValueTypeEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "", UndefinedValueType.String())
	assert.Equal(t, "string", String.String())
	assert.Equal(t, "set", Set.String())
	assert.Equal(t, "ip_address", IpAddress.String())
	assert.Equal(t, "semantic_equality", SemanticEquality.String())
	assert.Equal(t, "object", Object.String())

	assert.Equal(t, String, test.MustUnmarshalText[ValueTypeEnum](t, "string"))
	assert.Equal(t, Set, test.MustUnmarshalText[ValueTypeEnum](t, "set"))
	assert.Equal(t, IpAddress, test.MustUnmarshalText[ValueTypeEnum](t, "ip_address"))
	assert.Equal(t, SemanticEquality, test.MustUnmarshalText[ValueTypeEnum](t, "semantic_equality"))
	assert.Equal(t, Object, test.MustUnmarshalText[ValueTypeEnum](t, "object"))

	// Empty and unknown both error: the zero value (UndefinedValueType) is reached by
	// omitting the YAML field entirely (setValueType then falls back to the meta-derived
	// type), so a field present-but-empty is treated as a typo.
	assert.Error(t, test.UnmarshalTextErr[ValueTypeEnum](""))
	assert.Error(t, test.UnmarshalTextErr[ValueTypeEnum]("weird"))
}

func TestRestrictionEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "", UndefinedRestriction.String())
	assert.Equal(t, "required", Required.String())
	assert.Equal(t, "optional", Optional.String())
	assert.Equal(t, "read_only", ReadOnly.String())
	assert.Equal(t, "exclude", Exclude.String())

	assert.Equal(t, Required, test.MustUnmarshalText[RestrictionEnum](t, "required"))
	assert.Equal(t, Optional, test.MustUnmarshalText[RestrictionEnum](t, "optional"))
	assert.Equal(t, ReadOnly, test.MustUnmarshalText[RestrictionEnum](t, "read_only"))
	assert.Equal(t, Exclude, test.MustUnmarshalText[RestrictionEnum](t, "exclude"))

	// Empty and unknown both error: the zero value (UndefinedRestriction) is
	// reached by omitting the YAML field entirely, so a field present-but-empty
	// is treated as a typo.
	assert.Error(t, test.UnmarshalTextErr[RestrictionEnum](""))
	assert.Error(t, test.UnmarshalTextErr[RestrictionEnum]("garbage"))
}

func TestLegacyAttributeTypeEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "", UndefinedLegacyAttributeType.String())
	assert.Equal(t, "string_attribute", StringAttribute.String())
	assert.Equal(t, "bool_attribute", BoolAttribute.String())
	assert.Equal(t, "int64_attribute", Int64Attribute.String())
	assert.Equal(t, "float64_attribute", Float64Attribute.String())
	assert.Equal(t, "list_attribute", ListAttribute.String())
	assert.Equal(t, "set_attribute", SetAttribute.String())
	assert.Equal(t, "map_attribute", MapAttribute.String())
	assert.Equal(t, "single_nested_attribute", SingleNestedAttribute.String())
	assert.Equal(t, "list_nested_attribute", ListNestedAttribute.String())
	assert.Equal(t, "set_nested_attribute", SetNestedAttribute.String())
	assert.Equal(t, "map_nested_attribute", MapNestedAttribute.String())

	assert.Equal(t, StringAttribute, test.MustUnmarshalText[LegacyAttributeTypeEnum](t, "string_attribute"))
	assert.Equal(t, BoolAttribute, test.MustUnmarshalText[LegacyAttributeTypeEnum](t, "bool_attribute"))
	assert.Equal(t, Int64Attribute, test.MustUnmarshalText[LegacyAttributeTypeEnum](t, "int64_attribute"))
	assert.Equal(t, Float64Attribute, test.MustUnmarshalText[LegacyAttributeTypeEnum](t, "float64_attribute"))
	assert.Equal(t, ListAttribute, test.MustUnmarshalText[LegacyAttributeTypeEnum](t, "list_attribute"))
	assert.Equal(t, SetAttribute, test.MustUnmarshalText[LegacyAttributeTypeEnum](t, "set_attribute"))
	assert.Equal(t, MapAttribute, test.MustUnmarshalText[LegacyAttributeTypeEnum](t, "map_attribute"))
	assert.Equal(t, SingleNestedAttribute, test.MustUnmarshalText[LegacyAttributeTypeEnum](t, "single_nested_attribute"))
	assert.Equal(t, ListNestedAttribute, test.MustUnmarshalText[LegacyAttributeTypeEnum](t, "list_nested_attribute"))
	assert.Equal(t, SetNestedAttribute, test.MustUnmarshalText[LegacyAttributeTypeEnum](t, "set_nested_attribute"))
	assert.Equal(t, MapNestedAttribute, test.MustUnmarshalText[LegacyAttributeTypeEnum](t, "map_nested_attribute"))

	// Empty and unknown both error: the zero value (UndefinedLegacyAttributeType)
	// is reached by omitting the YAML field entirely, so a field present-but-empty
	// is treated as a typo.
	assert.Error(t, test.UnmarshalTextErr[LegacyAttributeTypeEnum](""))
	assert.Error(t, test.UnmarshalTextErr[LegacyAttributeTypeEnum]("garbage"))
}

func TestLegacyTypeToValueType(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	// Exhaustive over every LegacyAttributeTypeEnum constant. Adding a new
	// LegacyAttributeTypeEnum must add a case here so the collapse policy stays
	// explicit at the test level.
	cases := []test.TestCase{
		// Zero value: no override; collapses to the String default.
		{Name: "undefined_to_string", Input: UndefinedLegacyAttributeType, Expected: String},

		// Scalar legacy types collapse to String.
		{Name: "string_to_string", Input: StringAttribute, Expected: String},
		{Name: "bool_to_string", Input: BoolAttribute, Expected: String},
		{Name: "int64_to_string", Input: Int64Attribute, Expected: String},
		{Name: "float64_to_string", Input: Float64Attribute, Expected: String},

		// Collection-shaped legacy types collapse to Set.
		{Name: "list_to_set", Input: ListAttribute, Expected: Set},
		{Name: "set_to_set", Input: SetAttribute, Expected: Set},
		{Name: "list_nested_to_set", Input: ListNestedAttribute, Expected: Set},
		{Name: "set_nested_to_set", Input: SetNestedAttribute, Expected: Set},

		// Map and single-nested legacy types collapse to Object.
		{Name: "map_to_object", Input: MapAttribute, Expected: Object},
		{Name: "single_nested_to_object", Input: SingleNestedAttribute, Expected: Object},
		{Name: "map_nested_to_object", Input: MapNestedAttribute, Expected: Object},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.Expected.(ValueTypeEnum), legacyTypeToValueType(tc.Input.(LegacyAttributeTypeEnum)))
		})
	}
}

func TestLegacyStatusEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "functioning", Functioning.String())
	assert.Equal(t, "frozen", Frozen.String())
	assert.Equal(t, "removed", Removed.String())

	assert.Equal(t, Functioning, test.MustUnmarshalText[LegacyStatusEnum](t, "functioning"))
	assert.Equal(t, Frozen, test.MustUnmarshalText[LegacyStatusEnum](t, "frozen"))
	assert.Equal(t, Removed, test.MustUnmarshalText[LegacyStatusEnum](t, "removed"))

	// Empty errors out: the zero value (Functioning) is reached by omitting the YAML
	// field, so an explicit empty value is treated as a typo.
	assert.Error(t, test.UnmarshalTextErr[LegacyStatusEnum](""))
	assert.Error(t, test.UnmarshalTextErr[LegacyStatusEnum]("garbage"))
}

func TestMigrationSourceEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "", UndefinedMigrationSource.String())
	assert.Equal(t, "from_sdkv2", FromSDKv2.String())

	assert.Equal(t, FromSDKv2, test.MustUnmarshalText[MigrationSourceEnum](t, "from_sdkv2"))

	// Empty errors out: the zero value (UndefinedMigrationSource) is reached by
	// omitting the YAML field, so an explicit empty value is treated as a typo.
	assert.Error(t, test.UnmarshalTextErr[MigrationSourceEnum](""))
	assert.Error(t, test.UnmarshalTextErr[MigrationSourceEnum]("garbage"))
}
