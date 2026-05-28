package data

import (
	"bytes"
	"os"
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
						Type:      "explicit",
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
						Type:      "named",
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
				ErrorMsg: "undefined relationship type",
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
						Type:      "named",
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

			err := class.setParents()
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

	err := class.setParents()
	assert.NoError(t, err, test.MessageUnexpectedError(err))

	// Verify the warning was logged.
	logOutput := logBuffer.String()
	expectedWarning := "WARN: Parent class 'fvTenant' is defined in both IncludeParents and ExcludeParents for class 'testClass'. IncludeParents takes precedence."
	assert.Contains(t, logOutput, expectedWarning, test.MessageEqual(expectedWarning, logOutput, "warning log message"))
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

			err := class.setParents()
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
						"pcTag": {Restriction: "read_only"},
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
						"annotation": {Restriction: "exclude"},
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
						"value": {Restriction: "required"},
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
						"scope": {Restriction: "read_only"},
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
						"userdom": {Restriction: "optional"},
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
