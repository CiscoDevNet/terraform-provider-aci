package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciApplicationProfile_Basic(t *testing.T) {
	var application_profile models.ApplicationProfile
	description := "application_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciApplicationProfileConfig_basic(description, "unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists("aci_application_profile.fooapplication_profile", &application_profile),
					testAccCheckAciApplicationProfileAttributes(description, "unspecified", &application_profile),
				),
			},
			{
				ResourceName:      "aci_application_profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciApplicationProfile_update(t *testing.T) {
	var application_profile models.ApplicationProfile
	description := "application_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciApplicationProfileConfig_basic(description, "unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists("aci_application_profile.fooapplication_profile", &application_profile),
					testAccCheckAciApplicationProfileAttributes(description, "unspecified", &application_profile),
				),
			},
			{
				Config: testAccCheckAciApplicationProfileConfig_basic(description, "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationProfileExists("aci_application_profile.fooapplication_profile", &application_profile),
					testAccCheckAciApplicationProfileAttributes(description, "level2", &application_profile),
				),
			},
		},
	})
}

func testAccCheckAciApplicationProfileConfig_basic(description, prio string) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "tenant_for_ap" {
		name        = "tenant_for_ap"
		description = "This tenant is created by terraform ACI provider"
	}
	resource "aci_application_profile" "fooapplication_profile" {
		tenant_dn   = aci_tenant.tenant_for_ap.id
		description = "%s"
		name        = "demo_ap"
		annotation  = "tag_ap"
		name_alias  = "alias_ap"
		prio        = "%s"
	}
	`, description, prio)
}

func testAccCheckAciApplicationProfileExists(name string, application_profile *models.ApplicationProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Application Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Application Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		application_profileFound := models.ApplicationProfileFromContainer(cont)
		if application_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Application Profile %s not found", rs.Primary.ID)
		}
		*application_profile = *application_profileFound
		return nil
	}
}

func testAccCheckAciApplicationProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_application_profile" {
			cont, err := client.Get(rs.Primary.ID)
			application_profile := models.ApplicationProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Application Profile %s Still exists", application_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciApplicationProfileAttributes(description, prio string, application_profile *models.ApplicationProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != application_profile.Description {
			return fmt.Errorf("Bad application_profile Description %s", application_profile.Description)
		}

		if "demo_ap" != application_profile.Name {
			return fmt.Errorf("Bad application_profile name %s", application_profile.Name)
		}

		if "tag_ap" != application_profile.Annotation {
			return fmt.Errorf("Bad application_profile annotation %s", application_profile.Annotation)
		}

		if "alias_ap" != application_profile.NameAlias {
			return fmt.Errorf("Bad application_profile name_alias %s", application_profile.NameAlias)
		}

		if prio != application_profile.Prio {
			return fmt.Errorf("Bad application_profile prio %s", application_profile.Prio)
		}

		return nil
	}
}
