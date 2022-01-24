package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciManagedNodeConnectivityGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_managed_node_connectivity_group.test"
	dataSourceName := "data.aci_managed_node_connectivity_group.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciManagedNodeConnectivityGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateManagedNodeConnectivityGroupDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccManagedNodeConnectivityGroupConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccManagedNodeConnectivityGroupDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccManagedNodeConnectivityGroupDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccManagedNodeConnectivityGroupDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccManagedNodeConnectivityGroupConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing managed_node_connectivity_group Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
	
		name  = "%s"
	}

	data "aci_managed_node_connectivity_group" "test" {
	
		name  = aci_managed_node_connectivity_group.test.name
		depends_on = [ aci_managed_node_connectivity_group.test ]
	}
	`, rName)
	return resource
}

func CreateManagedNodeConnectivityGroupDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing managed_node_connectivity_group Data Source without ", attrName)
	rBlock := `
	
	resource "aci_managed_node_connectivity_group" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_managed_node_connectivity_group" "test" {
	
	#	name  = aci_managed_node_connectivity_group.test.name
		depends_on = [ aci_managed_node_connectivity_group.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccManagedNodeConnectivityGroupDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing managed_node_connectivity_group Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
	
		name  = "%s"
	}

	data "aci_managed_node_connectivity_group" "test" {
	
		name  = "${aci_managed_node_connectivity_group.test.name}_invalid"
		depends_on = [ aci_managed_node_connectivity_group.test ]
	}
	`, rName)
	return resource
}

func CreateAccManagedNodeConnectivityGroupDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing managed_node_connectivity_group Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
	
		name  = "%s"
	}

	data "aci_managed_node_connectivity_group" "test" {
	
		name  = aci_managed_node_connectivity_group.test.name
		%s = "%s"
		depends_on = [ aci_managed_node_connectivity_group.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccManagedNodeConnectivityGroupDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing managed_node_connectivity_group Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_managed_node_connectivity_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_managed_node_connectivity_group" "test" {
	
		name  = aci_managed_node_connectivity_group.test.name
		depends_on = [ aci_managed_node_connectivity_group.test ]
	}
	`, rName, key, value)
	return resource
}
