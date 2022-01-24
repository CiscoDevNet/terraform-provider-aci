package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciRecurringWindowDataSource_Basic(t *testing.T) {
	resourceName := "aci_recurring_window.test"
	dataSourceName := "data.aci_recurring_window.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	trigSchedPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRecurringWindowDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateRecurringWindowDSWithoutRequired(trigSchedPName, rName, "scheduler_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateRecurringWindowDSWithoutRequired(trigSchedPName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRecurringWindowConfigDataSource(trigSchedPName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "scheduler_dn", resourceName, "scheduler_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "concur_cap", resourceName, "concur_cap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "day", resourceName, "day"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hour", resourceName, "hour"),
					resource.TestCheckResourceAttrPair(dataSourceName, "minute", resourceName, "minute"),
					resource.TestCheckResourceAttrPair(dataSourceName, "node_upg_interval", resourceName, "node_upg_interval"),
					resource.TestCheckResourceAttrPair(dataSourceName, "proc_break", resourceName, "proc_break"),
					resource.TestCheckResourceAttrPair(dataSourceName, "proc_cap", resourceName, "proc_cap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "time_cap", resourceName, "time_cap"),
				),
			},
			{
				Config:      CreateAccRecurringWindowDataSourceUpdate(trigSchedPName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccRecurringWindowDSWithInvalidParentDn(trigSchedPName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccRecurringWindowDataSourceUpdatedResource(trigSchedPName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccRecurringWindowConfigDataSource(trigSchedPName, rName string) string {
	fmt.Println("=== STEP  testing recurring_window Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = "%s"
	}

	data "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = aci_recurring_window.test.name
		depends_on = [ aci_recurring_window.test ]
	}
	`, trigSchedPName, rName)
	return resource
}

func CreateRecurringWindowDSWithoutRequired(trigSchedPName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing recurring_window Data Source without ", attrName)
	rBlock := `
	
	resource "aci_trigger_scheduler" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "scheduler_dn":
		rBlock += `
	data "aci_recurring_window" "test" {
	#	scheduler_dn  = aci_trigger_scheduler.test.id
		name  = aci_recurring_window.test.name
		depends_on = [ aci_recurring_window.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
	#	name  = aci_recurring_window.test.name
		depends_on = [ aci_recurring_window.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, trigSchedPName, rName)
}

func CreateAccRecurringWindowDSWithInvalidParentDn(trigSchedPName, rName string) string {
	fmt.Println("=== STEP  testing recurring_window Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = "%s"
	}

	data "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = "${aci_recurring_window.test.name}_invalid"
		depends_on = [ aci_recurring_window.test ]
	}
	`, trigSchedPName, rName)
	return resource
}

func CreateAccRecurringWindowDataSourceUpdate(trigSchedPName, rName, key, value string) string {
	fmt.Println("=== STEP  testing recurring_window Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = "%s"
	}

	data "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = aci_recurring_window.test.name
		%s = "%s"
		depends_on = [ aci_recurring_window.test ]
	}
	`, trigSchedPName, rName, key, value)
	return resource
}

func CreateAccRecurringWindowDataSourceUpdatedResource(trigSchedPName, rName, key, value string) string {
	fmt.Println("=== STEP  testing recurring_window Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = aci_recurring_window.test.name
		depends_on = [ aci_recurring_window.test ]
	}
	`, trigSchedPName, rName, key, value)
	return resource
}
