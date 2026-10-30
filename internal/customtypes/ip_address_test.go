package customtypes

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

func TestIPAddressValueCanonicalString(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "test_ipv4", Input: "192.0.2.1", Expected: "192.0.2.1"},
		{Name: "test_ipv6", Input: "2001:0db8:0000:0000:0000:0000:0000:0001", Expected: "2001:db8::1"},
		{Name: "test_ipv6_cidr", Input: "2001:0db8::1/64", Expected: "2001:db8::1/64"},
		{Name: "test_whitespace", Input: " 2001:db8::1 ", Expected: "2001:db8::1"},
		{Name: "test_invalid", Input: "invalid", Expected: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := NewIPAddressValue(testCase.Input.(string)).CanonicalString()
			if actual != testCase.Expected.(string) {
				t.Fatalf("canonical address %q, expected %q", actual, testCase.Expected)
			}
		})
	}
}

func TestIPAddressValueSemanticEquality(t *testing.T) {
	equal, diagnostics := NewIPAddressValue("2001:0db8::1").StringSemanticEquals(
		context.Background(),
		NewIPAddressValue("2001:db8:0:0:0:0:0:1"),
	)
	if diagnostics.HasError() || !equal {
		t.Fatalf("equivalent IP addresses were not equal: equal=%t diagnostics=%v", equal, diagnostics)
	}

	equal, diagnostics = NewIPAddressValue("192.0.2.1").StringSemanticEquals(
		context.Background(),
		NewIPAddressValue("192.0.2.2"),
	)
	if diagnostics.HasError() || equal {
		t.Fatalf("different IP addresses comparison: equal=%t diagnostics=%v", equal, diagnostics)
	}
}

func TestIPAddressValueStatesAndTypeMismatch(t *testing.T) {
	if !NewIPAddressNull().IsNull() {
		t.Fatal("NewIPAddressNull did not return a null value")
	}
	if !NewIPAddressUnknown().IsUnknown() {
		t.Fatal("NewIPAddressUnknown did not return an unknown value")
	}

	_, diagnostics := NewIPAddressValue("192.0.2.1").StringSemanticEquals(
		context.Background(),
		basetypes.NewStringValue("192.0.2.1"),
	)
	if !diagnostics.HasError() {
		t.Fatal("different value type did not produce an error diagnostic")
	}
}
