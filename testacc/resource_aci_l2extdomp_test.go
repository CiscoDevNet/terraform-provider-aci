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

func TestAccAciL2Domain_Basic(t *testing.T) {
	var l2_domain_default models.L2Domain
	var l2_domain_updated models.L2Domain
	resourceName := "aci_l2_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2DomainDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateL2DomainWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL2DomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2DomainExists(resourceName, &l2_domain_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccL2DomainConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2DomainExists(resourceName, &l2_domain_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l2_domain"),

					testAccCheckAciL2DomainIdEqual(&l2_domain_default, &l2_domain_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL2DomainConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccL2DomainRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccL2DomainConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2DomainExists(resourceName, &l2_domain_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciL2DomainIdNotEqual(&l2_domain_default, &l2_domain_updated),
				),
			},
		},
	})
}

func TestAccAciL2Domain_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2DomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL2DomainConfig(rName),
			},
			{
				Config:      CreateAccL2DomainUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL2DomainUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccL2DomainUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL2DomainConfig(rName),
			},
		},
	})
}

func TestAccAciL2Domain_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2DomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL2DomainConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciL2DomainExists(name string, l2_domain *models.L2Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L2 Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L2 Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l2_domainFound := models.L2DomainFromContainer(cont)
		if l2_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L2 Domain %s not found", rs.Primary.ID)
		}
		*l2_domain = *l2_domainFound
		return nil
	}
}

func testAccCheckAciL2DomainDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l2_domain destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l2_domain" {
			cont, err := client.Get(rs.Primary.ID)
			l2_domain := models.L2DomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L2 Domain %s Still exists", l2_domain.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL2DomainIdEqual(m1, m2 *models.L2Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l2_domain DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL2DomainIdNotEqual(m1, m2 *models.L2Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l2_domain DNs are equal")
		}
		return nil
	}
}

func CreateL2DomainWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l2_domain creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_l2_domain" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccL2DomainConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing l2_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_domain" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccL2DomainConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing l2_domain creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_l2_domain" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccL2DomainConfig(rName string) string {
	fmt.Println("=== STEP  testing l2_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_domain" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccL2DomainConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple l2_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_domain" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccL2DomainConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing l2_domain creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_domain" "test" {
	
		name  = "%s"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l2_domain"
		
	}
	`, rName)

	return resource
}

func CreateAccL2DomainRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l2_domain updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_l2_domain" "test" {
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l2_domain"
		
	}
	`)

	return resource
}

func CreateAccL2DomainUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l2_domain attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_l2_domain" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
