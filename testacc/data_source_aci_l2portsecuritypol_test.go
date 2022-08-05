package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciPortSecurityPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_port_security_policy.test"
	dataSourceName := "data.aci_port_security_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPortSecurityPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreatePortSecurityPolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccPortSecurityPolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "maximum", resourceName, "maximum"),
					resource.TestCheckResourceAttrPair(dataSourceName, "timeout", resourceName, "timeout"),
					resource.TestCheckResourceAttrPair(dataSourceName, "violation", resourceName, "violation"),
				),
			},
			{
				Config:      CreateAccPortSecurityPolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccPortSecurityPolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccPortSecurityPolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccPortSecurityPolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing port_security_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_port_security_policy" "test" {
	
		name  = "%s"
	}

	data "aci_port_security_policy" "test" {
	
		name  = aci_port_security_policy.test.name
		depends_on = [ aci_port_security_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccPortSecurityPolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing port_security_policy Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_port_security_policy" "test" {
		name  = "%s"
	}

	data "aci_port_security_policy" "test" {
		name  = "${aci_port_security_policy.test.name}_invalid"
		depends_on = [ aci_port_security_policy.test ]
	}
	`, rName)
	return resource
}

func CreatePortSecurityPolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing port_security_policy data-source without ", attrName)
	rBlock := `
	
	resource "aci_port_security_policy" "test" {
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_port_security_policy" "test" {
	#	name  = "%s"
		depends_on = [ aci_port_security_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccPortSecurityPolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing port_security_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_port_security_policy" "test" {
	
		name  = "%s"
	}

	data "aci_port_security_policy" "test" {
	
		name  = aci_port_security_policy.test.name
		%s = "%s"
		depends_on = [ aci_port_security_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccPortSecurityPolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing port_security_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_port_security_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_port_security_policy" "test" {
	
		name  = aci_port_security_policy.test.name
		depends_on = [ aci_port_security_policy.test ]
	}
	`, rName, key, value)
	return resource
}
