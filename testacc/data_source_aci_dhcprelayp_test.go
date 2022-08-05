package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciDHCPRelayPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_dhcp_relay_policy.test"
	dataSourceName := "data.aci_dhcp_relay_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDHCPRelayPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateDHCPRelayPolicyDSWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateDHCPRelayPolicyDSWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDHCPRelayPolicyConfigDataSource(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mode", resourceName, "mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "owner", resourceName, "owner"),
				),
			},
			{
				Config:      CreateAccDHCPRelayPolicyDataSourceUpdate(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccDHCPRelayPolicyDSWithInvalidParentDn(fvTenantName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccDHCPRelayPolicyDataSourceUpdatedResource(fvTenantName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccDHCPRelayPolicyConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_relay_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_dhcp_relay_policy.test.name
		depends_on = [ aci_dhcp_relay_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateDHCPRelayPolicyDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing dhcp_relay_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_dhcp_relay_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = aci_dhcp_relay_policy.test.name
		depends_on = [ aci_dhcp_relay_policy.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = aci_dhcp_relay_policy.test.name
		depends_on = [ aci_dhcp_relay_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccDHCPRelayPolicyDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_relay_policy Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_dhcp_relay_policy.test.name}_invalid"
		depends_on = [ aci_dhcp_relay_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccDHCPRelayPolicyDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing dhcp_relay_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_dhcp_relay_policy.test.name
		%s = "%s"
		depends_on = [ aci_dhcp_relay_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccDHCPRelayPolicyDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing dhcp_relay_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_dhcp_relay_policy.test.name
		depends_on = [ aci_dhcp_relay_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
