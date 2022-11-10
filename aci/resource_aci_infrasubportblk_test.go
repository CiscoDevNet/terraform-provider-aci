package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciAccessSubPortBlock_Basic(t *testing.T) {
	var access_sub_port_block models.AccessSubPortBlock
	description := "access_sub_port_block created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessSubPortBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessSubPortBlockConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessSubPortBlockExists("aci_access_sub_port_block.fooaccess_sub_port_block", &access_sub_port_block),
					testAccCheckAciAccessSubPortBlockAttributes(description, &access_sub_port_block),
				),
			},
		},
	})
}

func TestAccAciAccessSubPortBlock_update(t *testing.T) {
	var access_sub_port_block models.AccessSubPortBlock
	description := "access_sub_port_block created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessSubPortBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessSubPortBlockConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessSubPortBlockExists("aci_access_sub_port_block.fooaccess_sub_port_block", &access_sub_port_block),
					testAccCheckAciAccessSubPortBlockAttributes(description, &access_sub_port_block),
				),
			},
			{
				Config: testAccCheckAciAccessSubPortBlockConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessSubPortBlockExists("aci_access_sub_port_block.fooaccess_sub_port_block", &access_sub_port_block),
					testAccCheckAciAccessSubPortBlockAttributes(description, &access_sub_port_block),
				),
			},
		},
	})
}

func testAccCheckAciAccessSubPortBlockConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_leaf_interface_profile" "example" {
		description = "s"
		name        = "example"
		annotation  = "tag_leaf"
		name_alias  = "s"
	}
	resource "aci_access_port_selector" "example" {
		leaf_interface_profile_dn = aci_leaf_interface_profile.example.id
		description               = "s"
		name                      = "example"
		access_port_selector_type = "ALL"
		annotation                = "tag_port_selector"
		name_alias                = "alias_port_selector"
	} 

	resource "aci_access_sub_port_block" "fooaccess_sub_port_block" {
		  access_port_selector_dn  = aci_access_port_selector.example.id
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  from_card  = "1"
		  from_port  = "1"
		  from_sub_port  = "1"
		  name_alias  = "example"
		  to_card  = "1"
		  to_port  = "1"
		  to_sub_port  = "1"
		}
	`, description)
}

func testAccCheckAciAccessSubPortBlockExists(name string, access_sub_port_block *models.AccessSubPortBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Access Sub Port Block %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access Sub Port Block dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		access_sub_port_blockFound := models.AccessSubPortBlockFromContainer(cont)
		if access_sub_port_blockFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Access Sub Port Block %s not found", rs.Primary.ID)
		}
		*access_sub_port_block = *access_sub_port_blockFound
		return nil
	}
}

func testAccCheckAciAccessSubPortBlockDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_access_sub_port_block" {
			cont, err := client.Get(rs.Primary.ID)
			access_sub_port_block := models.AccessSubPortBlockFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Access Sub Port Block %s Still exists", access_sub_port_block.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciAccessSubPortBlockAttributes(description string, access_sub_port_block *models.AccessSubPortBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != access_sub_port_block.Description {
			return fmt.Errorf("Bad access_sub_port_block Description %s", access_sub_port_block.Description)
		}

		if "example" != access_sub_port_block.Name {
			return fmt.Errorf("Bad access_sub_port_block name %s", access_sub_port_block.Name)
		}

		if "example" != access_sub_port_block.Annotation {
			return fmt.Errorf("Bad access_sub_port_block annotation %s", access_sub_port_block.Annotation)
		}

		if "1" != access_sub_port_block.FromCard {
			return fmt.Errorf("Bad access_sub_port_block from_card %s", access_sub_port_block.FromCard)
		}

		if "1" != access_sub_port_block.FromPort {
			return fmt.Errorf("Bad access_sub_port_block from_port %s", access_sub_port_block.FromPort)
		}

		if "1" != access_sub_port_block.FromSubPort {
			return fmt.Errorf("Bad access_sub_port_block from_sub_port %s", access_sub_port_block.FromSubPort)
		}

		if "example" != access_sub_port_block.NameAlias {
			return fmt.Errorf("Bad access_sub_port_block name_alias %s", access_sub_port_block.NameAlias)
		}

		if "1" != access_sub_port_block.ToCard {
			return fmt.Errorf("Bad access_sub_port_block to_card %s", access_sub_port_block.ToCard)
		}

		if "1" != access_sub_port_block.ToPort {
			return fmt.Errorf("Bad access_sub_port_block to_port %s", access_sub_port_block.ToPort)
		}

		if "1" != access_sub_port_block.ToSubPort {
			return fmt.Errorf("Bad access_sub_port_block to_sub_port %s", access_sub_port_block.ToSubPort)
		}

		return nil
	}
}
