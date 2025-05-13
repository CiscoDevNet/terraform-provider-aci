package data

import (
	"fmt"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

const (
	constTestClassNameSingleWordInShortName    = "fvTenant"
	constTestClassNameMultipleWordsInShortName = "fvRsIpslaMonPol"
	constTestClassNameErrorInShortName         = "error"
	constTestMetaFileContentForLabel           = "tenant"
)

func TestSplitClassNameToPackageNameAndShortNameSingle(t *testing.T) {
	test.InitializeTest(t)
	packageName, shortName, err := splitClassNameToPackageNameAndShortName(constTestClassNameSingleWordInShortName)
	assert.Equal(t, packageName, "fv", fmt.Sprintf("Expected package name to be 'fv', but got '%s'", packageName))
	assert.Equal(t, shortName, "Tenant", fmt.Sprintf("Expected short name to be 'Tenant', but got '%s'", shortName))
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
}

func TestSplitClassNameToPackageNameAndShortNameMultiple(t *testing.T) {
	test.InitializeTest(t)
	packageName, shortName, err := splitClassNameToPackageNameAndShortName(constTestClassNameMultipleWordsInShortName)
	assert.Equal(t, packageName, "fv", fmt.Sprintf("Expected package name to be 'fv', but got '%s'", packageName))
	assert.Equal(t, shortName, "RsIpslaMonPol", fmt.Sprintf("Expected short name to be 'RsIpslaMonPol', but got '%s'", shortName))
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
}

func TestSplitClassNameToPackageNameAndShortNameError(t *testing.T) {
	test.InitializeTest(t)
	packageName, shortName, err := splitClassNameToPackageNameAndShortName(constTestClassNameErrorInShortName)
	assert.Equal(t, packageName, "", fmt.Sprintf("Expected package name to be '', but got '%s'", packageName))
	assert.Equal(t, shortName, "", fmt.Sprintf("Expected short name to be '', but got '%s'", shortName))
	assert.Error(t, err, fmt.Sprintf("Expected error, but got '%s'", err))
}

func TestSetResourceNameFromLabel(t *testing.T) {
	test.InitializeTest(t)
	class := Class{ClassName: constTestClassNameSingleWordInShortName}
	class.MetaFileContent = map[string]interface{}{"label": constTestMetaFileContentForLabel}
	err := class.setResourceName()
	assert.NoError(t, err, fmt.Sprintf("Expected no error, but got '%s'", err))
	assert.Equal(t, class.ResourceName, constTestMetaFileContentForLabel, fmt.Sprintf("Expected resource name to be 'tenant', but got '%s'", class.ResourceName))
}

func TestSetResourceNameFromEmptyLabelError(t *testing.T) {
	test.InitializeTest(t)
	class := Class{ClassName: constTestClassNameSingleWordInShortName}
	class.MetaFileContent = map[string]interface{}{"label": ""}
	err := class.setResourceName()
	assert.Error(t, err, fmt.Sprintf("Expected error, but got '%s'", err))
}

func TestSetResourceNameFromNoLabelError(t *testing.T) {
	test.InitializeTest(t)
	class := Class{ClassName: constTestClassNameSingleWordInShortName}
	class.MetaFileContent = map[string]interface{}{"no_label": ""}
	err := class.setResourceName()
	assert.Error(t, err, fmt.Sprintf("Expected error, but got '%s'", err))
}

func TestSetResourceNameNested(t *testing.T) {
	test.InitializeTest(t)
	class := Class{ClassName: constTestClassNameSingleWordInShortName}

	tests := []map[string]interface{}{
		{"resource_name": "relation_from_bridge_domain_to_netflow_monitor_policy", "identifiers": []string{"tnNetflowMonitorPolName", "fltType"}, "expected": "relation_to_netflow_monitor_policies"},
		{"resource_name": "annotation", "identifiers": []string{"key"}, "expected": "annotations"},
		{"resource_name": "relation_to_consumed_contract", "identifiers": []string{"tnVzBrCPName"}, "expected": "relation_to_consumed_contracts"},
		{"resource_name": "associated_site", "identifiers": []string{}, "expected": "associated_site"},
		{"resource_name": "relation_from_netflow_exporter_to_vrf", "identifiers": []string{}, "expected": "relation_to_vrf"},
	}

	for _, test := range tests {
		genLogger.Info(fmt.Sprintf("Executing: %s' with input '%s' and expected output '%s'", t.Name(), test["resource_name"], test["expected"]))
		class.ResourceName = test["resource_name"].(string)
		class.IdentifiedBy = test["identifiers"].([]string)
		class.setResourceNameNested()
		assert.Equal(t, test["expected"], class.ResourceNameNested, fmt.Sprintf("Expected '%s', but got '%s'", test["expected"], class.ResourceNameNested))
	}

}
