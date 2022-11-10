package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciUserRole_Basic(t *testing.T) {
	var user_role models.UserRole
	description := "user_role created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciUserRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciUserRoleConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserRoleExists("aci_user_security_domain_role.foouser_role", &user_role),
					testAccCheckAciUserRoleAttributes(description, &user_role),
				),
			},
		},
	})
}

func TestAccAciUserRole_Update(t *testing.T) {
	var user_role models.UserRole
	description := "user_role created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciUserRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciUserRoleConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserRoleExists("aci_user_security_domain_role.foouser_role", &user_role),
					testAccCheckAciUserRoleAttributes(description, &user_role),
				),
			},
			{
				Config: testAccCheckAciUserRoleConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserRoleExists("aci_user_security_domain_role.foouser_role", &user_role),
					testAccCheckAciUserRoleAttributes(description, &user_role),
				),
			},
		},
	})
}

func testAccCheckAciUserRoleConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_user_security_domain_role" "foouser_role" {
		name 		= "example_test"
		description = "%s"
		user_domain_dn = aci_user_security_domain_role.foouser_domain.id
		annotation = "from_Terraform"
		name_alias = "user_security_alias"
	}

	`, description)
}

func testAccCheckAciUserRoleExists(name string, user_role *models.UserRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("User Role %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User Role dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		user_roleFound := models.UserRoleFromContainer(cont)
		if user_roleFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("User Role %s not found", rs.Primary.ID)
		}
		*user_role = *user_roleFound
		return nil
	}
}

func testAccCheckAciUserRoleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_user_security_domain_role" {
			cont, err := client.Get(rs.Primary.ID)
			user_role := models.UserRoleFromContainer(cont)
			if err == nil {
				return fmt.Errorf("User Role %s Still exists", user_role.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciUserRoleAttributes(description string, user_role *models.UserRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "example_test" != user_role.Name {
			return fmt.Errorf("Bad aaa_user_role %s", user_role.Name)
		}
		if description != user_role.Description {
			return fmt.Errorf("Bad user_role Description %s", user_role.Description)
		}
		if "from_Terraform" != user_role.Annotation {
			return fmt.Errorf("Bad user_role Annotation %s", user_role.Annotation)
		}
		if "user_security_alias" != user_role.NameAlias {
			return fmt.Errorf("Bad user_role Name Alias %s", user_role.NameAlias)
		}
		return nil
	}
}
