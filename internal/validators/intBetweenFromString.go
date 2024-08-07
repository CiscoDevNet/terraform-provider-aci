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
	min, max int
}

func (v InBetweenFromStringValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v InBetweenFromStringValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value must be between %d and %d", v.min, v.max)
}

func (v InBetweenFromStringValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	configValueInt, err := strconv.Atoi(request.ConfigValue.ValueString())
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			err.Error(),
			request.ConfigValue.ValueString(),
		))
	} else if configValueInt < int(v.min) || configValueInt > int(v.max) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("%d", configValueInt),
		))
	}
}

func InBetweenFromString(min, max int) validator.String {
	if min > max {
		return nil
	}

	return InBetweenFromStringValidator{
		min: min,
		max: max,
	}
}
