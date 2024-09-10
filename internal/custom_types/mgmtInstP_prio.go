package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// MgmtInstPPrio custom string type.

var _ basetypes.StringTypable = MgmtInstPPrioStringType{}

type MgmtInstPPrioStringType struct {
	basetypes.StringType
}

func (t MgmtInstPPrioStringType) Equal(o attr.Type) bool {
	other, ok := o.(MgmtInstPPrioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t MgmtInstPPrioStringType) String() string {
	return "MgmtInstPPrioStringType"
}

func (t MgmtInstPPrioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := MgmtInstPPrioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t MgmtInstPPrioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t MgmtInstPPrioStringType) ValueType(ctx context.Context) attr.Value {
	return MgmtInstPPrioStringValue{}
}

// MgmtInstPPrio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = MgmtInstPPrioStringValue{}

type MgmtInstPPrioStringValue struct {
	basetypes.StringValue
}

func (v MgmtInstPPrioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(MgmtInstPPrioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v MgmtInstPPrioStringValue) Type(ctx context.Context) attr.Type {
	return MgmtInstPPrioStringType{}
}

func (v MgmtInstPPrioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(MgmtInstPPrioStringValue)

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

	priorMappedValue := MgmtInstPPrioValueMap(v.StringValue)

	newMappedValue := MgmtInstPPrioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func MgmtInstPPrioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewMgmtInstPPrioStringNull() MgmtInstPPrioStringValue {
	return MgmtInstPPrioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewMgmtInstPPrioStringUnknown() MgmtInstPPrioStringValue {
	return MgmtInstPPrioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewMgmtInstPPrioStringValue(value string) MgmtInstPPrioStringValue {
	return MgmtInstPPrioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewMgmtInstPPrioStringPointerValue(value *string) MgmtInstPPrioStringValue {
	return MgmtInstPPrioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
