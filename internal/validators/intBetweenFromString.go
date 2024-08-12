package validators

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.String = InBetweenFromStringValidator{}

type InBetweenFromStringValidator struct {
	min, max float64
}

func (v InBetweenFromStringValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v InBetweenFromStringValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be between %.2f and %.2f", v.min, v.max)
}

func (v InBetweenFromStringValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	configValueFloat, err := strconv.ParseFloat(request.ConfigValue.ValueString(), 64)
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			err.Error(),
			request.ConfigValue.ValueString(),
		))
	} else if configValueFloat < float64(v.min) || configValueFloat > float64(v.max) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%.2f", configValueFloat),
		))
	}
}

func InBetweenFromString(min, max float64) validator.String {
	if min > max {
		return nil
	}

	return InBetweenFromStringValidator{
		min: min,
		max: max,
	}
}
