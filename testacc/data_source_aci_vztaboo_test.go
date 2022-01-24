package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTabooContractDataSource_Basic(t *testing.T) {
	resourceName := "aci_taboo_contract.test"
	dataSourceName := "data.aci_taboo_contract.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTabooContractDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateTabooContractDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateTabooContractDSWithoutTenantdn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTabooContractConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccTabooContractDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccTabooContractDSWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccTabooContractDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccTabooContractConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing taboo_contract creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_taboo_contract.test.name
		depends_on = [
			aci_taboo_contract.test
		]
	}
	`, fvTenantName, rName)
	return resource
}
func CreateAccTabooContractDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing taboo_contract creation with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_taboo_contract" "test" {

		tenant_dn  = "${aci_tenant.test.id}_invalid"
		name  = aci_taboo_contract.test.name
		depends_on = [
			aci_taboo_contract.test
		]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccTabooContractDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing taboo_contract Updation with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_taboo_contract.test.name
		depends_on = [
			aci_taboo_contract.test
		]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccTabooContractDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing taboo_contract Updation with random attributes")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_taboo_contract.test.name
		%s = "%s"
		depends_on = [
			aci_taboo_contract.test
		]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateTabooContractDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing Taboo Contract data source reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
		depends_on = [
			aci_taboo_contract.test
		]
	}
	`, rName, rName)
	return resource
}

func CreateTabooContractDSWithoutTenantdn(rName string) string {
	fmt.Println("=== STEP  Basic: testing Taboo Contract data source reading without giving Tenant Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_taboo_contract" "test" {
		name  = aci_taboo_contract.test.name
		depends_on = [
			aci_taboo_contract.test
		]
	}
	`, rName, rName)
	return resource
}
