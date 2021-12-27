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

func TestAccAciAAADomain_Basic(t *testing.T) {
	var aaa_domain_default models.SecurityDomain
	var aaa_domain_updated models.SecurityDomain
	resourceName := "aci_aaa_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAAADomainDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateAAADomainWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAAADomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAAADomainExists(resourceName, &aaa_domain_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccAAADomainConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAAADomainExists(resourceName, &aaa_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_aaa_domain"),
					testAccCheckAciAAADomainIdEqual(&aaa_domain_default, &aaa_domain_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccAAADomainConfig(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccAAADomainRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAAADomainConfig(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAAADomainExists(resourceName, &aaa_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciAAADomainIdNotEqual(&aaa_domain_default, &aaa_domain_updated),
				),
			},
		},
	})
}

func TestAccAciAAADomain_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAAADomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAAADomainConfig(rName),
			},
			{
				Config:      CreateAccAAADomainUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAAADomainUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAAADomainUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAAADomainUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)+ is not expected here.`),
			},
			{
				Config: CreateAccAAADomainConfig(rName),
			},
		},
	})
}

func TestAccAciAAADomain_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAAADomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAAADomainConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciAAADomainExists(name string, aaa_domain *models.SecurityDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("AAA Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AAA Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		aaa_domainFound := models.SecurityDomainFromContainer(cont)
		if aaa_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("AAA Domain %s not found", rs.Primary.ID)
		}
		*aaa_domain = *aaa_domainFound
		return nil
	}
}

func testAccCheckAciAAADomainDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing aaa_domain destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_aaa_domain" {
			cont, err := client.Get(rs.Primary.ID)
			aaa_domain := models.SecurityDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("AAA Domain %s Still exists", aaa_domain.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAAADomainIdEqual(m1, m2 *models.SecurityDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("aaa_domain DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciAAADomainIdNotEqual(m1, m2 *models.SecurityDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("aaa_domain DNs are equal")
		}
		return nil
	}
}

func CreateAAADomainWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing aaa_domain creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_aaa_domain" "test" {
	
	#	name  = "%s"
		description = "created while acceptance testing"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccAAADomainConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple aaa_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_aaa_domain" "test1" {
		name  = "%s"
	}

	resource "aci_aaa_domain" "test2" {
		name  = "%s"
	}

	resource "aci_aaa_domain" "test3" {
		name  = "%s"
	}
	`, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccAAADomainConfig(rName string) string {
	fmt.Println("=== STEP  testing aaa_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_aaa_domain" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccAAADomainConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing aaa_domain creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_aaa_domain" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_aaa_domain"
	}
	`, rName)

	return resource
}

func CreateAccAAADomainRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing aaa_domain creation with optional parameters")
	resource := fmt.Sprintln(`
	resource "aci_aaa_domain" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_aaa_domain"
	}
	`)

	return resource
}

func CreateAccAAADomainUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing aaa_domain attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_aaa_domain" "test" {
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
