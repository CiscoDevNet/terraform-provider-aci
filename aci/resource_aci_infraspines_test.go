package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSwitchSpineAssociation_Basic(t *testing.T) {
	var switch_association models.SwitchSpineAssociation
	description := "switch_association created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSwitchSpineAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSwitchSpineAssociationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchSpineAssociationExists("aci_spine_switch_association.fooswitch_association", &switch_association),
					testAccCheckAciSwitchSpineAssociationAttributes(description, &switch_association),
				),
			},
		},
	})
}

func TestAccAciSwitchSpineAssociation_update(t *testing.T) {
	var switch_association models.SwitchSpineAssociation
	description := "switch_association created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSwitchSpineAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSwitchSpineAssociationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchSpineAssociationExists("aci_spine_switch_association.fooswitch_association", &switch_association),
					testAccCheckAciSwitchSpineAssociationAttributes(description, &switch_association),
				),
			},
			{
				Config: testAccCheckAciSwitchSpineAssociationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchSpineAssociationExists("aci_spine_switch_association.fooswitch_association", &switch_association),
					testAccCheckAciSwitchSpineAssociationAttributes(description, &switch_association),
				),
			},
		},
	})
}

func testAccCheckAciSwitchSpineAssociationConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_spine_profile" "foospine_profile" {		
		name  = "spine_profile_1"
	}

	resource "aci_spine_switch_association" "fooswitch_association" {
		spine_profile_dn  = aci_spine_profile.foospine_profile.id
		description = "%s"
		name  = "spine_switch_association_1"
		spine_switch_association_type  = "range"
		annotation  = "spine_switch_association_tag"
		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciSwitchSpineAssociationExists(name string, switch_association *models.SwitchSpineAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Switch Association %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Switch Association dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		switch_associationFound := models.SwitchSpineAssociationFromContainer(cont)
		if switch_associationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Switch Association %s not found", rs.Primary.ID)
		}
		*switch_association = *switch_associationFound
		return nil
	}
}

func testAccCheckAciSwitchSpineAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_spine_switch_association" {
			cont, err := client.Get(rs.Primary.ID)
			switch_association := models.SwitchSpineAssociationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Switch Association %s Still exists", switch_association.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSwitchSpineAssociationAttributes(description string, switch_association *models.SwitchSpineAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != switch_association.Description {
			return fmt.Errorf("Bad switch_association Description %s", switch_association.Description)
		}

		if "spine_switch_association_1" != switch_association.Name {
			return fmt.Errorf("Bad switch_association name %s", switch_association.Name)
		}

		if "range" != switch_association.SwitchAssociationType {
			return fmt.Errorf("Bad switch_association switch_association_type %s", switch_association.SwitchAssociationType)
		}

		if "spine_switch_association_tag" != switch_association.Annotation {
			return fmt.Errorf("Bad switch_association annotation %s", switch_association.Annotation)
		}

		if "example" != switch_association.NameAlias {
			return fmt.Errorf("Bad switch_association name_alias %s", switch_association.NameAlias)
		}

		return nil
	}
}
