package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciUserManagement_Basic(t *testing.T) {
	var user_management models.UserManagement
	aaa_user_ep_name := acctest.RandString(5)
	description := "user_management created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciUserManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciUserManagementConfig_basic(aaa_user_ep_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserManagementExists("aci_user_management.foouser_management", &user_management),
					testAccCheckAciUserManagementAttributes(aaa_user_ep_name, description, &user_management),
				),
			},
		},
	})
}

func TestAccAciUserManagement_Update(t *testing.T) {
	var user_management models.UserManagement
	aaa_user_ep_name := acctest.RandString(5)
	description := "user_management created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciUserManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciUserManagementConfig_basic(aaa_user_ep_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserManagementExists("aci_user_management.foouser_management", &user_management),
					testAccCheckAciUserManagementAttributes(aaa_user_ep_name, description, &user_management),
				),
			},
			{
				Config: testAccCheckAciUserManagementConfig_basic(aaa_user_ep_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserManagementExists("aci_user_management.foouser_management", &user_management),
					testAccCheckAciUserManagementAttributes(aaa_user_ep_name, description, &user_management),
				),
			},
		},
	})
}

func testAccCheckAciUserManagementConfig_basic(aaa_user_ep_name string) string {
	return fmt.Sprintf(`

	resource "aci_user_management" "foouser_management" {
		name 		= "%s"
		description = "user_management created while acceptance testing"

	}

	`, aaa_user_ep_name)
}

func testAccCheckAciUserManagementExists(name string, user_management *models.UserManagement) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("User Management %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User Management dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		user_managementFound := models.UserManagementFromContainer(cont)
		if user_managementFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("User Management %s not found", rs.Primary.ID)
		}
		*user_management = *user_managementFound
		return nil
	}
}

func testAccCheckAciUserManagementDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_user_management" {
			cont, err := client.Get(rs.Primary.ID)
			user_management := models.UserManagementFromContainer(cont)
			if err == nil {
				return fmt.Errorf("User Management %s Still exists", user_management.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciUserManagementAttributes(aaa_user_ep_name, description string, user_management *models.UserManagement) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aaa_user_ep_name != GetMOName(user_management.DistinguishedName) {
			return fmt.Errorf("Bad aaa_user_ep %s", GetMOName(user_management.DistinguishedName))
		}

		if description != user_management.Description {
			return fmt.Errorf("Bad user_management Description %s", user_management.Description)
		}
		return nil
	}
}
