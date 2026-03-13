package data

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

type setSubCategoryInput struct {
	SubCategory                      SubCategoryEnum
	IsSingleNestedWhenDefinedAsChild bool
}

type setSubCategoryExpected struct {
	SubCategory SubCategoryEnum
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
				SubCategory: SubCategoryNetworking,
			},
			Expected: setSubCategoryExpected{
				SubCategory: SubCategoryNetworking,
			},
		},
		{
			Name: "test_valid_sub_category_aaa",
			Input: setSubCategoryInput{
				SubCategory: SubCategoryAAA,
			},
			Expected: setSubCategoryExpected{
				SubCategory: SubCategoryAAA,
			},
		},
		{
			Name: "test_valid_sub_category_application_mgmt",
			Input: setSubCategoryInput{
				SubCategory: SubCategoryApplicationMgmt,
			},
			Expected: setSubCategoryExpected{
				SubCategory: SubCategoryApplicationMgmt,
			},
		},
		{
			Name: "test_valid_sub_category_l4l7_services",
			Input: setSubCategoryInput{
				SubCategory: SubCategoryL4L7Services,
			},
			Expected: setSubCategoryExpected{
				SubCategory: SubCategoryL4L7Services,
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
