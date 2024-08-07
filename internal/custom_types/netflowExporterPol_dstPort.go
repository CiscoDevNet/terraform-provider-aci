package customTypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// NetflowExporterPoldstPort custom string type.

var _ basetypes.StringTypable = NetflowExporterPoldstPortStringType{}

type NetflowExporterPoldstPortStringType struct {
	basetypes.StringType
}

func (t NetflowExporterPoldstPortStringType) Equal(o attr.Type) bool {
	other, ok := o.(NetflowExporterPoldstPortStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t NetflowExporterPoldstPortStringType) String() string {
	return "NetflowExporterPoldstPortStringType"
}

func (t NetflowExporterPoldstPortStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := NetflowExporterPoldstPortStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t NetflowExporterPoldstPortStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t NetflowExporterPoldstPortStringType) ValueType(ctx context.Context) attr.Value {
	return NetflowExporterPoldstPortStringValue{}
}

// NetflowExporterPoldstPort custom string value.

var _ basetypes.StringValuableWithSemanticEquals = NetflowExporterPoldstPortStringValue{}

type NetflowExporterPoldstPortStringValue struct {
	basetypes.StringValue
}

func (v NetflowExporterPoldstPortStringValue) Equal(o attr.Value) bool {
	other, ok := o.(NetflowExporterPoldstPortStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v NetflowExporterPoldstPortStringValue) Type(ctx context.Context) attr.Type {
	return NetflowExporterPoldstPortStringType{}
}

func (v NetflowExporterPoldstPortStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(NetflowExporterPoldstPortStringValue)

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

	priorMappedValue := NetflowExporterPoldstPortValueMap(v.StringValue)

	newMappedValue := NetflowExporterPoldstPortValueMap(newValue.StringValue)

	return priorMappedValue.Equal(newMappedValue), diags
}

func NetflowExporterPoldstPortValueMap(value basetypes.StringValue) basetypes.StringValue {
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

func NewNetflowExporterPoldstPortStringNull() NetflowExporterPoldstPortStringValue {
	return NetflowExporterPoldstPortStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewNetflowExporterPoldstPortStringUnknown() NetflowExporterPoldstPortStringValue {
	return NetflowExporterPoldstPortStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewNetflowExporterPoldstPortStringValue(value string) NetflowExporterPoldstPortStringValue {
	return NetflowExporterPoldstPortStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewNetflowExporterPoldstPortStringPointerValue(value *string) NetflowExporterPoldstPortStringValue {
	return NetflowExporterPoldstPortStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
