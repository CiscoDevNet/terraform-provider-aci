package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciUserSecurityDomainDataSource_Basic(t *testing.T) {
	resourceName := "aci_user_security_domain.test"
	dataSourceName := "data.aci_user_security_domain.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	aaaUserName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciUserSecurityDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateUserSecurityDomainDSWithoutRequired(aaaUserName, rName, "local_user_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateUserSecurityDomainDSWithoutRequired(aaaUserName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccUserSecurityDomainConfigDataSource(aaaUserName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "local_user_dn", resourceName, "local_user_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccUserSecurityDomainDataSourceUpdate(aaaUserName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccUserSecurityDomainDSWithInvalidParentDn(aaaUserName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccUserSecurityDomainDataSourceUpdatedResource(aaaUserName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccUserSecurityDomainConfigDataSource(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing user_security_domain Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}

	data "aci_user_security_domain" "test" {
		local_user_dn = aci_user_security_domain.test.local_user_dn
		name  = aci_user_security_domain.test.name
		depends_on = [ aci_user_security_domain.test ]
	}
	`, aaaUserName, rName)
	return resource
}

func CreateUserSecurityDomainDSWithoutRequired(aaaUserName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing user_security_domain Data Source without ", attrName)
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
	case "local_user_dn":
		rBlock += `
	data "aci_user_security_domain" "test" {
	#	local_user_dn  = aci_local_user.test.id
		name  = aci_user_security_domain.test.name
		depends_on = [ aci_user_security_domain.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_user_security_domain" "test" {
		local_user_dn  = aci_local_user.test.id
	#	name  = aci_user_security_domain.test.name
		depends_on = [ aci_user_security_domain.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaUserName, rName)
}

func CreateAccUserSecurityDomainDSWithInvalidParentDn(aaaUserName, rName string) string {
	fmt.Println("=== STEP  testing user_security_domain Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}
	
	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}

	data "aci_user_security_domain" "test" {
		local_user_dn  = "${aci_local_user.test.id}_invalid"
		name  = aci_user_security_domain.test.name
		depends_on = [ aci_user_security_domain.test ]
	}
	`, aaaUserName, rName)
	return resource
}

func CreateAccUserSecurityDomainDataSourceUpdate(aaaUserName, rName, key, value string) string {
	fmt.Println("=== STEP  testing user_security_domain Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_local_user" "test" {
		name 		= "%s"
		pwd = "Test_coverage"
	}

	resource "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = "%s"
	}

	data "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = aci_user_security_domain.test.name
		%s = "%s"
		depends_on = [ aci_user_security_domain.test ]
	}
	`, aaaUserName, rName, key, value)
	return resource
}

func CreateAccUserSecurityDomainDataSourceUpdatedResource(aaaUserName, rName, key, value string) string {
	fmt.Println("=== STEP  testing user_security_domain Data Source with updated resource")
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

	data "aci_user_security_domain" "test" {
		local_user_dn = aci_local_user.test.id
		name  = aci_user_security_domain.test.name
		depends_on = [ aci_user_security_domain.test ]
	}
	`, aaaUserName, rName, key, value)
	return resource
}
