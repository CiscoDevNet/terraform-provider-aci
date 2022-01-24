package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciContractDataSource_Basic(t *testing.T) {
	resourceName := "aci_contract.test"
	dataSourceName := "data.aci_contract.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciContractDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateAccContractDSWithoutTenant(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccContractDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccContractDSConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "scope", resourceName, "scope"),
					resource.TestCheckResourceAttrPair(dataSourceName, "target_dscp", resourceName, "target_dscp"),
				),
			},
			{
				Config:      CreateAccContractUpdatedConfigDataSourceRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccContractUpdatedConfigDataSource(rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
			{
				Config:      CreateAccContractDataSourceWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
		},
	})
}

func CreateAccContractDSWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract data source creation without creating tenant")
	resource := fmt.Sprintf(`
	data "aci_contract" "test" {
		name = "%s"
	}
	`, rName)
	return resource
}

func CreateAccContractDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract data source creation without name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_contract" "test" {
		tenant_dn = aci_contract.test.tenant_dn
	}
	`, rName, rName)
	return resource
}

func CreateAccContractDSConfig(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract data source creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_contract" "test" {
		tenant_dn = aci_contract.test.tenant_dn
		name = aci_contract.test.name
	}
	`, rName, rName)
	return resource
}

func CreateAccContractUpdatedConfigDataSourceRandomAttr(rName, key, value string) string {
	fmt.Println("=== STEP  Basic: testing contract data source creation with random attributes")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_contract" "test" {
		tenant_dn = aci_contract.test.tenant_dn
		name = aci_contract.test.name
		%s = "%s"
	}
	`, rName, rName, key, value)
	return resource
}

func CreateAccContractUpdatedConfigDataSource(rName, key, value string) string {
	fmt.Println("=== STEP  Basic: testing contract data source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}

	data "aci_contract" "test" {
		tenant_dn = aci_contract.test.tenant_dn
		name = aci_contract.test.name
	}
	`, rName, rName, key, value)
	return resource
}

func CreateAccContractDataSourceWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract data source with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_contract" "test" {
		tenant_dn = aci_contract.test.tenant_dn
		name = "${aci_contract.test.name}abc"
	}
	`, rName, rName)
	return resource
}
