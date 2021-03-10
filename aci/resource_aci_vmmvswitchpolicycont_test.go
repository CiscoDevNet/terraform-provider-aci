package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciVSwitchPolicyGroup_Basic(t *testing.T) {
	var v_switch_policy_group models.VSwitchPolicyGroup
	vmm_prov_p_name := acctest.RandString(5)
	vmm_dom_p_name := acctest.RandString(5)
	vmm_v_switch_policy_cont_name := acctest.RandString(5)
	description := "v_switch_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVSwitchPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVSwitchPolicyGroupConfig_basic(vmm_prov_p_name, vmm_dom_p_name, vmm_v_switch_policy_cont_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVSwitchPolicyGroupExists("aci_v_switch_policy_group.foov_switch_policy_group", &v_switch_policy_group),
					testAccCheckAciVSwitchPolicyGroupAttributes(vmm_prov_p_name, vmm_dom_p_name, vmm_v_switch_policy_cont_name, description, &v_switch_policy_group),
				),
			},
		},
	})
}

func testAccCheckAciVSwitchPolicyGroupConfig_basic(vmm_prov_p_name, vmm_dom_p_name, vmm_v_switch_policy_cont_name string) string {
	return fmt.Sprintf(`

	resource "aci_provider_profile" "fooprovider_profile" {
		name 		= "%s"
		description = "provider_profile created while acceptance testing"

	}

	resource "aci_vmm_domain" "foovmm_domain" {
		name 		= "%s"
		description = "vmm_domain created while acceptance testing"
		provider_profile_dn = "${aci_provider_profile.fooprovider_profile.id}"
	}

	resource "aci_v_switch_policy_group" "foov_switch_policy_group" {
		name 		= "%s"
		description = "v_switch_policy_group created while acceptance testing"
		vmm_domain_dn = "${aci_vmm_domain.foovmm_domain.id}"
	}

	`, vmm_prov_p_name, vmm_dom_p_name, vmm_v_switch_policy_cont_name)
}

func testAccCheckAciVSwitchPolicyGroupExists(name string, v_switch_policy_group *models.VSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VSwitch Policy Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VSwitch Policy Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		v_switch_policy_groupFound := models.VSwitchPolicyGroupFromContainer(cont)
		if v_switch_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VSwitch Policy Group %s not found", rs.Primary.ID)
		}
		*v_switch_policy_group = *v_switch_policy_groupFound
		return nil
	}
}

func testAccCheckAciVSwitchPolicyGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_v_switch_policy_group" {
			cont, err := client.Get(rs.Primary.ID)
			v_switch_policy_group := models.VSwitchPolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VSwitch Policy Group %s Still exists", v_switch_policy_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciVSwitchPolicyGroupAttributes(vmm_prov_p_name, vmm_dom_p_name, vmm_v_switch_policy_cont_name, description string, v_switch_policy_group *models.VSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vmm_v_switch_policy_cont_name != GetMOName(v_switch_policy_group.DistinguishedName) {
			return fmt.Errorf("Bad vmm_v_switch_policy_cont %s", GetMOName(v_switch_policy_group.DistinguishedName))
		}

		if vmm_dom_p_name != GetMOName(GetParentDn(v_switch_policy_group.DistinguishedName)) {
			return fmt.Errorf(" Bad vmm_dom_p %s", GetMOName(GetParentDn(v_switch_policy_group.DistinguishedName)))
		}
		if description != v_switch_policy_group.Description {
			return fmt.Errorf("Bad v_switch_policy_group Description %s", v_switch_policy_group.Description)
		}

		return nil
	}
}
