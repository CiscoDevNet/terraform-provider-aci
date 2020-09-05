package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciSwitchAssociation_Basic(t *testing.T) {
	var switch_association models.SwitchAssociation
	description := "switch_association created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSwitchAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSwitchAssociationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchAssociationExists("aci_leaf_selector.fooswitch_association", &switch_association),
					testAccCheckAciSwitchAssociationAttributes(description, &switch_association),
				),
			},
			{
				ResourceName:      "aci_leaf_selector",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciSwitchAssociation_update(t *testing.T) {
	var switch_association models.SwitchAssociation
	description := "switch_association created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSwitchAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSwitchAssociationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchAssociationExists("aci_leaf_selector.fooswitch_association", &switch_association),
					testAccCheckAciSwitchAssociationAttributes(description, &switch_association),
				),
			},
			{
				Config: testAccCheckAciSwitchAssociationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSwitchAssociationExists("aci_leaf_selector.fooswitch_association", &switch_association),
					testAccCheckAciSwitchAssociationAttributes(description, &switch_association),
				),
			},
		},
	})
}

func testAccCheckAciSwitchAssociationConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_leaf_selector" "fooswitch_association" {
		  leaf_profile_dn  = "${aci_leaf_profile.example.id}"
		description = "%s"

		name  = "example"

		switch_association_type  = "example"
		  annotation  = "example"
		  name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciSwitchAssociationExists(name string, switch_association *models.SwitchAssociation) resource.TestCheckFunc {
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

		switch_associationFound := models.SwitchAssociationFromContainer(cont)
		if switch_associationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Switch Association %s not found", rs.Primary.ID)
		}
		*switch_association = *switch_associationFound
		return nil
	}
}

func testAccCheckAciSwitchAssociationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_leaf_selector" {
			cont, err := client.Get(rs.Primary.ID)
			switch_association := models.SwitchAssociationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Switch Association %s Still exists", switch_association.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSwitchAssociationAttributes(description string, switch_association *models.SwitchAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != switch_association.Description {
			return fmt.Errorf("Bad switch_association Description %s", switch_association.Description)
		}

		if "example" != switch_association.Name {
			return fmt.Errorf("Bad switch_association name %s", switch_association.Name)
		}

		if "example" != switch_association.Switch_association_type {
			return fmt.Errorf("Bad switch_association switch_association_type %s", switch_association.Switch_association_type)
		}

		if "example" != switch_association.Annotation {
			return fmt.Errorf("Bad switch_association annotation %s", switch_association.Annotation)
		}

		if "example" != switch_association.NameAlias {
			return fmt.Errorf("Bad switch_association name_alias %s", switch_association.NameAlias)
		}

		return nil
	}
}
