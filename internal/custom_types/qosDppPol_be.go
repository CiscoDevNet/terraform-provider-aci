package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDppPolBe custom string type.

var _ basetypes.StringTypable = QosDppPolBeStringType{}

type QosDppPolBeStringType struct {
	basetypes.StringType
}

func (t QosDppPolBeStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDppPolBeStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDppPolBeStringType) String() string {
	return "QosDppPolBeStringType"
}

func (t QosDppPolBeStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDppPolBeStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDppPolBeStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDppPolBeStringType) ValueType(ctx context.Context) attr.Value {
	return QosDppPolBeStringValue{}
}

// QosDppPolBe custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDppPolBeStringValue{}

type QosDppPolBeStringValue struct {
	basetypes.StringValue
}

func (v QosDppPolBeStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDppPolBeStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDppPolBeStringValue) Type(ctx context.Context) attr.Type {
	return QosDppPolBeStringType{}
}

func (v QosDppPolBeStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDppPolBeStringValue)

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

	priorMappedValue := QosDppPolBeValueMap(v.StringValue)

	newMappedValue := QosDppPolBeValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func QosDppPolBeValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffffffffffffffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewQosDppPolBeStringNull() QosDppPolBeStringValue {
	return QosDppPolBeStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDppPolBeStringUnknown() QosDppPolBeStringValue {
	return QosDppPolBeStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDppPolBeStringValue(value string) QosDppPolBeStringValue {
	return QosDppPolBeStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDppPolBeStringPointerValue(value *string) QosDppPolBeStringValue {
	return QosDppPolBeStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
