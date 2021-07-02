package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciFirmwareGroup_Basic(t *testing.T) {
	var firmware_group models.FirmwareGroup
	description := "firmware_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFirmwareGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFirmwareGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareGroupExists("aci_firmware_group.foofirmware_group", &firmware_group),
					testAccCheckAciFirmwareGroupAttributes(description, &firmware_group),
				),
			},
		},
	})
}

func TestAccAciFirmwareGroup_update(t *testing.T) {
	var firmware_group models.FirmwareGroup
	description := "firmware_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFirmwareGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFirmwareGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareGroupExists("aci_firmware_group.foofirmware_group", &firmware_group),
					testAccCheckAciFirmwareGroupAttributes(description, &firmware_group),
				),
			},
			{
				Config: testAccCheckAciFirmwareGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareGroupExists("aci_firmware_group.foofirmware_group", &firmware_group),
					testAccCheckAciFirmwareGroupAttributes(description, &firmware_group),
				),
			},
		},
	})
}

func testAccCheckAciFirmwareGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_firmware_group" "foofirmware_group" {
		description = "%s"
		
		  name  = "example"
		  annotation  = "example"
		  name_alias  = "example"
		  firmware_group_type  = "ALL"
		}
	`, description)
}

func testAccCheckAciFirmwareGroupExists(name string, firmware_group *models.FirmwareGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Firmware Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Firmware Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		firmware_groupFound := models.FirmwareGroupFromContainer(cont)
		if firmware_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Firmware Group %s not found", rs.Primary.ID)
		}
		*firmware_group = *firmware_groupFound
		return nil
	}
}

func testAccCheckAciFirmwareGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_firmware_group" {
			cont, err := client.Get(rs.Primary.ID)
			firmware_group := models.FirmwareGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Firmware Group %s Still exists", firmware_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciFirmwareGroupAttributes(description string, firmware_group *models.FirmwareGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != firmware_group.Description {
			return fmt.Errorf("Bad firmware_group Description %s", firmware_group.Description)
		}

		if "example" != firmware_group.Name {
			return fmt.Errorf("Bad firmware_group name %s", firmware_group.Name)
		}

		if "example" != firmware_group.Annotation {
			return fmt.Errorf("Bad firmware_group annotation %s", firmware_group.Annotation)
		}

		if "example" != firmware_group.NameAlias {
			return fmt.Errorf("Bad firmware_group name_alias %s", firmware_group.NameAlias)
		}

		if "ALL" != firmware_group.FirmwareGroup_type {
			return fmt.Errorf("Bad firmware_group firmware_group_type %s", firmware_group.FirmwareGroup_type)
		}

		return nil
	}
}
