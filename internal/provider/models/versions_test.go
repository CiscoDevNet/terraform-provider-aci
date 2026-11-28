package models

import (
	"slices"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type classVersionsExpected struct {
	Versions string
	Found    bool
}

func TestClassVersions(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "test_top_level_class", Input: "fvTenant", Expected: classVersionsExpected{Versions: "1.0(1e)-", Found: true}},
		{Name: "test_child_only_class", Input: "cloudCertStore", Expected: classVersionsExpected{Versions: "4.1(1i)-", Found: true}},
		{Name: "test_unknown_class", Input: "missingClass", Expected: classVersionsExpected{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			versions, found := ClassVersions(testCase.Input.(string))
			expected := testCase.Expected.(classVersionsExpected)
			if versions != expected.Versions || found != expected.Found {
				t.Fatalf("ClassVersions returned (%q, %t), expected (%q, %t)", versions, found, expected.Versions, expected.Found)
			}
		})
	}
}

type supportedSubtreeClassesInput struct {
	APICVersion    string
	SubtreeClasses []string
}

func TestSupportedSubtreeClasses(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "test_filters_unsupported_unknown_and_duplicate_classes",
			Input: supportedSubtreeClassesInput{
				APICVersion:    "4.0(1a)",
				SubtreeClasses: []string{"cloudCertStore", "fvTenant", "fvTenant", "missingClass"},
			},
			Expected: []string{"fvTenant"},
		},
		{
			Name: "test_preserves_supported_input_order",
			Input: supportedSubtreeClassesInput{
				APICVersion:    "4.1(1i)",
				SubtreeClasses: []string{"cloudCertStore", "fvTenant"},
			},
			Expected: []string{"cloudCertStore", "fvTenant"},
		},
		{
			Name: "test_invalid_apic_version",
			Input: supportedSubtreeClassesInput{
				APICVersion:    "invalid",
				SubtreeClasses: []string{"fvTenant"},
			},
			Expected: []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(supportedSubtreeClassesInput)
			actual := supportedSubtreeClasses(input.APICVersion, input.SubtreeClasses)
			expected := testCase.Expected.([]string)
			if !slices.Equal(actual, expected) {
				t.Fatalf("supportedSubtreeClasses returned %v, expected %v", actual, expected)
			}
		})
	}
}
