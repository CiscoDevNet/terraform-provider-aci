package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// StormctrlIfPolMcBurstPps custom string type.

var _ basetypes.StringTypable = StormctrlIfPolMcBurstPpsStringType{}

type StormctrlIfPolMcBurstPpsStringType struct {
	basetypes.StringType
}

func (t StormctrlIfPolMcBurstPpsStringType) Equal(o attr.Type) bool {
	other, ok := o.(StormctrlIfPolMcBurstPpsStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t StormctrlIfPolMcBurstPpsStringType) String() string {
	return "StormctrlIfPolMcBurstPpsStringType"
}

func (t StormctrlIfPolMcBurstPpsStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := StormctrlIfPolMcBurstPpsStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t StormctrlIfPolMcBurstPpsStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t StormctrlIfPolMcBurstPpsStringType) ValueType(ctx context.Context) attr.Value {
	return StormctrlIfPolMcBurstPpsStringValue{}
}

// StormctrlIfPolMcBurstPps custom string value.

var _ basetypes.StringValuableWithSemanticEquals = StormctrlIfPolMcBurstPpsStringValue{}

type StormctrlIfPolMcBurstPpsStringValue struct {
	basetypes.StringValue
}

func (v StormctrlIfPolMcBurstPpsStringValue) Equal(o attr.Value) bool {
	other, ok := o.(StormctrlIfPolMcBurstPpsStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v StormctrlIfPolMcBurstPpsStringValue) Type(ctx context.Context) attr.Type {
	return StormctrlIfPolMcBurstPpsStringType{}
}

func (v StormctrlIfPolMcBurstPpsStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(StormctrlIfPolMcBurstPpsStringValue)

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

	priorMappedValue := StormctrlIfPolMcBurstPpsValueMap(v.StringValue)

	newMappedValue := StormctrlIfPolMcBurstPpsValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v StormctrlIfPolMcBurstPpsStringValue) NamedValueString() string {
	return StormctrlIfPolMcBurstPpsValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func StormctrlIfPolMcBurstPpsValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0xffffffff": "unspecified",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewStormctrlIfPolMcBurstPpsStringNull() StormctrlIfPolMcBurstPpsStringValue {
	return StormctrlIfPolMcBurstPpsStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewStormctrlIfPolMcBurstPpsStringUnknown() StormctrlIfPolMcBurstPpsStringValue {
	return StormctrlIfPolMcBurstPpsStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewStormctrlIfPolMcBurstPpsStringValue(value string) StormctrlIfPolMcBurstPpsStringValue {
	return StormctrlIfPolMcBurstPpsStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewStormctrlIfPolMcBurstPpsStringPointerValue(value *string) StormctrlIfPolMcBurstPpsStringValue {
	return StormctrlIfPolMcBurstPpsStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
