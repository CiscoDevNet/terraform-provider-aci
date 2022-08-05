package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFabricNodeDataSource_Basic(t *testing.T) {
	resourceName := "aci_logical_node_to_fabric_node.test"
	dataSourceName := "data.aci_logical_node_to_fabric_node.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	rtrid, _ := acctest.RandIpAddress("10.1.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateFabricNodeDSWithoutRequired(rName, rName, rName, fabDn1, "logical_node_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateFabricNodeDSWithoutRequired(rName, rName, rName, fabDn1, "tdn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricNodeConfigDataSource(rName, rName, rName, fabDn1, rtrid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "logical_node_profile_dn", resourceName, "logical_node_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tdn", resourceName, "tdn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "config_issues", resourceName, "config_issues"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rtr_id", resourceName, "rtr_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rtr_id_loop_back", resourceName, "rtr_id_loop_back"),
				),
			},
			{
				Config:      CreateAccFabricNodeDataSourceUpdate(rName, rName, rName, fabDn1, rtrid, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccFabricNodeDSWithInvalidParentDn(rName, rName, rName, fabDn1, rtrid),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccFabricNodeDataSourceUpdatedResource(rName, rName, rName, fabDn1, rtrid, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFabricNodeConfigDataSource(fvTenantName, l3extOutName, l3extLNodePName, tDn, ip string) string {
	fmt.Println("=== STEP  testing logical_node_to_fabric_node Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = "%s"
		rtr_id = "%s"
	}

	data "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = aci_logical_node_to_fabric_node.test.tdn
		depends_on = [ aci_logical_node_to_fabric_node.test ]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, tDn, ip)
	return resource
}

func CreateFabricNodeDSWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, tDn, attrName string) string {
	fmt.Println("=== STEP  Basic: testing logical_node_to_fabric_node Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = "%s"
		rtr_id = "10.1.0.0"
	}
	`
	switch attrName {
	case "logical_node_profile_dn":
		rBlock += `
	data "aci_logical_node_to_fabric_node" "test" {
	#	logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = aci_logical_node_to_fabric_node.test.tdn
		depends_on = [ aci_logical_node_to_fabric_node.test ]
	}
		`
	case "tdn":
		rBlock += `
	data "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
	#	tdn  = aci_logical_node_to_fabric_node.test.tdn
		depends_on = [ aci_logical_node_to_fabric_node.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName, tDn)
}

func CreateAccFabricNodeDSWithInvalidParentDn(fvTenantName, l3extOutName, l3extLNodePName, tDn, ip string) string {
	fmt.Println("=== STEP  testing logical_node_to_fabric_node Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = "%s"
		rtr_id = "%s"
	}

	data "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = "${aci_logical_node_profile.test.id}_invalie"
		tdn  = aci_logical_node_to_fabric_node.test.tdn
		depends_on = [ aci_logical_node_to_fabric_node.test ]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, tDn, ip)
	return resource
}

func CreateAccFabricNodeDataSourceUpdate(fvTenantName, l3extOutName, l3extLNodePName, tDn, ip, key, value string) string {
	fmt.Println("=== STEP  testing logical_node_to_fabric_node Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = "%s"
		rtr_id = "%s"
	}

	data "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = aci_logical_node_to_fabric_node.test.tdn
		%s = "%s"
		depends_on = [ aci_logical_node_to_fabric_node.test ]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, tDn, ip, key, value)
	return resource
}

func CreateAccFabricNodeDataSourceUpdatedResource(fvTenantName, l3extOutName, l3extLNodePName, tDn, ip, key, value string) string {
	fmt.Println("=== STEP  testing logical_node_to_fabric_node Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = "%s"
		rtr_id = "%s"
		%s = "%s"
	}

	data "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = aci_logical_node_to_fabric_node.test.tdn
		depends_on = [ aci_logical_node_to_fabric_node.test ]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, tDn, ip, key, value)
	return resource
}
