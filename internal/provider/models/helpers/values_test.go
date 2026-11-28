package helpers

import (
	"context"
	"slices"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

func mustValueAttributes(t *testing.T, value string) Attributes {
	t.Helper()

	object, err := container.ParseJSON([]byte(value))
	if err != nil {
		t.Fatalf("parse attributes object: %v", err)
	}

	var diagnostics diag.Diagnostics
	attributes, ok := AttributesFromObject(context.Background(), &diagnostics, object, "fvTest")
	if !ok || diagnostics.HasError() {
		t.Fatalf("read attributes object: %v", diagnostics)
	}
	return attributes
}

func TestObjectIsEmpty(t *testing.T) {
	attributeTypes := map[string]attr.Type{"name": types.StringType}
	testCases := []test.TestCase{
		{
			Name: "test_null_attribute",
			Input: types.ObjectValueMust(attributeTypes, map[string]attr.Value{
				"name": types.StringNull(),
			}),
			Expected: true,
		},
		{
			Name: "test_known_attribute",
			Input: types.ObjectValueMust(attributeTypes, map[string]attr.Value{
				"name": types.StringValue("example"),
			}),
			Expected: false,
		},
		{
			Name: "test_unknown_attribute",
			Input: types.ObjectValueMust(attributeTypes, map[string]attr.Value{
				"name": types.StringUnknown(),
			}),
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if actual := objectIsEmpty(testCase.Input.(types.Object)); actual != testCase.Expected.(bool) {
				t.Fatalf("objectIsEmpty returned %t, expected %t", actual, testCase.Expected)
			}
		})
	}
}

type decodeStringExpected struct {
	Found bool
	Value string
	Null  bool
	Error bool
}

func TestDecodeStringAttribute(t *testing.T) {
	attributes := mustValueAttributes(t, `{"attributes":{"name":"example","number":42}}`)
	testCases := []test.TestCase{
		{Name: "test_present", Input: "name", Expected: decodeStringExpected{Found: true, Value: "example"}},
		{Name: "test_absent", Input: "missing", Expected: decodeStringExpected{Null: true}},
		{Name: "test_non_string", Input: "number", Expected: decodeStringExpected{Found: true, Null: true, Error: true}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			target := types.StringNull()
			found := DecodeStringAttribute(context.Background(), &diagnostics, attributes, testCase.Input.(string), &target)
			expected := testCase.Expected.(decodeStringExpected)

			if found != expected.Found || diagnostics.HasError() != expected.Error {
				t.Fatalf("found=%t error=%t, expected found=%t error=%t", found, diagnostics.HasError(), expected.Found, expected.Error)
			}
			if target.IsNull() != expected.Null || (!target.IsNull() && target.ValueString() != expected.Value) {
				t.Fatalf("decoded value %s, expected value=%q null=%t", target.String(), expected.Value, expected.Null)
			}
		})
	}
}

func TestDecodeStringAttributeWithEmptyValue(t *testing.T) {
	attributes := mustValueAttributes(t, `{"attributes":{"empty":"","value":"configured"}}`)
	testCases := []test.TestCase{
		{Name: "test_empty_uses_fallback", Input: "empty", Expected: "fallback"},
		{Name: "test_configured_is_preserved", Input: "value", Expected: "configured"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			target := types.StringNull()
			if !DecodeStringAttributeWithEmptyValue(context.Background(), &diagnostics, attributes, testCase.Input.(string), "fallback", &target) {
				t.Fatal("expected attribute to be found")
			}
			if diagnostics.HasError() || target.ValueString() != testCase.Expected.(string) {
				t.Fatalf("decoded value %q, expected %q; diagnostics=%v", target.ValueString(), testCase.Expected, diagnostics)
			}
		})
	}
}

type decodedCustomString struct {
	Value string
}

func TestDecodeCustomStringAttribute(t *testing.T) {
	attributes := mustValueAttributes(t, `{"attributes":{"name":"example"}}`)
	var diagnostics diag.Diagnostics
	var target decodedCustomString

	found := DecodeCustomStringAttribute(
		context.Background(),
		&diagnostics,
		attributes,
		"name",
		&target,
		func(value string) decodedCustomString { return decodedCustomString{Value: value} },
	)
	if !found || diagnostics.HasError() || target.Value != "example" {
		t.Fatalf("found=%t value=%q diagnostics=%v", found, target.Value, diagnostics)
	}
}

type decodeStringSetExpected struct {
	Found  bool
	Values []string
	Null   bool
	Error  bool
}

func TestDecodeStringSetAttribute(t *testing.T) {
	attributes := mustValueAttributes(t, `{"attributes":{"values":"one,two","empty":"","number":42}}`)
	testCases := []test.TestCase{
		{Name: "test_values", Input: "values", Expected: decodeStringSetExpected{Found: true, Values: []string{"one", "two"}}},
		{Name: "test_empty", Input: "empty", Expected: decodeStringSetExpected{Found: true, Values: []string{}}},
		{Name: "test_absent", Input: "missing", Expected: decodeStringSetExpected{Null: true}},
		{Name: "test_non_string", Input: "number", Expected: decodeStringSetExpected{Found: true, Null: true, Error: true}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			target := types.SetNull(types.StringType)
			found := DecodeStringSetAttribute(context.Background(), &diagnostics, attributes, testCase.Input.(string), &target)
			expected := testCase.Expected.(decodeStringSetExpected)

			if found != expected.Found || diagnostics.HasError() != expected.Error || target.IsNull() != expected.Null {
				t.Fatalf("found=%t null=%t error=%t, expected found=%t null=%t error=%t", found, target.IsNull(), diagnostics.HasError(), expected.Found, expected.Null, expected.Error)
			}
			if target.IsNull() {
				return
			}

			var actual []string
			valueDiagnostics := target.ElementsAs(context.Background(), &actual, false)
			if valueDiagnostics.HasError() {
				t.Fatalf("decode resulting set: %v", valueDiagnostics)
			}
			slices.Sort(actual)
			slices.Sort(expected.Values)
			if !slices.Equal(actual, expected.Values) {
				t.Fatalf("decoded values %v, expected %v", actual, expected.Values)
			}
		})
	}
}

type payloadStringExpected struct {
	Added bool
	Value string
}

func TestAddPayloadStringAttribute(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "test_value", Input: types.StringValue("example"), Expected: payloadStringExpected{Added: true, Value: "example"}},
		{Name: "test_null", Input: types.StringNull(), Expected: payloadStringExpected{}},
		{Name: "test_unknown", Input: types.StringUnknown(), Expected: payloadStringExpected{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			attributes := map[string]any{}
			added := AddPayloadStringAttribute(context.Background(), &diagnostics, attributes, "name", testCase.Input.(types.String))
			expected := testCase.Expected.(payloadStringExpected)

			if added != expected.Added || diagnostics.HasError() {
				t.Fatalf("added=%t diagnostics=%v, expected added=%t", added, diagnostics, expected.Added)
			}
			if actual, ok := attributes["name"]; ok != expected.Added || (ok && actual != expected.Value) {
				t.Fatalf("payload value=%v present=%t, expected value=%q present=%t", actual, ok, expected.Value, expected.Added)
			}
		})
	}
}

type payloadStringSetExpected struct {
	Added bool
	Value string
	Error bool
}

func TestAddPayloadStringSetAttribute(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "test_values",
			Input: types.SetValueMust(types.StringType, []attr.Value{
				types.StringValue("one"),
				types.StringValue("two"),
			}),
			Expected: payloadStringSetExpected{Added: true, Value: "one,two"},
		},
		{Name: "test_empty", Input: types.SetValueMust(types.StringType, []attr.Value{}), Expected: payloadStringSetExpected{Added: true}},
		{Name: "test_null", Input: types.SetNull(types.StringType), Expected: payloadStringSetExpected{}},
		{Name: "test_unknown", Input: types.SetUnknown(types.StringType), Expected: payloadStringSetExpected{}},
		{
			Name:     "test_invalid_element_type",
			Input:    types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(1)}),
			Expected: payloadStringSetExpected{Error: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			attributes := map[string]any{}
			added := AddPayloadStringSetAttribute(context.Background(), &diagnostics, attributes, "values", testCase.Input.(types.Set))
			expected := testCase.Expected.(payloadStringSetExpected)

			if added != expected.Added || diagnostics.HasError() != expected.Error {
				t.Fatalf("added=%t error=%t, expected added=%t error=%t", added, diagnostics.HasError(), expected.Added, expected.Error)
			}
			if actual, ok := attributes["values"]; ok != expected.Added || (ok && actual != expected.Value) {
				t.Fatalf("payload value=%v present=%t, expected value=%q present=%t", actual, ok, expected.Value, expected.Added)
			}
		})
	}
}
