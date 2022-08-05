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

func TestAccAciFirmwarePolicy_Basic(t *testing.T) {
	var firmware_policy_default models.FirmwarePolicy
	var firmware_policy_updated models.FirmwarePolicy
	resourceName := "aci_firmware_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwarePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFirmwarePolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFirmwarePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwarePolicyExists(resourceName, &firmware_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "effective_on_reboot", "no"),
					resource.TestCheckResourceAttr(resourceName, "ignore_compat", "no"),
					resource.TestCheckResourceAttr(resourceName, "internal_label", ""),
					resource.TestCheckResourceAttr(resourceName, "version", ""),
					resource.TestCheckResourceAttr(resourceName, "version_check_override", "untriggered"),
				),
			},
			{
				Config: CreateAccFirmwarePolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwarePolicyExists(resourceName, &firmware_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_firmware_policy"),
					resource.TestCheckResourceAttr(resourceName, "effective_on_reboot", "yes"),
					resource.TestCheckResourceAttr(resourceName, "ignore_compat", "yes"),
					resource.TestCheckResourceAttr(resourceName, "internal_label", "internal_label_test"),
					resource.TestCheckResourceAttr(resourceName, "version", "n9000-14.2(3q)"),
					resource.TestCheckResourceAttr(resourceName, "version_check_override", "trigger"),
					testAccCheckAciFirmwarePolicyIdEqual(&firmware_policy_default, &firmware_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFirmwarePolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccFirmwarePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccFirmwarePolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwarePolicyExists(resourceName, &firmware_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciFirmwarePolicyIdNotEqual(&firmware_policy_default, &firmware_policy_updated),
				),
			},
		},
	})
}

func TestAccAciFirmwarePolicy_Update(t *testing.T) {
	var firmware_policy_default models.FirmwarePolicy
	var firmware_policy_updated models.FirmwarePolicy
	resourceName := "aci_firmware_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwarePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFirmwarePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwarePolicyExists(resourceName, &firmware_policy_default),
				),
			},

			{
				Config: CreateAccFirmwarePolicyUpdatedAttr(rName, "version_check_override", "trigger-immediate"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwarePolicyExists(resourceName, &firmware_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "version_check_override", "trigger-immediate"),
					testAccCheckAciFirmwarePolicyIdEqual(&firmware_policy_default, &firmware_policy_updated),
				),
			},
			{
				Config: CreateAccFirmwarePolicyUpdatedAttr(rName, "version_check_override", "triggered"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwarePolicyExists(resourceName, &firmware_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "version_check_override", "triggered"),
					testAccCheckAciFirmwarePolicyIdEqual(&firmware_policy_default, &firmware_policy_updated),
				),
			},
			{
				Config: CreateAccFirmwarePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciFirmwarePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwarePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFirmwarePolicyConfig(rName),
			},

			{
				Config:      CreateAccFirmwarePolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFirmwarePolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFirmwarePolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccFirmwarePolicyUpdatedAttr(rName, "effective_on_reboot", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFirmwarePolicyUpdatedAttr(rName, "ignore_compat", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFirmwarePolicyUpdatedAttr(rName, "internal_label", acctest.RandString(513)),
				ExpectError: regexp.MustCompile(`failed validation for value`),
			},

			{
				Config:      CreateAccFirmwarePolicyUpdatedAttr(rName, "version", acctest.RandString(513)),
				ExpectError: regexp.MustCompile(`failed validation for value`),
			},

			{
				Config:      CreateAccFirmwarePolicyUpdatedAttr(rName, "version_check_override", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFirmwarePolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFirmwarePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciFirmwarePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwarePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFirmwarePolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciFirmwarePolicyExists(name string, firmware_policy *models.FirmwarePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Firmware Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Firmware Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		firmware_policyFound := models.FirmwarePolicyFromContainer(cont)
		if firmware_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Firmware Policy %s not found", rs.Primary.ID)
		}
		*firmware_policy = *firmware_policyFound
		return nil
	}
}

func testAccCheckAciFirmwarePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing firmware_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_firmware_policy" {
			cont, err := client.Get(rs.Primary.ID)
			firmware_policy := models.FirmwarePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Firmware Policy %s Still exists", firmware_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFirmwarePolicyIdEqual(m1, m2 *models.FirmwarePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("firmware_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFirmwarePolicyIdNotEqual(m1, m2 *models.FirmwarePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("firmware_policy DNs are equal")
		}
		return nil
	}
}

func CreateFirmwarePolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing firmware_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_firmware_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFirmwarePolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing firmware_policy creation with name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccFirmwarePolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing firmware_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFirmwarePolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing firmware_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFirmwarePolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple firmware_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccFirmwarePolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing firmware_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_firmware_policy"
		effective_on_reboot = "yes"
		ignore_compat = "yes"
		internal_label = "internal_label_test"
		version = "n9000-14.2(3q)"
		version_check_override = "trigger"
		
	}
	`, rName)

	return resource
}

func CreateAccFirmwarePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing firmware_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_firmware_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_firmware_policy"
		effective_on_reboot = "yes"
		ignore_compat = "yes"
		internal_label = ""
		version = ""
		version_check_override = "trigger"
		
	}
	`)

	return resource
}

func CreateAccFirmwarePolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing firmware_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccFirmwarePolicyUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing firmware_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_policy" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
