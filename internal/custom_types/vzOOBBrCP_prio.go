package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// VzOOBBrCPPrio custom string type.

var _ basetypes.StringTypable = VzOOBBrCPPrioStringType{}

type VzOOBBrCPPrioStringType struct {
	basetypes.StringType
}

func (t VzOOBBrCPPrioStringType) Equal(o attr.Type) bool {
	other, ok := o.(VzOOBBrCPPrioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t VzOOBBrCPPrioStringType) String() string {
	return "VzOOBBrCPPrioStringType"
}

func (t VzOOBBrCPPrioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := VzOOBBrCPPrioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t VzOOBBrCPPrioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t VzOOBBrCPPrioStringType) ValueType(ctx context.Context) attr.Value {
	return VzOOBBrCPPrioStringValue{}
}

// VzOOBBrCPPrio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = VzOOBBrCPPrioStringValue{}

type VzOOBBrCPPrioStringValue struct {
	basetypes.StringValue
}

func (v VzOOBBrCPPrioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(VzOOBBrCPPrioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v VzOOBBrCPPrioStringValue) Type(ctx context.Context) attr.Type {
	return VzOOBBrCPPrioStringType{}
}

func (v VzOOBBrCPPrioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(VzOOBBrCPPrioStringValue)

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

	priorMappedValue := VzOOBBrCPPrioValueMap(v.StringValue)

	newMappedValue := VzOOBBrCPPrioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func VzOOBBrCPPrioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewVzOOBBrCPPrioStringNull() VzOOBBrCPPrioStringValue {
	return VzOOBBrCPPrioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewVzOOBBrCPPrioStringUnknown() VzOOBBrCPPrioStringValue {
	return VzOOBBrCPPrioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewVzOOBBrCPPrioStringValue(value string) VzOOBBrCPPrioStringValue {
	return VzOOBBrCPPrioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewVzOOBBrCPPrioStringPointerValue(value *string) VzOOBBrCPPrioStringValue {
	return VzOOBBrCPPrioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
