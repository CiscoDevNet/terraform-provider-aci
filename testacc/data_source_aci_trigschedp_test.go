package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTriggerSchedulerDataSource_Basic(t *testing.T) {
	resourceName := "aci_trigger_scheduler.test"
	dataSourceName := "data.aci_trigger_scheduler.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTriggerSchedulerDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateTriggerSchedulerDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTriggerSchedulerConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccTriggerSchedulerDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccTriggerSchedulerDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccTriggerSchedulerDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccTriggerSchedulerConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing trigger_scheduler Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
	
		name  = "%s"
	}

	data "aci_trigger_scheduler" "test" {
	
		name  = aci_trigger_scheduler.test.name
		depends_on = [ aci_trigger_scheduler.test ]
	}
	`, rName)
	return resource
}

func CreateTriggerSchedulerDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing trigger_scheduler Data Source without ", attrName)
	rBlock := `
	
	resource "aci_trigger_scheduler" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_trigger_scheduler" "test" {
	
	#	name  = aci_trigger_scheduler.test.name
		depends_on = [ aci_trigger_scheduler.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccTriggerSchedulerDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing trigger_scheduler Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
	
		name  = "%s"
	}

	data "aci_trigger_scheduler" "test" {
	
		name  = "${aci_trigger_scheduler.test.name}_invalid"
		depends_on = [ aci_trigger_scheduler.test ]
	}
	`, rName)
	return resource
}

func CreateAccTriggerSchedulerDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing trigger_scheduler Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
	
		name  = "%s"
	}

	data "aci_trigger_scheduler" "test" {
	
		name  = aci_trigger_scheduler.test.name
		%s = "%s"
		depends_on = [ aci_trigger_scheduler.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccTriggerSchedulerDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing trigger_scheduler Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_trigger_scheduler" "test" {
	
		name  = aci_trigger_scheduler.test.name
		depends_on = [ aci_trigger_scheduler.test ]
	}
	`, rName, key, value)
	return resource
}
