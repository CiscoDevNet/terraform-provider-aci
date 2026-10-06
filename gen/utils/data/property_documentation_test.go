package data

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type setPropertyDescriptionInput struct {
	DefinitionDescription string
	MetaComment           []any
	MetaLabel             string
}

func TestSetPropertyDescription(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_definition_override",
			Input: setPropertyDescriptionInput{
				DefinitionDescription: "Override description.",
				MetaComment:           []any{"Meta comment."},
				MetaLabel:             "Label",
			},
			Expected: "Override description.",
		},
		{
			Name: "test_meta_comment_joined",
			Input: setPropertyDescriptionInput{
				MetaComment: []any{"First part.", "Second part."},
			},
			Expected: "First part. Second part.",
		},
		{
			Name: "test_meta_comment_whitespace_collapsed",
			Input: setPropertyDescriptionInput{
				MetaComment: []any{"A comment\n  with   extra \t whitespace."},
			},
			Expected: "A comment with extra whitespace.",
		},
		{
			Name: "test_meta_label_fallback",
			Input: setPropertyDescriptionInput{
				MetaLabel: "Annotation",
			},
			Expected: "Annotation",
		},
		{
			Name: "test_empty_comment_falls_back_to_label",
			Input: setPropertyDescriptionInput{
				MetaComment: []any{},
				MetaLabel:   "Fallback Label",
			},
			Expected: "Fallback Label",
		},
		{
			Name:     "test_all_empty",
			Input:    setPropertyDescriptionInput{},
			Expected: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPropertyDescriptionInput)
			expected := testCase.Expected.(string)

			metaDetails := map[string]any{}
			if input.MetaComment != nil {
				metaDetails["comment"] = input.MetaComment
			}
			if input.MetaLabel != "" {
				metaDetails["label"] = input.MetaLabel
			}

			p := &Property{
				PropertyName: "testProp",
				propertyDefinition: PropertyDefinition{
					Documentation: ArtifactDocumentationDefinition{
						Description: input.DefinitionDescription,
					},
				},
				metaDetails: metaDetails,
			}

			p.Documentation.setDescription(p)

			assert.Equal(t, expected, p.Documentation.Description, test.MessageEqual(expected, p.Documentation.Description, testCase.Name))
		})
	}
}

type setPropertyDocumentationReferenceInput struct {
	Reference        *PropertyReference
	TargetArtifacts  []ArtifactEnum
	TargetPresent    bool
	TargetDefinition *ClassDefinition
	Description      string
}

func TestSetPropertyDocumentationReferences(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	resourceReference := "a [resource](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3_outside) with `aci_l3_outside.example.name`"
	datasourceReference := "a [datasource](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/data-sources/l3_outside) with `data.aci_l3_outside.example.name`"
	reference := func() *PropertyReference {
		return &PropertyReference{
			ExampleClass:         testClassName("l3extOut"),
			ExampleAttributeName: "name",
		}
	}

	testCases := []test.TestCase{
		{
			Name:     "test_without_reference",
			Input:    setPropertyDocumentationReferenceInput{},
			Expected: "The target name.",
		},
		{
			Name: "test_missing_target_class",
			Input: setPropertyDocumentationReferenceInput{
				Reference: reference(),
			},
			Expected: "The target name.",
		},
		{
			Name: "test_target_without_artifact",
			Input: setPropertyDocumentationReferenceInput{
				Reference:     reference(),
				TargetPresent: true,
			},
			Expected: "The target name.",
		},
		{
			Name: "test_resource_reference",
			Input: setPropertyDocumentationReferenceInput{
				Reference:       reference(),
				TargetArtifacts: []ArtifactEnum{ResourceArtifact},
				TargetPresent:   true,
			},
			Expected: "The target name. This attribute can be referenced from " + resourceReference + ".",
		},
		{
			Name: "test_datasource_reference",
			Input: setPropertyDocumentationReferenceInput{
				Reference:       reference(),
				TargetArtifacts: []ArtifactEnum{DatasourceArtifact},
				TargetPresent:   true,
			},
			Expected: "The target name. This attribute can be referenced from " + datasourceReference + ".",
		},
		{
			Name: "test_resource_and_datasource_reference",
			Input: setPropertyDocumentationReferenceInput{
				Reference:       reference(),
				TargetArtifacts: []ArtifactEnum{ResourceArtifact, DatasourceArtifact},
				TargetPresent:   true,
			},
			Expected: "The target name. This attribute can be referenced from " + resourceReference + " or from " + datasourceReference + ".",
		},
		{
			Name: "test_definition_only_target_defaults_to_both_artifacts",
			Input: setPropertyDocumentationReferenceInput{
				Reference:        reference(),
				TargetDefinition: &ClassDefinition{ResourceName: "l3_outside"},
			},
			Expected: "The target name. This attribute can be referenced from " + resourceReference + " or from " + datasourceReference + ".",
		},
		{
			Name: "test_definition_only_datasource_target",
			Input: setPropertyDocumentationReferenceInput{
				Reference:        reference(),
				TargetDefinition: &ClassDefinition{ResourceName: "l3_outside", Artifacts: []ArtifactEnum{DatasourceArtifact}},
			},
			Expected: "The target name. This attribute can be referenced from " + datasourceReference + ".",
		},
		{
			Name: "test_definition_only_target_without_artifacts",
			Input: setPropertyDocumentationReferenceInput{
				Reference:        reference(),
				TargetDefinition: &ClassDefinition{ResourceName: "l3_outside", Artifacts: []ArtifactEnum{}},
			},
			Expected: "The target name.",
		},
		{
			Name: "test_existing_reference_description_is_not_duplicated",
			Input: setPropertyDocumentationReferenceInput{
				Reference:       reference(),
				TargetArtifacts: []ArtifactEnum{ResourceArtifact, DatasourceArtifact},
				TargetPresent:   true,
				Description:     "This attribute can be referenced from multiple target resources.",
			},
			Expected: "This attribute can be referenced from multiple target resources.",
		},
		{
			Name: "test_existing_registry_link_is_not_duplicated",
			Input: setPropertyDocumentationReferenceInput{
				Reference:       reference(),
				TargetArtifacts: []ArtifactEnum{ResourceArtifact, DatasourceArtifact},
				TargetPresent:   true,
				Description:     "The target can be configured using [aci_l3_outside](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3_outside).",
			},
			Expected: "The target can be configured using [aci_l3_outside](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3_outside).",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPropertyDocumentationReferenceInput)
			expected := testCase.Expected.(string)
			description := input.Description
			if description == "" {
				description = "The target name."
			}
			property := &Property{
				PropertyName:  "tnL3extOutName",
				Reference:     input.Reference,
				Documentation: PropertyDocumentation{Description: description},
			}
			class := Class{
				PropertiesAll: []string{"tnL3extOutName"},
				Properties:    map[string]*Property{"tnL3extOutName": property},
			}
			classes := map[string]Class{}
			if input.TargetPresent {
				classes["l3extOut"] = Class{
					ResourceName: "l3_outside",
					Artifacts:    input.TargetArtifacts,
				}
			}

			definitions := map[string]ClassDefinition{}
			if input.TargetDefinition != nil {
				definitions["l3extOut"] = *input.TargetDefinition
			}

			class.setPropertyDocumentationReferences(&DataStore{Classes: classes, ClassDefinitions: definitions})

			assert.Equal(t, expected, property.Documentation.Description, test.MessageEqual(expected, property.Documentation.Description, testCase.Name))
		})
	}
}

type setPropertyDocumentationDescriptionsInput struct {
	GlobalOverrides map[string]string
	Properties      map[string]*Property
	Label           string
}

type setPropertyDocumentationDescriptionsExpected struct {
	Descriptions map[string]string
}

func TestSetPropertyDocumentationDescriptions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	newProperty := func(name, perClassDescription, initialDescription string) *Property {
		return &Property{
			PropertyName:  name,
			Documentation: PropertyDocumentation{Description: initialDescription},
			propertyDefinition: PropertyDefinition{
				Documentation: ArtifactDocumentationDefinition{Description: perClassDescription},
			},
		}
	}

	testCases := []test.TestCase{
		{
			// No global overrides: every existing Description is left intact.
			Name: "test_no_overrides_noop",
			Input: setPropertyDocumentationDescriptionsInput{
				Label: "Application EPG",
				Properties: map[string]*Property{
					"name": newProperty("name", "", "Meta-derived description."),
				},
			},
			Expected: setPropertyDocumentationDescriptionsExpected{
				Descriptions: map[string]string{"name": "Meta-derived description."},
			},
		},
		{
			// Override with %s gets interpolated against the class Label.
			Name: "test_template_interpolates_label",
			Input: setPropertyDocumentationDescriptionsInput{
				Label: "Application EPG",
				GlobalOverrides: map[string]string{
					"annotation": "The annotation of the %s object.",
				},
				Properties: map[string]*Property{
					"annotation": newProperty("annotation", "", "Meta comment for annotation."),
				},
			},
			Expected: setPropertyDocumentationDescriptionsExpected{
				Descriptions: map[string]string{"annotation": "The annotation of the Application EPG object."},
			},
		},
		{
			// Override without %s is used verbatim.
			Name: "test_template_no_placeholder_verbatim",
			Input: setPropertyDocumentationDescriptionsInput{
				Label: "Application EPG",
				GlobalOverrides: map[string]string{
					"nameAlias": "A name alias for this object.",
				},
				Properties: map[string]*Property{
					"nameAlias": newProperty("nameAlias", "", "Meta-derived alias text."),
				},
			},
			Expected: setPropertyDocumentationDescriptionsExpected{
				Descriptions: map[string]string{"nameAlias": "A name alias for this object."},
			},
		},
		{
			// Per-class override wins — the global lookup must NOT replace it.
			Name: "test_per_class_override_wins",
			Input: setPropertyDocumentationDescriptionsInput{
				Label: "Application EPG",
				GlobalOverrides: map[string]string{
					"annotation": "Global override for %s.",
				},
				Properties: map[string]*Property{
					"annotation": newProperty("annotation", "Per-class override wins.", "Per-class override wins."),
				},
			},
			Expected: setPropertyDocumentationDescriptionsExpected{
				Descriptions: map[string]string{"annotation": "Per-class override wins."},
			},
		},
		{
			// Per-class templates retain precedence and still resolve their class label.
			Name: "test_per_class_template_interpolates_label",
			Input: setPropertyDocumentationDescriptionsInput{
				Label: "Application EPG",
				GlobalOverrides: map[string]string{
					"contractExceptionTag": "Global override for %s.",
				},
				Properties: map[string]*Property{
					"contractExceptionTag": newProperty("contractExceptionTag", "The contract exception tag of the %s object.", "The contract exception tag of the %s object."),
				},
			},
			Expected: setPropertyDocumentationDescriptionsExpected{
				Descriptions: map[string]string{"contractExceptionTag": "The contract exception tag of the Application EPG object."},
			},
		},
		{
			// Mixed: one property has a global match, one does not.
			Name: "test_mixed_match_and_no_match",
			Input: setPropertyDocumentationDescriptionsInput{
				Label: "Application EPG",
				GlobalOverrides: map[string]string{
					"annotation": "The annotation of the %s object.",
				},
				Properties: map[string]*Property{
					"annotation": newProperty("annotation", "", "Old annotation."),
					"name":       newProperty("name", "", "Meta name comment."),
				},
			},
			Expected: setPropertyDocumentationDescriptionsExpected{
				Descriptions: map[string]string{
					"annotation": "The annotation of the Application EPG object.",
					"name":       "Meta name comment.",
				},
			},
		},
		{
			// Descriptions use the same terminal-period normalization as the legacy generator.
			Name: "test_missing_terminal_period",
			Input: setPropertyDocumentationDescriptionsInput{
				Label: "Application EPG",
				Properties: map[string]*Property{
					"name": newProperty("name", "", "Name"),
				},
			},
			Expected: setPropertyDocumentationDescriptionsExpected{
				Descriptions: map[string]string{"name": "Name."},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPropertyDocumentationDescriptionsInput)
			expected := testCase.Expected.(setPropertyDocumentationDescriptionsExpected)

			class := Class{
				Name:       testClassName("fvAEPg"),
				Properties: input.Properties,
			}
			class.Documentation.Label = input.Label

			ds := &DataStore{
				GlobalMetaDefinition: GlobalMetaDefinition{
					PropertyDocumentationOverrides: input.GlobalOverrides,
				},
			}

			class.setPropertyDocumentationDescriptions(ds)

			for name, want := range expected.Descriptions {
				got := class.Properties[name].Documentation.Description
				assert.Equal(t, want, got, test.MessageEqual(want, got, testCase.Name+"/"+name))
			}
		})
	}
}

type setPropertyDocumentationSupportedVersionsInput struct {
	ClassVersions    string
	PropertyVersions string
}

func TestSetPropertyDocumentationSupportedVersions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_property_covers_class",
			Input: setPropertyDocumentationSupportedVersionsInput{
				ClassVersions:    "1.1(1j)-",
				PropertyVersions: "1.0(1e)-",
			},
			Expected: "",
		},
		{
			Name: "test_property_starts_later",
			Input: setPropertyDocumentationSupportedVersionsInput{
				ClassVersions:    "1.0(1e)-",
				PropertyVersions: "3.2(1l)-",
			},
			Expected: "3.2(1l) and later.",
		},
		{
			Name: "test_property_ends_earlier",
			Input: setPropertyDocumentationSupportedVersionsInput{
				ClassVersions:    "1.0(1e)-",
				PropertyVersions: "1.0(1e)-6.1(1e)",
			},
			Expected: "1.0(1e) to 6.1(1e)",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPropertyDocumentationSupportedVersionsInput)
			expected := testCase.Expected.(string)

			classVersions, err := NewVersions(input.ClassVersions)
			assert.NoError(t, err)
			propertyVersions, err := NewVersions(input.PropertyVersions)
			assert.NoError(t, err)
			property := &Property{SupportedVersions: propertyVersions}
			class := Class{
				SupportedVersions: classVersions,
				Properties:        map[string]*Property{"property": property},
			}

			class.setPropertyDocumentationSupportedVersions()

			assert.Equal(t, expected, property.Documentation.SupportedVersions, test.MessageEqual(expected, property.Documentation.SupportedVersions, testCase.Name))
		})
	}
}

func TestSetPropertyNotes(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_notes_set",
			Input: []string{
				"Note one.",
				"Note two.",
			},
			Expected: []string{
				"Note one.",
				"Note two.",
			},
		},
		{
			Name:     "test_nil_notes",
			Input:    []string(nil),
			Expected: []string(nil),
		},
		{
			Name:     "test_empty_notes",
			Input:    []string{},
			Expected: []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.([]string)
			expected := testCase.Expected.([]string)

			p := &Property{
				PropertyName: "testProp",
				propertyDefinition: PropertyDefinition{
					Documentation: ArtifactDocumentationDefinition{
						Notes: input,
					},
				},
			}

			p.Documentation.setNotes(p)

			assert.Equal(t, expected, p.Documentation.Notes, test.MessageEqual(expected, p.Documentation.Notes, testCase.Name))
		})
	}
}

func TestSetPropertyWarnings(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_warnings_set",
			Input: []string{
				"Warning one.",
				"Warning two.",
			},
			Expected: []string{
				"Warning one.",
				"Warning two.",
			},
		},
		{
			Name:     "test_nil_warnings",
			Input:    []string(nil),
			Expected: []string(nil),
		},
		{
			Name:     "test_empty_warnings",
			Input:    []string{},
			Expected: []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.([]string)
			expected := testCase.Expected.([]string)

			p := &Property{
				PropertyName: "testProp",
				propertyDefinition: PropertyDefinition{
					Documentation: ArtifactDocumentationDefinition{
						Warnings: input,
					},
				},
			}

			p.Documentation.setWarnings(p)

			assert.Equal(t, expected, p.Documentation.Warnings, test.MessageEqual(expected, p.Documentation.Warnings, testCase.Name))
		})
	}
}

type setPropertyDefaultValuesInput struct {
	DefinitionDefaultValues map[string]string
	MetaDefault             any
	ValidValues             ValidValues
	ValueType               ValueTypeEnum
}

type setPropertyDefaultValuesExpected struct {
	DefaultValues []DefaultValue
	Error         bool
	ErrorMsg      string
}

func TestSetPropertyDefaultValues(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_definition_override_no_versions",
			Input: setPropertyDefaultValuesInput{
				DefinitionDefaultValues: map[string]string{"enabled": "", "disabled": ""},
				MetaDefault:             "1",
				ValidValues:             ValidValues{"1": {LocalName: "enabled"}},
			},
			Expected: setPropertyDefaultValuesExpected{
				DefaultValues: []DefaultValue{{Value: "disabled"}, {Value: "enabled"}},
			},
		},
		{
			Name: "test_definition_override_with_versions",
			Input: setPropertyDefaultValuesInput{
				DefinitionDefaultValues: map[string]string{"enabled": "5.0(1a)-"},
			},
			Expected: setPropertyDefaultValuesExpected{
				DefaultValues: []DefaultValue{{Value: "enabled", Versions: testVersions(t, "5.0(1a)-")}},
			},
		},
		{
			Name: "test_definition_override_multi_version",
			Input: setPropertyDefaultValuesInput{
				DefinitionDefaultValues: map[string]string{
					"legacy":  "-4.2(7w)",
					"current": "5.0(1a)-6.0(2h)",
					"future":  "6.1(1a)-",
				},
			},
			Expected: setPropertyDefaultValuesExpected{
				DefaultValues: []DefaultValue{
					{Value: "current", Versions: testVersions(t, "5.0(1a)-6.0(2h)")},
					{Value: "future", Versions: testVersions(t, "6.1(1a)-")},
					{Value: "legacy", Versions: testVersions(t, "-4.2(7w)")},
				},
			},
		},
		{
			Name: "test_definition_override_mixed_versioned_and_unversioned",
			Input: setPropertyDefaultValuesInput{
				DefinitionDefaultValues: map[string]string{
					"baseline": "",
					"v5_plus":  "5.0(1a)-",
				},
			},
			Expected: setPropertyDefaultValuesExpected{
				DefaultValues: []DefaultValue{
					{Value: "baseline"},
					{Value: "v5_plus", Versions: testVersions(t, "5.0(1a)-")},
				},
			},
		},
		{
			Name: "test_definition_override_invalid_versions",
			Input: setPropertyDefaultValuesInput{
				DefinitionDefaultValues: map[string]string{"enabled": "invalid"},
			},
			Expected: setPropertyDefaultValuesExpected{
				Error:    true,
				ErrorMsg: "failed to parse default value versions for property 'testProp', value 'enabled'",
			},
		},
		{
			Name: "test_meta_string_default_with_valid_values_lookup",
			Input: setPropertyDefaultValuesInput{
				MetaDefault: "1",
				ValidValues: ValidValues{"1": {LocalName: "enabled"}, "0": {LocalName: "disabled"}},
			},
			Expected: setPropertyDefaultValuesExpected{
				DefaultValues: []DefaultValue{{Value: "enabled"}},
			},
		},
		{
			Name: "test_meta_string_default_no_valid_values",
			Input: setPropertyDefaultValuesInput{
				MetaDefault: "default_value",
			},
			Expected: setPropertyDefaultValuesExpected{
				DefaultValues: []DefaultValue{{Value: "default_value"}},
			},
		},
		{
			Name: "test_meta_float64_default",
			Input: setPropertyDefaultValuesInput{
				MetaDefault: float64(10),
			},
			Expected: setPropertyDefaultValuesExpected{
				DefaultValues: []DefaultValue{{Value: "10"}},
			},
		},
		{
			Name: "test_meta_float64_default_decimal",
			Input: setPropertyDefaultValuesInput{
				MetaDefault: float64(3.5),
			},
			Expected: setPropertyDefaultValuesExpected{
				DefaultValues: []DefaultValue{{Value: "3.5"}},
			},
		},
		{
			Name: "test_set_type_none_not_in_valid_values",
			Input: setPropertyDefaultValuesInput{
				MetaDefault: "none",
				ValidValues: ValidValues{"read": {LocalName: "read"}, "write": {LocalName: "write"}},
				ValueType:   Set,
			},
			Expected: setPropertyDefaultValuesExpected{
				DefaultValues: []DefaultValue{{Value: ""}},
			},
		},
		{
			Name: "test_set_type_none_in_valid_values",
			Input: setPropertyDefaultValuesInput{
				MetaDefault: "none",
				ValidValues: ValidValues{"none": {LocalName: "none"}, "read": {LocalName: "read"}},
				ValueType:   Set,
			},
			Expected: setPropertyDefaultValuesExpected{
				DefaultValues: []DefaultValue{{Value: "none"}},
			},
		},
		{
			Name: "test_no_meta_default",
			Input: setPropertyDefaultValuesInput{
				MetaDefault: nil,
			},
			Expected: setPropertyDefaultValuesExpected{
				DefaultValues: []DefaultValue(nil),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPropertyDefaultValuesInput)
			expected := testCase.Expected.(setPropertyDefaultValuesExpected)

			metaDetails := map[string]any{}
			if input.MetaDefault != nil {
				metaDetails["default"] = input.MetaDefault
			}

			p := &Property{
				PropertyName: "testProp",
				propertyDefinition: PropertyDefinition{
					DefaultValues: input.DefinitionDefaultValues,
				},
				metaDetails: metaDetails,
				ValidValues: input.ValidValues,
				ValueType:   input.ValueType,
			}

			err := p.Documentation.setDefaultValues(p)

			if expected.Error {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), expected.ErrorMsg)
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, expected.DefaultValues, p.Documentation.DefaultValues, test.MessageEqual(expected.DefaultValues, p.Documentation.DefaultValues, testCase.Name))
			}
		})
	}
}

// testVersions is a helper that parses a version string or fails the test.
func testVersions(t *testing.T, versionStr string) *Versions {
	t.Helper()
	v, err := NewVersions(versionStr)
	if err != nil {
		t.Fatalf("testVersions: failed to parse %q: %v", versionStr, err)
	}
	return v
}
