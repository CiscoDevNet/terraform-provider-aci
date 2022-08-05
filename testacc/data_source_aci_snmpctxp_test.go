package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciVRFSnmpContextDataSource_Basic(t *testing.T) {
	resourceName := "aci_vrf_snmp_context.test"
	dataSourceName := "data.aci_vrf_snmp_context.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	fvTenantName := makeTestVariable(acctest.RandString(5))
	fvCtxName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVRFSnmpContextDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVRFSnmpContextDSWithoutRequired(fvTenantName, fvCtxName, "vrf_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVRFSnmpContextConfigDataSource(fvTenantName, fvCtxName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "vrf_dn", resourceName, "vrf_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccVRFSnmpContextDataSourceUpdate(fvTenantName, fvCtxName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccVRFSnmpContextDSWithInvalidParentDn(fvTenantName, fvCtxName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccVRFSnmpContextDataSourceUpdatedResource(fvTenantName, fvCtxName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccVRFSnmpContextConfigDataSource(fvTenantName, fvCtxName string) string {
	fmt.Println("=== STEP  testing vrf_snmp_context Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_vrf.test.id
	}

	data "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_vrf.test.id
		depends_on = [ aci_vrf_snmp_context.test ]
	}
	`, fvTenantName, fvCtxName)
	return resource
}

func CreateVRFSnmpContextDSWithoutRequired(fvTenantName, fvCtxName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vrf_snmp_context Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_vrf.test.id
	}
	`
	switch attrName {
	case "vrf_dn":
		rBlock += `
	data "aci_vrf_snmp_context" "test" {
	#	vrf_dn  = aci_vrf.test.id
	
		depends_on = [ aci_vrf_snmp_context.test ]
	}
		`

	}
	return fmt.Sprintf(rBlock, fvTenantName, fvCtxName)
}

func CreateAccVRFSnmpContextDSWithInvalidParentDn(fvTenantName, fvCtxName string) string {
	fmt.Println("=== STEP  testing vrf_snmp_context Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_vrf.test.id
	}

	data "aci_vrf_snmp_context" "test" {
		vrf_dn  = "${aci_vrf.test.id}invalid"
		depends_on = [ aci_vrf_snmp_context.test ]
	}
	`, fvTenantName, fvCtxName)
	return resource
}

func CreateAccVRFSnmpContextDataSourceUpdate(fvTenantName, fvCtxName, key, value string) string {
	fmt.Println("=== STEP  testing vrf_snmp_context Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_vrf.test.id
	}

	data "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_vrf.test.id
		%s = "%s"
		depends_on = [ aci_vrf_snmp_context.test ]
	}
	`, fvTenantName, fvCtxName, key, value)
	return resource
}

func CreateAccVRFSnmpContextDataSourceUpdatedResource(fvTenantName, fvCtxName, key, value string) string {
	fmt.Println("=== STEP  testing vrf_snmp_context Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_vrf.test.id
		%s = "%s"
	}

	data "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_vrf.test.id
		depends_on = [ aci_vrf_snmp_context.test ]
	}
	`, fvTenantName, fvCtxName, key, value)
	return resource
}
