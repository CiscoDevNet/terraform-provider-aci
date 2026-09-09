package data

import (
	"fmt"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

// TestParseGlobalMetaDefinition_HappyPath exercises the strict YAML unmarshal
// for every field on GlobalMetaDefinition.
func TestParseGlobalMetaDefinition_HappyPath(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	yamlBytes := []byte(`
always_include_as_child:
  - fvRsBd
attribute_name_overrides:
  descr: description
exclude_parents:
  - polUni
exclude_properties:
  - childAction
documentation_label_overrides:
  Bgp: BGP
property_documentation_overrides:
  descr: The description of the %s object.
  nameAlias: The name alias of the %s object.
`)

	parsedDefinition, err := parseGlobalMetaDefinition(yamlBytes)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, []string{"fvRsBd"}, parsedDefinition.AlwaysIncludeAsChild)
	assert.Equal(t, "description", parsedDefinition.AttributeNameOverrides["descr"])
	assert.Equal(t, []string{"polUni"}, parsedDefinition.ExcludeParents)
	assert.Equal(t, []string{"childAction"}, parsedDefinition.ExcludeProperties)
	assert.Equal(t, "BGP", parsedDefinition.DocumentationLabelOverrides["Bgp"])
	assert.Equal(t, "The description of the %s object.", parsedDefinition.PropertyDocumentationOverrides["descr"])
	assert.Equal(t, "The name alias of the %s object.", parsedDefinition.PropertyDocumentationOverrides["nameAlias"])
}

// TestParseGlobalMetaDefinition_UnknownField verifies UnmarshalStrict rejects
// unknown YAML keys so stale/renamed fields fail loudly during generation.
func TestParseGlobalMetaDefinition_UnknownField(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	yamlBytes := []byte(`
always_include_as_child:
  - fvRsBd
this_field_does_not_exist: oops
`)

	_, err := parseGlobalMetaDefinition(yamlBytes)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "this_field_does_not_exist")
}

// TestParseGlobalMetaDefinition_MalformedYAML verifies a syntactically invalid
// YAML payload is rejected by the parser.
func TestParseGlobalMetaDefinition_MalformedYAML(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	_, err := parseGlobalMetaDefinition([]byte("always_include_as_child: [unterminated"))
	assert.Error(t, err)
}

// TestParseClassDefinition_HappyPath exercises a moderately complete class
// definition including a property with test_config and a class-level test_config.
func TestParseClassDefinition_HappyPath(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	yamlBytes := []byte(`
resource_name: tenant
resource_name_nested: tenants
identified_by:
  - name
artifacts:
  - resource
  - datasource
example_files:
  - certificate.pem
parent_dn_variants:
  - parent_class: fvTenant
    rn_prepend: certstore
    wrapper_class: cloudCertStore
    test_platform: apic
  - parent_class: pkiEp
    rn_prepend: userext/pkiext
    test_platform: cloud
documentation:
  label: Tenant
  description: A tenant.
  example_parent_classes:
    - fvTenant
    - fvAEPg
relation_info:
  type: named
  to_classes:
    - vz:BrCP
properties:
  name:
    restriction: required
    reference:
      classes: [fvTenant]
      example_class: fvTenant
      example_attribute: name
    example_config:
      minimum:
        - value: aci_tenant.example.id
          value_type: reference
          parent_class: fvTenant
      full: []
    test_config:
      create:
        - config_value: test_tenant
          config_include: true
  descr:
    test_config:
      ignore_in_test: true
      update:
        - config_value: updated
          assert_value: updated
  arpLearning:
    value_type: vmm_arp_learning
test_config:
  replace_auto_resolved: true
  ignore_tests:
    - child
  ignore_import_state_verify: true
  dependencies:
    - class_name: fvTenant
      reference: aci_tenant.test.id
      reference_type: resource
      role: parent
      target_classes:
        - fvBD
`)

	parsedDefinition, err := parseClassDefinition(yamlBytes)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, "tenant", parsedDefinition.ResourceName)
	assert.Equal(t, "tenants", parsedDefinition.ResourceNameNested)
	assert.Equal(t, []string{"name"}, parsedDefinition.IdentifiedBy)
	assert.Equal(t, []ArtifactEnum{ResourceArtifact, DatasourceArtifact}, parsedDefinition.Artifacts)
	assert.Equal(t, []string{"certificate.pem"}, parsedDefinition.ExampleFiles)
	assert.Len(t, parsedDefinition.ParentDnVariants, 2)
	assert.Equal(t, "fvTenant", parsedDefinition.ParentDnVariants[0].ParentClass)
	assert.Equal(t, "certstore", parsedDefinition.ParentDnVariants[0].RnPrepend)
	assert.Equal(t, "cloudCertStore", parsedDefinition.ParentDnVariants[0].WrapperClass)
	assert.Equal(t, Apic, parsedDefinition.ParentDnVariants[0].TestPlatform)
	assert.Equal(t, "pkiEp", parsedDefinition.ParentDnVariants[1].ParentClass)
	assert.Equal(t, "userext/pkiext", parsedDefinition.ParentDnVariants[1].RnPrepend)
	assert.Equal(t, "", parsedDefinition.ParentDnVariants[1].WrapperClass)
	assert.Equal(t, Cloud, parsedDefinition.ParentDnVariants[1].TestPlatform)
	assert.Equal(t, "Tenant", parsedDefinition.Documentation.Label)
	assert.Equal(t, []string{"fvTenant", "fvAEPg"}, parsedDefinition.Documentation.ExampleParentClasses)
	assert.Equal(t, Named, parsedDefinition.RelationInfo.Type)
	assert.Equal(t, []string{"vz:BrCP"}, parsedDefinition.RelationInfo.ToClasses)

	nameProp, ok := parsedDefinition.Properties["name"]
	assert.True(t, ok)
	assert.Equal(t, Required, nameProp.Restriction)
	assert.True(t, nameProp.TestConfig.Create.Defined)
	assert.Len(t, nameProp.TestConfig.Create.Entries, 1)
	assert.Equal(t, "test_tenant", nameProp.TestConfig.Create.Entries[0].ConfigValue)
	assert.NotNil(t, nameProp.TestConfig.Create.Entries[0].ConfigInclude)
	assert.True(t, *nameProp.TestConfig.Create.Entries[0].ConfigInclude)
	assert.True(t, nameProp.ExampleConfig.Minimum.Defined)
	assert.Len(t, nameProp.ExampleConfig.Minimum.Entries, 1)
	assert.Equal(t, "aci_tenant.example.id", nameProp.ExampleConfig.Minimum.Entries[0].Value)
	assert.Equal(t, ReferenceValue, nameProp.ExampleConfig.Minimum.Entries[0].ValueType)
	assert.Equal(t, "fvTenant", nameProp.ExampleConfig.Minimum.Entries[0].ParentClass)
	assert.Equal(t, "name", nameProp.Reference.ExampleAttributeName)
	assert.True(t, nameProp.ExampleConfig.Full.Defined)
	assert.Empty(t, nameProp.ExampleConfig.Full.Entries)

	descrProp, ok := parsedDefinition.Properties["descr"]
	assert.True(t, ok)
	assert.True(t, descrProp.TestConfig.IgnoreInTest)
	assert.True(t, descrProp.TestConfig.Update.Defined)
	assert.Len(t, descrProp.TestConfig.Update.Entries, 1)
	assert.Equal(t, "updated", descrProp.TestConfig.Update.Entries[0].AssertValue)

	arpLearningProp, ok := parsedDefinition.Properties["arpLearning"]
	assert.True(t, ok)
	assert.Equal(t, VMMArpLearning, arpLearningProp.ValueType)

	assert.True(t, parsedDefinition.TestConfig.ReplaceAutoResolved)
	assert.Equal(t, []IgnoreTestEnum{ChildIgnoreTest}, parsedDefinition.TestConfig.IgnoreTests)
	assert.True(t, parsedDefinition.TestConfig.IgnoreImportStateVerify)
	assert.Len(t, parsedDefinition.TestConfig.Dependencies, 1)
	assert.Equal(t, "fvTenant", parsedDefinition.TestConfig.Dependencies[0].ClassName)
	assert.Equal(t, ResourceReference, parsedDefinition.TestConfig.Dependencies[0].ReferenceType)
	assert.Equal(t, Parent, parsedDefinition.TestConfig.Dependencies[0].Role)
	assert.Equal(t, []string{"fvBD"}, parsedDefinition.TestConfig.Dependencies[0].TargetClasses)
}

// TestParseClassDefinition_UnknownField at the top level.
func TestParseClassDefinition_UnknownField(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	yamlBytes := []byte(`
resource_name: tenant
bogus_top_level_field: 42
`)

	_, err := parseClassDefinition(yamlBytes)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bogus_top_level_field")
}

// TestParseClassDefinition_UnknownFieldNestedProperty verifies UnmarshalStrict
// propagates into nested PropertyDefinition entries.
func TestParseClassDefinition_UnknownFieldNestedProperty(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	yamlBytes := []byte(`
properties:
  name:
    restriction: required
    unknown_property_field: x
`)

	_, err := parseClassDefinition(yamlBytes)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown_property_field")
}

// TestParseClassDefinition_UnknownFieldNestedTestConfig verifies UnmarshalStrict
// propagates into deeply nested test_config dependency entries.
func TestParseClassDefinition_UnknownFieldNestedTestConfig(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	yamlBytes := []byte(`
test_config:
  dependencies:
    - class_name: fvTenant
      reference: aci_tenant.test.id
      role: parent
      unknown_dep_field: x
`)

	_, err := parseClassDefinition(yamlBytes)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown_dep_field")
}

// TestParseClassDefinition_InvalidEnumValues exercises every enum surface
// reachable through YAML to confirm typos are rejected at parse time.
func TestParseClassDefinition_InvalidEnumValues(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "invalid_restriction",
			Input: `
properties:
  name:
    restriction: not_a_restriction
`,
			Expected: "unknown restriction",
		},
		{
			Name: "invalid_value_type",
			Input: `
properties:
  name:
    value_type: not_a_type
`,
			Expected: "unknown value_type",
		},
		{
			Name: "invalid_relationship_type",
			Input: `
relation_info:
  type: not_a_type
`,
			Expected: "unknown relationship type",
		},
		{
			Name: "invalid_reference_type",
			Input: `
test_config:
  dependencies:
    - class_name: fvTenant
      reference: aci_tenant.test.id
      reference_type: not_a_reference
      role: parent
`,
			Expected: "unknown reference_type",
		},
		{
			Name: "invalid_role",
			Input: `
test_config:
  dependencies:
    - class_name: fvTenant
      reference: aci_tenant.test.id
      role: not_a_role
`,
			Expected: "unknown role",
		},
		{
			Name: "invalid_test_value_render_type",
			Input: `
properties:
  name:
    test_config:
      create:
        - config_value: x
          value_type: not_a_render
`,
			Expected: "unknown value_type",
		},
		{
			Name: "invalid_example_value_render_type",
			Input: `
properties:
  name:
    example_config:
      full:
        - value: x
          value_type: not_a_render
`,
			Expected: "unknown value_type",
		},
		{
			Name: "empty_restriction_is_typo",
			Input: `
properties:
  name:
    restriction: ""
`,
			Expected: "unknown restriction",
		},
		{
			Name: "invalid_artifact",
			Input: `
artifacts:
  - not_a_real_artifact
`,
			Expected: "unknown artifact",
		},
		{
			Name: "empty_artifact_is_typo",
			Input: `
artifacts:
  - ""
`,
			Expected: "unknown artifact",
		},
		{
			Name: "invalid_ignore_test",
			Input: `
test_config:
  ignore_tests:
    - not_a_real_bucket
`,
			Expected: "unknown ignore_tests entry",
		},
		{
			Name: "invalid_parent_dn_variant_test_platform",
			Input: `
parent_dn_variants:
  - parent_class: fvTenant
    test_platform: not_a_real_platform
`,
			Expected: "unknown platform_type",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			_, err := parseClassDefinition([]byte(testCase.Input.(string)))
			assert.Error(t, err, testCase.Name)
			assert.Contains(t, err.Error(), testCase.Expected.(string), testCase.Name)
		})
	}
}

// TestParseClassDefinition_ConfigIncludePointerSemantics verifies that the
// *bool pointer for config_include distinguishes "unset" (nil) from
// "explicit false" — a guarantee documented in test_configuration.md §2.6.
func TestParseClassDefinition_ConfigIncludePointerSemantics(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	yamlBytes := []byte(`
properties:
  name:
    test_config:
      create:
        - config_value: a
        - config_value: b
          config_include: true
        - config_value: c
          config_include: false
`)

	parsedDefinition, err := parseClassDefinition(yamlBytes)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	entries := parsedDefinition.Properties["name"].TestConfig.Create.Entries
	assert.Len(t, entries, 3)
	assert.Nil(t, entries[0].ConfigInclude, "omitted config_include must decode to nil")
	if assert.NotNil(t, entries[1].ConfigInclude) {
		assert.True(t, *entries[1].ConfigInclude)
	}
	if assert.NotNil(t, entries[2].ConfigInclude) {
		assert.False(t, *entries[2].ConfigInclude, "explicit false must survive as *bool=false, not nil")
	}
}

func TestParseClassDefinition_TestValueScenarioPresence(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	parsedDefinition, err := parseClassDefinition([]byte(`
properties:
  tags:
    test_config:
      update: []
  name:
    test_config: {}
`))
	assert.NoError(t, err, test.MessageUnexpectedError(err))

	tags := parsedDefinition.Properties["tags"].TestConfig
	assert.True(t, tags.Update.Defined, "an explicit empty list must remain defined")
	assert.Empty(t, tags.Update.Entries)
	assert.False(t, tags.Create.Defined, "an omitted scenario must remain undefined")

	name := parsedDefinition.Properties["name"].TestConfig
	assert.False(t, name.Update.Defined)
	assert.Nil(t, name.Update.Entries)
}

// TestParseClassDefinition_RequiresReplacePointerSemantics mirrors the
// ConfigInclude check for the top-level *bool requires_replace override.
func TestParseClassDefinition_RequiresReplacePointerSemantics(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "omitted", Input: "properties:\n  name:\n    restriction: required\n", Expected: (*bool)(nil)},
		{Name: "true", Input: "properties:\n  name:\n    requires_replace: true\n", Expected: boolPtr(true)},
		{Name: "false", Input: "properties:\n  name:\n    requires_replace: false\n", Expected: boolPtr(false)},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			parsedDefinition, err := parseClassDefinition([]byte(testCase.Input.(string)))
			assert.NoError(t, err, test.MessageUnexpectedError(err))
			expected := testCase.Expected.(*bool)
			actual := parsedDefinition.Properties["name"].RequiresReplace
			if expected == nil {
				assert.Nil(t, actual, testCase.Name)
			} else {
				if assert.NotNil(t, actual, testCase.Name) {
					assert.Equal(t, *expected, *actual, testCase.Name)
				}
			}
		})
	}
}

// TestParseClassDefinition_EmptyYAML verifies that an empty payload decodes
// to a zero-value ClassDefinition without error.
func TestParseClassDefinition_EmptyYAML(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	parsedDefinition, err := parseClassDefinition(nil)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, ClassDefinition{}, parsedDefinition)
}

func boolPtr(b bool) *bool { return &b }

// TestParseClassDefinition_MigrationSource exercises the top-level migration_source
// field in isolation: defaulted to UndefinedMigrationSource when omitted, and
// FromSDKv2 when explicitly set. Kept separate from the state_upgrades tests
// since migration_source is an independent top-level field decoded by
// parseClassDefinition regardless of whether state_upgrades is present.
func TestParseClassDefinition_MigrationSource(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "omitted", Input: "resource_name: tenant\n", Expected: UndefinedMigrationSource},
		{Name: "from_sdkv2", Input: "migration_source: from_sdkv2\n", Expected: FromSDKv2},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			parsedDefinition, err := parseClassDefinition([]byte(testCase.Input.(string)))
			assert.NoError(t, err, test.MessageUnexpectedError(err))
			assert.Equal(t, testCase.Expected.(MigrationSourceEnum), parsedDefinition.MigrationSource, testCase.Name)
		})
	}
}

// TestParseClassDefinition_StateUpgrades is the canonical happy-path decode for
// state_upgrades. It mixes every scenario the parser needs to handle into one
// realistic fvBD (Bridge Domain) example:
//
//   - Two state_upgrades entries at distinct prior_schema_versions (1 and 3),
//     proving slice decode at non-zero unique values.
//   - An attribute with only legacy_attribute set, proving the trio of Undefined
//     defaults (LegacyType, LegacyRestriction, LegacyStatus) for omitted fields.
//   - An attribute with all four fields explicitly set, including legacy_status=frozen.
//   - An attribute with legacy_status=removed (the case where legacy_type and
//     legacy_restriction must be supplied because there is no current attribute
//     to inherit from).
//   - A child block (fvSubnet) carrying its own legacy_attribute (block rename)
//     alongside inner attributes AND its own children (fvRsBDSubnetToOut),
//     proving 3-level recursion: children -> children -> attributes.
//   - A child block (fvRsCtx) without a block-level legacy_attribute, proving
//     the empty-string decode for an unchanged block whose inner attributes
//     are the only diff.
//   - An attribute with legacy_status explicitly set to functioning, proving
//     that the explicit value round-trips identically to the default.
func TestParseClassDefinition_StateUpgrades(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	yamlBytes := []byte(`
state_upgrades:
  - prior_schema_version: 1
    attributes:
      mac:
        legacy_attribute: mac_address
      arpFlood:
        legacy_attribute: arp_flooding
        legacy_type: bool_attribute
        legacy_restriction: optional
        legacy_status: frozen
      multiDstPktAct:
        legacy_attribute: multi_dst_pkt_act
        legacy_type: string_attribute
        legacy_restriction: optional
        legacy_status: removed
    children:
      fvSubnet:
        legacy_attribute: subnets
        attributes:
          ip:
            legacy_attribute: subnet_ip
        children:
          fvRsBDSubnetToOut:
            attributes:
              tnL3extOutName:
                legacy_attribute: l3out_name
                legacy_type: string_attribute
      fvRsCtx:
        attributes:
          tnFvCtxName:
            legacy_attribute: vrf_name
            legacy_status: functioning
  - prior_schema_version: 3
    attributes:
      descr:
        legacy_attribute: description
`)

	parsedDefinition, err := parseClassDefinition(yamlBytes)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, UndefinedMigrationSource, parsedDefinition.MigrationSource, "migration_source is decoded by a dedicated test; absent here must default to UndefinedMigrationSource")

	if !assert.Len(t, parsedDefinition.StateUpgrades, 2, "two state_upgrades entries must round-trip") {
		return
	}

	v1 := parsedDefinition.StateUpgrades[0]
	assert.Equal(t, 1, v1.PriorSchemaVersion, "non-zero prior_schema_version must round-trip; a 0 here would mask a silently-dropped field")
	assert.Len(t, v1.Attributes, 3)

	// Defaults proof: only legacy_attribute set, the trio of Undefined defaults applies.
	mac := v1.Attributes["mac"]
	assert.Equal(t, "mac_address", mac.LegacyAttribute)
	assert.Equal(t, UndefinedLegacyAttributeType, mac.LegacyType)
	assert.Equal(t, UndefinedRestriction, mac.LegacyRestriction)
	assert.Equal(t, Functioning, mac.LegacyStatus, "Functioning is the LegacyStatus zero value reached by omitting the YAML field")

	// All four fields explicitly set, legacy_status=frozen.
	arp := v1.Attributes["arpFlood"]
	assert.Equal(t, "arp_flooding", arp.LegacyAttribute)
	assert.Equal(t, BoolAttribute, arp.LegacyType)
	assert.Equal(t, Optional, arp.LegacyRestriction)
	assert.Equal(t, Frozen, arp.LegacyStatus)

	// legacy_status=removed: the attribute is gone from the current schema, so
	// legacy_type and legacy_restriction must be supplied here.
	rem := v1.Attributes["multiDstPktAct"]
	assert.Equal(t, "multi_dst_pkt_act", rem.LegacyAttribute)
	assert.Equal(t, StringAttribute, rem.LegacyType)
	assert.Equal(t, Optional, rem.LegacyRestriction)
	assert.Equal(t, Removed, rem.LegacyStatus)

	// Block with its own legacy_attribute (rename) + inner attributes + nested children.
	subnet, ok := v1.Children["fvSubnet"]
	if assert.True(t, ok, "outer child must decode") {
		assert.Equal(t, "subnets", subnet.LegacyAttribute, "block-level legacy_attribute records a block rename")
		assert.Equal(t, "subnet_ip", subnet.Attributes["ip"].LegacyAttribute)

		rsToOut, ok := subnet.Children["fvRsBDSubnetToOut"]
		if assert.True(t, ok, "grandchild block must decode (children -> children)") {
			inner, ok := rsToOut.Attributes["tnL3extOutName"]
			if assert.True(t, ok, "3-level recursion: children -> children -> attributes") {
				assert.Equal(t, "l3out_name", inner.LegacyAttribute)
				assert.Equal(t, StringAttribute, inner.LegacyType)
			}
		}
	}

	// Block without its own legacy_attribute (block unchanged, only inner attribute changed).
	rsCtx, ok := v1.Children["fvRsCtx"]
	if assert.True(t, ok) {
		assert.Equal(t, "", rsCtx.LegacyAttribute, "absent block-level legacy_attribute decodes to empty")
		inner := rsCtx.Attributes["tnFvCtxName"]
		assert.Equal(t, "vrf_name", inner.LegacyAttribute)
		assert.Equal(t, Functioning, inner.LegacyStatus, "explicit 'functioning' must round-trip identically to the omitted default")
	}

	v3 := parsedDefinition.StateUpgrades[1]
	assert.Equal(t, 3, v3.PriorSchemaVersion, "second entry's prior_schema_version must round-trip independently")
	assert.Equal(t, "description", v3.Attributes["descr"].LegacyAttribute)
}

// TestParseClassDefinition_StateUpgradesInvalidEnumValues confirms strict
// decoding rejects unknown values for each enum reachable through state_upgrades.
func TestParseClassDefinition_StateUpgradesInvalidEnumValues(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "invalid_legacy_type",
			Input: `
state_upgrades:
  - prior_schema_version: 0
    attributes:
      name:
        legacy_type: not_a_type
`,
			Expected: "unknown legacy_type",
		},
		{
			Name: "invalid_legacy_status",
			Input: `
state_upgrades:
  - prior_schema_version: 0
    attributes:
      name:
        legacy_status: weird
`,
			Expected: "unknown legacy_status",
		},
		{
			Name: "invalid_legacy_restriction",
			Input: `
state_upgrades:
  - prior_schema_version: 0
    attributes:
      name:
        legacy_restriction: weird
`,
			Expected: "unknown restriction",
		},
		{
			Name:     "invalid_migration_source",
			Input:    "migration_source: bogus\n",
			Expected: "unknown migration_source",
		},
		{
			Name: "empty_legacy_status_is_typo",
			Input: `
state_upgrades:
  - prior_schema_version: 0
    attributes:
      name:
        legacy_status: ""
`,
			Expected: "unknown legacy_status",
		},
		{
			Name: "empty_legacy_restriction_is_typo",
			Input: `
state_upgrades:
  - prior_schema_version: 0
    attributes:
      name:
        legacy_restriction: ""
`,
			Expected: "unknown restriction",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			_, err := parseClassDefinition([]byte(testCase.Input.(string)))
			assert.Error(t, err, testCase.Name)
			assert.Contains(t, err.Error(), testCase.Expected.(string), testCase.Name)
		})
	}
}

// TestParseClassDefinition_StateUpgradesUnknownField verifies UnmarshalStrict
// propagates into AttributeUpgradeDefinition nodes at any nesting depth.
func TestParseClassDefinition_StateUpgradesUnknownField(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "unknown_field_on_entry",
			Input: `
state_upgrades:
  - prior_schema_version: 0
    bogus_upgrade_field: true
`,
			Expected: "bogus_upgrade_field",
		},
		{
			Name: "unknown_field_on_attribute_node",
			Input: `
state_upgrades:
  - prior_schema_version: 0
    attributes:
      name:
        bogus_attribute_field: true
`,
			Expected: "bogus_attribute_field",
		},
		{
			Name: "unknown_field_on_inner_child_node",
			Input: `
state_upgrades:
  - prior_schema_version: 0
    children:
      fvRsBd:
        attributes:
          name:
            bogus_inner_field: true
`,
			Expected: "bogus_inner_field",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			_, err := parseClassDefinition([]byte(testCase.Input.(string)))
			assert.Error(t, err, testCase.Name)
			assert.Contains(t, err.Error(), testCase.Expected.(string), testCase.Name)
		})
	}
}

type attributeUpgradeValidationExpected struct {
	ErrorCount int
	Substrings []string
}

// TestAttributeUpgradeDefinition_Validate covers the attributes-bucket validator.
// Only Removed nodes are checked; the three legacy_* fields are independently
// required so a fully-empty Removed node aggregates three diagnostics in one pass.
func TestAttributeUpgradeDefinition_Validate(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	const path = `Class 'fvCtx': prior_schema_version 0: attributes["name"]`

	testCases := []test.TestCase{
		{
			Name:     "test_non_removed_status_skips_all_checks",
			Input:    AttributeUpgradeDefinition{LegacyStatus: Functioning},
			Expected: attributeUpgradeValidationExpected{},
		},
		{
			Name: "test_removed_with_all_legacy_fields_set_passes",
			Input: AttributeUpgradeDefinition{
				LegacyAttribute:   "old_name",
				LegacyType:        StringAttribute,
				LegacyRestriction: Optional,
				LegacyStatus:      Removed,
			},
			Expected: attributeUpgradeValidationExpected{},
		},
		{
			Name: "test_removed_missing_legacy_attribute_only",
			Input: AttributeUpgradeDefinition{
				LegacyType:        StringAttribute,
				LegacyRestriction: Optional,
				LegacyStatus:      Removed,
			},
			Expected: attributeUpgradeValidationExpected{ErrorCount: 1, Substrings: []string{"requires legacy_attribute"}},
		},
		{
			Name: "test_removed_missing_legacy_type_only",
			Input: AttributeUpgradeDefinition{
				LegacyAttribute:   "old_name",
				LegacyRestriction: Optional,
				LegacyStatus:      Removed,
			},
			Expected: attributeUpgradeValidationExpected{ErrorCount: 1, Substrings: []string{"requires legacy_type"}},
		},
		{
			Name: "test_removed_missing_legacy_restriction_only",
			Input: AttributeUpgradeDefinition{
				LegacyAttribute: "old_name",
				LegacyType:      StringAttribute,
				LegacyStatus:    Removed,
			},
			Expected: attributeUpgradeValidationExpected{ErrorCount: 1, Substrings: []string{"requires legacy_restriction"}},
		},
		{
			Name:  "test_removed_all_fields_missing_aggregates_three_diagnostics",
			Input: AttributeUpgradeDefinition{LegacyStatus: Removed},
			Expected: attributeUpgradeValidationExpected{
				ErrorCount: 3,
				Substrings: []string{"requires legacy_attribute", "requires legacy_type", "requires legacy_restriction"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			expected := testCase.Expected.(attributeUpgradeValidationExpected)
			node := testCase.Input.(AttributeUpgradeDefinition)
			ctx := NewContext()
			node.validate(ctx, path)
			err := ctx.Diagnostics.Error()
			if expected.ErrorCount == 0 {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
			msg := err.Error()
			assert.Contains(t, msg, fmt.Sprintf("encountered %d error(s)", expected.ErrorCount))
			assert.Contains(t, msg, path)
			for _, substr := range expected.Substrings {
				assert.Contains(t, msg, substr)
			}
		})
	}
}

// TestAttributeUpgradeDefinition_ValidateChild covers the children-bucket validator.
// In addition to the Removed-only legacy_attribute check, it enforces the
// scalar-wrap shape (the only allowed way for a child entry to omit both
// legacy_attribute and legacy_type) and recurses into inner attributes/children.
func TestAttributeUpgradeDefinition_ValidateChild(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	const path = `Class 'fvCtx': prior_schema_version 0: children["fvRsBd"]`

	testCases := []test.TestCase{
		{
			Name:     "test_block_rename_only",
			Input:    AttributeUpgradeDefinition{LegacyAttribute: "old_block"},
			Expected: attributeUpgradeValidationExpected{},
		},
		{
			Name: "test_scalar_wrap_inner_attribute_carries_legacy",
			Input: AttributeUpgradeDefinition{
				Attributes: map[string]AttributeUpgradeDefinition{
					"value": {LegacyAttribute: "old_flat"},
				},
			},
			Expected: attributeUpgradeValidationExpected{},
		},
		{
			Name: "test_inner_attributes_only_no_legacy_on_block_or_inner",
			Input: AttributeUpgradeDefinition{
				Attributes: map[string]AttributeUpgradeDefinition{
					"value": {LegacyType: StringAttribute},
				},
			},
			Expected: attributeUpgradeValidationExpected{},
		},
		{
			Name:     "test_orphan_block_no_legacy_no_inner_errors",
			Input:    AttributeUpgradeDefinition{},
			Expected: attributeUpgradeValidationExpected{ErrorCount: 1, Substrings: []string{"neither legacy_attribute / legacy_type"}},
		},
		{
			Name: "test_removed_block_missing_legacy_attribute",
			Input: AttributeUpgradeDefinition{
				LegacyType:        StringAttribute,
				LegacyRestriction: Optional,
				LegacyStatus:      Removed,
			},
			Expected: attributeUpgradeValidationExpected{ErrorCount: 1, Substrings: []string{"requires legacy_attribute"}},
		},
		{
			Name: "test_recurses_into_inner_attribute_validate",
			Input: AttributeUpgradeDefinition{
				LegacyAttribute: "old_block",
				Attributes: map[string]AttributeUpgradeDefinition{
					"inner": {LegacyStatus: Removed},
				},
			},
			Expected: attributeUpgradeValidationExpected{
				ErrorCount: 3,
				Substrings: []string{`attributes["inner"]`, "requires legacy_attribute", "requires legacy_type", "requires legacy_restriction"},
			},
		},
		{
			Name: "test_recurses_into_inner_child_validate_child",
			Input: AttributeUpgradeDefinition{
				LegacyAttribute: "old_block",
				Children: map[string]AttributeUpgradeDefinition{
					"grand": {},
				},
			},
			Expected: attributeUpgradeValidationExpected{
				ErrorCount: 1,
				Substrings: []string{`children["grand"]`, "neither legacy_attribute / legacy_type"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			node := testCase.Input.(AttributeUpgradeDefinition)
			expected := testCase.Expected.(attributeUpgradeValidationExpected)
			ctx := NewContext()
			node.validateChild(ctx, path)
			err := ctx.Diagnostics.Error()
			if expected.ErrorCount == 0 {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
			msg := err.Error()
			assert.Contains(t, msg, fmt.Sprintf("encountered %d error(s)", expected.ErrorCount))
			for _, substr := range expected.Substrings {
				assert.Contains(t, msg, substr)
			}
		})
	}
}

// TestAttributeUpgradeDefinition_HasInnerLegacyAttribute covers the recursive
// detector that backs the scalar-wrap shape check in validateChild.
func TestAttributeUpgradeDefinition_HasInnerLegacyAttribute(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_empty_node",
			Input:    AttributeUpgradeDefinition{},
			Expected: false,
		},
		{
			Name:     "test_self_legacy_attribute_does_not_count_as_inner",
			Input:    AttributeUpgradeDefinition{LegacyAttribute: "old_self"},
			Expected: false,
		},
		{
			Name: "test_inner_attribute_with_legacy_attribute",
			Input: AttributeUpgradeDefinition{
				Attributes: map[string]AttributeUpgradeDefinition{
					"value": {LegacyAttribute: "old_flat"},
				},
			},
			Expected: true,
		},
		{
			Name: "test_inner_attribute_without_legacy_attribute",
			Input: AttributeUpgradeDefinition{
				Attributes: map[string]AttributeUpgradeDefinition{
					"value": {LegacyType: StringAttribute},
				},
			},
			Expected: false,
		},
		{
			Name: "test_inner_child_with_legacy_attribute",
			Input: AttributeUpgradeDefinition{
				Children: map[string]AttributeUpgradeDefinition{
					"sub": {LegacyAttribute: "old_block"},
				},
			},
			Expected: true,
		},
		{
			Name: "test_nested_child_carries_legacy_attribute",
			Input: AttributeUpgradeDefinition{
				Children: map[string]AttributeUpgradeDefinition{
					"sub": {
						Children: map[string]AttributeUpgradeDefinition{
							"grand": {LegacyAttribute: "old_deep"},
						},
					},
				},
			},
			Expected: true,
		},
		{
			Name: "test_nested_child_no_legacy_anywhere",
			Input: AttributeUpgradeDefinition{
				Children: map[string]AttributeUpgradeDefinition{
					"sub": {
						Children: map[string]AttributeUpgradeDefinition{
							"grand": {LegacyType: StringAttribute},
						},
					},
				},
			},
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, testCase.Expected, testCase.Input.(AttributeUpgradeDefinition).hasInnerLegacyAttribute())
		})
	}
}
