package data

import (
	"fmt"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

var relationInfoFvRsCons = map[string]interface{}{
	"label": "contract",
	"relationInfo": map[string]interface{}{
		"type":   "named",
		"fromMo": "fv:EPg",
		"toMo":   "vz:BrCP",
	},
}

var relationInfoNetflowRsExporterToCtx = map[string]interface{}{
	"label": "netflow exporter",
	"relationInfo": map[string]interface{}{
		"type":   "explicit",
		"fromMo": "netflow:AExporterPol",
		"toMo":   "fv:Ctx",
	},
	"isCreatableDeletable": "always",
}

func TestSplitClassNameToPackageNameAndShortNameSingle(t *testing.T) {
	test.InitializeTest(t)
	packageName, shortName, err := splitClassNameToPackageNameAndShortName("fvTenant")
	assert.Equal(t, packageName, "fv", fmt.Sprintf("Expected package name to be 'fv', but got '%s'", packageName))
	assert.Equal(t, shortName, "Tenant", fmt.Sprintf("Expected short name to be 'Tenant', but got '%s'", shortName))
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
}

func TestSplitClassNameToPackageNameAndShortNameMultiple(t *testing.T) {
	test.InitializeTest(t)
	packageName, shortName, err := splitClassNameToPackageNameAndShortName("fvRsIpslaMonPol")
	assert.Equal(t, packageName, "fv", fmt.Sprintf("Expected package name to be 'fv', but got '%s'", packageName))
	assert.Equal(t, shortName, "RsIpslaMonPol", fmt.Sprintf("Expected short name to be 'RsIpslaMonPol', but got '%s'", shortName))
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
}

func TestSplitClassNameToPackageNameAndShortNameError(t *testing.T) {
	test.InitializeTest(t)
	packageName, shortName, err := splitClassNameToPackageNameAndShortName("error")
	assert.Equal(t, packageName, "", fmt.Sprintf("Expected package name to be '', but got '%s'", packageName))
	assert.Equal(t, shortName, "", fmt.Sprintf("Expected short name to be '', but got '%s'", shortName))
	assert.Error(t, err, fmt.Sprintf("Expected error, but got '%s'", err))
}

func TestSetResourceNameFromLabelNoRelationWithIdentifier(t *testing.T) {
	ds := initializeDataStoreTest(t)
	class := Class{ClassName: "tagAnnotation", IdentifiedBy: []string{"key"}}
	class.MetaFileContent = map[string]interface{}{
		"label": "annotation",
	}
	err := class.setResourceName(ds)
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
	assert.Equal(t, class.ResourceName, "annotation", fmt.Sprintf("Expected resource name to be 'annotation', but got '%s'", class.ResourceName))
	assert.Equal(t, class.ResourceNameNested, "annotations", fmt.Sprintf("Expected nested resource name to be 'annotations', but got '%s'", class.ResourceNameNested))
}

func TestSetResourceNameFromLabelNoRelationWithoutIdentifier(t *testing.T) {
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
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
	err = class.setResourceName(ds)
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
	assert.Equal(t, class.ResourceName, "relation_to_vrf", fmt.Sprintf("Expected resource name to be 'relation_to_vrf', but got '%s'", class.ResourceName))
	assert.Equal(t, class.ResourceNameNested, "relation_to_vrf", fmt.Sprintf("Expected nested resource name to be 'relation_to_vrf', but got '%s'", class.ResourceNameNested))
}

func TestSetResourceNameToRelation(t *testing.T) {
	ds := initializeDataStoreTest(t)
	ds.GlobalMetaDefinition = GlobalMetaDefinition{
		NoMetaFile: map[string]string{
			"vzBrCP": "contract",
		},
	}
	class := Class{ClassName: "fvRsCons", IdentifiedBy: []string{"tnVzBrCPName"}}
	class.MetaFileContent = relationInfoFvRsCons
	err := class.setRelation()
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
	err = class.setResourceName(ds)
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
	assert.Equal(t, class.ResourceName, "relation_to_contract", fmt.Sprintf("Expected resource name to be 'relation_to_contract', but got '%s'", class.ResourceName))
	assert.Equal(t, class.ResourceNameNested, "relation_to_contracts", fmt.Sprintf("Expected nested resource name to be 'relation_to_contracts', but got '%s'", class.ResourceName))
}

func TestSetResourceNameFromToRelation(t *testing.T) {
	ds := initializeDataStoreTest(t)
	ds.GlobalMetaDefinition = GlobalMetaDefinition{
		NoMetaFile: map[string]string{
			"fvCtx":               "vrf",
			"netflowAExporterPol": "netflow_exporter_policy",
		},
	}
	class := Class{ClassName: "netflowRsExporterToCtx"}
	class.MetaFileContent = relationInfoNetflowRsExporterToCtx
	err := class.setRelation()
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
	err = class.setResourceName(ds)
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
	assert.Equal(t, class.ResourceName, "relation_from_netflow_exporter_policy_to_vrf", fmt.Sprintf("Expected resource name to be 'relation_from_netflow_exporter_policy_to_vrf', but got '%s'", class.ResourceName))
	assert.Equal(t, class.ResourceNameNested, "relation_to_vrf", fmt.Sprintf("Expected nested resource name to be 'relation_to_vrf', but got '%s'", class.ResourceNameNested))
}

func TestSetResourceNameFromEmptyLabelError(t *testing.T) {
	ds := initializeDataStoreTest(t)
	class := Class{}
	class.MetaFileContent = map[string]interface{}{"label": ""}
	err := class.setResourceName(ds)
	assert.Error(t, err, fmt.Sprintf("Expected error, but got '%s'", err))
}

func TestSetResourceNameFromNoLabelError(t *testing.T) {
	ds := initializeDataStoreTest(t)
	class := Class{}
	class.MetaFileContent = map[string]interface{}{"no_label": ""}
	err := class.setResourceName(ds)
	assert.Error(t, err, fmt.Sprintf("Expected error, but got '%s'", err))
}

func TestSetRelationNoRelation(t *testing.T) {
	test.InitializeTest(t)
	class := Class{ClassName: "fvTenant"}
	err := class.setRelation()
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
	assert.Equal(t, class.Relation.FromClass, "", fmt.Sprintf("Expected FromClass to be 'fvEPg', but got '%s'", class.Relation.FromClass))
	assert.Equal(t, class.Relation.ToClass, "", fmt.Sprintf("Expected ToClass to be 'vzBrCP', but got '%s'", class.Relation.ToClass))
	assert.Equal(t, class.Relation.Type, RelationshipTypeEnum(0), fmt.Sprintf("Expected Type to be 'Undefined', but got '%v'", class.Relation.Type))
	assert.Equal(t, class.Relation.IncludeFrom, false, fmt.Sprintf("Expected IncludeFrom to be 'false', but got '%v'", class.Relation.IncludeFrom))
	assert.Equal(t, class.Relation.RelationalClass, false, fmt.Sprintf("Expected RelationalClass to be 'false', but got '%v'", class.Relation.RelationalClass))
}

func TestSetRelationIncludeFromFalseAndTypeNamed(t *testing.T) {
	test.InitializeTest(t)
	class := Class{ClassName: "fvRsCons"}
	class.MetaFileContent = relationInfoFvRsCons
	err := class.setRelation()
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
	assert.Equal(t, class.Relation.FromClass, "fvEPg", fmt.Sprintf("Expected FromClass to be 'fvEPg', but got '%s'", class.Relation.FromClass))
	assert.Equal(t, class.Relation.ToClass, "vzBrCP", fmt.Sprintf("Expected ToClass to be 'vzBrCP', but got '%s'", class.Relation.ToClass))
	assert.Equal(t, class.Relation.Type, Named, fmt.Sprintf("Expected Type to be 'Named', but got '%v'", class.Relation.Type))
	assert.Equal(t, class.Relation.IncludeFrom, false, fmt.Sprintf("Expected IncludeFrom to be 'false', but got '%v'", class.Relation.IncludeFrom))
	assert.Equal(t, class.Relation.RelationalClass, true, fmt.Sprintf("Expected RelationalClass to be 'true', but got '%v'", class.Relation.RelationalClass))
}

func TestSetRelationIncludeFromTrueAndTypeExplicit(t *testing.T) {
	test.InitializeTest(t)
	class := Class{ClassName: "netflowRsExporterToCtx"}
	class.MetaFileContent = relationInfoNetflowRsExporterToCtx
	err := class.setRelation()
	class.setAllowDelete()
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
	assert.Equal(t, class.Relation.FromClass, "netflowAExporterPol", fmt.Sprintf("Expected FromClass to be 'netflowAExporterPol', but got '%s'", class.Relation.FromClass))
	assert.Equal(t, class.Relation.ToClass, "fvCtx", fmt.Sprintf("Expected ToClass to be 'fvCtx', but got '%s'", class.Relation.ToClass))
	assert.Equal(t, class.Relation.Type, Explicit, fmt.Sprintf("Expected Type to be 'Explicit', but got '%v'", class.Relation.Type))
	assert.Equal(t, class.Relation.IncludeFrom, true, fmt.Sprintf("Expected IncludeFrom to be 'false', but got '%v'", class.Relation.IncludeFrom))
	assert.Equal(t, class.Relation.RelationalClass, true, fmt.Sprintf("Expected RelationalClass to be 'true', but got '%v'", class.Relation.RelationalClass))
}

func TestSetRelationUndefinedType(t *testing.T) {
	test.InitializeTest(t)
	class := Class{ClassName: "netflowRsExporterToCtx"}
	class.MetaFileContent = map[string]interface{}{
		"relationInfo": map[string]interface{}{
			"type":   "undefinedType",
			"fromMo": "netflow:AExporterPol",
			"toMo":   "fv:Ctx",
		},
	}
	err := class.setRelation()
	assert.Error(t, err, fmt.Sprintf("Expected error, but got '%s'", err))
}

func TestSetAllowDelete(t *testing.T) {
	test.InitializeTest(t)
	class := Class{ClassName: "fvPeeringP"}

	//Case 1: Verify that AllowDelete is false when fvPeeringP's definitions.allow_delete is an empty string and meta.isCreatableDeletable is set to 'never'.
	class.MetaFileContent = map[string]interface{}{
		"isCreatableDeletable": "never",
	}
	class.setAllowDelete()
	assert.Equal(t, class.AllowDelete, false, fmt.Sprintf("Expected isCreatableDeletable to be 'false', but got '%t'", class.AllowDelete))

	//Case 2: Verify that AllowDelete is true when fvPeeringP's definitions.allow_delete is an empty string and meta.isCreatableDeletable is not set to 'never'.
	class.MetaFileContent = map[string]interface{}{
		"isCreatableDeletable": "always",
	}
	class.setAllowDelete()
	assert.Equal(t, class.AllowDelete, true, fmt.Sprintf("Expected isCreatableDeletable to be 'true', but got '%t'", class.AllowDelete))

	// Reset class.AllowDelete to false for subsequent definition overwrite checks.
	class.AllowDelete = false

	//Case 3: Verify that AllowDelete is false if fvPeeringP's definitions.allow_delete is explicitly set to 'never'.
	class.ClassDefinition = ClassDefinition{AllowDelete: "never"}
	class.setAllowDelete()
	assert.Equal(t, class.AllowDelete, false, fmt.Sprintf("Expected isCreatableDeletable to be 'false', but got '%t'", class.AllowDelete))

	//Case 4: Verify that AllowDelete is true if fvPeeringP's definitions.allow_delete is explicitly not set to 'never'.
	class.ClassDefinition = ClassDefinition{AllowDelete: "always"}
	class.setAllowDelete()
	assert.Equal(t, class.AllowDelete, true, fmt.Sprintf("Expected isCreatableDeletable to be 'true', but got '%t'", class.AllowDelete))
}
