package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciCloudEPg_Basic(t *testing.T) {
	var cloud_epg models.CloudEPg
	description := "cloud_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudEPgConfig_basic(description, "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists("aci_cloud_epg.foocloud_epg", &cloud_epg),
					testAccCheckAciCloudEPgAttributes(description, "All", &cloud_epg),
				),
			},
		},
	})
}

func TestAccAciCloudEPg_update(t *testing.T) {
	var cloud_epg models.CloudEPg
	description := "cloud_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudEPgConfig_basic(description, "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists("aci_cloud_epg.foocloud_epg", &cloud_epg),
					testAccCheckAciCloudEPgAttributes(description, "All", &cloud_epg),
				),
			},
			{
				Config: testAccCheckAciCloudEPgConfig_basic(description, "AtleastOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudEPgExists("aci_cloud_epg.foocloud_epg", &cloud_epg),
					testAccCheckAciCloudEPgAttributes(description, "AtleastOne", &cloud_epg),
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
		tenant_dn   = aci_tenant.footenant.id
		name        = "demo_app"
		annotation  = "tag_app"
	}

	resource "aci_cloud_epg" "foocloud_epg" {
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.foocloud_applicationcontainer.id
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

func testAccCheckAciCloudEPgExists(name string, cloud_epg *models.CloudEPg) resource.TestCheckFunc {
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

		cloud_epgFound := models.CloudEPgFromContainer(cont)
		if cloud_epgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud EPg %s not found", rs.Primary.ID)
		}
		*cloud_epg = *cloud_epgFound
		return nil
	}
}

func testAccCheckAciCloudEPgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_epg" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_epg := models.CloudEPgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud EPg %s Still exists", cloud_epg.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudEPgAttributes(description, match_t string, cloud_epg *models.CloudEPg) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_epg.Description {
			return fmt.Errorf("Bad cloud_epg Description %s", cloud_epg.Description)
		}

		if "cloud_epg" != cloud_epg.Name {
			return fmt.Errorf("Bad cloud_epg name %s", cloud_epg.Name)
		}

		if "tag_epg" != cloud_epg.Annotation {
			return fmt.Errorf("Bad cloud_epg annotation %s", cloud_epg.Annotation)
		}

		if "0" != cloud_epg.ExceptionTag {
			return fmt.Errorf("Bad cloud_epg exception_tag %s", cloud_epg.ExceptionTag)
		}

		if "disabled" != cloud_epg.FloodOnEncap {
			return fmt.Errorf("Bad cloud_epg flood_on_encap %s", cloud_epg.FloodOnEncap)
		}

		if match_t != cloud_epg.MatchT {
			return fmt.Errorf("Bad cloud_epg match_t %s", cloud_epg.MatchT)
		}

		if "alias_epg" != cloud_epg.NameAlias {
			return fmt.Errorf("Bad cloud_epg name_alias %s", cloud_epg.NameAlias)
		}

		if "exclude" != cloud_epg.PrefGrMemb {
			return fmt.Errorf("Bad cloud_epg pref_gr_memb %s", cloud_epg.PrefGrMemb)
		}

		if "unspecified" != cloud_epg.Prio {
			return fmt.Errorf("Bad cloud_epg prio %s", cloud_epg.Prio)
		}

		return nil
	}
}
