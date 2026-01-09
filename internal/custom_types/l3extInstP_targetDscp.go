package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// L3extInstPTargetDscp custom string type.

var _ basetypes.StringTypable = L3extInstPTargetDscpStringType{}

type L3extInstPTargetDscpStringType struct {
	basetypes.StringType
}

func (t L3extInstPTargetDscpStringType) Equal(o attr.Type) bool {
	other, ok := o.(L3extInstPTargetDscpStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t L3extInstPTargetDscpStringType) String() string {
	return "L3extInstPTargetDscpStringType"
}

func (t L3extInstPTargetDscpStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := L3extInstPTargetDscpStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t L3extInstPTargetDscpStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t L3extInstPTargetDscpStringType) ValueType(ctx context.Context) attr.Value {
	return L3extInstPTargetDscpStringValue{}
}

// L3extInstPTargetDscp custom string value.

var _ basetypes.StringValuableWithSemanticEquals = L3extInstPTargetDscpStringValue{}

type L3extInstPTargetDscpStringValue struct {
	basetypes.StringValue
}

func (v L3extInstPTargetDscpStringValue) Equal(o attr.Value) bool {
	other, ok := o.(L3extInstPTargetDscpStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v L3extInstPTargetDscpStringValue) Type(ctx context.Context) attr.Type {
	return L3extInstPTargetDscpStringType{}
}

func (v L3extInstPTargetDscpStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(L3extInstPTargetDscpStringValue)

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

	priorMappedValue := L3extInstPTargetDscpValueMap(v.StringValue)

	newMappedValue := L3extInstPTargetDscpValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v L3extInstPTargetDscpStringValue) NamedValueString() string {
	return L3extInstPTargetDscpValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func L3extInstPTargetDscpValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewL3extInstPTargetDscpStringNull() L3extInstPTargetDscpStringValue {
	return L3extInstPTargetDscpStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewL3extInstPTargetDscpStringUnknown() L3extInstPTargetDscpStringValue {
	return L3extInstPTargetDscpStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewL3extInstPTargetDscpStringValue(value string) L3extInstPTargetDscpStringValue {
	return L3extInstPTargetDscpStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewL3extInstPTargetDscpStringPointerValue(value *string) L3extInstPTargetDscpStringValue {
	return L3extInstPTargetDscpStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
