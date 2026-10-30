package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type stringNotNullInput struct {
	Path  path.Path
	Value types.String
}

type stringNotNullExpected struct {
	HasError bool
	Detail   string
}

func TestStringNotNull(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "test_configured_value",
			Input: stringNotNullInput{
				Path:  path.Root("name"),
				Value: types.StringValue("configured"),
			},
			Expected: stringNotNullExpected{},
		},
		{
			Name: "test_unknown_value",
			Input: stringNotNullInput{
				Path:  path.Root("name"),
				Value: types.StringUnknown(),
			},
			Expected: stringNotNullExpected{},
		},
		{
			Name: "test_null_root_value",
			Input: stringNotNullInput{
				Path:  path.Root("name"),
				Value: types.StringNull(),
			},
			Expected: stringNotNullExpected{
				HasError: true,
				Detail:   "Attribute is required for path: name",
			},
		},
		{
			Name: "test_null_set_element_value",
			Input: stringNotNullInput{
				Path:  path.Root("children").AtSetValue(types.StringValue("element")).AtName("name"),
				Value: types.StringNull(),
			},
			Expected: stringNotNullExpected{
				HasError: true,
				Detail:   "Inappropriate value for attribute \"children\": attribute \"name\" is required.",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(stringNotNullInput)
			expected := testCase.Expected.(stringNotNullExpected)
			request := validator.StringRequest{Path: input.Path, ConfigValue: input.Value}
			var response validator.StringResponse

			StringNotNull().ValidateString(context.Background(), request, &response)

			if response.Diagnostics.HasError() != expected.HasError {
				t.Fatalf("diagnostic error state %t, expected %t: %v", response.Diagnostics.HasError(), expected.HasError, response.Diagnostics)
			}
			if expected.Detail != "" && response.Diagnostics[0].Detail() != expected.Detail {
				t.Fatalf("diagnostic detail %q, expected %q", response.Diagnostics[0].Detail(), expected.Detail)
			}
		})
	}
}
