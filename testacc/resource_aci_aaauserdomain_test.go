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

func TestAccAciUserSecurityDomain_Basic(t *testing.T) {
	var user_security_domain_default models.UserDomain
	var user_security_domain_updated models.UserDomain
	resourceName := "aci_user_security_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	aaaUserName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciUserSecurityDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateUserSecurityDomainWithoutRequired(aaaUserName, rName, "local_user_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateUserSecurityDomainWithoutRequired(aaaUserName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccUserSecurityDomainConfig(aaaUserName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserSecurityDomainExists(resourceName, &user_security_domain_default),
					resource.TestCheckResourceAttr(resourceName, "local_user_dn", fmt.Sprintf("uni/userext/user-%s", aaaUserName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccUserSecurityDomainConfigWithOptionalValues(aaaUserName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserSecurityDomainExists(resourceName, &user_security_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "local_user_dn", fmt.Sprintf("uni/userext/user-%s", aaaUserName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_user_security_domain"),
					testAccCheckAciUserSecurityDomainIdEqual(&user_security_domain_default, &user_security_domain_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccUserSecurityDomainConfigUpdatedName(aaaUserName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccUserSecurityDomainRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccUserSecurityDomainConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserSecurityDomainExists(resourceName, &user_security_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "local_user_dn", fmt.Sprintf("uni/userext/user-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciUserSecurityDomainIdNotEqual(&user_security_domain_default, &user_security_domain_updated),
				),
			},
			{
				Config: CreateAccUserSecurityDomainConfig(aaaUserName, rName),
			},
			{
				Config: CreateAccUserSecurityDomainConfigWithRequiredParams(aaaUserName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserSecurityDomainExists(resourceName, &user_security_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "local_user_dn", fmt.Sprintf("uni/userext/user-%s", aaaUserName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciUserSecurityDomainIdNotEqual(&user_security_domain_default, &user_security_domain_updated),
				),
			},
		},
	})
}

func TestAccAciUserSecurityDomain_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaUserName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciUserSecurityDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccUserSecurityDomainConfig(aaaUserName, rName),
			},
			{
				Config:      CreateAccUserSecurityDomainWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccUserSecurityDomainUpdatedAttr(aaaUserName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccUserSecurityDomainUpdatedAttr(aaaUserName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccUserSecurityDomainUpdatedAttr(aaaUserName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccUserSecurityDomainUpdatedAttr(aaaUserName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccUserSecurityDomainConfig(aaaUserName, rName),
			},
		},
	})
}

func TestAccAciUserSecurityDomain_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaUserName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciUserSecurityDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccUserSecurityDomainConfigMultiple(aaaUserName, rName),
			},
		},
	})
}

func testAccCheckAciUserSecurityDomainExists(name string, user_security_domain *models.UserDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("User Security Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User Security Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		user_security_domainFound := models.UserDomainFromContainer(cont)
		if user_security_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("User Security Domain %s not found", rs.Primary.ID)
		}
		*user_security_domain = *user_security_domainFound
		return nil
	}
}

func testAccCheckAciUserSecurityDomainDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing user_security_domain destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_user_security_domain" {
			cont, err := client.Get(rs.Primary.ID)
			user_security_domain := models.UserDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("User Security Domain %s Still exists", user_security_domain.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciUserSecurityDomainIdEqual(m1, m2 *models.UserDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("user_security_domain DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciUserSecurityDomainIdNotEqual(m1, m2 *models.UserDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("user_security_domain DNs are equal")
		}
		return nil
	}
}

func CreateUserSecurityDomainWithoutRequired(aaaUserName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing user_security_domain creation without ", attrName)
	rBlock := `
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	`
	switch attrName {
	case "local_user_dn":
		rBlock += `
	resource "aci_user_security_domain" "test" {
	#	local_user_dn = aci_local_user.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaUserName, rName)
}

func CreateAccUserSecurityDomainConfigWithRequiredParams(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing user_security_domain creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}
	`, aaaUserName, rName)
	return resource
}
func CreateAccUserSecurityDomainConfigUpdatedName(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing user_security_domain creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}
	`, aaaUserName, rName)
	return resource
}

func CreateAccUserSecurityDomainConfig(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing user_security_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}
	`, aaaUserName, rName)
	return resource
}

func CreateAccUserSecurityDomainConfigMultiple(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing multiple user_security_domain creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, aaaUserName, rName)
	return resource
}

func CreateAccUserSecurityDomainWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing user_security_domain creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_user_security_domain" "test" {
		local_user_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccUserSecurityDomainConfigWithOptionalValues(aaaUserName, rName string) string {
	fmt.Println("=== STEP  Basic: testing user_security_domain creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
		
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_user_security_domain"
		
	}
	`, aaaUserName, rName)

	return resource
}

func CreateAccUserSecurityDomainRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing user_security_domain updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_user_security_domain" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_user_security_domain"
	}
	`)

	return resource
}

func CreateAccUserSecurityDomainUpdatedAttr(aaaUserName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing user_security_domain attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
		%s = "%s"
	}
	`, aaaUserName, rName, attribute, value)
	return resource
}

func CreateAccUserSecurityDomainUpdatedAttrList(aaaUserName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing user_security_domain attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
		%s = %s
	}
	`, aaaUserName, rName, attribute, value)
	return resource
}
