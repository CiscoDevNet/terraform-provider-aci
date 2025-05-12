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
