package data

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

// testClassName creates a ClassName for testing purposes.
// Panics if the class name is invalid, which is acceptable in tests.
func testClassName(name string) *ClassName {
	cn, err := NewClassName(name)
	if err != nil {
		panic(err)
	}
	return cn
}

// classNamesToStrings converts a slice of ClassName pointers to a slice of strings for test assertions.
func classNamesToStrings(names []*ClassName) []string {
	result := make([]string, len(names))
	for i, n := range names {
		result[i] = n.String()
	}
	return result
}

func TestSetResourceNameFromLabelNoRelationWithIdentifier(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	class := Class{Name: testClassName("tagAnnotation"), IdentifiedBy: []string{"key"}}
	class.MetaFileContent = map[string]any{
		"label": "annotation",
	}

	err := class.setResourceName(ds)

	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, "annotation", class.ResourceName, test.MessageEqual("annotation", class.ResourceName, t.Name()))
	assert.Equal(t, "annotations", class.ResourceNameNested, test.MessageEqual("annotations", class.ResourceNameNested, t.Name()))
}

func TestSetResourceNameFromLabelNoRelationWithoutIdentifier(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	ds.GlobalMetaDefinition = GlobalMetaDefinition{
		NoMetaFile: map[string]string{
			"fvCtx": "vrf",
		},
	}
	class := Class{Name: testClassName("fvRsScope")}
	class.MetaFileContent = map[string]any{
		"label": "Private Network",
		"relationInfo": map[string]any{
			"type":   "named",
			"fromMo": "fv:ESg",
			"toMo":   "fv:Ctx",
		},
	}

	err := class.setRelation()
	assert.NoError(t, err, test.MessageUnexpectedError(err))

	err = class.setResourceName(ds)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, "relation_to_vrf", class.ResourceName, test.MessageEqual("relation_to_vrf", class.ResourceName, t.Name()))
	assert.Equal(t, "relation_to_vrf", class.ResourceNameNested, test.MessageEqual("relation_to_vrf", class.ResourceNameNested, t.Name()))
}

func TestSetResourceNameToRelation(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	ds.GlobalMetaDefinition = GlobalMetaDefinition{
		NoMetaFile: map[string]string{
			"vzBrCP": "contract",
		},
	}
	class := Class{Name: testClassName("fvRsCons"), IdentifiedBy: []string{"tnVzBrCPName"}}
	class.MetaFileContent = map[string]any{
		"label": "contract",
		"relationInfo": map[string]any{
			"type":   "named",
			"fromMo": "fv:EPg",
			"toMo":   "vz:BrCP",
		},
	}

	err := class.setRelation()
	assert.NoError(t, err, test.MessageUnexpectedError(err))

	err = class.setResourceName(ds)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, "relation_to_contract", class.ResourceName, test.MessageEqual("relation_to_contract", class.ResourceName, t.Name()))
	assert.Equal(t, "relation_to_contracts", class.ResourceNameNested, test.MessageEqual("relation_to_contracts", class.ResourceNameNested, t.Name()))
}

func TestSetResourceNameFromToRelation(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	ds.GlobalMetaDefinition = GlobalMetaDefinition{
		NoMetaFile: map[string]string{
			"fvCtx":               "vrf",
			"netflowAExporterPol": "netflow_exporter_policy",
		},
	}
	class := Class{Name: testClassName("netflowRsExporterToCtx")}
	class.MetaFileContent = map[string]any{
		"label": "netflow exporter",
		"relationInfo": map[string]any{
			"type":   "explicit",
			"fromMo": "netflow:AExporterPol",
			"toMo":   "fv:Ctx",
		},
		"isCreatableDeletable": "always",
	}

	err := class.setRelation()
	assert.NoError(t, err, test.MessageUnexpectedError(err))

	err = class.setResourceName(ds)
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, "relation_from_netflow_exporter_policy_to_vrf", class.ResourceName, test.MessageEqual("relation_from_netflow_exporter_policy_to_vrf", class.ResourceName, t.Name()))
	assert.Equal(t, "relation_to_vrf", class.ResourceNameNested, test.MessageEqual("relation_to_vrf", class.ResourceNameNested, t.Name()))
}

func TestSetResourceNameFromEmptyLabelError(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	class := Class{Name: testClassName("fvTenant")}
	class.MetaFileContent = map[string]any{"label": ""}

	err := class.setResourceName(ds)

	assert.EqualError(t, err, "failed to set resource name for class 'fvTenant': resource_name not defined and label not found")
}

func TestSetResourceNameFromNoLabelError(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	class := Class{Name: testClassName("fvTenant")}
	class.MetaFileContent = map[string]any{"no_label": ""}

	err := class.setResourceName(ds)

	assert.EqualError(t, err, "failed to set resource name for class 'fvTenant': resource_name not defined and label not found")
}

func TestSetResourceNameFromDefinitionOverride(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	class := Class{Name: testClassName("fvCtx"), IdentifiedBy: []string{"name"}}
	class.MetaFileContent = map[string]any{
		"label": "context",
	}
	class.ClassDefinition = ClassDefinition{ResourceName: "vrf"}

	err := class.setResourceName(ds)

	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, "vrf", class.ResourceName, test.MessageEqual("vrf", class.ResourceName, t.Name()))
	assert.Equal(t, "vrfs", class.ResourceNameNested, test.MessageEqual("vrfs", class.ResourceNameNested, t.Name()))
}

func TestSetResourceNameFromDefinitionOverrideWithoutIdentifier(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	class := Class{Name: testClassName("fvRsCtx")}
	class.MetaFileContent = map[string]any{
		"label": "context",
	}
	class.ClassDefinition = ClassDefinition{ResourceName: "vrf"}

	err := class.setResourceName(ds)

	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, "vrf", class.ResourceName, test.MessageEqual("vrf", class.ResourceName, t.Name()))
	assert.Equal(t, "vrf", class.ResourceNameNested, test.MessageEqual("vrf", class.ResourceNameNested, t.Name()))
}

func TestSetResourceNameFromDefinitionOverrideWithoutLabel(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	class := Class{Name: testClassName("fvRsCtx")}
	class.MetaFileContent = map[string]any{}
	class.ClassDefinition = ClassDefinition{ResourceName: "vrf"}

	err := class.setResourceName(ds)

	assert.NoError(t, err, test.MessageUnexpectedError(err))
	assert.Equal(t, "vrf", class.ResourceName, test.MessageEqual("vrf", class.ResourceName, t.Name()))
	assert.Equal(t, "vrf", class.ResourceNameNested, test.MessageEqual("vrf", class.ResourceNameNested, t.Name()))
}

type setRelationInput struct {
	ClassName       string
	MetaFileContent map[string]any
	ClassDefinition ClassDefinition
}

type setRelationExpected struct {
	FromClass       *ClassName
	ToClasses       []*ClassName
	Type            RelationshipTypeEnum
	IncludeFrom     bool
	RelationalClass bool
	Error           bool
	ErrorMsg        string
}

func TestSetRelation(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_no_relation",
			Input: setRelationInput{
				ClassName:       "fvTenant",
				MetaFileContent: nil,
			},
			Expected: setRelationExpected{
				FromClass:       nil,
				ToClasses:       nil,
				Type:            RelationshipTypeEnum(0),
				IncludeFrom:     false,
				RelationalClass: false,
			},
		},
		{
			Name: "test_include_from_false_type_named",
			Input: setRelationInput{
				ClassName: "fvRsCons",
				MetaFileContent: map[string]any{
					"label": "contract",
					"relationInfo": map[string]any{
						"type":   "named",
						"fromMo": "fv:EPg",
						"toMo":   "vz:BrCP",
					},
				},
			},
			Expected: setRelationExpected{
				FromClass:       testClassName("fvEPg"),
				ToClasses:       []*ClassName{testClassName("vzBrCP")},
				Type:            Named,
				IncludeFrom:     false,
				RelationalClass: true,
			},
		},
		{
			Name: "test_include_from_true_type_explicit",
			Input: setRelationInput{
				ClassName: "netflowRsExporterToCtx",
				MetaFileContent: map[string]any{
					"label": "netflow exporter",
					"relationInfo": map[string]any{
						"type":   "explicit",
						"fromMo": "netflow:AExporterPol",
						"toMo":   "fv:Ctx",
					},
					"isCreatableDeletable": "always",
				},
			},
			Expected: setRelationExpected{
				FromClass:       testClassName("netflowAExporterPol"),
				ToClasses:       []*ClassName{testClassName("fvCtx")},
				Type:            Explicit,
				IncludeFrom:     true,
				RelationalClass: true,
			},
		},
		{
			Name: "test_relation_to_classes_override",
			Input: setRelationInput{
				ClassName: "fvRsDomAtt",
				MetaFileContent: map[string]any{
					"label": "domain attachment",
					"relationInfo": map[string]any{
						"type":   "named",
						"fromMo": "fv:AEPg",
						"toMo":   "infra:DomP",
					},
				},
				ClassDefinition: ClassDefinition{
					RelationInfo: RelationInfoDefinition{
						ToClasses: []string{"vmm:DomP", "phys:DomP", "fc:DomP", "l2ext:DomP"},
					},
				},
			},
			Expected: setRelationExpected{
				FromClass:       testClassName("fvAEPg"),
				ToClasses:       []*ClassName{testClassName("fcDomP"), testClassName("l2extDomP"), testClassName("physDomP"), testClassName("vmmDomP")},
				Type:            Named,
				IncludeFrom:     false,
				RelationalClass: true,
			},
		},
		{
			Name: "test_relation_info_field_override",
			Input: setRelationInput{
				ClassName: "fvRsCons",
				MetaFileContent: map[string]any{
					"label": "contract",
					"relationInfo": map[string]any{
						"type":   "named",
						"fromMo": "fv:EPg",
						"toMo":   "vz:BrCP",
					},
				},
				ClassDefinition: ClassDefinition{
					RelationInfo: RelationInfoDefinition{ToClasses: []string{"vz:OOBBrCP"}},
				},
			},
			Expected: setRelationExpected{
				FromClass:       testClassName("fvEPg"),
				ToClasses:       []*ClassName{testClassName("vzOOBBrCP")},
				Type:            Named,
				IncludeFrom:     false,
				RelationalClass: true,
			},
		},
		{
			Name: "test_relation_info_full_definition_only",
			Input: setRelationInput{
				ClassName:       "fakeRsCustom",
				MetaFileContent: map[string]any{},
				ClassDefinition: ClassDefinition{
					RelationInfo: RelationInfoDefinition{
						Type:      Explicit,
						FromClass: "fv:AEPg",
						ToClasses: []string{"vz:BrCP"},
					},
				},
			},
			Expected: setRelationExpected{
				FromClass:       testClassName("fvAEPg"),
				ToClasses:       []*ClassName{testClassName("vzBrCP")},
				Type:            Explicit,
				IncludeFrom:     false,
				RelationalClass: true,
			},
		},
		{
			Name: "test_relation_info_missing_to_class_error",
			Input: setRelationInput{
				ClassName:       "fakeRsCustom",
				MetaFileContent: map[string]any{},
				ClassDefinition: ClassDefinition{
					RelationInfo: RelationInfoDefinition{
						Type:      Named,
						FromClass: "fv:AEPg",
					},
				},
			},
			Expected: setRelationExpected{
				Error:    true,
				ErrorMsg: "missing required relation_info fields: to_classes",
			},
		},
		{
			Name: "test_undefined_type_error",
			Input: setRelationInput{
				ClassName: "netflowRsExporterToCtx",
				MetaFileContent: map[string]any{
					"relationInfo": map[string]any{
						"type":   "undefinedType",
						"fromMo": "netflow:AExporterPol",
						"toMo":   "fv:Ctx",
					},
				},
			},
			Expected: setRelationExpected{
				Error:    true,
				ErrorMsg: `unknown relationship type "undefinedType"`,
			},
		},
		{
			Name: "test_relation_info_disabled_overrides_meta",
			Input: setRelationInput{
				ClassName: "fvRsCtx",
				MetaFileContent: map[string]any{
					"relationInfo": map[string]any{
						"type":   "named",
						"fromMo": "fv:EPg",
						"toMo":   "vz:BrCP",
					},
				},
				ClassDefinition: ClassDefinition{
					RelationInfo: RelationInfoDefinition{Disabled: true},
				},
			},
			Expected: setRelationExpected{
				FromClass:       nil,
				ToClasses:       nil,
				Type:            RelationshipTypeEnum(0),
				IncludeFrom:     false,
				RelationalClass: false,
			},
		},
		{
			Name: "test_relation_info_disabled_conflict_error",
			Input: setRelationInput{
				ClassName:       "fvRsCtx",
				MetaFileContent: map[string]any{},
				ClassDefinition: ClassDefinition{
					RelationInfo: RelationInfoDefinition{
						Disabled:  true,
						Type:      Named,
						FromClass: "fv:EPg",
						ToClasses: []string{"vz:BrCP"},
					},
				},
			},
			Expected: setRelationExpected{
				Error:    true,
				ErrorMsg: "relation_info.disabled is mutually exclusive with type, from_class, and to_classes",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setRelationInput)
			expected := testCase.Expected.(setRelationExpected)
			class := Class{Name: testClassName(input.ClassName)}
			class.MetaFileContent = input.MetaFileContent
			class.ClassDefinition = input.ClassDefinition

			err := class.setRelation()

			if expected.Error {
				assert.Error(t, err)
				if expected.ErrorMsg != "" {
					assert.ErrorContains(t, err, expected.ErrorMsg)
				}
				return
			}

			assert.NoError(t, err, test.MessageUnexpectedError(err))
			assert.Equal(t, expected.FromClass, class.Relation.FromClass, test.MessageEqual(expected.FromClass, class.Relation.FromClass, testCase.Name))
			assert.Equal(t, expected.ToClasses, class.Relation.ToClasses, test.MessageEqual(expected.ToClasses, class.Relation.ToClasses, testCase.Name))
			assert.Equal(t, expected.Type, class.Relation.Type, test.MessageEqual(expected.Type, class.Relation.Type, testCase.Name))
			assert.Equal(t, expected.IncludeFrom, class.Relation.IncludeFrom, test.MessageEqual(expected.IncludeFrom, class.Relation.IncludeFrom, testCase.Name))
			assert.Equal(t, expected.RelationalClass, class.Relation.RelationalClass, test.MessageEqual(expected.RelationalClass, class.Relation.RelationalClass, testCase.Name))
		})
	}
}

type setAllowDeleteInput struct {
	MetaFileContent map[string]any
	ClassDefinition ClassDefinition
}

func TestSetAllowDelete(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_isCreatableDeletable_never",
			Input: setAllowDeleteInput{
				MetaFileContent: map[string]any{"isCreatableDeletable": "never"},
				ClassDefinition: ClassDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_isCreatableDeletable_always",
			Input: setAllowDeleteInput{
				MetaFileContent: map[string]any{"isCreatableDeletable": "always"},
				ClassDefinition: ClassDefinition{},
			},
			Expected: true,
		},
		{
			Name: "test_definition_allow_delete_never",
			Input: setAllowDeleteInput{
				MetaFileContent: map[string]any{},
				ClassDefinition: ClassDefinition{AllowDelete: "never"},
			},
			Expected: false,
		},
		{
			Name: "test_definition_allow_delete_always",
			Input: setAllowDeleteInput{
				MetaFileContent: map[string]any{},
				ClassDefinition: ClassDefinition{AllowDelete: "always"},
			},
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setAllowDeleteInput)
			class := Class{Name: testClassName("fvPeeringP")}
			class.MetaFileContent = input.MetaFileContent
			class.ClassDefinition = input.ClassDefinition
			class.AllowDelete = false

			class.setAllowDelete()

			assert.Equal(t, testCase.Expected, class.AllowDelete, test.MessageEqual(testCase.Expected, class.AllowDelete, testCase.Name))
		})
	}
}

type shouldIncludeChildInput struct {
	RN                          string
	ClassName                   string
	ExcludeChildrenFromClassDef []string
	AlwaysIncludeFromGlobalDef  []string
}

func TestShouldIncludeChild(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_exclude_when_in_excludeChildrenFromClassDef",
			Input: shouldIncludeChildInput{
				RN:                          "ctx-",
				ClassName:                   "fvCtx",
				ExcludeChildrenFromClassDef: []string{"fvCtx"},
				AlwaysIncludeFromGlobalDef:  []string{},
			},
			Expected: false,
		},
		{
			Name: "test_excludeChildrenFromClassDef_takes_precedence_over_alwaysIncludeFromGlobalDef",
			Input: shouldIncludeChildInput{
				RN:                          "ctx-",
				ClassName:                   "fvCtx",
				ExcludeChildrenFromClassDef: []string{"fvCtx"},
				AlwaysIncludeFromGlobalDef:  []string{"fvCtx"},
			},
			Expected: false,
		},
		{
			Name: "test_include_when_in_alwaysIncludeFromGlobalDef",
			Input: shouldIncludeChildInput{
				RN:                          "annotationKey-",
				ClassName:                   "tagAnnotation",
				ExcludeChildrenFromClassDef: []string{},
				AlwaysIncludeFromGlobalDef:  []string{"tagAnnotation", "tagTag"},
			},
			Expected: true,
		},
		{
			Name: "test_alwaysIncludeFromGlobalDef_overrides_default_exclude",
			Input: shouldIncludeChildInput{
				RN:                          "tagKey-",
				ClassName:                   "tagTag",
				ExcludeChildrenFromClassDef: []string{},
				AlwaysIncludeFromGlobalDef:  []string{"tagAnnotation", "tagTag"},
			},
			Expected: true,
		},
		{
			Name: "test_include_when_rn_starts_with_rs",
			Input: shouldIncludeChildInput{
				RN:                          "rsBDToOut",
				ClassName:                   "fvRsBDToOut",
				ExcludeChildrenFromClassDef: []string{},
				AlwaysIncludeFromGlobalDef:  []string{},
			},
			Expected: true,
		},
		{
			Name: "test_include_when_rn_starts_with_rs_even_if_ends_with_dash",
			Input: shouldIncludeChildInput{
				RN:                          "rsCtx-",
				ClassName:                   "fvRsCtx",
				ExcludeChildrenFromClassDef: []string{},
				AlwaysIncludeFromGlobalDef:  []string{},
			},
			Expected: true,
		},
		{
			Name: "test_exclude_when_rn_ends_with_dash",
			Input: shouldIncludeChildInput{
				RN:                          "ctx-",
				ClassName:                   "fvCtx",
				ExcludeChildrenFromClassDef: []string{},
				AlwaysIncludeFromGlobalDef:  []string{},
			},
			Expected: false,
		},
		{
			Name: "test_include_when_rn_does_not_end_with_dash",
			Input: shouldIncludeChildInput{
				RN:                          "subnet",
				ClassName:                   "fvSubnet",
				ExcludeChildrenFromClassDef: []string{},
				AlwaysIncludeFromGlobalDef:  []string{},
			},
			Expected: true,
		},
		{
			Name: "test_include_when_rn_does_not_end_with_dash_and_nil_lists",
			Input: shouldIncludeChildInput{
				RN:                          "ap",
				ClassName:                   "fvAp",
				ExcludeChildrenFromClassDef: nil,
				AlwaysIncludeFromGlobalDef:  nil,
			},
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(shouldIncludeChildInput)

			result := shouldIncludeChild(input.RN, input.ClassName, input.ExcludeChildrenFromClassDef, input.AlwaysIncludeFromGlobalDef)

			assert.Equal(t, testCase.Expected, result, test.MessageEqual(testCase.Expected, result, testCase.Name))
		})
	}
}

type setChildrenInput struct {
	IncludeChildren      []string
	ExcludeChildren      []string
	AlwaysIncludeAsChild []string
	RnMap                map[string]any
}

func TestSetChildren(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_empty_rnMap_returns_only_includeChildren",
			Input: setChildrenInput{
				IncludeChildren:      []string{"fvSubnet"},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{},
				RnMap:                map[string]any{},
			},
			Expected: []string{"fvSubnet"},
		},
		{
			Name: "test_includes_rs_prefixed_classes_from_rnMap",
			Input: setChildrenInput{
				IncludeChildren:      []string{},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{},
				RnMap: map[string]any{
					"rsBDToOut": "fv:RsBDToOut",
					"rsCtx":     "fv:RsCtx",
				},
			},
			Expected: []string{"fvRsBDToOut", "fvRsCtx"},
		},
		{
			Name: "test_excludes_classes_with_rn_ending_in_dash",
			Input: setChildrenInput{
				IncludeChildren:      []string{},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{},
				RnMap: map[string]any{
					"ctx-":      "fv:Ctx",
					"rsBDToOut": "fv:RsBDToOut",
				},
			},
			Expected: []string{"fvRsBDToOut"},
		},
		{
			Name: "test_removes_colon_from_class_names",
			Input: setChildrenInput{
				IncludeChildren:      []string{},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{"tagAnnotation", "tagTag"},
				RnMap: map[string]any{
					"annotationKey-": "tag:Annotation",
				},
			},
			Expected: []string{"tagAnnotation"},
		},
		{
			Name: "test_respects_excludeChildren",
			Input: setChildrenInput{
				IncludeChildren:      []string{},
				ExcludeChildren:      []string{"fvRsBDToOut"},
				AlwaysIncludeAsChild: []string{},
				RnMap: map[string]any{
					"rsBDToOut": "fv:RsBDToOut",
					"rsCtx":     "fv:RsCtx",
				},
			},
			Expected: []string{"fvRsCtx"},
		},
		{
			Name: "test_combines_includeChildren_with_rnMap_results",
			Input: setChildrenInput{
				IncludeChildren:      []string{"fvSubnet"},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{},
				RnMap: map[string]any{
					"rsBDToOut": "fv:RsBDToOut",
				},
			},
			Expected: []string{"fvRsBDToOut", "fvSubnet"},
		},
		{
			Name: "test_removes_duplicates",
			Input: setChildrenInput{
				IncludeChildren:      []string{"fvRsBDToOut"},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{},
				RnMap: map[string]any{
					"rsBDToOut": "fv:RsBDToOut",
				},
			},
			Expected: []string{"fvRsBDToOut"},
		},
		{
			Name: "test_sorts_children_alphabetically",
			Input: setChildrenInput{
				IncludeChildren:      []string{"fvZone", "fvAp"},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{},
				RnMap: map[string]any{
					"rsBDToOut": "fv:RsBDToOut",
				},
			},
			Expected: []string{"fvAp", "fvRsBDToOut", "fvZone"},
		},
		{
			Name: "test_nil_rnMap_returns_only_includeChildren",
			Input: setChildrenInput{
				IncludeChildren:      []string{"fvSubnet"},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{},
				RnMap:                nil,
			},
			Expected: []string{"fvSubnet"},
		},
		{
			Name: "test_empty_inputs_returns_empty_children",
			Input: setChildrenInput{
				IncludeChildren:      []string{},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{},
				RnMap:                map[string]any{},
			},
			Expected: []string{},
		},
		{
			Name: "test_includeChildren_takes_precedence_over_excludeChildren",
			Input: setChildrenInput{
				IncludeChildren:      []string{"fvSubnet"},
				ExcludeChildren:      []string{"fvSubnet"},
				AlwaysIncludeAsChild: []string{},
				RnMap:                map[string]any{},
			},
			Expected: []string{"fvSubnet"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setChildrenInput)
			expected := testCase.Expected.([]string)

			class := &Class{
				Name: testClassName("testClass"),
				ClassDefinition: ClassDefinition{
					IncludeChildren: input.IncludeChildren,
					ExcludeChildren: input.ExcludeChildren,
				},
				MetaFileContent: map[string]any{
					"rnMap": input.RnMap,
				},
			}

			ds := &DataStore{
				GlobalMetaDefinition: GlobalMetaDefinition{
					AlwaysIncludeAsChild: input.AlwaysIncludeAsChild,
				},
			}

			err := class.setChildren(ds)
			assert.NoError(t, err, test.MessageUnexpectedError(err))

			if len(expected) == 0 {
				assert.Empty(t, class.Children, test.MessageEqual(expected, class.Children, testCase.Name))
			} else {
				assert.Equal(t, expected, classNamesToStrings(class.Children), test.MessageEqual(expected, class.Children, testCase.Name))
			}
		})
	}
}

func TestSetChildrenWarnsWhenClassInBothIncludeAndExclude(t *testing.T) {
	test.InitializeTest(t)

	// Capture log output using a buffer.
	var logBuffer bytes.Buffer
	genLogger.SetOutputForTesting(&logBuffer)
	genLogger.SetLogLevel("WARN")

	// Restore original log output after test.
	defer func() {
		genLogger.SetOutputForTesting(os.Stdout)
	}()

	class := &Class{
		Name: testClassName("testClass"),
		ClassDefinition: ClassDefinition{
			IncludeChildren: []string{"fvSubnet"},
			ExcludeChildren: []string{"fvSubnet"},
		},
		MetaFileContent: map[string]any{
			"rnMap": map[string]any{},
		},
	}

	ds := &DataStore{
		GlobalMetaDefinition: GlobalMetaDefinition{
			AlwaysIncludeAsChild: []string{},
		},
	}

	err := class.setChildren(ds)
	assert.NoError(t, err, test.MessageUnexpectedError(err))

	// Verify the warning was logged.
	logOutput := logBuffer.String()
	expectedWarning := "WARN: Child class 'fvSubnet' is defined in both IncludeChildren and ExcludeChildren for class 'testClass'. IncludeChildren takes precedence."
	assert.Contains(t, logOutput, expectedWarning, test.MessageEqual(expectedWarning, logOutput, "warning log message"))
}

type setChildrenErrorExpected struct {
	ErrorMsg string
}

func TestSetChildrenErrors(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_error_multiple_colons_in_rnMap",
			Input: setChildrenInput{
				IncludeChildren:      []string{},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{},
				RnMap: map[string]any{
					"rsBDToOut": "fv:Rs:BDToOut",
				},
			},
			Expected: setChildrenErrorExpected{
				ErrorMsg: "invalid class name 'fv:Rs:BDToOut': multiple colons detected",
			},
		},
		{
			Name: "test_error_invalid_class_name_in_includeChildren",
			Input: setChildrenInput{
				IncludeChildren:      []string{"invalidclass"},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{},
				RnMap:                map[string]any{},
			},
			Expected: setChildrenErrorExpected{
				ErrorMsg: "failed to split class name 'invalidclass' for name space separation",
			},
		},
		{
			Name: "test_error_empty_class_name_in_includeChildren",
			Input: setChildrenInput{
				IncludeChildren:      []string{""},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{},
				RnMap:                map[string]any{},
			},
			Expected: setChildrenErrorExpected{
				ErrorMsg: "failed to split class name '' for name space separation",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setChildrenInput)
			expected := testCase.Expected.(setChildrenErrorExpected)

			class := &Class{
				Name: testClassName("testClass"),
				ClassDefinition: ClassDefinition{
					IncludeChildren: input.IncludeChildren,
					ExcludeChildren: input.ExcludeChildren,
				},
				MetaFileContent: map[string]any{
					"rnMap": input.RnMap,
				},
			}

			ds := &DataStore{
				GlobalMetaDefinition: GlobalMetaDefinition{
					AlwaysIncludeAsChild: input.AlwaysIncludeAsChild,
				},
			}

			err := class.setChildren(ds)
			assert.EqualError(t, err, expected.ErrorMsg)
		})
	}
}

type setParentsInput struct {
	IncludeParents []string
	ExcludeParents []string
	ContainedBy    map[string]any
}

func TestSetParents(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_empty_containedBy_returns_only_includeParents",
			Input: setParentsInput{
				IncludeParents: []string{"fvTenant"},
				ExcludeParents: []string{},
				ContainedBy:    map[string]any{},
			},
			Expected: []string{"fvTenant"},
		},
		{
			Name: "test_includes_all_classes_from_containedBy",
			Input: setParentsInput{
				IncludeParents: []string{},
				ExcludeParents: []string{},
				ContainedBy: map[string]any{
					"fv:AEPg": "",
					"fv:BD":   "",
				},
			},
			Expected: []string{"fvAEPg", "fvBD"},
		},
		{
			Name: "test_removes_colon_from_class_names",
			Input: setParentsInput{
				IncludeParents: []string{},
				ExcludeParents: []string{},
				ContainedBy: map[string]any{
					"fv:Tenant": "",
				},
			},
			Expected: []string{"fvTenant"},
		},
		{
			Name: "test_respects_excludeParents",
			Input: setParentsInput{
				IncludeParents: []string{},
				ExcludeParents: []string{"fvAEPg"},
				ContainedBy: map[string]any{
					"fv:AEPg": "",
					"fv:BD":   "",
				},
			},
			Expected: []string{"fvBD"},
		},
		{
			Name: "test_combines_includeParents_with_containedBy_results",
			Input: setParentsInput{
				IncludeParents: []string{"fvTenant"},
				ExcludeParents: []string{},
				ContainedBy: map[string]any{
					"fv:AEPg": "",
				},
			},
			Expected: []string{"fvAEPg", "fvTenant"},
		},
		{
			Name: "test_removes_duplicates",
			Input: setParentsInput{
				IncludeParents: []string{"fvAEPg"},
				ExcludeParents: []string{},
				ContainedBy: map[string]any{
					"fv:AEPg": "",
				},
			},
			Expected: []string{"fvAEPg"},
		},
		{
			Name: "test_sorts_parents_alphabetically",
			Input: setParentsInput{
				IncludeParents: []string{"fvCtx", "fvAp"},
				ExcludeParents: []string{},
				ContainedBy: map[string]any{
					"fv:BD": "",
				},
			},
			Expected: []string{"fvAp", "fvBD", "fvCtx"},
		},
		{
			Name: "test_nil_containedBy_returns_only_includeParents",
			Input: setParentsInput{
				IncludeParents: []string{"fvTenant"},
				ExcludeParents: []string{},
				ContainedBy:    nil,
			},
			Expected: []string{"fvTenant"},
		},
		{
			Name: "test_empty_inputs_returns_empty_parents",
			Input: setParentsInput{
				IncludeParents: []string{},
				ExcludeParents: []string{},
				ContainedBy:    map[string]any{},
			},
			Expected: []string{},
		},
		{
			Name: "test_includeParents_takes_precedence_over_excludeParents",
			Input: setParentsInput{
				IncludeParents: []string{"fvTenant"},
				ExcludeParents: []string{"fvTenant"},
				ContainedBy:    map[string]any{},
			},
			Expected: []string{"fvTenant"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setParentsInput)
			expected := testCase.Expected.([]string)

			class := &Class{
				Name: testClassName("testClass"),
				ClassDefinition: ClassDefinition{
					IncludeParents: input.IncludeParents,
					ExcludeParents: input.ExcludeParents,
				},
				MetaFileContent: map[string]any{
					"containedBy": input.ContainedBy,
				},
			}

			err := class.setParents(&DataStore{})
			assert.NoError(t, err, test.MessageUnexpectedError(err))

			if len(expected) == 0 {
				assert.Empty(t, class.Parents, test.MessageEqual(expected, class.Parents, testCase.Name))
			} else {
				assert.Equal(t, expected, classNamesToStrings(class.Parents), test.MessageEqual(expected, class.Parents, testCase.Name))
			}
		})
	}
}

func TestSetParentsWarnsWhenClassInBothIncludeAndExclude(t *testing.T) {
	test.InitializeTest(t)

	// Capture log output using a buffer.
	var logBuffer bytes.Buffer
	genLogger.SetOutputForTesting(&logBuffer)
	genLogger.SetLogLevel("WARN")

	// Restore original log output after test.
	defer func() {
		genLogger.SetOutputForTesting(os.Stdout)
	}()

	class := &Class{
		Name: testClassName("testClass"),
		ClassDefinition: ClassDefinition{
			IncludeParents: []string{"fvTenant"},
			ExcludeParents: []string{"fvTenant"},
		},
		MetaFileContent: map[string]any{
			"containedBy": map[string]any{},
		},
	}

	err := class.setParents(&DataStore{})
	assert.NoError(t, err, test.MessageUnexpectedError(err))

	// Verify the warning was logged.
	logOutput := logBuffer.String()
	expectedWarning := "WARN: Parent class 'fvTenant' is defined in both IncludeParents and ExcludeParents for class 'testClass'. IncludeParents takes precedence."
	assert.Contains(t, logOutput, expectedWarning, test.MessageEqual(expectedWarning, logOutput, "warning log message"))
}

func TestSetParentsGlobalExcludeParents(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_global_exclude_filters_containedBy_class",
			Input: setParentsInput{
				IncludeParents: []string{},
				ExcludeParents: []string{},
				ContainedBy: map[string]any{
					"pol:Uni": "",
					"fv:AEPg": "",
				},
			},
			Expected: []string{"fvAEPg"},
		},
		{
			Name: "test_global_exclude_filters_multiple_classes",
			Input: setParentsInput{
				IncludeParents: []string{},
				ExcludeParents: []string{},
				ContainedBy: map[string]any{
					"pol:Uni":     "",
					"fabric:Inst": "",
					"fv:Tenant":   "",
				},
			},
			Expected: []string{"fvTenant"},
		},
		{
			Name: "test_global_exclude_does_not_affect_includeParents",
			Input: setParentsInput{
				IncludeParents: []string{"polUni"},
				ExcludeParents: []string{},
				ContainedBy:    map[string]any{},
			},
			Expected: []string{"polUni"},
		},
		{
			Name: "test_global_and_class_exclude_combined",
			Input: setParentsInput{
				IncludeParents: []string{},
				ExcludeParents: []string{"fvAEPg"},
				ContainedBy: map[string]any{
					"pol:Uni": "",
					"fv:AEPg": "",
					"fv:BD":   "",
				},
			},
			Expected: []string{"fvBD"},
		},
		{
			Name: "test_global_exclude_all_containedBy_returns_empty",
			Input: setParentsInput{
				IncludeParents: []string{},
				ExcludeParents: []string{},
				ContainedBy: map[string]any{
					"pol:Uni":     "",
					"fabric:Inst": "",
				},
			},
			Expected: []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setParentsInput)
			expected := testCase.Expected.([]string)

			class := &Class{
				Name: testClassName("testClass"),
				ClassDefinition: ClassDefinition{
					IncludeParents: input.IncludeParents,
					ExcludeParents: input.ExcludeParents,
				},
				MetaFileContent: map[string]any{
					"containedBy": input.ContainedBy,
				},
			}

			ds := &DataStore{
				GlobalMetaDefinition: GlobalMetaDefinition{
					ExcludeParents: []string{"polUni", "fabricInst"},
				},
			}

			err := class.setParents(ds)
			assert.NoError(t, err, test.MessageUnexpectedError(err))

			if len(expected) == 0 {
				assert.Empty(t, class.Parents, test.MessageEqual(expected, class.Parents, testCase.Name))
			} else {
				assert.Equal(t, expected, classNamesToStrings(class.Parents), test.MessageEqual(expected, class.Parents, testCase.Name))
			}
		})
	}
}

func TestSetParentsExcludeParents(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_empty_exclude_lists_preserve_full_containedBy",
			Input: setParentsInput{
				IncludeParents: []string{},
				ExcludeParents: []string{},
				ContainedBy: map[string]any{
					"fv:Tenant": "",
					"fv:AEPg":   "",
				},
			},
			Expected: []string{"fvAEPg", "fvTenant"},
		},
		{
			Name: "test_excluded_entry_not_in_containedBy_is_noop",
			Input: setParentsInput{
				IncludeParents: []string{},
				ExcludeParents: []string{"fvCtx"},
				ContainedBy: map[string]any{
					"fv:Tenant": "",
				},
			},
			Expected: []string{"fvTenant"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setParentsInput)
			expected := testCase.Expected.([]string)

			class := &Class{
				Name: testClassName("testClass"),
				ClassDefinition: ClassDefinition{
					IncludeParents: input.IncludeParents,
					ExcludeParents: input.ExcludeParents,
				},
				MetaFileContent: map[string]any{
					"containedBy": input.ContainedBy,
				},
			}

			err := class.setParents(&DataStore{})
			assert.NoError(t, err, test.MessageUnexpectedError(err))

			if len(expected) == 0 {
				assert.Empty(t, class.Parents, test.MessageEqual(expected, class.Parents, testCase.Name))
			} else {
				assert.Equal(t, expected, classNamesToStrings(class.Parents), test.MessageEqual(expected, class.Parents, testCase.Name))
			}
		})
	}
}

type setParentsErrorExpected struct {
	ErrorMsg string
}

func TestSetParentsErrors(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_error_multiple_colons_in_containedBy",
			Input: setParentsInput{
				IncludeParents: []string{},
				ExcludeParents: []string{},
				ContainedBy: map[string]any{
					"fv:AE:Pg": "",
				},
			},
			Expected: setParentsErrorExpected{
				ErrorMsg: "invalid class name 'fv:AE:Pg': multiple colons detected",
			},
		},
		{
			Name: "test_error_invalid_class_name_in_includeParents",
			Input: setParentsInput{
				IncludeParents: []string{"invalidparent"},
				ExcludeParents: []string{},
				ContainedBy:    map[string]any{},
			},
			Expected: setParentsErrorExpected{
				ErrorMsg: "failed to split class name 'invalidparent' for name space separation",
			},
		},
		{
			Name: "test_error_empty_class_name_in_includeParents",
			Input: setParentsInput{
				IncludeParents: []string{""},
				ExcludeParents: []string{},
				ContainedBy:    map[string]any{},
			},
			Expected: setParentsErrorExpected{
				ErrorMsg: "failed to split class name '' for name space separation",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setParentsInput)
			expected := testCase.Expected.(setParentsErrorExpected)

			class := &Class{
				Name: testClassName("testClass"),
				ClassDefinition: ClassDefinition{
					IncludeParents: input.IncludeParents,
					ExcludeParents: input.ExcludeParents,
				},
				MetaFileContent: map[string]any{
					"containedBy": input.ContainedBy,
				},
			}

			err := class.setParents(&DataStore{})
			assert.EqualError(t, err, expected.ErrorMsg)
		})
	}
}

type sortAndConvertToClassNamesExpected struct {
	ClassNames []string
	Error      bool
	ErrorMsg   string
}

func TestSortAndConvertToClassNames(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_empty_slice_returns_empty_result",
			Input:    []string{},
			Expected: sortAndConvertToClassNamesExpected{ClassNames: []string{}},
		},
		{
			Name:     "test_single_valid_class_name",
			Input:    []string{"fvTenant"},
			Expected: sortAndConvertToClassNamesExpected{ClassNames: []string{"fvTenant"}},
		},
		{
			Name:     "test_multiple_class_names_sorted",
			Input:    []string{"fvCtx", "fvAp", "fvBD"},
			Expected: sortAndConvertToClassNamesExpected{ClassNames: []string{"fvAp", "fvBD", "fvCtx"}},
		},
		{
			Name:     "test_duplicates_removed",
			Input:    []string{"fvTenant", "fvBD", "fvTenant"},
			Expected: sortAndConvertToClassNamesExpected{ClassNames: []string{"fvBD", "fvTenant"}},
		},
		{
			Name:  "test_error_empty_string",
			Input: []string{""},
			Expected: sortAndConvertToClassNamesExpected{
				Error:    true,
				ErrorMsg: "failed to split class name '' for name space separation",
			},
		},
		{
			Name:  "test_error_invalid_class_name_no_uppercase",
			Input: []string{"invalidclass"},
			Expected: sortAndConvertToClassNamesExpected{
				Error:    true,
				ErrorMsg: "failed to split class name 'invalidclass' for name space separation",
			},
		},
		{
			Name:  "test_error_invalid_class_name_no_package",
			Input: []string{"Tenant"},
			Expected: sortAndConvertToClassNamesExpected{
				Error:    true,
				ErrorMsg: "failed to split class name 'Tenant' for name space separation",
			},
		},
		{
			Name:  "test_error_mixed_valid_and_invalid",
			Input: []string{"fvTenant", "invalidclass"},
			Expected: sortAndConvertToClassNamesExpected{
				Error:    true,
				ErrorMsg: "failed to split class name 'invalidclass' for name space separation",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.([]string)
			expected := testCase.Expected.(sortAndConvertToClassNamesExpected)

			result, err := sortAndConvertToClassNames(input)

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else {
				assert.Equal(t, expected.ClassNames, classNamesToStrings(result), test.MessageEqual(expected.ClassNames, result, testCase.Name))
			}
		})
	}
}

type setClassSupportedVersionsInput struct {
	MetaFileContent map[string]any
	ClassDefinition ClassDefinition
}

type setClassSupportedVersionsExpected struct {
	Raw      string
	String   string
	Error    bool
	ErrorMsg string
}

func TestSetClassSupportedVersions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_versions_from_meta_file",
			Input: setClassSupportedVersionsInput{
				MetaFileContent: map[string]any{"versions": "4.2(7f)-"},
				ClassDefinition: ClassDefinition{},
			},
			Expected: setClassSupportedVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name: "test_versions_from_class_definition_override",
			Input: setClassSupportedVersionsInput{
				MetaFileContent: map[string]any{"versions": "4.2(7f)-"},
				ClassDefinition: ClassDefinition{SupportedVersions: "5.0(1a)-"},
			},
			Expected: setClassSupportedVersionsExpected{
				Raw:    "5.0(1a)-",
				String: "5.0(1a) and later",
			},
		},
		{
			Name: "test_versions_from_class_definition_when_meta_empty",
			Input: setClassSupportedVersionsInput{
				MetaFileContent: map[string]any{"versions": ""},
				ClassDefinition: ClassDefinition{SupportedVersions: "5.0(1a)-"},
			},
			Expected: setClassSupportedVersionsExpected{
				Raw:    "5.0(1a)-",
				String: "5.0(1a) and later",
			},
		},
		{
			Name: "test_versions_from_class_definition_when_meta_nil",
			Input: setClassSupportedVersionsInput{
				MetaFileContent: map[string]any{},
				ClassDefinition: ClassDefinition{SupportedVersions: "5.0(1a)-"},
			},
			Expected: setClassSupportedVersionsExpected{
				Raw:    "5.0(1a)-",
				String: "5.0(1a) and later",
			},
		},
		{
			Name: "test_multiple_ranges",
			Input: setClassSupportedVersionsInput{
				MetaFileContent: map[string]any{"versions": "3.2(10e)-3.2(10g),4.2(7f)-"},
				ClassDefinition: ClassDefinition{},
			},
			Expected: setClassSupportedVersionsExpected{
				Raw:    "3.2(10e)-3.2(10g),4.2(7f)-",
				String: "3.2(10e) to 3.2(10g), 4.2(7f) and later",
			},
		},
		{
			Name: "test_multiple_ranges_sorted",
			Input: setClassSupportedVersionsInput{
				MetaFileContent: map[string]any{"versions": "5.2(1g)-,3.2(10e)-3.2(10g)"},
				ClassDefinition: ClassDefinition{},
			},
			Expected: setClassSupportedVersionsExpected{
				Raw:    "5.2(1g)-,3.2(10e)-3.2(10g)",
				String: "3.2(10e) to 3.2(10g), 5.2(1g) and later",
			},
		},
		{
			Name: "test_error_empty_versions",
			Input: setClassSupportedVersionsInput{
				MetaFileContent: map[string]any{"versions": ""},
				ClassDefinition: ClassDefinition{},
			},
			Expected: setClassSupportedVersionsExpected{
				Error:    true,
				ErrorMsg: "versions not specified for class 'fvTenant': add versions to the class definition file",
			},
		},
		{
			Name: "test_error_nil_versions",
			Input: setClassSupportedVersionsInput{
				MetaFileContent: map[string]any{},
				ClassDefinition: ClassDefinition{},
			},
			Expected: setClassSupportedVersionsExpected{
				Error:    true,
				ErrorMsg: "versions not specified for class 'fvTenant': add versions to the class definition file",
			},
		},
		{
			Name: "test_error_invalid_version",
			Input: setClassSupportedVersionsInput{
				MetaFileContent: map[string]any{"versions": "invalid"},
				ClassDefinition: ClassDefinition{},
			},
			Expected: setClassSupportedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse versions for class 'fvTenant': invalid version 'invalid': unknown",
			},
		},
		{
			Name: "test_error_invalid_version_in_range",
			Input: setClassSupportedVersionsInput{
				MetaFileContent: map[string]any{"versions": "4.2(7f)-,invalid"},
				ClassDefinition: ClassDefinition{},
			},
			Expected: setClassSupportedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse versions for class 'fvTenant': invalid version 'invalid': unknown",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setClassSupportedVersionsInput)
			expected := testCase.Expected.(setClassSupportedVersionsExpected)

			class := Class{
				Name:            testClassName("fvTenant"),
				ClassDefinition: input.ClassDefinition,
				MetaFileContent: input.MetaFileContent,
			}

			err := class.setSupportedVersions()

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, expected.Raw, class.SupportedVersions.Raw(), test.MessageEqual(expected.Raw, class.SupportedVersions.Raw(), testCase.Name))
				assert.Equal(t, expected.String, class.SupportedVersions.String(), test.MessageEqual(expected.String, class.SupportedVersions.String(), testCase.Name))
			}
		})
	}
}

type setClassDeprecatedVersionsInput struct {
	ClassDefinitionVersions string
	MetaDeprecatedSince     any
}

type setClassDeprecatedVersionsExpected struct {
	Raw      string
	String   string
	Nil      bool
	Error    bool
	ErrorMsg string
}

func TestSetDeprecatedVersions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_deprecated_versions_not_set",
			Input: setClassDeprecatedVersionsInput{
				ClassDefinitionVersions: "",
			},
			Expected: setClassDeprecatedVersionsExpected{
				Nil: true,
			},
		},
		{
			Name: "test_deprecated_versions_single_range",
			Input: setClassDeprecatedVersionsInput{
				ClassDefinitionVersions: "4.2(7f)-",
			},
			Expected: setClassDeprecatedVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name: "test_deprecated_versions_bounded_range",
			Input: setClassDeprecatedVersionsInput{
				ClassDefinitionVersions: "3.2(10e)-4.2(7f)",
			},
			Expected: setClassDeprecatedVersionsExpected{
				Raw:    "3.2(10e)-4.2(7f)",
				String: "3.2(10e) to 4.2(7f)",
			},
		},
		{
			Name: "test_deprecated_versions_multiple_ranges",
			Input: setClassDeprecatedVersionsInput{
				ClassDefinitionVersions: "3.2(10e)-3.2(10g),4.2(7f)-",
			},
			Expected: setClassDeprecatedVersionsExpected{
				Raw:    "3.2(10e)-3.2(10g),4.2(7f)-",
				String: "3.2(10e) to 3.2(10g), 4.2(7f) and later",
			},
		},
		{
			Name: "test_error_invalid_deprecated_version",
			Input: setClassDeprecatedVersionsInput{
				ClassDefinitionVersions: "invalid",
			},
			Expected: setClassDeprecatedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse deprecated versions for class 'fvTenant': invalid version 'invalid': unknown",
			},
		},
		{
			Name: "test_error_invalid_deprecated_version_in_range",
			Input: setClassDeprecatedVersionsInput{
				ClassDefinitionVersions: "4.2(7f)-,invalid",
			},
			Expected: setClassDeprecatedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse deprecated versions for class 'fvTenant': invalid version 'invalid': unknown",
			},
		},
		{
			Name: "test_meta_deprecated_since_single_range",
			Input: setClassDeprecatedVersionsInput{
				MetaDeprecatedSince: "5.2(1g)-",
			},
			Expected: setClassDeprecatedVersionsExpected{
				Raw:    "5.2(1g)-",
				String: "5.2(1g) and later",
			},
		},
		{
			Name: "test_meta_deprecated_since_wrong_type",
			Input: setClassDeprecatedVersionsInput{
				MetaDeprecatedSince: 123,
			},
			Expected: setClassDeprecatedVersionsExpected{
				Nil: true,
			},
		},
		{
			Name: "test_override_replaces_meta",
			Input: setClassDeprecatedVersionsInput{
				ClassDefinitionVersions: "4.2(7f)-",
				MetaDeprecatedSince:     "5.2(1g)-",
			},
			Expected: setClassDeprecatedVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name: "test_meta_parse_error",
			Input: setClassDeprecatedVersionsInput{
				MetaDeprecatedSince: "invalid",
			},
			Expected: setClassDeprecatedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse deprecated versions for class 'fvTenant': invalid version 'invalid': unknown",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setClassDeprecatedVersionsInput)
			expected := testCase.Expected.(setClassDeprecatedVersionsExpected)

			metaContent := map[string]any{}
			if input.MetaDeprecatedSince != nil {
				metaContent["deprecatedSince"] = input.MetaDeprecatedSince
			}

			class := Class{
				Name:            testClassName("fvTenant"),
				MetaFileContent: metaContent,
				ClassDefinition: ClassDefinition{
					DeprecatedVersions: input.ClassDefinitionVersions,
				},
			}

			err := class.setDeprecatedVersions()

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else if expected.Nil {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Nil(t, class.DeprecatedVersions, "expected DeprecatedVersions to be nil")
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, expected.Raw, class.DeprecatedVersions.Raw(), test.MessageEqual(expected.Raw, class.DeprecatedVersions.Raw(), testCase.Name))
				assert.Equal(t, expected.String, class.DeprecatedVersions.String(), test.MessageEqual(expected.String, class.DeprecatedVersions.String(), testCase.Name))
			}
		})
	}
}

type setIdentifiedByInput struct {
	ClassDefinitionIdentifiedBy []string
	MetaIdentifiedBy            any
}

func TestSetIdentifiedBy(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_identified_by_from_meta_file",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: nil,
				MetaIdentifiedBy:            []any{"name"},
			},
			Expected: []string{"name"},
		},
		{
			Name: "test_identified_by_from_meta_file_multiple",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: nil,
				MetaIdentifiedBy:            []any{"key", "value"},
			},
			Expected: []string{"key", "value"},
		},
		{
			Name: "test_identified_by_from_meta_file_sorted",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: nil,
				MetaIdentifiedBy:            []any{"value", "key"},
			},
			Expected: []string{"key", "value"},
		},
		{
			Name: "test_identified_by_from_meta_file_deduplicated",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: nil,
				MetaIdentifiedBy:            []any{"name", "name"},
			},
			Expected: []string{"name"},
		},
		{
			Name: "test_class_definition_overrides_meta_file",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: []string{"overrideKey"},
				MetaIdentifiedBy:            []any{"metaKey"},
			},
			Expected: []string{"overrideKey"},
		},
		{
			Name: "test_class_definition_overrides_when_meta_nil",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: []string{"overrideKey"},
				MetaIdentifiedBy:            nil,
			},
			Expected: []string{"overrideKey"},
		},
		{
			Name: "test_class_definition_multiple_sorted",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: []string{"zAttr", "aAttr"},
				MetaIdentifiedBy:            nil,
			},
			Expected: []string{"aAttr", "zAttr"},
		},
		{
			Name: "test_empty_class_definition_falls_back_to_meta",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: []string{},
				MetaIdentifiedBy:            []any{"metaKey"},
			},
			Expected: []string{"metaKey"},
		},
		{
			Name: "test_nil_class_definition_and_nil_meta_returns_empty",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: nil,
				MetaIdentifiedBy:            nil,
			},
			Expected: []string{},
		},
		{
			Name: "test_empty_class_definition_and_empty_meta_returns_empty",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: []string{},
				MetaIdentifiedBy:            []any{},
			},
			Expected: []string{},
		},
		{
			Name: "test_nil_class_definition_and_missing_meta_key_returns_empty",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: nil,
				MetaIdentifiedBy:            "not_a_slice",
			},
			Expected: []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setIdentifiedByInput)
			expected := testCase.Expected.([]string)

			class := Class{
				Name: testClassName("fvTenant"),
				ClassDefinition: ClassDefinition{
					IdentifiedBy: input.ClassDefinitionIdentifiedBy,
				},
				MetaFileContent: map[string]any{},
			}
			if input.MetaIdentifiedBy != nil {
				class.MetaFileContent["identifiedBy"] = input.MetaIdentifiedBy
			}

			class.setIdentifiedBy()

			if len(expected) == 0 {
				assert.Empty(t, class.IdentifiedBy, test.MessageEqual(expected, class.IdentifiedBy, testCase.Name))
			} else {
				assert.Equal(t, expected, class.IdentifiedBy, test.MessageEqual(expected, class.IdentifiedBy, testCase.Name))
			}
		})
	}
}

type setIsSingleNestedWhenDefinedAsChildInput struct {
	ClassDefinitionIsSingleNestedWhenDefinedAsChild bool
	IdentifiedBy                                    []string
}

func TestSetIsSingleNestedWhenDefinedAsChild(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_no_identifiers_returns_true",
			Input: setIsSingleNestedWhenDefinedAsChildInput{
				ClassDefinitionIsSingleNestedWhenDefinedAsChild: false,
				IdentifiedBy: []string{},
			},
			Expected: true,
		},
		{
			Name: "test_nil_identifiers_returns_true",
			Input: setIsSingleNestedWhenDefinedAsChildInput{
				ClassDefinitionIsSingleNestedWhenDefinedAsChild: false,
				IdentifiedBy: nil,
			},
			Expected: true,
		},
		{
			Name: "test_with_identifiers_returns_false",
			Input: setIsSingleNestedWhenDefinedAsChildInput{
				ClassDefinitionIsSingleNestedWhenDefinedAsChild: false,
				IdentifiedBy: []string{"name"},
			},
			Expected: false,
		},
		{
			Name: "test_class_definition_override_true",
			Input: setIsSingleNestedWhenDefinedAsChildInput{
				ClassDefinitionIsSingleNestedWhenDefinedAsChild: true,
				IdentifiedBy: []string{"name"},
			},
			Expected: true,
		},
		{
			Name: "test_class_definition_override_true_no_identifiers",
			Input: setIsSingleNestedWhenDefinedAsChildInput{
				ClassDefinitionIsSingleNestedWhenDefinedAsChild: true,
				IdentifiedBy: []string{},
			},
			Expected: true,
		},
		{
			Name: "test_multiple_identifiers_returns_false",
			Input: setIsSingleNestedWhenDefinedAsChildInput{
				ClassDefinitionIsSingleNestedWhenDefinedAsChild: false,
				IdentifiedBy: []string{"key", "value"},
			},
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setIsSingleNestedWhenDefinedAsChildInput)

			class := Class{
				Name: testClassName("fvTenant"),
				ClassDefinition: ClassDefinition{
					IsSingleNestedWhenDefinedAsChild: input.ClassDefinitionIsSingleNestedWhenDefinedAsChild,
				},
				IdentifiedBy: input.IdentifiedBy,
			}

			class.setIsSingleNestedWhenDefinedAsChild()

			assert.Equal(t, testCase.Expected, class.IsSingleNestedWhenDefinedAsChild, test.MessageEqual(testCase.Expected, class.IsSingleNestedWhenDefinedAsChild, testCase.Name))
		})
	}
}

type setClassDeprecatedInput struct {
	ClassDefinitionDeprecated bool
	MetaIsDeprecated          any
}

func TestSetDeprecated(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_meta_missing_no_override",
			Input:    setClassDeprecatedInput{},
			Expected: false,
		},
		{
			Name:     "test_meta_false_no_override",
			Input:    setClassDeprecatedInput{MetaIsDeprecated: false},
			Expected: false,
		},
		{
			Name:     "test_meta_true_no_override",
			Input:    setClassDeprecatedInput{MetaIsDeprecated: true},
			Expected: true,
		},
		{
			Name:     "test_meta_wrong_type",
			Input:    setClassDeprecatedInput{MetaIsDeprecated: "yes"},
			Expected: false,
		},
		{
			Name:     "test_override_true_meta_missing",
			Input:    setClassDeprecatedInput{ClassDefinitionDeprecated: true},
			Expected: true,
		},
		{
			Name:     "test_override_true_meta_false",
			Input:    setClassDeprecatedInput{ClassDefinitionDeprecated: true, MetaIsDeprecated: false},
			Expected: true,
		},
		{
			Name:     "test_override_false_meta_true",
			Input:    setClassDeprecatedInput{ClassDefinitionDeprecated: false, MetaIsDeprecated: true},
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setClassDeprecatedInput)

			metaContent := map[string]any{}
			if input.MetaIsDeprecated != nil {
				metaContent["isDeprecated"] = input.MetaIsDeprecated
			}

			class := Class{
				Name:            testClassName("fvTenant"),
				MetaFileContent: metaContent,
				ClassDefinition: ClassDefinition{
					Deprecated: input.ClassDefinitionDeprecated,
				},
			}

			class.setDeprecated()

			assert.Equal(t, testCase.Expected, class.Deprecated, test.MessageEqual(testCase.Expected, class.Deprecated, testCase.Name))
		})
	}
}

type setClassHiddenInput struct {
	ClassDefinitionHidden bool
	MetaIsHidden          any
}

func TestSetHidden(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_meta_missing_no_override",
			Input:    setClassHiddenInput{},
			Expected: false,
		},
		{
			Name:     "test_meta_false_no_override",
			Input:    setClassHiddenInput{MetaIsHidden: false},
			Expected: false,
		},
		{
			Name:     "test_meta_true_no_override",
			Input:    setClassHiddenInput{MetaIsHidden: true},
			Expected: true,
		},
		{
			Name:     "test_meta_wrong_type",
			Input:    setClassHiddenInput{MetaIsHidden: "yes"},
			Expected: false,
		},
		{
			Name:     "test_override_true_meta_missing",
			Input:    setClassHiddenInput{ClassDefinitionHidden: true},
			Expected: true,
		},
		{
			Name:     "test_override_true_meta_false",
			Input:    setClassHiddenInput{ClassDefinitionHidden: true, MetaIsHidden: false},
			Expected: true,
		},
		{
			Name:     "test_override_false_meta_true",
			Input:    setClassHiddenInput{ClassDefinitionHidden: false, MetaIsHidden: true},
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setClassHiddenInput)

			metaContent := map[string]any{}
			if input.MetaIsHidden != nil {
				metaContent["isHidden"] = input.MetaIsHidden
			}

			class := Class{
				Name:            testClassName("fvTenant"),
				MetaFileContent: metaContent,
				ClassDefinition: ClassDefinition{
					Hidden: input.ClassDefinitionHidden,
				},
			}

			class.setHidden()

			assert.Equal(t, testCase.Expected, class.Hidden, test.MessageEqual(testCase.Expected, class.Hidden, testCase.Name))
		})
	}
}

type setClassHiddenVersionsInput struct {
	ClassDefinitionVersions string
	MetaHiddenSince         any
}

type setClassHiddenVersionsExpected struct {
	Raw      string
	String   string
	Nil      bool
	Error    bool
	ErrorMsg string
}

func TestSetHiddenVersions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:  "test_hidden_versions_not_set",
			Input: setClassHiddenVersionsInput{},
			Expected: setClassHiddenVersionsExpected{
				Nil: true,
			},
		},
		{
			Name: "test_definition_single_range",
			Input: setClassHiddenVersionsInput{
				ClassDefinitionVersions: "4.2(7f)-",
			},
			Expected: setClassHiddenVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name: "test_meta_hidden_since_single_range",
			Input: setClassHiddenVersionsInput{
				MetaHiddenSince: "5.2(1g)-",
			},
			Expected: setClassHiddenVersionsExpected{
				Raw:    "5.2(1g)-",
				String: "5.2(1g) and later",
			},
		},
		{
			Name: "test_meta_hidden_since_wrong_type",
			Input: setClassHiddenVersionsInput{
				MetaHiddenSince: 123,
			},
			Expected: setClassHiddenVersionsExpected{
				Nil: true,
			},
		},
		{
			Name: "test_override_replaces_meta",
			Input: setClassHiddenVersionsInput{
				ClassDefinitionVersions: "4.2(7f)-",
				MetaHiddenSince:         "5.2(1g)-",
			},
			Expected: setClassHiddenVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name: "test_definition_parse_error",
			Input: setClassHiddenVersionsInput{
				ClassDefinitionVersions: "invalid",
			},
			Expected: setClassHiddenVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse hidden versions for class 'fvTenant': invalid version 'invalid': unknown",
			},
		},
		{
			Name: "test_meta_parse_error",
			Input: setClassHiddenVersionsInput{
				MetaHiddenSince: "invalid",
			},
			Expected: setClassHiddenVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse hidden versions for class 'fvTenant': invalid version 'invalid': unknown",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setClassHiddenVersionsInput)
			expected := testCase.Expected.(setClassHiddenVersionsExpected)

			metaContent := map[string]any{}
			if input.MetaHiddenSince != nil {
				metaContent["hiddenSince"] = input.MetaHiddenSince
			}

			class := Class{
				Name:            testClassName("fvTenant"),
				MetaFileContent: metaContent,
				ClassDefinition: ClassDefinition{
					HiddenVersions: input.ClassDefinitionVersions,
				},
			}

			err := class.setHiddenVersions()

			if expected.Error {
				assert.EqualError(t, err, expected.ErrorMsg)
			} else if expected.Nil {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Nil(t, class.HiddenVersions, "expected HiddenVersions to be nil")
			} else {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, expected.Raw, class.HiddenVersions.Raw(), test.MessageEqual(expected.Raw, class.HiddenVersions.Raw(), testCase.Name))
				assert.Equal(t, expected.String, class.HiddenVersions.String(), test.MessageEqual(expected.String, class.HiddenVersions.String(), testCase.Name))
			}
		})
	}
}

type setPlatformTypeInput struct {
	ClassDefinitionPlatformType string
	PlatformFlavors             any
}

func TestSetPlatformType(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_apic_only_from_meta",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []any{"apic"},
			},
			Expected: Apic,
		},
		{
			Name: "test_cloud_only_from_meta",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []any{"capic"},
			},
			Expected: Cloud,
		},
		{
			Name: "test_both_from_meta",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []any{"apic", "capic"},
			},
			Expected: Both,
		},
		{
			Name: "test_both_from_meta_reverse_order",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []any{"capic", "apic"},
			},
			Expected: Both,
		},
		{
			Name: "test_definition_overrides_meta_apic",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "apic",
				PlatformFlavors:             []any{"capic"},
			},
			Expected: Apic,
		},
		{
			Name: "test_definition_overrides_meta_cloud",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "cloud",
				PlatformFlavors:             []any{"apic"},
			},
			Expected: Cloud,
		},
		{
			Name: "test_definition_overrides_meta_both",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "both",
				PlatformFlavors:             nil,
			},
			Expected: Both,
		},
		{
			Name:     "test_default_zero_value_is_apic",
			Input:    setPlatformTypeInput{},
			Expected: Apic,
		},
		{
			Name: "test_nil_platform_flavors_no_definition_defaults_to_apic",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             nil,
			},
			Expected: Apic,
		},
		{
			Name: "test_empty_platform_flavors_defaults_to_apic",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []any{},
			},
			Expected: Apic,
		},
		{
			Name: "test_unknown_flavor_defaults_to_apic",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []any{"unknown"},
			},
			Expected: Apic,
		},
		{
			Name: "test_unknown_flavor_with_apic",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []any{"apic", "unknown"},
			},
			Expected: Apic,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPlatformTypeInput)

			class := Class{
				Name: testClassName("fvTenant"),
				ClassDefinition: ClassDefinition{
					PlatformType: input.ClassDefinitionPlatformType,
				},
				MetaFileContent: map[string]any{},
			}
			if input.PlatformFlavors != nil {
				class.MetaFileContent["platformFlavors"] = input.PlatformFlavors
			}

			class.setPlatformType()

			assert.Equal(t, testCase.Expected, class.PlatformType, test.MessageEqual(testCase.Expected, class.PlatformType, testCase.Name))
		})
	}
}

func TestSetPlatformTypeWarnsOnUnknownFlavor(t *testing.T) {
	test.InitializeTest(t)

	// Capture log output using a buffer.
	var logBuffer bytes.Buffer
	genLogger.SetOutputForTesting(&logBuffer)
	genLogger.SetLogLevel("WARN")

	// Restore original log output after test.
	defer func() {
		genLogger.SetOutputForTesting(os.Stdout)
	}()

	class := Class{
		Name: testClassName("fvTenant"),
		MetaFileContent: map[string]any{
			"platformFlavors": []any{"unknownFlavor"},
		},
	}

	class.setPlatformType()

	// Verify the warning was logged.
	logOutput := logBuffer.String()
	expectedWarning := "WARN: Unknown platform flavor 'unknownFlavor' found for class 'fvTenant'."
	assert.Contains(t, logOutput, expectedWarning, test.MessageEqual(expectedWarning, logOutput, "warning log message"))
}

func TestSetRequiredAsChild(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_required_as_child_false",
			Input:    false,
			Expected: false,
		},
		{
			Name:     "test_required_as_child_true",
			Input:    true,
			Expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			class := Class{
				Name: testClassName("fvTenant"),
				ClassDefinition: ClassDefinition{
					RequiredAsChild: testCase.Input.(bool),
				},
			}

			class.setRequiredAsChild()

			assert.Equal(t, testCase.Expected, class.RequiredAsChild, test.MessageEqual(testCase.Expected, class.RequiredAsChild, testCase.Name))
		})
	}
}

type setRnFormatInput struct {
	MetaFileContent map[string]any
	ClassDefinition ClassDefinition
}

func TestSetRnFormat(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_rn_format_from_meta_file",
			Input: setRnFormatInput{
				MetaFileContent: map[string]any{"rnFormat": "tn-{name}"},
				ClassDefinition: ClassDefinition{},
			},
			Expected: "tn-{name}",
		},
		{
			Name: "test_rn_format_definition_override",
			Input: setRnFormatInput{
				MetaFileContent: map[string]any{"rnFormat": "tn-{name}"},
				ClassDefinition: ClassDefinition{RnFormat: "custom-{name}"},
			},
			Expected: "custom-{name}",
		},
		{
			Name: "test_rn_format_definition_override_without_meta",
			Input: setRnFormatInput{
				MetaFileContent: map[string]any{},
				ClassDefinition: ClassDefinition{RnFormat: "custom-{name}"},
			},
			Expected: "custom-{name}",
		},
		{
			Name: "test_rn_prepend_with_meta_format",
			Input: setRnFormatInput{
				MetaFileContent: map[string]any{"rnFormat": "sdifpol-{name}"},
				ClassDefinition: ClassDefinition{RnPrepend: "infra"},
			},
			Expected: "infra/sdifpol-{name}",
		},
		{
			Name: "test_rn_prepend_with_definition_override",
			Input: setRnFormatInput{
				MetaFileContent: map[string]any{"rnFormat": "tn-{name}"},
				ClassDefinition: ClassDefinition{RnFormat: "custom-{name}", RnPrepend: "bar"},
			},
			Expected: "bar/custom-{name}",
		},
		{
			Name: "test_rn_prepend_multi_segment",
			Input: setRnFormatInput{
				MetaFileContent: map[string]any{"rnFormat": "instP-{name}"},
				ClassDefinition: ClassDefinition{RnPrepend: "tn-mgmt/extmgmt-default"},
			},
			Expected: "tn-mgmt/extmgmt-default/instP-{name}",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setRnFormatInput)
			class := Class{Name: testClassName("fvTenant")}
			class.MetaFileContent = input.MetaFileContent
			class.ClassDefinition = input.ClassDefinition

			err := class.setRnFormat()

			assert.NoError(t, err, test.MessageUnexpectedError(err))
			assert.Equal(t, testCase.Expected, class.RnFormat, test.MessageEqual(testCase.Expected, class.RnFormat, testCase.Name))
		})
	}
}

func TestSetRnFormatErrors(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_error_no_rn_format_from_any_source",
			Input: setRnFormatInput{
				MetaFileContent: map[string]any{},
				ClassDefinition: ClassDefinition{},
			},
			Expected: "rnFormat not specified for class 'fvTenant': add rn_format to the class definition file",
		},
		{
			Name: "test_error_rn_prepend_only_without_rn_format",
			Input: setRnFormatInput{
				MetaFileContent: map[string]any{},
				ClassDefinition: ClassDefinition{RnPrepend: "infra"},
			},
			Expected: "rnFormat not specified for class 'fvTenant': add rn_format to the class definition file",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setRnFormatInput)
			class := Class{Name: testClassName("fvTenant")}
			class.MetaFileContent = input.MetaFileContent
			class.ClassDefinition = input.ClassDefinition

			err := class.setRnFormat()

			assert.EqualError(t, err, testCase.Expected.(string))
		})
	}
}

type setPropertiesInput struct {
	MetaFileContent      map[string]any
	ClassDefinition      ClassDefinition
	GlobalMetaDefinition GlobalMetaDefinition
}

type setPropertiesExpected struct {
	PropertiesAll      []string
	PropertiesRequired []string
	PropertiesOptional []string
	PropertiesReadOnly []string
	PropertyDetails    map[string]struct {
		AttributeName string
		Required      bool
		Optional      bool
		ReadOnly      bool
	}
}

func TestSetProperties(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_configurable_naming_property_is_required",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{
					"properties": map[string]any{
						"name": map[string]any{
							"isConfigurable": true,
							"isNaming":       true,
						},
					},
				},
				ClassDefinition: ClassDefinition{},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{"name"},
				PropertiesRequired: []string{"name"},
				PropertiesOptional: []string{},
				PropertiesReadOnly: []string{},
				PropertyDetails: map[string]struct {
					AttributeName string
					Required      bool
					Optional      bool
					ReadOnly      bool
				}{
					"name": {AttributeName: "name", Required: true},
				},
			},
		},
		{
			Name: "test_configurable_non_naming_property_is_optional",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{
					"properties": map[string]any{
						"nameAlias": map[string]any{
							"isConfigurable": true,
							"isNaming":       false,
						},
					},
				},
				ClassDefinition: ClassDefinition{},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{"nameAlias"},
				PropertiesRequired: []string{},
				PropertiesOptional: []string{"nameAlias"},
				PropertiesReadOnly: []string{},
				PropertyDetails: map[string]struct {
					AttributeName string
					Required      bool
					Optional      bool
					ReadOnly      bool
				}{
					"nameAlias": {AttributeName: "name_alias", Optional: true},
				},
			},
		},
		{
			Name: "test_non_configurable_excluded_by_default",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{
					"properties": map[string]any{
						"pcTag": map[string]any{
							"isConfigurable": false,
							"isNaming":       false,
						},
					},
				},
				ClassDefinition: ClassDefinition{},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{},
				PropertiesRequired: []string{},
				PropertiesOptional: []string{},
				PropertiesReadOnly: []string{},
			},
		},
		{
			Name: "test_non_configurable_with_read_only_restriction",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{
					"properties": map[string]any{
						"pcTag": map[string]any{
							"isConfigurable": false,
							"isNaming":       false,
						},
					},
				},
				ClassDefinition: ClassDefinition{
					Properties: map[string]PropertyDefinition{
						"pcTag": {Restriction: ReadOnly},
					},
				},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{"pcTag"},
				PropertiesRequired: []string{},
				PropertiesOptional: []string{},
				PropertiesReadOnly: []string{"pcTag"},
				PropertyDetails: map[string]struct {
					AttributeName string
					Required      bool
					Optional      bool
					ReadOnly      bool
				}{
					"pcTag": {AttributeName: "pc_tag", ReadOnly: true},
				},
			},
		},
		{
			Name: "test_configurable_with_exclude_restriction",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{
					"properties": map[string]any{
						"annotation": map[string]any{
							"isConfigurable": true,
							"isNaming":       false,
						},
					},
				},
				ClassDefinition: ClassDefinition{
					Properties: map[string]PropertyDefinition{
						"annotation": {Restriction: Exclude},
					},
				},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{},
				PropertiesRequired: []string{},
				PropertiesOptional: []string{},
				PropertiesReadOnly: []string{},
			},
		},
		{
			Name: "test_configurable_with_required_restriction_override",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{
					"properties": map[string]any{
						"value": map[string]any{
							"isConfigurable": true,
							"isNaming":       false,
						},
					},
				},
				ClassDefinition: ClassDefinition{
					Properties: map[string]PropertyDefinition{
						"value": {Restriction: Required},
					},
				},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{"value"},
				PropertiesRequired: []string{"value"},
				PropertiesOptional: []string{},
				PropertiesReadOnly: []string{},
				PropertyDetails: map[string]struct {
					AttributeName string
					Required      bool
					Optional      bool
					ReadOnly      bool
				}{
					"value": {AttributeName: "value", Required: true},
				},
			},
		},
		{
			Name: "test_attribute_name_override",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{
					"properties": map[string]any{
						"pcEnfPref": map[string]any{
							"isConfigurable": true,
							"isNaming":       false,
						},
					},
				},
				ClassDefinition: ClassDefinition{
					Properties: map[string]PropertyDefinition{
						"pcEnfPref": {AttributeName: "policy_control_enforcement"},
					},
				},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{"pcEnfPref"},
				PropertiesRequired: []string{},
				PropertiesOptional: []string{"pcEnfPref"},
				PropertiesReadOnly: []string{},
				PropertyDetails: map[string]struct {
					AttributeName string
					Required      bool
					Optional      bool
					ReadOnly      bool
				}{
					"pcEnfPref": {AttributeName: "policy_control_enforcement", Optional: true},
				},
			},
		},
		{
			Name: "test_global_attribute_name_override",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{
					"properties": map[string]any{
						"descr": map[string]any{
							"isConfigurable": true,
							"isNaming":       false,
						},
					},
				},
				ClassDefinition: ClassDefinition{},
				GlobalMetaDefinition: GlobalMetaDefinition{
					AttributeNameOverrides: map[string]string{
						"descr": "description",
					},
				},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{"descr"},
				PropertiesRequired: []string{},
				PropertiesOptional: []string{"descr"},
				PropertiesReadOnly: []string{},
				PropertyDetails: map[string]struct {
					AttributeName string
					Required      bool
					Optional      bool
					ReadOnly      bool
				}{
					"descr": {AttributeName: "description", Optional: true},
				},
			},
		},
		{
			Name: "test_multiple_properties_sorted_alphabetically",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{
					"properties": map[string]any{
						"nameAlias": map[string]any{
							"isConfigurable": true,
							"isNaming":       false,
						},
						"annotation": map[string]any{
							"isConfigurable": true,
							"isNaming":       false,
						},
						"name": map[string]any{
							"isConfigurable": true,
							"isNaming":       true,
						},
						"pcTag": map[string]any{
							"isConfigurable": false,
							"isNaming":       false,
						},
						"scope": map[string]any{
							"isConfigurable": false,
							"isNaming":       false,
						},
					},
				},
				ClassDefinition: ClassDefinition{
					Properties: map[string]PropertyDefinition{
						"scope": {Restriction: ReadOnly},
					},
				},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{"annotation", "name", "nameAlias", "scope"},
				PropertiesRequired: []string{"name"},
				PropertiesOptional: []string{"annotation", "nameAlias"},
				PropertiesReadOnly: []string{"scope"},
				PropertyDetails: map[string]struct {
					AttributeName string
					Required      bool
					Optional      bool
					ReadOnly      bool
				}{
					"annotation": {AttributeName: "annotation", Optional: true},
					"nameAlias":  {AttributeName: "name_alias", Optional: true},
					"name":       {AttributeName: "name", Required: true},
					"scope":      {AttributeName: "scope", ReadOnly: true},
				},
			},
		},
		{
			Name: "test_empty_properties",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{},
				ClassDefinition: ClassDefinition{},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{},
				PropertiesRequired: []string{},
				PropertiesOptional: []string{},
				PropertiesReadOnly: []string{},
			},
		},
		{
			Name: "test_no_properties_key_in_meta",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{"label": "test"},
				ClassDefinition: ClassDefinition{},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{},
				PropertiesRequired: []string{},
				PropertiesOptional: []string{},
				PropertiesReadOnly: []string{},
			},
		},
		{
			Name: "test_global_exclude_property",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{
					"properties": map[string]any{
						"userdom": map[string]any{
							"isConfigurable": true,
							"isNaming":       false,
						},
						"name": map[string]any{
							"isConfigurable": true,
							"isNaming":       true,
						},
					},
				},
				ClassDefinition: ClassDefinition{},
				GlobalMetaDefinition: GlobalMetaDefinition{
					ExcludeProperties: []string{"userdom"},
				},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{"name"},
				PropertiesRequired: []string{"name"},
				PropertiesOptional: []string{},
				PropertiesReadOnly: []string{},
				PropertyDetails: map[string]struct {
					AttributeName string
					Required      bool
					Optional      bool
					ReadOnly      bool
				}{
					"name": {AttributeName: "name", Required: true},
				},
			},
		},
		{
			Name: "test_global_exclude_overridden_by_class_definition",
			Input: setPropertiesInput{
				MetaFileContent: map[string]any{
					"properties": map[string]any{
						"userdom": map[string]any{
							"isConfigurable": true,
							"isNaming":       false,
						},
					},
				},
				ClassDefinition: ClassDefinition{
					Properties: map[string]PropertyDefinition{
						"userdom": {Restriction: Optional},
					},
				},
				GlobalMetaDefinition: GlobalMetaDefinition{
					ExcludeProperties: []string{"userdom"},
				},
			},
			Expected: setPropertiesExpected{
				PropertiesAll:      []string{"userdom"},
				PropertiesRequired: []string{},
				PropertiesOptional: []string{"userdom"},
				PropertiesReadOnly: []string{},
				PropertyDetails: map[string]struct {
					AttributeName string
					Required      bool
					Optional      bool
					ReadOnly      bool
				}{
					"userdom": {AttributeName: "userdom", Optional: true},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setPropertiesInput)
			expected := testCase.Expected.(setPropertiesExpected)

			class := &Class{
				Name:            testClassName("testClass"),
				ClassDefinition: input.ClassDefinition,
				MetaFileContent: input.MetaFileContent,
				Properties:      make(map[string]*Property),
			}

			ds := &DataStore{
				GlobalMetaDefinition: input.GlobalMetaDefinition,
			}

			err := class.setProperties(ds)
			assert.NoError(t, err, test.MessageUnexpectedError(err))

			test.AssertStringSlice(t, expected.PropertiesAll, class.PropertiesAll, testCase.Name)
			test.AssertStringSlice(t, expected.PropertiesRequired, class.PropertiesRequired, testCase.Name)
			test.AssertStringSlice(t, expected.PropertiesOptional, class.PropertiesOptional, testCase.Name)
			test.AssertStringSlice(t, expected.PropertiesReadOnly, class.PropertiesReadOnly, testCase.Name)

			// Check individual property details when specified.
			if expected.PropertyDetails != nil {
				for propName, details := range expected.PropertyDetails {
					prop, ok := class.Properties[propName]
					assert.True(t, ok, test.MessageContains(class.Properties, propName, testCase.Name))
					if ok {
						assert.Equal(t, details.AttributeName, prop.AttributeName, test.MessageEqual(details.AttributeName, prop.AttributeName, testCase.Name))
						assert.Equal(t, details.Required, prop.Required, test.MessageEqual(details.Required, prop.Required, testCase.Name))
						assert.Equal(t, details.Optional, prop.Optional, test.MessageEqual(details.Optional, prop.Optional, testCase.Name))
						assert.Equal(t, details.ReadOnly, prop.ReadOnly, test.MessageEqual(details.ReadOnly, prop.ReadOnly, testCase.Name))
					}
				}
			}
		})
	}
}

type setTestDependenciesInput struct {
	Classes     map[string]Class
	TargetClass string
}

type setTestDependenciesExpected struct {
	DependencyCount int
	Dependencies    []expectedDependency
}

type expectedDependency struct {
	Role            TestDependencyRoleEnum
	Reference       string
	ReferenceType   ReferenceTypeEnum
	Class           string
	NestedCount     int
	NestedRefs      []string
	ConfigOverrides map[string]string
}

func TestSetTestDependencies(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_auto_resolve_parents",
			Input: setTestDependenciesInput{
				TargetClass: "fvAp",
				Classes: map[string]Class{
					"fvTenant": {
						Name:         testClassName("fvTenant"),
						ResourceName: "tenant",
					},
					"fvAp": {
						Name:            testClassName("fvAp"),
						ResourceName:    "application_profile",
						Parents:         []*ClassName{testClassName("fvTenant")},
						ClassDefinition: ClassDefinition{},
					},
				},
			},
			Expected: setTestDependenciesExpected{
				DependencyCount: 2,
				Dependencies: []expectedDependency{
					{Role: Parent, Reference: "aci_tenant.test.id", ReferenceType: ResourceReference},
					{Role: Parent, Reference: "aci_tenant.test_2.id", ReferenceType: ResourceReference},
				},
			},
		},
		{
			Name: "test_auto_resolve_targets_single_target",
			Input: setTestDependenciesInput{
				TargetClass: "fvRsBd",
				Classes: map[string]Class{
					"fvBD": {
						Name:         testClassName("fvBD"),
						ResourceName: "bridge_domain",
					},
					"fvRsBd": {
						Name:         testClassName("fvRsBd"),
						ResourceName: "relation_to_bridge_domain",
						Relation: Relation{
							RelationalClass: true,
							ToClasses:       []*ClassName{testClassName("fvBD")},
						},
						ClassDefinition: ClassDefinition{},
					},
				},
			},
			Expected: setTestDependenciesExpected{
				DependencyCount: 2,
				Dependencies: []expectedDependency{
					{Role: Target, Reference: "aci_bridge_domain.test.id", ReferenceType: ResourceReference},
					{Role: Target, Reference: "aci_bridge_domain.test_2.id", ReferenceType: ResourceReference},
				},
			},
		},
		{
			Name: "test_auto_resolve_recursive",
			Input: setTestDependenciesInput{
				TargetClass: "fvAEPg",
				Classes: map[string]Class{
					"fvTenant": {
						Name:         testClassName("fvTenant"),
						ResourceName: "tenant",
					},
					"fvAp": {
						Name:         testClassName("fvAp"),
						ResourceName: "application_profile",
						Parents:      []*ClassName{testClassName("fvTenant")},
					},
					"fvAEPg": {
						Name:            testClassName("fvAEPg"),
						ResourceName:    "application_epg",
						Parents:         []*ClassName{testClassName("fvAp")},
						ClassDefinition: ClassDefinition{},
					},
				},
			},
			Expected: setTestDependenciesExpected{
				DependencyCount: 2,
				Dependencies: []expectedDependency{
					{Role: Parent, Reference: "aci_application_profile.test.id", ReferenceType: ResourceReference, NestedCount: 1, NestedRefs: []string{"aci_tenant.test.id"}},
					{Role: Parent, Reference: "aci_application_profile.test_2.id", ReferenceType: ResourceReference, NestedCount: 1, NestedRefs: []string{"aci_tenant.test.id"}},
				},
			},
		},
		{
			Name: "test_explicit_replace_auto_resolved",
			Input: setTestDependenciesInput{
				TargetClass: "testClass",
				Classes: map[string]Class{
					"fvTenant": {
						Name:         testClassName("fvTenant"),
						ResourceName: "tenant",
					},
					"testClass": {
						Name:         testClassName("testClass"),
						ResourceName: "test_resource",
						Parents:      []*ClassName{testClassName("fvTenant")},
						ClassDefinition: ClassDefinition{
							TestConfig: ClassTestConfigDefinition{
								ReplaceAutoResolved: true,
								Dependencies: []TestDependencyDefinition{
									{
										ClassName:     "fvTenant",
										Reference:     "aci_tenant.test.id",
										ReferenceType: ResourceReference,
										Role:          Parent,
									},
								},
							},
						},
					},
				},
			},
			Expected: setTestDependenciesExpected{
				DependencyCount: 1,
				Dependencies: []expectedDependency{
					{Role: Parent, Reference: "aci_tenant.test.id", ReferenceType: ResourceReference, Class: "fvTenant"},
				},
			},
		},
		{
			Name: "test_skip_root_level_parent",
			Input: setTestDependenciesInput{
				TargetClass: "fvTenant",
				Classes: map[string]Class{
					"fvTenant": {
						Name:            testClassName("fvTenant"),
						ResourceName:    "tenant",
						Parents:         []*ClassName{testClassName("polUni")},
						ClassDefinition: ClassDefinition{},
					},
				},
			},
			Expected: setTestDependenciesExpected{
				DependencyCount: 0,
			},
		},
		{
			Name: "test_multi_target_requires_explicit",
			Input: setTestDependenciesInput{
				TargetClass: "fvRsProv",
				Classes: map[string]Class{
					"vzBrCP": {
						Name:         testClassName("vzBrCP"),
						ResourceName: "contract",
					},
					"vzTaboo": {
						Name:         testClassName("vzTaboo"),
						ResourceName: "taboo_contract",
					},
					"fvRsProv": {
						Name:         testClassName("fvRsProv"),
						ResourceName: "relation_to_provided_contract",
						Relation: Relation{
							RelationalClass: true,
							ToClasses:       []*ClassName{testClassName("vzBrCP"), testClassName("vzTaboo")},
						},
						ClassDefinition: ClassDefinition{},
					},
				},
			},
			Expected: setTestDependenciesExpected{
				DependencyCount: 0,
			},
		},
		{
			Name: "test_config_override_placeholder_resolution",
			Input: setTestDependenciesInput{
				TargetClass: "fvRsBd",
				Classes: map[string]Class{
					"fvTenant": {
						Name:         testClassName("fvTenant"),
						ResourceName: "tenant",
					},
					"fvBD": {
						Name:         testClassName("fvBD"),
						ResourceName: "bridge_domain",
						Parents:      []*ClassName{testClassName("fvTenant")},
					},
					"fvRsBd": {
						Name:         testClassName("fvRsBd"),
						ResourceName: "relation_to_bridge_domain",
						ClassDefinition: ClassDefinition{
							TestConfig: ClassTestConfigDefinition{
								ReplaceAutoResolved: true,
								Dependencies: []TestDependencyDefinition{
									{
										ClassName:     "fvTenant",
										Reference:     "aci_tenant.test.id",
										ReferenceType: ResourceReference,
										Role:          Parent,
									},
									{
										ClassName:       "fvBD",
										Reference:       "aci_bridge_domain.test.id",
										ReferenceType:   ResourceReference,
										Role:            Target,
										ConfigOverrides: map[string]string{"tenant_dn": "{{aci_tenant.test.id}}"},
										Dependencies: []TestDependencyDefinition{
											{
												ClassName:     "fvTenant",
												Reference:     "aci_tenant.test.id",
												ReferenceType: ResourceReference,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Expected: setTestDependenciesExpected{
				DependencyCount: 2,
				Dependencies: []expectedDependency{
					{Role: Parent, Reference: "aci_tenant.test.id", ReferenceType: ResourceReference},
					{Role: Target, Reference: "aci_bridge_domain.test.id", ReferenceType: ResourceReference, ConfigOverrides: map[string]string{"tenant_dn": "aci_tenant.test.id"}},
				},
			},
		},
		{
			Name: "test_explicit_additive",
			Input: setTestDependenciesInput{
				TargetClass: "fvRsBd",
				Classes: map[string]Class{
					"fvTenant": {
						Name:         testClassName("fvTenant"),
						ResourceName: "tenant",
					},
					"fvAp": {
						Name:         testClassName("fvAp"),
						ResourceName: "application_profile",
						Parents:      []*ClassName{testClassName("fvTenant")},
					},
					"fvBD": {
						Name:         testClassName("fvBD"),
						ResourceName: "bridge_domain",
						Parents:      []*ClassName{testClassName("fvTenant")},
					},
					"fvRsBd": {
						Name:         testClassName("fvRsBd"),
						ResourceName: "relation_to_bridge_domain",
						Parents:      []*ClassName{testClassName("fvAp")},
						ClassDefinition: ClassDefinition{
							TestConfig: ClassTestConfigDefinition{
								Dependencies: []TestDependencyDefinition{
									{
										ClassName:     "fvBD",
										Reference:     "aci_bridge_domain.extra.id",
										ReferenceType: ResourceReference,
										Role:          Target,
									},
								},
							},
						},
					},
				},
			},
			Expected: setTestDependenciesExpected{
				DependencyCount: 3,
				Dependencies: []expectedDependency{
					// Explicit definitions are processed first.
					{Role: Target, Reference: "aci_bridge_domain.extra.id", ReferenceType: ResourceReference, Class: "fvBD"},
					// Auto-resolved remainder fills in parents.
					{Role: Parent, Reference: "aci_application_profile.test.id", ReferenceType: ResourceReference},
					{Role: Parent, Reference: "aci_application_profile.test_2.id", ReferenceType: ResourceReference},
				},
			},
		},
		{
			Name: "test_explicit_additive_dedup",
			Input: setTestDependenciesInput{
				TargetClass: "fvAp",
				Classes: map[string]Class{
					"fvTenant": {
						Name:         testClassName("fvTenant"),
						ResourceName: "tenant",
					},
					"fvAp": {
						Name:         testClassName("fvAp"),
						ResourceName: "application_profile",
						Parents:      []*ClassName{testClassName("fvTenant")},
						ClassDefinition: ClassDefinition{
							TestConfig: ClassTestConfigDefinition{
								Dependencies: []TestDependencyDefinition{
									{
										ClassName:     "fvTenant",
										Reference:     "aci_tenant.test.id",
										ReferenceType: ResourceReference,
										Role:          Parent,
									},
								},
							},
						},
					},
				},
			},
			Expected: setTestDependenciesExpected{
				// Explicit defines aci_tenant.test.id first; auto-resolve adds aci_tenant.test_2.id.
				DependencyCount: 2,
				Dependencies: []expectedDependency{
					{Role: Parent, Reference: "aci_tenant.test.id", ReferenceType: ResourceReference},
					{Role: Parent, Reference: "aci_tenant.test_2.id", ReferenceType: ResourceReference},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setTestDependenciesInput)
			expected := testCase.Expected.(setTestDependenciesExpected)

			ds := &DataStore{
				Classes:              input.Classes,
				GlobalMetaDefinition: GlobalMetaDefinition{},
				ctx:                  NewContext(),
			}

			class := ds.Classes[input.TargetClass]
			class.setTestDependencies(ds)

			assert.Len(t, class.TestDependencies, expected.DependencyCount, test.MessageEqual(expected.DependencyCount, len(class.TestDependencies), testCase.Name))

			for i, expectedDep := range expected.Dependencies {
				if i >= len(class.TestDependencies) {
					break
				}
				actual := class.TestDependencies[i]
				assert.Equal(t, expectedDep.Role, actual.Role, test.MessageEqual(expectedDep.Role, actual.Role, testCase.Name))
				assert.Equal(t, expectedDep.Reference, actual.Reference, test.MessageEqual(expectedDep.Reference, actual.Reference, testCase.Name))
				assert.Equal(t, expectedDep.ReferenceType, actual.ReferenceType, test.MessageEqual(expectedDep.ReferenceType, actual.ReferenceType, testCase.Name))

				if expectedDep.Class != "" {
					assert.Equal(t, expectedDep.Class, actual.Class.String(), test.MessageEqual(expectedDep.Class, actual.Class.String(), testCase.Name))
				}

				if expectedDep.NestedCount > 0 {
					assert.Len(t, actual.Dependencies, expectedDep.NestedCount, test.MessageEqual(expectedDep.NestedCount, len(actual.Dependencies), testCase.Name))
					for j, nestedRef := range expectedDep.NestedRefs {
						if j < len(actual.Dependencies) {
							assert.Equal(t, nestedRef, actual.Dependencies[j].Reference, test.MessageEqual(nestedRef, actual.Dependencies[j].Reference, testCase.Name))
						}
					}
				}

				if expectedDep.ConfigOverrides != nil {
					for key, val := range expectedDep.ConfigOverrides {
						assert.Equal(t, val, actual.ConfigOverrides[key], test.MessageEqual(val, actual.ConfigOverrides[key], testCase.Name))
					}
				}
			}
		})
	}
}

func TestSetTestDependenciesDagDedup(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	ds := &DataStore{
		Classes: map[string]Class{
			"fvTenant": {
				Name:         testClassName("fvTenant"),
				ResourceName: "tenant",
			},
			"fvAp": {
				Name:         testClassName("fvAp"),
				ResourceName: "application_profile",
				Parents:      []*ClassName{testClassName("fvTenant")},
			},
			"fvBD": {
				Name:         testClassName("fvBD"),
				ResourceName: "bridge_domain",
				Parents:      []*ClassName{testClassName("fvTenant")},
			},
			"fvRsBd": {
				Name:         testClassName("fvRsBd"),
				ResourceName: "relation_to_bridge_domain",
				Parents:      []*ClassName{testClassName("fvAp")},
				Relation: Relation{
					RelationalClass: true,
					ToClasses:       []*ClassName{testClassName("fvBD")},
				},
				ClassDefinition: ClassDefinition{},
			},
		},
		GlobalMetaDefinition: GlobalMetaDefinition{},
		ctx:                  NewContext(),
	}

	class := ds.Classes["fvRsBd"]
	class.setTestDependencies(ds)

	var tenantRefs []*TestDependency
	collectAllDeps(class.TestDependencies, &tenantRefs, "aci_tenant.test.id")
	if len(tenantRefs) > 1 {
		for i := 1; i < len(tenantRefs); i++ {
			assert.Same(t, tenantRefs[0], tenantRefs[i], "DAG dedup should reuse the same pointer")
		}
	}
}

type parsePlaceholderExpected struct {
	Reference string
	Ok        bool
}

func TestParsePlaceholder(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_valid_placeholder",
			Input:    "{{aci_bridge_domain.test.id}}",
			Expected: parsePlaceholderExpected{Reference: "aci_bridge_domain.test.id", Ok: true},
		},
		{
			Name:     "test_leading_whitespace",
			Input:    "{{ aci_bridge_domain.test.id}}",
			Expected: parsePlaceholderExpected{Reference: "aci_bridge_domain.test.id", Ok: true},
		},
		{
			Name:     "test_trailing_whitespace",
			Input:    "{{aci_bridge_domain.test.id }}",
			Expected: parsePlaceholderExpected{Reference: "aci_bridge_domain.test.id", Ok: true},
		},
		{
			Name:     "test_both_whitespace",
			Input:    "{{ aci_bridge_domain.test.id }}",
			Expected: parsePlaceholderExpected{Reference: "aci_bridge_domain.test.id", Ok: true},
		},
		{
			Name:     "test_not_a_placeholder_plain_string",
			Input:    "aci_tenant.test.id",
			Expected: parsePlaceholderExpected{Reference: "", Ok: false},
		},
		{
			Name:     "test_not_a_placeholder_only_prefix",
			Input:    "{{aci_tenant.test.id",
			Expected: parsePlaceholderExpected{Reference: "", Ok: false},
		},
		{
			Name:     "test_not_a_placeholder_only_suffix",
			Input:    "aci_tenant.test.id}}",
			Expected: parsePlaceholderExpected{Reference: "", Ok: false},
		},
		{
			Name:     "test_empty_string",
			Input:    "",
			Expected: parsePlaceholderExpected{Reference: "", Ok: false},
		},
		{
			Name:     "test_prefix_not_at_start",
			Input:    "prefix{{aci_tenant.test.id}}",
			Expected: parsePlaceholderExpected{Reference: "", Ok: false},
		},
		{
			Name:     "test_suffix_not_at_end",
			Input:    "{{aci_tenant.test.id}}suffix",
			Expected: parsePlaceholderExpected{Reference: "", Ok: false},
		},
		{
			Name:     "test_single_braces",
			Input:    "{aci_tenant.test.id}",
			Expected: parsePlaceholderExpected{Reference: "", Ok: false},
		},
		{
			Name:     "test_tab_and_newline_whitespace_trimmed",
			Input:    "{{\t\naci_tenant.test.id\n\t}}",
			Expected: parsePlaceholderExpected{Reference: "aci_tenant.test.id", Ok: true},
		},
		{
			Name: "test_empty_placeholder_documents_current_behavior",
			// Known edge: "{{}}" parses as a placeholder with an empty reference.
			// validateTestDependencyPlaceholders / lookup-by-reference will then fail
			// loudly downstream, so this is preferred over silently returning false here.
			Input:    "{{}}",
			Expected: parsePlaceholderExpected{Reference: "", Ok: true},
		},
		{
			Name:     "test_whitespace_only_placeholder",
			Input:    "{{   }}",
			Expected: parsePlaceholderExpected{Reference: "", Ok: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(string)
			expected := testCase.Expected.(parsePlaceholderExpected)

			ref, ok := parsePlaceholder(input)

			assert.Equal(t, expected.Ok, ok, test.MessageEqual(expected.Ok, ok, testCase.Name))
			assert.Equal(t, expected.Reference, ref, test.MessageEqual(expected.Reference, ref, testCase.Name))
		})
	}
}

// TestIsPlaceholder exercises the boolean predicate directly. parsePlaceholder
// delegates to isPlaceholder for its prefix/suffix check, but isPlaceholder is
// also called on its own (e.g. validateTestDependencyPlaceholders) so its
// behavior is worth pinning independently.
func TestIsPlaceholder(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	cases := []test.TestCase{
		{Name: "valid", Input: "{{aci_tenant.test.id}}", Expected: true},
		{Name: "valid_with_whitespace", Input: "{{ aci_tenant.test.id }}", Expected: true},
		{Name: "empty_placeholder", Input: "{{}}", Expected: true},
		{Name: "plain_string", Input: "aci_tenant.test.id", Expected: false},
		{Name: "empty_string", Input: "", Expected: false},
		{Name: "only_prefix", Input: "{{aci_tenant.test.id", Expected: false},
		{Name: "only_suffix", Input: "aci_tenant.test.id}}", Expected: false},
		{Name: "prefix_not_at_start", Input: "prefix{{x}}", Expected: false},
		{Name: "suffix_not_at_end", Input: "{{x}}suffix", Expected: false},
		{Name: "single_braces", Input: "{x}", Expected: false},
	}

	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			got := isPlaceholder(testCase.Input.(string))
			assert.Equal(t, testCase.Expected.(bool), got, test.MessageEqual(testCase.Expected, got, testCase.Name))
		})
	}
}

func TestConfigOverridePlaceholderWhitespace(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	ds := &DataStore{
		Classes: map[string]Class{
			"fvTenant": {
				Name:         testClassName("fvTenant"),
				ResourceName: "tenant",
			},
			"fvBD": {
				Name:         testClassName("fvBD"),
				ResourceName: "bridge_domain",
				Parents:      []*ClassName{testClassName("fvTenant")},
			},
			"fvRsBd": {
				Name:         testClassName("fvRsBd"),
				ResourceName: "relation_to_bridge_domain",
				ClassDefinition: ClassDefinition{
					TestConfig: ClassTestConfigDefinition{
						ReplaceAutoResolved: true,
						Dependencies: []TestDependencyDefinition{
							{
								ClassName:     "fvTenant",
								Reference:     "aci_tenant.test.id",
								ReferenceType: ResourceReference,
								Role:          Parent,
							},
							{
								ClassName:       "fvBD",
								Reference:       "aci_bridge_domain.test.id",
								ReferenceType:   ResourceReference,
								Role:            Target,
								ConfigOverrides: map[string]string{"tenant_dn": "{{ aci_tenant.test.id }}"},
								Dependencies: []TestDependencyDefinition{
									{
										ClassName:     "fvTenant",
										Reference:     "aci_tenant.test.id",
										ReferenceType: ResourceReference,
									},
								},
							},
						},
					},
				},
			},
		},
		GlobalMetaDefinition: GlobalMetaDefinition{},
		ctx:                  NewContext(),
	}

	class := ds.Classes["fvRsBd"]
	class.setTestDependencies(ds)

	// ConfigOverrides placeholder with whitespace should resolve correctly.
	assert.Len(t, class.TestDependencies, 2)
	assert.Equal(t, "aci_tenant.test.id", class.TestDependencies[1].ConfigOverrides["tenant_dn"])
}

// collectAllDeps recursively collects all TestDependency nodes with the given reference.
func collectAllDeps(deps []*TestDependency, result *[]*TestDependency, ref string) {
	for _, d := range deps {
		if d.Reference == ref {
			*result = append(*result, d)
		}
		collectAllDeps(d.Dependencies, result, ref)
	}
}

type resolvePropertyTestValuesInput struct {
	Properties       map[string]*Property
	TestDependencies []*TestDependency
	Relation         Relation
}

type resolvePropertyTestValuesExpected struct {
	PropertyChecks     map[string]expectedPropertyTestValues
	DiagnosticContains string
}

type expectedPropertyTestValues struct {
	Nil            bool
	CreateValue    string
	CreateType     ValueRenderTypeEnum
	UpdateValue    string
	UpdateType     ValueRenderTypeEnum
	DefaultValue   string
	DefaultInclude bool
	ForceNewValue  string
	ForceNewType   ValueRenderTypeEnum
	ForceNewNil    bool
}

func TestResolvePropertyTestValues(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	boolTrue := true

	testCases := []test.TestCase{
		{
			Name: "test_auto_wire_target_dn",
			Input: resolvePropertyTestValuesInput{
				Properties: map[string]*Property{
					"tDn": {
						PropertyName:  "tDn",
						AttributeName: "target_dn",
						Required:      true,
					},
				},
				Relation: Relation{RelationalClass: true, Type: Explicit},
				TestDependencies: []*TestDependency{
					{Class: testClassName("fvBD"), Reference: "aci_bridge_domain.test.id", ReferenceType: ResourceReference, Role: Target},
					{Class: testClassName("fvBD"), Reference: "aci_bridge_domain.test_2.id", ReferenceType: ResourceReference, Role: Target},
				},
			},
			Expected: resolvePropertyTestValuesExpected{
				PropertyChecks: map[string]expectedPropertyTestValues{
					"tDn": {CreateValue: "aci_bridge_domain.test.id", CreateType: ReferenceValue, UpdateValue: "aci_bridge_domain.test_2.id", UpdateType: ReferenceValue, ForceNewValue: "aci_bridge_domain.test.id", ForceNewType: ReferenceValue},
				},
			},
		},
		{
			Name: "test_placeholder_resolution",
			Input: resolvePropertyTestValuesInput{
				Properties: func() map[string]*Property {
					prop := &Property{
						PropertyName:  "tDn",
						AttributeName: "target_dn",
						propertyDefinition: PropertyDefinition{
							TestConfig: TestConfigDefinition{
								Create: []TestValueEntryDefinition{
									{ConfigValue: "{{aci_tenant.test.id}}", ConfigInclude: &boolTrue, ValueType: ReferenceValue},
								},
							},
						},
					}
					prop.setTestValues()
					return map[string]*Property{"tDn": prop}
				}(),
				TestDependencies: []*TestDependency{
					{Class: testClassName("fvTenant"), Reference: "aci_tenant.test.id", ReferenceType: ResourceReference, Role: Parent},
				},
			},
			Expected: resolvePropertyTestValuesExpected{
				PropertyChecks: map[string]expectedPropertyTestValues{
					"tDn": {CreateValue: "aci_tenant.test.id", CreateType: ReferenceValue},
				},
			},
		},
		{
			Name: "test_auto_wire_parent_dn",
			Input: resolvePropertyTestValuesInput{
				Properties: map[string]*Property{
					"parentDn": {
						PropertyName:  "parentDn",
						AttributeName: "parent_dn",
						Required:      true,
					},
				},
				TestDependencies: []*TestDependency{
					{Class: testClassName("fvTenant"), Reference: "aci_tenant.test.id", ReferenceType: ResourceReference, Role: Parent},
					{Class: testClassName("fvTenant"), Reference: "aci_tenant.test_2.id", ReferenceType: ResourceReference, Role: Parent},
				},
			},
			Expected: resolvePropertyTestValuesExpected{
				PropertyChecks: map[string]expectedPropertyTestValues{
					"parentDn": {CreateValue: "aci_tenant.test.id", CreateType: ReferenceValue, UpdateValue: "aci_tenant.test.id", UpdateType: ReferenceValue, DefaultValue: "aci_tenant.test.id", DefaultInclude: true, ForceNewValue: "aci_tenant.test_2.id", ForceNewType: ReferenceValue},
				},
			},
		},
		{
			Name: "test_auto_wire_parent_dn_skip_if_no_parent_deps",
			Input: resolvePropertyTestValuesInput{
				Properties: map[string]*Property{
					"parentDn": {
						PropertyName:  "parentDn",
						AttributeName: "parent_dn",
						Required:      true,
					},
				},
				TestDependencies: []*TestDependency{},
			},
			Expected: resolvePropertyTestValuesExpected{
				PropertyChecks: map[string]expectedPropertyTestValues{
					"parentDn": {Nil: true},
				},
			},
		},
		{
			Name: "test_auto_wire_parent_dn_single_parent_dep",
			Input: resolvePropertyTestValuesInput{
				Properties: map[string]*Property{
					"parentDn": {
						PropertyName:  "parentDn",
						AttributeName: "parent_dn",
						Required:      true,
					},
				},
				TestDependencies: []*TestDependency{
					{Class: testClassName("fvTenant"), Reference: "aci_tenant.test.id", ReferenceType: ResourceReference, Role: Parent},
				},
			},
			Expected: resolvePropertyTestValuesExpected{
				PropertyChecks: map[string]expectedPropertyTestValues{
					"parentDn": {
						CreateValue:    "aci_tenant.test.id",
						CreateType:     ReferenceValue,
						UpdateValue:    "aci_tenant.test.id",
						UpdateType:     ReferenceValue,
						DefaultValue:   "aci_tenant.test.id",
						DefaultInclude: true,
						ForceNewNil:    true,
					},
				},
			},
		},
		{
			Name: "test_auto_wire_parent_dn_skip_if_explicit_config",
			Input: resolvePropertyTestValuesInput{
				Properties: func() map[string]*Property {
					prop := &Property{
						PropertyName:  "parentDn",
						AttributeName: "parent_dn",
						Required:      true,
						propertyDefinition: PropertyDefinition{
							TestConfig: TestConfigDefinition{
								Create: []TestValueEntryDefinition{
									{ConfigValue: "aci_tenant.custom.id", ConfigInclude: &boolTrue},
								},
							},
						},
					}
					prop.setTestValues()
					return map[string]*Property{"parentDn": prop}
				}(),
				TestDependencies: []*TestDependency{
					{Class: testClassName("fvTenant"), Reference: "aci_tenant.test.id", ReferenceType: ResourceReference, Role: Parent},
				},
			},
			Expected: resolvePropertyTestValuesExpected{
				PropertyChecks: map[string]expectedPropertyTestValues{
					"parentDn": {CreateValue: "aci_tenant.custom.id"},
				},
			},
		},
		{
			Name: "test_named_relation_two_targets_wires_name",
			Input: resolvePropertyTestValuesInput{
				Properties: map[string]*Property{
					"tnFvBDName": {PropertyName: "tnFvBDName", AttributeName: "tn_fv_bd_name", Required: true},
				},
				Relation: Relation{
					RelationalClass: true,
					Type:            Named,
					ToClasses:       []*ClassName{testClassName("fvBD")},
				},
				TestDependencies: []*TestDependency{
					{Class: testClassName("fvBD"), Reference: "aci_bridge_domain.test.id", ReferenceType: ResourceReference, Role: Target},
					{Class: testClassName("fvBD"), Reference: "aci_bridge_domain.test_2.id", ReferenceType: ResourceReference, Role: Target},
				},
			},
			Expected: resolvePropertyTestValuesExpected{
				PropertyChecks: map[string]expectedPropertyTestValues{
					"tnFvBDName": {
						CreateValue:   "aci_bridge_domain.test.name",
						CreateType:    ReferenceValue,
						UpdateValue:   "aci_bridge_domain.test_2.name",
						UpdateType:    ReferenceValue,
						ForceNewValue: "aci_bridge_domain.test.name",
						ForceNewType:  ReferenceValue,
					},
				},
			},
		},
		{
			Name: "test_named_relation_single_target_reuses_create",
			Input: resolvePropertyTestValuesInput{
				Properties: map[string]*Property{
					"tnFvBDName": {PropertyName: "tnFvBDName", AttributeName: "tn_fv_bd_name", Required: true},
				},
				Relation: Relation{
					RelationalClass: true,
					Type:            Named,
					ToClasses:       []*ClassName{testClassName("fvBD")},
				},
				TestDependencies: []*TestDependency{
					{Class: testClassName("fvBD"), Reference: "aci_bridge_domain.test.id", ReferenceType: ResourceReference, Role: Target},
				},
			},
			Expected: resolvePropertyTestValuesExpected{
				PropertyChecks: map[string]expectedPropertyTestValues{
					"tnFvBDName": {
						CreateValue:   "aci_bridge_domain.test.name",
						CreateType:    ReferenceValue,
						UpdateValue:   "aci_bridge_domain.test.name",
						UpdateType:    ReferenceValue,
						ForceNewValue: "aci_bridge_domain.test.name",
						ForceNewType:  ReferenceValue,
					},
				},
			},
		},
		{
			Name: "test_named_relation_no_property_no_op",
			Input: resolvePropertyTestValuesInput{
				Properties: map[string]*Property{},
				Relation: Relation{
					RelationalClass: true,
					Type:            Named,
					ToClasses:       []*ClassName{testClassName("fvBD")},
				},
				TestDependencies: []*TestDependency{
					{Class: testClassName("fvBD"), Reference: "aci_bridge_domain.test.id", ReferenceType: ResourceReference, Role: Target},
				},
			},
			Expected: resolvePropertyTestValuesExpected{
				PropertyChecks: map[string]expectedPropertyTestValues{},
			},
		},
		{
			Name: "test_named_relation_static_reference_diagnostic",
			Input: resolvePropertyTestValuesInput{
				Properties: map[string]*Property{
					"tnFvBDName": {PropertyName: "tnFvBDName", AttributeName: "tn_fv_bd_name", Required: true},
				},
				Relation: Relation{
					RelationalClass: true,
					Type:            Named,
					ToClasses:       []*ClassName{testClassName("fvBD")},
				},
				TestDependencies: []*TestDependency{
					{Class: testClassName("fvBD"), Reference: "uni/tn-common/BD-default", ReferenceType: StaticReference, Role: Target},
				},
			},
			Expected: resolvePropertyTestValuesExpected{
				PropertyChecks: map[string]expectedPropertyTestValues{
					"tnFvBDName": {Nil: true},
				},
				DiagnosticContains: "static reference",
			},
		},
		{
			Name: "test_named_relation_with_test_config_skipped",
			Input: resolvePropertyTestValuesInput{
				Properties: func() map[string]*Property {
					prop := &Property{
						PropertyName:  "tnFvBDName",
						AttributeName: "tn_fv_bd_name",
						propertyDefinition: PropertyDefinition{
							TestConfig: TestConfigDefinition{
								Create: []TestValueEntryDefinition{
									{ConfigValue: "explicit_name", ConfigInclude: &boolTrue},
								},
							},
						},
					}
					prop.setTestValues()
					return map[string]*Property{"tnFvBDName": prop}
				}(),
				Relation: Relation{
					RelationalClass: true,
					Type:            Named,
					ToClasses:       []*ClassName{testClassName("fvBD")},
				},
				TestDependencies: []*TestDependency{
					{Class: testClassName("fvBD"), Reference: "aci_bridge_domain.test.id", ReferenceType: ResourceReference, Role: Target},
				},
			},
			Expected: resolvePropertyTestValuesExpected{
				PropertyChecks: map[string]expectedPropertyTestValues{
					"tnFvBDName": {CreateValue: "explicit_name"},
				},
			},
		},
		{
			Name: "test_relational_unknown_type_diagnostic",
			Input: resolvePropertyTestValuesInput{
				Properties: map[string]*Property{
					"tDn": {PropertyName: "tDn", AttributeName: "target_dn", Required: true},
				},
				Relation: Relation{RelationalClass: true, Type: UndefinedRelationshipType},
				TestDependencies: []*TestDependency{
					{Class: testClassName("fvBD"), Reference: "aci_bridge_domain.test.id", ReferenceType: ResourceReference, Role: Target},
				},
			},
			Expected: resolvePropertyTestValuesExpected{
				PropertyChecks: map[string]expectedPropertyTestValues{
					"tDn": {Nil: true},
				},
				DiagnosticContains: "unsupported relationship type",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(resolvePropertyTestValuesInput)
			expected := testCase.Expected.(resolvePropertyTestValuesExpected)

			class := Class{
				Relation:         input.Relation,
				Properties:       input.Properties,
				TestDependencies: input.TestDependencies,
			}

			ds := &DataStore{ctx: NewContext()}
			class.setPropertyTestValues(ds)

			if expected.DiagnosticContains != "" {
				err := ds.ctx.Diagnostics.Error()
				assert.Error(t, err, testCase.Name+": expected a diagnostic")
				if err != nil {
					assert.True(t, strings.Contains(err.Error(), expected.DiagnosticContains), testCase.Name+": diagnostic should contain %q, got: %s", expected.DiagnosticContains, err.Error())
				}
			}

			for propName, check := range expected.PropertyChecks {
				prop := class.Properties[propName]
				if check.Nil {
					assert.Nil(t, prop.TestValues, testCase.Name+": "+propName+" TestValues should be nil")
					continue
				}
				assert.NotNil(t, prop.TestValues, testCase.Name+": "+propName+" TestValues should not be nil")
				if prop.TestValues == nil {
					continue
				}
				if check.CreateValue != "" {
					assert.Equal(t, check.CreateValue, prop.TestValues.Create[0].ConfigValue, test.MessageEqual(check.CreateValue, prop.TestValues.Create[0].ConfigValue, testCase.Name))
					if check.CreateType != 0 {
						assert.Equal(t, check.CreateType, prop.TestValues.Create[0].ValueType, test.MessageEqual(check.CreateType, prop.TestValues.Create[0].ValueType, testCase.Name))
					}
				}
				if check.UpdateValue != "" {
					assert.Equal(t, check.UpdateValue, prop.TestValues.Update[0].ConfigValue, test.MessageEqual(check.UpdateValue, prop.TestValues.Update[0].ConfigValue, testCase.Name))
					if check.UpdateType != 0 {
						assert.Equal(t, check.UpdateType, prop.TestValues.Update[0].ValueType, test.MessageEqual(check.UpdateType, prop.TestValues.Update[0].ValueType, testCase.Name))
					}
				}
				if check.DefaultValue != "" {
					assert.Equal(t, check.DefaultValue, prop.TestValues.Default[0].ConfigValue, test.MessageEqual(check.DefaultValue, prop.TestValues.Default[0].ConfigValue, testCase.Name))
					assert.Equal(t, check.DefaultInclude, prop.TestValues.Default[0].ConfigInclude, test.MessageEqual(check.DefaultInclude, prop.TestValues.Default[0].ConfigInclude, testCase.Name))
				}
				if check.ForceNewNil {
					assert.Nil(t, prop.TestValues.ForceNew, testCase.Name+": "+propName+" ForceNew should be nil")
				} else if check.ForceNewValue != "" {
					assert.Equal(t, check.ForceNewValue, prop.TestValues.ForceNew[0].ConfigValue, test.MessageEqual(check.ForceNewValue, prop.TestValues.ForceNew[0].ConfigValue, testCase.Name))
					if check.ForceNewType != 0 {
						assert.Equal(t, check.ForceNewType, prop.TestValues.ForceNew[0].ValueType, test.MessageEqual(check.ForceNewType, prop.TestValues.ForceNew[0].ValueType, testCase.Name))
					}
				}
			}
		})
	}
}

type resolveChildTestValuesInput struct {
	Classes     map[string]Class
	TargetClass string
}

type resolveChildTestValuesExpected struct {
	ChildCount int
	Children   []expectedTestChild
}

type expectedTestChild struct {
	ClassName     string
	InstanceCount int
	Instances     []expectedTestChildInstance
}

type expectedTestChildInstance struct {
	Properties map[string]expectedTestChildProperty
}

type expectedTestChildProperty struct {
	ConfigValue string
	ValueType   ValueRenderTypeEnum
}

func TestResolveChildTestValues(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_single_nested",
			Input: resolveChildTestValuesInput{
				TargetClass: "fvAEPg",
				Classes: map[string]Class{
					"fvRsBd": {
						Name:                             testClassName("fvRsBd"),
						ResourceName:                     "relation_to_bridge_domain",
						IsSingleNestedWhenDefinedAsChild: true,
						Properties: map[string]*Property{
							"tDn": {
								PropertyName:  "tDn",
								AttributeName: "target_dn",
								Required:      true,
								TestValues: &TestValues{
									Create: []TestValueEntry{{ConfigValue: "aci_bd.test.id", ConfigInclude: true, AssertValue: "aci_bd.test.id", ValueType: ReferenceValue}},
									Update: []TestValueEntry{{ConfigValue: "aci_bd.test_2.id", ConfigInclude: true, AssertValue: "aci_bd.test_2.id", ValueType: ReferenceValue}},
								},
							},
						},
					},
					"fvAEPg": {
						Name:            testClassName("fvAEPg"),
						ResourceName:    "application_epg",
						Children:        []*ClassName{testClassName("fvRsBd")},
						ClassDefinition: ClassDefinition{},
					},
				},
			},
			Expected: resolveChildTestValuesExpected{
				ChildCount: 1,
				Children: []expectedTestChild{
					{
						ClassName:     "fvRsBd",
						InstanceCount: 1,
						Instances: []expectedTestChildInstance{
							{Properties: map[string]expectedTestChildProperty{"target_dn": {ConfigValue: "aci_bd.test.id", ValueType: ReferenceValue}}},
						},
					},
				},
			},
		},
		{
			Name: "test_list_type",
			Input: resolveChildTestValuesInput{
				TargetClass: "fvTenant",
				Classes: map[string]Class{
					"tagAnnotation": {
						Name:                             testClassName("tagAnnotation"),
						ResourceName:                     "annotation",
						IsSingleNestedWhenDefinedAsChild: false,
						Properties: map[string]*Property{
							"key": {
								PropertyName:  "key",
								AttributeName: "key",
								Required:      true,
								TestValues: &TestValues{
									Create: []TestValueEntry{{ConfigValue: "key_0", ConfigInclude: true, AssertValue: "key_0", ValueType: StringValue}},
									Update: []TestValueEntry{{ConfigValue: "key_1", ConfigInclude: true, AssertValue: "key_1", ValueType: StringValue}},
								},
							},
							"value": {
								PropertyName:  "value",
								AttributeName: "value",
								Optional:      true,
								TestValues: &TestValues{
									Create: []TestValueEntry{{ConfigValue: "value_1", ConfigInclude: true, AssertValue: "value_1", ValueType: StringValue}},
									Update: []TestValueEntry{{ConfigValue: "value_2", ConfigInclude: true, AssertValue: "value_2", ValueType: StringValue}},
								},
							},
						},
					},
					"fvTenant": {
						Name:            testClassName("fvTenant"),
						ResourceName:    "tenant",
						Children:        []*ClassName{testClassName("tagAnnotation")},
						ClassDefinition: ClassDefinition{},
					},
				},
			},
			Expected: resolveChildTestValuesExpected{
				ChildCount: 1,
				Children: []expectedTestChild{
					{
						ClassName:     "tagAnnotation",
						InstanceCount: 2,
						Instances: []expectedTestChildInstance{
							{Properties: map[string]expectedTestChildProperty{"key": {ConfigValue: "key_0", ValueType: StringValue}, "value": {ConfigValue: "value_1", ValueType: StringValue}}},
							{Properties: map[string]expectedTestChildProperty{"key": {ConfigValue: "key_1", ValueType: StringValue}, "value": {ConfigValue: "value_2", ValueType: StringValue}}},
						},
					},
				},
			},
		},
		{
			Name: "test_with_override",
			Input: resolveChildTestValuesInput{
				TargetClass: "fvTenant",
				Classes: map[string]Class{
					"tagAnnotation": {
						Name:                             testClassName("tagAnnotation"),
						ResourceName:                     "annotation",
						IsSingleNestedWhenDefinedAsChild: false,
						Properties: map[string]*Property{
							"key": {
								PropertyName:  "key",
								AttributeName: "key",
								Required:      true,
								TestValues: &TestValues{
									Create: []TestValueEntry{{ConfigValue: "auto_key", ConfigInclude: true, AssertValue: "auto_key", ValueType: StringValue}},
									Update: []TestValueEntry{{ConfigValue: "auto_key_2", ConfigInclude: true, AssertValue: "auto_key_2", ValueType: StringValue}},
								},
							},
						},
					},
					"fvTenant": {
						Name:         testClassName("fvTenant"),
						ResourceName: "tenant",
						Children:     []*ClassName{testClassName("tagAnnotation")},
						ClassDefinition: ClassDefinition{
							TestConfig: ClassTestConfigDefinition{
								Children: map[string]ChildTestOverrideDefinition{
									"tagAnnotation": {
										Instances: []ChildTestInstanceOverrideDefinition{
											{Properties: map[string]string{"key": "override_key_0"}},
											{Properties: map[string]string{"key": "override_key_1"}},
										},
									},
								},
							},
						},
					},
				},
			},
			Expected: resolveChildTestValuesExpected{
				ChildCount: 1,
				Children: []expectedTestChild{
					{
						ClassName:     "tagAnnotation",
						InstanceCount: 2,
						Instances: []expectedTestChildInstance{
							{Properties: map[string]expectedTestChildProperty{"key": {ConfigValue: "override_key_0", ValueType: StringValue}}},
							{Properties: map[string]expectedTestChildProperty{"key": {ConfigValue: "override_key_1", ValueType: StringValue}}},
						},
					},
				},
			},
		},
		{
			Name: "test_placeholder_resolution",
			Input: resolveChildTestValuesInput{
				TargetClass: "fvAEPg",
				Classes: map[string]Class{
					"fvRsBd": {
						Name:                             testClassName("fvRsBd"),
						ResourceName:                     "relation_to_bridge_domain",
						IsSingleNestedWhenDefinedAsChild: true,
						Properties: map[string]*Property{
							"tDn": {
								PropertyName:  "tDn",
								AttributeName: "target_dn",
								Required:      true,
								TestValues: &TestValues{
									Create: []TestValueEntry{{ConfigValue: "aci_bridge_domain.test.id", ConfigInclude: true, AssertValue: "aci_bridge_domain.test.id", ValueType: ReferenceValue}},
									Update: []TestValueEntry{{ConfigValue: "aci_bridge_domain.test_2.id", ConfigInclude: true, AssertValue: "aci_bridge_domain.test_2.id", ValueType: ReferenceValue}},
								},
							},
						},
					},
					"fvAEPg": {
						Name:         testClassName("fvAEPg"),
						ResourceName: "application_epg",
						Children:     []*ClassName{testClassName("fvRsBd")},
						ClassDefinition: ClassDefinition{
							TestConfig: ClassTestConfigDefinition{
								Children: map[string]ChildTestOverrideDefinition{
									"fvRsBd": {
										Instances: []ChildTestInstanceOverrideDefinition{
											{Properties: map[string]string{"target_dn": "{{aci_bridge_domain.test.id}}"}},
										},
									},
								},
							},
						},
						TestDependencies: []*TestDependency{
							{Class: testClassName("fvBD"), Reference: "aci_bridge_domain.test.id", ReferenceType: ResourceReference, Role: Target},
						},
					},
				},
			},
			Expected: resolveChildTestValuesExpected{
				ChildCount: 1,
				Children: []expectedTestChild{
					{
						ClassName:     "fvRsBd",
						InstanceCount: 1,
						Instances: []expectedTestChildInstance{
							{Properties: map[string]expectedTestChildProperty{"target_dn": {ConfigValue: "aci_bridge_domain.test.id", ValueType: ReferenceValue}}},
						},
					},
				},
			},
		},
		{
			// Item 8 case 3: list-type child where a property has only Create (no Update).
			// instance[1] must fall back to the Create value for that property.
			Name: "test_list_child_update_fallback_to_create",
			Input: resolveChildTestValuesInput{
				TargetClass: "fvTenant",
				Classes: map[string]Class{
					"tagAnnotation": {
						Name:                             testClassName("tagAnnotation"),
						ResourceName:                     "annotation",
						IsSingleNestedWhenDefinedAsChild: false,
						Properties: map[string]*Property{
							"key": {
								PropertyName:  "key",
								AttributeName: "key",
								Required:      true,
								TestValues: &TestValues{
									Create: []TestValueEntry{{ConfigValue: "key_0", ConfigInclude: true, AssertValue: "key_0", ValueType: StringValue}},
									Update: []TestValueEntry{{ConfigValue: "key_1", ConfigInclude: true, AssertValue: "key_1", ValueType: StringValue}},
								},
							},
							"value": {
								PropertyName:  "value",
								AttributeName: "value",
								Optional:      true,
								// Update is intentionally missing — instance[1] must reuse Create.
								TestValues: &TestValues{
									Create: []TestValueEntry{{ConfigValue: "value_only_create", ConfigInclude: true, AssertValue: "value_only_create", ValueType: StringValue}},
								},
							},
						},
					},
					"fvTenant": {
						Name:            testClassName("fvTenant"),
						ResourceName:    "tenant",
						Children:        []*ClassName{testClassName("tagAnnotation")},
						ClassDefinition: ClassDefinition{},
					},
				},
			},
			Expected: resolveChildTestValuesExpected{
				ChildCount: 1,
				Children: []expectedTestChild{
					{
						ClassName:     "tagAnnotation",
						InstanceCount: 2,
						Instances: []expectedTestChildInstance{
							{Properties: map[string]expectedTestChildProperty{"key": {ConfigValue: "key_0", ValueType: StringValue}, "value": {ConfigValue: "value_only_create", ValueType: StringValue}}},
							{Properties: map[string]expectedTestChildProperty{"key": {ConfigValue: "key_1", ValueType: StringValue}, "value": {ConfigValue: "value_only_create", ValueType: StringValue}}},
						},
					},
				},
			},
		},
		{
			// Item 8 case 4: child has a mix of normal, IgnoreInTest, and ReadOnly properties.
			// Only the normal property must appear in the instance.
			Name: "test_child_skips_ignored_and_readonly",
			Input: resolveChildTestValuesInput{
				TargetClass: "fvTenant",
				Classes: map[string]Class{
					"tagAnnotation": {
						Name:                             testClassName("tagAnnotation"),
						ResourceName:                     "annotation",
						IsSingleNestedWhenDefinedAsChild: true,
						Properties: map[string]*Property{
							"key": {
								PropertyName:  "key",
								AttributeName: "key",
								Required:      true,
								TestValues: &TestValues{
									Create: []TestValueEntry{{ConfigValue: "k", ConfigInclude: true, AssertValue: "k", ValueType: StringValue}},
								},
							},
							"ignored": {
								PropertyName:  "ignored",
								AttributeName: "ignored",
								Optional:      true,
								IgnoreInTest:  true,
								TestValues: &TestValues{
									Create: []TestValueEntry{{ConfigValue: "skip_me", ConfigInclude: true, AssertValue: "skip_me", ValueType: StringValue}},
								},
							},
							"readOnly": {
								PropertyName:  "readOnly",
								AttributeName: "read_only",
								ReadOnly:      true,
								TestValues: &TestValues{
									Create: []TestValueEntry{{ConfigValue: "skip_me_too", ConfigInclude: true, AssertValue: "skip_me_too", ValueType: StringValue}},
								},
							},
						},
					},
					"fvTenant": {
						Name:         testClassName("fvTenant"),
						ResourceName: "tenant",
						Children:     []*ClassName{testClassName("tagAnnotation")},
					},
				},
			},
			Expected: resolveChildTestValuesExpected{
				ChildCount: 1,
				Children: []expectedTestChild{
					{
						ClassName:     "tagAnnotation",
						InstanceCount: 1,
						Instances: []expectedTestChildInstance{
							{Properties: map[string]expectedTestChildProperty{"key": {ConfigValue: "k", ValueType: StringValue}}},
						},
					},
				},
			},
		},
		{
			// Item 8 case 8: override value is a {{placeholder}} but no matching TestDependency
			// exists, so resolution leaves the value as-is. ValueType must still be inferred
			// as ReferenceValue by buildOverrideInstances based on the placeholder shape.
			Name: "test_child_override_with_placeholder_inferred_as_reference",
			Input: resolveChildTestValuesInput{
				TargetClass: "fvAEPg",
				Classes: map[string]Class{
					"fvRsBd": {
						Name:                             testClassName("fvRsBd"),
						ResourceName:                     "relation_to_bridge_domain",
						IsSingleNestedWhenDefinedAsChild: true,
						Properties: map[string]*Property{
							"tDn": {
								PropertyName:  "tDn",
								AttributeName: "target_dn",
								Required:      true,
								TestValues: &TestValues{
									Create: []TestValueEntry{{ConfigValue: "aci_bridge_domain.test.id", ConfigInclude: true, AssertValue: "aci_bridge_domain.test.id", ValueType: ReferenceValue}},
								},
							},
						},
					},
					"fvAEPg": {
						Name:         testClassName("fvAEPg"),
						ResourceName: "application_epg",
						Children:     []*ClassName{testClassName("fvRsBd")},
						ClassDefinition: ClassDefinition{
							TestConfig: ClassTestConfigDefinition{
								Children: map[string]ChildTestOverrideDefinition{
									"fvRsBd": {
										Instances: []ChildTestInstanceOverrideDefinition{
											{Properties: map[string]string{"target_dn": "{{aci_bridge_domain.unresolved.id}}"}},
										},
									},
								},
							},
						},
						// No TestDependency matching "aci_bridge_domain.unresolved.id" → placeholder stays.
					},
				},
			},
			Expected: resolveChildTestValuesExpected{
				ChildCount: 1,
				Children: []expectedTestChild{
					{
						ClassName:     "fvRsBd",
						InstanceCount: 1,
						Instances: []expectedTestChildInstance{
							{Properties: map[string]expectedTestChildProperty{"target_dn": {ConfigValue: "{{aci_bridge_domain.unresolved.id}}", ValueType: ReferenceValue}}},
						},
					},
				},
			},
		},
		{
			// Verifies the visited map in setChildTestValues prevents infinite recursion
			// when child relationships form a cycle (A → B → A).
			Name: "test_circular_children_skipped",
			Input: resolveChildTestValuesInput{
				TargetClass: "aaA",
				Classes: map[string]Class{
					"aaA": {
						Name:         testClassName("aaA"),
						ResourceName: "class_a",
						Children:     []*ClassName{testClassName("bbB")},
					},
					"bbB": {
						Name:         testClassName("bbB"),
						ResourceName: "class_b",
						Children:     []*ClassName{testClassName("aaA")},
					},
				},
			},
			Expected: resolveChildTestValuesExpected{
				// classA's TestChildren includes classB; classB's recursion would re-enter
				// classA but visited blocks it. Asserting we get at least one entry without
				// hanging is the safety check.
				ChildCount: 1,
			},
		},
		{
			// Override targets a child class that is NOT in the parent's resolved children.
			// applyChildOverrides logs a warning and continues without panicking.
			Name: "test_child_override_no_match_warns",
			Input: resolveChildTestValuesInput{
				TargetClass: "fvTenant",
				Classes: map[string]Class{
					"fvTenant": {
						Name:         testClassName("fvTenant"),
						ResourceName: "tenant",
						Children:     []*ClassName{},
						ClassDefinition: ClassDefinition{
							TestConfig: ClassTestConfigDefinition{
								Children: map[string]ChildTestOverrideDefinition{
									"unrelatedClass": {
										Instances: []ChildTestInstanceOverrideDefinition{
											{Properties: map[string]string{"foo": "bar"}},
										},
									},
								},
							},
						},
					},
				},
			},
			Expected: resolveChildTestValuesExpected{
				ChildCount: 0,
			},
		},
		{
			// Verifies item 2: child-driven dependency collection walks the child class's
			// full TestDependencies DAG (not just top-level), so a reference inside a
			// child instance property to a deeply nested dependency is still picked up.
			Name: "test_child_driven_dependency_collected_recursively",
			Input: resolveChildTestValuesInput{
				TargetClass: "fvAEPg",
				Classes: map[string]Class{
					"fvRsBd": {
						Name:                             testClassName("fvRsBd"),
						ResourceName:                     "relation_to_bridge_domain",
						IsSingleNestedWhenDefinedAsChild: true,
						Properties: map[string]*Property{
							"tDn": {
								PropertyName:  "tDn",
								AttributeName: "target_dn",
								Required:      true,
								TestValues: &TestValues{
									Create: []TestValueEntry{{ConfigValue: "aci_bridge_domain.deep.id", ConfigInclude: true, AssertValue: "aci_bridge_domain.deep.id", ValueType: ReferenceValue}},
								},
							},
						},
						// Reference appears nested inside fvBD's own dependencies (depth 1),
						// not at top level. Item 2's recursive lookup must find it.
						TestDependencies: []*TestDependency{
							{
								Class:         testClassName("fvTenant"),
								Reference:     "aci_tenant.test.id",
								ReferenceType: ResourceReference,
								Role:          Parent,
								Dependencies: []*TestDependency{
									{Class: testClassName("fvBD"), Reference: "aci_bridge_domain.deep.id", ReferenceType: ResourceReference},
								},
							},
						},
					},
					"fvAEPg": {
						Name:         testClassName("fvAEPg"),
						ResourceName: "application_epg",
						Children:     []*ClassName{testClassName("fvRsBd")},
					},
				},
			},
			Expected: resolveChildTestValuesExpected{
				ChildCount: 1,
				Children: []expectedTestChild{
					{
						ClassName:     "fvRsBd",
						InstanceCount: 1,
						Instances: []expectedTestChildInstance{
							{Properties: map[string]expectedTestChildProperty{"target_dn": {ConfigValue: "aci_bridge_domain.deep.id", ValueType: ReferenceValue}}},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(resolveChildTestValuesInput)
			expected := testCase.Expected.(resolveChildTestValuesExpected)

			ds := &DataStore{
				Classes:              input.Classes,
				GlobalMetaDefinition: GlobalMetaDefinition{},
				ctx:                  NewContext(),
			}

			class := ds.Classes[input.TargetClass]
			class.setChildTestValues(ds)

			assert.Len(t, class.TestChildren, expected.ChildCount, test.MessageEqual(expected.ChildCount, len(class.TestChildren), testCase.Name))

			for i, expectedChild := range expected.Children {
				if i >= len(class.TestChildren) {
					break
				}
				actual := class.TestChildren[i]
				assert.Equal(t, expectedChild.ClassName, actual.Class.String(), test.MessageEqual(expectedChild.ClassName, actual.Class.String(), testCase.Name))
				assert.Len(t, actual.Instances, expectedChild.InstanceCount, test.MessageEqual(expectedChild.InstanceCount, len(actual.Instances), testCase.Name))

				for j, expectedInstance := range expectedChild.Instances {
					if j >= len(actual.Instances) {
						break
					}
					for propName, expectedProp := range expectedInstance.Properties {
						actualEntry, exists := actual.Instances[j].Properties[propName]
						assert.True(t, exists, testCase.Name+": property "+propName+" should exist in instance")
						if exists {
							assert.Equal(t, expectedProp.ConfigValue, actualEntry.ConfigValue, test.MessageEqual(expectedProp.ConfigValue, actualEntry.ConfigValue, testCase.Name))
							assert.Equal(t, expectedProp.ValueType, actualEntry.ValueType, test.MessageEqual(expectedProp.ValueType, actualEntry.ValueType, testCase.Name))
						}
					}
				}
			}
		})
	}
}

// Item 8 case 10: placeholder inside a grandchild (nested override Children at level 2)
// must resolve against the parent class's TestDependencies via the recursive
// resolvePlaceholdersInTestChildren call.
func TestResolveChildTestValuesGrandchildPlaceholder(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	// Grandchild class fvSubnet (under fvBD) carries a placeholder property; the
	// override-driven grandchild instance references aci_bridge_domain.test.id which is
	// declared at fvAEPg's TestDependencies level.
	classes := map[string]Class{
		"fvSubnet": {
			Name:                             testClassName("fvSubnet"),
			ResourceName:                     "subnet",
			IsSingleNestedWhenDefinedAsChild: true,
			Properties: map[string]*Property{
				"ip": {
					PropertyName:  "ip",
					AttributeName: "ip",
					Required:      true,
					TestValues: &TestValues{
						Create: []TestValueEntry{{ConfigValue: "10.0.0.1/24", ConfigInclude: true, AssertValue: "10.0.0.1/24", ValueType: StringValue}},
					},
				},
			},
		},
		"fvBD": {
			Name:         testClassName("fvBD"),
			ResourceName: "bridge_domain",
			Children:     []*ClassName{testClassName("fvSubnet")},
		},
		"fvAEPg": {
			Name:         testClassName("fvAEPg"),
			ResourceName: "application_epg",
			Children:     []*ClassName{testClassName("fvBD")},
			ClassDefinition: ClassDefinition{
				TestConfig: ClassTestConfigDefinition{
					Children: map[string]ChildTestOverrideDefinition{
						"fvBD": {
							Instances: []ChildTestInstanceOverrideDefinition{
								{
									Properties: map[string]string{},
									Children: map[string]ChildTestOverrideDefinition{
										"fvSubnet": {
											Instances: []ChildTestInstanceOverrideDefinition{
												{Properties: map[string]string{"ip": "{{aci_bridge_domain.test.id}}"}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			TestDependencies: []*TestDependency{
				{Class: testClassName("fvBD"), Reference: "aci_bridge_domain.test.id", ReferenceType: ResourceReference, Role: Target},
			},
		},
	}

	ds := &DataStore{
		Classes:              classes,
		GlobalMetaDefinition: GlobalMetaDefinition{},
		ctx:                  NewContext(),
	}

	class := ds.Classes["fvAEPg"]
	class.setChildTestValues(ds)

	// Drill down: fvAEPg.TestChildren[0] = fvBD, .Instances[0].Children[0] = fvSubnet, .Instances[0].Properties["ip"]
	assert.Len(t, class.TestChildren, 1, "expected one top-level child (fvBD)")
	if len(class.TestChildren) == 0 {
		return
	}
	bd := class.TestChildren[0]
	assert.Equal(t, "fvBD", bd.Class.String())
	assert.Len(t, bd.Instances, 1, "expected one fvBD instance from override")
	if len(bd.Instances) == 0 {
		return
	}
	assert.Len(t, bd.Instances[0].Children, 1, "expected one grandchild (fvSubnet)")
	if len(bd.Instances[0].Children) == 0 {
		return
	}
	subnet := bd.Instances[0].Children[0]
	assert.Equal(t, "fvSubnet", subnet.Class.String())
	assert.Len(t, subnet.Instances, 1)
	if len(subnet.Instances) == 0 {
		return
	}
	ip, ok := subnet.Instances[0].Properties["ip"]
	assert.True(t, ok, "grandchild property 'ip' must exist")
	assert.Equal(t, "aci_bridge_domain.test.id", ip.ConfigValue, "grandchild placeholder must resolve against parent TestDependencies")
	assert.Equal(t, ReferenceValue, ip.ValueType, "resolved grandchild value must be typed as ReferenceValue")
}

type validateTestCompletenessInput struct {
	TestDependencies []*TestDependency
	Properties       map[string]*Property
	TestChildren     []*TestChild
}

func TestValidateTestCompleteness(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_no_errors_when_all_resolved",
			Input: validateTestCompletenessInput{
				TestDependencies: []*TestDependency{
					{Reference: "aci_tenant.test.id", ConfigOverrides: map[string]string{"key": "resolved_value"}},
				},
				Properties: map[string]*Property{
					"name": {
						AttributeName: "name",
						TestValues: &TestValues{
							Create: []TestValueEntry{{ConfigValue: "test_name"}},
						},
					},
				},
			},
			Expected: 0,
		},
		{
			Name: "test_unresolved_config_override_placeholder",
			Input: validateTestCompletenessInput{
				TestDependencies: []*TestDependency{
					{Reference: "aci_bd.test.id", ConfigOverrides: map[string]string{"tenant_dn": "{{aci_tenant.test.id}}"}},
				},
			},
			Expected: 1,
		},
		{
			Name: "test_unresolved_nested_dependency_placeholder",
			Input: validateTestCompletenessInput{
				TestDependencies: []*TestDependency{
					{
						Reference: "aci_bd.test.id",
						Dependencies: []*TestDependency{
							{Reference: "aci_tenant.test.id", ConfigOverrides: map[string]string{"name": "{{aci_other.test.id}}"}},
						},
					},
				},
			},
			Expected: 1,
		},
		{
			Name: "test_unresolved_property_placeholder",
			Input: validateTestCompletenessInput{
				Properties: map[string]*Property{
					"tDn": {
						AttributeName: "target_dn",
						TestValues: &TestValues{
							Create: []TestValueEntry{{ConfigValue: "{{aci_bd.test.id}}"}},
							Update: []TestValueEntry{{ConfigValue: "{{aci_bd.test_2.id}}"}},
						},
					},
				},
			},
			Expected: 2,
		},
		{
			Name: "test_unresolved_child_placeholder",
			Input: validateTestCompletenessInput{
				TestChildren: []*TestChild{
					{
						Class: testClassName("fvRsBd"),
						Instances: []TestChildInstance{
							{Properties: map[string]TestValueEntry{"target_dn": {ConfigValue: "{{aci_bd.test.id}}"}}},
						},
					},
				},
			},
			Expected: 1,
		},
		{
			Name: "test_multiple_errors_accumulated",
			Input: validateTestCompletenessInput{
				TestDependencies: []*TestDependency{
					{Reference: "aci_bd.test.id", ConfigOverrides: map[string]string{"tenant_dn": "{{aci_tenant.test.id}}"}},
				},
				Properties: map[string]*Property{
					"tDn": {
						AttributeName: "target_dn",
						TestValues: &TestValues{
							Create:   []TestValueEntry{{ConfigValue: "{{aci_bd.test.id}}"}},
							ForceNew: []TestValueEntry{{ConfigValue: "{{aci_bd.test.id}}"}},
						},
					},
				},
				TestChildren: []*TestChild{
					{
						Class: testClassName("fvRsBd"),
						Instances: []TestChildInstance{
							{Properties: map[string]TestValueEntry{"target_dn": {ConfigValue: "{{aci_bd.test.id}}"}}},
						},
					},
				},
			},
			Expected: 4,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(validateTestCompletenessInput)
			expectedErrorCount := testCase.Expected.(int)

			ctx := NewContext()
			class := &Class{
				Name:             testClassName("testClass"),
				TestDependencies: input.TestDependencies,
				Properties:       input.Properties,
				TestChildren:     input.TestChildren,
			}

			class.validateTestCompleteness(ctx)

			assert.Len(t, ctx.Diagnostics.errors, expectedErrorCount, test.MessageEqual(expectedErrorCount, len(ctx.Diagnostics.errors), testCase.Name))
		})
	}
}

// TestValidateTestDependencyPlaceholders_CycleProtection verifies that the
// `visited` map prevents infinite recursion when the dependency DAG forms a
// cycle (A → B → A). The unresolved placeholder on each node must be reported
// exactly once.
func TestValidateTestDependencyPlaceholders_CycleProtection(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	dependencyA := &TestDependency{
		Reference:       "aci_a.test.id",
		ConfigOverrides: map[string]string{"k": "{{unresolved.a}}"},
	}
	dependencyB := &TestDependency{
		Reference:       "aci_b.test.id",
		ConfigOverrides: map[string]string{"k": "{{unresolved.b}}"},
	}
	dependencyA.Dependencies = []*TestDependency{dependencyB}
	dependencyB.Dependencies = []*TestDependency{dependencyA} // cycle

	ctx := NewContext()
	class := &Class{Name: testClassName("testClass"), TestDependencies: []*TestDependency{dependencyA}}

	class.validateTestDependencyPlaceholders(ctx, class.TestDependencies, make(map[*TestDependency]bool))

	assert.Len(t, ctx.Diagnostics.errors, 2,
		"each node in the cycle must be visited exactly once and emit one diagnostic")
}

// TestValidateTestDependencyPlaceholders_PerDependencyChildren verifies that
// unresolved placeholders inside a dependency's `Children` overrides are also
// surfaced (not only those inside `ConfigOverrides`).
func TestValidateTestDependencyPlaceholders_PerDependencyChildren(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	dependency := &TestDependency{
		Reference: "aci_bd.test.id",
		Children: map[string]*TestChild{
			"fvRsBd": {
				Class: testClassName("fvRsBd"),
				Instances: []TestChildInstance{
					{Properties: map[string]TestValueEntry{
						"target_dn": {ConfigValue: "{{aci_unknown.test.id}}"},
					}},
				},
			},
		},
	}

	ctx := NewContext()
	class := &Class{Name: testClassName("testClass"), TestDependencies: []*TestDependency{dependency}}

	class.validateTestDependencyPlaceholders(ctx, class.TestDependencies, make(map[*TestDependency]bool))

	assert.Len(t, ctx.Diagnostics.errors, 1)
	assert.Contains(t, ctx.Diagnostics.errors[0], "target_dn")
}

// TestResolvePlaceholdersInDependencyChildren_CycleProtection verifies the
// placeholder resolver also terminates on cyclic dependency DAGs.
func TestResolvePlaceholdersInDependencyChildren_CycleProtection(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	dependencyA := &TestDependency{
		Reference: "aci_a.test.id",
		Children: map[string]*TestChild{
			"fvRsBd": {
				Class: testClassName("fvRsBd"),
				Instances: []TestChildInstance{
					{Properties: map[string]TestValueEntry{
						"k": {ConfigValue: "{{aci_a.test.id}}"},
					}},
				},
			},
		},
	}
	dependencyB := &TestDependency{Reference: "aci_b.test.id"}
	dependencyA.Dependencies = []*TestDependency{dependencyB}
	dependencyB.Dependencies = []*TestDependency{dependencyA} // cycle

	class := &Class{
		Name:             testClassName("testClass"),
		TestDependencies: []*TestDependency{dependencyA},
	}

	// Must terminate; the placeholder resolves to dependencyA's own reference.
	class.resolvePlaceholdersInDependencyChildren(class.TestDependencies, make(map[*TestDependency]bool))

	entry := dependencyA.Children["fvRsBd"].Instances[0].Properties["k"]
	assert.Equal(t, "aci_a.test.id", entry.ConfigValue)
	assert.Equal(t, ReferenceValue, entry.ValueType)
}

// TestMergeOverrideChildren covers the three branches:
//   - base entry kept as-is (no overlay key)
//   - base entry replaced (matching overlay key with instances)
//   - overlay-only entry appended (no matching base)
//
// Plus edge cases: both empty → nil, invalid overlay class name → skipped with
// warning, overlay key with empty Instances → ignored.
func TestMergeOverrideChildren(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	t.Run("both_empty_returns_nil", func(t *testing.T) {
		t.Parallel()
		ds := &DataStore{Classes: map[string]Class{}}
		got := mergeOverrideChildren(ds, nil, nil)
		assert.Nil(t, got)
	})

	t.Run("base_only_kept_as_is", func(t *testing.T) {
		t.Parallel()
		ds := &DataStore{Classes: map[string]Class{}}
		base := []*TestChild{{Class: testClassName("fvRsBd")}}
		got := mergeOverrideChildren(ds, base, nil)
		assert.Len(t, got, 1)
		assert.Same(t, base[0], got[0], "base entry must be kept by-pointer when no overlay matches")
	})

	t.Run("base_replaced_by_overlay", func(t *testing.T) {
		t.Parallel()
		ds := &DataStore{Classes: map[string]Class{}}
		base := []*TestChild{{
			Class: testClassName("fvRsBd"),
			Instances: []TestChildInstance{
				{Properties: map[string]TestValueEntry{"k": {ConfigValue: "old"}}},
			},
		}}
		overlay := map[string]ChildTestOverrideDefinition{
			"fvRsBd": {Instances: []ChildTestInstanceOverrideDefinition{
				{Properties: map[string]string{"k": "new"}},
			}},
		}

		got := mergeOverrideChildren(ds, base, overlay)
		assert.Len(t, got, 1)
		assert.Equal(t, "fvRsBd", got[0].Class.String())
		assert.Len(t, got[0].Instances, 1)
		assert.Equal(t, "new", got[0].Instances[0].Properties["k"].ConfigValue)
	})

	t.Run("overlay_only_appended", func(t *testing.T) {
		t.Parallel()
		ds := &DataStore{Classes: map[string]Class{}}
		overlay := map[string]ChildTestOverrideDefinition{
			"fvRsBd": {Instances: []ChildTestInstanceOverrideDefinition{
				{Properties: map[string]string{"k": "v"}},
			}},
		}

		got := mergeOverrideChildren(ds, nil, overlay)
		assert.Len(t, got, 1)
		assert.Equal(t, "fvRsBd", got[0].Class.String())
		assert.Equal(t, "v", got[0].Instances[0].Properties["k"].ConfigValue)
	})

	t.Run("overlay_only_with_empty_instances_skipped", func(t *testing.T) {
		t.Parallel()
		ds := &DataStore{Classes: map[string]Class{}}
		overlay := map[string]ChildTestOverrideDefinition{
			"fvRsBd": {Instances: nil},
		}

		got := mergeOverrideChildren(ds, nil, overlay)
		assert.Empty(t, got, "overlay-only entry with no instances must not be appended")
	})

	t.Run("overlay_invalid_class_name_skipped", func(t *testing.T) {
		t.Parallel()
		ds := &DataStore{Classes: map[string]Class{}}
		overlay := map[string]ChildTestOverrideDefinition{
			// lowercase first letter of short name => NewClassName fails.
			"badclassname": {Instances: []ChildTestInstanceOverrideDefinition{
				{Properties: map[string]string{"k": "v"}},
			}},
		}

		got := mergeOverrideChildren(ds, nil, overlay)
		assert.Empty(t, got, "overlay entry with invalid class name must be skipped (logged as warning)")
	})

	t.Run("base_kept_when_overlay_has_empty_instances_for_same_class", func(t *testing.T) {
		t.Parallel()
		ds := &DataStore{Classes: map[string]Class{}}
		base := []*TestChild{{Class: testClassName("fvRsBd")}}
		overlay := map[string]ChildTestOverrideDefinition{
			// Match exists but Instances is empty -> base entry must be preserved unchanged.
			"fvRsBd": {Instances: nil},
		}

		got := mergeOverrideChildren(ds, base, overlay)
		assert.Len(t, got, 1)
		assert.Same(t, base[0], got[0])
	})
}

// TestSetTestDependencies_ParentCycleProtection verifies setTestDependencies
// terminates when two classes are each other's parent (A.Parents=[B],
// B.Parents=[A]). Cycle protection lives in buildDependency via the
// testDependencies map keyed by Reference.
func TestSetTestDependencies_ParentCycleProtection(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	classNameA := testClassName("fvA")
	classNameB := testClassName("fvB")

	classA := Class{
		Name:         classNameA,
		ResourceName: "a",
		Parents:      []*ClassName{classNameB},
	}
	classB := Class{
		Name:         classNameB,
		ResourceName: "b",
		Parents:      []*ClassName{classNameA},
	}

	ctx := NewContext()
	ds := &DataStore{
		Classes: map[string]Class{
			"fvA": classA,
			"fvB": classB,
		},
		ctx: ctx,
	}

	// Must terminate (no stack overflow / hang).
	classA.setTestDependencies(ds)

	// First parent of A is B → two top-level Parent deps: aci_b.test.id and aci_b.test_2.id.
	assert.Len(t, classA.TestDependencies, 2)
	for _, testDependency := range classA.TestDependencies {
		assert.Equal(t, Parent, testDependency.Role)
		assert.Equal(t, "fvB", testDependency.Class.String())
	}
	// Each B-dep recurses into A's parent (B) which dedupes against the map.
	// aci_b.test.id has Dependencies=[aci_a.test.id]; that nested aci_a.test.id
	// in turn points back at the already-built aci_b.test.id node (same pointer)
	// rather than building a fresh one — proving the cycle terminated.
	assert.Len(t, classA.TestDependencies[0].Dependencies, 1)
	nestedDependencyA := classA.TestDependencies[0].Dependencies[0]
	assert.Equal(t, "aci_a.test.id", nestedDependencyA.Reference)
	if assert.Len(t, nestedDependencyA.Dependencies, 1) {
		assert.Same(t, classA.TestDependencies[0], nestedDependencyA.Dependencies[0],
			"cycle must short-circuit by returning the already-built dependency pointer")
	}

	assert.NoError(t, ctx.Diagnostics.Error())
}

// TestGetTestDependenciesFromDefinitions_MissingRoleAtTopLevel verifies the
// depth==0 + UndefinedRole validation path: a top-level entry without a role
// must produce a diagnostic and be skipped.
func TestGetTestDependenciesFromDefinitions_MissingRoleAtTopLevel(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	ctx := NewContext()
	ds := &DataStore{Classes: map[string]Class{}, ctx: ctx}
	class := &Class{Name: testClassName("testClass")}

	dependencyDefinitions := []TestDependencyDefinition{
		{ClassName: "fvTenant", Reference: "aci_tenant.test.id"}, // Role omitted
	}

	got := class.getTestDependenciesFromDefinitions(dependencyDefinitions, ds, map[string]*TestDependency{}, 0)
	assert.Empty(t, got, "top-level entry without role must be skipped")
	assert.Len(t, ctx.Diagnostics.errors, 1)
	assert.Contains(t, ctx.Diagnostics.errors[0], "missing required 'role'")
}

// TestGetTestDependenciesFromDefinitions_DuplicateConflictingRoles verifies the
// duplicate-detection branch: two top-level entries with the same Reference but
// different non-Undefined Roles must emit a "conflicting roles" diagnostic.
func TestGetTestDependenciesFromDefinitions_DuplicateConflictingRoles(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	ctx := NewContext()
	ds := &DataStore{Classes: map[string]Class{}, ctx: ctx}
	class := &Class{Name: testClassName("testClass")}

	dependencyDefinitions := []TestDependencyDefinition{
		{ClassName: "fvTenant", Reference: "aci_tenant.test.id", Role: Parent},
		{ClassName: "fvTenant", Reference: "aci_tenant.test.id", Role: Target},
	}

	class.getTestDependenciesFromDefinitions(dependencyDefinitions, ds, map[string]*TestDependency{}, 0)
	if assert.Len(t, ctx.Diagnostics.errors, 1) {
		assert.Contains(t, ctx.Diagnostics.errors[0], "conflicting roles")
	}
}

// TestGetTestDependenciesFromDefinitions_DuplicateWithConfigOverrides verifies
// that a duplicate top-level entry carrying ConfigOverrides is flagged so the
// author folds them into the first declaration.
func TestGetTestDependenciesFromDefinitions_DuplicateWithConfigOverrides(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	ctx := NewContext()
	ds := &DataStore{Classes: map[string]Class{}, ctx: ctx}
	class := &Class{Name: testClassName("testClass")}

	dependencyDefinitions := []TestDependencyDefinition{
		{ClassName: "fvTenant", Reference: "aci_tenant.test.id", Role: Parent},
		{ClassName: "fvTenant", Reference: "aci_tenant.test.id", Role: Parent,
			ConfigOverrides: map[string]string{"name": "override"}},
	}

	class.getTestDependenciesFromDefinitions(dependencyDefinitions, ds, map[string]*TestDependency{}, 0)
	if assert.Len(t, ctx.Diagnostics.errors, 1) {
		assert.Contains(t, ctx.Diagnostics.errors[0], "carries config_overrides")
	}
}

// TestGetTestDependenciesFromDefinitions_RolePromotion verifies the
// nested-then-top-level promotion branch: a dep first introduced nested
// (Role=UndefinedRole) is promoted in place when later declared at top level
// with a real role — no diagnostic should fire.
func TestGetTestDependenciesFromDefinitions_RolePromotion(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	ctx := NewContext()
	ds := &DataStore{Classes: map[string]Class{}, ctx: ctx}
	class := &Class{Name: testClassName("testClass")}

	existingDependencies := map[string]*TestDependency{}
	// Seed the map with a nested-style entry (UndefinedRole).
	existingDependencies["aci_tenant.test.id"] = &TestDependency{
		Class:     testClassName("fvTenant"),
		Reference: "aci_tenant.test.id",
		Role:      UndefinedRole,
	}

	dependencyDefinitions := []TestDependencyDefinition{
		{ClassName: "fvTenant", Reference: "aci_tenant.test.id", Role: Parent},
	}

	class.getTestDependenciesFromDefinitions(dependencyDefinitions, ds, existingDependencies, 0)
	assert.NoError(t, ctx.Diagnostics.Error())
	assert.Equal(t, Parent, existingDependencies["aci_tenant.test.id"].Role,
		"nested entry must be promoted to Parent in place")
}

// classWithProperties builds a synthetic Class with a Properties map sufficient
// for setStateUpgrades / setPropertyStateUpgradeValues tests. The Properties map
// is keyed by meta PropertyName and contains the minimum fields touched by
// state_upgrades validation and distribution.
func classWithProperties(name string, props map[string]*Property, children []string) *Class {
	classNames := make([]*ClassName, 0, len(children))
	for _, child := range children {
		classNames = append(classNames, testClassName(child))
	}
	return &Class{
		Name:       testClassName(name),
		Properties: props,
		Children:   classNames,
	}
}

func TestSetStateUpgrades_HappyPathRenameAndTypeChange(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	class := classWithProperties("fvCtx", map[string]*Property{
		"pcEnfPref": {PropertyName: "pcEnfPref", AttributeName: "pc_enforcement_preference", Optional: true, Computed: true, ValueType: String},
		"name":      {PropertyName: "name", AttributeName: "name", Required: true, ValueType: String},
	}, nil)
	class.ClassDefinition.StateUpgrades = []StateUpgradeDefinition{
		{
			PriorSchemaVersion: 0,
			Attributes: map[string]AttributeUpgradeDefinition{
				"pcEnfPref": {LegacyAttribute: "pc_enf_pref"},
				"name":      {LegacyType: StringAttribute},
			},
		},
	}

	ds := &DataStore{ctx: NewContext()}
	class.setStateUpgrades(ds)
	err := ds.ctx.Diagnostics.Error()
	assert.NoError(t, err, test.MessageUnexpectedError(err))
	if assert.Len(t, class.StateUpgrades, 1) {
		entry := class.StateUpgrades[0]
		assert.Equal(t, "pc_enf_pref", entry.Attributes["pcEnfPref"].LegacyAttribute)
		assert.Equal(t, StringAttribute, entry.Attributes["name"].LegacyType)
	}

	class.setPropertyStateUpgradeValues()
	pcEnfPref := class.Properties["pcEnfPref"]
	if assert.Contains(t, pcEnfPref.StateUpgradeValues, 0) {
		v := pcEnfPref.StateUpgradeValues[0]
		assert.Equal(t, "pc_enf_pref", v.AttributeName)
		assert.Equal(t, String, v.Type)
		assert.True(t, v.Optional)
		assert.True(t, v.Computed)
	}
	nameProp := class.Properties["name"]
	if assert.Contains(t, nameProp.StateUpgradeValues, 0) {
		v := nameProp.StateUpgradeValues[0]
		assert.Equal(t, "name", v.AttributeName, "unchanged name carries forward")
		assert.True(t, v.Required)
	}
}

// TestSetStateUpgrades_DuplicatePriorSchemaVersion covers the uniqueness rule
// on prior_schema_version across entries. Each saved-state version must map to
// exactly one upgrader, so a repeated value is rejected even when the entries
// themselves are otherwise valid.
func TestSetStateUpgrades_DuplicatePriorSchemaVersion(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	class := classWithProperties("fvCtx", map[string]*Property{"name": {PropertyName: "name"}}, nil)
	class.ClassDefinition.StateUpgrades = []StateUpgradeDefinition{
		{PriorSchemaVersion: 0, Attributes: map[string]AttributeUpgradeDefinition{"name": {LegacyAttribute: "n"}}},
		{PriorSchemaVersion: 0, Attributes: map[string]AttributeUpgradeDefinition{"name": {LegacyAttribute: "m"}}},
	}

	ds := &DataStore{ctx: NewContext()}
	class.setStateUpgrades(ds)
	err := ds.ctx.Diagnostics.Error()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "duplicate prior_schema_version 0")
	}
}

// TestValidateStateUpgrades_MigrationSourceCoherence covers the rule that a
// non-zero MigrationSource requires at least one state_upgrades entry. The
// specific prior_schema_version on that entry does not matter — SDKv2 resources
// may carry their own prior-provider upgrade history into the migration hop.
// Conversely, framework-native v0->v1 upgrades may declare state_upgrades
// entries without setting a MigrationSource.
func TestValidateStateUpgrades_MigrationSourceCoherence(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	cases := []struct {
		name     string
		source   MigrationSourceEnum
		upgrades []StateUpgradeDefinition
		wantErr  string // empty == expect no error
	}{
		{
			name:    "source_set_no_entries_errors",
			source:  FromSDKv2,
			wantErr: "requires at least one state_upgrades entry",
		},
		{
			name:   "source_set_v0_entry_ok",
			source: FromSDKv2,
			upgrades: []StateUpgradeDefinition{
				{PriorSchemaVersion: 0, Attributes: map[string]AttributeUpgradeDefinition{"name": {LegacyAttribute: "n"}}},
			},
		},
		{
			name:   "source_set_non_zero_prior_version_ok",
			source: FromSDKv2,
			upgrades: []StateUpgradeDefinition{
				{PriorSchemaVersion: 2, Attributes: map[string]AttributeUpgradeDefinition{"name": {LegacyAttribute: "n"}}},
			},
		},
		{
			name: "no_source_v0_entry_ok",
			upgrades: []StateUpgradeDefinition{
				{PriorSchemaVersion: 0, Attributes: map[string]AttributeUpgradeDefinition{"name": {LegacyAttribute: "n"}}},
			},
		},
		{
			name: "no_source_no_entries_ok",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			class := classWithProperties("fvCtx", map[string]*Property{"name": {PropertyName: "name"}}, nil)
			class.ClassDefinition.MigrationSource = tc.source
			class.ClassDefinition.StateUpgrades = tc.upgrades

			ds := &DataStore{ctx: NewContext()}
			class.setStateUpgrades(ds)
			err := ds.ctx.Diagnostics.Error()
			if tc.wantErr == "" {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				assert.Equal(t, tc.source, class.MigrationSource)
				return
			}
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}

// TestValidateStateUpgrades_Exhaustiveness covers the top-level key resolution
// rule: every Attributes / Children key in an entry must resolve to a current
// Property / child class on the class, OR be marked LegacyStatus == Removed in
// which case the key may reference a meta-only / excluded name. The mixed case
// proves resolved + Removed allow-through paths are silent and only the truly
// unknown non-Removed key surfaces a diagnostic.
func TestValidateStateUpgrades_Exhaustiveness(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	cases := []struct {
		name     string
		props    map[string]*Property
		children []string
		upgrades []StateUpgradeDefinition
		wantErr  string
		wantNot  []string
	}{
		{
			name:  "attribute_resolves_ok",
			props: map[string]*Property{"name": {PropertyName: "name"}},
			upgrades: []StateUpgradeDefinition{
				{PriorSchemaVersion: 0, Attributes: map[string]AttributeUpgradeDefinition{"name": {LegacyAttribute: "n"}}},
			},
		},
		{
			name:  "attribute_unknown_errors",
			props: map[string]*Property{"name": {PropertyName: "name"}},
			upgrades: []StateUpgradeDefinition{
				{PriorSchemaVersion: 0, Attributes: map[string]AttributeUpgradeDefinition{
					"notARealProperty": {LegacyAttribute: "legacy_name"},
				}},
			},
			wantErr: "not found in resolved properties",
		},
		{
			name:  "attribute_unknown_but_removed_ok",
			props: map[string]*Property{"name": {PropertyName: "name"}},
			upgrades: []StateUpgradeDefinition{
				{PriorSchemaVersion: 0, Attributes: map[string]AttributeUpgradeDefinition{
					"gone": {LegacyAttribute: "old_gone", LegacyType: StringAttribute, LegacyRestriction: Optional, LegacyStatus: Removed},
				}},
			},
		},
		{
			name:     "child_resolves_ok",
			children: []string{"fvRsBd"},
			upgrades: []StateUpgradeDefinition{
				{PriorSchemaVersion: 0, Children: map[string]AttributeUpgradeDefinition{
					"fvRsBd": {LegacyAttribute: "renamed"},
				}},
			},
		},
		{
			name:     "child_unknown_errors",
			children: []string{"fvRsBd"},
			upgrades: []StateUpgradeDefinition{
				{PriorSchemaVersion: 0, Children: map[string]AttributeUpgradeDefinition{
					"notAChildClass": {LegacyAttribute: "x"},
				}},
			},
			wantErr: "not found in resolved children",
		},
		{
			name:     "child_unknown_but_removed_ok",
			children: []string{"fvRsBd"},
			upgrades: []StateUpgradeDefinition{
				{PriorSchemaVersion: 0, Children: map[string]AttributeUpgradeDefinition{
					"fvRsCustQosPol": {LegacyAttribute: "old_qos", LegacyType: StringAttribute, LegacyRestriction: Optional, LegacyStatus: Removed},
				}},
			},
		},
		{
			name:     "mixed_resolved_unknown_removed_only_unknown_errors",
			children: []string{"fvRsBd"},
			upgrades: []StateUpgradeDefinition{
				{PriorSchemaVersion: 0, Children: map[string]AttributeUpgradeDefinition{
					"fvRsBd":         {LegacyAttribute: "renamed"},
					"fvRsDomAtt":     {LegacyAttribute: "relation_to_domain"},
					"fvRsCustQosPol": {LegacyAttribute: "old_qos", LegacyType: StringAttribute, LegacyRestriction: Optional, LegacyStatus: Removed},
				}},
			},
			wantErr: `child "fvRsDomAtt" not found in resolved children`,
			wantNot: []string{`"fvRsBd"`, `"fvRsCustQosPol"`},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			props := tc.props
			if props == nil {
				props = map[string]*Property{}
			}
			class := classWithProperties("fvAEPg", props, tc.children)
			class.ClassDefinition.StateUpgrades = tc.upgrades

			ds := &DataStore{ctx: NewContext()}
			class.setStateUpgrades(ds)
			err := ds.ctx.Diagnostics.Error()
			if tc.wantErr == "" {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				return
			}
			if !assert.Error(t, err) {
				return
			}
			msg := err.Error()
			assert.Contains(t, msg, tc.wantErr)
			for _, forbidden := range tc.wantNot {
				assert.NotContains(t, msg, forbidden)
			}
		})
	}
}

// TestValidateStateUpgrades_AccumulatesAllDiagnostics proves the single-pass
// policy: no check short-circuits on the first failure, so a definition that
// violates multiple rules surfaces every violation in one diagnostics summary.
func TestValidateStateUpgrades_AccumulatesAllDiagnostics(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	class := classWithProperties("fvAEPg", map[string]*Property{"name": {PropertyName: "name"}}, []string{"fvRsBd"})
	class.ClassDefinition.StateUpgrades = []StateUpgradeDefinition{
		{
			PriorSchemaVersion: 0,
			Attributes: map[string]AttributeUpgradeDefinition{
				"notARealProperty": {LegacyAttribute: "legacy_name"}, // triggers exhaustiveness error
				"gone":             {LegacyStatus: Removed},          // triggers Removed-missing-fields error
			},
		},
		{
			PriorSchemaVersion: 0, // duplicate of the entry above
			Attributes: map[string]AttributeUpgradeDefinition{
				"name": {LegacyAttribute: "n"},
			},
		},
	}

	ds := &DataStore{ctx: NewContext()}
	class.setStateUpgrades(ds)
	err := ds.ctx.Diagnostics.Error()
	if !assert.Error(t, err) {
		return
	}
	msg := err.Error()
	assert.Contains(t, msg, "duplicate prior_schema_version 0", "uniqueness rule must surface")
	assert.Contains(t, msg, "not found in resolved properties", "exhaustiveness rule must surface in the same pass")
	assert.Contains(t, msg, "legacy_status: removed requires legacy_attribute", "Removed-required-fields rule must surface in the same pass")
}

// TestValidateStateUpgradeEntry_LegacyAttributeCollisions covers the per-entry
// collision detector that flags two attributes sharing the same legacy_attribute
// inside one prior_schema_version. The detector is scoped to a single entry, so
// reusing the same legacy_attribute across distinct prior_schema_version entries
// is allowed; empty legacy_attribute values are not tracked.
func TestValidateStateUpgradeEntry_LegacyAttributeCollisions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	cases := []struct {
		name     string
		upgrades []StateUpgradeDefinition
		wantErr  string
	}{
		{
			name: "same_entry_duplicate_errors",
			upgrades: []StateUpgradeDefinition{
				{
					PriorSchemaVersion: 0,
					Attributes: map[string]AttributeUpgradeDefinition{
						"a": {LegacyAttribute: "shared_legacy"},
						"b": {LegacyAttribute: "shared_legacy"},
					},
				},
			},
			wantErr: "duplicate legacy_attribute",
		},
		{
			name: "cross_entry_repeat_ok",
			upgrades: []StateUpgradeDefinition{
				{
					PriorSchemaVersion: 0,
					Attributes:         map[string]AttributeUpgradeDefinition{"a": {LegacyAttribute: "shared_legacy"}},
				},
				{
					PriorSchemaVersion: 2,
					Attributes:         map[string]AttributeUpgradeDefinition{"b": {LegacyAttribute: "shared_legacy"}},
				},
			},
		},
		{
			name: "empty_legacy_attribute_not_counted",
			upgrades: []StateUpgradeDefinition{
				{
					PriorSchemaVersion: 0,
					Attributes: map[string]AttributeUpgradeDefinition{
						"a": {LegacyType: StringAttribute},
						"b": {LegacyType: BoolAttribute},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			class := classWithProperties("fvCtx", map[string]*Property{
				"a": {PropertyName: "a"},
				"b": {PropertyName: "b"},
			}, nil)
			class.ClassDefinition.StateUpgrades = tc.upgrades

			ds := &DataStore{ctx: NewContext()}
			class.setStateUpgrades(ds)
			err := ds.ctx.Diagnostics.Error()
			if tc.wantErr == "" {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				return
			}
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}

// TestAttributeUpgradeDefinition_Validate_Removed covers the attributes-bucket
// Removed-required-fields rule: a node with LegacyStatus == Removed must supply
// legacy_attribute, legacy_type, and legacy_restriction because there is no
// current attribute to inherit from. The negative controls prove the requirement
// is gated on the Removed status itself — a fully-populated Removed node and a
// non-Removed node missing the same fields both pass.
func TestAttributeUpgradeDefinition_Validate_Removed(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	cases := []struct {
		name    string
		node    AttributeUpgradeDefinition
		wantErr string
	}{
		{
			name: "missing_legacy_attribute",
			node: AttributeUpgradeDefinition{
				LegacyType:        StringAttribute,
				LegacyRestriction: Optional,
				LegacyStatus:      Removed,
			},
			wantErr: "legacy_attribute",
		},
		{
			name: "missing_legacy_type",
			node: AttributeUpgradeDefinition{
				LegacyAttribute:   "old",
				LegacyRestriction: Optional,
				LegacyStatus:      Removed,
			},
			wantErr: "legacy_type",
		},
		{
			name: "missing_legacy_restriction",
			node: AttributeUpgradeDefinition{
				LegacyAttribute: "old",
				LegacyType:      StringAttribute,
				LegacyStatus:    Removed,
			},
			wantErr: "legacy_restriction",
		},
		{
			name: "removed_all_legacy_fields_present_ok",
			node: AttributeUpgradeDefinition{
				LegacyAttribute:   "old",
				LegacyType:        StringAttribute,
				LegacyRestriction: Optional,
				LegacyStatus:      Removed,
			},
		},
		{
			name: "not_removed_missing_fields_ok",
			node: AttributeUpgradeDefinition{LegacyStatus: Functioning},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Use a key that resolves so the exhaustiveness check stays silent
			// regardless of the node's LegacyStatus; only the Removed-required
			// rule (or its absence) drives the assertion.
			class := classWithProperties("fvCtx", map[string]*Property{
				"someProp": {PropertyName: "someProp"},
			}, nil)
			class.ClassDefinition.StateUpgrades = []StateUpgradeDefinition{
				{
					PriorSchemaVersion: 0,
					Attributes:         map[string]AttributeUpgradeDefinition{"someProp": tc.node},
				},
			}

			ds := &DataStore{ctx: NewContext()}
			class.setStateUpgrades(ds)
			err := ds.ctx.Diagnostics.Error()
			if tc.wantErr == "" {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				return
			}
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}

// TestValidateChild_ShapeRules covers the children-bucket shape validator:
// block rename, scalar-wrap via inner attributes, scalar-wrap via inner child
// (the recursive hasInnerLegacyAttribute path), the empty-block error, the
// Removed-without-legacy_attribute error, and recursive propagation of inner
// failures with full path prefixes (children["X"]: attributes["Y"], and
// children["X"]: children["Y"]).
func TestValidateChild_ShapeRules(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	cases := []struct {
		name       string
		childKey   string
		definition AttributeUpgradeDefinition
		wantErr    string
	}{
		{
			name:       "block_rename_only_ok",
			childKey:   "fvRsBd",
			definition: AttributeUpgradeDefinition{LegacyAttribute: "relation_to_bridge_domain"},
		},
		{
			name:     "scalar_wrap_via_inner_attribute_ok",
			childKey: "fvRsNodeAtt",
			definition: AttributeUpgradeDefinition{
				Attributes: map[string]AttributeUpgradeDefinition{
					"tDn":   {LegacyAttribute: "node_dn"},
					"encap": {LegacyAttribute: "node_encap"},
				},
			},
		},
		{
			name:     "scalar_wrap_via_inner_child_ok",
			childKey: "fvSubnet",
			definition: AttributeUpgradeDefinition{
				Children: map[string]AttributeUpgradeDefinition{
					"fvRsBDSubnetToOut": {LegacyAttribute: "deep_legacy"},
				},
			},
		},
		{
			name:       "empty_block_no_inner_errors",
			childKey:   "fvRsBd",
			definition: AttributeUpgradeDefinition{},
			wantErr:    "neither legacy_attribute / legacy_type",
		},
		{
			name:       "removed_block_missing_legacy_attribute_errors",
			childKey:   "fvRsBd",
			definition: AttributeUpgradeDefinition{LegacyStatus: Removed},
			wantErr:    "removed requires legacy_attribute",
		},
		{
			name:     "removed_block_with_legacy_attribute_ok",
			childKey: "fvRsBd",
			definition: AttributeUpgradeDefinition{
				LegacyAttribute: "old",
				LegacyStatus:    Removed,
			},
		},
		{
			name:     "inner_attribute_failure_propagates_with_path",
			childKey: "fvRsBd",
			definition: AttributeUpgradeDefinition{
				LegacyAttribute: "renamed",
				Attributes: map[string]AttributeUpgradeDefinition{
					"foo": {LegacyStatus: Removed}, // missing required fields
				},
			},
			wantErr: `children["fvRsBd"]: attributes["foo"]`,
		},
		{
			name:     "inner_child_failure_propagates_with_path",
			childKey: "fvRsBd",
			definition: AttributeUpgradeDefinition{
				LegacyAttribute: "renamed",
				Children: map[string]AttributeUpgradeDefinition{
					"fvSubnet": {}, // empty inner block triggers the shape error
				},
			},
			wantErr: `children["fvRsBd"]: children["fvSubnet"]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			class := classWithProperties("fvAEPg", map[string]*Property{}, []string{tc.childKey})
			class.ClassDefinition.StateUpgrades = []StateUpgradeDefinition{
				{
					PriorSchemaVersion: 0,
					Children:           map[string]AttributeUpgradeDefinition{tc.childKey: tc.definition},
				},
			}

			ds := &DataStore{ctx: NewContext()}
			class.setStateUpgrades(ds)
			err := ds.ctx.Diagnostics.Error()
			if tc.wantErr == "" {
				assert.NoError(t, err, test.MessageUnexpectedError(err))
				return
			}
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tc.wantErr)
			}
		})
	}
}

func TestHasChild(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	class := &Class{
		Name: testClassName("fvAEPg"),
		Children: []*ClassName{
			testClassName("fvRsBd"),
			testClassName("fvRsDomAtt"),
		},
	}

	assert.True(t, class.hasChild("fvRsBd"), "exact match on first child")
	assert.True(t, class.hasChild("fvRsDomAtt"), "exact match on second child")
	assert.False(t, class.hasChild("fvRsCustQosPol"), "absent child returns false")
	assert.False(t, class.hasChild(""), "empty needle returns false")

	empty := &Class{Name: testClassName("fvAEPg")}
	assert.False(t, empty.hasChild("fvRsBd"), "class with no children returns false without panicking")
}

// TestSetPropertyStateUpgradeValues_Distribution covers the per-version
// per-Property distribution of the class-level StateUpgrades tree into
// Property.StateUpgradeValues maps: single-version write, multi-version
// accumulation under one PropertyName, untouched-property nil-map preservation,
// silent skip for keys with no resolved Property (LegacyStatus == Removed), and
// parentDn (the synthetic parent-DN attribute) handled the same as any
// meta-derived property name.
func TestSetPropertyStateUpgradeValues_Distribution(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	cases := []struct {
		name   string
		setup  func() *Class
		assert func(t *testing.T, class *Class)
	}{
		{
			name: "single_version_single_property",
			setup: func() *Class {
				class := classWithProperties("fvCtx", map[string]*Property{
					"name": {PropertyName: "name", AttributeName: "name", Required: true, ValueType: String},
				}, nil)
				class.ClassDefinition.StateUpgrades = []StateUpgradeDefinition{
					{
						PriorSchemaVersion: 0,
						Attributes:         map[string]AttributeUpgradeDefinition{"name": {LegacyAttribute: "old_name"}},
					},
				}
				return class
			},
			assert: func(t *testing.T, class *Class) {
				name := class.Properties["name"]
				if assert.Len(t, name.StateUpgradeValues, 1) {
					assert.Equal(t, "old_name", name.StateUpgradeValues[0].AttributeName)
				}
			},
		},
		{
			name: "same_property_two_versions_accumulates",
			setup: func() *Class {
				class := classWithProperties("fvCtx", map[string]*Property{
					"name": {PropertyName: "name", AttributeName: "name", Required: true, ValueType: String},
				}, nil)
				class.ClassDefinition.StateUpgrades = []StateUpgradeDefinition{
					{
						PriorSchemaVersion: 0,
						Attributes:         map[string]AttributeUpgradeDefinition{"name": {LegacyAttribute: "v0_name"}},
					},
					{
						PriorSchemaVersion: 2,
						Attributes:         map[string]AttributeUpgradeDefinition{"name": {LegacyAttribute: "v2_name"}},
					},
				}
				return class
			},
			assert: func(t *testing.T, class *Class) {
				name := class.Properties["name"]
				if assert.Len(t, name.StateUpgradeValues, 2) {
					assert.Equal(t, "v0_name", name.StateUpgradeValues[0].AttributeName)
					assert.Equal(t, "v2_name", name.StateUpgradeValues[2].AttributeName)
				}
			},
		},
		{
			name: "untouched_property_keeps_nil_map",
			setup: func() *Class {
				class := classWithProperties("fvCtx", map[string]*Property{
					"name":  {PropertyName: "name", AttributeName: "name", Required: true, ValueType: String},
					"descr": {PropertyName: "descr", AttributeName: "description", Optional: true, ValueType: String},
				}, nil)
				class.ClassDefinition.StateUpgrades = []StateUpgradeDefinition{
					{
						PriorSchemaVersion: 0,
						Attributes:         map[string]AttributeUpgradeDefinition{"name": {LegacyAttribute: "old_name"}},
					},
				}
				return class
			},
			assert: func(t *testing.T, class *Class) {
				assert.NotNil(t, class.Properties["name"].StateUpgradeValues)
				assert.Nil(t, class.Properties["descr"].StateUpgradeValues, "property the upgrade tree never names keeps a nil map")
			},
		},
		{
			name: "removed_key_with_no_property_silent_skip",
			setup: func() *Class {
				class := classWithProperties("fvCtx", map[string]*Property{}, nil)
				class.ClassDefinition.StateUpgrades = []StateUpgradeDefinition{
					{
						PriorSchemaVersion: 0,
						Attributes: map[string]AttributeUpgradeDefinition{
							"oldAttr": {
								LegacyAttribute:   "old_attr",
								LegacyType:        StringAttribute,
								LegacyRestriction: Optional,
								LegacyStatus:      Removed,
							},
						},
					},
				}
				return class
			},
			assert: func(t *testing.T, class *Class) {
				assert.Empty(t, class.Properties, "no Property exists for the removed key, nothing to distribute")
			},
		},
		{
			name: "parentDn_distribution",
			setup: func() *Class {
				class := classWithProperties("fvCtx", map[string]*Property{
					"parentDn": {PropertyName: "parentDn", AttributeName: "parent_dn", Required: true, ValueType: String},
				}, nil)
				class.ClassDefinition.StateUpgrades = []StateUpgradeDefinition{
					{
						PriorSchemaVersion: 0,
						Attributes:         map[string]AttributeUpgradeDefinition{"parentDn": {LegacyAttribute: "old_parent_dn"}},
					},
				}
				return class
			},
			assert: func(t *testing.T, class *Class) {
				parentDn := class.Properties["parentDn"]
				if assert.Contains(t, parentDn.StateUpgradeValues, 0) {
					assert.Equal(t, "old_parent_dn", parentDn.StateUpgradeValues[0].AttributeName)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			class := tc.setup()
			ds := &DataStore{ctx: NewContext()}
			class.setStateUpgrades(ds)
			err := ds.ctx.Diagnostics.Error()
			assert.NoError(t, err, test.MessageUnexpectedError(err))
			class.setPropertyStateUpgradeValues()
			tc.assert(t, class)
		})
	}
}

// TestBuildStateUpgradeValue_Overlays covers the seed-then-overlay semantics of
// buildStateUpgradeValue: the StateUpgradeValue starts as a copy of the current
// Property's AttributeName / Type / Required-Optional-Computed triplet, then each
// legacy_* override (when non-zero) rewrites the matching field. Cases prove
// every overlay direction, the no-op (zero overlay leaves the seed in place),
// every collapse arm of legacyTypeToValueType the overlay can reach (String /
// Set / Object), every arm of the restriction-triplet rewrite, and a combined
// case proving the three overlays are independent.
func TestBuildStateUpgradeValue_Overlays(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	cases := []struct {
		name     string
		seed     *Property
		overlay  AttributeUpgradeDefinition
		expected StateUpgradeValue
	}{
		{
			name:     "no_overlays_seed_carries_through",
			seed:     &Property{PropertyName: "name", AttributeName: "name", Optional: true, Computed: true, ValueType: String},
			expected: StateUpgradeValue{AttributeName: "name", Optional: true, Computed: true, Type: String},
		},
		{
			name:     "legacy_attribute_rewrites_name",
			seed:     &Property{PropertyName: "name", AttributeName: "name", Required: true, ValueType: String},
			overlay:  AttributeUpgradeDefinition{LegacyAttribute: "old_name"},
			expected: StateUpgradeValue{AttributeName: "old_name", Required: true, Type: String},
		},
		{
			name:     "legacy_type_string_to_set",
			seed:     &Property{PropertyName: "tags", AttributeName: "tags", Optional: true, ValueType: String},
			overlay:  AttributeUpgradeDefinition{LegacyType: SetAttribute},
			expected: StateUpgradeValue{AttributeName: "tags", Optional: true, Type: Set},
		},
		{
			name:     "legacy_type_set_to_string",
			seed:     &Property{PropertyName: "tags", AttributeName: "tags", Optional: true, ValueType: Set},
			overlay:  AttributeUpgradeDefinition{LegacyType: StringAttribute},
			expected: StateUpgradeValue{AttributeName: "tags", Optional: true, Type: String},
		},
		{
			name:     "legacy_type_map_to_object",
			seed:     &Property{PropertyName: "labels", AttributeName: "labels", Optional: true, ValueType: String},
			overlay:  AttributeUpgradeDefinition{LegacyType: MapAttribute},
			expected: StateUpgradeValue{AttributeName: "labels", Optional: true, Type: Object},
		},
		{
			name:     "legacy_type_single_nested_to_object",
			seed:     &Property{PropertyName: "spec", AttributeName: "spec", Optional: true, ValueType: String},
			overlay:  AttributeUpgradeDefinition{LegacyType: SingleNestedAttribute},
			expected: StateUpgradeValue{AttributeName: "spec", Optional: true, Type: Object},
		},
		{
			name:     "legacy_type_zero_keeps_seed_set",
			seed:     &Property{PropertyName: "tags", AttributeName: "tags", Optional: true, ValueType: Set},
			overlay:  AttributeUpgradeDefinition{LegacyType: UndefinedLegacyAttributeType},
			expected: StateUpgradeValue{AttributeName: "tags", Optional: true, Type: Set},
		},
		{
			name:     "legacy_restriction_required",
			seed:     &Property{PropertyName: "name", AttributeName: "name", Optional: true, Computed: true, ValueType: String},
			overlay:  AttributeUpgradeDefinition{LegacyRestriction: Required},
			expected: StateUpgradeValue{AttributeName: "name", Required: true, Type: String},
		},
		{
			name:     "legacy_restriction_optional",
			seed:     &Property{PropertyName: "name", AttributeName: "name", Required: true, ValueType: String},
			overlay:  AttributeUpgradeDefinition{LegacyRestriction: Optional},
			expected: StateUpgradeValue{AttributeName: "name", Optional: true, Computed: true, Type: String},
		},
		{
			name:     "legacy_restriction_read_only",
			seed:     &Property{PropertyName: "name", AttributeName: "name", Required: true, ValueType: String},
			overlay:  AttributeUpgradeDefinition{LegacyRestriction: ReadOnly},
			expected: StateUpgradeValue{AttributeName: "name", Computed: true, Type: String},
		},
		{
			name: "all_three_overlays_combined",
			seed: &Property{PropertyName: "x", AttributeName: "x", Required: true, ValueType: String},
			overlay: AttributeUpgradeDefinition{
				LegacyAttribute:   "old_x",
				LegacyType:        SetAttribute,
				LegacyRestriction: Optional,
			},
			expected: StateUpgradeValue{AttributeName: "old_x", Optional: true, Computed: true, Type: Set},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, buildStateUpgradeValue(tc.seed, tc.overlay))
		})
	}
}
