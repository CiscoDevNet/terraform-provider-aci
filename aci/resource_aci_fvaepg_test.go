package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciApplicationEPG_Basic(t *testing.T) {
	var application_epg models.ApplicationEPG
	description := "application_epg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationEPGDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciApplicationEPGConfig_basic(description, "unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists("aci_application_epg.fooapplication_epg", &application_epg),
					testAccCheckAciApplicationEPGAttributes(description, "unspecified", &application_epg),
				),
			},
		},
	})
}

func TestAccAciApplicationEPG_update(t *testing.T) {
	var application_epg models.ApplicationEPG
	description := "application_epg created while acceptance testing"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationEPGDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciApplicationEPGConfig_basic(description, "unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists("aci_application_epg.fooapplication_epg", &application_epg),
					testAccCheckAciApplicationEPGAttributes(description, "unspecified", &application_epg),
				),
			},
			{
				Config: testAccCheckAciApplicationEPGConfig_basic(description, "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists("aci_application_epg.fooapplication_epg", &application_epg),
					testAccCheckAciApplicationEPGAttributes(description, "level3", &application_epg),
				),
			},
		},
	})
}

func testAccCheckAciApplicationEPGConfig_basic(description, prio string) string {
	return fmt.Sprintf(`
		resource "aci_tenant" "tenant_for_epg" {
			name        = "tenant_for_epg"
			description = "This tenant is created by terraform ACI provider"
		}
	  
		resource "aci_application_profile" "app_profile_for_epg" {
			tenant_dn   = aci_tenant.tenant_for_epg.id
			name        = "ap_for_epg"
			description = "This app profile is created by terraform ACI providers"
		}

		resource "aci_application_epg" "fooapplication_epg" {
			application_profile_dn  = aci_application_profile.app_profile_for_epg.id
			name  					= "demo_epg"
			description 			= "%s"
			annotation  			= "tag_epg"
			exception_tag 			= "0"
			flood_on_encap  		= "disabled"
			fwd_ctrl  				= "none"
			has_mcast_source  		= "no"
			is_attr_based_epg  	    = "no"
			match_t  				= "AtleastOne"
			name_alias  			= "alias_epg"
			pc_enf_pref  			= "unenforced"
			pref_gr_memb  			= "exclude"
			prio  					= "%s"
			shutdown  				= "no"
		}
	`, description, prio)
}

func testAccCheckAciApplicationEPGExists(name string, application_epg *models.ApplicationEPG) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Application EPG %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Application EPG dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		application_epgFound := models.ApplicationEPGFromContainer(cont)
		if application_epgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Application EPG %s not found", rs.Primary.ID)
		}
		*application_epg = *application_epgFound
		return nil
	}
}

func testAccCheckAciApplicationEPGDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_application_epg" {
			cont, err := client.Get(rs.Primary.ID)
			application_epg := models.ApplicationEPGFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Application EPG %s Still exists", application_epg.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciApplicationEPGAttributes(description, prio string, application_epg *models.ApplicationEPG) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != application_epg.Description {
			return fmt.Errorf("Bad application_epg Description %s", application_epg.Description)
		}

		if "demo_epg" != application_epg.Name {
			return fmt.Errorf("Bad application_epg name %s", application_epg.Name)
		}

		if "tag_epg" != application_epg.Annotation {
			return fmt.Errorf("Bad application_epg annotation %s", application_epg.Annotation)
		}

		if "0" != application_epg.ExceptionTag {
			return fmt.Errorf("Bad application_epg exception_tag %s", application_epg.ExceptionTag)
		}

		if "disabled" != application_epg.FloodOnEncap {
			return fmt.Errorf("Bad application_epg flood_on_encap %s", application_epg.FloodOnEncap)
		}

		if "no" != application_epg.HasMcastSource {
			return fmt.Errorf("Bad application_epg has_mcast_source %s", application_epg.HasMcastSource)
		}

		if "no" != application_epg.IsAttrBasedEPg {
			return fmt.Errorf("Bad application_epg is_attr_based_epg %s", application_epg.IsAttrBasedEPg)
		}

		if "AtleastOne" != application_epg.MatchT {
			return fmt.Errorf("Bad application_epg match_t %s", application_epg.MatchT)
		}

		if "alias_epg" != application_epg.NameAlias {
			return fmt.Errorf("Bad application_epg name_alias %s", application_epg.NameAlias)
		}

		if "unenforced" != application_epg.PcEnfPref {
			return fmt.Errorf("Bad application_epg pc_enf_pref %s", application_epg.PcEnfPref)
		}

		if "exclude" != application_epg.PrefGrMemb {
			return fmt.Errorf("Bad application_epg pref_gr_memb %s", application_epg.PrefGrMemb)
		}

		if prio != application_epg.Prio {
			return fmt.Errorf("Bad application_epg prio %s", application_epg.Prio)
		}

		if "no" != application_epg.Shutdown {
			return fmt.Errorf("Bad application_epg shutdown %s", application_epg.Shutdown)
		}
		return nil
	}
}
