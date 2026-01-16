package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// StormctrlIfPolBcBurstPps custom string type.

var _ basetypes.StringTypable = StormctrlIfPolBcBurstPpsStringType{}

type StormctrlIfPolBcBurstPpsStringType struct {
	basetypes.StringType
}

func (t StormctrlIfPolBcBurstPpsStringType) Equal(o attr.Type) bool {
	other, ok := o.(StormctrlIfPolBcBurstPpsStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t StormctrlIfPolBcBurstPpsStringType) String() string {
	return "StormctrlIfPolBcBurstPpsStringType"
}

func (t StormctrlIfPolBcBurstPpsStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := StormctrlIfPolBcBurstPpsStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t StormctrlIfPolBcBurstPpsStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t StormctrlIfPolBcBurstPpsStringType) ValueType(ctx context.Context) attr.Value {
	return StormctrlIfPolBcBurstPpsStringValue{}
}

// StormctrlIfPolBcBurstPps custom string value.

var _ basetypes.StringValuableWithSemanticEquals = StormctrlIfPolBcBurstPpsStringValue{}

type StormctrlIfPolBcBurstPpsStringValue struct {
	basetypes.StringValue
}

func (v StormctrlIfPolBcBurstPpsStringValue) Equal(o attr.Value) bool {
	other, ok := o.(StormctrlIfPolBcBurstPpsStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v StormctrlIfPolBcBurstPpsStringValue) Type(ctx context.Context) attr.Type {
	return StormctrlIfPolBcBurstPpsStringType{}
}

func (v StormctrlIfPolBcBurstPpsStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(StormctrlIfPolBcBurstPpsStringValue)

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

	priorMappedValue := StormctrlIfPolBcBurstPpsValueMap(v.StringValue)

	newMappedValue := StormctrlIfPolBcBurstPpsValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v StormctrlIfPolBcBurstPpsStringValue) NamedValueString() string {
	return StormctrlIfPolBcBurstPpsValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func StormctrlIfPolBcBurstPpsValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffffffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewStormctrlIfPolBcBurstPpsStringNull() StormctrlIfPolBcBurstPpsStringValue {
	return StormctrlIfPolBcBurstPpsStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewStormctrlIfPolBcBurstPpsStringUnknown() StormctrlIfPolBcBurstPpsStringValue {
	return StormctrlIfPolBcBurstPpsStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewStormctrlIfPolBcBurstPpsStringValue(value string) StormctrlIfPolBcBurstPpsStringValue {
	return StormctrlIfPolBcBurstPpsStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewStormctrlIfPolBcBurstPpsStringPointerValue(value *string) StormctrlIfPolBcBurstPpsStringValue {
	return StormctrlIfPolBcBurstPpsStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
