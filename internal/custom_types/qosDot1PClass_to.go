package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDot1PClassTo custom string type.

var _ basetypes.StringTypable = QosDot1PClassToStringType{}

type QosDot1PClassToStringType struct {
	basetypes.StringType
}

func (t QosDot1PClassToStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDot1PClassToStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDot1PClassToStringType) String() string {
	return "QosDot1PClassToStringType"
}

func (t QosDot1PClassToStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDot1PClassToStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDot1PClassToStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDot1PClassToStringType) ValueType(ctx context.Context) attr.Value {
	return QosDot1PClassToStringValue{}
}

// QosDot1PClassTo custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDot1PClassToStringValue{}

type QosDot1PClassToStringValue struct {
	basetypes.StringValue
}

func (v QosDot1PClassToStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDot1PClassToStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDot1PClassToStringValue) Type(ctx context.Context) attr.Type {
	return QosDot1PClassToStringType{}
}

func (v QosDot1PClassToStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDot1PClassToStringValue)

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

	priorMappedValue := QosDot1PClassToValueMap(v.StringValue)

	newMappedValue := QosDot1PClassToValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v QosDot1PClassToStringValue) NamedValueString() string {
	return QosDot1PClassToValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func QosDot1PClassToValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewQosDot1PClassToStringNull() QosDot1PClassToStringValue {
	return QosDot1PClassToStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDot1PClassToStringUnknown() QosDot1PClassToStringValue {
	return QosDot1PClassToStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDot1PClassToStringValue(value string) QosDot1PClassToStringValue {
	return QosDot1PClassToStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDot1PClassToStringPointerValue(value *string) QosDot1PClassToStringValue {
	return QosDot1PClassToStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
