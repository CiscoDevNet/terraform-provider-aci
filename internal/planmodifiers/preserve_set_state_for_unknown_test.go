package planmodifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type preserveSetInput struct {
	RawState    tftypes.Value
	ConfigValue types.Set
	StateValue  types.Set
	PlanValue   types.Set
}

func TestPreserveSetStateForUnknown(t *testing.T) {
	nonNullRawState := tftypes.NewValue(tftypes.String, "state")
	nullRawState := tftypes.NewValue(tftypes.String, nil)
	stateSet := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("state")})
	testCases := []test.TestCase{
		{
			Name: "test_update_null_state_unknown_plan",
			Input: preserveSetInput{
				RawState:    nonNullRawState,
				ConfigValue: types.SetNull(types.StringType),
				StateValue:  types.SetNull(types.StringType),
				PlanValue:   types.SetUnknown(types.StringType),
			},
			Expected: types.SetNull(types.StringType),
		},
		{
			Name: "test_update_null_state_unknown_configuration",
			Input: preserveSetInput{
				RawState:    nonNullRawState,
				ConfigValue: types.SetUnknown(types.StringType),
				StateValue:  types.SetNull(types.StringType),
				PlanValue:   types.SetUnknown(types.StringType),
			},
			Expected: types.SetNull(types.StringType),
		},
		{
			Name: "test_update_known_state_omitted_configuration",
			Input: preserveSetInput{
				RawState:    nonNullRawState,
				ConfigValue: types.SetNull(types.StringType),
				StateValue:  stateSet,
				PlanValue:   types.SetUnknown(types.StringType),
			},
			Expected: stateSet,
		},
		{
			Name: "test_update_known_state_unknown_configuration",
			Input: preserveSetInput{
				RawState:    nonNullRawState,
				ConfigValue: types.SetUnknown(types.StringType),
				StateValue:  stateSet,
				PlanValue:   types.SetUnknown(types.StringType),
			},
			Expected: types.SetUnknown(types.StringType),
		},
		{
			Name: "test_create_null_state_unknown_plan",
			Input: preserveSetInput{
				RawState:    nullRawState,
				ConfigValue: types.SetNull(types.StringType),
				StateValue:  types.SetNull(types.StringType),
				PlanValue:   types.SetUnknown(types.StringType),
			},
			Expected: types.SetUnknown(types.StringType),
		},
		{
			Name: "test_update_known_plan",
			Input: preserveSetInput{
				RawState:    nonNullRawState,
				ConfigValue: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("configured")}),
				StateValue:  types.SetNull(types.StringType),
				PlanValue:   types.SetValueMust(types.StringType, []attr.Value{types.StringValue("configured")}),
			},
			Expected: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("configured")}),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(preserveSetInput)
			request := planmodifier.SetRequest{
				State:       tfsdk.State{Raw: input.RawState},
				ConfigValue: input.ConfigValue,
				StateValue:  input.StateValue,
				PlanValue:   input.PlanValue,
			}
			response := planmodifier.SetResponse{PlanValue: input.PlanValue}

			PreserveSetStateForUnknown().PlanModifySet(context.Background(), request, &response)

			if !response.PlanValue.Equal(testCase.Expected.(types.Set)) {
				t.Fatalf("plan value %#v, expected %#v", response.PlanValue, testCase.Expected)
			}
		})
	}
}
