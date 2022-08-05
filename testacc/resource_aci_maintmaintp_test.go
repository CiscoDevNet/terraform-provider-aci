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

func TestAccAciMaintenancePolicy_Basic(t *testing.T) {
	var maintenance_policy_default models.MaintenancePolicy
	var maintenance_policy_updated models.MaintenancePolicy
	resourceName := "aci_maintenance_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMaintenancePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateMaintenancePolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMaintenancePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMaintenancePolicyExists(resourceName, &maintenance_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					// resource.TestCheckResourceAttr(resourceName, "admin_st", "untriggered"),
					resource.TestCheckResourceAttr(resourceName, "graceful", "no"),
					resource.TestCheckResourceAttr(resourceName, "ignore_compat", "no"),
					resource.TestCheckResourceAttr(resourceName, "internal_label", ""),
					resource.TestCheckResourceAttr(resourceName, "notif_cond", "notifyOnlyOnFailures"),
					resource.TestCheckResourceAttr(resourceName, "run_mode", "pauseOnlyOnFailures"),
					resource.TestCheckResourceAttr(resourceName, "version", ""),
					resource.TestCheckResourceAttr(resourceName, "version_check_override", "untriggered"),
				),
			},
			{
				Config: CreateAccMaintenancePolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMaintenancePolicyExists(resourceName, &maintenance_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_maintenance_policy"),

					resource.TestCheckResourceAttr(resourceName, "admin_st", "triggered"),

					resource.TestCheckResourceAttr(resourceName, "graceful", "yes"),

					resource.TestCheckResourceAttr(resourceName, "ignore_compat", "yes"),

					resource.TestCheckResourceAttr(resourceName, "internal_label", ""),

					resource.TestCheckResourceAttr(resourceName, "notif_cond", "notifyAlwaysBetweenSets"),

					resource.TestCheckResourceAttr(resourceName, "run_mode", "pauseAlwaysBetweenSets"),

					resource.TestCheckResourceAttr(resourceName, "version", ""),

					resource.TestCheckResourceAttr(resourceName, "version_check_override", "trigger"),

					testAccCheckAciMaintenancePolicyIdEqual(&maintenance_policy_default, &maintenance_policy_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_st"},
			},
			{
				Config:      CreateAccMaintenancePolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccMaintenancePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccMaintenancePolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMaintenancePolicyExists(resourceName, &maintenance_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciMaintenancePolicyIdNotEqual(&maintenance_policy_default, &maintenance_policy_updated),
				),
			},
		},
	})
}

func TestAccAciMaintenancePolicy_Update(t *testing.T) {
	var maintenance_policy_default models.MaintenancePolicy
	var maintenance_policy_updated models.MaintenancePolicy
	resourceName := "aci_maintenance_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMaintenancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMaintenancePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMaintenancePolicyExists(resourceName, &maintenance_policy_default),
				),
			},

			{
				Config: CreateAccMaintenancePolicyUpdatedAttr(rName, "version_check_override", "trigger-immediate"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMaintenancePolicyExists(resourceName, &maintenance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "version_check_override", "trigger-immediate"),
					testAccCheckAciMaintenancePolicyIdEqual(&maintenance_policy_default, &maintenance_policy_updated),
				),
			},
			{
				Config: CreateAccMaintenancePolicyUpdatedAttr(rName, "version_check_override", "triggered"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMaintenancePolicyExists(resourceName, &maintenance_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "version_check_override", "triggered"),
					testAccCheckAciMaintenancePolicyIdEqual(&maintenance_policy_default, &maintenance_policy_updated),
				),
			},
			{
				Config: CreateAccMaintenancePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciMaintenancePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMaintenancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMaintenancePolicyConfig(rName),
			},

			{
				Config:      CreateAccMaintenancePolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMaintenancePolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMaintenancePolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccMaintenancePolicyUpdatedAttr(rName, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccMaintenancePolicyUpdatedAttr(rName, "graceful", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccMaintenancePolicyUpdatedAttr(rName, "ignore_compat", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccMaintenancePolicyUpdatedAttr(rName, "notif_cond", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccMaintenancePolicyUpdatedAttr(rName, "run_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccMaintenancePolicyUpdatedAttr(rName, "version_check_override", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccMaintenancePolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccMaintenancePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciMaintenancePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMaintenancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMaintenancePolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciMaintenancePolicyExists(name string, maintenance_policy *models.MaintenancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Maintenance Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Maintenance Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		maintenance_policyFound := models.MaintenancePolicyFromContainer(cont)
		if maintenance_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Maintenance Policy %s not found", rs.Primary.ID)
		}
		*maintenance_policy = *maintenance_policyFound
		return nil
	}
}

func testAccCheckAciMaintenancePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing maintenance_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_maintenance_policy" {
			cont, err := client.Get(rs.Primary.ID)
			maintenance_policy := models.MaintenancePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Maintenance Policy %s Still exists", maintenance_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciMaintenancePolicyIdEqual(m1, m2 *models.MaintenancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("maintenance_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciMaintenancePolicyIdNotEqual(m1, m2 *models.MaintenancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("maintenance_policy DNs are equal")
		}
		return nil
	}
}

func CreateMaintenancePolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing maintenance_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_maintenance_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccMaintenancePolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing maintenance_policy creation with name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_maintenance_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccMaintenancePolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing maintenance_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_maintenance_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccMaintenancePolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing maintenance_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_maintenance_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccMaintenancePolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple maintenance_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_maintenance_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccMaintenancePolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing maintenance_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_maintenance_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_maintenance_policy"
		admin_st = "triggered"
		graceful = "yes"
		ignore_compat = "yes"
		internal_label = ""
		notif_cond = "notifyAlwaysBetweenSets"
		run_mode = "pauseAlwaysBetweenSets"
		version = ""
		version_check_override = "trigger"
		
	}
	`, rName)

	return resource
}

func CreateAccMaintenancePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing maintenance_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_maintenance_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_maintenance_policy"
		admin_st = "triggered"
		graceful = "yes"
		ignore_compat = "yes"
		internal_label = ""
		notif_cond = "notifyAlwaysBetweenSets"
		run_mode = "pauseAlwaysBetweenSets"
		version = ""
		version_check_override = "trigger"
		
	}
	`)

	return resource
}

func CreateAccMaintenancePolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing maintenance_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_maintenance_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
