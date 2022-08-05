package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outPathAttachmentDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_path_attachment.test"
	dataSourceName := "data.aci_l3out_path_attachment.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outPathAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outPathAttachmentDSWithoutRequired(rName, rName, rName, rName, pathEp5, "logical_interface_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outPathAttachmentDSWithoutRequired(rName, rName, rName, rName, pathEp5, "target_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outPathAttachmentConfigDataSource(rName, rName, rName, rName, pathEp5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "logical_interface_profile_dn", resourceName, "logical_interface_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "target_dn", resourceName, "target_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "addr", resourceName, "addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "autostate", resourceName, "autostate"),
					resource.TestCheckResourceAttrPair(dataSourceName, "encap", resourceName, "encap"),
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
				Config:      CreateAccL3outPathAttachmentDataSourceUpdate(rName, rName, rName, rName, pathEp5, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3outPathAttachmentDSWithInvalidParentDn(rName, rName, rName, rName, pathEp5),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccL3outPathAttachmentDataSourceUpdatedResource(rName, rName, rName, rName, pathEp5, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL3outPathAttachmentConfigDataSource(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment Data Source with required arguments only")
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
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"
	}

	data "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = aci_l3out_path_attachment.test.target_dn
		depends_on = [ aci_l3out_path_attachment.test ]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn)
	return resource
}

func CreateL3outPathAttachmentDSWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_path_attachment creation without ", attrName)
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
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"
	}
	`
	switch attrName {
	case "logical_interface_profile_dn":
		rBlock += `
	data "aci_l3out_path_attachment" "test" {
	#	logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = aci_l3out_path_attachment.test.target_dn
		depends_on = [ aci_l3out_path_attachment.test ]
	}
		`
	case "target_dn":
		rBlock += `
	data "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	#	target_dn  = aci_l3out_path_attachment.test.target_dn
		depends_on = [ aci_l3out_path_attachment.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn)
}

func CreateAccL3outPathAttachmentDSWithInvalidParentDn(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment Data Source with Invalid Parent Dn")
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
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"
	}

	data "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = "${aci_logical_interface_profile.test.id}_invalid"
		target_dn  = aci_l3out_path_attachment.test.target_dn
		depends_on = [ aci_l3out_path_attachment.test ]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn)
	return resource
}

func CreateAccL3outPathAttachmentDataSourceUpdate(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, key, value string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment Data Source with random attribute")
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
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"
	}

	data "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = aci_l3out_path_attachment.test.target_dn
		if_inst_t = "ext-svi"
		%s = "%s"
		depends_on = [ aci_l3out_path_attachment.test ]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, key, value)
	return resource
}

func CreateAccL3outPathAttachmentDataSourceUpdatedResource(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, key, value string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment Data Source with updated resource")
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
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = "%s"
		if_inst_t = "ext-svi"
		%s = "%s"
	}

	data "aci_l3out_path_attachment" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		target_dn  = aci_l3out_path_attachment.test.target_dn
		depends_on = [ aci_l3out_path_attachment.test ]
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, tDn, key, value)
	return resource
}
