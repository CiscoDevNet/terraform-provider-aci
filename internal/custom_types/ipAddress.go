package customTypes

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// IPAddress custom string type.
var _ basetypes.StringTypable = IPAddressStringType{}

type IPAddressStringType struct {
	basetypes.StringType
}

func (t IPAddressStringType) Equal(o attr.Type) bool {
	other, ok := o.(IPAddressStringType)
	if !ok {
		return false
	}
	return t.StringType.Equal(other.StringType)
}

func (t IPAddressStringType) String() string {
	return "IPAddressStringType"
}

func (t IPAddressStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := IPAddressStringValue{
		StringValue: in,
	}
	return value, nil
}

func (t IPAddressStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t IPAddressStringType) ValueType(ctx context.Context) attr.Value {
	return IPAddressStringValue{}
}

// IPAddress custom string value.
var _ basetypes.StringValuableWithSemanticEquals = IPAddressStringValue{}

type IPAddressStringValue struct {
	basetypes.StringValue
}

func (v IPAddressStringValue) Equal(o attr.Value) bool {
	other, ok := o.(IPAddressStringValue)
	if !ok {
		return false
	}
	return v.StringValue.Equal(other.StringValue)
}

func (v IPAddressStringValue) Type(ctx context.Context) attr.Type {
	return IPAddressStringType{}
}

func (v IPAddressStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	newValue, ok := newValuable.(IPAddressStringValue)

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

	priorMappedValue := ParseIPAddressValue(v.StringValue)
	newMappedValue := ParseIPAddressValue(newValue.StringValue)
	return priorMappedValue.Equal(newMappedValue), diags
}

func (v IPAddressStringValue) NamedValueString() string {
	return ParseIPAddressValue(v.StringValue).ValueString()
}

func ParseIPAddressValue(value basetypes.StringValue) basetypes.StringValue {
	ipAddress := net.ParseIP(strings.TrimSpace(value.ValueString()))
	if ipAddress == nil {
		ip, ipnet, err := net.ParseCIDR(strings.TrimSpace(value.ValueString()))
		address := ip.String()
		ones, _ := ipnet.Mask.Size()
		if err != nil {
			return basetypes.NewStringNull()
		}
		return basetypes.NewStringValue(fmt.Sprintf("%s/%d", address, ones))
	}
	return basetypes.NewStringValue(ipAddress.String())
}

func NewIPAddressStringNull() IPAddressStringValue {
	return IPAddressStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewIPAddressStringUnknown() IPAddressStringValue {
	return IPAddressStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewIPAddressStringValue(value string) IPAddressStringValue {
	return IPAddressStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewIPAddressStringPointerValue(value *string) IPAddressStringValue {
	return IPAddressStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
