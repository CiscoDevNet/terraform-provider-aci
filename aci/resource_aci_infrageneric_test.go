package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciAccessGeneric_Basic(t *testing.T) {
	var access_generic models.AccessGeneric
	description := "access_generic created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessGenericDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessGenericConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessGenericExists("aci_access_generic.fooaccess_generic", &access_generic),
					testAccCheckAciAccessGenericAttributes(description, &access_generic),
				),
			},
		},
	})
}

func TestAccAciAccessGeneric_update(t *testing.T) {
	var access_generic models.AccessGeneric
	description := "access_generic created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessGenericDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessGenericConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessGenericExists("aci_access_generic.fooaccess_generic", &access_generic),
					testAccCheckAciAccessGenericAttributes(description, &access_generic),
				),
			},
			{
				Config: testAccCheckAciAccessGenericConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessGenericExists("aci_access_generic.fooaccess_generic", &access_generic),
					testAccCheckAciAccessGenericAttributes(description, &access_generic),
				),
			},
		},
	})
}

func testAccCheckAciAccessGenericConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_access_generic" "fooaccess_generic" {
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.example.id
		description = "%s"
		name        = "default"
		annotation  = "example"
		name_alias  = "access_generic"
	}
	`, description)
}

func testAccCheckAciAccessGenericExists(name string, access_generic *models.AccessGeneric) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Access Generic %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access Generic dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		access_genericFound := models.AccessGenericFromContainer(cont)
		if access_genericFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Access Generic %s not found", rs.Primary.ID)
		}
		*access_generic = *access_genericFound
		return nil
	}
}

func testAccCheckAciAccessGenericDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_access_generic" {
			cont, err := client.Get(rs.Primary.ID)
			access_generic := models.AccessGenericFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Access Generic %s Still exists", access_generic.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciAccessGenericAttributes(description string, access_generic *models.AccessGeneric) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != access_generic.Description {
			return fmt.Errorf("Bad access_generic Description %s", access_generic.Description)
		}

		if "default" != access_generic.Name {
			return fmt.Errorf("Bad access_generic name %s", access_generic.Name)
		}

		if "example" != access_generic.Annotation {
			return fmt.Errorf("Bad access_generic annotation %s", access_generic.Annotation)
		}

		if "access_generic" != access_generic.NameAlias {
			return fmt.Errorf("Bad access_generic name_alias %s", access_generic.NameAlias)
		}

		return nil
	}
}
