package helpers

import (
	"context"
	"slices"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

func mustResponseContainer(t *testing.T, value string) *container.Container {
	t.Helper()

	result, err := container.ParseJSON([]byte(value))
	if err != nil {
		t.Fatalf("parse APIC JSON: %v", err)
	}
	return result
}

type attributesFromObjectInput struct {
	JSON string
	Nil  bool
}

type attributesFromObjectExpected struct {
	Success    bool
	Error      bool
	ValueCount int
}

func TestAttributesFromObject(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:     "test_nil_object",
			Input:    attributesFromObjectInput{Nil: true},
			Expected: attributesFromObjectExpected{Error: true},
		},
		{
			Name:     "test_missing_attributes",
			Input:    attributesFromObjectInput{JSON: `{}`},
			Expected: attributesFromObjectExpected{Error: true},
		},
		{
			Name:     "test_attributes_are_not_an_object",
			Input:    attributesFromObjectInput{JSON: `{"attributes":[]}`},
			Expected: attributesFromObjectExpected{Error: true},
		},
		{
			Name:     "test_valid_attributes",
			Input:    attributesFromObjectInput{JSON: `{"attributes":{"dn":"uni/test-example","name":"example"}}`},
			Expected: attributesFromObjectExpected{Success: true, ValueCount: 2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(attributesFromObjectInput)
			expected := testCase.Expected.(attributesFromObjectExpected)
			var object *container.Container
			if !input.Nil {
				object = mustResponseContainer(t, input.JSON)
			}
			var diagnostics diag.Diagnostics
			attributes, success := AttributesFromObject(context.Background(), &diagnostics, object, "fvTest")

			if success != expected.Success || diagnostics.HasError() != expected.Error || len(attributes.values) != expected.ValueCount {
				t.Fatalf(
					"success=%t error=%t values=%d, expected success=%t error=%t values=%d; diagnostics=%v",
					success, diagnostics.HasError(), len(attributes.values),
					expected.Success, expected.Error, expected.ValueCount,
					diagnostics,
				)
			}
		})
	}
}

type childObjectsInput struct {
	JSON         string
	KnownClasses []string
}

type childObjectsExpected struct {
	Counts map[string]int
	Error  bool
}

func TestChildObjectsByClass(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:     "test_missing_children",
			Input:    childObjectsInput{JSON: `{"attributes":{}}`, KnownClasses: []string{"fvChild"}},
			Expected: childObjectsExpected{Counts: map[string]int{"fvChild": 0}},
		},
		{
			Name:     "test_children_are_not_an_array",
			Input:    childObjectsInput{JSON: `{"children":{}}`, KnownClasses: []string{"fvChild"}},
			Expected: childObjectsExpected{Counts: map[string]int{"fvChild": 0}, Error: true},
		},
		{
			Name: "test_known_children_are_grouped",
			Input: childObjectsInput{
				JSON:         `{"children":[{"fvChild":{"attributes":{"name":"one"}}},{"ignored":{"attributes":{}}},{"fvChild":{"attributes":{"name":"two"}}}]}`,
				KnownClasses: []string{"fvChild"},
			},
			Expected: childObjectsExpected{Counts: map[string]int{"fvChild": 2}},
		},
		{
			Name: "test_malformed_child_envelope",
			Input: childObjectsInput{
				JSON:         `{"children":[42]}`,
				KnownClasses: []string{"fvChild"},
			},
			Expected: childObjectsExpected{Counts: map[string]int{"fvChild": 0}, Error: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(childObjectsInput)
			expected := testCase.Expected.(childObjectsExpected)
			var diagnostics diag.Diagnostics
			objects := ChildObjectsByClass(
				context.Background(),
				&diagnostics,
				mustResponseContainer(t, input.JSON),
				"fvParent",
				input.KnownClasses,
			)

			if diagnostics.HasError() != expected.Error {
				t.Fatalf("error=%t, expected %t; diagnostics=%v", diagnostics.HasError(), expected.Error, diagnostics)
			}
			for className, expectedCount := range expected.Counts {
				if actualCount := len(objects[className]); actualCount != expectedCount {
					t.Fatalf("class %q count %d, expected %d", className, actualCount, expectedCount)
				}
			}
		})
	}
}

type responseTestModel struct {
	Name types.String `tfsdk:"name"`
}

var responseTestAttributeTypes = map[string]attr.Type{"name": types.StringType}

func responseTestObject(name string) types.Object {
	return types.ObjectValueMust(responseTestAttributeTypes, map[string]attr.Value{"name": types.StringValue(name)})
}

func decodeResponseTestObject(ctx context.Context, object *container.Container, fallbackModel *responseTestModel) (responseTestModel, diag.Diagnostics) {
	model := responseTestModel{Name: types.StringNull()}
	if fallbackModel != nil {
		model = *fallbackModel
	}

	var diagnostics diag.Diagnostics
	attributes, ok := AttributesFromObject(ctx, &diagnostics, object, "fvChild")
	if ok {
		DecodeStringAttribute(ctx, &diagnostics, attributes, "name", &model.Name)
	}
	return model, diagnostics
}

type singletonDecodeInput struct {
	Objects      []string
	FallbackName string
}

type singletonDecodeExpected struct {
	Name         string
	Null         bool
	FallbackName string
	Warnings     int
}

func TestDecodeSingletonChild(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:     "test_missing_child",
			Input:    singletonDecodeInput{},
			Expected: singletonDecodeExpected{Null: true},
		},
		{
			Name:     "test_decode_child_with_fallback",
			Input:    singletonDecodeInput{Objects: []string{`{"attributes":{"name":"returned"}}`}, FallbackName: "fallback"},
			Expected: singletonDecodeExpected{Name: "returned", FallbackName: "fallback"},
		},
		{
			Name: "test_multiple_children_use_first",
			Input: singletonDecodeInput{Objects: []string{
				`{"attributes":{"name":"first"}}`,
				`{"attributes":{"name":"second"}}`,
			}},
			Expected: singletonDecodeExpected{Name: "first", Warnings: 1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(singletonDecodeInput)
			expected := testCase.Expected.(singletonDecodeExpected)
			objects := make([]*container.Container, 0, len(input.Objects))
			for _, object := range input.Objects {
				objects = append(objects, mustResponseContainer(t, object))
			}

			var fallback *types.Object
			if input.FallbackName != "" {
				value := responseTestObject(input.FallbackName)
				fallback = &value
			}
			fallbackName := ""
			decode := func(ctx context.Context, object *container.Container, fallbackModel *responseTestModel) (responseTestModel, diag.Diagnostics) {
				if fallbackModel != nil {
					fallbackName = fallbackModel.Name.ValueString()
				}
				return decodeResponseTestObject(ctx, object, fallbackModel)
			}

			var diagnostics diag.Diagnostics
			target := types.ObjectNull(responseTestAttributeTypes)
			DecodeSingletonChild(
				context.Background(),
				&diagnostics,
				objects,
				fallback,
				responseTestAttributeTypes,
				"fvParent",
				"fvChild",
				func() responseTestModel { return responseTestModel{Name: types.StringNull()} },
				decode,
				&target,
			)

			if diagnostics.HasError() || diagnostics.WarningsCount() != expected.Warnings || fallbackName != expected.FallbackName {
				t.Fatalf("error=%t warnings=%d fallback=%q, expected warnings=%d fallback=%q; diagnostics=%v", diagnostics.HasError(), diagnostics.WarningsCount(), fallbackName, expected.Warnings, expected.FallbackName, diagnostics)
			}
			var model responseTestModel
			valueDiagnostics := target.As(context.Background(), &model, basetypes.ObjectAsOptions{})
			if valueDiagnostics.HasError() {
				t.Fatalf("decode target object: %v", valueDiagnostics)
			}
			if model.Name.IsNull() != expected.Null || (!model.Name.IsNull() && model.Name.ValueString() != expected.Name) {
				t.Fatalf("decoded name %s, expected name=%q null=%t", model.Name.String(), expected.Name, expected.Null)
			}
		})
	}
}

func TestDecodeRepeatedChildren(t *testing.T) {
	objects := []*container.Container{
		mustResponseContainer(t, `{"attributes":{"name":"one"}}`),
		mustResponseContainer(t, `{"attributes":{"name":"two"}}`),
	}
	var diagnostics diag.Diagnostics
	target := types.SetNull(types.ObjectType{AttrTypes: responseTestAttributeTypes})
	DecodeRepeatedChildren(
		context.Background(),
		&diagnostics,
		objects,
		responseTestAttributeTypes,
		decodeResponseTestObject,
		&target,
	)
	if diagnostics.HasError() {
		t.Fatalf("decode repeated children: %v", diagnostics)
	}

	var models []responseTestModel
	valueDiagnostics := target.ElementsAs(context.Background(), &models, false)
	if valueDiagnostics.HasError() {
		t.Fatalf("decode target set: %v", valueDiagnostics)
	}
	names := make([]string, 0, len(models))
	for _, model := range models {
		names = append(names, model.Name.ValueString())
	}
	slices.Sort(names)
	if !slices.Equal(names, []string{"one", "two"}) {
		t.Fatalf("decoded names %v, expected [one two]", names)
	}
}

type modelFromResponseInput struct {
	JSON         string
	Nil          bool
	DecodeError  bool
	FallbackName string
}

type modelFromResponseExpected struct {
	Found        bool
	DN           string
	Name         string
	Error        bool
	DecodeCount  int
	FallbackName string
}

func TestModelFromResponse(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "test_nil_response", Input: modelFromResponseInput{Nil: true}, Expected: modelFromResponseExpected{Error: true}},
		{Name: "test_missing_imdata", Input: modelFromResponseInput{JSON: `{}`}, Expected: modelFromResponseExpected{Error: true}},
		{Name: "test_imdata_is_not_an_array", Input: modelFromResponseInput{JSON: `{"imdata":{}}`}, Expected: modelFromResponseExpected{Error: true}},
		{Name: "test_empty_response", Input: modelFromResponseInput{JSON: `{"imdata":[]}`}, Expected: modelFromResponseExpected{}},
		{
			Name:     "test_unrelated_class",
			Input:    modelFromResponseInput{JSON: `{"imdata":[{"other":{"attributes":{"dn":"uni/other"}}}]}`},
			Expected: modelFromResponseExpected{},
		},
		{
			Name:     "test_malformed_object_envelope",
			Input:    modelFromResponseInput{JSON: `{"imdata":[42]}`},
			Expected: modelFromResponseExpected{Error: true},
		},
		{
			Name: "test_duplicate_objects",
			Input: modelFromResponseInput{JSON: `{"imdata":[
				{"fvTest":{"attributes":{"dn":"uni/test-one"}}},
				{"fvTest":{"attributes":{"dn":"uni/test-two"}}}
			]}`},
			Expected: modelFromResponseExpected{Error: true},
		},
		{
			Name:     "test_missing_attributes",
			Input:    modelFromResponseInput{JSON: `{"imdata":[{"fvTest":{}}]}`},
			Expected: modelFromResponseExpected{Error: true},
		},
		{
			Name:     "test_missing_dn",
			Input:    modelFromResponseInput{JSON: `{"imdata":[{"fvTest":{"attributes":{"name":"example"}}}]}`},
			Expected: modelFromResponseExpected{Error: true},
		},
		{
			Name:     "test_dn_is_not_a_string",
			Input:    modelFromResponseInput{JSON: `{"imdata":[{"fvTest":{"attributes":{"dn":42}}}]}`},
			Expected: modelFromResponseExpected{Error: true},
		},
		{
			Name: "test_valid_response_with_fallback",
			Input: modelFromResponseInput{
				JSON:         `{"imdata":[{"fvTest":{"attributes":{"dn":"uni/test-example","name":"returned"}}}]}`,
				FallbackName: "fallback",
			},
			Expected: modelFromResponseExpected{
				Found:        true,
				DN:           "uni/test-example",
				Name:         "returned",
				DecodeCount:  1,
				FallbackName: "fallback",
			},
		},
		{
			Name: "test_decode_diagnostics_are_propagated",
			Input: modelFromResponseInput{
				JSON:        `{"imdata":[{"fvTest":{"attributes":{"dn":"uni/test-example","name":"returned"}}}]}`,
				DecodeError: true,
			},
			Expected: modelFromResponseExpected{Found: true, DN: "uni/test-example", Name: "returned", Error: true, DecodeCount: 1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(modelFromResponseInput)
			expected := testCase.Expected.(modelFromResponseExpected)
			var response *container.Container
			if !input.Nil {
				response = mustResponseContainer(t, input.JSON)
			}

			var fallback *responseTestModel
			if input.FallbackName != "" {
				fallback = &responseTestModel{Name: types.StringValue(input.FallbackName)}
			}
			decodeCount := 0
			fallbackName := ""
			decode := func(ctx context.Context, object *container.Container, fallbackModel *responseTestModel) (responseTestModel, diag.Diagnostics) {
				decodeCount++
				if fallbackModel != nil {
					fallbackName = fallbackModel.Name.ValueString()
				}
				model, diagnostics := decodeResponseTestObject(ctx, object, fallbackModel)
				if input.DecodeError {
					diagnostics.AddError("decode failed", "test error")
				}
				return model, diagnostics
			}

			var diagnostics diag.Diagnostics
			model, dn := ModelFromResponse(
				context.Background(),
				&diagnostics,
				response,
				"fvTest",
				fallback,
				decode,
			)

			if (model != nil) != expected.Found || dn != expected.DN || diagnostics.HasError() != expected.Error || decodeCount != expected.DecodeCount || fallbackName != expected.FallbackName {
				t.Fatalf(
					"found=%t dn=%q error=%t decodes=%d fallback=%q, expected found=%t dn=%q error=%t decodes=%d fallback=%q; diagnostics=%v",
					model != nil, dn, diagnostics.HasError(), decodeCount, fallbackName,
					expected.Found, expected.DN, expected.Error, expected.DecodeCount, expected.FallbackName,
					diagnostics,
				)
			}
			if model != nil && model.Name.ValueString() != expected.Name {
				t.Fatalf("decoded name %q, expected %q", model.Name.ValueString(), expected.Name)
			}
		})
	}
}
