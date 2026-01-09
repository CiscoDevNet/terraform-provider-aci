package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDot1PClassTarget custom string type.

var _ basetypes.StringTypable = QosDot1PClassTargetStringType{}

type QosDot1PClassTargetStringType struct {
	basetypes.StringType
}

func (t QosDot1PClassTargetStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDot1PClassTargetStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDot1PClassTargetStringType) String() string {
	return "QosDot1PClassTargetStringType"
}

func (t QosDot1PClassTargetStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDot1PClassTargetStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDot1PClassTargetStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDot1PClassTargetStringType) ValueType(ctx context.Context) attr.Value {
	return QosDot1PClassTargetStringValue{}
}

// QosDot1PClassTarget custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDot1PClassTargetStringValue{}

type QosDot1PClassTargetStringValue struct {
	basetypes.StringValue
}

func (v QosDot1PClassTargetStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDot1PClassTargetStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDot1PClassTargetStringValue) Type(ctx context.Context) attr.Type {
	return QosDot1PClassTargetStringType{}
}

func (v QosDot1PClassTargetStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDot1PClassTargetStringValue)

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

	priorMappedValue := QosDot1PClassTargetValueMap(v.StringValue)

	newMappedValue := QosDot1PClassTargetValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v QosDot1PClassTargetStringValue) NamedValueString() string {
	return QosDot1PClassTargetValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func QosDot1PClassTargetValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewQosDot1PClassTargetStringNull() QosDot1PClassTargetStringValue {
	return QosDot1PClassTargetStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDot1PClassTargetStringUnknown() QosDot1PClassTargetStringValue {
	return QosDot1PClassTargetStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDot1PClassTargetStringValue(value string) QosDot1PClassTargetStringValue {
	return QosDot1PClassTargetStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDot1PClassTargetStringPointerValue(value *string) QosDot1PClassTargetStringValue {
	return QosDot1PClassTargetStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
