package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciStaticNodeMgmtAddressDataSource_Basic(t *testing.T) {
	resourceName := "aci_static_node_mgmt_address.test"
	dataSourceName := "data.aci_static_node_mgmt_address.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	tDn := "topology/pod-1/node-1"
	nodeType := "in_band"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciStaticNodeMgmtAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateStaticNodeMgmtAddressDSWithoutRequired(nodeType, rName, tDn, "management_epg_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateStaticNodeMgmtAddressDSWithoutRequired(nodeType, rName, tDn, "t_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccStaticNodeMgmtAddressConfigDataSource(nodeType, rName, tDn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "management_epg_dn", resourceName, "management_epg_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "t_dn", resourceName, "t_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "addr", resourceName, "addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "gw", resourceName, "gw"),
					resource.TestCheckResourceAttrPair(dataSourceName, "v6_addr", resourceName, "v6_addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "v6_gw", resourceName, "v6_gw"),
				),
			},
			{
				Config:      CreateAccStaticNodeMgmtAddressDataSourceUpdate(nodeType, rName, tDn, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccStaticNodeMgmtAddressDSWithInvalidParentDn(nodeType, rName, tDn),
				ExpectError: regexp.MustCompile(`Invalid RN`),
			},

			{
				Config: CreateAccStaticNodeMgmtAddressDataSourceUpdatedResource(nodeType, rName, tDn, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccStaticNodeMgmtAddressConfigDataSource(nodeType, mgmtInBName, tDn string) string {
	fmt.Println("=== STEP  testing static_node_mgmt_address Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type = "%s"
		name = "%s"
	}
	
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		type = "%s"
		t_dn  = "%s"
	}

	data "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		t_dn  = aci_static_node_mgmt_address.test.t_dn
		type = aci_static_node_mgmt_address.test.type
		depends_on = [ aci_static_node_mgmt_address.test ]
	}
	`, nodeType, mgmtInBName, nodeType, tDn)
	return resource
}

func CreateStaticNodeMgmtAddressDSWithoutRequired(nodeType, mgmtInBName, tDn, attrName string) string {
	fmt.Println("=== STEP  Basic: testing static_node_mgmt_address Data Source without ", attrName)
	rBlock := `
	
	
	resource "aci_node_mgmt_epg" "test" {
		name = "%s"
		type = "%s"
	}
	
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		t_dn  = "%s"
		type = "%s"
	}
	`
	switch attrName {
	case "management_epg_dn":
		rBlock += `
	data "aci_static_node_mgmt_address" "test" {
	#	management_epg_dn  = aci_node_mgmt_epg.test.id
		type = aci_static_node_mgmt_address.test.type
		t_dn  = aci_static_node_mgmt_address.test.t_dn
		depends_on = [ aci_static_node_mgmt_address.test ]
	}
		`
	case "t_dn":
		rBlock += `
	data "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		type = aci_static_node_mgmt_address.test.type
	#	t_dn  = aci_static_node_mgmt_address.test.t_dn
		depends_on = [ aci_static_node_mgmt_address.test ]
	}
		`
	case "type":
		rBlock += `
	data "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		t_dn  = aci_static_node_mgmt_address.test.t_dn
	#	type = aci_static_node_mgmt_address.test.type	
		depends_on = [ aci_static_node_mgmt_address.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, mgmtInBName, nodeType, tDn, nodeType)
}

func CreateAccStaticNodeMgmtAddressDSWithInvalidParentDn(nodeType, mgmtInBName, tDn string) string {
	fmt.Println("=== STEP  testing static_node_mgmt_address Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		name = "%s"
		type = "%s"
	}
	
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		t_dn  = "%s"
		type = "%s"
	}

	data "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		t_dn  = "${aci_static_node_mgmt_address.test.t_dn}_invalid"
		type = aci_static_node_mgmt_address.test.type
		depends_on = [ aci_static_node_mgmt_address.test ]
	}
	`, mgmtInBName, nodeType, tDn, nodeType)
	return resource
}

func CreateAccStaticNodeMgmtAddressDataSourceUpdate(nodeType, mgmtInBName, tDn, key, value string) string {
	fmt.Println("=== STEP  testing static_node_mgmt_address Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		name = "%s"
		type = "%s"
	}
	
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		t_dn  = "%s"
		type = "%s"
	}

	data "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		t_dn  = aci_static_node_mgmt_address.test.t_dn
		type = aci_static_node_mgmt_address.test.type
		%s = "%s"
		depends_on = [ aci_static_node_mgmt_address.test ]
	}
	`, mgmtInBName, nodeType, tDn, nodeType, key, value)
	return resource
}

func CreateAccStaticNodeMgmtAddressDataSourceUpdatedResource(nodeType, mgmtInBName, tDn, key, value string) string {
	fmt.Println("=== STEP  testing static_node_mgmt_address Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		name = "%s"
		type = "%s"
	}
	
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		t_dn  = "%s"
		type = "%s"
		%s = "%s"
	}

	data "aci_static_node_mgmt_address" "test" {
		management_epg_dn  = aci_node_mgmt_epg.test.id
		t_dn  = aci_static_node_mgmt_address.test.t_dn
		type  = aci_static_node_mgmt_address.test.type
		depends_on = [ aci_static_node_mgmt_address.test ]
	}
	`, mgmtInBName, nodeType, tDn, nodeType, key, value)
	return resource
}
