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

func TestAccAciImportedContract_Basic(t *testing.T) {
	var imported_contract_default models.ImportedContract
	var imported_contract_updated models.ImportedContract
	resourceName := "aci_imported_contract.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	prOther := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciImportedContractDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccImportedContractWithoutTenant(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccImportedContractWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccImportedContractConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciImportedContractExists(resourceName, &imported_contract_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_if", ""),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
				),
			},
			{
				Config: CreateAccImportedContractConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciImportedContractExists(resourceName, &imported_contract_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag_imported_contract"),
					resource.TestCheckResourceAttr(resourceName, "description", "from terraform"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "example"),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciImportedContractIdEqual(&imported_contract_default, &imported_contract_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccImportedContractRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccImportedContractConfigWithParentAndName(rName, longrName),
				ExpectError: regexp.MustCompile(`failed validation for value`),
			},
			{
				Config: CreateAccImportedContractConfigWithParentAndName(rName, rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciImportedContractExists(resourceName, &imported_contract_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rOther),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciImportedContractIdNotEqual(&imported_contract_default, &imported_contract_updated),
				),
			},
			{
				Config: CreateAccImportedContractConfig(rName),
			},
			{
				Config: CreateAccImportedContractConfigWithParentAndName(prOther, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciImportedContractExists(resourceName, &imported_contract_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", prOther)),
					testAccCheckAciImportedContractIdNotEqual(&imported_contract_default, &imported_contract_updated),
				),
			},
		},
	})
}

func TestAccAciImportedContract_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciImportedContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccImportedContractConfig(rName),
			},
			{
				Config:      CreateAccImportedContractWithInValidTenantDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class vzCPIf (.)+`),
			},
			{
				Config:      CreateAccImportedContractUpdatedAttr(rName, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccImportedContractUpdatedAttr(rName, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccImportedContractUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccImportedContractUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config: CreateAccImportedContractConfig(rName),
			},
		},
	})
}

func TestAccAciImportedContract_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciImportedContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccImportedContractsConfig(rName),
			},
		},
	})
}

func CreateAccImportedContractsConfig(rName string) string {
	fmt.Println("=== STEP  creating multiple imported contract")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_imported_contract" "test1"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_imported_contract" "test2"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_imported_contract" "test3"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	`, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccImportedContractUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing imported contract attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_imported_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccImportedContractWithInValidTenantDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing imported contract creation with invalid tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_imported_contract" "test"{
		tenant_dn = aci_application_profile.test.id
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}
func testAccCheckAciImportedContractExists(name string, imported_contract *models.ImportedContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Imported Contract %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Imported Contract dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		imported_contractFound := models.ImportedContractFromContainer(cont)
		if imported_contractFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Imported Contract %s not found", rs.Primary.ID)
		}
		*imported_contract = *imported_contractFound
		return nil
	}
}

func testAccCheckAciImportedContractDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing imported contract destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_imported_contract" {
			cont, err := client.Get(rs.Primary.ID)
			imported_contract := models.ImportedContractFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Imported Contract %s Still exists", imported_contract.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func CreateAccImportedContractWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing imported_contract creation without creating tenant")
	resource := fmt.Sprintf(`
	resource "aci_imported_contract" "test" {
		name = "%s"
	}
	`, rName)
	return resource
}
func CreateAccImportedContractWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing imported_contract creation without name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_imported_contract" "test" {
		tenant_dn = aci_tenant.test.id
	}
	`, rName)
	return resource
}
func CreateAccImportedContractConfig(rName string) string {
	fmt.Println("=== STEP  testing imported contract creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_imported_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName)
	return resource
}
func CreateAccImportedContractConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  testing imported contract creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_imported_contract" "test" {
		tenant_dn   = aci_tenant.test.id
		name        = "%s"
		annotation  = "tag_imported_contract"
		description = "from terraform"
		name_alias  = "example"
	  }
	`, rName, rName)
	return resource
}
func testAccCheckAciImportedContractIdEqual(ic1, ic2 *models.ImportedContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ic1.DistinguishedName != ic2.DistinguishedName {
			return fmt.Errorf("imported contract imported contracts are not equal")
		}
		return nil
	}
}
func CreateAccImportedContractRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing imported contract updation without required fields")
	resource := fmt.Sprintln(`

	resource "aci_imported_contract" "test" {
		annotation  = "tag_imported_contract"
		description = "from terraform"
		name_alias  = "example"
	  }
	`)
	return resource
}

func testAccCheckAciImportedContractIdNotEqual(c1, c2 *models.ImportedContract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if c1.DistinguishedName == c2.DistinguishedName {
			return fmt.Errorf("imported contract DNs are equal")
		}
		return nil
	}
}

func CreateAccImportedContractConfigWithParentAndName(prName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing imported contract creation with tenant name %s name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_imported_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, prName, rName)
	return resource
}
