package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// FvApPrio custom string type.

var _ basetypes.StringTypable = FvApPrioStringType{}

type FvApPrioStringType struct {
	basetypes.StringType
}

func (t FvApPrioStringType) Equal(o attr.Type) bool {
	other, ok := o.(FvApPrioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t FvApPrioStringType) String() string {
	return "FvApPrioStringType"
}

func (t FvApPrioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := FvApPrioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t FvApPrioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t FvApPrioStringType) ValueType(ctx context.Context) attr.Value {
	return FvApPrioStringValue{}
}

// FvApPrio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = FvApPrioStringValue{}

type FvApPrioStringValue struct {
	basetypes.StringValue
}

func (v FvApPrioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(FvApPrioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v FvApPrioStringValue) Type(ctx context.Context) attr.Type {
	return FvApPrioStringType{}
}

func (v FvApPrioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(FvApPrioStringValue)

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

	priorMappedValue := FvApPrioValueMap(v.StringValue)

	newMappedValue := FvApPrioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v FvApPrioStringValue) NamedValueString() string {
	return FvApPrioValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func FvApPrioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewFvApPrioStringNull() FvApPrioStringValue {
	return FvApPrioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewFvApPrioStringUnknown() FvApPrioStringValue {
	return FvApPrioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewFvApPrioStringValue(value string) FvApPrioStringValue {
	return FvApPrioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewFvApPrioStringPointerValue(value *string) FvApPrioStringValue {
	return FvApPrioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
