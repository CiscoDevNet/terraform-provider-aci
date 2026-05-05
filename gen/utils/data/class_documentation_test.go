package data

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

type setSubCategoryInput struct {
	SubCategory                      string
	IsSingleNestedWhenDefinedAsChild bool
}

type setSubCategoryExpected struct {
	SubCategory string
	Error       bool
	ErrorMsg    string
}

func TestSetSubCategory(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_valid_sub_category_networking",
			Input: setSubCategoryInput{
				SubCategory: "Networking",
			},
			Expected: setSubCategoryExpected{
				SubCategory: "Networking",
			},
		},
		{
			Name: "test_valid_sub_category_aaa",
			Input: setSubCategoryInput{
				SubCategory: "AAA",
			},
			Expected: setSubCategoryExpected{
				SubCategory: "AAA",
			},
		},
		{
			Name: "test_valid_sub_category_application_mgmt",
			Input: setSubCategoryInput{
				SubCategory: "Application Management",
			},
			Expected: setSubCategoryExpected{
				SubCategory: "Application Management",
			},
		},
		{
			Name: "test_valid_sub_category_l4l7_services",
			Input: setSubCategoryInput{
				SubCategory: "L4-L7",
			},
			Expected: setSubCategoryExpected{
				SubCategory: "L4-L7",
			},
		},
		{
			Name: "test_empty_sub_category_not_single_nested",
			Input: setSubCategoryInput{
				SubCategory:                      "",
				IsSingleNestedWhenDefinedAsChild: false,
			},
			Expected: setSubCategoryExpected{
				Error:    true,
				ErrorMsg: "class 'fvTenant': sub_category not specified for class 'fvTenant': add documentation.sub_category to the class definition file",
			},
		},
		{
			Name: "test_empty_sub_category_single_nested",
			Input: setSubCategoryInput{
				SubCategory:                      "",
				IsSingleNestedWhenDefinedAsChild: true,
			},
			Expected: setSubCategoryExpected{
				SubCategory: "",
			},
		},
		{
			Name: "test_invalid_sub_category",
			Input: setSubCategoryInput{
				SubCategory: "Invalid Category",
			},
			Expected: setSubCategoryExpected{
				Error:    true,
				ErrorMsg: "class 'fvTenant': invalid sub_category 'Invalid Category'",
			},
		},
		{
			Name: "test_invalid_sub_category_typo",
			Input: setSubCategoryInput{
				SubCategory: "networking",
			},
			Expected: setSubCategoryExpected{
				Error:    true,
				ErrorMsg: "class 'fvTenant': invalid sub_category 'networking'",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setSubCategoryInput)
			expected := testCase.Expected.(setSubCategoryExpected)

			class := Class{
				Name:                             testClassName("fvTenant"),
				IsSingleNestedWhenDefinedAsChild: input.IsSingleNestedWhenDefinedAsChild,
				ClassDefinition: ClassDefinition{
					Documentation: ClassDocumentationDefinition{
						SubCategory: input.SubCategory,
					},
				},
			}

			err := class.Documentation.setSubCategory(&class)

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, expected.SubCategory, class.Documentation.SubCategory, test.MessageEqual(expected.SubCategory, class.Documentation.SubCategory, testCase.Name))
			}
		})
	}
}

type setUiLocationsInput struct {
	UiLocations                      []string
	IsSingleNestedWhenDefinedAsChild bool
}

type setUiLocationsExpected struct {
	UiLocations []string
	Error       bool
	ErrorMsg    string
}

func TestSetUiLocations(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_valid_single_ui_location",
			Input: setUiLocationsInput{
				UiLocations: []string{"Tenants -> Networking -> VRFs"},
			},
			Expected: setUiLocationsExpected{
				UiLocations: []string{"Tenants -> Networking -> VRFs"},
			},
		},
		{
			Name: "test_valid_multiple_ui_locations",
			Input: setUiLocationsInput{
				UiLocations: []string{
					"Tenants -> Application Profiles -> Application EPGs -> Contracts",
					"Tenants -> Networking -> L3Outs -> External EPGs -> Contracts",
				},
			},
			Expected: setUiLocationsExpected{
				UiLocations: []string{
					"Tenants -> Application Profiles -> Application EPGs -> Contracts",
					"Tenants -> Networking -> L3Outs -> External EPGs -> Contracts",
				},
			},
		},
		{
			Name: "test_not_shown_in_ui",
			Input: setUiLocationsInput{
				UiLocations: []string{"Not shown in UI"},
			},
			Expected: setUiLocationsExpected{
				UiLocations: []string{"Not shown in UI"},
			},
		},
		{
			Name: "test_empty_ui_locations_single_nested",
			Input: setUiLocationsInput{
				UiLocations:                      nil,
				IsSingleNestedWhenDefinedAsChild: true,
			},
			Expected: setUiLocationsExpected{
				UiLocations: nil,
			},
		},
		{
			Name: "test_empty_ui_locations_not_single_nested",
			Input: setUiLocationsInput{
				UiLocations:                      nil,
				IsSingleNestedWhenDefinedAsChild: false,
			},
			Expected: setUiLocationsExpected{
				Error:    true,
				ErrorMsg: "class 'fvTenant': ui_locations not specified: add documentation.ui_locations to the class definition file",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setUiLocationsInput)
			expected := testCase.Expected.(setUiLocationsExpected)

			class := Class{
				Name:                             testClassName("fvTenant"),
				IsSingleNestedWhenDefinedAsChild: input.IsSingleNestedWhenDefinedAsChild,
				ClassDefinition: ClassDefinition{
					Documentation: ClassDocumentationDefinition{
						UiLocations: input.UiLocations,
					},
				},
			}

			err := class.Documentation.setUiLocations(&class)

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, expected.UiLocations, class.Documentation.UiLocations, test.MessageEqual(expected.UiLocations, class.Documentation.UiLocations, testCase.Name))
			}
		})
	}
}

type setNotesInput struct {
	Shared     []string
	Resource   []string
	Datasource []string
}

type setNotesExpected struct {
	ResourceNotes   []string
	DatasourceNotes []string
}

func TestSetNotes(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_all_empty",
			Input:    setNotesInput{},
			Expected: setNotesExpected{ResourceNotes: nil, DatasourceNotes: nil},
		},
		{
			Name: "test_only_shared",
			Input: setNotesInput{
				Shared: []string{"shared note 1", "shared note 2"},
			},
			Expected: setNotesExpected{
				ResourceNotes:   []string{"shared note 1", "shared note 2"},
				DatasourceNotes: []string{"shared note 1", "shared note 2"},
			},
		},
		{
			Name: "test_only_resource",
			Input: setNotesInput{
				Resource: []string{"resource note"},
			},
			Expected: setNotesExpected{
				ResourceNotes:   []string{"resource note"},
				DatasourceNotes: nil,
			},
		},
		{
			Name: "test_only_datasource",
			Input: setNotesInput{
				Datasource: []string{"datasource note"},
			},
			Expected: setNotesExpected{
				ResourceNotes:   nil,
				DatasourceNotes: []string{"datasource note"},
			},
		},
		{
			Name: "test_shared_and_resource",
			Input: setNotesInput{
				Shared:   []string{"shared note"},
				Resource: []string{"resource note"},
			},
			Expected: setNotesExpected{
				ResourceNotes:   []string{"shared note", "resource note"},
				DatasourceNotes: []string{"shared note"},
			},
		},
		{
			Name: "test_shared_resource_and_datasource",
			Input: setNotesInput{
				Shared:     []string{"shared note"},
				Resource:   []string{"resource note"},
				Datasource: []string{"datasource note 1", "datasource note 2"},
			},
			Expected: setNotesExpected{
				ResourceNotes:   []string{"shared note", "resource note"},
				DatasourceNotes: []string{"shared note", "datasource note 1", "datasource note 2"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setNotesInput)
			expected := testCase.Expected.(setNotesExpected)

			class := Class{
				Name: testClassName("fvTenant"),
				ClassDefinition: ClassDefinition{
					Documentation: ClassDocumentationDefinition{
						Notes:      input.Shared,
						Resource:   ArtifactDocumentationDefinition{Notes: input.Resource},
						Datasource: ArtifactDocumentationDefinition{Notes: input.Datasource},
					},
				},
			}

			class.Documentation.setNotes(&class)

			assert.Equal(t, expected.ResourceNotes, class.Documentation.ResourceNotes, test.MessageEqual(expected.ResourceNotes, class.Documentation.ResourceNotes, testCase.Name))
			assert.Equal(t, expected.DatasourceNotes, class.Documentation.DatasourceNotes, test.MessageEqual(expected.DatasourceNotes, class.Documentation.DatasourceNotes, testCase.Name))
		})
	}
}

type setWarningsInput struct {
	Shared     []string
	Resource   []string
	Datasource []string
}

type setWarningsExpected struct {
	ResourceWarnings   []string
	DatasourceWarnings []string
}

func TestSetWarnings(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_all_empty",
			Input:    setWarningsInput{},
			Expected: setWarningsExpected{ResourceWarnings: nil, DatasourceWarnings: nil},
		},
		{
			Name: "test_only_shared",
			Input: setWarningsInput{
				Shared: []string{"shared warning"},
			},
			Expected: setWarningsExpected{
				ResourceWarnings:   []string{"shared warning"},
				DatasourceWarnings: []string{"shared warning"},
			},
		},
		{
			Name: "test_only_resource",
			Input: setWarningsInput{
				Resource: []string{"resource warning"},
			},
			Expected: setWarningsExpected{
				ResourceWarnings:   []string{"resource warning"},
				DatasourceWarnings: nil,
			},
		},
		{
			Name: "test_only_datasource",
			Input: setWarningsInput{
				Datasource: []string{"datasource warning"},
			},
			Expected: setWarningsExpected{
				ResourceWarnings:   nil,
				DatasourceWarnings: []string{"datasource warning"},
			},
		},
		{
			Name: "test_shared_and_resource",
			Input: setWarningsInput{
				Shared:   []string{"shared warning"},
				Resource: []string{"resource warning"},
			},
			Expected: setWarningsExpected{
				ResourceWarnings:   []string{"shared warning", "resource warning"},
				DatasourceWarnings: []string{"shared warning"},
			},
		},
		{
			Name: "test_shared_resource_and_datasource",
			Input: setWarningsInput{
				Shared:     []string{"shared warning"},
				Resource:   []string{"resource warning"},
				Datasource: []string{"datasource warning 1", "datasource warning 2"},
			},
			Expected: setWarningsExpected{
				ResourceWarnings:   []string{"shared warning", "resource warning"},
				DatasourceWarnings: []string{"shared warning", "datasource warning 1", "datasource warning 2"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setWarningsInput)
			expected := testCase.Expected.(setWarningsExpected)

			class := Class{
				Name: testClassName("fvTenant"),
				ClassDefinition: ClassDefinition{
					Documentation: ClassDocumentationDefinition{
						Warnings:   input.Shared,
						Resource:   ArtifactDocumentationDefinition{Warnings: input.Resource},
						Datasource: ArtifactDocumentationDefinition{Warnings: input.Datasource},
					},
				},
			}

			class.Documentation.setWarnings(&class)

			assert.Equal(t, expected.ResourceWarnings, class.Documentation.ResourceWarnings, test.MessageEqual(expected.ResourceWarnings, class.Documentation.ResourceWarnings, testCase.Name))
			assert.Equal(t, expected.DatasourceWarnings, class.Documentation.DatasourceWarnings, test.MessageEqual(expected.DatasourceWarnings, class.Documentation.DatasourceWarnings, testCase.Name))
		})
	}
}
