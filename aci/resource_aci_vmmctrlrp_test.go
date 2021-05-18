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

func TestAccAciVMMController_Basic(t *testing.T) {
	var vmm_controller models.VMMController
	vmm_prov_p_name := acctest.RandString(5)
	vmm_dom_p_name := acctest.RandString(5)
	vmm_ctrlr_p_name := acctest.RandString(5)
	description := "vmm_controller created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVMMControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVMMControllerConfig_basic(vmm_prov_p_name, vmm_dom_p_name, vmm_ctrlr_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMControllerExists("aci_vmm_controller.foovmm_controller", &vmm_controller),
					testAccCheckAciVMMControllerAttributes(vmm_prov_p_name, vmm_dom_p_name, vmm_ctrlr_p_name, description, &vmm_controller),
				),
			},
		},
	})
}

func testAccCheckAciVMMControllerConfig_basic(vmm_prov_p_name, vmm_dom_p_name, vmm_ctrlr_p_name string) string {
	return fmt.Sprintf(`

	resource "aci_provider_profile" "fooprovider_profile" {
		name 		= "%s"
		description = "provider_profile created while acceptance testing"

	}

	resource "aci_vmm_domain" "foovmm_domain" {
		name 		= "%s"
		description = "vmm_domain created while acceptance testing"
		provider_profile_dn = aci_provider_profile.fooprovider_profile.id
	}

	resource "aci_vmm_controller" "foovmm_controller" {
		name 		= "%s"
		description = "vmm_controller created while acceptance testing"
		vmm_domain_dn = aci_vmm_domain.foovmm_domain.id
	}

	`, vmm_prov_p_name, vmm_dom_p_name, vmm_ctrlr_p_name)
}

func testAccCheckAciVMMControllerExists(name string, vmm_controller *models.VMMController) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VMM Controller %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VMM Controller dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vmm_controllerFound := models.VMMControllerFromContainer(cont)
		if vmm_controllerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VMM Controller %s not found", rs.Primary.ID)
		}
		*vmm_controller = *vmm_controllerFound
		return nil
	}
}

func testAccCheckAciVMMControllerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vmm_controller" {
			cont, err := client.Get(rs.Primary.ID)
			vmm_controller := models.VMMControllerFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VMM Controller %s Still exists", vmm_controller.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVMMControllerAttributes(vmm_prov_p_name, vmm_dom_p_name, vmm_ctrlr_p_name, description string, vmm_controller *models.VMMController) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vmm_ctrlr_p_name != GetMOName(vmm_controller.DistinguishedName) {
			return fmt.Errorf("Bad vmm_ctrlr_p %s", GetMOName(vmm_controller.DistinguishedName))
		}

		if vmm_dom_p_name != GetMOName(GetParentDn(vmm_controller.DistinguishedName)) {
			return fmt.Errorf(" Bad vmm_dom_p %s", GetMOName(GetParentDn(vmm_controller.DistinguishedName)))
		}
		if description != vmm_controller.Description {
			return fmt.Errorf("Bad vmm_controller Description %s", vmm_controller.Description)
		}
		return nil
	}
}
