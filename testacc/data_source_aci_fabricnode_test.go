package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTopologyFabricNodeDataSource_Basic(t *testing.T) {
	dataSourceName := "data.aci_fabric_node.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      CreateTopologyFabricNodeDSWithoutRequired(fabricPodDn, fabricNodeId, "fabric_pod_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateTopologyFabricNodeDSWithoutRequired(fabricPodDn, fabricNodeId, "fabric_node_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTopologyFabricNodeConfigDataSource(fabricPodDn, fabricNodeId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "fabric_pod_dn", fabricPodDn),
					resource.TestCheckResourceAttr(dataSourceName, "fabric_node_id", fabricNodeId),
					resource.TestCheckResourceAttrSet(dataSourceName, "ad_st"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "apic_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "fabric_st"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "node_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "role"),
				),
			},
			{
				Config:      CreateAccTopologyFabricNodeDataSourceUpdate(fabricPodDn, fabricNodeId, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccTopologyFabricNodeDSWithInvalidNodeId(fabricPodDn, randomValue),
				ExpectError: regexp.MustCompile(`Invalid RN`),
			},
			{
				Config: CreateAccTopologyFabricNodeConfigDataSource(fabricPodDn, fabricNodeId),
			},
		},
	})
}

func CreateAccTopologyFabricNodeConfigDataSource(podId, nodeId string) string {
	fmt.Println("=== STEP  testing fabric_node Data Source with required arguments only")
	resource := fmt.Sprintf(`

	data "aci_fabric_node" "test" {
		fabric_pod_dn  = "%s"
		fabric_node_id  = "%s"
	}
	`, podId, nodeId)
	return resource
}

func CreateTopologyFabricNodeDSWithoutRequired(podDn, nodeId, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_node Data Source without ", attrName)
	rBlock := `
	`
	switch attrName {
	case "fabric_pod_dn":
		rBlock += `
	data "aci_fabric_node" "test" {
	#	fabric_pod_dn  = "%s"
		fabric_node_id  = "%s"
	}
		`
	case "fabric_node_id":
		rBlock += `
	data "aci_fabric_node" "test" {
		fabric_pod_dn  = "%s"
	#	fabric_node_id  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, podDn, nodeId)
}

func CreateAccTopologyFabricNodeDSWithInvalidNodeId(podId, nodeId string) string {
	fmt.Println("=== STEP  testing fabric_node Data Source with invalid fabric_node_id")
	resource := fmt.Sprintf(`

	data "aci_fabric_node" "test" {
		fabric_pod_dn  = "%s"
		fabric_node_id  = "%s"
	}
	`, podId, nodeId)
	return resource
}

func CreateAccTopologyFabricNodeDataSourceUpdate(podId, nodeId, key, value string) string {
	fmt.Println("=== STEP  testing fabric_node Data Source with random attribute")
	resource := fmt.Sprintf(`

	data "aci_fabric_node" "test" {
		fabric_pod_dn  = "%s"
		fabric_node_id  = "%s"
		%s = "%s"
	}
	`, podId, nodeId, key, value)
	return resource
}
