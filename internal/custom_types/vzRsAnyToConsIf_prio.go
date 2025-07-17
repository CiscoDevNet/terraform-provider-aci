package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// VzRsAnyToConsIfPrio custom string type.

var _ basetypes.StringTypable = VzRsAnyToConsIfPrioStringType{}

type VzRsAnyToConsIfPrioStringType struct {
	basetypes.StringType
}

func (t VzRsAnyToConsIfPrioStringType) Equal(o attr.Type) bool {
	other, ok := o.(VzRsAnyToConsIfPrioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t VzRsAnyToConsIfPrioStringType) String() string {
	return "VzRsAnyToConsIfPrioStringType"
}

func (t VzRsAnyToConsIfPrioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := VzRsAnyToConsIfPrioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t VzRsAnyToConsIfPrioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t VzRsAnyToConsIfPrioStringType) ValueType(ctx context.Context) attr.Value {
	return VzRsAnyToConsIfPrioStringValue{}
}

// VzRsAnyToConsIfPrio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = VzRsAnyToConsIfPrioStringValue{}

type VzRsAnyToConsIfPrioStringValue struct {
	basetypes.StringValue
}

func (v VzRsAnyToConsIfPrioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(VzRsAnyToConsIfPrioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v VzRsAnyToConsIfPrioStringValue) Type(ctx context.Context) attr.Type {
	return VzRsAnyToConsIfPrioStringType{}
}

func (v VzRsAnyToConsIfPrioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(VzRsAnyToConsIfPrioStringValue)

	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	priorMappedValue := VzRsAnyToConsIfPrioValueMap(v.StringValue)

	newMappedValue := VzRsAnyToConsIfPrioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v VzRsAnyToConsIfPrioStringValue) NamedValueString() string {
	return VzRsAnyToConsIfPrioValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func VzRsAnyToConsIfPrioValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0": "unspecified",
		"1": "level3",
		"2": "level2",
		"3": "level1",
		"7": "level6",
		"8": "level5",
		"9": "level4",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewVzRsAnyToConsIfPrioStringNull() VzRsAnyToConsIfPrioStringValue {
	return VzRsAnyToConsIfPrioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewVzRsAnyToConsIfPrioStringUnknown() VzRsAnyToConsIfPrioStringValue {
	return VzRsAnyToConsIfPrioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewVzRsAnyToConsIfPrioStringValue(value string) VzRsAnyToConsIfPrioStringValue {
	return VzRsAnyToConsIfPrioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewVzRsAnyToConsIfPrioStringPointerValue(value *string) VzRsAnyToConsIfPrioStringValue {
	return VzRsAnyToConsIfPrioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
