package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFabricNodeMemberDataSource_Basic(t *testing.T) {
	resourceName := "aci_fabric_node_member.test"
	dataSourceName := "data.aci_fabric_node_member.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	serial := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeMemberDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFabricNodeMemberDSWithoutRequired(rName, serial, fabricNodeMemNodeId5, "serial"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricNodeMemberConfigDataSource(rName, serial, fabricNodeMemNodeId5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "serial", resourceName, "serial"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ext_pool_id", resourceName, "ext_pool_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fabric_id", resourceName, "fabric_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "node_id", resourceName, "node_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "node_type", resourceName, "node_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pod_id", resourceName, "pod_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "role", resourceName, "role"),
				),
			},
			{
				Config:      CreateAccFabricNodeMemberDataSourceUpdate(rName, serial, fabricNodeMemNodeId5, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccFabricNodeMemberDSWithInvalidSerial(rName, serial, fabricNodeMemNodeId5),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccFabricNodeMemberDataSourceUpdatedResource(rName, serial, fabricNodeMemNodeId5, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFabricNodeMemberConfigDataSource(name, serial, node string) string {
	fmt.Println("=== STEP  testing fabric_node_member Data Source with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_fabric_node_member" "test" {
		name = "%s"
		serial  = "%s"
		node_id = "%s"
	}

	data "aci_fabric_node_member" "test" {

		serial  = aci_fabric_node_member.test.serial
		depends_on = [ aci_fabric_node_member.test ]
	}
	`, name, serial, node)
	return resource
}

func CreateFabricNodeMemberDSWithoutRequired(name, serial, node, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_node_member Data Source without ", attrName)
	rBlock := `

	resource "aci_fabric_node_member" "test" {
		name = "%s"
		serial  = "%s"
		node_id = "%s"
	}
	`
	switch attrName {
	case "serial":
		rBlock += `
	data "aci_fabric_node_member" "test" {

	#	serial  = aci_fabric_node_member.test.serial
		depends_on = [ aci_fabric_node_member.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, name, serial, node)
}

func CreateAccFabricNodeMemberDSWithInvalidSerial(name, serial, node string) string {
	fmt.Println("=== STEP  testing fabric_node_member Data Source with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_fabric_node_member" "test" {
		name = "%s"
		serial  = "%s"
		node_id = "%s"
	}

	data "aci_fabric_node_member" "test" {

		serial  = "${aci_fabric_node_member.test.serial}_invalid"
		depends_on = [ aci_fabric_node_member.test ]
	}
	`, name, serial, node)
	return resource
}

func CreateAccFabricNodeMemberDataSourceUpdate(name, serial, node, key, value string) string {
	fmt.Println("=== STEP  testing fabric_node_member Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_fabric_node_member" "test" {
		name = "%s"
		serial  = "%s"
		node_id = "%s"
	}

	data "aci_fabric_node_member" "test" {

		serial  = aci_fabric_node_member.test.serial
		%s = "%s"
		depends_on = [ aci_fabric_node_member.test ]
	}
	`, name, serial, node, key, value)
	return resource
}

func CreateAccFabricNodeMemberDataSourceUpdatedResource(rName, serial, node, key, value string) string {
	fmt.Println("=== STEP  testing fabric_node_member Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_fabric_node_member" "test" {
		name = "%s"
		serial  = "%s"
		node_id = "%s"
		%s = "%s"
	}

	data "aci_fabric_node_member" "test" {

		serial  = aci_fabric_node_member.test.serial
		depends_on = [ aci_fabric_node_member.test ]
	}
	`, name, serial, node, key, value)
	return resource
}
