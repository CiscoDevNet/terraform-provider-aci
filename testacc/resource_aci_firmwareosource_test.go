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

func TestAccAciFirmwareDownloadTask_Basic(t *testing.T) {
	var firmware_download_task_default models.FirmwareDownloadTask
	var firmware_download_task_updated models.FirmwareDownloadTask
	resourceName := "aci_firmware_download_task.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwareDownloadTaskDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFirmwareDownloadTaskWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFirmwareDownloadTaskConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareDownloadTaskExists(resourceName, &firmware_download_task_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "auth_pass", "password"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "usePassword"),
					resource.TestCheckResourceAttr(resourceName, "dnld_task_flip", "yes"),
					// resource.TestCheckResourceAttr(resourceName, "identity_private_key_contents", ""),
					// resource.TestCheckResourceAttr(resourceName, "identity_private_key_passphrase", ""),
					// resource.TestCheckResourceAttr(resourceName, "identity_public_key_contents", ""),
					resource.TestCheckResourceAttr(resourceName, "load_catalog_if_exists_and_newer", "yes"),
					// resource.TestCheckResourceAttr(resourceName, "password", ""),
					resource.TestCheckResourceAttr(resourceName, "polling_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "proto", "scp"),
					resource.TestCheckResourceAttr(resourceName, "url", ""),
					resource.TestCheckResourceAttr(resourceName, "user", ""),
				),
			},
			{
				Config: CreateAccFirmwareDownloadTaskConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareDownloadTaskExists(resourceName, &firmware_download_task_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_firmware_download_task"),

					resource.TestCheckResourceAttr(resourceName, "auth_pass", "key"),

					resource.TestCheckResourceAttr(resourceName, "auth_type", "useSshKeyContents"),

					resource.TestCheckResourceAttr(resourceName, "dnld_task_flip", "no"),

					resource.TestCheckResourceAttr(resourceName, "identity_private_key_contents", ""),

					resource.TestCheckResourceAttr(resourceName, "identity_private_key_passphrase", ""),

					resource.TestCheckResourceAttr(resourceName, "identity_public_key_contents", ""),

					resource.TestCheckResourceAttr(resourceName, "load_catalog_if_exists_and_newer", "no"),

					resource.TestCheckResourceAttr(resourceName, "password", ""),

					resource.TestCheckResourceAttr(resourceName, "proto", "http"),

					resource.TestCheckResourceAttr(resourceName, "url", ""),

					resource.TestCheckResourceAttr(resourceName, "user", ""),

					testAccCheckAciFirmwareDownloadTaskIdEqual(&firmware_download_task_default, &firmware_download_task_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"identity_private_key_contents",
					"identity_private_key_passphrase",
					"identity_public_key_contents",
					"password",
				},
			},
			{
				Config:      CreateAccFirmwareDownloadTaskConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccFirmwareDownloadTaskRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccFirmwareDownloadTaskConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareDownloadTaskExists(resourceName, &firmware_download_task_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciFirmwareDownloadTaskIdNotEqual(&firmware_download_task_default, &firmware_download_task_updated),
				),
			},
		},
	})
}

func TestAccAciFirmwareDownloadTask_Update(t *testing.T) {
	var firmware_download_task_default models.FirmwareDownloadTask
	var firmware_download_task_updated models.FirmwareDownloadTask
	resourceName := "aci_firmware_download_task.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwareDownloadTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFirmwareDownloadTaskConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareDownloadTaskExists(resourceName, &firmware_download_task_default),
				),
			},

			{
				Config: CreateAccFirmwareDownloadTaskUpdatedAttr(rName, "proto", "local"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareDownloadTaskExists(resourceName, &firmware_download_task_updated),
					resource.TestCheckResourceAttr(resourceName, "proto", "local"),
					testAccCheckAciFirmwareDownloadTaskIdEqual(&firmware_download_task_default, &firmware_download_task_updated),
				),
			},
			{
				Config: CreateAccFirmwareDownloadTaskUpdatedAttr(rName, "proto", "usbkey"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwareDownloadTaskExists(resourceName, &firmware_download_task_updated),
					resource.TestCheckResourceAttr(resourceName, "proto", "usbkey"),
					testAccCheckAciFirmwareDownloadTaskIdEqual(&firmware_download_task_default, &firmware_download_task_updated),
				),
			},
			{
				Config: CreateAccFirmwareDownloadTaskConfig(rName),
			},
		},
	})
}

func TestAccAciFirmwareDownloadTask_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwareDownloadTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFirmwareDownloadTaskConfig(rName),
			},

			{
				Config:      CreateAccFirmwareDownloadTaskUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFirmwareDownloadTaskUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFirmwareDownloadTaskUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccFirmwareDownloadTaskUpdatedAttr(rName, "auth_pass", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFirmwareDownloadTaskUpdatedAttr(rName, "auth_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFirmwareDownloadTaskUpdatedAttr(rName, "dnld_task_flip", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFirmwareDownloadTaskUpdatedAttr(rName, "load_catalog_if_exists_and_newer", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFirmwareDownloadTaskUpdatedAttr(rName, "polling_interval", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccFirmwareDownloadTaskUpdatedAttr(rName, "proto", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFirmwareDownloadTaskUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFirmwareDownloadTaskConfig(rName),
			},
		},
	})
}

func TestAccAciFirmwareDownloadTask_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwareDownloadTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFirmwareDownloadTaskConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciFirmwareDownloadTaskExists(name string, firmware_download_task *models.FirmwareDownloadTask) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Firmware Download Task %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Firmware Download Task dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		firmware_download_taskFound := models.FirmwareDownloadTaskFromContainer(cont)
		if firmware_download_taskFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Firmware Download Task %s not found", rs.Primary.ID)
		}
		*firmware_download_task = *firmware_download_taskFound
		return nil
	}
}

func testAccCheckAciFirmwareDownloadTaskDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing firmware_download_task destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_firmware_download_task" {
			cont, err := client.Get(rs.Primary.ID)
			firmware_download_task := models.FirmwareDownloadTaskFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Firmware Download Task %s Still exists", firmware_download_task.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFirmwareDownloadTaskIdEqual(m1, m2 *models.FirmwareDownloadTask) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("firmware_download_task DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFirmwareDownloadTaskIdNotEqual(m1, m2 *models.FirmwareDownloadTask) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("firmware_download_task DNs are equal")
		}
		return nil
	}
}

func CreateFirmwareDownloadTaskWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing firmware_download_task creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_firmware_download_task" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFirmwareDownloadTaskConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing firmware_download_task creation with resource name =", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_download_task" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccFirmwareDownloadTaskConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing firmware_download_task creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_download_task" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFirmwareDownloadTaskConfig(rName string) string {
	fmt.Println("=== STEP  testing firmware_download_task creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_download_task" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFirmwareDownloadTaskConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple firmware_download_task creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_download_task" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccFirmwareDownloadTaskConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing firmware_download_task creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_download_task" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_firmware_download_task"
		auth_pass = "key"
		auth_type = "useSshKeyContents"
		dnld_task_flip = "no"
		identity_private_key_contents = ""
		identity_private_key_passphrase = ""
		identity_public_key_contents = ""
		load_catalog_if_exists_and_newer = "no"
		password = ""
		proto = "http"
		url = ""
		user = ""
		
	}
	`, rName)

	return resource
}

func CreateAccFirmwareDownloadTaskRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing firmware_download_task updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_firmware_download_task" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_firmware_download_task"
		auth_pass = "key"
		auth_type = "useSshKeyContents"
		dnld_task_flip = "no"
		identity_private_key_contents = ""
		identity_private_key_passphrase = ""
		identity_public_key_contents = ""
		load_catalog_if_exists_and_newer = "no"
		password = ""
		proto = "http"
		url = ""
		user = ""
		
	}
	`)

	return resource
}

func CreateAccFirmwareDownloadTaskUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing firmware_download_task attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_download_task" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
