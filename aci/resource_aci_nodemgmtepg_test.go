package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciInBandManagementEPg_Basic(t *testing.T) {
	var in_band_management_e_pg models.InBandManagementEPg
	description := "node_mgmt_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInBandManagementEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInBandManagementEPgConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInBandManagementEPgExists("aci_node_mgmt_epg.foonode_mgmt_epg", &in_band_management_e_pg),
					testAccCheckAciInBandManagementEPgAttributes(description, &in_band_management_e_pg),
				),
			},
		},
	})
}

func TestAccAciInBandManagementEPg_update(t *testing.T) {
	var in_band_management_e_pg models.InBandManagementEPg
	description := "node_mgmt_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInBandManagementEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInBandManagementEPgConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInBandManagementEPgExists("aci_node_mgmt_epg.foonode_mgmt_epg", &in_band_management_e_pg),
					testAccCheckAciInBandManagementEPgAttributes(description, &in_band_management_e_pg),
				),
			},
			{
				Config: testAccCheckAciInBandManagementEPgConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInBandManagementEPgExists("aci_node_mgmt_epg.foonode_mgmt_epg", &in_band_management_e_pg),
					testAccCheckAciInBandManagementEPgAttributes(description, &in_band_management_e_pg),
				),
			},
		},
	})
}

func testAccCheckAciInBandManagementEPgConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_node_mgmt_epg" "foonode_mgmt_epg" {
		type = "in_band"
		management_profile_dn  = "uni/tn-mgmt/mgmtp-default"
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

func testAccCheckAciInBandManagementEPgExists(name string, in_band_management_e_pg *models.InBandManagementEPg) resource.TestCheckFunc {
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

		in_band_management_e_pgFound := models.InBandManagementEPgFromContainer(cont)
		if in_band_management_e_pgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("In-Band Management EPg %s not found", rs.Primary.ID)
		}
		*in_band_management_e_pg = *in_band_management_e_pgFound
		return nil
	}
}

func testAccCheckAciInBandManagementEPgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_node_mgmt_epg" {
			cont, err := client.Get(rs.Primary.ID)
			in_band_management_e_pg := models.InBandManagementEPgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("In-Band Management EPg %s Still exists", in_band_management_e_pg.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciInBandManagementEPgAttributes(description string, in_band_management_e_pg *models.InBandManagementEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != in_band_management_e_pg.Description {
			return fmt.Errorf("Bad in_band_management_e_pg Description %s", in_band_management_e_pg.Description)
		}

		if "example" != in_band_management_e_pg.Name {
			return fmt.Errorf("Bad in_band_management_e_pg name %s", in_band_management_e_pg.Name)
		}

		if "example" != in_band_management_e_pg.Annotation {
			return fmt.Errorf("Bad in_band_management_e_pg annotation %s", in_band_management_e_pg.Annotation)
		}

		if "vlan-1" != in_band_management_e_pg.Encap {
			return fmt.Errorf("Bad in_band_management_e_pg encap %s", in_band_management_e_pg.Encap)
		}

		if "example" != in_band_management_e_pg.ExceptionTag {
			return fmt.Errorf("Bad in_band_management_e_pg exception_tag %s", in_band_management_e_pg.ExceptionTag)
		}

		if "disabled" != in_band_management_e_pg.FloodOnEncap {
			return fmt.Errorf("Bad in_band_management_e_pg flood_on_encap %s", in_band_management_e_pg.FloodOnEncap)
		}

		if "All" != in_band_management_e_pg.MatchT {
			return fmt.Errorf("Bad in_band_management_e_pg match_t %s", in_band_management_e_pg.MatchT)
		}

		if "example" != in_band_management_e_pg.NameAlias {
			return fmt.Errorf("Bad in_band_management_e_pg name_alias %s", in_band_management_e_pg.NameAlias)
		}

		if "exclude" != in_band_management_e_pg.PrefGrMemb {
			return fmt.Errorf("Bad in_band_management_e_pg pref_gr_memb %s", in_band_management_e_pg.PrefGrMemb)
		}

		if "level1" != in_band_management_e_pg.Prio {
			return fmt.Errorf("Bad in_band_management_e_pg prio %s", in_band_management_e_pg.Prio)
		}

		return nil
	}
}

func TestAccAciOutOfBandManagementEPg_Basic(t *testing.T) {
	var out_of_band_management_e_pg models.OutOfBandManagementEPg
	description := "node_mgmt_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOutOfBandManagementEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOutOfBandManagementEPgConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutOfBandManagementEPgExists("aci_node_mgmt_epg.foonode_mgmt_epg", &out_of_band_management_e_pg),
					testAccCheckAciOutOfBandManagementEPgAttributes(description, &out_of_band_management_e_pg),
				),
			},
		},
	})
}

func TestAccAciOutOfBandManagementEPg_update(t *testing.T) {
	var out_of_band_management_e_pg models.OutOfBandManagementEPg
	description := "node_mgmt_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOutOfBandManagementEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOutOfBandManagementEPgConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutOfBandManagementEPgExists("aci_node_mgmt_epg.foonode_mgmt_epg", &out_of_band_management_e_pg),
					testAccCheckAciOutOfBandManagementEPgAttributes(description, &out_of_band_management_e_pg),
				),
			},
			{
				Config: testAccCheckAciOutOfBandManagementEPgConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutOfBandManagementEPgExists("aci_node_mgmt_epg.foonode_mgmt_epg", &out_of_band_management_e_pg),
					testAccCheckAciOutOfBandManagementEPgAttributes(description, &out_of_band_management_e_pg),
				),
			},
		},
	})
}

func testAccCheckAciOutOfBandManagementEPgConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_node_mgmt_epg" "foonode_mgmt_epg" {
		type = "out_of_band"
		management_profile_dn  = "uni/tn-mgmt/mgmtp-default"
		description = "%s"
		name  = "example"
  		annotation  = "example"
  		name_alias  = "example"
  		prio = "level1"
	}
	`, description)
}

func testAccCheckAciOutOfBandManagementEPgExists(name string, out_of_band_management_e_pg *models.OutOfBandManagementEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Out-Of-Band Management EPg %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Out-Of-Band Management EPg dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		out_of_band_management_e_pgFound := models.OutOfBandManagementEPgFromContainer(cont)
		if out_of_band_management_e_pgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Out-Of-Band Management EPg %s not found", rs.Primary.ID)
		}
		*out_of_band_management_e_pg = *out_of_band_management_e_pgFound
		return nil
	}
}

func testAccCheckAciOutOfBandManagementEPgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_node_mgmt_epg" {
			cont, err := client.Get(rs.Primary.ID)
			out_of_band_management_e_pg := models.OutOfBandManagementEPgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Out-Of-Band Management EPg %s Still exists", out_of_band_management_e_pg.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciOutOfBandManagementEPgAttributes(description string, out_of_band_management_e_pg *models.OutOfBandManagementEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != out_of_band_management_e_pg.Description {
			return fmt.Errorf("Bad out_of_band_management_e_pg Description %s", out_of_band_management_e_pg.Description)
		}

		if "example" != out_of_band_management_e_pg.Name {
			return fmt.Errorf("Bad out_of_band_management_e_pg name %s", out_of_band_management_e_pg.Name)
		}

		if "example" != out_of_band_management_e_pg.Annotation {
			return fmt.Errorf("Bad out_of_band_management_e_pg annotation %s", out_of_band_management_e_pg.Annotation)
		}

		if "example" != out_of_band_management_e_pg.NameAlias {
			return fmt.Errorf("Bad out_of_band_management_e_pg name_alias %s", out_of_band_management_e_pg.NameAlias)
		}

		if "level1" != out_of_band_management_e_pg.Prio {
			return fmt.Errorf("Bad out_of_band_management_e_pg prio %s", out_of_band_management_e_pg.Prio)
		}

		return nil
	}
}
