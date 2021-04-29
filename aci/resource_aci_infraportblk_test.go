package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciAccessPortBlock_Basic(t *testing.T) {
	var access_port_block models.AccessPortBlock
	description := "access_port_block created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessPortBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessPortBlockConfig_basic(description, "1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortBlockExists("aci_access_port_block.fooaccess_port_block", &access_port_block),
					testAccCheckAciAccessPortBlockAttributes(description, "1", &access_port_block),
				),
			},
			{
				ResourceName:      "aci_access_port_block",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciAccessPortBlock_update(t *testing.T) {
	var access_port_block models.AccessPortBlock
	description := "access_port_block created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessPortBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessPortBlockConfig_basic(description, "1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortBlockExists("aci_access_port_block.fooaccess_port_block", &access_port_block),
					testAccCheckAciAccessPortBlockAttributes(description, "1", &access_port_block),
				),
			},
			{
				Config: testAccCheckAciAccessPortBlockConfig_basic(description, "2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortBlockExists("aci_access_port_block.fooaccess_port_block", &access_port_block),
					testAccCheckAciAccessPortBlockAttributes(description, "2", &access_port_block),
				),
			},
		},
	})
}

func testAccCheckAciAccessPortBlockConfig_basic(description, from_port string) string {
	return fmt.Sprintf(`

	resource "aci_access_port_block" "fooaccess_port_block" {
		access_port_selector_dn = "${aci_access_port_selector.example.id}"
		description             = "%s"
		name                    = "demo_port_block"
		annotation              = "tag_port_block"
		from_card               = "1"
		from_port               = "%s"
		name_alias              = "alias_port_block"
		to_card                 = "3"
		to_port                 = "3"
	}
	  
	`, description, from_port)
}

func testAccCheckAciAccessPortBlockExists(name string, access_port_block *models.AccessPortBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Access Port Block %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access Port Block dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		access_port_blockFound := models.AccessPortBlockFromContainer(cont)
		if access_port_blockFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Access Port Block %s not found", rs.Primary.ID)
		}
		*access_port_block = *access_port_blockFound
		return nil
	}
}

func testAccCheckAciAccessPortBlockDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_access_port_block" {
			cont, err := client.Get(rs.Primary.ID)
			access_port_block := models.AccessPortBlockFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Access Port Block %s Still exists", access_port_block.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciAccessPortBlockAttributes(description, from_port string, access_port_block *models.AccessPortBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != access_port_block.Description {
			return fmt.Errorf("Bad access_port_block Description %s", access_port_block.Description)
		}

		if "demo_port_block" != access_port_block.Name {
			return fmt.Errorf("Bad access_port_block name %s", access_port_block.Name)
		}

		if "tag_port_block" != access_port_block.Annotation {
			return fmt.Errorf("Bad access_port_block annotation %s", access_port_block.Annotation)
		}

		if "1" != access_port_block.FromCard {
			return fmt.Errorf("Bad access_port_block from_card %s", access_port_block.FromCard)
		}

		if from_port != access_port_block.FromPort {
			return fmt.Errorf("Bad access_port_block from_port %s", access_port_block.FromPort)
		}

		if "alias_port_block" != access_port_block.NameAlias {
			return fmt.Errorf("Bad access_port_block name_alias %s", access_port_block.NameAlias)
		}

		if "3" != access_port_block.ToCard {
			return fmt.Errorf("Bad access_port_block to_card %s", access_port_block.ToCard)
		}

		if "3" != access_port_block.ToPort {
			return fmt.Errorf("Bad access_port_block to_port %s", access_port_block.ToPort)
		}

		return nil
	}
}
