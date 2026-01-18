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

// IPv6Address custom string type.
var _ basetypes.StringTypable = IPv6AddressStringType{}

type IPv6AddressStringType struct {
	basetypes.StringType
}

func (t IPv6AddressStringType) Equal(o attr.Type) bool {
	other, ok := o.(IPv6AddressStringType)
	if !ok {
		return false
	}
	return t.StringType.Equal(other.StringType)
}

func (t IPv6AddressStringType) String() string {
	return "IPv6AddressStringType"
}

func (t IPv6AddressStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := IPv6AddressStringValue{
		StringValue: in,
	}
	return value, nil
}

func (t IPv6AddressStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t IPv6AddressStringType) ValueType(ctx context.Context) attr.Value {
	return IPv6AddressStringValue{}
}

// IPv6Address custom string value.
var _ basetypes.StringValuableWithSemanticEquals = IPv6AddressStringValue{}

type IPv6AddressStringValue struct {
	basetypes.StringValue
}

func (v IPv6AddressStringValue) Equal(o attr.Value) bool {
	other, ok := o.(IPv6AddressStringValue)
	if !ok {
		return false
	}
	return v.StringValue.Equal(other.StringValue)
}

func (v IPv6AddressStringValue) Type(ctx context.Context) attr.Type {
	return IPv6AddressStringType{}
}

func (v IPv6AddressStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	newValue, ok := newValuable.(IPv6AddressStringValue)

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

	priorMappedValue := ParseIPv6AddressValue(v.StringValue)
	newMappedValue := ParseIPv6AddressValue(newValue.StringValue)
	return priorMappedValue.Equal(newMappedValue), diags
}

func (v IPv6AddressStringValue) NamedValueString() string {
	return ParseIPv6AddressValue(v.StringValue).ValueString()
}

func ParseIPv6AddressValue(value basetypes.StringValue) basetypes.StringValue {
	ipv6Address := net.ParseIP(strings.TrimSpace(value.ValueString()))
	if ipv6Address == nil {
		return basetypes.NewStringNull()
	}
	return basetypes.NewStringValue(ipv6Address.String())
}

func NewIPv6AddressStringNull() IPv6AddressStringValue {
	return IPv6AddressStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewIPv6AddressStringUnknown() IPv6AddressStringValue {
	return IPv6AddressStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewIPv6AddressStringValue(value string) IPv6AddressStringValue {
	return IPv6AddressStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewIPv6AddressStringPointerValue(value *string) IPv6AddressStringValue {
	return IPv6AddressStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
