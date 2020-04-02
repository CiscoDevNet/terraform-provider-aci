package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciLocalUser_Basic(t *testing.T) {
	var local_user models.LocalUser
	description := "local_user created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLocalUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLocalUserConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLocalUserExists("aci_local_user.foolocal_user", &local_user),
					testAccCheckAciLocalUserAttributes(description, &local_user),
				),
			},
			{
				ResourceName:      "aci_local_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciLocalUser_update(t *testing.T) {
	var local_user models.LocalUser
	description := "local_user created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLocalUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLocalUserConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLocalUserExists("aci_local_user.foolocal_user", &local_user),
					testAccCheckAciLocalUserAttributes(description, &local_user),
				),
			},
			{
				Config: testAccCheckAciLocalUserConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLocalUserExists("aci_local_user.foolocal_user", &local_user),
					testAccCheckAciLocalUserAttributes(description, &local_user),
				),
			},
		},
	})
}

func testAccCheckAciLocalUserConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_local_user" "foolocal_user" {
		description = "%s"
		
		name  = "example"
		  account_status  = "example"
		  annotation  = "example"
		  cert_attribute  = "example"
		  clear_pwd_history  = "example"
		  email  = "example"
		  expiration  = "example"
		  expires  = "example"
		  first_name  = "example"
		  last_name  = "example"
		  name_alias  = "example"
		  otpenable  = "example"
		  otpkey  = "example"
		  phone  = "example"
		  pwd  = "example"
		  pwd_life_time  = "example"
		  pwd_update_required  = "example"
		  rbac_string  = "example"
		  unix_user_id  = "example"
		}
	`, description)
}

func testAccCheckAciLocalUserExists(name string, local_user *models.LocalUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Local User %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Local User dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		local_userFound := models.LocalUserFromContainer(cont)
		if local_userFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Local User %s not found", rs.Primary.ID)
		}
		*local_user = *local_userFound
		return nil
	}
}

func testAccCheckAciLocalUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_local_user" {
			cont, err := client.Get(rs.Primary.ID)
			local_user := models.LocalUserFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Local User %s Still exists", local_user.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLocalUserAttributes(description string, local_user *models.LocalUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != local_user.Description {
			return fmt.Errorf("Bad local_user Description %s", local_user.Description)
		}

		if "example" != local_user.Name {
			return fmt.Errorf("Bad local_user name %s", local_user.Name)
		}

		if "example" != local_user.AccountStatus {
			return fmt.Errorf("Bad local_user account_status %s", local_user.AccountStatus)
		}

		if "example" != local_user.Annotation {
			return fmt.Errorf("Bad local_user annotation %s", local_user.Annotation)
		}

		if "example" != local_user.CertAttribute {
			return fmt.Errorf("Bad local_user cert_attribute %s", local_user.CertAttribute)
		}

		if "example" != local_user.ClearPwdHistory {
			return fmt.Errorf("Bad local_user clear_pwd_history %s", local_user.ClearPwdHistory)
		}

		if "example" != local_user.Email {
			return fmt.Errorf("Bad local_user email %s", local_user.Email)
		}

		if "example" != local_user.Expiration {
			return fmt.Errorf("Bad local_user expiration %s", local_user.Expiration)
		}

		if "example" != local_user.Expires {
			return fmt.Errorf("Bad local_user expires %s", local_user.Expires)
		}

		if "example" != local_user.FirstName {
			return fmt.Errorf("Bad local_user first_name %s", local_user.FirstName)
		}

		if "example" != local_user.LastName {
			return fmt.Errorf("Bad local_user last_name %s", local_user.LastName)
		}

		if "example" != local_user.NameAlias {
			return fmt.Errorf("Bad local_user name_alias %s", local_user.NameAlias)
		}

		if "example" != local_user.Otpenable {
			return fmt.Errorf("Bad local_user otpenable %s", local_user.Otpenable)
		}

		if "example" != local_user.Otpkey {
			return fmt.Errorf("Bad local_user otpkey %s", local_user.Otpkey)
		}

		if "example" != local_user.Phone {
			return fmt.Errorf("Bad local_user phone %s", local_user.Phone)
		}

		if "example" != local_user.Pwd {
			return fmt.Errorf("Bad local_user pwd %s", local_user.Pwd)
		}

		if "example" != local_user.PwdLifeTime {
			return fmt.Errorf("Bad local_user pwd_life_time %s", local_user.PwdLifeTime)
		}

		if "example" != local_user.PwdUpdateRequired {
			return fmt.Errorf("Bad local_user pwd_update_required %s", local_user.PwdUpdateRequired)
		}

		if "example" != local_user.RbacString {
			return fmt.Errorf("Bad local_user rbac_string %s", local_user.RbacString)
		}

		if "example" != local_user.UnixUserId {
			return fmt.Errorf("Bad local_user unix_user_id %s", local_user.UnixUserId)
		}

		return nil
	}
}
