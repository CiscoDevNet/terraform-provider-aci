package testacc

import (
	"io/ioutil"
	"log"
	"os"

	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

//TODO: check password is not showing in state file

var testAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider
var systemInfo *models.System

func init() {
	testAccProvider = aci.Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"aci": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
	log.SetOutput(ioutil.Discard)
	systemInfo = fetchSysInfo()
}
func TestProvider(t *testing.T) {
	if err := aci.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func sharedAciClient() *client.Client {
	config := aci.Config{
		Username:   os.Getenv("ACI_USERNAME"),
		Password:   os.Getenv("ACI_PASSWORD"),
		URL:        os.Getenv("ACI_URL"),
		PrivateKey: os.Getenv("ACI_PRIVATE_KEY"),
		Certname:   os.Getenv("ACI_CERT_NAME"),
		ProxyUrl:   os.Getenv("ACI_PROXY_URL"),
		ProxyCreds: os.Getenv("ACI_PROXY_CREDS"),
		IsInsecure: true,
	}
	return config.GetClient().(*client.Client)
}

func fetchSysInfo() *models.System {

	aciClient := sharedAciClient()
	topSystemCont, err := aciClient.GetViaURL("/api/node/class/topSystem.json")
	if err != nil {
		log.Panic("System info not found:", err)
	}

	return models.SystemListFromContainer(topSystemCont)[0]
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
