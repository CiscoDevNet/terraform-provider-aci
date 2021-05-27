package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLeafInterfaceProfile_Basic(t *testing.T) {
	var leaf_interface_profile models.LeafInterfaceProfile
	description := "leaf_interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLeafInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLeafInterfaceProfileConfig_basic(description, "alias_leaf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafInterfaceProfileExists("aci_leaf_interface_profile.fooleaf_interface_profile", &leaf_interface_profile),
					testAccCheckAciLeafInterfaceProfileAttributes(description, "alias_leaf", &leaf_interface_profile),
				),
			},
		},
	})
}

func TestAccAciLeafInterfaceProfile_update(t *testing.T) {
	var leaf_interface_profile models.LeafInterfaceProfile
	description := "leaf_interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLeafInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLeafInterfaceProfileConfig_basic(description, "alias_leaf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafInterfaceProfileExists("aci_leaf_interface_profile.fooleaf_interface_profile", &leaf_interface_profile),
					testAccCheckAciLeafInterfaceProfileAttributes(description, "alias_leaf", &leaf_interface_profile),
				),
			},
			{
				Config: testAccCheckAciLeafInterfaceProfileConfig_basic(description, "alias_update"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafInterfaceProfileExists("aci_leaf_interface_profile.fooleaf_interface_profile", &leaf_interface_profile),
					testAccCheckAciLeafInterfaceProfileAttributes(description, "alias_update", &leaf_interface_profile),
				),
			},
		},
	})
}

func testAccCheckAciLeafInterfaceProfileConfig_basic(description, name_alias string) string {
	return fmt.Sprintf(`

	resource "aci_leaf_interface_profile" "fooleaf_interface_profile" {
		description = "%s"
		name        = "demo_leaf_profile"
		annotation  = "tag_leaf"
		name_alias  = "%s"
	}
	`, description, name_alias)
}

func testAccCheckAciLeafInterfaceProfileExists(name string, leaf_interface_profile *models.LeafInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Leaf Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Leaf Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		leaf_interface_profileFound := models.LeafInterfaceProfileFromContainer(cont)
		if leaf_interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Leaf Interface Profile %s not found", rs.Primary.ID)
		}
		*leaf_interface_profile = *leaf_interface_profileFound
		return nil
	}
}

func testAccCheckAciLeafInterfaceProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_leaf_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			leaf_interface_profile := models.LeafInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Leaf Interface Profile %s Still exists", leaf_interface_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLeafInterfaceProfileAttributes(description, name_alias string, leaf_interface_profile *models.LeafInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != leaf_interface_profile.Description {
			return fmt.Errorf("Bad leaf_interface_profile Description %s", leaf_interface_profile.Description)
		}

		if "demo_leaf_profile" != leaf_interface_profile.Name {
			return fmt.Errorf("Bad leaf_interface_profile name %s", leaf_interface_profile.Name)
		}

		if "tag_leaf" != leaf_interface_profile.Annotation {
			return fmt.Errorf("Bad leaf_interface_profile annotation %s", leaf_interface_profile.Annotation)
		}

		if name_alias != leaf_interface_profile.NameAlias {
			return fmt.Errorf("Bad leaf_interface_profile name_alias %s", leaf_interface_profile.NameAlias)
		}

		return nil
	}
}
