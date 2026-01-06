package data

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type splitClassNameExpected struct {
	PackageName string
	ShortName   string
	Error       bool
	ErrorMsg    string
}

func TestSplitClassNameToPackageNameAndShortName(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	testCases := []test.TestCase{
		{
			Name:  "test_single_word",
			Input: "fvTenant",
			Expected: splitClassNameExpected{
				PackageName: "fv",
				ShortName:   "Tenant",
			},
		},
		{
			Name:  "test_multiple_words",
			Input: "fvRsIpslaMonPol",
			Expected: splitClassNameExpected{
				PackageName: "fv",
				ShortName:   "RsIpslaMonPol",
			},
		},
		{
			Name:  "test_error_no_uppercase",
			Input: "error",
			Expected: splitClassNameExpected{
				Error:    true,
				ErrorMsg: "failed to split class name",
			},
		},
		{
			Name:  "test_empty_string",
			Input: "",
			Expected: splitClassNameExpected{
				Error:    true,
				ErrorMsg: "failed to split class name",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			expected := testCase.Expected.(splitClassNameExpected)
			packageName, shortName, err := splitClassNameToPackageNameAndShortName(testCase.Input.(string))

			if expected.Error {
				require.Error(t, err)
				assert.ErrorContains(t, err, expected.ErrorMsg)
			} else {
				require.NoError(t, err, test.MessageUnexpectedError(err))
			}
			assert.Equal(t, expected.PackageName, packageName, test.MessageEqual(expected.PackageName, packageName, testCase.Name))
			assert.Equal(t, expected.ShortName, shortName, test.MessageEqual(expected.ShortName, shortName, testCase.Name))
		})
	}
}

func TestSetResourceNameFromLabelNoRelationWithIdentifier(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	class := Class{ClassName: "tagAnnotation", IdentifiedBy: []string{"key"}}
	class.MetaFileContent = map[string]interface{}{
		"label": "annotation",
	}

	err := class.setResourceName(ds)

	require.NoError(t, err, test.MessageUnexpectedError(err))
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
	class := Class{ClassName: "fvRsScope"}
	class.MetaFileContent = map[string]interface{}{
		"label": "Private Network",
		"relationInfo": map[string]interface{}{
			"type":   "named",
			"fromMo": "fv:ESg",
			"toMo":   "fv:Ctx",
		},
	}

	err := class.setRelation()
	require.NoError(t, err, test.MessageUnexpectedError(err))

	err = class.setResourceName(ds)
	require.NoError(t, err, test.MessageUnexpectedError(err))
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
	class := Class{ClassName: "fvRsCons", IdentifiedBy: []string{"tnVzBrCPName"}}
	class.MetaFileContent = map[string]interface{}{
		"label": "contract",
		"relationInfo": map[string]interface{}{
			"type":   "named",
			"fromMo": "fv:EPg",
			"toMo":   "vz:BrCP",
		},
	}

	err := class.setRelation()
	require.NoError(t, err, test.MessageUnexpectedError(err))

	err = class.setResourceName(ds)
	require.NoError(t, err, test.MessageUnexpectedError(err))
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
	class := Class{ClassName: "netflowRsExporterToCtx"}
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
	require.NoError(t, err, test.MessageUnexpectedError(err))

	err = class.setResourceName(ds)
	require.NoError(t, err, test.MessageUnexpectedError(err))
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
			class := Class{ClassName: input.ClassName}
			class.MetaFileContent = input.MetaFileContent

			err := class.setRelation()

			if expected.Error {
				require.Error(t, err)
				if expected.ErrorMsg != "" {
					assert.ErrorContains(t, err, expected.ErrorMsg)
				}
				return
			}

			require.NoError(t, err, test.MessageUnexpectedError(err))
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
			class := Class{ClassName: "fvPeeringP"}
			class.MetaFileContent = input.MetaFileContent
			class.ClassDefinition = input.ClassDefinition
			class.AllowDelete = false

			class.setAllowDelete()

			assert.Equal(t, testCase.Expected, class.AllowDelete, test.MessageEqual(testCase.Expected, class.AllowDelete, testCase.Name))
		})
	}
}
