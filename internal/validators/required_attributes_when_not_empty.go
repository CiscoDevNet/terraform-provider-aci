package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Object = requiredAttributesWhenNotEmptyValidator{}

// requiredAttributesWhenNotEmptyValidator permits an empty object to remove a
// child while requiring its identifying attributes when the object is set.
type requiredAttributesWhenNotEmptyValidator struct {
	attributeName      string
	requiredAttributes []string
}

func (v requiredAttributesWhenNotEmptyValidator) Description(_ context.Context) string {
	return "is required"
}

func (v requiredAttributesWhenNotEmptyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v requiredAttributesWhenNotEmptyValidator) ValidateObject(_ context.Context, request validator.ObjectRequest, response *validator.ObjectResponse) {
	if request.ConfigValue.IsNull() {
		return
	}

	attributes := request.ConfigValue.Attributes()
	if IsEmptySingleNestedAttribute(attributes) || !requiredAttributesMissing(attributes, v.requiredAttributes) {
		return
	}

	errorMessage := fmt.Sprintf(
		"Inappropriate value for attribute \"%s\": attribute \"%s\" is required.",
		v.attributeName,
		strings.Join(v.requiredAttributes, ", "),
	)
	if len(v.requiredAttributes) > 1 {
		errorMessage = fmt.Sprintf(
			"Inappropriate values for attribute \"%s\": attributes \"%s\" are required.",
			v.attributeName,
			strings.Join(v.requiredAttributes, "\", \""),
		)
	}

	response.Diagnostics.AddError("Incorrect attribute value type", errorMessage)
}

// IsEmptySingleNestedAttribute reports whether every attribute in a nested
// object is null.
func IsEmptySingleNestedAttribute(attributes map[string]attr.Value) bool {
	for _, value := range attributes {
		if !value.IsNull() {
			return false
		}
	}
	return true
}

func requiredAttributesMissing(attributes map[string]attr.Value, requiredAttributes []string) bool {
	for _, requiredAttribute := range requiredAttributes {
		value, ok := attributes[requiredAttribute]
		if !ok || value.IsNull() {
			return true
		}
	}
	return false
}

// RequiredAttributesWhenNotEmpty returns a validator that permits an empty
// object but rejects partially configured objects that omit required values.
func RequiredAttributesWhenNotEmpty(attributeName string, requiredAttributes []string) validator.Object {
	return requiredAttributesWhenNotEmptyValidator{
		attributeName:      attributeName,
		requiredAttributes: requiredAttributes,
	}
}
