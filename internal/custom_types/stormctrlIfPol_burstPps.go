package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// StormctrlIfPolBurstPps custom string type.

var _ basetypes.StringTypable = StormctrlIfPolBurstPpsStringType{}

type StormctrlIfPolBurstPpsStringType struct {
	basetypes.StringType
}

func (t StormctrlIfPolBurstPpsStringType) Equal(o attr.Type) bool {
	other, ok := o.(StormctrlIfPolBurstPpsStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t StormctrlIfPolBurstPpsStringType) String() string {
	return "StormctrlIfPolBurstPpsStringType"
}

func (t StormctrlIfPolBurstPpsStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := StormctrlIfPolBurstPpsStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t StormctrlIfPolBurstPpsStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t StormctrlIfPolBurstPpsStringType) ValueType(ctx context.Context) attr.Value {
	return StormctrlIfPolBurstPpsStringValue{}
}

// StormctrlIfPolBurstPps custom string value.

var _ basetypes.StringValuableWithSemanticEquals = StormctrlIfPolBurstPpsStringValue{}

type StormctrlIfPolBurstPpsStringValue struct {
	basetypes.StringValue
}

func (v StormctrlIfPolBurstPpsStringValue) Equal(o attr.Value) bool {
	other, ok := o.(StormctrlIfPolBurstPpsStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v StormctrlIfPolBurstPpsStringValue) Type(ctx context.Context) attr.Type {
	return StormctrlIfPolBurstPpsStringType{}
}

func (v StormctrlIfPolBurstPpsStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(StormctrlIfPolBurstPpsStringValue)

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

	priorMappedValue := StormctrlIfPolBurstPpsValueMap(v.StringValue)

	newMappedValue := StormctrlIfPolBurstPpsValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v StormctrlIfPolBurstPpsStringValue) NamedValueString() string {
	return StormctrlIfPolBurstPpsValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func StormctrlIfPolBurstPpsValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffffffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewStormctrlIfPolBurstPpsStringNull() StormctrlIfPolBurstPpsStringValue {
	return StormctrlIfPolBurstPpsStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewStormctrlIfPolBurstPpsStringUnknown() StormctrlIfPolBurstPpsStringValue {
	return StormctrlIfPolBurstPpsStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewStormctrlIfPolBurstPpsStringValue(value string) StormctrlIfPolBurstPpsStringValue {
	return StormctrlIfPolBurstPpsStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewStormctrlIfPolBurstPpsStringPointerValue(value *string) StormctrlIfPolBurstPpsStringValue {
	return StormctrlIfPolBurstPpsStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
