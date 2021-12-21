package acctest

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciEPGToContract_Basic(t *testing.T) {
	var epg_to_contract_default_cons models.ContractConsumer
	var epg_to_contract_updated_cons models.ContractConsumer
	var epg_to_contract_default_prov models.ContractProvider
	// var epg_to_contract_updated_prov models.ContractProvider
	resourceName := "aci_epg_to_contract.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOtherName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEPGToContractDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccEPGToContractWithoutApplicationEPG(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccEPGToContractWithoutContract(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccEPGToContractWithoutContractType(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEPGToContractConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_default_cons),
					//resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "application_epg_dn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "contract_dn", fmt.Sprintf("uni/tn-%s/brc-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "contract_type", "consumer"),
					//resource.TestCheckResourceAttr(resourceName, "match_t", ""),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
				),
			},
			{
				Config: CreateAccEPGToContractConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_updated_cons),
					resource.TestCheckResourceAttr(resourceName, "annotation", "test_annotation"),
					resource.TestCheckResourceAttr(resourceName, "application_epg_dn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "contract_dn", fmt.Sprintf("uni/tn-%s/brc-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "contract_type", "consumer"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtleastOne"),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					testAccCheckAciEPGToContractIdEqual(&epg_to_contract_default_cons, &epg_to_contract_updated_cons),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"contract_dn", "match_t"},
			},
			{
				Config: CreateAccEPGToContractConfigWithUpdatedContract(rName, rOtherName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_updated_cons),
					resource.TestCheckResourceAttr(resourceName, "contract_dn", fmt.Sprintf("uni/tn-%s/brc-%s", rName, rOtherName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rscons-%s", rName, rName, rName, rOtherName)),
					testAccCheckAciEPGToContractIdNotEqual(&epg_to_contract_default_cons, &epg_to_contract_updated_cons),
				),
			},
			{
				Config: CreateAccEPGToContractConfig(rName),
			},
			{
				Config: CreateAccEPGToContractConfigWithUpdatedApplicationEPG(rName, rOtherName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_updated_cons),
					resource.TestCheckResourceAttr(resourceName, "application_epg_dn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rOtherName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rscons-%s", rName, rName, rOtherName, rName)),
					testAccCheckAciEPGToContractIdNotEqual(&epg_to_contract_default_cons, &epg_to_contract_updated_cons),
				),
			},
			{
				Config: CreateAccEPGToContractConfig(rName),
			},
			{
				Config: CreateAccEPGToContractConfigWithUpdatedContractType(rName, "provider"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractProviderExists(resourceName, &epg_to_contract_default_prov),
					resource.TestCheckResourceAttr(resourceName, "contract_type", "provider"),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rsprov-%s", rName, rName, rName, rName)),
				),
			},
		},
	})
}

func TestAccAciEPGToContract_Update(t *testing.T) {
	var epg_to_contract_default_cons models.ContractConsumer
	var epg_to_contract_updated_cons models.ContractConsumer
	resourceName := "aci_epg_to_contract.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEPGToContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEPGToContractConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_default_cons),
				),
			},
			{
				Config: CreateAccEPGToContractUpdatedAttr(rName, "match_t", "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_updated_cons),
					resource.TestCheckResourceAttr(resourceName, "match_t", "All"),
					testAccCheckAciEPGToContractIdEqual(&epg_to_contract_default_cons, &epg_to_contract_updated_cons),
				),
			},
			{
				Config: CreateAccEPGToContractUpdatedAttr(rName, "match_t", "AtmostOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_updated_cons),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtmostOne"),
					testAccCheckAciEPGToContractIdEqual(&epg_to_contract_default_cons, &epg_to_contract_updated_cons),
				),
			},
			{
				Config: CreateAccEPGToContractUpdatedAttr(rName, "match_t", "None"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_updated_cons),
					resource.TestCheckResourceAttr(resourceName, "match_t", "None"),
					testAccCheckAciEPGToContractIdEqual(&epg_to_contract_default_cons, &epg_to_contract_updated_cons),
				),
			},
			{
				Config: CreateAccEPGToContractUpdatedAttr(rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_updated_cons),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciEPGToContractIdEqual(&epg_to_contract_default_cons, &epg_to_contract_updated_cons),
				),
			},
			{
				Config: CreateAccEPGToContractUpdatedAttr(rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_updated_cons),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciEPGToContractIdEqual(&epg_to_contract_default_cons, &epg_to_contract_updated_cons),
				),
			},
			{
				Config: CreateAccEPGToContractUpdatedAttr(rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_updated_cons),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciEPGToContractIdEqual(&epg_to_contract_default_cons, &epg_to_contract_updated_cons),
				),
			},
			{
				Config: CreateAccEPGToContractUpdatedAttr(rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_updated_cons),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciEPGToContractIdEqual(&epg_to_contract_default_cons, &epg_to_contract_updated_cons),
				),
			},
			{
				Config: CreateAccEPGToContractUpdatedAttr(rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGToContractExists(resourceName, &epg_to_contract_updated_cons),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciEPGToContractIdEqual(&epg_to_contract_default_cons, &epg_to_contract_updated_cons),
				),
			},
		},
	})

}

func TestAccAciEPGToContract_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longAnnotationDesc := acctest.RandString(129)
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEPGToContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEPGToContractConfig(rName),
			},
			{
				Config:      CreateAccEPGToContractWithInvalidApplicationEPG(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+ name dn, class fvRsCons (.)+`),
			},
			{
				Config:      CreateAccEPGToContractWithInvalidContractType(rName, randomValue),
				ExpectError: regexp.MustCompile(`Contract Type: Value must be from (.)+`),
			},
			{
				Config:      CreateAccEPGToContractUpdatedAttr(rName, "annotation", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEPGToContractUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccEPGToContractUpdatedAttr(rName, "match_t", randomValue),
				ExpectError: regexp.MustCompile(`expected match_t to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccEPGToContractUpdatedAttr(rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`expected prio to be one of (.)+, got (.)+`),
			},
			{
				Config: CreateAccEPGToContractConfig(rName),
			},
		},
	})
}

func TestAccAciEPGToContracts_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEPGToContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEPGToContractsConfig(rName),
			},
		},
	})
}

func CreateAccEPGToContractWithoutApplicationEPG(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg_to_contract without creating application_epg")
	resource := fmt.Sprint(`
	resource "aci_epg_to_contract" "test" {
		contract_type = "consumer"
	}
	`)
	return resource

}
func CreateAccEPGToContractWithoutContract(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg_to_contract without creating contract")
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
		contract_type = "consumer"
		application_epg_dn = aci_application_epg.test.id
	}
	`, rName, rName, rName, rName)
	return resource

}

func CreateAccEPGToContractWithoutContractType(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg_to_contract without passing contract type attribute")
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
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccEPGToContractConfig(rName string) string {
	fmt.Println("=== STEP  Basic: testing EPG to Contract creation with required paramters only")
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
		contract_type = "consumer"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccEPGToContractConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing EPG to Contract creation with optional paramters")
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
		contract_type = "consumer"
		annotation = "test_annotation"
		match_t = "AtleastOne"
		prio = "level1"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccEPGToContractConfigWithUpdatedContract(rName, rOtherName string) string {
	fmt.Println("=== STEP  Basic: testing EPG to Contract creation with Updated Contract")
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
		contract_type = "consumer"
	}
	`, rName, rName, rName, rOtherName)
	return resource
}

func CreateAccEPGToContractConfigWithUpdatedApplicationEPG(rName, rOtherName string) string {
	fmt.Println("=== STEP  Basic: testing EPG to Contract creation with Updated Contract")
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
		contract_type = "consumer"
	}
	`, rName, rName, rOtherName, rName)
	return resource
}

func CreateAccEPGToContractConfigWithUpdatedContractType(rName, contractType string) string {
	fmt.Println("=== STEP  Basic: testing EPG to Contract creation with Updated Contract Type")
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
		contract_type = "%s"
	}
	`, rName, rName, rName, rName, contractType)
	return resource
}

func CreateAccEPGToContractUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing EPG to Contract %s = %s\n", attribute, value)
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
		contract_type = "consumer"
		%s = "%s"
	}
	`, rName, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccEPGToContractWithInvalidContractType(rName, randomValue string) string {
	fmt.Println("=== STEP  Basic: testing EPG to Contract updation with Invalid application_epg_dn")
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
		application_epg_dn = aci_application_profile.test.id
		contract_type = "%s"
	}
	`, rName, rName, rName, rName, randomValue)
	return resource
}
func CreateAccEPGToContractWithInvalidApplicationEPG(rName string) string {
	fmt.Println("=== STEP  Basic: testing EPG to Contract updation with Invalid application_epg_dn")
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
		application_epg_dn = aci_application_profile.test.id
		contract_type = "consumer"
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccEPGToContractWithInvalidContract(rName string) string {
	fmt.Println("=== STEP  Basic: testing EPG to Contract updation with Invalid contract_dn")
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
		contract_dn = aci_tenant.test.id
		application_epg_dn = aci_application_epg.test.id
		contract_type = "consumer"
	}
	`, rName, rName, rName, rName)
	return resource
}
func CreateAccEPGToContractsConfig(rName string) string {
	fmt.Println("=== STEP  creating multiple epg_to_contracts")
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
	
	resource "aci_application_epg" "test1" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract" "test1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_epg_to_contract" "test" {
		contract_dn = aci_contract.test.id
		application_epg_dn = aci_application_epg.test.id
		contract_type = "consumer"
	}

	resource "aci_epg_to_contract" "test1" {
		contract_dn = aci_contract.test.id
		application_epg_dn = aci_application_epg.test.id
		contract_type = "provider"
	}

	resource "aci_epg_to_contract" "test2" {
		contract_dn = aci_contract.test1.id
		application_epg_dn = aci_application_epg.test.id
		contract_type = "consumer"
	}

	resource "aci_epg_to_contract" "test3" {
		contract_dn = aci_contract.test.id
		application_epg_dn = aci_application_epg.test1.id
		contract_type = "consumer"
	}
	`, rName, rName, rName, rName, rName, rName)
	return resource

}

func testAccCheckAciEPGToContractDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing EPG To Contract destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_epg_to_contract" {
			cont, err := client.Get(rs.Primary.ID)
			epg_to_contract := models.ContractConsumerFromContainer(cont)
			if err == nil {
				return fmt.Errorf("EPG To Contract %s still exists", epg_to_contract.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciEPGToContractProviderExists(name string, epg_to_contract *models.ContractProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("EPG To Contract %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EPG to Contract Dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		epg_to_contractFound := models.ContractProviderFromContainer(cont)
		if epg_to_contractFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("EPG to Contract %s not found", rs.Primary.ID)
		}
		*epg_to_contract = *epg_to_contractFound
		return nil
	}
}

func testAccCheckAciEPGToContractExists(name string, epg_to_contract *models.ContractConsumer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("EPG To Contract %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Filter Entry Dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		epg_to_contractFound := models.ContractConsumerFromContainer(cont)
		if epg_to_contractFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("EPG to Contract %s not found", rs.Primary.ID)
		}
		*epg_to_contract = *epg_to_contractFound
		return nil
	}
}

func testAccCheckAciEPGToContractIdNotEqual(cc1, cc2 *models.ContractConsumer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cc1.DistinguishedName == cc2.DistinguishedName {
			return fmt.Errorf("Filter Entry DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciEPGToContractIdEqual(cc1, cc2 *models.ContractConsumer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cc1.DistinguishedName != cc2.DistinguishedName {
			return fmt.Errorf("EPG to Contract DNs are no equal")
		}
		return nil
	}
}
