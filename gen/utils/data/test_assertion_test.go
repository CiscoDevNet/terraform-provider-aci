package data

import (
	"fmt"
	"slices"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestDataSourceLegacyNamedRelationAssertions(t *testing.T) {
	test.InitializeTest(t)
	type input struct {
		Named        bool
		Type         LegacyAttributeTypeEnum
		Block        bool
		Unconfigured bool
		Ignored      bool
		Versions     string
	}
	testCases := []test.TestCase{
		{Name: "test_named_scalar", Input: input{Named: true, Type: StringAttribute}, Expected: testAssertionAttributeNotNull},
		{Name: "test_named_scalar_version_gate", Input: input{Named: true, Type: StringAttribute, Versions: "1.0(1e)-6.1(1e)"}, Expected: testAssertionAttributeNotNull},
		{Name: "test_explicit_target", Input: input{Type: StringAttribute}, Expected: testAssertionAttributePair},
		{Name: "test_named_collection", Input: input{Named: true, Type: SetAttribute}, Expected: testAssertionSetPair},
		{Name: "test_collection_version_gate", Input: input{Type: SetAttribute, Versions: "2.0(1a)-"}, Expected: testAssertionSetPair},
		{Name: "test_legacy_block", Input: input{Named: true, Type: SetNestedAttribute, Block: true}, Expected: testAssertionSetPair},
		{Name: "test_unconfigured_child", Input: input{Named: true, Type: StringAttribute, Unconfigured: true}, Expected: testAssertionAttributePair},
		{Name: "test_ignored_child", Input: input{Named: true, Type: StringAttribute, Ignored: true}, Expected: testAssertionAttributePair},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			value := testCase.Input.(input)
			class, store := resourceScenarioFixture()
			childName := class.TestChildren[1].Class
			class.LegacyAttributes = []MigrationAttribute{
				{AttributeName: "legacy_monitor", ChildName: childName, Named: value.Named, Type: value.Type, Block: value.Block, Status: LegacyDeprecated},
				{AttributeName: "legacy_name", ChildAttribute: "name", Type: StringAttribute, Status: LegacyDeprecated},
			}
			if value.Unconfigured {
				class.TestChildren = class.TestChildren[:1]
			}
			child := store.Classes[childName.String()]
			if value.Ignored {
				child.TestConfig.IgnoreTests = []IgnoreTestEnum{ChildIgnoreTest}
			}
			if value.Versions != "" {
				versions, err := NewVersions(value.Versions)
				if !assert.NoError(t, err) {
					return
				}
				child.SupportedVersions = versions
			}
			store.Classes[childName.String()] = child
			assertions := class.DataSourceTestAssertions(store)
			index := slices.IndexFunc(assertions, func(a TestAssertion) bool { return a.AttributePath == "legacy_monitor" })
			if !assert.NotEqual(t, -1, index) {
				return
			}
			assert.Equal(t, testCase.Expected, assertions[index].assertionType)
			if assertions[index].IsSetPair() {
				assert.Equal(t, "legacy_monitor", assertions[index].PairAttribute)
				assert.Equal(t, "aci_example.test_1", assertions[index].PairAddress)
			}
			if value.Versions != "" {
				assert.Equal(t, []string{value.Versions}, assertions[index].Versions)
			}
			aliasIndex := slices.IndexFunc(assertions, func(a TestAssertion) bool { return a.AttributePath == "legacy_name" })
			if assert.NotEqual(t, -1, aliasIndex) {
				assert.True(t, assertions[aliasIndex].IsAttributePair())
			}
		})
	}
}

func TestDataSourceFrozenAttributes(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "string", Input: StringAttribute, Expected: testAssertionAttributeValue},
		{Name: "set", Input: SetAttribute, Expected: testAssertionSetCount},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			test.InitializeTest(t)
			class, store := resourceScenarioFixture()
			class.LegacyAttributes = []MigrationAttribute{{AttributeName: "frozen", Type: testCase.Input.(LegacyAttributeTypeEnum), Status: Frozen}}
			assertions := class.DataSourceTestAssertions(store)
			index := slices.IndexFunc(assertions, func(a TestAssertion) bool { return a.AttributePath == "frozen" })
			if assert.NotEqual(t, -1, index) {
				assert.Equal(t, testCase.Expected, assertions[index].assertionType)
				assert.Empty(t, assertions[index].ExpectedValue)
				assert.Zero(t, assertions[index].ExpectedCount)
			}
		})
	}
}

func TestTestValueAssertionStrategy(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	versions := []string{"3.0(1k)-"}
	testCases := []test.TestCase{
		{Name: "test_literal_value", Input: TestValueEntry{AssertValue: "0"}, Expected: testAssertionAttributeValue},
		{Name: "test_literal_presence", Input: TestValueEntry{AssertValue: "0", AssertStrategy: AssertPresenceOnly}, Expected: testAssertionAttributeSet},
		{Name: "test_reference_value", Input: TestValueEntry{AssertValue: "aci_tenant.test_1.id", ValueType: ReferenceValue}, Expected: testAssertionAttributePair},
		{Name: "test_reference_presence", Input: TestValueEntry{AssertValue: "aci_tenant.test_1.id", ValueType: ReferenceValue, AssertStrategy: AssertPresenceOnly}, Expected: testAssertionAttributeSet},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			actual := testAttributeAssertion("site_id", testCase.Input.(TestValueEntry), versions)
			assert.Equal(t, testCase.Expected.(testAssertionType), actual.assertionType)
			assert.Equal(t, versions, actual.Versions)
		})
	}
}

func TestNestedPresenceAssertionScope(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	class, store := resourceScenarioFixture()
	child := &class.TestChildren[0].Instances[0]
	scenario := child.Properties["value"]
	scenario.Entries[0].AssertStrategy = AssertPresenceOnly
	child.Properties["value"] = scenario
	assertions := class.resourceTestAssertions(store)
	index := slices.IndexFunc(assertions, func(a TestAssertion) bool { return a.IsAttributeSet() && a.AttributePath == "annotations.*.value" })
	if assert.NotEqual(t, -1, index) {
		assert.Equal(t, []TestAssertionParent{{AttributePath: "annotations.*", Values: []TestAssertionValue{{AttributePath: "key", ExpectedValue: "key_1"}}}}, assertions[index].Parents)
	}
}

func TestDependencyVersionsCoveredByAssertion(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	dependencyVersions, err := NewVersions("3.0(1a)-")
	if err != nil {
		t.Fatal(err)
	}
	testCases := []test.TestCase{
		{Name: "test_no_assertion_gate", Input: []string(nil), Expected: false},
		{Name: "test_equal_gate", Input: []string{"3.0(1a)-"}, Expected: true},
		{Name: "test_narrower_gate", Input: []string{"4.0(1a)-"}, Expected: true},
		{Name: "test_broader_gate", Input: []string{"2.0(1a)-"}, Expected: false},
		{Name: "test_one_of_multiple_gates_is_narrower", Input: []string{"2.0(1a)-", "4.0(1a)-"}, Expected: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			actual := dependencyVersionsCoveredByAssertion(dependencyVersions, testCase.Input.([]string))
			assert.Equal(t, testCase.Expected.(bool), actual, test.MessageEqual(testCase.Expected, actual, testCase.Name))
		})
	}
}

func TestRepeatedGrandchildAssertionParents(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	childVersions, err := NewVersions("2.0(1a)-")
	if !assert.NoError(t, err) {
		return
	}
	grandchildVersions, err := NewVersions("4.0(1a)-")
	if !assert.NoError(t, err) {
		return
	}
	sourceVersions, err := NewVersions("5.0(1a)-")
	if !assert.NoError(t, err) {
		return
	}
	child := dataSourceTestClass("exampleChild", "example_child")
	child.ResourceNameNested = "children"
	child.IdentifiedBy = []string{"target"}
	child.SupportedVersions = childVersions
	child.PropertiesAll = []string{"target"}
	child.Properties = map[string]*Property{"target": dataSourceTestProperty("target", "target", nil)}
	grandchild := dataSourceTestClass("tagAnnotation", "annotation")
	grandchild.ResourceNameNested = "annotations"
	grandchild.IdentifiedBy = []string{"key"}
	grandchild.SupportedVersions = grandchildVersions
	grandchild.PropertiesAll = []string{"key", "value"}
	grandchild.Properties = map[string]*Property{"key": dataSourceTestProperty("key", "key", nil), "value": dataSourceTestProperty("value", "value", nil)}
	greatGrandchild := dataSourceTestClass("exampleRule", "example_rule")
	greatGrandchild.ResourceNameNested = "rules"
	greatGrandchild.IdentifiedBy = []string{"code"}
	greatGrandchild.SupportedVersions = grandchildVersions
	greatGrandchild.PropertiesAll = []string{"code"}
	greatGrandchild.Properties = map[string]*Property{"code": dataSourceTestProperty("code", "code", nil)}
	source := dataSourceTestClass("vzBrCP", "contract", ResourceArtifact)
	source.SupportedVersions = sourceVersions
	root := dataSourceTestClass("exampleRoot", "example_root", DatasourceArtifact)
	root.TestChildren = []*TestChild{{Class: child.Name}}
	for index := 1; index <= 2; index++ {
		reference := fmt.Sprintf("aci_contract.test_%d.name", index)
		root.TestDependencies = append(root.TestDependencies, &TestDependency{Class: source.Name, Reference: reference, ReferenceType: ResourceReference, Role: Target})
		root.TestChildren[0].Instances = append(root.TestChildren[0].Instances, TestChildInstance{
			Properties: map[string]TestValueScenario{"target": dataSourceTestScenario(reference, reference, ReferenceValue, true)},
			Children: []*TestChild{{Class: grandchild.Name, Instances: []TestChildInstance{{
				Properties: map[string]TestValueScenario{
					"key":   dataSourceTestScenario("key_1", "key_1", StringValue, true),
					"value": dataSourceTestScenario("value_1", "value_1", StringValue, true),
				},
				Children: []*TestChild{{Class: greatGrandchild.Name, Instances: []TestChildInstance{{Properties: map[string]TestValueScenario{
					"code": dataSourceTestScenario("code_1", "code_1", StringValue, true),
				}}}}},
			}}}},
		})
	}
	store := &DataStore{Classes: map[string]Class{child.Name.String(): child, grandchild.Name.String(): grandchild, greatGrandchild.Name.String(): greatGrandchild, source.Name.String(): source}}
	assertions := root.DataSourceTestAssertions(store)
	assert.Equal(t, assertions, root.DataSourceTestAssertions(store), "rendering parent references must not mutate shared assertion data")
	testCases := []test.TestCase{
		{Name: "test_first_parent", Input: 1, Expected: "aci_contract.test_1.0"},
		{Name: "test_second_parent", Input: 2, Expected: "aci_contract.test_2.0"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			identity := TestAssertionValue{AttributePath: "target", PairAddress: testCase.Expected.(string), PairAttribute: "name"}
			countIndex := slices.IndexFunc(assertions, func(a TestAssertion) bool {
				return a.AttributePath == "children.*" && slices.Contains(a.NestedValues, identity) && slices.Contains(a.NestedValues, TestAssertionValue{AttributePath: "annotations.#", ExpectedValue: "1"})
			})
			if !assert.NotEqual(t, -1, countIndex, "a count is required for each parent's annotations") {
				return
			}
			assert.Equal(t, []string{"2.0(1a)-", "4.0(1a)-", "5.0(1a)-"}, assertions[countIndex].Versions)
			grandchildIndex := slices.IndexFunc(assertions, func(a TestAssertion) bool {
				return a.AttributePath == "children.*.annotations.*" && len(a.Parents) == 1 && slices.Contains(a.Parents[0].Values, identity) && slices.Contains(a.NestedValues, TestAssertionValue{AttributePath: "value", ExpectedValue: "value_1"})
			})
			if !assert.NotEqual(t, -1, grandchildIndex) {
				return
			}
			assert.Equal(t, "children.*", assertions[grandchildIndex].Parents[0].AttributePath)
			assert.Equal(t, []string{"2.0(1a)-", "4.0(1a)-", "5.0(1a)-"}, assertions[grandchildIndex].Versions)
			greatGrandchildIndex := slices.IndexFunc(assertions, func(a TestAssertion) bool {
				return a.AttributePath == "children.*.annotations.*.rules.*" && len(a.Parents) == 2 && slices.Contains(a.Parents[0].Values, identity)
			})
			if !assert.NotEqual(t, -1, greatGrandchildIndex) {
				return
			}
			assert.Equal(t, TestAssertionParent{AttributePath: "children.*.annotations.*", Values: []TestAssertionValue{{AttributePath: "key", ExpectedValue: "key_1"}}}, assertions[greatGrandchildIndex].Parents[1])
		})
	}
}

func TestDataSourceTestAssertions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	classVersions, err := NewVersions("1.0(1e)-")
	if err != nil {
		t.Fatal(err)
	}
	propertyVersions, err := NewVersions("2.0(1a)-")
	if err != nil {
		t.Fatal(err)
	}

	class := dataSourceTestClass("exampleAssertions", "example_assertions", DatasourceArtifact)
	class.SupportedVersions = classVersions
	class.Properties = map[string]*Property{
		"literal":   dataSourceTestProperty("literal", "literal", &TestValues{Create: dataSourceTestScenario("input", "normalized", StringValue, true)}),
		"reference": dataSourceTestProperty("reference", "reference", &TestValues{Create: dataSourceTestScenario("aci_tenant.test.id", "aci_tenant.test.id", ReferenceValue, true)}),
		"members": {
			PropertyName:  "members",
			AttributeName: "members",
			ValueType:     Set,
			TestValues: &TestValues{Create: newTestValueScenario([]TestValueEntry{
				{ConfigValue: "first", ConfigInclude: true, AssertValue: "first", ValueType: StringValue},
				{ConfigValue: "aci_tenant.test.id", ConfigInclude: true, AssertValue: "aci_tenant.test.id", ValueType: ReferenceValue},
			})},
		},
		"computed": {
			PropertyName:  "computed",
			AttributeName: "computed",
			ValueType:     String,
			ReadOnly:      true,
		},
		"versioned": {
			PropertyName:      "versioned",
			AttributeName:     "versioned",
			ValueType:         String,
			SupportedVersions: propertyVersions,
			TestValues:        &TestValues{Create: dataSourceTestScenario("value", "value", StringValue, true)},
		},
		"sensitive": {
			PropertyName:  "sensitive",
			AttributeName: "sensitive",
			ValueType:     String,
			Sensitive:     true,
			TestValues:    &TestValues{Create: dataSourceTestScenario("secret", "secret", StringValue, true)},
		},
		"datasourceIgnored": {
			PropertyName:           "datasourceIgnored",
			AttributeName:          "datasource_ignored",
			ValueType:              String,
			IgnoreInDataSourceTest: true,
			TestValues:             &TestValues{Create: dataSourceTestScenario("ignored", "ignored", StringValue, true)},
		},
	}
	class.PropertiesAll = []string{"literal", "reference", "members", "computed", "versioned", "sensitive", "datasourceIgnored"}
	tenant := dataSourceTestClass("fvTenant", "tenant", ResourceArtifact)
	class.TestDependencies = []*TestDependency{{
		Class:         tenant.Name,
		Reference:     "aci_tenant.test.id",
		ReferenceType: ResourceReference,
		Role:          Target,
	}}

	expected := []TestAssertion{
		{assertionType: testAssertionAttributeValue, AttributePath: "literal", ExpectedValue: "normalized"},
		{assertionType: testAssertionAttributePair, AttributePath: "reference", PairAddress: "aci_tenant.test.0", PairAttribute: "id"},
		{assertionType: testAssertionSetCount, AttributePath: "members", ExpectedCount: 2},
		{assertionType: testAssertionSetElementValue, AttributePath: "members.*", ExpectedValue: "first"},
		{assertionType: testAssertionSetElementPair, AttributePath: "members.*", PairAddress: "aci_tenant.test.0", PairAttribute: "id"},
		{assertionType: testAssertionAttributeSet, AttributePath: "computed"},
		{assertionType: testAssertionAttributeValue, AttributePath: "versioned", ExpectedValue: "value", Versions: []string{"2.0(1a)-"}},
	}
	store := &DataStore{Classes: map[string]Class{tenant.Name.String(): tenant}}
	assert.Equal(t, expected, class.DataSourceTestAssertions(store))
}

func TestDataSourceTestChildAssertions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	parentVersions, err := NewVersions("1.0(1e)-")
	if err != nil {
		t.Fatal(err)
	}
	childVersions, err := NewVersions("2.0(1a)-")
	if err != nil {
		t.Fatal(err)
	}
	referenceVersions, err := NewVersions("3.0(1a)-")
	if err != nil {
		t.Fatal(err)
	}

	grandchild := dataSourceTestClass("exampleGrandchild", "example_grandchild")
	grandchild.ResourceNameNested = "details"
	grandchild.IsSingleNestedWhenDefinedAsChild = true
	grandchild.SupportedVersions = childVersions
	grandchild.Properties = map[string]*Property{
		"value": dataSourceTestProperty("value", "value", nil),
	}
	grandchild.PropertiesAll = []string{"value"}

	child := dataSourceTestClass("exampleChild", "example_child")
	child.ResourceNameNested = "children"
	child.IdentifiedBy = []string{"key"}
	child.SupportedVersions = childVersions
	child.Properties = map[string]*Property{
		"key": dataSourceTestProperty("key", "key", nil),
		"target": {
			PropertyName:      "target",
			AttributeName:     "target",
			ValueType:         String,
			SupportedVersions: referenceVersions,
		},
		"labels": {
			PropertyName:  "labels",
			AttributeName: "labels",
			ValueType:     Set,
		},
		"secret": {
			PropertyName:  "secret",
			AttributeName: "secret",
			ValueType:     String,
			Sensitive:     true,
		},
	}
	child.PropertiesAll = []string{"key", "target", "labels", "secret"}

	childInstances := []TestChildInstance{
		{
			Properties: map[string]TestValueScenario{
				"key":    dataSourceTestScenario("key_1", "key_1", StringValue, true),
				"target": dataSourceTestScenario("aci_tenant.test_1.id", "aci_tenant.test_1.id", ReferenceValue, true),
				"labels": dataSourceTestScenario("label_1", "label_1", StringValue, true),
				"secret": dataSourceTestScenario("secret_1", "secret_1", StringValue, true),
			},
			Children: []*TestChild{{
				Class: grandchild.Name,
				Instances: []TestChildInstance{{Properties: map[string]TestValueScenario{
					"value": dataSourceTestScenario("detail_1", "detail_1", StringValue, true),
				}}},
			}},
		},
		{
			Properties: map[string]TestValueScenario{
				"key":    dataSourceTestScenario("key_2", "key_2", StringValue, true),
				"target": dataSourceTestScenario("aci_tenant.test_2.id", "aci_tenant.test_2.id", ReferenceValue, true),
				"labels": dataSourceTestScenario("label_2", "label_2", StringValue, true),
				"secret": dataSourceTestScenario("secret_2", "secret_2", StringValue, true),
			},
			Children: []*TestChild{{
				Class: grandchild.Name,
				Instances: []TestChildInstance{{Properties: map[string]TestValueScenario{
					"value": dataSourceTestScenario("detail_2", "detail_2", StringValue, true),
				}}},
			}},
		},
	}

	parent := dataSourceTestClass("exampleParent", "example_parent", DatasourceArtifact)
	parent.SupportedVersions = parentVersions
	parent.TestChildren = []*TestChild{{Class: child.Name, Instances: childInstances}}
	tenant := dataSourceTestClass("fvTenant", "tenant", ResourceArtifact)
	parent.TestDependencies = []*TestDependency{
		{Class: tenant.Name, Reference: "aci_tenant.test_1.id", ReferenceType: ResourceReference, Role: Target},
		{Class: tenant.Name, Reference: "aci_tenant.test_2.id", ReferenceType: ResourceReference, Role: Target},
	}
	store := &DataStore{Classes: map[string]Class{
		child.Name.String():      child,
		grandchild.Name.String(): grandchild,
		tenant.Name.String():     tenant,
	}}

	expected := []TestAssertion{
		{assertionType: testAssertionSetCount, AttributePath: "children", ExpectedCount: 2, Versions: []string{"2.0(1a)-"}},
		{assertionType: testAssertionNestedSetElement, AttributePath: "children.*", NestedValues: []TestAssertionValue{{AttributePath: "key", ExpectedValue: "key_1"}, {AttributePath: "labels.#", ExpectedValue: "1"}}, Versions: []string{"2.0(1a)-"}},
		{assertionType: testAssertionNestedSetElement, AttributePath: "children.*", NestedValues: []TestAssertionValue{{AttributePath: "key", ExpectedValue: "key_1"}, {AttributePath: "target", PairAddress: "aci_tenant.test_1.0", PairAttribute: "id"}}, Versions: []string{"2.0(1a)-", "3.0(1a)-"}},
		{assertionType: testAssertionSetElementValue, AttributePath: "children.*.labels.*", ExpectedValue: "label_1", Parents: []TestAssertionParent{{AttributePath: "children.*", Values: []TestAssertionValue{{AttributePath: "key", ExpectedValue: "key_1"}}}}, Versions: []string{"2.0(1a)-"}},
		{assertionType: testAssertionNestedSetElement, AttributePath: "children.*", NestedValues: []TestAssertionValue{{AttributePath: "key", ExpectedValue: "key_1"}, {AttributePath: "details.value", ExpectedValue: "detail_1"}}, Versions: []string{"2.0(1a)-"}},
		{assertionType: testAssertionNestedSetElement, AttributePath: "children.*", NestedValues: []TestAssertionValue{{AttributePath: "key", ExpectedValue: "key_2"}, {AttributePath: "labels.#", ExpectedValue: "1"}}, Versions: []string{"2.0(1a)-"}},
		{assertionType: testAssertionNestedSetElement, AttributePath: "children.*", NestedValues: []TestAssertionValue{{AttributePath: "key", ExpectedValue: "key_2"}, {AttributePath: "target", PairAddress: "aci_tenant.test_2.0", PairAttribute: "id"}}, Versions: []string{"2.0(1a)-", "3.0(1a)-"}},
		{assertionType: testAssertionSetElementValue, AttributePath: "children.*.labels.*", ExpectedValue: "label_2", Parents: []TestAssertionParent{{AttributePath: "children.*", Values: []TestAssertionValue{{AttributePath: "key", ExpectedValue: "key_2"}}}}, Versions: []string{"2.0(1a)-"}},
		{assertionType: testAssertionNestedSetElement, AttributePath: "children.*", NestedValues: []TestAssertionValue{{AttributePath: "key", ExpectedValue: "key_2"}, {AttributePath: "details.value", ExpectedValue: "detail_2"}}, Versions: []string{"2.0(1a)-"}},
	}
	assert.Equal(t, expected, parent.DataSourceTestAssertions(store))
}

func TestReferenceChildIdentityAssertions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	childVersions, err := NewVersions("2.0(1a)-")
	if !assert.NoError(t, err) {
		return
	}
	sourceVersions, err := NewVersions("3.0(1a)-")
	if !assert.NoError(t, err) {
		return
	}
	propertyVersions, err := NewVersions("4.0(1a)-")
	if !assert.NoError(t, err) {
		return
	}
	testCases := []test.TestCase{
		{Name: "test_metadata_source_adds_dependency_gate", Input: true, Expected: []string{"2.0(1a)-", "3.0(1a)-"}},
		{Name: "test_definition_only_source_keeps_owner_gate", Input: false, Expected: []string{"2.0(1a)-"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			source := dataSourceTestClass("fvTenant", "tenant", ResourceArtifact)
			source.SupportedVersions = sourceVersions
			details := dataSourceTestClass("exampleDetails", "example_details")
			details.ResourceNameNested = "details"
			details.IsSingleNestedWhenDefinedAsChild = true
			details.SupportedVersions = childVersions
			details.PropertiesAll = []string{"value"}
			details.Properties = map[string]*Property{"value": dataSourceTestProperty("value", "value", nil)}
			child := dataSourceTestClass("exampleChild", "example_child")
			child.ResourceNameNested = "children"
			child.IdentifiedBy = []string{"target"}
			child.SupportedVersions = childVersions
			child.PropertiesAll = []string{"target", "priority"}
			child.Properties = map[string]*Property{
				"target":   dataSourceTestProperty("target", "target", nil),
				"priority": dataSourceTestProperty("priority", "priority", nil),
			}
			child.Properties["priority"].SupportedVersions = propertyVersions
			parent := dataSourceTestClass("exampleParent", "example_parent", DatasourceArtifact)
			parent.TestDependencies = []*TestDependency{{Class: source.Name, Reference: "aci_tenant.test_1.id", ReferenceType: ResourceReference, Role: Target}}
			parent.TestChildren = []*TestChild{{Class: child.Name, Instances: []TestChildInstance{{
				Properties: map[string]TestValueScenario{
					"target":   dataSourceTestScenario("aci_tenant.test_1.id", "aci_tenant.test_1.id", ReferenceValue, true),
					"priority": dataSourceTestScenario("level1", "level1", StringValue, true),
				},
				Children: []*TestChild{{Class: details.Name, Instances: []TestChildInstance{{Properties: map[string]TestValueScenario{
					"value": dataSourceTestScenario("detail_1", "detail_1", StringValue, true),
				}}}}},
			}}}}
			store := &DataStore{Classes: map[string]Class{child.Name.String(): child, details.Name.String(): details}}
			if testCase.Input.(bool) {
				store.Classes[source.Name.String()] = source
			}
			identity := TestAssertionValue{AttributePath: "target", PairAddress: "aci_tenant.test_1.0", PairAttribute: "id"}
			expected := []TestAssertion{
				{assertionType: testAssertionSetCount, AttributePath: "children", ExpectedCount: 1, Versions: []string{"2.0(1a)-"}},
				{assertionType: testAssertionNestedSetElement, AttributePath: "children.*", NestedValues: []TestAssertionValue{identity}, Versions: testCase.Expected.([]string)},
				{assertionType: testAssertionNestedSetElement, AttributePath: "children.*", NestedValues: []TestAssertionValue{identity, {AttributePath: "priority", ExpectedValue: "level1"}}, Versions: []string{"2.0(1a)-", "4.0(1a)-"}},
				{assertionType: testAssertionNestedSetElement, AttributePath: "children.*", NestedValues: []TestAssertionValue{identity, {AttributePath: "details.value", ExpectedValue: "detail_1"}}, Versions: testCase.Expected.([]string)},
			}
			assert.Equal(t, expected, parent.DataSourceTestAssertions(store))
		})
	}
}

func TestPrimitiveSetAssertionParents(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	type primitiveSetInput struct {
		Singleton bool
		Entries   []TestValueEntry
	}
	testCases := []test.TestCase{
		{Name: "test_empty_set", Input: primitiveSetInput{}, Expected: "0"},
		{Name: "test_literal_members", Input: primitiveSetInput{Entries: []TestValueEntry{
			{ConfigValue: "first", AssertValue: "first", ConfigInclude: true},
			{ConfigValue: "second", AssertValue: "second", ConfigInclude: true},
		}}, Expected: "2"},
		{Name: "test_reference_member", Input: primitiveSetInput{Entries: []TestValueEntry{
			{ConfigValue: "aci_tenant.target_1.id", AssertValue: "aci_tenant.target_1.id", ConfigInclude: true, ValueType: ReferenceValue},
		}}, Expected: "1"},
		{Name: "test_singleton_inside_repeated_child", Input: primitiveSetInput{Singleton: true, Entries: []TestValueEntry{
			{ConfigValue: "first", AssertValue: "first", ConfigInclude: true},
		}}, Expected: "1"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(primitiveSetInput)
			versions, err := NewVersions("4.0(1a)-")
			if !assert.NoError(t, err) {
				return
			}
			owner := dataSourceTestClass("exampleOwner", "example_owner")
			owner.ResourceNameNested = "children"
			owner.IdentifiedBy = []string{"key"}
			owner.PropertiesAll = []string{"members", "key"}
			members := &Property{PropertyName: "members", AttributeName: "members", ValueType: Set, SupportedVersions: versions}
			owner.Properties = map[string]*Property{"key": dataSourceTestProperty("key", "key", nil), "members": members}
			instance := TestChildInstance{Properties: map[string]TestValueScenario{
				"key":     dataSourceTestScenario("aci_tenant.owner_1.id", "aci_tenant.owner_1.id", ReferenceValue, true),
				"members": newTestValueScenario(input.Entries),
			}}
			root := dataSourceTestClass("exampleRoot", "example_root", DatasourceArtifact)
			source := dataSourceTestClass("fvTenant", "tenant", ResourceArtifact)
			store := &DataStore{Classes: map[string]Class{owner.Name.String(): owner, source.Name.String(): source}}
			path := "children.*.members.*"
			countPath := "members.#"
			if input.Singleton {
				details := dataSourceTestClass("exampleDetails", "example_details")
				details.ResourceNameNested = "details"
				details.IsSingleNestedWhenDefinedAsChild = true
				details.PropertiesAll = []string{"members"}
				details.Properties = map[string]*Property{"members": members}
				instance.Children = []*TestChild{{Class: details.Name, Instances: []TestChildInstance{{Properties: map[string]TestValueScenario{"members": instance.Properties["members"]}}}}}
				delete(instance.Properties, "members")
				store.Classes[details.Name.String()] = details
				path = "children.*.details.members.*"
				countPath = "details.members.#"
			}
			root.TestChildren = []*TestChild{{Class: owner.Name, Instances: []TestChildInstance{instance}}}
			root.TestDependencies = []*TestDependency{
				{Class: source.Name, Reference: "aci_tenant.owner_1.id", ReferenceType: ResourceReference},
				{Class: source.Name, Reference: "aci_tenant.target_1.id", ReferenceType: ResourceReference},
			}

			assertions := root.DataSourceTestAssertions(store)

			identity := TestAssertionValue{AttributePath: "key", PairAddress: "aci_tenant.owner_1.0", PairAttribute: "id"}
			countIndex := slices.IndexFunc(assertions, func(a TestAssertion) bool {
				return a.IsNestedSetElement() && a.AttributePath == "children.*" && slices.Contains(a.NestedValues, TestAssertionValue{AttributePath: countPath, ExpectedValue: testCase.Expected.(string)})
			})
			if !assert.NotEqual(t, -1, countIndex) {
				return
			}
			assert.Contains(t, assertions[countIndex].NestedValues, identity)
			assert.Equal(t, []string{"4.0(1a)-"}, assertions[countIndex].Versions)
			var memberAssertions []TestAssertion
			for _, a := range assertions {
				if a.AttributePath == path {
					memberAssertions = append(memberAssertions, a)
					assert.Equal(t, []TestAssertionParent{{AttributePath: "children.*", Values: []TestAssertionValue{identity}}}, a.Parents)
					assert.Equal(t, []string{"4.0(1a)-"}, a.Versions)
				}
			}
			if !assert.Len(t, memberAssertions, len(input.Entries)) {
				return
			}
			if len(input.Entries) > 0 && input.Entries[0].ValueType == ReferenceValue {
				assert.True(t, memberAssertions[0].IsSetElementPair())
				assert.Equal(t, "aci_tenant.target_1.0", memberAssertions[0].PairAddress)
				assert.Equal(t, "id", memberAssertions[0].PairAttribute)
			}
		})
	}
}

func TestResourceTestScenarioAssertStrategies(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	class, store := resourceScenarioFixture()
	property := class.Properties["descr"]
	property.TestValues.Default.Entries[0].AssertStrategy = AssertPresenceOnly
	scenarios, err := class.ResourceTestScenarios(store)
	if !assert.NoError(t, err) {
		return
	}
	testCases := []test.TestCase{
		{Name: "Minimum", Expected: testAssertionAttributeSet},
		{Name: "Full", Expected: testAssertionAttributeValue},
		{Name: "Update", Expected: testAssertionAttributeValue},
		{Name: "OmitOptional", Expected: testAssertionAttributeValue},
		{Name: "Reset", Expected: testAssertionAttributeSet},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			index := slices.IndexFunc(scenarios, func(s ResourceTestScenario) bool { return s.Name == testCase.Name })
			if !assert.NotEqual(t, -1, index) {
				return
			}
			assertions := scenarios[index].Assertions
			index = slices.IndexFunc(assertions, func(a TestAssertion) bool { return a.AttributePath == "description" })
			if assert.NotEqual(t, -1, index) {
				assert.Equal(t, testCase.Expected.(testAssertionType), assertions[index].assertionType)
			}
		})
	}
}

func TestResourceTestSensitiveAssertions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	testCases := []test.TestCase{
		{Name: "test_optional_secret", Input: false, Expected: TestAssertion{assertionType: testAssertionAttributeAbsent, AttributePath: "password"}},
		{Name: "test_computed_secret", Input: true, Expected: TestAssertion{assertionType: testAssertionAttributeValue, AttributePath: "password", ExpectedValue: "password_2"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			class, store := resourceScenarioFixture()
			class.PropertiesAll = append(class.PropertiesAll, "password")
			class.Properties["password"] = dataSourceTestProperty("password", "password", dataSourceTestValues("password_1", "password_2", "", false))
			class.Properties["password"].Sensitive = true
			class.Properties["password"].Computed = testCase.Input.(bool)
			scenarios, err := class.ResourceTestScenarios(store)
			if !assert.NoError(t, err) {
				return
			}
			checks := []test.TestCase{
				{Name: "Minimum", Expected: TestAssertion{assertionType: testAssertionAttributeAbsent, AttributePath: "password"}},
				{Name: "Full", Expected: TestAssertion{assertionType: testAssertionAttributeValue, AttributePath: "password", ExpectedValue: "password_1"}},
				{Name: "Update", Expected: TestAssertion{assertionType: testAssertionAttributeValue, AttributePath: "password", ExpectedValue: "password_2"}},
				{Name: "OmitOptional", Expected: testCase.Expected},
			}
			for _, check := range checks {
				t.Run(check.Name, func(t *testing.T) {
					t.Parallel()
					index := slices.IndexFunc(scenarios, func(scenario ResourceTestScenario) bool { return scenario.Name == check.Name })
					if assert.NotEqual(t, -1, index) {
						assert.Contains(t, scenarios[index].Assertions, check.Expected.(TestAssertion))
						if check.Name == "Minimum" || check.Name == "OmitOptional" {
							assert.NotRegexp(t, `(?m)^\s*password\s*=`, scenarios[index].Configuration)
						}
					}
				})
			}
			assert.Contains(t, class.ResourceTestImportIgnore(store), "password")
		})
	}
}

func TestResourceTestRuntimeExpectations(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	tenant := dataSourceTestClass("fvTenant", "tenant", ResourceArtifact)
	tenant.IdentifiedBy = []string{"name"}
	tenant.RnFormat = "tn-{name}"
	tenant.MetaFileContent = map[string]any{"dnFormats": []any{"uni/tn-{name}"}}
	tenant.PropertiesAll = []string{"name"}
	tenant.Properties = map[string]*Property{"name": dataSourceTestProperty("name", "name", dataSourceTestValues("tenant_1", "tenant_2", "tenant_1", true))}
	root := dataSourceTestClass("exampleRoot", "example_root", ResourceArtifact)
	root.PropertiesAll = []string{"target"}
	root.Properties = map[string]*Property{"target": dataSourceTestProperty("target", "target_dn", dataSourceTestValues("uni/tn-tenant_1/ap-example", "uni/tn-tenant_1/ap-other", "", false))}
	root.Properties["target"].Reference = &PropertyReference{Classes: []*ClassName{tenant.Name}}
	root.TestDependencies = []*TestDependency{{Class: tenant.Name, Reference: "aci_tenant.test_1.id", ReferenceType: ResourceReference, Role: Target, TestInstance: 1}}
	store := &DataStore{Classes: map[string]Class{"fvTenant": tenant, "exampleRoot": root}}
	target := dataSourceTestClass("exampleTarget", "target", ResourceArtifact)
	target.IdentifiedBy = []string{"name"}
	target.RnFormat = "target-{name}"
	target.MetaFileContent = map[string]any{"dnFormats": []any{"uni/target-{name}"}}
	target.PropertiesAll = []string{"name"}
	target.Properties = map[string]*Property{"name": dataSourceTestProperty("name", "name", dataSourceTestValues("target_1", "target_2", "target_1", true))}
	store.Classes["exampleTarget"] = target
	nested := root.resourceTestCopy()
	nested.Properties["target"].Reference = &PropertyReference{Classes: []*ClassName{target.Name}}
	nested.TestDependencies = []*TestDependency{{Class: target.Name, Reference: "aci_target.test_1.id", ReferenceType: ResourceReference, Role: Target, Dependencies: root.TestDependencies}}
	testCases := []test.TestCase{
		{Name: "test_prefixed_name", Input: tenant, Expected: TestAssertion{assertionType: testAssertionAttributeValue, AttributePath: "name", ExpectedValue: "tenant_1", ExpectedExpression: `fmt.Sprintf("%s-%s", testPrefix, "tenant_1")`}},
		{Name: "test_formatted_dependency_dn", Input: root, Expected: TestAssertion{assertionType: testAssertionAttributePair, AttributePath: "target_dn", ExpectedValue: "uni/tn-tenant_1/ap-example", PairAddress: "aci_tenant.test_1.0", PairAttribute: "name", PairFormat: "uni/tn-%s/ap-example"}},
		{Name: "test_formatted_prerequisite_dn", Input: nested, Expected: TestAssertion{assertionType: testAssertionAttributePair, AttributePath: "target_dn", ExpectedValue: "uni/tn-tenant_1/ap-example", PairAddress: "aci_tenant.test_1.0", PairAttribute: "name", PairFormat: "uni/tn-%s/ap-example"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			class := testCase.Input.(Class)
			assertions := class.resourceTestAssertions(store)
			assert.Contains(t, assertions, testCase.Expected.(TestAssertion))
		})
	}
}

func TestResourceTestSemanticStringAssertions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	class, store := resourceScenarioFixture()
	property := class.Properties["descr"]
	property.AttributeName = "arp_learning"
	property.ValueType = VMMArpLearning
	property.TestValues = dataSourceTestValues("enabled", "disabled", "disabled", false)
	scenarios, err := class.ResourceTestScenarios(store)
	if !assert.NoError(t, err) {
		return
	}
	minimum := scenarios[0]
	assert.NotContains(t, minimum.Configuration, "arp_learning")
	assert.True(t, slices.ContainsFunc(minimum.Assertions, func(a TestAssertion) bool {
		return a.AttributePath == "arp_learning" && a.ExpectedValue == "disabled" && a.SemanticStringType == "customtypes.VMMArpLearningType{}"
	}))
	for _, assertion := range minimum.Assertions {
		if assertion.AttributePath != "arp_learning" {
			assert.Empty(t, assertion.SemanticStringType)
		}
	}
}
