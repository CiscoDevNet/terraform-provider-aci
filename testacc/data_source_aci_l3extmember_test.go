package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outVPCMemberDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_vpc_member.test"
	dataSourceName := "data.aci_l3out_vpc_member.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outVPCMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outVPCMemberDSWithoutRequired(rName, "leaf_port_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outVPCMemberDSWithoutRequired(rName, "side"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outVPCMemberConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "leaf_port_dn", resourceName, "leaf_port_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "side", resourceName, "side"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "addr", resourceName, "addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ipv6_dad", resourceName, "ipv6_dad"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ll_addr", resourceName, "ll_addr"),
				),
			},
			{
				Config:      CreateAccL3outVPCMemberDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccL3outVPCMemberDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3outVPCMemberDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL3outVPCMemberConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing l3out_vpc_member Data Source with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_logical_node_profile" "test"{
		l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
	}
	resource "aci_logical_interface_profile" "test"{
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = "%s"
	}
	resource "aci_l3out_path_attachment" "test"{
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
		target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
		if_inst_t = "ext-svi"
	}
	resource "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_path_attachment.test.id
		side  = "A"
	}
	data "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_vpc_member.test.leaf_port_dn
		side  = aci_l3out_vpc_member.test.side
		depends_on = [ aci_l3out_vpc_member.test ]
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateL3outVPCMemberDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_vpc_member Data Source without ", attrName)
	rBlock := `
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_logical_node_profile" "test"{
		l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
	}
	resource "aci_logical_interface_profile" "test"{
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = "%s"
	}
	resource "aci_l3out_path_attachment" "test"{
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
		target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
		if_inst_t = "ext-svi"
	}
	resource "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_path_attachment.test.id
		side  = "A"
	}
	`
	switch attrName {
	case "side":
		rBlock += `
	data "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_vpc_member.test.leaf_port_dn
	#	side  = aci_l3out_vpc_member.test.side
		depends_on = [ aci_l3out_vpc_member.test ]
	}
		`
	case "leaf_port_dn":
		rBlock += `
	data "aci_l3out_vpc_member" "test" {
	#	leaf_port_dn = aci_l3out_vpc_member.test.leaf_port_dn
		side  = aci_l3out_vpc_member.test.side
		depends_on = [ aci_l3out_vpc_member.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, rName, rName, rName)
}

func CreateAccL3outVPCMemberDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing l3out_vpc_member Data Source with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_logical_node_profile" "test"{
		l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
	}
	resource "aci_logical_interface_profile" "test"{
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = "%s"
	}
	resource "aci_l3out_path_attachment" "test"{
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
		target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
		if_inst_t = "ext-svi"
	}
	resource "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_path_attachment.test.id
		side  = "A"
	}
	data "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_vpc_member.test.leaf_port_dn
		side  = aci_l3out_vpc_member.test.side
		depends_on = [ aci_l3out_vpc_member.test ]
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccL3outVPCMemberDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_vpc_member Data Source with random attribute")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_logical_node_profile" "test"{
		l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
	}
	resource "aci_logical_interface_profile" "test"{
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = "%s"
	}
	resource "aci_l3out_path_attachment" "test"{
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
		target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
		if_inst_t = "ext-svi"
	}
	resource "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_path_attachment.test.id
		side  = "A"
	}
	data "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_vpc_member.test.leaf_port_dn
		side  = aci_l3out_vpc_member.test.side
		%s = "%s"
		depends_on = [ aci_l3out_vpc_member.test ]
	}
	`, rName, rName, rName, rName, key, value)
	return resource
}

func CreateAccL3outVPCMemberDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_vpc_member Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l3_outside" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_logical_node_profile" "test"{
		l3_outside_dn = aci_l3_outside.test.id
		name = "%s"
	}
	resource "aci_logical_interface_profile" "test"{
		logical_node_profile_dn = aci_logical_node_profile.test.id
		name = "%s"
	}
	resource "aci_l3out_path_attachment" "test"{
		logical_interface_profile_dn = aci_logical_interface_profile.test.id
		target_dn = "topology/pod-1/paths-101/pathep-[eth1/1]"
		if_inst_t = "ext-svi"
	}
	resource "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_path_attachment.test.id
		side  = "A"
		%s = "%s"
	}
	data "aci_l3out_vpc_member" "test" {
		leaf_port_dn = aci_l3out_vpc_member.test.leaf_port_dn
		side  = aci_l3out_vpc_member.test.side
		depends_on = [ aci_l3out_vpc_member.test ]
	}
	`, rName, rName, rName, rName, key, value)
	return resource
}
