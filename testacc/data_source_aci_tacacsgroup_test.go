package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTACACSMonitoringDestinationGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_tacacs_accounting.test"
	dataSourceName := "data.aci_tacacs_accounting.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSMonitoringDestinationGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateTACACSMonitoringDestinationGroupDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSMonitoringDestinationGroupConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccTACACSMonitoringDestinationGroupDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccTACACSMonitoringDestinationGroupDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccTACACSMonitoringDestinationGroupDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccTACACSMonitoringDestinationGroupConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing tacacs_accounting Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
	}

	data "aci_tacacs_accounting" "test" {
	
		name  = aci_tacacs_accounting.test.name
		depends_on = [ aci_tacacs_accounting.test ]
	}
	`, rName)
	return resource
}

func CreateTACACSMonitoringDestinationGroupDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_accounting Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_tacacs_accounting" "test" {
	
	#	name  = aci_tacacs_accounting.test.name
		depends_on = [ aci_tacacs_accounting.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccTACACSMonitoringDestinationGroupDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing tacacs_accounting Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
	}

	data "aci_tacacs_accounting" "test" {
	
		name  = "${aci_tacacs_accounting.test.name}_invalid"
		depends_on = [ aci_tacacs_accounting.test ]
	}
	`, rName)
	return resource
}

func CreateAccTACACSMonitoringDestinationGroupDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing tacacs_accounting Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
	}

	data "aci_tacacs_accounting" "test" {
	
		name  = aci_tacacs_accounting.test.name
		%s = "%s"
		depends_on = [ aci_tacacs_accounting.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccTACACSMonitoringDestinationGroupDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing tacacs_accounting Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_accounting" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_tacacs_accounting" "test" {
	
		name  = aci_tacacs_accounting.test.name
		depends_on = [ aci_tacacs_accounting.test ]
	}
	`, rName, key, value)
	return resource
}
