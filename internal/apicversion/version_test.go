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

func TestNormalizedVersionRanges(t *testing.T) {
	type input struct {
		Version string
		Ranges  []VersionRange
	}
	intervals := []VersionRange{{Min: "4.2(0)", Max: "5.2(10)", MaxExclusive: true}, {Min: "6.0(0)", Max: "6.2(1)", MaxExclusive: true}}
	testCases := []test.TestCase{
		{Name: "before_start", Input: input{"4.1(1)", intervals}, Expected: false},
		{Name: "inclusive_start", Input: input{"4.2(0)", intervals}, Expected: true},
		{Name: "before_backport_cutoff", Input: input{"5.2(9z)", intervals}, Expected: true},
		{Name: "at_backport_cutoff", Input: input{"5.2(10)", intervals}, Expected: false},
		{Name: "inside_gap", Input: input{"5.2(11)", intervals}, Expected: false},
		{Name: "next_interval", Input: input{"6.0(0)", intervals}, Expected: true},
		{Name: "deprecated_but_available", Input: input{"6.1(5e)", intervals}, Expected: true},
		{Name: "at_hidden_cutoff", Input: input{"6.2(1)", intervals}, Expected: false},
		{Name: "after_hidden_cutoff", Input: input{"6.2(2e)", intervals}, Expected: false},
		{Name: "empty_intersection", Input: input{"6.2(2e)", []VersionRange{}}, Expected: false},
		{Name: "unchanged_inclusive_end", Input: input{"6.2(1)", []VersionRange{{Min: "4.2(0)", Max: "6.2(1)"}}}, Expected: true},
		{Name: "exclusive_start", Input: input{"5.2(10)", []VersionRange{{Min: "5.2(10)", MinExclusive: true}}}, Expected: false},
		{Name: "after_exclusive_start", Input: input{"5.2(10a)", []VersionRange{{Min: "5.2(10)", MinExclusive: true}}}, Expected: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			value := testCase.Input.(input)
			actual, err := IsVersionInRanges(value.Version, value.Ranges)
			if err != nil {
				t.Fatal(err)
			}
			if actual != testCase.Expected {
				t.Fatalf("got %v, want %v", actual, testCase.Expected)
			}
		})
	}
}
