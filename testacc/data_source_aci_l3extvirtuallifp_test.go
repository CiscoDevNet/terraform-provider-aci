package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outFloatingSVIDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_floating_svi.test"
	dataSourceName := "data.aci_l3out_floating_svi.test"
	node_dn := "topology/pod-1/node-111"
	encap := "vlan-20"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outFloatingSVIDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccL3outFloatingSVIConfigDataSourceWithoutParentDn(rName, rName, rName, rName, node_dn, encap),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccL3outFloatingSVIConfigDataSourceWithoutEncap(rName, rName, rName, rName, node_dn, encap),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config:      CreateAccL3outFloatingSVIConfigDataSourceWithoutNodeDn(rName, rName, rName, rName, node_dn, encap),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outFloatingSVIConfigDataSource(rName, rName, rName, rName, node_dn, encap),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "logical_interface_profile_dn", resourceName, "logical_interface_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "addr", resourceName, "addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "autostate", resourceName, "autostate"),
					resource.TestCheckResourceAttrPair(dataSourceName, "encap_scope", resourceName, "encap_scope"),
					resource.TestCheckResourceAttrPair(dataSourceName, "if_inst_t", resourceName, "if_inst_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ipv6_dad", resourceName, "ipv6_dad"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ll_addr", resourceName, "ll_addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mac", resourceName, "mac"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mode", resourceName, "mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mtu", resourceName, "mtu"),
					resource.TestCheckResourceAttrPair(dataSourceName, "target_dscp", resourceName, "target_dscp"),
				),
			},
			{
				Config:      CreateAccL3outFloatingSVIDataSourceUpdate(rName, rName, rName, rName, node_dn, encap, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3outFloatingSVIDSWithInvalidParentDn(rName, rName, rName, rName, node_dn, encap),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3outFloatingSVIDataSourceUpdate(rName, rName, rName, rName, node_dn, encap, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL3outFloatingSVIConfigDataSource(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, node_dn, encap string) string {
	fmt.Println("=== STEP  testing l3out_floating_svi creation with required arguments only")
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
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}

	data "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = aci_l3out_floating_svi.test.node_dn
		encap  = aci_l3out_floating_svi.test.encap
		depends_on = [
			aci_l3out_floating_svi.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, node_dn, encap)
	return resource
}

func CreateAccL3outFloatingSVIConfigDataSourceWithoutParentDn(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, node_dn, encap string) string {
	fmt.Println("=== STEP  testing l3out_floating_svi creation without logical_interface_profile_dn")
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
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}

	data "aci_l3out_floating_svi" "test" {
		node_dn  = aci_l3out_floating_svi.test.node_dn
		encap  = aci_l3out_floating_svi.test.encap
		depends_on = [
			aci_l3out_floating_svi.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, node_dn, encap)
	return resource
}

func CreateAccL3outFloatingSVIConfigDataSourceWithoutNodeDn(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, node_dn, encap string) string {
	fmt.Println("=== STEP  testing l3out_floating_svi creation without node_dn")
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
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}

	data "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		encap  = aci_l3out_floating_svi.test.encap
		depends_on = [
			aci_l3out_floating_svi.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, node_dn, encap)
	return resource
}

func CreateAccL3outFloatingSVIConfigDataSourceWithoutEncap(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, node_dn, encap string) string {
	fmt.Println("=== STEP  testing l3out_floating_svi creation without encap")
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
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}

	data "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = aci_l3out_floating_svi.test.node_dn
		depends_on = [
			aci_l3out_floating_svi.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, node_dn, encap)
	return resource
}
func CreateAccL3outFloatingSVIDSWithInvalidParentDn(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, node_dn, encap string) string {
	fmt.Println("=== STEP  testing l3out_floating_svi creation with Invalid Parent Dn")
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
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}

	data "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = "${aci_logical_interface_profile.test.id}invalid"
		node_dn  = aci_l3out_floating_svi.test.node_dn
		encap  = aci_l3out_floating_svi.test.encap
		depends_on = [
			aci_l3out_floating_svi.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, node_dn, encap)
	return resource
}

func CreateAccL3outFloatingSVIDataSourceUpdate(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, node_dn, encap, key, value string) string {
	fmt.Println("=== STEP  testing l3out_floating_svi updation with required arguments only")
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
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}

	data "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = aci_l3out_floating_svi.test.node_dn
		encap  = aci_l3out_floating_svi.test.encap
		%s = "%s"
		depends_on = [
			aci_l3out_floating_svi.test
		]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, node_dn, encap, key, value)
	return resource
}
