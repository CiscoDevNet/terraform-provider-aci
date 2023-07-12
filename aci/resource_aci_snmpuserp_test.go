package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciUserProfile_Basic(t *testing.T) {
	var user_profile models.UserProfile
	snmp_pol_name := acctest.RandString(5)
	snmp_user_p_name := acctest.RandString(5)
	description := "user_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciUserProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciUserProfileConfig_basic(snmp_pol_name, snmp_user_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserProfileExists("aci_user_profile.foo_user_profile", &user_profile),
					testAccCheckAciUserProfileAttributes(snmp_pol_name, snmp_user_p_name, description, &user_profile),
				),
			},
		},
	})
}

func TestAccAciUserProfile_Update(t *testing.T) {
	var user_profile models.UserProfile
	snmp_pol_name := acctest.RandString(5)
	snmp_user_p_name := acctest.RandString(5)
	description := "user_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciUserProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciUserProfileConfig_basic(snmp_pol_name, snmp_user_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserProfileExists("aci_user_profile.foo_user_profile", &user_profile),
					testAccCheckAciUserProfileAttributes(snmp_pol_name, snmp_user_p_name, description, &user_profile),
				),
			},
			{
				Config: testAccCheckAciUserProfileConfig_basic(snmp_pol_name, snmp_user_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserProfileExists("aci_user_profile.foo_user_profile", &user_profile),
					testAccCheckAciUserProfileAttributes(snmp_pol_name, snmp_user_p_name, description, &user_profile),
				),
			},
		},
	})
}

func testAccCheckAciUserProfileConfig_basic(snmp_pol_name, snmp_user_p_name string) string {
	return fmt.Sprintf(`

	resource "aci_snmppolicy" "foo_snmppolicy" {
		name 		= "%s"
		description = "snmppolicy created while acceptance testing"

	}

	resource "aci_user_profile" "foo_user_profile" {
		name 		= "%s"
		description = "user_profile created while acceptance testing"
		snmppolicy_dn = aci_snmppolicy.foo_snmppolicy.id
	}

	`, snmp_pol_name, snmp_user_p_name)
}

func testAccCheckAciUserProfileExists(name string, user_profile *models.UserProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("User Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		user_profileFound := models.UserProfileFromContainer(cont)
		if user_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("User Profile %s not found", rs.Primary.ID)
		}
		*user_profile = *user_profileFound
		return nil
	}
}

func testAccCheckAciUserProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_user_profile" {
			cont, err := client.Get(rs.Primary.ID)
			user_profile := models.UserProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("User Profile %s Still exists", user_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciUserProfileAttributes(snmp_pol_name, snmp_user_p_name, description string, user_profile *models.UserProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if snmp_user_p_name != GetMOName(user_profile.DistinguishedName) {
			return fmt.Errorf("Bad snmpuser_p %s", GetMOName(user_profile.DistinguishedName))
		}

		if snmp_pol_name != GetMOName(GetParentDn(user_profile.DistinguishedName, user_profile.Rn)) {
			return fmt.Errorf(" Bad snmppol %s", GetMOName(GetParentDn(user_profile.DistinguishedName, user_profile.Rn)))
		}
		if description != user_profile.Description {
			return fmt.Errorf("Bad user_profile Description %s", user_profile.Description)
		}
		return nil
	}
}
