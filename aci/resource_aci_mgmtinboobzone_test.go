package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciInBManagedNodesZone_Basic(t *testing.T) {
	var in_b_managed_nodes_zone models.InBManagedNodesZone
	mgmt_grp_name := acctest.RandString(5)
	mgmt_in_b_zone_name := acctest.RandString(5)
	description := "in_b_managed_nodes_zone created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInBManagedNodesZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInBManagedNodesZoneConfig_basic(mgmt_grp_name, mgmt_in_b_zone_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInBManagedNodesZoneExists("aci_mgmt_zone.fooin_b_managed_nodes_zone", &in_b_managed_nodes_zone),
					testAccCheckAciInBManagedNodesZoneAttributes(mgmt_grp_name, mgmt_in_b_zone_name, description, &in_b_managed_nodes_zone),
				),
			},
		},
	})
}

func TestAccAciInBManagedNodesZone_Update(t *testing.T) {
	var in_b_managed_nodes_zone models.InBManagedNodesZone
	mgmt_grp_name := acctest.RandString(5)
	mgmt_in_b_zone_name := acctest.RandString(5)
	description := "in_b_managed_nodes_zone created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInBManagedNodesZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInBManagedNodesZoneConfig_basic(mgmt_grp_name, mgmt_in_b_zone_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInBManagedNodesZoneExists("aci_mgmt_zone.fooin_b_managed_nodes_zone", &in_b_managed_nodes_zone),
					testAccCheckAciInBManagedNodesZoneAttributes(mgmt_grp_name, mgmt_in_b_zone_name, description, &in_b_managed_nodes_zone),
				),
			},
			{
				Config: testAccCheckAciInBManagedNodesZoneConfig_basic(mgmt_grp_name, mgmt_in_b_zone_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInBManagedNodesZoneExists("aci_mgmt_zone.fooin_b_managed_nodes_zone", &in_b_managed_nodes_zone),
					testAccCheckAciInBManagedNodesZoneAttributes(mgmt_grp_name, mgmt_in_b_zone_name, description, &in_b_managed_nodes_zone),
				),
			},
		},
	})
}

func testAccCheckAciInBManagedNodesZoneConfig_basic(mgmt_grp_name, mgmt_in_b_zone_name string) string {
	return fmt.Sprintf(`

	resource "aci_managed_node_connectivity_group" "foomanaged_node_connectivity_group" {
		name 		= "%s"

	}

	resource "aci_mgmt_zone" "fooin_b_managed_nodes_zone" {
		name 		= "%s"
		type = "in_band"
		description = "in_b_managed_nodes_zone created while acceptance testing"
		managed_node_connectivity_group_dn = aci_managed_node_connectivity_group.foomanaged_node_connectivity_group.id
	}

	`, mgmt_grp_name, mgmt_in_b_zone_name)
}

func testAccCheckAciInBManagedNodesZoneExists(name string, in_b_managed_nodes_zone *models.InBManagedNodesZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("InB Managed Nodes Zone %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No InB Managed Nodes Zone dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		in_b_managed_nodes_zoneFound := models.InBManagedNodesZoneFromContainer(cont)
		if in_b_managed_nodes_zoneFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("InB Managed Nodes Zone %s not found", rs.Primary.ID)
		}
		*in_b_managed_nodes_zone = *in_b_managed_nodes_zoneFound
		return nil
	}
}

func testAccCheckAciInBManagedNodesZoneDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_in_b_managed_nodes_zone" {
			cont, err := client.Get(rs.Primary.ID)
			in_b_managed_nodes_zone := models.InBManagedNodesZoneFromContainer(cont)
			if err == nil {
				return fmt.Errorf("InB Managed Nodes Zone %s Still exists", in_b_managed_nodes_zone.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciInBManagedNodesZoneAttributes(mgmt_grp_name, mgmt_in_b_zone_name, description string, in_b_managed_nodes_zone *models.InBManagedNodesZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// if mgmt_in_b_zone_name != GetMOName(in_b_managed_nodes_zone.DistinguishedName) {
		// 	return fmt.Errorf("Bad mgmt_in_b_zone %s", GetMOName(in_b_managed_nodes_zone.DistinguishedName))
		// }
		// if mgmt_grp_name != GetMOName(GetParentDn(in_b_managed_nodes_zone.DistinguishedName)) {
		// 	return fmt.Errorf(" Bad mgmt_grp %s", GetMOName(GetParentDn(in_b_managed_nodes_zone.DistinguishedName)))
		// }
		if description != in_b_managed_nodes_zone.Description {
			return fmt.Errorf("Bad in_b_managed_nodes_zone Description %s", in_b_managed_nodes_zone.Description)
		}
		return nil
	}
}
