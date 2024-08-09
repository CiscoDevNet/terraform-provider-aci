package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDppPolViolateMarkCos custom string type.

var _ basetypes.StringTypable = QosDppPolViolateMarkCosStringType{}

type QosDppPolViolateMarkCosStringType struct {
	basetypes.StringType
}

func (t QosDppPolViolateMarkCosStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDppPolViolateMarkCosStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDppPolViolateMarkCosStringType) String() string {
	return "QosDppPolViolateMarkCosStringType"
}

func (t QosDppPolViolateMarkCosStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDppPolViolateMarkCosStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDppPolViolateMarkCosStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDppPolViolateMarkCosStringType) ValueType(ctx context.Context) attr.Value {
	return QosDppPolViolateMarkCosStringValue{}
}

// QosDppPolViolateMarkCos custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDppPolViolateMarkCosStringValue{}

type QosDppPolViolateMarkCosStringValue struct {
	basetypes.StringValue
}

func (v QosDppPolViolateMarkCosStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDppPolViolateMarkCosStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDppPolViolateMarkCosStringValue) Type(ctx context.Context) attr.Type {
	return QosDppPolViolateMarkCosStringType{}
}

func (v QosDppPolViolateMarkCosStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDppPolViolateMarkCosStringValue)

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

	priorMappedValue := QosDppPolViolateMarkCosValueMap(v.StringValue)

	newMappedValue := QosDppPolViolateMarkCosValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func QosDppPolViolateMarkCosValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewQosDppPolViolateMarkCosStringNull() QosDppPolViolateMarkCosStringValue {
	return QosDppPolViolateMarkCosStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDppPolViolateMarkCosStringUnknown() QosDppPolViolateMarkCosStringValue {
	return QosDppPolViolateMarkCosStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDppPolViolateMarkCosStringValue(value string) QosDppPolViolateMarkCosStringValue {
	return QosDppPolViolateMarkCosStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDppPolViolateMarkCosStringPointerValue(value *string) QosDppPolViolateMarkCosStringValue {
	return QosDppPolViolateMarkCosStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
