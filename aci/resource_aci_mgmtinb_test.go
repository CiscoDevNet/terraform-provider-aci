package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciInBandManagementEPg_Basic(t *testing.T) {
	var node_inb_mgmt_epg models.InBandManagementEPg
	description := "node_inb_mgmt_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInBandManagementEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInBandManagementEPgConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInBandManagementEPgExists("aci_node_inb_mgmt_epg.foonode_inb_mgmt_epg", &node_inb_mgmt_epg),
					testAccCheckAciInBandManagementEPgAttributes(description, &node_inb_mgmt_epg),
				),
			},
		},
	})
}

func TestAccAciInBandManagementEPg_update(t *testing.T) {
	var node_inb_mgmt_epg models.InBandManagementEPg
	description := "node_inb_mgmt_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInBandManagementEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInBandManagementEPgConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInBandManagementEPgExists("aci_node_inb_mgmt_epg.foonode_inb_mgmt_epg", &node_inb_mgmt_epg),
					testAccCheckAciInBandManagementEPgAttributes(description, &node_inb_mgmt_epg),
				),
			},
			{
				Config: testAccCheckAciInBandManagementEPgConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInBandManagementEPgExists("aci_node_inb_mgmt_epg.foonode_inb_mgmt_epg", &node_inb_mgmt_epg),
					testAccCheckAciInBandManagementEPgAttributes(description, &node_inb_mgmt_epg),
				),
			},
		},
	})
}

func testAccCheckAciInBandManagementEPgConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_node_inb_mgmt_epg" "foonode_inb_mgmt_epg" {
		management_profile_dn  = "${aci_management_profile.example.id}"
		description = "%s"
		name  = "example"
  		annotation  = "example"
  		encap  = "vlan-1"
  		exception_tag  = "example"
  		flood_on_encap = "disabled"
  		match_t = "All"
  		name_alias  = "example"
  		pref_gr_memb = "exclude"
  		prio = "level1"
	}
	`, description)
}

func testAccCheckAciInBandManagementEPgExists(name string, node_inb_mgmt_epg *models.InBandManagementEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("In-Band Management EPg %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No In-Band Management EPg dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		node_inb_mgmt_epgFound := models.InBandManagementEPgFromContainer(cont)
		if node_inb_mgmt_epgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("In-Band Management EPg %s not found", rs.Primary.ID)
		}
		*node_inb_mgmt_epg = *node_inb_mgmt_epgFound
		return nil
	}
}

func testAccCheckAciInBandManagementEPgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_node_inb_mgmt_epg" {
			cont, err := client.Get(rs.Primary.ID)
			node_inb_mgmt_epg := models.InBandManagementEPgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("In-Band Management EPg %s Still exists", node_inb_mgmt_epg.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciInBandManagementEPgAttributes(description string, node_inb_mgmt_epg *models.InBandManagementEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != node_inb_mgmt_epg.Description {
			return fmt.Errorf("Bad node_inb_mgmt_epg Description %s", node_inb_mgmt_epg.Description)
		}

		if "example" != node_inb_mgmt_epg.Name {
			return fmt.Errorf("Bad node_inb_mgmt_epg name %s", node_inb_mgmt_epg.Name)
		}

		if "example" != node_inb_mgmt_epg.Annotation {
			return fmt.Errorf("Bad node_inb_mgmt_epg annotation %s", node_inb_mgmt_epg.Annotation)
		}

		if "example" != node_inb_mgmt_epg.Encap {
			return fmt.Errorf("Bad node_inb_mgmt_epg encap %s", node_inb_mgmt_epg.Encap)
		}

		if "example" != node_inb_mgmt_epg.ExceptionTag {
			return fmt.Errorf("Bad node_inb_mgmt_epg exception_tag %s", node_inb_mgmt_epg.ExceptionTag)
		}

		if "disabled" != node_inb_mgmt_epg.FloodOnEncap {
			return fmt.Errorf("Bad node_inb_mgmt_epg flood_on_encap %s", node_inb_mgmt_epg.FloodOnEncap)
		}

		if "All" != node_inb_mgmt_epg.MatchT {
			return fmt.Errorf("Bad node_inb_mgmt_epg match_t %s", node_inb_mgmt_epg.MatchT)
		}

		if "example" != node_inb_mgmt_epg.NameAlias {
			return fmt.Errorf("Bad node_inb_mgmt_epg name_alias %s", node_inb_mgmt_epg.NameAlias)
		}

		if "exclude" != node_inb_mgmt_epg.PrefGrMemb {
			return fmt.Errorf("Bad node_inb_mgmt_epg pref_gr_memb %s", node_inb_mgmt_epg.PrefGrMemb)
		}

		if "level1" != node_inb_mgmt_epg.Prio {
			return fmt.Errorf("Bad node_inb_mgmt_epg prio %s", node_inb_mgmt_epg.Prio)
		}

		return nil
	}
}
