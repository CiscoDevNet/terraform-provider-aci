package data

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
property_example_value_overrides:
  accessControlAllowOrigins: https://example.com
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
	assert.Equal(t, "https://example.com", parsedDefinition.PropertyExampleValueOverrides["accessControlAllowOrigins"])
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
example_parent_class: fvTenant
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
  registry_path: tenants
  description: A tenant.
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
        - reference:
            class_name: fvTenant
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
  run_serially: true
  post_test_delay_seconds: 15
  replace_auto_resolved: true
  ignore_tests:
    - child
  ignore_import_state_verify: true
  dependencies:
    - class_name: fvTenant
      reference: aci_tenant.test.id
      reference_type: resource
      role: parent
      ignore_changes:
        - relation_l3extrs_redistribute_pol
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
	assert.Equal(t, "fvTenant", parsedDefinition.ExampleParentClass)
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
	assert.Equal(t, "tenants", parsedDefinition.Documentation.RegistryPath)
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
	assert.Equal(t, &ExampleReferenceDefinition{ClassName: "fvTenant"}, nameProp.ExampleConfig.Minimum.Entries[0].Reference)
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
	assert.True(t, parsedDefinition.TestConfig.RunSerially)
	assert.Equal(t, 15, parsedDefinition.TestConfig.PostTestDelaySeconds)
	assert.Equal(t, []IgnoreTestEnum{ChildIgnoreTest}, parsedDefinition.TestConfig.IgnoreTests)
	assert.True(t, parsedDefinition.TestConfig.IgnoreImportStateVerify)
	assert.Len(t, parsedDefinition.TestConfig.Dependencies, 1)
	assert.Equal(t, "fvTenant", parsedDefinition.TestConfig.Dependencies[0].ClassName)
	assert.Equal(t, ResourceReference, parsedDefinition.TestConfig.Dependencies[0].ReferenceType)
	assert.Equal(t, Parent, parsedDefinition.TestConfig.Dependencies[0].Role)
	assert.Equal(t, []string{"relation_l3extrs_redistribute_pol"}, parsedDefinition.TestConfig.Dependencies[0].IgnoreChanges)
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

func TestParseClassDefinition_AssertStrategy(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	testCases := []test.TestCase{
		{Name: "test_omitted", Input: "", Expected: AssertValueStrategy},
		{Name: "test_value", Input: "\n          assert_strategy: value", Expected: AssertValueStrategy},
		{Name: "test_presence_only", Input: "\n          assert_strategy: presence_only", Expected: AssertPresenceOnly},
		{Name: "test_unknown", Input: "\n          assert_strategy: presence"},
		{Name: "test_empty", Input: "\n          assert_strategy: \"\""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			definition, err := parseClassDefinition([]byte("properties:\n  siteId:\n    test_config:\n      update:\n        - config_value: \"1\"" + testCase.Input.(string)))
			if testCase.Expected == nil {
				assert.Error(t, err)
				return
			}
			if assert.NoError(t, err) {
				assert.Equal(t, testCase.Expected.(AssertStrategyEnum), definition.Properties["siteId"].TestConfig.Update.Entries[0].AssertStrategy)
			}
		})
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

func TestParseClassDefinition_PropertyBooleanOverrides(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "omitted", Input: "", Expected: (*bool)(nil)},
		{Name: "true", Input: "true", Expected: boolPtr(true)},
		{Name: "false", Input: "false", Expected: boolPtr(false)},
	}

	for _, field := range []string{"computed", "deprecated", "requires_replace"} {
		t.Run(field, func(t *testing.T) {
			t.Parallel()
			for _, testCase := range testCases {
				t.Run(testCase.Name, func(t *testing.T) {
					t.Parallel()
					yamlText := "properties:\n  name:\n    restriction: optional\n"
					if value := testCase.Input.(string); value != "" {
						yamlText += fmt.Sprintf("    %s: %s\n", field, value)
					}
					definition, err := parseClassDefinition([]byte(yamlText))
					if !assert.NoError(t, err, test.MessageUnexpectedError(err)) {
						return
					}
					property := definition.Properties["name"]
					actual := map[string]*bool{"computed": property.Computed, "deprecated": property.Deprecated, "requires_replace": property.RequiresReplace}[field]
					assert.Equal(t, testCase.Expected.(*bool), actual, testCase.Name)
				})
			}
		})
	}
}

// TestParseClassDefinition_EmptyYAML verifies that an empty payload decodes
// to a zero-value ClassDefinition without error.
func TestParseClassDefinition_TestPrefix(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	testCases := []test.TestCase{
		{Name: "test_automatic", Input: "", Expected: (*bool)(nil)},
		{Name: "test_enabled", Input: "prefix: {enabled: true}", Expected: true},
		{Name: "test_excluded", Input: "prefix: {enabled: false}", Expected: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			definition, err := parseClassDefinition([]byte(fmt.Sprintf(`
test_config:
  prefix:
    characters: abcdef
    length: 12
  dependencies:
  - class_name: fvTenant
    %s
properties:
  name:
    test_config:
      %s
`, testCase.Input, testCase.Input)))
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, TestPrefixDefinition{Characters: "abcdef", Length: 12}, definition.TestConfig.Prefix)
			for _, enabled := range []*bool{definition.TestConfig.Dependencies[0].Prefix.Enabled, definition.Properties["name"].TestConfig.Prefix.Enabled} {
				if expected, ok := testCase.Expected.(bool); ok {
					if assert.NotNil(t, enabled) {
						assert.Equal(t, expected, *enabled)
					}
				} else {
					assert.Nil(t, enabled)
				}
			}
			class := Class{ClassDefinition: definition}
			class.setTestConfig()
			assert.Equal(t, definition.TestConfig.Prefix, class.TestConfig.Prefix)
			property := Property{PropertyName: "name", AttributeName: "name", ValueType: String, propertyDefinition: definition.Properties["name"]}
			property.setTestValues()
			assert.Equal(t, definition.Properties["name"].TestConfig.Prefix.Enabled, property.PrefixEnabled)
		})
	}
}

func TestParseClassDefinition_EmptyYAML(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	parsedDefinition, err := parseClassDefinition(nil)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, ClassDefinition{}, parsedDefinition)
}

func boolPtr(b bool) *bool { return &b }

// TestParseClassDefinition_MigrationSource decodes origin per saved schema version.
func TestParseClassDefinition_MigrationSource(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "omitted", Input: "resource_name: tenant\n", Expected: UndefinedMigrationSource},
		{Name: "from_sdkv2", Input: "state_upgrades:\n  0:\n    source: sdkv2\n", Expected: FromSDKv2},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			parsedDefinition, err := parseClassDefinition([]byte(testCase.Input.(string)))
			assert.NoError(t, err, test.MessageUnexpectedError(err))
			assert.Equal(t, testCase.Expected.(MigrationSourceEnum), parsedDefinition.StateUpgrades[0].Source, testCase.Name)
		})
	}
}

// TestParseClassDefinition_StateUpgrades is the canonical happy-path decode for
// state_upgrades. It mixes every scenario the parser needs to handle into one
// realistic fvBD (Bridge Domain) example:
//
//   - Two state_upgrades entries at distinct schema versions (1 and 3),
//     proving slice decode at non-zero unique values.
//   - An attribute with only attribute_name set, proving the trio of Undefined
//     defaults (Type, Restriction, Status) for omitted fields.
//   - An attribute with all four fields explicitly set, including status=frozen.
//   - An attribute with status=removed (the case where type and
//     restriction must be supplied because there is no current attribute
//     to inherit from).
//   - A child block (fvSubnet) carrying its own attribute_name (block rename)
//     alongside inner attributes AND its own children (fvRsBDSubnetToOut),
//     proving 3-level recursion: children -> children -> attributes.
//   - A child block (fvRsCtx) without a block-level attribute_name, proving
//     the empty-string decode for an unchanged block whose inner attributes
//     are the only diff.
//   - An attribute with status explicitly set to functioning, proving
//     that the explicit value round-trips identically to the default.
func TestParseClassDefinition_StateUpgrades(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	yamlBytes := []byte(`state_upgrades:
  1:
    attributes:
      arpFlood:
        attribute_name: arp_flooding
        restriction: optional
        status: frozen
        type: bool_attribute
      mac:
        attribute_name: mac_address
      multiDstPktAct:
        attribute_name: multi_dst_pkt_act
        restriction: optional
        status: removed
        type: string_attribute
    children:
      fvRsCtx:
        attributes:
          tnFvCtxName:
            attribute_name: vrf_name
            status: deprecated
      fvSubnet:
        attributes:
          ip:
            attribute_name: subnet_ip
        children:
          fvRsBDSubnetToOut:
            attributes:
              tnL3extOutName:
                attribute_name: l3out_name
                type: string_attribute
        attribute_name: subnets
  3:
    attributes:
      descr:
        attribute_name: description
`,
	)

	parsedDefinition, err := parseClassDefinition(yamlBytes)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, UndefinedMigrationSource, parsedDefinition.StateUpgrades[1].Source, "an omitted version source remains undefined")

	if !assert.Len(t, parsedDefinition.StateUpgrades, 2, "two state_upgrades entries must round-trip") {
		return
	}

	v1 := parsedDefinition.StateUpgrades[1]

	assert.Len(t, v1.Attributes, 3)

	// Defaults proof: only attribute_name set, the trio of Undefined defaults applies.
	mac := v1.Attributes["mac"]
	assert.Equal(t, "mac_address", mac.AttributeName)
	assert.Equal(t, UndefinedLegacyAttributeType, mac.Type)
	assert.Equal(t, UndefinedLegacyRestriction, mac.Restriction)
	assert.Equal(t, UndefinedLegacyStatus, mac.Status, "omission must remain distinguishable for status inheritance")

	// All four fields explicitly set, status=frozen.
	arp := v1.Attributes["arpFlood"]
	assert.Equal(t, "arp_flooding", arp.AttributeName)
	assert.Equal(t, BoolAttribute, arp.Type)
	assert.Equal(t, LegacyOptional, arp.Restriction)
	assert.Equal(t, Frozen, arp.Status)

	// status=removed: the attribute is gone from the current schema, so
	// type and restriction must be supplied here.
	rem := v1.Attributes["multiDstPktAct"]
	assert.Equal(t, "multi_dst_pkt_act", rem.AttributeName)
	assert.Equal(t, StringAttribute, rem.Type)
	assert.Equal(t, LegacyOptional, rem.Restriction)
	assert.Equal(t, Removed, rem.Status)

	// Block with its own attribute_name (rename) + inner attributes + nested children.
	subnet, ok := v1.Children["fvSubnet"]
	if assert.True(t, ok, "outer child must decode") {
		assert.Equal(t, "subnets", subnet.AttributeName, "block-level attribute_name records a block rename")
		assert.Equal(t, "subnet_ip", subnet.Attributes["ip"].AttributeName)

		rsToOut, ok := subnet.Children["fvRsBDSubnetToOut"]
		if assert.True(t, ok, "grandchild block must decode (children -> children)") {
			inner, ok := rsToOut.Attributes["tnL3extOutName"]
			if assert.True(t, ok, "3-level recursion: children -> children -> attributes") {
				assert.Equal(t, "l3out_name", inner.AttributeName)
				assert.Equal(t, StringAttribute, inner.Type)
			}
		}
	}

	// Block without its own attribute_name (block unchanged, only inner attribute changed).
	rsCtx, ok := v1.Children["fvRsCtx"]
	if assert.True(t, ok) {
		assert.Equal(t, "", rsCtx.AttributeName, "absent block-level attribute_name decodes to empty")
		inner := rsCtx.Attributes["tnFvCtxName"]
		assert.Equal(t, "vrf_name", inner.AttributeName)
		assert.Equal(t, LegacyDeprecated, inner.Status, "explicit 'functioning' must round-trip identically to the omitted default")
	}

	v3 := parsedDefinition.StateUpgrades[3]

	assert.Equal(t, "description", v3.Attributes["descr"].AttributeName)
}

// TestTopSystemMigrationDefinition protects the datasource's documented DN and
// complete deferred legacy-alias mapping.
func TestTopSystemMigrationDefinition(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("determine definitions test path")
	}
	definitionPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "definitions", "topSystem.yaml")
	contents, err := os.ReadFile(definitionPath)
	if err != nil {
		t.Fatalf("read topSystem definition: %v", err)
	}
	definition, err := parseClassDefinition(contents)
	if err != nil {
		t.Fatalf("parse topSystem definition: %v", err)
	}

	assert.Equal(t, []string{"topology/pod-{id}/node-{id}/sys"}, definition.Documentation.DnFormats)
	assert.Equal(t, FromSDKv2, definition.StateUpgrades[1].Source)
	assert.Equal(t, "The `id` attribute of the parent APIC `fabricNode` class.", definition.Properties["id"].Documentation.Description)
	assert.Equal(t, "id", definition.Properties["podId"].DocumentationPropertyName)
	assert.Equal(t, "The `id` attribute of the parent APIC `fabricPod` class.", definition.Properties["podId"].Documentation.Description)
	if !assert.Len(t, definition.StateUpgrades, 1) {
		return
	}

	upgrade := definition.StateUpgrades[1]

	expectedAttributes := map[string]AttributeUpgradeDefinition{
		"bootstrapState":          {AttributeName: "boot_strap_tate"},
		"etepAddr":                {AttributeName: "etep_addr"},
		"inbMgmtAddr":             {AttributeName: "inb_mgmt_addr"},
		"inbMgmtAddr6":            {AttributeName: "inb_mgmt_addr6"},
		"inbMgmtAddr6Mask":        {AttributeName: "inb_mgmt_addr6_mask"},
		"inbMgmtAddrMask":         {AttributeName: "inb_mgmt_addr_mask"},
		"inbMgmtGateway":          {AttributeName: "inb_mgmt_gateway"},
		"inbMgmtGateway6":         {AttributeName: "inb_mgmt_gateway6"},
		"monPolDn":                {AttributeName: "mon_pol_dn"},
		"oobMgmtAddr":             {AttributeName: "oob_mgmt_addr"},
		"oobMgmtAddr6":            {AttributeName: "oob_mgmt_addr6"},
		"oobMgmtAddr6Mask":        {AttributeName: "oob_mgmt_addr6_mask"},
		"oobMgmtAddrMask":         {AttributeName: "oob_mgmt_addr_mask"},
		"oobMgmtGateway":          {AttributeName: "oob_mgmt_gateway"},
		"oobMgmtGateway6":         {AttributeName: "oob_mgmt_gateway6"},
		"rlAutoMode":              {AttributeName: "rl_auto_mode"},
		"rlGroupId":               {AttributeName: "rl_group_id"},
		"rlOperPodId":             {AttributeName: "rl_oper_pod_id"},
		"rlRoutableMode":          {AttributeName: "rl_routable_mode"},
		"rldirectMode":            {AttributeName: "rldirect_mode"},
		"unicastXrEpLearnDisable": {AttributeName: "unicast_xr_ep_learn_disable"},
	}
	for key, attribute := range expectedAttributes {
		attribute.Status = LegacyDeprecated
		if key != "rlAutoMode" && key != "rlGroupId" {
			attribute.Restriction = LegacyOptionalComputed
		}
		expectedAttributes[key] = attribute
	}
	assert.Equal(t, expectedAttributes, upgrade.Attributes)
}

func TestHistoricalDefinitionDifferences(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("determine definitions test path")
	}
	definitionsPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "definitions")
	definitions := map[string]ClassDefinition{}
	for _, className := range []string{"fvAEPg", "fvBD", "fvCtx", "fvESg", "fvTenant", "vmmDomP", "stpIfPol", "netflowMonitorPol", "infraAccPortP", "infraFexP", "infraSpAccPortP", "rtctrlProfile", "tagTag"} {
		contents, err := os.ReadFile(filepath.Join(definitionsPath, className+".yaml"))
		if !assert.NoError(t, err) {
			continue
		}
		definition, err := parseClassDefinition(contents)
		assert.NoError(t, err)
		definitions[className] = definition
		for _, upgrade := range definition.StateUpgrades {
			for _, attribute := range []string{"id", "name", "nameAlias", "annotation", "descr"} {
				assert.NotContains(t, upgrade.Attributes, attribute, "%s must inherit unchanged fields", className)
			}
		}
	}

	testCases := []test.TestCase{
		{Name: "STP preserves the old list", Input: definitions["stpIfPol"].StateUpgrades[1].Attributes["ctrl"], Expected: AttributeUpgradeDefinition{AttributeName: "ctrl", Type: ListAttribute, ElementType: StringAttribute, Status: LegacyDeprecated}},
		{Name: "EPG preserves the old set of DNs", Input: definitions["fvAEPg"].StateUpgrades[1].Children["fvRsCons"].Attributes["tnVzBrCPName"], Expected: AttributeUpgradeDefinition{AttributeName: "relation_fv_rs_cons", Type: SetAttribute, ElementType: StringAttribute, Restriction: LegacyOptional, Status: LegacyDeprecated}},
		{Name: "VRF preserves a computed frozen string", Input: definitions["fvCtx"].StateUpgrades[1].Attributes["knw_mcast_act"], Expected: AttributeUpgradeDefinition{Type: StringAttribute, Restriction: LegacyOptionalComputed, Status: Frozen}},
		{Name: "old block inner name maps to a named relation", Input: definitions["fvESg"].StateUpgrades[1].Children["fvRsCons"].Attributes["tnVzBrCPName"], Expected: AttributeUpgradeDefinition{AttributeName: "target_dn"}},
		{Name: "unchanged block members retain their historical names", Input: definitions["fvAEPg"].StateUpgrades[1].Children["fvRsNodeAtt"].Attributes["mode"], Expected: AttributeUpgradeDefinition{AttributeName: "mode"}},
		{Name: "Netflow keeps its framework collection shape", Input: definitions["netflowMonitorPol"].StateUpgrades[0].Children["netflowRsMonitorToRecord"], Expected: AttributeUpgradeDefinition{Type: SetNestedAttribute}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.Equal(t, testCase.Expected, testCase.Input)
		})
	}
	assert.Equal(t, FromFramework, definitions["netflowMonitorPol"].StateUpgrades[0].Source)
	assert.Equal(t, "relation_fv_rs_node_att", definitions["fvAEPg"].StateUpgrades[1].Children["fvRsNodeAtt"].AttributeName)
	assert.Equal(t, LegacyDeprecated, definitions["fvAEPg"].StateUpgrades[1].Children["fvRsNodeAtt"].Status)

	// Every retained orphan needs its own old shape; none can inherit it from current metadata.
	frozenFields := map[string][]string{
		"fvAEPg":   {"relation_fv_rs_path_att", "relation_fv_rs_prov_def"},
		"fvBD":     {"relation_fv_rs_bd_flood_to", "relation_fv_rs_bd_to_profile"},
		"fvCtx":    {"knw_mcast_act", "relation_fv_rs_ctx_mcast_to", "relation_fv_rs_vrf_validation_pol"},
		"fvESg":    {"relation_fv_rs_cust_qos_pol", "relation_fv_rs_prot_by"},
		"fvTenant": {"relation_fv_rs_tn_deny_rule"},
		"vmmDomP":  {"relation_infra_rs_dom_vxlan_ns_def", "relation_infra_rs_vlan_ns_def", "relation_vmm_rs_default_cdp_if_pol", "relation_vmm_rs_default_fw_pol", "relation_vmm_rs_default_l2_inst_pol", "relation_vmm_rs_default_lacp_lag_pol", "relation_vmm_rs_default_lldp_if_pol", "relation_vmm_rs_default_stp_if_pol"},
	}
	for className, fields := range frozenFields {
		count := 0
		for _, attribute := range definitions[className].StateUpgrades[1].Attributes {
			if attribute.Status == Frozen {
				count++
			}
		}
		assert.Equal(t, len(fields), count, "%s frozen fields", className)
		for _, field := range fields {
			attribute := definitions[className].StateUpgrades[1].Attributes[field]
			assert.Equal(t, Frozen, attribute.Status, "%s.%s", className, field)
			assert.NotEqual(t, UndefinedLegacyAttributeType, attribute.Type, "%s.%s", className, field)
			assert.NotEqual(t, UndefinedLegacyRestriction, attribute.Restriction, "%s.%s", className, field)
			if attribute.Type == SetAttribute {
				assert.Equal(t, StringAttribute, attribute.ElementType, "%s.%s", className, field)
			}
		}
	}
	for _, className := range []string{"infraAccPortP", "infraFexP", "infraSpAccPortP", "rtctrlProfile", "tagTag"} {
		assert.Empty(t, definitions[className].StateUpgrades, "%s needs only the migration notice", className)
		assert.Contains(t, definitions[className].Documentation.Resource.Warnings, SDKv2MigrationWarning)
	}
}

func TestPublishedResourceSchemaVersions(t *testing.T) {
	t.Parallel()
	// Preserve the published develop contract independently of migration implementation progress.
	testCases := []test.TestCase{
		{Name: "cdpIfPol", Expected: 2},
		{Name: "fvAEPg", Expected: 2},
		{Name: "fvAp", Expected: 2},
		{Name: "fvBD", Expected: 2},
		{Name: "fvCtx", Expected: 2},
		{Name: "fvESg", Expected: 2},
		{Name: "fvTenant", Expected: 2},
		{Name: "infraAccPortP", Expected: 1},
		{Name: "infraAttEntityP", Expected: 2},
		{Name: "infraFexP", Expected: 1},
		{Name: "infraHPortS", Expected: 2},
		{Name: "infraPortBlk", Expected: 2},
		{Name: "infraSHPortS", Expected: 2},
		{Name: "infraSpAccPortP", Expected: 1},
		{Name: "l2IfPol", Expected: 2},
		{Name: "lldpIfPol", Expected: 2},
		{Name: "netflowMonitorPol", Expected: 1},
		{Name: "rtctrlProfile", Expected: 1},
		{Name: "stpIfPol", Expected: 2},
		{Name: "tagTag", Expected: 1},
		{Name: "vmmDomP", Expected: 2},
		{Name: "vmmUsrAccP", Expected: 2},
		{Name: "vzAny", Expected: 2},
		{Name: "vzTaboo", Expected: 2},
	}
	_, source, _, _ := runtime.Caller(0)
	definitionsPath := filepath.Join(filepath.Dir(source), "..", "..", "definitions")
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			payload, err := os.ReadFile(filepath.Join(definitionsPath, testCase.Name+".yaml"))
			if !assert.NoError(t, err) {
				return
			}
			definition, err := parseClassDefinition(payload)
			if !assert.NoError(t, err) || !assert.NotNil(t, definition.SchemaVersion) {
				return
			}
			assert.Equal(t, testCase.Expected, *definition.SchemaVersion)
			if testCase.Expected == 1 && testCase.Name != "netflowMonitorPol" {
				assert.Empty(t, definition.StateUpgrades, "an unchanged published version needs no artificial upgrade entry")
			}
		})
	}
}

// TestParseClassDefinition_StateUpgradesInvalidEnumValues confirms strict
// decoding rejects unknown values for each enum reachable through state_upgrades.
func TestParseClassDefinition_StateUpgradesInvalidEnumValues(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "invalid_type",
			Input: `state_upgrades:
  0:
    attributes:
      name:
        type: not_a_type
`,

			Expected: "unknown type",
		},
		{
			Name: "invalid_status",
			Input: `state_upgrades:
  0:
    attributes:
      name:
        status: weird
`,

			Expected: "unknown status",
		},
		{
			Name: "invalid_restriction",
			Input: `state_upgrades:
  0:
    attributes:
      name:
        restriction: weird
`,

			Expected: "unknown restriction",
		},
		{
			Name:     "invalid_migration_source",
			Input:    "state_upgrades:\n  0:\n    source: bogus\n",
			Expected: "unknown source",
		},
		{
			Name: "empty_status_is_typo",
			Input: `state_upgrades:
  0:
    attributes:
      name:
        status: ""
`,

			Expected: "unknown status",
		},
		{
			Name: "empty_restriction_is_typo",
			Input: `state_upgrades:
  0:
    attributes:
      name:
        restriction: ""
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
			Input: `state_upgrades:
  0:
    bogus_upgrade_field: true
`,

			Expected: "bogus_upgrade_field",
		},
		{
			Name: "unknown_field_on_attribute_node",
			Input: `state_upgrades:
  0:
    attributes:
      name:
        bogus_attribute_field: true
`,

			Expected: "bogus_attribute_field",
		},
		{
			Name: "unknown_field_on_inner_child_node",
			Input: `state_upgrades:
  0:
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

	const path = `Class 'fvCtx': schema version 0: attributes["name"]`

	testCases := []test.TestCase{
		{
			Name:     "test_non_removed_status_skips_all_checks",
			Input:    AttributeUpgradeDefinition{Status: LegacyDeprecated},
			Expected: attributeUpgradeValidationExpected{},
		},
		{
			Name: "test_removed_with_all_legacy_fields_set_passes",
			Input: AttributeUpgradeDefinition{
				AttributeName: "old_name",
				Type:          StringAttribute,
				Restriction:   LegacyOptional,
				Status:        Removed,
			},
			Expected: attributeUpgradeValidationExpected{},
		},
		{
			Name: "test_removed_missing_attribute_name_only",
			Input: AttributeUpgradeDefinition{
				Type:        StringAttribute,
				Restriction: LegacyOptional,
				Status:      Removed,
			},
			Expected: attributeUpgradeValidationExpected{ErrorCount: 1, Substrings: []string{"requires attribute_name"}},
		},
		{
			Name: "test_removed_missing_type_only",
			Input: AttributeUpgradeDefinition{
				AttributeName: "old_name",
				Restriction:   LegacyOptional,
				Status:        Removed,
			},
			Expected: attributeUpgradeValidationExpected{ErrorCount: 1, Substrings: []string{"requires type"}},
		},
		{
			Name: "test_removed_missing_restriction_only",
			Input: AttributeUpgradeDefinition{
				AttributeName: "old_name",
				Type:          StringAttribute,
				Status:        Removed,
			},
			Expected: attributeUpgradeValidationExpected{ErrorCount: 1, Substrings: []string{"requires restriction"}},
		},
		{
			Name:  "test_removed_all_fields_missing_aggregates_three_diagnostics",
			Input: AttributeUpgradeDefinition{Status: Removed},
			Expected: attributeUpgradeValidationExpected{
				ErrorCount: 3,
				Substrings: []string{"requires attribute_name", "requires type", "requires restriction"},
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
// In addition to the Removed-only attribute_name check, it enforces the
// scalar-wrap shape (the only allowed way for a child entry to omit both
// attribute_name and type) and recurses into inner attributes/children.
func TestAttributeUpgradeDefinition_ValidateChild(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	const path = `Class 'fvCtx': schema version 0: children["fvRsBd"]`

	testCases := []test.TestCase{
		{
			Name:     "test_block_rename_only",
			Input:    AttributeUpgradeDefinition{AttributeName: "old_block"},
			Expected: attributeUpgradeValidationExpected{},
		},
		{
			Name: "test_scalar_wrap_inner_attribute_carries_legacy",
			Input: AttributeUpgradeDefinition{
				Attributes: map[string]AttributeUpgradeDefinition{
					"value": {AttributeName: "old_flat"},
				},
			},
			Expected: attributeUpgradeValidationExpected{},
		},
		{
			Name: "test_inner_attributes_only_no_legacy_on_block_or_inner",
			Input: AttributeUpgradeDefinition{
				Attributes: map[string]AttributeUpgradeDefinition{
					"value": {Type: StringAttribute},
				},
			},
			Expected: attributeUpgradeValidationExpected{},
		},
		{
			Name:     "test_orphan_block_no_legacy_no_inner_errors",
			Input:    AttributeUpgradeDefinition{},
			Expected: attributeUpgradeValidationExpected{ErrorCount: 1, Substrings: []string{"neither attribute_name / type"}},
		},
		{
			Name: "test_removed_block_missing_attribute_name",
			Input: AttributeUpgradeDefinition{
				Type:        StringAttribute,
				Restriction: LegacyOptional,
				Status:      Removed,
			},
			Expected: attributeUpgradeValidationExpected{ErrorCount: 1, Substrings: []string{"requires attribute_name"}},
		},
		{
			Name: "test_recurses_into_inner_attribute_validate",
			Input: AttributeUpgradeDefinition{
				AttributeName: "old_block",
				Attributes: map[string]AttributeUpgradeDefinition{
					"inner": {Status: Removed},
				},
			},
			Expected: attributeUpgradeValidationExpected{
				ErrorCount: 3,
				Substrings: []string{`attributes["inner"]`, "requires attribute_name", "requires type", "requires restriction"},
			},
		},
		{
			Name: "test_recurses_into_inner_child_validate_child",
			Input: AttributeUpgradeDefinition{
				AttributeName: "old_block",
				Children: map[string]AttributeUpgradeDefinition{
					"grand": {},
				},
			},
			Expected: attributeUpgradeValidationExpected{
				ErrorCount: 1,
				Substrings: []string{`children["grand"]`, "neither attribute_name / type"},
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
			Name:     "test_self_attribute_name_does_not_count_as_inner",
			Input:    AttributeUpgradeDefinition{AttributeName: "old_self"},
			Expected: false,
		},
		{
			Name: "test_inner_attribute_with_attribute_name",
			Input: AttributeUpgradeDefinition{
				Attributes: map[string]AttributeUpgradeDefinition{
					"value": {AttributeName: "old_flat"},
				},
			},
			Expected: true,
		},
		{
			Name: "test_inner_attribute_without_attribute_name",
			Input: AttributeUpgradeDefinition{
				Attributes: map[string]AttributeUpgradeDefinition{
					"value": {Type: StringAttribute},
				},
			},
			Expected: false,
		},
		{
			Name: "test_inner_child_with_attribute_name",
			Input: AttributeUpgradeDefinition{
				Children: map[string]AttributeUpgradeDefinition{
					"sub": {AttributeName: "old_block"},
				},
			},
			Expected: true,
		},
		{
			Name: "test_nested_child_carries_attribute_name",
			Input: AttributeUpgradeDefinition{
				Children: map[string]AttributeUpgradeDefinition{
					"sub": {
						Children: map[string]AttributeUpgradeDefinition{
							"grand": {AttributeName: "old_deep"},
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
							"grand": {Type: StringAttribute},
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

func TestTestConfigScenarioPresence(t *testing.T) {
	t.Parallel()
	testCases := []test.TestCase{
		{Name: "test_include", Input: "include_in: [minimum, reset]"},
		{Name: "test_exclude", Input: "exclude_from: [update]"},
		{Name: "test_distinct_scenarios", Input: "include_in: [minimum]\n      exclude_from: [update]"},
		{Name: "test_conflicting_scenario", Input: "include_in: [reset]\n      exclude_from: [reset]", Expected: "appears in both include_in and exclude_from"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			_, err := parseClassDefinition([]byte("properties:\n  descr:\n    test_config:\n      " + testCase.Input.(string) + "\n"))
			if expected, _ := testCase.Expected.(string); expected != "" {
				assert.ErrorContains(t, err, expected)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseClassDefinitionRejectsDuplicateSchemaVersions(t *testing.T) {
	test.InitializeTest(t)
	_, err := parseClassDefinition([]byte("state_upgrades:\n  1: {source: sdkv2}\n  1: {source: framework}\n"))
	assert.Error(t, err)
}
