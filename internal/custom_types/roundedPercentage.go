package customTypes

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// RoundedPercentage custom string type.

var _ basetypes.StringTypable = RoundedPercentageStringType{}

type RoundedPercentageStringType struct {
	basetypes.StringType
}

func (t RoundedPercentageStringType) Equal(o attr.Type) bool {
	other, ok := o.(RoundedPercentageStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t RoundedPercentageStringType) String() string {
	return "RoundedPercentageStringType"
}

func (t RoundedPercentageStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := RoundedPercentageStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t RoundedPercentageStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t RoundedPercentageStringType) ValueType(ctx context.Context) attr.Value {
	return RoundedPercentageStringValue{}
}

// RoundedPercentage custom string value.

var _ basetypes.StringValuableWithSemanticEquals = RoundedPercentageStringValue{}

type RoundedPercentageStringValue struct {
	basetypes.StringValue
}

func (v RoundedPercentageStringValue) Equal(o attr.Value) bool {
	other, ok := o.(RoundedPercentageStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v RoundedPercentageStringValue) Type(ctx context.Context) attr.Type {
	return RoundedPercentageStringType{}
}

func (v RoundedPercentageStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(RoundedPercentageStringValue)

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

	priorMappedValue, err := RoundPercentageValue(v.StringValue)

	if err != nil {
		addRoudingError(err, &diags)
	}

	newMappedValue, err := RoundPercentageValue(newValue.StringValue)

	if err != nil {
		addRoudingError(err, &diags)
	}

	return priorMappedValue.Equal(newMappedValue), diags
}

func addRoudingError(err error, diags *diag.Diagnostics) {

	if err != nil {
		diags.AddError(
			"Semantic Equality Check Error",
			"Rounding of percentage value failed while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Error: "+fmt.Sprintf("%v", err),
		)
	}
}

func (v RoundedPercentageStringValue) NamedValueString() string {
	roundedPercentage, _ := RoundPercentageValue(v.StringValue)
	return roundedPercentage.ValueString()
}

func RoundPercentageValue(value basetypes.StringValue) (basetypes.StringValue, error) {

	percentage, err := strconv.ParseFloat(strings.TrimSpace(value.ValueString()), 64)

	if err != nil {
		return basetypes.NewStringNull(), err
	}

	return basetypes.NewStringValue(strconv.FormatFloat(percentage, 'f', 2, 64)), err
}

func NewRoundedPercentageStringNull() RoundedPercentageStringValue {
	return RoundedPercentageStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewRoundedPercentageStringUnknown() RoundedPercentageStringValue {
	return RoundedPercentageStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewRoundedPercentageStringValue(value string) RoundedPercentageStringValue {
	return RoundedPercentageStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewRoundedPercentageStringPointerValue(value *string) RoundedPercentageStringValue {
	return RoundedPercentageStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
