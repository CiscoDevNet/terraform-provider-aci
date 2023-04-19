package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciNodeBlock_Basic_Firmware(t *testing.T) {
	var node_block models.NodeBlockFW
	description := "node_block created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciNodeBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciNodeBlockConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockExists("aci_node_block_firmware.foonode_block_firmware", &node_block),
					testAccCheckAciNodeBlockAttributes(description, &node_block),
				),
			},
		},
	})
}

func TestAccAciNodeBlock_update(t *testing.T) {
	var node_block models.NodeBlockFW
	description := "node_block created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciNodeBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciNodeBlockConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockExists("aci_node_block_firmware.foonode_block_firmware", &node_block),
					testAccCheckAciNodeBlockAttributes(description, &node_block),
				),
			},
			{
				Config: testAccCheckAciNodeBlockConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockExists("aci_node_block_firmware.foonode_block_firmware", &node_block),
					testAccCheckAciNodeBlockAttributes(description, &node_block),
				),
			},
		},
	})
}

func testAccCheckAciNodeBlockConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_firmware_group" "example" {
		name                 = "example"
		annotation           = "example"
		description          = "from terraform"
		name_alias           = "example"
		firmware_group_type  = "range"
	}

	resource "aci_node_block_firmware" "foonode_block_firmware" {
		firmware_group_dn = aci_firmware_group.example.id
		description       = "%s"
		name              = "crest_test_vishwa"
		annotation        = "example"
		from_             = "1"
		name_alias        = "example"
		to_               = "5"
	}
	`, description)
}

func testAccCheckAciNodeBlockExists(name string, node_block *models.NodeBlockFW) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Node Block %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Node Block dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		node_blockFound := models.NodeBlockFromContainer(cont)
		if node_blockFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Node Block %s not found", rs.Primary.ID)
		}
		*node_block = *node_blockFound
		return nil
	}
}

func testAccCheckAciNodeBlockDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_node_block_firmware" {
			cont, err := client.Get(rs.Primary.ID)
			node_block := models.NodeBlockFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Node Block %s Still exists", node_block.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciNodeBlockAttributes(description string, node_block *models.NodeBlockFW) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != node_block.Description {
			return fmt.Errorf("Bad node_block Description %s", node_block.Description)
		}

		if "crest_test_vishwa" != node_block.Name {
			return fmt.Errorf("Bad node_block name %s", node_block.Name)
		}

		if "example" != node_block.Annotation {
			return fmt.Errorf("Bad node_block annotation %s", node_block.Annotation)
		}

		if "1" != node_block.From_ {
			return fmt.Errorf("Bad node_block from_ %s", node_block.From_)
		}

		if "example" != node_block.NameAlias {
			return fmt.Errorf("Bad node_block name_alias %s", node_block.NameAlias)
		}

		if "5" != node_block.To_ {
			return fmt.Errorf("Bad node_block to_ %s", node_block.To_)
		}

		return nil
	}
}
