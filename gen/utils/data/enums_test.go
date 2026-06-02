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

	assert.Equal(t, Apic, test.MustUnmarshalText[PlatformTypeEnum](t, ""))
	assert.Equal(t, Apic, test.MustUnmarshalText[PlatformTypeEnum](t, "apic"))
	assert.Equal(t, Cloud, test.MustUnmarshalText[PlatformTypeEnum](t, "cloud"))
	assert.Equal(t, Both, test.MustUnmarshalText[PlatformTypeEnum](t, "both"))
	assert.Error(t, test.UnmarshalTextErr[PlatformTypeEnum]("garbage"))
}

func TestRelationshipTypeEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "", UndefinedRelationshipType.String())
	assert.Equal(t, "named", Named.String())
	assert.Equal(t, "explicit", Explicit.String())

	assert.Equal(t, UndefinedRelationshipType, test.MustUnmarshalText[RelationshipTypeEnum](t, ""))
	assert.Equal(t, Named, test.MustUnmarshalText[RelationshipTypeEnum](t, "named"))
	assert.Equal(t, Explicit, test.MustUnmarshalText[RelationshipTypeEnum](t, "explicit"))
	assert.Error(t, test.UnmarshalTextErr[RelationshipTypeEnum]("garbage"))
}

func TestTestDependencyRoleEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "", UndefinedRole.String())
	assert.Equal(t, "parent", Parent.String())
	assert.Equal(t, "target", Target.String())

	assert.Equal(t, UndefinedRole, test.MustUnmarshalText[TestDependencyRoleEnum](t, ""))
	assert.Equal(t, Parent, test.MustUnmarshalText[TestDependencyRoleEnum](t, "parent"))
	assert.Equal(t, Target, test.MustUnmarshalText[TestDependencyRoleEnum](t, "target"))
	assert.Error(t, test.UnmarshalTextErr[TestDependencyRoleEnum]("garbage"))
}

func TestReferenceTypeEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "", UndefinedReferenceType.String())
	assert.Equal(t, "static", StaticReference.String())
	assert.Equal(t, "resource", ResourceReference.String())
	assert.Equal(t, "data_source", DataSourceReference.String())

	// Empty defaults to ResourceReference (documented default when reference_type is omitted).
	assert.Equal(t, ResourceReference, test.MustUnmarshalText[ReferenceTypeEnum](t, ""))
	assert.Equal(t, StaticReference, test.MustUnmarshalText[ReferenceTypeEnum](t, "static"))
	assert.Equal(t, ResourceReference, test.MustUnmarshalText[ReferenceTypeEnum](t, "resource"))
	assert.Equal(t, DataSourceReference, test.MustUnmarshalText[ReferenceTypeEnum](t, "data_source"))
	assert.Error(t, test.UnmarshalTextErr[ReferenceTypeEnum]("garbage"))
}

func TestRegexStatementTypeEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "include", Include.String())
	assert.Equal(t, Include, test.MustUnmarshalText[RegexStatementTypeEnum](t, "include"))
	assert.Error(t, test.UnmarshalTextErr[RegexStatementTypeEnum](""))
	assert.Error(t, test.UnmarshalTextErr[RegexStatementTypeEnum]("garbage"))
}

func TestValueRenderTypeEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "string", StringValue.String())
	assert.Equal(t, "reference", ReferenceValue.String())

	// Empty and "string" both decode to StringValue (the default render mode).
	assert.Equal(t, StringValue, test.MustUnmarshalText[ValueRenderTypeEnum](t, ""))
	assert.Equal(t, StringValue, test.MustUnmarshalText[ValueRenderTypeEnum](t, "string"))
	assert.Equal(t, ReferenceValue, test.MustUnmarshalText[ValueRenderTypeEnum](t, "reference"))
	assert.Error(t, test.UnmarshalTextErr[ValueRenderTypeEnum]("garbage"))
}

func TestValueTypeEnum(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t, "string", String.String())
	assert.Equal(t, "set", Set.String())
	assert.Equal(t, "ip_address", IpAddress.String())
	assert.Equal(t, "semantic_equality", SemanticEquality.String())

	// Empty leaves the zero value (meaning "not overridden"); callers use that to fall back to meta-derived type.
	assert.Equal(t, ValueTypeEnum(0), test.MustUnmarshalText[ValueTypeEnum](t, ""))
	assert.Equal(t, String, test.MustUnmarshalText[ValueTypeEnum](t, "string"))
	assert.Equal(t, Set, test.MustUnmarshalText[ValueTypeEnum](t, "set"))
	assert.Equal(t, IpAddress, test.MustUnmarshalText[ValueTypeEnum](t, "ip_address"))
	assert.Equal(t, SemanticEquality, test.MustUnmarshalText[ValueTypeEnum](t, "semantic_equality"))
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

	assert.Equal(t, UndefinedRestriction, test.MustUnmarshalText[RestrictionEnum](t, ""))
	assert.Equal(t, Required, test.MustUnmarshalText[RestrictionEnum](t, "required"))
	assert.Equal(t, Optional, test.MustUnmarshalText[RestrictionEnum](t, "optional"))
	assert.Equal(t, ReadOnly, test.MustUnmarshalText[RestrictionEnum](t, "read_only"))
	assert.Equal(t, Exclude, test.MustUnmarshalText[RestrictionEnum](t, "exclude"))
	assert.Error(t, test.UnmarshalTextErr[RestrictionEnum]("garbage"))
}
