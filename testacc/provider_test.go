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

var testAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider
var systemInfo *models.System
var providerProfileDn string = "uni/vmmp-VMware"
var certificate_terraformuser string = `
-----BEGIN CERTIFICATE-----
MIICMjCCAZugAwIBAgIJAMqfErBhyJ/MMA0GCSqGSIb3DQEBCwUAMDIxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMQ4wDAYDVQQKDAVjaXNjbzAeFw0y
MjAyMDIwMjE4MzBaFw0yMjAzMDQwMjE4MzBaMDIxCzAJBgNVBAYTAkFVMRMwEQYD
VQQIDApTb21lLVN0YXRlMQ4wDAYDVQQKDAVjaXNjbzCBnzANBgkqhkiG9w0BAQEF
AAOBjQAwgYkCgYEAqsd9v1p8RhFJ0kFsIRhLmTsk71ezZWKXfYVH/dvOHUgiRASY
c/XxrKf6OHSzDVayqrRl+kH1s08qdEVR+CzghZerYfCF6tXmfE2sAkkLQJStPpUW
8x8oGFGXsqE8FGwoVM0NYhNk8Mb5qGYG3RkrUe/iA7ciGWc5XXplGW4t7CMCAwEA
AaNQME4wHQYDVR0OBBYEFBnezWx5WfUKJ7TVq9GaCK7keOZIMB8GA1UdIwQYMBaA
FBnezWx5WfUKJ7TVq9GaCK7keOZIMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEL
BQADgYEAOwT6w5d5GeZo8dFtWDLQMPC2zG8G5gZKe8BDbkoo1iYPDEr/mwd40FDI
iPS+W0nUlIiYmbdg8oBOvPOdaPxdvlzunRHenDg3rrUoT6v+fEG84PTKMKFveP/P
6QmNWpq9VtLOWnOWCJSuojYd5Ss//hDLbgmTcCmK0VOmQlAW9qA=
-----END CERTIFICATE-----`

const region = "us-west-1"
const cloudVendor = "aws"
const cloudProvPName = "aws"
const name = "us-east-1"
const cloudProviderRegion = "uni/clouddomp/provp-aws/region-us-east-1"
const zoneName = "us-east-1a"
const fabricNodeId = "101"
const fabricPodDn = "topology/pod-1"
const podId = "1"
const nodeId = "101"
const pathEpName = "eth1/12"
const fabDn1 = "topology/pod-1/node-101"
const systemPodId = "1"
const systemNodeId = "1"
const hubNetwork = "uni/tn-infra/gwrouterp-default"
const tdn1 = "topology/pod-1/paths-101/pathep-[eth1/12]"
const tdn2 = "topology/pod-1/paths-101/pathep-[eth1/6]"
const multdn1 = "topology/pod-1/paths-101/pathep-[eth1/21]"
const multdn2 = "topology/pod-1/paths-101/pathep-[eth1/30]"
const multdn3 = "topology/pod-1/paths-101/pathep-[eth1/1]"
const fabDn2 = "topology/pod-1/node-201"
const fabDn3 = "topology/pod-1/node-111"
const fabDn4 = "topology/pod-1/node-1"
const vmmProvProfileDn = "uni/vmmp-VMware"
const vmmProvProfileDnOther = "uni/vmmp-OpenShift"
const fabricNodeMemNodeId1 = "301"
const fabricNodeMemNodeId2 = "302"
const fabricNodeMemNodeId3 = "303"
const fabricNodeMemNodeId4 = "304"
const fabricNodeMemNodeId5 = "306"
const pathEp1 = "topology/pod-1/paths-101/pathep-[eth1/1]"
const pathEp2 = "topology/pod-1/paths-101/pathep-[eth1/12]"
const pathEp3 = "topology/pod-1/paths-101/pathep-[eth1/6]"
const pathEp4 = "topology/pod-1/paths-101/pathep-[eth1/17]"
const pathEp5 = "topology/pod-1/paths-101/pathep-[eth1/34]"

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
		log.Panic("System info not found: ", err)
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
