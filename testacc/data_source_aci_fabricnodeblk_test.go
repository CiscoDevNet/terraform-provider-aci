package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciNodeBlockFirmwareDataSource_Basic(t *testing.T) {
	resourceName := "aci_node_block_firmware.test"
	dataSourceName := "data.aci_node_block_firmware.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	firmwareFwGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockFWDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateNodeBlockFirmwareDSWithoutRequired(firmwareFwGrpName, rName, "firmware_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateNodeBlockFirmwareDSWithoutRequired(firmwareFwGrpName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccNodeBlockFirmwareConfigDataSource(firmwareFwGrpName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "firmware_group_dn", resourceName, "firmware_group_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "from_", resourceName, "from_"),
					resource.TestCheckResourceAttrPair(dataSourceName, "to_", resourceName, "to_"),
				),
			},
			{
				Config:      CreateAccNodeBlockFirmwareDataSourceUpdate(firmwareFwGrpName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccNodeBlockFirmwareDSWithInvalidParentDn(firmwareFwGrpName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccNodeBlockFirmwareDataSourceUpdatedResource(firmwareFwGrpName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccNodeBlockFirmwareConfigDataSource(firmwareFwGrpName, rName string) string {
	fmt.Println("=== STEP  testing node_block_firmware Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = "%s"
	}

	data "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = aci_node_block_firmware.test.name
		depends_on = [ aci_node_block_firmware.test ]
	}
	`, firmwareFwGrpName, rName)
	return resource
}

func CreateNodeBlockFirmwareDSWithoutRequired(firmwareFwGrpName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing node_block_firmware Data Source without ", attrName)
	rBlock := `
	
	resource "aci_firmware_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "firmware_group_dn":
		rBlock += `
	data "aci_node_block_firmware" "test" {
	#	firmware_group_dn  = aci_firmware_group.test.id
		name  = aci_node_block_firmware.test.name
		depends_on = [ aci_node_block_firmware.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
	#	name  = aci_node_block_firmware.test.name
		depends_on = [ aci_node_block_firmware.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, firmwareFwGrpName, rName)
}

func CreateAccNodeBlockFirmwareDSWithInvalidParentDn(firmwareFwGrpName, rName string) string {
	fmt.Println("=== STEP  testing node_block_firmware Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = "%s"
	}

	data "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = "${aci_node_block_firmware.test.name}_invalid"
		depends_on = [ aci_node_block_firmware.test ]
	}
	`, firmwareFwGrpName, rName)
	return resource
}

func CreateAccNodeBlockFirmwareDataSourceUpdate(firmwareFwGrpName, rName, key, value string) string {
	fmt.Println("=== STEP  testing node_block_firmware Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = "%s"
	}

	data "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = aci_node_block_firmware.test.name
		%s = "%s"
		depends_on = [ aci_node_block_firmware.test ]
	}
	`, firmwareFwGrpName, rName, key, value)
	return resource
}

func CreateAccNodeBlockFirmwareDataSourceUpdatedResource(firmwareFwGrpName, rName, key, value string) string {
	fmt.Println("=== STEP  testing node_block_firmware Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = aci_node_block_firmware.test.name
		depends_on = [ aci_node_block_firmware.test ]
	}
	`, firmwareFwGrpName, rName, key, value)
	return resource
}
