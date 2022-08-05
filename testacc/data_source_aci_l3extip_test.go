package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outPathAttachmentSecondaryIpDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_path_attachment_secondary_ip.test"
	dataSourceName := "data.aci_l3out_path_attachment_secondary_ip.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	addr, _ := acctest.RandIpAddress("10.4.0.0/16")
	addr = fmt.Sprintf("%s/16", addr)
	addrUpdated, _ := acctest.RandIpAddress("10.5.0.0/16")
	addrUpdated = fmt.Sprintf("%s/16", addrUpdated)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outPathAttachmentSecondaryIpDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateL3outPathAttachmentSecondaryIpDSWithoutRequired(rName, pathEp1, addr, "addr"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outPathAttachmentSecondaryIpDSWithoutRequired(rName, pathEp1, addr, "l3out_path_attachment_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outPathAttachmentSecondaryIpConfigDataSource(rName, pathEp1, addr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "l3out_path_attachment_dn", resourceName, "l3out_path_attachment_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "addr", resourceName, "addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ipv6_dad", resourceName, "ipv6_dad"),
				),
			},
			{
				Config:      CreateAccL3outPathAttachmentSecondaryIpDataSourceUpdate(rName, pathEp1, addr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3outPathAttachmentSecondaryIpDSWithInvalidIp(rName, pathEp1, addr, addrUpdated),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3outPathAttachmentSecondaryIpDataSourceUpdatedResource(rName, pathEp1, addr, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL3outPathAttachmentSecondaryIpConfigDataSource(name, tdn, addr string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment_secondary_ip Data Source with required arguments only")
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

	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
		addr  = "%s"
	}

	data "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment_secondary_ip.test.l3out_path_attachment_dn
		addr  = aci_l3out_path_attachment_secondary_ip.test.addr
		depends_on = [ aci_l3out_path_attachment_secondary_ip.test ]
	}
	`, name, name, name, name, tdn, addr)
	return resource
}

func CreateL3outPathAttachmentSecondaryIpDSWithoutRequired(name, tdn, addr, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_path_attachment_secondary_ip Data Source without ", attrName)
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

	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
		addr  = "%s"
	}
	`
	switch attrName {
	case "addr":
		rBlock += `
	data "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment_secondary_ip.test.l3out_path_attachment_dn
	#	addr  = aci_l3out_path_attachment_secondary_ip.test.addr
		depends_on = [ aci_l3out_path_attachment_secondary_ip.test ]
	}
		`
	case "l3out_path_attachment_dn":
		rBlock += `
	data "aci_l3out_path_attachment_secondary_ip" "test" {
	#	l3out_path_attachment_dn = aci_l3out_path_attachment_secondary_ip.test.l3out_path_attachment_dn
		addr  = aci_l3out_path_attachment_secondary_ip.test.addr
		depends_on = [ aci_l3out_path_attachment_secondary_ip.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, name, name, name, name, tdn, addr)
}

func CreateAccL3outPathAttachmentSecondaryIpDSWithInvalidIp(name, tdn, addr, addrOther string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment_secondary_ip Data Source with invalid addr")
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

	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
		addr  = "%s"
	}

	data "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment_secondary_ip.test.l3out_path_attachment_dn
		addr  = "%s"
		depends_on = [ aci_l3out_path_attachment_secondary_ip.test ]
	}
	`, name, name, name, name, tdn, addr, addrOther)
	return resource
}

func CreateAccL3outPathAttachmentSecondaryIpDataSourceUpdate(name, tdn, addr, key, value string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment_secondary_ip Data Source with random attribute")
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

	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
		addr  = "%s"
	}

	data "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment_secondary_ip.test.l3out_path_attachment_dn
		addr  = aci_l3out_path_attachment_secondary_ip.test.addr
		%s = "%s"
		depends_on = [ aci_l3out_path_attachment_secondary_ip.test ]
	}
	`, name, name, name, name, tdn, addr, key, value)
	return resource
}

func CreateAccL3outPathAttachmentSecondaryIpDataSourceUpdatedResource(name, tdn, addr, key, value string) string {
	fmt.Println("=== STEP  testing l3out_path_attachment_secondary_ip Data Source with updated resource")
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
	
	resource "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment.test.id
		addr  = "%s"
		%s = "%s"
	}

	data "aci_l3out_path_attachment_secondary_ip" "test" {
		l3out_path_attachment_dn = aci_l3out_path_attachment_secondary_ip.test.l3out_path_attachment_dn
		addr  = aci_l3out_path_attachment_secondary_ip.test.addr
		depends_on = [ aci_l3out_path_attachment_secondary_ip.test ]
	}
	`, name, name, name, name, tdn, addr, key, value)
	return resource
}
