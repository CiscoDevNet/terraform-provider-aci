package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// VzOOBBrCPTargetDscp custom string type.

var _ basetypes.StringTypable = VzOOBBrCPTargetDscpStringType{}

type VzOOBBrCPTargetDscpStringType struct {
	basetypes.StringType
}

func (t VzOOBBrCPTargetDscpStringType) Equal(o attr.Type) bool {
	other, ok := o.(VzOOBBrCPTargetDscpStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t VzOOBBrCPTargetDscpStringType) String() string {
	return "VzOOBBrCPTargetDscpStringType"
}

func (t VzOOBBrCPTargetDscpStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := VzOOBBrCPTargetDscpStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t VzOOBBrCPTargetDscpStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t VzOOBBrCPTargetDscpStringType) ValueType(ctx context.Context) attr.Value {
	return VzOOBBrCPTargetDscpStringValue{}
}

// VzOOBBrCPTargetDscp custom string value.

var _ basetypes.StringValuableWithSemanticEquals = VzOOBBrCPTargetDscpStringValue{}

type VzOOBBrCPTargetDscpStringValue struct {
	basetypes.StringValue
}

func (v VzOOBBrCPTargetDscpStringValue) Equal(o attr.Value) bool {
	other, ok := o.(VzOOBBrCPTargetDscpStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v VzOOBBrCPTargetDscpStringValue) Type(ctx context.Context) attr.Type {
	return VzOOBBrCPTargetDscpStringType{}
}

func (v VzOOBBrCPTargetDscpStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(VzOOBBrCPTargetDscpStringValue)

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

	priorMappedValue := VzOOBBrCPTargetDscpValueMap(v.StringValue)

	newMappedValue := VzOOBBrCPTargetDscpValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v VzOOBBrCPTargetDscpStringValue) NamedValueString() string {
	return VzOOBBrCPTargetDscpValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func VzOOBBrCPTargetDscpValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0":  "CS0",
		"10": "AF11",
		"12": "AF12",
		"14": "AF13",
		"16": "CS2",
		"18": "AF21",
		"20": "AF22",
		"22": "AF23",
		"24": "CS3",
		"26": "AF31",
		"28": "AF32",
		"30": "AF33",
		"32": "CS4",
		"34": "AF41",
		"36": "AF42",
		"38": "AF43",
		"40": "CS5",
		"44": "VA",
		"46": "EF",
		"48": "CS6",
		"56": "CS7",
		"64": "unspecified",
		"8":  "CS1",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewVzOOBBrCPTargetDscpStringNull() VzOOBBrCPTargetDscpStringValue {
	return VzOOBBrCPTargetDscpStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewVzOOBBrCPTargetDscpStringUnknown() VzOOBBrCPTargetDscpStringValue {
	return VzOOBBrCPTargetDscpStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewVzOOBBrCPTargetDscpStringValue(value string) VzOOBBrCPTargetDscpStringValue {
	return VzOOBBrCPTargetDscpStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewVzOOBBrCPTargetDscpStringPointerValue(value *string) VzOOBBrCPTargetDscpStringValue {
	return VzOOBBrCPTargetDscpStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
