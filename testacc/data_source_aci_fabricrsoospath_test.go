package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciOutofServiceFabricPathlistDataSource_Basic(t *testing.T) {
	resourceName := "aci_interface_blacklist.test"
	dataSourceName := "data.aci_interface_blacklist.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	podId := "1"
	nodeId := "201"
	interfaceName := "eth1/1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOutofServiceFabricPathDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateInterfaceBlacklistDSWithoutRequired(podId, nodeId, interfaceName, "pod_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateInterfaceBlacklistDSWithoutRequired(podId, nodeId, interfaceName, "node_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateInterfaceBlacklistDSWithoutRequired(podId, nodeId, interfaceName, "interface"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccInterfaceBlacklistConfigDataSource(podId, nodeId, interfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "pod_id", resourceName, "pod_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "node_id", resourceName, "node_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "interface", resourceName, "interface"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config:      CreateAccInterfaceBlacklistDataSourceUpdate(podId, nodeId, interfaceName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccInterfaceBlacklistDSWithInvalidPodID(podId, nodeId, interfaceName),
				ExpectError: regexp.MustCompile(`Object may not exists`),
			},
			{
				Config: CreateAccInterfaceBlacklistDataSourceUpdatedResource(podId, nodeId, interfaceName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccInterfaceBlacklistConfigDataSource(podId, nodeId, interfaceName string) string {
	fmt.Println("=== STEP  testing interface_blacklist Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		pod_id  = %s
  		node_id = %s
 		interface = "%s"
	}

	data "aci_interface_blacklist" "test" {
	
		pod_id  = aci_interface_blacklist.test.pod_id
		node_id  = aci_interface_blacklist.test.node_id
		interface  = aci_interface_blacklist.test.interface
		depends_on = [ aci_interface_blacklist.test ]
	}
	`, podId, nodeId, interfaceName)
	return resource
}

func CreateInterfaceBlacklistDSWithoutRequired(podId, nodeId, interfaceName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing interface_blacklist Data Source without ", attrName)
	rBlock := `
	
	resource "aci_interface_blacklist" "test" {
	
		pod_id  = %s
  		node_id = %s
 		interface = "%s"
	}
	`
	switch attrName {
	case "pod_id":
		rBlock += `
	data "aci_interface_blacklist" "test" {
	
	#	pod_id  = aci_interface_blacklist.test.pod_id
		node_id  = aci_interface_blacklist.test.node_id
		interface  = aci_interface_blacklist.test.interface
		depends_on = [ aci_interface_blacklist.test ]
	}
		`
	case "node_id":
		rBlock += `
	data "aci_interface_blacklist" "test" {
	
		pod_id  = aci_interface_blacklist.test.pod_id
	#	node_id  = aci_interface_blacklist.test.node_id
		interface  = aci_interface_blacklist.test.interface
		depends_on = [ aci_interface_blacklist.test ]
	}
		`
	case "interface":
		rBlock += `
	data "aci_interface_blacklist" "test" {
	
		pod_id  = aci_interface_blacklist.test.pod_id
		node_id  = aci_interface_blacklist.test.node_id
	#	interface  = aci_interface_blacklist.test.interface
		depends_on = [ aci_interface_blacklist.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, podId, nodeId, interfaceName)
}

func CreateAccInterfaceBlacklistDSWithInvalidPodID(podId, nodeId, interfaceName string) string {
	fmt.Println("=== STEP  testing interface_blacklist Data Source with invalid pod_id")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		pod_id  = %s
  		node_id = %s
 		interface = "%s"
	}

	data "aci_interface_blacklist" "test" {
	
		pod_id  = aci_interface_blacklist.test.pod_id+1
		node_id  = aci_interface_blacklist.test.node_id
		interface  = aci_interface_blacklist.test.interface
		depends_on = [ aci_interface_blacklist.test ]
	}
	`, podId, nodeId, interfaceName)
	return resource
}

func CreateAccInterfaceBlacklistDataSourceUpdate(podId, nodeId, interfaceName, key, value string) string {
	fmt.Println("=== STEP  testing interface_blacklist Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		pod_id  = %s
  		node_id = %s
 		interface = "%s"
	}

	data "aci_interface_blacklist" "test" {
	
		pod_id  = aci_interface_blacklist.test.pod_id
		node_id  = aci_interface_blacklist.test.node_id
		interface  = aci_interface_blacklist.test.interface
		%s = "%s"
		depends_on = [ aci_interface_blacklist.test ]
	}
	`, podId, nodeId, interfaceName, key, value)
	return resource
}

func CreateAccInterfaceBlacklistDataSourceUpdatedResource(podId, nodeId, interfaceName, key, value string) string {
	fmt.Println("=== STEP  testing interface_blacklist Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_blacklist" "test" {
	
		pod_id  = %s
  		node_id = %s
 		interface = "%s"
		%s = "%s"
	}

	data "aci_interface_blacklist" "test" {
	
		pod_id  = aci_interface_blacklist.test.pod_id
		node_id  = aci_interface_blacklist.test.node_id
		interface  = aci_interface_blacklist.test.interface
		depends_on = [ aci_interface_blacklist.test ]
	}
	`, podId, nodeId, interfaceName, key, value)
	return resource
}
