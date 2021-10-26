package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciAccessAccessGroup_Basic(t *testing.T) {
	var access_access_group models.AccessAccessGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessAccessGroupConfig_basic("101"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessAccessGroupExists("aci_access_group.fooaccess_access_group", &access_access_group),
					testAccCheckAciAccessAccessGroupAttributes("101", &access_access_group),
				),
			},
		},
	})
}

func TestAccAciAccessAccessGroup_update(t *testing.T) {
	var access_access_group models.AccessAccessGroup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessAccessGroupConfig_basic("101"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessAccessGroupExists("aci_access_group.fooaccess_access_group", &access_access_group),
					testAccCheckAciAccessAccessGroupAttributes("101", &access_access_group),
				),
			},
			{
				Config: testAccCheckAciAccessAccessGroupConfig_basic("102"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessAccessGroupExists("aci_access_group.fooaccess_access_group", &access_access_group),
					testAccCheckAciAccessAccessGroupAttributes("102", &access_access_group),
				),
			},
		},
	})
}

func testAccCheckAciAccessAccessGroupConfig_basic(fexID string) string {
	return fmt.Sprintf(`
	resource "aci_leaf_interface_profile" "example" {
		name  = "foo_leaf_int_prof"
	}

	resource "aci_access_port_selector" "example" {
		leaf_interface_profile_dn = aci_leaf_interface_profile.example.id
		description               = "from terraform"
		name                      = "demo_port_selector"
		access_port_selector_type  = "ALL"
	}
	resource "aci_fex_bundle_group" "example" {
		fex_profile_dn  = aci_fex_profile.example.id
		name            = "example"
	}
	resource "aci_fex_profile" "example" {
		name        = "fex_prof"
	}
	resource "aci_access_group" "fooaccess_access_group" {
		access_port_selector_dn  = aci_access_port_selector.example.id
		annotation = "check"
		fex_id  = "%s"
		tdn  = aci_fex_bundle_group.example.id
	}
	`, fexID)
}

func testAccCheckAciAccessAccessGroupExists(name string, access_access_group *models.AccessAccessGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Access Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		access_access_groupFound := models.AccessAccessGroupFromContainer(cont)
		if access_access_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Access Group %s not found", rs.Primary.ID)
		}
		*access_access_group = *access_access_groupFound
		return nil
	}
}

func testAccCheckAciAccessAccessGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_access_group" {
			cont, err := client.Get(rs.Primary.ID)
			access_access_group := models.AccessAccessGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Access Group %s Still exists", access_access_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciAccessAccessGroupAttributes(fexID string, access_access_group *models.AccessAccessGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if "check" != access_access_group.Annotation {
			return fmt.Errorf("Bad aci_access_group annotation %s", access_access_group.Annotation)
		}

		if fexID != access_access_group.FexId {
			return fmt.Errorf("Bad aci_access_group fex_id %s", access_access_group.FexId)
		}

		return nil
	}
}
