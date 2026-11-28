package customtypes

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type semanticEqualityInput struct {
	Prior string
	New   string
}

func TestSemanticStringValueSemanticEquality(t *testing.T) {
	semanticType := NewSemanticStringType(map[string]string{
		"1": "level3",
		"3": "level1",
	})
	testCases := []test.TestCase{
		{Name: "test_wire_and_local_name", Input: semanticEqualityInput{Prior: "3", New: "level1"}, Expected: true},
		{Name: "test_matching_wire_values", Input: semanticEqualityInput{Prior: "1", New: "1"}, Expected: true},
		{Name: "test_unmapped_matching_values", Input: semanticEqualityInput{Prior: "other", New: "other"}, Expected: true},
		{Name: "test_different_values", Input: semanticEqualityInput{Prior: "1", New: "3"}, Expected: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(semanticEqualityInput)
			equal, diagnostics := semanticType.Value(input.Prior).StringSemanticEquals(
				context.Background(),
				semanticType.Value(input.New),
			)
			if diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diagnostics)
			}
			if equal != testCase.Expected.(bool) {
				t.Fatalf("semantic equality returned %t, expected %t", equal, testCase.Expected)
			}
		})
	}
}

func TestSemanticStringTypeRetainsDefinition(t *testing.T) {
	valueMap := map[string]string{"3": "level1"}
	semanticType := NewSemanticStringType(valueMap)
	valueMap["3"] = "changed"

	value := semanticType.Value("3")
	if actual := value.CanonicalString(); actual != "level1" {
		t.Fatalf("canonical value %q, expected level1", actual)
	}
	if !value.Type(context.Background()).Equal(NewSemanticStringType(map[string]string{"3": "level1"})) {
		t.Fatal("equivalent semantic type definitions did not compare equally")
	}
	if value.Type(context.Background()).Equal(NewSemanticStringType(map[string]string{"3": "other"})) {
		t.Fatal("different semantic type definitions compared equally")
	}
	if value.Equal(semanticType.Value("level1")) {
		t.Fatal("exact value equality used semantic normalization")
	}
}

func TestSemanticStringTypeValueStates(t *testing.T) {
	semanticType := NewSemanticStringType(map[string]string{"3": "level1"})
	if !semanticType.Null().IsNull() {
		t.Fatal("Null did not return a null value")
	}
	if !semanticType.Unknown().IsUnknown() {
		t.Fatal("Unknown did not return an unknown value")
	}
	if actual := semanticType.Value("3").ValueString(); actual != "3" {
		t.Fatalf("raw value %q, expected 3", actual)
	}
}

func TestSemanticStringTypeTerraformConversion(t *testing.T) {
	semanticType := NewSemanticStringType(map[string]string{"3": "level1"})
	attributeValue, err := semanticType.ValueFromTerraform(
		context.Background(),
		tftypes.NewValue(tftypes.String, "3"),
	)
	if err != nil {
		t.Fatalf("convert Terraform value: %v", err)
	}

	value, ok := attributeValue.(SemanticStringValue)
	if !ok {
		t.Fatalf("converted value has type %T", attributeValue)
	}
	if actual := value.CanonicalString(); actual != "level1" {
		t.Fatalf("canonical value %q, expected level1", actual)
	}
}

func TestSemanticStringValueRejectsDifferentTypes(t *testing.T) {
	semanticType := NewSemanticStringType(map[string]string{"3": "level1"})
	_, diagnostics := semanticType.Value("3").StringSemanticEquals(
		context.Background(),
		basetypes.NewStringValue("level1"),
	)
	if !diagnostics.HasError() {
		t.Fatal("different value type did not produce an error diagnostic")
	}

	differentType := NewSemanticStringType(map[string]string{"3": "other"})
	_, diagnostics = semanticType.Value("3").StringSemanticEquals(
		context.Background(),
		differentType.Value("3"),
	)
	if !diagnostics.HasError() {
		t.Fatal("different semantic definition did not produce an error diagnostic")
	}
}

func TestSemanticStringTypeImplementsAttributeType(t *testing.T) {
	var attributeType attr.Type = NewSemanticStringType(map[string]string{"3": "level1"})
	if attributeType.String() != "customtypes.SemanticStringType" {
		t.Fatalf("type string %q", attributeType.String())
	}
}

func TestSemanticStringValueModelConversion(t *testing.T) {
	type model struct {
		Priority SemanticStringValue `tfsdk:"priority"`
	}

	ctx := context.Background()
	semanticType := NewSemanticStringType(map[string]string{"3": "level1"})
	attributeTypes := map[string]attr.Type{"priority": semanticType}
	value, diagnostics := types.ObjectValueFrom(ctx, attributeTypes, model{Priority: semanticType.Value("3")})
	if diagnostics.HasError() {
		t.Fatalf("convert model to object: %v", diagnostics)
	}

	var decoded model
	diagnostics = value.As(ctx, &decoded, basetypes.ObjectAsOptions{})
	if diagnostics.HasError() {
		t.Fatalf("convert object to model: %v", diagnostics)
	}
	if actual := decoded.Priority.CanonicalString(); actual != "level1" {
		t.Fatalf("canonical decoded value %q, expected level1", actual)
	}
}
