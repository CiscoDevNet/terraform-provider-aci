package acceptancetest

import (
	"maps"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/customtypes"
)

func TestCheckSetPair(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	type input struct {
		Attributes map[string]string
		Reference  map[string]string
	}
	primitive := map[string]string{"expected.#": "2", "expected.100": "first", "expected.200": "second"}
	objects := map[string]string{
		"expected.#":        "2",
		"expected.100.name": "first", "expected.100.value": "one",
		"expected.200.name": "second", "expected.200.value": "two",
	}
	// Synthetic entries share an identifier to expose partial object matching.
	overlappingObjects := map[string]string{
		"expected.#":           "2",
		"expected.100.node_dn": "node-101", "expected.100.mode": "regular", "expected.100.encap": "vlan-100",
		"expected.200.node_dn": "node-101", "expected.200.mode": "regular",
	}
	testCases := []test.TestCase{
		{
			Name: "test_primitive_members_ignore_hashes_and_order",
			Input: input{Attributes: map[string]string{
				"legacy.#": "2", "legacy.300": "second", "legacy.400": "first",
			}, Reference: primitive},
		},
		{
			Name: "test_different_counts_fail",
			Input: input{Attributes: map[string]string{
				"legacy.#": "1", "legacy.300": "first",
			}, Reference: primitive}, Expected: "expected",
		},
		{
			Name: "test_different_primitive_members_fail",
			Input: input{Attributes: map[string]string{
				"legacy.#": "2", "legacy.300": "wrong", "legacy.400": "first",
			}, Reference: primitive}, Expected: "no TypeSet element",
		},
		{
			Name: "test_block_members_ignore_hashes_and_order",
			Input: input{Attributes: map[string]string{
				"legacy.#":        "2",
				"legacy.300.name": "second", "legacy.300.value": "two",
				"legacy.400.name": "first", "legacy.400.value": "one",
			}, Reference: objects},
		},
		{
			Name: "test_values_on_different_objects_cannot_match",
			Input: input{Attributes: map[string]string{
				"legacy.#":        "2",
				"legacy.300.name": "second", "legacy.300.value": "one",
				"legacy.400.name": "first", "legacy.400.value": "two",
			}, Reference: objects}, Expected: "no TypeSet element",
		},
		{
			Name: "test_missing_block_member_fails",
			Input: input{Attributes: map[string]string{
				"legacy.#": "2", "legacy.300.name": "second", "legacy.400.name": "first", "legacy.400.value": "one",
			}, Reference: objects}, Expected: "no TypeSet element",
		},
		{
			Name: "test_unexpected_block_member_fails",
			Input: input{Attributes: map[string]string{
				"legacy.#":        "2",
				"legacy.300.name": "second", "legacy.300.value": "two", "legacy.300.extra": "unexpected",
				"legacy.400.name": "first", "legacy.400.value": "one",
			}, Reference: objects}, Expected: "no TypeSet element",
		},
		{
			Name: "test_overlapping_partial_objects_fail",
			Input: input{Attributes: map[string]string{
				"legacy.#":           "2",
				"legacy.300.node_dn": "node-101", "legacy.300.mode": "regular", "legacy.300.encap": "vlan-100",
				"legacy.400.node_dn": "node-101",
			}, Reference: overlappingObjects}, Expected: "no TypeSet element",
		},
		{
			Name: "test_identical_overlapping_objects_ignore_hashes_and_order",
			Input: input{Attributes: map[string]string{
				"legacy.#":           "2",
				"legacy.300.node_dn": "node-101", "legacy.300.mode": "regular",
				"legacy.400.node_dn": "node-101", "legacy.400.mode": "regular", "legacy.400.encap": "vlan-100",
			}, Reference: overlappingObjects},
		},
		{
			Name: "test_empty_value_differs_from_missing_member",
			Input: input{Attributes: map[string]string{
				"legacy.#": "1", "legacy.300.name": "first",
			}, Reference: map[string]string{
				"expected.#": "1", "expected.100.name": "first", "expected.100.value": "",
			}}, Expected: "no TypeSet element",
		},
		{
			Name:  "test_empty_sets",
			Input: input{Attributes: map[string]string{"legacy.#": "0"}, Reference: map[string]string{"expected.#": "0"}},
		},
		{
			Name:  "test_null_and_empty_keep_framework_count_semantics",
			Input: input{Attributes: map[string]string{}, Reference: map[string]string{"expected.#": "0"}},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			value := testCase.Input.(input)
			state := nestedSetCheckState(nestedSetCheckInput{
				Attributes: value.Attributes,
				Sources:    map[string]map[string]string{"aci_example.test_1.0": value.Reference},
			})
			err := CheckSetPair("data.aci_example.test", "legacy", "aci_example.test_1.0", "expected")(state)
			if expected, _ := testCase.Expected.(string); expected != "" {
				assert.ErrorContains(t, err, expected)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckPresent(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	type input struct {
		Path       string
		Attributes map[string]string
		Parents    []ParentSetElement
	}
	testCases := []test.TestCase{
		{Name: "test_default_site", Input: input{Path: "site_id", Attributes: map[string]string{"site_id": "0"}}},
		{Name: "test_mso_site", Input: input{Path: "site_id", Attributes: map[string]string{"site_id": "4"}}},
		{Name: "test_missing", Input: input{Path: "site_id"}, Expected: true},
		{Name: "test_empty", Input: input{Path: "site_id", Attributes: map[string]string{"site_id": ""}}, Expected: true},
		{Name: "test_scoped_child", Input: input{Path: "children.*.value", Attributes: map[string]string{"children.101.key": "key_1", "children.101.value": "server_value"}, Parents: []ParentSetElement{{Path: "children.*", Values: map[string]AttributeValue{"key": {ExpectedValue: "key_1"}}}}}},
		{Name: "test_other_child_cannot_satisfy_presence", Input: input{Path: "children.*.value", Attributes: map[string]string{"children.101.key": "key_1", "children.202.key": "key_2", "children.202.value": "server_value"}, Parents: []ParentSetElement{{Path: "children.*", Values: map[string]AttributeValue{"key": {ExpectedValue: "key_1"}}}}}, Expected: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			value := testCase.Input.(input)
			state := nestedSetCheckState(nestedSetCheckInput{Attributes: value.Attributes})
			err := CheckPresent("data.aci_example.test", value.Path, value.Parents...)(state)
			if failed, _ := testCase.Expected.(bool); failed {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckSemanticStringAttribute(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	type input struct {
		Attributes map[string]string
		Address    string
		ValueType  basetypes.StringTypable
	}
	testCases := []test.TestCase{
		{Name: "test_empty_equals_disabled", Input: input{Attributes: map[string]string{"arp_learning": ""}}},
		{Name: "test_disabled_equals_disabled", Input: input{Attributes: map[string]string{"arp_learning": "disabled"}}},
		{Name: "test_enabled_differs", Input: input{Attributes: map[string]string{"arp_learning": "enabled"}}, Expected: "expected semantic value"},
		{Name: "test_missing_attribute", Input: input{Attributes: map[string]string{}}, Expected: "missing from state"},
		{Name: "test_missing_resource", Input: input{Address: "aci_missing.test"}, Expected: "missing from state"},
		{Name: "test_type_without_semantic_equality", Input: input{Attributes: map[string]string{"arp_learning": "disabled"}, ValueType: basetypes.StringType{}}, Expected: "does not support string semantic equality"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			value := testCase.Input.(input)
			if value.Address == "" {
				value.Address = "data.aci_example.test"
			}
			if value.ValueType == nil {
				value.ValueType = customtypes.VMMArpLearningType{}
			}
			state := nestedSetCheckState(nestedSetCheckInput{Attributes: value.Attributes})
			err := CheckSemanticStringAttribute(value.Address, "arp_learning", "disabled", value.ValueType)(state)
			if expected, _ := testCase.Expected.(string); expected != "" {
				assert.ErrorContains(t, err, expected)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type nestedSetCheckInput struct {
	Path       string
	Attributes map[string]string
	Sources    map[string]map[string]string
	Values     map[string]AttributeValue
	Value      AttributeValue
	Parents    []ParentSetElement
}

func TestCheckAttributeAndAbsent(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	state := nestedSetCheckState(nestedSetCheckInput{
		Attributes: map[string]string{"target_dn": "uni/tn-prefixed/ap-example", "children.#": "1", "children.101.name": "child_1", "children.101.details.%": "0"},
		Sources:    map[string]map[string]string{"aci_tenant.test_1.0": {"name": "prefixed"}},
	})
	parents := []ParentSetElement{{Path: "children.*", Values: map[string]AttributeValue{"name": {ExpectedValue: "child_1"}}}}
	testCases := []test.TestCase{
		{Name: "test_formatted_reference", Input: CheckAttribute("data.aci_example.test", "target_dn", AttributeValue{PairAddress: "aci_tenant.test_1.0", PairAttribute: "name", PairFormat: "uni/tn-%s/ap-example"})},
		{Name: "test_wrong_formatted_reference", Input: CheckAttribute("data.aci_example.test", "target_dn", AttributeValue{PairAddress: "aci_tenant.test_1.0", PairAttribute: "name", PairFormat: "uni/tn-%s/ap-other"}), Expected: true},
		{Name: "test_absent_singleton", Input: CheckAbsent("data.aci_example.test", "details")},
		{Name: "test_present_singleton", Input: CheckAbsent("data.aci_example.test", "children"), Expected: true},
		{Name: "test_empty_singleton_in_set", Input: CheckAbsent("data.aci_example.test", "children.*.details", parents...)},
		{Name: "test_wrong_parent", Input: CheckAbsent("data.aci_example.test", "children.*.details", ParentSetElement{Path: "children.*", Values: map[string]AttributeValue{"name": {ExpectedValue: "child_2"}}}), Expected: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			err := testCase.Input.(resource.TestCheckFunc)(state)
			if fail, _ := testCase.Expected.(bool); fail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckAbsentEmptyObjects(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	testCases := []test.TestCase{
		{Name: "test_absent", Input: map[string]string{}},
		{Name: "test_all_null_singleton", Input: map[string]string{"details.%": "4"}},
		{Name: "test_nested_empty_containers", Input: map[string]string{"details.%": "4", "details.nested.%": "2", "details.tags.#": "0"}},
		{Name: "test_populated_scalar", Input: map[string]string{"details.%": "4", "details.target_dn": "uni/tn-example"}, Expected: true},
		{Name: "test_explicit_empty_string", Input: map[string]string{"details.%": "4", "details.target_dn": ""}, Expected: true},
		{Name: "test_populated_nested_scalar", Input: map[string]string{"details.%": "4", "details.nested.%": "2", "details.nested.name": "example"}, Expected: true},
		{Name: "test_populated_set", Input: map[string]string{"details.%": "4", "details.tags.#": "1"}, Expected: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			state := nestedSetCheckState(nestedSetCheckInput{Attributes: testCase.Input.(map[string]string)})
			err := CheckAbsent("data.aci_example.test", "details")(state)
			if fail, _ := testCase.Expected.(bool); fail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
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
