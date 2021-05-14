package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSpineAccessPortPolicyGroup_Basic(t *testing.T) {
	var spine_access_port_policy_group models.SpineAccessPortPolicyGroup
	description := "spine_port_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpineAccessPortPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpineAccessPortPolicyGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineAccessPortPolicyGroupExists("aci_spine_port_policy_group.foospine_access_port_policy_group", &spine_access_port_policy_group),
					testAccCheckAciSpineAccessPortPolicyGroupAttributes(description, &spine_access_port_policy_group),
				),
			},
		},
	})
}

func TestAccAciSpineAccessPortPolicyGroup_update(t *testing.T) {
	var spine_access_port_policy_group models.SpineAccessPortPolicyGroup
	description := "spine_port_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpineAccessPortPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpineAccessPortPolicyGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineAccessPortPolicyGroupExists("aci_spine_port_policy_group.foospine_access_port_policy_group", &spine_access_port_policy_group),
					testAccCheckAciSpineAccessPortPolicyGroupAttributes(description, &spine_access_port_policy_group),
				),
			},
			{
				Config: testAccCheckAciSpineAccessPortPolicyGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineAccessPortPolicyGroupExists("aci_spine_port_policy_group.foospine_access_port_policy_group", &spine_access_port_policy_group),
					testAccCheckAciSpineAccessPortPolicyGroupAttributes(description, &spine_access_port_policy_group),
				),
			},
		},
	})
}

func testAccCheckAciSpineAccessPortPolicyGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_spine_port_policy_group" "foospine_access_port_policy_group" {
		description = "%s"
		name  = "example"
		annotation  = "example"
		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciSpineAccessPortPolicyGroupExists(name string, spine_access_port_policy_group *models.SpineAccessPortPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spine Access Port Policy Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spine Access Port Policy Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spine_access_port_policy_groupFound := models.SpineAccessPortPolicyGroupFromContainer(cont)
		if spine_access_port_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spine Access Port Policy Group %s not found", rs.Primary.ID)
		}
		*spine_access_port_policy_group = *spine_access_port_policy_groupFound
		return nil
	}
}

func testAccCheckAciSpineAccessPortPolicyGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_spine_port_policy_group" {
			cont, err := client.Get(rs.Primary.ID)
			spine_access_port_policy_group := models.SpineAccessPortPolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spine Access Port Policy Group %s Still exists", spine_access_port_policy_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSpineAccessPortPolicyGroupAttributes(description string, spine_access_port_policy_group *models.SpineAccessPortPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != spine_access_port_policy_group.Description {
			return fmt.Errorf("Bad spine_access_port_policy_group Description %s", spine_access_port_policy_group.Description)
		}

		if "example" != spine_access_port_policy_group.Name {
			return fmt.Errorf("Bad spine_access_port_policy_group name %s", spine_access_port_policy_group.Name)
		}

		if "example" != spine_access_port_policy_group.Annotation {
			return fmt.Errorf("Bad spine_access_port_policy_group annotation %s", spine_access_port_policy_group.Annotation)
		}

		if "example" != spine_access_port_policy_group.NameAlias {
			return fmt.Errorf("Bad spine_access_port_policy_group name_alias %s", spine_access_port_policy_group.NameAlias)
		}

		return nil
	}
}
