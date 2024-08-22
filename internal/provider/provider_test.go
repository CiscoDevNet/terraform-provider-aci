package provider

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/aci"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

func testAccPreCheck(t *testing.T) {
	// globalAnnotation is explicitly set here, as the provider initialization in acceptance tests differs from the standard initialization process.
	setGlobalAnnotation(basetypes.NewStringNull(), "ACI_ANNOTATION")

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
}

func setEnvVariable(t *testing.T, key, value string) {
	t.Setenv(key, value)
}
