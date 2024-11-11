package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// FvRsProvPrio custom string type.

var _ basetypes.StringTypable = FvRsProvPrioStringType{}

type FvRsProvPrioStringType struct {
	basetypes.StringType
}

func (t FvRsProvPrioStringType) Equal(o attr.Type) bool {
	other, ok := o.(FvRsProvPrioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t FvRsProvPrioStringType) String() string {
	return "FvRsProvPrioStringType"
}

func (t FvRsProvPrioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := FvRsProvPrioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t FvRsProvPrioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t FvRsProvPrioStringType) ValueType(ctx context.Context) attr.Value {
	return FvRsProvPrioStringValue{}
}

// FvRsProvPrio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = FvRsProvPrioStringValue{}

type FvRsProvPrioStringValue struct {
	basetypes.StringValue
}

func (v FvRsProvPrioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(FvRsProvPrioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v FvRsProvPrioStringValue) Type(ctx context.Context) attr.Type {
	return FvRsProvPrioStringType{}
}

func (v FvRsProvPrioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(FvRsProvPrioStringValue)

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

	priorMappedValue := FvRsProvPrioValueMap(v.StringValue)

	newMappedValue := FvRsProvPrioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v FvRsProvPrioStringValue) NamedValueString() string {
	return FvRsProvPrioValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func FvRsProvPrioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewFvRsProvPrioStringNull() FvRsProvPrioStringValue {
	return FvRsProvPrioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewFvRsProvPrioStringUnknown() FvRsProvPrioStringValue {
	return FvRsProvPrioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewFvRsProvPrioStringValue(value string) FvRsProvPrioStringValue {
	return FvRsProvPrioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewFvRsProvPrioStringPointerValue(value *string) FvRsProvPrioStringValue {
	return FvRsProvPrioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
