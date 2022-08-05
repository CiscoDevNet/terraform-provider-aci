package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFirmwareDownloadTaskDataSource_Basic(t *testing.T) {
	resourceName := "aci_firmware_download_task.test"
	dataSourceName := "data.aci_firmware_download_task.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwareDownloadTaskDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFirmwareDownloadTaskDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFirmwareDownloadTaskConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_pass", resourceName, "auth_pass"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_type", resourceName, "auth_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dnld_task_flip", resourceName, "dnld_task_flip"),
					resource.TestCheckResourceAttrPair(dataSourceName, "identity_private_key_contents", resourceName, "identity_private_key_contents"),
					resource.TestCheckResourceAttrPair(dataSourceName, "identity_private_key_passphrase", resourceName, "identity_private_key_passphrase"),
					resource.TestCheckResourceAttrPair(dataSourceName, "identity_public_key_contents", resourceName, "identity_public_key_contents"),
					resource.TestCheckResourceAttrPair(dataSourceName, "load_catalog_if_exists_and_newer", resourceName, "load_catalog_if_exists_and_newer"),
					resource.TestCheckResourceAttrPair(dataSourceName, "password", resourceName, "password"),
					resource.TestCheckResourceAttrPair(dataSourceName, "polling_interval", resourceName, "polling_interval"),
					resource.TestCheckResourceAttrPair(dataSourceName, "proto", resourceName, "proto"),
					resource.TestCheckResourceAttrPair(dataSourceName, "url", resourceName, "url"),
					resource.TestCheckResourceAttrPair(dataSourceName, "user", resourceName, "user"),
				),
			},
			{
				Config:      CreateAccFirmwareDownloadTaskDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccFirmwareDownloadTaskDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccFirmwareDownloadTaskDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFirmwareDownloadTaskConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing firmware_download_task Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_download_task" "test" {
	
		name  = "%s"
	}

	data "aci_firmware_download_task" "test" {
	
		name  = aci_firmware_download_task.test.name
		depends_on = [ aci_firmware_download_task.test ]
	}
	`, rName)
	return resource
}

func CreateFirmwareDownloadTaskDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing firmware_download_task Data Source without ", attrName)
	rBlock := `
	
	resource "aci_firmware_download_task" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_firmware_download_task" "test" {
	
	#	name  = aci_firmware_download_task.test.name
		depends_on = [ aci_firmware_download_task.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFirmwareDownloadTaskDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing firmware_download_task Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_download_task" "test" {
	
		name  = "%s"
	}

	data "aci_firmware_download_task" "test" {
	
		name  = "${aci_firmware_download_task.test.name}_invalid"
		depends_on = [ aci_firmware_download_task.test ]
	}
	`, rName)
	return resource
}

func CreateAccFirmwareDownloadTaskDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing firmware_download_task Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_download_task" "test" {
	
		name  = "%s"
	}

	data "aci_firmware_download_task" "test" {
	
		name  = aci_firmware_download_task.test.name
		%s = "%s"
		depends_on = [ aci_firmware_download_task.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccFirmwareDownloadTaskDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing firmware_download_task Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_download_task" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_firmware_download_task" "test" {
	
		name  = aci_firmware_download_task.test.name
		depends_on = [ aci_firmware_download_task.test ]
	}
	`, rName, key, value)
	return resource
}
