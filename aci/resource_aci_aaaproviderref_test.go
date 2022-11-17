package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciProviderGroupMember_Basic(t *testing.T) {
	var login_domain_provider models.ProviderGroupMember
	description := "login_domain_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciProviderGroupMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciProviderGroupMemberConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciProviderGroupMemberExists("aci_login_domain_provider.foologin_domain_provider", &login_domain_provider),
					testAccCheckAciProviderGroupMemberAttributes(description, &login_domain_provider),
				),
			},
		},
	})
}

func TestAccAciProviderGroupMember_Update(t *testing.T) {
	var login_domain_provider models.ProviderGroupMember
	description := "login_domain_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciProviderGroupMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciProviderGroupMemberConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciProviderGroupMemberExists("aci_login_domain_provider.foologin_domain_provider", &login_domain_provider),
					testAccCheckAciProviderGroupMemberAttributes(description, &login_domain_provider),
				),
			},
			{
				Config: testAccCheckAciProviderGroupMemberConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciProviderGroupMemberExists("aci_login_domain_provider.foologin_domain_provider", &login_domain_provider),
					testAccCheckAciProviderGroupMemberAttributes(description, &login_domain_provider),
				),
			},
		},
	})
}

func testAccCheckAciProviderGroupMemberConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_login_domain_provider" "foologin_domain_provider" {
		parent_dn = aci_tacacs_provider_group.test.id
		name = "test"
		description = "%s"
		name_alias = "test_name_alias"
		annotation = "test_annotation"
		order = "0"
	}

	`, description)
}

func testAccCheckAciProviderGroupMemberExists(name string, login_domain_provider *models.ProviderGroupMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Provider Group Member %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Provider Group Member dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		login_domain_providerFound := models.ProviderGroupMemberFromContainer(cont)
		if login_domain_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Provider Group Member %s not found", rs.Primary.ID)
		}
		*login_domain_provider = *login_domain_providerFound
		return nil
	}
}

func testAccCheckAciProviderGroupMemberDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_login_domain_provider" {
			cont, err := client.Get(rs.Primary.ID)
			login_domain_provider := models.ProviderGroupMemberFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Provider Group Member %s Still exists", login_domain_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciProviderGroupMemberAttributes(description string, login_domain_provider *models.ProviderGroupMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != GetMOName(login_domain_provider.DistinguishedName) {
			return fmt.Errorf("Bad aaa_provider_ref %s", GetMOName(login_domain_provider.DistinguishedName))
		}

		if description != login_domain_provider.Description {
			return fmt.Errorf("Bad login_domain_provider Description %s", login_domain_provider.Description)
		}

		if "test_name_alias" != login_domain_provider.NameAlias {
			return fmt.Errorf("Bad login_domain_provider NameAlias %s", login_domain_provider.NameAlias)
		}

		if "test_annotation" != login_domain_provider.Annotation {
			return fmt.Errorf("Bad login_domain_provider Annotation %s", login_domain_provider.Annotation)
		}

		if "0" != login_domain_provider.Order {
			return fmt.Errorf("Bad login_domain_provider Order %s", login_domain_provider.Order)
		}
		return nil
	}
}
