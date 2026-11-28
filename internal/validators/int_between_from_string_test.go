package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type inBetweenFromStringInput struct {
	Min   float64
	Max   float64
	Value types.String
}

type inBetweenFromStringExpected struct {
	NilValidator bool
	HasError     bool
}

func TestInBetweenFromString(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "test_value_in_range",
			Input: inBetweenFromStringInput{
				Min:   1,
				Max:   10,
				Value: types.StringValue("5"),
			},
			Expected: inBetweenFromStringExpected{},
		},
		{
			Name: "test_value_outside_range",
			Input: inBetweenFromStringInput{
				Min:   1,
				Max:   10,
				Value: types.StringValue("11"),
			},
			Expected: inBetweenFromStringExpected{HasError: true},
		},
		{
			Name: "test_non_numeric_value",
			Input: inBetweenFromStringInput{
				Min:   1,
				Max:   10,
				Value: types.StringValue("invalid"),
			},
			Expected: inBetweenFromStringExpected{HasError: true},
		},
		{
			Name: "test_null_value",
			Input: inBetweenFromStringInput{
				Min:   1,
				Max:   10,
				Value: types.StringNull(),
			},
			Expected: inBetweenFromStringExpected{},
		},
		{
			Name: "test_unknown_value",
			Input: inBetweenFromStringInput{
				Min:   1,
				Max:   10,
				Value: types.StringUnknown(),
			},
			Expected: inBetweenFromStringExpected{},
		},
		{
			Name: "test_invalid_range",
			Input: inBetweenFromStringInput{
				Min:   10,
				Max:   1,
				Value: types.StringValue("5"),
			},
			Expected: inBetweenFromStringExpected{NilValidator: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(inBetweenFromStringInput)
			expected := testCase.Expected.(inBetweenFromStringExpected)
			stringValidator := InBetweenFromString(input.Min, input.Max)
			if (stringValidator == nil) != expected.NilValidator {
				t.Fatalf("nil validator state %t, expected %t", stringValidator == nil, expected.NilValidator)
			}
			if stringValidator == nil {
				return
			}

			request := validator.StringRequest{Path: path.Root("value"), ConfigValue: input.Value}
			var response validator.StringResponse
			stringValidator.ValidateString(context.Background(), request, &response)
			if response.Diagnostics.HasError() != expected.HasError {
				t.Fatalf("diagnostic error state %t, expected %t: %v", response.Diagnostics.HasError(), expected.HasError, response.Diagnostics)
			}
		})
	}
}
