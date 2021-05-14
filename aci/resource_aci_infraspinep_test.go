package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSpineProfile_Basic(t *testing.T) {
	var spine_profile models.SpineProfile
	description := "spine_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpineProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpineProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineProfileExists("aci_spine_profile.foospine_profile", &spine_profile),
					testAccCheckAciSpineProfileAttributes(description, &spine_profile),
				),
			},
		},
	})
}

func TestAccAciSpineProfile_update(t *testing.T) {
	var spine_profile models.SpineProfile
	description := "spine_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpineProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpineProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineProfileExists("aci_spine_profile.foospine_profile", &spine_profile),
					testAccCheckAciSpineProfileAttributes(description, &spine_profile),
				),
			},
			{
				Config: testAccCheckAciSpineProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineProfileExists("aci_spine_profile.foospine_profile", &spine_profile),
					testAccCheckAciSpineProfileAttributes(description, &spine_profile),
				),
			},
		},
	})
}

func testAccCheckAciSpineProfileConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_spine_profile" "foospine_profile" {
		description = "%s"
		name  = "example"
		annotation  = "example"
		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciSpineProfileExists(name string, spine_profile *models.SpineProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spine Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spine Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spine_profileFound := models.SpineProfileFromContainer(cont)
		if spine_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spine Profile %s not found", rs.Primary.ID)
		}
		*spine_profile = *spine_profileFound
		return nil
	}
}

func testAccCheckAciSpineProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_spine_profile" {
			cont, err := client.Get(rs.Primary.ID)
			spine_profile := models.SpineProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spine Profile %s Still exists", spine_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSpineProfileAttributes(description string, spine_profile *models.SpineProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != spine_profile.Description {
			return fmt.Errorf("Bad spine_profile Description %s", spine_profile.Description)
		}

		if "example" != spine_profile.Name {
			return fmt.Errorf("Bad spine_profile name %s", spine_profile.Name)
		}

		if "example" != spine_profile.Annotation {
			return fmt.Errorf("Bad spine_profile annotation %s", spine_profile.Annotation)
		}

		if "example" != spine_profile.NameAlias {
			return fmt.Errorf("Bad spine_profile name_alias %s", spine_profile.NameAlias)
		}

		return nil
	}
}
