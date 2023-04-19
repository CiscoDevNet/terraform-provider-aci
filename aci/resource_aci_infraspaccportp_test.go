package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSpineInterfaceProfile_Basic(t *testing.T) {
	var spine_interface_profile models.SpineInterfaceProfile
	description := "spine_interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpineInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpineInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineInterfaceProfileExists("aci_spine_interface_profile.foospine_interface_profile", &spine_interface_profile),
					testAccCheckAciSpineInterfaceProfileAttributes(description, &spine_interface_profile),
				),
			},
		},
	})
}

func TestAccAciSpineInterfaceProfile_update(t *testing.T) {
	var spine_interface_profile models.SpineInterfaceProfile
	description := "spine_interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpineInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpineInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineInterfaceProfileExists("aci_spine_interface_profile.foospine_interface_profile", &spine_interface_profile),
					testAccCheckAciSpineInterfaceProfileAttributes(description, &spine_interface_profile),
				),
			},
			{
				Config: testAccCheckAciSpineInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineInterfaceProfileExists("aci_spine_interface_profile.foospine_interface_profile", &spine_interface_profile),
					testAccCheckAciSpineInterfaceProfileAttributes(description, &spine_interface_profile),
				),
			},
		},
	})
}

func testAccCheckAciSpineInterfaceProfileConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_spine_interface_profile" "foospine_interface_profile" {
		description = "%s"	
		name  = "example"
		annotation  = "example"
		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciSpineInterfaceProfileExists(name string, spine_interface_profile *models.SpineInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spine Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spine Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spine_interface_profileFound := models.SpineInterfaceProfileFromContainer(cont)
		if spine_interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spine Interface Profile %s not found", rs.Primary.ID)
		}
		*spine_interface_profile = *spine_interface_profileFound
		return nil
	}
}

func testAccCheckAciSpineInterfaceProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_spine_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			spine_interface_profile := models.SpineInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spine Interface Profile %s Still exists", spine_interface_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSpineInterfaceProfileAttributes(description string, spine_interface_profile *models.SpineInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != spine_interface_profile.Description {
			return fmt.Errorf("Bad spine_interface_profile Description %s", spine_interface_profile.Description)
		}

		if "example" != spine_interface_profile.Name {
			return fmt.Errorf("Bad spine_interface_profile name %s", spine_interface_profile.Name)
		}

		if "example" != spine_interface_profile.Annotation {
			return fmt.Errorf("Bad spine_interface_profile annotation %s", spine_interface_profile.Annotation)
		}

		if "example" != spine_interface_profile.NameAlias {
			return fmt.Errorf("Bad spine_interface_profile name_alias %s", spine_interface_profile.NameAlias)
		}

		return nil
	}
}
