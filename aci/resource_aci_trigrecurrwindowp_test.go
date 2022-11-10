package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciRecurringWindow_Basic(t *testing.T) {
	var recurring_window models.RecurringWindow
	annotation := "testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRecurringWindowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRecurringWindowConfig_basic(annotation),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists("aci_recurring_window.test", &recurring_window),
					testAccCheckAciRecurringWindowAttributes(annotation, &recurring_window),
				),
			},
		},
	})
}

func TestAccAciRecurringWindowTrigger_update(t *testing.T) {
	var recurring_window models.RecurringWindow
	annotation := "testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRecurringWindowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRecurringWindowConfig_basic(annotation),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists("aci_recurring_window.test", &recurring_window),
					testAccCheckAciRecurringWindowAttributes(annotation, &recurring_window),
				),
			},
			{
				Config: testAccCheckAciRecurringWindowConfig_basic(annotation),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRecurringWindowExists("aci_recurring_window.test", &recurring_window),
					testAccCheckAciRecurringWindowAttributes(annotation, &recurring_window),
				),
			},
		},
	})
}

func testAccCheckAciRecurringWindowConfig_basic(annotation string) string {
	return fmt.Sprintf(`

	resource "aci_recurring_window" "test" {
		name 		= "test"
		scheduler_dn = "uni/fabric/schedp-dm"
		concur_cap = "unlimited"
		day = "every-day"
		hour = "0"
		minute = "0"
		node_upg_interval = "0"
		proc_break = "none"
		proc_cap = "unlimited"
		time_cap = "unlimited"
		annotation = "%s"
		name_alias = "example"
	}

	`, annotation)
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

func testAccCheckAciRecurringWindowAttributes(annotation string, recurring_window *models.RecurringWindow) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != GetMOName(recurring_window.DistinguishedName) {
			return fmt.Errorf("Bad trig_recurr_window_p %s", GetMOName(recurring_window.DistinguishedName))
		}

		if "unlimited" != recurring_window.ConcurCap {
			return fmt.Errorf("Bad recurring_window ConcurCap %s", recurring_window.ConcurCap)
		}

		if "every-day" != recurring_window.Day {
			return fmt.Errorf("Bad recurring_window Day %s", recurring_window.Day)
		}

		if "0" != recurring_window.Hour {
			return fmt.Errorf("Bad recurring_window Hour %s", recurring_window.Hour)
		}

		if "0" != recurring_window.Minute {
			return fmt.Errorf("Bad recurring_window Minute %s", recurring_window.Minute)
		}

		if "0" != recurring_window.NodeUpgInterval {
			return fmt.Errorf("Bad recurring_window NodeUpgInterval %s", recurring_window.NodeUpgInterval)
		}

		if "none" != recurring_window.ProcBreak {
			return fmt.Errorf("Bad recurring_window ProcBreak %s", recurring_window.ProcBreak)
		}

		if "unlimited" != recurring_window.ProcCap {
			return fmt.Errorf("Bad recurring_window ProcCap %s", recurring_window.ProcCap)
		}

		if "unlimited" != recurring_window.TimeCap {
			return fmt.Errorf("Bad recurring_window TimeCap %s", recurring_window.TimeCap)
		}

		if annotation != recurring_window.Annotation {
			return fmt.Errorf("Bad recurring_window Annotation %s", recurring_window.Annotation)
		}

		if "example" != recurring_window.NameAlias {
			return fmt.Errorf("Bad recurring_window NameAlias %s", recurring_window.NameAlias)
		}
		return nil
	}
}
