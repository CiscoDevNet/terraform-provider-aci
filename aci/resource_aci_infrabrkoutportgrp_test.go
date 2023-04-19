package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLeafBreakoutPortGroup_Basic(t *testing.T) {
	var leaf_breakout_port_group models.LeafBreakoutPortGroup
	description := "leaf_breakout_port_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLeafBreakoutPortGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLeafBreakoutPortGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafBreakoutPortGroupExists("aci_leaf_breakout_port_group.test", &leaf_breakout_port_group),
					testAccCheckAciLeafBreakoutPortGroupAttributes(description, &leaf_breakout_port_group),
				),
			},
		},
	})
}

func TestAccAciLeafBreakoutPortGroup_update(t *testing.T) {
	var leaf_breakout_port_group models.LeafBreakoutPortGroup
	description := "leaf_breakout_port_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLeafBreakoutPortGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLeafBreakoutPortGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafBreakoutPortGroupExists("aci_leaf_breakout_port_group.test", &leaf_breakout_port_group),
					testAccCheckAciLeafBreakoutPortGroupAttributes(description, &leaf_breakout_port_group),
				),
			},
			{
				Config: testAccCheckAciLeafBreakoutPortGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafBreakoutPortGroupExists("aci_leaf_breakout_port_group.test", &leaf_breakout_port_group),
					testAccCheckAciLeafBreakoutPortGroupAttributes(description, &leaf_breakout_port_group),
				),
			},
		},
	})
}

func testAccCheckAciLeafBreakoutPortGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_leaf_breakout_port_group" "test" {
		description = "%s"
		name  = "test"
		annotation  = "example"
		brkout_map  = "100g-4x"
		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciLeafBreakoutPortGroupExists(name string, leaf_breakout_port_group *models.LeafBreakoutPortGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Leaf Breakout Port Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Leaf Breakout Port Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		leaf_breakout_port_groupFound := models.LeafBreakoutPortGroupFromContainer(cont)
		if leaf_breakout_port_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Leaf Breakout Port Group %s not found", rs.Primary.ID)
		}
		*leaf_breakout_port_group = *leaf_breakout_port_groupFound
		return nil
	}
}

func testAccCheckAciLeafBreakoutPortGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_leaf_breakout_port_group" {
			cont, err := client.Get(rs.Primary.ID)
			leaf_breakout_port_group := models.LeafBreakoutPortGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Leaf Breakout Port Group %s Still exists", leaf_breakout_port_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLeafBreakoutPortGroupAttributes(description string, leaf_breakout_port_group *models.LeafBreakoutPortGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != leaf_breakout_port_group.Description {
			return fmt.Errorf("Bad leaf_breakout_port_group Description %s", leaf_breakout_port_group.Description)
		}

		if "test" != leaf_breakout_port_group.Name {
			return fmt.Errorf("Bad leaf_breakout_port_group name %s", leaf_breakout_port_group.Name)
		}

		if "example" != leaf_breakout_port_group.Annotation {
			return fmt.Errorf("Bad leaf_breakout_port_group annotation %s", leaf_breakout_port_group.Annotation)
		}

		if "100g-4x" != leaf_breakout_port_group.BrkoutMap {
			return fmt.Errorf("Bad leaf_breakout_port_group brkout_map %s", leaf_breakout_port_group.BrkoutMap)
		}

		if "example" != leaf_breakout_port_group.NameAlias {
			return fmt.Errorf("Bad leaf_breakout_port_group name_alias %s", leaf_breakout_port_group.NameAlias)
		}

		return nil
	}
}
