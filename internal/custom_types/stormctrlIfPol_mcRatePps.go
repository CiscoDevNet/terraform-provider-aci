package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// StormctrlIfPolMcRatePps custom string type.

var _ basetypes.StringTypable = StormctrlIfPolMcRatePpsStringType{}

type StormctrlIfPolMcRatePpsStringType struct {
	basetypes.StringType
}

func (t StormctrlIfPolMcRatePpsStringType) Equal(o attr.Type) bool {
	other, ok := o.(StormctrlIfPolMcRatePpsStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t StormctrlIfPolMcRatePpsStringType) String() string {
	return "StormctrlIfPolMcRatePpsStringType"
}

func (t StormctrlIfPolMcRatePpsStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := StormctrlIfPolMcRatePpsStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t StormctrlIfPolMcRatePpsStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t StormctrlIfPolMcRatePpsStringType) ValueType(ctx context.Context) attr.Value {
	return StormctrlIfPolMcRatePpsStringValue{}
}

// StormctrlIfPolMcRatePps custom string value.

var _ basetypes.StringValuableWithSemanticEquals = StormctrlIfPolMcRatePpsStringValue{}

type StormctrlIfPolMcRatePpsStringValue struct {
	basetypes.StringValue
}

func (v StormctrlIfPolMcRatePpsStringValue) Equal(o attr.Value) bool {
	other, ok := o.(StormctrlIfPolMcRatePpsStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v StormctrlIfPolMcRatePpsStringValue) Type(ctx context.Context) attr.Type {
	return StormctrlIfPolMcRatePpsStringType{}
}

func (v StormctrlIfPolMcRatePpsStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(StormctrlIfPolMcRatePpsStringValue)

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

	priorMappedValue := StormctrlIfPolMcRatePpsValueMap(v.StringValue)

	newMappedValue := StormctrlIfPolMcRatePpsValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v StormctrlIfPolMcRatePpsStringValue) NamedValueString() string {
	return StormctrlIfPolMcRatePpsValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func StormctrlIfPolMcRatePpsValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffffffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewStormctrlIfPolMcRatePpsStringNull() StormctrlIfPolMcRatePpsStringValue {
	return StormctrlIfPolMcRatePpsStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewStormctrlIfPolMcRatePpsStringUnknown() StormctrlIfPolMcRatePpsStringValue {
	return StormctrlIfPolMcRatePpsStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewStormctrlIfPolMcRatePpsStringValue(value string) StormctrlIfPolMcRatePpsStringValue {
	return StormctrlIfPolMcRatePpsStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewStormctrlIfPolMcRatePpsStringPointerValue(value *string) StormctrlIfPolMcRatePpsStringValue {
	return StormctrlIfPolMcRatePpsStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
