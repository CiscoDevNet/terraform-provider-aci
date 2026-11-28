package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ planmodifier.Set = preserveSetStateForUnknownModifier{}

type preserveSetStateForUnknownModifier struct{}

func (m preserveSetStateForUnknownModifier) Description(_ context.Context) string {
	return "Preserves the prior set value when the planned value is unknown, including null returned by an APIC version that does not support the attribute."
}

func (m preserveSetStateForUnknownModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m preserveSetStateForUnknownModifier) PlanModifySet(_ context.Context, request planmodifier.SetRequest, response *planmodifier.SetResponse) {
	if request.State.Raw.IsNull() || !request.PlanValue.IsUnknown() {
		return
	}

	// A null state value indicates that the connected APIC version does not
	// expose this attribute. Preserve it even when configuration is unknown.
	if request.StateValue.IsNull() {
		response.PlanValue = request.StateValue
		return
	}

	if request.ConfigValue.IsUnknown() {
		return
	}

	response.PlanValue = request.StateValue
}

// PreserveSetStateForUnknown preserves the established update behavior for
// optional and computed set attributes across different APIC versions.
func PreserveSetStateForUnknown() planmodifier.Set {
	return preserveSetStateForUnknownModifier{}
}
