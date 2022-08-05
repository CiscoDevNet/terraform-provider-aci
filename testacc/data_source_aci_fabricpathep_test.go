package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFabricPathEpDataSource_Basic(t *testing.T) {
	dataSourceName := "data.aci_fabric_path_ep.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      CreateFabricPathEpDSWithoutRequired("pod_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateFabricPathEpDSWithoutRequired("name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateFabricPathEpDSWithoutRequired("node_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricPathEpConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "pod_id", podId),
					resource.TestCheckResourceAttr(dataSourceName, "name", pathEpName),
					resource.TestCheckResourceAttr(dataSourceName, "node_id", nodeId),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc"),
				),
			},
			{
				Config:      CreateAccFabricPathEpDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccFabricPathEpDSWithInvalidParentDn(randomValue),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccFabricPathEpConfigDataSource(),
			},
		},
	})
}

func CreateAccFabricPathEpConfigDataSource() string {
	fmt.Println("=== STEP  testing fabric_path_ep Data Source with required arguments only")
	resource := fmt.Sprintf(`
	data "aci_fabric_path_ep" "test" {
		pod_id  = "%s"
		name  = "%s"
		node_id = "%s"
	}
	`, podId, pathEpName, nodeId)
	return resource
}

func CreateFabricPathEpDSWithoutRequired(attrName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_path_ep Data Source without ", attrName)
	rBlock := `
	`
	switch attrName {
	case "pod_id":
		rBlock += `
	data "aci_fabric_path_ep" "test" {
	#	pod_id  = "%s"
		name  = "%s"
		node_id = "%s"
	}
		`
	case "name":
		rBlock += `
	data "aci_fabric_path_ep" "test" {
		pod_id  = "%s"
	#	name  = "%s"
		node_id = "%s"
	}
		`
	case "node_id":
		rBlock += `
	data "aci_fabric_path_ep" "test" {
		pod_id  = "%s"
		name  = "%s"
	#	node_id = "%s"
	}
	`
	}
	return fmt.Sprintf(rBlock, podId, pathEpName, nodeId)
}

func CreateAccFabricPathEpDSWithInvalidParentDn(value string) string {
	fmt.Println("=== STEP  testing fabric_path_ep Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`

	data "aci_fabric_path_ep" "test" {
		pod_id  = "%s"
		name  = "%s"
		node_id = "%s"
	}
	`, podId, value, nodeId)
	return resource
}

func CreateAccFabricPathEpDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing fabric_path_ep Data Source with random attribute")
	resource := fmt.Sprintf(`

	data "aci_fabric_path_ep" "test" {
		pod_id  = "%s"
		name  = "%s"
		node_id = "%s"
		%s = "%s"
	}
	`, podId, pathEpName, nodeId, key, value)
	return resource
}
