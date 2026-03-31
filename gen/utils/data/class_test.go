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
	class.MetaFileContent = map[string]interface{}{
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
	class.MetaFileContent = map[string]interface{}{
		"label": "Private Network",
		"relationInfo": map[string]interface{}{
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
	class.MetaFileContent = map[string]interface{}{
		"label": "contract",
		"relationInfo": map[string]interface{}{
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
	class.MetaFileContent = map[string]interface{}{
		"label": "netflow exporter",
		"relationInfo": map[string]interface{}{
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
	class := Class{}
	class.MetaFileContent = map[string]interface{}{"label": ""}

	err := class.setResourceName(ds)

	assert.Error(t, err)
}

func TestSetResourceNameFromNoLabelError(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	class := Class{}
	class.MetaFileContent = map[string]interface{}{"no_label": ""}

	err := class.setResourceName(ds)

	assert.Error(t, err)
}

type setRelationInput struct {
	ClassName       string
	MetaFileContent map[string]interface{}
}

type setRelationExpected struct {
	FromClass       string
	ToClass         string
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
				FromClass:       "",
				ToClass:         "",
				Type:            RelationshipTypeEnum(0),
				IncludeFrom:     false,
				RelationalClass: false,
			},
		},
		{
			Name: "test_include_from_false_type_named",
			Input: setRelationInput{
				ClassName: "fvRsCons",
				MetaFileContent: map[string]interface{}{
					"label": "contract",
					"relationInfo": map[string]interface{}{
						"type":   "named",
						"fromMo": "fv:EPg",
						"toMo":   "vz:BrCP",
					},
				},
			},
			Expected: setRelationExpected{
				FromClass:       "fvEPg",
				ToClass:         "vzBrCP",
				Type:            Named,
				IncludeFrom:     false,
				RelationalClass: true,
			},
		},
		{
			Name: "test_include_from_true_type_explicit",
			Input: setRelationInput{
				ClassName: "netflowRsExporterToCtx",
				MetaFileContent: map[string]interface{}{
					"label": "netflow exporter",
					"relationInfo": map[string]interface{}{
						"type":   "explicit",
						"fromMo": "netflow:AExporterPol",
						"toMo":   "fv:Ctx",
					},
					"isCreatableDeletable": "always",
				},
			},
			Expected: setRelationExpected{
				FromClass:       "netflowAExporterPol",
				ToClass:         "fvCtx",
				Type:            Explicit,
				IncludeFrom:     true,
				RelationalClass: true,
			},
		},
		{
			Name: "test_undefined_type_error",
			Input: setRelationInput{
				ClassName: "netflowRsExporterToCtx",
				MetaFileContent: map[string]interface{}{
					"relationInfo": map[string]interface{}{
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setRelationInput)
			expected := testCase.Expected.(setRelationExpected)
			class := Class{Name: testClassName(input.ClassName)}
			class.MetaFileContent = input.MetaFileContent

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
			assert.Equal(t, expected.ToClass, class.Relation.ToClass, test.MessageEqual(expected.ToClass, class.Relation.ToClass, testCase.Name))
			assert.Equal(t, expected.Type, class.Relation.Type, test.MessageEqual(expected.Type, class.Relation.Type, testCase.Name))
			assert.Equal(t, expected.IncludeFrom, class.Relation.IncludeFrom, test.MessageEqual(expected.IncludeFrom, class.Relation.IncludeFrom, testCase.Name))
			assert.Equal(t, expected.RelationalClass, class.Relation.RelationalClass, test.MessageEqual(expected.RelationalClass, class.Relation.RelationalClass, testCase.Name))
		})
	}
}

type setAllowDeleteInput struct {
	MetaFileContent map[string]interface{}
	ClassDefinition ClassDefinition
}

func TestSetAllowDelete(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_isCreatableDeletable_never",
			Input: setAllowDeleteInput{
				MetaFileContent: map[string]interface{}{"isCreatableDeletable": "never"},
				ClassDefinition: ClassDefinition{},
			},
			Expected: false,
		},
		{
			Name: "test_isCreatableDeletable_always",
			Input: setAllowDeleteInput{
				MetaFileContent: map[string]interface{}{"isCreatableDeletable": "always"},
				ClassDefinition: ClassDefinition{},
			},
			Expected: true,
		},
		{
			Name: "test_definition_allow_delete_never",
			Input: setAllowDeleteInput{
				MetaFileContent: map[string]interface{}{},
				ClassDefinition: ClassDefinition{AllowDelete: "never"},
			},
			Expected: false,
		},
		{
			Name: "test_definition_allow_delete_always",
			Input: setAllowDeleteInput{
				MetaFileContent: map[string]interface{}{},
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
	RnMap                map[string]interface{}
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
				RnMap:                map[string]interface{}{},
			},
			Expected: []string{"fvSubnet"},
		},
		{
			Name: "test_includes_rs_prefixed_classes_from_rnMap",
			Input: setChildrenInput{
				IncludeChildren:      []string{},
				ExcludeChildren:      []string{},
				AlwaysIncludeAsChild: []string{},
				RnMap: map[string]interface{}{
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
				RnMap: map[string]interface{}{
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
				RnMap: map[string]interface{}{
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
				RnMap: map[string]interface{}{
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
				RnMap: map[string]interface{}{
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
				RnMap: map[string]interface{}{
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
				RnMap: map[string]interface{}{
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
				RnMap:                map[string]interface{}{},
			},
			Expected: []string{},
		},
		{
			Name: "test_includeChildren_takes_precedence_over_excludeChildren",
			Input: setChildrenInput{
				IncludeChildren:      []string{"fvSubnet"},
				ExcludeChildren:      []string{"fvSubnet"},
				AlwaysIncludeAsChild: []string{},
				RnMap:                map[string]interface{}{},
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
				MetaFileContent: map[string]interface{}{
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
		MetaFileContent: map[string]interface{}{
			"rnMap": map[string]interface{}{},
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
				RnMap: map[string]interface{}{
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
				RnMap:                map[string]interface{}{},
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
				RnMap:                map[string]interface{}{},
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
				MetaFileContent: map[string]interface{}{
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
	ContainedBy    map[string]interface{}
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
				ContainedBy:    map[string]interface{}{},
			},
			Expected: []string{"fvTenant"},
		},
		{
			Name: "test_includes_all_classes_from_containedBy",
			Input: setParentsInput{
				IncludeParents: []string{},
				ExcludeParents: []string{},
				ContainedBy: map[string]interface{}{
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
				ContainedBy: map[string]interface{}{
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
				ContainedBy: map[string]interface{}{
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
				ContainedBy: map[string]interface{}{
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
				ContainedBy: map[string]interface{}{
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
				ContainedBy: map[string]interface{}{
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
				ContainedBy:    map[string]interface{}{},
			},
			Expected: []string{},
		},
		{
			Name: "test_includeParents_takes_precedence_over_excludeParents",
			Input: setParentsInput{
				IncludeParents: []string{"fvTenant"},
				ExcludeParents: []string{"fvTenant"},
				ContainedBy:    map[string]interface{}{},
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
				MetaFileContent: map[string]interface{}{
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
		MetaFileContent: map[string]interface{}{
			"containedBy": map[string]interface{}{},
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
				ContainedBy: map[string]interface{}{
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
				ContainedBy:    map[string]interface{}{},
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
				ContainedBy:    map[string]interface{}{},
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
				MetaFileContent: map[string]interface{}{
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

type setSupportedVersionsInput struct {
	ClassDefinitionVersions string
	MetaVersions            interface{}
}

type setSupportedVersionsExpected struct {
	Raw      string
	String   string
	Error    bool
	ErrorMsg string
}

func TestSetSupportedVersions(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_versions_from_meta_file",
			Input: setSupportedVersionsInput{
				ClassDefinitionVersions: "",
				MetaVersions:            "4.2(7f)-",
			},
			Expected: setSupportedVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name: "test_versions_from_class_definition_override",
			Input: setSupportedVersionsInput{
				ClassDefinitionVersions: "5.0(1a)-",
				MetaVersions:            "4.2(7f)-",
			},
			Expected: setSupportedVersionsExpected{
				Raw:    "5.0(1a)-",
				String: "5.0(1a) and later",
			},
		},
		{
			Name: "test_versions_from_class_definition_when_meta_empty",
			Input: setSupportedVersionsInput{
				ClassDefinitionVersions: "5.0(1a)-",
				MetaVersions:            "",
			},
			Expected: setSupportedVersionsExpected{
				Raw:    "5.0(1a)-",
				String: "5.0(1a) and later",
			},
		},
		{
			Name: "test_versions_from_class_definition_when_meta_nil",
			Input: setSupportedVersionsInput{
				ClassDefinitionVersions: "5.0(1a)-",
				MetaVersions:            nil,
			},
			Expected: setSupportedVersionsExpected{
				Raw:    "5.0(1a)-",
				String: "5.0(1a) and later",
			},
		},
		{
			Name: "test_multiple_ranges",
			Input: setSupportedVersionsInput{
				ClassDefinitionVersions: "",
				MetaVersions:            "3.2(10e)-3.2(10g),4.2(7f)-",
			},
			Expected: setSupportedVersionsExpected{
				Raw:    "3.2(10e)-3.2(10g),4.2(7f)-",
				String: "3.2(10e) to 3.2(10g), 4.2(7f) and later",
			},
		},
		{
			Name: "test_multiple_ranges_sorted",
			Input: setSupportedVersionsInput{
				ClassDefinitionVersions: "",
				MetaVersions:            "5.2(1g)-,3.2(10e)-3.2(10g)",
			},
			Expected: setSupportedVersionsExpected{
				Raw:    "5.2(1g)-,3.2(10e)-3.2(10g)",
				String: "3.2(10e) to 3.2(10g), 5.2(1g) and later",
			},
		},
		{
			Name: "test_error_empty_versions",
			Input: setSupportedVersionsInput{
				ClassDefinitionVersions: "",
				MetaVersions:            "",
			},
			Expected: setSupportedVersionsExpected{
				Error:    true,
				ErrorMsg: "versions not specified for class 'fvTenant': add versions to the class definition file",
			},
		},
		{
			Name: "test_error_nil_versions",
			Input: setSupportedVersionsInput{
				ClassDefinitionVersions: "",
				MetaVersions:            nil,
			},
			Expected: setSupportedVersionsExpected{
				Error:    true,
				ErrorMsg: "versions not specified for class 'fvTenant': add versions to the class definition file",
			},
		},
		{
			Name: "test_error_invalid_version",
			Input: setSupportedVersionsInput{
				ClassDefinitionVersions: "",
				MetaVersions:            "invalid",
			},
			Expected: setSupportedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse versions for class 'fvTenant': invalid version 'invalid': unknown",
			},
		},
		{
			Name: "test_error_invalid_version_in_range",
			Input: setSupportedVersionsInput{
				ClassDefinitionVersions: "",
				MetaVersions:            "4.2(7f)-,invalid",
			},
			Expected: setSupportedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse versions for class 'fvTenant': invalid version 'invalid': unknown",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setSupportedVersionsInput)
			expected := testCase.Expected.(setSupportedVersionsExpected)

			class := Class{
				Name: testClassName("fvTenant"),
				ClassDefinition: ClassDefinition{
					SupportedVersions: input.ClassDefinitionVersions,
				},
				MetaFileContent: map[string]interface{}{},
			}
			if input.MetaVersions != nil {
				class.MetaFileContent["versions"] = input.MetaVersions
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

type setDeprecatedVersionsInput struct {
	ClassDefinitionVersions string
}

type setDeprecatedVersionsExpected struct {
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
			Input: setDeprecatedVersionsInput{
				ClassDefinitionVersions: "",
			},
			Expected: setDeprecatedVersionsExpected{
				Nil: true,
			},
		},
		{
			Name: "test_deprecated_versions_single_range",
			Input: setDeprecatedVersionsInput{
				ClassDefinitionVersions: "4.2(7f)-",
			},
			Expected: setDeprecatedVersionsExpected{
				Raw:    "4.2(7f)-",
				String: "4.2(7f) and later",
			},
		},
		{
			Name: "test_deprecated_versions_bounded_range",
			Input: setDeprecatedVersionsInput{
				ClassDefinitionVersions: "3.2(10e)-4.2(7f)",
			},
			Expected: setDeprecatedVersionsExpected{
				Raw:    "3.2(10e)-4.2(7f)",
				String: "3.2(10e) to 4.2(7f)",
			},
		},
		{
			Name: "test_deprecated_versions_multiple_ranges",
			Input: setDeprecatedVersionsInput{
				ClassDefinitionVersions: "3.2(10e)-3.2(10g),4.2(7f)-",
			},
			Expected: setDeprecatedVersionsExpected{
				Raw:    "3.2(10e)-3.2(10g),4.2(7f)-",
				String: "3.2(10e) to 3.2(10g), 4.2(7f) and later",
			},
		},
		{
			Name: "test_error_invalid_deprecated_version",
			Input: setDeprecatedVersionsInput{
				ClassDefinitionVersions: "invalid",
			},
			Expected: setDeprecatedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse deprecated versions for class 'fvTenant': invalid version 'invalid': unknown",
			},
		},
		{
			Name: "test_error_invalid_deprecated_version_in_range",
			Input: setDeprecatedVersionsInput{
				ClassDefinitionVersions: "4.2(7f)-,invalid",
			},
			Expected: setDeprecatedVersionsExpected{
				Error:    true,
				ErrorMsg: "failed to parse deprecated versions for class 'fvTenant': invalid version 'invalid': unknown",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(setDeprecatedVersionsInput)
			expected := testCase.Expected.(setDeprecatedVersionsExpected)

			class := Class{
				Name: testClassName("fvTenant"),
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
	MetaIdentifiedBy            interface{}
}

func TestSetIdentifiedBy(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_identified_by_from_meta_file",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: nil,
				MetaIdentifiedBy:            []interface{}{"name"},
			},
			Expected: []string{"name"},
		},
		{
			Name: "test_identified_by_from_meta_file_multiple",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: nil,
				MetaIdentifiedBy:            []interface{}{"key", "value"},
			},
			Expected: []string{"key", "value"},
		},
		{
			Name: "test_identified_by_from_meta_file_sorted",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: nil,
				MetaIdentifiedBy:            []interface{}{"value", "key"},
			},
			Expected: []string{"key", "value"},
		},
		{
			Name: "test_identified_by_from_meta_file_deduplicated",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: nil,
				MetaIdentifiedBy:            []interface{}{"name", "name"},
			},
			Expected: []string{"name"},
		},
		{
			Name: "test_class_definition_overrides_meta_file",
			Input: setIdentifiedByInput{
				ClassDefinitionIdentifiedBy: []string{"overrideKey"},
				MetaIdentifiedBy:            []interface{}{"metaKey"},
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
				MetaIdentifiedBy:            []interface{}{"metaKey"},
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
				MetaIdentifiedBy:            []interface{}{},
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
				MetaFileContent: map[string]interface{}{},
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

func TestSetDeprecated(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:     "test_deprecated_false",
			Input:    false,
			Expected: false,
		},
		{
			Name:     "test_deprecated_true",
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
					Deprecated: testCase.Input.(bool),
				},
			}

			class.setDeprecated()

			assert.Equal(t, testCase.Expected, class.Deprecated, test.MessageEqual(testCase.Expected, class.Deprecated, testCase.Name))
		})
	}
}

type setPlatformTypeInput struct {
	ClassDefinitionPlatformType string
	PlatformFlavors             interface{}
}

func TestSetPlatformType(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name: "test_apic_only_from_meta",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []interface{}{"apic"},
			},
			Expected: Apic,
		},
		{
			Name: "test_cloud_only_from_meta",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []interface{}{"capic"},
			},
			Expected: Cloud,
		},
		{
			Name: "test_both_from_meta",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []interface{}{"apic", "capic"},
			},
			Expected: Both,
		},
		{
			Name: "test_both_from_meta_reverse_order",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []interface{}{"capic", "apic"},
			},
			Expected: Both,
		},
		{
			Name: "test_definition_overrides_meta_apic",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "apic",
				PlatformFlavors:             []interface{}{"capic"},
			},
			Expected: Apic,
		},
		{
			Name: "test_definition_overrides_meta_cloud",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "cloud",
				PlatformFlavors:             []interface{}{"apic"},
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
				PlatformFlavors:             []interface{}{},
			},
			Expected: Apic,
		},
		{
			Name: "test_unknown_flavor_defaults_to_apic",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []interface{}{"unknown"},
			},
			Expected: Apic,
		},
		{
			Name: "test_unknown_flavor_with_apic",
			Input: setPlatformTypeInput{
				ClassDefinitionPlatformType: "",
				PlatformFlavors:             []interface{}{"apic", "unknown"},
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
				MetaFileContent: map[string]interface{}{},
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
		MetaFileContent: map[string]interface{}{
			"platformFlavors": []interface{}{"unknownFlavor"},
		},
	}

	class.setPlatformType()

	// Verify the warning was logged.
	logOutput := logBuffer.String()
	expectedWarning := "WARN: Unknown platform flavor 'unknownFlavor' found for class 'fvTenant'."
	assert.Contains(t, logOutput, expectedWarning, test.MessageEqual(expectedWarning, logOutput, "warning log message"))
}
