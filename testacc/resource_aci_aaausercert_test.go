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

func TestAccAcix509Certificate_Basic(t *testing.T) {
	var x509_certificate_default models.X509Certificate
	var x509_certificate_updated models.X509Certificate
	resourceName := "aci_x509_certificate.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	aaaUserName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAcix509CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      Createx509CertificateWithoutRequired(aaaUserName, rName, "local_user_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      Createx509CertificateWithoutRequired(aaaUserName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      Createx509CertificateWithoutRequired(aaaUserName, rName, "data"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccx509CertificateConfig(aaaUserName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAcix509CertificateExists(resourceName, &x509_certificate_default),
					resource.TestCheckResourceAttr(resourceName, "local_user_dn", fmt.Sprintf("uni/userext/user-%s", aaaUserName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccx509CertificateConfigWithOptionalValues(aaaUserName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAcix509CertificateExists(resourceName, &x509_certificate_updated),
					resource.TestCheckResourceAttr(resourceName, "local_user_dn", fmt.Sprintf("uni/userext/user-%s", aaaUserName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_x509_certificate"),
					testAccCheckAcix509CertificateIdEqual(&x509_certificate_default, &x509_certificate_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccx509CertificateConfigUpdatedName(aaaUserName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccx509CertificateConfigUpdatedData(aaaUserName, rName),
				ExpectError: regexp.MustCompile(`Validation of (.)+ X509 Certificate submitted failed - error returned by low level validation routine is`),
			},
			{
				Config:      CreateAccx509CertificateRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccx509CertificateConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAcix509CertificateExists(resourceName, &x509_certificate_updated),
					resource.TestCheckResourceAttr(resourceName, "local_user_dn", fmt.Sprintf("uni/userext/user-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAcix509CertificateIdNotEqual(&x509_certificate_default, &x509_certificate_updated),
				),
			},
			{
				Config: CreateAccx509CertificateConfig(aaaUserName, rName),
			},
			{
				Config: CreateAccx509CertificateConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAcix509CertificateExists(resourceName, &x509_certificate_updated),
					resource.TestCheckResourceAttr(resourceName, "local_user_dn", fmt.Sprintf("uni/userext/user-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAcix509CertificateIdNotEqual(&x509_certificate_default, &x509_certificate_updated),
				),
			},
		},
	})
}

func TestAccAcix509Certificate_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaUserName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAcix509CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccx509CertificateConfig(aaaUserName, rName),
			},
			{
				Config:      CreateAccx509CertificateWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccx509CertificateUpdatedAttr(aaaUserName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccx509CertificateUpdatedAttr(aaaUserName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccx509CertificateUpdatedAttr(aaaUserName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccx509CertificateUpdatedAttr(aaaUserName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccx509CertificateConfig(aaaUserName, rName),
			},
		},
	})
}

func TestAccAcix509Certificate_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	aaaUserName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAcix509CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccx509CertificateConfigMultiple(aaaUserName, rName),
			},
		},
	})
}

func testAccCheckAcix509CertificateExists(name string, x509_certificate *models.X509Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("x509 Certificate %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No x509 Certificate dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		x509_certificateFound := models.X509CertificateFromContainer(cont)
		if x509_certificateFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("x509 Certificate %s not found", rs.Primary.ID)
		}
		*x509_certificate = *x509_certificateFound
		return nil
	}
}

func testAccCheckAcix509CertificateDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing x509_certificate destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_x509_certificate" {
			cont, err := client.Get(rs.Primary.ID)
			x509_certificate := models.X509CertificateFromContainer(cont)
			if err == nil {
				return fmt.Errorf("x509 Certificate %s Still exists", x509_certificate.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAcix509CertificateIdEqual(m1, m2 *models.X509Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("x509_certificate DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAcix509CertificateIdNotEqual(m1, m2 *models.X509Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("x509_certificate DNs are equal")
		}
		return nil
	}
}

func Createx509CertificateWithoutRequired(aaaUserName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing x509_certificate creation without ", attrName)
	rBlock := `
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"	
	}
	
	`
	switch attrName {
	case "local_user_dn":
		rBlock += `
	resource "aci_x509_certificate" "test" {
	#	local_user_dn  = aci_local_user.test.id
		name  = "%s"
		data = <<EOF%s
		EOF 
	}
		`
	case "name":
		rBlock += `
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
	#	name  = "%s"
		data = <<EOF%s
		EOF
	}
		`
	case "data":
		rBlock += `
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
	}
		`
		return fmt.Sprintf(rBlock, aaaUserName, rName)
	}
	return fmt.Sprintf(rBlock, aaaUserName, rName, certificate_terraformuser)
}

func CreateAccx509CertificateConfigWithRequiredParams(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing x509_certificate creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
		data = <<EOF%s
		EOF
	}
	`, aaaUserName, rName, certificate_terraformuser)
	return resource
}
func CreateAccx509CertificateConfigUpdatedName(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing x509_certificate creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
		data = <<EOF%s
		EOF
	}
	`, aaaUserName, rName, certificate_terraformuser)
	return resource
}

func CreateAccx509CertificateConfigUpdatedData(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing x509_certificate creation with invalid Data = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
		data = "%s"
	}
	`, aaaUserName, rName, rName)
	return resource
}

func CreateAccx509CertificateConfig(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing x509_certificate creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
		data = <<EOF%s
		EOF
	}
	`, aaaUserName, rName, certificate_terraformuser)
	return resource
}

func CreateAccx509CertificateConfigMultiple(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing multiple x509_certificate creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s_${count.index}"
		count = 5
		data = <<EOF%s
		EOF
	}
	`, aaaUserName, rName, certificate_terraformuser)
	return resource
}

func CreateAccx509CertificateWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing x509_certificate creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_tenant.test.id
		name  = "%s"	
		data = <<EOF%s
		EOF
	}
	`, rName, rName, certificate_terraformuser)
	return resource
}

func CreateAccx509CertificateConfigWithOptionalValues(aaaUserName, rName string) string {
	fmt.Println("=== STEP  Basic: testing x509_certificate creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = "${aci_local_user.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_x509_certificate"
		data = <<EOF%s
		EOF
	}
	`, aaaUserName, rName, certificate_terraformuser)

	return resource
}

func CreateAccx509CertificateRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing x509_certificate updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_x509_certificate" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_x509_certificate"
	}
	`)

	return resource
}

func CreateAccx509CertificateUpdatedAttr(aaaUserName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing x509_certificate attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "%s"
		data = <<EOF%s
		EOF
		%s = "%s"
	}
	`, aaaUserName, rName, certificate_terraformuser, attribute, value)
	return resource
}
