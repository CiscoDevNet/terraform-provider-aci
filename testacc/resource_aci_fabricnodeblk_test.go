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

func TestAccAciNodeBlockFW_Basic(t *testing.T) {
	var node_block_firmware_default models.NodeBlockFW
	var node_block_firmware_updated models.NodeBlockFW
	resourceName := "aci_node_block_firmware.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	firmwareFwGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockFWDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateNodeBlockFWWithoutRequired(firmwareFwGrpName, rName, "firmware_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateNodeBlockFWWithoutRequired(firmwareFwGrpName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccNodeBlockFWConfig(firmwareFwGrpName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockFWExists(resourceName, &node_block_firmware_default),
					resource.TestCheckResourceAttr(resourceName, "firmware_group_dn", fmt.Sprintf("uni/fabric/fwgrp-%s", firmwareFwGrpName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "from_", "1"),
					resource.TestCheckResourceAttr(resourceName, "to_", "1"),
				),
			},
			{
				Config: CreateAccNodeBlockFWConfigWithOptionalValues(firmwareFwGrpName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockFWExists(resourceName, &node_block_firmware_updated),
					resource.TestCheckResourceAttr(resourceName, "firmware_group_dn", fmt.Sprintf("uni/fabric/fwgrp-%s", firmwareFwGrpName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_node_block_firmware"),
					resource.TestCheckResourceAttr(resourceName, "from_", "2"),
					resource.TestCheckResourceAttr(resourceName, "to_", "2"),

					testAccCheckAciNodeBlockFWIdEqual(&node_block_firmware_default, &node_block_firmware_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccNodeBlockFWConfigUpdatedName(firmwareFwGrpName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccNodeBlockFWRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccNodeBlockFWConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockFWExists(resourceName, &node_block_firmware_updated),
					resource.TestCheckResourceAttr(resourceName, "firmware_group_dn", fmt.Sprintf("uni/fabric/fwgrp-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciNodeBlockFWIdNotEqual(&node_block_firmware_default, &node_block_firmware_updated),
				),
			},
			{
				Config: CreateAccNodeBlockFWConfig(firmwareFwGrpName, rName),
			},
			{
				Config: CreateAccNodeBlockFWConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockFWExists(resourceName, &node_block_firmware_updated),
					resource.TestCheckResourceAttr(resourceName, "firmware_group_dn", fmt.Sprintf("uni/fabric/fwgrp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciNodeBlockFWIdNotEqual(&node_block_firmware_default, &node_block_firmware_updated),
				),
			},
		},
	})
}

func TestAccAciNodeBlockFW_Update(t *testing.T) {
	var node_block_firmware_default models.NodeBlockFW
	var node_block_firmware_updated models.NodeBlockFW
	resourceName := "aci_node_block_firmware.test"
	rName := makeTestVariable(acctest.RandString(5))

	firmwareFwGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockFWDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccNodeBlockFWConfig(firmwareFwGrpName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockFWExists(resourceName, &node_block_firmware_default),
				),
			},
			{
				Config: CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "to_", "16000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockFWExists(resourceName, &node_block_firmware_updated),
					resource.TestCheckResourceAttr(resourceName, "to_", "16000"),
					testAccCheckAciNodeBlockFWIdEqual(&node_block_firmware_default, &node_block_firmware_updated),
				),
			},
			{
				Config: CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "from_", "16000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockFWExists(resourceName, &node_block_firmware_updated),
					resource.TestCheckResourceAttr(resourceName, "from_", "16000"),
					testAccCheckAciNodeBlockFWIdEqual(&node_block_firmware_default, &node_block_firmware_updated),
				),
			},
			{
				Config: CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "from_", "7999"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockFWExists(resourceName, &node_block_firmware_updated),
					resource.TestCheckResourceAttr(resourceName, "from_", "7999"),
					testAccCheckAciNodeBlockFWIdEqual(&node_block_firmware_default, &node_block_firmware_updated),
				),
			},
			{
				Config: CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "to_", "7999"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockFWExists(resourceName, &node_block_firmware_updated),
					resource.TestCheckResourceAttr(resourceName, "to_", "7999"),
					testAccCheckAciNodeBlockFWIdEqual(&node_block_firmware_default, &node_block_firmware_updated),
				),
			},
			{
				Config: CreateAccNodeBlockFWConfig(firmwareFwGrpName, rName),
			},
		},
	})
}

func TestAccAciNodeBlockFW_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	firmwareFwGrpName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockFWDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccNodeBlockFWConfig(firmwareFwGrpName, rName),
			},
			{
				Config:      CreateAccNodeBlockFWWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "from_", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "from_", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "from_", "16001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "to_", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "to_", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, "to_", "16001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccNodeBlockFWConfig(firmwareFwGrpName, rName),
			},
		},
	})
}

func TestAccAciNodeBlockFW_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	firmwareFwGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockFWDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccNodeBlockFWConfigMultiple(firmwareFwGrpName, rName),
			},
		},
	})
}

func testAccCheckAciNodeBlockFWExists(name string, node_block_firmware *models.NodeBlockFW) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Node Block Firmware %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Node Block Firmware dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		node_block_firmwareFound := models.NodeBlockFromContainer(cont)
		if node_block_firmwareFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Node Block Firmware %s not found", rs.Primary.ID)
		}
		*node_block_firmware = *node_block_firmwareFound
		return nil
	}
}

func testAccCheckAciNodeBlockFWDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing node_block_firmware destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_node_block_firmware" {
			cont, err := client.Get(rs.Primary.ID)
			node_block_firmware := models.NodeBlockFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Node Block Firmware %s Still exists", node_block_firmware.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciNodeBlockFWIdEqual(m1, m2 *models.NodeBlockFW) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("node_block_firmware DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciNodeBlockFWIdNotEqual(m1, m2 *models.NodeBlockFW) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("node_block_firmware DNs are equal")
		}
		return nil
	}
}

func CreateNodeBlockFWWithoutRequired(firmwareFwGrpName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing node_block_firmware creation without ", attrName)
	rBlock := `
	
	resource "aci_firmware_group" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "firmware_group_dn":
		rBlock += `
	resource "aci_node_block_firmware" "test" {
	#	firmware_group_dn  = aci_firmware_group.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, firmwareFwGrpName, rName)
}

func CreateAccNodeBlockFWConfigWithRequiredParams(firmwareFwGrpName, rName string) string {
	fmt.Printf("=== STEP  testing node_block_firmware creation with parent resource name %s and resource name %s\n", firmwareFwGrpName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = "%s"
	}
	`, firmwareFwGrpName, rName)
	return resource
}
func CreateAccNodeBlockFWConfigUpdatedName(firmwareFwGrpName, rName string) string {
	fmt.Println("=== STEP  testing node_block_firmware creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = "%s"
	}
	`, firmwareFwGrpName, rName)
	return resource
}

func CreateAccNodeBlockFWConfig(firmwareFwGrpName, rName string) string {
	fmt.Println("=== STEP  testing node_block_firmware creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = "%s"
	}
	`, firmwareFwGrpName, rName)
	return resource
}

func CreateAccNodeBlockFWConfigMultiple(firmwareFwGrpName, rName string) string {
	fmt.Println("=== STEP  testing multiple node_block_firmware creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = "%s_${count.index}"
		from_ = (count.index+1)*10 
		to_ = (count.index+1)*10+5
		count = 5
	}
	`, firmwareFwGrpName, rName)
	return resource
}

func CreateAccNodeBlockFWWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing node_block_firmware creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccNodeBlockFWConfigWithOptionalValues(firmwareFwGrpName, rName string) string {
	fmt.Println("=== STEP  Basic: testing node_block_firmware creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = "${aci_firmware_group.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_node_block_firmware"
		from_ = "2"
		to_ = "2"
		
	}
	`, firmwareFwGrpName, rName)

	return resource
}

func CreateAccNodeBlockFWRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing node_block_firmware updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_node_block_firmware" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_node_block_firmware"
		from_ = "2"
		to_ = "2"
		
	}
	`)

	return resource
}

func CreateAccNodeBlockFWUpdatedAttr(firmwareFwGrpName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing node_block_firmware attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_firmware_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_node_block_firmware" "test" {
		firmware_group_dn  = aci_firmware_group.test.id
		name  = "%s"
		%s = "%s"
	}
	`, firmwareFwGrpName, rName, attribute, value)
	return resource
}
