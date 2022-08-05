package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciRecurringWindow_Basic(t *testing.T) {
	var recurring_window_default models.RecurringWindow
	var recurring_window_updated models.RecurringWindow
	resourceName := "aci_recurring_window.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	trigSchedPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRecurringWindowDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateRecurringWindowWithoutRequired(trigSchedPName, rName, "scheduler_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateRecurringWindowWithoutRequired(trigSchedPName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRecurringWindowConfig(trigSchedPName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_default),
					resource.TestCheckResourceAttr(resourceName, "scheduler_dn", fmt.Sprintf("uni/fabric/schedp-%s", trigSchedPName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "concur_cap", "unlimited"),
					resource.TestCheckResourceAttr(resourceName, "day", "every-day"),
					resource.TestCheckResourceAttr(resourceName, "hour", "0"),
					resource.TestCheckResourceAttr(resourceName, "minute", "0"),
					resource.TestCheckResourceAttr(resourceName, "node_upg_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "proc_break", "none"),
					resource.TestCheckResourceAttr(resourceName, "proc_cap", "unlimited"),
					resource.TestCheckResourceAttr(resourceName, "time_cap", "unlimited"),
				),
			},
			{
				Config: CreateAccRecurringWindowConfigWithOptionalValues(trigSchedPName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "scheduler_dn", fmt.Sprintf("uni/fabric/schedp-%s", trigSchedPName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_recurring_window"),
					resource.TestCheckResourceAttr(resourceName, "concur_cap", "1"),

					resource.TestCheckResourceAttr(resourceName, "day", "Friday"),
					resource.TestCheckResourceAttr(resourceName, "hour", "1"),
					resource.TestCheckResourceAttr(resourceName, "minute", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_upg_interval", "1"),

					resource.TestCheckResourceAttr(resourceName, "proc_break", "none"),
					resource.TestCheckResourceAttr(resourceName, "proc_cap", "1"),

					resource.TestCheckResourceAttr(resourceName, "time_cap", "00:00:00:00.020"),

					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccRecurringWindowConfigUpdatedName(trigSchedPName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccRecurringWindowRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccRecurringWindowConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "scheduler_dn", fmt.Sprintf("uni/fabric/schedp-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciRecurringWindowIdNotEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowConfig(trigSchedPName, rName),
			},
			{
				Config: CreateAccRecurringWindowConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "scheduler_dn", fmt.Sprintf("uni/fabric/schedp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciRecurringWindowIdNotEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
		},
	})
}

func TestAccAciRecurringWindow_Update(t *testing.T) {
	var recurring_window_default models.RecurringWindow
	var recurring_window_updated models.RecurringWindow
	resourceName := "aci_recurring_window.test"
	rName := makeTestVariable(acctest.RandString(5))

	trigSchedPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRecurringWindowDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRecurringWindowConfig(trigSchedPName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_default),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "concur_cap", "65535"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "concur_cap", "65535"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "concur_cap", "32767"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "concur_cap", "32767"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},

			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "day", "Monday"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "day", "Monday"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "day", "Saturday"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "day", "Saturday"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "day", "Sunday"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "day", "Sunday"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "day", "Thursday"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "day", "Thursday"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "day", "Tuesday"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "day", "Tuesday"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "day", "Wednesday"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "day", "Wednesday"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "day", "even-day"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "day", "even-day"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "day", "odd-day"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "day", "odd-day"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			}, {
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "hour", "23"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "hour", "23"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "hour", "11"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "hour", "11"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "minute", "59"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "minute", "59"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "minute", "29"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "minute", "29"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "node_upg_interval", "18000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "node_upg_interval", "18000"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "node_upg_interval", "9000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "node_upg_interval", "9000"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "proc_cap", "65535"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "proc_cap", "65535"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},
			{
				Config: CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "proc_cap", "32767"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists(resourceName, &recurring_window_updated),
					resource.TestCheckResourceAttr(resourceName, "proc_cap", "32767"),
					testAccCheckAciRecurringWindowIdEqual(&recurring_window_default, &recurring_window_updated),
				),
			},

			{
				Config: CreateAccRecurringWindowConfig(trigSchedPName, rName),
			},
		},
	})
}

func TestAccAciRecurringWindow_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	trigSchedPName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRecurringWindowDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRecurringWindowConfig(trigSchedPName, rName),
			},
			{
				Config:      CreateAccRecurringWindowWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "concur_cap", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "concur_cap", "-1"),
				ExpectError: regexp.MustCompile(`unknown property`),
			},
			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "concur_cap", "65536"),
				ExpectError: regexp.MustCompile(`unknown property`),
			},

			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "day", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "hour", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "hour", "-1"),
				ExpectError: regexp.MustCompile(`unknown property`),
			},
			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "hour", "24"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "minute", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "minute", "-1"),
				ExpectError: regexp.MustCompile(`unknown property`),
			},
			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "minute", "60"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "node_upg_interval", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "node_upg_interval", "-1"),
				ExpectError: regexp.MustCompile(`unknown property`),
			},
			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "node_upg_interval", "18001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "proc_break", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "proc_cap", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "proc_cap", "-1"),
				ExpectError: regexp.MustCompile(`unknown property`),
			},
			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "proc_cap", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "time_cap", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, "time_cap", "01:23:59:59.000"),
				ExpectError: regexp.MustCompile(`Max Duration cannot exceed 24 hours`),
			},

			{
				Config:      CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccRecurringWindowConfig(trigSchedPName, rName),
			},
		},
	})
}

func TestAccAciRecurringWindow_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	trigSchedPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciRecurringWindowDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccRecurringWindowConfigMultiple(trigSchedPName, rName),
			},
		},
	})
}

func testAccCheckAciRecurringWindowExists(name string, recurring_window *models.RecurringWindow) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Recurring Window %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Recurring Window dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		recurring_windowFound := models.RecurringWindowFromContainer(cont)
		if recurring_windowFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Recurring Window %s not found", rs.Primary.ID)
		}
		*recurring_window = *recurring_windowFound
		return nil
	}
}

func testAccCheckAciRecurringWindowDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing recurring_window destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_recurring_window" {
			cont, err := client.Get(rs.Primary.ID)
			recurring_window := models.RecurringWindowFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Recurring Window %s Still exists", recurring_window.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRecurringWindowIdEqual(m1, m2 *models.RecurringWindow) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("recurring_window DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciRecurringWindowIdNotEqual(m1, m2 *models.RecurringWindow) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("recurring_window DNs are equal")
		}
		return nil
	}
}

func CreateRecurringWindowWithoutRequired(trigSchedPName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing recurring_window creation without ", attrName)
	rBlock := `
	
	resource "aci_trigger_scheduler" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "scheduler_dn":
		rBlock += `
	resource "aci_recurring_window" "test" {
	#	scheduler_dn  = aci_trigger_scheduler.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, trigSchedPName, rName)
}

func CreateAccRecurringWindowConfigWithRequiredParams(trigSchedPName, rName string) string {
	fmt.Printf("=== STEP  testing recurring_window creation with parent resource name %s and resource name %s\n", trigSchedPName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = "%s"
	}
	`, trigSchedPName, rName)
	return resource
}
func CreateAccRecurringWindowConfigUpdatedName(trigSchedPName, rName string) string {
	fmt.Println("=== STEP  testing recurring_window creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = "%s"
	}
	`, trigSchedPName, rName)
	return resource
}

func CreateAccRecurringWindowConfig(trigSchedPName, rName string) string {
	fmt.Println("=== STEP  testing recurring_window creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = "%s"
	}
	`, trigSchedPName, rName)
	return resource
}

func CreateAccRecurringWindowConfigMultiple(trigSchedPName, rName string) string {
	fmt.Println("=== STEP  testing multiple recurring_window creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, trigSchedPName, rName)
	return resource
}

func CreateAccRecurringWindowWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing recurring_window creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_recurring_window" "test" {
		scheduler_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccRecurringWindowConfigWithOptionalValues(trigSchedPName, rName string) string {
	fmt.Println("=== STEP  Basic: testing recurring_window creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_recurring_window" "test" {
		scheduler_dn  = "${aci_trigger_scheduler.test.id}"
		name  = "%s"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_recurring_window"
		concur_cap = "1"
		day = "Friday"
		hour = "1"
		minute = "1"
		node_upg_interval = "1"
		proc_break = "none"
		proc_cap = "1"
		time_cap = "00:00:00:00.020"
		
	}
	`, trigSchedPName, rName)

	return resource
}

func CreateAccRecurringWindowRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing recurring_window updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_recurring_window" "test" {
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_recurring_window"
		concur_cap = "1"
		day = "Friday"
		hour = "1"
		minute = "1"
		node_upg_interval = "1"
		proc_break = ""
		proc_cap = "1"
		time_cap = ""
		
	}
	`)

	return resource
}

func CreateAccRecurringWindowUpdatedAttr(trigSchedPName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing recurring_window attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_recurring_window" "test" {
		scheduler_dn  = aci_trigger_scheduler.test.id
		name  = "%s"
		%s = "%s"
	}
	`, trigSchedPName, rName, attribute, value)
	return resource
}
