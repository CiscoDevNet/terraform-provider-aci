package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciNodeMgmtEpgDataSource_Basic(t *testing.T) {
	resourceName := "aci_node_mgmt_epg.test"
	dataSourceName := "data.aci_node_mgmt_epg.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	nodeType := "in_band"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeMgmtEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateNodeMgmtEpgDSWithoutRequired(nodeType, rName, "type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateNodeMgmtEpgDSWithoutRequired(nodeType, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccNodeMgmtEpgConfigDataSource(nodeType, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "management_profile_dn", resourceName, "management_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "encap", resourceName, "encap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "flood_on_encap", resourceName, "flood_on_encap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_t", resourceName, "match_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pref_gr_memb", resourceName, "pref_gr_memb"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
				),
			},
			{
				Config:      CreateAccNodeMgmtEpgDataSourceUpdate(nodeType, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccNodeMgmtEpgDSWithInvalidParentDn(nodeType, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccNodeMgmtEpgDataSourceUpdatedResource(nodeType, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccNodeMgmtEpgConfigDataSource(nodeType, rName string) string {
	fmt.Println("=== STEP  testing node_mgmt_epg Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type  = "%s"
		name  = "%s"
	}

	data "aci_node_mgmt_epg" "test" {
		management_profile_dn  = aci_node_mgmt_epg.test.management_profile_dn
		name  = aci_node_mgmt_epg.test.name
		type  = aci_node_mgmt_epg.test.type
		depends_on = [ aci_node_mgmt_epg.test ]
	}
	`, nodeType, rName)
	return resource
}

func CreateNodeMgmtEpgDSWithoutRequired(nodeType, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing node_mgmt_epg Data Source without ", attrName)
	rBlock := `
	
	resource "aci_node_mgmt_epg" "test" {
		type = "%s"
		name  = "%s"
	}
	`
	switch attrName {
	case "type":
		rBlock += `
	data "aci_node_mgmt_epg" "test" {
		management_profile_dn  = aci_node_mgmt_epg.test.management_profile_dn
	#	type = aci_node_mgmt_epg.test.type
		name  = aci_node_mgmt_epg.test.name
		depends_on = [ aci_node_mgmt_epg.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_node_mgmt_epg" "test" {
		management_profile_dn  = aci_node_mgmt_epg.test.management_profile_dn
		type = aci_node_mgmt_epg.test.type
	#	name  = aci_node_mgmt_epg.test.name
		depends_on = [ aci_node_mgmt_epg.test ]
	}
		`
	case "management_profile_dn":
		rBlock += `
	data "aci_node_mgmt_epg" "test" {
	#	management_profile_dn  = aci_node_mgmt_epg.test.management_profile_dn
		type = aci_node_mgmt_epg.test.type
		name  = aci_node_mgmt_epg.test.name
		depends_on = [ aci_node_mgmt_epg.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, nodeType, rName)
}

func CreateAccNodeMgmtEpgDSWithInvalidParentDn(nodeType, rName string) string {
	fmt.Println("=== STEP  testing node_mgmt_epg Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`

	resource "aci_node_mgmt_epg" "test" {
		type  = "%s"
		name  = "%s"
	}

	data "aci_node_mgmt_epg" "test" {
		management_profile_dn  = aci_node_mgmt_epg.test.management_profile_dn
		name  = "${aci_node_mgmt_epg.test.name}_invalid"
		type  = aci_node_mgmt_epg.test.type
		depends_on = [ aci_node_mgmt_epg.test ]
	}
	`, nodeType, rName)
	return resource
}

func CreateAccNodeMgmtEpgDataSourceUpdate(nodeType, rName, key, value string) string {
	fmt.Println("=== STEP  testing node_mgmt_epg Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type  = "%s"
		name  = "%s"
	}

	data "aci_node_mgmt_epg" "test" {
		management_profile_dn  =  aci_node_mgmt_epg.test.management_profile_dn
		name  = aci_node_mgmt_epg.test.name
		type  = aci_node_mgmt_epg.test.type
		%s = "%s"
		depends_on = [ aci_node_mgmt_epg.test ]
	}
	`, nodeType, rName, key, value)
	return resource
}

func CreateAccNodeMgmtEpgDataSourceUpdatedResource(nodeType, rName, key, value string) string {
	fmt.Println("=== STEP  testing node_mgmt_epg Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_node_mgmt_epg" "test" {
		type  = "%s"
		name  = "%s"
		%s = "%s"
	}

	data "aci_node_mgmt_epg" "test" {
		management_profile_dn  = aci_node_mgmt_epg.test.management_profile_dn
		name  = aci_node_mgmt_epg.test.name
		type  = aci_node_mgmt_epg.test.type
		depends_on = [ aci_node_mgmt_epg.test ]
	}
	`, nodeType, rName, key, value)
	return resource
}
