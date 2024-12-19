package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDscpClassTargetCos custom string type.

var _ basetypes.StringTypable = QosDscpClassTargetCosStringType{}

type QosDscpClassTargetCosStringType struct {
	basetypes.StringType
}

func (t QosDscpClassTargetCosStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDscpClassTargetCosStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDscpClassTargetCosStringType) String() string {
	return "QosDscpClassTargetCosStringType"
}

func (t QosDscpClassTargetCosStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDscpClassTargetCosStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDscpClassTargetCosStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDscpClassTargetCosStringType) ValueType(ctx context.Context) attr.Value {
	return QosDscpClassTargetCosStringValue{}
}

// QosDscpClassTargetCos custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDscpClassTargetCosStringValue{}

type QosDscpClassTargetCosStringValue struct {
	basetypes.StringValue
}

func (v QosDscpClassTargetCosStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDscpClassTargetCosStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDscpClassTargetCosStringValue) Type(ctx context.Context) attr.Type {
	return QosDscpClassTargetCosStringType{}
}

func (v QosDscpClassTargetCosStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDscpClassTargetCosStringValue)

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

	priorMappedValue := QosDscpClassTargetCosValueMap(v.StringValue)

	newMappedValue := QosDscpClassTargetCosValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v QosDscpClassTargetCosStringValue) NamedValueString() string {
	return QosDscpClassTargetCosValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func QosDscpClassTargetCosValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewQosDscpClassTargetCosStringNull() QosDscpClassTargetCosStringValue {
	return QosDscpClassTargetCosStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDscpClassTargetCosStringUnknown() QosDscpClassTargetCosStringValue {
	return QosDscpClassTargetCosStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDscpClassTargetCosStringValue(value string) QosDscpClassTargetCosStringValue {
	return QosDscpClassTargetCosStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDscpClassTargetCosStringPointerValue(value *string) QosDscpClassTargetCosStringValue {
	return QosDscpClassTargetCosStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
