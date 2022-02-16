package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFirmwarePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_firmware_policy.test"
	dataSourceName := "data.aci_firmware_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFirmwarePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFirmwarePolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFirmwarePolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "effective_on_reboot", resourceName, "effective_on_reboot"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ignore_compat", resourceName, "ignore_compat"),
					resource.TestCheckResourceAttrPair(dataSourceName, "internal_label", resourceName, "internal_label"),
					resource.TestCheckResourceAttrPair(dataSourceName, "version", resourceName, "version"),
					resource.TestCheckResourceAttrPair(dataSourceName, "version_check_override", resourceName, "version_check_override"),
				),
			},
			{
				Config:      CreateAccFirmwarePolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccFirmwarePolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccFirmwarePolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFirmwarePolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing firmware_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_policy" "test" {
	
		name  = "%s"
	}

	data "aci_firmware_policy" "test" {
	
		name  = aci_firmware_policy.test.name
		depends_on = [ aci_firmware_policy.test ]
	}
	`, rName)
	return resource
}

func CreateFirmwarePolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing firmware_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_firmware_policy" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_firmware_policy" "test" {
	
	#	name  = aci_firmware_policy.test.name
		depends_on = [ aci_firmware_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFirmwarePolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing firmware_policy Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_policy" "test" {
	
		name  = "%s"
	}

	data "aci_firmware_policy" "test" {
	
		name  = "${aci_firmware_policy.test.name}_invalid"
		depends_on = [ aci_firmware_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccFirmwarePolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing firmware_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_policy" "test" {
	
		name  = "%s"
	}

	data "aci_firmware_policy" "test" {
	
		name  = aci_firmware_policy.test.name
		%s = "%s"
		depends_on = [ aci_firmware_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccFirmwarePolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing firmware_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_firmware_policy" "test" {
	
		name  = aci_firmware_policy.test.name
		depends_on = [ aci_firmware_policy.test ]
	}
	`, rName, key, value)
	return resource
}
