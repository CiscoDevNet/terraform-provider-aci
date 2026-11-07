package data

import (
	"fmt"
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
