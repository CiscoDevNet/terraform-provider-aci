package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// StormctrlIfPolBcRatePps custom string type.

var _ basetypes.StringTypable = StormctrlIfPolBcRatePpsStringType{}

type StormctrlIfPolBcRatePpsStringType struct {
	basetypes.StringType
}

func (t StormctrlIfPolBcRatePpsStringType) Equal(o attr.Type) bool {
	other, ok := o.(StormctrlIfPolBcRatePpsStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t StormctrlIfPolBcRatePpsStringType) String() string {
	return "StormctrlIfPolBcRatePpsStringType"
}

func (t StormctrlIfPolBcRatePpsStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := StormctrlIfPolBcRatePpsStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t StormctrlIfPolBcRatePpsStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t StormctrlIfPolBcRatePpsStringType) ValueType(ctx context.Context) attr.Value {
	return StormctrlIfPolBcRatePpsStringValue{}
}

// StormctrlIfPolBcRatePps custom string value.

var _ basetypes.StringValuableWithSemanticEquals = StormctrlIfPolBcRatePpsStringValue{}

type StormctrlIfPolBcRatePpsStringValue struct {
	basetypes.StringValue
}

func (v StormctrlIfPolBcRatePpsStringValue) Equal(o attr.Value) bool {
	other, ok := o.(StormctrlIfPolBcRatePpsStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v StormctrlIfPolBcRatePpsStringValue) Type(ctx context.Context) attr.Type {
	return StormctrlIfPolBcRatePpsStringType{}
}

func (v StormctrlIfPolBcRatePpsStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(StormctrlIfPolBcRatePpsStringValue)

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

	priorMappedValue := StormctrlIfPolBcRatePpsValueMap(v.StringValue)

	newMappedValue := StormctrlIfPolBcRatePpsValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v StormctrlIfPolBcRatePpsStringValue) NamedValueString() string {
	return StormctrlIfPolBcRatePpsValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func StormctrlIfPolBcRatePpsValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffffffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewStormctrlIfPolBcRatePpsStringNull() StormctrlIfPolBcRatePpsStringValue {
	return StormctrlIfPolBcRatePpsStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewStormctrlIfPolBcRatePpsStringUnknown() StormctrlIfPolBcRatePpsStringValue {
	return StormctrlIfPolBcRatePpsStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewStormctrlIfPolBcRatePpsStringValue(value string) StormctrlIfPolBcRatePpsStringValue {
	return StormctrlIfPolBcRatePpsStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewStormctrlIfPolBcRatePpsStringPointerValue(value *string) StormctrlIfPolBcRatePpsStringValue {
	return StormctrlIfPolBcRatePpsStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
