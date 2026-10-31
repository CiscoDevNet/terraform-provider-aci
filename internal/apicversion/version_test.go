package apicversion

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type compareVersionsRangeInput struct {
	Version  string
	Ranges   string
	Operator string
}

type compareVersionsRangeExpected struct {
	Result bool
	Error  string
}

func TestCompareVersionsRange(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:     "test_inside_bounded_range",
			Input:    compareVersionsRangeInput{Version: "4.2(7g)", Ranges: "4.2(7f)-4.2(7w)", Operator: "inside"},
			Expected: compareVersionsRangeExpected{Result: true},
		},
		{
			Name:     "test_inside_multiple_ranges",
			Input:    compareVersionsRangeInput{Version: "5.2(1h)", Ranges: "4.2(7f)-4.2(7w),5.2(1g)-", Operator: "inside"},
			Expected: compareVersionsRangeExpected{Result: true},
		},
		{
			Name:     "test_outside_range",
			Input:    compareVersionsRangeInput{Version: "4.2(7z)", Ranges: "4.2(7f)-4.2(7w)", Operator: "outside"},
			Expected: compareVersionsRangeExpected{Result: true},
		},
		{
			Name:     "test_invalid_version",
			Input:    compareVersionsRangeInput{Version: "invalid", Ranges: "4.2(7f)-", Operator: "inside"},
			Expected: compareVersionsRangeExpected{Error: "Invalid APIC version provided: invalid"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(compareVersionsRangeInput)
			actual, err := CompareVersionsRange(input.Version, input.Ranges, input.Operator)
			expected := testCase.Expected.(compareVersionsRangeExpected)
			if expected.Error != "" {
				if err == nil || err.Error() != expected.Error {
					t.Fatalf("CompareVersionsRange returned error %v, expected %q", err, expected.Error)
				}
				return
			}
			if err != nil {
				t.Fatalf("CompareVersionsRange returned unexpected error: %v", err)
			}
			if actual != expected.Result {
				t.Fatalf("CompareVersionsRange returned %t, expected %t", actual, expected.Result)
			}
		})
	}
}
