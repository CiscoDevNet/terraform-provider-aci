package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ planmodifier.Object = reconcileSingletonChildModifier{}

type reconcileSingletonChildModifier struct {
	schema ReconciliationSchema
}

func (m reconcileSingletonChildModifier) Description(_ context.Context) string {
	return "Reconciles a singleton child with its prior state while preserving an explicitly empty object as a deletion request."
}

func (m reconcileSingletonChildModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m reconcileSingletonChildModifier) PlanModifyObject(ctx context.Context, request planmodifier.ObjectRequest, response *planmodifier.ObjectResponse) {
	planValue, diagnostics := m.schema.reconcileSingletonValue(
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

// ReconcileSingletonChild reconciles one optional or required child object and
// recursively owns state preservation for its complete subtree.
func ReconcileSingletonChild(schema ReconciliationSchema) planmodifier.Object {
	return reconcileSingletonChildModifier{schema: schema}
}
