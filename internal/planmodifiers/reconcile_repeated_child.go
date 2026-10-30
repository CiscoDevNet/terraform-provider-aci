package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ planmodifier.Set = reconcileRepeatedChildModifier{}

type reconcileRepeatedChildModifier struct {
	schema ReconciliationSchema
}

func (m reconcileRepeatedChildModifier) Description(_ context.Context) string {
	return "Reconciles repeated children with prior state by all of their semantic identifiers instead of by set position."
}

func (m reconcileRepeatedChildModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m reconcileRepeatedChildModifier) PlanModifySet(ctx context.Context, request planmodifier.SetRequest, response *planmodifier.SetResponse) {
	planValue, diagnostics := m.schema.reconcileRepeatedValue(
		ctx,
		request.ConfigValue,
		request.PlanValue,
		request.StateValue,
		!request.State.Raw.IsNull(),
	)
	response.Diagnostics.Append(diagnostics...)
	if !diagnostics.HasError() {
		response.PlanValue = planValue
	}
}

// ReconcileRepeatedChild reconciles an unordered set of child objects and
// recursively owns state preservation for every matched child's subtree.
func ReconcileRepeatedChild(schema ReconciliationSchema) planmodifier.Set {
	return reconcileRepeatedChildModifier{schema: schema}
}
