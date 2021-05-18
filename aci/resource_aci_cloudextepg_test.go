package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciCloudExternalEPg_Basic(t *testing.T) {
	var cloud_external_epg models.CloudExternalEPg
	description := "cloud_external_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudExternalEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudExternalEPgConfig_basic(description, "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists("aci_cloud_external_epg.foocloud_external_epg", &cloud_external_epg),
					testAccCheckAciCloudExternalEPgAttributes(description, "All", &cloud_external_epg),
				),
			},
		},
	})
}

func TestAccAciCloudExternalEPg_update(t *testing.T) {
	var cloud_external_epg models.CloudExternalEPg
	description := "cloud_external_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudExternalEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudExternalEPgConfig_basic(description, "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists("aci_cloud_external_epg.foocloud_external_epg", &cloud_external_epg),
					testAccCheckAciCloudExternalEPgAttributes(description, "All", &cloud_external_epg),
				),
			},
			{
				Config: testAccCheckAciCloudExternalEPgConfig_basic(description, "AtleastOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudExternalEPgExists("aci_cloud_external_epg.foocloud_external_epg", &cloud_external_epg),
					testAccCheckAciCloudExternalEPgAttributes(description, "AtleastOne", &cloud_external_epg),
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

	resource "aci_cloud_external_epg" "foocloud_external_epg" {
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

func testAccCheckAciCloudExternalEPgExists(name string, cloud_external_epg *models.CloudExternalEPg) resource.TestCheckFunc {
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

		cloud_external_epgFound := models.CloudExternalEPgFromContainer(cont)
		if cloud_external_epgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud External EPg %s not found", rs.Primary.ID)
		}
		*cloud_external_epg = *cloud_external_epgFound
		return nil
	}
}

func testAccCheckAciCloudExternalEPgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_external_epg" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_external_epg := models.CloudExternalEPgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud External EPg %s Still exists", cloud_external_epg.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudExternalEPgAttributes(description, match_t string, cloud_external_epg *models.CloudExternalEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_external_epg.Description {
			return fmt.Errorf("Bad cloud_external_epg Description %s", cloud_external_epg.Description)
		}

		if "cloud_ext_epg" != cloud_external_epg.Name {
			return fmt.Errorf("Bad cloud_external_epg name %s", cloud_external_epg.Name)
		}

		if "tag_ext_epg" != cloud_external_epg.Annotation {
			return fmt.Errorf("Bad cloud_external_epg annotation %s", cloud_external_epg.Annotation)
		}

		if "0" != cloud_external_epg.ExceptionTag {
			return fmt.Errorf("Bad cloud_external_epg exception_tag %s", cloud_external_epg.ExceptionTag)
		}

		if "disabled" != cloud_external_epg.FloodOnEncap {
			return fmt.Errorf("Bad cloud_external_epg flood_on_encap %s", cloud_external_epg.FloodOnEncap)
		}

		if match_t != cloud_external_epg.MatchT {
			return fmt.Errorf("Bad cloud_external_epg match_t %s", cloud_external_epg.MatchT)
		}

		if "alias_ext" != cloud_external_epg.NameAlias {
			return fmt.Errorf("Bad cloud_external_epg name_alias %s", cloud_external_epg.NameAlias)
		}

		if "exclude" != cloud_external_epg.PrefGrMemb {
			return fmt.Errorf("Bad cloud_external_epg pref_gr_memb %s", cloud_external_epg.PrefGrMemb)
		}

		if "unspecified" != cloud_external_epg.Prio {
			return fmt.Errorf("Bad cloud_external_epg prio %s", cloud_external_epg.Prio)
		}

		if "inter-site" != cloud_external_epg.RouteReachability {
			return fmt.Errorf("Bad cloud_external_epg route_reachability %s", cloud_external_epg.RouteReachability)
		}

		return nil
	}
}
