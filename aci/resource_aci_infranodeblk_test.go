package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciInfraNodeBlock_Basic(t *testing.T) {
	var node_block models.NodeBlock
	description := "node_block created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciNodeBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInfraNodeBlockConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInfraNodeBlockExists("aci_node_block.foonode_block", &node_block),
					testAccCheckAciInfraNodeBlockAttributes(description, &node_block),
				),
			},
			{
				ResourceName:      "aci_node_block",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciInfraNodeBlock_update(t *testing.T) {
	var node_block models.NodeBlock
	description := "node_block created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciNodeBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInfraNodeBlockConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInfraNodeBlockExists("aci_node_block.foonode_block", &node_block),
					testAccCheckAciInfraNodeBlockAttributes(description, &node_block),
				),
			},
			{
				Config: testAccCheckAciNodeBlockConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockExists("aci_node_block.foonode_block", &node_block),
					testAccCheckAciNodeBlockAttributes(description, &node_block),
				),
			},
		},
	})
}

func testAccCheckAciInfraNodeBlockConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_node_block" "foonode_block" {
		  switch_association_dn  = "${aci_switch_association.example.id}"
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  from_  = "105"
		  name_alias  = "example"
		  to_  = "106"
		}
	`, description)
}

func testAccCheckAciInfraNodeBlockExists(name string, node_block *models.NodeBlock) resource.TestCheckFunc {
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

func testAccCheckAciInfraNodeBlockDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_node_block" {
			cont, err := client.Get(rs.Primary.ID)
			node_block := models.NodeBlockFromContainerBLK(cont)
			if err == nil {
				return fmt.Errorf("Node Block %s Still exists", node_block.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciInfraNodeBlockAttributes(description string, node_block *models.NodeBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != node_block.Description {
			return fmt.Errorf("Bad node_block Description %s", node_block.Description)
		}

		if "example" != node_block.Name {
			return fmt.Errorf("Bad node_block name %s", node_block.Name)
		}

		if "example" != node_block.Annotation {
			return fmt.Errorf("Bad node_block annotation %s", node_block.Annotation)
		}

		if "105" != node_block.From_ {
			return fmt.Errorf("Bad node_block from_ %s", node_block.From_)
		}

		if "example" != node_block.NameAlias {
			return fmt.Errorf("Bad node_block name_alias %s", node_block.NameAlias)
		}

		if "106" != node_block.To_ {
			return fmt.Errorf("Bad node_block to_ %s", node_block.To_)
		}

		return nil
	}
}
