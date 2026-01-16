package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// StormctrlIfPolRatePps custom string type.

var _ basetypes.StringTypable = StormctrlIfPolRatePpsStringType{}

type StormctrlIfPolRatePpsStringType struct {
	basetypes.StringType
}

func (t StormctrlIfPolRatePpsStringType) Equal(o attr.Type) bool {
	other, ok := o.(StormctrlIfPolRatePpsStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t StormctrlIfPolRatePpsStringType) String() string {
	return "StormctrlIfPolRatePpsStringType"
}

func (t StormctrlIfPolRatePpsStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := StormctrlIfPolRatePpsStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t StormctrlIfPolRatePpsStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t StormctrlIfPolRatePpsStringType) ValueType(ctx context.Context) attr.Value {
	return StormctrlIfPolRatePpsStringValue{}
}

// StormctrlIfPolRatePps custom string value.

var _ basetypes.StringValuableWithSemanticEquals = StormctrlIfPolRatePpsStringValue{}

type StormctrlIfPolRatePpsStringValue struct {
	basetypes.StringValue
}

func (v StormctrlIfPolRatePpsStringValue) Equal(o attr.Value) bool {
	other, ok := o.(StormctrlIfPolRatePpsStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v StormctrlIfPolRatePpsStringValue) Type(ctx context.Context) attr.Type {
	return StormctrlIfPolRatePpsStringType{}
}

func (v StormctrlIfPolRatePpsStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(StormctrlIfPolRatePpsStringValue)

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

	priorMappedValue := StormctrlIfPolRatePpsValueMap(v.StringValue)

	newMappedValue := StormctrlIfPolRatePpsValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v StormctrlIfPolRatePpsStringValue) NamedValueString() string {
	return StormctrlIfPolRatePpsValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func StormctrlIfPolRatePpsValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffffffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewStormctrlIfPolRatePpsStringNull() StormctrlIfPolRatePpsStringValue {
	return StormctrlIfPolRatePpsStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewStormctrlIfPolRatePpsStringUnknown() StormctrlIfPolRatePpsStringValue {
	return StormctrlIfPolRatePpsStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewStormctrlIfPolRatePpsStringValue(value string) StormctrlIfPolRatePpsStringValue {
	return StormctrlIfPolRatePpsStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewStormctrlIfPolRatePpsStringPointerValue(value *string) StormctrlIfPolRatePpsStringValue {
	return StormctrlIfPolRatePpsStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
