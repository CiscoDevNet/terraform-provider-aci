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
	resource "aci_vmm_domain" "foovmm_domain" {
		name 		= "%s"
		provider_profile_dn = "uni/vmmp-VMware"
	}

	resource "aci_vmm_controller" "foovmm_controller" {
		name 		= "%s"
		vmm_domain_dn = aci_vmm_domain.foovmm_domain.id
		annotation = "orchestrator:terraform"
		dvs_version = "unmanaged"
		host_or_ip = "10.10.10.10"
		inventory_trig_st = "untriggered"
		mode = "default"
		msft_config_err_msg = "Error"
		n1kv_stats_mode = "enabled"
		port = "0"
		root_cont_name = "vmmdc"
		scope = "vm"
		seq_num = "0"
		stats_mode = "disabled"
		vxlan_depl_pref = "vxlan"
	}

	`, vmm_dom_p_name, vmm_ctrlr_p_name)
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
		return nil
	}
}
