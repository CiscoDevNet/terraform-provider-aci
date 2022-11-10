package aci

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciDuoProviderGroup_Basic(t *testing.T) {
	var duo_provider_group models.DuoProviderGroup
	description := "duo_provider_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDuoProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDuoProviderGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists("aci_duo_provider_group.fooduo_provider_group", &duo_provider_group),
					testAccCheckAciDuoProviderGroupAttributes(description, &duo_provider_group),
				),
			},
		},
	})
}

func TestAccAciDuoProviderGroup_Update(t *testing.T) {
	var duo_provider_group models.DuoProviderGroup
	description := "duo_provider_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDuoProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDuoProviderGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists("aci_duo_provider_group.fooduo_provider_group", &duo_provider_group),
					testAccCheckAciDuoProviderGroupAttributes(description, &duo_provider_group),
				),
			},
			{
				Config: testAccCheckAciDuoProviderGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDuoProviderGroupExists("aci_duo_provider_group.fooduo_provider_group", &duo_provider_group),
					testAccCheckAciDuoProviderGroupAttributes(description, &duo_provider_group),
				),
			},
		},
	})
}

func testAccCheckAciDuoProviderGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_duo_provider_group" "fooduo_provider_group" {
		name 		= "test"
		description = "%s"
		auth_choice = "CiscoAVPair"
		ldap_group_map_ref = "100"
		provider_type = "radius"
		sec_fac_auth_methods = ["auto"]
	}

	`, description)
}

func testAccCheckAciDuoProviderGroupExists(name string, duo_provider_group *models.DuoProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Duo Provider Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Duo Provider Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		duo_provider_groupFound := models.DuoProviderGroupFromContainer(cont)
		if duo_provider_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Duo Provider Group %s not found", rs.Primary.ID)
		}
		*duo_provider_group = *duo_provider_groupFound
		return nil
	}
}

func testAccCheckAciDuoProviderGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_duo_provider_group" {
			cont, err := client.Get(rs.Primary.ID)
			duo_provider_group := models.DuoProviderGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Duo Provider Group %s Still exists", duo_provider_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciDuoProviderGroupAttributes(description string, duo_provider_group *models.DuoProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != GetMOName(duo_provider_group.DistinguishedName) {
			return fmt.Errorf("Bad aaa_duo_provider_group %s", GetMOName(duo_provider_group.DistinguishedName))
		}

		if description != duo_provider_group.Description {
			return fmt.Errorf("Bad duo_provider_group Description %s", duo_provider_group.Description)
		}

		if "CiscoAVPair" != duo_provider_group.AuthChoice {
			return fmt.Errorf("Bad duo_provider_group AuthChoice %s", duo_provider_group.AuthChoice)
		}

		if "100" != duo_provider_group.LdapGroupMapRef {
			return fmt.Errorf("Bad duo_provider_group LdapGroupMapRef %s", duo_provider_group.LdapGroupMapRef)
		}

		if "radius" != duo_provider_group.ProviderType {
			return fmt.Errorf("Bad duo_provider_group ProviderType %s", duo_provider_group.ProviderType)
		}

		secFecAuthMet := strings.Split(duo_provider_group.SecFacAuthMethods, ",")
		if "auto" != secFecAuthMet[0] {
			return fmt.Errorf("Bad duo_provider_group secFecAuthMet %s", secFecAuthMet)
		}
		return nil
	}
}
