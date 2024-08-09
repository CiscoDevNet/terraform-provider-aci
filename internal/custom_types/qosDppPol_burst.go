package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDppPolBurst custom string type.

var _ basetypes.StringTypable = QosDppPolBurstStringType{}

type QosDppPolBurstStringType struct {
	basetypes.StringType
}

func (t QosDppPolBurstStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDppPolBurstStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDppPolBurstStringType) String() string {
	return "QosDppPolBurstStringType"
}

func (t QosDppPolBurstStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDppPolBurstStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDppPolBurstStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDppPolBurstStringType) ValueType(ctx context.Context) attr.Value {
	return QosDppPolBurstStringValue{}
}

// QosDppPolBurst custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDppPolBurstStringValue{}

type QosDppPolBurstStringValue struct {
	basetypes.StringValue
}

func (v QosDppPolBurstStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDppPolBurstStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDppPolBurstStringValue) Type(ctx context.Context) attr.Type {
	return QosDppPolBurstStringType{}
}

func (v QosDppPolBurstStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDppPolBurstStringValue)

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

	priorMappedValue := QosDppPolBurstValueMap(v.StringValue)

	newMappedValue := QosDppPolBurstValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func QosDppPolBurstValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffffffffffffffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewQosDppPolBurstStringNull() QosDppPolBurstStringValue {
	return QosDppPolBurstStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDppPolBurstStringUnknown() QosDppPolBurstStringValue {
	return QosDppPolBurstStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDppPolBurstStringValue(value string) QosDppPolBurstStringValue {
	return QosDppPolBurstStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDppPolBurstStringPointerValue(value *string) QosDppPolBurstStringValue {
	return QosDppPolBurstStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
