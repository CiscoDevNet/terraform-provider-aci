package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// VzOOBBrCPprio custom string type.

var _ basetypes.StringTypable = VzOOBBrCPprioStringType{}

type VzOOBBrCPprioStringType struct {
	basetypes.StringType
}

func (t VzOOBBrCPprioStringType) Equal(o attr.Type) bool {
	other, ok := o.(VzOOBBrCPprioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t VzOOBBrCPprioStringType) String() string {
	return "VzOOBBrCPprioStringType"
}

func (t VzOOBBrCPprioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := VzOOBBrCPprioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t VzOOBBrCPprioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t VzOOBBrCPprioStringType) ValueType(ctx context.Context) attr.Value {
	return VzOOBBrCPprioStringValue{}
}

// VzOOBBrCPprio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = VzOOBBrCPprioStringValue{}

type VzOOBBrCPprioStringValue struct {
	basetypes.StringValue
}

func (v VzOOBBrCPprioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(VzOOBBrCPprioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v VzOOBBrCPprioStringValue) Type(ctx context.Context) attr.Type {
	return VzOOBBrCPprioStringType{}
}

func (v VzOOBBrCPprioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(VzOOBBrCPprioStringValue)

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

	priorMappedValue := VzOOBBrCPprioValueMap(v.StringValue)

	newMappedValue := VzOOBBrCPprioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func VzOOBBrCPprioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewVzOOBBrCPprioStringNull() VzOOBBrCPprioStringValue {
	return VzOOBBrCPprioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewVzOOBBrCPprioStringUnknown() VzOOBBrCPprioStringValue {
	return VzOOBBrCPprioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewVzOOBBrCPprioStringValue(value string) VzOOBBrCPprioStringValue {
	return VzOOBBrCPprioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewVzOOBBrCPprioStringPointerValue(value *string) VzOOBBrCPprioStringValue {
	return VzOOBBrCPprioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
