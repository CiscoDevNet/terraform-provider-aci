package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// PoeIfPolMax custom string type.

var _ basetypes.StringTypable = PoeIfPolMaxStringType{}

type PoeIfPolMaxStringType struct {
	basetypes.StringType
}

func (t PoeIfPolMaxStringType) Equal(o attr.Type) bool {
	other, ok := o.(PoeIfPolMaxStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t PoeIfPolMaxStringType) String() string {
	return "PoeIfPolMaxStringType"
}

func (t PoeIfPolMaxStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := PoeIfPolMaxStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t PoeIfPolMaxStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t PoeIfPolMaxStringType) ValueType(ctx context.Context) attr.Value {
	return PoeIfPolMaxStringValue{}
}

// PoeIfPolMax custom string value.

var _ basetypes.StringValuableWithSemanticEquals = PoeIfPolMaxStringValue{}

type PoeIfPolMaxStringValue struct {
	basetypes.StringValue
}

func (v PoeIfPolMaxStringValue) Equal(o attr.Value) bool {
	other, ok := o.(PoeIfPolMaxStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v PoeIfPolMaxStringValue) Type(ctx context.Context) attr.Type {
	return PoeIfPolMaxStringType{}
}

func (v PoeIfPolMaxStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(PoeIfPolMaxStringValue)

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

	priorMappedValue := PoeIfPolMaxValueMap(v.StringValue)

	newMappedValue := PoeIfPolMaxValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v PoeIfPolMaxStringValue) NamedValueString() string {
	return PoeIfPolMaxValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func PoeIfPolMaxValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"15400": "15400",
		"30000": "30000",
		"4000":  "4000",
		"45000": "45000",
		"60000": "60000",
		"7000":  "7000",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewPoeIfPolMaxStringNull() PoeIfPolMaxStringValue {
	return PoeIfPolMaxStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewPoeIfPolMaxStringUnknown() PoeIfPolMaxStringValue {
	return PoeIfPolMaxStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewPoeIfPolMaxStringValue(value string) PoeIfPolMaxStringValue {
	return PoeIfPolMaxStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewPoeIfPolMaxStringPointerValue(value *string) PoeIfPolMaxStringValue {
	return PoeIfPolMaxStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
