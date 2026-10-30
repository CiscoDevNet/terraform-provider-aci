package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

func TestRequiredAttributesWhenNotEmpty(t *testing.T) {
	attributeTypes := map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
	}
	testCases := []test.TestCase{
		{
			Name:     "test_null_object",
			Input:    types.ObjectNull(attributeTypes),
			Expected: false,
		},
		{
			Name:     "test_unknown_object",
			Input:    types.ObjectUnknown(attributeTypes),
			Expected: false,
		},
		{
			Name: "test_empty_object_for_deletion",
			Input: types.ObjectValueMust(attributeTypes, map[string]attr.Value{
				"name":        types.StringNull(),
				"description": types.StringNull(),
			}),
			Expected: false,
		},
		{
			Name: "test_partial_object",
			Input: types.ObjectValueMust(attributeTypes, map[string]attr.Value{
				"name":        types.StringNull(),
				"description": types.StringValue("configured"),
			}),
			Expected: true,
		},
		{
			Name: "test_complete_object",
			Input: types.ObjectValueMust(attributeTypes, map[string]attr.Value{
				"name":        types.StringValue("configured"),
				"description": types.StringNull(),
			}),
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			request := validator.ObjectRequest{ConfigValue: testCase.Input.(types.Object)}
			var response validator.ObjectResponse

			RequiredAttributesWhenNotEmpty("child", []string{"name"}).ValidateObject(context.Background(), request, &response)

			if response.Diagnostics.HasError() != testCase.Expected.(bool) {
				t.Fatalf("diagnostic error state %t, expected %t: %v", response.Diagnostics.HasError(), testCase.Expected.(bool), response.Diagnostics)
			}
		})
	}
}

func TestRequiredAttributesWhenNotEmptyDiagnostic(t *testing.T) {
	attributeTypes := map[string]attr.Type{
		"name":   types.StringType,
		"target": types.StringType,
		"value":  types.StringType,
	}
	request := validator.ObjectRequest{
		ConfigValue: types.ObjectValueMust(attributeTypes, map[string]attr.Value{
			"name":   types.StringNull(),
			"target": types.StringNull(),
			"value":  types.StringValue("configured"),
		}),
	}
	var response validator.ObjectResponse

	RequiredAttributesWhenNotEmpty("child", []string{"name", "target"}).ValidateObject(context.Background(), request, &response)

	if !response.Diagnostics.HasError() {
		t.Fatal("expected missing required attributes to produce a diagnostic")
	}
	expected := "Inappropriate values for attribute \"child\": attributes \"name\", \"target\" are required."
	if response.Diagnostics[0].Detail() != expected {
		t.Fatalf("diagnostic detail %q, expected %q", response.Diagnostics[0].Detail(), expected)
	}
}
