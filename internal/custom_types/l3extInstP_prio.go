package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// L3extInstPPrio custom string type.

var _ basetypes.StringTypable = L3extInstPPrioStringType{}

type L3extInstPPrioStringType struct {
	basetypes.StringType
}

func (t L3extInstPPrioStringType) Equal(o attr.Type) bool {
	other, ok := o.(L3extInstPPrioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t L3extInstPPrioStringType) String() string {
	return "L3extInstPPrioStringType"
}

func (t L3extInstPPrioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := L3extInstPPrioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t L3extInstPPrioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t L3extInstPPrioStringType) ValueType(ctx context.Context) attr.Value {
	return L3extInstPPrioStringValue{}
}

// L3extInstPPrio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = L3extInstPPrioStringValue{}

type L3extInstPPrioStringValue struct {
	basetypes.StringValue
}

func (v L3extInstPPrioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(L3extInstPPrioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v L3extInstPPrioStringValue) Type(ctx context.Context) attr.Type {
	return L3extInstPPrioStringType{}
}

func (v L3extInstPPrioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(L3extInstPPrioStringValue)

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

	priorMappedValue := L3extInstPPrioValueMap(v.StringValue)

	newMappedValue := L3extInstPPrioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v L3extInstPPrioStringValue) NamedValueString() string {
	return L3extInstPPrioValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func L3extInstPPrioValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0": "unspecified",
		"1": "level3",
		"2": "level2",
		"3": "level1",
		"7": "level6",
		"8": "level5",
		"9": "level4",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewL3extInstPPrioStringNull() L3extInstPPrioStringValue {
	return L3extInstPPrioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewL3extInstPPrioStringUnknown() L3extInstPPrioStringValue {
	return L3extInstPPrioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewL3extInstPPrioStringValue(value string) L3extInstPPrioStringValue {
	return L3extInstPPrioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewL3extInstPPrioStringPointerValue(value *string) L3extInstPPrioStringValue {
	return L3extInstPPrioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
