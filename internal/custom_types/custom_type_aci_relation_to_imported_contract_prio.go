package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// FvRsConsIfprio custom string type.

var _ basetypes.StringTypable = FvRsConsIfprioStringType{}

type FvRsConsIfprioStringType struct {
	basetypes.StringType
}

func (t FvRsConsIfprioStringType) Equal(o attr.Type) bool {
	other, ok := o.(FvRsConsIfprioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t FvRsConsIfprioStringType) String() string {
	return "FvRsConsIfprioStringType"
}

func (t FvRsConsIfprioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := FvRsConsIfprioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t FvRsConsIfprioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t FvRsConsIfprioStringType) ValueType(ctx context.Context) attr.Value {
	return FvRsConsIfprioStringValue{}
}

// FvRsConsIfprio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = FvRsConsIfprioStringValue{}

type FvRsConsIfprioStringValue struct {
	basetypes.StringValue
}

func (v FvRsConsIfprioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(FvRsConsIfprioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v FvRsConsIfprioStringValue) Type(ctx context.Context) attr.Type {
	return FvRsConsIfprioStringType{}
}

func (v FvRsConsIfprioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(FvRsConsIfprioStringValue)

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

	priorMappedValue := FvRsConsIfprioValueMap(v.StringValue)

	newMappedValue := FvRsConsIfprioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func FvRsConsIfprioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewFvRsConsIfprioStringNull() FvRsConsIfprioStringValue {
	return FvRsConsIfprioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewFvRsConsIfprioStringUnknown() FvRsConsIfprioStringValue {
	return FvRsConsIfprioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewFvRsConsIfprioStringValue(value string) FvRsConsIfprioStringValue {
	return FvRsConsIfprioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewFvRsConsIfprioStringPointerValue(value *string) FvRsConsIfprioStringValue {
	return FvRsConsIfprioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
