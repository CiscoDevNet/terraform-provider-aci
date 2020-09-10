package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciPODMaintenanceGroup_Basic(t *testing.T) {
	var pod_maintenance_group models.PODMaintenanceGroup
	description := "pod_maintenance_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPODMaintenanceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPODMaintenanceGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists("aci_pod_maintenance_group.foopod_maintenance_group", &pod_maintenance_group),
					testAccCheckAciPODMaintenanceGroupAttributes(description, &pod_maintenance_group),
				),
			},
			{
				ResourceName:      "aci_pod_maintenance_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciPODMaintenanceGroup_update(t *testing.T) {
	var pod_maintenance_group models.PODMaintenanceGroup
	description := "pod_maintenance_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPODMaintenanceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPODMaintenanceGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists("aci_pod_maintenance_group.foopod_maintenance_group", &pod_maintenance_group),
					testAccCheckAciPODMaintenanceGroupAttributes(description, &pod_maintenance_group),
				),
			},
			{
				Config: testAccCheckAciPODMaintenanceGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPODMaintenanceGroupExists("aci_pod_maintenance_group.foopod_maintenance_group", &pod_maintenance_group),
					testAccCheckAciPODMaintenanceGroupAttributes(description, &pod_maintenance_group),
				),
			},
		},
	})
}

func testAccCheckAciPODMaintenanceGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_pod_maintenance_group" "foopod_maintenance_group" {
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  fwtype  = "switch"
		  name_alias  = "example"
		  pod_maintenance_group_type  = "example"
		}
	`, description)
}

func testAccCheckAciPODMaintenanceGroupExists(name string, pod_maintenance_group *models.PODMaintenanceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("POD Maintenance Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No POD Maintenance Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		pod_maintenance_groupFound := models.PODMaintenanceGroupFromContainer(cont)
		if pod_maintenance_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("POD Maintenance Group %s not found", rs.Primary.ID)
		}
		*pod_maintenance_group = *pod_maintenance_groupFound
		return nil
	}
}

func testAccCheckAciPODMaintenanceGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_pod_maintenance_group" {
			cont, err := client.Get(rs.Primary.ID)
			pod_maintenance_group := models.PODMaintenanceGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("POD Maintenance Group %s Still exists", pod_maintenance_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciPODMaintenanceGroupAttributes(description string, pod_maintenance_group *models.PODMaintenanceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != pod_maintenance_group.Description {
			return fmt.Errorf("Bad pod_maintenance_group Description %s", pod_maintenance_group.Description)
		}

		if "example" != pod_maintenance_group.Name {
			return fmt.Errorf("Bad pod_maintenance_group name %s", pod_maintenance_group.Name)
		}

		if "example" != pod_maintenance_group.Annotation {
			return fmt.Errorf("Bad pod_maintenance_group annotation %s", pod_maintenance_group.Annotation)
		}

		if "switch" != pod_maintenance_group.Fwtype {
			return fmt.Errorf("Bad pod_maintenance_group fwtype %s", pod_maintenance_group.Fwtype)
		}

		if "example" != pod_maintenance_group.NameAlias {
			return fmt.Errorf("Bad pod_maintenance_group name_alias %s", pod_maintenance_group.NameAlias)
		}

		if "example" != pod_maintenance_group.PODMaintenanceGroup_type {
			return fmt.Errorf("Bad pod_maintenance_group pod_maintenance_group_type %s", pod_maintenance_group.PODMaintenanceGroup_type)
		}

		return nil
	}
}
