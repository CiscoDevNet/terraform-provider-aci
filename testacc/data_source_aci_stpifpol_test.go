package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSpanningTreeInterfacePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_spanning_tree_interface_policy.test"
	dataSourceName := "data.aci_spanning_tree_interface_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpanningTreeInterfacePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateSpanningTreeInterfacePolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpanningTreeInterfacePolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.#", resourceName, "ctrl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.0", resourceName, "ctrl.0"),
				),
			},
			{
				Config:      CreateAccSpanningTreeInterfacePolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccSpanningTreeInterfacePolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccSpanningTreeInterfacePolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccSpanningTreeInterfacePolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing spanning_tree_interface_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spanning_tree_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_spanning_tree_interface_policy" "test" {
	
		name  = aci_spanning_tree_interface_policy.test.name
		depends_on = [ aci_spanning_tree_interface_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccSpanningTreeInterfacePolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing spanning_tree_interface_policy Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_spanning_tree_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_spanning_tree_interface_policy" "test" {
	
		name  = "${aci_spanning_tree_interface_policy.test.name}_invalid"
		depends_on = [ aci_spanning_tree_interface_policy.test ]
	}
	`, rName)
	return resource
}

func CreateSpanningTreeInterfacePolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spanning_tree_interface_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_spanning_tree_interface_policy" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_spanning_tree_interface_policy" "test" {
	
	#	name  = "%s"
		depends_on = [ aci_spanning_tree_interface_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSpanningTreeInterfacePolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing spanning_tree_interface_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_spanning_tree_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_spanning_tree_interface_policy" "test" {
	
		name  = aci_spanning_tree_interface_policy.test.name
		%s = "%s"
		depends_on = [ aci_spanning_tree_interface_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccSpanningTreeInterfacePolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing spanning_tree_interface_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_spanning_tree_interface_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_spanning_tree_interface_policy" "test" {
	
		name  = aci_spanning_tree_interface_policy.test.name
		depends_on = [ aci_spanning_tree_interface_policy.test ]
	}
	`, rName, key, value)
	return resource
}
