package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciCloudEPg_Basic(t *testing.T) {
	var cloud_e_pg models.CloudEPg
	description := "cloud_e_pg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudEPgConfig_basic(description, "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists("aci_cloud_e_pg.foocloud_e_pg", &cloud_e_pg),
					testAccCheckAciCloudEPgAttributes(description, "All", &cloud_e_pg),
				),
			},
			{
				ResourceName:      "aci_cloud_e_pg",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciCloudEPg_update(t *testing.T) {
	var cloud_e_pg models.CloudEPg
	description := "cloud_e_pg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudEPgConfig_basic(description, "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists("aci_cloud_e_pg.foocloud_e_pg", &cloud_e_pg),
					testAccCheckAciCloudEPgAttributes(description, "All", &cloud_e_pg),
				),
			},
			{
				Config: testAccCheckAciCloudEPgConfig_basic(description, "AtleastOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists("aci_cloud_e_pg.foocloud_e_pg", &cloud_e_pg),
					testAccCheckAciCloudEPgAttributes(description, "AtleastOne", &cloud_e_pg),
				),
			},
		},
	})
}

func testAccCheckAciCloudEPgConfig_basic(description, match_t string) string {
	return fmt.Sprintf(`
	
	resource "aci_tenant" "footenant" {
		description = "Tenant created while acceptance testing"
		name        = "demo_tenant"
	}

	resource "aci_cloud_applicationcontainer" "foocloud_applicationcontainer" {
		tenant_dn   = "${aci_tenant.footenant.id}"
		name        = "demo_app"
		annotation  = "tag_app"
	}

	resource "aci_cloud_e_pg" "foocloud_e_pg" {
		cloud_applicationcontainer_dn = "${aci_cloud_applicationcontainer.foocloud_applicationcontainer.id}"
		description                   = "%s"
		name                          = "cloud_epg"
		annotation                    = "tag_epg"
		exception_tag                 = "0"
		flood_on_encap                = "disabled"
		match_t                       = "%s"
		name_alias                    = "alias_epg"
		pref_gr_memb                  = "exclude"
		prio                          = "unspecified"
	}
	  
	`, description, match_t)
}

func testAccCheckAciCloudEPgExists(name string, cloud_e_pg *models.CloudEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud EPg %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud EPg dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_e_pgFound := models.CloudEPgFromContainer(cont)
		if cloud_e_pgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud EPg %s not found", rs.Primary.ID)
		}
		*cloud_e_pg = *cloud_e_pgFound
		return nil
	}
}

func testAccCheckAciCloudEPgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_e_pg" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_e_pg := models.CloudEPgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud EPg %s Still exists", cloud_e_pg.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudEPgAttributes(description, match_t string, cloud_e_pg *models.CloudEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_e_pg.Description {
			return fmt.Errorf("Bad cloud_e_pg Description %s", cloud_e_pg.Description)
		}

		if "cloud_epg" != cloud_e_pg.Name {
			return fmt.Errorf("Bad cloud_e_pg name %s", cloud_e_pg.Name)
		}

		if "tag_epg" != cloud_e_pg.Annotation {
			return fmt.Errorf("Bad cloud_e_pg annotation %s", cloud_e_pg.Annotation)
		}

		if "0" != cloud_e_pg.ExceptionTag {
			return fmt.Errorf("Bad cloud_e_pg exception_tag %s", cloud_e_pg.ExceptionTag)
		}

		if "disabled" != cloud_e_pg.FloodOnEncap {
			return fmt.Errorf("Bad cloud_e_pg flood_on_encap %s", cloud_e_pg.FloodOnEncap)
		}

		if match_t != cloud_e_pg.MatchT {
			return fmt.Errorf("Bad cloud_e_pg match_t %s", cloud_e_pg.MatchT)
		}

		if "alias_epg" != cloud_e_pg.NameAlias {
			return fmt.Errorf("Bad cloud_e_pg name_alias %s", cloud_e_pg.NameAlias)
		}

		if "exclude" != cloud_e_pg.PrefGrMemb {
			return fmt.Errorf("Bad cloud_e_pg pref_gr_memb %s", cloud_e_pg.PrefGrMemb)
		}

		if "unspecified" != cloud_e_pg.Prio {
			return fmt.Errorf("Bad cloud_e_pg prio %s", cloud_e_pg.Prio)
		}

		return nil
	}
}
