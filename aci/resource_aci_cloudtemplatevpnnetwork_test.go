package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciCloudTemplateforVPNNetwork_Basic(t *testing.T) {
	var templatefor_vpn_network models.TemplateforVPNNetwork
	fv_tenant_name := acctest.RandString(5)
	cloudtemplate_infra_network_name := acctest.RandString(5)
	cloudtemplate_ext_network_name := acctest.RandString(5)
	cloudtemplate_vpn_network_name := acctest.RandString(5)
	description := "templatefor_vpn_network created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudTemplateforVPNNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudTemplateforVPNNetworkConfig_basic(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, cloudtemplate_vpn_network_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudTemplateforVPNNetworkExists("aci_cloud_external_network_vpn_network.footemplatefor_vpn_network", &templatefor_vpn_network),
					testAccCheckAciCloudTemplateforVPNNetworkAttributes(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, cloudtemplate_vpn_network_name, description, &templatefor_vpn_network),
				),
			},
		},
	})
}

func TestAccAciCloudTemplateforVPNNetwork_Update(t *testing.T) {
	var templatefor_vpn_network models.TemplateforVPNNetwork
	fv_tenant_name := acctest.RandString(5)
	cloudtemplate_infra_network_name := acctest.RandString(5)
	cloudtemplate_ext_network_name := acctest.RandString(5)
	cloudtemplate_vpn_network_name := acctest.RandString(5)
	description := "templatefor_vpn_network created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudTemplateforVPNNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudTemplateforVPNNetworkConfig_basic(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, cloudtemplate_vpn_network_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudTemplateforVPNNetworkExists("aci_cloud_external_network_vpn_network.footemplatefor_vpn_network", &templatefor_vpn_network),
					testAccCheckAciCloudTemplateforVPNNetworkAttributes(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, cloudtemplate_vpn_network_name, description, &templatefor_vpn_network),
				),
			},
			{
				Config: testAccCheckAciCloudTemplateforVPNNetworkConfig_basic(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, cloudtemplate_vpn_network_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudTemplateforVPNNetworkExists("aci_cloud_external_network_vpn_network.footemplatefor_vpn_network", &templatefor_vpn_network),
					testAccCheckAciCloudTemplateforVPNNetworkAttributes(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, cloudtemplate_vpn_network_name, description, &templatefor_vpn_network),
				),
			},
		},
	})
}

func testAccCheckAciCloudTemplateforVPNNetworkConfig_basic(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, cloudtemplate_vpn_network_name string) string {
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

	resource "aci_cloud_external_network" "footemplatefor_external_network" {
		name 		= "%s"
		description = "templatefor_external_network created while acceptance testing"
		infra_network_template_dn = aci_infra_network_template.fooinfra_network_template.id
	}

	resource "aci_cloud_external_network_vpn_network" "footemplatefor_vpn_network" {
		name 		= "%s"
		description = "templatefor_vpn_network created while acceptance testing"
		aci_cloud_external_network_dn = aci_cloud_external_network.footemplatefor_external_network.id
	}

	`, fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, cloudtemplate_vpn_network_name)
}

func testAccCheckAciCloudTemplateforVPNNetworkExists(name string, templatefor_vpn_network *models.TemplateforVPNNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Template for VPN Network %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Template for VPN Network dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		templatefor_vpn_networkFound := models.TemplateforVPNNetworkFromContainer(cont)
		if templatefor_vpn_networkFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Template for VPN Network %s not found", rs.Primary.ID)
		}
		*templatefor_vpn_network = *templatefor_vpn_networkFound
		return nil
	}
}

func testAccCheckAciCloudTemplateforVPNNetworkDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_external_network_vpn_network" {
			cont, err := client.Get(rs.Primary.ID)
			templatefor_vpn_network := models.TemplateforVPNNetworkFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Template for VPN Network %s Still exists", templatefor_vpn_network.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudTemplateforVPNNetworkAttributes(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, cloudtemplate_vpn_network_name, description string, templatefor_vpn_network *models.TemplateforVPNNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloudtemplate_vpn_network_name != GetMOName(templatefor_vpn_network.DistinguishedName) {
			return fmt.Errorf("Bad cloudtemplate_vpn_network %s", GetMOName(templatefor_vpn_network.DistinguishedName))
		}

		if cloudtemplate_ext_network_name != GetMOName(GetParentDn(templatefor_vpn_network.DistinguishedName)) {
			return fmt.Errorf(" Bad cloudtemplate_ext_network %s", GetMOName(GetParentDn(templatefor_vpn_network.DistinguishedName)))
		}
		if description != templatefor_vpn_network.Description {
			return fmt.Errorf("Bad templatefor_vpn_network Description %s", templatefor_vpn_network.Description)
		}
		return nil
	}
}
