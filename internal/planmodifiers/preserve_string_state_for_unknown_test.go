package planmodifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type preserveStringInput struct {
	RawState    tftypes.Value
	ConfigValue types.String
	StateValue  types.String
	PlanValue   types.String
}

func TestPreserveStringStateForUnknown(t *testing.T) {
	nonNullRawState := tftypes.NewValue(tftypes.String, "state")
	nullRawState := tftypes.NewValue(tftypes.String, nil)
	testCases := []test.TestCase{
		{
			Name: "test_update_null_state_unknown_plan",
			Input: preserveStringInput{
				RawState:    nonNullRawState,
				ConfigValue: types.StringNull(),
				StateValue:  types.StringNull(),
				PlanValue:   types.StringUnknown(),
			},
			Expected: types.StringNull(),
		},
		{
			Name: "test_update_null_state_unknown_configuration",
			Input: preserveStringInput{
				RawState:    nonNullRawState,
				ConfigValue: types.StringUnknown(),
				StateValue:  types.StringNull(),
				PlanValue:   types.StringUnknown(),
			},
			Expected: types.StringNull(),
		},
		{
			Name: "test_update_known_state_omitted_configuration",
			Input: preserveStringInput{
				RawState:    nonNullRawState,
				ConfigValue: types.StringNull(),
				StateValue:  types.StringValue("state"),
				PlanValue:   types.StringUnknown(),
			},
			Expected: types.StringValue("state"),
		},
		{
			Name: "test_update_known_state_unknown_configuration",
			Input: preserveStringInput{
				RawState:    nonNullRawState,
				ConfigValue: types.StringUnknown(),
				StateValue:  types.StringValue("state"),
				PlanValue:   types.StringUnknown(),
			},
			Expected: types.StringUnknown(),
		},
		{
			Name: "test_create_null_state_unknown_plan",
			Input: preserveStringInput{
				RawState:    nullRawState,
				ConfigValue: types.StringNull(),
				StateValue:  types.StringNull(),
				PlanValue:   types.StringUnknown(),
			},
			Expected: types.StringUnknown(),
		},
		{
			Name: "test_update_known_plan",
			Input: preserveStringInput{
				RawState:    nonNullRawState,
				ConfigValue: types.StringValue("configured"),
				StateValue:  types.StringNull(),
				PlanValue:   types.StringValue("configured"),
			},
			Expected: types.StringValue("configured"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(preserveStringInput)
			request := planmodifier.StringRequest{
				State:       tfsdk.State{Raw: input.RawState},
				ConfigValue: input.ConfigValue,
				StateValue:  input.StateValue,
				PlanValue:   input.PlanValue,
			}
			response := planmodifier.StringResponse{PlanValue: input.PlanValue}

			PreserveStringStateForUnknown().PlanModifyString(context.Background(), request, &response)

			if !response.PlanValue.Equal(testCase.Expected.(types.String)) {
				t.Fatalf("plan value %#v, expected %#v", response.PlanValue, testCase.Expected)
			}
		})
	}
}
