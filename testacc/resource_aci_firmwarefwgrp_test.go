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

func TestAccAciFirmwareGroup_Basic(t *testing.T) {
	var firmware_group_default models.FirmwareGroup
	var firmware_group_updated models.FirmwareGroup
	resourceName := "aci_firmware_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwareGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFirmwareGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFirmwareGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareGroupExists(resourceName, &firmware_group_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "firmware_group_type", "range"),
				),
			},
			{
				Config: CreateAccFirmwareGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareGroupExists(resourceName, &firmware_group_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_firmware_group"),
					resource.TestCheckResourceAttr(resourceName, "firmware_group_type", "ALL"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFirmwareGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccFirmwareGroupConfigInvalidName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccFirmwareGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccFirmwareGroupConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareGroupExists(resourceName, &firmware_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciFirmwareGroupIdNotEqual(&firmware_group_default, &firmware_group_updated),
				),
			},
		},
	})
}

func TestAccAciFirmwareGroup_Update(t *testing.T) {
	var firmware_group_default models.FirmwareGroup
	var firmware_group_updated models.FirmwareGroup
	resourceName := "aci_firmware_group.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwareGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFirmwareGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareGroupExists(resourceName, &firmware_group_default),
				),
			},

			{
				Config: CreateAccFirmwareGroupUpdatedAttr(rName, "firmware_group_type", "ALL_IN_POD"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareGroupExists(resourceName, &firmware_group_updated),
					resource.TestCheckResourceAttr(resourceName, "firmware_group_type", "ALL_IN_POD"),
					testAccCheckAciFirmwareGroupIdEqual(&firmware_group_default, &firmware_group_updated),
				),
			},
			{
				Config: CreateAccFirmwareGroupUpdatedAttr(rName, "firmware_group_type", "range"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareGroupExists(resourceName, &firmware_group_updated),
					resource.TestCheckResourceAttr(resourceName, "firmware_group_type", "range"),
					testAccCheckAciFirmwareGroupIdEqual(&firmware_group_default, &firmware_group_updated),
				),
			},
			{
				Config: CreateAccFirmwareGroupConfig(rName),
			},
		},
	})
}

func TestAccAciFirmwareGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwareGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFirmwareGroupConfig(rName),
			},

			{
				Config:      CreateAccFirmwareGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFirmwareGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFirmwareGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFirmwareGroupUpdatedAttr(rName, "firmware_group_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},

			{
				Config:      CreateAccFirmwareGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)*is not expected here.`),
			},
			{
				Config: CreateAccFirmwareGroupConfig(rName),
			},
		},
	})
}

func TestAccAciFirmwareGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwareGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFirmwareGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciFirmwareGroupExists(name string, firmware_group *models.FirmwareGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Firmware Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Firmware Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		firmware_groupFound := models.FirmwareGroupFromContainer(cont)
		if firmware_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Firmware Group %s not found", rs.Primary.ID)
		}
		*firmware_group = *firmware_groupFound
		return nil
	}
}

func testAccCheckAciFirmwareGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing firmware_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_firmware_group" {
			cont, err := client.Get(rs.Primary.ID)
			firmware_group := models.FirmwareGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Firmware Group %s Still exists", firmware_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFirmwareGroupIdEqual(m1, m2 *models.FirmwareGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("firmware_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFirmwareGroupIdNotEqual(m1, m2 *models.FirmwareGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("firmware_group DNs are equal")
		}
		return nil
	}
}

func CreateFirmwareGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing firmware_group creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_firmware_group" "test" {
	
	#	name  = "%s"
		description = "created while acceptance testing"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFirmwareGroupConfigInvalidName(rName string) string {
	fmt.Println("=== STEP  testing firmware_group creation with invalid name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFirmwareGroupConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing firmware_group creation with resource name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFirmwareGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing firmware_group creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFirmwareGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing firmware_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_firmware_group"
		firmware_group_type = "ALL"
	}
	`, rName)

	return resource
}

func CreateAccFirmwareGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing firmware_group update without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_firmware_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_firmware_group"
		firmware_group_type = "ALL"
	}
	`)

	return resource
}

func CreateAccFirmwareGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing firmware_group attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccFirmwareGroupConfigUpdatedName(longerName string) string {
	fmt.Println("=== STEP  Basic: testing Firmware Group creation with invalid name with long length")
	resource := fmt.Sprintf(`
	resource "aci_firmware_group" "test" {
	name  = "%s"
	}
	`, longerName)
	return resource
}

func CreateAccFirmwareGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple firmware_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test1" {
		name  = "%s"
	}

	resource "aci_firmware_group" "test2" {
		name  = "%s"
	}

	resource "aci_firmware_group" "test3" {
		name  = "%s"
	}
	`, rName+"1", rName+"2", rName+"3")
	return resource
}
