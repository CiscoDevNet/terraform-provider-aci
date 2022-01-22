package aci

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"aci": testAccProvider,
	}
}
func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	// We will use this function later on to make sure our test environment is valid.
	// For example, you can make sure here that some environment variables are set.
	if v := os.Getenv("ACI_USERNAME"); v == "" {
		t.Fatal("ACI_USERNAME env variable must be set for acceptance tests")
	}
	if v := os.Getenv("ACI_PASSWORD"); v == "" {
		privateKey := os.Getenv("ACI_PRIVATE_KEY")
		certName := os.Getenv("ACI_CERT_NAME")
		if privateKey == "" && certName == "" {
			t.Fatal("Either of ACI_PASSWORD or ACI_PRIVATE_KEY/ACI_CERT_NAME env variables must be set for acceptance tests")
		}
	}
	if v := os.Getenv("ACI_URL"); v == "" {
		t.Fatal("ACI_URL env variable must be set for acceptance tests")
	}
}

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"aci": func() (*schema.Provider, error) {
		return testAccProvider, nil
	},
}
