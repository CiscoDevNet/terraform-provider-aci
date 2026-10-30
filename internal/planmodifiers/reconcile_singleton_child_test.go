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

type singletonChildModifierInput struct {
	RawState    tftypes.Value
	ConfigValue types.Object
	PlanValue   types.Object
	StateValue  types.Object
}

func TestReconcileSingletonChild(t *testing.T) {
	attributeTypes := map[string]attr.Type{"value": types.StringType}
	stateValue := types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"value": types.StringValue("state"),
	})
	unknownPlan := types.ObjectUnknown(attributeTypes)
	testCases := []test.TestCase{
		{
			Name: "test_existing_resource",
			Input: singletonChildModifierInput{
				RawState:    tftypes.NewValue(tftypes.String, "state"),
				ConfigValue: types.ObjectNull(attributeTypes),
				PlanValue:   unknownPlan,
				StateValue:  stateValue,
			},
			Expected: stateValue,
		},
		{
			Name: "test_new_resource",
			Input: singletonChildModifierInput{
				RawState:    tftypes.NewValue(tftypes.String, nil),
				ConfigValue: types.ObjectNull(attributeTypes),
				PlanValue:   unknownPlan,
				StateValue:  types.ObjectNull(attributeTypes),
			},
			Expected: unknownPlan,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(singletonChildModifierInput)
			request := planmodifier.ObjectRequest{
				State:       tfsdk.State{Raw: input.RawState},
				ConfigValue: input.ConfigValue,
				PlanValue:   input.PlanValue,
				StateValue:  input.StateValue,
			}
			response := planmodifier.ObjectResponse{PlanValue: input.PlanValue}

			ReconcileSingletonChild(ReconciliationSchema{}).PlanModifyObject(context.Background(), request, &response)

			if response.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", response.Diagnostics)
			}
			if !response.PlanValue.Equal(testCase.Expected.(types.Object)) {
				t.Fatalf("plan value %#v, expected %#v", response.PlanValue, testCase.Expected)
			}
		})
	}
}
