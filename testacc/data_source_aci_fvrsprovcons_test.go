package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciEPGToContractDataSource_Basic(t *testing.T) {
	resourceName := "aci_epg_to_contract.test"
	dataSourceName := "data.aci_epg_to_contract.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEPGToContractDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccEPGToContractDSWithoutApplicationEPG(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccEPGToContractDSWithoutContract(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccEPGToContractDSWithoutContractType(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEPGToContractConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "contract_dn", resourceName, "contract_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "application_epg_dn", resourceName, "application_epg_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "contract_type", resourceName, "contract_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_t", resourceName, "match_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
				),
			},
			{
				Config:      CreateAccEPGToContractUpdatedConfigDataSourceRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEPGToContractUpdatedConfigDataSource(rName, "annotation", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccEPGToContractUpdatedConfigDataSourceRandomAttr(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing EPG to Contract data source with Random Attribute")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}
	
	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_epg_to_contract" "test" {
		contract_dn = aci_contract.test.id
		application_epg_dn = aci_application_epg.test.id
		contract_type = "provider"
	}

	data "aci_epg_to_contract" "test" {
		contract_dn = aci_epg_to_contract.test.contract_dn
		application_epg_dn = aci_epg_to_contract.test.application_epg_dn
		contract_type = aci_epg_to_contract.test.contract_type
		%s = "%s"
	}
	`, rName, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccEPGToContractUpdatedConfigDataSource(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing EPG to Contract data source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}
	
	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_epg_to_contract" "test" {
		contract_dn = aci_contract.test.id
		application_epg_dn = aci_application_epg.test.id
		contract_type = "provider"
		%s = "%s"
	}

	data "aci_epg_to_contract" "test" {
		contract_dn = aci_epg_to_contract.test.contract_dn
		application_epg_dn = aci_epg_to_contract.test.application_epg_dn
		contract_type = aci_epg_to_contract.test.contract_type
	}
	`, rName, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccEPGToContractConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: Testing EPG to Contract data source")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}
	
	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_epg_to_contract" "test" {
		contract_dn = aci_contract.test.id
		application_epg_dn = aci_application_epg.test.id
		contract_type = "provider"
	}

	data "aci_epg_to_contract" "test" {
		contract_dn = aci_epg_to_contract.test.contract_dn
		application_epg_dn = aci_epg_to_contract.test.application_epg_dn
		contract_type = aci_epg_to_contract.test.contract_type
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccEPGToContractDSWithoutApplicationEPG(rName string) string {
	fmt.Println("=== STEP  Basic: Testing EPG to Contract data source without Application EPG")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}
	
	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_epg_to_contract" "test" {
		contract_dn = aci_contract.test.id
		application_epg_dn = aci_application_epg.test.id
		contract_type = "provider"
	}

	data "aci_epg_to_contract" "test" {
		contract_dn = aci_epg_to_contract.test.contract_dn
		contract_type = aci_epg_to_contract.test.contract_type
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccEPGToContractDSWithInvalidApplicationEPG(rName string) string {
	fmt.Println("=== STEP  Basic: Testing EPG to Contract data source with Invalid Application EPG")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}
	
	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_epg_to_contract" "test" {
		contract_dn = aci_contract.test.id
		application_epg_dn = aci_application_epg.test.id
		contract_type = "provider"
	}

	data "aci_epg_to_contract" "test" {
		contract_dn = aci_epg_to_contract.test.contract_dn
		contract_type = aci_epg_to_contract.test.contract_type
		application_epg_dn = "${aci_epg_to_contract.test.application_epg_dn}abc"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccEPGToContractDSWithoutContract(rName string) string {
	fmt.Println("=== STEP  Basic: Testing EPG to Contract data source without Contract")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}
	
	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_epg_to_contract" "test" {
		contract_dn = aci_contract.test.id
		application_epg_dn = aci_application_epg.test.id
		contract_type = "provider"
	}

	data "aci_epg_to_contract" "test" {
		application_epg_dn = aci_epg_to_contract.test.application_epg_dn
		contract_type = aci_epg_to_contract.test.contract_type
	}
	`, rName, rName, rName, rName)
	return resource
}
func CreateAccEPGToContractDSWithoutContractType(rName string) string {
	fmt.Println("=== STEP  Basic: Testing EPG to Contract data source without Contract Type")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}
	
	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_epg_to_contract" "test" {
		contract_dn = aci_contract.test.id
		application_epg_dn = aci_application_epg.test.id
		contract_type = "provider"
	}

	data "aci_epg_to_contract" "test" {
		contract_dn = aci_epg_to_contract.test.contract_dn
		application_epg_dn = aci_epg_to_contract.test.application_epg_dn
	}
	`, rName, rName, rName, rName)
	return resource
}
