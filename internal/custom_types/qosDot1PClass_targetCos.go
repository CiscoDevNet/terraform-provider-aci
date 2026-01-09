package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDot1PClassTargetCos custom string type.

var _ basetypes.StringTypable = QosDot1PClassTargetCosStringType{}

type QosDot1PClassTargetCosStringType struct {
	basetypes.StringType
}

func (t QosDot1PClassTargetCosStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDot1PClassTargetCosStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDot1PClassTargetCosStringType) String() string {
	return "QosDot1PClassTargetCosStringType"
}

func (t QosDot1PClassTargetCosStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDot1PClassTargetCosStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDot1PClassTargetCosStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDot1PClassTargetCosStringType) ValueType(ctx context.Context) attr.Value {
	return QosDot1PClassTargetCosStringValue{}
}

// QosDot1PClassTargetCos custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDot1PClassTargetCosStringValue{}

type QosDot1PClassTargetCosStringValue struct {
	basetypes.StringValue
}

func (v QosDot1PClassTargetCosStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDot1PClassTargetCosStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDot1PClassTargetCosStringValue) Type(ctx context.Context) attr.Type {
	return QosDot1PClassTargetCosStringType{}
}

func (v QosDot1PClassTargetCosStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDot1PClassTargetCosStringValue)

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

	priorMappedValue := QosDot1PClassTargetCosValueMap(v.StringValue)

	newMappedValue := QosDot1PClassTargetCosValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v QosDot1PClassTargetCosStringValue) NamedValueString() string {
	return QosDot1PClassTargetCosValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func QosDot1PClassTargetCosValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0": "0",
		"1": "1",
		"2": "2",
		"3": "3",
		"4": "4",
		"5": "5",
		"6": "6",
		"7": "7",
		"8": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewQosDot1PClassTargetCosStringNull() QosDot1PClassTargetCosStringValue {
	return QosDot1PClassTargetCosStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDot1PClassTargetCosStringUnknown() QosDot1PClassTargetCosStringValue {
	return QosDot1PClassTargetCosStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDot1PClassTargetCosStringValue(value string) QosDot1PClassTargetCosStringValue {
	return QosDot1PClassTargetCosStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDot1PClassTargetCosStringPointerValue(value *string) QosDot1PClassTargetCosStringValue {
	return QosDot1PClassTargetCosStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
