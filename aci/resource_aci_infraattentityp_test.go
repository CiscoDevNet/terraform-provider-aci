package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciAttachableAccessEntityProfile_Basic(t *testing.T) {
	var attachable_access_entity_profile models.AttachableAccessEntityProfile
	description := "attachable_access_entity_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAttachableAccessEntityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAttachableAccessEntityProfileConfig_basic(description, "alias_entity_prof"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAttachableAccessEntityProfileExists("aci_attachable_access_entity_profile.fooattachable_access_entity_profile", &attachable_access_entity_profile),
					testAccCheckAciAttachableAccessEntityProfileAttributes(description, "alias_entity_prof", &attachable_access_entity_profile),
				),
			},
		},
	})
}

func TestAccAciAttachableAccessEntityProfile_update(t *testing.T) {
	var attachable_access_entity_profile models.AttachableAccessEntityProfile
	description := "attachable_access_entity_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAttachableAccessEntityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAttachableAccessEntityProfileConfig_basic(description, "alias_entity_prof"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAttachableAccessEntityProfileExists("aci_attachable_access_entity_profile.fooattachable_access_entity_profile", &attachable_access_entity_profile),
					testAccCheckAciAttachableAccessEntityProfileAttributes(description, "alias_entity_prof", &attachable_access_entity_profile),
				),
			},
			{
				Config: testAccCheckAciAttachableAccessEntityProfileConfig_basic(description, "updated_alias"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAttachableAccessEntityProfileExists("aci_attachable_access_entity_profile.fooattachable_access_entity_profile", &attachable_access_entity_profile),
					testAccCheckAciAttachableAccessEntityProfileAttributes(description, "updated_alias", &attachable_access_entity_profile),
				),
			},
		},
	})
}

func testAccCheckAciAttachableAccessEntityProfileConfig_basic(description, name_alias string) string {
	return fmt.Sprintf(`
	resource "aci_attachable_access_entity_profile" "fooattachable_access_entity_profile" {
		description = "%s"
		name        = "demo_entity_prof"
		annotation  = "tag_entity"
		name_alias  = "%s"
	}
	`, description, name_alias)
}

func testAccCheckAciAttachableAccessEntityProfileExists(name string, attachable_access_entity_profile *models.AttachableAccessEntityProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Attachable Access Entity Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Attachable Access Entity Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		attachable_access_entity_profileFound := models.AttachableAccessEntityProfileFromContainer(cont)
		if attachable_access_entity_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Attachable Access Entity Profile %s not found", rs.Primary.ID)
		}
		*attachable_access_entity_profile = *attachable_access_entity_profileFound
		return nil
	}
}

func testAccCheckAciAttachableAccessEntityProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_attachable_access_entity_profile" {
			cont, err := client.Get(rs.Primary.ID)
			attachable_access_entity_profile := models.AttachableAccessEntityProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Attachable Access Entity Profile %s Still exists", attachable_access_entity_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciAttachableAccessEntityProfileAttributes(description, name_alias string, attachable_access_entity_profile *models.AttachableAccessEntityProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != attachable_access_entity_profile.Description {
			return fmt.Errorf("Bad attachable_access_entity_profile Description %s", attachable_access_entity_profile.Description)
		}

		if "demo_entity_prof" != attachable_access_entity_profile.Name {
			return fmt.Errorf("Bad attachable_access_entity_profile name %s", attachable_access_entity_profile.Name)
		}

		if "tag_entity" != attachable_access_entity_profile.Annotation {
			return fmt.Errorf("Bad attachable_access_entity_profile annotation %s", attachable_access_entity_profile.Annotation)
		}

		if name_alias != attachable_access_entity_profile.NameAlias {
			return fmt.Errorf("Bad attachable_access_entity_profile name_alias %s", attachable_access_entity_profile.NameAlias)
		}

		return nil
	}
}
