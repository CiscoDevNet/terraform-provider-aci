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

func TestAccAciFCDomain_Basic(t *testing.T) {
	var fc_domain_default models.FCDomain
	var fc_domain_updated models.FCDomain
	resourceName := "aci_fc_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFCDomainDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFCDomainWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFCDomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFCDomainExists(resourceName, &fc_domain_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccFCDomainConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFCDomainExists(resourceName, &fc_domain_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_fc_domain"),

					testAccCheckAciFCDomainIdEqual(&fc_domain_default, &fc_domain_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFCDomainConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccFCDomainRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccFCDomainConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFCDomainExists(resourceName, &fc_domain_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciFCDomainIdNotEqual(&fc_domain_default, &fc_domain_updated),
				),
			},
		},
	})
}

func TestAccAciFCDomain_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFCDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFCDomainConfig(rName),
			},
			{
				Config:      CreateAccFCDomainUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFCDomainUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccFCDomainUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFCDomainConfig(rName),
			},
		},
	})
}

func TestAccAciFCDomain_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFCDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFCDomainConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciFCDomainExists(name string, fc_domain *models.FCDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("FC Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No FC Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fc_domainFound := models.FCDomainFromContainer(cont)
		if fc_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("FC Domain %s not found", rs.Primary.ID)
		}
		*fc_domain = *fc_domainFound
		return nil
	}
}

func testAccCheckAciFCDomainDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing fc_domain destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_fc_domain" {
			cont, err := client.Get(rs.Primary.ID)
			fc_domain := models.FCDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("FC Domain %s Still exists", fc_domain.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFCDomainIdEqual(m1, m2 *models.FCDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("fc_domain DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFCDomainIdNotEqual(m1, m2 *models.FCDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("fc_domain DNs are equal")
		}
		return nil
	}
}

func CreateFCDomainWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fc_domain creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_fc_domain" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFCDomainConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing fc_domain creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_fc_domain" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccFCDomainConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing fc_domain creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_fc_domain" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFCDomainConfig(rName string) string {
	fmt.Println("=== STEP  testing fc_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fc_domain" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFCDomainConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple fc_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fc_domain" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccFCDomainConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing fc_domain creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_fc_domain" "test" {
	
		name  = "%s"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fc_domain"
		
	}
	`, rName)

	return resource
}

func CreateAccFCDomainRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing fc_domain updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_fc_domain" "test" {
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fc_domain"
	}
	`)

	return resource
}

func CreateAccFCDomainUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing fc_domain attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fc_domain" "test" {
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
