package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciEndpointSecurityGroupSelector_Basic(t *testing.T) {
	var endpoint_security_group_selector models.EndpointSecurityGroupSelector
	fv_tenant_name := acctest.RandString(5)
	fv_ap_name := acctest.RandString(5)
	fv_e_sg_name := acctest.RandString(5)
	fv_ep_selector_name := acctest.RandString(5)
	description := "endpoint_security_group_selector created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciEndpointSecurityGroupSelectorConfig_basic(fv_tenant_name, fv_ap_name, fv_e_sg_name, fv_ep_selector_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupSelectorExists("aci_endpoint_security_group_selector.fooendpoint_security_group_selector", &endpoint_security_group_selector),
					testAccCheckAciEndpointSecurityGroupSelectorAttributes(fv_tenant_name, fv_ap_name, fv_e_sg_name, fv_ep_selector_name, description, &endpoint_security_group_selector),
				),
			},
		},
	})
}

func testAccCheckAciEndpointSecurityGroupSelectorConfig_basic(fv_tenant_name, fv_ap_name, fv_e_sg_name, fv_ep_selector_name string) string {
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

	resource "aci_endpoint_security_group_selector" "fooendpoint_security_group_selector" {
		name 		= "%s"
		description = "endpoint_security_group_selector created while acceptance testing"
		endpoint_security_group_dn = aci_endpoint_security_group.fooendpoint_security_group.id
	}

	`, fv_tenant_name, fv_ap_name, fv_e_sg_name, fv_ep_selector_name)
}

func testAccCheckAciEndpointSecurityGroupSelectorExists(name string, endpoint_security_group_selector *models.EndpointSecurityGroupSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Endpoint Security Group Selector %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Endpoint Security Group Selector dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		endpoint_security_group_selectorFound := models.EndpointSecurityGroupSelectorFromContainer(cont)
		if endpoint_security_group_selectorFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Endpoint Security Group Selector %s not found", rs.Primary.ID)
		}
		*endpoint_security_group_selector = *endpoint_security_group_selectorFound
		return nil
	}
}

func testAccCheckAciEndpointSecurityGroupSelectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_endpoint_security_group_selector" {
			cont, err := client.Get(rs.Primary.ID)
			endpoint_security_group_selector := models.EndpointSecurityGroupSelectorFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Endpoint Security Group Selector %s Still exists", endpoint_security_group_selector.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciEndpointSecurityGroupSelectorAttributes(fv_tenant_name, fv_ap_name, fv_e_sg_name, fv_ep_selector_name, description string, endpoint_security_group_selector *models.EndpointSecurityGroupSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fv_ep_selector_name != GetMOName(endpoint_security_group_selector.DistinguishedName) {
			return fmt.Errorf("Bad fvep_selector %s", GetMOName(endpoint_security_group_selector.DistinguishedName))
		}

		if fv_e_sg_name != GetMOName(GetParentDn(endpoint_security_group_selector.DistinguishedName)) {
			return fmt.Errorf(" Bad fve_sg %s", GetMOName(GetParentDn(endpoint_security_group_selector.DistinguishedName)))
		}
		if description != endpoint_security_group_selector.Description {
			return fmt.Errorf("Bad endpoint_security_group_selector Description %s", endpoint_security_group_selector.Description)
		}
		return nil
	}
}
