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
)

func TestSplitClassNameToPackageNameAndShortNameSingle(t *testing.T) {
	test.InitializeTest(t, 0)
	packageName, shortName := splitClassNameToPackageNameAndShortName(constTestClassNameSingleWordInShortName)
	assert.Equal(t, packageName, "fv", fmt.Sprintf("Expected package name to be 'fv', but got '%s'", packageName))
	assert.Equal(t, shortName, "Tenant", fmt.Sprintf("Expected short name to be 'Tenant', but got '%s'", shortName))
}

func TestSplitClassNameToPackageNameAndShortNameMultiple(t *testing.T) {
	test.InitializeTest(t, 0)
	packageName, shortName := splitClassNameToPackageNameAndShortName(constTestClassNameMultipleWordsInShortName)
	assert.Equal(t, packageName, "fv", fmt.Sprintf("Expected package name to be 'fv', but got '%s'", packageName))
	assert.Equal(t, shortName, "RsIpslaMonPol", fmt.Sprintf("Expected short name to be 'RsIpslaMonPol', but got '%s'", shortName))
}
