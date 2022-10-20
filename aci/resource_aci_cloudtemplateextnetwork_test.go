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

func TestAccAciCloudTemplateforExternalNetwork_Basic(t *testing.T) {
	var templatefor_external_network models.CloudTemplateforExternalNetwork
	fv_tenant_name := acctest.RandString(5)
	cloudtemplate_infra_network_name := acctest.RandString(5)
	cloudtemplate_ext_network_name := acctest.RandString(5)
	description := "templatefor_external_network created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudTemplateforExternalNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudTemplateforExternalNetworkConfig_basic(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudTemplateforExternalNetworkExists("aci_cloud_external_network.footemplatefor_external_network", &templatefor_external_network),
					testAccCheckAciCloudTemplateforExternalNetworkAttributes(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, description, &templatefor_external_network),
				),
			},
		},
	})
}

func TestAccAciCloudTemplateforExternalNetwork_Update(t *testing.T) {
	var templatefor_external_network models.CloudTemplateforExternalNetwork
	fv_tenant_name := acctest.RandString(5)
	cloudtemplate_infra_network_name := acctest.RandString(5)
	cloudtemplate_ext_network_name := acctest.RandString(5)
	description := "templatefor_external_network created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudTemplateforExternalNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudTemplateforExternalNetworkConfig_basic(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudTemplateforExternalNetworkExists("aci_cloud_external_network.footemplatefor_external_network", &templatefor_external_network),
					testAccCheckAciCloudTemplateforExternalNetworkAttributes(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, description, &templatefor_external_network),
				),
			},
			{
				Config: testAccCheckAciCloudTemplateforExternalNetworkConfig_basic(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudTemplateforExternalNetworkExists("aci_cloud_external_network.footemplatefor_external_network", &templatefor_external_network),
					testAccCheckAciCloudTemplateforExternalNetworkAttributes(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, description, &templatefor_external_network),
				),
			},
		},
	})
}

func testAccCheckAciCloudTemplateforExternalNetworkConfig_basic(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name string) string {
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

	`, fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name)
}

func testAccCheckAciCloudTemplateforExternalNetworkExists(name string, templatefor_external_network *models.CloudTemplateforExternalNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Template for External Network %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Template for External Network dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		templatefor_external_networkFound := models.CloudTemplateforExternalNetworkFromContainer(cont)
		if templatefor_external_networkFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Template for External Network %s not found", rs.Primary.ID)
		}
		*templatefor_external_network = *templatefor_external_networkFound
		return nil
	}
}

func testAccCheckAciCloudTemplateforExternalNetworkDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cloud_external_network" {
			cont, err := client.Get(rs.Primary.ID)
			templatefor_external_network := models.CloudTemplateforExternalNetworkFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Template for External Network %s Still exists", templatefor_external_network.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCloudTemplateforExternalNetworkAttributes(fv_tenant_name, cloudtemplate_infra_network_name, cloudtemplate_ext_network_name, description string, templatefor_external_network *models.CloudTemplateforExternalNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloudtemplate_ext_network_name != GetMOName(templatefor_external_network.DistinguishedName) {
			return fmt.Errorf("Bad cloudtemplate_ext_network %s", GetMOName(templatefor_external_network.DistinguishedName))
		}

		if cloudtemplate_infra_network_name != GetMOName(GetParentDn(templatefor_external_network.DistinguishedName)) {
			return fmt.Errorf(" Bad cloudtemplate_infra_network %s", GetMOName(GetParentDn(templatefor_external_network.DistinguishedName)))
		}
		if description != templatefor_external_network.Description {
			return fmt.Errorf("Bad templatefor_external_network Description %s", templatefor_external_network.Description)
		}
		return nil
	}
}
