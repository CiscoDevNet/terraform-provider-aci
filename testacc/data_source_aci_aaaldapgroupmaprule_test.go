package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLDAPGroupMapRuleDataSource_Basic(t *testing.T) {
	resourceName := "aci_ldap_group_map_rule.test"
	dataSourceName := "data.aci_ldap_group_map_rule.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPGroupMapRuleDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLDAPGroupMapRuleDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLDAPGroupMapRuleDSWithoutRequired(rName, "type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLDAPGroupMapRuleConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "groupdn", resourceName, "groupdn"),
				),
			},
			{
				Config:      CreateAccLDAPGroupMapRuleDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccLDAPGroupMapRuleDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccLDAPGroupMapRuleDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLDAPGroupMapRuleConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map_rule" "test" {
		type = "duo"
		name  = "%s"
	}

	data "aci_ldap_group_map_rule" "test" {
		type = aci_ldap_group_map_rule.test.type
		name  = aci_ldap_group_map_rule.test.name
		depends_on = [ aci_ldap_group_map_rule.test ]
	}
	`, rName)
	return resource
}

func CreateLDAPGroupMapRuleDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ldap_group_map_rule Data Source without ", attrName)
	rBlock := `
	
	resource "aci_ldap_group_map_rule" "test" {
		type = "duo"
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_ldap_group_map_rule" "test" {
		type = aci_ldap_group_map_rule.test.type
	#	name  = aci_ldap_group_map_rule.test.name
		depends_on = [ aci_ldap_group_map_rule.test ]
	}
		`
	case "type":
		rBlock += `
	data "aci_ldap_group_map_rule" "test" {
	#	type = aci_ldap_group_map_rule.test.type
		name  = aci_ldap_group_map_rule.test.name
		depends_on = [ aci_ldap_group_map_rule.test ]
	}
	`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLDAPGroupMapRuleDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map_rule" "test" {
		type = "duo"
		name  = "%s"
	}

	data "aci_ldap_group_map_rule" "test" {
		type = aci_ldap_group_map_rule.test.type
		name  = "${aci_ldap_group_map_rule.test.name}_invalid"
		depends_on = [ aci_ldap_group_map_rule.test ]
	}
	`, rName)
	return resource
}

func CreateAccLDAPGroupMapRuleDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map_rule" "test" {
		type = "duo"
		name  = "%s"
	}

	data "aci_ldap_group_map_rule" "test" {
		type = aci_ldap_group_map_rule.test.type
		name  = aci_ldap_group_map_rule.test.name
		%s = "%s"
		depends_on = [ aci_ldap_group_map_rule.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccLDAPGroupMapRuleDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map_rule" "test" {
		type = "duo"
		name  = "%s"
		%s = "%s"
	}

	data "aci_ldap_group_map_rule" "test" {
		type = aci_ldap_group_map_rule.test.type
		name  = aci_ldap_group_map_rule.test.name
		depends_on = [ aci_ldap_group_map_rule.test ]
	}
	`, rName, key, value)
	return resource
}
