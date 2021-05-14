package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLeafProfile_Basic(t *testing.T) {
	var leaf_profile models.LeafProfile
	description := "leaf_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLeafProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLeafProfileConfig_basic(description, "alias_node_ep"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafProfileExists("aci_leaf_profile.fooleaf_profile", &leaf_profile),
					testAccCheckAciLeafProfileAttributes(description, "alias_node_ep", &leaf_profile),
				),
			},
		},
	})
}

func TestAccAciLeafProfile_update(t *testing.T) {
	var leaf_profile models.LeafProfile
	description := "leaf_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLeafProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLeafProfileConfig_basic(description, "alias_node_ep"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafProfileExists("aci_leaf_profile.fooleaf_profile", &leaf_profile),
					testAccCheckAciLeafProfileAttributes(description, "alias_node_ep", &leaf_profile),
				),
			},
			{
				Config: testAccCheckAciLeafProfileConfig_basic(description, "alias_update_node"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafProfileExists("aci_leaf_profile.fooleaf_profile", &leaf_profile),
					testAccCheckAciLeafProfileAttributes(description, "alias_update_node", &leaf_profile),
				),
			},
		},
	})
}

func testAccCheckAciLeafProfileConfig_basic(description, name_alias string) string {
	return fmt.Sprintf(`

	resource "aci_leaf_profile" "fooleaf_profile" {
		description = "%s"
		name        = "demo_node_ep"
		annotation  = "tag_node_ep"
		name_alias  = "%s"
	}  
	`, description, name_alias)
}

func testAccCheckAciLeafProfileExists(name string, leaf_profile *models.LeafProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Leaf Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Leaf Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		leaf_profileFound := models.LeafProfileFromContainer(cont)
		if leaf_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Leaf Profile %s not found", rs.Primary.ID)
		}
		*leaf_profile = *leaf_profileFound
		return nil
	}
}

func testAccCheckAciLeafProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_leaf_profile" {
			cont, err := client.Get(rs.Primary.ID)
			leaf_profile := models.LeafProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Leaf Profile %s Still exists", leaf_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLeafProfileAttributes(description, name_alias string, leaf_profile *models.LeafProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != leaf_profile.Description {
			return fmt.Errorf("Bad leaf_profile Description %s", leaf_profile.Description)
		}

		if "demo_node_ep" != leaf_profile.Name {
			return fmt.Errorf("Bad leaf_profile name %s", leaf_profile.Name)
		}

		if "tag_node_ep" != leaf_profile.Annotation {
			return fmt.Errorf("Bad leaf_profile annotation %s", leaf_profile.Annotation)
		}

		if name_alias != leaf_profile.NameAlias {
			return fmt.Errorf("Bad leaf_profile name_alias %s", leaf_profile.NameAlias)
		}

		return nil
	}
}
