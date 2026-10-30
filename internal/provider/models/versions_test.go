package models

import (
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
