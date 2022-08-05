package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciBGPAddressFamilyContextDataSource_Basic(t *testing.T) {
	resourceName := "aci_bgp_address_family_context.test"
	dataSourceName := "data.aci_bgp_address_family_context.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPAddressFamilyContextPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBGPAddressFamilyContextDSWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBGPAddressFamilyContextDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBGPAddressFamilyContextConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl", resourceName, "ctrl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "e_dist", resourceName, "e_dist"),
					resource.TestCheckResourceAttrPair(dataSourceName, "i_dist", resourceName, "i_dist"),
					resource.TestCheckResourceAttrPair(dataSourceName, "local_dist", resourceName, "local_dist"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_ecmp", resourceName, "max_ecmp"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_ecmp_ibgp", resourceName, "max_ecmp_ibgp"),
				),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccBGPAddressFamilyContextDSWithInvalidName(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccBGPAddressFamilyContextDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccBGPAddressFamilyContextConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_address_family_context Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_address_family_context.test.name
		depends_on = [ aci_bgp_address_family_context.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateBGPAddressFamilyContextDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_address_family_context Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_bgp_address_family_context" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
		depends_on = [ aci_bgp_address_family_context.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
		depends_on = [ aci_bgp_address_family_context.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBGPAddressFamilyContextDSWithInvalidName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_address_family_context Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_bgp_address_family_context.test.name}_invalid"
		depends_on = [ aci_bgp_address_family_context.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBGPAddressFamilyContextDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bgp_address_family_context Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_address_family_context.test.name
		%s = "%s"
		depends_on = [ aci_bgp_address_family_context.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccBGPAddressFamilyContextDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bgp_address_family_context Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_address_family_context.test.name
		depends_on = [ aci_bgp_address_family_context.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
