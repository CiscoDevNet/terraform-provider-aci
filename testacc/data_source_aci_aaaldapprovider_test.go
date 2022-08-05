package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLDAPProviderDataSource_Basic(t *testing.T) {
	resourceName := "aci_ldap_provider.test"
	dataSourceName := "data.aci_ldap_provider.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLDAPProviderDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLDAPProviderDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLDAPProviderDSWithoutRequired(rName, "type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLDAPProviderConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ssl_validation_level", resourceName, "ssl_validation_level"),
					resource.TestCheckResourceAttrPair(dataSourceName, "attribute", resourceName, "attribute"),
					resource.TestCheckResourceAttrPair(dataSourceName, "basedn", resourceName, "basedn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enable_ssl", resourceName, "enable_ssl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "filter", resourceName, "filter"),
					resource.TestCheckResourceAttrPair(dataSourceName, "monitor_server", resourceName, "monitor_server"),
					resource.TestCheckResourceAttrPair(dataSourceName, "monitoring_user", resourceName, "monitoring_user"),
					resource.TestCheckResourceAttrPair(dataSourceName, "port", resourceName, "port"),
					resource.TestCheckResourceAttrPair(dataSourceName, "retries", resourceName, "retries"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rootdn", resourceName, "rootdn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "timeout", resourceName, "timeout"),
				),
			},
			{
				Config:      CreateAccLDAPProviderDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccLDAPProviderDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccLDAPProviderDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLDAPProviderConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing ldap_provider Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
		type = "duo"
	}

	data "aci_ldap_provider" "test" {
	
		name  = aci_ldap_provider.test.name
		type = aci_ldap_provider.test.type
		depends_on = [ aci_ldap_provider.test ]
	}
	`, rName)
	return resource
}

func CreateLDAPProviderDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ldap_provider Data Source without ", attrName)
	rBlock := `
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
		type = "duo"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_ldap_provider" "test" {
	
	#	name  = aci_ldap_provider.test.name
		type = aci_ldap_provider.test.type
		depends_on = [ aci_ldap_provider.test ]
	}
		`
	case "type":
		rBlock += `
	data "aci_ldap_provider" "test" {
	
		name  = aci_ldap_provider.test.name
	#	type = aci_ldap_provider.test.type
		depends_on = [ aci_ldap_provider.test ]
	}
	`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLDAPProviderDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing ldap_provider Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
		type = "duo"
	}

	data "aci_ldap_provider" "test" {
	
		name  = "${aci_ldap_provider.test.name}_invalid"
		type = aci_ldap_provider.test.type
		depends_on = [ aci_ldap_provider.test ]
	}
	`, rName)
	return resource
}

func CreateAccLDAPProviderDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing ldap_provider Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
		type = "duo"
	}

	data "aci_ldap_provider" "test" {
	
		name  = aci_ldap_provider.test.name
		type = aci_ldap_provider.test.type
		%s = "%s"
		depends_on = [ aci_ldap_provider.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccLDAPProviderDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing ldap_provider Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_ldap_provider" "test" {
	
		name  = "%s"
		type = "duo"
		%s = "%s"
	}

	data "aci_ldap_provider" "test" {
	
		name  = aci_ldap_provider.test.name
		type = aci_ldap_provider.test.type
		depends_on = [ aci_ldap_provider.test ]
	}
	`, rName, key, value)
	return resource
}
