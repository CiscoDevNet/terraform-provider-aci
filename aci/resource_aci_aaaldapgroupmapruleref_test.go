package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLDAPGroupMapruleref_Basic(t *testing.T) {
	var ldap_group_mapruleref models.LDAPGroupMapruleref
	description := "ldap_group_mapruleref created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLDAPGroupMaprulerefDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLDAPGroupMaprulerefConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMaprulerefExists("aci_ldap_group_map_rule_to_group_map.fooldap_group_mapruleref", &ldap_group_mapruleref),
					testAccCheckAciLDAPGroupMaprulerefAttributes(description, &ldap_group_mapruleref),
				),
			},
		},
	})
}

func TestAccAciLDAPGroupMapruleref_Update(t *testing.T) {
	var ldap_group_mapruleref models.LDAPGroupMapruleref
	description := "ldap_group_mapruleref created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLDAPGroupMaprulerefDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLDAPGroupMaprulerefConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMaprulerefExists("aci_ldap_group_map_rule_to_group_map.fooldap_group_mapruleref", &ldap_group_mapruleref),
					testAccCheckAciLDAPGroupMaprulerefAttributes(description, &ldap_group_mapruleref),
				),
			},
			{
				Config: testAccCheckAciLDAPGroupMaprulerefConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMaprulerefExists("aci_ldap_group_map_rule_to_group_map.fooldap_group_mapruleref", &ldap_group_mapruleref),
					testAccCheckAciLDAPGroupMaprulerefAttributes(description, &ldap_group_mapruleref),
				),
			},
		},
	})
}

func testAccCheckAciLDAPGroupMaprulerefConfig_basic(description string) string {
	return fmt.Sprintf(`	

	resource "aci_ldap_group_map_rule_to_group_map" "fooldap_group_mapruleref" {
		name = "test"
		description = "%s"
		ldap_group_map_dn = aci_ldap_group_map.test.id
		name_alias = "test_name_alias_value"
		annotation = "test_annotation_value"
	}
	`, description)
}

func testAccCheckAciLDAPGroupMaprulerefExists(name string, ldap_group_mapruleref *models.LDAPGroupMapruleref) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LDAP Group Map rule ref %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LDAP Group Map rule ref dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ldap_group_maprulerefFound := models.LDAPGroupMaprulerefFromContainer(cont)
		if ldap_group_maprulerefFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LDAP Group Map rule ref %s not found", rs.Primary.ID)
		}
		*ldap_group_mapruleref = *ldap_group_maprulerefFound
		return nil
	}
}

func testAccCheckAciLDAPGroupMaprulerefDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ldap_group_map_rule_to_group_map" {
			cont, err := client.Get(rs.Primary.ID)
			ldap_group_mapruleref := models.LDAPGroupMaprulerefFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LDAP Group Map rule ref %s Still exists", ldap_group_mapruleref.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLDAPGroupMaprulerefAttributes(description string, ldap_group_mapruleref *models.LDAPGroupMapruleref) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != GetMOName(ldap_group_mapruleref.DistinguishedName) {
			return fmt.Errorf("Bad aaa_ldap_group_map_rule_ref %s", GetMOName(ldap_group_mapruleref.DistinguishedName))
		}

		if description != ldap_group_mapruleref.Description {
			return fmt.Errorf("Bad ldap_group_mapruleref Description %s", ldap_group_mapruleref.Description)
		}

		if "test_name_alias_value" != ldap_group_mapruleref.NameAlias {
			return fmt.Errorf("Bad ldap_group_mapruleref NameAlias %s", ldap_group_mapruleref.NameAlias)
		}

		if "test_annotation_value" != ldap_group_mapruleref.Annotation {
			return fmt.Errorf("Bad ldap_group_mapruleref NameAlias %s", ldap_group_mapruleref.Annotation)
		}
		return nil
	}
}
