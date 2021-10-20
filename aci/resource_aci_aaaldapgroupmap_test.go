package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLDAPGroupMap_Basic(t *testing.T) {
	var ldap_group_map models.LDAPGroupMap
	aaa_ldap_group_map_name := acctest.RandString(5)
	description := "ldap_group_map created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLDAPGroupMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLDAPGroupMapConfig_basic(aaa_ldap_group_map_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapExists("aci_ldap_group_map.fooldap_group_map", &ldap_group_map),
					testAccCheckAciLDAPGroupMapAttributes(aaa_ldap_group_map_name, description, &ldap_group_map),
				),
			},
		},
	})
}

func TestAccAciLDAPGroupMap_Update(t *testing.T) {
	var ldap_group_map models.LDAPGroupMap
	aaa_ldap_group_map_name := acctest.RandString(5)
	description := "ldap_group_map created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLDAPGroupMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLDAPGroupMapConfig_basic(aaa_ldap_group_map_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapExists("aci_ldap_group_map.fooldap_group_map", &ldap_group_map),
					testAccCheckAciLDAPGroupMapAttributes(aaa_ldap_group_map_name, description, &ldap_group_map),
				),
			},
			{
				Config: testAccCheckAciLDAPGroupMapConfig_basic(aaa_ldap_group_map_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapExists("aci_ldap_group_map.fooldap_group_map", &ldap_group_map),
					testAccCheckAciLDAPGroupMapAttributes(aaa_ldap_group_map_name, description, &ldap_group_map),
				),
			},
		},
	})
}

func testAccCheckAciLDAPGroupMapConfig_basic(aaa_ldap_group_map_name string) string {
	return fmt.Sprintf(`

	resource "aci_ldap_group_map" "fooldap_group_map" {
		name 		= "%s"
		description = "ldap_group_map created while acceptance testing"

	}

	`, aaa_ldap_group_map_name)
}

func testAccCheckAciLDAPGroupMapExists(name string, ldap_group_map *models.LDAPGroupMap) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LDAP Group Map %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LDAP Group Map dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ldap_group_mapFound := models.LDAPGroupMapFromContainer(cont)
		if ldap_group_mapFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LDAP Group Map %s not found", rs.Primary.ID)
		}
		*ldap_group_map = *ldap_group_mapFound
		return nil
	}
}

func testAccCheckAciLDAPGroupMapDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ldap_group_map" {
			cont, err := client.Get(rs.Primary.ID)
			ldap_group_map := models.LDAPGroupMapFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LDAP Group Map %s Still exists", ldap_group_map.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLDAPGroupMapAttributes(aaa_ldap_group_map_name, description string, ldap_group_map *models.LDAPGroupMap) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aaa_ldap_group_map_name != GetMOName(ldap_group_map.DistinguishedName) {
			return fmt.Errorf("Bad aaa_ldap_group_map %s", GetMOName(ldap_group_map.DistinguishedName))
		}

		if description != ldap_group_map.Description {
			return fmt.Errorf("Bad ldap_group_map Description %s", ldap_group_map.Description)
		}
		return nil
	}
}
