package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDscpClassTarget custom string type.

var _ basetypes.StringTypable = QosDscpClassTargetStringType{}

type QosDscpClassTargetStringType struct {
	basetypes.StringType
}

func (t QosDscpClassTargetStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDscpClassTargetStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDscpClassTargetStringType) String() string {
	return "QosDscpClassTargetStringType"
}

func (t QosDscpClassTargetStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDscpClassTargetStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDscpClassTargetStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDscpClassTargetStringType) ValueType(ctx context.Context) attr.Value {
	return QosDscpClassTargetStringValue{}
}

// QosDscpClassTarget custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDscpClassTargetStringValue{}

type QosDscpClassTargetStringValue struct {
	basetypes.StringValue
}

func (v QosDscpClassTargetStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDscpClassTargetStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDscpClassTargetStringValue) Type(ctx context.Context) attr.Type {
	return QosDscpClassTargetStringType{}
}

func (v QosDscpClassTargetStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDscpClassTargetStringValue)

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

	priorMappedValue := QosDscpClassTargetValueMap(v.StringValue)

	newMappedValue := QosDscpClassTargetValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v QosDscpClassTargetStringValue) NamedValueString() string {
	return QosDscpClassTargetValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func QosDscpClassTargetValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewQosDscpClassTargetStringNull() QosDscpClassTargetStringValue {
	return QosDscpClassTargetStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDscpClassTargetStringUnknown() QosDscpClassTargetStringValue {
	return QosDscpClassTargetStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDscpClassTargetStringValue(value string) QosDscpClassTargetStringValue {
	return QosDscpClassTargetStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDscpClassTargetStringPointerValue(value *string) QosDscpClassTargetStringValue {
	return QosDscpClassTargetStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
