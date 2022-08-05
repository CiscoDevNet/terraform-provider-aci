package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLDAPGroupMapDataSource_Basic(t *testing.T) {
	resourceName := "aci_ldap_group_map.test"
	dataSourceName := "data.aci_ldap_group_map.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPGroupMapDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLDAPGroupMapDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLDAPGroupMapDSWithoutRequired(rName, "type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLDAPGroupMapConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
				),
			},
			{
				Config:      CreateAccLDAPGroupMapDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccLDAPGroupMapDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccLDAPGroupMapDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLDAPGroupMapConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing ldap_group_map Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
	
		name  = "%s"
		type = "duo"
	}

	data "aci_ldap_group_map" "test" {
	
		name  = aci_ldap_group_map.test.name
		type = aci_ldap_group_map.test.type
		depends_on = [ aci_ldap_group_map.test ]
	}
	`, rName)
	return resource
}

func CreateLDAPGroupMapDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ldap_group_map Data Source without ", attrName)
	rBlock := `
	
	resource "aci_ldap_group_map" "test" {
	
		name  = "%s"
		type = "duo"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_ldap_group_map" "test" {
	
	#	name  = aci_ldap_group_map.test.name
		type = aci_ldap_group_map.test.type
		depends_on = [ aci_ldap_group_map.test ]
	}
		`
	case "type":
		rBlock += `
	data "aci_ldap_group_map" "test" {
	
		name  = aci_ldap_group_map.test.name
	#	type = aci_ldap_group_map.test.type
		depends_on = [ aci_ldap_group_map.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLDAPGroupMapDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing ldap_group_map Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
	
		name  = "%s"
		type = "duo"
	}

	data "aci_ldap_group_map" "test" {
	
		name  = "${aci_ldap_group_map.test.name}_invalid"
		type = aci_ldap_group_map.test.type
		depends_on = [ aci_ldap_group_map.test ]
	}
	`, rName)
	return resource
}

func CreateAccLDAPGroupMapDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing ldap_group_map Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
	
		name  = "%s"
		type = "duo"
	}

	data "aci_ldap_group_map" "test" {
	
		name  = aci_ldap_group_map.test.name
		type = aci_ldap_group_map.test.type
		%s = "%s"
		depends_on = [ aci_ldap_group_map.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccLDAPGroupMapDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing ldap_group_map Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_group_map" "test" {
	
		name  = "%s"
		type = "duo"
		%s = "%s"
	}

	data "aci_ldap_group_map" "test" {
	
		name  = aci_ldap_group_map.test.name
		type = aci_ldap_group_map.test.type
		depends_on = [ aci_ldap_group_map.test ]
	}
	`, rName, key, value)
	return resource
}
