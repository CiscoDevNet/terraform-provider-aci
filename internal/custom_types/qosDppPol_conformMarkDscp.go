package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDppPolConformMarkDscp custom string type.

var _ basetypes.StringTypable = QosDppPolConformMarkDscpStringType{}

type QosDppPolConformMarkDscpStringType struct {
	basetypes.StringType
}

func (t QosDppPolConformMarkDscpStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDppPolConformMarkDscpStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDppPolConformMarkDscpStringType) String() string {
	return "QosDppPolConformMarkDscpStringType"
}

func (t QosDppPolConformMarkDscpStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDppPolConformMarkDscpStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDppPolConformMarkDscpStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDppPolConformMarkDscpStringType) ValueType(ctx context.Context) attr.Value {
	return QosDppPolConformMarkDscpStringValue{}
}

// QosDppPolConformMarkDscp custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDppPolConformMarkDscpStringValue{}

type QosDppPolConformMarkDscpStringValue struct {
	basetypes.StringValue
}

func (v QosDppPolConformMarkDscpStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDppPolConformMarkDscpStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDppPolConformMarkDscpStringValue) Type(ctx context.Context) attr.Type {
	return QosDppPolConformMarkDscpStringType{}
}

func (v QosDppPolConformMarkDscpStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDppPolConformMarkDscpStringValue)

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

	priorMappedValue := QosDppPolConformMarkDscpValueMap(v.StringValue)

	newMappedValue := QosDppPolConformMarkDscpValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func QosDppPolConformMarkDscpValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewQosDppPolConformMarkDscpStringNull() QosDppPolConformMarkDscpStringValue {
	return QosDppPolConformMarkDscpStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDppPolConformMarkDscpStringUnknown() QosDppPolConformMarkDscpStringValue {
	return QosDppPolConformMarkDscpStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDppPolConformMarkDscpStringValue(value string) QosDppPolConformMarkDscpStringValue {
	return QosDppPolConformMarkDscpStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDppPolConformMarkDscpStringPointerValue(value *string) QosDppPolConformMarkDscpStringValue {
	return QosDppPolConformMarkDscpStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
