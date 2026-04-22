package data

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

type setAttributeNameInput struct {
	PropertyName       string
	PropertyDefinition PropertyDefinition
	GlobalDefinition   GlobalMetaDefinition
}

func TestSetAttributeName(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_default_snake_case_from_camelCase",
			Input: setAttributeNameInput{
				PropertyName:       "ipLearning",
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: "ip_learning",
		},
		{
			Name: "test_default_snake_case_from_PascalCase",
			Input: setAttributeNameInput{
				PropertyName:       "OptimizeWanBandwidth",
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: "optimize_wan_bandwidth",
		},
		{
			Name: "test_default_snake_case_single_word",
			Input: setAttributeNameInput{
				PropertyName:       "name",
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: "name",
		},
		{
			Name: "test_definition_override",
			Input: setAttributeNameInput{
				PropertyName:       "pcEnfPref",
				PropertyDefinition: PropertyDefinition{AttributeName: "policy_control_enforcement"},
			},
			Expected: "policy_control_enforcement",
		},
		{
			Name: "test_definition_override_takes_precedence",
			Input: setAttributeNameInput{
				PropertyName:       "ipLearning",
				PropertyDefinition: PropertyDefinition{AttributeName: "custom_ip_learning"},
			},
			Expected: "custom_ip_learning",
		},
		{
			Name: "test_global_override",
			Input: setAttributeNameInput{
				PropertyName:       "descr",
				PropertyDefinition: PropertyDefinition{},
				GlobalDefinition: GlobalMetaDefinition{
					AttributeNameOverrides: map[string]string{
						"descr": "description",
					},
				},
			},
			Expected: "description",
		},
		{
			Name: "test_definition_override_takes_precedence_over_global",
			Input: setAttributeNameInput{
				PropertyName:       "descr",
				PropertyDefinition: PropertyDefinition{AttributeName: "custom_description"},
				GlobalDefinition: GlobalMetaDefinition{
					AttributeNameOverrides: map[string]string{
						"descr": "description",
					},
				},
			},
			Expected: "custom_description",
		},
		{
			Name: "test_global_override_no_match_falls_back_to_snake_case",
			Input: setAttributeNameInput{
				PropertyName:       "ipLearning",
				PropertyDefinition: PropertyDefinition{},
				GlobalDefinition: GlobalMetaDefinition{
					AttributeNameOverrides: map[string]string{
						"descr": "description",
					},
				},
			},
			Expected: "ip_learning",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setAttributeNameInput)
			expected := testCase.Expected.(string)

			property := &Property{
				PropertyName:       input.PropertyName,
				globalDefinition:   input.GlobalDefinition,
				propertyDefinition: input.PropertyDefinition,
			}

			property.setAttributeName()

			assert.Equal(t, expected, property.AttributeName, test.MessageEqual(expected, property.AttributeName, testCase.Name))
		})
	}
}

type setRequiredInput struct {
	MetaDetails        map[string]interface{}
	PropertyDefinition PropertyDefinition
}

func TestSetRequired(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_isConfigurable_and_isNaming_true",
			Input: setRequiredInput{
				MetaDetails:        map[string]interface{}{"isConfigurable": true, "isNaming": true},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: true,
		},
		{
			Name: "test_isConfigurable_true_isNaming_false",
			Input: setRequiredInput{
				MetaDetails:        map[string]interface{}{"isConfigurable": true, "isNaming": false},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_isConfigurable_false_isNaming_true",
			Input: setRequiredInput{
				MetaDetails:        map[string]interface{}{"isConfigurable": false, "isNaming": true},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_isConfigurable_false_isNaming_false",
			Input: setRequiredInput{
				MetaDetails:        map[string]interface{}{"isConfigurable": false, "isNaming": false},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_definition_restriction_required",
			Input: setRequiredInput{
				MetaDetails:        map[string]interface{}{"isConfigurable": true, "isNaming": false},
				PropertyDefinition: PropertyDefinition{Restriction: "required"},
			},
			Expected: true,
		},
		{
			Name: "test_definition_restriction_required_overrides_non_naming",
			Input: setRequiredInput{
				MetaDetails:        map[string]interface{}{"isConfigurable": false, "isNaming": false},
				PropertyDefinition: PropertyDefinition{Restriction: "required"},
			},
			Expected: true,
		},
		{
			Name: "test_empty_meta_details",
			Input: setRequiredInput{
				MetaDetails:        map[string]interface{}{},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setRequiredInput)
			expected := testCase.Expected.(bool)

			property := &Property{
				PropertyName:       "testProp",
				metaDetails:        input.MetaDetails,
				propertyDefinition: input.PropertyDefinition,
			}

			property.setRequired()

			assert.Equal(t, expected, property.Required, test.MessageEqual(expected, property.Required, testCase.Name))
		})
	}
}

type setOptionalInput struct {
	MetaDetails        map[string]interface{}
	PropertyDefinition PropertyDefinition
	Required           bool
}

func TestSetOptional(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_isConfigurable_true_not_required",
			Input: setOptionalInput{
				MetaDetails:        map[string]interface{}{"isConfigurable": true},
				PropertyDefinition: PropertyDefinition{},
				Required:           false,
			},
			Expected: true,
		},
		{
			Name: "test_isConfigurable_true_already_required",
			Input: setOptionalInput{
				MetaDetails:        map[string]interface{}{"isConfigurable": true},
				PropertyDefinition: PropertyDefinition{},
				Required:           true,
			},
			Expected: false,
		},
		{
			Name: "test_isConfigurable_false",
			Input: setOptionalInput{
				MetaDetails:        map[string]interface{}{"isConfigurable": false},
				PropertyDefinition: PropertyDefinition{},
				Required:           false,
			},
			Expected: false,
		},
		{
			Name: "test_definition_restriction_optional",
			Input: setOptionalInput{
				MetaDetails:        map[string]interface{}{"isConfigurable": false},
				PropertyDefinition: PropertyDefinition{Restriction: "optional"},
				Required:           false,
			},
			Expected: true,
		},
		{
			Name: "test_definition_restriction_read_only_isConfigurable_true",
			Input: setOptionalInput{
				MetaDetails:        map[string]interface{}{"isConfigurable": true},
				PropertyDefinition: PropertyDefinition{Restriction: "read_only"},
				Required:           false,
			},
			Expected: false,
		},
		{
			Name: "test_empty_meta_details",
			Input: setOptionalInput{
				MetaDetails:        map[string]interface{}{},
				PropertyDefinition: PropertyDefinition{},
				Required:           false,
			},
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setOptionalInput)
			expected := testCase.Expected.(bool)

			property := &Property{
				PropertyName:       "testProp",
				metaDetails:        input.MetaDetails,
				propertyDefinition: input.PropertyDefinition,
				Required:           input.Required,
			}

			property.setOptional()

			assert.Equal(t, expected, property.Optional, test.MessageEqual(expected, property.Optional, testCase.Name))
		})
	}
}

type setReadOnlyInput struct {
	PropertyDefinition PropertyDefinition
}

func TestSetReadOnly(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_definition_restriction_read_only",
			Input: setReadOnlyInput{
				PropertyDefinition: PropertyDefinition{Restriction: "read_only"},
			},
			Expected: true,
		},
		{
			Name: "test_no_restriction",
			Input: setReadOnlyInput{
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_definition_restriction_required",
			Input: setReadOnlyInput{
				PropertyDefinition: PropertyDefinition{Restriction: "required"},
			},
			Expected: false,
		},
		{
			Name: "test_definition_restriction_optional",
			Input: setReadOnlyInput{
				PropertyDefinition: PropertyDefinition{Restriction: "optional"},
			},
			Expected: false,
		},
		{
			Name: "test_definition_restriction_exclude",
			Input: setReadOnlyInput{
				PropertyDefinition: PropertyDefinition{Restriction: "exclude"},
			},
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setReadOnlyInput)
			expected := testCase.Expected.(bool)

			property := &Property{
				PropertyName:       "testProp",
				metaDetails:        map[string]interface{}{},
				propertyDefinition: input.PropertyDefinition,
			}

			property.setReadOnly()

			assert.Equal(t, expected, property.ReadOnly, test.MessageEqual(expected, property.ReadOnly, testCase.Name))
		})
	}
}

type setComputedInput struct {
	Required bool
}

func TestSetComputed(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_not_required_is_computed",
			Input: setComputedInput{
				Required: false,
			},
			Expected: true,
		},
		{
			Name: "test_required_is_not_computed",
			Input: setComputedInput{
				Required: true,
			},
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setComputedInput)
			expected := testCase.Expected.(bool)

			property := &Property{
				PropertyName: "testProp",
				Required:     input.Required,
			}

			property.setComputed()

			assert.Equal(t, expected, property.Computed, test.MessageEqual(expected, property.Computed, testCase.Name))
		})
	}
}

type setSensitiveInput struct {
	MetaDetails        map[string]interface{}
	PropertyDefinition PropertyDefinition
}

func TestSetSensitive(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_meta_secure_true",
			Input: setSensitiveInput{
				MetaDetails:        map[string]interface{}{"secure": true},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: true,
		},
		{
			Name: "test_meta_secure_false",
			Input: setSensitiveInput{
				MetaDetails:        map[string]interface{}{"secure": false},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_definition_override_true",
			Input: setSensitiveInput{
				MetaDetails:        map[string]interface{}{"secure": false},
				PropertyDefinition: PropertyDefinition{Sensitive: true},
			},
			Expected: true,
		},
		{
			Name: "test_definition_override_true_no_meta_secure",
			Input: setSensitiveInput{
				MetaDetails:        map[string]interface{}{},
				PropertyDefinition: PropertyDefinition{Sensitive: true},
			},
			Expected: true,
		},
		{
			Name: "test_no_override_no_meta_secure",
			Input: setSensitiveInput{
				MetaDetails:        map[string]interface{}{},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setSensitiveInput)
			expected := testCase.Expected.(bool)

			property := &Property{
				PropertyName:       "testProp",
				metaDetails:        input.MetaDetails,
				propertyDefinition: input.PropertyDefinition,
			}

			property.setSensitive()

			assert.Equal(t, expected, property.Sensitive, test.MessageEqual(expected, property.Sensitive, testCase.Name))
		})
	}
}

func TestPropertySetDeprecated(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_deprecated_false",
			Input:    PropertyDefinition{},
			Expected: false,
		},
		{
			Name:     "test_deprecated_true",
			Input:    PropertyDefinition{Deprecated: true},
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			property := &Property{
				PropertyName:       "testProp",
				propertyDefinition: testCase.Input.(PropertyDefinition),
			}

			property.setDeprecated()

			assert.Equal(t, testCase.Expected, property.Deprecated, test.MessageEqual(testCase.Expected, property.Deprecated, testCase.Name))
		})
	}
}

type setPropertyDeprecatedVersionsExpected struct {
	Raw      string
	String   string
	Nil      bool
	Error    bool
	ErrorMsg string
}

func TestPropertySetDeprecatedVersions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:  "test_deprecated_versions_not_set",
			Input: PropertyDefinition{},
			Expected: setPropertyDeprecatedVersionsExpected{
				Nil: true,
			},
		},
		{
			Name:  "test_deprecated_versions_single_range",
			Input: PropertyDefinition{DeprecatedVersions: "4.2(7f)-"},
			Expected: setPropertyDeprecatedVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name:  "test_deprecated_versions_bounded_range",
			Input: PropertyDefinition{DeprecatedVersions: "3.2(10e)-4.2(7f)"},
			Expected: setPropertyDeprecatedVersionsExpected{
				Raw:    "3.2(10e)-4.2(7f)",
				String: "3.2(10e) to 4.2(7f)",
			},
		},
		{
			Name:  "test_deprecated_versions_multiple_ranges",
			Input: PropertyDefinition{DeprecatedVersions: "3.2(10e)-3.2(10g),4.2(7f)-"},
			Expected: setPropertyDeprecatedVersionsExpected{
				Raw:    "3.2(10e)-3.2(10g),4.2(7f)-",
				String: "3.2(10e) to 3.2(10g), 4.2(7f) and later",
			},
		},
		{
			Name:  "test_error_invalid_deprecated_version",
			Input: PropertyDefinition{DeprecatedVersions: "invalid"},
			Expected: setPropertyDeprecatedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse deprecated versions for property 'testProp': invalid version 'invalid': unknown",
			},
		},
		{
			Name:  "test_error_invalid_deprecated_version_in_range",
			Input: PropertyDefinition{DeprecatedVersions: "4.2(7f)-,invalid"},
			Expected: setPropertyDeprecatedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse deprecated versions for property 'testProp': invalid version 'invalid': unknown",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			expected := testCase.Expected.(setPropertyDeprecatedVersionsExpected)

			property := &Property{
				PropertyName:       "testProp",
				propertyDefinition: testCase.Input.(PropertyDefinition),
			}

			err := property.setDeprecatedVersions()

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else if expected.Nil {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Nil(t, property.DeprecatedVersions, "expected DeprecatedVersions to be nil")
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, expected.Raw, property.DeprecatedVersions.Raw(), test.MessageEqual(expected.Raw, property.DeprecatedVersions.Raw(), testCase.Name))
				assert.Equal(t, expected.String, property.DeprecatedVersions.String(), test.MessageEqual(expected.String, property.DeprecatedVersions.String(), testCase.Name))
			}
		})
	}
}

type setPropertySupportedVersionsInput struct {
	MetaDetails        map[string]interface{}
	PropertyDefinition PropertyDefinition
}

type setPropertySupportedVersionsExpected struct {
	Raw      string
	String   string
	Nil      bool
	Error    bool
	ErrorMsg string
}

func TestPropertySetSupportedVersions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_versions_not_set",
			Input: setPropertySupportedVersionsInput{
				MetaDetails:        map[string]interface{}{},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: setPropertySupportedVersionsExpected{
				Nil: true,
			},
		},
		{
			Name: "test_versions_from_meta_file",
			Input: setPropertySupportedVersionsInput{
				MetaDetails:        map[string]interface{}{"versions": "4.2(7f)-"},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: setPropertySupportedVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name: "test_versions_from_definition_override",
			Input: setPropertySupportedVersionsInput{
				MetaDetails:        map[string]interface{}{"versions": "4.2(7f)-"},
				PropertyDefinition: PropertyDefinition{SupportedVersions: "5.0(1a)-"},
			},
			Expected: setPropertySupportedVersionsExpected{
				Raw:    "5.0(1a)-",
				String: "5.0(1a) and later",
			},
		},
		{
			Name: "test_versions_from_definition_when_meta_empty",
			Input: setPropertySupportedVersionsInput{
				MetaDetails:        map[string]interface{}{},
				PropertyDefinition: PropertyDefinition{SupportedVersions: "5.0(1a)-"},
			},
			Expected: setPropertySupportedVersionsExpected{
				Raw:    "5.0(1a)-",
				String: "5.0(1a) and later",
			},
		},
		{
			Name: "test_bounded_range",
			Input: setPropertySupportedVersionsInput{
				MetaDetails:        map[string]interface{}{"versions": "3.2(10e)-4.2(7f)"},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: setPropertySupportedVersionsExpected{
				Raw:    "3.2(10e)-4.2(7f)",
				String: "3.2(10e) to 4.2(7f)",
			},
		},
		{
			Name: "test_multiple_ranges",
			Input: setPropertySupportedVersionsInput{
				MetaDetails:        map[string]interface{}{"versions": "3.2(10e)-3.2(10g),4.2(7f)-"},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: setPropertySupportedVersionsExpected{
				Raw:    "3.2(10e)-3.2(10g),4.2(7f)-",
				String: "3.2(10e) to 3.2(10g), 4.2(7f) and later",
			},
		},
		{
			Name: "test_error_invalid_version",
			Input: setPropertySupportedVersionsInput{
				MetaDetails:        map[string]interface{}{"versions": "invalid"},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: setPropertySupportedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse versions for property 'testProp': invalid version 'invalid': unknown",
			},
		},
		{
			Name: "test_error_invalid_version_in_range",
			Input: setPropertySupportedVersionsInput{
				MetaDetails:        map[string]interface{}{"versions": "4.2(7f)-,invalid"},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: setPropertySupportedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse versions for property 'testProp': invalid version 'invalid': unknown",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPropertySupportedVersionsInput)
			expected := testCase.Expected.(setPropertySupportedVersionsExpected)

			property := &Property{
				PropertyName:       "testProp",
				metaDetails:        input.MetaDetails,
				propertyDefinition: input.PropertyDefinition,
			}

			err := property.setSupportedVersions()

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else if expected.Nil {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Nil(t, property.SupportedVersions, "expected SupportedVersions to be nil")
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, expected.Raw, property.SupportedVersions.Raw(), test.MessageEqual(expected.Raw, property.SupportedVersions.Raw(), testCase.Name))
				assert.Equal(t, expected.String, property.SupportedVersions.String(), test.MessageEqual(expected.String, property.SupportedVersions.String(), testCase.Name))
			}
		})
	}
}
