package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outBGPExternalPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_bgp_external_policy.test"
	dataSourceName := "data.aci_l3out_bgp_external_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	fvTenantName := makeTestVariable(acctest.RandString(5))
	l3extOutName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outBGPExternalPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outBGPExternalPolicyDSWithoutRequired(fvTenantName, l3extOutName, "l3_outside_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outBGPExternalPolicyConfigDataSource(fvTenantName, l3extOutName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "l3_outside_dn", resourceName, "l3_outside_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccL3outBGPExternalPolicyDataSourceUpdate(fvTenantName, l3extOutName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3outBGPExternalPolicyDSWithInvalidParentDn(fvTenantName, l3extOutName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccL3outBGPExternalPolicyDataSourceUpdatedResource(fvTenantName, l3extOutName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL3outBGPExternalPolicyConfigDataSource(fvTenantName, l3extOutName string) string {
	fmt.Println("=== STEP  testing l3out_bgp_external_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}

	data "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		depends_on = [ aci_l3out_bgp_external_policy.test ]
	}
	`, fvTenantName, l3extOutName)
	return resource
}

func CreateL3outBGPExternalPolicyDSWithoutRequired(fvTenantName, l3extOutName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_bgp_external_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}
	`
	switch attrName {
	case "l3_outside_dn":
		rBlock += `
	data "aci_l3out_bgp_external_policy" "test" {
	#	l3_outside_dn  = aci_l3_outside.test.id
		depends_on = [ aci_l3out_bgp_external_policy.test ]
	}
	`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName)
}

func CreateAccL3outBGPExternalPolicyDSWithInvalidParentDn(fvTenantName, l3extOutName string) string {
	fmt.Println("=== STEP  testing l3out_bgp_external_policy Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}

	data "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = "${aci_l3_outside.test.id}_invalid"
		depends_on = [ aci_l3out_bgp_external_policy.test ]
	}
	`, fvTenantName, l3extOutName)
	return resource
}

func CreateAccL3outBGPExternalPolicyDataSourceUpdate(fvTenantName, l3extOutName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_bgp_external_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}

	data "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		%s = "%s"
		depends_on = [ aci_l3out_bgp_external_policy.test ]
	}
	`, fvTenantName, l3extOutName, key, value)
	return resource
}

func CreateAccL3outBGPExternalPolicyDataSourceUpdatedResource(fvTenantName, l3extOutName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_bgp_external_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		%s = "%s"
	}

	data "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		depends_on = [ aci_l3out_bgp_external_policy.test ]
	}
	`, fvTenantName, l3extOutName, key, value)
	return resource
}
