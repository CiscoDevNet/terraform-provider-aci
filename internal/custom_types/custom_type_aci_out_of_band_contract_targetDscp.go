package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// VzOOBBrCPtargetDscp custom string type.

var _ basetypes.StringTypable = VzOOBBrCPtargetDscpStringType{}

type VzOOBBrCPtargetDscpStringType struct {
	basetypes.StringType
}

func (t VzOOBBrCPtargetDscpStringType) Equal(o attr.Type) bool {
	other, ok := o.(VzOOBBrCPtargetDscpStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t VzOOBBrCPtargetDscpStringType) String() string {
	return "VzOOBBrCPtargetDscpStringType"
}

func (t VzOOBBrCPtargetDscpStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := VzOOBBrCPtargetDscpStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t VzOOBBrCPtargetDscpStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t VzOOBBrCPtargetDscpStringType) ValueType(ctx context.Context) attr.Value {
	return VzOOBBrCPtargetDscpStringValue{}
}

// VzOOBBrCPtargetDscp custom string value.

var _ basetypes.StringValuableWithSemanticEquals = VzOOBBrCPtargetDscpStringValue{}

type VzOOBBrCPtargetDscpStringValue struct {
	basetypes.StringValue
}

func (v VzOOBBrCPtargetDscpStringValue) Equal(o attr.Value) bool {
	other, ok := o.(VzOOBBrCPtargetDscpStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v VzOOBBrCPtargetDscpStringValue) Type(ctx context.Context) attr.Type {
	return VzOOBBrCPtargetDscpStringType{}
}

func (v VzOOBBrCPtargetDscpStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(VzOOBBrCPtargetDscpStringValue)

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

	priorMappedValue := VzOOBBrCPtargetDscpValueMap(v.StringValue)

	newMappedValue := VzOOBBrCPtargetDscpValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func VzOOBBrCPtargetDscpValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewVzOOBBrCPtargetDscpStringNull() VzOOBBrCPtargetDscpStringValue {
	return VzOOBBrCPtargetDscpStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewVzOOBBrCPtargetDscpStringUnknown() VzOOBBrCPtargetDscpStringValue {
	return VzOOBBrCPtargetDscpStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewVzOOBBrCPtargetDscpStringValue(value string) VzOOBBrCPtargetDscpStringValue {
	return VzOOBBrCPtargetDscpStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewVzOOBBrCPtargetDscpStringPointerValue(value *string) VzOOBBrCPtargetDscpStringValue {
	return VzOOBBrCPtargetDscpStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
