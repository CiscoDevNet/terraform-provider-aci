package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciBGPPeerPrefixDataSource_Basic(t *testing.T) {
	resourceName := "aci_bgp_peer_prefix.test"
	dataSourceName := "data.aci_bgp_peer_prefix.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPPeerPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBGPPeerPrefixDSWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBGPPeerPrefixDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBGPPeerPrefixConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "action", resourceName, "action"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_pfx", resourceName, "max_pfx"),
					resource.TestCheckResourceAttrPair(dataSourceName, "restart_time", resourceName, "restart_time"),
					resource.TestCheckResourceAttrPair(dataSourceName, "thresh", resourceName, "thresh"),
				),
			},
			{
				Config:      CreateAccBGPPeerPrefixDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccBGPPeerPrefixDSWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccBGPPeerPrefixDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccBGPPeerPrefixConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_peer_prefix Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_peer_prefix.test.name
		depends_on = [ aci_bgp_peer_prefix.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateBGPPeerPrefixDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_peer_prefix Data Source without ", attrName)
	rBlock := `
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_bgp_peer_prefix" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
		depends_on = [ aci_bgp_peer_prefix.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
		depends_on = [ aci_bgp_peer_prefix.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBGPPeerPrefixDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_peer_prefix Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_peer_prefix" "test" {
		tenant_dn  = "${aci_tenant.test.id}_invalid"
		name  = aci_bgp_peer_prefix.test.name
		depends_on = [ aci_bgp_peer_prefix.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBGPPeerPrefixDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bgp_peer_prefix Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_peer_prefix.test.name
		%s = "%s"
		depends_on = [ aci_bgp_peer_prefix.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccBGPPeerPrefixDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bgp_peer_prefix Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_bgp_peer_prefix" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_peer_prefix.test.name
		depends_on = [ aci_bgp_peer_prefix.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
