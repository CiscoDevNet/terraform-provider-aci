package acctest

import (
	// "fmt"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aci/aci"
	"os"
	// "regexp"
	"testing"
)

//TODO: check password is not showing in state file

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = aci.Provider()
	testAccProviders = map[string]*schema.Provider{
		"aci": testAccProvider,
	}
}
func TestProvider(t *testing.T) {
	if err := aci.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = aci.Provider()
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
