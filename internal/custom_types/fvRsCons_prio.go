package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// FvRsConsprio custom string type.

var _ basetypes.StringTypable = FvRsConsprioStringType{}

type FvRsConsprioStringType struct {
	basetypes.StringType
}

func (t FvRsConsprioStringType) Equal(o attr.Type) bool {
	other, ok := o.(FvRsConsprioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t FvRsConsprioStringType) String() string {
	return "FvRsConsprioStringType"
}

func (t FvRsConsprioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := FvRsConsprioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t FvRsConsprioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t FvRsConsprioStringType) ValueType(ctx context.Context) attr.Value {
	return FvRsConsprioStringValue{}
}

// FvRsConsprio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = FvRsConsprioStringValue{}

type FvRsConsprioStringValue struct {
	basetypes.StringValue
}

func (v FvRsConsprioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(FvRsConsprioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v FvRsConsprioStringValue) Type(ctx context.Context) attr.Type {
	return FvRsConsprioStringType{}
}

func (v FvRsConsprioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(FvRsConsprioStringValue)

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

	priorMappedValue := FvRsConsprioValueMap(v.StringValue)

	newMappedValue := FvRsConsprioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func FvRsConsprioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewFvRsConsprioStringNull() FvRsConsprioStringValue {
	return FvRsConsprioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewFvRsConsprioStringUnknown() FvRsConsprioStringValue {
	return FvRsConsprioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewFvRsConsprioStringValue(value string) FvRsConsprioStringValue {
	return FvRsConsprioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewFvRsConsprioStringPointerValue(value *string) FvRsConsprioStringValue {
	return FvRsConsprioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
