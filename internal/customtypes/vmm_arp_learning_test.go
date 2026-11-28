package customtypes

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

func TestVMMArpLearningValueCanonicalString(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "test_empty", Input: "", Expected: "disabled"},
		{Name: "test_disabled", Input: "disabled", Expected: "disabled"},
		{Name: "test_enabled", Input: "enabled", Expected: "enabled"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := NewVMMArpLearningValue(testCase.Input.(string)).CanonicalString()
			if actual != testCase.Expected.(string) {
				t.Fatalf("canonical value %q, expected %q", actual, testCase.Expected)
			}
		})
	}
}

func TestVMMArpLearningValueSemanticEquality(t *testing.T) {
	equal, diagnostics := NewVMMArpLearningValue("").StringSemanticEquals(
		context.Background(),
		NewVMMArpLearningValue("disabled"),
	)
	if diagnostics.HasError() || !equal {
		t.Fatalf("empty and disabled were not equal: equal=%t diagnostics=%v", equal, diagnostics)
	}

	equal, diagnostics = NewVMMArpLearningValue("enabled").StringSemanticEquals(
		context.Background(),
		NewVMMArpLearningValue("disabled"),
	)
	if diagnostics.HasError() || equal {
		t.Fatalf("enabled and disabled comparison: equal=%t diagnostics=%v", equal, diagnostics)
	}
}

func TestVMMArpLearningValueStatesAndTypeMismatch(t *testing.T) {
	if !NewVMMArpLearningNull().IsNull() {
		t.Fatal("NewVMMArpLearningNull did not return a null value")
	}
	if !NewVMMArpLearningUnknown().IsUnknown() {
		t.Fatal("NewVMMArpLearningUnknown did not return an unknown value")
	}

	_, diagnostics := NewVMMArpLearningValue("disabled").StringSemanticEquals(
		context.Background(),
		basetypes.NewStringValue("disabled"),
	)
	if !diagnostics.HasError() {
		t.Fatal("different value type did not produce an error diagnostic")
	}
}
