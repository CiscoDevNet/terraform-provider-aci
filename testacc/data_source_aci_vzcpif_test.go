package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciImportedContractDataSource_Basic(t *testing.T) {
	resourceName := "aci_imported_contract.test"
	dataSourceName := "data.aci_imported_contract.test"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciImportedContractDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccImportedContractDSWithoutTenant(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccImportedContractDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccImportedContractConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
				),
			},
			{
				Config:      CreateAccImportedContractDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccImportedContractDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccImportedContractDataSourceUpdate(rName, "description", "description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccImportedContractDataSourceUpdateRandomAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing imported contract data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_imported_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	data "aci_imported_contract" "test" {
		name = aci_imported_contract.test.name
		tenant_dn = aci_imported_contract.test.tenant_dn
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccImportedContractDataSourceUpdate(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing imported contract data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_imported_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}
	data "aci_imported_contract" "test" {
		name = aci_imported_contract.test.name
		tenant_dn = aci_imported_contract.test.tenant_dn
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccImportedContractConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: testing imported contract creation for data source test")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_imported_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_imported_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = aci_imported_contract.test.name
	}
	`, rName, rName)
	return resource
}

func CreateAccImportedContractDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing imported contract reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_imported_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_imported_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "${aci_imported_contract.test.name}xyz"
	}
	`, rName, rName)
	return resource
}

func CreateAccImportedContractDSWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing imported contract reading without giving tenant_dn")
	resource := fmt.Sprintf(`
	data "aci_imported_contract" "test" {
		name = "%s"
	}
	`, rName)
	return resource
}

func CreateAccImportedContractDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing imported contract reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	data "aci_imported_contract" "test" {
		tenant_dn = aci_tenant.test.id
	}
	`, rName)
	return resource
}
