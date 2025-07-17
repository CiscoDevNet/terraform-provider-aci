package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// VzRsAnyToConsPrio custom string type.

var _ basetypes.StringTypable = VzRsAnyToConsPrioStringType{}

type VzRsAnyToConsPrioStringType struct {
	basetypes.StringType
}

func (t VzRsAnyToConsPrioStringType) Equal(o attr.Type) bool {
	other, ok := o.(VzRsAnyToConsPrioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t VzRsAnyToConsPrioStringType) String() string {
	return "VzRsAnyToConsPrioStringType"
}

func (t VzRsAnyToConsPrioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := VzRsAnyToConsPrioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t VzRsAnyToConsPrioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t VzRsAnyToConsPrioStringType) ValueType(ctx context.Context) attr.Value {
	return VzRsAnyToConsPrioStringValue{}
}

// VzRsAnyToConsPrio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = VzRsAnyToConsPrioStringValue{}

type VzRsAnyToConsPrioStringValue struct {
	basetypes.StringValue
}

func (v VzRsAnyToConsPrioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(VzRsAnyToConsPrioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v VzRsAnyToConsPrioStringValue) Type(ctx context.Context) attr.Type {
	return VzRsAnyToConsPrioStringType{}
}

func (v VzRsAnyToConsPrioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(VzRsAnyToConsPrioStringValue)

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

	priorMappedValue := VzRsAnyToConsPrioValueMap(v.StringValue)

	newMappedValue := VzRsAnyToConsPrioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v VzRsAnyToConsPrioStringValue) NamedValueString() string {
	return VzRsAnyToConsPrioValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func VzRsAnyToConsPrioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewVzRsAnyToConsPrioStringNull() VzRsAnyToConsPrioStringValue {
	return VzRsAnyToConsPrioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewVzRsAnyToConsPrioStringUnknown() VzRsAnyToConsPrioStringValue {
	return VzRsAnyToConsPrioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewVzRsAnyToConsPrioStringValue(value string) VzRsAnyToConsPrioStringValue {
	return VzRsAnyToConsPrioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewVzRsAnyToConsPrioStringPointerValue(value *string) VzRsAnyToConsPrioStringValue {
	return VzRsAnyToConsPrioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
