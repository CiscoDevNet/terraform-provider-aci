package helpers

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type payloadTestModel struct {
	Name types.String `tfsdk:"name"`
}

var payloadTestAttributeTypes = map[string]attr.Type{"name": types.StringType}

func payloadTestObject(name types.String) types.Object {
	return types.ObjectValueMust(payloadTestAttributeTypes, map[string]attr.Value{"name": name})
}

func payloadTestObjectPointer(value types.Object) *types.Object {
	return &value
}

func payloadTestSet(names ...string) types.Set {
	elements := make([]attr.Value, 0, len(names))
	for _, name := range names {
		elements = append(elements, payloadTestObject(types.StringValue(name)))
	}
	return types.SetValueMust(types.ObjectType{AttrTypes: payloadTestAttributeTypes}, elements)
}

func payloadTestSetPointer(value types.Set) *types.Set {
	return &value
}

type singletonPayloadInput struct {
	Desired         types.Object
	Prior           *types.Object
	DeleteSupported bool
	BuildError      bool
}

type singletonPayloadExpected struct {
	Success      bool
	Error        bool
	PayloadCount int
	BuildCount   int
	DeleteCount  int
	PriorName    string
}

func TestBuildSingletonChildPayload(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:     "test_null_desired_value",
			Input:    singletonPayloadInput{Desired: types.ObjectNull(payloadTestAttributeTypes)},
			Expected: singletonPayloadExpected{Success: true},
		},
		{
			Name:     "test_unknown_desired_value",
			Input:    singletonPayloadInput{Desired: types.ObjectUnknown(payloadTestAttributeTypes)},
			Expected: singletonPayloadExpected{Success: true},
		},
		{
			Name:     "test_create_child",
			Input:    singletonPayloadInput{Desired: payloadTestObject(types.StringValue("created"))},
			Expected: singletonPayloadExpected{Success: true, PayloadCount: 1, BuildCount: 1},
		},
		{
			Name: "test_update_child",
			Input: singletonPayloadInput{
				Desired: payloadTestObject(types.StringValue("updated")),
				Prior:   payloadTestObjectPointer(payloadTestObject(types.StringValue("prior"))),
			},
			Expected: singletonPayloadExpected{Success: true, PayloadCount: 1, BuildCount: 1, PriorName: "prior"},
		},
		{
			Name:     "test_empty_child_without_prior_value",
			Input:    singletonPayloadInput{Desired: payloadTestObject(types.StringNull())},
			Expected: singletonPayloadExpected{Success: true},
		},
		{
			Name: "test_delete_child",
			Input: singletonPayloadInput{
				Desired:         payloadTestObject(types.StringNull()),
				Prior:           payloadTestObjectPointer(payloadTestObject(types.StringValue("prior"))),
				DeleteSupported: true,
			},
			Expected: singletonPayloadExpected{Success: true, PayloadCount: 1, DeleteCount: 1},
		},
		{
			Name: "test_reject_child_deletion",
			Input: singletonPayloadInput{
				Desired: payloadTestObject(types.StringNull()),
				Prior:   payloadTestObjectPointer(payloadTestObject(types.StringValue("prior"))),
			},
			Expected: singletonPayloadExpected{Error: true},
		},
		{
			Name:     "test_propagate_build_error",
			Input:    singletonPayloadInput{Desired: payloadTestObject(types.StringValue("invalid")), BuildError: true},
			Expected: singletonPayloadExpected{Error: true, BuildCount: 1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(singletonPayloadInput)
			expected := testCase.Expected.(singletonPayloadExpected)
			var diagnostics diag.Diagnostics
			buildCount := 0
			deleteCount := 0
			priorName := ""

			build := func(model *payloadTestModel, _ context.Context, prior *payloadTestModel, create bool, defaultAnnotation string) (map[string]any, diag.Diagnostics) {
				buildCount++
				if !create || defaultAnnotation != "provider-default" {
					t.Fatalf("create=%t defaultAnnotation=%q", create, defaultAnnotation)
				}
				if prior != nil {
					priorName = prior.Name.ValueString()
				}
				if input.BuildError {
					var buildDiagnostics diag.Diagnostics
					buildDiagnostics.AddError("build failed", "test error")
					return nil, buildDiagnostics
				}
				return map[string]any{
					"attributes": map[string]any{"name": model.Name.ValueString()},
				}, nil
			}

			var buildDelete func(*payloadTestModel, context.Context, *diag.Diagnostics) map[string]any
			if input.DeleteSupported {
				buildDelete = func(model *payloadTestModel, _ context.Context, _ *diag.Diagnostics) map[string]any {
					deleteCount++
					return map[string]any{
						"attributes": map[string]any{"name": model.Name.ValueString(), "status": "deleted"},
					}
				}
			}

			payloads, success := BuildSingletonChildPayload(
				context.Background(),
				&diagnostics,
				input.Desired,
				input.Prior,
				"fvChild",
				"Cannot delete child",
				"provider-default",
				build,
				buildDelete,
			)

			if success != expected.Success || diagnostics.HasError() != expected.Error {
				t.Fatalf("success=%t error=%t, expected success=%t error=%t; diagnostics=%v", success, diagnostics.HasError(), expected.Success, expected.Error, diagnostics)
			}
			if len(payloads) != expected.PayloadCount || buildCount != expected.BuildCount || deleteCount != expected.DeleteCount || priorName != expected.PriorName {
				t.Fatalf(
					"payloads=%d builds=%d deletes=%d prior=%q, expected payloads=%d builds=%d deletes=%d prior=%q",
					len(payloads), buildCount, deleteCount, priorName,
					expected.PayloadCount, expected.BuildCount, expected.DeleteCount, expected.PriorName,
				)
			}
			if len(payloads) > 0 {
				if _, ok := payloads[0]["fvChild"]; !ok {
					t.Fatalf("payload does not contain child envelope: %#v", payloads[0])
				}
			}
		})
	}
}

type repeatedPayloadInput struct {
	Desired         types.Set
	Prior           *types.Set
	DeleteSupported bool
	BuildErrorName  string
}

type repeatedPayloadExpected struct {
	Success      bool
	Error        bool
	PayloadCount int
	BuildCount   int
	DeleteCount  int
	MatchedCount int
}

func TestBuildRepeatedChildPayloads(t *testing.T) {
	childObjectType := types.ObjectType{AttrTypes: payloadTestAttributeTypes}
	testCases := []test.TestCase{
		{
			Name:     "test_null_desired_value",
			Input:    repeatedPayloadInput{Desired: types.SetNull(childObjectType)},
			Expected: repeatedPayloadExpected{Success: true},
		},
		{
			Name:     "test_unknown_desired_value",
			Input:    repeatedPayloadInput{Desired: types.SetUnknown(childObjectType)},
			Expected: repeatedPayloadExpected{Success: true},
		},
		{
			Name: "test_reconcile_children",
			Input: repeatedPayloadInput{
				Desired:         payloadTestSet("existing", "created"),
				Prior:           payloadTestSetPointer(payloadTestSet("existing", "stale")),
				DeleteSupported: true,
			},
			Expected: repeatedPayloadExpected{Success: true, PayloadCount: 3, BuildCount: 2, DeleteCount: 1, MatchedCount: 1},
		},
		{
			Name: "test_reject_child_deletion",
			Input: repeatedPayloadInput{
				Desired: payloadTestSet("existing"),
				Prior:   payloadTestSetPointer(payloadTestSet("existing", "stale")),
			},
			Expected: repeatedPayloadExpected{Error: true, PayloadCount: 1, BuildCount: 1, MatchedCount: 1},
		},
		{
			Name: "test_propagate_build_error",
			Input: repeatedPayloadInput{
				Desired:        payloadTestSet("invalid"),
				BuildErrorName: "invalid",
			},
			Expected: repeatedPayloadExpected{Error: true, BuildCount: 1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(repeatedPayloadInput)
			expected := testCase.Expected.(repeatedPayloadExpected)
			var diagnostics diag.Diagnostics
			buildCount := 0
			deleteCount := 0
			matchedCount := 0

			build := func(model *payloadTestModel, _ context.Context, prior *payloadTestModel, create bool, defaultAnnotation string) (map[string]any, diag.Diagnostics) {
				buildCount++
				if !create || defaultAnnotation != "provider-default" {
					t.Fatalf("create=%t defaultAnnotation=%q", create, defaultAnnotation)
				}
				if prior != nil {
					matchedCount++
				}
				if model.Name.ValueString() == input.BuildErrorName {
					var buildDiagnostics diag.Diagnostics
					buildDiagnostics.AddError("build failed", "test error")
					return nil, buildDiagnostics
				}
				return map[string]any{
					"attributes": map[string]any{"name": model.Name.ValueString()},
				}, nil
			}

			var buildDelete func(*payloadTestModel, context.Context, *diag.Diagnostics) map[string]any
			if input.DeleteSupported {
				buildDelete = func(model *payloadTestModel, _ context.Context, _ *diag.Diagnostics) map[string]any {
					deleteCount++
					return map[string]any{
						"attributes": map[string]any{"name": model.Name.ValueString(), "status": "deleted"},
					}
				}
			}

			payloads, success := BuildRepeatedChildPayloads(
				context.Background(),
				&diagnostics,
				input.Desired,
				input.Prior,
				"fvChild",
				"Cannot delete child",
				"provider-default",
				func(model *payloadTestModel) string { return model.Name.ValueString() },
				build,
				buildDelete,
			)

			if success != expected.Success || diagnostics.HasError() != expected.Error {
				t.Fatalf("success=%t error=%t, expected success=%t error=%t; diagnostics=%v", success, diagnostics.HasError(), expected.Success, expected.Error, diagnostics)
			}
			if len(payloads) != expected.PayloadCount || buildCount != expected.BuildCount || deleteCount != expected.DeleteCount || matchedCount != expected.MatchedCount {
				t.Fatalf(
					"payloads=%d builds=%d deletes=%d matches=%d, expected payloads=%d builds=%d deletes=%d matches=%d",
					len(payloads), buildCount, deleteCount, matchedCount,
					expected.PayloadCount, expected.BuildCount, expected.DeleteCount, expected.MatchedCount,
				)
			}
			for _, payload := range payloads {
				if _, ok := payload["fvChild"]; !ok {
					t.Fatalf("payload does not contain child envelope: %#v", payload)
				}
			}
		})
	}
}

func TestNewPayloadObject(t *testing.T) {
	attributes := map[string]any{"name": "example"}
	children := []map[string]any{{"fvChild": map[string]any{"attributes": map[string]any{"name": "child"}}}}
	testCases := []test.TestCase{
		{
			Name:     "test_without_children",
			Input:    false,
			Expected: map[string]any{"attributes": attributes},
		},
		{
			Name:     "test_with_children",
			Input:    true,
			Expected: map[string]any{"attributes": attributes, "children": children},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := NewPayloadObject(context.Background(), nil, attributes, children, testCase.Input.(bool))
			if !reflect.DeepEqual(actual, testCase.Expected) {
				t.Fatalf("payload object %#v, expected %#v", actual, testCase.Expected)
			}
		})
	}
}

func TestNewNestedDeletePayloadObject(t *testing.T) {
	attributes := map[string]any{"name": "example"}
	actual := NewNestedDeletePayloadObject(context.Background(), nil, attributes)
	expected := map[string]any{
		"attributes": map[string]any{"name": "example", "status": "deleted"},
		"children":   []map[string]any{},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("delete payload object %#v, expected %#v", actual, expected)
	}
}

type payloadContainerExpected struct {
	JSON  string
	Error bool
}

func TestNewPayloadContainer(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "test_valid_payload",
			Input: map[string]any{
				"fvTest": map[string]any{"attributes": map[string]any{"name": "example"}},
			},
			Expected: payloadContainerExpected{JSON: `{"fvTest":{"attributes":{"name":"example"}}}`},
		},
		{
			Name:     "test_unmarshalable_payload",
			Input:    map[string]any{"invalid": func() {}},
			Expected: payloadContainerExpected{Error: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			payload := NewPayloadContainer(context.Background(), &diagnostics, testCase.Input.(map[string]any), "test payload")
			expected := testCase.Expected.(payloadContainerExpected)

			if diagnostics.HasError() != expected.Error || (payload == nil) != expected.Error {
				t.Fatalf("payload nil=%t error=%t, expected error=%t; diagnostics=%v", payload == nil, diagnostics.HasError(), expected.Error, diagnostics)
			}
			if payload != nil {
				assertPayloadJSONEqual(t, payload, expected.JSON)
			}
		})
	}
}

func TestNewDeletePayload(t *testing.T) {
	var diagnostics diag.Diagnostics
	payload := NewDeletePayload(context.Background(), &diagnostics, "fvTest", "uni/test-example")
	if diagnostics.HasError() || payload == nil {
		t.Fatalf("construct delete payload: %v", diagnostics)
	}
	assertPayloadJSONEqual(t, payload, `{"fvTest":{"attributes":{"dn":"uni/test-example","status":"deleted"}}}`)
}

func assertPayloadJSONEqual(t *testing.T, actual *container.Container, expectedJSON string) {
	t.Helper()

	var actualValue any
	if err := json.Unmarshal([]byte(actual.String()), &actualValue); err != nil {
		t.Fatalf("unmarshal actual payload: %v", err)
	}
	var expectedValue any
	if err := json.Unmarshal([]byte(expectedJSON), &expectedValue); err != nil {
		t.Fatalf("unmarshal expected payload: %v", err)
	}
	if !reflect.DeepEqual(actualValue, expectedValue) {
		t.Fatalf("payload %#v, expected %#v", actualValue, expectedValue)
	}
}
