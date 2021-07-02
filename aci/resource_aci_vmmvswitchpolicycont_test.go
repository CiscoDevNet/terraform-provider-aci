package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciVSwitchPolicyGroup_Basic(t *testing.T) {
	var v_switch_policy_group models.VSwitchPolicyGroup
	description := "v_switch_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVSwitchPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVSwitchPolicyGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVSwitchPolicyGroupExists("aci_vswitch_policy.foov_switch_policy_group", &v_switch_policy_group),
					testAccCheckAciVSwitchPolicyGroupAttributes(description, &v_switch_policy_group),
				),
			},
		},
	})
}

func testAccCheckAciVSwitchPolicyGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_vlan_pool" "vmm_vlan_pool" {
		name       = "vlan_pool_P12"
		alloc_mode = "dynamic"
	  }	  
	  
	  resource "aci_vmm_domain" "ave" {
		provider_profile_dn = "uni/vmmp-VMware"
		relation_infra_rs_vlan_ns = aci_vlan_pool.vmm_vlan_pool.id
		name                = "example"
		enable_ave = "yes"
		mcast_addr = "239.10.10.10"		
	  }
	resource "aci_vswitch_policy" "foov_switch_policy_group" {
		vmm_domain_dn  = aci_vmm_domain.ave.id
		description = "%s"
  		annotation  = "example"
  		name_alias  = "example"
	}
	`, description)
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

		if rs.Type == "aci_vswitch_policy" {
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

func testAccCheckAciVSwitchPolicyGroupAttributes(description string, v_switch_policy_group *models.VSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != v_switch_policy_group.Description {
			return fmt.Errorf("Bad v_switch_policy_group Description %s", v_switch_policy_group.Description)
		}

		if "example" != v_switch_policy_group.Annotation {
			return fmt.Errorf("Bad v_switch_policy_group annotation %s", v_switch_policy_group.Annotation)
		}

		if "example" != v_switch_policy_group.NameAlias {
			return fmt.Errorf("Bad v_switch_policy_group name_alias %s", v_switch_policy_group.NameAlias)
		}

		return nil
	}
}
