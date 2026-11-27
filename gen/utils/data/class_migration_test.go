package data

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestSchemaVersionIndependentOfUpgradeHistory(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "fvTenant", Input: 2, Expected: true},
		{Name: "fvAp", Input: 2, Expected: true},
		{Name: "fvAEPg", Input: 2, Expected: true},
		{Name: "vzTaboo", Input: 2, Expected: true},
		{Name: "fvBD", Input: 2, Expected: true},
		{Name: "fvCtx", Input: 2, Expected: true},
		{Name: "fvESg", Input: 2, Expected: true},
		{Name: "vmmDomP", Input: 2, Expected: true},
		{Name: "vmmUsrAccP", Input: 2, Expected: true},
		{Name: "infraAttEntityP", Input: 2, Expected: true},
		{Name: "netflowMonitorPol", Input: 1, Expected: true},
		{Name: "stpIfPol", Input: 2, Expected: true},
		{Name: "infraFexP", Input: 1, Expected: false},
		{Name: "fvSubnet", Input: nil, Expected: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			class := classWithProperties(testCase.Name, map[string]*Property{"name": {PropertyName: "name", AttributeName: "name", ValueType: String}}, nil)
			expectedVersion := 0
			if testCase.Input != nil {
				expectedVersion = testCase.Input.(int)
				class.ClassDefinition.SchemaVersion = &expectedVersion
			}
			class.PropertiesAll = []string{"name"}
			if testCase.Expected.(bool) {
				class.StateUpgrades = map[int]StateUpgradeDefinition{0: {Attributes: map[string]AttributeUpgradeDefinition{"name": {AttributeName: "old_name"}}}}
			}
			store := &DataStore{Classes: map[string]Class{}, ctx: NewContext()}
			class.setMigrationData(store)
			assert.NoError(t, store.ctx.Diagnostics.Error())
			assert.Equal(t, expectedVersion, class.SchemaVersion)
			if testCase.Expected.(bool) {
				assert.Len(t, class.LegacyAttributes, 1)
				assert.Len(t, class.MigrationSchemas, 1)
			} else {
				assert.Empty(t, class.LegacyAttributes)
				assert.Empty(t, class.MigrationSchemas)
			}
		})
	}
}

func TestDataSourceOnlyLegacyAttributes(t *testing.T) {
	test.InitializeTest(t)
	class := classWithProperties("topSystem", map[string]*Property{
		"bootstrapState": {PropertyName: "bootstrapState", AttributeName: "bootstrap_state", ValueType: String, ReadOnly: true, Computed: true},
	}, nil)
	class.Artifacts = []ArtifactEnum{DatasourceArtifact}
	class.StateUpgrades = map[int]StateUpgradeDefinition{1: {Source: FromSDKv2, Attributes: map[string]AttributeUpgradeDefinition{
		"bootstrapState": {AttributeName: "boot_strap_tate", Status: LegacyDeprecated},
	}}}
	store := &DataStore{Classes: map[string]Class{}, ctx: NewContext()}
	class.setMigrationData(store)
	assert.NoError(t, store.ctx.Diagnostics.Error())
	assert.Len(t, class.LegacyAttributes, 1)
	assert.Empty(t, class.MigrationSchemas, "datasource compatibility must not generate resource upgraders")
	assert.Zero(t, class.SchemaVersion)
}

func TestMigrationStatusInheritance(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "new_alias_defaults_to_deprecated", Input: []LegacyStatusEnum{UndefinedLegacyStatus, UndefinedLegacyStatus}, Expected: LegacyDeprecated},
		{Name: "frozen_inherited", Input: []LegacyStatusEnum{Frozen, UndefinedLegacyStatus}, Expected: Frozen},
		{Name: "removed_inherited", Input: []LegacyStatusEnum{Removed, UndefinedLegacyStatus}, Expected: Removed},
		{Name: "explicit_removal", Input: []LegacyStatusEnum{LegacyDeprecated, Removed}, Expected: Removed},
		{Name: "explicit_freezing", Input: []LegacyStatusEnum{LegacyDeprecated, Frozen}, Expected: Frozen},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			version := 5
			class := classWithProperties("fvTenant", map[string]*Property{"name": {PropertyName: "name", AttributeName: "name", ValueType: String}}, nil)
			class.ClassDefinition.SchemaVersion = &version
			statuses := testCase.Input.([]LegacyStatusEnum)
			// Supply history out of order and skip a version: resolution follows version numbers.
			class.StateUpgrades = map[int]StateUpgradeDefinition{
				3: {Attributes: map[string]AttributeUpgradeDefinition{"name": {AttributeName: "old_name", Type: StringAttribute, Restriction: LegacyOptionalComputed, Status: statuses[1]}}},
				1: {Attributes: map[string]AttributeUpgradeDefinition{"name": {AttributeName: "old_name", Type: StringAttribute, Restriction: LegacyOptionalComputed, Status: statuses[0]}}},
			}
			store := &DataStore{Classes: map[string]Class{}, ctx: NewContext()}
			class.setPropertyStateUpgradeValues()
			class.setMigrationData(store)
			assert.NoError(t, store.ctx.Diagnostics.Error())
			if assert.Len(t, class.LegacyAttributes, 1) {
				assert.Equal(t, testCase.Expected, class.LegacyAttributes[0].Status)
				assert.Equal(t, testCase.Expected, class.Properties["name"].StateUpgradeValues[1].Status)
			}
			if assert.Len(t, class.MigrationSchemas, 2) {
				assert.Equal(t, 1, class.MigrationSchemas[0].PriorSchemaVersion)
				for _, attribute := range class.MigrationSchemas[0].Attributes {
					if attribute.AttributeName == "old_name" {
						assert.Equal(t, testCase.Expected, attribute.Status)
					}
				}
			}
		})
	}
}

func TestMigrationBlockStatusInheritance(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "new_block_defaults_to_deprecated", Input: []LegacyStatusEnum{UndefinedLegacyStatus, UndefinedLegacyStatus}, Expected: LegacyDeprecated},
		{Name: "frozen_inherited", Input: []LegacyStatusEnum{Frozen, UndefinedLegacyStatus}, Expected: Frozen},
		{Name: "removed_inherited", Input: []LegacyStatusEnum{Removed, UndefinedLegacyStatus}, Expected: Removed},
		{Name: "explicit_removal", Input: []LegacyStatusEnum{LegacyDeprecated, Removed}, Expected: Removed},
		{Name: "explicit_freezing", Input: []LegacyStatusEnum{LegacyDeprecated, Frozen}, Expected: Frozen},
		{Name: "explicit_deprecation", Input: []LegacyStatusEnum{Frozen, LegacyDeprecated}, Expected: LegacyDeprecated},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			child := classWithProperties("fvRsNodeAtt", map[string]*Property{
				"tDn": {PropertyName: "tDn", AttributeName: "target_dn", ValueType: String, Required: true},
			}, nil)
			child.IdentifiedBy = []string{"tDn"}
			child.ResourceNameNested = "relation_to_static_leafs"
			class := classWithProperties("fvAEPg", map[string]*Property{}, []string{"fvRsNodeAtt"})
			version := 5
			class.ClassDefinition.SchemaVersion = &version
			statuses := testCase.Input.([]LegacyStatusEnum)
			members := map[string]AttributeUpgradeDefinition{"tDn": {AttributeName: "node_dn"}}
			// Supply history out of order and skip a version: resolution follows version numbers.
			class.StateUpgrades = map[int]StateUpgradeDefinition{
				3: {Children: map[string]AttributeUpgradeDefinition{
					"fvRsNodeAtt": {AttributeName: "relation_fv_rs_node_att", Status: statuses[1], Attributes: members},
				}},
				1: {Children: map[string]AttributeUpgradeDefinition{
					"fvRsNodeAtt": {AttributeName: "relation_fv_rs_node_att", Status: statuses[0], Attributes: members},
				}},
			}
			store := &DataStore{Classes: map[string]Class{"fvRsNodeAtt": *child}, ctx: NewContext()}
			class.setMigrationData(store)
			assert.NoError(t, store.ctx.Diagnostics.Error())
			if assert.Len(t, class.LegacyAttributes, 1) {
				assert.Equal(t, testCase.Expected, class.LegacyAttributes[0].Status)
			}
			if assert.Len(t, class.MigrationSchemas, 2) {
				assert.Equal(t, 1, class.MigrationSchemas[0].PriorSchemaVersion)
				assert.Equal(t, 3, class.MigrationSchemas[1].PriorSchemaVersion)
				for _, schema := range class.MigrationSchemas {
					blocks := schema.Blocks()
					if assert.Len(t, blocks, 1) {
						assert.Equal(t, testCase.Expected, blocks[0].Status)
						assert.Len(t, blocks[0].Attributes, 1)
					}
				}
			}
		})
	}
}

func TestHistoricalRestrictions(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "required", Input: LegacyRequired, Expected: []bool{true, false, false}},
		{Name: "optional", Input: LegacyOptional, Expected: []bool{false, true, false}},
		{Name: "optional_computed", Input: LegacyOptionalComputed, Expected: []bool{false, true, true}},
		{Name: "read_only", Input: LegacyReadOnly, Expected: []bool{false, false, true}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			class := classWithProperties("fvTenant", map[string]*Property{"name": {PropertyName: "name", AttributeName: "name", ValueType: String}}, nil)
			store := &DataStore{ctx: NewContext()}
			result := class.migrationAttribute(store, "name", AttributeUpgradeDefinition{AttributeName: "name", Type: StringAttribute, Restriction: testCase.Input.(LegacyRestrictionEnum)}, nil)
			assert.NoError(t, store.ctx.Diagnostics.Error())
			assert.Equal(t, testCase.Expected, []bool{result.Required, result.Optional, result.Computed})
			assert.False(t, result.IsAlias(), "unchanged historical attributes are not deprecated aliases")
		})
	}
}

func TestMigrationRootAliasTargets(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "parent", Input: &Property{PropertyName: "parentDn", AttributeName: "parent_dn", ValueType: String}, Expected: "ParentDN"},
		{Name: "semantic_string", Input: &Property{PropertyName: "prio", AttributeName: "priority", ValueType: SemanticEquality}, Expected: "Prio"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			property := testCase.Input.(*Property)
			class := classWithProperties("fvAEPg", map[string]*Property{property.PropertyName: property}, nil)
			store := &DataStore{ctx: NewContext()}
			attribute := class.migrationAttribute(store, property.PropertyName, AttributeUpgradeDefinition{AttributeName: "old_attribute"}, nil)
			assert.NoError(t, store.ctx.Diagnostics.Error())
			assert.Equal(t, testCase.Expected, attribute.TargetField)
			assert.Equal(t, property.ModelValueType(class.Name), attribute.TargetModelType)
			assert.Equal(t, property.ModelAttributeType(class.Name), attribute.TargetAttributeType)
			assert.Equal(t, "types.String", attribute.ModelValueType(), "old values keep their historical decoding type")
			assert.Equal(t, property.PropertyName == "parentDn", attribute.RequiresReplace)
		})
	}
}

func TestMigrationHistoricalValueTypes(t *testing.T) {
	type input struct {
		Current    ValueTypeEnum
		Historical ValueTypeEnum
	}
	type expected struct {
		ValueType     string
		AttributeType string
		NullValue     string
		CustomType    string
		HasError      bool
	}
	testCases := []test.TestCase{
		{Name: "no_override", Input: input{Current: IpAddress}, Expected: expected{ValueType: "types.String", AttributeType: "types.StringType", NullValue: "types.StringNull()"}},
		{Name: "historical_ip", Input: input{Current: IpAddress, Historical: IpAddress}, Expected: expected{ValueType: "customtypes.IPAddressValue", AttributeType: "customtypes.IPAddressType{}", NullValue: "customtypes.NewIPAddressNull()", CustomType: "customtypes.IPAddressType{}"}},
		{Name: "current_type_changed", Input: input{Current: String, Historical: IpAddress}, Expected: expected{ValueType: "customtypes.IPAddressValue", AttributeType: "customtypes.IPAddressType{}", NullValue: "customtypes.NewIPAddressNull()", CustomType: "customtypes.IPAddressType{}"}},
		{Name: "historical_plain_string", Input: input{Current: IpAddress, Historical: String}, Expected: expected{ValueType: "types.String", AttributeType: "types.StringType", NullValue: "types.StringNull()"}},
		{Name: "unsupported_collection", Input: input{Current: String, Historical: Set}, Expected: expected{HasError: true}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			values := testCase.Input.(input)
			want := testCase.Expected.(expected)
			property := &Property{PropertyName: "llAddr", AttributeName: "link_local_ipv6_address", ValueType: values.Current}
			class := classWithProperties("fvBD", map[string]*Property{"llAddr": property}, nil)
			store := &DataStore{ctx: NewContext()}
			definition := AttributeUpgradeDefinition{AttributeName: "ll_addr", ValueType: values.Historical}
			attribute := class.migrationAttribute(store, "llAddr", definition, nil)
			assert.Equal(t, want.HasError, store.ctx.Diagnostics.Error() != nil)
			if want.HasError {
				return
			}
			assert.Equal(t, want.ValueType, attribute.ModelValueType())
			assert.Equal(t, want.AttributeType, attribute.ModelAttributeType())
			assert.Equal(t, want.NullValue, attribute.NullValue())
			assert.Equal(t, want.CustomType, attribute.ModelCustomStringType())
			assert.Equal(t, "StringAttribute", attribute.SchemaAttributeType())
			assert.Equal(t, property.ModelValueType(class.Name), attribute.TargetModelType)
			if values.Historical != UndefinedValueType {
				assert.Equal(t, values.Historical, buildStateUpgradeValue(property, definition).Type)
			}
		})
	}
}

func TestMigrationCollectionTargets(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "named", Input: Named, Expected: true},
		{Name: "explicit_dn", Input: Explicit, Expected: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			child := classWithProperties("fvRsCons", map[string]*Property{"target": {PropertyName: "target", AttributeName: "target_name", ValueType: String}}, nil)
			child.ResourceNameNested = "consumer_contracts"
			child.Relation.Type = testCase.Input.(RelationshipTypeEnum)
			version := 2
			class := classWithProperties("fvAEPg", map[string]*Property{}, []string{"fvRsCons"})
			class.ClassDefinition.SchemaVersion = &version
			class.StateUpgrades = map[int]StateUpgradeDefinition{1: {
				Source:   FromSDKv2,
				Children: map[string]AttributeUpgradeDefinition{"fvRsCons": {Attributes: map[string]AttributeUpgradeDefinition{"target": {AttributeName: "old_contracts", Type: SetAttribute, ElementType: StringAttribute}}}},
			}}
			store := &DataStore{Classes: map[string]Class{"fvRsCons": *child}, ctx: NewContext()}
			class.setMigrationData(store)
			assert.NoError(t, store.ctx.Diagnostics.Error())
			if assert.Len(t, class.LegacyAttributes, 1) {
				assert.True(t, class.LegacyAttributes[0].Repeated)
				assert.Equal(t, testCase.Expected, class.LegacyAttributes[0].Named)
				assert.Equal(t, "target_name", class.LegacyAttributes[0].ChildAttribute)
			}
			for _, attribute := range class.MigrationSchemas[0].Attributes {
				assert.NotEqual(t, "consumer_contracts", attribute.AttributeName, "the flat legacy collection replaces the current child set in historical state")
			}
		})
	}
}

func TestMigrationBlockMembers(t *testing.T) {
	test.InitializeTest(t)
	child := classWithProperties("fvRsNodeAtt", map[string]*Property{
		"tDn":        {PropertyName: "tDn", AttributeName: "target_dn", ValueType: String, Required: true},
		"descr":      {PropertyName: "descr", AttributeName: "description", ValueType: String, Optional: true, Computed: true},
		"annotation": {PropertyName: "annotation", AttributeName: "annotation", ValueType: String, Optional: true, Computed: true},
	}, nil)
	child.IdentifiedBy = []string{"tDn"}
	child.ResourceNameNested = "relation_to_static_leafs"
	version := 2
	class := classWithProperties("fvAEPg", map[string]*Property{}, []string{"fvRsNodeAtt"})
	class.ClassDefinition.SchemaVersion = &version
	class.StateUpgrades = map[int]StateUpgradeDefinition{1: {Source: FromSDKv2, Children: map[string]AttributeUpgradeDefinition{
		"fvRsNodeAtt": {AttributeName: "relation_fv_rs_node_att", Attributes: map[string]AttributeUpgradeDefinition{
			"tDn": {AttributeName: "node_dn"}, "descr": {AttributeName: "description"},
		}},
	}}}
	store := &DataStore{Classes: map[string]Class{"fvRsNodeAtt": *child}, ctx: NewContext()}
	class.setMigrationData(store)
	assert.NoError(t, store.ctx.Diagnostics.Error())
	if assert.Len(t, class.LegacyAttributes, 1) {
		block := class.LegacyAttributes[0]
		assert.True(t, block.Block)
		assert.Equal(t, "target_dn", block.ChildAttribute)
		assert.Equal(t, []string{"target_dn"}, block.IdentifyingAttributes)
		assert.Len(t, block.Attributes, 2)
		assert.NotContains(t, block.ModelAttributeType(), "annotation", "newer fields must not enter the old block")
		assert.Contains(t, block.ModelAttributeType(), `"description": types.StringType`)
		assert.Contains(t, block.ModelAttributeType(), `"node_dn": types.StringType`)
	}
}

func TestMigrationBlockCombinedIdentifiers(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "policy_then_filter", Input: []string{"policy", "filter"}, Expected: []string{"policy_name", "filter_type"}},
		{Name: "filter_then_policy", Input: []string{"filter", "policy"}, Expected: []string{"filter_type", "policy_name"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			child := classWithProperties("exampleRelation", map[string]*Property{
				"policy": {PropertyName: "policy", AttributeName: "policy_name", ValueType: String, Required: true},
				"filter": {PropertyName: "filter", AttributeName: "filter_type", ValueType: String, Required: true},
			}, nil)
			child.IdentifiedBy = testCase.Input.([]string)
			child.ResourceNameNested = "policies"
			version := 2
			class := classWithProperties("fvTenant", map[string]*Property{}, []string{"exampleRelation"})
			class.ClassDefinition.SchemaVersion = &version
			class.StateUpgrades = map[int]StateUpgradeDefinition{1: {Source: FromSDKv2, Children: map[string]AttributeUpgradeDefinition{
				"exampleRelation": {AttributeName: "old_policies", Attributes: map[string]AttributeUpgradeDefinition{
					"policy": {AttributeName: "old_policy"}, "filter": {AttributeName: "old_filter"},
				}},
			}}}
			store := &DataStore{Classes: map[string]Class{"exampleRelation": *child}, ctx: NewContext()}
			class.setMigrationData(store)
			assert.NoError(t, store.ctx.Diagnostics.Error())
			if assert.Len(t, class.LegacyAttributes, 1) {
				assert.Equal(t, testCase.Expected, class.LegacyAttributes[0].IdentifyingAttributes)
			}
		})
	}
}

func TestMigrationBlockNamedTarget(t *testing.T) {
	test.InitializeTest(t)
	child := classWithProperties("fvRsBDToNetflowMonitorPol", map[string]*Property{
		"tnNetflowMonitorPolName": {PropertyName: "tnNetflowMonitorPolName", AttributeName: "policy_name", ValueType: String, Required: true},
		"fltType":                 {PropertyName: "fltType", AttributeName: "filter_type", ValueType: String, Required: true},
	}, nil)
	child.IdentifiedBy = []string{"fltType", "tnNetflowMonitorPolName"}
	child.Relation = Relation{Type: Named, ToClasses: []*ClassName{testClassName("netflowMonitorPol")}}
	child.ResourceNameNested = "netflow_policies"
	version := 2
	class := classWithProperties("fvBD", map[string]*Property{}, []string{child.Name.String()})
	class.ClassDefinition.SchemaVersion = &version
	class.StateUpgrades = map[int]StateUpgradeDefinition{1: {Source: FromSDKv2, Children: map[string]AttributeUpgradeDefinition{
		child.Name.String(): {AttributeName: "old_policies", Attributes: map[string]AttributeUpgradeDefinition{
			"tnNetflowMonitorPolName": {AttributeName: "old_policy"}, "fltType": {AttributeName: "old_filter"},
		}},
	}}}
	store := &DataStore{Classes: map[string]Class{child.Name.String(): *child}, ctx: NewContext()}
	class.setMigrationData(store)
	assert.NoError(t, store.ctx.Diagnostics.Error())
	if !assert.Len(t, class.LegacyAttributes, 1) {
		return
	}
	block := class.LegacyAttributes[0]
	assert.True(t, block.Named)
	assert.Equal(t, "policy_name", block.ChildAttribute)
	assert.Equal(t, []string{"filter_type", "policy_name"}, block.IdentifyingAttributes)
	testCases := []test.TestCase{
		{Name: "old_policy", Input: true, Expected: "tDn"},
		{Name: "old_filter", Input: false, Expected: "fltType"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			for _, member := range block.Attributes {
				if member.AttributeName == testCase.Name {
					assert.Equal(t, testCase.Input, member.Named)
					assert.Equal(t, testCase.Expected, member.ResponseAttribute)
					return
				}
			}
			t.Fatal("historical block member was not normalized")
		})
	}
}

func TestMigrationChildTypeOverride(t *testing.T) {
	test.InitializeTest(t)
	child := classWithProperties("netflowRsMonitorToRecord", map[string]*Property{
		"tnNetflowRecPolName": {PropertyName: "tnNetflowRecPolName", AttributeName: "netflow_record_policy_name", ValueType: String},
	}, nil)
	child.ResourceNameNested, child.IsSingleNestedWhenDefinedAsChild = "relation_to_netflow_record", true
	class := classWithProperties("netflowMonitorPol", map[string]*Property{}, []string{child.Name.String()})
	version := 1
	class.ClassDefinition.SchemaVersion = &version
	class.StateUpgrades = map[int]StateUpgradeDefinition{0: {Source: FromFramework, Children: map[string]AttributeUpgradeDefinition{
		child.Name.String(): {Type: SetNestedAttribute},
	}}}
	store := &DataStore{Classes: map[string]Class{child.Name.String(): *child}, ctx: NewContext()}
	class.setMigrationData(store)
	assert.NoError(t, store.ctx.Diagnostics.Error())
	assert.Empty(t, class.LegacyAttributes, "a historical outer shape must not expose a second current attribute")
	if assert.Len(t, class.MigrationSchemas, 1) {
		for _, attribute := range class.MigrationSchemas[0].Attributes {
			if attribute.AttributeName == child.ResourceNameNested {
				assert.True(t, attribute.IsChildTypeChange())
				assert.False(t, attribute.IsAlias())
				assert.Equal(t, "types.Set", attribute.ModelValueType())
				assert.Equal(t, "types.Object", attribute.TargetModelType)
				assert.Contains(t, attribute.SchemaExpression, "LegacyChildSetAttribute")
				return
			}
		}
		t.Fatal("historical child missing from schema")
	}
}

func TestMigrationTypeOverrideDoesNotExposeLegacyAttribute(t *testing.T) {
	test.InitializeTest(t)
	version := 2
	class := classWithProperties("fvTenant", map[string]*Property{"controls": {PropertyName: "controls", AttributeName: "controls", ValueType: Set, Optional: true, Computed: true}}, nil)
	class.PropertiesAll = []string{"controls"}
	class.ClassDefinition.SchemaVersion = &version
	class.StateUpgrades = map[int]StateUpgradeDefinition{1: {Source: FromFramework, Attributes: map[string]AttributeUpgradeDefinition{"controls": {Type: ListAttribute}}}}
	store := &DataStore{Classes: map[string]Class{}, ctx: NewContext()}
	class.setMigrationData(store)
	assert.NoError(t, store.ctx.Diagnostics.Error())
	assert.Empty(t, class.LegacyAttributes)
	for _, attribute := range class.MigrationSchemas[0].Attributes {
		if attribute.AttributeName == "controls" {
			assert.Equal(t, ListAttribute, attribute.Type)
			assert.Equal(t, StringAttribute, attribute.ElementType)
			assert.False(t, attribute.IsAlias())
		}
	}
}

func TestMigrationOrphanStatusInheritsEarlierShape(t *testing.T) {
	test.InitializeTest(t)
	version := 4
	class := classWithProperties("fvTenant", map[string]*Property{}, nil)
	class.ClassDefinition.SchemaVersion = &version
	class.ClassDefinition.StateUpgrades = map[int]StateUpgradeDefinition{
		1: {Source: FromSDKv2, Attributes: map[string]AttributeUpgradeDefinition{"old_field": {Type: SetAttribute, ElementType: StringAttribute, Restriction: LegacyOptional, Status: Frozen}}},
		3: {Source: FromFramework, Attributes: map[string]AttributeUpgradeDefinition{"old_field": {Status: Removed}}},
	}
	store := &DataStore{Classes: map[string]Class{}, ctx: NewContext()}
	class.setStateUpgrades(store)
	class.setMigrationData(store)
	assert.NoError(t, store.ctx.Diagnostics.Error())
	if assert.Len(t, class.LegacyAttributes, 1) {
		assert.Equal(t, Removed, class.LegacyAttributes[0].Status)
	}
	for _, schema := range class.MigrationSchemas {
		for _, attribute := range schema.Attributes {
			if attribute.AttributeName == "old_field" {
				assert.Equal(t, SetAttribute, attribute.Type)
				assert.Equal(t, StringAttribute, attribute.ElementType)
				assert.Equal(t, Removed, attribute.Status)
			}
		}
	}
}

func TestRetiredStoredAliasesPreserveModernFields(t *testing.T) {
	type input struct {
		Name   string
		Source MigrationSourceEnum
	}
	testCases := []test.TestCase{
		{Name: "framework_property_name_collision", Input: input{Name: "prio", Source: FromFramework}, Expected: true},
		{Name: "framework_renamed_attribute", Input: input{Name: "old_priority", Source: FromFramework}, Expected: true},
		{Name: "sdkv2_property_name_collision", Input: input{Name: "prio", Source: FromSDKv2}, Expected: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			version := 3
			values := testCase.Input.(input)
			name := values.Name
			class := classWithProperties("fvAEPg", map[string]*Property{
				"prio": {PropertyName: "prio", AttributeName: "priority", ValueType: String, Optional: true, Computed: true},
			}, nil)
			class.PropertiesAll = []string{"prio"}
			class.ClassDefinition.SchemaVersion = &version
			class.StateUpgrades = map[int]StateUpgradeDefinition{
				1: {Attributes: map[string]AttributeUpgradeDefinition{"prio": {AttributeName: name}}},
				2: {Source: values.Source, Attributes: map[string]AttributeUpgradeDefinition{name: {AttributeName: name, Type: StringAttribute, Restriction: LegacyOptionalComputed, Status: Removed}}},
			}
			store := &DataStore{Classes: map[string]Class{}, ctx: NewContext()}
			class.setMigrationData(store)
			assert.NoError(t, store.ctx.Diagnostics.Error())
			assert.False(t, class.HasCurrentLegacyAttributes())
			if !assert.Len(t, class.MigrationSchemas, 2) {
				return
			}
			attributes := map[string]MigrationAttribute{}
			for _, attribute := range class.MigrationSchemas[1].Attributes {
				attributes[attribute.AttributeName] = attribute
			}
			if testCase.Expected.(bool) {
				assert.Contains(t, attributes, "priority")
				assert.Empty(t, attributes[name].TargetField)
				assert.Equal(t, "Prio", attributes["priority"].TargetField)
			} else {
				assert.NotContains(t, attributes, "priority")
				assert.Equal(t, "Prio", attributes[name].TargetField)
			}
			assert.Contains(t, attributes, name)
			assert.Equal(t, Removed, attributes[name].Status)
		})
	}
}

func TestRetiredStoredBlockPreservesModernChild(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "inherited_shape", Input: AttributeUpgradeDefinition{}},
		{Name: "member_rename", Input: AttributeUpgradeDefinition{Attributes: map[string]AttributeUpgradeDefinition{"tDn": {AttributeName: "old_node_dn"}}}, Expected: "cannot override members or nested children"},
		{Name: "member_type", Input: AttributeUpgradeDefinition{Attributes: map[string]AttributeUpgradeDefinition{"tDn": {Type: SetAttribute, ElementType: StringAttribute}}}, Expected: "cannot override members or nested children"},
		{Name: "grandchild", Input: AttributeUpgradeDefinition{Children: map[string]AttributeUpgradeDefinition{"tagAnnotation": {AttributeName: "old_annotations"}}}, Expected: "cannot override members or nested children"},
		{Name: "block_type", Input: AttributeUpgradeDefinition{Type: StringAttribute}, Expected: "cannot change its inherited type"},
		{Name: "element_type", Input: AttributeUpgradeDefinition{ElementType: StringAttribute}, Expected: "cannot change its inherited type"},
		{Name: "value_type", Input: AttributeUpgradeDefinition{ValueType: String}, Expected: "cannot change its inherited type"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			version := 3
			child := classWithProperties("fvRsNodeAtt", map[string]*Property{
				"tDn": {PropertyName: "tDn", AttributeName: "target_dn", ValueType: String, Required: true},
			}, nil)
			child.PropertiesAll, child.IdentifiedBy = []string{"tDn"}, []string{"tDn"}
			child.ResourceNameNested = "relation_to_static_leafs"
			class := classWithProperties("fvAEPg", map[string]*Property{}, []string{"fvRsNodeAtt"})
			class.ClassDefinition.SchemaVersion = &version
			definition := testCase.Input.(AttributeUpgradeDefinition)
			definition.AttributeName, definition.Status, definition.Restriction = "relation_fv_rs_node_att", Removed, LegacyOptional
			if definition.Type == UndefinedLegacyAttributeType {
				definition.Type = SetNestedAttribute
			}
			class.StateUpgrades = map[int]StateUpgradeDefinition{
				1: {Children: map[string]AttributeUpgradeDefinition{"fvRsNodeAtt": {
					AttributeName: "relation_fv_rs_node_att", Attributes: map[string]AttributeUpgradeDefinition{"tDn": {AttributeName: "node_dn"}},
				}}},
				2: {Source: FromFramework, Attributes: map[string]AttributeUpgradeDefinition{"relation_fv_rs_node_att": definition}},
			}
			store := &DataStore{Classes: map[string]Class{"fvRsNodeAtt": *child}, ctx: NewContext()}
			class.setMigrationData(store)
			if expected, _ := testCase.Expected.(string); expected != "" {
				assert.ErrorContains(t, store.ctx.Diagnostics.Error(), expected)
				assert.ErrorContains(t, store.ctx.Diagnostics.Error(), "schema version 2")
				return
			}
			assert.NoError(t, store.ctx.Diagnostics.Error())
			assert.False(t, class.HasCurrentLegacyAttributes())
			if !assert.Len(t, class.MigrationSchemas, 2) {
				return
			}
			attributes := map[string]MigrationAttribute{}
			for _, attribute := range class.MigrationSchemas[1].Attributes {
				attributes[attribute.AttributeName] = attribute
			}
			assert.Contains(t, attributes, "relation_to_static_leafs")
			assert.Contains(t, attributes, "relation_fv_rs_node_att")
			block := attributes["relation_fv_rs_node_att"]
			assert.True(t, block.Block)
			assert.Nil(t, block.ChildName)
			if assert.Len(t, block.Attributes, 1) {
				assert.Equal(t, "node_dn", block.Attributes[0].AttributeName)
			}
		})
	}
}

func TestMigrationRejectsGrandchildChanges(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "attribute_rename", Input: AttributeUpgradeDefinition{AttributeName: "old_value"}},
		{Name: "attribute_type", Input: AttributeUpgradeDefinition{Type: SetAttribute, ElementType: StringAttribute}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			version := 2
			child := classWithProperties("fvAp", map[string]*Property{}, nil)
			child.ResourceNameNested = "applications"
			class := classWithProperties("fvTenant", map[string]*Property{}, []string{"fvAp"})
			class.ClassDefinition.SchemaVersion = &version
			class.StateUpgrades = map[int]StateUpgradeDefinition{1: {
				Source: FromFramework,
				Children: map[string]AttributeUpgradeDefinition{"fvAp": {
					Children: map[string]AttributeUpgradeDefinition{"tagAnnotation": {
						Attributes: map[string]AttributeUpgradeDefinition{"value": testCase.Input.(AttributeUpgradeDefinition)},
					}},
				}},
			}}
			store := &DataStore{Classes: map[string]Class{"fvAp": *child}, ctx: NewContext()}
			class.setMigrationData(store)
			assert.ErrorContains(t, store.ctx.Diagnostics.Error(), "needs nested-schema migration support")
		})
	}
}

func TestMigrationNestedHistoricalShapes(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "singleton_children", Input: true},
		{Name: "repeated_children", Input: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			leaf := classWithProperties("fvAEPg", map[string]*Property{
				"name":  {PropertyName: "name", AttributeName: "name", ValueType: String, Required: true},
				"flags": {PropertyName: "flags", AttributeName: "flags", ValueType: Set, Optional: true, Computed: true},
			}, nil)
			leaf.PropertiesAll, leaf.IdentifiedBy = []string{"name", "flags"}, []string{"name"}
			leaf.ResourceNameNested = "epgs"
			leaf.IsSingleNestedWhenDefinedAsChild = testCase.Input.(bool)
			application := classWithProperties("fvAp", map[string]*Property{}, []string{"fvAEPg"})
			application.ResourceNameNested, application.IsSingleNestedWhenDefinedAsChild = "applications", testCase.Input.(bool)
			version := 5
			class := classWithProperties("fvTenant", map[string]*Property{}, []string{"fvAp"})
			class.ClassDefinition.SchemaVersion = &version
			class.StateUpgrades = map[int]StateUpgradeDefinition{
				1: {Source: FromFramework, Children: map[string]AttributeUpgradeDefinition{"fvAp": {
					Children: map[string]AttributeUpgradeDefinition{"fvAEPg": {Attributes: map[string]AttributeUpgradeDefinition{
						"name":  {AttributeName: "old_name", Status: Removed, Type: StringAttribute, Restriction: LegacyRequired},
						"flags": {Type: ListAttribute, ElementType: StringAttribute},
					}}},
				}}},
				3: {Source: FromFramework, Children: map[string]AttributeUpgradeDefinition{"fvAp": {
					Children: map[string]AttributeUpgradeDefinition{"fvAEPg": {Attributes: map[string]AttributeUpgradeDefinition{
						"name": {AttributeName: "previous_name", Status: Removed, Type: StringAttribute, Restriction: LegacyRequired},
					}}},
				}}},
			}
			store := &DataStore{Classes: map[string]Class{"fvAp": *application, "fvAEPg": *leaf}, ctx: NewContext()}
			class.setMigrationData(store)
			if !assert.NoError(t, store.ctx.Diagnostics.Error()) || !assert.Len(t, class.MigrationSchemas, 2) {
				return
			}
			for index, oldName := range []string{"old_name", "previous_name"} {
				if !assert.Len(t, class.MigrationSchemas[index].Attributes, 2) {
					continue
				}
				var parent MigrationAttribute
				for _, attribute := range class.MigrationSchemas[index].Attributes {
					if attribute.AttributeName == "applications" {
						parent = attribute
					}
				}
				assert.True(t, parent.Nested)
				assert.False(t, parent.IsAlias())
				if !assert.Len(t, parent.Attributes, 1) {
					continue
				}
				child := parent.Attributes[0]
				assert.True(t, child.Nested)
				assert.Equal(t, "epgs", child.AttributeName)
				assert.NotContains(t, child.ModelAttributeType(), `"id"`)
				assert.Contains(t, child.ModelAttributeType(), oldName)
				assert.Contains(t, child.UpgradeSchemaExpression(), `Replacement: "name"`)
				assert.Equal(t, testCase.Input.(bool), child.Type == SingleNestedAttribute)
				if index == 0 {
					assert.Contains(t, child.ModelAttributeType(), "types.ListType")
				} else {
					assert.Contains(t, child.ModelAttributeType(), "types.SetType")
				}
			}
			assert.Empty(t, store.Classes["fvAEPg"].MigrationSchemas, "normalizing a parent's historical view must not rewrite the child's own history")
		})
	}
}

func TestNestedLegacyAttributesOwnedByClass(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "shared_alias", Input: LegacyDeprecated},
		{Name: "inconsistent_parent_exposure", Input: Removed, Expected: "inconsistent current exposure"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			leaf := classWithProperties("fvAEPg", map[string]*Property{"name": {PropertyName: "name", AttributeName: "name", ValueType: String, Required: true}}, nil)
			leaf.PropertiesAll, leaf.IdentifiedBy, leaf.ResourceNameNested = []string{"name"}, []string{"name"}, "epgs"
			leaf.LegacyAttributes = []MigrationAttribute{{AttributeName: "root_only", FieldName: "LegacyRootOnly", Type: StringAttribute, Status: Frozen}}
			store := &DataStore{Classes: map[string]Class{"fvAEPg": *leaf}, ctx: NewContext()}
			for index, name := range []string{"fvTenant", "fvBD"} {
				class := classWithProperties(name, map[string]*Property{}, []string{"fvAEPg"})
				version, savedVersion := 4+index*5, 1+index*6
				class.ClassDefinition.SchemaVersion = &version
				status := LegacyDeprecated
				if index == 1 {
					status = testCase.Input.(LegacyStatusEnum)
				}
				class.StateUpgrades = map[int]StateUpgradeDefinition{savedVersion: {
					Source: FromFramework, Children: map[string]AttributeUpgradeDefinition{"fvAEPg": {Attributes: map[string]AttributeUpgradeDefinition{
						"name": {AttributeName: "old_name", Type: StringAttribute, Restriction: LegacyRequired, Status: status},
					}}},
				}}
				class.setMigrationData(store)
				store.Classes[name] = *class
			}
			store.setNestedLegacyAttributes()
			if message, expectedError := testCase.Expected.(string); expectedError {
				assert.ErrorContains(t, store.ctx.Diagnostics.Error(), message)
				return
			}
			assert.NoError(t, store.ctx.Diagnostics.Error())
			class := store.Classes["fvAEPg"]
			if assert.Len(t, class.NestedLegacyAttributes, 1) {
				assert.Equal(t, "old_name", class.NestedLegacyAttributes[0].AttributeName)
			}
			assert.Len(t, class.RootLegacyAttributes(), 1)
			assert.Equal(t, "root_only", class.RootLegacyAttributes()[0].AttributeName)
			assert.True(t, store.Classes["fvTenant"].HasNestedLegacyChanges)
			assert.True(t, store.Classes["fvBD"].HasNestedLegacyChanges)
		})
	}
}

func TestNestedLegacyAttributesUseCurrentRestrictions(t *testing.T) {
	type input struct {
		Property Property
		History  []LegacyRestrictionEnum
	}
	testCases := []test.TestCase{
		{Name: "required_current_property", Input: input{Property: Property{Required: true}, History: []LegacyRestrictionEnum{LegacyReadOnly}}},
		{Name: "optional_current_property", Input: input{Property: Property{Optional: true, Computed: true}, History: []LegacyRestrictionEnum{LegacyReadOnly}}},
		{Name: "read_only_current_property", Input: input{Property: Property{ReadOnly: true, Computed: true}, History: []LegacyRestrictionEnum{LegacyOptionalComputed}}},
		{Name: "different_parent_histories", Input: input{Property: Property{Optional: true, Computed: true}, History: []LegacyRestrictionEnum{LegacyRequired, LegacyOptionalComputed, LegacyReadOnly}}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			values := testCase.Input.(input)
			property := values.Property
			property.PropertyName, property.AttributeName, property.ValueType = "value", "value", String
			child := classWithProperties("tagAnnotation", map[string]*Property{"value": &property}, nil)
			child.PropertiesAll, child.ResourceNameNested = []string{"value"}, "annotations"
			store := &DataStore{ctx: NewContext(), Classes: map[string]Class{"tagAnnotation": *child}}
			parents := []string{"fvTenant", "fvBD", "fvAp"}
			for index, restriction := range values.History {
				parent := classWithProperties(parents[index], map[string]*Property{}, []string{"tagAnnotation"})
				version, savedVersion := 4+index*5, 1+index*6
				parent.ClassDefinition.SchemaVersion = &version
				parent.StateUpgrades = map[int]StateUpgradeDefinition{savedVersion: {
					Source: FromFramework, Children: map[string]AttributeUpgradeDefinition{"tagAnnotation": {Attributes: map[string]AttributeUpgradeDefinition{
						"value": {AttributeName: "old_value", Restriction: restriction, Status: LegacyDeprecated},
					}}},
				}}
				parent.setMigrationData(store)
				store.Classes[parents[index]] = *parent
			}
			store.setNestedLegacyAttributes()
			if !assert.NoError(t, store.ctx.Diagnostics.Error()) {
				return
			}
			current := store.Classes["tagAnnotation"]
			if !assert.Len(t, current.NestedLegacyAttributes, 1) {
				return
			}
			alias := current.NestedLegacyAttributes[0]
			assert.False(t, alias.Required)
			assert.Equal(t, !property.ReadOnly, alias.Optional)
			assert.True(t, alias.Computed)
			assert.Equal(t, property, *current.Properties["value"])
			// Collecting the shared alias must not rewrite any parent's historical decoding flags.
			for index, restriction := range values.History {
				parent := store.Classes[parents[index]]
				assert.Equal(t, restriction, parent.StateUpgrades[1+index*6].Children["tagAnnotation"].Attributes["value"].Restriction)
				if !assert.Len(t, parent.MigrationSchemas, 1) {
					continue
				}
				found := false
				for _, attribute := range parent.MigrationSchemas[0].Attributes {
					if attribute.AttributeName != "annotations" {
						continue
					}
					for _, member := range attribute.Attributes {
						if member.AttributeName == "old_value" {
							found = true
							assert.Equal(t, restriction == LegacyRequired, member.Required)
							assert.Equal(t, restriction == LegacyOptionalComputed, member.Optional)
							assert.Equal(t, restriction != LegacyRequired, member.Computed)
						}
					}
				}
				assert.True(t, found, "historical child must retain its old_value field")
			}
		})
	}
}

func TestMigrationUnsupportedPrimitiveConversions(t *testing.T) {
	type input struct {
		Current ValueTypeEnum
		Prior   LegacyAttributeTypeEnum
	}
	testCases := []test.TestCase{
		{Name: "string_to_set", Input: input{Current: Set, Prior: StringAttribute}, Expected: true},
		{Name: "set_to_string", Input: input{Current: String, Prior: SetAttribute}, Expected: true},
		{Name: "list_to_string", Input: input{Current: String, Prior: ListAttribute}, Expected: true},
		{Name: "list_to_set", Input: input{Current: Set, Prior: ListAttribute}, Expected: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			values := testCase.Input.(input)
			class := classWithProperties("fvAp", map[string]*Property{"value": {PropertyName: "value", AttributeName: "value", ValueType: values.Current, Optional: true, Computed: true}}, nil)
			version := 2
			class.PropertiesAll, class.ClassDefinition.SchemaVersion = []string{"value"}, &version
			class.StateUpgrades = map[int]StateUpgradeDefinition{1: {Source: FromFramework, Attributes: map[string]AttributeUpgradeDefinition{"value": {Type: values.Prior, ElementType: StringAttribute}}}}
			store := &DataStore{Classes: map[string]Class{}, ctx: NewContext()}
			class.setMigrationData(store)
			if testCase.Expected.(bool) {
				assert.ErrorContains(t, store.ctx.Diagnostics.Error(), "explicit conversion")
			} else {
				assert.NoError(t, store.ctx.Diagnostics.Error())
			}
		})
	}
}

func TestLegacyLookupAttributes(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "identifier", Input: MigrationAttribute{AttributeName: "old_name", ChildAttribute: "name", Type: StringAttribute, Status: LegacyDeprecated}, Expected: true},
		{Name: "parent", Input: MigrationAttribute{AttributeName: "tenant_dn", ChildAttribute: "parent_dn", Type: StringAttribute, Status: LegacyDeprecated}, Expected: true},
		{Name: "ordinary_output", Input: MigrationAttribute{AttributeName: "descr", ChildAttribute: "description", Type: StringAttribute, Status: LegacyDeprecated}, Expected: false},
		{Name: "frozen", Input: MigrationAttribute{AttributeName: "old_name", ChildAttribute: "name", Type: StringAttribute, Status: Frozen}, Expected: false},
		{Name: "removed", Input: MigrationAttribute{AttributeName: "old_name", ChildAttribute: "name", Type: StringAttribute, Status: Removed}, Expected: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			class, _ := resourceScenarioFixture()
			class.LegacyAttributes = []MigrationAttribute{testCase.Input.(MigrationAttribute)}
			if actual := len(class.LegacyLookupAttributes()) > 0; actual != testCase.Expected.(bool) {
				t.Fatalf("lookup alias eligibility: %t", actual)
			}
		})
	}
}
