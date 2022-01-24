package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciDuoProviderGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_duo_provider_group.test"
	dataSourceName := "data.aci_duo_provider_group.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:	  func(){ testAccPreCheck(t) },
		ProviderFactories:    testAccProviders,
		CheckDestroy: testAccCheckAciDuoProviderGroupDestroy,
		Steps: []resource.TestStep{
			
			{
				Config:      CreateDuoProviderGroupDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDuoProviderGroupConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_choice", resourceName, "auth_choice"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ldap_group_map_ref", resourceName, "ldap_group_map_ref"),
					resource.TestCheckResourceAttrPair(dataSourceName, "provider_type", resourceName, "provider_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "sec_fac_auth_methods.#", resourceName, "sec_fac_auth_methods.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "sec_fac_auth_methods.0", resourceName, "sec_fac_auth_methods.0"),
					
				),
			},
			{
				Config:      CreateAccDuoProviderGroupDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			
			{
				Config:      CreateAccDuoProviderGroupDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccDuoProviderGroupDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}


func CreateAccDuoProviderGroupConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing duo_provider_group Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
	
		name  = "%s"
	}

	data "aci_duo_provider_group" "test" {
	
		name  = aci_duo_provider_group.test.name
		depends_on = [ aci_duo_provider_group.test ]
	}
	`, rName)
	return resource
}

func CreateDuoProviderGroupDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing duo_provider_group Data Source without ",attrName)
	rBlock := `
	
	resource "aci_duo_provider_group" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_duo_provider_group" "test" {
	
	#	name  = aci_duo_provider_group.test.name
		depends_on = [ aci_duo_provider_group.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock,rName)
}


func CreateAccDuoProviderGroupDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing duo_provider_group Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
	
		name  = "%s"
	}

	data "aci_duo_provider_group" "test" {
	
		name  = "${aci_duo_provider_group.test.name}_invalid"
		depends_on = [ aci_duo_provider_group.test ]
	}
	`, rName)
	return resource
}

func CreateAccDuoProviderGroupDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing duo_provider_group Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
	
		name  = "%s"
	}

	data "aci_duo_provider_group" "test" {
	
		name  = aci_duo_provider_group.test.name
		%s = "%s"
		depends_on = [ aci_duo_provider_group.test ]
	}
	`, rName,key,value)
	return resource
}

func CreateAccDuoProviderGroupDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing duo_provider_group Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_duo_provider_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_duo_provider_group" "test" {
	
		name  = aci_duo_provider_group.test.name
		depends_on = [ aci_duo_provider_group.test ]
	}
	`, rName,key,value)
	return resource
}