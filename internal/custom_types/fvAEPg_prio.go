package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// FvAEPgPrio custom string type.

var _ basetypes.StringTypable = FvAEPgPrioStringType{}

type FvAEPgPrioStringType struct {
	basetypes.StringType
}

func (t FvAEPgPrioStringType) Equal(o attr.Type) bool {
	other, ok := o.(FvAEPgPrioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t FvAEPgPrioStringType) String() string {
	return "FvAEPgPrioStringType"
}

func (t FvAEPgPrioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := FvAEPgPrioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t FvAEPgPrioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t FvAEPgPrioStringType) ValueType(ctx context.Context) attr.Value {
	return FvAEPgPrioStringValue{}
}

// FvAEPgPrio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = FvAEPgPrioStringValue{}

type FvAEPgPrioStringValue struct {
	basetypes.StringValue
}

func (v FvAEPgPrioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(FvAEPgPrioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v FvAEPgPrioStringValue) Type(ctx context.Context) attr.Type {
	return FvAEPgPrioStringType{}
}

func (v FvAEPgPrioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(FvAEPgPrioStringValue)

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

	priorMappedValue := FvAEPgPrioValueMap(v.StringValue)

	newMappedValue := FvAEPgPrioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func FvAEPgPrioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewFvAEPgPrioStringNull() FvAEPgPrioStringValue {
	return FvAEPgPrioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewFvAEPgPrioStringUnknown() FvAEPgPrioStringValue {
	return FvAEPgPrioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewFvAEPgPrioStringValue(value string) FvAEPgPrioStringValue {
	return FvAEPgPrioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewFvAEPgPrioStringPointerValue(value *string) FvAEPgPrioStringValue {
	return FvAEPgPrioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
