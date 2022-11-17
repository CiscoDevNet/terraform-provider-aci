package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciTriggerScheduler_Basic(t *testing.T) {
	var trigger_scheduler models.TriggerScheduler
	description := "trigger_scheduler created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTriggerSchedulerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTriggerSchedulerConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTriggerSchedulerExists("aci_trigger_scheduler.footrigger_scheduler", &trigger_scheduler),
					testAccCheckAciTriggerSchedulerAttributes(description, &trigger_scheduler),
				),
			},
		},
	})
}

func TestAccAciTriggerScheduler_update(t *testing.T) {
	var trigger_scheduler models.TriggerScheduler
	description := "trigger_scheduler created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTriggerSchedulerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTriggerSchedulerConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTriggerSchedulerExists("aci_trigger_scheduler.footrigger_scheduler", &trigger_scheduler),
					testAccCheckAciTriggerSchedulerAttributes(description, &trigger_scheduler),
				),
			},
			{
				Config: testAccCheckAciTriggerSchedulerConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTriggerSchedulerExists("aci_trigger_scheduler.footrigger_scheduler", &trigger_scheduler),
					testAccCheckAciTriggerSchedulerAttributes(description, &trigger_scheduler),
				),
			},
		},
	})
}

func testAccCheckAciTriggerSchedulerConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_trigger_scheduler" "footrigger_scheduler" {
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciTriggerSchedulerExists(name string, trigger_scheduler *models.TriggerScheduler) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Trigger Scheduler %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Trigger Scheduler dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		trigger_schedulerFound := models.TriggerSchedulerFromContainer(cont)
		if trigger_schedulerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Trigger Scheduler %s not found", rs.Primary.ID)
		}
		*trigger_scheduler = *trigger_schedulerFound
		return nil
	}
}

func testAccCheckAciTriggerSchedulerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_trigger_scheduler" {
			cont, err := client.Get(rs.Primary.ID)
			trigger_scheduler := models.TriggerSchedulerFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Trigger Scheduler %s Still exists", trigger_scheduler.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciTriggerSchedulerAttributes(description string, trigger_scheduler *models.TriggerScheduler) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != trigger_scheduler.Description {
			return fmt.Errorf("Bad trigger_scheduler Description %s", trigger_scheduler.Description)
		}

		if "example" != trigger_scheduler.Name {
			return fmt.Errorf("Bad trigger_scheduler name %s", trigger_scheduler.Name)
		}

		if "example" != trigger_scheduler.Annotation {
			return fmt.Errorf("Bad trigger_scheduler annotation %s", trigger_scheduler.Annotation)
		}

		if "example" != trigger_scheduler.NameAlias {
			return fmt.Errorf("Bad trigger_scheduler name_alias %s", trigger_scheduler.NameAlias)
		}

		return nil
	}
}
