package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLeafAccessPortPolicyGroup_Basic(t *testing.T) {
	var leaf_access_port_policy_group models.LeafAccessPortPolicyGroup
	description := "leaf_access_port_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLeafAccessPortPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLeafAccessPortPolicyGroupConfig_basic(description, "alias_port"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafAccessPortPolicyGroupExists("aci_leaf_access_port_policy_group.fooleaf_access_port_policy_group", &leaf_access_port_policy_group),
					testAccCheckAciLeafAccessPortPolicyGroupAttributes(description, "alias_port", &leaf_access_port_policy_group),
				),
			},
			{
				ResourceName:      "aci_leaf_access_port_policy_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciLeafAccessPortPolicyGroup_update(t *testing.T) {
	var leaf_access_port_policy_group models.LeafAccessPortPolicyGroup
	description := "leaf_access_port_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLeafAccessPortPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLeafAccessPortPolicyGroupConfig_basic(description, "alias_port"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafAccessPortPolicyGroupExists("aci_leaf_access_port_policy_group.fooleaf_access_port_policy_group", &leaf_access_port_policy_group),
					testAccCheckAciLeafAccessPortPolicyGroupAttributes(description, "alias_port", &leaf_access_port_policy_group),
				),
			},
			{
				Config: testAccCheckAciLeafAccessPortPolicyGroupConfig_basic(description, "alias_updated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafAccessPortPolicyGroupExists("aci_leaf_access_port_policy_group.fooleaf_access_port_policy_group", &leaf_access_port_policy_group),
					testAccCheckAciLeafAccessPortPolicyGroupAttributes(description, "alias_updated", &leaf_access_port_policy_group),
				),
			},
		},
	})
}

func testAccCheckAciLeafAccessPortPolicyGroupConfig_basic(description, name_alias string) string {
	return fmt.Sprintf(`

	resource "aci_leaf_access_port_policy_group" "fooleaf_access_port_policy_group" {
		description = "%s"
		name        = "demo_access_port"
		annotation  = "tag_ports"
		name_alias  = "%s"
	}  
	`, description, name_alias)
}

func testAccCheckAciLeafAccessPortPolicyGroupExists(name string, leaf_access_port_policy_group *models.LeafAccessPortPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Leaf Access Port Policy Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Leaf Access Port Policy Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		leaf_access_port_policy_groupFound := models.LeafAccessPortPolicyGroupFromContainer(cont)
		if leaf_access_port_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Leaf Access Port Policy Group %s not found", rs.Primary.ID)
		}
		*leaf_access_port_policy_group = *leaf_access_port_policy_groupFound
		return nil
	}
}

func testAccCheckAciLeafAccessPortPolicyGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_leaf_access_port_policy_group" {
			cont, err := client.Get(rs.Primary.ID)
			leaf_access_port_policy_group := models.LeafAccessPortPolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Leaf Access Port Policy Group %s Still exists", leaf_access_port_policy_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLeafAccessPortPolicyGroupAttributes(description, name_alias string, leaf_access_port_policy_group *models.LeafAccessPortPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != leaf_access_port_policy_group.Description {
			return fmt.Errorf("Bad leaf_access_port_policy_group Description %s", leaf_access_port_policy_group.Description)
		}

		if "demo_access_port" != leaf_access_port_policy_group.Name {
			return fmt.Errorf("Bad leaf_access_port_policy_group name %s", leaf_access_port_policy_group.Name)
		}

		if "tag_ports" != leaf_access_port_policy_group.Annotation {
			return fmt.Errorf("Bad leaf_access_port_policy_group annotation %s", leaf_access_port_policy_group.Annotation)
		}

		if name_alias != leaf_access_port_policy_group.NameAlias {
			return fmt.Errorf("Bad leaf_access_port_policy_group name_alias %s", leaf_access_port_policy_group.NameAlias)
		}

		return nil
	}
}
