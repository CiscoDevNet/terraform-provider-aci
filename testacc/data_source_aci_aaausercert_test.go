package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAcix509CertificateDataSource_Basic(t *testing.T) {
	resourceName := "aci_x509_certificate.test"
	dataSourceName := "data.aci_x509_certificate.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	aaaUserName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAcix509CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      Createx509CertificateDSWithoutRequired(aaaUserName, rName, "local_user_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      Createx509CertificateDSWithoutRequired(aaaUserName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccx509CertificateConfigDataSource(aaaUserName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "local_user_dn", resourceName, "local_user_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "data", resourceName, "data"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccx509CertificateDataSourceUpdate(aaaUserName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccx509CertificateDSWithInvalidParentDn(aaaUserName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccx509CertificateDataSourceUpdatedResource(aaaUserName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccx509CertificateConfigDataSource(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing x509_certificate Data Source with required arguments only")
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

	data "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = aci_x509_certificate.test.name
		depends_on = [ aci_x509_certificate.test ]
	}
	`, aaaUserName, rName, certificate_terraformuser)
	return resource
}

func Createx509CertificateDSWithoutRequired(aaaUserName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing x509_certificate Data Source without ", attrName)
	rBlock := `
	
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
	`
	switch attrName {
	case "local_user_dn":
		rBlock += `
	data "aci_x509_certificate" "test" {
	#	local_user_dn  = aci_local_user.test.id
		name  = aci_x509_certificate.test.name
		depends_on = [ aci_x509_certificate.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
	#	name  = aci_x509_certificate.test.name
		depends_on = [ aci_x509_certificate.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaUserName, rName, certificate_terraformuser)
}

func CreateAccx509CertificateDSWithInvalidParentDn(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing x509_certificate Data Source with Invalid Parent Dn")
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

	data "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = "${aci_x509_certificate.test.name}_invalid"
		depends_on = [ aci_x509_certificate.test ]
	}
	`, aaaUserName, rName, certificate_terraformuser)
	return resource
}

func CreateAccx509CertificateDataSourceUpdate(aaaUserName, rName, key, value string) string {
	fmt.Println("=== STEP  testing x509_certificate Data Source with random attribute")
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

	data "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = aci_x509_certificate.test.name
		%s = "%s"
		depends_on = [ aci_x509_certificate.test ]
	}
	`, aaaUserName, rName, certificate_terraformuser, key, value)
	return resource
}

func CreateAccx509CertificateDataSourceUpdatedResource(aaaUserName, rName, key, value string) string {
	fmt.Println("=== STEP  testing x509_certificate Data Source with updated resource")
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

	data "aci_x509_certificate" "test" {
		local_user_dn  = aci_local_user.test.id
		name  = aci_x509_certificate.test.name
		depends_on = [ aci_x509_certificate.test ]
	}
	`, aaaUserName, rName, certificate_terraformuser, key, value)
	return resource
}
