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

func TestAccAciPhysicalDomain_Basic(t *testing.T) {
	var physical_domain_default models.PhysicalDomain
	var physical_domain_updated models.PhysicalDomain
	resourceName := "aci_physical_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPhysicalDomainDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreatePhysicalDomainWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccPhysicalDomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPhysicalDomainExists(resourceName, &physical_domain_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccPhysicalDomainConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPhysicalDomainExists(resourceName, &physical_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_physical_domain"),
					testAccCheckAciPhysicalDomainIdEqual(&physical_domain_default, &physical_domain_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccPhysicalDomainConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccPhysicalDomainRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccPhysicalDomainConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPhysicalDomainExists(resourceName, &physical_domain_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciPhysicalDomainIdNotEqual(&physical_domain_default, &physical_domain_updated),
				),
			},
		},
	})
}

func TestAccAciPhysicalDomain_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPhysicalDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPhysicalDomainConfig(rName),
			},
			{
				Config:      CreateAccPhysicalDomainUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPhysicalDomainUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPhysicalDomainUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccPhysicalDomainConfig(rName),
			},
		},
	})
}

func TestAccAciPhysicalDomain_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPhysicalDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPhysicalDomainConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciPhysicalDomainExists(name string, physical_domain *models.PhysicalDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Physical Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Physical Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		physical_domainFound := models.PhysicalDomainFromContainer(cont)
		if physical_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Physical Domain %s not found", rs.Primary.ID)
		}
		*physical_domain = *physical_domainFound
		return nil
	}
}

func testAccCheckAciPhysicalDomainDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing physical_domain destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_physical_domain" {
			cont, err := client.Get(rs.Primary.ID)
			physical_domain := models.PhysicalDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Physical Domain %s Still exists", physical_domain.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciPhysicalDomainIdEqual(m1, m2 *models.PhysicalDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("physical_domain DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciPhysicalDomainIdNotEqual(m1, m2 *models.PhysicalDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("physical_domain DNs are equal")
		}
		return nil
	}
}

func CreatePhysicalDomainWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing physical_domain creation without ", attrName)
	rBlock := `

	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_physical_domain" "test" {

	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccPhysicalDomainConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing physical_domain creation with Updated required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_physical_domain" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccPhysicalDomainConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing physical_domain creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	resource "aci_physical_domain" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccPhysicalDomainConfig(rName string) string {
	fmt.Println("=== STEP  testing physical_domain creation with required arguments", rName)
	resource := fmt.Sprintf(`
	resource "aci_physical_domain" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccPhysicalDomainConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple physical_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_physical_domain" "test" {
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccPhysicalDomainConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing physical_domain creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_physical_domain" "test" {
		name  = "%s"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_physical_domain"
	}
	`, rName)

	return resource
}

func CreateAccPhysicalDomainRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing physical_domain updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_physical_domain" "test" {
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_physical_domain"
	}
	`)
	return resource
}

func CreateAccPhysicalDomainUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing physical_domain attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_physical_domain" "test" {
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
