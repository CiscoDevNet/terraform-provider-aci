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

type repeatedChildModifierInput struct {
	RawState    tftypes.Value
	ConfigValue types.Set
	PlanValue   types.Set
	StateValue  types.Set
}

func TestReconcileRepeatedChild(t *testing.T) {
	attributeTypes := map[string]attr.Type{"id": types.StringType}
	objectType := types.ObjectType{AttrTypes: attributeTypes}
	stateValue := types.SetValueMust(objectType, []attr.Value{
		types.ObjectValueMust(attributeTypes, map[string]attr.Value{"id": types.StringValue("one")}),
	})
	unknownPlan := types.SetUnknown(objectType)
	testCases := []test.TestCase{
		{
			Name: "test_existing_resource",
			Input: repeatedChildModifierInput{
				RawState:    tftypes.NewValue(tftypes.String, "state"),
				ConfigValue: types.SetNull(objectType),
				PlanValue:   unknownPlan,
				StateValue:  stateValue,
			},
			Expected: stateValue,
		},
		{
			Name: "test_new_resource",
			Input: repeatedChildModifierInput{
				RawState:    tftypes.NewValue(tftypes.String, nil),
				ConfigValue: types.SetNull(objectType),
				PlanValue:   unknownPlan,
				StateValue:  types.SetNull(objectType),
			},
			Expected: unknownPlan,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(repeatedChildModifierInput)
			request := planmodifier.SetRequest{
				State:       tfsdk.State{Raw: input.RawState},
				ConfigValue: input.ConfigValue,
				PlanValue:   input.PlanValue,
				StateValue:  input.StateValue,
			}
			response := planmodifier.SetResponse{PlanValue: input.PlanValue}

			ReconcileRepeatedChild(ReconciliationSchema{Identifiers: []string{"id"}}).PlanModifySet(context.Background(), request, &response)

			if response.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", response.Diagnostics)
			}
			if !response.PlanValue.Equal(testCase.Expected.(types.Set)) {
				t.Fatalf("plan value %#v, expected %#v", response.PlanValue, testCase.Expected)
			}
		})
	}
}
