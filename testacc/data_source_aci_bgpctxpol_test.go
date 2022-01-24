package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciBGPTimersPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_bgp_timers.test"
	dataSourceName := "data.aci_bgp_timers.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPTimersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBGPTimersPolicyDSWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBGPTimersPolicyDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBGPTimersPolicyConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "gr_ctrl", resourceName, "gr_ctrl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hold_intvl", resourceName, "hold_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ka_intvl", resourceName, "ka_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_as_limit", resourceName, "max_as_limit"),
					resource.TestCheckResourceAttrPair(dataSourceName, "stale_intvl", resourceName, "stale_intvl"),
				),
			},
			{
				Config:      CreateAccBGPTimersPolicyDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccBGPTimersPolicyDSWithInvalidName(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccBGPTimersPolicyDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateBGPTimersPolicyDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_timers_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}

	resource "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_bgp_timers" "test" {
	#	tenant_dn  = aci_bgp_timers.test.tenant_dn
		name  = aci_bgp_timers.test.name
	}
		`
	case "name":
		rBlock += `
	data "aci_bgp_timers" "test" {
		tenant_dn  = aci_bgp_timers.test.tenant_dn
	#	name  = aci_bgp_timers.test.name
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBGPTimersPolicyConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_timers_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_timers.test.name
		depends_on = [
			aci_bgp_timers.test
		]
	}
	`, fvTenantName, rName)
	return resource
}
func CreateAccBGPTimersPolicyDSWithInvalidName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_timers_policy Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_bgp_timers.test.name}_invalid"
		depends_on = [
			aci_bgp_timers.test
		]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBGPTimersPolicyDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bgp_timers_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_timers.test.name
		%s = "%s"
		depends_on = [
			aci_bgp_timers.test
		]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccBGPTimersPolicyDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bgp_timers_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_timers.test.name
		depends_on = [
			aci_bgp_timers.test
		]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
