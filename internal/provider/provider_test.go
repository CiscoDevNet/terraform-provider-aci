package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/aci"
	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"aci": func() (tfprotov6.ProviderServer, error) {
		ctx := context.Background()

		upgradedSdkServer, err := tf5to6server.UpgradeServer(
			ctx,
			aci.Provider().GRPCProvider,
		)

		if err != nil {
			return nil, err
		}

		providers := []func() tfprotov6.ProviderServer{
			providerserver.NewProtocol6(New("test")()),
			func() tfprotov6.ProviderServer {
				return upgradedSdkServer
			},
		}

		muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	},
}

var (
	aciClientTest     *client.Client
	aciClientTestOnce sync.Once
)

func getACIClientTest(t *testing.T) *client.Client {
	aciClientTestOnce.Do(func() {
		var aci_url, aci_username, aci_password string
		if v := os.Getenv("ACI_USERNAME"); v == "" {
			t.Fatal("ACI_USERNAME must be set for acceptance tests")
		} else {
			aci_username = v
		}
		if v := os.Getenv("ACI_PASSWORD"); v == "" {
			t.Fatal("ACI_PASSWORD must be set for acceptance tests")
		} else {
			aci_password = v
		}
		if v := os.Getenv("ACI_URL"); v == "" {
			t.Fatal("ACI_URL must be set for acceptance tests")
		} else {
			aci_url = v
		}
		if v := os.Getenv("ACI_VAL_REL_DN"); v == "" {
			t.Fatal("ACI_VAL_REL_DN must be set for acceptance tests")
			boolValue, err := strconv.ParseBool(v)
			if err != nil || boolValue == true {
				t.Fatal("ACI_VAL_REL_DN must be a 'false' boolean value")
			}
		}

		aciClientTest = client.NewClient(aci_url, aci_username, client.Password(aci_password), client.Insecure(true))
	})
	return aciClientTest
}

func testAccPreCheck(t *testing.T, testType string) {
	infoCloud, _ := getACIClientTest(t).GetViaURL("/api/node/class/cloudProvP.json")
	environment, _ := extractEnvironmentValue(infoCloud)
	if environment == "public-cloud" && testType == "apic" {
		t.Skip("[WARNING] Skipping the test because the test cannot be run on a cloud APIC")
	} else if environment != "public-cloud" && testType == "cloud" {
		t.Skip("[WARNING] Skipping the test because the test cannot be run on an on-prem APIC")
	}
}

func extractEnvironmentValue(requestData *container.Container) (string, error) {
	if requestData.Search("imdata").Search("cloudProvP").Data() != nil {
		classReadInfo := requestData.Search("imdata").Search("cloudProvP").Data().([]interface{})
		if len(classReadInfo) == 1 {
			attributes := classReadInfo[0].(map[string]interface{})["attributes"].(map[string]interface{})
			for attributeName, attributeValue := range attributes {
				if attributeName == "environment" {
					return attributeValue.(string), nil
				}
			}
		}
	}
	return "", fmt.Errorf("no cloudProvP instances found in the response")

}

// waitForApicBeforeRefresh ensures the APIC is available by polling the API until it responds or times out before Terraform refresh is applied in the test.
func waitForApicBeforeRefresh(s *terraform.State) error {
	aciClient := getACIClientTest(nil)

	timeoutTimer := time.NewTimer(50 * time.Second)
	defer timeoutTimer.Stop()

	pollTimer := time.NewTimer(5 * time.Second)
	defer pollTimer.Stop()

	for {
		select {
		case <-timeoutTimer.C:
			return fmt.Errorf("timeout reached while waiting for APIC to become available")
		case <-pollTimer.C:
			_, err := aciClient.GetViaURL("/api/aaaListDomains.json")
			if err != nil {
				pollTimer.Reset(5 * time.Second)
			} else {
				return nil
			}
		}
	}
}

func testCheckResourceAttr(resourceName, attribute, value1 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		attrValue, ok := rs.Primary.Attributes[attribute]
		if !ok {
			return nil
		}

		if attrValue != value1 {
			return fmt.Errorf("attribute %s in resource %s should be %s , but got %s", attribute, resourceName, value1, attrValue)
		}

		return nil
	}
}

func testCheckResourceDestroy(s *terraform.State) error {
	aciClient := getACIClientTest(nil)

	for name, rs := range s.RootModule().Resources {
		if !strings.HasPrefix(name, "data.") {
			_, err := aciClient.Get(rs.Primary.ID)
			if err != nil {
				if strings.Contains(err.Error(), "Error retrieving Object: Object may not exist") {
					continue
				} else {
					return fmt.Errorf("error checking if resource '%s' with ID '%s' still exists: %s", rs.Type, rs.Primary.ID, err)
				}
			}
			return fmt.Errorf("terraform destroy was unsuccessful. The resource '%s' with ID '%s' still exists", rs.Type, rs.Primary.ID)
		}
	}
	return nil
}

func setGlobalAnnotationEnvVariable(t *testing.T, annotation string) {
	t.Setenv("ACI_ANNOTATION", annotation)
}
