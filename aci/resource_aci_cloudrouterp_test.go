package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciCloudVpnGateway_Basic(t *testing.T) {
	var cloud_vpn_gateway models.CloudVpnGateway
	description := "cloud_vpn_gateway created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudVpnGatewayConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudVpnGatewayExists("aci_cloud_vpn_gateway.foocloud_vpn_gateway", &cloud_vpn_gateway),
					testAccCheckAciCloudVpnGatewayAttributes(description, &cloud_vpn_gateway),
				),
			},
		},
	})
}

func TestAccAciCloudVpnGateway_update(t *testing.T) {
	var cloud_vpn_gateway models.CloudVpnGateway
	description := "cloud_vpn_gateway created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudVpnGatewayConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudVpnGatewayExists("aci_cloud_vpn_gateway.foocloud_vpn_gateway", &cloud_vpn_gateway),
					testAccCheckAciCloudVpnGatewayAttributes(description, &cloud_vpn_gateway),
				),
			},
			{
				Config: testAccCheckAciCloudVpnGatewayConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudVpnGatewayExists("aci_cloud_vpn_gateway.foocloud_vpn_gateway", &cloud_vpn_gateway),
					testAccCheckAciCloudVpnGatewayAttributes(description, &cloud_vpn_gateway),
				),
			},
		},
	})
}

func testAccCheckAciCloudVpnGatewayConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_cloud_vpn_gateway" "foocloud_vpn_gateway" {
		cloud_context_profile_dn  = aci_cloud_context_profile.foocloud_context_profile.id
		description    = "%s"
		name           = "example_name"
		annotation     = "example_annotation"
		name_alias     = "example_alias"
		num_instances  = "1"
		cloud_router_profile_type  = "vpn-gw"
	}
	`, description)
}

func testAccCheckAciCloudVpnGatewayExists(name string, cloud_vpn_gateway *models.CloudVpnGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Router Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Router Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_vpn_gatewayFound := models.CloudVpnGatewayFromContainer(cont)
		if cloud_vpn_gatewayFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Router Profile %s not found", rs.Primary.ID)
		}
		*cloud_vpn_gateway = *cloud_vpn_gatewayFound
		return nil
	}
}

func testAccCheckAciCloudVpnGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_vpn_gateway" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_vpn_gateway := models.CloudVpnGatewayFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Router Profile %s Still exists", cloud_vpn_gateway.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudVpnGatewayAttributes(description string, cloud_vpn_gateway *models.CloudVpnGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_vpn_gateway.Description {
			return fmt.Errorf("Bad cloud_vpn_gateway Description %s", cloud_vpn_gateway.Description)
		}

		if "example_name" != cloud_vpn_gateway.Name {
			return fmt.Errorf("Bad cloud_vpn_gateway name %s", cloud_vpn_gateway.Name)
		}

		if "example_annotation" != cloud_vpn_gateway.Annotation {
			return fmt.Errorf("Bad cloud_vpn_gateway annotation %s", cloud_vpn_gateway.Annotation)
		}

		if "example_alias" != cloud_vpn_gateway.NameAlias {
			return fmt.Errorf("Bad cloud_vpn_gateway name_alias %s", cloud_vpn_gateway.NameAlias)
		}

		if "1" != cloud_vpn_gateway.NumInstances {
			return fmt.Errorf("Bad cloud_vpn_gateway num_instances %s", cloud_vpn_gateway.NumInstances)
		}

		if "vpn-gw" != cloud_vpn_gateway.CloudVpnGateway_type {
			return fmt.Errorf("Bad cloud_vpn_gateway cloud_vpn_gateway_type %s", cloud_vpn_gateway.CloudVpnGateway_type)
		}

		return nil
	}
}
