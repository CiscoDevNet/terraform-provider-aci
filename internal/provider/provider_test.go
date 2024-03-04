package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/aci"
	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
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

func testAccPreCheck(t *testing.T, testType string) {
	if v := os.Getenv("ACI_USERNAME"); v == "" {
		t.Fatal("ACI_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("ACI_PASSWORD"); v == "" {
		t.Fatal("ACI_PASSWORD must be set for acceptance tests")
	}
	if v := os.Getenv("ACI_URL"); v == "" {
		t.Fatal("ACI_URL must be set for acceptance tests")
	}
	if v := os.Getenv("ACI_VAL_REL_DN"); v == "" {
		t.Fatal("ACI_VAL_REL_DN must be set for acceptance tests")
		boolValue, err := strconv.ParseBool(v)
		if err != nil || boolValue == true {
			t.Fatal("ACI_VAL_REL_DN must be a 'false' boolean value")
		}
	}
	aci_url := os.Getenv("ACI_URL")
	aci_username := os.Getenv("ACI_USERNAME")
	aci_password := os.Getenv("ACI_PASSWORD")
	aciClient := client.NewClient(aci_url, aci_username, client.Password(aci_password), client.Insecure(true))

	cloudProviderUrl := "/api/node/class/cloudProvP.json"
	infoCloud, _ := aciClient.GetViaURL(cloudProviderUrl)
	environment, _ := extractEnvironmentValue(infoCloud)
	if environment == "public-cloud" && testType == "apic" {
		t.Skip("[WARNING] Skipping test because the test cannot be run on a cloud APIC")
	} else if environment != "public-cloud" && testType == "cloud" {
		t.Skip("[WARNING] Skipping test because the test cannot be run on an on-prem APIC")
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

func setGlobalAnnotationEnvVariable(t *testing.T, annotation string) {
	t.Setenv("ACI_ANNOTATION", annotation)
}
