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

func TestAccAciTriggerScheduler_Basic(t *testing.T) {
	var trigger_scheduler_default models.TriggerScheduler
	var trigger_scheduler_updated models.TriggerScheduler
	resourceName := "aci_trigger_scheduler.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTriggerSchedulerDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateTriggerSchedulerWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTriggerSchedulerConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTriggerSchedulerExists(resourceName, &trigger_scheduler_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccTriggerSchedulerConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTriggerSchedulerExists(resourceName, &trigger_scheduler_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_trigger_scheduler"),

					testAccCheckAciTriggerSchedulerIdEqual(&trigger_scheduler_default, &trigger_scheduler_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccTriggerSchedulerConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccTriggerSchedulerRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccTriggerSchedulerConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTriggerSchedulerExists(resourceName, &trigger_scheduler_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciTriggerSchedulerIdNotEqual(&trigger_scheduler_default, &trigger_scheduler_updated),
				),
			},
		},
	})
}

func TestAccAcitrigSchedP_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTriggerSchedulerDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTriggerSchedulerConfig(rName),
			},

			{
				Config:      CreateAccTriggerSchedulerUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTriggerSchedulerUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccTriggerSchedulerUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccTriggerSchedulerUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccTriggerSchedulerConfig(rName),
			},
		},
	})
}

func TestAccAciTriggerScheduler_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTriggerSchedulerDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccTriggerSchedulerConfigMultiple(rName),
			},
		},
	})
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
	fmt.Println("=== STEP  testing trigger_scheduler destroy")
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

func testAccCheckAciTriggerSchedulerIdEqual(m1, m2 *models.TriggerScheduler) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("trigger_scheduler DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciTriggerSchedulerIdNotEqual(m1, m2 *models.TriggerScheduler) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("trigger_scheduler DNs are equal")
		}
		return nil
	}
}

func CreateTriggerSchedulerWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing trigger_scheduler creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_trigger_scheduler" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccTriggerSchedulerConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing trigger_scheduler creation with name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccTriggerSchedulerConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing trigger_scheduler creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccTriggerSchedulerConfig(rName string) string {
	fmt.Println("=== STEP  testing trigger_scheduler creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccTriggerSchedulerConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple trigger_scheduler creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccTriggerSchedulerConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing trigger_scheduler creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_trigger_scheduler"
		
	}
	`, rName)

	return resource
}

func CreateAccTriggerSchedulerRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing trigger_scheduler updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_trigger_scheduler" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_trigger_scheduler"
		
	}
	`)

	return resource
}

func CreateAccTriggerSchedulerUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing trigger_scheduler attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_trigger_scheduler" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
