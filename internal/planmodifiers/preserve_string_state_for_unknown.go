package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ planmodifier.String = preserveStringStateForUnknownModifier{}

type preserveStringStateForUnknownModifier struct{}

func (m preserveStringStateForUnknownModifier) Description(_ context.Context) string {
	return "Preserves the prior string value when the planned value is unknown, including null returned by an APIC version that does not support the attribute."
}

func (m preserveStringStateForUnknownModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m preserveStringStateForUnknownModifier) PlanModifyString(_ context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
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

// PreserveStringStateForUnknown preserves the established update behavior for
// optional and computed string attributes across different APIC versions.
func PreserveStringStateForUnknown() planmodifier.String {
	return preserveStringStateForUnknownModifier{}
}
