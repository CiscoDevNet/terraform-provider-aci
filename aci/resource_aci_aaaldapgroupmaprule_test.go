package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLDAPGroupMapRule_Basic(t *testing.T) {
	var ldap_group_map_rule models.LDAPGroupMapRule
	description := "ldap_group_map_rule created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLDAPGroupMapRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLDAPGroupMapRuleConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapRuleExists("aci_ldap_group_map_rule.fooldap_group_map_rule", &ldap_group_map_rule),
					testAccCheckAciLDAPGroupMapRuleAttributes(description, &ldap_group_map_rule),
				),
			},
		},
	})
}

func TestAccAciLDAPGroupMapRule_Update(t *testing.T) {
	var ldap_group_map_rule models.LDAPGroupMapRule
	description := "ldap_group_map_rule created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLDAPGroupMapRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLDAPGroupMapRuleConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapRuleExists("aci_ldap_group_map_rule.fooldap_group_map_rule", &ldap_group_map_rule),
					testAccCheckAciLDAPGroupMapRuleAttributes(description, &ldap_group_map_rule),
				),
			},
			{
				Config: testAccCheckAciLDAPGroupMapRuleConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPGroupMapRuleExists("aci_ldap_group_map_rule.fooldap_group_map_rule", &ldap_group_map_rule),
					testAccCheckAciLDAPGroupMapRuleAttributes(description, &ldap_group_map_rule),
				),
			},
		},
	})
}

func testAccCheckAciLDAPGroupMapRuleConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_ldap_group_map_rule" "fooldap_group_map_rule" {
		name 		= "test"
		description = "%s"
		name_alias  = "ldap_group_map_rule_alias"
		annotation  = "example"
		groupdn     = "example"
		type        = "duo"
	}

	`, description)
}

func testAccCheckAciLDAPGroupMapRuleExists(name string, ldap_group_map_rule *models.LDAPGroupMapRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LDAP Group Map Rule %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LDAP Group Map Rule dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ldap_group_map_ruleFound := models.LDAPGroupMapRuleFromContainer(cont)
		if ldap_group_map_ruleFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LDAP Group Map Rule %s not found", rs.Primary.ID)
		}
		*ldap_group_map_rule = *ldap_group_map_ruleFound
		return nil
	}
}

func testAccCheckAciLDAPGroupMapRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ldap_group_map_rule" {
			cont, err := client.Get(rs.Primary.ID)
			ldap_group_map_rule := models.LDAPGroupMapRuleFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LDAP Group Map Rule %s Still exists", ldap_group_map_rule.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLDAPGroupMapRuleAttributes(description string, ldap_group_map_rule *models.LDAPGroupMapRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != ldap_group_map_rule.Name {
			return fmt.Errorf("Bad aaa_ldap_group_map_rule %s", ldap_group_map_rule.Name)
		}

		if description != ldap_group_map_rule.Description {
			return fmt.Errorf("Bad ldap_group_map_rule Description %s", ldap_group_map_rule.Description)
		}

		if "ldap_group_map_rule_alias" != ldap_group_map_rule.NameAlias {
			return fmt.Errorf("Bad ldap_group_map_rule Name Alias %s", ldap_group_map_rule.NameAlias)
		}

		if "example" != ldap_group_map_rule.Annotation {
			return fmt.Errorf("Bad ldap_group_map_rule Annotation %s", ldap_group_map_rule.Annotation)
		}

		if "example" != ldap_group_map_rule.Groupdn {
			return fmt.Errorf("Bad ldap_group_map_rule Group DN %s", ldap_group_map_rule.Groupdn)
		}
		return nil
	}
}
