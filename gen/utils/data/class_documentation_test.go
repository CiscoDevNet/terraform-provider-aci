package data

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestSetClassName(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_basic_class_name",
			Input:    "fvTenant",
			Expected: "[fvTenant](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvTenant/overview)",
		},
		{
			Name:     "test_long_class_name",
			Input:    "l3extRsLblToInstP",
			Expected: "[l3extRsLblToInstP](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/l3extRsLblToInstP/overview)",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(string)
			expected := testCase.Expected.(string)

			class := Class{Name: testClassName(input)}

			class.Documentation.setClassName(&class)

			assert.Equal(t, expected, class.Documentation.ClassName, test.MessageEqual(expected, class.Documentation.ClassName, testCase.Name))
		})
	}
}

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

type setDocumentationChildrenInput struct {
	RnMap                      map[string]interface{}
	ChildrenIncludedInResource []string
	StoreClasses               map[string]Class
}

type setDocumentationChildrenExpected struct {
	Children []string
}

func TestSetDocumentationChildren(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_no_rnmap",
			Input:    setDocumentationChildrenInput{},
			Expected: setDocumentationChildrenExpected{Children: nil},
		},
		{
			Name: "test_child_included_in_resource_excluded",
			Input: setDocumentationChildrenInput{
				RnMap: map[string]interface{}{
					"rstoEpg": "fv:RsToEpg",
				},
				ChildrenIncludedInResource: []string{"fvRsToEpg"},
				StoreClasses: map[string]Class{
					"fvRsToEpg": {ResourceName: "relation_to_epg"},
				},
			},
			Expected: setDocumentationChildrenExpected{Children: []string{}},
		},
		{
			Name: "test_unknown_child_excluded",
			Input: setDocumentationChildrenInput{
				RnMap: map[string]interface{}{
					"unknown-{name}": "foo:Bar",
				},
				StoreClasses: map[string]Class{},
			},
			Expected: setDocumentationChildrenExpected{Children: []string{}},
		},
		{
			Name: "test_child_without_resource_name_excluded",
			Input: setDocumentationChildrenInput{
				RnMap: map[string]interface{}{
					"noresource-{name}": "foo:NoResource",
				},
				StoreClasses: map[string]Class{
					"fooNoResource": {ResourceName: ""},
				},
			},
			Expected: setDocumentationChildrenExpected{Children: []string{}},
		},
		{
			Name: "test_two_valid_children_sorted",
			Input: setDocumentationChildrenInput{
				RnMap: map[string]interface{}{
					"zeta-{name}":  "fv:Zeta",
					"alpha-{name}": "fv:Alpha",
				},
				StoreClasses: map[string]Class{
					"fvZeta":  {ResourceName: "zeta"},
					"fvAlpha": {ResourceName: "alpha"},
				},
			},
			Expected: setDocumentationChildrenExpected{Children: []string{
				"[aci_alpha](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/alpha)",
				"[aci_zeta](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/zeta)",
			}},
		},
		{
			Name: "test_mixed_included_in_resource_and_valid",
			Input: setDocumentationChildrenInput{
				RnMap: map[string]interface{}{
					"included-{name}": "fv:Included",
					"valid-{name}":    "fv:Valid",
					"missing-{name}":  "fv:Missing",
				},
				ChildrenIncludedInResource: []string{"fvIncluded"},
				StoreClasses: map[string]Class{
					"fvIncluded": {ResourceName: "included"},
					"fvValid":    {ResourceName: "valid"},
				},
			},
			Expected: setDocumentationChildrenExpected{Children: []string{
				"[aci_valid](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/valid)",
			}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setDocumentationChildrenInput)
			expected := testCase.Expected.(setDocumentationChildrenExpected)

			childrenIncludedInResource := make([]*ClassName, 0, len(input.ChildrenIncludedInResource))
			for _, name := range input.ChildrenIncludedInResource {
				childrenIncludedInResource = append(childrenIncludedInResource, testClassName(name))
			}

			class := Class{
				Name:     testClassName("fvTenant"),
				Children: childrenIncludedInResource,
			}
			if input.RnMap != nil {
				class.MetaFileContent = map[string]interface{}{"rnMap": input.RnMap}
			}

			ds := &DataStore{Classes: input.StoreClasses}

			class.Documentation.setChildren(&class, ds)

			assert.Equal(t, expected.Children, class.Documentation.Children, test.MessageEqual(expected.Children, class.Documentation.Children, testCase.Name))
		})
	}
}

func TestSetDeprecationWarning(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_not_deprecated",
			Input:    false,
			Expected: "",
		},
		{
			Name:     "test_deprecated",
			Input:    true,
			Expected: "The [fvTenant](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvTenant/overview) class is deprecated and will be removed in a future release.",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			class := Class{
				Name:       testClassName("fvTenant"),
				Deprecated: testCase.Input.(bool),
			}
			class.Documentation.setClassName(&class)
			class.Documentation.setDeprecationWarning(&class)

			expected := testCase.Expected.(string)
			assert.Equal(t, expected, class.Documentation.DeprecationWarning, test.MessageEqual(expected, class.Documentation.DeprecationWarning, testCase.Name))
		})
	}
}

type setLabelInput struct {
	ResourceName string
	Label        string
	Overrides    map[string]string
}

func TestSetLabel(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_label_override_used",
			Input: setLabelInput{
				ResourceName: "application_epg",
				Label:        "Custom Label",
				Overrides:    map[string]string{"Application": "App"},
			},
			Expected: "Custom Label",
		},
		{
			Name: "test_humanize_no_overrides",
			Input: setLabelInput{
				ResourceName: "application_epg",
			},
			Expected: "Application Epg",
		},
		{
			Name: "test_humanize_single_word_override",
			Input: setLabelInput{
				ResourceName: "bgp_timers",
				Overrides:    map[string]string{"Bgp": "BGP"},
			},
			Expected: "BGP Timers",
		},
		{
			Name: "test_humanize_multi_word_override",
			Input: setLabelInput{
				ResourceName: "external_network_instance_profile",
				Overrides:    map[string]string{"External Network Instance Profile": "External EPG"},
			},
			Expected: "External EPG",
		},
		{
			Name: "test_single_word_override_no_partial_match",
			Input: setLabelInput{
				ResourceName: "applications_epg",
				Overrides:    map[string]string{"Application": "App"},
			},
			Expected: "Applications Epg",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setLabelInput)
			expected := testCase.Expected.(string)

			class := Class{
				Name:         testClassName("fvFoo"),
				ResourceName: input.ResourceName,
				ClassDefinition: ClassDefinition{
					Documentation: ClassDocumentationDefinition{Label: input.Label},
				},
			}
			ds := &DataStore{
				GlobalMetaDefinition: GlobalMetaDefinition{
					DocumentationLabelOverrides: input.Overrides,
				},
			}

			class.Documentation.setLabel(&class, ds)

			assert.Equal(t, expected, class.Documentation.Label, test.MessageEqual(expected, class.Documentation.Label, testCase.Name))
		})
	}
}

type setDescriptionInput struct {
	Label               string
	SharedDescription   string
	ResourceDescription string
	DataDescription     string
}

type setDescriptionExpected struct {
	ResourceDescription   string
	DatasourceDescription string
}

func TestSetDescription(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_no_appendix",
			Input: setDescriptionInput{
				Label: "Application EPG",
			},
			Expected: setDescriptionExpected{
				ResourceDescription:   "Manages ACI Application EPG.",
				DatasourceDescription: "Data source for ACI Application EPG.",
			},
		},
		{
			Name: "test_shared_only",
			Input: setDescriptionInput{
				Label:             "Application EPG",
				SharedDescription: "Provides custom routing.",
			},
			Expected: setDescriptionExpected{
				ResourceDescription:   "Manages ACI Application EPG. Provides custom routing.",
				DatasourceDescription: "Data source for ACI Application EPG. Provides custom routing.",
			},
		},
		{
			Name: "test_resource_only",
			Input: setDescriptionInput{
				Label:               "Application EPG",
				ResourceDescription: "Use with caution.",
			},
			Expected: setDescriptionExpected{
				ResourceDescription:   "Manages ACI Application EPG. Use with caution.",
				DatasourceDescription: "Data source for ACI Application EPG.",
			},
		},
		{
			Name: "test_datasource_only",
			Input: setDescriptionInput{
				Label:           "Application EPG",
				DataDescription: "Returns matching objects.",
			},
			Expected: setDescriptionExpected{
				ResourceDescription:   "Manages ACI Application EPG.",
				DatasourceDescription: "Data source for ACI Application EPG. Returns matching objects.",
			},
		},
		{
			Name: "test_shared_plus_resource",
			Input: setDescriptionInput{
				Label:               "Application EPG",
				SharedDescription:   "Provides custom routing.",
				ResourceDescription: "Use with caution.",
			},
			Expected: setDescriptionExpected{
				ResourceDescription:   "Manages ACI Application EPG. Provides custom routing. Use with caution.",
				DatasourceDescription: "Data source for ACI Application EPG. Provides custom routing.",
			},
		},
		{
			Name: "test_shared_plus_datasource",
			Input: setDescriptionInput{
				Label:             "Application EPG",
				SharedDescription: "Provides custom routing.",
				DataDescription:   "Returns matching objects.",
			},
			Expected: setDescriptionExpected{
				ResourceDescription:   "Manages ACI Application EPG. Provides custom routing.",
				DatasourceDescription: "Data source for ACI Application EPG. Provides custom routing. Returns matching objects.",
			},
		},
		{
			Name: "test_all_three",
			Input: setDescriptionInput{
				Label:               "Application EPG",
				SharedDescription:   "Provides custom routing.",
				ResourceDescription: "Use with caution.",
				DataDescription:     "Returns matching objects.",
			},
			Expected: setDescriptionExpected{
				ResourceDescription:   "Manages ACI Application EPG. Provides custom routing. Use with caution.",
				DatasourceDescription: "Data source for ACI Application EPG. Provides custom routing. Returns matching objects.",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setDescriptionInput)
			expected := testCase.Expected.(setDescriptionExpected)

			class := Class{
				Name: testClassName("fvFoo"),
				ClassDefinition: ClassDefinition{
					Documentation: ClassDocumentationDefinition{
						Description: input.SharedDescription,
						Resource:    ArtifactDocumentationDefinition{Description: input.ResourceDescription},
						Datasource:  ArtifactDocumentationDefinition{Description: input.DataDescription},
					},
				},
			}
			class.Documentation.Label = input.Label

			class.Documentation.setDescription(&class)

			assert.Equal(t, expected.ResourceDescription, class.Documentation.ResourceDescription, test.MessageEqual(expected.ResourceDescription, class.Documentation.ResourceDescription, testCase.Name+"/resource"))
			assert.Equal(t, expected.DatasourceDescription, class.Documentation.DatasourceDescription, test.MessageEqual(expected.DatasourceDescription, class.Documentation.DatasourceDescription, testCase.Name+"/datasource"))
		})
	}
}

type setDnFormatsInput struct {
	MetaFormats     []interface{}
	OverrideFormats []string
}

func TestSetDnFormats(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	tooManyNotice := "Too many DN formats to display, see model documentation for all possible parents of [fvTenant](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvTenant/overview)."

	testCases := []test.TestCase{
		{
			Name:     "test_no_meta_no_override",
			Input:    setDnFormatsInput{},
			Expected: []string(nil),
		},
		{
			Name: "test_meta_single",
			Input: setDnFormatsInput{
				MetaFormats: []interface{}{"uni/tn-{name}/ap-{name}/epg-{name}"},
			},
			Expected: []string{"uni/tn-{name}/ap-{name}/epg-{name}"},
		},
		{
			Name: "test_meta_multiple_sorted",
			Input: setDnFormatsInput{
				MetaFormats: []interface{}{
					"uni/tn-{name}/certstore/keyring-{name}",
					"uni/userext/pkiext/keyring-{name}",
				},
			},
			Expected: []string{
				"uni/tn-{name}/certstore/keyring-{name}",
				"uni/userext/pkiext/keyring-{name}",
			},
		},
		{
			Name: "test_meta_unsorted_input_is_sorted",
			Input: setDnFormatsInput{
				MetaFormats: []interface{}{"b/{name}", "a/{name}", "c/{name}"},
			},
			Expected: []string{"a/{name}", "b/{name}", "c/{name}"},
		},
		{
			Name: "test_override_replaces_meta",
			Input: setDnFormatsInput{
				MetaFormats:     []interface{}{"meta/{name}"},
				OverrideFormats: []string{"override/{name}"},
			},
			Expected: []string{"override/{name}"},
		},
		{
			Name: "test_override_is_sorted",
			Input: setDnFormatsInput{
				OverrideFormats: []string{"z/{name}", "a/{name}"},
			},
			Expected: []string{"a/{name}", "z/{name}"},
		},
		{
			Name: "test_at_cap_no_notice",
			Input: setDnFormatsInput{
				MetaFormats: []interface{}{"e/{name}", "d/{name}", "c/{name}", "b/{name}", "a/{name}"},
			},
			Expected: []string{"a/{name}", "b/{name}", "c/{name}", "d/{name}", "e/{name}"},
		},
		{
			Name: "test_over_cap_prepends_notice_and_truncates",
			Input: setDnFormatsInput{
				MetaFormats: []interface{}{"f/{name}", "e/{name}", "d/{name}", "c/{name}", "b/{name}", "a/{name}"},
			},
			Expected: []string{
				tooManyNotice,
				"a/{name}", "b/{name}", "c/{name}", "d/{name}", "e/{name}",
			},
		},
		{
			Name: "test_meta_non_string_entries_skipped",
			Input: setDnFormatsInput{
				MetaFormats: []interface{}{"a/{name}", 42, "b/{name}"},
			},
			Expected: []string{"a/{name}", "b/{name}"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setDnFormatsInput)
			expected := testCase.Expected.([]string)

			class := Class{
				Name: testClassName("fvTenant"),
				ClassDefinition: ClassDefinition{
					Documentation: ClassDocumentationDefinition{DnFormats: input.OverrideFormats},
				},
			}
			if input.MetaFormats != nil {
				class.MetaFileContent = map[string]interface{}{"dnFormats": input.MetaFormats}
			}

			class.Documentation.setClassName(&class)
			class.Documentation.setDnFormats(&class)

			assert.Equal(t, expected, class.Documentation.DnFormats, test.MessageEqual(expected, class.Documentation.DnFormats, testCase.Name))
		})
	}
}
