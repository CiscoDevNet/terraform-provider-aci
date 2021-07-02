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

func TestAccAciEndpointSecurityGroup_Basic(t *testing.T) {
	var endpoint_security_group models.EndpointSecurityGroup
	fv_tenant_name := acctest.RandString(5)
	fv_ap_name := acctest.RandString(5)
	fv_e_sg_name := acctest.RandString(5)
	description := "endpoint_security_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciEndpointSecurityGroupConfig_basic(fv_tenant_name, fv_ap_name, fv_e_sg_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists("aci_endpoint_security_group.fooendpoint_security_group", &endpoint_security_group),
					testAccCheckAciEndpointSecurityGroupAttributes(fv_tenant_name, fv_ap_name, fv_e_sg_name, description, &endpoint_security_group),
				),
			},
		},
	})
}

func testAccCheckAciEndpointSecurityGroupConfig_basic(fv_tenant_name, fv_ap_name, fv_e_sg_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_application_profile" "fooapplication_profile" {
		name 		= "%s"
		description = "application_profile created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	resource "aci_endpoint_security_group" "fooendpoint_security_group" {
		name 		= "%s"
		description = "endpoint_security_group created while acceptance testing"
		application_profile_dn = aci_application_profile.fooapplication_profile.id
	}

	`, fv_tenant_name, fv_ap_name, fv_e_sg_name)
}

func testAccCheckAciEndpointSecurityGroupExists(name string, endpoint_security_group *models.EndpointSecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Endpoint Security Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Endpoint Security Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		endpoint_security_groupFound := models.EndpointSecurityGroupFromContainer(cont)
		if endpoint_security_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Endpoint Security Group %s not found", rs.Primary.ID)
		}
		*endpoint_security_group = *endpoint_security_groupFound
		return nil
	}
}

func testAccCheckAciEndpointSecurityGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_endpoint_security_group" {
			cont, err := client.Get(rs.Primary.ID)
			endpoint_security_group := models.EndpointSecurityGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Endpoint Security Group %s Still exists", endpoint_security_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciEndpointSecurityGroupAttributes(fv_tenant_name, fv_ap_name, fv_e_sg_name, description string, endpoint_security_group *models.EndpointSecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if description != endpoint_security_group.Description {
			return fmt.Errorf("Bad endpoint_security_group Description %s", endpoint_security_group.Description)
		}
		return nil
	}
}
