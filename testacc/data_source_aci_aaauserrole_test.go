package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciUserSecurityDomainRoleDataSource_Basic(t *testing.T) {
	resourceName := "aci_user_security_domain_role.test"
	dataSourceName := "data.aci_user_security_domain_role.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	aaaUserName := makeTestVariable(acctest.RandString(5))
	aaaUserDomainName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciUserSecurityDomainRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateUserSecurityDomainRoleDSWithoutRequired(aaaUserName, aaaUserDomainName, rName, "user_domain_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateUserSecurityDomainRoleDSWithoutRequired(aaaUserName, aaaUserDomainName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccUserSecurityDomainRoleConfigDataSource(aaaUserName, aaaUserDomainName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "user_domain_dn", resourceName, "user_domain_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "priv_type", resourceName, "priv_type"),
				),
			},
			{
				Config:      CreateAccUserSecurityDomainRoleDataSourceUpdate(aaaUserName, aaaUserDomainName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccUserSecurityDomainRoleDSWithInvalidParentDn(aaaUserName, aaaUserDomainName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccUserSecurityDomainRoleDataSourceUpdatedResource(aaaUserName, aaaUserDomainName, rName, "description", "description_testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccUserSecurityDomainRoleConfigDataSource(aaaUserName, aaaUserDomainName, rName string) string {
	fmt.Println("=== STEP  testing user_security_domain_role Data Source with required arguments only")
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

	data "aci_user_security_domain_role" "test" {
		user_domain_dn  =  aci_user_security_domain.test.id
		name  = aci_user_security_domain_role.test.name
		depends_on = [ aci_user_security_domain_role.test ]
	}
	`, aaaUserName, aaaUserDomainName, rName)
	return resource
}

func CreateUserSecurityDomainRoleDSWithoutRequired(aaaUserName, aaaUserDomainName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing user_security_domain_role Data Source without ", attrName)
	rBlock := `
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
	`
	switch attrName {
	case "user_domain_dn":
		rBlock += `
	data "aci_user_security_domain_role" "test" {
	#	user_domain_dn  =  aci_user_security_domain.test.id
		name  = aci_user_security_domain_role.test.name
		depends_on = [ aci_user_security_domain_role.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_user_security_domain_role" "test" {
		user_domain_dn  =  aci_user_security_domain.test.id
	#	name  = aci_user_security_domain_role.test.name
		depends_on = [ aci_user_security_domain_role.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, aaaUserName, aaaUserDomainName, rName)
}

func CreateAccUserSecurityDomainRoleDSWithInvalidParentDn(aaaUserName, aaaUserDomainName, rName string) string {
	fmt.Println("=== STEP  testing user_security_domain_role Data Source with Invalid Parent Dn")
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

	data "aci_user_security_domain_role" "test" {
		user_domain_dn  =  aci_user_security_domain.test.id
		name  = "${aci_user_security_domain_role.test.name}_invalid"
		depends_on = [ aci_user_security_domain_role.test ]
	}
	`, aaaUserName, aaaUserDomainName, rName)
	return resource
}

func CreateAccUserSecurityDomainRoleDataSourceUpdate(aaaUserName, aaaUserDomainName, rName, key, value string) string {
	fmt.Println("=== STEP  testing user_security_domain_role Data Source with random attribute")
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

	data "aci_user_security_domain_role" "test" {
		user_domain_dn  =  aci_user_security_domain.test.id
		name  = aci_user_security_domain_role.test.name
		%s = "%s"
		depends_on = [ aci_user_security_domain_role.test ]
	}
	`, aaaUserName, aaaUserDomainName, rName, key, value)
	return resource
}

func CreateAccUserSecurityDomainRoleDataSourceUpdatedResource(aaaUserName, aaaUserDomainName, rName, key, value string) string {
	fmt.Println("=== STEP  testing user_security_domain_role Data Source with updated resource")
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

	data "aci_user_security_domain_role" "test" {
		user_domain_dn  =  aci_user_security_domain.test.id
		name  = aci_user_security_domain_role.test.name
		depends_on = [ aci_user_security_domain_role.test ]
	}
	`, aaaUserName, aaaUserDomainName, rName, key, value)
	return resource
}
