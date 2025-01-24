package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDscpClassFrom custom string type.

var _ basetypes.StringTypable = QosDscpClassFromStringType{}

type QosDscpClassFromStringType struct {
	basetypes.StringType
}

func (t QosDscpClassFromStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDscpClassFromStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDscpClassFromStringType) String() string {
	return "QosDscpClassFromStringType"
}

func (t QosDscpClassFromStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDscpClassFromStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDscpClassFromStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDscpClassFromStringType) ValueType(ctx context.Context) attr.Value {
	return QosDscpClassFromStringValue{}
}

// QosDscpClassFrom custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDscpClassFromStringValue{}

type QosDscpClassFromStringValue struct {
	basetypes.StringValue
}

func (v QosDscpClassFromStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDscpClassFromStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDscpClassFromStringValue) Type(ctx context.Context) attr.Type {
	return QosDscpClassFromStringType{}
}

func (v QosDscpClassFromStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDscpClassFromStringValue)

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

	priorMappedValue := QosDscpClassFromValueMap(v.StringValue)

	newMappedValue := QosDscpClassFromValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v QosDscpClassFromStringValue) NamedValueString() string {
	return QosDscpClassFromValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func QosDscpClassFromValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewQosDscpClassFromStringNull() QosDscpClassFromStringValue {
	return QosDscpClassFromStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDscpClassFromStringUnknown() QosDscpClassFromStringValue {
	return QosDscpClassFromStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDscpClassFromStringValue(value string) QosDscpClassFromStringValue {
	return QosDscpClassFromStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDscpClassFromStringPointerValue(value *string) QosDscpClassFromStringValue {
	return QosDscpClassFromStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
