package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDot1PClassFrom custom string type.

var _ basetypes.StringTypable = QosDot1PClassFromStringType{}

type QosDot1PClassFromStringType struct {
	basetypes.StringType
}

func (t QosDot1PClassFromStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDot1PClassFromStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDot1PClassFromStringType) String() string {
	return "QosDot1PClassFromStringType"
}

func (t QosDot1PClassFromStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDot1PClassFromStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDot1PClassFromStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDot1PClassFromStringType) ValueType(ctx context.Context) attr.Value {
	return QosDot1PClassFromStringValue{}
}

// QosDot1PClassFrom custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDot1PClassFromStringValue{}

type QosDot1PClassFromStringValue struct {
	basetypes.StringValue
}

func (v QosDot1PClassFromStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDot1PClassFromStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDot1PClassFromStringValue) Type(ctx context.Context) attr.Type {
	return QosDot1PClassFromStringType{}
}

func (v QosDot1PClassFromStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDot1PClassFromStringValue)

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

	priorMappedValue := QosDot1PClassFromValueMap(v.StringValue)

	newMappedValue := QosDot1PClassFromValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v QosDot1PClassFromStringValue) NamedValueString() string {
	return QosDot1PClassFromValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func QosDot1PClassFromValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewQosDot1PClassFromStringNull() QosDot1PClassFromStringValue {
	return QosDot1PClassFromStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDot1PClassFromStringUnknown() QosDot1PClassFromStringValue {
	return QosDot1PClassFromStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDot1PClassFromStringValue(value string) QosDot1PClassFromStringValue {
	return QosDot1PClassFromStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDot1PClassFromStringPointerValue(value *string) QosDot1PClassFromStringValue {
	return QosDot1PClassFromStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
