package data

import (
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
no_meta_file:
  fvCtx: vrf
documentation_label_overrides:
  Bgp: BGP
`)

	parsedDefinition, err := parseGlobalMetaDefinition(yamlBytes)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, []string{"fvRsBd"}, parsedDefinition.AlwaysIncludeAsChild)
	assert.Equal(t, "description", parsedDefinition.AttributeNameOverrides["descr"])
	assert.Equal(t, []string{"polUni"}, parsedDefinition.ExcludeParents)
	assert.Equal(t, []string{"childAction"}, parsedDefinition.ExcludeProperties)
	assert.Equal(t, "vrf", parsedDefinition.NoMetaFile["fvCtx"])
	assert.Equal(t, "BGP", parsedDefinition.DocumentationLabelOverrides["Bgp"])
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
identified_by:
  - name
documentation:
  label: Tenant
  description: A tenant.
relation_info:
  type: named
  to_classes:
    - vz:BrCP
properties:
  name:
    restriction: required
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
test_config:
  replace_auto_resolved: true
  dependencies:
    - class_name: fvTenant
      reference: aci_tenant.test.id
      reference_type: resource
      role: parent
`)

	parsedDefinition, err := parseClassDefinition(yamlBytes)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, "tenant", parsedDefinition.ResourceName)
	assert.Equal(t, []string{"name"}, parsedDefinition.IdentifiedBy)
	assert.Equal(t, "Tenant", parsedDefinition.Documentation.Label)
	assert.Equal(t, Named, parsedDefinition.RelationInfo.Type)
	assert.Equal(t, []string{"vz:BrCP"}, parsedDefinition.RelationInfo.ToClasses)

	nameProp, ok := parsedDefinition.Properties["name"]
	assert.True(t, ok)
	assert.Equal(t, Required, nameProp.Restriction)
	assert.Len(t, nameProp.TestConfig.Create, 1)
	assert.Equal(t, "test_tenant", nameProp.TestConfig.Create[0].ConfigValue)
	assert.NotNil(t, nameProp.TestConfig.Create[0].ConfigInclude)
	assert.True(t, *nameProp.TestConfig.Create[0].ConfigInclude)

	descrProp, ok := parsedDefinition.Properties["descr"]
	assert.True(t, ok)
	assert.True(t, descrProp.TestConfig.IgnoreInTest)
	assert.Len(t, descrProp.TestConfig.Update, 1)
	assert.Equal(t, "updated", descrProp.TestConfig.Update[0].AssertValue)

	assert.True(t, parsedDefinition.TestConfig.ReplaceAutoResolved)
	assert.Len(t, parsedDefinition.TestConfig.Dependencies, 1)
	assert.Equal(t, "fvTenant", parsedDefinition.TestConfig.Dependencies[0].ClassName)
	assert.Equal(t, ResourceReference, parsedDefinition.TestConfig.Dependencies[0].ReferenceType)
	assert.Equal(t, Parent, parsedDefinition.TestConfig.Dependencies[0].Role)
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

	cases := []test.TestCase{
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
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			_, err := parseClassDefinition([]byte(tc.Input.(string)))
			assert.Error(t, err, tc.Name)
			assert.Contains(t, err.Error(), tc.Expected.(string), tc.Name)
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
	entries := parsedDefinition.Properties["name"].TestConfig.Create
	assert.Len(t, entries, 3)
	assert.Nil(t, entries[0].ConfigInclude, "omitted config_include must decode to nil")
	if assert.NotNil(t, entries[1].ConfigInclude) {
		assert.True(t, *entries[1].ConfigInclude)
	}
	if assert.NotNil(t, entries[2].ConfigInclude) {
		assert.False(t, *entries[2].ConfigInclude, "explicit false must survive as *bool=false, not nil")
	}
}

// TestParseClassDefinition_RequiresReplacePointerSemantics mirrors the
// ConfigInclude check for the top-level *bool requires_replace override.
func TestParseClassDefinition_RequiresReplacePointerSemantics(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	cases := []test.TestCase{
		{Name: "omitted", Input: "properties:\n  name:\n    restriction: required\n", Expected: (*bool)(nil)},
		{Name: "true", Input: "properties:\n  name:\n    requires_replace: true\n", Expected: boolPtr(true)},
		{Name: "false", Input: "properties:\n  name:\n    requires_replace: false\n", Expected: boolPtr(false)},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			parsedDefinition, err := parseClassDefinition([]byte(tc.Input.(string)))
			assert.NoError(t, err, test.MessageUnexpectedError(err))
			expected := tc.Expected.(*bool)
			actual := parsedDefinition.Properties["name"].RequiresReplace
			if expected == nil {
				assert.Nil(t, actual, tc.Name)
			} else {
				if assert.NotNil(t, actual, tc.Name) {
					assert.Equal(t, *expected, *actual, tc.Name)
				}
			}
		})
	}
}

// TestParseClassDefinition_EmptyYAML verifies that an empty payload decodes
// to a zero-value ClassDefinition without error — the same behavior used by
// loadClassDefinition when the class has no override file.
func TestParseClassDefinition_EmptyYAML(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	parsedDefinition, err := parseClassDefinition(nil)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, ClassDefinition{}, parsedDefinition)
}

func boolPtr(b bool) *bool { return &b }
