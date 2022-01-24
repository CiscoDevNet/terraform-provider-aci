package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciDHCPOptionPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_dhcp_option_policy.test"
	dataSourceName := "data.aci_dhcp_option_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDHCPOptionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateDHCPOptionPolicyDSWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateDHCPOptionPolicyDSWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDHCPOptionPolicyConfigDataSource(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccDHCPOptionPolicyDataSourceUpdate(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccDHCPOptionPolicyDSWithInvalidName(fvTenantName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccDHCPOptionPolicyDataSourceUpdatedResource(fvTenantName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccDHCPOptionPolicyConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_option_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_dhcp_option_policy.test.name
		depends_on = [ aci_dhcp_option_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateDHCPOptionPolicyDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing dhcp_option_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_dhcp_option_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = aci_dhcp_option_policy.test.name
		depends_on = [ aci_dhcp_option_policy.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = aci_dhcp_option_policy.test.name
		depends_on = [ aci_dhcp_option_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccDHCPOptionPolicyDSWithInvalidName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_option_policy Data Source with Invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_dhcp_option_policy.test.name}_invalid"
		depends_on = [ aci_dhcp_option_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccDHCPOptionPolicyDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing dhcp_option_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_dhcp_option_policy.test.name
		%s = "%s"
		depends_on = [ aci_dhcp_option_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccDHCPOptionPolicyDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing dhcp_option_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_dhcp_option_policy.test.name
		depends_on = [ aci_dhcp_option_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
