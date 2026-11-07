package data

import (
	"fmt"
	"slices"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

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
