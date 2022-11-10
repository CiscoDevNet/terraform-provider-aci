package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSAMLProviderGroup_Basic(t *testing.T) {
	var saml_provider_group models.SAMLProviderGroup
	description := "saml_provider_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSAMLProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSAMLProviderGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderGroupExists("aci_saml_provider_group.foosaml_provider_group", &saml_provider_group),
					testAccCheckAciSAMLProviderGroupAttributes(description, &saml_provider_group),
				),
			},
		},
	})
}

func TestAccAciSAMLProviderGroup_Update(t *testing.T) {
	var saml_provider_group models.SAMLProviderGroup
	description := "saml_provider_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSAMLProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSAMLProviderGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderGroupExists("aci_saml_provider_group.foosaml_provider_group", &saml_provider_group),
					testAccCheckAciSAMLProviderGroupAttributes(description, &saml_provider_group),
				),
			},
			{
				Config: testAccCheckAciSAMLProviderGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderGroupExists("aci_saml_provider_group.foosaml_provider_group", &saml_provider_group),
					testAccCheckAciSAMLProviderGroupAttributes(description, &saml_provider_group),
				),
			},
		},
	})
}

func testAccCheckAciSAMLProviderGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_saml_provider_group" "foosaml_provider_group" {
		name 		= "test"
		description = "%s"
		annotation = "example"
		name_alias = "saml_provider_group_alias"

	}

	`, description)
}

func testAccCheckAciSAMLProviderGroupExists(name string, saml_provider_group *models.SAMLProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("SAML Provider Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SAML Provider Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		saml_provider_groupFound := models.SAMLProviderGroupFromContainer(cont)
		if saml_provider_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("SAML Provider Group %s not found", rs.Primary.ID)
		}
		*saml_provider_group = *saml_provider_groupFound
		return nil
	}
}

func testAccCheckAciSAMLProviderGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_saml_provider_group" {
			cont, err := client.Get(rs.Primary.ID)
			saml_provider_group := models.SAMLProviderGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("SAML Provider Group %s Still exists", saml_provider_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSAMLProviderGroupAttributes(description string, saml_provider_group *models.SAMLProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != saml_provider_group.Name {
			return fmt.Errorf("Bad aaa_saml_provider_group %s", saml_provider_group.Name)
		}

		if description != saml_provider_group.Description {
			return fmt.Errorf("Bad saml_provider_group Description %s", saml_provider_group.Description)
		}

		if "example" != saml_provider_group.Annotation {
			return fmt.Errorf("Bad saml_provider_group Annotation %s", saml_provider_group.Annotation)
		}

		if "saml_provider_group_alias" != saml_provider_group.NameAlias {
			return fmt.Errorf("Bad saml_provider_group Name Alias %s", saml_provider_group.NameAlias)
		}
		return nil
	}
}
