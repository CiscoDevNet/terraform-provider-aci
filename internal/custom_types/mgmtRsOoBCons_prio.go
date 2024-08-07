package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// MgmtRsOoBConsprio custom string type.

var _ basetypes.StringTypable = MgmtRsOoBConsprioStringType{}

type MgmtRsOoBConsprioStringType struct {
	basetypes.StringType
}

func (t MgmtRsOoBConsprioStringType) Equal(o attr.Type) bool {
	other, ok := o.(MgmtRsOoBConsprioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t MgmtRsOoBConsprioStringType) String() string {
	return "MgmtRsOoBConsprioStringType"
}

func (t MgmtRsOoBConsprioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := MgmtRsOoBConsprioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t MgmtRsOoBConsprioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t MgmtRsOoBConsprioStringType) ValueType(ctx context.Context) attr.Value {
	return MgmtRsOoBConsprioStringValue{}
}

// MgmtRsOoBConsprio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = MgmtRsOoBConsprioStringValue{}

type MgmtRsOoBConsprioStringValue struct {
	basetypes.StringValue
}

func (v MgmtRsOoBConsprioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(MgmtRsOoBConsprioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v MgmtRsOoBConsprioStringValue) Type(ctx context.Context) attr.Type {
	return MgmtRsOoBConsprioStringType{}
}

func (v MgmtRsOoBConsprioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(MgmtRsOoBConsprioStringValue)

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

	priorMappedValue := MgmtRsOoBConsprioValueMap(v.StringValue)

	newMappedValue := MgmtRsOoBConsprioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func MgmtRsOoBConsprioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewMgmtRsOoBConsprioStringNull() MgmtRsOoBConsprioStringValue {
	return MgmtRsOoBConsprioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewMgmtRsOoBConsprioStringUnknown() MgmtRsOoBConsprioStringValue {
	return MgmtRsOoBConsprioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewMgmtRsOoBConsprioStringValue(value string) MgmtRsOoBConsprioStringValue {
	return MgmtRsOoBConsprioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewMgmtRsOoBConsprioStringPointerValue(value *string) MgmtRsOoBConsprioStringValue {
	return MgmtRsOoBConsprioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
