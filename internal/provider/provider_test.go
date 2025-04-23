package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/aci"
	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

var testAccProtoV6ProviderFactoriesFunction = map[string]func() (tfprotov6.ProviderServer, error){
	"aci": providerserver.NewProtocol6WithError(New("test")()),
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

		aciClientTest = client.NewClient(aci_url, aci_username, client.Password(aci_password), client.Insecure(true), client.MaxRetries(2))
	})
	return aciClientTest
}

func testAccPreCheck(t *testing.T, testType string, testApplicableFromVersion string) {
	// globalAnnotation is explicitly set here, as the provider initialization in acceptance tests differs from the standard initialization process.
	setGlobalAnnotation(basetypes.NewStringNull(), "ACI_ANNOTATION")

	if testType == "cloud" {
		t.Skip("[WARNING] Test skipped because it is not supported on an on-prem APIC")
	}

	infoController, err := getACIClientTest(t).GetViaURL("/api/node/class/firmwareCtrlrRunning.json")
	if err != nil {
		t.Fatalf("Error fetching APIC controller information: %v", err)
	}
	apicVersion := extractControllerVersion(infoController)
	// TODO process the version when it has ranges associated with it
	if apicVersion != strings.TrimSuffix(testApplicableFromVersion, "-") {
		if IsVersionGreater(*ParseVersion(testApplicableFromVersion).Version, *ParseVersion(apicVersion).Version) {
			t.Skip("[WARNING] Test skipped because it is not supported on APIC version:", apicVersion)
		}
	}

}

func composeAggregateTestCheckFuncWithVersion(t *testing.T, propertyVersion string, operator string, checks ...resource.TestCheckFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.aci_system.version"]

		if !ok {
			return fmt.Errorf("data source aci_system was not found in the test configuration")
		}

		apicVersion := rs.Primary.Attributes["version"]

		var result []error

		comparisonResult, err := CompareVersionsRange(apicVersion, propertyVersion, operator)
		if err != nil {
			return fmt.Errorf("Failed to compare versions in the test check function: %w", err)
		}

		if !comparisonResult {
			t.Logf("[WARNING] Test for checking the resource attribute value skipped due to version incompatibility between Property version: %s and APIC version: %s.", propertyVersion, apicVersion)
			return nil
		}

		for i, check := range checks {
			if err := check(s); err != nil {
				result = append(result, fmt.Errorf("Check %d/%d error: %w", i+1, len(checks), err))
			}

		}
		return errors.Join(result...)
	}
}

func extractControllerVersion(requestData *container.Container) string {
	classReadInfo := requestData.Search("imdata").Search("firmwareCtrlrRunning").Data().([]interface{})
	if len(classReadInfo) == 1 {
		attributes := classReadInfo[0].(map[string]interface{})["attributes"].(map[string]interface{})
		for attributeName, attributeValue := range attributes {
			if attributeName == "version" {
				return attributeValue.(string)
			}
		}
	}
	return ""
}

func testCheckResourceDestroy(s *terraform.State) error {
	aciClient := getACIClientTest(nil)
	for name, rs := range s.RootModule().Resources {
		if strings.HasPrefix(name, "aci_") {
			_, err := aciClient.Get(rs.Primary.ID)
			if err != nil {
				if strings.Contains(err.Error(), "Error retrieving Object: Object may not exist") {
					continue
				} else {
					return fmt.Errorf("error checking if resource '%s' with ID '%s' still exists: %s", name, rs.Primary.ID, err)
				}
			}
			return fmt.Errorf("terraform destroy was unsuccessful. The resource '%s' with ID '%s' still exists", name, rs.Primary.ID)
		}
	}
	return nil
}

func setEnvVariable(t *testing.T, key, value string) {
	t.Setenv(key, value)
}

func CheckOutputBool(name string, value bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()
		rs, ok := ms.Outputs[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		boolValue, err := strconv.ParseBool(rs.Value.(string))
		if err != nil {
			return fmt.Errorf("Error Parsing value: %s", err)
		}
		if boolValue != value {
			return fmt.Errorf(
				"Output '%s': expected %#v, got %#v",
				name,
				value,
				rs.Value)
		}

		return nil
	}
}
