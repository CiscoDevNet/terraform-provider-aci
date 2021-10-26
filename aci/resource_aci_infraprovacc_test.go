package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciVlanEncapsulationforVxlanTraffic_Basic(t *testing.T) {
	var vlan_encapsulationfor_vxlan_traffic models.VlanEncapsulationforVxlanTraffic
	description := "vlan_encapsulationfor_vxlan_traffic created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVlanEncapsulationforVxlanTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVlanEncapsulationforVxlanTrafficConfig_basic(description, "alias_heavy_traffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVlanEncapsulationforVxlanTrafficExists("aci_vlan_encapsulationfor_vxlan_traffic.foovlan_encapsulationfor_vxlan_traffic", &vlan_encapsulationfor_vxlan_traffic),
					testAccCheckAciVlanEncapsulationforVxlanTrafficAttributes(description, "alias_heavy_traffic", &vlan_encapsulationfor_vxlan_traffic),
				),
			},
		},
	})
}

func TestAccAciVlanEncapsulationforVxlanTraffic_update(t *testing.T) {
	var vlan_encapsulationfor_vxlan_traffic models.VlanEncapsulationforVxlanTraffic
	description := "vlan_encapsulationfor_vxlan_traffic created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVlanEncapsulationforVxlanTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVlanEncapsulationforVxlanTrafficConfig_basic(description, "alias_heavy_traffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVlanEncapsulationforVxlanTrafficExists("aci_vlan_encapsulationfor_vxlan_traffic.foovlan_encapsulationfor_vxlan_traffic", &vlan_encapsulationfor_vxlan_traffic),
					testAccCheckAciVlanEncapsulationforVxlanTrafficAttributes(description, "alias_heavy_traffic", &vlan_encapsulationfor_vxlan_traffic),
				),
			},
			{
				Config: testAccCheckAciVlanEncapsulationforVxlanTrafficConfig_basic(description, "alias_low_traffic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVlanEncapsulationforVxlanTrafficExists("aci_vlan_encapsulationfor_vxlan_traffic.foovlan_encapsulationfor_vxlan_traffic", &vlan_encapsulationfor_vxlan_traffic),
					testAccCheckAciVlanEncapsulationforVxlanTrafficAttributes(description, "alias_low_traffic", &vlan_encapsulationfor_vxlan_traffic),
				),
			},
		},
	})
}

func testAccCheckAciVlanEncapsulationforVxlanTrafficConfig_basic(description, name_alias string) string {
	return fmt.Sprintf(`

	resource "aci_attachable_access_entity_profile" "example" {
		description = "AAEP description"
		name        = "demo_entity_prof"
		annotation  = "tag_entity"
		name_alias  = "alias_entity"
	}
	
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "foovlan_encapsulationfor_vxlan_traffic" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.example.id
		description                         = "%s"
		annotation                          = "tag_traffic"
		name_alias                          = "%s"
	}
	`, description, name_alias)
}

func testAccCheckAciVlanEncapsulationforVxlanTrafficExists(name string, vlan_encapsulationfor_vxlan_traffic *models.VlanEncapsulationforVxlanTraffic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Vlan Encapsulation for Vxlan Traffic %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vlan Encapsulation for Vxlan Traffic dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vlan_encapsulationfor_vxlan_trafficFound := models.VlanEncapsulationforVxlanTrafficFromContainer(cont)
		if vlan_encapsulationfor_vxlan_trafficFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Vlan Encapsulation for Vxlan Traffic %s not found", rs.Primary.ID)
		}
		*vlan_encapsulationfor_vxlan_traffic = *vlan_encapsulationfor_vxlan_trafficFound
		return nil
	}
}

func testAccCheckAciVlanEncapsulationforVxlanTrafficDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_vlan_encapsulationfor_vxlan_traffic" {
			cont, err := client.Get(rs.Primary.ID)
			vlan_encapsulationfor_vxlan_traffic := models.VlanEncapsulationforVxlanTrafficFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Vlan Encapsulation for Vxlan Traffic %s Still exists", vlan_encapsulationfor_vxlan_traffic.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciVlanEncapsulationforVxlanTrafficAttributes(description, name_alias string, vlan_encapsulationfor_vxlan_traffic *models.VlanEncapsulationforVxlanTraffic) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != vlan_encapsulationfor_vxlan_traffic.Description {
			return fmt.Errorf("Bad vlan_encapsulationfor_vxlan_traffic Description %s", vlan_encapsulationfor_vxlan_traffic.Description)
		}

		if "tag_traffic" != vlan_encapsulationfor_vxlan_traffic.Annotation {
			return fmt.Errorf("Bad vlan_encapsulationfor_vxlan_traffic annotation %s", vlan_encapsulationfor_vxlan_traffic.Annotation)
		}

		if name_alias != vlan_encapsulationfor_vxlan_traffic.NameAlias {
			return fmt.Errorf("Bad vlan_encapsulationfor_vxlan_traffic name_alias %s", vlan_encapsulationfor_vxlan_traffic.NameAlias)
		}

		return nil
	}
}
