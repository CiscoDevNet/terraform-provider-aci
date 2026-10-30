package customtypes

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

var _ basetypes.StringTypable = IPAddressType{}

// IPAddressType is a string type that compares canonical IPv4, IPv6, and CIDR values.
type IPAddressType struct {
	basetypes.StringType
}

func (t IPAddressType) Equal(other attr.Type) bool {
	_, ok := other.(IPAddressType)
	return ok
}

func (t IPAddressType) String() string {
	return "customtypes.IPAddressType"
}

func (t IPAddressType) ValueFromString(_ context.Context, value basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return IPAddressValue{StringValue: value}, nil
}

func (t IPAddressType) ValueFromTerraform(ctx context.Context, value tftypes.Value) (attr.Value, error) {
	attributeValue, err := t.StringType.ValueFromTerraform(ctx, value)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attributeValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attributeValue)
	}

	return IPAddressValue{StringValue: stringValue}, nil
}

func (t IPAddressType) ValueType(context.Context) attr.Value {
	return IPAddressValue{}
}

var _ basetypes.StringValuableWithSemanticEquals = IPAddressValue{}

// IPAddressValue preserves the configured value while comparing canonical IP notation.
type IPAddressValue struct {
	basetypes.StringValue
}

func (v IPAddressValue) Equal(other attr.Value) bool {
	otherValue, ok := other.(IPAddressValue)
	return ok && v.StringValue.Equal(otherValue.StringValue)
}

func (v IPAddressValue) Type(context.Context) attr.Type {
	return IPAddressType{}
}

func (v IPAddressValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	newValue, ok := newValuable.(IPAddressValue)
	if !ok {
		return false, unexpectedSemanticStringValueDiagnostics(v, newValuable)
	}

	return canonicalIPAddressValue(v.StringValue).Equal(canonicalIPAddressValue(newValue.StringValue)), nil
}

// CanonicalString returns the normalized IPv4, IPv6, or CIDR representation.
func (v IPAddressValue) CanonicalString() string {
	return canonicalIPAddressValue(v.StringValue).ValueString()
}

func canonicalIPAddressValue(value basetypes.StringValue) basetypes.StringValue {
	rawValue := strings.TrimSpace(value.ValueString())
	if ipAddress := net.ParseIP(rawValue); ipAddress != nil {
		return basetypes.NewStringValue(ipAddress.String())
	}

	ipAddress, network, err := net.ParseCIDR(rawValue)
	if err != nil {
		return basetypes.NewStringNull()
	}
	maskSize, _ := network.Mask.Size()
	return basetypes.NewStringValue(fmt.Sprintf("%s/%d", ipAddress.String(), maskSize))
}

func NewIPAddressNull() IPAddressValue {
	return IPAddressValue{StringValue: basetypes.NewStringNull()}
}

func NewIPAddressUnknown() IPAddressValue {
	return IPAddressValue{StringValue: basetypes.NewStringUnknown()}
}

func NewIPAddressValue(value string) IPAddressValue {
	return IPAddressValue{StringValue: basetypes.NewStringValue(value)}
}
