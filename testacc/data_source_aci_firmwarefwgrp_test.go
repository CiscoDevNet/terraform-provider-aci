package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFirmwareGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_firmware_group.test"
	dataSourceName := "data.aci_firmware_group.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwareGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFirmwareGroupDSWithoutRequired(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFirmwareGroupConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "firmware_group_type", resourceName, "firmware_group_type"),
				),
			},
			{
				Config:      CreateAccFirmwareGroupDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccFirmwareGroupDSConfigWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccFirmwareGroupDataSourceUpdate(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFirmwareGroupConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing firmware_group creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
	
		name  = "%s"
	}

	data "aci_firmware_group" "test" {
	
		name  = aci_firmware_group.test.name
		depends_on = [
			aci_firmware_group.test
		]
	}
	`, rName)
	return resource
}

func CreateAccFirmwareGroupDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing firmware_group Data Source with random attributes")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
	
		name  = "%s"
	}

	data "aci_firmware_group" "test" {
	
		name  = aci_firmware_group.test.name
		%s = "%s"
		depends_on = [
			aci_firmware_group.test
		]
	}
	`, rName, key, value)
	return resource
}

func CreateFirmwareGroupDSWithoutRequired(rName string) string {
	fmt.Println("=== STEP  Basic: testing firmware_group Data Source without name")
	resource := fmt.Sprintf(`

	resource "aci_firmware_group" "test" {
		name  = "%s"
	}
	  data "aci_firmware_group" "test" {
	    depends_on = [
			aci_firmware_group.test
		]
	  }
	`, rName)
	return resource
}

func CreateAccFirmwareGroupDSConfigWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing firmware_group data source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
		name  = "%s"
	}

	data "aci_firmware_group" "test" {
		name  = "${aci_firmware_group.test.name}xyz"
	}
	`, rName)
	return resource
}
