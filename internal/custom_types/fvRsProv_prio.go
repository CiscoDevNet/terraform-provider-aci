package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// FvRsProvprio custom string type.

var _ basetypes.StringTypable = FvRsProvprioStringType{}

type FvRsProvprioStringType struct {
	basetypes.StringType
}

func (t FvRsProvprioStringType) Equal(o attr.Type) bool {
	other, ok := o.(FvRsProvprioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t FvRsProvprioStringType) String() string {
	return "FvRsProvprioStringType"
}

func (t FvRsProvprioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := FvRsProvprioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t FvRsProvprioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t FvRsProvprioStringType) ValueType(ctx context.Context) attr.Value {
	return FvRsProvprioStringValue{}
}

// FvRsProvprio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = FvRsProvprioStringValue{}

type FvRsProvprioStringValue struct {
	basetypes.StringValue
}

func (v FvRsProvprioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(FvRsProvprioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v FvRsProvprioStringValue) Type(ctx context.Context) attr.Type {
	return FvRsProvprioStringType{}
}

func (v FvRsProvprioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(FvRsProvprioStringValue)

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

	priorMappedValue := FvRsProvprioValueMap(v.StringValue)

	newMappedValue := FvRsProvprioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func FvRsProvprioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewFvRsProvprioStringNull() FvRsProvprioStringValue {
	return FvRsProvprioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewFvRsProvprioStringUnknown() FvRsProvprioStringValue {
	return FvRsProvprioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewFvRsProvprioStringValue(value string) FvRsProvprioStringValue {
	return FvRsProvprioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewFvRsProvprioStringPointerValue(value *string) FvRsProvprioStringValue {
	return FvRsProvprioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
