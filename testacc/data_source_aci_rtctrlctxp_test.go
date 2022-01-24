package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciRouteControlContextDataSource_Basic(t *testing.T) {
	resourceName := "aci_route_control_context.test"
	dataSourceName := "data.aci_route_control_context.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	rtctrlProfileName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRouteControlContextDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateRouteControlContextDSWithoutRequired(fvTenantName, rtctrlProfileName, rName, "route_control_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateRouteControlContextDSWithoutRequired(fvTenantName, rtctrlProfileName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRouteControlContextConfigDataSource(fvTenantName, rtctrlProfileName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "route_control_profile_dn", resourceName, "route_control_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "action", resourceName, "action"),
					resource.TestCheckResourceAttrPair(dataSourceName, "order", resourceName, "order"),
				),
			},
			{
				Config:      CreateAccRouteControlContextDataSourceUpdate(fvTenantName, rtctrlProfileName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccRouteControlContextDSWithInvalidParentDn(fvTenantName, rtctrlProfileName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccRouteControlContextDataSourceUpdatedResource(fvTenantName, rtctrlProfileName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccRouteControlContextConfigDataSource(fvTenantName, rtctrlProfileName, rName string) string {
	fmt.Println("=== STEP  testing route_control_context Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		name 		= "%s"
		parent_dn = aci_tenant.test.id
	}
	
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = "%s"
	}

	data "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = aci_route_control_context.test.name
		depends_on = [ aci_route_control_context.test ]
	}
	`, fvTenantName, rtctrlProfileName, rName)
	return resource
}

func CreateRouteControlContextDSWithoutRequired(fvTenantName, rtctrlProfileName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing route_control_context Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		name 		= "%s"
		parent_dn = aci_tenant.test.id
	}
	
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "route_control_profile_dn":
		rBlock += `
	data "aci_route_control_context" "test" {
	#	route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = aci_route_control_context.test.name
		depends_on = [ aci_route_control_context.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
	#	name  = aci_route_control_context.test.name
		depends_on = [ aci_route_control_context.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rtctrlProfileName, rName)
}

func CreateAccRouteControlContextDSWithInvalidParentDn(fvTenantName, rtctrlProfileName, rName string) string {
	fmt.Println("=== STEP  testing route_control_context Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		name 		= "%s"
		parent_dn = aci_tenant.test.id
	}
	
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = "%s"
	}

	data "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = "${aci_route_control_context.test.name}_invalid"
		depends_on = [ aci_route_control_context.test ]
	}
	`, fvTenantName, rtctrlProfileName, rName)
	return resource
}

func CreateAccRouteControlContextDataSourceUpdate(fvTenantName, rtctrlProfileName, rName, key, value string) string {
	fmt.Println("=== STEP  testing route_control_context Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		name 		= "%s"
		parent_dn = aci_tenant.test.id
	}
	
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = "%s"
	}

	data "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = aci_route_control_context.test.name
		%s = "%s"
		depends_on = [ aci_route_control_context.test ]
	}
	`, fvTenantName, rtctrlProfileName, rName, key, value)
	return resource
}

func CreateAccRouteControlContextDataSourceUpdatedResource(fvTenantName, rtctrlProfileName, rName, key, value string) string {
	fmt.Println("=== STEP  testing route_control_context Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_control_profile" "test" {
		name 		= "%s"
		parent_dn = aci_tenant.test.id
	}
	
	resource "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_route_control_context" "test" {
		route_control_profile_dn  = aci_bgp_route_control_profile.test.id
		name  = aci_route_control_context.test.name
		depends_on = [ aci_route_control_context.test ]
	}
	`, fvTenantName, rtctrlProfileName, rName, key, value)
	return resource
}
