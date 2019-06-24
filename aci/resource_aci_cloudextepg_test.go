package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciCloudExternalEPg_Basic(t *testing.T) {
	var cloud_external_e_pg models.CloudExternalEPg
	description := "cloud_external_e_pg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudExternalEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudExternalEPgConfig_basic(description, "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists("aci_cloud_external_e_pg.foocloud_external_e_pg", &cloud_external_e_pg),
					testAccCheckAciCloudExternalEPgAttributes(description, "All", &cloud_external_e_pg),
				),
			},
			{
				ResourceName:      "aci_cloud_external_e_pg",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciCloudExternalEPg_update(t *testing.T) {
	var cloud_external_e_pg models.CloudExternalEPg
	description := "cloud_external_e_pg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudExternalEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudExternalEPgConfig_basic(description, "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists("aci_cloud_external_e_pg.foocloud_external_e_pg", &cloud_external_e_pg),
					testAccCheckAciCloudExternalEPgAttributes(description, "All", &cloud_external_e_pg),
				),
			},
			{
				Config: testAccCheckAciCloudExternalEPgConfig_basic(description, "AtleastOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists("aci_cloud_external_e_pg.foocloud_external_e_pg", &cloud_external_e_pg),
					testAccCheckAciCloudExternalEPgAttributes(description, "AtleastOne", &cloud_external_e_pg),
				),
			},
		},
	})
}

func testAccCheckAciCloudExternalEPgConfig_basic(description, match_t string) string {
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

	resource "aci_cloud_external_e_pg" "foocloud_external_e_pg" {
		cloud_applicationcontainer_dn = "${aci_cloud_applicationcontainer.foocloud_applicationcontainer.id}"
		description                   = "%s"
		name                          = "cloud_ext_epg"
		annotation                    = "tag_ext_epg"
		exception_tag                 = "0"
		flood_on_encap                = "disabled"
		match_t                       = "%s"
		name_alias                    = "alias_ext"
		pref_gr_memb                  = "exclude"
		prio                          = "unspecified"
		route_reachability            = "inter-site"
	}  
	`, description, match_t)
}

func testAccCheckAciCloudExternalEPgExists(name string, cloud_external_e_pg *models.CloudExternalEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud External EPg %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud External EPg dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_external_e_pgFound := models.CloudExternalEPgFromContainer(cont)
		if cloud_external_e_pgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud External EPg %s not found", rs.Primary.ID)
		}
		*cloud_external_e_pg = *cloud_external_e_pgFound
		return nil
	}
}

func testAccCheckAciCloudExternalEPgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_external_e_pg" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_external_e_pg := models.CloudExternalEPgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud External EPg %s Still exists", cloud_external_e_pg.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudExternalEPgAttributes(description, match_t string, cloud_external_e_pg *models.CloudExternalEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_external_e_pg.Description {
			return fmt.Errorf("Bad cloud_external_e_pg Description %s", cloud_external_e_pg.Description)
		}

		if "cloud_ext_epg" != cloud_external_e_pg.Name {
			return fmt.Errorf("Bad cloud_external_e_pg name %s", cloud_external_e_pg.Name)
		}

		if "tag_ext_epg" != cloud_external_e_pg.Annotation {
			return fmt.Errorf("Bad cloud_external_e_pg annotation %s", cloud_external_e_pg.Annotation)
		}

		if "0" != cloud_external_e_pg.ExceptionTag {
			return fmt.Errorf("Bad cloud_external_e_pg exception_tag %s", cloud_external_e_pg.ExceptionTag)
		}

		if "disabled" != cloud_external_e_pg.FloodOnEncap {
			return fmt.Errorf("Bad cloud_external_e_pg flood_on_encap %s", cloud_external_e_pg.FloodOnEncap)
		}

		if match_t != cloud_external_e_pg.MatchT {
			return fmt.Errorf("Bad cloud_external_e_pg match_t %s", cloud_external_e_pg.MatchT)
		}

		if "alias_ext" != cloud_external_e_pg.NameAlias {
			return fmt.Errorf("Bad cloud_external_e_pg name_alias %s", cloud_external_e_pg.NameAlias)
		}

		if "exclude" != cloud_external_e_pg.PrefGrMemb {
			return fmt.Errorf("Bad cloud_external_e_pg pref_gr_memb %s", cloud_external_e_pg.PrefGrMemb)
		}

		if "unspecified" != cloud_external_e_pg.Prio {
			return fmt.Errorf("Bad cloud_external_e_pg prio %s", cloud_external_e_pg.Prio)
		}

		if "inter-site" != cloud_external_e_pg.RouteReachability {
			return fmt.Errorf("Bad cloud_external_e_pg route_reachability %s", cloud_external_e_pg.RouteReachability)
		}

		return nil
	}
}
