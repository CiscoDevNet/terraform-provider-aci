package data

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/stretchr/testify/assert"
)

const (
	metaHost = "10.0.0.1"
)

func initializeDataStoreTest(t *testing.T) *DataStore {
	test.InitializeTest(t)
	return &DataStore{}
}

func TestSetHostDefault(t *testing.T) {
	ds := initializeDataStoreTest(t)
	ds.setMetaHost()
	assert.Equal(t, ds.metaHost, constPubhubDevnetHost, "Expected meta host to be set to the default value %s, but got %s", constPubhubDevnetHost, ds.metaHost)
}

func TestSetHostFromEnvironmentVariable(t *testing.T) {
	ds := initializeDataStoreTest(t)
	t.Setenv(constEnvMetaHost, metaHost)
	ds.setMetaHost()
	assert.Equal(t, ds.metaHost, metaHost, "Expected meta host to be set to the custom value %s, but got %s", metaHost, ds.metaHost)
}
