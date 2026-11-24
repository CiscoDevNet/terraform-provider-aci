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
	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/acceptancetest"
	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/assert"
)

const testAccConfig = `
variable "apic_version" {
  type = string
}

variable "test_prefix" {
  type = string
}
`

var testAccRunner acceptancetest.Runner

func TestMain(m *testing.M) {
	exitCode := m.Run()
	// Finalize only in the parent process; unused runners and worker processes have no report to write.
	if err := testAccRunner.Finish(); err != nil {
		fmt.Fprintf(os.Stderr, "Write acceptance-test report: %v\n", err)
		if exitCode == 0 {
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"aci": func() (tfprotov6.ProviderServer, error) {
		// The harness recreates provider servers for each Terraform command, repeating
		// client initialization. See gen/docs/test_configuration.md#provider-authentication-limitation.
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
	aciClientTest       *client.Client
	aciClientTestOnce   sync.Once
	apicVersionTest     string
	apicVersionTestErr  error
	apicVersionTestOnce sync.Once
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

func getAPICVersionTest(t *testing.T) string {
	t.Helper()
	aciClient := getACIClientTest(t)
	// The test client targets one APIC endpoint for the process, so every test
	// configuration and assertion can reuse the same controller-version lookup.
	apicVersionTestOnce.Do(func() {
		infoController, err := aciClient.GetViaURL("/api/node/class/firmwareCtrlrRunning.json")
		if err != nil {
			apicVersionTestErr = err
			return
		}
		apicVersionTest = extractControllerVersion(infoController)
	})
	if apicVersionTestErr != nil {
		t.Fatalf("Error fetching APIC controller information: %v", apicVersionTestErr)
	}
	return apicVersionTest
}

func testAccPreCheck(t *testing.T, testType string, testApplicableFromVersion string) {
	// globalAnnotation is explicitly set here, as the provider initialization in acceptance tests differs from the standard initialization process.
	setGlobalAnnotation(basetypes.NewStringNull(), "ACI_ANNOTATION")

	if testType == "cloud" {
		acceptancetest.Skip(t, "This class is not supported on an on-prem APIC")
	}

	apicVersion := getAPICVersionTest(t)

	comparisonResult, err := CompareVersionsRange(apicVersion, testApplicableFromVersion, "inside")
	if err != nil {
		acceptancetest.Skip(t, fmt.Sprintf("The supported range %q could not be checked: %s", testApplicableFromVersion, err))
	} else if !comparisonResult {
		acceptancetest.Skip(t, fmt.Sprintf("APIC %s is outside the supported range: %s", apicVersion, testApplicableFromVersion))
	}

}

func composeAggregateTestCheckFuncWithVersion(t *testing.T, propertyVersion string, operator string, checks ...resource.TestCheckFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		apicVersion := getAPICVersionTest(t)

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

type versionedStateCheck struct {
	t           *testing.T
	versions    string
	apicVersion func() string
	check       statecheck.StateCheck
}

func stateCheckWithVersion(t *testing.T, versions string, check statecheck.StateCheck) statecheck.StateCheck {
	return versionedStateCheck{
		t:           t,
		versions:    versions,
		apicVersion: func() string { return getAPICVersionTest(t) },
		check:       check,
	}
}

func (c versionedStateCheck) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	// Resolve the controller version when the check runs, after the acceptance precheck.
	apicVersion := c.apicVersion()
	supported, err := CompareVersionsRange(apicVersion, c.versions, "inside")
	if err != nil {
		resp.Error = fmt.Errorf("failed to compare versions in the state check: %w", err)
		return
	}
	if !supported {
		c.t.Logf("[WARNING] State check skipped due to version incompatibility between property version: %s and APIC version: %s.", c.versions, apicVersion)
		return
	}
	c.check.CheckState(ctx, req, resp)
}

func TestStateCheckWithVersion(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	type input struct {
		APICVersion string
		Versions    string
		Target      any
	}
	testCases := []test.TestCase{
		{Name: "test_supported_empty_target", Input: input{APICVersion: "6.2(2e)", Versions: "1.0(1e)-", Target: ""}},
		{Name: "test_supported_resolved_target", Input: input{APICVersion: "6.2(2e)", Versions: "1.0(1e)-", Target: "uni/tn-common/monepg-policy_1"}},
		{Name: "test_supported_null_target", Input: input{APICVersion: "6.2(2e)", Versions: "1.0(1e)-"}, Expected: "expected non-nil value"},
		{Name: "test_below_supported_range", Input: input{APICVersion: "2.0(1e)", Versions: "3.0(1k)-"}},
		{Name: "test_above_supported_range", Input: input{APICVersion: "6.2(2e)", Versions: "1.0(1e)-6.1(1e)"}},
		{Name: "test_invalid_range", Input: input{APICVersion: "6.2(2e)", Versions: "invalid"}, Expected: "failed to compare versions"},
		{Name: "test_invalid_apic_version", Input: input{APICVersion: "invalid", Versions: "1.0(1e)-"}, Expected: "failed to compare versions"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			value := testCase.Input.(input)
			check := stateCheckWithVersion(t, value.Versions, statecheck.ExpectKnownValue("data.aci_tenant.test", tfjsonpath.New("relation"), knownvalue.NotNull())).(versionedStateCheck)
			versionReads := 0
			check.apicVersion = func() string {
				versionReads++
				return value.APICVersion
			}
			state := &tfjson.State{Values: &tfjson.StateValues{RootModule: &tfjson.StateModule{Resources: []*tfjson.StateResource{{
				Address:         "data.aci_tenant.test",
				AttributeValues: map[string]any{"relation": value.Target},
			}}}}}
			var response statecheck.CheckStateResponse
			check.CheckState(context.Background(), statecheck.CheckStateRequest{State: state}, &response)
			assert.Equal(t, 1, versionReads)
			if expected, _ := testCase.Expected.(string); expected != "" {
				assert.ErrorContains(t, response.Error, expected)
			} else {
				assert.NoError(t, response.Error)
			}
		})
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
