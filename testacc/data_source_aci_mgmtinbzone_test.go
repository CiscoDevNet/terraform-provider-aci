package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciMgmtZoneDataSource_Basic(t *testing.T) {
	resourceName := "aci_mgmt_zone.test"
	dataSourceName := "data.aci_mgmt_zone.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	mgmtGrpName := makeTestVariable(acctest.RandString(5))
	rName := makeTestVariable(acctest.RandString(5))
	zoneType := "in_band"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMgmtZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateMgmtZoneDSWithoutRequired(mgmtGrpName, zoneType, rName, "managed_node_connectivity_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMgmtZoneDSWithoutRequired(mgmtGrpName, zoneType, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMgmtZoneDSWithoutRequired(mgmtGrpName, zoneType, rName, "type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMgmtZoneConfigDataSource(mgmtGrpName, zoneType, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "managed_node_connectivity_group_dn", resourceName, "managed_node_connectivity_group_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
				),
			},
			{
				Config:      CreateAccMgmtZoneDataSourceUpdate(mgmtGrpName, zoneType, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccMgmtZoneDSWithInvalidParentDn(mgmtGrpName, zoneType, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccMgmtZoneDataSourceUpdatedResource(mgmtGrpName, zoneType, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccMgmtZoneConfigDataSource(mgmtGrpName, zoneType, rName string) string {
	fmt.Println("=== STEP  testing mgmt_zone Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		type = "%s"
		name = "%s"
	}

	data "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		name = aci_mgmt_zone.test.name
		type = aci_mgmt_zone.test.type
		depends_on = [ aci_mgmt_zone.test ]
	}
	`, mgmtGrpName, zoneType, rName)
	return resource
}

func CreateMgmtZoneDSWithoutRequired(mgmtGrpName, zoneType, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing mgmt_zone Data Source without ", attrName)
	rBlock := `
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		type = "%s"
		name = "%s"
	}
	`
	switch attrName {
	case "managed_node_connectivity_group_dn":
		rBlock += `
	data "aci_mgmt_zone" "test" {
	#	managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		name = aci_mgmt_zone.test.name
		type = aci_mgmt_zone.test.type
		depends_on = [ aci_mgmt_zone.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
	#	name = aci_mgmt_zone.test.name
		type = aci_mgmt_zone.test.type
		depends_on = [ aci_mgmt_zone.test ]
	}
		`
	case "type":
		rBlock += `
	data "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		name = aci_mgmt_zone.test.name
	#	type = aci_mgmt_zone.test.type
		depends_on = [ aci_mgmt_zone.test ]
	}
		`

	}
	return fmt.Sprintf(rBlock, mgmtGrpName, zoneType, rName)
}

func CreateAccMgmtZoneDSWithInvalidParentDn(mgmtGrpName, zoneType, rName string) string {
	fmt.Println("=== STEP  testing mgmt_zone Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		type = "%s"
		name = "%s"
	}

	data "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = "${aci_managed_node_connectivity_group.test.id}_invalid"
		name = aci_mgmt_zone.test.name
		type = aci_mgmt_zone.test.type
		depends_on = [ aci_mgmt_zone.test ]
	}
	`, mgmtGrpName, zoneType, rName)
	return resource
}

func CreateAccMgmtZoneDataSourceUpdate(mgmtGrpName, zoneType, rName, key, value string) string {
	fmt.Println("=== STEP  testing mgmt_zone Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		type = "%s"
		name = "%s"
	}

	data "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		name = aci_mgmt_zone.test.name
		type = aci_mgmt_zone.test.type
		%s = "%s"
		depends_on = [ aci_mgmt_zone.test ]
	}
	`, mgmtGrpName, zoneType, rName, key, value)
	return resource
}

func CreateAccMgmtZoneDataSourceUpdatedResource(mgmtGrpName, zoneType, rName, key, value string) string {
	fmt.Println("=== STEP  testing mgmt_zone Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		type = "%s"
		name = "%s"
		%s = "%s"
	}

	data "aci_mgmt_zone" "test" {
		managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.test.id
		name = aci_mgmt_zone.test.name
		type = aci_mgmt_zone.test.type
		depends_on = [ aci_mgmt_zone.test ]
	}
	`, mgmtGrpName, zoneType, rName, key, value)
	return resource
}
