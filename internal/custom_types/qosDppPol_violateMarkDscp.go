package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// QosDppPolViolateMarkDscp custom string type.

var _ basetypes.StringTypable = QosDppPolViolateMarkDscpStringType{}

type QosDppPolViolateMarkDscpStringType struct {
	basetypes.StringType
}

func (t QosDppPolViolateMarkDscpStringType) Equal(o attr.Type) bool {
	other, ok := o.(QosDppPolViolateMarkDscpStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t QosDppPolViolateMarkDscpStringType) String() string {
	return "QosDppPolViolateMarkDscpStringType"
}

func (t QosDppPolViolateMarkDscpStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := QosDppPolViolateMarkDscpStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t QosDppPolViolateMarkDscpStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t QosDppPolViolateMarkDscpStringType) ValueType(ctx context.Context) attr.Value {
	return QosDppPolViolateMarkDscpStringValue{}
}

// QosDppPolViolateMarkDscp custom string value.

var _ basetypes.StringValuableWithSemanticEquals = QosDppPolViolateMarkDscpStringValue{}

type QosDppPolViolateMarkDscpStringValue struct {
	basetypes.StringValue
}

func (v QosDppPolViolateMarkDscpStringValue) Equal(o attr.Value) bool {
	other, ok := o.(QosDppPolViolateMarkDscpStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v QosDppPolViolateMarkDscpStringValue) Type(ctx context.Context) attr.Type {
	return QosDppPolViolateMarkDscpStringType{}
}

func (v QosDppPolViolateMarkDscpStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(QosDppPolViolateMarkDscpStringValue)

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

	priorMappedValue := QosDppPolViolateMarkDscpValueMap(v.StringValue)

	newMappedValue := QosDppPolViolateMarkDscpValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func QosDppPolViolateMarkDscpValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewQosDppPolViolateMarkDscpStringNull() QosDppPolViolateMarkDscpStringValue {
	return QosDppPolViolateMarkDscpStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewQosDppPolViolateMarkDscpStringUnknown() QosDppPolViolateMarkDscpStringValue {
	return QosDppPolViolateMarkDscpStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewQosDppPolViolateMarkDscpStringValue(value string) QosDppPolViolateMarkDscpStringValue {
	return QosDppPolViolateMarkDscpStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewQosDppPolViolateMarkDscpStringPointerValue(value *string) QosDppPolViolateMarkDscpStringValue {
	return QosDppPolViolateMarkDscpStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
