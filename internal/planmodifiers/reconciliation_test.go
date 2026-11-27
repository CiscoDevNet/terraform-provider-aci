package planmodifiers

import (
	"context"
	"maps"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/customtypes"
)

type singletonReconciliationInput struct {
	ConfigValue types.Object
	PlanValue   types.Object
	StateValue  types.Object
}

func TestReconciliationSchemaReconcileSingletonValue(t *testing.T) {
	attributeTypes := map[string]attr.Type{
		"id":    types.StringType,
		"value": types.StringType,
	}
	objectValue := func(id, value types.String) types.Object {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"id":    id,
			"value": value,
		})
	}
	stateValue := objectValue(types.StringValue("one"), types.StringValue("state"))
	testCases := []test.TestCase{
		{
			Name: "test_unknown_configuration",
			Input: singletonReconciliationInput{
				ConfigValue: types.ObjectUnknown(attributeTypes),
				PlanValue:   types.ObjectUnknown(attributeTypes),
				StateValue:  stateValue,
			},
			Expected: types.ObjectUnknown(attributeTypes),
		},
		{
			Name: "test_explicit_empty_object",
			Input: singletonReconciliationInput{
				ConfigValue: objectValue(types.StringNull(), types.StringNull()),
				PlanValue:   objectValue(types.StringUnknown(), types.StringUnknown()),
				StateValue:  stateValue,
			},
			Expected: objectValue(types.StringNull(), types.StringNull()),
		},
		{
			Name: "test_explicit_empty_object_with_unknown_plan",
			Input: singletonReconciliationInput{
				ConfigValue: objectValue(types.StringNull(), types.StringNull()),
				PlanValue:   types.ObjectUnknown(attributeTypes),
				StateValue:  stateValue,
			},
			Expected: objectValue(types.StringNull(), types.StringNull()),
		},
		{
			Name: "test_omitted_non_null_property",
			Input: singletonReconciliationInput{
				ConfigValue: objectValue(types.StringValue("one"), types.StringNull()),
				PlanValue:   objectValue(types.StringValue("one"), types.StringUnknown()),
				StateValue:  stateValue,
			},
			Expected: stateValue,
		},
		{
			Name: "test_null_state_property",
			Input: singletonReconciliationInput{
				ConfigValue: objectValue(types.StringValue("one"), types.StringNull()),
				PlanValue:   objectValue(types.StringValue("one"), types.StringUnknown()),
				StateValue:  objectValue(types.StringValue("one"), types.StringNull()),
			},
			Expected: objectValue(types.StringValue("one"), types.StringUnknown()),
		},
		{
			Name: "test_unknown_configured_property",
			Input: singletonReconciliationInput{
				ConfigValue: objectValue(types.StringValue("one"), types.StringUnknown()),
				PlanValue:   objectValue(types.StringValue("one"), types.StringUnknown()),
				StateValue:  stateValue,
			},
			Expected: objectValue(types.StringValue("one"), types.StringUnknown()),
		},
	}

	schema := ReconciliationSchema{}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(singletonReconciliationInput)
			actual, diagnostics := schema.reconcileSingletonValue(
				context.Background(),
				input.ConfigValue,
				input.PlanValue,
				input.StateValue,
				true,
			)

			if diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diagnostics)
			}
			if !actual.Equal(testCase.Expected.(types.Object)) {
				t.Fatalf("reconciled value %#v, expected %#v", actual, testCase.Expected)
			}
		})
	}
}

type repeatedReconciliationInput struct {
	ConfigValue types.Set
	PlanValue   types.Set
	StateValue  types.Set
	Schema      ReconciliationSchema
}

func TestReconciliationSchemaReconcileRepeatedValue(t *testing.T) {
	attributeTypes := map[string]attr.Type{
		"id":    types.StringType,
		"value": types.StringType,
	}
	objectType := types.ObjectType{AttrTypes: attributeTypes}
	objectValue := func(id, value types.String) types.Object {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"id":    id,
			"value": value,
		})
	}
	setValue := func(objects ...types.Object) types.Set {
		elements := make([]attr.Value, len(objects))
		for index, object := range objects {
			elements[index] = object
		}
		return types.SetValueMust(objectType, elements)
	}
	schema := ReconciliationSchema{Identifiers: []string{"id"}}
	testCases := []test.TestCase{
		{
			Name: "test_reordered_elements",
			Input: repeatedReconciliationInput{
				ConfigValue: setValue(
					objectValue(types.StringValue("one"), types.StringNull()),
					objectValue(types.StringValue("two"), types.StringNull()),
				),
				PlanValue: setValue(
					objectValue(types.StringValue("one"), types.StringUnknown()),
					objectValue(types.StringValue("two"), types.StringUnknown()),
				),
				StateValue: setValue(
					objectValue(types.StringValue("two"), types.StringValue("state-two")),
					objectValue(types.StringValue("one"), types.StringValue("state-one")),
				),
				Schema: schema,
			},
			Expected: setValue(
				objectValue(types.StringValue("one"), types.StringValue("state-one")),
				objectValue(types.StringValue("two"), types.StringValue("state-two")),
			),
		},
		{
			Name: "test_cardinality_change",
			Input: repeatedReconciliationInput{
				ConfigValue: setValue(
					objectValue(types.StringValue("two"), types.StringNull()),
					objectValue(types.StringValue("three"), types.StringValue("configured-three")),
				),
				PlanValue: setValue(
					objectValue(types.StringValue("two"), types.StringUnknown()),
					objectValue(types.StringValue("three"), types.StringValue("configured-three")),
				),
				StateValue: setValue(
					objectValue(types.StringValue("one"), types.StringValue("state-one")),
					objectValue(types.StringValue("two"), types.StringValue("state-two")),
				),
				Schema: schema,
			},
			Expected: setValue(
				objectValue(types.StringValue("two"), types.StringValue("state-two")),
				objectValue(types.StringValue("three"), types.StringValue("configured-three")),
			),
		},
		{
			Name: "test_null_state_property",
			Input: repeatedReconciliationInput{
				ConfigValue: setValue(objectValue(types.StringValue("one"), types.StringValue("configured"))),
				PlanValue:   setValue(objectValue(types.StringValue("one"), types.StringValue("configured"))),
				StateValue:  setValue(objectValue(types.StringValue("one"), types.StringNull())),
				Schema:      schema,
			},
			Expected: setValue(objectValue(types.StringValue("one"), types.StringNull())),
		},
		{
			Name: "test_unknown_configured_property",
			Input: repeatedReconciliationInput{
				ConfigValue: setValue(objectValue(types.StringValue("one"), types.StringUnknown())),
				PlanValue:   setValue(objectValue(types.StringValue("one"), types.StringUnknown())),
				StateValue:  setValue(objectValue(types.StringValue("one"), types.StringValue("state"))),
				Schema:      schema,
			},
			Expected: setValue(objectValue(types.StringValue("one"), types.StringUnknown())),
		},
		{
			Name: "test_unknown_identifier",
			Input: repeatedReconciliationInput{
				ConfigValue: setValue(objectValue(types.StringUnknown(), types.StringNull())),
				PlanValue:   setValue(objectValue(types.StringUnknown(), types.StringUnknown())),
				StateValue:  setValue(objectValue(types.StringValue("one"), types.StringValue("state"))),
				Schema:      schema,
			},
			Expected: setValue(objectValue(types.StringUnknown(), types.StringUnknown())),
		},
		{
			Name: "test_null_repeated_state",
			Input: repeatedReconciliationInput{
				ConfigValue: types.SetUnknown(objectType),
				PlanValue:   types.SetUnknown(objectType),
				StateValue:  types.SetNull(objectType),
				Schema:      schema,
			},
			Expected: types.SetNull(objectType),
		},
		{
			Name: "test_no_identifiers_does_not_use_position",
			Input: repeatedReconciliationInput{
				ConfigValue: setValue(objectValue(types.StringValue("one"), types.StringNull())),
				PlanValue:   setValue(objectValue(types.StringValue("one"), types.StringUnknown())),
				StateValue:  setValue(objectValue(types.StringValue("one"), types.StringValue("state"))),
				Schema:      ReconciliationSchema{},
			},
			Expected: setValue(objectValue(types.StringValue("one"), types.StringUnknown())),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(repeatedReconciliationInput)
			actual, diagnostics := input.Schema.reconcileRepeatedValue(
				context.Background(),
				input.ConfigValue,
				input.PlanValue,
				input.StateValue,
				true,
			)

			if diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diagnostics)
			}
			if !actual.Equal(testCase.Expected.(types.Set)) {
				t.Fatalf("reconciled value %#v, expected %#v", actual, testCase.Expected)
			}
		})
	}
}

func TestReconciliationSchemaMatchesAllIdentifiers(t *testing.T) {
	attributeTypes := map[string]attr.Type{
		"first":  types.StringType,
		"second": types.StringType,
		"value":  types.StringType,
	}
	objectType := types.ObjectType{AttrTypes: attributeTypes}
	objectValue := func(first, second string, value types.String) types.Object {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"first":  types.StringValue(first),
			"second": types.StringValue(second),
			"value":  value,
		})
	}
	setValue := func(objects ...types.Object) types.Set {
		elements := make([]attr.Value, len(objects))
		for index, object := range objects {
			elements[index] = object
		}
		return types.SetValueMust(objectType, elements)
	}

	configValue := setValue(
		objectValue("same", "one", types.StringNull()),
		objectValue("same", "two", types.StringNull()),
	)
	planValue := setValue(
		objectValue("same", "one", types.StringUnknown()),
		objectValue("same", "two", types.StringUnknown()),
	)
	stateValue := setValue(
		objectValue("same", "two", types.StringValue("state-two")),
		objectValue("same", "one", types.StringValue("state-one")),
	)
	expected := setValue(
		objectValue("same", "one", types.StringValue("state-one")),
		objectValue("same", "two", types.StringValue("state-two")),
	)
	schema := ReconciliationSchema{Identifiers: []string{"first", "second"}}

	actual, diagnostics := schema.reconcileRepeatedValue(context.Background(), configValue, planValue, stateValue, true)

	if diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diagnostics)
	}
	if !actual.Equal(expected) {
		t.Fatalf("reconciled value %#v, expected %#v", actual, expected)
	}
}

func TestReconciliationSchemaMatchesSemanticIdentifiers(t *testing.T) {
	fromType := customtypes.NewSemanticStringType(map[string]string{"10": "AF11"})
	toType := customtypes.NewSemanticStringType(map[string]string{"20": "AF22"})
	attributeTypes := map[string]attr.Type{
		"from":  fromType,
		"to":    toType,
		"value": types.StringType,
	}
	objectType := types.ObjectType{AttrTypes: attributeTypes}
	objectValue := func(from, to string, value types.String) types.Object {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"from":  fromType.Value(from),
			"to":    toType.Value(to),
			"value": value,
		})
	}
	setValue := func(object types.Object) types.Set {
		return types.SetValueMust(objectType, []attr.Value{object})
	}

	configValue := setValue(objectValue("10", "20", types.StringNull()))
	planValue := setValue(objectValue("10", "20", types.StringUnknown()))
	stateValue := setValue(objectValue("AF11", "AF22", types.StringValue("state")))
	expected := setValue(objectValue("10", "20", types.StringValue("state")))
	schema := ReconciliationSchema{Identifiers: []string{"from", "to"}}

	actual, diagnostics := schema.reconcileRepeatedValue(context.Background(), configValue, planValue, stateValue, true)

	if diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diagnostics)
	}
	if !actual.Equal(expected) {
		t.Fatalf("reconciled value %#v, expected %#v", actual, expected)
	}
}

func TestReconciliationSchemaConsumesSemanticMatchOnce(t *testing.T) {
	fromType := customtypes.NewSemanticStringType(map[string]string{"10": "AF11"})
	toType := customtypes.NewSemanticStringType(map[string]string{"20": "AF22"})
	attributeTypes := map[string]attr.Type{
		"from":  fromType,
		"to":    toType,
		"value": types.StringType,
	}
	objectType := types.ObjectType{AttrTypes: attributeTypes}
	objectValue := func(from string, value types.String) types.Object {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"from":  fromType.Value(from),
			"to":    toType.Value("20"),
			"value": value,
		})
	}
	setValue := func(objects ...types.Object) types.Set {
		elements := make([]attr.Value, len(objects))
		for index, object := range objects {
			elements[index] = object
		}
		return types.SetValueMust(objectType, elements)
	}

	configValue := setValue(
		objectValue("10", types.StringNull()),
		objectValue("AF11", types.StringNull()),
	)
	planValue := setValue(
		objectValue("10", types.StringUnknown()),
		objectValue("AF11", types.StringUnknown()),
	)
	stateValue := setValue(objectValue("AF11", types.StringValue("state")))
	schema := ReconciliationSchema{Identifiers: []string{"from", "to"}}

	actual, diagnostics := schema.reconcileRepeatedValue(context.Background(), configValue, planValue, stateValue, true)

	if diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diagnostics)
	}
	knownStateValues := 0
	unknownValues := 0
	for _, element := range actual.Elements() {
		value := element.(types.Object).Attributes()["value"].(types.String)
		if value.IsUnknown() {
			unknownValues++
		} else if value.ValueString() == "state" {
			knownStateValues++
		}
	}
	if knownStateValues != 1 || unknownValues != 1 {
		t.Fatalf("matched state %d times and retained %d unknown values, expected one of each", knownStateValues, unknownValues)
	}
}

func TestReconciliationSchemaPreservesNestedSingletonDeletion(t *testing.T) {
	childAttributeTypes := map[string]attr.Type{"detail": types.StringType}
	childType := types.ObjectType{AttrTypes: childAttributeTypes}
	childValue := func(detail types.String) types.Object {
		return types.ObjectValueMust(childAttributeTypes, map[string]attr.Value{"detail": detail})
	}
	parentAttributeTypes := map[string]attr.Type{
		"id":    types.StringType,
		"child": childType,
	}
	parentType := types.ObjectType{AttrTypes: parentAttributeTypes}
	parentValue := func(child types.Object) types.Object {
		return types.ObjectValueMust(parentAttributeTypes, map[string]attr.Value{
			"id":    types.StringValue("one"),
			"child": child,
		})
	}
	setValue := func(parent types.Object) types.Set {
		return types.SetValueMust(parentType, []attr.Value{parent})
	}

	configValue := setValue(parentValue(childValue(types.StringNull())))
	planValue := setValue(parentValue(childValue(types.StringUnknown())))
	stateValue := setValue(parentValue(childValue(types.StringValue("state"))))
	schema := ReconciliationSchema{
		Identifiers: []string{"id"},
		Children: map[string]ReconciliationChild{
			"child": {Schema: ReconciliationSchema{}},
		},
	}

	actual, diagnostics := schema.reconcileRepeatedValue(context.Background(), configValue, planValue, stateValue, true)

	if diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diagnostics)
	}
	if !actual.Equal(configValue) {
		t.Fatalf("reconciled value %#v, expected explicit empty child %#v", actual, configValue)
	}
}

func TestReconciliationSchemaReconcilesNestedRepeatedChildren(t *testing.T) {
	childAttributeTypes := map[string]attr.Type{
		"id":    types.StringType,
		"value": types.StringType,
	}
	childType := types.ObjectType{AttrTypes: childAttributeTypes}
	childValue := func(id string, value types.String) types.Object {
		return types.ObjectValueMust(childAttributeTypes, map[string]attr.Value{
			"id":    types.StringValue(id),
			"value": value,
		})
	}
	childSet := func(children ...types.Object) types.Set {
		elements := make([]attr.Value, len(children))
		for index, child := range children {
			elements[index] = child
		}
		return types.SetValueMust(childType, elements)
	}
	parentAttributeTypes := map[string]attr.Type{
		"name":     types.StringType,
		"children": types.SetType{ElemType: childType},
	}
	parentValue := func(children types.Set) types.Object {
		return types.ObjectValueMust(parentAttributeTypes, map[string]attr.Value{
			"name":     types.StringValue("parent"),
			"children": children,
		})
	}

	configValue := parentValue(childSet(
		childValue("one", types.StringNull()),
		childValue("two", types.StringNull()),
	))
	planValue := parentValue(childSet(
		childValue("one", types.StringUnknown()),
		childValue("two", types.StringUnknown()),
	))
	stateValue := parentValue(childSet(
		childValue("two", types.StringValue("state-two")),
		childValue("one", types.StringValue("state-one")),
	))
	expected := parentValue(childSet(
		childValue("one", types.StringValue("state-one")),
		childValue("two", types.StringValue("state-two")),
	))
	schema := ReconciliationSchema{
		Children: map[string]ReconciliationChild{
			"children": {
				Repeated: true,
				Schema:   ReconciliationSchema{Identifiers: []string{"id"}},
			},
		},
	}

	actual, diagnostics := schema.reconcileSingletonValue(context.Background(), configValue, planValue, stateValue, true)

	if diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diagnostics)
	}
	if !actual.Equal(expected) {
		t.Fatalf("reconciled value %#v, expected %#v", actual, expected)
	}
}

func TestReconciliationSchemaReconcilesDeepChildren(t *testing.T) {
	ctx := context.Background()
	protocolType := customtypes.NewSemanticStringType(map[string]string{"6": "tcp"})
	attributeTypes := map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"flags":       types.SetType{ElemType: types.StringType},
		"address":     customtypes.IPAddressType{},
		"protocol":    protocolType,
		"learning":    customtypes.VMMArpLearningType{},
	}
	stateAttributes := map[string]attr.Value{
		"name":        types.StringValue("configured-name"),
		"description": types.StringValue("retained-description"),
		"flags":       types.SetValueMust(types.StringType, []attr.Value{types.StringValue("one"), types.StringValue("two")}),
		"address":     customtypes.NewIPAddressValue("2001:db8::1"),
		"protocol":    protocolType.Value("tcp"),
		"learning":    customtypes.NewVMMArpLearningValue("disabled"),
	}
	type input struct {
		Depth    int
		ClearSet bool
	}
	testCases := []test.TestCase{
		{Name: "grandchildren", Input: input{Depth: 3}},
		{Name: "deeper_children", Input: input{Depth: 6}},
		{Name: "explicit_empty_set", Input: input{Depth: 4, ClearSet: true}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			values := testCase.Input.(input)
			configAttributes := map[string]attr.Value{}
			planAttributes := map[string]attr.Value{}
			expectedAttributes := maps.Clone(stateAttributes)
			for name, attributeType := range attributeTypes {
				var err error
				configAttributes[name], err = attributeType.ValueFromTerraform(ctx, tftypes.NewValue(attributeType.TerraformType(ctx), nil))
				if err != nil {
					t.Fatal(err)
				}
				planAttributes[name], err = attributeType.ValueFromTerraform(ctx, tftypes.NewValue(attributeType.TerraformType(ctx), tftypes.UnknownValue))
				if err != nil {
					t.Fatal(err)
				}
			}
			// Keep the singleton configured; an object with only null fields requests deletion.
			configAttributes["name"], planAttributes["name"] = stateAttributes["name"], stateAttributes["name"]
			if values.ClearSet {
				empty := types.SetValueMust(types.StringType, []attr.Value{})
				configAttributes["flags"], planAttributes["flags"], expectedAttributes["flags"] = empty, empty, empty
			}
			config := types.ObjectValueMust(attributeTypes, configAttributes)
			plan := types.ObjectValueMust(attributeTypes, planAttributes)
			state := types.ObjectValueMust(attributeTypes, stateAttributes)
			expected := types.ObjectValueMust(attributeTypes, expectedAttributes)
			schema := ReconciliationSchema{}

			// Alternate singleton objects and unordered sets without changing the leaf value types.
			for depth := range values.Depth {
				childType := types.ObjectType{AttrTypes: config.AttributeTypes(ctx)}
				if depth%2 == 0 {
					memberTypes := map[string]attr.Type{"name": types.StringType, "observed": types.StringType, "child": childType}
					memberType := types.ObjectType{AttrTypes: memberTypes}
					rootTypes := map[string]attr.Type{"children": types.SetType{ElemType: memberType}}
					wrap := func(child types.Object, first, second types.String, reverse bool) types.Object {
						members := []attr.Value{
							types.ObjectValueMust(memberTypes, map[string]attr.Value{"name": types.StringValue("one"), "observed": first, "child": child}),
							types.ObjectValueMust(memberTypes, map[string]attr.Value{"name": types.StringValue("two"), "observed": second, "child": child}),
						}
						if reverse {
							members[0], members[1] = members[1], members[0]
						}
						return types.ObjectValueMust(rootTypes, map[string]attr.Value{"children": types.SetValueMust(memberType, members)})
					}
					config = wrap(config, types.StringNull(), types.StringNull(), false)
					plan = wrap(plan, types.StringUnknown(), types.StringUnknown(), false)
					state = wrap(state, types.StringValue("state-one"), types.StringValue("state-two"), true)
					expected = wrap(expected, types.StringValue("state-one"), types.StringValue("state-two"), false)
					schema = ReconciliationSchema{Children: map[string]ReconciliationChild{"children": {
						Repeated: true,
						Schema: ReconciliationSchema{Identifiers: []string{"name"}, Children: map[string]ReconciliationChild{
							"child": {Schema: schema},
						}},
					}}}
					continue
				}
				rootTypes := map[string]attr.Type{"child": childType}
				wrap := func(child types.Object) types.Object {
					return types.ObjectValueMust(rootTypes, map[string]attr.Value{"child": child})
				}
				config, plan, state, expected = wrap(config), wrap(plan), wrap(state), wrap(expected)
				schema = ReconciliationSchema{Children: map[string]ReconciliationChild{"child": {Schema: schema}}}
			}

			actual, diagnostics := schema.reconcileSingletonValue(ctx, config, plan, state, true)
			if diagnostics.HasError() {
				t.Fatal(diagnostics)
			}
			if !actual.Equal(expected) {
				t.Fatalf("reconciled value %s, expected %s", actual, expected)
			}
		})
	}
}

func TestReconciliationSchemaMalformedValues(t *testing.T) {
	t.Run("test_missing_identifier", func(t *testing.T) {
		left := types.ObjectValueMust(
			map[string]attr.Type{"value": types.StringType},
			map[string]attr.Value{"value": types.StringValue("left")},
		)
		right := types.ObjectValueMust(
			map[string]attr.Type{"id": types.StringType},
			map[string]attr.Value{"id": types.StringValue("one")},
		)

		_, diagnostics := (ReconciliationSchema{Identifiers: []string{"id"}}).identifiersEqual(context.Background(), left, right)

		if !diagnostics.HasError() {
			t.Fatal("expected a missing identifier diagnostic")
		}
	})

	t.Run("test_unexpected_repeated_child_type", func(t *testing.T) {
		_, diagnostics := (ReconciliationChild{Repeated: true}).reconcile(
			context.Background(),
			"child",
			types.StringNull(),
			types.StringNull(),
			types.StringNull(),
		)

		if !diagnostics.HasError() {
			t.Fatal("expected an unexpected child type diagnostic")
		}
	})
}
