package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// MgmtRsOoBConsPrio custom string type.

var _ basetypes.StringTypable = MgmtRsOoBConsPrioStringType{}

type MgmtRsOoBConsPrioStringType struct {
	basetypes.StringType
}

func (t MgmtRsOoBConsPrioStringType) Equal(o attr.Type) bool {
	other, ok := o.(MgmtRsOoBConsPrioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t MgmtRsOoBConsPrioStringType) String() string {
	return "MgmtRsOoBConsPrioStringType"
}

func (t MgmtRsOoBConsPrioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := MgmtRsOoBConsPrioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t MgmtRsOoBConsPrioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t MgmtRsOoBConsPrioStringType) ValueType(ctx context.Context) attr.Value {
	return MgmtRsOoBConsPrioStringValue{}
}

// MgmtRsOoBConsPrio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = MgmtRsOoBConsPrioStringValue{}

type MgmtRsOoBConsPrioStringValue struct {
	basetypes.StringValue
}

func (v MgmtRsOoBConsPrioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(MgmtRsOoBConsPrioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v MgmtRsOoBConsPrioStringValue) Type(ctx context.Context) attr.Type {
	return MgmtRsOoBConsPrioStringType{}
}

func (v MgmtRsOoBConsPrioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(MgmtRsOoBConsPrioStringValue)

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

	priorMappedValue := MgmtRsOoBConsPrioValueMap(v.StringValue)

	newMappedValue := MgmtRsOoBConsPrioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func MgmtRsOoBConsPrioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewMgmtRsOoBConsPrioStringNull() MgmtRsOoBConsPrioStringValue {
	return MgmtRsOoBConsPrioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewMgmtRsOoBConsPrioStringUnknown() MgmtRsOoBConsPrioStringValue {
	return MgmtRsOoBConsPrioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewMgmtRsOoBConsPrioStringValue(value string) MgmtRsOoBConsPrioStringValue {
	return MgmtRsOoBConsPrioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewMgmtRsOoBConsPrioStringPointerValue(value *string) MgmtRsOoBConsPrioStringValue {
	return MgmtRsOoBConsPrioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
