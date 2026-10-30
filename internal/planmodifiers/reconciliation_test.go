package planmodifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
