package data

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitClassNameToPackageNameAndShortName(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	tests := []struct {
		name                string
		className           string
		expectedPackageName string
		expectedShortName   string
		expectError         bool
		expectedErrorMsg    string
	}{
		{
			name:                "single_word",
			className:           "fvTenant",
			expectedPackageName: "fv",
			expectedShortName:   "Tenant",
			expectError:         false,
		},
		{
			name:                "multiple_words",
			className:           "fvRsIpslaMonPol",
			expectedPackageName: "fv",
			expectedShortName:   "RsIpslaMonPol",
			expectError:         false,
		},
		{
			name:                "error_no_uppercase",
			className:           "error",
			expectedPackageName: "",
			expectedShortName:   "",
			expectError:         true,
			expectedErrorMsg:    "failed to split class name",
		},
		{
			name:                "empty_string",
			className:           "",
			expectedPackageName: "",
			expectedShortName:   "",
			expectError:         true,
			expectedErrorMsg:    "failed to split class name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			packageName, shortName, err := splitClassNameToPackageNameAndShortName(tc.className)

			if tc.expectError {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedErrorMsg)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.expectedPackageName, packageName)
			assert.Equal(t, tc.expectedShortName, shortName)
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

	require.NoError(t, err)
	assert.Equal(t, "annotation", class.ResourceName)
	assert.Equal(t, "annotations", class.ResourceNameNested)
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
	require.NoError(t, err)

	err = class.setResourceName(ds)
	require.NoError(t, err)
	assert.Equal(t, "relation_to_vrf", class.ResourceName)
	assert.Equal(t, "relation_to_vrf", class.ResourceNameNested)
}

func TestSetResourceNameToRelation(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	ds.GlobalMetaDefinition = GlobalMetaDefinition{
		NoMetaFile: test.GlobalMetaDefinitionContract(),
	}
	class := Class{ClassName: "fvRsCons", IdentifiedBy: []string{"tnVzBrCPName"}}
	class.MetaFileContent = test.RelationInfoFvRsCons

	err := class.setRelation()
	require.NoError(t, err)

	err = class.setResourceName(ds)
	require.NoError(t, err)
	assert.Equal(t, "relation_to_contract", class.ResourceName)
	assert.Equal(t, "relation_to_contracts", class.ResourceNameNested)
}

func TestSetResourceNameFromToRelation(t *testing.T) {
	t.Parallel()
	ds := initializeDataStoreTest(t)
	ds.GlobalMetaDefinition = GlobalMetaDefinition{
		NoMetaFile: test.GlobalMetaDefinitionNetflow(),
	}
	class := Class{ClassName: "netflowRsExporterToCtx"}
	class.MetaFileContent = test.RelationInfoNetflowRsExporterToCtx

	err := class.setRelation()
	require.NoError(t, err)

	err = class.setResourceName(ds)
	require.NoError(t, err)
	assert.Equal(t, "relation_from_netflow_exporter_policy_to_vrf", class.ResourceName)
	assert.Equal(t, "relation_to_vrf", class.ResourceNameNested)
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

func TestSetRelation(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	tests := []struct {
		name                    string
		className               string
		metaFileContent         map[string]interface{}
		expectedFromClass       string
		expectedToClass         string
		expectedType            RelationshipTypeEnum
		expectedIncludeFrom     bool
		expectedRelationalClass bool
		expectError             bool
		expectedErrorMsg        string
	}{
		{
			name:                    "no_relation",
			className:               "fvTenant",
			metaFileContent:         nil,
			expectedFromClass:       "",
			expectedToClass:         "",
			expectedType:            RelationshipTypeEnum(0),
			expectedIncludeFrom:     false,
			expectedRelationalClass: false,
			expectError:             false,
		},
		{
			name:                    "include_from_false_type_named",
			className:               "fvRsCons",
			metaFileContent:         test.RelationInfoFvRsCons,
			expectedFromClass:       "fvEPg",
			expectedToClass:         "vzBrCP",
			expectedType:            Named,
			expectedIncludeFrom:     false,
			expectedRelationalClass: true,
			expectError:             false,
		},
		{
			name:                    "include_from_true_type_explicit",
			className:               "netflowRsExporterToCtx",
			metaFileContent:         test.RelationInfoNetflowRsExporterToCtx,
			expectedFromClass:       "netflowAExporterPol",
			expectedToClass:         "fvCtx",
			expectedType:            Explicit,
			expectedIncludeFrom:     true,
			expectedRelationalClass: true,
			expectError:             false,
		},
		{
			name:      "undefined_type_error",
			className: "netflowRsExporterToCtx",
			metaFileContent: map[string]interface{}{
				"relationInfo": map[string]interface{}{
					"type":   "undefinedType",
					"fromMo": "netflow:AExporterPol",
					"toMo":   "fv:Ctx",
				},
			},
			expectError:      true,
			expectedErrorMsg: "undefined relationship type",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			class := Class{ClassName: tc.className}
			class.MetaFileContent = tc.metaFileContent

			err := class.setRelation()

			if tc.expectError {
				require.Error(t, err)
				if tc.expectedErrorMsg != "" {
					assert.ErrorContains(t, err, tc.expectedErrorMsg)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedFromClass, class.Relation.FromClass)
			assert.Equal(t, tc.expectedToClass, class.Relation.ToClass)
			assert.Equal(t, tc.expectedType, class.Relation.Type)
			assert.Equal(t, tc.expectedIncludeFrom, class.Relation.IncludeFrom)
			assert.Equal(t, tc.expectedRelationalClass, class.Relation.RelationalClass)
		})
	}
}

func TestSetAllowDelete(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)

	tests := []struct {
		name            string
		metaFileContent map[string]interface{}
		classDefinition ClassDefinition
		expected        bool
	}{
		{
			name:            "isCreatableDeletable_never",
			metaFileContent: map[string]interface{}{"isCreatableDeletable": "never"},
			classDefinition: ClassDefinition{},
			expected:        false,
		},
		{
			name:            "isCreatableDeletable_always",
			metaFileContent: map[string]interface{}{"isCreatableDeletable": "always"},
			classDefinition: ClassDefinition{},
			expected:        true,
		},
		{
			name:            "definition_allow_delete_never",
			metaFileContent: map[string]interface{}{},
			classDefinition: ClassDefinition{AllowDelete: "never"},
			expected:        false,
		},
		{
			name:            "definition_allow_delete_always",
			metaFileContent: map[string]interface{}{},
			classDefinition: ClassDefinition{AllowDelete: "always"},
			expected:        true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			class := Class{ClassName: "fvPeeringP"}
			class.MetaFileContent = tc.metaFileContent
			class.ClassDefinition = tc.classDefinition
			class.AllowDelete = false

			class.setAllowDelete()

			assert.Equal(t, tc.expected, class.AllowDelete)
		})
	}
}
