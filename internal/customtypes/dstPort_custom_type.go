package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// DstPort custom string type.

var _ basetypes.StringTypable = DstPortStringType{}

type DstPortStringType struct {
	basetypes.StringType
}

func (t DstPortStringType) Equal(o attr.Type) bool {
	other, ok := o.(DstPortStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t DstPortStringType) String() string {
	return "DstPortStringType"
}

func (t DstPortStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := DstPortStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t DstPortStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t DstPortStringType) ValueType(ctx context.Context) attr.Value {
	return DstPortStringValue{}
}

// DstPort custom string value.

var _ basetypes.StringValuableWithSemanticEquals = DstPortStringValue{}

type DstPortStringValue struct {
	basetypes.StringValue
}

func (v DstPortStringValue) Equal(o attr.Value) bool {
	other, ok := o.(DstPortStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v DstPortStringValue) Type(ctx context.Context) attr.Type {
	return DstPortStringType{}
}

func (v DstPortStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(DstPortStringValue)

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

	priorMappedValue := DstPortValueMap(v.StringValue)

	newMappedValue := DstPortValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func DstPortValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{"443": "https"}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewDstPortStringNull() DstPortStringValue {
	return DstPortStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewDstPortStringUnknown() DstPortStringValue {
	return DstPortStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewDstPortStringValue(value string) DstPortStringValue {
	return DstPortStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewDstPortStringPointerValue(value *string) DstPortStringValue {
	return DstPortStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
