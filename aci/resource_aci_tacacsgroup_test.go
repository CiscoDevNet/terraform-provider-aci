package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciTACACSMonitoringDestinationGroup_Basic(t *testing.T) {
	var tacacs_monitoring_destination_group models.TACACSMonitoringDestinationGroup
	description := "tacacs_monitoring_destination_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTACACSMonitoringDestinationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTACACSMonitoringDestinationGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSMonitoringDestinationGroupExists("aci_tacacs_accounting.footacacs_monitoring_destination_group", &tacacs_monitoring_destination_group),
					testAccCheckAciTACACSMonitoringDestinationGroupAttributes(description, &tacacs_monitoring_destination_group),
				),
			},
		},
	})
}

func TestAccAciTACACSMonitoringDestinationGroup_Update(t *testing.T) {
	var tacacs_monitoring_destination_group models.TACACSMonitoringDestinationGroup
	description := "tacacs_monitoring_destination_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTACACSMonitoringDestinationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTACACSMonitoringDestinationGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSMonitoringDestinationGroupExists("aci_tacacs_accounting.footacacs_monitoring_destination_group", &tacacs_monitoring_destination_group),
					testAccCheckAciTACACSMonitoringDestinationGroupAttributes(description, &tacacs_monitoring_destination_group),
				),
			},
			{
				Config: testAccCheckAciTACACSMonitoringDestinationGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSMonitoringDestinationGroupExists("aci_tacacs_accounting.footacacs_monitoring_destination_group", &tacacs_monitoring_destination_group),
					testAccCheckAciTACACSMonitoringDestinationGroupAttributes(description, &tacacs_monitoring_destination_group),
				),
			},
		},
	})
}

func testAccCheckAciTACACSMonitoringDestinationGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_tacacs_accounting" "footacacs_monitoring_destination_group" {
		name 		= "test"
		description = "%s"
		annotation  = "example"
		name_alias	= "tacacs_accounting_alias"
	}

	`, description)
}

func testAccCheckAciTACACSMonitoringDestinationGroupExists(name string, tacacs_monitoring_destination_group *models.TACACSMonitoringDestinationGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("TACACS Monitoring Destination Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TACACS Monitoring Destination Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tacacs_monitoring_destination_groupFound := models.TACACSMonitoringDestinationGroupFromContainer(cont)
		if tacacs_monitoring_destination_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("TACACS Monitoring Destination Group %s not found", rs.Primary.ID)
		}
		*tacacs_monitoring_destination_group = *tacacs_monitoring_destination_groupFound
		return nil
	}
}

func testAccCheckAciTACACSMonitoringDestinationGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tacacs_accounting" {
			cont, err := client.Get(rs.Primary.ID)
			tacacs_monitoring_destination_group := models.TACACSMonitoringDestinationGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("TACACS Monitoring Destination Group %s Still exists", tacacs_monitoring_destination_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTACACSMonitoringDestinationGroupAttributes(description string, tacacs_monitoring_destination_group *models.TACACSMonitoringDestinationGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != tacacs_monitoring_destination_group.Name {
			return fmt.Errorf("Bad tacacs_group %s", tacacs_monitoring_destination_group.Name)
		}

		if description != tacacs_monitoring_destination_group.Description {
			return fmt.Errorf("Bad tacacs_monitoring_destination_group Description %s", tacacs_monitoring_destination_group.Description)
		}

		if "example" != tacacs_monitoring_destination_group.Annotation {
			return fmt.Errorf("Bad tacacs_monitoring_destination_group Annotation %s", tacacs_monitoring_destination_group.Annotation)
		}

		if "tacacs_accounting_alias" != tacacs_monitoring_destination_group.NameAlias {
			return fmt.Errorf("Bad tacacs_monitoring_destination_group Name Alias %s", tacacs_monitoring_destination_group.NameAlias)
		}
		return nil
	}
}
