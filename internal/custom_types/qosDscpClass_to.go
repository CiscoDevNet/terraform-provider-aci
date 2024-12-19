package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDscpClassTo custom string type.

var _ basetypes.StringTypable = QosDscpClassToStringType{}

type QosDscpClassToStringType struct {
	basetypes.StringType
}

func (t QosDscpClassToStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDscpClassToStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDscpClassToStringType) String() string {
	return "QosDscpClassToStringType"
}

func (t QosDscpClassToStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDscpClassToStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDscpClassToStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDscpClassToStringType) ValueType(ctx context.Context) attr.Value {
	return QosDscpClassToStringValue{}
}

// QosDscpClassTo custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDscpClassToStringValue{}

type QosDscpClassToStringValue struct {
	basetypes.StringValue
}

func (v QosDscpClassToStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDscpClassToStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDscpClassToStringValue) Type(ctx context.Context) attr.Type {
	return QosDscpClassToStringType{}
}

func (v QosDscpClassToStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDscpClassToStringValue)

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

	priorMappedValue := QosDscpClassToValueMap(v.StringValue)

	newMappedValue := QosDscpClassToValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v QosDscpClassToStringValue) NamedValueString() string {
	return QosDscpClassToValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func QosDscpClassToValueMap(value basetypes.StringValue) basetypes.StringValue {
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
		"8":  "CS1",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewQosDscpClassToStringNull() QosDscpClassToStringValue {
	return QosDscpClassToStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDscpClassToStringUnknown() QosDscpClassToStringValue {
	return QosDscpClassToStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDscpClassToStringValue(value string) QosDscpClassToStringValue {
	return QosDscpClassToStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDscpClassToStringPointerValue(value *string) QosDscpClassToStringValue {
	return QosDscpClassToStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
