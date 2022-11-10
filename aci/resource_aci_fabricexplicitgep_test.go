package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciVPCExplicitProtectionGroup_Basic(t *testing.T) {
	var vpc_explicit_protection_group models.VPCExplicitProtectionGroup
	description := "vpc_explicit_protection_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVPCExplicitProtectionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVPCExplicitProtectionGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists("aci_vpc_explicit_protection_group.foovpc_explicit_protection_group", &vpc_explicit_protection_group),
					testAccCheckAciVPCExplicitProtectionGroupAttributes(description, &vpc_explicit_protection_group),
				),
			},
		},
	})
}

func TestAccAciVPCExplicitProtectionGroup_update(t *testing.T) {
	var vpc_explicit_protection_group models.VPCExplicitProtectionGroup
	description := "vpc_explicit_protection_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVPCExplicitProtectionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVPCExplicitProtectionGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists("aci_vpc_explicit_protection_group.foovpc_explicit_protection_group", &vpc_explicit_protection_group),
					testAccCheckAciVPCExplicitProtectionGroupAttributes(description, &vpc_explicit_protection_group),
				),
			},
			{
				Config: testAccCheckAciVPCExplicitProtectionGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCExplicitProtectionGroupExists("aci_vpc_explicit_protection_group.foovpc_explicit_protection_group", &vpc_explicit_protection_group),
					testAccCheckAciVPCExplicitProtectionGroupAttributes(description, &vpc_explicit_protection_group),
				),
			},
		},
	})
}

func testAccCheckAciVPCExplicitProtectionGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_vpc_explicit_protection_group" "foovpc_explicit_protection_group" {
		description = "%s"
		name  = "example2"
		switch1 = "203"
		switch2 = "204"
		annotation  = "example"
		vpc_explicit_protection_group_id  = "5"
		}
	`, description)
}

func testAccCheckAciVPCExplicitProtectionGroupExists(name string, vpc_explicit_protection_group *models.VPCExplicitProtectionGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VPC Explicit Protection Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC Explicit Protection Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vpc_explicit_protection_groupFound := models.VPCExplicitProtectionGroupFromContainer(cont)
		if vpc_explicit_protection_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VPC Explicit Protection Group %s not found", rs.Primary.ID)
		}
		*vpc_explicit_protection_group = *vpc_explicit_protection_groupFound
		return nil
	}
}

func testAccCheckAciVPCExplicitProtectionGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_vpc_explicit_protection_group" {
			cont, err := client.Get(rs.Primary.ID)
			vpc_explicit_protection_group := models.VPCExplicitProtectionGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VPC Explicit Protection Group %s Still exists", vpc_explicit_protection_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciVPCExplicitProtectionGroupAttributes(description string, vpc_explicit_protection_group *models.VPCExplicitProtectionGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if "example2" != vpc_explicit_protection_group.Name {
			return fmt.Errorf("Bad vpc_explicit_protection_group name %s", vpc_explicit_protection_group.Name)
		}

		if "example" != vpc_explicit_protection_group.Annotation {
			return fmt.Errorf("Bad vpc_explicit_protection_group annotation %s", vpc_explicit_protection_group.Annotation)
		}

		if "5" != vpc_explicit_protection_group.VPCExplicitProtectionGroup_id {
			return fmt.Errorf("Bad vpc_explicit_protection_group vpc_explicit_protection_group_id %s", vpc_explicit_protection_group.VPCExplicitProtectionGroup_id)
		}

		return nil
	}
}
