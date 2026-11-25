package data

import (
	"encoding/json"
	"fmt"
	"maps"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type dataSourceDependencyRenderingInput struct {
	Class     Class
	DataStore *DataStore
}

type dataSourceDependencyRenderingExpected struct {
	Contains    []string
	NotContains []string
	Error       string
}

func TestMigrationTestCaseValues(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "values", Expected: 2},
		{Name: "null_relation", Expected: 0},
		{Name: "empty_relation", Expected: 0},
		{Name: "missing_optional_fields", Expected: 2},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			child := classWithProperties("fvRsNodeAtt", map[string]*Property{
				"tDn":   {PropertyName: "tDn", AttributeName: "target_dn", ValueType: String, Required: true},
				"descr": {PropertyName: "descr", AttributeName: "description", ValueType: String},
			}, nil)
			child.PropertiesAll = []string{"tDn", "descr"}
			child.IdentifiedBy = []string{"tDn"}
			child.ResourceNameNested = "nodes"
			class := classWithProperties("fvTenant", map[string]*Property{
				"name":       {PropertyName: "name", AttributeName: "name", ValueType: String, Required: true},
				"descr":      {PropertyName: "descr", AttributeName: "description", ValueType: String},
				"annotation": {PropertyName: "annotation", AttributeName: "annotation", ValueType: String},
			}, []string{"fvRsNodeAtt"})
			class.PropertiesAll = []string{"name", "descr", "annotation"}
			version := 2
			class.ClassDefinition.SchemaVersion = &version
			class.StateUpgrades = map[int]StateUpgradeDefinition{1: {
				Attributes: map[string]AttributeUpgradeDefinition{
					"descr":         {AttributeName: "old_description"},
					"removed_field": {Type: StringAttribute, Restriction: LegacyOptional, Status: Removed},
				},
				Children: map[string]AttributeUpgradeDefinition{"fvRsNodeAtt": {
					AttributeName: "old_nodes", Attributes: map[string]AttributeUpgradeDefinition{
						"tDn": {AttributeName: "node_dn"}, "descr": {AttributeName: "description"},
					},
				}},
			}}
			store := &DataStore{Classes: map[string]Class{"fvRsNodeAtt": *child}, ctx: NewContext()}
			class.setMigrationData(store)
			if !assert.NoError(t, store.ctx.Diagnostics.Error()) {
				return
			}
			cases, err := class.MigrationTestCases(store)
			if !assert.NoError(t, err) {
				return
			}
			for _, fixture := range cases {
				if fixture.Name != "v1_"+testCase.Name {
					continue
				}
				var prior, expected map[string]any
				if !assert.NoError(t, json.Unmarshal([]byte(fixture.State), &prior)) || !assert.NoError(t, json.Unmarshal([]byte(fixture.Expected), &expected)) {
					return
				}
				assert.Equal(t, "name_value_1", prior["name"])
				assert.Equal(t, "old_description_value_1", prior["old_description"])
				assert.Equal(t, "old_description_value_1", expected["description"])
				assert.Equal(t, "removed_field_value_1", prior["removed_field"])
				assert.NotContains(t, expected, "removed_field")
				if testCase.Name != "missing_optional_fields" {
					assert.Equal(t, "annotation_value_1", prior["annotation"])
				}
				if testCase.Expected.(int) == 0 {
					assert.Empty(t, prior["old_nodes"])
					assert.Empty(t, expected["nodes"])
					return
				}
				blocks := prior["old_nodes"].([]any)
				children := expected["nodes"].([]any)
				if !assert.Len(t, blocks, 2) || !assert.Len(t, children, 2) {
					return
				}
				assert.Equal(t, "old_nodes_node_dn_value_1", blocks[0].(map[string]any)["node_dn"])
				assert.Equal(t, "old_nodes_description_value_2", blocks[1].(map[string]any)["description"])
				assert.Equal(t, "old_nodes_node_dn_value_2", children[1].(map[string]any)["target_dn"])
				return
			}
			t.Fatal("migration fixture missing")
		})
	}
}

func TestDependencyIgnoredChanges(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	type input struct {
		Loaded        bool
		IgnoreChanges []string
		ReferenceType ReferenceTypeEnum
	}
	testCases := []test.TestCase{
		{Name: "test_legacy_parent", Input: input{IgnoreChanges: []string{"relation_l3extrs_redistribute_pol"}}, Expected: true},
		{Name: "test_loaded_parent", Input: input{Loaded: true, IgnoreChanges: []string{"relation_l3extrs_redistribute_pol"}}, Expected: true},
		{Name: "test_without_override", Input: input{}, Expected: false},
		{Name: "test_datasource_has_no_lifecycle", Input: input{ReferenceType: DataSourceReference, IgnoreChanges: []string{"relation_l3extrs_redistribute_pol"}}, Expected: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(input)
			definition := ClassDefinition{TestConfig: ClassTestConfigDefinition{Dependencies: []TestDependencyDefinition{{
				ClassName: "l3extOut", Reference: "aci_l3_outside.test.id", ReferenceType: input.ReferenceType, Role: Parent,
				IgnoreChanges: input.IgnoreChanges, ConfigOverrides: map[string]string{"name": "name_1"},
			}}}}
			if input.ReferenceType == DataSourceReference {
				definition.TestConfig.Dependencies[0].Reference = "data.aci_l3_outside.test.id"
			}
			root := dataSourceTestClass("exampleRoot", "example_root", ResourceArtifact, DatasourceArtifact)
			root.ClassDefinition = definition
			store := &DataStore{ctx: NewContext(), Classes: map[string]Class{}, ClassDefinitions: map[string]ClassDefinition{"l3extOut": {ResourceName: "l3_outside"}}}
			if input.Loaded {
				parent := dataSourceTestClass("l3extOut", "l3_outside", ResourceArtifact)
				parent.PropertiesAll = []string{"name"}
				parent.Properties = map[string]*Property{"name": dataSourceTestProperty("name", "name", dataSourceTestValues("name_1", "name_2", "name_1", true))}
				store.Classes["l3extOut"] = parent
			}
			root.TestDependencies = root.getTestDependenciesFromDefinitions(definition.TestConfig.Dependencies, store, map[string]*TestDependency{}, 0)
			if !assert.NoError(t, store.ctx.Diagnostics.Error()) || !assert.Len(t, root.TestDependencies, 1) {
				return
			}
			assert.Equal(t, input.IgnoreChanges, root.TestDependencies[0].IgnoreChanges)
			for _, resourceTest := range []bool{false, true} {
				renderer := newTestConfigurationRenderer(&root, store)
				renderer.resourceTest = resourceTest
				configuration, err := renderer.dependenciesHCL()
				if !assert.NoError(t, err) {
					continue
				}
				assert.NoError(t, validateDataSourceTestConfiguration(configuration, nil))
				if testCase.Expected.(bool) {
					assert.Contains(t, configuration, "lifecycle {")
					assert.Contains(t, configuration, "ignore_changes = [relation_l3extrs_redistribute_pol]")
					assert.Contains(t, configuration, "standalone resources manage them")
				} else {
					assert.NotContains(t, configuration, "lifecycle {")
				}
			}
		})
	}
}

func TestDataSourceTestPropertyOverrides(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_attribute_name",
			Input:    map[string]string{"display_name": "attribute_value"},
			Expected: dataSourceDependencyRenderingExpected{Contains: []string{`display_name = "attribute_value"`}},
		},
		{
			Name:     "test_property_name",
			Input:    map[string]string{"displayName": "property_value"},
			Expected: dataSourceDependencyRenderingExpected{Contains: []string{`display_name = "property_value"`}},
		},
		{
			Name:     "test_attribute_name_takes_precedence",
			Input:    map[string]string{"display_name": "attribute_value", "displayName": "property_value"},
			Expected: dataSourceDependencyRenderingExpected{Contains: []string{`display_name = "attribute_value"`}, NotContains: []string{"property_value"}},
		},
		{
			Name:     "test_unknown_override",
			Input:    map[string]string{"unknown": "value"},
			Expected: dataSourceDependencyRenderingExpected{Error: `config override "unknown" does not match`},
		},
	}
	for _, blockType := range []string{"resource", "data"} {
		for _, testCase := range testCases {
			t.Run(blockType+"/"+testCase.Name, func(t *testing.T) {
				t.Parallel()
				expected := testCase.Expected.(dataSourceDependencyRenderingExpected)
				class := dataSourceTestClass("exampleClass", "example_class", ResourceArtifact, DatasourceArtifact)
				class.IdentifiedBy = []string{"displayName"}
				class.PropertiesAll = []string{"displayName"}
				class.Properties = map[string]*Property{"displayName": dataSourceTestProperty("displayName", "display_name", dataSourceTestValues("value_1", "value_2", "value_1", true))}
				renderer := newTestConfigurationRenderer(&class, &DataStore{})
				var contents strings.Builder
				address := terraformTestAddress{BlockType: blockType, Type: "aci_example_class", Name: "test_1"}
				err := renderer.writeClassBlock(&contents, address, &class, fullTestResourceScenario, testCase.Input.(map[string]string), nil, "", nil)
				if expected.Error != "" {
					if assert.Error(t, err) {
						assert.Contains(t, err.Error(), expected.Error)
					}
					return
				}
				if !assert.NoError(t, err) {
					return
				}
				for _, value := range expected.Contains {
					assert.Contains(t, contents.String(), value)
				}
				for _, value := range expected.NotContains {
					assert.NotContains(t, contents.String(), value)
				}
			})
		}
	}
}

func TestDataSourceTestPrefixSelection(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	enabled, disabled := true, false
	testCases := []test.TestCase{
		{Name: "test_free_name", Input: Property{PropertyName: "name", ValueType: String}, Expected: true},
		{Name: "test_property_exclusion", Input: Property{PropertyName: "name", ValueType: String, PrefixEnabled: &disabled}, Expected: false},
		{Name: "test_explicit_other_identifier", Input: Property{PropertyName: "label", ValueType: String, PrefixEnabled: &enabled}, Expected: true},
		{Name: "test_other_identifier_is_not_automatically_changed", Input: Property{PropertyName: "label", ValueType: String}, Expected: false},
		{Name: "test_non_identifier", Input: Property{PropertyName: "description", ValueType: String, PrefixEnabled: &enabled}, Expected: false},
		{Name: "test_ip_identifier", Input: Property{PropertyName: "address", ValueType: IpAddress, PrefixEnabled: &enabled}, Expected: false},
		{Name: "test_mac_identifier", Input: Property{PropertyName: "mac", ValueType: String}, Expected: false},
		{Name: "test_numeric_identifier", Input: Property{PropertyName: "name", ValueType: String, metaDetails: map[string]any{"uitype": "number"}}, Expected: false},
		{Name: "test_enum_identifier", Input: Property{PropertyName: "name", ValueType: String, ValidValues: ValidValues{"in": {LocalName: "in"}}}, Expected: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			property := testCase.Input.(Property)
			class := dataSourceTestClass("exampleRoot", "example_root", ResourceArtifact)
			class.IdentifiedBy = []string{"name", "label", "address", "mac"}
			renderer := newTestConfigurationRenderer(&class, &DataStore{})
			assert.Equal(t, testCase.Expected, renderer.shouldPrefixProperty(&class, &property))
		})
	}
}

func TestDataSourceTestPrefixConfiguration(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	enabled, disabled := true, false
	testCases := []test.TestCase{
		{
			Name:     "test_parent_root_is_prefixed_and_constrained_child_identifier_is_unchanged",
			Input:    (*bool)(nil),
			Expected: `format("%s-%s", var.test_prefix, "tenant_1")`,
		},
		{
			Name:     "test_dependency_can_exclude_automatic_prefix",
			Input:    &disabled,
			Expected: `name  = "tenant_1"`,
		},
		{
			Name:     "test_dependency_can_request_prefix_explicitly",
			Input:    &enabled,
			Expected: `format("%s-%s", var.test_prefix, "tenant_1")`,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			tenant := dataSourceTestClass("fvTenant", "tenant", ResourceArtifact)
			tenant.IdentifiedBy = []string{"name"}
			tenant.Properties = map[string]*Property{"name": dataSourceTestProperty("name", "name", dataSourceTestValues("tenant_1", "tenant_2", "tenant_1", true))}
			tenant.PropertiesAll = []string{"name"}
			class := dataSourceTestClass("exampleSubnet", "example_subnet", ResourceArtifact, DatasourceArtifact)
			class.IdentifiedBy = []string{"name", "ip"}
			class.Properties = map[string]*Property{
				"parentDn": dataSourceTestProperty("parentDn", "parent_dn", dataSourceTestReferenceValues("aci_tenant.test_1.id", "aci_tenant.test_1.id", "aci_tenant.test_1.id")),
				"name":     dataSourceTestProperty("name", "name", dataSourceTestValues("subnet_1", "subnet_2", "subnet_1", true)),
				"ip":       dataSourceTestProperty("ip", "ip", dataSourceTestValues("192.0.2.1/24", "192.0.2.2/24", "192.0.2.1/24", true)),
			}
			class.Properties["ip"].ValueType = IpAddress
			class.PropertiesAll = []string{"parentDn", "name", "ip"}
			class.TestDependencies = []*TestDependency{{Class: tenant.Name, Reference: "aci_tenant.test_1.id", ReferenceType: ResourceReference, Role: Parent, PrefixEnabled: testCase.Input.(*bool)}}
			store := &DataStore{Classes: map[string]Class{tenant.Name.String(): tenant, class.Name.String(): class}}
			dependencies, err := class.DataSourceTestDependenciesHCL(store)
			if !assert.NoError(t, err) {
				return
			}
			assert.Contains(t, dependencies, testCase.Expected)
			existing, err := class.DataSourceTestExistingHCL(store)
			if !assert.NoError(t, err) {
				return
			}
			file, diagnostics := hclsyntax.ParseConfig([]byte(existing), "fixture.tf", hcl.InitialPos)
			if !assert.False(t, diagnostics.HasErrors(), diagnostics.Error()) {
				return
			}
			for _, block := range file.Body.(*hclsyntax.Body).Blocks {
				name, diagnostics := block.Body.Attributes["name"].Expr.Value(nil)
				assert.False(t, diagnostics.HasErrors(), diagnostics.Error())
				assert.Equal(t, cty.StringVal("subnet_1"), name)
				ip, diagnostics := block.Body.Attributes["ip"].Expr.Value(nil)
				assert.False(t, diagnostics.HasErrors(), diagnostics.Error())
				assert.Equal(t, cty.StringVal("192.0.2.1/24"), ip)
			}
			assert.Equal(t, "tenant_1", tenant.Properties["name"].TestValues.Default.Entries[0].ConfigValue)
			assert.Equal(t, StringValue, tenant.Properties["name"].TestValues.Default.Entries[0].ValueType)
		})
	}
}

func TestDataSourceTestPrefixExpression(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	class := dataSourceTestClass("fvTenant", "tenant", ResourceArtifact, DatasourceArtifact)
	class.IdentifiedBy = []string{"name"}
	property := dataSourceTestProperty("name", "name", dataSourceTestValues("name_1", "name_2", "name_1", true))
	class.Properties = map[string]*Property{"name": property}
	class.PropertiesAll = []string{"name"}
	store := &DataStore{Classes: map[string]Class{class.Name.String(): class}}
	existing, err := class.DataSourceTestExistingHCL(store)
	if !assert.NoError(t, err) {
		return
	}
	file, diagnostics := hclsyntax.ParseConfig([]byte(existing), "fixture.tf", hcl.InitialPos)
	if !assert.False(t, diagnostics.HasErrors(), diagnostics.Error()) {
		return
	}
	testCases := []test.TestCase{
		{Name: "test_first_prefix", Input: "abc123", Expected: "abc123-name_1"},
		{Name: "test_second_prefix", Input: "xyz789", Expected: "xyz789-name_1"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			ctx := &hcl.EvalContext{
				Variables: map[string]cty.Value{"var": cty.ObjectVal(map[string]cty.Value{"test_prefix": cty.StringVal(testCase.Input.(string))})},
				Functions: map[string]function.Function{"format": stdlib.FormatFunc},
			}
			for _, block := range file.Body.(*hclsyntax.Body).Blocks {
				value, diagnostics := block.Body.Attributes["name"].Expr.Value(ctx)
				assert.False(t, diagnostics.HasErrors(), diagnostics.Error())
				assert.Equal(t, cty.StringVal(testCase.Expected.(string)), value)
			}
		})
	}
	assertions := class.DataSourceTestAssertions(store)
	if assert.Len(t, assertions, 1) {
		assert.True(t, assertions[0].IsAttributePair())
		assert.Equal(t, "aci_tenant.test_1", assertions[0].PairAddress)
	}
	assert.Equal(t, "name_1", property.TestValues.Create.Entries[0].ConfigValue)
}

func TestDataSourceTestPrefixExpressionOverride(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	disabled := false
	class := dataSourceTestClass("fvTenant", "tenant", ResourceArtifact, DatasourceArtifact)
	class.IdentifiedBy = []string{"name"}
	expression := `format("%sname1", var.test_prefix)`
	class.Properties = map[string]*Property{"name": {
		PropertyName: "name", AttributeName: "name", ValueType: String, PrefixEnabled: &disabled,
		TestValues: &TestValues{Create: dataSourceTestScenario(expression, expression, ExpressionValue, true)},
	}}
	class.PropertiesAll = []string{"name"}
	store := &DataStore{Classes: map[string]Class{class.Name.String(): class}}
	existing, err := class.DataSourceTestExistingHCL(store)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, strings.Count(existing, expression), "resource and datasource must use the explicit expression unchanged")
		assert.NotContains(t, existing, `format("%s-%s"`)
	}
	assertions := class.DataSourceTestAssertions(store)
	if assert.Len(t, assertions, 1) {
		assert.True(t, assertions[0].IsAttributePair(), "compare the evaluated name instead of the HCL expression text")
	}
}

func TestDataSourceTestDependencyRendering(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	tenant := dataSourceTestClass("fvTenant", "tenant", ResourceArtifact)
	tenant.IdentifiedBy = []string{"name"}
	tenant.Properties = map[string]*Property{
		"name": dataSourceTestProperty("name", "name", dataSourceTestValues("tenant", "tenant_2", "tenant", true)),
	}
	tenant.PropertiesAll = []string{"name"}

	applicationProfile := dataSourceTestClass("fvAp", "application_profile", ResourceArtifact)
	applicationProfile.Properties = map[string]*Property{
		"name": dataSourceTestProperty("name", "name", dataSourceTestValues("application", "application_2", "application", true)),
		"parentDn": dataSourceTestProperty("parentDn", "parent_dn", dataSourceTestReferenceValues(
			"aci_tenant.test.id",
			"aci_tenant.test.id",
			"aci_tenant.test.id",
		)),
		"descr": dataSourceTestProperty("descr", "description", &TestValues{
			Create:  dataSourceTestScenario("full description", "full description", StringValue, true),
			Default: dataSourceTestScenario("", "default description", StringValue, false),
		}),
	}
	applicationProfile.PropertiesAll = []string{"name", "parentDn", "descr"}
	applicationProfile.TestDependencies = []*TestDependency{{
		Class:         tenant.Name,
		Reference:     "aci_tenant.test.id",
		ReferenceType: ResourceReference,
		Role:          Parent,
		OrderOnly:     true,
	}}

	frameworkSubject := dataSourceTestClass("exampleFramework", "example_framework", DatasourceArtifact)
	frameworkSubject.Properties = map[string]*Property{
		"parentDn": dataSourceTestProperty("parentDn", "parent_dn", dataSourceTestReferenceValues(
			"aci_application_profile.test.id",
			"aci_application_profile.test.id",
			"aci_application_profile.test.id",
		)),
	}
	frameworkSubject.PropertiesAll = []string{"parentDn"}
	frameworkSubject.TestDependencies = []*TestDependency{{
		Class:         applicationProfile.Name,
		Reference:     "aci_application_profile.test.id",
		ReferenceType: ResourceReference,
		Role:          Parent,
		Dependencies:  applicationProfile.TestDependencies,
	}}
	frameworkStore := &DataStore{Classes: map[string]Class{
		tenant.Name.String():             tenant,
		applicationProfile.Name.String(): applicationProfile,
		frameworkSubject.Name.String():   frameworkSubject,
	}}

	legacySubject := dataSourceTestClass("exampleLegacy", "example_legacy", DatasourceArtifact)
	legacySubject.Properties = map[string]*Property{
		"parentDn": dataSourceTestProperty("parentDn", "parent_dn", dataSourceTestReferenceValues(
			"aci_contract.test.id",
			"aci_contract.test.id",
			"aci_contract.test.id",
		)),
	}
	legacySubject.PropertiesAll = []string{"parentDn"}
	legacySubject.TestDependencies = []*TestDependency{{
		Class:         testClassName("vzBrCP"),
		Reference:     "aci_contract.test.id",
		ReferenceType: ResourceReference,
		Role:          Parent,
		ConfigOverrides: map[string]string{
			"name":      "contract",
			"tenant_dn": "aci_tenant.test.id",
		},
		Dependencies: []*TestDependency{{
			Class:         tenant.Name,
			Reference:     "aci_tenant.test.id",
			ReferenceType: ResourceReference,
			OrderOnly:     true,
		}},
	}}
	legacyStore := &DataStore{Classes: map[string]Class{
		tenant.Name.String():        tenant,
		legacySubject.Name.String(): legacySubject,
	}}

	incompleteLegacySubject := legacySubject
	incompleteLegacySubject.Name = testClassName("exampleIncompleteLegacy")
	incompleteLegacySubject.TestDependencies = []*TestDependency{{
		Class:         testClassName("vzBrCP"),
		Reference:     "aci_contract.test.id",
		ReferenceType: ResourceReference,
		Role:          Parent,
	}}
	incompleteLegacyStore := &DataStore{Classes: map[string]Class{
		incompleteLegacySubject.Name.String(): incompleteLegacySubject,
	}}

	testCases := []test.TestCase{
		{
			Name:  "test_framework_dependency_uses_required_values_and_infers_its_parent",
			Input: dataSourceDependencyRenderingInput{Class: frameworkSubject, DataStore: frameworkStore},
			Expected: dataSourceDependencyRenderingExpected{
				Contains: []string{
					`resource "aci_tenant" "test"`,
					`resource "aci_application_profile" "test"`,
					`count      = 1`,
					`name       = "application"`,
					`parent_dn  = aci_tenant.test[0].id`,
					`depends_on = [aci_tenant.test]`,
				},
				NotContains: []string{"description"},
			},
		},
		{
			Name:  "test_legacy_dependency_uses_complete_config_overrides",
			Input: dataSourceDependencyRenderingInput{Class: legacySubject, DataStore: legacyStore},
			Expected: dataSourceDependencyRenderingExpected{
				Contains: []string{
					`resource "aci_tenant" "test"`,
					`resource "aci_contract" "test"`,
					`count      = 1`,
					`name       = "contract"`,
					`tenant_dn  = aci_tenant.test[0].id`,
					`depends_on = [aci_tenant.test]`,
				},
			},
		},
		{
			Name:  "test_legacy_dependency_requires_complete_config_overrides",
			Input: dataSourceDependencyRenderingInput{Class: incompleteLegacySubject, DataStore: incompleteLegacyStore},
			Expected: dataSourceDependencyRenderingExpected{
				Error: `legacy class "vzBrCP" requires complete config_overrides`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(dataSourceDependencyRenderingInput)
			expected := testCase.Expected.(dataSourceDependencyRenderingExpected)
			actual, err := input.Class.DataSourceTestDependenciesHCL(input.DataStore)
			if expected.Error != "" {
				if err == nil || !strings.Contains(err.Error(), expected.Error) {
					t.Fatalf("expected error containing %q, got %v", expected.Error, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("render dependencies: %v", err)
			}
			for _, value := range expected.Contains {
				assert.Contains(t, actual, value, test.MessageContains(actual, value, testCase.Name))
			}
			for _, value := range expected.NotContains {
				assert.NotContains(t, actual, value, test.MessageNotContains(actual, value, testCase.Name))
			}
		})
	}
}

func TestDataSourceTestCountedDataSourceDependency(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	tenant := dataSourceTestClass("fvTenant", "tenant", ResourceArtifact, DatasourceArtifact)
	tenant.IdentifiedBy = []string{"name"}
	tenant.Properties = map[string]*Property{
		"name": dataSourceTestProperty("name", "name", dataSourceTestValues("tenant_1", "tenant_2", "tenant_1", true)),
	}
	tenant.PropertiesAll = []string{"name"}

	class := dataSourceTestClass("exampleDataDependency", "example_data_dependency", ResourceArtifact, DatasourceArtifact)
	class.Properties = map[string]*Property{
		"parentDn": dataSourceTestProperty("parentDn", "parent_dn", dataSourceTestReferenceValues(
			"data.aci_tenant.test_1.id",
			"data.aci_tenant.test_1.id",
			"data.aci_tenant.test_1.id",
		)),
	}
	class.PropertiesAll = []string{"parentDn"}
	class.TestDependencies = []*TestDependency{{
		Class:         tenant.Name,
		Reference:     "data.aci_tenant.test_1.id",
		ReferenceType: DataSourceReference,
		Role:          Parent,
	}}
	store := &DataStore{Classes: map[string]Class{
		class.Name.String():  class,
		tenant.Name.String(): tenant,
	}}

	dependencies, err := class.DataSourceTestDependenciesHCL(store)
	if err != nil {
		t.Fatalf("render dependencies: %v", err)
	}
	existing, err := class.DataSourceTestExistingHCL(store)
	if err != nil {
		t.Fatalf("render existing configuration: %v", err)
	}
	assertions := class.DataSourceTestAssertions(store)

	testCases := []test.TestCase{
		{Name: "test_dependency_block", Input: dependencies, Expected: `data "aci_tenant" "test_1" {`},
		{Name: "test_dependency_count", Input: dependencies, Expected: `count = 1`},
		{Name: "test_dependency_reference", Input: existing, Expected: `parent_dn = data.aci_tenant.test_1[0].id`},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			actual := testCase.Input.(string)
			expected := testCase.Expected.(string)
			assert.Contains(t, actual, expected, test.MessageContains(actual, expected, testCase.Name))
		})
	}
	if assert.Len(t, assertions, 1) {
		assert.Equal(t, "data.aci_tenant.test_1.0", assertions[0].PairAddress)
	}
}

func TestDataSourceTestExistingAndMissingConfiguration(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	child := dataSourceTestClass("exampleChild", "example_child")
	child.ResourceNameNested = "children"
	child.Properties = map[string]*Property{
		"key": dataSourceTestProperty("key", "key", nil),
	}
	child.PropertiesAll = []string{"key"}

	identified := dataSourceTestClass("exampleIdentified", "example_identified", ResourceArtifact, DatasourceArtifact)
	identified.IdentifiedBy = []string{"name"}
	identified.Properties = map[string]*Property{
		"name":  dataSourceTestProperty("name", "name", dataSourceTestValues("created", "missing", "created", true)),
		"descr": dataSourceTestProperty("descr", "description", dataSourceTestValues("full description", "updated description", "", false)),
	}
	identified.PropertiesAll = []string{"name", "descr"}
	identified.TestChildren = []*TestChild{{
		Class: child.Name,
		Instances: []TestChildInstance{{Properties: map[string]TestValueScenario{
			"key": dataSourceTestScenario("child key", "child key", StringValue, true),
		}}},
	}}
	identifiedStore := &DataStore{Classes: map[string]Class{
		identified.Name.String(): identified,
		child.Name.String():      child,
	}}

	tenant := dataSourceTestClass("fvTenant", "tenant", ResourceArtifact)
	tenant.IdentifiedBy = []string{"name"}
	tenant.Properties = map[string]*Property{
		"name": dataSourceTestProperty("name", "name", dataSourceTestValues("tenant", "tenant_2", "tenant", true)),
	}
	tenant.PropertiesAll = []string{"name"}

	constantRN := dataSourceTestClass("exampleConstantRN", "example_constant_rn", ResourceArtifact, DatasourceArtifact)
	constantRN.Properties = map[string]*Property{
		"parentDn": dataSourceTestProperty("parentDn", "parent_dn", dataSourceTestReferenceValues(
			"aci_tenant.test.id",
			"aci_tenant.test.id",
			"aci_tenant.test_2.id",
		)),
	}
	constantRN.PropertiesAll = []string{"parentDn"}
	constantRN.TestDependencies = []*TestDependency{
		{Class: tenant.Name, Reference: "aci_tenant.test.id", ReferenceType: ResourceReference, Role: Parent},
		{Class: tenant.Name, Reference: "aci_tenant.test_2.id", ReferenceType: ResourceReference, Role: Parent},
	}
	constantRNStore := &DataStore{Classes: map[string]Class{
		tenant.Name.String():     tenant,
		constantRN.Name.String(): constantRN,
	}}
	explicitParentDN := constantRN
	explicitParentDN.Properties = map[string]*Property{
		"parentDn": dataSourceTestProperty("parentDn", "parent_dn", dataSourceTestReferenceValues(
			"aci_tenant.test.id",
			"aci_tenant.test.id",
			"aci_tenant.test_2.id",
		)),
	}
	explicitParentDN.Properties["parentDn"].TestValues.DataSourceMissing = dataSourceTestScenario("uni/tn-non_existing", "uni/tn-non_existing", StringValue, true)
	explicitParentDNStore := &DataStore{Classes: map[string]Class{
		tenant.Name.String():           tenant,
		explicitParentDN.Name.String(): explicitParentDN,
	}}

	defaultParent := dataSourceTestClass("exampleDefaultParent", "example_default_parent", ResourceArtifact, DatasourceArtifact)
	defaultParent.IdentifiedBy = []string{"name"}
	defaultParent.Properties = map[string]*Property{
		"name": dataSourceTestProperty("name", "name", dataSourceTestValues("created", "missing", "created", true)),
		"parentDn": {
			PropertyName:  "parentDn",
			AttributeName: "parent_dn",
			ValueType:     String,
			Documentation: PropertyDocumentation{DefaultValues: []DefaultValue{{Value: "uni"}}},
		},
	}
	defaultParent.PropertiesAll = []string{"name", "parentDn"}
	defaultParentStore := &DataStore{Classes: map[string]Class{defaultParent.Name.String(): defaultParent}}

	sameIdentifier := dataSourceTestClass("exampleSameIdentifier", "example_same_identifier", ResourceArtifact, DatasourceArtifact)
	sameIdentifier.IdentifiedBy = []string{"name"}
	sameIdentifier.Properties = map[string]*Property{
		"name": dataSourceTestProperty("name", "name", dataSourceTestValues("name_1", "name_1", "name_1", true)),
	}
	sameIdentifier.PropertiesAll = []string{"name"}
	sameIdentifierStore := &DataStore{Classes: map[string]Class{sameIdentifier.Name.String(): sameIdentifier}}

	numericIdentifier := dataSourceTestClass("exampleNumericIdentifier", "example_numeric_identifier", ResourceArtifact, DatasourceArtifact)
	numericIdentifier.IdentifiedBy = []string{"siteId"}
	siteID := dataSourceTestProperty("siteId", "site_id", dataSourceTestValues("1", "1", "1", true))
	siteID.metaDetails = map[string]any{"uitype": "number"}
	siteID.Validators = []Validator{{Min: 0, Max: 1000}}
	numericIdentifier.Properties = map[string]*Property{"siteId": siteID}
	numericIdentifier.PropertiesAll = []string{"siteId"}
	numericIdentifierStore := &DataStore{Classes: map[string]Class{numericIdentifier.Name.String(): numericIdentifier}}

	dnIdentifier := dataSourceTestClass("exampleDNIdentifier", "example_dn_identifier", ResourceArtifact, DatasourceArtifact)
	dnIdentifier.IdentifiedBy = []string{"tDn"}
	dnIdentifier.Properties = map[string]*Property{
		"tDn": dataSourceTestProperty("tDn", "target_dn", dataSourceTestValues("uni/tn-name/object-name_1", "uni/tn-name/object-name_1", "uni/tn-name/object-name_1", true)),
	}
	dnIdentifier.PropertiesAll = []string{"tDn"}
	dnIdentifierStore := &DataStore{Classes: map[string]Class{dnIdentifier.Name.String(): dnIdentifier}}

	semanticDNIdentifier := dataSourceTestClass("exampleSemanticDNIdentifier", "example_semantic_dn_identifier", ResourceArtifact, DatasourceArtifact)
	semanticDNIdentifier.IdentifiedBy = []string{"ldev"}
	logicalDevice := dataSourceTestProperty("ldev", "logical_device", dataSourceTestValues("uni/tn-name/lDevVip-device_1", "logical_device_2", "uni/tn-name/lDevVip-device_1", true))
	logicalDevice.Reference = &PropertyReference{ExampleAttributeName: "id"}
	semanticDNIdentifier.Properties = map[string]*Property{"ldev": logicalDevice}
	semanticDNIdentifier.PropertiesAll = []string{"ldev"}
	semanticDNIdentifierStore := &DataStore{Classes: map[string]Class{semanticDNIdentifier.Name.String(): semanticDNIdentifier}}

	explicitDNIdentifier := dataSourceTestClass("exampleExplicitDNIdentifier", "example_explicit_dn_identifier", ResourceArtifact, DatasourceArtifact)
	explicitDNIdentifier.IdentifiedBy = []string{"tDn"}
	explicitTarget := dataSourceTestProperty("tDn", "target_dn", dataSourceTestValues("uni/tn-name/object-name_1", "uni/tn-name/object-name_1", "uni/tn-name/object-name_1", true))
	explicitTarget.TestValues.DataSourceMissing = dataSourceTestScenario("uni/tn-name/object-missing", "uni/tn-name/object-missing", StringValue, true)
	explicitDNIdentifier.Properties = map[string]*Property{"tDn": explicitTarget}
	explicitDNIdentifier.PropertiesAll = []string{"tDn"}
	explicitDNIdentifierStore := &DataStore{Classes: map[string]Class{explicitDNIdentifier.Name.String(): explicitDNIdentifier}}

	testCases := []test.TestCase{
		{
			Name:  "test_identified_object_uses_full_resource_and_changed_identifier",
			Input: dataSourceDependencyRenderingInput{Class: identified, DataStore: identifiedStore},
			Expected: dataSourceDependencyRenderingExpected{
				Contains: []string{
					`resource "aci_example_identified" "test_1"`,
					`description = "full description"`,
					`children = [`,
					`key = "child key"`,
					`data "aci_example_identified" "test"`,
					`depends_on = [aci_example_identified.test_1]`,
					`data "aci_example_identified" "test_non_existing"`,
					`name = format("%s-%s", var.test_prefix, "missing")`,
				},
			},
		},
		{
			Name:  "test_constant_rn_uses_second_valid_parent_for_missing_lookup",
			Input: dataSourceDependencyRenderingInput{Class: constantRN, DataStore: constantRNStore},
			Expected: dataSourceDependencyRenderingExpected{
				Contains: []string{
					`resource "aci_tenant" "test_2"`,
					`data "aci_example_constant_rn" "test_non_existing"`,
					`parent_dn = aci_tenant.test_2[0].id`,
				},
			},
		},
		{
			Name:  "test_constant_rn_prefers_explicit_missing_parent_without_unused_dependency",
			Input: dataSourceDependencyRenderingInput{Class: explicitParentDN, DataStore: explicitParentDNStore},
			Expected: dataSourceDependencyRenderingExpected{
				Contains: []string{
					`data "aci_example_constant_rn" "test_non_existing"`,
					`parent_dn = "uni/tn-non_existing"`,
				},
				NotContains: []string{`resource "aci_tenant" "test_2"`},
			},
		},
		{
			Name:  "test_schema_default_parent_is_omitted",
			Input: dataSourceDependencyRenderingInput{Class: defaultParent, DataStore: defaultParentStore},
			Expected: dataSourceDependencyRenderingExpected{
				NotContains: []string{"parent_dn"},
			},
		},
		{
			Name:  "test_equal_string_identifier_derives_canonical_missing_value",
			Input: dataSourceDependencyRenderingInput{Class: sameIdentifier, DataStore: sameIdentifierStore},
			Expected: dataSourceDependencyRenderingExpected{Contains: []string{
				`data "aci_example_same_identifier" "test_non_existing"`,
				`name = format("%s-%s", var.test_prefix, "name_2")`,
			}},
		},
		{
			Name:  "test_equal_numeric_identifier_advances_within_range",
			Input: dataSourceDependencyRenderingInput{Class: numericIdentifier, DataStore: numericIdentifierStore},
			Expected: dataSourceDependencyRenderingExpected{Contains: []string{
				`data "aci_example_numeric_identifier" "test_non_existing"`,
				`site_id = "2"`,
			}},
		},
		{
			Name:  "test_equal_dn_identifier_does_not_invent_missing_dn",
			Input: dataSourceDependencyRenderingInput{Class: dnIdentifier, DataStore: dnIdentifierStore},
			Expected: dataSourceDependencyRenderingExpected{NotContains: []string{
				`data "aci_example_dn_identifier" "test_non_existing"`,
			}},
		},
		{
			Name:  "test_semantic_dn_identifier_does_not_use_synthetic_update_value",
			Input: dataSourceDependencyRenderingInput{Class: semanticDNIdentifier, DataStore: semanticDNIdentifierStore},
			Expected: dataSourceDependencyRenderingExpected{NotContains: []string{
				`data "aci_example_semantic_dn_identifier" "test_non_existing"`,
			}},
		},
		{
			Name:  "test_dn_identifier_uses_explicit_missing_value",
			Input: dataSourceDependencyRenderingInput{Class: explicitDNIdentifier, DataStore: explicitDNIdentifierStore},
			Expected: dataSourceDependencyRenderingExpected{Contains: []string{
				`data "aci_example_explicit_dn_identifier" "test_non_existing"`,
				`target_dn = "uni/tn-name/object-missing"`,
			}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(dataSourceDependencyRenderingInput)
			expected := testCase.Expected.(dataSourceDependencyRenderingExpected)
			dependencies, err := input.Class.DataSourceTestDependenciesHCL(input.DataStore)
			if err != nil {
				t.Fatalf("render dependencies: %v", err)
			}
			existing, err := input.Class.DataSourceTestExistingHCL(input.DataStore)
			if err != nil {
				t.Fatalf("render existing configuration: %v", err)
			}
			missing, err := input.Class.DataSourceTestMissingHCL(input.DataStore)
			if err != nil {
				t.Fatalf("render missing configuration: %v", err)
			}
			actual := dependencies + existing + missing
			for _, value := range expected.Contains {
				assert.Contains(t, actual, value, test.MessageContains(actual, value, testCase.Name))
			}
			for _, value := range expected.NotContains {
				assert.NotContains(t, actual, value, test.MessageNotContains(actual, value, testCase.Name))
			}
		})
	}
}

func TestDataSourceTestVersionedConfiguration(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	classVersions, err := NewVersions("1.0(1e)-")
	if err != nil {
		t.Fatal(err)
	}
	scalarVersions, err := NewVersions("2.0(1a)-")
	if err != nil {
		t.Fatal(err)
	}
	singletonVersions, err := NewVersions("3.0(1a)-")
	if err != nil {
		t.Fatal(err)
	}
	nestedPropertyVersions, err := NewVersions("4.0(1a)-")
	if err != nil {
		t.Fatal(err)
	}
	repeatedVersions, err := NewVersions("5.0(1a)-")
	if err != nil {
		t.Fatal(err)
	}

	singleton := dataSourceTestClass("exampleSingleton", "example_singleton")
	singleton.ResourceNameNested = "singleton"
	singleton.IsSingleNestedWhenDefinedAsChild = true
	singleton.SupportedVersions = singletonVersions
	singleton.Properties = map[string]*Property{
		"value": {
			PropertyName:      "value",
			AttributeName:     "value",
			ValueType:         String,
			SupportedVersions: nestedPropertyVersions,
		},
	}
	singleton.PropertiesAll = []string{"value"}

	repeated := dataSourceTestClass("exampleRepeated", "example_repeated")
	repeated.ResourceNameNested = "repeated"
	repeated.SupportedVersions = repeatedVersions
	repeated.Properties = map[string]*Property{
		"key": dataSourceTestProperty("key", "key", nil),
	}
	repeated.PropertiesAll = []string{"key"}

	class := dataSourceTestClass("exampleVersioned", "example_versioned", ResourceArtifact, DatasourceArtifact)
	class.SupportedVersions = classVersions
	class.Properties = map[string]*Property{
		"name": dataSourceTestProperty("name", "name", dataSourceTestValues("name_1", "name_2", "name_1", true)),
		"versioned": {
			PropertyName:      "versioned",
			AttributeName:     "versioned",
			ValueType:         String,
			SupportedVersions: scalarVersions,
			TestValues:        &TestValues{Create: dataSourceTestScenario("value_1", "value_1", StringValue, true)},
		},
	}
	class.PropertiesAll = []string{"name", "versioned"}
	class.TestChildren = []*TestChild{
		{
			Class: singleton.Name,
			Instances: []TestChildInstance{{Properties: map[string]TestValueScenario{
				"value": dataSourceTestScenario("nested_value_1", "nested_value_1", StringValue, true),
			}}},
		},
		{
			Class: repeated.Name,
			Instances: []TestChildInstance{{Properties: map[string]TestValueScenario{
				"key": dataSourceTestScenario("key_1", "key_1", StringValue, true),
			}}},
		},
	}
	store := &DataStore{Classes: map[string]Class{
		class.Name.String():     class,
		singleton.Name.String(): singleton,
		repeated.Name.String():  repeated,
	}}

	dependencies, err := class.DataSourceTestDependenciesHCL(store)
	if err != nil {
		t.Fatalf("render dependencies: %v", err)
	}
	existing, err := class.DataSourceTestExistingHCL(store)
	if err != nil {
		t.Fatalf("render existing configuration: %v", err)
	}
	actual := dependencies + existing
	testCases := []test.TestCase{
		{Name: "test_scalar_property", Expected: `versioned = provider::aci::compare_versions(var.apic_version, "inside", "2.0(1a)-") ? "value_1" : null`},
		{Name: "test_singleton_child", Expected: `singleton = provider::aci::compare_versions(var.apic_version, "inside", "3.0(1a)-") ? {`},
		{Name: "test_nested_property", Expected: `value = provider::aci::compare_versions(var.apic_version, "inside", "4.0(1a)-") ? "nested_value_1" : null`},
		{Name: "test_repeated_child", Expected: `repeated = provider::aci::compare_versions(var.apic_version, "inside", "5.0(1a)-") ? [`},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			expected := testCase.Expected.(string)
			assert.Contains(t, actual, expected, test.MessageContains(actual, expected, testCase.Name))
		})
	}
}

func TestDataSourceTestVersionedDependency(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	classVersions, err := NewVersions("1.0(1e)-")
	if err != nil {
		t.Fatal(err)
	}
	dependencyVersions, err := NewVersions("3.0(1a)-")
	if err != nil {
		t.Fatal(err)
	}

	dependency := dataSourceTestClass("exampleDependency", "example_dependency", ResourceArtifact)
	dependency.SupportedVersions = dependencyVersions
	dependency.Properties = map[string]*Property{
		"name": dataSourceTestProperty("name", "name", dataSourceTestValues("dependency_1", "dependency_2", "dependency_1", true)),
	}
	dependency.PropertiesAll = []string{"name"}

	class := dataSourceTestClass("exampleVersionedDependency", "example_versioned_dependency", ResourceArtifact, DatasourceArtifact)
	class.SupportedVersions = classVersions
	target := dataSourceTestProperty("target", "target", dataSourceTestReferenceValues(
		"aci_example_dependency.test_1.id",
		"aci_example_dependency.test_1.id",
		"aci_example_dependency.test_1.id",
	))
	target.SupportedVersions = dependencyVersions
	class.Properties = map[string]*Property{"target": target}
	class.PropertiesAll = []string{"target"}
	class.TestDependencies = []*TestDependency{{
		Class:         dependency.Name,
		Reference:     "aci_example_dependency.test_1.id",
		ReferenceType: ResourceReference,
		Role:          Target,
	}}
	store := &DataStore{Classes: map[string]Class{
		class.Name.String():      class,
		dependency.Name.String(): dependency,
	}}

	dependencies, err := class.DataSourceTestDependenciesHCL(store)
	if err != nil {
		t.Fatalf("render dependencies: %v", err)
	}
	existing, err := class.DataSourceTestExistingHCL(store)
	if err != nil {
		t.Fatalf("render existing configuration: %v", err)
	}
	assertions := class.DataSourceTestAssertions(store)

	testCases := []test.TestCase{
		{Name: "test_conditional_dependency_count", Input: dependencies, Expected: `count = provider::aci::compare_versions(var.apic_version, "inside", "3.0(1a)-") ? 1 : 0`},
		{Name: "test_counted_dependency_reference", Input: existing, Expected: `target = provider::aci::compare_versions(var.apic_version, "inside", "3.0(1a)-") ? aci_example_dependency.test_1[0].id : null`},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			actual := testCase.Input.(string)
			expected := testCase.Expected.(string)
			assert.Contains(t, actual, expected, test.MessageContains(actual, expected, testCase.Name))
		})
	}
	if assert.Len(t, assertions, 1) {
		assert.Equal(t, "aci_example_dependency.test_1.0", assertions[0].PairAddress)
		assert.Equal(t, []string{"3.0(1a)-"}, assertions[0].Versions)
	}
}

func TestAssociatedSiteDataSourceTestConfiguration(t *testing.T) {
	test.InitializeTest(t)

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("determine data test path")
	}
	t.Chdir(filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..")))

	ctx := NewContext()
	store, err := NewDataStore(ctx)
	if err != nil {
		t.Fatalf("initialize data store: %v", err)
	}
	if err := ctx.Diagnostics.Error(); err != nil {
		t.Fatalf("initialize test data: %v", err)
	}
	class := store.Classes["fvSiteAssociated"]
	dependencies, err := class.DataSourceTestDependenciesHCL(store)
	if err != nil {
		t.Fatalf("render dependencies: %v", err)
	}
	assert.Contains(t, dependencies, `resource "aci_bridge_domain" "test_1"`)
	assert.Contains(t, dependencies, `resource "aci_bridge_domain" "test_2"`)
	assert.Equal(t, 1, strings.Count(dependencies, `resource "aci_tenant" "test_1"`))

	existing, err := class.DataSourceTestExistingHCL(store)
	if err != nil {
		t.Fatalf("render existing configuration: %v", err)
	}
	missing, err := class.DataSourceTestMissingHCL(store)
	if err != nil {
		t.Fatalf("render missing configuration: %v", err)
	}
	assert.NotEmpty(t, missing, "associated site has no missing-object lookup")
	testCases := []test.TestCase{
		{
			Name:  "test_create_and_read_under_first_bridge_domain",
			Input: existing,
			Expected: dataSourceDependencyRenderingExpected{
				Contains: []string{
					`resource "aci_associated_site" "test_1"`,
					`data "aci_associated_site" "test"`,
					"aci_bridge_domain.test_1[0].id",
				},
				NotContains: []string{"aci_bridge_domain.test_2[0].id"},
			},
		},
		{
			Name:  "test_missing_lookup_under_second_bridge_domain_without_creating_stasc",
			Input: missing,
			Expected: dataSourceDependencyRenderingExpected{
				Contains: []string{
					`data "aci_associated_site" "test_non_existing"`,
					"aci_bridge_domain.test_2[0].id",
				},
				NotContains: []string{`resource "aci_associated_site"`, "aci_bridge_domain.test_1[0].id"},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			configuration := testCase.Input.(string)
			expected := testCase.Expected.(dataSourceDependencyRenderingExpected)
			for _, value := range expected.Contains {
				assert.Contains(t, configuration, value)
			}
			for _, value := range expected.NotContains {
				assert.NotContains(t, configuration, value)
			}
			assert.NoError(t, validateDataSourceTestConfiguration(dependencies+"\n"+configuration, nil))
		})
	}
}

func TestDataSourceSharedContractDependencies(t *testing.T) {
	test.InitializeTest(t)

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("determine data test path")
	}
	t.Chdir(filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..")))
	ctx := NewContext()
	store, err := NewDataStore(ctx)
	if !assert.NoError(t, err) || !assert.NoError(t, ctx.Diagnostics.Error()) {
		return
	}
	testCases := []test.TestCase{
		{Name: "test_esg_relations_share_contracts", Input: "fvESg", Expected: []string{"test_1", "test_2"}},
		{Name: "test_epg_relations_share_contracts", Input: "fvAEPg", Expected: []string{"test_1", "test_2"}},
		{Name: "test_external_epg_relations_share_contracts", Input: "l3extInstP", Expected: []string{"test_1", "test_2"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			class := store.Classes[testCase.Input.(string)]
			dependencies, err := class.DataSourceTestDependenciesHCL(store)
			if !assert.NoError(t, err) {
				return
			}
			file, diagnostics := hclsyntax.ParseConfig([]byte(dependencies), "dependencies.tf", hcl.InitialPos)
			if !assert.False(t, diagnostics.HasErrors(), diagnostics.Error()) {
				return
			}
			var labels []string
			names := make(map[string]string)
			for _, block := range file.Body.(*hclsyntax.Body).Blocks {
				if block.Type != "resource" || block.Labels[0] != "aci_contract" {
					continue
				}
				labels = append(labels, block.Labels[1])
				value, diagnostics := block.Body.Attributes["name"].Expr.Value(nil)
				if !assert.False(t, diagnostics.HasErrors(), diagnostics.Error()) {
					return
				}
				name := value.AsString()
				assert.NotContains(t, names, name, "contract %s repeats the APIC name used by %s", block.Labels[1], names[name])
				names[name] = block.Labels[1]
			}
			assert.Equal(t, testCase.Expected, labels)
			existing, err := class.DataSourceTestExistingHCL(store)
			if !assert.NoError(t, err) {
				return
			}
			assert.NotContains(t, existing, "aci_contract.test_3")
			assert.NotContains(t, existing, "aci_contract.test_4")
			assert.NoError(t, validateDataSourceTestConfiguration(dependencies+"\n"+existing, class.DataSourceTestAssertions(store)))
		})
	}
}

func TestAllDataSourceTestConfigurationsAreSelfContained(t *testing.T) {
	test.InitializeTest(t)

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("determine data test path")
	}
	t.Chdir(filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..")))

	ctx := NewContext()
	store, err := NewDataStore(ctx)
	if err != nil {
		t.Fatalf("initialize data store: %v", err)
	}
	if err := ctx.Diagnostics.Error(); err != nil {
		t.Fatalf("initialize test data: %v", err)
	}

	classNames := make([]string, 0, len(store.Classes))
	for className := range store.Classes {
		classNames = append(classNames, className)
	}
	sort.Strings(classNames)
	for _, className := range classNames {
		class := store.Classes[className]
		if !class.HasDatasourceArtifact() || slices.Contains(class.TestConfig.IgnoreTests, DatasourceIgnoreTest) {
			continue
		}
		t.Run(className, func(t *testing.T) {
			dependencies, err := class.DataSourceTestDependenciesHCL(store)
			if err != nil {
				t.Fatalf("render dependencies: %v", err)
			}
			if err := validateDataSourceTestDependencyCounts(dependencies); err != nil {
				t.Fatalf("validate dependency counts: %v", err)
			}
			existing, err := class.DataSourceTestExistingHCL(store)
			if err != nil {
				t.Fatalf("render existing configuration: %v", err)
			}
			if err := validateDataSourceTestConfiguration(dependencies+"\n"+existing, class.DataSourceTestAssertions(store)); err != nil {
				t.Fatalf("validate existing configuration: %v", err)
			}

			missing, err := class.DataSourceTestMissingHCL(store)
			if err != nil {
				t.Fatalf("render missing configuration: %v", err)
			}
			if missing != "" {
				if err := validateDataSourceTestConfiguration(dependencies+"\n"+missing, nil); err != nil {
					t.Fatalf("validate missing configuration: %v", err)
				}
			}
		})
	}
}

func validateDataSourceTestDependencyCounts(configuration string) error {
	if strings.TrimSpace(configuration) == "" {
		return nil
	}
	file, diagnostics := hclsyntax.ParseConfig([]byte(configuration), "generated.tf", hcl.InitialPos)
	if diagnostics.HasErrors() {
		return diagnostics
	}
	for _, block := range file.Body.(*hclsyntax.Body).Blocks {
		if len(block.Labels) != 2 || block.Type != "resource" && block.Type != "data" {
			continue
		}
		if _, ok := block.Body.Attributes["count"]; !ok {
			return fmt.Errorf("dependency %q has no count", strings.Join(block.Labels, "."))
		}
	}
	return nil
}

func validateDataSourceTestConfiguration(configuration string, assertions []TestAssertion) error {
	file, diagnostics := hclsyntax.ParseConfig([]byte(configuration), "generated.tf", hcl.InitialPos)
	if diagnostics.HasErrors() {
		return diagnostics
	}
	body := file.Body.(*hclsyntax.Body)
	declared := make(map[string]bool)
	counted := make(map[string]bool)
	for _, block := range body.Blocks {
		if len(block.Labels) != 2 || (block.Type != "resource" && block.Type != "data") {
			continue
		}
		address := block.Labels[0] + "." + block.Labels[1]
		if block.Type == "data" {
			address = "data." + address
		}
		if declared[address] {
			return fmt.Errorf("duplicate Terraform object %q", address)
		}
		declared[address] = true
		if _, hasCount := block.Body.Attributes["count"]; hasCount {
			declared[address+".0"] = true
			counted[address] = true
		}
	}

	missing := make(map[string]bool)
	uncounted := make(map[string]bool)
	selfReferences := make(map[string]bool)
	var walk func(*hclsyntax.Body, string)
	walk = func(current *hclsyntax.Body, owner string) {
		for attributeName, attribute := range current.Attributes {
			for _, traversal := range attribute.Expr.Variables() {
				address := dataSourceTestTraversalAddress(traversal)
				if address != "" && address == owner {
					selfReferences[address] = true
				}
				if address != "" && !declared[address] {
					missing[address] = true
				}
				if attributeName != "depends_on" && counted[address] && !dataSourceTestTraversalHasIndex(traversal) {
					uncounted[address] = true
				}
			}
		}
		for _, block := range current.Blocks {
			// Removed targets refer to prior state and deliberately have no resource block here.
			if block.Type == "removed" {
				continue
			}
			blockOwner := owner
			if len(block.Labels) == 2 && (block.Type == "resource" || block.Type == "data") {
				blockOwner = strings.Join(block.Labels, ".")
				if block.Type == "data" {
					blockOwner = "data." + blockOwner
				}
			}
			walk(block.Body, blockOwner)
		}
	}
	walk(body, "")
	if len(selfReferences) > 0 {
		addresses := make([]string, 0, len(selfReferences))
		for address := range selfReferences {
			addresses = append(addresses, address)
		}
		sort.Strings(addresses)
		return fmt.Errorf("self-referential Terraform objects: %s", strings.Join(addresses, ", "))
	}
	for _, assertion := range assertions {
		if (assertion.IsAttributePair() || assertion.IsSetElementPair()) && !declared[assertion.PairAddress] {
			missing[assertion.PairAddress] = true
		}
		for _, value := range assertion.NestedValues {
			if value.PairAddress != "" && !declared[value.PairAddress] {
				missing[value.PairAddress] = true
			}
		}
		for _, parent := range assertion.Parents {
			for _, value := range parent.Values {
				if value.PairAddress != "" && !declared[value.PairAddress] {
					missing[value.PairAddress] = true
				}
			}
		}
	}
	if len(missing) == 0 {
		if len(uncounted) == 0 {
			return nil
		}
		addresses := make([]string, 0, len(uncounted))
		for address := range uncounted {
			addresses = append(addresses, address)
		}
		sort.Strings(addresses)
		return fmt.Errorf("uncounted Terraform references: %s", strings.Join(addresses, ", "))
	}
	addresses := make([]string, 0, len(missing))
	for address := range missing {
		addresses = append(addresses, address)
	}
	sort.Strings(addresses)
	return fmt.Errorf("undeclared Terraform objects: %s", strings.Join(addresses, ", "))
}

func TestDataSourceTestConfigurationSelfReferences(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_counted_dependency_cannot_be_its_own_parent",
			Input: `resource "aci_epg_useg_sub_block_statement" "parent_1" {
  count = 1
  parent_dn = aci_epg_useg_sub_block_statement.parent_1[0].id
}`,
			Expected: "self-referential Terraform objects: aci_epg_useg_sub_block_statement.parent_1",
		},
		{
			Name: "test_parent_and_child_of_same_type_are_distinct_objects",
			Input: `resource "aci_epg_useg_sub_block_statement" "parent_1" {
  count = 1
  parent_dn = "uni/tn-name_1/ap-name_1/epg-name_1/crtrn"
}
resource "aci_epg_useg_sub_block_statement" "test_1" {
  parent_dn = aci_epg_useg_sub_block_statement.parent_1[0].id
}`,
			Expected: "",
		},
		{
			Name: "test_removed_target_refers_to_prior_state",
			Input: `removed {
  from = aci_tenant.existing_on_create
  lifecycle {
    destroy = false
  }
}`,
			Expected: "",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			err := validateDataSourceTestConfiguration(testCase.Input.(string), nil)
			if expected := testCase.Expected.(string); expected != "" {
				assert.EqualError(t, err, expected)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func dataSourceTestTraversalHasIndex(traversal hcl.Traversal) bool {
	return slices.ContainsFunc(traversal, func(traverser hcl.Traverser) bool {
		_, ok := traverser.(hcl.TraverseIndex)
		return ok
	})
}

func dataSourceTestTraversalAddress(traversal hcl.Traversal) string {
	if len(traversal) < 2 {
		return ""
	}
	root, ok := traversal[0].(hcl.TraverseRoot)
	if !ok {
		return ""
	}
	if root.Name == "data" {
		if len(traversal) < 3 {
			return ""
		}
		resourceType, typeOK := traversal[1].(hcl.TraverseAttr)
		resourceName, nameOK := traversal[2].(hcl.TraverseAttr)
		if !typeOK || !nameOK || !strings.HasPrefix(resourceType.Name, "aci_") {
			return ""
		}
		return "data." + resourceType.Name + "." + resourceName.Name
	}
	if !strings.HasPrefix(root.Name, "aci_") {
		return ""
	}
	resourceName, ok := traversal[1].(hcl.TraverseAttr)
	if !ok {
		return ""
	}
	return root.Name + "." + resourceName.Name
}

func dataSourceTestClass(className, resourceName string, artifacts ...ArtifactEnum) Class {
	return Class{
		Name:         testClassName(className),
		ResourceName: resourceName,
		Artifacts:    artifacts,
	}
}

func dataSourceTestProperty(propertyName, attributeName string, values *TestValues) *Property {
	return &Property{
		PropertyName:  propertyName,
		AttributeName: attributeName,
		ValueType:     String,
		TestValues:    values,
	}
}

func dataSourceTestValues(create, update, defaultValue string, includeDefault bool) *TestValues {
	return &TestValues{
		Create:   dataSourceTestScenario(create, create, StringValue, true),
		Update:   dataSourceTestScenario(update, update, StringValue, true),
		Default:  dataSourceTestScenario(defaultValue, defaultValue, StringValue, includeDefault),
		ForceNew: dataSourceTestScenario(create, create, StringValue, true),
	}
}

func dataSourceTestReferenceValues(create, update, forceNew string) *TestValues {
	return &TestValues{
		Create:   dataSourceTestScenario(create, create, ReferenceValue, true),
		Update:   dataSourceTestScenario(update, update, ReferenceValue, true),
		Default:  dataSourceTestScenario(create, create, ReferenceValue, true),
		ForceNew: dataSourceTestScenario(forceNew, forceNew, ReferenceValue, true),
	}
}

func dataSourceTestScenario(configValue, assertValue string, valueType ValueRenderTypeEnum, include bool) TestValueScenario {
	return newTestValueScenario([]TestValueEntry{{
		ConfigValue:   configValue,
		ConfigInclude: include,
		AssertValue:   assertValue,
		ValueType:     valueType,
	}})
}

type resourceScenarioExpected struct {
	Contains    []string
	NotContains []string
	Action      string
	IDChanges   bool
	PlanOnly    bool
}

func TestResourceTestCreateOnlyProperties(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	class, store := resourceScenarioFixture()
	class.Properties["descr"].CreateOnly = true
	child := store.Classes["tagAnnotation"]
	child.Properties["value"].CreateOnly = true
	store.Classes["tagAnnotation"] = child
	scenarios, err := class.ResourceTestScenarios(store)
	if !assert.NoError(t, err) {
		return
	}
	testCases := []test.TestCase{
		{Name: "Minimum"}, {Name: "Full"}, {Name: "Update"}, {Name: "OmitOptional"}, {Name: "Reset"}, {Name: "UpdateChildren"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			index := slices.IndexFunc(scenarios, func(s ResourceTestScenario) bool { return s.Name == testCase.Name })
			if assert.NotEqual(t, -1, index) {
				text := strings.Join(strings.Fields(scenarios[index].Configuration), " ")
				assert.Contains(t, text, `description = "description_1"`)
				assert.NotContains(t, text, `description = "description_2"`)
				if testCase.Name == "UpdateChildren" {
					assert.Contains(t, text, `value = "value_1"`)
				}
			}
		})
	}
}

func TestResourceTestScenarioPresence(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	type input struct {
		Include  []string
		Exclude  []string
		Step     string
		Present  bool
		Expected string
		Required bool
	}
	testCases := []test.TestCase{
		{Name: "test_include_minimum", Input: input{Include: []string{"minimum"}, Step: "Minimum", Present: true, Expected: "description_1"}},
		{Name: "test_include_omitted", Input: input{Include: []string{"omit_optional"}, Step: "OmitOptional", Present: true, Expected: "description_2"}},
		{Name: "test_include_required_reset", Input: input{Include: []string{"reset"}, Step: "Reset", Present: true, Required: true}},
		{Name: "test_exclude_update_retains_state", Input: input{Exclude: []string{"update"}, Step: "Update", Expected: "description_1"}},
		{Name: "test_repeat_exclusion", Input: input{Exclude: []string{"update"}, Step: "UpdateUnchanged", Expected: "description_1"}},
		{Name: "test_exclude_full_default", Input: input{Exclude: []string{"full"}, Step: "Full", Expected: ""}},
		{Name: "test_exclude_full_repeat", Input: input{Exclude: []string{"full"}, Step: "FullUnchanged", Expected: ""}},
		{Name: "test_exclude_reset_retains_state", Input: input{Exclude: []string{"reset"}, Step: "Reset", Expected: "description_2"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			value := testCase.Input.(input)
			class, store := resourceScenarioFixture()
			property := class.Properties["descr"]
			property.Required = value.Required
			property.propertyDefinition.TestConfig.IncludeIn = value.Include
			property.propertyDefinition.TestConfig.ExcludeFrom = value.Exclude
			scenarios, err := class.ResourceTestScenarios(store)
			if !assert.NoError(t, err) {
				return
			}
			index := slices.IndexFunc(scenarios, func(s ResourceTestScenario) bool { return s.Name == value.Step })
			if !assert.NotEqual(t, -1, index) {
				return
			}
			scenario := scenarios[index]
			configuration := strings.Join(strings.Fields(scenario.Configuration), " ")
			assert.Equal(t, value.Present, strings.Contains(configuration, "description ="))
			assert.True(t, slices.ContainsFunc(scenario.Assertions, func(a TestAssertion) bool {
				return a.AttributePath == "description" && a.ExpectedValue == value.Expected
			}))
		})
	}
}

func TestResourceTestExplicitSensitiveReset(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	type input struct {
		Include, Explicit bool
		Expected          string
	}
	testCases := []test.TestCase{
		{Name: "test_explicit_reset", Input: input{Include: true, Explicit: true, Expected: ""}},
		{Name: "test_no_inclusion", Input: input{Explicit: true, Expected: "password_2"}},
		{Name: "test_no_explicit_value", Input: input{Include: true, Expected: "password_2"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			value := testCase.Input.(input)
			class, store := resourceScenarioFixture()
			property := dataSourceTestProperty("password", "password", dataSourceTestValues("password_1", "password_2", "", false))
			property.Sensitive = true
			property.propertyDefinition.TestConfig.Default.Defined = value.Explicit
			if value.Include {
				property.propertyDefinition.TestConfig.IncludeIn = []string{"reset"}
			}
			class.PropertiesAll = append(class.PropertiesAll, "password")
			class.Properties["password"] = property
			scenarios, err := class.ResourceTestScenarios(store)
			if !assert.NoError(t, err) {
				return
			}
			index := slices.IndexFunc(scenarios, func(s ResourceTestScenario) bool { return s.Name == "Reset" })
			if assert.NotEqual(t, -1, index) {
				configuration := strings.Join(strings.Fields(scenarios[index].Configuration), " ")
				assert.Contains(t, configuration, `password = "`+value.Expected+`"`)
			}
		})
	}
}

func TestResourceTestFrozenAttributes(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "string", Input: StringAttribute, Expected: []string{`frozen = "example"`, `frozen = ""`}},
		{Name: "set", Input: SetAttribute, Expected: []string{`frozen = ["example"]`, `frozen = []`}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			class, store := resourceScenarioFixture()
			class.LegacyAttributes = []MigrationAttribute{{AttributeName: "frozen", Type: testCase.Input.(LegacyAttributeTypeEnum), Status: Frozen}}
			scenarios, err := class.ResourceTestScenarios(store)
			if !assert.NoError(t, err) {
				return
			}
			for _, name := range []string{"Legacy", "LegacyUnchanged", "CurrentAfterLegacy", "LegacyConflict"} {
				assert.Equal(t, -1, slices.IndexFunc(scenarios, func(s ResourceTestScenario) bool { return s.Name == name }), "frozen-only configuration must not generate %s", name)
			}
			assert.NotEqual(t, -1, slices.IndexFunc(scenarios, func(s ResourceTestScenario) bool { return s.Name == "FrozenUnchanged" }))
			for index, name := range []string{"Frozen", "FrozenCleared"} {
				scenarioIndex := slices.IndexFunc(scenarios, func(s ResourceTestScenario) bool { return s.Name == name })
				if !assert.NotEqual(t, -1, scenarioIndex) {
					continue
				}
				scenario := scenarios[scenarioIndex]
				assert.Contains(t, strings.Join(strings.Fields(scenario.Configuration), " "), testCase.Expected.([]string)[index])
				assertionIndex := slices.IndexFunc(scenario.Assertions, func(a TestAssertion) bool { return a.AttributePath == "frozen" })
				if !assert.NotEqual(t, -1, assertionIndex) {
					continue
				}
				assertion := scenario.Assertions[assertionIndex]
				if testCase.Input == StringAttribute {
					assert.Equal(t, testAssertionAttributeValue, assertion.assertionType)
					assert.Equal(t, []string{"example", ""}[index], assertion.ExpectedValue)
				} else {
					assert.Equal(t, testAssertionSetCount, assertion.assertionType)
					assert.Equal(t, []int{1, 0}[index], assertion.ExpectedCount)
				}
			}
		})
	}
}

func TestResourceTestScenarios(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	class, store := resourceScenarioFixture()
	original := class.resourceTestCopy()
	scenarios, err := class.ResourceTestScenarios(store)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, original, class.resourceTestCopy(), "scenario projection must not mutate canonical test values")
	testCases := []test.TestCase{
		{Name: "Minimum", Expected: resourceScenarioExpected{Contains: []string{`name = "name_1"`, `parent_dn = aci_parent.test_1[0].id`}, NotContains: []string{"description", "labels", "annotations", "monitoring"}, Action: "Create"}},
		{Name: "Full", Expected: resourceScenarioExpected{Contains: []string{`description = "description_1"`, `labels`, `annotations`, `resource_only = "resource_only_1"`}}},
		{Name: "FullUnchanged", Expected: resourceScenarioExpected{PlanOnly: true}},
		{Name: "Update", Expected: resourceScenarioExpected{Contains: []string{`name = "name_1"`, `description = "description_2"`}, NotContains: []string{`name = "name_2"`, "aci_parent.test_2[0].id"}, Action: "Update"}},
		{Name: "OmitOptional", Expected: resourceScenarioExpected{Contains: []string{`name = "name_1"`}, NotContains: []string{"description", "labels", "annotation ="}}},
		{Name: "Reset", Expected: resourceScenarioExpected{Contains: []string{`description = ""`, `labels`, `[]`}, Action: "Update"}},
		{Name: "ClearChildren", Expected: resourceScenarioExpected{Contains: []string{`annotations = []`, `monitoring`}, Action: "Update"}},
		{Name: "UpdateChildren", Expected: resourceScenarioExpected{Contains: []string{`key   = "key_1"`, `value = "value_2"`}, Action: "Update"}},
		{Name: "RemoveChild", Expected: resourceScenarioExpected{NotContains: []string{`key   = "key_2"`}, Action: "Update"}},
		{Name: "ChangeChildIdentifier", Expected: resourceScenarioExpected{Contains: []string{`key   = "key_2"`}, NotContains: []string{`key   = "key_1"`}, Action: "Update"}},
		{Name: "OmitChildren", Expected: resourceScenarioExpected{NotContains: []string{"annotations", "monitoring"}}},
		{Name: "ReorderChildren", Expected: resourceScenarioExpected{PlanOnly: true}},
		{Name: "ExternalDeletion", Expected: resourceScenarioExpected{Action: "Create", PlanOnly: true}},
		{Name: "ReplaceIdentifier", Expected: resourceScenarioExpected{Contains: []string{`name = "name_2"`}, Action: "Replace", IDChanges: true}},
		{Name: "ReplaceParent", Expected: resourceScenarioExpected{Contains: []string{`aci_parent.test_2[0].id`}, Action: "Replace", IDChanges: true}},
		{Name: "Destroy", Expected: resourceScenarioExpected{Action: "Destroy"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			index := slices.IndexFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.Name == testCase.Name })
			if !assert.NotEqual(t, -1, index) {
				return
			}
			scenario := scenarios[index]
			expected := testCase.Expected.(resourceScenarioExpected)
			// Ignore alignment added by hclwrite when checking lifecycle values.
			configuration := strings.Join(strings.Fields(scenario.Configuration), " ")
			for _, value := range expected.Contains {
				assert.Contains(t, configuration, strings.Join(strings.Fields(value), " "))
			}
			for _, value := range expected.NotContains {
				assert.NotContains(t, configuration, strings.Join(strings.Fields(value), " "))
			}
			assert.Equal(t, expected.Action, scenario.Action)
			assert.Equal(t, expected.IDChanges, scenario.IDChanges)
			assert.Equal(t, expected.PlanOnly, scenario.PlanOnly)
		})
	}
	assert.True(t, slices.ContainsFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.Import }))
	importIndex := slices.IndexFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.Import })
	resetIndex := slices.IndexFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.Name == "Reset" })
	childrenIndex := slices.IndexFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.Name == "ClearChildren" })
	assert.Greater(t, importIndex, resetIndex, "import compares state persisted by the update/reset applies")
	assert.Less(t, importIndex, childrenIndex, "import precedes child lifecycle changes, matching the legacy ordering")
	assert.True(t, slices.ContainsFunc(scenarios, func(scenario ResourceTestScenario) bool {
		return scenario.ExpectError == "ExampleMonitor object defined by monitoring cannot be deleted"
	}))
	for _, name := range []string{"OmitOptional", "OmitChildren", "ClearChildren", "ChangeChildIdentifier"} {
		index := slices.IndexFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.Name == name })
		if !assert.NotEqual(t, -1, index) {
			continue
		}
		assertions := scenarios[index].Assertions
		if name == "OmitOptional" {
			assert.True(t, slices.ContainsFunc(assertions, func(a TestAssertion) bool {
				return a.AttributePath == "description" && a.ExpectedValue == "description_2"
			}))
			assert.True(t, slices.ContainsFunc(assertions, func(a TestAssertion) bool {
				return a.AttributePath == "annotation" && a.ExpectedExpression == "globalAnnotation"
			}))
		}
		if name == "OmitChildren" || name == "ClearChildren" {
			count := 2
			if name == "ClearChildren" {
				count = 0
			}
			assert.True(t, slices.ContainsFunc(assertions, func(a TestAssertion) bool {
				return a.AttributePath == "annotations" && a.IsSetCount() && a.ExpectedCount == count
			}))
		}
	}
}

func TestResourceTestExistingObjects(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	class, store := resourceScenarioFixture()
	enabled := true
	class.Properties["name"].PrefixEnabled = &enabled
	class.Properties["descr"].CreateOnly = true
	class.Properties["resourceOnly"].propertyDefinition.TestConfig.IncludeIn = []string{"minimum"}
	class.Properties["resourceOnly"].SupportedVersions, _ = NewVersions("6.0(1a)-")
	monitor := store.Classes["exampleMonitor"]
	monitor.RequiredAsChild = true
	store.Classes[monitor.Name.String()] = monitor
	class.TestChildren[1].Instances[0].Children = []*TestChild{class.TestChildren[0].clone()}
	original := class.resourceTestCopy()
	scenarios, err := class.ResourceTestScenarios(store)
	if !assert.NoError(t, err) || !assert.GreaterOrEqual(t, len(scenarios), 5) {
		return
	}
	assert.Equal(t, original, class.resourceTestCopy())
	assert.Equal(t, []string{"Minimum", "RejectExisting", "AllowExisting", "ForgetExisting", "Full"}, []string{scenarios[0].Name, scenarios[1].Name, scenarios[2].Name, scenarios[3].Name, scenarios[4].Name})
	minimum, diagnostics := hclsyntax.ParseConfig([]byte(scenarios[0].Configuration), "minimum.tf", hcl.InitialPos)
	if !assert.False(t, diagnostics.HasErrors(), diagnostics.Error()) {
		return
	}
	minimumBody := minimum.Body.(*hclsyntax.Body).Blocks[0].Body.SrcRange.SliceBytes([]byte(scenarios[0].Configuration))
	assert.Contains(t, string(minimumBody), "var.test_prefix")
	assert.Contains(t, string(minimumBody), "parent_dn")
	assert.Contains(t, string(minimumBody), "description_1")
	assert.Contains(t, string(minimumBody), "monitoring = {")
	assert.Contains(t, string(minimumBody), "6.0(1a)-")
	assert.NotContains(t, string(minimumBody), "labels")

	type expected struct {
		AllowExisting bool
		Forget        bool
		Error         string
		Action        string
	}
	testCases := []test.TestCase{
		{Name: "RejectExisting", Expected: expected{Error: "Object Already Exists"}},
		{Name: "AllowExisting", Expected: expected{AllowExisting: true, Action: "Create"}},
		{Name: "ForgetExisting", Expected: expected{AllowExisting: true, Forget: true}},
	}
	for index, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			scenario := scenarios[index+1]
			want := testCase.Expected.(expected)
			assert.Equal(t, want.Error, scenario.ExpectError)
			assert.Equal(t, want.Action, scenario.Action)
			assert.Equal(t, !want.Forget, scenario.ExistingObject)
			assert.Equal(t, want.Forget, scenario.ForgetExisting)
			assert.Empty(t, scenario.Assertions, "these steps use small identity/state checks")
			file, diagnostics := hclsyntax.ParseConfig([]byte(scenario.Configuration), "existing.tf", hcl.InitialPos)
			if !assert.False(t, diagnostics.HasErrors(), diagnostics.Error()) {
				return
			}
			resources, providers, removed := 0, 0, 0
			for _, block := range file.Body.(*hclsyntax.Body).Blocks {
				switch block.Type {
				case "resource":
					resources++
					assert.Equal(t, "aci_example", block.Labels[0])
					assert.Contains(t, []string{"test_1", "existing_on_create"}, block.Labels[1])
					assert.Equal(t, string(minimumBody), string(block.Body.SrcRange.SliceBytes([]byte(scenario.Configuration))))
				case "provider":
					providers++
					assert.Equal(t, []string{"aci"}, block.Labels)
					value, diagnostics := block.Body.Attributes["allow_existing_on_create"].Expr.Value(nil)
					assert.False(t, diagnostics.HasErrors(), diagnostics.Error())
					assert.True(t, value.RawEquals(cty.BoolVal(want.AllowExisting)))
				case "removed":
					removed++
					traversal, diagnostics := hcl.AbsTraversalForExpr(block.Body.Attributes["from"].Expr)
					assert.False(t, diagnostics.HasErrors(), diagnostics.Error())
					assert.Equal(t, "aci_example.existing_on_create", dataSourceTestTraversalAddress(traversal))
					if assert.Len(t, block.Body.Blocks, 1) {
						assert.Equal(t, "lifecycle", block.Body.Blocks[0].Type)
						value, diagnostics := block.Body.Blocks[0].Body.Attributes["destroy"].Expr.Value(nil)
						assert.False(t, diagnostics.HasErrors(), diagnostics.Error())
						assert.True(t, value.RawEquals(cty.False))
					}
				}
			}
			assert.Equal(t, 1, providers)
			if want.Forget {
				assert.Equal(t, 1, resources)
				assert.Equal(t, 1, removed)
			} else {
				assert.Equal(t, 2, resources)
				assert.Zero(t, removed)
			}
		})
	}
}

func TestResourceTestMinimumRequiredChildren(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	class, store := resourceScenarioFixture()
	monitor := store.Classes["exampleMonitor"]
	monitor.RequiredAsChild = true
	store.Classes[monitor.Name.String()] = monitor
	detail := store.Classes["tagAnnotation"]
	detail.Name = testClassName("requiredDetail")
	detail.ResourceNameNested = "required_details"
	detail.RequiredAsChild = true
	store.Classes[detail.Name.String()] = detail
	required := class.TestChildren[0].clone()
	required.Class = detail.Name
	class.TestChildren[1].Instances[0].Children = []*TestChild{required, class.TestChildren[0].clone()}
	original := class.resourceTestCopy()

	scenarios, err := class.ResourceTestScenarios(store)
	if !assert.NoError(t, err) {
		return
	}
	minimum := scenarios[0]
	assert.Equal(t, "Minimum", minimum.Name)
	assert.Contains(t, minimum.Configuration, "monitoring = {")
	assert.Contains(t, minimum.Configuration, "required_details = [")
	assert.NotContains(t, minimum.Configuration, "annotations")
	assert.True(t, slices.ContainsFunc(minimum.Assertions, func(assertion TestAssertion) bool {
		return assertion.AttributePath == "monitoring.required_details" && assertion.IsSetCount() && assertion.ExpectedCount == 2
	}))
	assert.Equal(t, original, class.resourceTestCopy())
}

func TestResourceTestRequiredChildScenarios(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	testCases := []test.TestCase{
		{Name: "test_non_deletable_required_child", Input: false},
		{Name: "test_deletable_required_child", Input: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			class, store := resourceScenarioFixture()
			monitor := store.Classes["exampleMonitor"]
			monitor.RequiredAsChild = true
			monitor.AllowDelete = testCase.Input.(bool)
			store.Classes[monitor.Name.String()] = monitor
			scenarios, err := class.ResourceTestScenarios(store)
			if !assert.NoError(t, err) {
				return
			}
			for _, name := range []string{"ClearChildren", "OmitChildren"} {
				index := slices.IndexFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.Name == name })
				if assert.NotEqual(t, -1, index) {
					assert.Contains(t, scenarios[index].Configuration, "monitoring = {")
				}
			}
			index := slices.IndexFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.Name == "RejectDeleteExampleMonitor" })
			if assert.NotEqual(t, -1, index) {
				assert.NotContains(t, scenarios[index].Configuration, "monitoring =")
				assert.Equal(t, `The argument "monitoring" is required`, scenarios[index].ExpectError)
				assert.Contains(t, scenarios[index+1].Configuration, "monitoring = {")
			}
		})
	}
}

func TestResourceTestChildIdentifierPreservesDescendants(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	testCases := []test.TestCase{
		{Name: "test_children", Input: 1},
		{Name: "test_deep_children", Input: 3},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			class, store := resourceScenarioFixture()
			child := class.TestChildren[0].clone()
			for _, name := range []string{"exampleLayerOne", "exampleLayerTwo", "exampleLayerThree"}[:testCase.Input.(int)] {
				parent := store.Classes["tagAnnotation"]
				parent.Name = testClassName(name)
				parent.ResourceNameNested = name
				store.Classes[name] = parent
				instances := []TestChildInstance{buildChildInstance(&parent, 0), buildChildInstance(&parent, 1)}
				for index := range instances {
					instances[index].Children = []*TestChild{child.clone()}
				}
				child = &TestChild{Class: parent.Name, Instances: instances}
			}
			class.TestChildren = []*TestChild{child}
			original := class.resourceTestCopy()
			changed := class.resourceTestCopy()
			builder := resourceTestBuilder{class: &class, dataStore: store}
			builder.changeChildren(&changed, "trim")
			builder.changeChildren(&changed, "identifier")
			if assert.Len(t, changed.TestChildren[0].Instances, 1) {
				assert.Equal(t, child.Instances[1], changed.TestChildren[0].Instances[0])
			}
			assert.Equal(t, original, class.resourceTestCopy())
		})
	}
}

func TestResourceTestApplicability(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	testCases := []test.TestCase{
		{Name: "test_non_deletable_root", Input: false, Expected: false},
		{Name: "test_deletable_root", Input: true, Expected: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			class := dataSourceTestClass("constantRoot", "constant", ResourceArtifact)
			class.AllowDelete = testCase.Input.(bool)
			scenarios, err := class.ResourceTestScenarios(&DataStore{})
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, testCase.Expected, slices.ContainsFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.ExternalDelete }))
			assert.False(t, slices.ContainsFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.IDChanges || scenario.Action == "Update" }))
		})
	}
}

func TestResourceTestChangedVersionsAndSingletonClear(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	class, store := resourceScenarioFixture()
	monitor := store.Classes["exampleMonitor"]
	monitor.AllowDelete = true
	monitor.SupportedVersions, _ = NewVersions("5.0(1a)-")
	store.Classes["exampleMonitor"] = monitor
	for name, property := range class.Properties {
		if name != "name" && name != "parentDn" {
			property.SupportedVersions, _ = NewVersions("6.0(1a)-")
		}
	}
	scenarios, err := class.ResourceTestScenarios(store)
	if !assert.NoError(t, err) {
		return
	}
	testCases := []test.TestCase{
		{Name: "Update", Expected: [][]string{{"6.0(1a)-"}}},
		{Name: "Reset", Expected: [][]string{{"6.0(1a)-"}}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			index := slices.IndexFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.Name == testCase.Name })
			if assert.NotEqual(t, -1, index) {
				assert.Equal(t, testCase.Expected, scenarios[index].ChangeVersions)
			}
		})
	}
	index := slices.IndexFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.Name == "ClearChildren" })
	if assert.NotEqual(t, -1, index) {
		assert.Contains(t, scenarios[index].Configuration, `? {} : null`)
		assert.True(t, slices.ContainsFunc(scenarios[index].Assertions, func(a TestAssertion) bool { return a.AttributePath == "monitoring" && a.IsAttributeAbsent() }))
	}
}

func TestResourceTestDependencyUnion(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	class, store := resourceScenarioFixture()
	configuration, err := class.ResourceTestDependenciesHCL(store)
	assert.NoError(t, err)
	assert.Contains(t, configuration, `resource "aci_parent" "test_1"`)
	assert.Contains(t, configuration, `resource "aci_parent" "test_2"`)
	class.Properties["name"].PrefixEnabled = nil
	assertions := class.resourceTestAssertions(store)
	assert.False(t, slices.ContainsFunc(assertions, func(a TestAssertion) bool { return a.PairAddress == "aci_example.test_1" }), "resource expectations must not compare the resource to itself")
	assert.True(t, slices.ContainsFunc(assertions, func(a TestAssertion) bool {
		return a.AttributePath == "resource_only" && a.ExpectedValue == "resource_only_1"
	}))
}

func resourceScenarioFixture() (Class, *DataStore) {
	disabled := false
	class := dataSourceTestClass("exampleRoot", "example", ResourceArtifact)
	class.AllowDelete = true
	class.IdentifiedBy = []string{"name"}
	class.PropertiesAll = []string{"annotation", "descr", "labels", "name", "parentDn", "resourceOnly"}
	class.Properties = map[string]*Property{
		"annotation":   dataSourceTestProperty("annotation", "annotation", dataSourceTestValues("annotation_1", "annotation_2", "", false)),
		"descr":        dataSourceTestProperty("descr", "description", dataSourceTestValues("description_1", "description_2", "", false)),
		"labels":       dataSourceTestProperty("labels", "labels", dataSourceTestValues("label_1", "label_2", "", false)),
		"name":         dataSourceTestProperty("name", "name", dataSourceTestValues("name_1", "name_2", "name_1", true)),
		"parentDn":     dataSourceTestProperty("parentDn", "parent_dn", dataSourceTestReferenceValues("aci_parent.test_1.id", "aci_parent.test_1.id", "aci_parent.test_2.id")),
		"resourceOnly": dataSourceTestProperty("resourceOnly", "resource_only", dataSourceTestValues("resource_only_1", "resource_only_2", "", false)),
	}
	for _, property := range class.Properties {
		property.Computed = true
	}
	class.Properties["name"].Required = true
	class.Properties["name"].RequiresReplace = true
	class.Properties["name"].PrefixEnabled = &disabled
	class.Properties["parentDn"].Required = true
	class.Properties["parentDn"].RequiresReplace = true
	class.Properties["labels"].ValueType = Set
	class.Properties["labels"].TestValues.Default = newTestValueScenario(nil)
	class.Properties["resourceOnly"].IgnoreInDataSourceTest = true
	class.Properties["descr"].propertyDefinition.TestConfig.Default.Defined = true
	class.Properties["labels"].propertyDefinition.TestConfig.Default.Defined = true
	parent := dataSourceTestClass("exampleParent", "parent", ResourceArtifact)
	parent.AllowDelete = true
	parent.IdentifiedBy = []string{"name"}
	parent.PropertiesAll = []string{"name"}
	parent.Properties = map[string]*Property{"name": dataSourceTestProperty("name", "name", dataSourceTestValues("parent_1", "parent_2", "parent_1", true))}
	for _, reference := range []string{"aci_parent.test_1.id", "aci_parent.test_2.id"} {
		class.TestDependencies = append(class.TestDependencies, &TestDependency{Class: parent.Name, Reference: reference, ReferenceType: ResourceReference, Role: Parent})
	}
	child := dataSourceTestClass("tagAnnotation", "annotation")
	child.AllowDelete = true
	child.ResourceNameNested = "annotations"
	child.IdentifiedBy = []string{"key"}
	child.PropertiesAll = []string{"key", "value"}
	child.Properties = map[string]*Property{
		"key":   dataSourceTestProperty("key", "key", dataSourceTestValues("key_1", "key_2", "key_1", true)),
		"value": dataSourceTestProperty("value", "value", dataSourceTestValues("value_1", "value_2", "", false)),
	}
	monitor := dataSourceTestClass("exampleMonitor", "monitor")
	monitor.IsSingleNestedWhenDefinedAsChild = true
	monitor.ResourceNameNested = "monitoring"
	monitor.PropertiesAll = []string{"name"}
	monitor.Properties = map[string]*Property{"name": dataSourceTestProperty("name", "name", dataSourceTestValues("monitor_1", "monitor_2", "monitor_1", true))}
	class.TestChildren = []*TestChild{{Class: child.Name, Instances: []TestChildInstance{buildChildInstance(&child, 0), buildChildInstance(&child, 1)}}, {Class: monitor.Name, Instances: []TestChildInstance{buildChildInstance(&monitor, 0)}}}
	store := &DataStore{Classes: map[string]Class{class.Name.String(): class, parent.Name.String(): parent, child.Name.String(): child, monitor.Name.String(): monitor}}
	return class, store
}

func TestAllResourceTestConfigurationsAreSelfContained(t *testing.T) {
	test.InitializeTest(t)
	_, file, _, _ := runtime.Caller(0)
	t.Chdir(filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..")))
	ctx := NewContext()
	store, err := NewDataStore(ctx)
	if !assert.NoError(t, err) || !assert.NoError(t, ctx.Diagnostics.Error()) {
		return
	}
	for _, name := range slices.Sorted(maps.Keys(store.Classes)) {
		class := store.Classes[name]
		if !class.HasResourceArtifact() || slices.Contains(class.TestConfig.IgnoreTests, ResourceIgnoreTest) {
			continue
		}
		t.Run(name, func(t *testing.T) {
			dependencies, err := class.ResourceTestDependenciesHCL(store)
			if !assert.NoError(t, err) {
				return
			}
			assert.NoError(t, validateDataSourceTestDependencyCounts(dependencies))
			scenarios, err := class.ResourceTestScenarios(store)
			if !assert.NoError(t, err) {
				return
			}
			for _, scenario := range scenarios {
				t.Run(scenario.Name, func(t *testing.T) {
					assert.NoError(t, validateDataSourceTestConfiguration(dependencies+"\n"+scenario.Configuration, scenario.Assertions))
				})
			}
		})
	}
}

func TestResourceTestLegacyScenariosRequireConfiguredAliases(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "configured_alias", Input: "name", Expected: true},
		{Name: "unused_alias", Input: "missing_attribute", Expected: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			class, store := resourceScenarioFixture()
			class.LegacyAttributes = []MigrationAttribute{
				{AttributeName: "old_name", ChildAttribute: testCase.Input.(string), Type: StringAttribute, Status: LegacyDeprecated},
				{AttributeName: "frozen", Type: StringAttribute, Status: Frozen},
			}
			scenarios, err := class.ResourceTestScenarios(store)
			if !assert.NoError(t, err) {
				return
			}
			for _, name := range []string{"Legacy", "LegacyUnchanged", "CurrentAfterLegacy", "LegacyConflict"} {
				index := slices.IndexFunc(scenarios, func(s ResourceTestScenario) bool { return s.Name == name })
				assert.Equal(t, testCase.Expected, index >= 0, "%s", name)
			}
			for _, name := range []string{"Frozen", "FrozenUnchanged", "FrozenCleared"} {
				assert.NotEqual(t, -1, slices.IndexFunc(scenarios, func(s ResourceTestScenario) bool { return s.Name == name }))
			}
		})
	}
}
