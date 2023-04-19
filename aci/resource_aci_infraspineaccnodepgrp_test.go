package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSpineSwitchPolicyGroup_Basic(t *testing.T) {
	var spine_switch_policy_group models.SpineSwitchPolicyGroup
	description := "spine_switch_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpineSwitchPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpineSwitchPolicyGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineSwitchPolicyGroupExists("aci_spine_switch_policy_group.test", &spine_switch_policy_group),
					testAccCheckAciSpineSwitchPolicyGroupAttributes(description, &spine_switch_policy_group),
				),
			},
		},
	})
}

func TestAccAciSpineSwitchPolicyGroup_update(t *testing.T) {
	var spine_switch_policy_group models.SpineSwitchPolicyGroup
	description := "spine_switch_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpineSwitchPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpineSwitchPolicyGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineSwitchPolicyGroupExists("aci_spine_switch_policy_group.test", &spine_switch_policy_group),
					testAccCheckAciSpineSwitchPolicyGroupAttributes(description, &spine_switch_policy_group),
				),
			},
			{
				Config: testAccCheckAciSpineSwitchPolicyGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineSwitchPolicyGroupExists("aci_spine_switch_policy_group.test", &spine_switch_policy_group),
					testAccCheckAciSpineSwitchPolicyGroupAttributes(description, &spine_switch_policy_group),
				),
			},
		},
	})
}

func testAccCheckAciSpineSwitchPolicyGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_spine_switch_policy_group" "test" {
		name 		= "test"
		description = "%s"
		annotation = "test_annotation"
		name_alias = "test_name_alias"
	}

	`, description)
}

func testAccCheckAciSpineSwitchPolicyGroupExists(name string, spine_switch_policy_group *models.SpineSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spine Switch Policy Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spine Switch Policy Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spine_switch_policy_groupFound := models.SpineSwitchPolicyGroupFromContainer(cont)
		if spine_switch_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spine Switch Policy Group %s not found", rs.Primary.ID)
		}
		*spine_switch_policy_group = *spine_switch_policy_groupFound
		return nil
	}
}

func testAccCheckAciSpineSwitchPolicyGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_spine_switch_policy_group" {
			cont, err := client.Get(rs.Primary.ID)
			spine_switch_policy_group := models.SpineSwitchPolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spine Switch Policy Group %s Still exists", spine_switch_policy_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSpineSwitchPolicyGroupAttributes(description string, spine_switch_policy_group *models.SpineSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != GetMOName(spine_switch_policy_group.DistinguishedName) {
			return fmt.Errorf("Bad infra_spine_acc_node_p_grp %s", GetMOName(spine_switch_policy_group.DistinguishedName))
		}

		if description != spine_switch_policy_group.Description {
			return fmt.Errorf("Bad spine_switch_policy_group Description %s", spine_switch_policy_group.Description)
		}

		if "test_annotation" != spine_switch_policy_group.Annotation {
			return fmt.Errorf("Bad spine_switch_policy_group Annotation %s", spine_switch_policy_group.Annotation)
		}

		if "test_name_alias" != spine_switch_policy_group.NameAlias {
			return fmt.Errorf("Bad spine_switch_policy_group NameAlias %s", spine_switch_policy_group.NameAlias)
		}
		return nil
	}
}
