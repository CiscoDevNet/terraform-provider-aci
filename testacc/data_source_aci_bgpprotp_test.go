package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outBGPProtocolProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_bgp_protocol_profile.test"
	dataSourceName := "data.aci_l3out_bgp_protocol_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outBGPProtocolProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outBGPProtocolProfileDSWithoutRequired(rName, rName, rName, "logical_node_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outBGPProtocolProfileConfigDataSource(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "logical_node_profile_dn", resourceName, "logical_node_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccL3outBGPProtocolProfileDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccL3outBGPProtocolProfileDSWithInvalidParentDn(rName, rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3outBGPProtocolProfileDataSourceUpdate(rName, rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateL3outBGPProtocolProfileDSWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, attribute string) string {
	fmt.Println("=== STEP  testing l3out_bgp_protocol_profile data source without required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
	}

	data "aci_l3out_bgp_protocol_profile" "test" {
		depends_on = [
			aci_l3out_bgp_protocol_profile.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName)
	return resource
}

func CreateAccL3outBGPProtocolProfileConfigDataSource(fvTenantName, l3extOutName, l3extLNodePName string) string {
	fmt.Println("=== STEP  testing l3out_bgp_protocol_profile data source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
	}

	data "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		depends_on = [
			aci_l3out_bgp_protocol_profile.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName)
	return resource
}
func CreateAccL3outBGPProtocolProfileDSWithInvalidParentDn(fvTenantName, l3extOutName, l3extLNodePName string) string {
	fmt.Println("=== STEP  testing l3out_bgp_protocol_profile data source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
	}

	data "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = "${aci_logical_node_profile.test.id}invalid"
		depends_on = [
			aci_l3out_bgp_protocol_profile.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName)
	return resource
}

func CreateAccL3outBGPProtocolProfileDataSourceUpdate(fvTenantName, l3extOutName, l3extLNodePName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_bgp_protocol_profile data source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		%s = "%s"

	}

	data "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		depends_on = [
			aci_l3out_bgp_protocol_profile.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, key, value)
	return resource
}

func CreateAccL3outBGPProtocolProfileDataSourceUpdateRandomAttr(rName, key, value string) string {
	fmt.Println("=== STEP  testing L3out BGP Protocol Profile data source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
	}

	data "aci_l3out_bgp_protocol_profile" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		%s = "%s"
		depends_on = [
			aci_l3out_bgp_protocol_profile.test
		]
	}
	`, rName, rName, rName, key, value)
	return resource
}
