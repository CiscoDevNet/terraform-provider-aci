package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLoginDomain_Basic(t *testing.T) {
	var login_domain models.LoginDomain
	description := "login_domain created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLoginDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLoginDomainConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainExists("aci_login_domain.foologin_domain", &login_domain),
					testAccCheckAciLoginDomainAttributes(description, &login_domain),
				),
			},
		},
	})
}

func TestAccAciLoginDomain_Update(t *testing.T) {
	var login_domain models.LoginDomain
	description := "login_domain created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLoginDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLoginDomainConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainExists("aci_login_domain.foologin_domain", &login_domain),
					testAccCheckAciLoginDomainAttributes(description, &login_domain),
				),
			},
			{
				Config: testAccCheckAciLoginDomainConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoginDomainExists("aci_login_domain.foologin_domain", &login_domain),
					testAccCheckAciLoginDomainAttributes(description, &login_domain),
				),
			},
		},
	})
}

func testAccCheckAciLoginDomainConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_login_domain" "foologin_domain" {
		name 		= "test"
		description = "%s"
		name_alias  = "login_domain_alias"
		annotation  = "example"

	}

	`, description)
}

func testAccCheckAciLoginDomainExists(name string, login_domain *models.LoginDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Login Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Login Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		login_domainFound := models.LoginDomainFromContainer(cont)
		if login_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Login Domain %s not found", rs.Primary.ID)
		}
		*login_domain = *login_domainFound
		return nil
	}
}

func testAccCheckAciLoginDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_login_domain" {
			cont, err := client.Get(rs.Primary.ID)
			login_domain := models.LoginDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Login Domain %s Still exists", login_domain.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLoginDomainAttributes(description string, login_domain *models.LoginDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != login_domain.Name {
			return fmt.Errorf("Bad aaa_login_domain %s", login_domain.Name)
		}

		if description != login_domain.Description {
			return fmt.Errorf("Bad login_domain Description %s", login_domain.Description)
		}

		if "example" != login_domain.Annotation {
			return fmt.Errorf("Bad login_domain Annotation %s", login_domain.Annotation)
		}

		if "login_domain_alias" != login_domain.NameAlias {
			return fmt.Errorf("Bad login_domain Name Alias %s", login_domain.NameAlias)
		}
		return nil
	}
}
