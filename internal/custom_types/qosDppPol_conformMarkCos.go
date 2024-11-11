package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDppPolConformMarkCos custom string type.

var _ basetypes.StringTypable = QosDppPolConformMarkCosStringType{}

type QosDppPolConformMarkCosStringType struct {
	basetypes.StringType
}

func (t QosDppPolConformMarkCosStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDppPolConformMarkCosStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDppPolConformMarkCosStringType) String() string {
	return "QosDppPolConformMarkCosStringType"
}

func (t QosDppPolConformMarkCosStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDppPolConformMarkCosStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDppPolConformMarkCosStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDppPolConformMarkCosStringType) ValueType(ctx context.Context) attr.Value {
	return QosDppPolConformMarkCosStringValue{}
}

// QosDppPolConformMarkCos custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDppPolConformMarkCosStringValue{}

type QosDppPolConformMarkCosStringValue struct {
	basetypes.StringValue
}

func (v QosDppPolConformMarkCosStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDppPolConformMarkCosStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDppPolConformMarkCosStringValue) Type(ctx context.Context) attr.Type {
	return QosDppPolConformMarkCosStringType{}
}

func (v QosDppPolConformMarkCosStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDppPolConformMarkCosStringValue)

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

	priorMappedValue := QosDppPolConformMarkCosValueMap(v.StringValue)

	newMappedValue := QosDppPolConformMarkCosValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v QosDppPolConformMarkCosStringValue) NamedValueString() string {
	return QosDppPolConformMarkCosValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func QosDppPolConformMarkCosValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewQosDppPolConformMarkCosStringNull() QosDppPolConformMarkCosStringValue {
	return QosDppPolConformMarkCosStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDppPolConformMarkCosStringUnknown() QosDppPolConformMarkCosStringValue {
	return QosDppPolConformMarkCosStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDppPolConformMarkCosStringValue(value string) QosDppPolConformMarkCosStringValue {
	return QosDppPolConformMarkCosStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDppPolConformMarkCosStringPointerValue(value *string) QosDppPolConformMarkCosStringValue {
	return QosDppPolConformMarkCosStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
