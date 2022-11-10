package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciUserDomain_Basic(t *testing.T) {
	var user_domain models.UserDomain
	description := "user_domain created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciUserDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciUserDomainConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserDomainExists("aci_user_security_domain.foouser_domain", &user_domain),
					testAccCheckAciUserDomainAttributes(description, &user_domain),
				),
			},
		},
	})
}

func TestAccAciUserDomain_Update(t *testing.T) {
	var user_domain models.UserDomain
	description := "user_domain created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciUserDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciUserDomainConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserDomainExists("aci_user_security_domain.foouser_domain", &user_domain),
					testAccCheckAciUserDomainAttributes(description, &user_domain),
				),
			},
			{
				Config: testAccCheckAciUserDomainConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserDomainExists("aci_user_security_domain.foouser_domain", &user_domain),
					testAccCheckAciUserDomainAttributes(description, &user_domain),
				),
			},
		},
	})
}

func testAccCheckAciUserDomainConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_user_security_domain" "foouser_domain" {
		name 		= "test"
		description = "%s"
		local_user_dn = "uni/userext/user-dhaval"
		annotation = "test_annotation"
		name_alias = "test_name_alias"
	}
	`, description)
}

func testAccCheckAciUserDomainExists(name string, user_domain *models.UserDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("User Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		user_domainFound := models.UserDomainFromContainer(cont)
		if user_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("User Domain %s not found", rs.Primary.ID)
		}
		*user_domain = *user_domainFound
		return nil
	}
}

func testAccCheckAciUserDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_user_security_domain" {
			cont, err := client.Get(rs.Primary.ID)
			user_domain := models.UserDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("User Domain %s Still exists", user_domain.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciUserDomainAttributes(description string, user_domain *models.UserDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != user_domain.Name {
			return fmt.Errorf("Bad User Domain Name %s", user_domain.Name)
		}

		if description != user_domain.Description {
			return fmt.Errorf("Bad User Domain Description %s", user_domain.Description)
		}

		if "test_annotation" != user_domain.Annotation {
			return fmt.Errorf("Bad User Domain Annotation %s", user_domain.Annotation)
		}

		if "test_name_alias" != user_domain.NameAlias {
			return fmt.Errorf("Bad User Domain NameAlias %s", user_domain.NameAlias)
		}
		return nil
	}
}
