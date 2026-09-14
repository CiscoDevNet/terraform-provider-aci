package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

const (
	constTestDirectoryForGetFileNamesFromDirectory = "./test/data/test_get_file_names_from_directory"
	constTestFile1WithoutExtension                 = "file1"
	constTestFile1WithExtension                    = "file1.json"
	constTestFile2WithoutExtension                 = "file2"
	constTestFile2WithExtension                    = "file2.json"
	constTestFile3                                 = "file3"
	constTestDir1                                  = "dir1"
)

func TestGetFileNamesFromDirectoryWithExtension(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	filenames := GetFileNamesFromDirectory(constTestDirectoryForGetFileNamesFromDirectory, false)

	assert.NotEmpty(t, filenames, test.MessageNotEmpty(filenames, t.Name()))
	assert.Len(t, filenames, 3)
	assert.Contains(t, filenames, constTestFile1WithExtension, test.MessageContains(filenames, constTestFile1WithExtension, t.Name()))
	assert.Contains(t, filenames, constTestFile2WithExtension, test.MessageContains(filenames, constTestFile2WithExtension, t.Name()))
	assert.Contains(t, filenames, constTestFile3, test.MessageContains(filenames, constTestFile3, t.Name()))
	assert.NotContains(t, filenames, constTestDir1, test.MessageNotContains(filenames, constTestDir1, t.Name()))
}

func TestGetFileNamesFromDirectoryWithoutExtension(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	filenames := GetFileNamesFromDirectory(constTestDirectoryForGetFileNamesFromDirectory, true)

	assert.NotEmpty(t, filenames, test.MessageNotEmpty(filenames, t.Name()))
	assert.Len(t, filenames, 3)
	assert.Contains(t, filenames, constTestFile1WithoutExtension, test.MessageContains(filenames, constTestFile1WithoutExtension, t.Name()))
	assert.Contains(t, filenames, constTestFile2WithoutExtension, test.MessageContains(filenames, constTestFile2WithoutExtension, t.Name()))
	assert.Contains(t, filenames, constTestFile3, test.MessageContains(filenames, constTestFile3, t.Name()))
	assert.NotContains(t, filenames, constTestDir1, test.MessageNotContains(filenames, constTestDir1, t.Name()))
}

func TestUnderscore(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_lowercase", Input: "tenant", Expected: "tenant"},
		{Name: "test_capitalized", Input: "Tenant", Expected: "tenant"},
		{Name: "test_with_number", Input: "Tenant1", Expected: "tenant1"},
		{Name: "test_camel_case", Input: "ApplicationEndpointGroup", Expected: "application_endpoint_group"},
		{Name: "test_space_separated", Input: "Application Endpoint Group", Expected: "application_endpoint_group"},
		{Name: "test_hyphen_separated", Input: "Application-Endpoint-Group", Expected: "application_endpoint_group"},
		{Name: "test_mixed_separators", Input: "Application Endpoint-Group", Expected: "application_endpoint_group"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			result := Underscore(testCase.Input.(string))
			assert.Equal(t, testCase.Expected, result, test.MessageEqual(testCase.Expected, result, testCase.Name))
		})
	}
}

func TestPlural(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_policy_to_policies", Input: "monitor_policy", Expected: "monitor_policies"},
		{Name: "test_add_s", Input: "annotation", Expected: "annotations"},
		{Name: "test_already_plural_s", Input: "contracts", Expected: "contracts"},
		{Name: "test_already_plural_mappings", Input: "remote_site_id_mappings", Expected: "remote_site_id_mappings"},
		{Name: "test_already_plural_ies", Input: "monitor_policies", Expected: "monitor_policies"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			result := Plural(testCase.Input.(string))
			assert.Equal(t, testCase.Expected, result, test.MessageEqual(testCase.Expected, result, testCase.Name))
		})
	}
}

func TestUnderscoreEdgeCases(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_empty_string", Input: "", Expected: ""},
		{Name: "test_single_lowercase", Input: "a", Expected: "a"},
		{Name: "test_single_uppercase", Input: "A", Expected: "a"},
		{Name: "test_numbers_only", Input: "123", Expected: "123"},
		{Name: "test_leading_number", Input: "1Tenant", Expected: "1_tenant"},
		{Name: "test_underscore_input", Input: "already_snake", Expected: "already_snake"},
		{Name: "test_multiple_underscores", Input: "a__b", Expected: "a__b"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			result := Underscore(testCase.Input.(string))
			assert.Equal(t, testCase.Expected, result, test.MessageEqual(testCase.Expected, result, testCase.Name))
		})
	}
}

func TestGetFileNamesFromDirectoryNonExistent(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	filenames := GetFileNamesFromDirectory("./non_existent_directory", false)

	assert.Empty(t, filenames, test.MessageEmpty(filenames, t.Name()))
}

func TestGetValueFromMapWithOverrideString(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_nil_map_empty_override", Input: getStringValueFromMapWithOverrideInput{Key: "x"}, Expected: ""},
		{Name: "test_nil_map_with_override", Input: getStringValueFromMapWithOverrideInput{Key: "x", Override: "v"}, Expected: "v"},
		{Name: "test_missing_key_empty_override", Input: getStringValueFromMapWithOverrideInput{Values: map[string]any{}, Key: "x"}, Expected: ""},
		{Name: "test_missing_key_with_override", Input: getStringValueFromMapWithOverrideInput{Values: map[string]any{}, Key: "x", Override: "v"}, Expected: "v"},
		{Name: "test_string_value_no_override", Input: getStringValueFromMapWithOverrideInput{Values: map[string]any{"x": "m"}, Key: "x"}, Expected: "m"},
		{Name: "test_string_value_override_wins", Input: getStringValueFromMapWithOverrideInput{Values: map[string]any{"x": "m"}, Key: "x", Override: "v"}, Expected: "v"},
		{Name: "test_non_string_value_returns_zero", Input: getStringValueFromMapWithOverrideInput{Values: map[string]any{"x": 42}, Key: "x"}, Expected: ""},
		{Name: "test_nil_value_returns_zero", Input: getStringValueFromMapWithOverrideInput{Values: map[string]any{"x": nil}, Key: "x"}, Expected: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(getStringValueFromMapWithOverrideInput)
			actual := GetValueFromMapWithOverride(input.Values, input.Key, input.Override)
			assert.Equal(t, testCase.Expected, actual, test.MessageEqual(testCase.Expected, actual, testCase.Name))
		})
	}
}

func TestGetValueFromMapWithOverrideOtherTypes(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	value := "v"
	mapValue := "m"
	testCases := []struct {
		Name string
		Run  func(*testing.T)
	}{
		{Name: "test_int_value_no_override", Run: func(t *testing.T) { assert.Equal(t, 7, GetValueFromMapWithOverride(map[string]any{"x": 7}, "x", 0)) }},
		{Name: "test_int_override_wins", Run: func(t *testing.T) { assert.Equal(t, 9, GetValueFromMapWithOverride(map[string]any{"x": 7}, "x", 9)) }},
		{Name: "test_bool_true_value_no_override", Run: func(t *testing.T) {
			assert.Equal(t, true, GetValueFromMapWithOverride(map[string]any{"x": true}, "x", false))
		}},
		{Name: "test_pointer_override_wins", Run: func(t *testing.T) { assert.Equal(t, &value, GetValueFromMapWithOverride[*string](nil, "x", &value)) }},
		{Name: "test_pointer_nil_override_falls_back", Run: func(t *testing.T) {
			assert.Equal(t, &mapValue, GetValueFromMapWithOverride[*string](map[string]any{"x": &mapValue}, "x", nil))
		}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			testCase.Run(t)
		})
	}
}

type getStringValueFromMapWithOverrideInput struct {
	Values   map[string]any
	Key      string
	Override string
}
