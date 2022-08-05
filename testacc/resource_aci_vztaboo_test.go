package testacc

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

func TestAccAciTabooContract_Basic(t *testing.T) {
	var taboo_contract_default models.TabooContract
	var taboo_contract_updated models.TabooContract
	resourceName := "aci_taboo_contract.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTabooContractDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateTabooContractWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateTabooContractWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTabooContractConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTabooContractExists(resourceName, &taboo_contract_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccTabooContractConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTabooContractExists(resourceName, &taboo_contract_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_taboo_contract"),
					testAccCheckAciTabooContractIdEqual(&taboo_contract_default, &taboo_contract_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccTabooContractConfigUpdatedName(rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config: CreateAccTabooContractConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTabooContractExists(resourceName, &taboo_contract_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciTabooContractIdNotEqual(&taboo_contract_default, &taboo_contract_updated),
				),
			},
			{
				Config: CreateAccTabooContractConfig(rName, rName),
			},
			{
				Config: CreateAccTabooContractConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTabooContractExists(resourceName, &taboo_contract_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciTabooContractIdNotEqual(&taboo_contract_default, &taboo_contract_updated),
				),
			},
		},
	})
}

func TestAccAciTabooContract_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTabooContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTabooContractConfig(rName, rName),
			},
			{
				Config:      CreateAccTabooContractWithInValidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccTabooContractUpdatedAttr(rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTabooContractUpdatedAttr(rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTabooContractUpdatedAttr(rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccTabooContractUpdatedAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+is not expected here.`),
			},
			{
				Config: CreateAccTabooContractConfig(rName, rName),
			},
		},
	})
}

func TestAccAciTabooContract_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTabooContractDestroy,
		Steps: []resource.TestStep{

			{
				Config: CreateAccTabooContractConfigs(rName),
			},
		},
	})
}

func testAccCheckAciTabooContractExists(name string, taboo_contract *models.TabooContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Taboo Contract %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Taboo Contract dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		taboo_contractFound := models.TabooContractFromContainer(cont)
		if taboo_contractFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Taboo Contract %s not found", rs.Primary.ID)
		}
		*taboo_contract = *taboo_contractFound
		return nil
	}
}

func testAccCheckAciTabooContractDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing taboo_contract destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_taboo_contract" {
			cont, err := client.Get(rs.Primary.ID)
			taboo_contract := models.TabooContractFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Taboo Contract %s Still exists", taboo_contract.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTabooContractIdEqual(m1, m2 *models.TabooContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("taboo_contract DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciTabooContractIdNotEqual(m1, m2 *models.TabooContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("taboo_contract DNs are equal")
		}
		return nil
	}
}

func CreateTabooContractWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing taboo_contract creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_taboo_contract" "test" {
#		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		
	}
		`
	case "name":
		rBlock += `
	resource "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
#		name = "%s"

	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccTabooContractConfigWithRequiredParams(fvTenantName, rName string) string {
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
	`, fvTenantName, rName)
	return resource
}

func CreateAccTabooContractConfig(fvTenantName, rName string) string {
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
	`, fvTenantName, rName)
	return resource
}

func CreateAccTabooContractConfigs(rName string) string {
	fmt.Println("=== STEP  testing taboo_contract multiple creation")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"

	}
	resource "aci_taboo_contract" "test1" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"

	}
	resource "aci_taboo_contract" "test2" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"

	}
	`, rName, rName, rName+"1", rName+"2")
	return resource
}

func CreateAccTabooContractWithInValidParentDn(rName, prName string) string {
	fmt.Println("=== STEP  Negative Case: testing taboo_contract creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_aaa_domain" "test" {
		name        = "%s"
	}
	
	resource "aci_taboo_contract" "test" {
		tenant_dn  = aci_aaa_domain.test.id
		name  = "%s"	
	
	}
	`, prName, rName)
	return resource
}

func CreateAccTabooContractConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing taboo_contract creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_taboo_contract" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_taboo_contract"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccTabooContractRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing taboo_contract creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_taboo_contract" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_taboo_contract"
		
	}
	`)

	return resource
}

func CreateAccTabooContractConfigUpdatedName(rName, longerName string) string {
	fmt.Println("=== STEP  Basic: testing Tabboo Contract Updation with invalid name of long length")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_taboo_contract" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, rName, longerName)
	return resource
}

func CreateAccTabooContractUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing taboo_contract attribute: %s=%s \n", attribute, value)
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
	`, fvTenantName, rName, attribute, value)
	return resource
}
