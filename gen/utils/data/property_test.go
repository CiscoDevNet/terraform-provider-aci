package data

import (
	"bytes"
	"os"
	"strings"
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
	MetaDetails        map[string]any
	PropertyDefinition PropertyDefinition
}

func TestSetRequired(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_isConfigurable_and_isNaming_true",
			Input: setRequiredInput{
				MetaDetails:        map[string]any{"isConfigurable": true, "isNaming": true},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: true,
		},
		{
			Name: "test_isConfigurable_true_isNaming_false",
			Input: setRequiredInput{
				MetaDetails:        map[string]any{"isConfigurable": true, "isNaming": false},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_isConfigurable_false_isNaming_true",
			Input: setRequiredInput{
				MetaDetails:        map[string]any{"isConfigurable": false, "isNaming": true},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_isConfigurable_false_isNaming_false",
			Input: setRequiredInput{
				MetaDetails:        map[string]any{"isConfigurable": false, "isNaming": false},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_definition_restriction_required",
			Input: setRequiredInput{
				MetaDetails:        map[string]any{"isConfigurable": true, "isNaming": false},
				PropertyDefinition: PropertyDefinition{Restriction: Required},
			},
			Expected: true,
		},
		{
			Name: "test_definition_restriction_required_overrides_non_naming",
			Input: setRequiredInput{
				MetaDetails:        map[string]any{"isConfigurable": false, "isNaming": false},
				PropertyDefinition: PropertyDefinition{Restriction: Required},
			},
			Expected: true,
		},
		{
			Name: "test_empty_meta_details",
			Input: setRequiredInput{
				MetaDetails:        map[string]any{},
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

type setRequiresReplaceInput struct {
	MetaDetails        map[string]any
	PropertyDefinition PropertyDefinition
}

func TestSetRequiresReplace(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	boolTrue := true
	boolFalse := false

	testCases := []test.TestCase{
		{
			Name: "test_isNaming_true",
			Input: setRequiresReplaceInput{
				MetaDetails:        map[string]any{"isNaming": true},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: true,
		},
		{
			Name: "test_isNaming_false",
			Input: setRequiresReplaceInput{
				MetaDetails:        map[string]any{"isNaming": false},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_isNaming_missing",
			Input: setRequiresReplaceInput{
				MetaDetails:        map[string]any{},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_definition_override_true",
			Input: setRequiresReplaceInput{
				MetaDetails:        map[string]any{"isNaming": false},
				PropertyDefinition: PropertyDefinition{RequiresReplace: &boolTrue},
			},
			Expected: true,
		},
		{
			Name: "test_definition_override_false_suppresses_naming",
			Input: setRequiresReplaceInput{
				MetaDetails:        map[string]any{"isNaming": true},
				PropertyDefinition: PropertyDefinition{RequiresReplace: &boolFalse},
			},
			Expected: false,
		},
		{
			Name: "test_definition_override_true_without_naming",
			Input: setRequiresReplaceInput{
				MetaDetails:        map[string]any{},
				PropertyDefinition: PropertyDefinition{RequiresReplace: &boolTrue},
			},
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setRequiresReplaceInput)
			expected := testCase.Expected.(bool)

			property := &Property{
				PropertyName:       "testProp",
				metaDetails:        input.MetaDetails,
				propertyDefinition: input.PropertyDefinition,
			}

			property.setRequiresReplace()

			assert.Equal(t, expected, property.RequiresReplace, test.MessageEqual(expected, property.RequiresReplace, testCase.Name))
		})
	}
}

type setOptionalInput struct {
	MetaDetails        map[string]any
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
				MetaDetails:        map[string]any{"isConfigurable": true},
				PropertyDefinition: PropertyDefinition{},
				Required:           false,
			},
			Expected: true,
		},
		{
			Name: "test_isConfigurable_true_already_required",
			Input: setOptionalInput{
				MetaDetails:        map[string]any{"isConfigurable": true},
				PropertyDefinition: PropertyDefinition{},
				Required:           true,
			},
			Expected: false,
		},
		{
			Name: "test_isConfigurable_false",
			Input: setOptionalInput{
				MetaDetails:        map[string]any{"isConfigurable": false},
				PropertyDefinition: PropertyDefinition{},
				Required:           false,
			},
			Expected: false,
		},
		{
			Name: "test_definition_restriction_optional",
			Input: setOptionalInput{
				MetaDetails:        map[string]any{"isConfigurable": false},
				PropertyDefinition: PropertyDefinition{Restriction: Optional},
				Required:           false,
			},
			Expected: true,
		},
		{
			Name: "test_definition_restriction_read_only_isConfigurable_true",
			Input: setOptionalInput{
				MetaDetails:        map[string]any{"isConfigurable": true},
				PropertyDefinition: PropertyDefinition{Restriction: ReadOnly},
				Required:           false,
			},
			Expected: false,
		},
		{
			Name: "test_empty_meta_details",
			Input: setOptionalInput{
				MetaDetails:        map[string]any{},
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
				PropertyDefinition: PropertyDefinition{Restriction: ReadOnly},
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
				PropertyDefinition: PropertyDefinition{Restriction: Required},
			},
			Expected: false,
		},
		{
			Name: "test_definition_restriction_optional",
			Input: setReadOnlyInput{
				PropertyDefinition: PropertyDefinition{Restriction: Optional},
			},
			Expected: false,
		},
		{
			Name: "test_definition_restriction_exclude",
			Input: setReadOnlyInput{
				PropertyDefinition: PropertyDefinition{Restriction: Exclude},
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
				metaDetails:        map[string]any{},
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
	MetaDetails        map[string]any
	PropertyDefinition PropertyDefinition
}

func TestSetSensitive(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_meta_secure_true",
			Input: setSensitiveInput{
				MetaDetails:        map[string]any{"secure": true},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: true,
		},
		{
			Name: "test_meta_secure_false",
			Input: setSensitiveInput{
				MetaDetails:        map[string]any{"secure": false},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_definition_override_true",
			Input: setSensitiveInput{
				MetaDetails:        map[string]any{"secure": false},
				PropertyDefinition: PropertyDefinition{Sensitive: true},
			},
			Expected: true,
		},
		{
			Name: "test_definition_override_true_no_meta_secure",
			Input: setSensitiveInput{
				MetaDetails:        map[string]any{},
				PropertyDefinition: PropertyDefinition{Sensitive: true},
			},
			Expected: true,
		},
		{
			Name: "test_no_override_no_meta_secure",
			Input: setSensitiveInput{
				MetaDetails:        map[string]any{},
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

type setPropertyDeprecatedInput struct {
	PropertyDefinition PropertyDefinition
	MetaDetails        map[string]any
}

func TestPropertySetDeprecated(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_meta_missing_no_override",
			Input:    setPropertyDeprecatedInput{},
			Expected: false,
		},
		{
			Name:     "test_meta_false_no_override",
			Input:    setPropertyDeprecatedInput{MetaDetails: map[string]any{"isDeprecated": false}},
			Expected: false,
		},
		{
			Name:     "test_meta_true_no_override",
			Input:    setPropertyDeprecatedInput{MetaDetails: map[string]any{"isDeprecated": true}},
			Expected: true,
		},
		{
			Name:     "test_meta_wrong_type",
			Input:    setPropertyDeprecatedInput{MetaDetails: map[string]any{"isDeprecated": "yes"}},
			Expected: false,
		},
		{
			Name:     "test_override_true_meta_missing",
			Input:    setPropertyDeprecatedInput{PropertyDefinition: PropertyDefinition{Deprecated: true}},
			Expected: true,
		},
		{
			Name:     "test_override_true_meta_false",
			Input:    setPropertyDeprecatedInput{PropertyDefinition: PropertyDefinition{Deprecated: true}, MetaDetails: map[string]any{"isDeprecated": false}},
			Expected: true,
		},
		{
			Name:     "test_override_false_meta_true",
			Input:    setPropertyDeprecatedInput{PropertyDefinition: PropertyDefinition{Deprecated: false}, MetaDetails: map[string]any{"isDeprecated": true}},
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPropertyDeprecatedInput)

			property := &Property{
				PropertyName:       "testProp",
				propertyDefinition: input.PropertyDefinition,
				metaDetails:        input.MetaDetails,
			}

			property.setDeprecated()

			assert.Equal(t, testCase.Expected, property.Deprecated, test.MessageEqual(testCase.Expected, property.Deprecated, testCase.Name))
		})
	}
}

type setPropertyDeprecatedVersionsInput struct {
	PropertyDefinition PropertyDefinition
	MetaDetails        map[string]any
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
			Input: setPropertyDeprecatedVersionsInput{},
			Expected: setPropertyDeprecatedVersionsExpected{
				Nil: true,
			},
		},
		{
			Name:  "test_deprecated_versions_single_range",
			Input: setPropertyDeprecatedVersionsInput{PropertyDefinition: PropertyDefinition{DeprecatedVersions: "4.2(7f)-"}},
			Expected: setPropertyDeprecatedVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name:  "test_deprecated_versions_bounded_range",
			Input: setPropertyDeprecatedVersionsInput{PropertyDefinition: PropertyDefinition{DeprecatedVersions: "3.2(10e)-4.2(7f)"}},
			Expected: setPropertyDeprecatedVersionsExpected{
				Raw:    "3.2(10e)-4.2(7f)",
				String: "3.2(10e) to 4.2(7f)",
			},
		},
		{
			Name:  "test_deprecated_versions_multiple_ranges",
			Input: setPropertyDeprecatedVersionsInput{PropertyDefinition: PropertyDefinition{DeprecatedVersions: "3.2(10e)-3.2(10g),4.2(7f)-"}},
			Expected: setPropertyDeprecatedVersionsExpected{
				Raw:    "3.2(10e)-3.2(10g),4.2(7f)-",
				String: "3.2(10e) to 3.2(10g), 4.2(7f) and later",
			},
		},
		{
			Name:  "test_error_invalid_deprecated_version",
			Input: setPropertyDeprecatedVersionsInput{PropertyDefinition: PropertyDefinition{DeprecatedVersions: "invalid"}},
			Expected: setPropertyDeprecatedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse deprecated versions for property 'testProp': invalid version 'invalid': unknown",
			},
		},
		{
			Name:  "test_error_invalid_deprecated_version_in_range",
			Input: setPropertyDeprecatedVersionsInput{PropertyDefinition: PropertyDefinition{DeprecatedVersions: "4.2(7f)-,invalid"}},
			Expected: setPropertyDeprecatedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse deprecated versions for property 'testProp': invalid version 'invalid': unknown",
			},
		},
		{
			Name:  "test_meta_deprecated_since_single_range",
			Input: setPropertyDeprecatedVersionsInput{MetaDetails: map[string]any{"deprecatedSince": "5.2(1g)-"}},
			Expected: setPropertyDeprecatedVersionsExpected{
				Raw:    "5.2(1g)-",
				String: "5.2(1g) and later",
			},
		},
		{
			Name:  "test_meta_deprecated_since_wrong_type",
			Input: setPropertyDeprecatedVersionsInput{MetaDetails: map[string]any{"deprecatedSince": 123}},
			Expected: setPropertyDeprecatedVersionsExpected{
				Nil: true,
			},
		},
		{
			Name: "test_override_replaces_meta",
			Input: setPropertyDeprecatedVersionsInput{
				PropertyDefinition: PropertyDefinition{DeprecatedVersions: "4.2(7f)-"},
				MetaDetails:        map[string]any{"deprecatedSince": "5.2(1g)-"},
			},
			Expected: setPropertyDeprecatedVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name:  "test_meta_parse_error",
			Input: setPropertyDeprecatedVersionsInput{MetaDetails: map[string]any{"deprecatedSince": "invalid"}},
			Expected: setPropertyDeprecatedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse deprecated versions for property 'testProp': invalid version 'invalid': unknown",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPropertyDeprecatedVersionsInput)
			expected := testCase.Expected.(setPropertyDeprecatedVersionsExpected)

			property := &Property{
				PropertyName:       "testProp",
				propertyDefinition: input.PropertyDefinition,
				metaDetails:        input.MetaDetails,
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

type setPropertyHiddenInput struct {
	PropertyDefinition PropertyDefinition
	MetaDetails        map[string]any
}

func TestPropertySetHidden(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_meta_missing_no_override",
			Input:    setPropertyHiddenInput{},
			Expected: false,
		},
		{
			Name:     "test_meta_false_no_override",
			Input:    setPropertyHiddenInput{MetaDetails: map[string]any{"isHidden": false}},
			Expected: false,
		},
		{
			Name:     "test_meta_true_no_override",
			Input:    setPropertyHiddenInput{MetaDetails: map[string]any{"isHidden": true}},
			Expected: true,
		},
		{
			Name:     "test_meta_wrong_type",
			Input:    setPropertyHiddenInput{MetaDetails: map[string]any{"isHidden": "yes"}},
			Expected: false,
		},
		{
			Name:     "test_override_true_meta_missing",
			Input:    setPropertyHiddenInput{PropertyDefinition: PropertyDefinition{Hidden: true}},
			Expected: true,
		},
		{
			Name:     "test_override_true_meta_false",
			Input:    setPropertyHiddenInput{PropertyDefinition: PropertyDefinition{Hidden: true}, MetaDetails: map[string]any{"isHidden": false}},
			Expected: true,
		},
		{
			Name:     "test_override_false_meta_true",
			Input:    setPropertyHiddenInput{PropertyDefinition: PropertyDefinition{Hidden: false}, MetaDetails: map[string]any{"isHidden": true}},
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPropertyHiddenInput)

			property := &Property{
				PropertyName:       "testProp",
				propertyDefinition: input.PropertyDefinition,
				metaDetails:        input.MetaDetails,
			}

			property.setHidden()

			assert.Equal(t, testCase.Expected, property.Hidden, test.MessageEqual(testCase.Expected, property.Hidden, testCase.Name))
		})
	}
}

type setPropertyHiddenVersionsInput struct {
	PropertyDefinition PropertyDefinition
	MetaDetails        map[string]any
}

type setPropertyHiddenVersionsExpected struct {
	Raw      string
	String   string
	Nil      bool
	Error    bool
	ErrorMsg string
}

func TestPropertySetHiddenVersions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:  "test_hidden_versions_not_set",
			Input: setPropertyHiddenVersionsInput{},
			Expected: setPropertyHiddenVersionsExpected{
				Nil: true,
			},
		},
		{
			Name:  "test_definition_single_range",
			Input: setPropertyHiddenVersionsInput{PropertyDefinition: PropertyDefinition{HiddenVersions: "4.2(7f)-"}},
			Expected: setPropertyHiddenVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name:  "test_meta_hidden_since_single_range",
			Input: setPropertyHiddenVersionsInput{MetaDetails: map[string]any{"hiddenSince": "5.2(1g)-"}},
			Expected: setPropertyHiddenVersionsExpected{
				Raw:    "5.2(1g)-",
				String: "5.2(1g) and later",
			},
		},
		{
			Name:  "test_meta_hidden_since_wrong_type",
			Input: setPropertyHiddenVersionsInput{MetaDetails: map[string]any{"hiddenSince": 123}},
			Expected: setPropertyHiddenVersionsExpected{
				Nil: true,
			},
		},
		{
			Name: "test_override_replaces_meta",
			Input: setPropertyHiddenVersionsInput{
				PropertyDefinition: PropertyDefinition{HiddenVersions: "4.2(7f)-"},
				MetaDetails:        map[string]any{"hiddenSince": "5.2(1g)-"},
			},
			Expected: setPropertyHiddenVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name:  "test_definition_parse_error",
			Input: setPropertyHiddenVersionsInput{PropertyDefinition: PropertyDefinition{HiddenVersions: "invalid"}},
			Expected: setPropertyHiddenVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse hidden versions for property 'testProp': invalid version 'invalid': unknown",
			},
		},
		{
			Name:  "test_meta_parse_error",
			Input: setPropertyHiddenVersionsInput{MetaDetails: map[string]any{"hiddenSince": "invalid"}},
			Expected: setPropertyHiddenVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse hidden versions for property 'testProp': invalid version 'invalid': unknown",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPropertyHiddenVersionsInput)
			expected := testCase.Expected.(setPropertyHiddenVersionsExpected)

			property := &Property{
				PropertyName:       "testProp",
				propertyDefinition: input.PropertyDefinition,
				metaDetails:        input.MetaDetails,
			}

			err := property.setHiddenVersions()

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else if expected.Nil {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Nil(t, property.HiddenVersions, "expected HiddenVersions to be nil")
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, expected.Raw, property.HiddenVersions.Raw(), test.MessageEqual(expected.Raw, property.HiddenVersions.Raw(), testCase.Name))
				assert.Equal(t, expected.String, property.HiddenVersions.String(), test.MessageEqual(expected.String, property.HiddenVersions.String(), testCase.Name))
			}
		})
	}
}

type setPropertyValidatorsInput struct {
	PropertyDefinition PropertyDefinition
	MetaDetails        map[string]any
}

type setPropertyValidatorsExpected struct {
	Validators []Validator
	Error      bool
	ErrorMsg   string
}

func TestPropertySetValidators(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	regexInclude := []any{
		map[string]any{"regex": "^[a-zA-Z0-9_.-]+$", "type": "include"},
	}

	testCases := []test.TestCase{
		{
			Name:     "test_meta_missing",
			Input:    setPropertyValidatorsInput{},
			Expected: setPropertyValidatorsExpected{Validators: nil},
		},
		{
			Name: "test_meta_empty_array",
			Input: setPropertyValidatorsInput{
				MetaDetails: map[string]any{"validators": []any{}},
			},
			Expected: setPropertyValidatorsExpected{Validators: nil},
		},
		{
			Name: "test_meta_min_max_only",
			Input: setPropertyValidatorsInput{
				MetaDetails: map[string]any{
					"validators": []any{
						map[string]any{"min": float64(0), "max": float64(7)},
					},
				},
			},
			Expected: setPropertyValidatorsExpected{
				Validators: []Validator{{Min: 0, Max: 7}},
			},
		},
		{
			Name: "test_meta_min_max_with_regex_include",
			Input: setPropertyValidatorsInput{
				MetaDetails: map[string]any{
					"validators": []any{
						map[string]any{
							"min":    float64(0),
							"max":    float64(63),
							"regexs": regexInclude,
						},
					},
				},
			},
			Expected: setPropertyValidatorsExpected{
				Validators: []Validator{{
					Min: 0, Max: 63,
					RegexList: []RegexStatement{{Regex: "^[a-zA-Z0-9_.-]+$", Type: Include}},
				}},
			},
		},
		{
			Name: "test_meta_multiple_validators",
			Input: setPropertyValidatorsInput{
				MetaDetails: map[string]any{
					"validators": []any{
						map[string]any{"min": float64(0), "max": float64(7)},
						map[string]any{"min": float64(0), "max": float64(63), "regexs": regexInclude},
					},
				},
			},
			Expected: setPropertyValidatorsExpected{
				Validators: []Validator{
					{Min: 0, Max: 7},
					{Min: 0, Max: 63, RegexList: []RegexStatement{{Regex: "^[a-zA-Z0-9_.-]+$", Type: Include}}},
				},
			},
		},
		{
			Name: "test_meta_wrong_top_type",
			Input: setPropertyValidatorsInput{
				MetaDetails: map[string]any{"validators": "not-a-list"},
			},
			Expected: setPropertyValidatorsExpected{
				Error:    true,
				ErrorMsg: "failed to parse validators for property 'testProp': expected validators to be a list, got string",
			},
		},
		{
			Name: "test_meta_entry_wrong_type",
			Input: setPropertyValidatorsInput{
				MetaDetails: map[string]any{"validators": []any{42}},
			},
			Expected: setPropertyValidatorsExpected{
				Error:    true,
				ErrorMsg: "failed to parse validators for property 'testProp': expected validator entry 0 to be a map, got int",
			},
		},
		{
			Name: "test_meta_unknown_regex_type",
			Input: setPropertyValidatorsInput{
				MetaDetails: map[string]any{
					"validators": []any{
						map[string]any{
							"min": float64(0),
							"max": float64(63),
							"regexs": []any{
								map[string]any{"regex": ".+", "type": "exclude"},
							},
						},
					},
				},
			},
			Expected: setPropertyValidatorsExpected{
				Error:    true,
				ErrorMsg: `failed to parse validators for property 'testProp': validator entry 0 regex 0: unknown regex statement type "exclude" (expected one of: include)`,
			},
		},
		{
			Name: "test_meta_min_wrong_type",
			Input: setPropertyValidatorsInput{
				MetaDetails: map[string]any{
					"validators": []any{
						map[string]any{"min": "low", "max": float64(7)},
					},
				},
			},
			Expected: setPropertyValidatorsExpected{
				Error:    true,
				ErrorMsg: "failed to parse validators for property 'testProp': validator entry 0: expected min to be a number, got string",
			},
		},
		{
			Name: "test_meta_min_not_integer",
			Input: setPropertyValidatorsInput{
				MetaDetails: map[string]any{
					"validators": []any{
						map[string]any{"min": float64(0.5), "max": float64(7)},
					},
				},
			},
			Expected: setPropertyValidatorsExpected{
				Error:    true,
				ErrorMsg: "failed to parse validators for property 'testProp': validator entry 0: expected min to be an integer, got 0.5",
			},
		},
		{
			Name: "test_meta_min_not_integer",
			Input: setPropertyValidatorsInput{
				MetaDetails: map[string]any{
					"validators": []any{
						map[string]any{"min": float64(0.5), "max": float64(7)},
					},
				},
			},
			Expected: setPropertyValidatorsExpected{
				Error:    true,
				ErrorMsg: "failed to parse validators for property 'testProp': validator entry 0: expected min to be an integer, got 0.5",
			},
		},
		{
			Name: "test_definition_only_meta_missing",
			Input: setPropertyValidatorsInput{
				PropertyDefinition: PropertyDefinition{
					Validators: []ValidatorDefinition{
						{Min: 1, Max: 10},
					},
				},
			},
			Expected: setPropertyValidatorsExpected{
				Validators: []Validator{{Min: 1, Max: 10}},
			},
		},
		{
			Name: "test_definition_replaces_meta",
			Input: setPropertyValidatorsInput{
				PropertyDefinition: PropertyDefinition{
					Validators: []ValidatorDefinition{
						{Min: 1, Max: 10, RegexList: []RegexStatementDefinition{{Regex: "^x$", Type: "include"}}},
					},
				},
				MetaDetails: map[string]any{
					"validators": []any{
						map[string]any{"min": float64(0), "max": float64(7)},
					},
				},
			},
			Expected: setPropertyValidatorsExpected{
				Validators: []Validator{{
					Min: 1, Max: 10,
					RegexList: []RegexStatement{{Regex: "^x$", Type: Include}},
				}},
			},
		},
		{
			Name: "test_definition_unknown_regex_type",
			Input: setPropertyValidatorsInput{
				PropertyDefinition: PropertyDefinition{
					Validators: []ValidatorDefinition{
						{Min: 0, Max: 1, RegexList: []RegexStatementDefinition{{Regex: ".+", Type: "exclude"}}},
					},
				},
			},
			Expected: setPropertyValidatorsExpected{
				Error:    true,
				ErrorMsg: `failed to parse validators for property 'testProp': validator entry 0 regex 0: unknown regex statement type "exclude" (expected one of: include)`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPropertyValidatorsInput)
			expected := testCase.Expected.(setPropertyValidatorsExpected)

			property := &Property{
				PropertyName:       "testProp",
				propertyDefinition: input.PropertyDefinition,
				metaDetails:        input.MetaDetails,
			}

			err := property.setValidators()

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, expected.Validators, property.Validators, test.MessageEqual(expected.Validators, property.Validators, testCase.Name))
			}
		})
	}
}

type setPropertySupportedVersionsInput struct {
	MetaDetails        map[string]any
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
				MetaDetails:        map[string]any{},
				PropertyDefinition: PropertyDefinition{},
			},
			Expected: setPropertySupportedVersionsExpected{
				Nil: true,
			},
		},
		{
			Name: "test_versions_from_meta_file",
			Input: setPropertySupportedVersionsInput{
				MetaDetails:        map[string]any{"versions": "4.2(7f)-"},
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
				MetaDetails:        map[string]any{"versions": "4.2(7f)-"},
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
				MetaDetails:        map[string]any{},
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
				MetaDetails:        map[string]any{"versions": "3.2(10e)-4.2(7f)"},
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
				MetaDetails:        map[string]any{"versions": "3.2(10e)-3.2(10g),4.2(7f)-"},
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
				MetaDetails:        map[string]any{"versions": "invalid"},
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
				MetaDetails:        map[string]any{"versions": "4.2(7f)-,invalid"},
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

type setPropertyValidValuesInput struct {
	PropertyDefinition PropertyDefinition
	MetaDetails        map[string]any
}

type setPropertyValidValuesExpected struct {
	ValidValues ValidValues
	Error       bool
	ErrorMsg    string
	Warning     string
}

func TestPropertySetValidValues(t *testing.T) {
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:  "test_meta_missing",
			Input: setPropertyValidValuesInput{},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{},
			},
		},
		{
			Name: "test_meta_simple_enum_skips_default_value",
			Input: setPropertyValidValuesInput{
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"localName": "defaultValue", "value": "1"},
						map[string]any{"localName": "level1", "value": "1"},
						map[string]any{"localName": "level3", "value": "3"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{
					"1": ValidValue{LocalName: "level1"},
					"3": ValidValue{LocalName: "level3"},
				},
			},
		},
		{
			Name: "test_meta_only_default_value",
			Input: setPropertyValidValuesInput{
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"localName": "defaultValue", "value": "1"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{},
			},
		},
		{
			Name: "test_meta_bitmask",
			Input: setPropertyValidValuesInput{
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"localName": "read", "value": "1"},
						map[string]any{"localName": "write", "value": "2"},
						map[string]any{"localName": "execute", "value": "4"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{
					"1": ValidValue{LocalName: "read"},
					"2": ValidValue{LocalName: "write"},
					"4": ValidValue{LocalName: "execute"},
				},
			},
		},
		{
			Name: "test_meta_wrong_top_type",
			Input: setPropertyValidValuesInput{
				MetaDetails: map[string]any{"validValues": "not-a-list"},
			},
			Expected: setPropertyValidValuesExpected{
				Error:    true,
				ErrorMsg: "failed to parse validValues for property 'testProp': expected validValues to be a list, got string",
			},
		},
		{
			Name: "test_meta_entry_wrong_type",
			Input: setPropertyValidValuesInput{
				MetaDetails: map[string]any{"validValues": []any{42}},
			},
			Expected: setPropertyValidValuesExpected{
				Error:    true,
				ErrorMsg: "failed to parse validValues for property 'testProp': expected validValues entry 0 to be a map, got int",
			},
		},
		{
			Name: "test_meta_entry_missing_value",
			Input: setPropertyValidValuesInput{
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"localName": "level1"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				Error:    true,
				ErrorMsg: "failed to parse validValues for property 'testProp': validValues entry 0 is missing or has non-string localName/value",
			},
		},
		{
			Name: "test_meta_entry_missing_local_name",
			Input: setPropertyValidValuesInput{
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"value": "1"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				Error:    true,
				ErrorMsg: "failed to parse validValues for property 'testProp': validValues entry 0 is missing or has non-string localName/value",
			},
		},
		{
			Name: "test_meta_duplicate_values",
			Input: setPropertyValidValuesInput{
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"localName": "alpha", "value": "1"},
						map[string]any{"localName": "beta", "value": "1"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{
					"1":    {LocalName: "alpha"},
					"beta": {LocalName: "beta"},
				},
				Warning: `Duplicate validValues value "1" for property "testProp": keeping localName "alpha" under value key, registering alias "beta" under its localName key.`,
			},
		},
		{
			Name: "test_meta_duplicate_value_and_localname_collision",
			Input: setPropertyValidValuesInput{
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"localName": "alpha", "value": "1"},
						map[string]any{"localName": "beta", "value": "1"},
						map[string]any{"localName": "beta", "value": "1"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{
					"1":    {LocalName: "alpha"},
					"beta": {LocalName: "beta"},
				},
				Warning: `Duplicate validValues value "1" for property "testProp": keeping localName "alpha", skipping alias "beta" (localName key already in use).`,
			},
		},
		{
			Name: "test_definition_remove_only",
			Input: setPropertyValidValuesInput{
				PropertyDefinition: PropertyDefinition{
					RemoveValidValues: []string{"level1"},
				},
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"localName": "level1", "value": "1"},
						map[string]any{"localName": "level3", "value": "3"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{
					"3": ValidValue{LocalName: "level3"},
				},
			},
		},
		{
			Name: "test_definition_remove_non_existent_warns",
			Input: setPropertyValidValuesInput{
				PropertyDefinition: PropertyDefinition{
					RemoveValidValues: []string{"missing"},
				},
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"localName": "level1", "value": "1"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{
					"1": ValidValue{LocalName: "level1"},
				},
				Warning: `RemoveValidValues "missing" not found in meta for property "testProp"`,
			},
		},
		{
			Name: "test_definition_add_only_meta_missing",
			Input: setPropertyValidValuesInput{
				PropertyDefinition: PropertyDefinition{
					AddValidValues: []string{"custom"},
				},
			},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{
					"custom": ValidValue{LocalName: "custom"},
				},
			},
		},
		{
			Name: "test_definition_add_appends_to_meta",
			Input: setPropertyValidValuesInput{
				PropertyDefinition: PropertyDefinition{
					AddValidValues: []string{"extra"},
				},
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"localName": "level1", "value": "1"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{
					"1":     ValidValue{LocalName: "level1"},
					"extra": ValidValue{LocalName: "extra"},
				},
			},
		},
		{
			Name: "test_definition_add_overlaps_meta_warns",
			Input: setPropertyValidValuesInput{
				PropertyDefinition: PropertyDefinition{
					AddValidValues: []string{"1"},
				},
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"localName": "level1", "value": "1"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{
					"1": ValidValue{LocalName: "1"},
				},
				Warning: `AddValidValues "1" already present in meta for property "testProp"`,
			},
		},
		{
			Name: "test_definition_add_and_remove",
			Input: setPropertyValidValuesInput{
				PropertyDefinition: PropertyDefinition{
					AddValidValues:    []string{"custom"},
					RemoveValidValues: []string{"level1"},
				},
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"localName": "level1", "value": "1"},
						map[string]any{"localName": "level3", "value": "3"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{
					"3":      ValidValue{LocalName: "level3"},
					"custom": ValidValue{LocalName: "custom"},
				},
			},
		},
		{
			Name: "test_definition_add_and_remove_overlap_warns",
			Input: setPropertyValidValuesInput{
				PropertyDefinition: PropertyDefinition{
					AddValidValues:    []string{"level1"},
					RemoveValidValues: []string{"level1"},
				},
				MetaDetails: map[string]any{
					"validValues": []any{
						map[string]any{"localName": "level1", "value": "1"},
					},
				},
			},
			Expected: setPropertyValidValuesExpected{
				ValidValues: ValidValues{
					"level1": ValidValue{LocalName: "level1"},
				},
				Warning: `AddValidValues "level1" also listed in RemoveValidValues for property "testProp"`,
			},
		},
	}

	// Capture warnings via the package logger; restored after the test.
	var logBuffer bytes.Buffer
	genLogger.SetOutputForTesting(&logBuffer)
	genLogger.SetLogLevel("WARN")
	defer func() {
		genLogger.SetOutputForTesting(os.Stdout)
	}()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(setPropertyValidValuesInput)
			expected := testCase.Expected.(setPropertyValidValuesExpected)

			logBuffer.Reset()

			property := &Property{
				PropertyName:       "testProp",
				propertyDefinition: input.PropertyDefinition,
				metaDetails:        input.MetaDetails,
			}

			err := property.setValidValues()

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
				return
			}

			assert.NoError(t, err, test.MessageUnexpectedError(err))
			assert.Equal(t, expected.ValidValues, property.ValidValues, test.MessageEqual(expected.ValidValues, property.ValidValues, testCase.Name))

			logOutput := logBuffer.String()
			if expected.Warning == "" {
				assert.False(t, strings.Contains(logOutput, "WARN:"), "unexpected warning logged: %s", logOutput)
			} else {
				assert.Contains(t, logOutput, expected.Warning, test.MessageEqual(expected.Warning, logOutput, "warning log message"))
			}
		})
	}
}

func TestValidValuesMethods(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	vv := ValidValues{
		"3": ValidValue{LocalName: "level3"},
		"1": ValidValue{LocalName: "level1"},
		"2": ValidValue{LocalName: "level2"},
	}

	assert.Equal(t, []string{"level1", "level2", "level3"}, vv.LocalNamesList())
	assert.Equal(t, []string{"1", "2", "3"}, vv.ValuesList())
	assert.Equal(t, map[string]string{"1": "level1", "2": "level2", "3": "level3"}, vv.ValueLocalNameMap())

	empty := ValidValues{}
	assert.Equal(t, []string{}, empty.LocalNamesList())
	assert.Equal(t, []string{}, empty.ValuesList())
	assert.Equal(t, map[string]string{}, empty.ValueLocalNameMap())
}

type setValueTypeInput struct {
	PropertyDefinition PropertyDefinition
	MetaDetails        map[string]any
	ValidValues        ValidValues
	Validators         []Validator
}

type setValueTypeExpected struct {
	ValueType ValueTypeEnum
	Error     bool
	ErrorMsg  string
	Warning   string
}

func TestPropertySetValueType(t *testing.T) {
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_meta_bitmask_to_set",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "bitmask"},
			},
			Expected: setValueTypeExpected{ValueType: Set},
		},
		{
			Name: "test_meta_string_to_string",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "string"},
			},
			Expected: setValueTypeExpected{ValueType: String},
		},
		{
			Name: "test_meta_enum_no_warn",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "enum"},
			},
			Expected: setValueTypeExpected{ValueType: String},
		},
		{
			Name: "test_meta_auto_no_warn",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "auto"},
			},
			Expected: setValueTypeExpected{ValueType: String},
		},
		{
			Name: "test_meta_number_no_warn",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "number"},
			},
			Expected: setValueTypeExpected{ValueType: String},
		},
		{
			Name: "test_meta_boolean_no_warn",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "boolean"},
			},
			Expected: setValueTypeExpected{ValueType: String},
		},
		{
			Name: "test_meta_password_no_warn",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "password"},
			},
			Expected: setValueTypeExpected{ValueType: String},
		},
		{
			Name:     "test_meta_missing_uitype_defaults_string",
			Input:    setValueTypeInput{MetaDetails: map[string]any{}},
			Expected: setValueTypeExpected{ValueType: String},
		},
		{
			Name: "test_meta_unknown_uitype_warns",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "weirdtype"},
			},
			Expected: setValueTypeExpected{
				ValueType: String,
				Warning:   `Unmapped meta uiType "weirdtype" for property "testProp"`,
			},
		},
		{
			Name: "test_meta_validate_as_ip_to_ip_address",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "string", "validateAsIPv4OrIPv6": true},
			},
			Expected: setValueTypeExpected{ValueType: IpAddress},
		},
		{
			Name: "test_meta_bitmask_outranks_ip",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "bitmask", "validateAsIPv4OrIPv6": true},
			},
			Expected: setValueTypeExpected{ValueType: Set},
		},
		{
			Name: "test_meta_ip_outranks_semantic_equality",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "string", "validateAsIPv4OrIPv6": true},
				ValidValues: ValidValues{"1": ValidValue{LocalName: "one"}},
				Validators:  []Validator{{Min: 0, Max: 10}},
			},
			Expected: setValueTypeExpected{ValueType: IpAddress},
		},
		{
			Name: "test_validators_and_valid_values_to_semantic_equality",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "number"},
				ValidValues: ValidValues{"22": ValidValue{LocalName: "ssh"}},
				Validators:  []Validator{{Min: 0, Max: 65535}},
			},
			Expected: setValueTypeExpected{ValueType: SemanticEquality},
		},
		{
			Name: "test_only_validators_no_semantic_equality",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "number"},
				Validators:  []Validator{{Min: 0, Max: 65535}},
			},
			Expected: setValueTypeExpected{ValueType: String},
		},
		{
			Name: "test_only_valid_values_no_semantic_equality",
			Input: setValueTypeInput{
				MetaDetails: map[string]any{"uitype": "enum"},
				ValidValues: ValidValues{"1": ValidValue{LocalName: "one"}},
			},
			Expected: setValueTypeExpected{ValueType: String},
		},
		{
			Name: "test_definition_override_set",
			Input: setValueTypeInput{
				PropertyDefinition: PropertyDefinition{ValueType: Set},
				MetaDetails:        map[string]any{"uitype": "string"},
			},
			Expected: setValueTypeExpected{ValueType: Set},
		},
		{
			Name: "test_definition_override_ip_address",
			Input: setValueTypeInput{
				PropertyDefinition: PropertyDefinition{ValueType: IpAddress},
				MetaDetails:        map[string]any{},
			},
			Expected: setValueTypeExpected{ValueType: IpAddress},
		},
		{
			Name: "test_definition_override_semantic_equality",
			Input: setValueTypeInput{
				PropertyDefinition: PropertyDefinition{ValueType: SemanticEquality},
				MetaDetails:        map[string]any{"uitype": "bitmask"},
			},
			Expected: setValueTypeExpected{ValueType: SemanticEquality},
		},
	}

	// Capture warnings via the package logger; restored after the test.
	var logBuffer bytes.Buffer
	genLogger.SetOutputForTesting(&logBuffer)
	genLogger.SetLogLevel("WARN")
	defer func() {
		genLogger.SetOutputForTesting(os.Stdout)
	}()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(setValueTypeInput)
			expected := testCase.Expected.(setValueTypeExpected)

			logBuffer.Reset()

			property := &Property{
				PropertyName:       "testProp",
				propertyDefinition: input.PropertyDefinition,
				metaDetails:        input.MetaDetails,
				ValidValues:        input.ValidValues,
				Validators:         input.Validators,
			}

			err := property.setValueType()

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
				return
			}

			assert.NoError(t, err, test.MessageUnexpectedError(err))
			assert.Equal(t, expected.ValueType, property.ValueType, test.MessageEqual(expected.ValueType, property.ValueType, testCase.Name))

			logOutput := logBuffer.String()
			if expected.Warning == "" {
				assert.False(t, strings.Contains(logOutput, "WARN:"), "unexpected warning logged: %s", logOutput)
			} else {
				assert.Contains(t, logOutput, expected.Warning, test.MessageEqual(expected.Warning, logOutput, "warning log message"))
			}
		})
	}
}

type setTestValuesInput struct {
	PropertyName       string
	AttributeName      string
	PropertyDefinition PropertyDefinition
	ValidValues        ValidValues
	ValueType          ValueTypeEnum
	Documentation      PropertyDocumentation
	Required           bool
	Optional           bool
	ReadOnly           bool
}

type setTestValuesExpected struct {
	IgnoreInTest   bool
	TestValuesNil  bool
	CreateValues   []string
	UpdateValues   []string
	DefaultInclude []bool
	DefaultAssert  []string
	CreateTypes    []ValueRenderTypeEnum
	ForceNewValues []string
}

func TestSetTestValues(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	boolTrue := true
	boolFalse := false

	testCases := []test.TestCase{
		{
			Name: "test_ignore_in_test",
			Input: setTestValuesInput{
				PropertyName: "testProp",
				PropertyDefinition: PropertyDefinition{
					TestConfig: TestConfigDefinition{
						IgnoreInTest: true,
					},
				},
			},
			Expected: setTestValuesExpected{
				IgnoreInTest:  true,
				TestValuesNil: true,
			},
		},
		{
			Name: "test_read_only_property",
			Input: setTestValuesInput{
				PropertyName: "pcTag",
				ReadOnly:     true,
			},
			Expected: setTestValuesExpected{
				TestValuesNil: true,
			},
		},
		{
			Name: "test_explicit_config",
			Input: setTestValuesInput{
				PropertyName: "testProp",
				PropertyDefinition: PropertyDefinition{
					TestConfig: TestConfigDefinition{
						Create: []TestValueEntryDefinition{
							{ConfigValue: "create_val", ConfigInclude: &boolTrue, AssertValue: "assert_val", ValueType: StringValue},
						},
						Update: []TestValueEntryDefinition{
							{ConfigValue: "update_val", ConfigInclude: &boolTrue},
						},
						Default: []TestValueEntryDefinition{
							{ConfigValue: "", ConfigInclude: &boolFalse, AssertValue: "default_val"},
						},
						ForceNew: []TestValueEntryDefinition{
							{ConfigValue: "force_new_val", ConfigInclude: &boolTrue},
						},
					},
				},
			},
			Expected: setTestValuesExpected{
				CreateValues:   []string{"create_val"},
				UpdateValues:   []string{"update_val"},
				DefaultInclude: []bool{false},
				DefaultAssert:  []string{"default_val"},
				ForceNewValues: []string{"force_new_val"},
			},
		},
		{
			Name: "test_explicit_reference_type",
			Input: setTestValuesInput{
				PropertyName: "tDn",
				PropertyDefinition: PropertyDefinition{
					TestConfig: TestConfigDefinition{
						Create: []TestValueEntryDefinition{
							{ConfigValue: "aci_tenant.test.id", ConfigInclude: &boolTrue, ValueType: ReferenceValue},
						},
					},
				},
			},
			Expected: setTestValuesExpected{
				CreateValues: []string{"aci_tenant.test.id"},
				CreateTypes:  []ValueRenderTypeEnum{ReferenceValue},
			},
		},
		{
			Name: "test_auto_derive_from_valid_values",
			Input: setTestValuesInput{
				PropertyName: "mode",
				Optional:     true,
				ValidValues: ValidValues{
					"1": {LocalName: "regular"},
					"2": {LocalName: "native"},
					"3": {LocalName: "untagged"},
				},
				Documentation: PropertyDocumentation{
					DefaultValues: []DefaultValue{{Value: "regular"}},
				},
			},
			Expected: setTestValuesExpected{
				CreateValues:   []string{"native"},
				UpdateValues:   []string{"regular"},
				DefaultInclude: []bool{false},
				DefaultAssert:  []string{"regular"},
				ForceNewValues: []string{"native"},
			},
		},
		{
			Name: "test_auto_derive_from_valid_values_single_value",
			Input: setTestValuesInput{
				PropertyName: "mode",
				Optional:     true,
				ValidValues: ValidValues{
					"1": {LocalName: "enabled"},
				},
			},
			Expected: setTestValuesExpected{
				CreateValues:   []string{"enabled"},
				UpdateValues:   []string{"enabled"},
				DefaultInclude: []bool{false},
				DefaultAssert:  []string{""},
			},
		},
		{
			Name: "test_auto_derive_free_form_optional",
			Input: setTestValuesInput{
				PropertyName:  "descr",
				AttributeName: "description",
				Optional:      true,
			},
			Expected: setTestValuesExpected{
				CreateValues:   []string{"description_1"},
				UpdateValues:   []string{"description_2"},
				DefaultInclude: []bool{false},
				DefaultAssert:  []string{""},
				ForceNewValues: []string{"description_1"},
			},
		},
		{
			Name: "test_auto_derive_free_form_required",
			Input: setTestValuesInput{
				PropertyName:  "name",
				AttributeName: "name",
				Required:      true,
			},
			Expected: setTestValuesExpected{
				CreateValues:   []string{"name_1"},
				UpdateValues:   []string{"name_2"},
				DefaultInclude: []bool{true},
				DefaultAssert:  []string{"name_1"},
				ForceNewValues: []string{"name_1"},
			},
		},
		{
			Name: "test_auto_derive_set_type",
			Input: setTestValuesInput{
				PropertyName:  "matchT",
				AttributeName: "match_type",
				Optional:      true,
				ValueType:     Set,
				ValidValues: ValidValues{
					"1": {LocalName: "all"},
					"2": {LocalName: "atleast_one"},
					"3": {LocalName: "atmost_one"},
					"4": {LocalName: "none"},
				},
			},
			Expected: setTestValuesExpected{
				CreateValues:   []string{"all", "atleast_one"},
				UpdateValues:   []string{"atmost_one", "none"},
				DefaultInclude: []bool{false},
				DefaultAssert:  []string{""},
				ForceNewValues: []string{"all", "atleast_one"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setTestValuesInput)
			expected := testCase.Expected.(setTestValuesExpected)

			property := &Property{
				PropertyName:       input.PropertyName,
				AttributeName:      input.AttributeName,
				propertyDefinition: input.PropertyDefinition,
				ValidValues:        input.ValidValues,
				ValueType:          input.ValueType,
				Documentation:      input.Documentation,
				Required:           input.Required,
				Optional:           input.Optional,
				ReadOnly:           input.ReadOnly,
			}

			property.setTestValues()

			assert.Equal(t, expected.IgnoreInTest, property.IgnoreInTest, test.MessageEqual(expected.IgnoreInTest, property.IgnoreInTest, testCase.Name))

			if expected.TestValuesNil {
				assert.Nil(t, property.TestValues, testCase.Name+": TestValues should be nil")
				return
			}

			assert.NotNil(t, property.TestValues, testCase.Name+": TestValues should not be nil")

			if expected.CreateValues != nil {
				assert.Len(t, property.TestValues.Create, len(expected.CreateValues), testCase.Name+": Create length mismatch")
				for i, v := range expected.CreateValues {
					assert.Equal(t, v, property.TestValues.Create[i].ConfigValue, test.MessageEqual(v, property.TestValues.Create[i].ConfigValue, testCase.Name))
				}
			}

			if expected.UpdateValues != nil {
				assert.Len(t, property.TestValues.Update, len(expected.UpdateValues), testCase.Name+": Update length mismatch")
				for i, v := range expected.UpdateValues {
					assert.Equal(t, v, property.TestValues.Update[i].ConfigValue, test.MessageEqual(v, property.TestValues.Update[i].ConfigValue, testCase.Name))
				}
			}

			if expected.DefaultInclude != nil {
				assert.Len(t, property.TestValues.Default, len(expected.DefaultInclude), testCase.Name+": Default length mismatch")
				for i, inc := range expected.DefaultInclude {
					assert.Equal(t, inc, property.TestValues.Default[i].ConfigInclude, test.MessageEqual(inc, property.TestValues.Default[i].ConfigInclude, testCase.Name))
				}
			}

			if expected.DefaultAssert != nil {
				for i, v := range expected.DefaultAssert {
					assert.Equal(t, v, property.TestValues.Default[i].AssertValue, test.MessageEqual(v, property.TestValues.Default[i].AssertValue, testCase.Name))
				}
			}

			if expected.CreateTypes != nil {
				for i, vt := range expected.CreateTypes {
					assert.Equal(t, vt, property.TestValues.Create[i].ValueType, test.MessageEqual(vt, property.TestValues.Create[i].ValueType, testCase.Name))
				}
			}

			if expected.ForceNewValues != nil {
				assert.Len(t, property.TestValues.ForceNew, len(expected.ForceNewValues), testCase.Name+": ForceNew length mismatch")
				for i, v := range expected.ForceNewValues {
					assert.Equal(t, v, property.TestValues.ForceNew[i].ConfigValue, test.MessageEqual(v, property.TestValues.ForceNew[i].ConfigValue, testCase.Name))
				}
			}
		})
	}
}

type convertTestValueEntriesExpected struct {
	Nil           bool
	Length        int
	ConfigInclude []bool
	AssertValues  []string
}

func TestConvertTestValueEntries(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	boolTrue := true

	testCases := []test.TestCase{
		{
			Name:  "test_nil_input",
			Input: []TestValueEntryDefinition(nil),
			Expected: convertTestValueEntriesExpected{
				Nil: true,
			},
		},
		{
			Name:  "test_empty_input",
			Input: []TestValueEntryDefinition{},
			Expected: convertTestValueEntriesExpected{
				Nil: true,
			},
		},
		{
			Name: "test_nil_config_include_defaults_to_true",
			Input: []TestValueEntryDefinition{
				{ConfigValue: "value1", ConfigInclude: nil},
			},
			Expected: convertTestValueEntriesExpected{
				Length:        1,
				ConfigInclude: []bool{true},
				AssertValues:  []string{"value1"},
			},
		},
		{
			Name: "test_explicit_config_include",
			Input: []TestValueEntryDefinition{
				{ConfigValue: "value1", ConfigInclude: &boolTrue, AssertValue: "custom_assert"},
			},
			Expected: convertTestValueEntriesExpected{
				Length:        1,
				ConfigInclude: []bool{true},
				AssertValues:  []string{"custom_assert"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.([]TestValueEntryDefinition)
			expected := testCase.Expected.(convertTestValueEntriesExpected)

			result := convertTestValueEntries(input)

			if expected.Nil {
				assert.Nil(t, result, testCase.Name+": result should be nil")
				return
			}

			assert.Len(t, result, expected.Length, test.MessageEqual(expected.Length, len(result), testCase.Name))

			if expected.ConfigInclude != nil {
				for i, inc := range expected.ConfigInclude {
					assert.Equal(t, inc, result[i].ConfigInclude, test.MessageEqual(inc, result[i].ConfigInclude, testCase.Name))
				}
			}

			if expected.AssertValues != nil {
				for i, v := range expected.AssertValues {
					assert.Equal(t, v, result[i].AssertValue, test.MessageEqual(v, result[i].AssertValue, testCase.Name))
				}
			}
		})
	}
}

type setLegacyTestValuesInput struct {
	PropertyName       string
	AttributeName      string
	ValueType          ValueTypeEnum
	StateUpgradeValues map[int]StateUpgradeValue
	ExistingTestValues *TestValues
	PropertyDefinition PropertyDefinition
	IgnoreInTest       bool
	ReadOnly           bool
}

type setLegacyTestValuesExpected struct {
	LegacyNil      bool
	LegacyValues   []string
	LegacyAsserts  []string
	LegacyIncludes []bool
	LegacyTypes    []ValueRenderTypeEnum
	LegacyVersions []bool // true = Versions pointer preserved on that entry
	Warning        string
}

// TestSetLegacyTestValues drives Property.setLegacyTestValues across every
// branch: gating (nil TestValues, IgnoreInTest, ReadOnly), YAML override
// precedence (even when type would have diverged), the no-alias and
// all-removed skip paths, the same-name skip path, the divergent-type warn,
// and the clone path with Versions preservation (which generateForceNew
// currently drops).
func TestSetLegacyTestValues(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	boolFalse := false
	versionsAll, _ := NewVersions("1.0(1e)-")

	testCases := []test.TestCase{
		{
			Name: "test_skip_when_test_values_nil",
			Input: setLegacyTestValuesInput{
				PropertyName:       "name",
				AttributeName:      "name",
				StateUpgradeValues: map[int]StateUpgradeValue{0: {AttributeName: "old_name", Type: String, Status: Functioning}},
			},
			Expected: setLegacyTestValuesExpected{LegacyNil: true},
		},
		{
			Name: "test_skip_when_ignore_in_test",
			Input: setLegacyTestValuesInput{
				PropertyName:       "name",
				AttributeName:      "name",
				IgnoreInTest:       true,
				StateUpgradeValues: map[int]StateUpgradeValue{0: {AttributeName: "old_name", Type: String, Status: Functioning}},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "name_1", ConfigInclude: true, AssertValue: "name_1", ValueType: StringValue}}},
			},
			Expected: setLegacyTestValuesExpected{LegacyNil: true},
		},
		{
			Name: "test_skip_when_read_only",
			Input: setLegacyTestValuesInput{
				PropertyName:       "pcTag",
				AttributeName:      "pc_tag",
				ReadOnly:           true,
				StateUpgradeValues: map[int]StateUpgradeValue{0: {AttributeName: "old_pc_tag", Type: String, Status: Functioning}},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "1", ConfigInclude: true, AssertValue: "1", ValueType: StringValue}}},
			},
			Expected: setLegacyTestValuesExpected{LegacyNil: true},
		},
		{
			Name: "test_yaml_override_wins_even_when_type_diverges",
			Input: setLegacyTestValuesInput{
				PropertyName:  "tags",
				AttributeName: "tags",
				ValueType:     Set,
				StateUpgradeValues: map[int]StateUpgradeValue{
					0: {AttributeName: "tag_string", Type: String, Status: Functioning},
				},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "tag_a", ConfigInclude: true, AssertValue: "tag_a", ValueType: StringValue}}},
				PropertyDefinition: PropertyDefinition{
					TestConfig: TestConfigDefinition{
						Legacy: []TestValueEntryDefinition{
							{ConfigValue: "legacy_value", ConfigInclude: nil, AssertValue: "legacy_assert", ValueType: StringValue},
						},
					},
				},
			},
			Expected: setLegacyTestValuesExpected{
				LegacyValues:   []string{"legacy_value"},
				LegacyAsserts:  []string{"legacy_assert"},
				LegacyIncludes: []bool{true},
				LegacyTypes:    []ValueRenderTypeEnum{StringValue},
			},
		},
		{
			Name: "test_yaml_override_with_explicit_config_include_false",
			Input: setLegacyTestValuesInput{
				PropertyName:       "descr",
				AttributeName:      "description",
				ValueType:          String,
				StateUpgradeValues: map[int]StateUpgradeValue{0: {AttributeName: "old_description", Type: String, Status: Functioning}},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "description_1", ConfigInclude: true, AssertValue: "description_1", ValueType: StringValue}}},
				PropertyDefinition: PropertyDefinition{
					TestConfig: TestConfigDefinition{
						Legacy: []TestValueEntryDefinition{
							{ConfigValue: "", ConfigInclude: &boolFalse, AssertValue: "server_default"},
						},
					},
				},
			},
			Expected: setLegacyTestValuesExpected{
				LegacyValues:   []string{""},
				LegacyAsserts:  []string{"server_default"},
				LegacyIncludes: []bool{false},
			},
		},
		{
			Name: "test_skip_when_no_state_upgrade_values",
			Input: setLegacyTestValuesInput{
				PropertyName:       "name",
				AttributeName:      "name",
				ValueType:          String,
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "name_1", ConfigInclude: true, AssertValue: "name_1", ValueType: StringValue}}},
			},
			Expected: setLegacyTestValuesExpected{LegacyNil: true},
		},
		{
			Name: "test_skip_when_all_entries_removed",
			Input: setLegacyTestValuesInput{
				PropertyName:  "name",
				AttributeName: "name",
				ValueType:     String,
				StateUpgradeValues: map[int]StateUpgradeValue{
					0: {AttributeName: "old_name", Type: String, Status: Removed},
				},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "name_1", ConfigInclude: true, AssertValue: "name_1", ValueType: StringValue}}},
			},
			Expected: setLegacyTestValuesExpected{LegacyNil: true},
		},
		{
			Name: "test_skip_when_only_same_name_entry",
			Input: setLegacyTestValuesInput{
				PropertyName:  "name",
				AttributeName: "name",
				ValueType:     String,
				StateUpgradeValues: map[int]StateUpgradeValue{
					0: {AttributeName: "name", Type: String, Status: Functioning},
				},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "name_1", ConfigInclude: true, AssertValue: "name_1", ValueType: StringValue}}},
			},
			Expected: setLegacyTestValuesExpected{LegacyNil: true},
		},
		{
			Name: "test_skip_when_only_removed_distinct_name",
			Input: setLegacyTestValuesInput{
				PropertyName:  "name",
				AttributeName: "name",
				ValueType:     String,
				StateUpgradeValues: map[int]StateUpgradeValue{
					0: {AttributeName: "old_name", Type: String, Status: Removed},
					1: {AttributeName: "name", Type: String, Status: Functioning},
				},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "name_1", ConfigInclude: true, AssertValue: "name_1", ValueType: StringValue}}},
			},
			Expected: setLegacyTestValuesExpected{LegacyNil: true},
		},
		{
			Name: "test_clone_create_when_functioning_distinct_name_matching_type",
			Input: setLegacyTestValuesInput{
				PropertyName:  "name",
				AttributeName: "name",
				ValueType:     String,
				StateUpgradeValues: map[int]StateUpgradeValue{
					0: {AttributeName: "old_name", Type: String, Status: Functioning},
				},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "name_1", ConfigInclude: true, AssertValue: "name_1", ValueType: StringValue}}},
			},
			Expected: setLegacyTestValuesExpected{
				LegacyValues:   []string{"name_1"},
				LegacyAsserts:  []string{"name_1"},
				LegacyIncludes: []bool{true},
				LegacyTypes:    []ValueRenderTypeEnum{StringValue},
			},
		},
		{
			Name: "test_clone_create_when_frozen_distinct_name_matching_type",
			Input: setLegacyTestValuesInput{
				PropertyName:  "descr",
				AttributeName: "description",
				ValueType:     String,
				StateUpgradeValues: map[int]StateUpgradeValue{
					0: {AttributeName: "old_description", Type: String, Status: Frozen},
				},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "description_1", ConfigInclude: true, AssertValue: "description_1", ValueType: StringValue}}},
			},
			Expected: setLegacyTestValuesExpected{
				LegacyValues:   []string{"description_1"},
				LegacyAsserts:  []string{"description_1"},
				LegacyIncludes: []bool{true},
				LegacyTypes:    []ValueRenderTypeEnum{StringValue},
			},
		},
		{
			Name: "test_clone_preserves_reference_value_for_parent_dn",
			Input: setLegacyTestValuesInput{
				PropertyName:  "parentDn",
				AttributeName: "parent_dn",
				ValueType:     String,
				StateUpgradeValues: map[int]StateUpgradeValue{
					0: {AttributeName: "tenant_dn", Type: String, Status: Functioning},
				},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "aci_tenant.test.id", ConfigInclude: true, AssertValue: "aci_tenant.test.id", ValueType: ReferenceValue}}},
			},
			Expected: setLegacyTestValuesExpected{
				LegacyValues: []string{"aci_tenant.test.id"},
				LegacyTypes:  []ValueRenderTypeEnum{ReferenceValue},
			},
		},
		{
			Name: "test_clone_preserves_versions_pointer",
			Input: setLegacyTestValuesInput{
				PropertyName:  "vlanScope",
				AttributeName: "vlan_scope",
				ValueType:     String,
				StateUpgradeValues: map[int]StateUpgradeValue{
					0: {AttributeName: "old_vlan_scope", Type: String, Status: Functioning},
				},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "portlocal", ConfigInclude: true, AssertValue: "portlocal", ValueType: StringValue, Versions: versionsAll}}},
			},
			Expected: setLegacyTestValuesExpected{
				LegacyValues:   []string{"portlocal"},
				LegacyVersions: []bool{true},
			},
		},
		{
			Name: "test_warn_and_skip_when_type_diverges",
			Input: setLegacyTestValuesInput{
				PropertyName:  "tags",
				AttributeName: "tags",
				ValueType:     Set,
				StateUpgradeValues: map[int]StateUpgradeValue{
					0: {AttributeName: "tag_string", Type: String, Status: Functioning},
				},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "tag_a", ConfigInclude: true, AssertValue: "tag_a", ValueType: StringValue}}},
			},
			Expected: setLegacyTestValuesExpected{
				LegacyNil: true,
				Warning:   `Skipping Legacy test value auto-derivation for property "tags"`,
			},
		},
		{
			Name: "test_divergent_only_on_removed_entry_does_not_block",
			Input: setLegacyTestValuesInput{
				PropertyName:  "name",
				AttributeName: "name",
				ValueType:     String,
				StateUpgradeValues: map[int]StateUpgradeValue{
					0: {AttributeName: "old_name", Type: Set, Status: Removed},
					1: {AttributeName: "renamed_name", Type: String, Status: Functioning},
				},
				ExistingTestValues: &TestValues{Create: []TestValueEntry{{ConfigValue: "name_1", ConfigInclude: true, AssertValue: "name_1", ValueType: StringValue}}},
			},
			Expected: setLegacyTestValuesExpected{
				LegacyValues:   []string{"name_1"},
				LegacyIncludes: []bool{true},
			},
		},
	}

	// Capture warnings via the package logger; restored after the test.
	var logBuffer bytes.Buffer
	genLogger.SetOutputForTesting(&logBuffer)
	genLogger.SetLogLevel("WARN")
	defer func() {
		genLogger.SetOutputForTesting(os.Stdout)
	}()

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(setLegacyTestValuesInput)
			expected := testCase.Expected.(setLegacyTestValuesExpected)

			logBuffer.Reset()

			property := &Property{
				PropertyName:       input.PropertyName,
				AttributeName:      input.AttributeName,
				ValueType:          input.ValueType,
				StateUpgradeValues: input.StateUpgradeValues,
				TestValues:         input.ExistingTestValues,
				propertyDefinition: input.PropertyDefinition,
				IgnoreInTest:       input.IgnoreInTest,
				ReadOnly:           input.ReadOnly,
			}

			property.setLegacyTestValues()

			logOutput := logBuffer.String()
			if expected.Warning == "" {
				assert.False(t, strings.Contains(logOutput, "WARN:"), "unexpected warning logged: %s", logOutput)
			} else {
				assert.Contains(t, logOutput, expected.Warning, test.MessageEqual(expected.Warning, logOutput, "warning log message"))
			}

			if expected.LegacyNil {
				if property.TestValues != nil {
					assert.Nil(t, property.TestValues.Legacy, testCase.Name+": Legacy should be nil")
				}
				return
			}

			assert.NotNil(t, property.TestValues, testCase.Name+": TestValues should not be nil")
			assert.NotNil(t, property.TestValues.Legacy, testCase.Name+": Legacy should not be nil")

			if expected.LegacyValues != nil {
				assert.Len(t, property.TestValues.Legacy, len(expected.LegacyValues), testCase.Name+": Legacy length mismatch")
				for i, v := range expected.LegacyValues {
					assert.Equal(t, v, property.TestValues.Legacy[i].ConfigValue, test.MessageEqual(v, property.TestValues.Legacy[i].ConfigValue, testCase.Name))
				}
			}
			if expected.LegacyAsserts != nil {
				for i, v := range expected.LegacyAsserts {
					assert.Equal(t, v, property.TestValues.Legacy[i].AssertValue, test.MessageEqual(v, property.TestValues.Legacy[i].AssertValue, testCase.Name))
				}
			}
			if expected.LegacyIncludes != nil {
				for i, inc := range expected.LegacyIncludes {
					assert.Equal(t, inc, property.TestValues.Legacy[i].ConfigInclude, test.MessageEqual(inc, property.TestValues.Legacy[i].ConfigInclude, testCase.Name))
				}
			}
			if expected.LegacyTypes != nil {
				for i, vt := range expected.LegacyTypes {
					assert.Equal(t, vt, property.TestValues.Legacy[i].ValueType, test.MessageEqual(vt, property.TestValues.Legacy[i].ValueType, testCase.Name))
				}
			}
			if expected.LegacyVersions != nil {
				for i, hasVersions := range expected.LegacyVersions {
					if hasVersions {
						assert.NotNil(t, property.TestValues.Legacy[i].Versions, testCase.Name+": Versions should be preserved")
					} else {
						assert.Nil(t, property.TestValues.Legacy[i].Versions, testCase.Name+": Versions should be nil")
					}
				}
			}
		})
	}
}

type validateStandardBucketsCompleteInput struct {
	PropertyName       string
	PropertyDefinition PropertyDefinition
}

// TestValidateStandardBucketsComplete asserts the all-or-nothing rule on
// test_config: zero buckets and all four buckets are both valid (auto-derive
// and full manual override paths); any 1–3 buckets yields an error that
// names both the present and missing buckets. legacy is independent and
// supplying it alone is allowed.
func TestValidateStandardBucketsComplete(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	entry := []TestValueEntryDefinition{{ConfigValue: "v"}}

	testCases := []test.TestCase{
		{
			Name:  "test_no_buckets_ok",
			Input: validateStandardBucketsCompleteInput{PropertyName: "name", PropertyDefinition: PropertyDefinition{TestConfig: TestConfigDefinition{}}},
			Expected: struct {
				Error           bool
				ContainsPresent []string
				ContainsMissing []string
			}{Error: false},
		},
		{
			Name: "test_all_four_buckets_ok",
			Input: validateStandardBucketsCompleteInput{
				PropertyName: "name",
				PropertyDefinition: PropertyDefinition{TestConfig: TestConfigDefinition{
					Create: entry, Update: entry, Default: entry, ForceNew: entry,
				}},
			},
			Expected: struct {
				Error           bool
				ContainsPresent []string
				ContainsMissing []string
			}{Error: false},
		},
		{
			Name: "test_legacy_only_ok",
			Input: validateStandardBucketsCompleteInput{
				PropertyName:       "name",
				PropertyDefinition: PropertyDefinition{TestConfig: TestConfigDefinition{Legacy: entry}},
			},
			Expected: struct {
				Error           bool
				ContainsPresent []string
				ContainsMissing []string
			}{Error: false},
		},
		{
			Name: "test_create_only_errors",
			Input: validateStandardBucketsCompleteInput{
				PropertyName:       "name",
				PropertyDefinition: PropertyDefinition{TestConfig: TestConfigDefinition{Create: entry}},
			},
			Expected: struct {
				Error           bool
				ContainsPresent []string
				ContainsMissing []string
			}{Error: true, ContainsPresent: []string{"create"}, ContainsMissing: []string{"default", "force_new", "update"}},
		},
		{
			Name: "test_three_buckets_missing_default_errors",
			Input: validateStandardBucketsCompleteInput{
				PropertyName: "name",
				PropertyDefinition: PropertyDefinition{TestConfig: TestConfigDefinition{
					Create: entry, Update: entry, ForceNew: entry,
				}},
			},
			Expected: struct {
				Error           bool
				ContainsPresent []string
				ContainsMissing []string
			}{Error: true, ContainsPresent: []string{"create", "force_new", "update"}, ContainsMissing: []string{"default"}},
		},
		{
			Name: "test_partial_with_legacy_still_errors_on_standard",
			Input: validateStandardBucketsCompleteInput{
				PropertyName: "name",
				PropertyDefinition: PropertyDefinition{TestConfig: TestConfigDefinition{
					Create: entry, Legacy: entry,
				}},
			},
			Expected: struct {
				Error           bool
				ContainsPresent []string
				ContainsMissing []string
			}{Error: true, ContainsPresent: []string{"create"}, ContainsMissing: []string{"default", "force_new", "update"}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(validateStandardBucketsCompleteInput)
			expected := testCase.Expected.(struct {
				Error           bool
				ContainsPresent []string
				ContainsMissing []string
			})

			property := &Property{
				PropertyName:       input.PropertyName,
				propertyDefinition: input.PropertyDefinition,
			}

			err := property.validateStandardBucketsComplete()

			if !expected.Error {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				return
			}

			if assert.Error(t, err, testCase.Name+": expected error") {
				message := err.Error()
				for _, p := range expected.ContainsPresent {
					assert.Contains(t, message, p, test.MessageContains(message, p, testCase.Name))
				}
				for _, m := range expected.ContainsMissing {
					assert.Contains(t, message, m, test.MessageContains(message, m, testCase.Name))
				}
			}
		})
	}
}
