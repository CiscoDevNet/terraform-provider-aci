package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type stringValue interface {
	IsNull() bool
	IsUnknown() bool
	ValueString() string
}

func objectIsEmpty(value types.Object) bool {
	for _, attribute := range value.Attributes() {
		if !attribute.IsNull() {
			return false
		}
	}

	return true
}

func stringAttribute(
	diagnostics *diag.Diagnostics,
	attributes Attributes,
	propertyName string,
) (string, bool, bool) {
	attribute, found := attributes.values[propertyName]
	if !found {
		return "", false, false
	}

	rawValue, ok := attribute.Data().(string)
	if !ok {
		diagnostics.AddError(
			"Malformed APIC attribute",
			fmt.Sprintf(
				"APIC attribute %q on class %q must be a string, got %T.",
				propertyName,
				attributes.className,
				attribute.Data(),
			),
		)
		return "", true, false
	}

	return rawValue, true, true
}

func decodeStringAttribute[T any](
	diagnostics *diag.Diagnostics,
	attributes Attributes,
	propertyName string,
	emptyValue *string,
	target *T,
	constructor func(string) T,
) bool {
	rawValue, found, valid := stringAttribute(diagnostics, attributes, propertyName)
	if !found || !valid {
		return found
	}
	if emptyValue != nil && rawValue == "" {
		rawValue = *emptyValue
	}

	*target = constructor(rawValue)
	return true
}

// DecodeStringAttribute converts one APIC string attribute to a Terraform string value.
func DecodeStringAttribute(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	attributes Attributes,
	propertyName string,
	target *types.String,
) bool {
	return decodeStringAttribute(diagnostics, attributes, propertyName, nil, target, types.StringValue)
}

// DecodeStringAttributeWithEmptyValue substitutes empty APIC strings before conversion.
func DecodeStringAttributeWithEmptyValue(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	attributes Attributes,
	propertyName string,
	emptyValue string,
	target *types.String,
) bool {
	return decodeStringAttribute(diagnostics, attributes, propertyName, &emptyValue, target, types.StringValue)
}

// DecodeCustomStringAttribute converts one APIC string attribute with a concrete Terraform value constructor.
func DecodeCustomStringAttribute[T any](
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	attributes Attributes,
	propertyName string,
	target *T,
	constructor func(string) T,
) bool {
	return decodeStringAttribute(diagnostics, attributes, propertyName, nil, target, constructor)
}

// DecodeStringSetAttribute converts one comma-separated APIC attribute to a Terraform string set.
func DecodeStringSetAttribute(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	attributes Attributes,
	propertyName string,
	target *types.Set,
) bool {
	rawValue, found, valid := stringAttribute(diagnostics, attributes, propertyName)
	if !found || !valid {
		return found
	}

	values := []string{}
	if rawValue != "" {
		values = strings.Split(rawValue, ",")
	}

	value, valueDiagnostics := types.SetValueFrom(ctx, types.StringType, values)
	diagnostics.Append(valueDiagnostics...)
	if !valueDiagnostics.HasError() {
		*target = value
	}
	return true
}

// AddPayloadStringAttribute adds a known Terraform string value to APIC payload attributes.
func AddPayloadStringAttribute(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	attributes map[string]any,
	propertyName string,
	value stringValue,
) bool {
	if value.IsNull() || value.IsUnknown() {
		return false
	}

	attributes[propertyName] = value.ValueString()
	return true
}

// AddPayloadStringSetAttribute adds a known Terraform string set as a comma-separated APIC attribute.
func AddPayloadStringSetAttribute(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
	attributes map[string]any,
	propertyName string,
	value types.Set,
) bool {
	if value.IsNull() || value.IsUnknown() {
		return false
	}

	var values []string
	valueDiagnostics := value.ElementsAs(ctx, &values, false)
	diagnostics.Append(valueDiagnostics...)
	if valueDiagnostics.HasError() {
		return false
	}

	attributes[propertyName] = strings.Join(values, ",")
	return true
}
