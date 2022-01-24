package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciPodMaintenanceGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_pod_maintenance_group.test"
	dataSourceName := "data.aci_pod_maintenance_group.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPODMaintenanceGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreatePodMaintenanceGroupDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccPodMaintenanceGroupConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fwtype", resourceName, "fwtype"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pod_maintenance_group_type", resourceName, "pod_maintenance_group_type"),
				),
			},
			{
				Config:      CreateAccPodMaintenanceGroupDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccPodMaintenanceGroupDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccPodMaintenanceGroupDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccPodMaintenanceGroupConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing pod_maintenance_group Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
	
		name  = "%s"
	}

	data "aci_pod_maintenance_group" "test" {
	
		name  = aci_pod_maintenance_group.test.name
		depends_on = [ aci_pod_maintenance_group.test ]
	}
	`, rName)
	return resource
}

func CreatePodMaintenanceGroupDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing pod_maintenance_group Data Source without ", attrName)
	rBlock := `
	
	resource "aci_pod_maintenance_group" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_pod_maintenance_group" "test" {
	
	#	name  = aci_pod_maintenance_group.test.name
		depends_on = [ aci_pod_maintenance_group.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccPodMaintenanceGroupDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing pod_maintenance_group Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
	
		name  = "%s"
	}

	data "aci_pod_maintenance_group" "test" {
	
		name  = "${aci_pod_maintenance_group.test.name}_invalid"
		depends_on = [ aci_pod_maintenance_group.test ]
	}
	`, rName)
	return resource
}

func CreateAccPodMaintenanceGroupDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing pod_maintenance_group Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
	
		name  = "%s"
	}

	data "aci_pod_maintenance_group" "test" {
	
		name  = aci_pod_maintenance_group.test.name
		%s = "%s"
		depends_on = [ aci_pod_maintenance_group.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccPodMaintenanceGroupDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing pod_maintenance_group Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_pod_maintenance_group" "test" {
	
		name  = aci_pod_maintenance_group.test.name
		depends_on = [ aci_pod_maintenance_group.test ]
	}
	`, rName, key, value)
	return resource
}
