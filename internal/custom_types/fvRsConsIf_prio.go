package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// FvRsConsIfPrio custom string type.

var _ basetypes.StringTypable = FvRsConsIfPrioStringType{}

type FvRsConsIfPrioStringType struct {
	basetypes.StringType
}

func (t FvRsConsIfPrioStringType) Equal(o attr.Type) bool {
	other, ok := o.(FvRsConsIfPrioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t FvRsConsIfPrioStringType) String() string {
	return "FvRsConsIfPrioStringType"
}

func (t FvRsConsIfPrioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := FvRsConsIfPrioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t FvRsConsIfPrioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t FvRsConsIfPrioStringType) ValueType(ctx context.Context) attr.Value {
	return FvRsConsIfPrioStringValue{}
}

// FvRsConsIfPrio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = FvRsConsIfPrioStringValue{}

type FvRsConsIfPrioStringValue struct {
	basetypes.StringValue
}

func (v FvRsConsIfPrioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(FvRsConsIfPrioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v FvRsConsIfPrioStringValue) Type(ctx context.Context) attr.Type {
	return FvRsConsIfPrioStringType{}
}

func (v FvRsConsIfPrioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(FvRsConsIfPrioStringValue)

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

	priorMappedValue := FvRsConsIfPrioValueMap(v.StringValue)

	newMappedValue := FvRsConsIfPrioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func FvRsConsIfPrioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewFvRsConsIfPrioStringNull() FvRsConsIfPrioStringValue {
	return FvRsConsIfPrioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewFvRsConsIfPrioStringUnknown() FvRsConsIfPrioStringValue {
	return FvRsConsIfPrioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewFvRsConsIfPrioStringValue(value string) FvRsConsIfPrioStringValue {
	return FvRsConsIfPrioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewFvRsConsIfPrioStringPointerValue(value *string) FvRsConsIfPrioStringValue {
	return FvRsConsIfPrioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
