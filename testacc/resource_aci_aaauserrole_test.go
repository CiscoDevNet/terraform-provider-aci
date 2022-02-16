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

func TestAccAciUserSecurityDomainRole_Basic(t *testing.T) {
	var user_security_domain_role_default models.UserRole
	var user_security_domain_role_updated models.UserRole
	resourceName := "aci_user_security_domain_role.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	aaaUserName := makeTestVariable(acctest.RandString(5))
	aaaUserDomainName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciUserSecurityDomainRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateUserSecurityDomainRoleWithoutRequired(aaaUserName, aaaUserDomainName, rName, "user_domain_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateUserSecurityDomainRoleWithoutRequired(aaaUserName, aaaUserDomainName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccUserSecurityDomainRoleConfig(aaaUserName, aaaUserDomainName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserSecurityDomainRoleExists(resourceName, &user_security_domain_role_default),
					resource.TestCheckResourceAttr(resourceName, "user_domain_dn", fmt.Sprintf("uni/userext/user-%s/userdomain-%s", aaaUserName, aaaUserDomainName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "priv_type", "readPriv"),
				),
			},
			{
				Config: CreateAccUserSecurityDomainRoleConfigWithOptionalValues(aaaUserName, aaaUserDomainName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserSecurityDomainRoleExists(resourceName, &user_security_domain_role_updated),
					resource.TestCheckResourceAttr(resourceName, "user_domain_dn", fmt.Sprintf("uni/userext/user-%s/userdomain-%s", aaaUserName, aaaUserDomainName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_user_security_domain_role"),
					resource.TestCheckResourceAttr(resourceName, "priv_type", "writePriv"),
					testAccCheckAciUserSecurityDomainRoleIdEqual(&user_security_domain_role_default, &user_security_domain_role_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccUserSecurityDomainRoleConfigUpdatedName(aaaUserName, aaaUserDomainName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccUserSecurityDomainRoleRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccUserSecurityDomainRoleConfigWithRequiredParams(aaaUserName, rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserSecurityDomainRoleExists(resourceName, &user_security_domain_role_updated),
					resource.TestCheckResourceAttr(resourceName, "user_domain_dn", fmt.Sprintf("uni/userext/user-%s/userdomain-%s", aaaUserName, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciUserSecurityDomainRoleIdNotEqual(&user_security_domain_role_default, &user_security_domain_role_updated),
				),
			},
			{
				Config: CreateAccUserSecurityDomainRoleConfig(aaaUserName, aaaUserDomainName, rName),
			},
			{
				Config: CreateAccUserSecurityDomainRoleConfigWithRequiredParams(aaaUserName, aaaUserDomainName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciUserSecurityDomainRoleExists(resourceName, &user_security_domain_role_updated),
					resource.TestCheckResourceAttr(resourceName, "user_domain_dn", fmt.Sprintf("uni/userext/user-%s/userdomain-%s", aaaUserName, aaaUserDomainName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciUserSecurityDomainRoleIdNotEqual(&user_security_domain_role_default, &user_security_domain_role_updated),
				),
			},
		},
	})
}

func TestAccAciUserSecurityDomainRole_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaUserName := makeTestVariable(acctest.RandString(5))
	aaaUserDomainName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciUserSecurityDomainRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccUserSecurityDomainRoleConfig(aaaUserName, aaaUserDomainName, rName),
			},
			{
				Config:      CreateAccUserSecurityDomainRoleWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccUserSecurityDomainRoleUpdatedAttr(aaaUserName, aaaUserDomainName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccUserSecurityDomainRoleUpdatedAttr(aaaUserName, aaaUserDomainName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccUserSecurityDomainRoleUpdatedAttr(aaaUserName, aaaUserDomainName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccUserSecurityDomainRoleUpdatedAttr(aaaUserName, aaaUserDomainName, rName, "priv_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccUserSecurityDomainRoleUpdatedAttr(aaaUserName, aaaUserDomainName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccUserSecurityDomainRoleConfig(aaaUserName, aaaUserDomainName, rName),
			},
		},
	})
}

func TestAccAciUserSecurityDomainRole_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaUserName := makeTestVariable(acctest.RandString(5))
	aaaUserDomainName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciUserSecurityDomainRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccUserSecurityDomainRoleConfigMultiple(aaaUserName, aaaUserDomainName, rName),
			},
		},
	})
}

func testAccCheckAciUserSecurityDomainRoleExists(name string, user_security_domain_role *models.UserRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("User Security Domain Role %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User Security Domain Role dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		user_security_domain_roleFound := models.UserRoleFromContainer(cont)
		if user_security_domain_roleFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("User Security Domain Role %s not found", rs.Primary.ID)
		}
		*user_security_domain_role = *user_security_domain_roleFound
		return nil
	}
}

func testAccCheckAciUserSecurityDomainRoleDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing user_security_domain_role destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_user_security_domain_role" {
			cont, err := client.Get(rs.Primary.ID)
			user_security_domain_role := models.UserRoleFromContainer(cont)
			if err == nil {
				return fmt.Errorf("User Security Domain Role %s Still exists", user_security_domain_role.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciUserSecurityDomainRoleIdEqual(m1, m2 *models.UserRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("user_security_domain_role DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciUserSecurityDomainRoleIdNotEqual(m1, m2 *models.UserRole) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("user_security_domain_role DNs are equal")
		}
		return nil
	}
}

func CreateUserSecurityDomainRoleWithoutRequired(aaaUserName, aaaUserDomainName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing user_security_domain_role creation without ", attrName)
	rBlock := `
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}
	
	`
	switch attrName {
	case "user_domain_dn":
		rBlock += `
	resource "aci_user_security_domain_role" "test" {
	#	user_domain_dn  =  aci_user_security_domain.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_user_security_domain_role" "test" {
		user_domain_dn  =  aci_user_security_domain.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaUserName, aaaUserDomainName, rName)
}

func CreateAccUserSecurityDomainRoleConfigWithRequiredParams(aaaUserName, aaaUserDomainName, rName string) string {
	fmt.Println("=== STEP  testing user_security_domain_role creation with updated naming arguments")
	resource := fmt.Sprintf(`
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}
	
	resource "aci_user_security_domain_role" "test" {
		user_domain_dn  =  aci_user_security_domain.test.id
		name  = "%s"
	}
	`, aaaUserName, aaaUserDomainName, rName)
	return resource
}
func CreateAccUserSecurityDomainRoleConfigUpdatedName(aaaUserName, aaaUserDomainName, rName string) string {
	fmt.Println("=== STEP  testing user_security_domain_role creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}
	
	resource "aci_user_security_domain_role" "test" {
		user_domain_dn  =  aci_user_security_domain.test.id
		name  = "%s"
	}
	`, aaaUserName, aaaUserDomainName, rName)
	return resource
}

func CreateAccUserSecurityDomainRoleConfig(aaaUserName, aaaUserDomainName, rName string) string {
	fmt.Println("=== STEP  testing user_security_domain_role creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}

	resource "aci_user_security_domain_role" "test" {
		user_domain_dn  =  aci_user_security_domain.test.id
		name  = "%s"
	}
	`, aaaUserName, aaaUserDomainName, rName)
	return resource
}

func CreateAccUserSecurityDomainRoleConfigMultiple(aaaUserName, aaaUserDomainName, rName string) string {
	fmt.Println("=== STEP  testing multiple user_security_domain_role creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}
	
	resource "aci_user_security_domain_role" "test" {
		user_domain_dn  =  aci_user_security_domain.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, aaaUserName, aaaUserDomainName, rName)
	return resource
}

func CreateAccUserSecurityDomainRoleWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing user_security_domain_role creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_user_security_domain_role" "test" {
		user_domain_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccUserSecurityDomainRoleConfigWithOptionalValues(aaaUserName, aaaUserDomainName, rName string) string {
	fmt.Println("=== STEP  Basic: testing user_security_domain_role creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}
	
	resource "aci_user_security_domain_role" "test" {
		user_domain_dn  = "${ aci_user_security_domain.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_user_security_domain_role"
		priv_type = "writePriv"
		
	}
	`, aaaUserName, aaaUserDomainName, rName)

	return resource
}

func CreateAccUserSecurityDomainRoleRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing user_security_domain_role updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_user_security_domain_role" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_user_security_domain_role"
		priv_type = "writePriv"
	}
	`)

	return resource
}

func CreateAccUserSecurityDomainRoleUpdatedAttr(aaaUserName, aaaUserDomainName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing user_security_domain_role attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}

	resource "aci_user_security_domain_role" "test" {
		user_domain_dn  =  aci_user_security_domain.test.id
		name  = "%s"
		%s = "%s"
	}
	`, aaaUserName, aaaUserDomainName, rName, attribute, value)
	return resource
}

func CreateAccUserSecurityDomainRoleUpdatedAttrList(aaaUserName, aaaUserDomainName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing user_security_domain_role attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}
	
	resource "aci_user_security_domain_role" "test" {
		user_domain_dn  =  aci_user_security_domain.test.id
		name  = "%s"
		%s = %s
	}
	`, aaaUserName, aaaUserDomainName, rName, attribute, value)
	return resource
}
