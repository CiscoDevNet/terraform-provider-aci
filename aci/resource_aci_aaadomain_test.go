package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSecurityDomain_Basic(t *testing.T) {
	var security_domain models.SecurityDomain
	description := "aaa domain created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSecurityDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSecurityDomainConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSecurityDomainExists("aci_aaa_domain.foosecurity_domain", &security_domain),
					testAccCheckAciSecurityDomainAttributes(description, &security_domain),
				),
			},
		},
	})
}

func TestAccAciSecurityDomain_update(t *testing.T) {
	var security_domain models.SecurityDomain
	description := "aaa domain created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSecurityDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSecurityDomainConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSecurityDomainExists("aci_aaa_domain.foosecurity_domain", &security_domain),
					testAccCheckAciSecurityDomainAttributes(description, &security_domain),
				),
			},
			{
				Config: testAccCheckAciSecurityDomainConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSecurityDomainExists("aci_aaa_domain.foosecurity_domain", &security_domain),
					testAccCheckAciSecurityDomainAttributes(description, &security_domain),
				),
			},
		},
	})
}

func testAccCheckAciSecurityDomainConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_aaa_domain" "foosecurity_domain" {
		description = "%s"
		name  = "aaa_domain_1"
		annotation  = "example"
		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciSecurityDomainExists(name string, security_domain *models.SecurityDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("aaa domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No aaa domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		security_domainFound := models.SecurityDomainFromContainer(cont)
		if security_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("aaa domain %s not found", rs.Primary.ID)
		}
		*security_domain = *security_domainFound
		return nil
	}
}

func testAccCheckAciSecurityDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_aaa_domain" {
			cont, err := client.Get(rs.Primary.ID)
			security_domain := models.SecurityDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("aaa domain %s Still exists", security_domain.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSecurityDomainAttributes(description string, security_domain *models.SecurityDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != security_domain.Description {
			return fmt.Errorf("Bad aaa domain Description %s", security_domain.Description)
		}

		if "aaa_domain_1" != security_domain.Name {
			return fmt.Errorf("Bad aaa domain name %s", security_domain.Name)
		}

		if "example" != security_domain.Annotation {
			return fmt.Errorf("Bad aaa domain annotation %s", security_domain.Annotation)
		}

		if "example" != security_domain.NameAlias {
			return fmt.Errorf("Bad aaa domain name_alias %s", security_domain.NameAlias)
		}

		return nil
	}
}
