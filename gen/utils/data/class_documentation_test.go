package data

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
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
	RnMap                        map[string]any
	ChildrenIncludedInResource   []string
	ExcludeChildren              []string
	DocumentationExcludeChildren []string
	DocumentationIncludeChildren []string
	StoreClasses                 map[string]Class
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
				RnMap: map[string]any{
					"rstoEpg": "fv:RsToEpg",
				},
				ChildrenIncludedInResource: []string{"fvRsToEpg"},
				StoreClasses: map[string]Class{
					"fvRsToEpg": {ResourceName: "relation_to_epg", Artifacts: []ArtifactEnum{ResourceArtifact}},
				},
			},
			Expected: setDocumentationChildrenExpected{Children: []string{}},
		},
		{
			Name: "test_unknown_child_excluded",
			Input: setDocumentationChildrenInput{
				RnMap: map[string]any{
					"unknown-{name}": "foo:Bar",
				},
				StoreClasses: map[string]Class{},
			},
			Expected: setDocumentationChildrenExpected{Children: []string{}},
		},
		{
			Name: "test_child_without_resource_name_excluded",
			Input: setDocumentationChildrenInput{
				RnMap: map[string]any{
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
				RnMap: map[string]any{
					"zeta-{name}":  "fv:Zeta",
					"alpha-{name}": "fv:Alpha",
				},
				StoreClasses: map[string]Class{
					"fvZeta":  {ResourceName: "zeta", Artifacts: []ArtifactEnum{ResourceArtifact}},
					"fvAlpha": {ResourceName: "alpha", Artifacts: []ArtifactEnum{ResourceArtifact}},
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
				RnMap: map[string]any{
					"included-{name}": "fv:Included",
					"valid-{name}":    "fv:Valid",
					"missing-{name}":  "fv:Missing",
				},
				ChildrenIncludedInResource: []string{"fvIncluded"},
				StoreClasses: map[string]Class{
					"fvIncluded": {ResourceName: "included", Artifacts: []ArtifactEnum{ResourceArtifact}},
					"fvValid":    {ResourceName: "valid", Artifacts: []ArtifactEnum{ResourceArtifact}},
				},
			},
			Expected: setDocumentationChildrenExpected{Children: []string{
				"[aci_valid](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/valid)",
			}},
		},
		{
			// fvRsPathAtt-shape case: ExcludeChildren on the class definition must
			// also drop the entry from the docs Children link list (legacy key
			// remove_from_contains is unified into ExcludeChildren by the migration).
			Name: "test_exclude_children_dropped_from_docs",
			Input: setDocumentationChildrenInput{
				RnMap: map[string]any{
					"port-sec-{name}": "l2:PortSecurityPol",
					"keep-{name}":     "fv:Keep",
				},
				ExcludeChildren: []string{"l2PortSecurityPol"},
				StoreClasses: map[string]Class{
					"l2PortSecurityPol": {ResourceName: "port_security_policy", Artifacts: []ArtifactEnum{ResourceArtifact}},
					"fvKeep":            {ResourceName: "keep", Artifacts: []ArtifactEnum{ResourceArtifact}},
				},
			},
			Expected: setDocumentationChildrenExpected{Children: []string{
				"[aci_keep](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/keep)",
			}},
		},
		{
			Name: "test_documentation_exclude_children_dropped_from_docs",
			Input: setDocumentationChildrenInput{
				RnMap: map[string]any{
					"exclude-{name}": "l2:Exclude",
					"keep-{name}":    "fv:Keep",
				},
				DocumentationExcludeChildren: []string{"l2Exclude"},
				StoreClasses: map[string]Class{
					"l2Exclude": {ResourceName: "exclude", Artifacts: []ArtifactEnum{ResourceArtifact}},
					"fvKeep":    {ResourceName: "keep", Artifacts: []ArtifactEnum{ResourceArtifact}},
				},
			},
			Expected: setDocumentationChildrenExpected{Children: []string{
				"[aci_keep](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/keep)",
			}},
		},
		{
			Name: "test_documentation_include_children_without_rnmap",
			Input: setDocumentationChildrenInput{
				DocumentationIncludeChildren: []string{"l2:Included"},
				StoreClasses: map[string]Class{
					"l2Included": {ResourceName: "included", Artifacts: []ArtifactEnum{ResourceArtifact}},
				},
			},
			Expected: setDocumentationChildrenExpected{Children: []string{
				"[aci_included](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/included)",
			}},
		},
		{
			Name: "test_documentation_include_overrides_exclusions",
			Input: setDocumentationChildrenInput{
				RnMap: map[string]any{
					"included-{name}": "l2:Included",
				},
				ChildrenIncludedInResource:   []string{"l2Included"},
				ExcludeChildren:              []string{"l2Included"},
				DocumentationExcludeChildren: []string{"l2Included"},
				DocumentationIncludeChildren: []string{"l2Included"},
				StoreClasses: map[string]Class{
					"l2Included": {ResourceName: "included", Artifacts: []ArtifactEnum{ResourceArtifact}},
				},
			},
			Expected: setDocumentationChildrenExpected{Children: []string{
				"[aci_included](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/included)",
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
			class.ClassDefinition.ExcludeChildren = input.ExcludeChildren
			class.ClassDefinition.Documentation.ExcludeChildren = input.DocumentationExcludeChildren
			class.ClassDefinition.Documentation.IncludeChildren = input.DocumentationIncludeChildren
			if input.RnMap != nil {
				class.MetaFileContent = map[string]any{"rnMap": input.RnMap}
			}

			ds := &DataStore{Classes: input.StoreClasses}

			class.Documentation.setChildren(&class, ds)

			assert.Equal(t, expected.Children, class.Documentation.Children, test.MessageEqual(expected.Children, class.Documentation.Children, testCase.Name))
		})
	}
}

type documentationChildLinkInput struct {
	ChildName    string
	ForceInclude bool
}

type documentationChildLinkExpected struct {
	Link string
	OK   bool
}

func TestDocumentationChildLink(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	class := Class{
		ClassDefinition: ClassDefinition{
			ExcludeChildren: []string{"classExcluded"},
		},
	}
	ds := &DataStore{
		Classes: map[string]Class{
			"validChild":     {ResourceName: "valid_child", Artifacts: []ArtifactEnum{ResourceArtifact}},
			"excludedChild":  {ResourceName: "excluded_child", Artifacts: []ArtifactEnum{ResourceArtifact}},
			"embeddedChild":  {ResourceName: "embedded_child", Artifacts: []ArtifactEnum{ResourceArtifact}},
			"classExcluded":  {ResourceName: "class_excluded", Artifacts: []ArtifactEnum{ResourceArtifact}},
			"namedOnlyChild": {ResourceName: "named_only_child"},
		},
		ClassDefinitions: map[string]ClassDefinition{
			"legacyChild":         {ResourceName: "legacy_child"},
			"datasourceOnlyChild": {ResourceName: "datasource_only_child", Artifacts: []ArtifactEnum{DatasourceArtifact}},
		},
	}
	documentationExcludes := map[string]struct{}{"excludedChild": {}}
	childrenIncludedInResource := map[string]struct{}{"embeddedChild": {}}

	testCases := []test.TestCase{
		{
			Name:  "test_valid_child",
			Input: documentationChildLinkInput{ChildName: "validChild"},
			Expected: documentationChildLinkExpected{
				Link: "[aci_valid_child](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/valid_child)",
				OK:   true,
			},
		},
		{Name: "test_documentation_excluded_child", Input: documentationChildLinkInput{ChildName: "excludedChild"}, Expected: documentationChildLinkExpected{}},
		{Name: "test_embedded_child", Input: documentationChildLinkInput{ChildName: "embeddedChild"}, Expected: documentationChildLinkExpected{}},
		{Name: "test_class_excluded_child", Input: documentationChildLinkInput{ChildName: "classExcluded"}, Expected: documentationChildLinkExpected{}},
		{
			Name:  "test_forced_child_bypasses_exclusions",
			Input: documentationChildLinkInput{ChildName: "excludedChild", ForceInclude: true},
			Expected: documentationChildLinkExpected{
				Link: "[aci_excluded_child](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/excluded_child)",
				OK:   true,
			},
		},
		{Name: "test_unknown_child", Input: documentationChildLinkInput{}, Expected: documentationChildLinkExpected{}},
		{Name: "test_child_without_resource", Input: documentationChildLinkInput{ChildName: "noResourceChild"}, Expected: documentationChildLinkExpected{}},
		{Name: "test_child_without_resource_artifact", Input: documentationChildLinkInput{ChildName: "namedOnlyChild"}, Expected: documentationChildLinkExpected{}},
		{
			Name:  "test_definition_only_child_with_legacy_resource",
			Input: documentationChildLinkInput{ChildName: "legacyChild"},
			Expected: documentationChildLinkExpected{
				Link: "[aci_legacy_child](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/legacy_child)",
				OK:   true,
			},
		},
		{Name: "test_definition_only_child_without_resource_artifact", Input: documentationChildLinkInput{ChildName: "datasourceOnlyChild"}, Expected: documentationChildLinkExpected{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(documentationChildLinkInput)
			expected := testCase.Expected.(documentationChildLinkExpected)

			link, ok := documentationChildLink(input.ChildName, input.ForceInclude, &class, ds, documentationExcludes, childrenIncludedInResource)

			assert.Equal(t, expected.Link, link)
			assert.Equal(t, expected.OK, ok)
		})
	}
}

func TestDocumentationLinks(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	assert.Equal(t,
		"[aci_tenant](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/tenant)",
		getResourceDocumentationLink("tenant"),
	)
	assert.Equal(t,
		"[aci_tenant](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/data-sources/tenant)",
		getDatasourceDocumentationLink("tenant"),
	)
	assert.Equal(t, "", getTerraformRegistryDocumentationLink("tenant", ArtifactEnum(99)))
	ds := &DataStore{ClassDefinitions: map[string]ClassDefinition{
		"vnsLDevVip": {
			ResourceName:  "l4_l7_device",
			Documentation: ClassDocumentationDefinition{RegistryPath: "l4_l7_devices"},
		},
	}}
	assert.Equal(t,
		"[aci_l4_l7_device](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l4_l7_devices)",
		ds.getTerraformRegistryDocumentationLink("vnsLDevVip", ResourceArtifact),
	)
	assert.Equal(t,
		"[fvTenant](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvTenant/overview)",
		getDevnetObjectDocumentationLink("fvTenant"),
	)
}

type setDeprecationWarningInput struct {
	Deprecated         bool
	DeprecatedVersions string
	Hidden             bool
	HiddenVersions     string
}

func TestSetDeprecationWarning(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	classNameLink := "[fvTenant](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvTenant/overview)"
	terraformSentence := "The `aci_tenant` resource and datasource are deprecated and will be removed in a future release of the provider."

	testCases := []test.TestCase{
		{
			Name:     "test_none",
			Input:    setDeprecationWarningInput{},
			Expected: "",
		},
		{
			Name:     "test_deprecated_only_no_versions",
			Input:    setDeprecationWarningInput{Deprecated: true},
			Expected: "The " + classNameLink + " class is deprecated by the APIC API and may be removed in a future APIC release. " + terraformSentence,
		},
		{
			Name:     "test_deprecated_only_with_versions",
			Input:    setDeprecationWarningInput{Deprecated: true, DeprecatedVersions: "6.0(1a)-"},
			Expected: "The " + classNameLink + " class is deprecated by the APIC API as of version 6.0(1a) and later and may be removed in a future APIC release. " + terraformSentence,
		},
		{
			Name:     "test_hidden_only_no_versions",
			Input:    setDeprecationWarningInput{Hidden: true},
			Expected: "The " + classNameLink + " class is no longer accepted by the APIC API. " + terraformSentence,
		},
		{
			Name:     "test_hidden_only_with_versions",
			Input:    setDeprecationWarningInput{Hidden: true, HiddenVersions: "6.0(1a)-"},
			Expected: "The " + classNameLink + " class is no longer accepted by the APIC API as of version 6.0(1a) and later. " + terraformSentence,
		},
		{
			Name:     "test_hidden_supersedes_deprecated",
			Input:    setDeprecationWarningInput{Hidden: true, Deprecated: true, DeprecatedVersions: "5.0(1a)-"},
			Expected: "The " + classNameLink + " class is no longer accepted by the APIC API. " + terraformSentence,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setDeprecationWarningInput)

			class := Class{
				Name:         testClassName("fvTenant"),
				ResourceName: "tenant",
				Deprecated:   input.Deprecated,
				Hidden:       input.Hidden,
			}
			if input.DeprecatedVersions != "" {
				parsed, err := NewVersions(input.DeprecatedVersions)
				assert.NoError(t, err)
				class.DeprecatedVersions = parsed
			}
			if input.HiddenVersions != "" {
				parsed, err := NewVersions(input.HiddenVersions)
				assert.NoError(t, err)
				class.HiddenVersions = parsed
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
				ResourceDescription:   "Manages ACI Application EPG",
				DatasourceDescription: "Data source for ACI Application EPG",
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
				DatasourceDescription: "Data source for ACI Application EPG",
			},
		},
		{
			Name: "test_datasource_only",
			Input: setDescriptionInput{
				Label:           "Application EPG",
				DataDescription: "Returns matching objects.",
			},
			Expected: setDescriptionExpected{
				ResourceDescription:   "Manages ACI Application EPG",
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
	MetaFormats     []any
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
				MetaFormats: []any{"uni/tn-{name}/ap-{name}/epg-{name}"},
			},
			Expected: []string{"uni/tn-{name}/ap-{name}/epg-{name}"},
		},
		{
			Name: "test_meta_multiple_sorted",
			Input: setDnFormatsInput{
				MetaFormats: []any{
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
				MetaFormats: []any{"b/{name}", "a/{name}", "c/{name}"},
			},
			Expected: []string{"a/{name}", "b/{name}", "c/{name}"},
		},
		{
			Name: "test_override_replaces_meta",
			Input: setDnFormatsInput{
				MetaFormats:     []any{"meta/{name}"},
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
			Name: "test_override_over_cap_is_not_truncated",
			Input: setDnFormatsInput{
				OverrideFormats: []string{
					tooManyNotice,
					"f/{name}", "e/{name}", "d/{name}", "c/{name}", "b/{name}",
				},
			},
			Expected: []string{
				tooManyNotice,
				"b/{name}", "c/{name}", "d/{name}", "e/{name}", "f/{name}",
			},
		},
		{
			Name: "test_at_cap_no_notice",
			Input: setDnFormatsInput{
				MetaFormats: []any{"e/{name}", "d/{name}", "c/{name}", "b/{name}", "a/{name}"},
			},
			Expected: []string{"a/{name}", "b/{name}", "c/{name}", "d/{name}", "e/{name}"},
		},
		{
			Name: "test_over_cap_prepends_notice_and_truncates",
			Input: setDnFormatsInput{
				MetaFormats: []any{"f/{name}", "e/{name}", "d/{name}", "c/{name}", "b/{name}", "a/{name}"},
			},
			Expected: []string{
				tooManyNotice,
				"a/{name}", "b/{name}", "c/{name}", "d/{name}", "e/{name}",
			},
		},
		{
			Name: "test_meta_non_string_entries_skipped",
			Input: setDnFormatsInput{
				MetaFormats: []any{"a/{name}", 42, "b/{name}"},
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
				class.MetaFileContent = map[string]any{"dnFormats": input.MetaFormats}
			}

			class.Documentation.setClassName(&class)
			class.Documentation.setDnFormats(&class)

			assert.Equal(t, expected, class.Documentation.DnFormats, test.MessageEqual(expected, class.Documentation.DnFormats, testCase.Name))
			expectedImportDnFormat := ""
			for _, dnFormat := range expected {
				if !strings.HasPrefix(dnFormat, "Too many DN formats") {
					expectedImportDnFormat = dnFormat
					break
				}
			}
			assert.Equal(t, expectedImportDnFormat, class.Documentation.ImportDnFormat, test.MessageEqual(expectedImportDnFormat, class.Documentation.ImportDnFormat, testCase.Name+"/import"))
		})
	}
}

func TestSetImportDnFormat(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	tenantName := testClassName("fvTenant")
	pkiEpName := testClassName("pkiEp")
	aepgName := testClassName("fvAEPg")
	associatedSiteName := testClassName("fvSiteAssociated")
	ds := &DataStore{Classes: map[string]Class{
		"fvTenant": {Name: tenantName, MetaFileContent: map[string]any{"dnFormats": []any{"uni/tn-{name}"}}},
		"pkiEp":    {Name: pkiEpName, MetaFileContent: map[string]any{"dnFormats": []any{"uni/userext/pkiext"}}},
		"fvAEPg": {
			Name:            aepgName,
			MetaFileContent: map[string]any{"dnFormats": []any{"uni/tn-{name}/ap-{name}/epg-{name}"}},
		},
		"fvSiteAssociated": {
			Name:           associatedSiteName,
			RnFormat:       "stAsc",
			ExampleContext: ExampleContext{ParentClass: aepgName, ParentDnMode: ParentDnFromReference},
			MetaFileContent: map[string]any{"dnFormats": []any{
				"uni/tn-{name}/BD-{name}/stAsc",
				"uni/tn-{name}/ap-{name}/epg-{name}/stAsc",
			}},
		},
	}}

	testCases := []test.TestCase{
		{
			Name: "test_selected_parent_format",
			Input: Class{
				ExampleContext: ExampleContext{ParentClass: tenantName},
				RnFormat:       "annotationKey-[{key}]",
				MetaFileContent: map[string]any{"dnFormats": []any{
					"uni/tn-{name}/annotationKey-[{key}]",
					"uni/tn-{name}/BD-{name}/annotationKey-[{key}]",
				}},
			},
			Expected: "uni/tn-{name}/annotationKey-[{key}]",
		},
		{
			Name: "test_selected_parent_format_outside_displayed_formats",
			Input: Class{
				ExampleContext: ExampleContext{ParentClass: tenantName},
				RnFormat:       "annotationKey-[{key}]",
				MetaFileContent: map[string]any{"dnFormats": []any{
					"uni/a/annotationKey-[{key}]",
					"uni/b/annotationKey-[{key}]",
					"uni/c/annotationKey-[{key}]",
					"uni/d/annotationKey-[{key}]",
					"uni/e/annotationKey-[{key}]",
					"uni/tn-{name}/annotationKey-[{key}]",
				}},
			},
			Expected: "uni/tn-{name}/annotationKey-[{key}]",
		},
		{
			Name: "test_parent_dn_default_format",
			Input: Class{
				RnFormat:       "keyring-{name}",
				ExampleContext: ExampleContext{ParentDnMode: ParentDnFromSchemaDefault},
				Properties: map[string]*Property{
					"parentDn": {Documentation: PropertyDocumentation{DefaultValues: []DefaultValue{{Value: "uni/userext/pkiext"}}}},
				},
				MetaFileContent: map[string]any{"dnFormats": []any{
					"uni/tn-{name}/certstore/keyring-{name}",
					"uni/userext/pkiext/keyring-{name}",
				}},
			},
			Expected: "uni/userext/pkiext/keyring-{name}",
		},
		{
			Name: "test_nested_selected_parent_chain",
			Input: Class{
				RnFormat:       "site-{siteId}",
				ExampleContext: ExampleContext{ParentClass: associatedSiteName, ParentDnMode: ParentDnFromReference},
				MetaFileContent: map[string]any{"dnFormats": []any{
					"uni/tn-{name}/BD-{name}/stAsc/site-{siteId}",
					"uni/tn-{name}/ap-{name}/epg-{name}/stAsc/site-{siteId}",
				}},
			},
			Expected: "uni/tn-{name}/ap-{name}/epg-{name}/stAsc/site-{siteId}",
		},
		{
			Name: "test_parameterized_parent_dn_default_format",
			Input: Class{
				RnFormat:       "fabricExtConnP-{id}",
				ExampleContext: ExampleContext{ParentDnMode: ParentDnFromSchemaDefault},
				Properties: map[string]*Property{
					"parentDn": {Documentation: PropertyDocumentation{DefaultValues: []DefaultValue{{Value: "uni/tn-infra"}}}},
				},
				MetaFileContent: map[string]any{"dnFormats": []any{"uni/tn-{name}/fabricExtConnP-{id}"}},
			},
			Expected: "uni/tn-infra/fabricExtConnP-{id}",
		},
		{
			Name:     "test_first_actual_format_fallback",
			Input:    Class{MetaFileContent: map[string]any{"dnFormats": []any{"uni/infra/policy-{name}"}}},
			Expected: "uni/infra/policy-{name}",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			class := testCase.Input.(Class)
			class.Documentation.setImportDnFormat(&class, ds)
			assert.Equal(t, testCase.Expected, class.Documentation.ImportDnFormat, test.MessageEqual(testCase.Expected, class.Documentation.ImportDnFormat, testCase.Name))
		})
	}
}

func TestSetMigrationWarning(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	migrationWarning := "This resource has been migrated to the terraform plugin protocol version 6, refer to the [migration guide](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/guides/migration) for more details and implications for already managed resources."

	testCases := []test.TestCase{
		{
			Name:     "test_not_migration",
			Input:    false,
			Expected: "",
		},
		{
			Name:     "test_migration",
			Input:    true,
			Expected: migrationWarning,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			migrationSource := UndefinedMigrationSource
			if testCase.Input.(bool) {
				migrationSource = FromSDKv2
			}
			class := Class{
				Name:            testClassName("fvTenant"),
				MigrationSource: migrationSource,
			}
			class.Documentation.setMigrationWarning(&class)

			expected := testCase.Expected.(string)
			assert.Equal(t, expected, class.Documentation.MigrationWarning, test.MessageEqual(expected, class.Documentation.MigrationWarning, testCase.Name))
		})
	}
}

type setParentDnsInput struct {
	Parents          []string
	VariantParents   []string
	StoreClasses     map[string]Class
	StoreDefinitions map[string]ClassDefinition
}

type setParentDnsExpected struct {
	Notice    string
	Resources []string
	Classes   []string
}

func TestSetParentDns(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	link := func(className string) string {
		return fmt.Sprintf("[%s](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/%s/overview)", className, className)
	}
	resourceEntry := func(resourceName, className string) string {
		return fmt.Sprintf("[aci_%s](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/%s) (%s)", resourceName, resourceName, link(className))
	}
	additionalParents := "Additional parent options are available; see model documentation for all possible parents of " + link("fvFoo") + "."

	// Helpers to build large input sets for the over-cap cases.
	manyParentNames := func(prefix string, count int) []string {
		out := make([]string, count)
		for i := 0; i < count; i++ {
			out[i] = fmt.Sprintf("%s%02d", prefix, i)
		}
		return out
	}
	manyResourceStore := func(names []string) map[string]Class {
		store := make(map[string]Class, len(names))
		for _, n := range names {
			store[n] = Class{ResourceName: n, Artifacts: []ArtifactEnum{ResourceArtifact}}
		}
		return store
	}
	resourceEntries := func(names []string) []string {
		entries := make([]string, 0, len(names))
		for _, name := range names {
			entries = append(entries, resourceEntry(name, name))
		}
		return entries
	}
	classEntries := func(names []string) []string {
		entries := make([]string, 0, len(names))
		for _, name := range names {
			entries = append(entries, link(name))
		}
		return entries
	}

	testCases := []test.TestCase{
		{
			Name:     "test_no_parents",
			Input:    setParentDnsInput{},
			Expected: setParentDnsExpected{},
		},
		{
			Name: "test_single_resource",
			Input: setParentDnsInput{
				Parents: []string{"fvTenant"},
				StoreClasses: map[string]Class{
					"fvTenant": {ResourceName: "tenant", Artifacts: []ArtifactEnum{ResourceArtifact}},
				},
			},
			Expected: setParentDnsExpected{Resources: []string{resourceEntry("tenant", "fvTenant")}},
		},
		{
			Name: "test_definition_only_resource",
			Input: setParentDnsInput{
				Parents: []string{"l3extOut"},
				StoreDefinitions: map[string]ClassDefinition{
					"l3extOut": {ResourceName: "l3_outside"},
				},
			},
			Expected: setParentDnsExpected{Resources: []string{resourceEntry("l3_outside", "l3extOut")}},
		},
		{
			Name: "test_parent_dn_variants_replace_internal_wrapper_parents",
			Input: setParentDnsInput{
				Parents:        []string{"cloudCertStore", "pkiEp"},
				VariantParents: []string{"pkiEp", "fvTenant"},
				StoreClasses: map[string]Class{
					"fvTenant": {ResourceName: "tenant", Artifacts: []ArtifactEnum{ResourceArtifact}},
				},
			},
			Expected: setParentDnsExpected{
				Resources: []string{resourceEntry("tenant", "fvTenant")},
				Classes:   []string{link("pkiEp")},
			},
		},
		{
			Name: "test_multiple_resources_preserve_order",
			Input: setParentDnsInput{
				Parents: []string{"fvAlpha", "fvZeta"},
				StoreClasses: map[string]Class{
					"fvAlpha": {ResourceName: "alpha", Artifacts: []ArtifactEnum{ResourceArtifact}},
					"fvZeta":  {ResourceName: "zeta", Artifacts: []ArtifactEnum{ResourceArtifact}},
				},
			},
			Expected: setParentDnsExpected{Resources: []string{
				resourceEntry("alpha", "fvAlpha"),
				resourceEntry("zeta", "fvZeta"),
			}},
		},
		{
			Name: "test_only_class_only_with_note",
			Input: setParentDnsInput{
				Parents:      []string{"fvAlpha", "fvZeta"},
				StoreClasses: map[string]Class{},
			},
			Expected: setParentDnsExpected{Classes: []string{link("fvAlpha"), link("fvZeta")}},
		},
		{
			Name: "test_unknown_in_store_treated_as_class_only",
			Input: setParentDnsInput{
				Parents: []string{"fvKnown", "fvUnknown"},
				StoreClasses: map[string]Class{
					"fvKnown": {ResourceName: "known", Artifacts: []ArtifactEnum{ResourceArtifact}},
				},
			},
			Expected: setParentDnsExpected{
				Resources: []string{resourceEntry("known", "fvKnown")},
				Classes:   []string{link("fvUnknown")},
			},
		},
		{
			Name: "test_known_without_resource_name_treated_as_class_only",
			Input: setParentDnsInput{
				Parents: []string{"fvNoResource"},
				StoreClasses: map[string]Class{
					"fvNoResource": {ResourceName: ""},
				},
			},
			Expected: setParentDnsExpected{Classes: []string{link("fvNoResource")}},
		},
		{
			Name: "test_known_without_resource_artifact_treated_as_class_only",
			Input: setParentDnsInput{
				Parents: []string{"fvNoResource"},
				StoreClasses: map[string]Class{
					"fvNoResource": {ResourceName: "not_rendered"},
				},
			},
			Expected: setParentDnsExpected{Classes: []string{link("fvNoResource")}},
		},
		{
			Name: "test_resources_at_cap_have_no_notice",
			Input: setParentDnsInput{
				Parents:      manyParentNames("fvRes", 20),
				StoreClasses: manyResourceStore(manyParentNames("fvRes", 20)),
			},
			Expected: setParentDnsExpected{Resources: resourceEntries(manyParentNames("fvRes", 20))},
		},
		{
			Name: "test_resources_over_cap_are_sampled",
			Input: setParentDnsInput{
				Parents:      manyParentNames("fvRes", 21),
				StoreClasses: manyResourceStore(manyParentNames("fvRes", 21)),
			},
			Expected: setParentDnsExpected{
				Notice:    additionalParents,
				Resources: resourceEntries(manyParentNames("fvRes", 20)),
			},
		},
		{
			Name: "test_classes_over_cap_are_sampled",
			Input: setParentDnsInput{
				Parents:      manyParentNames("fvCls", 21),
				StoreClasses: map[string]Class{},
			},
			Expected: setParentDnsExpected{
				Notice:  additionalParents,
				Classes: classEntries(manyParentNames("fvCls", 20)),
			},
		},
		{
			Name: "test_combined_groups_over_cap_prioritize_resources",
			Input: setParentDnsInput{
				Parents:      append(manyParentNames("fvCls", 10), manyParentNames("fvRes", 11)...),
				StoreClasses: manyResourceStore(manyParentNames("fvRes", 11)),
			},
			Expected: setParentDnsExpected{
				Notice:    additionalParents,
				Resources: resourceEntries(manyParentNames("fvRes", 11)),
				Classes:   classEntries(manyParentNames("fvCls", 9)),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setParentDnsInput)
			expected := testCase.Expected.(setParentDnsExpected)

			parents := make([]*ClassName, 0, len(input.Parents))
			for _, name := range input.Parents {
				parents = append(parents, testClassName(name))
			}

			variants := make([]ParentDnVariant, 0, len(input.VariantParents))
			for _, name := range input.VariantParents {
				variants = append(variants, ParentDnVariant{ParentClass: testClassName(name)})
			}

			class := Class{
				Name:             testClassName("fvFoo"),
				Parents:          parents,
				ParentDnVariants: variants,
			}
			ds := &DataStore{Classes: input.StoreClasses, ClassDefinitions: input.StoreDefinitions}

			class.Documentation.setClassName(&class)
			class.Documentation.setParentDns(&class, ds)

			actual := setParentDnsExpected{
				Notice:    class.Documentation.ParentDnNotice,
				Resources: class.Documentation.ParentDnResources,
				Classes:   class.Documentation.ParentDnClasses,
			}
			assert.Equal(t, expected, actual, test.MessageEqual(expected, actual, testCase.Name))
		})
	}
}

func TestSetClassDocumentationSupportedVersions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{Name: "test_nil_versions", Input: "", Expected: ""},
		{Name: "test_single_version", Input: "1.0(1e)", Expected: "1.0(1e)"},
		{Name: "test_bounded_range", Input: "4.2(7f)-4.2(7w)", Expected: "4.2(7f) to 4.2(7w)"},
		{Name: "test_unbounded_max", Input: "5.2(1g)-", Expected: "5.2(1g) and later."},
		{Name: "test_unbounded_min", Input: "-4.2(7w)", Expected: "up to 4.2(7w)"},
		{Name: "test_multiple_ranges_preserve_source_order", Input: "5.2(1g)-,3.2(10e)-3.2(10g)", Expected: "5.2(1g) and later, 3.2(10e) to 3.2(10g)"},
		{Name: "test_mixed_single_and_range", Input: "1.0(1e),4.2(7f)-", Expected: "1.0(1e), 4.2(7f) and later."},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			raw := testCase.Input.(string)
			expected := testCase.Expected.(string)

			class := Class{Name: testClassName("fvFoo")}
			if raw != "" {
				versions, err := NewVersions(raw)
				assert.NoError(t, err)
				class.SupportedVersions = versions
			}

			class.Documentation.setSupportedVersions(&class)

			assert.Equal(t, expected, class.Documentation.SupportedVersions, test.MessageEqual(expected, class.Documentation.SupportedVersions, testCase.Name))
		})
	}
}

type setDescriptionWhenDefinedAsChildInput struct {
	Label              string
	ResourceName       string
	ResourceNameNested string
	Artifacts          []ArtifactEnum
	IsSingleNested     bool
	Relation           Relation
	ClassDefinition    ClassDefinition
	StoreClasses       map[string]Class
	StoreDefinitions   map[string]ClassDefinition
}

type setDescriptionWhenDefinedAsChildExpected struct {
	Description                string
	RelationTargetDescriptions []string
}

func TestSetDescriptionWhenDefinedAsChild(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	fooLink := "[fvFoo](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvFoo/overview)"
	bdLink := "[fvBD](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvBD/overview)"
	ctxLink := "[fvCtx](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvCtx/overview)"

	testCases := []test.TestCase{
		{
			Name: "test_non_relational_list",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Foo",
				ResourceName:       "foo",
				ResourceNameNested: "foo",
				Artifacts:          []ArtifactEnum{ResourceArtifact},
				IsSingleNested:     false,
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{Description: "`foo` - (list) A list of Foos (ACI object " + fooLink + ") which can also be configured using a separate [aci_foo](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/foo) resource."},
		},
		{
			Name: "test_non_relational_list_without_resource_artifact",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Foo",
				ResourceName:       "foo",
				ResourceNameNested: "foo",
				IsSingleNested:     false,
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{Description: "`foo` - (list) A list of Foos (ACI object " + fooLink + ")."},
		},
		{
			Name: "test_non_relational_map",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Foo",
				ResourceName:       "foo",
				ResourceNameNested: "foo",
				Artifacts:          []ArtifactEnum{ResourceArtifact},
				IsSingleNested:     true,
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{Description: "`foo` - (map) A map of Foo (ACI object " + fooLink + ") which can also be configured using a separate [aci_foo](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/foo) resource."},
		},
		{
			Name: "test_non_relational_map_without_resource_artifact",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Foo",
				ResourceName:       "foo",
				ResourceNameNested: "foo",
				IsSingleNested:     true,
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{Description: "`foo` - (map) A map of Foo (ACI object " + fooLink + ")."},
		},
		{
			Name: "test_relational_list",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Relation From Foo To BD",
				ResourceName:       "relation_from_foo_to_bd",
				ResourceNameNested: "relation_to_bd",
				IsSingleNested:     false,
				Relation:           Relation{RelationalClass: true, ToClasses: []*ClassName{testClassName("fvBD")}},
				StoreClasses: map[string]Class{
					"fvBD": {ResourceName: "bridge_domain", Artifacts: []ArtifactEnum{ResourceArtifact}, Documentation: ClassDocumentation{Label: "Bridge Domain"}},
				},
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{Description: "`relation_to_bd` - (list) A list of Relation From Foo To BDs (ACI object " + fooLink + ") pointing to Bridge Domain (ACI Object " + bdLink + ") which can be configured using the [aci_bridge_domain](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/bridge_domain) resource."},
		},
		{
			Name: "test_relational_map",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Relation From Foo To BD",
				ResourceName:       "relation_from_foo_to_bd",
				ResourceNameNested: "relation_to_bd",
				IsSingleNested:     true,
				Relation:           Relation{RelationalClass: true, ToClasses: []*ClassName{testClassName("fvBD")}},
				StoreClasses: map[string]Class{
					"fvBD": {ResourceName: "bridge_domain", Artifacts: []ArtifactEnum{ResourceArtifact}, Documentation: ClassDocumentation{Label: "Bridge Domain"}},
				},
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{Description: "`relation_to_bd` - (map) A map of Relation From Foo To BD (ACI object " + fooLink + ") pointing to Bridge Domain (ACI Object " + bdLink + ") which can be configured using the [aci_bridge_domain](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/bridge_domain) resource."},
		},
		{
			Name: "test_relational_map_without_resource_artifact",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Relation From Foo To BD",
				ResourceName:       "relation_from_foo_to_bd",
				ResourceNameNested: "relation_to_bd",
				IsSingleNested:     true,
				Relation:           Relation{RelationalClass: true, ToClasses: []*ClassName{testClassName("fvBD")}},
				StoreClasses: map[string]Class{
					"fvBD": {ResourceName: "bridge_domain", Documentation: ClassDocumentation{Label: "Bridge Domain"}},
				},
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{Description: "`relation_to_bd` - (map) A map of Relation From Foo To BD (ACI object " + fooLink + ") pointing to Bridge Domain (ACI Object " + bdLink + ")."},
		},
		{
			Name: "test_relational_toclass_missing_in_store",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Relation From Foo To BD",
				ResourceName:       "relation_from_foo_to_bd",
				ResourceNameNested: "relation_to_bd",
				IsSingleNested:     false,
				Relation:           Relation{RelationalClass: true, ToClasses: []*ClassName{testClassName("fvBD")}},
				StoreClasses:       map[string]Class{},
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{Description: "`relation_to_bd` - (list) A list of Relation From Foo To BDs (ACI object " + fooLink + ") pointing to fvBD (ACI Object " + bdLink + ")."},
		},
		{
			Name: "test_relational_toclass_uses_definition_label",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Relation From Foo To Monitoring Policy",
				ResourceName:       "relation_from_foo_to_monitoring_policy",
				ResourceNameNested: "relation_to_monitoring_policy",
				IsSingleNested:     true,
				Relation:           Relation{RelationalClass: true, ToClasses: []*ClassName{testClassName("monEPGPol")}},
				StoreDefinitions: map[string]ClassDefinition{
					"monEPGPol": {ResourceName: "monitoring_policy"},
				},
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{Description: "`relation_to_monitoring_policy` - (map) A map of Relation From Foo To Monitoring Policy (ACI object " + fooLink + ") pointing to Monitoring Policy (ACI Object [monEPGPol](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/monEPGPol/overview)) which can be configured using the [aci_monitoring_policy](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/monitoring_policy) resource."},
		},
		{
			Name: "test_relational_toclass_uses_definition_only_resource",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Relation From Foo To BD",
				ResourceName:       "relation_from_foo_to_bd",
				ResourceNameNested: "relation_to_bd",
				Relation:           Relation{RelationalClass: true, ToClasses: []*ClassName{testClassName("fvBD")}},
				StoreDefinitions: map[string]ClassDefinition{
					"fvBD": {ResourceName: "bridge_domain"},
				},
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{Description: "`relation_to_bd` - (list) A list of Relation From Foo To BDs (ACI object " + fooLink + ") pointing to Bridge Domain (ACI Object " + bdLink + ") which can be configured using the [aci_bridge_domain](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/bridge_domain) resource."},
		},
		{
			Name: "test_relational_toclass_without_resource_artifact",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Relation From Foo To BD",
				ResourceName:       "relation_from_foo_to_bd",
				ResourceNameNested: "relation_to_bd",
				IsSingleNested:     false,
				Relation:           Relation{RelationalClass: true, ToClasses: []*ClassName{testClassName("fvBD")}},
				StoreClasses: map[string]Class{
					"fvBD": {ResourceName: "bridge_domain", Documentation: ClassDocumentation{Label: "Bridge Domain"}},
				},
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{Description: "`relation_to_bd` - (list) A list of Relation From Foo To BDs (ACI object " + fooLink + ") pointing to Bridge Domain (ACI Object " + bdLink + ")."},
		},
		{
			Name: "test_relational_toclass_empty_label",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Relation From Foo To BD",
				ResourceName:       "relation_from_foo_to_bd",
				ResourceNameNested: "relation_to_bd",
				IsSingleNested:     false,
				Relation:           Relation{RelationalClass: true, ToClasses: []*ClassName{testClassName("fvBD")}},
				StoreClasses: map[string]Class{
					"fvBD": {ResourceName: "bridge_domain", Artifacts: []ArtifactEnum{ResourceArtifact}},
				},
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{Description: "`relation_to_bd` - (list) A list of Relation From Foo To BDs (ACI object " + fooLink + ") pointing to fvBD (ACI Object " + bdLink + ") which can be configured using the [aci_bridge_domain](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/bridge_domain) resource."},
		},
		{
			Name: "test_relational_multi_target_preserves_definition_order",
			Input: setDescriptionWhenDefinedAsChildInput{
				Label:              "Domain",
				ResourceName:       "relation_to_domain",
				ResourceNameNested: "relations_to_domains",
				Relation: Relation{RelationalClass: true, ToClasses: []*ClassName{
					testClassName("fvBD"),
					testClassName("fvCtx"),
				}},
				ClassDefinition: ClassDefinition{RelationInfo: RelationInfoDefinition{ToClasses: []string{"fvCtx", "fvBD"}}},
				StoreClasses: map[string]Class{
					"fvCtx": {ResourceName: "vrf", Artifacts: []ArtifactEnum{ResourceArtifact}},
					"fvBD":  {ResourceName: "bridge_domain"},
				},
			},
			Expected: setDescriptionWhenDefinedAsChildExpected{
				Description: "`relations_to_domains` - (list) A list of Domains (ACI object " + fooLink + ").",
				RelationTargetDescriptions: []string{
					ctxLink + " which can be configured using the [aci_vrf](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/vrf) resource",
					bdLink + ". Currently there is no resource available to configure this ACI object",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setDescriptionWhenDefinedAsChildInput)
			expected := testCase.Expected.(setDescriptionWhenDefinedAsChildExpected)

			class := Class{
				Name:                             testClassName("fvFoo"),
				ResourceName:                     input.ResourceName,
				ResourceNameNested:               input.ResourceNameNested,
				Artifacts:                        input.Artifacts,
				IsSingleNestedWhenDefinedAsChild: input.IsSingleNested,
				Relation:                         input.Relation,
				ClassDefinition:                  input.ClassDefinition,
			}
			class.Documentation.Label = input.Label
			ds := &DataStore{Classes: input.StoreClasses, ClassDefinitions: input.StoreDefinitions}

			class.Documentation.setClassName(&class)
			class.Documentation.setDescriptionWhenDefinedAsChild(&class, ds)

			assert.Equal(t, expected.Description, class.Documentation.DescriptionWhenDefinedAsChild, test.MessageEqual(expected.Description, class.Documentation.DescriptionWhenDefinedAsChild, testCase.Name))
			assert.Equal(t, expected.RelationTargetDescriptions, class.Documentation.RelationTargetDescriptions, test.MessageEqual(expected.RelationTargetDescriptions, class.Documentation.RelationTargetDescriptions, testCase.Name))
		})
	}
}

type nestedDocumentationClassExpected struct {
	Class          string
	Depth          int
	Indent         string
	PropertyIndent string
}

func TestNestedDocumentationClasses(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_without_children",
			Input: Class{
				Name: testClassName("fvTenant"),
			},
			Expected: []nestedDocumentationClassExpected{},
		},
		{
			Name: "test_depth_first_with_repeated_class",
			Input: Class{
				Name: testClassName("fvTenant"),
				Children: []*ClassName{
					testClassName("fvRsTenantMonPol"),
					testClassName("tagAnnotation"),
				},
			},
			Expected: []nestedDocumentationClassExpected{
				{Class: "fvRsTenantMonPol", Depth: 0, PropertyIndent: "    "},
				{Class: "tagAnnotation", Depth: 1, Indent: "    ", PropertyIndent: "        "},
				{Class: "tagTag", Depth: 1, Indent: "    ", PropertyIndent: "        "},
				{Class: "tagAnnotation", Depth: 0, PropertyIndent: "    "},
			},
		},
	}

	classes := map[string]Class{
		"fvRsTenantMonPol": {
			Name: testClassName("fvRsTenantMonPol"),
			Children: []*ClassName{
				testClassName("tagAnnotation"),
				testClassName("tagTag"),
			},
		},
		"tagAnnotation": {Name: testClassName("tagAnnotation")},
		"tagTag":        {Name: testClassName("tagTag")},
	}
	ds := &DataStore{Classes: classes}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			class := testCase.Input.(Class)
			expected := testCase.Expected.([]nestedDocumentationClassExpected)
			nestedClasses := class.NestedDocumentationClasses(ds)
			actual := make([]nestedDocumentationClassExpected, 0, len(nestedClasses))
			for _, nestedClass := range nestedClasses {
				actual = append(actual, nestedDocumentationClassExpected{
					Class:          nestedClass.Class.Name.String(),
					Depth:          nestedClass.Depth,
					Indent:         nestedClass.Indent(0),
					PropertyIndent: nestedClass.Indent(1),
				})
			}
			assert.Equal(t, expected, actual, test.MessageEqual(expected, actual, testCase.Name))
		})
	}
}

func TestNestedDocumentationClassesByResourceSection(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	parent := Class{
		Name: testClassName("fvTenant"),
		Children: []*ClassName{
			testClassName("requiredChild"),
			testClassName("optionalChild"),
		},
	}
	classes := map[string]Class{
		"requiredChild": {
			Name:            testClassName("requiredChild"),
			RequiredAsChild: true,
			Children:        []*ClassName{testClassName("optionalGrandchild")},
		},
		"optionalChild": {
			Name:     testClassName("optionalChild"),
			Children: []*ClassName{testClassName("requiredGrandchild")},
		},
		"optionalGrandchild": {Name: testClassName("optionalGrandchild")},
		"requiredGrandchild": {Name: testClassName("requiredGrandchild"), RequiredAsChild: true},
	}
	ds := &DataStore{Classes: classes}

	testCases := []test.TestCase{
		{
			Name:     "test_required_root_keeps_optional_descendant",
			Input:    parent.RequiredNestedDocumentationClasses(ds),
			Expected: []nestedDocumentationClassExpected{{Class: "requiredChild", Depth: 0, PropertyIndent: "    "}, {Class: "optionalGrandchild", Depth: 1, Indent: "    ", PropertyIndent: "        "}},
		},
		{
			Name:     "test_optional_root_keeps_required_descendant",
			Input:    parent.OptionalNestedDocumentationClasses(ds),
			Expected: []nestedDocumentationClassExpected{{Class: "optionalChild", Depth: 0, PropertyIndent: "    "}, {Class: "requiredGrandchild", Depth: 1, Indent: "    ", PropertyIndent: "        "}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			nestedClasses := testCase.Input.([]NestedDocumentationClass)
			actual := make([]nestedDocumentationClassExpected, 0, len(nestedClasses))
			for _, nestedClass := range nestedClasses {
				actual = append(actual, nestedDocumentationClassExpected{
					Class:          nestedClass.Class.Name.String(),
					Depth:          nestedClass.Depth,
					Indent:         nestedClass.Indent(0),
					PropertyIndent: nestedClass.Indent(1),
				})
			}
			expected := testCase.Expected.([]nestedDocumentationClassExpected)
			assert.Equal(t, expected, actual, test.MessageEqual(expected, actual, testCase.Name))
		})
	}
}

func TestDocumentationProperties(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	parentDn := &Property{PropertyName: "parentDn"}
	name := &Property{PropertyName: "name"}
	description := &Property{PropertyName: "descr"}
	status := &Property{PropertyName: "status"}
	class := Class{
		Properties: map[string]*Property{
			"parentDn": parentDn,
			"name":     name,
			"descr":    description,
			"status":   status,
		},
		PropertiesRequired: []string{"parentDn", "name"},
		PropertiesOptional: []string{"descr"},
		PropertiesReadOnly: []string{"status"},
	}

	testCases := []test.TestCase{
		{Name: "test_required", Input: class.RequiredDocumentationProperties(), Expected: []*Property{name}},
		{Name: "test_optional", Input: class.OptionalDocumentationProperties(), Expected: []*Property{description}},
		{Name: "test_read_only", Input: class.ReadOnlyDocumentationProperties(), Expected: []*Property{status}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			actual := testCase.Input.([]*Property)
			expected := testCase.Expected.([]*Property)
			assert.Equal(t, expected, actual, test.MessageEqual(expected, actual, testCase.Name))
		})
	}
}
