package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// NetflowExporterPolDstPort custom string type.

var _ basetypes.StringTypable = NetflowExporterPolDstPortStringType{}

type NetflowExporterPolDstPortStringType struct {
	basetypes.StringType
}

func (t NetflowExporterPolDstPortStringType) Equal(o attr.Type) bool {
	other, ok := o.(NetflowExporterPolDstPortStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t NetflowExporterPolDstPortStringType) String() string {
	return "NetflowExporterPolDstPortStringType"
}

func (t NetflowExporterPolDstPortStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := NetflowExporterPolDstPortStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t NetflowExporterPolDstPortStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t NetflowExporterPolDstPortStringType) ValueType(ctx context.Context) attr.Value {
	return NetflowExporterPolDstPortStringValue{}
}

// NetflowExporterPolDstPort custom string value.

var _ basetypes.StringValuableWithSemanticEquals = NetflowExporterPolDstPortStringValue{}

type NetflowExporterPolDstPortStringValue struct {
	basetypes.StringValue
}

func (v NetflowExporterPolDstPortStringValue) Equal(o attr.Value) bool {
	other, ok := o.(NetflowExporterPolDstPortStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v NetflowExporterPolDstPortStringValue) Type(ctx context.Context) attr.Type {
	return NetflowExporterPolDstPortStringType{}
}

func (v NetflowExporterPolDstPortStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(NetflowExporterPolDstPortStringValue)

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

	priorMappedValue := NetflowExporterPolDstPortValueMap(v.StringValue)

	newMappedValue := NetflowExporterPolDstPortValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func (v NetflowExporterPolDstPortStringValue) NamedValueString() string {
	return NetflowExporterPolDstPortValueMap(basetypes.NewStringValue(v.ValueString())).ValueString()
}

func NetflowExporterPolDstPortValueMap(value basetypes.StringValue) basetypes.StringValue {
	matchMap := map[string]string{
		"0":   "unspecified",
		"110": "pop3",
		"20":  "ftpData",
		"22":  "ssh",
		"25":  "smtp",
		"443": "https",
		"53":  "dns",
		"554": "rtsp",
		"80":  "http",
	}

	if val, ok := matchMap[value.ValueString()]; ok {
		return basetypes.NewStringValue(val)
	}

	return value
}

func NewNetflowExporterPolDstPortStringNull() NetflowExporterPolDstPortStringValue {
	return NetflowExporterPolDstPortStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewNetflowExporterPolDstPortStringUnknown() NetflowExporterPolDstPortStringValue {
	return NetflowExporterPolDstPortStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewNetflowExporterPolDstPortStringValue(value string) NetflowExporterPolDstPortStringValue {
	return NetflowExporterPolDstPortStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewNetflowExporterPolDstPortStringPointerValue(value *string) NetflowExporterPolDstPortStringValue {
	return NetflowExporterPolDstPortStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
