package data

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/provider"
	"github.com/stretchr/testify/assert"
)

type newVersionRangeExpected struct {
	Raw             string
	Min             *provider.Version
	Max             *provider.Version
	IsSingleVersion bool
	String          string
	Error           bool
	ErrorMsg        string
}

func TestNewVersionRange(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:  "test_single_version_with_tag",
			Input: "4.2(7f)",
			Expected: newVersionRangeExpected{
				Raw:             "4.2(7f)",
				Min:             &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				Max:             &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				IsSingleVersion: true,
				String:          "4.2(7f)",
			},
		},
		{
			Name:  "test_single_version_without_tag",
			Input: "4.2(7)",
			Expected: newVersionRangeExpected{
				Raw:             "4.2(7)",
				Min:             &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: 0},
				Max:             &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: 0},
				IsSingleVersion: true,
				String:          "4.2(7)",
			},
		},
		{
			Name:  "test_bounded_range",
			Input: "4.2(7f)-4.2(7w)",
			Expected: newVersionRangeExpected{
				Raw:             "4.2(7f)-4.2(7w)",
				Min:             &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				Max:             &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('w')},
				IsSingleVersion: false,
				String:          "4.2(7f) to 4.2(7w)",
			},
		},
		{
			Name:  "test_unbounded_upper",
			Input: "4.2(7f)-",
			Expected: newVersionRangeExpected{
				Raw:             "4.2(7f)-",
				Min:             &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				Max:             nil,
				IsSingleVersion: false,
				String:          "4.2(7f) and later",
			},
		},
		{
			Name:  "test_unbounded_lower",
			Input: "-4.2(7w)",
			Expected: newVersionRangeExpected{
				Raw:             "-4.2(7w)",
				Min:             nil,
				Max:             &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('w')},
				IsSingleVersion: false,
				String:          "up to 4.2(7w)",
			},
		},
		{
			Name:  "test_error_invalid_version",
			Input: "invalid",
			Expected: newVersionRangeExpected{
				Error:    true,
				ErrorMsg: "invalid version 'invalid': unknown",
			},
		},
		{
			Name:  "test_error_invalid_min_version",
			Input: "invalid-4.2(7w)",
			Expected: newVersionRangeExpected{
				Error:    true,
				ErrorMsg: "invalid minimum version 'invalid': unknown",
			},
		},
		{
			Name:  "test_error_invalid_max_version",
			Input: "4.2(7f)-invalid",
			Expected: newVersionRangeExpected{
				Error:    true,
				ErrorMsg: "invalid maximum version 'invalid': unknown",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			expected := testCase.Expected.(newVersionRangeExpected)
			versionRange, err := NewVersionRange(testCase.Input.(string))

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, expected.Raw, versionRange.Raw(), test.MessageEqual(expected.Raw, versionRange.Raw(), testCase.Name))
				assert.Equal(t, expected.Min, versionRange.Min(), test.MessageEqual(expected.Min, versionRange.Min(), testCase.Name))
				assert.Equal(t, expected.Max, versionRange.Max(), test.MessageEqual(expected.Max, versionRange.Max(), testCase.Name))
				assert.Equal(t, expected.IsSingleVersion, versionRange.IsSingleVersion(), test.MessageEqual(expected.IsSingleVersion, versionRange.IsSingleVersion(), testCase.Name))
				assert.Equal(t, expected.String, versionRange.String(), test.MessageEqual(expected.String, versionRange.String(), testCase.Name))
			}
		})
	}
}

type newVersionsExpected struct {
	Raw         string
	RangeCount  int
	String      string
	RangeValues []newVersionRangeExpected
	Error       bool
	ErrorMsg    string
}

func TestNewVersions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:  "test_single_version",
			Input: "4.2(7f)",
			Expected: newVersionsExpected{
				Raw:        "4.2(7f)",
				RangeCount: 1,
				String:     "4.2(7f)",
				RangeValues: []newVersionRangeExpected{
					{
						Raw:             "4.2(7f)",
						Min:             &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
						Max:             &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
						IsSingleVersion: true,
						String:          "4.2(7f)",
					},
				},
			},
		},
		{
			Name:  "test_unbounded_range",
			Input: "4.2(7f)-",
			Expected: newVersionsExpected{
				Raw:        "4.2(7f)-",
				RangeCount: 1,
				String:     "4.2(7f) and later",
				RangeValues: []newVersionRangeExpected{
					{
						Raw:             "4.2(7f)-",
						Min:             &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
						Max:             nil,
						IsSingleVersion: false,
						String:          "4.2(7f) and later",
					},
				},
			},
		},
		{
			Name:  "test_multiple_ranges",
			Input: "3.2(10e)-3.2(10g),4.2(7f)-",
			Expected: newVersionsExpected{
				Raw:        "3.2(10e)-3.2(10g),4.2(7f)-",
				RangeCount: 2,
				String:     "3.2(10e) to 3.2(10g), 4.2(7f) and later",
				RangeValues: []newVersionRangeExpected{
					{
						Raw:             "3.2(10e)-3.2(10g)",
						Min:             &provider.Version{Major: 3, Minor: 2, Patch: 10, Tag: int('e')},
						Max:             &provider.Version{Major: 3, Minor: 2, Patch: 10, Tag: int('g')},
						IsSingleVersion: false,
						String:          "3.2(10e) to 3.2(10g)",
					},
					{
						Raw:             "4.2(7f)-",
						Min:             &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
						Max:             nil,
						IsSingleVersion: false,
						String:          "4.2(7f) and later",
					},
				},
			},
		},
		{
			Name:  "test_multiple_ranges_sorted",
			Input: "5.2(1g)-,3.2(10e)-3.2(10g)",
			Expected: newVersionsExpected{
				Raw:        "5.2(1g)-,3.2(10e)-3.2(10g)",
				RangeCount: 2,
				String:     "3.2(10e) to 3.2(10g), 5.2(1g) and later",
				RangeValues: []newVersionRangeExpected{
					{
						Raw:             "3.2(10e)-3.2(10g)",
						Min:             &provider.Version{Major: 3, Minor: 2, Patch: 10, Tag: int('e')},
						Max:             &provider.Version{Major: 3, Minor: 2, Patch: 10, Tag: int('g')},
						IsSingleVersion: false,
						String:          "3.2(10e) to 3.2(10g)",
					},
					{
						Raw:             "5.2(1g)-",
						Min:             &provider.Version{Major: 5, Minor: 2, Patch: 1, Tag: int('g')},
						Max:             nil,
						IsSingleVersion: false,
						String:          "5.2(1g) and later",
					},
				},
			},
		},
		{
			Name:  "test_three_ranges_sorted",
			Input: "5.0(1a)-,3.0(1a)-3.0(1z),4.0(1a)-4.0(1z)",
			Expected: newVersionsExpected{
				Raw:        "5.0(1a)-,3.0(1a)-3.0(1z),4.0(1a)-4.0(1z)",
				RangeCount: 3,
				String:     "3.0(1a) to 3.0(1z), 4.0(1a) to 4.0(1z), 5.0(1a) and later",
				RangeValues: []newVersionRangeExpected{
					{
						Raw:             "3.0(1a)-3.0(1z)",
						Min:             &provider.Version{Major: 3, Minor: 0, Patch: 1, Tag: int('a')},
						Max:             &provider.Version{Major: 3, Minor: 0, Patch: 1, Tag: int('z')},
						IsSingleVersion: false,
						String:          "3.0(1a) to 3.0(1z)",
					},
					{
						Raw:             "4.0(1a)-4.0(1z)",
						Min:             &provider.Version{Major: 4, Minor: 0, Patch: 1, Tag: int('a')},
						Max:             &provider.Version{Major: 4, Minor: 0, Patch: 1, Tag: int('z')},
						IsSingleVersion: false,
						String:          "4.0(1a) to 4.0(1z)",
					},
					{
						Raw:             "5.0(1a)-",
						Min:             &provider.Version{Major: 5, Minor: 0, Patch: 1, Tag: int('a')},
						Max:             nil,
						IsSingleVersion: false,
						String:          "5.0(1a) and later",
					},
				},
			},
		},
		{
			Name:  "test_error_invalid_version",
			Input: "invalid",
			Expected: newVersionsExpected{
				Error:    true,
				ErrorMsg: "invalid version 'invalid': unknown",
			},
		},
		{
			Name:  "test_error_invalid_version_in_multiple",
			Input: "4.2(7f)-,invalid",
			Expected: newVersionsExpected{
				Error:    true,
				ErrorMsg: "invalid version 'invalid': unknown",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			expected := testCase.Expected.(newVersionsExpected)
			versions, err := NewVersions(testCase.Input.(string))

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, expected.Raw, versions.Raw(), test.MessageEqual(expected.Raw, versions.Raw(), testCase.Name))
				assert.Equal(t, expected.RangeCount, len(versions.Ranges()), test.MessageEqual(expected.RangeCount, len(versions.Ranges()), testCase.Name))
				assert.Equal(t, expected.String, versions.String(), test.MessageEqual(expected.String, versions.String(), testCase.Name))

				// Verify each range
				for i, expectedRange := range expected.RangeValues {
					actualRange := versions.Ranges()[i]
					assert.Equal(t, expectedRange.Raw, actualRange.Raw(), test.MessageEqual(expectedRange.Raw, actualRange.Raw(), testCase.Name))
					assert.Equal(t, expectedRange.Min, actualRange.Min(), test.MessageEqual(expectedRange.Min, actualRange.Min(), testCase.Name))
					assert.Equal(t, expectedRange.Max, actualRange.Max(), test.MessageEqual(expectedRange.Max, actualRange.Max(), testCase.Name))
					assert.Equal(t, expectedRange.IsSingleVersion, actualRange.IsSingleVersion(), test.MessageEqual(expectedRange.IsSingleVersion, actualRange.IsSingleVersion(), testCase.Name))
					assert.Equal(t, expectedRange.String, actualRange.String(), test.MessageEqual(expectedRange.String, actualRange.String(), testCase.Name))
				}
			}
		})
	}
}

type sortVersionsExpected struct {
	Result int
}

type sortVersionsInput struct {
	A          *provider.Version
	B          *provider.Version
	IsMinBound bool
}

func TestSortVersions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_both_nil_min_bound",
			Input: sortVersionsInput{
				A:          nil,
				B:          nil,
				IsMinBound: true,
			},
			Expected: sortVersionsExpected{Result: 0},
		},
		{
			Name: "test_both_nil_max_bound",
			Input: sortVersionsInput{
				A:          nil,
				B:          nil,
				IsMinBound: false,
			},
			Expected: sortVersionsExpected{Result: 0},
		},
		{
			Name: "test_a_nil_min_bound",
			Input: sortVersionsInput{
				A:          nil,
				B:          &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				IsMinBound: true,
			},
			Expected: sortVersionsExpected{Result: -1},
		},
		{
			Name: "test_a_nil_max_bound",
			Input: sortVersionsInput{
				A:          nil,
				B:          &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				IsMinBound: false,
			},
			Expected: sortVersionsExpected{Result: 1},
		},
		{
			Name: "test_b_nil_min_bound",
			Input: sortVersionsInput{
				A:          &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				B:          nil,
				IsMinBound: true,
			},
			Expected: sortVersionsExpected{Result: 1},
		},
		{
			Name: "test_b_nil_max_bound",
			Input: sortVersionsInput{
				A:          &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				B:          nil,
				IsMinBound: false,
			},
			Expected: sortVersionsExpected{Result: -1},
		},
		{
			Name: "test_equal_versions",
			Input: sortVersionsInput{
				A:          &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				B:          &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				IsMinBound: true,
			},
			Expected: sortVersionsExpected{Result: 0},
		},
		{
			Name: "test_a_less_than_b_major",
			Input: sortVersionsInput{
				A:          &provider.Version{Major: 3, Minor: 2, Patch: 7, Tag: int('f')},
				B:          &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				IsMinBound: true,
			},
			Expected: sortVersionsExpected{Result: -1},
		},
		{
			Name: "test_a_greater_than_b_major",
			Input: sortVersionsInput{
				A:          &provider.Version{Major: 5, Minor: 2, Patch: 7, Tag: int('f')},
				B:          &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				IsMinBound: true,
			},
			Expected: sortVersionsExpected{Result: 1},
		},
		{
			Name: "test_a_less_than_b_minor",
			Input: sortVersionsInput{
				A:          &provider.Version{Major: 4, Minor: 1, Patch: 7, Tag: int('f')},
				B:          &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				IsMinBound: true,
			},
			Expected: sortVersionsExpected{Result: -1},
		},
		{
			Name: "test_a_less_than_b_patch",
			Input: sortVersionsInput{
				A:          &provider.Version{Major: 4, Minor: 2, Patch: 6, Tag: int('f')},
				B:          &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				IsMinBound: true,
			},
			Expected: sortVersionsExpected{Result: -1},
		},
		{
			Name: "test_a_less_than_b_tag",
			Input: sortVersionsInput{
				A:          &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('e')},
				B:          &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				IsMinBound: true,
			},
			Expected: sortVersionsExpected{Result: -1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(sortVersionsInput)
			expected := testCase.Expected.(sortVersionsExpected)
			result := sortVersions(input.A, input.B, input.IsMinBound)

			assert.Equal(t, expected.Result, result, test.MessageEqual(expected.Result, result, testCase.Name))
		})
	}
}

type sortVersionRangesExpected struct {
	Result int
}

type sortVersionRangesInput struct {
	A VersionRange
	B VersionRange
}

func TestSortVersionRanges(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_equal_ranges",
			Input: sortVersionRangesInput{
				A: VersionRange{
					min: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
					max: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('w')},
				},
				B: VersionRange{
					min: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
					max: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('w')},
				},
			},
			Expected: sortVersionRangesExpected{Result: 0},
		},
		{
			Name: "test_a_min_less_than_b_min",
			Input: sortVersionRangesInput{
				A: VersionRange{
					min: &provider.Version{Major: 3, Minor: 2, Patch: 7, Tag: int('f')},
					max: &provider.Version{Major: 3, Minor: 2, Patch: 7, Tag: int('w')},
				},
				B: VersionRange{
					min: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
					max: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('w')},
				},
			},
			Expected: sortVersionRangesExpected{Result: -1},
		},
		{
			Name: "test_same_min_a_max_less_than_b_max",
			Input: sortVersionRangesInput{
				A: VersionRange{
					min: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
					max: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('g')},
				},
				B: VersionRange{
					min: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
					max: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('w')},
				},
			},
			Expected: sortVersionRangesExpected{Result: -1},
		},
		{
			Name: "test_same_min_a_unbounded_max",
			Input: sortVersionRangesInput{
				A: VersionRange{
					min: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
					max: nil,
				},
				B: VersionRange{
					min: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
					max: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('w')},
				},
			},
			Expected: sortVersionRangesExpected{Result: 1},
		},
		{
			Name: "test_a_unbounded_min",
			Input: sortVersionRangesInput{
				A: VersionRange{
					min: nil,
					max: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('w')},
				},
				B: VersionRange{
					min: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
					max: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('w')},
				},
			},
			Expected: sortVersionRangesExpected{Result: -1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(sortVersionRangesInput)
			expected := testCase.Expected.(sortVersionRangesExpected)
			result := sortVersionRanges(input.A, input.B)

			assert.Equal(t, expected.Result, result, test.MessageEqual(expected.Result, result, testCase.Name))
		})
	}
}

func TestVersionsSort(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	// Create unsorted versions
	versions := &Versions{
		raw: "5.0(1a)-,3.0(1a)-3.0(1z),4.0(1a)-4.0(1z)",
		ranges: []VersionRange{
			{
				raw: "5.0(1a)-",
				min: &provider.Version{Major: 5, Minor: 0, Patch: 1, Tag: int('a')},
				max: nil,
			},
			{
				raw: "3.0(1a)-3.0(1z)",
				min: &provider.Version{Major: 3, Minor: 0, Patch: 1, Tag: int('a')},
				max: &provider.Version{Major: 3, Minor: 0, Patch: 1, Tag: int('z')},
			},
			{
				raw: "4.0(1a)-4.0(1z)",
				min: &provider.Version{Major: 4, Minor: 0, Patch: 1, Tag: int('a')},
				max: &provider.Version{Major: 4, Minor: 0, Patch: 1, Tag: int('z')},
			},
		},
	}

	// Sort
	versions.Sort()

	// Verify order
	assert.Equal(t, 3, versions.Ranges()[0].Min().Major, "First range should have Major=3")
	assert.Equal(t, 4, versions.Ranges()[1].Min().Major, "Second range should have Major=4")
	assert.Equal(t, 5, versions.Ranges()[2].Min().Major, "Third range should have Major=5")
}

func TestVersionRangeString(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_single_version",
			Input: VersionRange{
				min: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				max: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
			},
			Expected: "4.2(7f)",
		},
		{
			Name: "test_single_version_no_tag",
			Input: VersionRange{
				min: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: 0},
				max: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: 0},
			},
			Expected: "4.2(7)",
		},
		{
			Name: "test_bounded_range",
			Input: VersionRange{
				min: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				max: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('w')},
			},
			Expected: "4.2(7f) to 4.2(7w)",
		},
		{
			Name: "test_unbounded_upper",
			Input: VersionRange{
				min: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('f')},
				max: nil,
			},
			Expected: "4.2(7f) and later",
		},
		{
			Name: "test_unbounded_lower",
			Input: VersionRange{
				min: nil,
				max: &provider.Version{Major: 4, Minor: 2, Patch: 7, Tag: int('w')},
			},
			Expected: "up to 4.2(7w)",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			versionRange := testCase.Input.(VersionRange)
			expected := testCase.Expected.(string)

			assert.Equal(t, expected, versionRange.String(), test.MessageEqual(expected, versionRange.String(), testCase.Name))
		})
	}
}
