package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLDAPGroupMaprulerefDataSource_Basic(t *testing.T) {
	resourceName := "aci_ldap_group_map_rule_to_group_map.test"
	dataSourceName := "data.aci_ldap_group_map_rule_to_group_map.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPGroupMaprulerefDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateLDAPGroupMaprulerefDSWithoutRequired(rName, rName, "ldap_group_map_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLDAPGroupMaprulerefDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLDAPGroupMaprulerefConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "ldap_group_map_dn", resourceName, "ldap_group_map_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccLDAPGroupMaprulerefDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccLDAPGroupMaprulerefDSWithInvalidName(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccLDAPGroupMaprulerefDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLDAPGroupMaprulerefConfigDataSource(aaaLdapGroupMapName, rName string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule_to_group_map Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
		type = "duo"
	}
	
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "%s"
	}

	data "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = aci_ldap_group_map_rule_to_group_map.test.name
		depends_on = [ aci_ldap_group_map_rule_to_group_map.test ]
	}
	`, aaaLdapGroupMapName, rName)
	return resource
}

func CreateLDAPGroupMaprulerefDSWithoutRequired(aaaLdapGroupMapName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ldap_group_map_rule_to_group_map Data Source without ", attrName)
	rBlock := `
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
		type = "duo"
	}
	
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "ldap_group_map_dn":
		rBlock += `
	data "aci_ldap_group_map_rule_to_group_map" "test" {
	#	ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = aci_ldap_group_map_rule_to_group_map.test.name
		depends_on = [ aci_ldap_group_map_rule_to_group_map.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
	#	name  = aci_ldap_group_map_rule_to_group_map.test.name
		depends_on = [ aci_ldap_group_map_rule_to_group_map.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaLdapGroupMapName, rName)
}

func CreateAccLDAPGroupMaprulerefDSWithInvalidName(aaaLdapGroupMapName, rName string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule_to_group_map Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
		type = "duo"
	}
	
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "%s"
	}

	data "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "${aci_ldap_group_map_rule_to_group_map.test.name}_invalid"
		depends_on = [ aci_ldap_group_map_rule_to_group_map.test ]
	}
	`, aaaLdapGroupMapName, rName)
	return resource
}

func CreateAccLDAPGroupMaprulerefDataSourceUpdate(aaaLdapGroupMapName, rName, key, value string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule_to_group_map Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
		type = "duo"
	}
	
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "%s"
	}

	data "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = aci_ldap_group_map_rule_to_group_map.test.name
		%s = "%s"
		depends_on = [ aci_ldap_group_map_rule_to_group_map.test ]
	}
	`, aaaLdapGroupMapName, rName, key, value)
	return resource
}

func CreateAccLDAPGroupMaprulerefDataSourceUpdatedResource(aaaLdapGroupMapName, rName, key, value string) string {
	fmt.Println("=== STEP  testing ldap_group_map_rule_to_group_map Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
		name 		= "%s"
		type = "duo"
	}
	
	resource "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_ldap_group_map_rule_to_group_map" "test" {
		ldap_group_map_dn  = aci_ldap_group_map.test.id
		name  = aci_ldap_group_map_rule_to_group_map.test.name
		depends_on = [ aci_ldap_group_map_rule_to_group_map.test ]
	}
	`, aaaLdapGroupMapName, rName, key, value)
	return resource
}
