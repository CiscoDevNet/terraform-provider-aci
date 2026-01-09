package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDppPolExceedMarkDscp custom string type.

var _ basetypes.StringTypable = QosDppPolExceedMarkDscpStringType{}

type QosDppPolExceedMarkDscpStringType struct {
	basetypes.StringType
}

func (t QosDppPolExceedMarkDscpStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDppPolExceedMarkDscpStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDppPolExceedMarkDscpStringType) String() string {
	return "QosDppPolExceedMarkDscpStringType"
}

func (t QosDppPolExceedMarkDscpStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDppPolExceedMarkDscpStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDppPolExceedMarkDscpStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDppPolExceedMarkDscpStringType) ValueType(ctx context.Context) attr.Value {
	return QosDppPolExceedMarkDscpStringValue{}
}

// QosDppPolExceedMarkDscp custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDppPolExceedMarkDscpStringValue{}

type QosDppPolExceedMarkDscpStringValue struct {
	basetypes.StringValue
}

func (v QosDppPolExceedMarkDscpStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDppPolExceedMarkDscpStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDppPolExceedMarkDscpStringValue) Type(ctx context.Context) attr.Type {
	return QosDppPolExceedMarkDscpStringType{}
}

func (v QosDppPolExceedMarkDscpStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDppPolExceedMarkDscpStringValue)

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

	priorMappedValue := QosDppPolExceedMarkDscpValueMap(v.StringValue)

	newMappedValue := QosDppPolExceedMarkDscpValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v QosDppPolExceedMarkDscpStringValue) NamedValueString() string {
	return QosDppPolExceedMarkDscpValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func QosDppPolExceedMarkDscpValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewQosDppPolExceedMarkDscpStringNull() QosDppPolExceedMarkDscpStringValue {
	return QosDppPolExceedMarkDscpStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDppPolExceedMarkDscpStringUnknown() QosDppPolExceedMarkDscpStringValue {
	return QosDppPolExceedMarkDscpStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDppPolExceedMarkDscpStringValue(value string) QosDppPolExceedMarkDscpStringValue {
	return QosDppPolExceedMarkDscpStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDppPolExceedMarkDscpStringPointerValue(value *string) QosDppPolExceedMarkDscpStringValue {
	return QosDppPolExceedMarkDscpStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
