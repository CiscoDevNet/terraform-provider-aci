package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciNodeBlockMGDataSource_Basic(t *testing.T) {
	resourceName := "aci_maintenance_group_node.test"
	dataSourceName := "data.aci_maintenance_group_node.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	maintMaintGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockMGDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateNodeBlockMGDSWithoutRequired(maintMaintGrpName, rName, "pod_maintenance_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateNodeBlockMGDSWithoutRequired(maintMaintGrpName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccNodeBlockMGConfigDataSource(maintMaintGrpName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "pod_maintenance_group_dn", resourceName, "pod_maintenance_group_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "from_", resourceName, "from_"),
					resource.TestCheckResourceAttrPair(dataSourceName, "to_", resourceName, "to_"),
				),
			},
			{
				Config:      CreateAccNodeBlockMGDataSourceUpdate(maintMaintGrpName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccNodeBlockMGDSWithInvalidName(maintMaintGrpName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccNodeBlockMGDataSourceUpdatedResource(maintMaintGrpName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccNodeBlockMGConfigDataSource(maintMaintGrpName, rName string) string {
	fmt.Println("=== STEP  testing maintenance_group_node Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "%s"
	}

	data "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = aci_maintenance_group_node.test.name
		depends_on = [ aci_maintenance_group_node.test ]
	}
	`, maintMaintGrpName, rName)
	return resource
}

func CreateNodeBlockMGDSWithoutRequired(maintMaintGrpName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing maintenance_group_node Data Source without ", attrName)
	rBlock := `
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "pod_maintenance_group_dn":
		rBlock += `
	data "aci_maintenance_group_node" "test" {
	#	pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = aci_maintenance_group_node.test.name
		depends_on = [ aci_maintenance_group_node.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
	#	name  = aci_maintenance_group_node.test.name
		depends_on = [ aci_maintenance_group_node.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, maintMaintGrpName, rName)
}

func CreateAccNodeBlockMGDSWithInvalidName(maintMaintGrpName, rName string) string {
	fmt.Println("=== STEP  testing maintenance_group_node Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "%s"
	}

	data "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "${aci_maintenance_group_node.test.name}_invalid"
		depends_on = [ aci_maintenance_group_node.test ]
	}
	`, maintMaintGrpName, rName)
	return resource
}

func CreateAccNodeBlockMGDataSourceUpdate(maintMaintGrpName, rName, key, value string) string {
	fmt.Println("=== STEP  testing maintenance_group_node Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "%s"
	}

	data "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = aci_maintenance_group_node.test.name
		%s = "%s"
		depends_on = [ aci_maintenance_group_node.test ]
	}
	`, maintMaintGrpName, rName, key, value)
	return resource
}

func CreateAccNodeBlockMGDataSourceUpdatedResource(maintMaintGrpName, rName, key, value string) string {
	fmt.Println("=== STEP  testing maintenance_group_node Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = aci_maintenance_group_node.test.name
		depends_on = [ aci_maintenance_group_node.test ]
	}
	`, maintMaintGrpName, rName, key, value)
	return resource
}
