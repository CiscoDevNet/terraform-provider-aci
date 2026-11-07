package acceptancetest

import (
	"maps"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type nestedSetCheckInput struct {
	Path       string
	Attributes map[string]string
	Sources    map[string]map[string]string
	Values     map[string]AttributeValue
	Value      AttributeValue
	Parents    []ParentSetElement
}

func TestCheckSetElement(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	parents := []ParentSetElement{{Path: "children.*", Values: map[string]AttributeValue{
		"contract_name": {PairAddress: "aci_contract.test_1.0", PairAttribute: "name"},
	}}}
	attributes := map[string]string{
		"children.#":                        "2",
		"children.101.contract_name":        "contract_1",
		"children.101.labels.#":             "1",
		"children.101.labels.10":            "first",
		"children.101.details.labels.#":     "1",
		"children.101.details.labels.10":    "first",
		"children.101.groups.#":             "1",
		"children.101.groups.300.name":      "group_1",
		"children.101.groups.300.labels.#":  "1",
		"children.101.groups.300.labels.10": "first",
		"children.202.contract_name":        "contract_2",
		"children.202.labels.#":             "1",
		"children.202.labels.20":            "second",
		"children.202.details.labels.#":     "1",
		"children.202.details.labels.20":    "second",
		"children.202.groups.#":             "1",
		"children.202.groups.300.name":      "group_1",
		"children.202.groups.300.labels.#":  "1",
		"children.202.groups.300.labels.20": "second",
	}
	sources := map[string]map[string]string{"aci_contract.test_1.0": {"name": "contract_1", "member": "first"}}
	deepParents := append([]ParentSetElement(nil), parents...)
	deepParents = append(deepParents, ParentSetElement{Path: "children.*.groups.*", Values: map[string]AttributeValue{"name": {ExpectedValue: "group_1"}}})
	testCases := []test.TestCase{
		{Name: "test_literal_member", Input: nestedSetCheckInput{Path: "children.*.labels.*", Attributes: attributes, Sources: sources, Value: AttributeValue{ExpectedValue: "first"}, Parents: parents}},
		{Name: "test_member_on_another_parent_fails", Input: nestedSetCheckInput{Path: "children.*.labels.*", Attributes: attributes, Sources: sources, Value: AttributeValue{ExpectedValue: "second"}, Parents: parents}, Expected: "no TypeSet element"},
		{Name: "test_reference_member", Input: nestedSetCheckInput{Path: "children.*.labels.*", Attributes: attributes, Sources: sources, Value: AttributeValue{PairAddress: "aci_contract.test_1.0", PairAttribute: "member"}, Parents: parents}},
		{Name: "test_missing_member_reference_fails", Input: nestedSetCheckInput{Path: "children.*.labels.*", Attributes: attributes, Sources: sources, Value: AttributeValue{PairAddress: "aci_contract.missing.0", PairAttribute: "member"}, Parents: parents}, Expected: "resolve aci_contract.missing.0.member"},
		{Name: "test_singleton_inside_parent", Input: nestedSetCheckInput{Path: "children.*.details.labels.*", Attributes: attributes, Sources: sources, Value: AttributeValue{ExpectedValue: "first"}, Parents: parents}},
		{Name: "test_singleton_inside_wrong_parent_fails", Input: nestedSetCheckInput{Path: "children.*.details.labels.*", Attributes: attributes, Sources: sources, Value: AttributeValue{ExpectedValue: "second"}, Parents: parents}, Expected: "no TypeSet element"},
		{Name: "test_multiple_ancestors", Input: nestedSetCheckInput{Path: "children.*.groups.*.labels.*", Attributes: attributes, Sources: sources, Value: AttributeValue{ExpectedValue: "first"}, Parents: deepParents}},
		{Name: "test_multiple_ancestors_wrong_parent_fails", Input: nestedSetCheckInput{Path: "children.*.groups.*.labels.*", Attributes: attributes, Sources: sources, Value: AttributeValue{ExpectedValue: "second"}, Parents: deepParents}, Expected: "no TypeSet element"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(nestedSetCheckInput)
			err := CheckSetElement("data.aci_example.test", input.Path, input.Value, input.Parents...)(nestedSetCheckState(input))
			if expected, _ := testCase.Expected.(string); expected != "" {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), expected)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckSetElementReadsCurrentState(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	check := CheckSetElement("data.aci_example.test", "children.*.labels.*", AttributeValue{PairAddress: "aci_contract.test_1.0", PairAttribute: "member"}, ParentSetElement{
		Path: "children.*", Values: map[string]AttributeValue{"contract_name": {PairAddress: "aci_contract.test_1.0", PairAttribute: "name"}},
	})
	for _, value := range []string{"first", "second"} {
		state := nestedSetCheckState(nestedSetCheckInput{
			Attributes: map[string]string{"children.#": "1", "children.101.contract_name": "contract_1", "children.101.labels.#": "1", "children.101.labels.10": value},
			Sources:    map[string]map[string]string{"aci_contract.test_1.0": {"name": "contract_1", "member": value}},
		})
		assert.NoError(t, check(state))
	}
}

func TestCheckNestedSetElementParents(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	sources := map[string]map[string]string{
		"aci_contract.test_1.0": {"name": "contract_1"},
		"aci_contract.test_2.0": {"name": "contract_2"},
	}
	parentValues := map[string]AttributeValue{
		"contract_name": {PairAddress: "aci_contract.test_1.0", PairAttribute: "name"},
		"direction":     {ExpectedValue: "in"},
	}
	parents := []ParentSetElement{{Path: "children.*", Values: parentValues}}
	values := map[string]AttributeValue{"key": {ExpectedValue: "key_1"}, "value": {ExpectedValue: "value_1"}}
	attributes := map[string]string{
		"children.#":                         "2",
		"children.101.contract_name":         "contract_1",
		"children.101.direction":             "in",
		"children.101.annotations.#":         "1",
		"children.101.annotations.100.key":   "key_1",
		"children.101.annotations.100.value": "value_1",
		"children.202.contract_name":         "contract_2",
		"children.202.direction":             "out",
		"children.202.annotations.#":         "1",
		"children.202.annotations.100.key":   "key_1",
		"children.202.annotations.100.value": "value_2",
	}
	missing := maps.Clone(attributes)
	missing["children.101.annotations.#"] = "0"
	delete(missing, "children.101.annotations.100.key")
	delete(missing, "children.101.annotations.100.value")
	swapped := maps.Clone(attributes)
	swapped["children.101.annotations.100.value"] = "value_2"
	swapped["children.202.annotations.100.value"] = "value_1"
	deep := maps.Clone(attributes)
	deep["children.101.details.members.#"] = "1"
	deep["children.101.details.members.300.name"] = "member_1"
	deep["children.101.details.members.300.annotations.#"] = "1"
	deep["children.101.details.members.300.annotations.100.key"] = "key_1"
	deep["children.101.details.members.300.annotations.100.value"] = "value_1"
	deep["children.202.details.members.#"] = "1"
	deep["children.202.details.members.300.name"] = "member_1"
	deep["children.202.details.members.300.annotations.#"] = "2"
	deep["children.202.details.members.300.annotations.100.key"] = "key_1"
	deep["children.202.details.members.300.annotations.100.value"] = "value_2"
	deepParents := append([]ParentSetElement(nil), parents...)
	deepParents = append(deepParents, ParentSetElement{
		Path: "children.*.details.members.*", Values: map[string]AttributeValue{"name": {ExpectedValue: "member_1"}},
	})
	testCases := []test.TestCase{
		{
			Name:  "test_correct_parent_with_reused_grandchild_hash",
			Input: nestedSetCheckInput{Path: "children.*.annotations.*", Attributes: attributes, Sources: sources, Values: values, Parents: parents},
		},
		{
			Name:  "test_missing_grandchild_fails",
			Input: nestedSetCheckInput{Path: "children.*.annotations.*", Attributes: missing, Sources: sources, Values: values, Parents: parents}, Expected: "no TypeSet element",
		},
		{
			Name:  "test_swapped_grandchildren_fail",
			Input: nestedSetCheckInput{Path: "children.*.annotations.*", Attributes: swapped, Sources: sources, Values: values, Parents: parents}, Expected: "no TypeSet element",
		},
		{
			Name: "test_missing_parent_identifier_fails",
			Input: nestedSetCheckInput{Path: "children.*.annotations.*", Attributes: attributes, Sources: sources, Values: values, Parents: []ParentSetElement{{
				Path: "children.*", Values: map[string]AttributeValue{"contract_name": {ExpectedValue: "missing"}},
			}}}, Expected: "no TypeSet element",
		},
		{
			Name: "test_all_parent_identifiers_must_match",
			Input: nestedSetCheckInput{Path: "children.*.annotations.*", Attributes: attributes, Sources: sources, Values: values, Parents: []ParentSetElement{{
				Path: "children.*", Values: map[string]AttributeValue{"contract_name": parentValues["contract_name"], "direction": {ExpectedValue: "out"}},
			}}}, Expected: "no TypeSet element",
		},
		{
			Name:  "test_parent_reference_must_exist",
			Input: nestedSetCheckInput{Path: "children.*.annotations.*", Attributes: attributes, Values: values, Parents: parents}, Expected: "resolve aci_contract.test_1.0.name",
		},
		{
			Name:  "test_multiple_ancestor_sets",
			Input: nestedSetCheckInput{Path: "children.*.details.members.*.annotations.*", Attributes: deep, Sources: sources, Values: values, Parents: deepParents},
		},
		{
			Name: "test_multiple_ancestor_sets_reject_sibling_value",
			Input: nestedSetCheckInput{Path: "children.*.details.members.*.annotations.*", Attributes: deep, Sources: sources, Values: map[string]AttributeValue{
				"key": {ExpectedValue: "key_1"}, "value": {ExpectedValue: "value_2"},
			}, Parents: deepParents}, Expected: "no TypeSet element",
		},
		{
			Name: "test_count_under_selected_ancestor",
			Input: nestedSetCheckInput{Path: "children.*.details.members.*", Attributes: deep, Sources: sources, Values: map[string]AttributeValue{
				"name": {ExpectedValue: "member_1"}, "annotations.#": {ExpectedValue: "1"},
			}, Parents: parents},
		},
		{
			Name: "test_count_from_another_ancestor_fails",
			Input: nestedSetCheckInput{Path: "children.*.details.members.*", Attributes: deep, Sources: sources, Values: map[string]AttributeValue{
				"name": {ExpectedValue: "member_1"}, "annotations.#": {ExpectedValue: "2"},
			}, Parents: parents}, Expected: "no TypeSet element",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(nestedSetCheckInput)
			check := CheckNestedSetElement("data.aci_example.test", input.Path, input.Values, input.Parents...)
			err := check(nestedSetCheckState(input))
			if expected, _ := testCase.Expected.(string); expected != "" {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), expected)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckNestedSetElementParentReadsCurrentState(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	check := CheckNestedSetElement("data.aci_example.test", "children.*.annotations.*", map[string]AttributeValue{
		"key": {ExpectedValue: "key_1"},
	}, ParentSetElement{Path: "children.*", Values: map[string]AttributeValue{
		"contract_name": {PairAddress: "aci_contract.test_1.0", PairAttribute: "name"},
	}})
	for _, name := range []string{"contract_1", "contract_2"} {
		state := nestedSetCheckState(nestedSetCheckInput{
			Attributes: map[string]string{"children.#": "1", "children.101.contract_name": name, "children.101.annotations.#": "1", "children.101.annotations.100.key": "key_1"},
			Sources:    map[string]map[string]string{"aci_contract.test_1.0": {"name": name}},
		})
		assert.NoError(t, check(state))
	}
}

func TestCheckNestedSetElement(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	references := map[string]AttributeValue{
		"contract_name": {PairAddress: "aci_contract.test_1.0", PairAttribute: "name"},
		"priority":      {ExpectedValue: "level1"},
	}
	sources := map[string]map[string]string{
		"aci_contract.test_1.0": {"name": "contract_1"},
		"aci_contract.test_2.0": {"name": "contract_2"},
	}
	attributes := map[string]string{
		"children.#":                 "2",
		"children.101.contract_name": "contract_1",
		"children.101.priority":      "level1",
		"children.101.details.value": "detail_1",
		"children.202.contract_name": "contract_2",
		"children.202.priority":      "unspecified",
		"children.202.details.value": "detail_2",
	}
	testCases := []test.TestCase{
		{
			Name:     "test_correct_child_assignments",
			Input:    nestedSetCheckInput{Attributes: attributes, Sources: sources, Values: references},
			Expected: "",
		},
		{
			Name: "test_swapped_references_fail",
			Input: nestedSetCheckInput{Sources: sources, Values: references, Attributes: map[string]string{
				"children.#":                 "2",
				"children.101.contract_name": "contract_2",
				"children.101.priority":      "level1",
				"children.202.contract_name": "contract_1",
				"children.202.priority":      "unspecified",
			}},
			Expected: "no TypeSet element",
		},
		{
			Name: "test_missing_source_resource_fails",
			Input: nestedSetCheckInput{
				Attributes: attributes,
				Values:     references,
			},
			Expected: "resolve aci_contract.test_1.0.name",
		},
		{
			Name: "test_missing_source_attribute_fails",
			Input: nestedSetCheckInput{
				Attributes: attributes,
				Sources:    map[string]map[string]string{"aci_contract.test_1.0": {}},
				Values:     references,
			},
			Expected: "resolve aci_contract.test_1.0.name",
		},
		{
			Name: "test_literal_only_group",
			Input: nestedSetCheckInput{Attributes: attributes, Values: map[string]AttributeValue{
				"contract_name": {ExpectedValue: "contract_1"},
				"priority":      {ExpectedValue: "level1"},
			}},
			Expected: "",
		},
		{
			Name: "test_empty_literal_retains_existing_helper_behavior",
			Input: nestedSetCheckInput{Attributes: attributes, Values: map[string]AttributeValue{
				"priority":   {ExpectedValue: "level1"},
				"annotation": {ExpectedValue: ""},
			}},
			Expected: "",
		},
		{
			Name: "test_singleton_descendant_uses_parent_reference_identity",
			Input: nestedSetCheckInput{Attributes: attributes, Sources: sources, Values: map[string]AttributeValue{
				"contract_name": {PairAddress: "aci_contract.test_1.0", PairAttribute: "name"},
				"details.value": {ExpectedValue: "detail_1"},
			}},
			Expected: "",
		},
		{
			Name: "test_singleton_descendant_on_another_parent_fails",
			Input: nestedSetCheckInput{Attributes: attributes, Sources: sources, Values: map[string]AttributeValue{
				"contract_name": {PairAddress: "aci_contract.test_1.0", PairAttribute: "name"},
				"details.value": {ExpectedValue: "detail_2"},
			}},
			Expected: "no TypeSet element",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(nestedSetCheckInput)
			check := CheckNestedSetElement("data.aci_example.test", "children.*", input.Values)
			err := check(nestedSetCheckState(input))
			if expected := testCase.Expected.(string); expected != "" {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), expected)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckNestedSetElementReadsCurrentState(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	check := CheckNestedSetElement("data.aci_example.test", "children.*", map[string]AttributeValue{
		"contract_name": {PairAddress: "aci_contract.test_1.0", PairAttribute: "name"},
	})
	for _, name := range []string{"contract_1", "contract_2"} {
		state := nestedSetCheckState(nestedSetCheckInput{
			Attributes: map[string]string{"children.#": "1", "children.101.contract_name": name},
			Sources:    map[string]map[string]string{"aci_contract.test_1.0": {"name": name}},
		})
		assert.NoError(t, check(state))
	}
}

func nestedSetCheckState(input nestedSetCheckInput) *terraform.State {
	resources := map[string]*terraform.ResourceState{
		"data.aci_example.test": {Type: "aci_example", Primary: &terraform.InstanceState{ID: "example", Attributes: input.Attributes}},
	}
	for address, attributes := range input.Sources {
		resources[address] = &terraform.ResourceState{Type: "aci_contract", Primary: &terraform.InstanceState{ID: address, Attributes: attributes}}
	}
	return &terraform.State{Modules: []*terraform.ModuleState{{Path: []string{"root"}, Resources: resources}}}
}
