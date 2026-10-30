package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = stringNotNullValidator{}

// stringNotNullValidator validates that a string is configured where the
// surrounding schema cannot declare the attribute as required.
type stringNotNullValidator struct{}

func (v stringNotNullValidator) Description(_ context.Context) string {
	return "is required"
}

func (v stringNotNullValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v stringNotNullValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if !request.ConfigValue.IsNull() {
		return
	}

	path := request.Path.String()
	if strings.Contains(path, "[") {
		setName := path[:strings.Index(path, "[")]
		attributeName := path[strings.Index(path, "]")+2:]
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Incorrect attribute value type",
			fmt.Sprintf("Inappropriate value for attribute \"%s\": attribute \"%s\" is required.", setName, attributeName),
		)
		return
	}

	response.Diagnostics.AddAttributeError(
		request.Path,
		"Incorrect attribute value type",
		fmt.Sprintf("Attribute is required for path: %s", path),
	)
}

// StringNotNull returns a validator for semantically required strings that
// must be represented as optional and computed in the Terraform schema.
func StringNotNull() validator.String {
	return stringNotNullValidator{}
}
