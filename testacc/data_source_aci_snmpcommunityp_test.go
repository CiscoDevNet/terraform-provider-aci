package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciVRFSnmpContextCommunityDataSource_Basic(t *testing.T) {
	resourceName := "aci_vrf_snmp_context_community.test"
	dataSourceName := "data.aci_vrf_snmp_context_community.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	fvCtxName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVRFSnmpContextCommunityDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVRFSnmpContextCommunityDSWithoutRequired(fvTenantName, fvCtxName, rName, "vrf_snmp_context_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVRFSnmpContextCommunityDSWithoutRequired(fvTenantName, fvCtxName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVRFSnmpContextCommunityConfigDataSource(fvTenantName, fvCtxName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "vrf_snmp_context_dn", resourceName, "vrf_snmp_context_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccVRFSnmpContextCommunityDataSourceUpdate(fvTenantName, fvCtxName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccVRFSnmpContextCommunityDSWithInvalidParentDn(fvTenantName, fvCtxName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccVRFSnmpContextCommunityDataSourceUpdatedResource(fvTenantName, fvCtxName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccVRFSnmpContextCommunityConfigDataSource(fvTenantName, fvCtxName, rName string) string {
	fmt.Println("=== STEP  testing vrf_snmp_context_community Data Source with required arguments only")
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
		name  = "%s"
	}
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn   = aci_vrf_snmp_context.test.id
		name  = "%s"
	}

	data "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn   = aci_vrf_snmp_context.test.id
		name  = aci_vrf_snmp_context_community.test.name
		depends_on = [ aci_vrf_snmp_context_community.test ]
	}
	`, fvTenantName, fvCtxName, fvCtxName, rName)
	return resource
}

func CreateVRFSnmpContextCommunityDSWithoutRequired(fvTenantName, fvCtxName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vrf_snmp_context_community Data Source without ", attrName)
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
		name  = "%s"
	}
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn   = aci_vrf_snmp_context.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "vrf_snmp_context_dn":
		rBlock += `

	data "aci_vrf_snmp_context_community" "test" {
	#	vrf_snmp_context_dn   = aci_vrf_snmp_context.test.id
		name  = aci_vrf_snmp_context_community.test.name
		depends_on = [ aci_vrf_snmp_context_community.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn  = aci_vrf_snmp_context.test.id
	#	name  = aci_vrf_snmp_context_community.test.name
		depends_on = [ aci_vrf_snmp_context_community.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, fvCtxName, fvCtxName, rName)
}

func CreateAccVRFSnmpContextCommunityDSWithInvalidParentDn(fvTenantName, fvCtxName, rName string) string {
	fmt.Println("=== STEP  testing vrf_snmp_context_community Data Source with Invalid Parent Dn")
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
		name  = "%s"
	}
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn   = aci_vrf_snmp_context.test.id
		name  = "%s"
	}

	data "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn  = aci_vrf_snmp_context.test.id
		name  = "${aci_vrf_snmp_context_community.test.name}_invalid"
		depends_on = [ aci_vrf_snmp_context_community.test ]
	}
	`, fvTenantName, fvCtxName, fvCtxName, rName)
	return resource
}

func CreateAccVRFSnmpContextCommunityDataSourceUpdate(fvTenantName, fvCtxName, rName, key, value string) string {
	fmt.Println("=== STEP  testing vrf_snmp_context_community Data Source with random attribute")
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
		name  = "%s"
	}
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn   = aci_vrf_snmp_context.test.id
		name  = "%s"
	}

	data "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn  = aci_vrf_snmp_context.test.id
		name  = aci_vrf_snmp_context_community.test.name
		%s = "%s"
		depends_on = [ aci_vrf_snmp_context_community.test ]
	}
	`, fvTenantName, fvCtxName, fvCtxName, rName, key, value)
	return resource
}

func CreateAccVRFSnmpContextCommunityDataSourceUpdatedResource(fvTenantName, fvCtxName, rName, key, value string) string {
	fmt.Println("=== STEP  testing vrf_snmp_context_community Data Source with updated resource")
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
		name  = "%s"
	}
	
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn   = aci_vrf_snmp_context.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn  = aci_vrf_snmp_context.test.id
		name  = aci_vrf_snmp_context_community.test.name
		depends_on = [ aci_vrf_snmp_context_community.test ]
	}
	`, fvTenantName, fvCtxName, fvCtxName, rName, key, value)
	return resource
}
