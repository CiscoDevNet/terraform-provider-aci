package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSubnetPoolforIpSecTunnels_Basic(t *testing.T) {
	var subnet_poolfor_ip_sec_tunnels models.SubnetPoolforIpSecTunnels
	fv_tenant_name := acctest.RandString(5)
	aci_cloud_external_network_vpn_network := acctest.RandString(5)
	aci_cloud_ipsec_tunnel_subnet_pool := acctest.RandString(5)
	description := "subnet_poolfor_ip_sec_tunnels created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetPoolforIpSecTunnelsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubnetPoolforIpSecTunnelsConfig_basic(fv_tenant_name, aci_cloud_external_network_vpn_network, aci_cloud_ipsec_tunnel_subnet_pool),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetPoolforIpSecTunnelsExists("aci_subnet_poolfor_ip_sec_tunnels.foosubnet_poolfor_ip_sec_tunnels", &subnet_poolfor_ip_sec_tunnels),
					testAccCheckAciSubnetPoolforIpSecTunnelsAttributes(fv_tenant_name, aci_cloud_external_network_vpn_network, aci_cloud_ipsec_tunnel_subnet_pool, description, &subnet_poolfor_ip_sec_tunnels),
				),
			},
		},
	})
}

func TestAccAciSubnetPoolforIpSecTunnels_Update(t *testing.T) {
	var subnet_poolfor_ip_sec_tunnels models.SubnetPoolforIpSecTunnels
	fv_tenant_name := acctest.RandString(5)
	aci_cloud_external_network_vpn_network := acctest.RandString(5)
	aci_cloud_ipsec_tunnel_subnet_pool := acctest.RandString(5)
	description := "subnet_poolfor_ip_sec_tunnels created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetPoolforIpSecTunnelsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubnetPoolforIpSecTunnelsConfig_basic(fv_tenant_name, aci_cloud_external_network_vpn_network, aci_cloud_ipsec_tunnel_subnet_pool),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetPoolforIpSecTunnelsExists("aci_subnet_poolfor_ip_sec_tunnels.foosubnet_poolfor_ip_sec_tunnels", &subnet_poolfor_ip_sec_tunnels),
					testAccCheckAciSubnetPoolforIpSecTunnelsAttributes(fv_tenant_name, aci_cloud_external_network_vpn_network, aci_cloud_ipsec_tunnel_subnet_pool, description, &subnet_poolfor_ip_sec_tunnels),
				),
			},
			{
				Config: testAccCheckAciSubnetPoolforIpSecTunnelsConfig_basic(fv_tenant_name, aci_cloud_external_network_vpn_network, aci_cloud_ipsec_tunnel_subnet_pool),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetPoolforIpSecTunnelsExists("aci_subnet_poolfor_ip_sec_tunnels.foosubnet_poolfor_ip_sec_tunnels", &subnet_poolfor_ip_sec_tunnels),
					testAccCheckAciSubnetPoolforIpSecTunnelsAttributes(fv_tenant_name, aci_cloud_external_network_vpn_network, aci_cloud_ipsec_tunnel_subnet_pool, description, &subnet_poolfor_ip_sec_tunnels),
				),
			},
		},
	})
}

func testAccCheckAciSubnetPoolforIpSecTunnelsConfig_basic(fv_tenant_name, aci_cloud_external_network_vpn_network, aci_cloud_ipsec_tunnel_subnet_pool string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_infra_network_template" "fooinfra_network_template" {
		name 		= "%s"
		description = "infra_network_template created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	resource "aci_subnet_poolfor_ip_sec_tunnels" "foosubnet_poolfor_ip_sec_tunnels" {
		name 		= "%s"
		description = "subnet_poolfor_ip_sec_tunnels created while acceptance testing"
		infra_network_template_dn = aci_infra_network_template.fooinfra_network_template.id
	}

	`, fv_tenant_name, aci_cloud_external_network_vpn_network, aci_cloud_ipsec_tunnel_subnet_pool)
}

func testAccCheckAciSubnetPoolforIpSecTunnelsExists(name string, subnet_poolfor_ip_sec_tunnels *models.SubnetPoolforIpSecTunnels) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Subnet Pool for IpSec Tunnels %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subnet Pool for IpSec Tunnels dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		subnet_poolfor_ip_sec_tunnelsFound := models.SubnetPoolforIpSecTunnelsFromContainer(cont)
		if subnet_poolfor_ip_sec_tunnelsFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Subnet Pool for IpSec Tunnels %s not found", rs.Primary.ID)
		}
		*subnet_poolfor_ip_sec_tunnels = *subnet_poolfor_ip_sec_tunnelsFound
		return nil
	}
}

func testAccCheckAciSubnetPoolforIpSecTunnelsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_subnet_poolfor_ip_sec_tunnels" {
			cont, err := client.Get(rs.Primary.ID)
			subnet_poolfor_ip_sec_tunnels := models.SubnetPoolforIpSecTunnelsFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Subnet Pool for IpSec Tunnels %s Still exists", subnet_poolfor_ip_sec_tunnels.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSubnetPoolforIpSecTunnelsAttributes(fv_tenant_name, aci_cloud_external_network_vpn_network, aci_cloud_ipsec_tunnel_subnet_pool, description string, subnet_poolfor_ip_sec_tunnels *models.SubnetPoolforIpSecTunnels) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aci_cloud_ipsec_tunnel_subnet_pool != GetMOName(subnet_poolfor_ip_sec_tunnels.DistinguishedName) {
			return fmt.Errorf("Bad cloudtemplate_ip_sec_tunnel_subnet_pool %s", GetMOName(subnet_poolfor_ip_sec_tunnels.DistinguishedName))
		}

		if aci_cloud_external_network_vpn_network != GetMOName(GetParentDn(subnet_poolfor_ip_sec_tunnels.DistinguishedName, models.RncloudtemplateIpSecTunnelSubnetPool)) {
			return fmt.Errorf(" Bad cloudtemplate_infra_network %s", GetMOName(GetParentDn(subnet_poolfor_ip_sec_tunnels.DistinguishedName, models.RncloudtemplateIpSecTunnelSubnetPool)))
		}
		if description != subnet_poolfor_ip_sec_tunnels.Description {
			return fmt.Errorf("Bad subnet_poolfor_ip_sec_tunnels Description %s", subnet_poolfor_ip_sec_tunnels.Description)
		}
		return nil
	}
}
