package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciContractSubjectDataSource_Basic(t *testing.T) {
	resourceName := "aci_contract_subject.test"
	dataSourceName := "data.aci_contract_subject.test"
	rName := acctest.RandString(5)
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciContractSubjectDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccContractSubjectDSWithoutContract(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccContractSubjectDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccContractSubjectConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cons_match_t", resourceName, "cons_match_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prov_match_t", resourceName, "prov_match_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rev_flt_ports", resourceName, "rev_flt_ports"),
					resource.TestCheckResourceAttrPair(dataSourceName, "target_dscp", resourceName, "target_dscp"),
					resource.TestCheckResourceAttrPair(dataSourceName, "contract_dn", resourceName, "contract_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
				),
			},
			{
				Config:      CreateAccContractSubjectUpdatedConfigDataSourceRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccContractSubjectUpdatedConfigDataSource(rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
			{
				Config:      CreateAccContractSubjectDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
		},
	})
}

func CreateAccContractSubjectUpdatedConfigDataSourceRandomAttr(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing contract subject data source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
	}

	data "aci_contract_subject" "test" {
		contract_dn = aci_contract.test.id
		name = aci_contract_subject.test.name
		%s = "%s"
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccContractSubjectUpdatedConfigDataSource(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing contract subject data source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
		%s = "%s"
	}

	data "aci_contract_subject" "test" {
		contract_dn = aci_contract.test.id
		name = aci_contract_subject.test.name
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccContractSubjectConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: Testing contract subject data source")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
		description = "test_description"
	}

	data "aci_contract_subject" "test" {
		contract_dn = aci_contract.test.id
		name = aci_contract_subject.test.name
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccContractSubjectDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract subject reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
	}

	data "aci_contract_subject" "test" {
		contract_dn = aci_contract.test.id
		name = "${aci_contract_subject.test.name}abc"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccContractSubjectDSWithoutContract(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract subject reading without giving contract_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
	}

	data "aci_contract_subject" "test" {
		name = "%s"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccContractSubjectDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract subject reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract_subject" "test"{
		contract_dn = aci_contract.test.id
		name = "%s"
	}

	data "aci_contract" "test" {
		contract_dn = aci_contract.test.id
	}
	`, rName, rName, rName)
	return resource
}
