package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFabricNodeControlDataSource_Basic(t *testing.T) {
	resourceName := "aci_fabric_node_control.test"
	dataSourceName := "data.aci_fabric_node_control.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeControlDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateFabricNodeControlDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricNodeControlConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "control", resourceName, "control"),
					resource.TestCheckResourceAttrPair(dataSourceName, "feature_sel", resourceName, "feature_sel"),
				),
			},
			{
				Config:      CreateAccFabricNodeControlDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccFabricNodeControlDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccFabricNodeControlDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFabricNodeControlConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing fabric_node_control Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_control" "test" {
		name  = "%s"
	}
	data "aci_fabric_node_control" "test" {
		name  = aci_fabric_node_control.test.name
		depends_on = [ aci_fabric_node_control.test ]
	}
	`, rName)
	return resource
}

func CreateFabricNodeControlDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_node_control Data Source without ", attrName)
	rBlock := `
	resource "aci_fabric_node_control" "test" {
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_fabric_node_control" "test" {
	#	name  = aci_fabric_node_control.test.name
		depends_on = [ aci_fabric_node_control.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFabricNodeControlDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing fabric_node_control Data Source with Invalid Name")
	resource := fmt.Sprintf(`
	resource "aci_fabric_node_control" "test" {
		name  = "%s"
	}

	data "aci_fabric_node_control" "test" {
		name  = "${aci_fabric_node_control.test.name}_invalid"
		depends_on = [ aci_fabric_node_control.test ]
	}
	`, rName)
	return resource
}

func CreateAccFabricNodeControlDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing fabric_node_control Data Source with random attribute")
	resource := fmt.Sprintf(`
	resource "aci_fabric_node_control" "test" {
		name  = "%s"
	}

	data "aci_fabric_node_control" "test" {
		name  = aci_fabric_node_control.test.name
		%s = "%s"
		depends_on = [ aci_fabric_node_control.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccFabricNodeControlDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing fabric_node_control Data Source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_fabric_node_control" "test" {
		name  = "%s"
		%s = "%s"
	}

	data "aci_fabric_node_control" "test" {
		name  = aci_fabric_node_control.test.name
		depends_on = [ aci_fabric_node_control.test ]
	}
	`, rName, key, value)
	return resource
}
