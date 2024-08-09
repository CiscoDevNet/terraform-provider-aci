package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// FvRsConsPrio custom string type.

var _ basetypes.StringTypable = FvRsConsPrioStringType{}

type FvRsConsPrioStringType struct {
	basetypes.StringType
}

func (t FvRsConsPrioStringType) Equal(o attr.Type) bool {
	other, ok := o.(FvRsConsPrioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t FvRsConsPrioStringType) String() string {
	return "FvRsConsPrioStringType"
}

func (t FvRsConsPrioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := FvRsConsPrioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t FvRsConsPrioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t FvRsConsPrioStringType) ValueType(ctx context.Context) attr.Value {
	return FvRsConsPrioStringValue{}
}

// FvRsConsPrio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = FvRsConsPrioStringValue{}

type FvRsConsPrioStringValue struct {
	basetypes.StringValue
}

func (v FvRsConsPrioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(FvRsConsPrioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v FvRsConsPrioStringValue) Type(ctx context.Context) attr.Type {
	return FvRsConsPrioStringType{}
}

func (v FvRsConsPrioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(FvRsConsPrioStringValue)

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

	priorMappedValue := FvRsConsPrioValueMap(v.StringValue)

	newMappedValue := FvRsConsPrioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func FvRsConsPrioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewFvRsConsPrioStringNull() FvRsConsPrioStringValue {
	return FvRsConsPrioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewFvRsConsPrioStringUnknown() FvRsConsPrioStringValue {
	return FvRsConsPrioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewFvRsConsPrioStringValue(value string) FvRsConsPrioStringValue {
	return FvRsConsPrioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewFvRsConsPrioStringPointerValue(value *string) FvRsConsPrioStringValue {
	return FvRsConsPrioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
