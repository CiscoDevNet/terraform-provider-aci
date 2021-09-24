package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciAccessSwitchPolicyGroup_Basic(t *testing.T) {
	var access_switch_policy_group models.AccessSwitchPolicyGroup
	description := "access_switch_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessSwitchPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessSwitchPolicyGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessSwitchPolicyGroupExists("aci_access_switch_policy_group.test", &access_switch_policy_group),
					testAccCheckAciAccessSwitchPolicyGroupAttributes(description, &access_switch_policy_group),
				),
			},
		},
	})
}

func TestAccAciAccessSwitchPolicyGroup_update(t *testing.T) {
	var access_switch_policy_group models.AccessSwitchPolicyGroup
	description := "access_switch_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessSwitchPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessSwitchPolicyGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessSwitchPolicyGroupExists("aci_access_switch_policy_group.test", &access_switch_policy_group),
					testAccCheckAciAccessSwitchPolicyGroupAttributes(description, &access_switch_policy_group),
				),
			},
			{
				Config: testAccCheckAciAccessSwitchPolicyGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessSwitchPolicyGroupExists("aci_access_switch_policy_group.test", &access_switch_policy_group),
					testAccCheckAciAccessSwitchPolicyGroupAttributes(description, &access_switch_policy_group),
				),
			},
		},
	})
}

func testAccCheckAciAccessSwitchPolicyGroupConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_access_switch_policy_group" "test" {
		name 		= "test"
		description = "%s"
		name_alias = "test_alias"
  		annotation = "test_annotation"
	}
	`, description)
}

func testAccCheckAciAccessSwitchPolicyGroupExists(name string, access_switch_policy_group *models.AccessSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Access Switch Policy Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access Switch Policy Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		access_switch_policy_groupFound := models.AccessSwitchPolicyGroupFromContainer(cont)
		if access_switch_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Access Switch Policy Group %s not found", rs.Primary.ID)
		}
		*access_switch_policy_group = *access_switch_policy_groupFound
		return nil
	}
}

func testAccCheckAciAccessSwitchPolicyGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_access_switch_policy_group" {
			cont, err := client.Get(rs.Primary.ID)
			access_switch_policy_group := models.AccessSwitchPolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Access Switch Policy Group %s Still exists", access_switch_policy_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAccessSwitchPolicyGroupAttributes(description string, access_switch_policy_group *models.AccessSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != GetMOName(access_switch_policy_group.DistinguishedName) {
			return fmt.Errorf("Bad infra_acc_node_p_grp %s", GetMOName(access_switch_policy_group.DistinguishedName))
		}

		if description != access_switch_policy_group.Description {
			return fmt.Errorf("Bad access_switch_policy_group Description %s", access_switch_policy_group.Description)
		}

		if "test_alias" != access_switch_policy_group.NameAlias {
			return fmt.Errorf("Bad access_switch_policy_group NameAlias %s", access_switch_policy_group.NameAlias)
		}

		if "test_annotation" != access_switch_policy_group.Annotation {
			return fmt.Errorf("Bad access_switch_policy_group Annotation %s", access_switch_policy_group.Annotation)
		}
		return nil
	}
}
