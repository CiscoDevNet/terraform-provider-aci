package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// MgmtInstPprio custom string type.

var _ basetypes.StringTypable = MgmtInstPprioStringType{}

type MgmtInstPprioStringType struct {
	basetypes.StringType
}

func (t MgmtInstPprioStringType) Equal(o attr.Type) bool {
	other, ok := o.(MgmtInstPprioStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t MgmtInstPprioStringType) String() string {
	return "MgmtInstPprioStringType"
}

func (t MgmtInstPprioStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := MgmtInstPprioStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t MgmtInstPprioStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t MgmtInstPprioStringType) ValueType(ctx context.Context) attr.Value {
	return MgmtInstPprioStringValue{}
}

// MgmtInstPprio custom string value.

var _ basetypes.StringValuableWithSemanticEquals = MgmtInstPprioStringValue{}

type MgmtInstPprioStringValue struct {
	basetypes.StringValue
}

func (v MgmtInstPprioStringValue) Equal(o attr.Value) bool {
	other, ok := o.(MgmtInstPprioStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v MgmtInstPprioStringValue) Type(ctx context.Context) attr.Type {
	return MgmtInstPprioStringType{}
}

func (v MgmtInstPprioStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(MgmtInstPprioStringValue)

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

	priorMappedValue := MgmtInstPprioValueMap(v.StringValue)

	newMappedValue := MgmtInstPprioValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func MgmtInstPprioValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewMgmtInstPprioStringNull() MgmtInstPprioStringValue {
	return MgmtInstPprioStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewMgmtInstPprioStringUnknown() MgmtInstPprioStringValue {
	return MgmtInstPprioStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewMgmtInstPprioStringValue(value string) MgmtInstPprioStringValue {
	return MgmtInstPprioStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewMgmtInstPprioStringPointerValue(value *string) MgmtInstPprioStringValue {
	return MgmtInstPprioStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
