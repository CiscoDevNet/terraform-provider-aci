package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// StormctrlIfPolUucRatePps custom string type.

var _ basetypes.StringTypable = StormctrlIfPolUucRatePpsStringType{}

type StormctrlIfPolUucRatePpsStringType struct {
	basetypes.StringType
}

func (t StormctrlIfPolUucRatePpsStringType) Equal(o attr.Type) bool {
	other, ok := o.(StormctrlIfPolUucRatePpsStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t StormctrlIfPolUucRatePpsStringType) String() string {
	return "StormctrlIfPolUucRatePpsStringType"
}

func (t StormctrlIfPolUucRatePpsStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := StormctrlIfPolUucRatePpsStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t StormctrlIfPolUucRatePpsStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t StormctrlIfPolUucRatePpsStringType) ValueType(ctx context.Context) attr.Value {
	return StormctrlIfPolUucRatePpsStringValue{}
}

// StormctrlIfPolUucRatePps custom string value.

var _ basetypes.StringValuableWithSemanticEquals = StormctrlIfPolUucRatePpsStringValue{}

type StormctrlIfPolUucRatePpsStringValue struct {
	basetypes.StringValue
}

func (v StormctrlIfPolUucRatePpsStringValue) Equal(o attr.Value) bool {
	other, ok := o.(StormctrlIfPolUucRatePpsStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v StormctrlIfPolUucRatePpsStringValue) Type(ctx context.Context) attr.Type {
	return StormctrlIfPolUucRatePpsStringType{}
}

func (v StormctrlIfPolUucRatePpsStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(StormctrlIfPolUucRatePpsStringValue)

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

	priorMappedValue := StormctrlIfPolUucRatePpsValueMap(v.StringValue)

	newMappedValue := StormctrlIfPolUucRatePpsValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v StormctrlIfPolUucRatePpsStringValue) NamedValueString() string {
	return StormctrlIfPolUucRatePpsValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func StormctrlIfPolUucRatePpsValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffffffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewStormctrlIfPolUucRatePpsStringNull() StormctrlIfPolUucRatePpsStringValue {
	return StormctrlIfPolUucRatePpsStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewStormctrlIfPolUucRatePpsStringUnknown() StormctrlIfPolUucRatePpsStringValue {
	return StormctrlIfPolUucRatePpsStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewStormctrlIfPolUucRatePpsStringValue(value string) StormctrlIfPolUucRatePpsStringValue {
	return StormctrlIfPolUucRatePpsStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewStormctrlIfPolUucRatePpsStringPointerValue(value *string) StormctrlIfPolUucRatePpsStringValue {
	return StormctrlIfPolUucRatePpsStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
