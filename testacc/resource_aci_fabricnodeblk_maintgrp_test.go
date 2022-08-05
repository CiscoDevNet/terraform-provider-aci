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

func TestAccAciNodeBlockMG_Basic(t *testing.T) {
	var maintenance_group_node_default models.NodeBlockMG
	var maintenance_group_node_updated models.NodeBlockMG
	resourceName := "aci_maintenance_group_node.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	maintMaintGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockMGDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateNodeBlockMGWithoutRequired(maintMaintGrpName, rName, "pod_maintenance_group_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateNodeBlockMGWithoutRequired(maintMaintGrpName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccNodeBlockMGConfig(maintMaintGrpName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockMGExists(resourceName, &maintenance_group_node_default),
					resource.TestCheckResourceAttr(resourceName, "pod_maintenance_group_dn", fmt.Sprintf("uni/fabric/maintgrp-%s", maintMaintGrpName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "from_", "1"),
					resource.TestCheckResourceAttr(resourceName, "to_", "1"),
				),
			},
			{
				Config: CreateAccNodeBlockMGConfigWithOptionalValues(maintMaintGrpName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockMGExists(resourceName, &maintenance_group_node_updated),
					resource.TestCheckResourceAttr(resourceName, "pod_maintenance_group_dn", fmt.Sprintf("uni/fabric/maintgrp-%s", maintMaintGrpName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_maintenance_group_node"),
					resource.TestCheckResourceAttr(resourceName, "from_", "2"),
					resource.TestCheckResourceAttr(resourceName, "to_", "2"),

					testAccCheckAciNodeBlockMGIdEqual(&maintenance_group_node_default, &maintenance_group_node_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccNodeBlockMGConfigUpdatedName(maintMaintGrpName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccNodeBlockMGRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccNodeBlockMGConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockMGExists(resourceName, &maintenance_group_node_updated),
					resource.TestCheckResourceAttr(resourceName, "pod_maintenance_group_dn", fmt.Sprintf("uni/fabric/maintgrp-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciNodeBlockMGIdNotEqual(&maintenance_group_node_default, &maintenance_group_node_updated),
				),
			},
			{
				Config: CreateAccNodeBlockMGConfig(maintMaintGrpName, rName),
			},
			{
				Config: CreateAccNodeBlockMGConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockMGExists(resourceName, &maintenance_group_node_updated),
					resource.TestCheckResourceAttr(resourceName, "pod_maintenance_group_dn", fmt.Sprintf("uni/fabric/maintgrp-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciNodeBlockMGIdNotEqual(&maintenance_group_node_default, &maintenance_group_node_updated),
				),
			},
		},
	})
}

func TestAccAciNodeBlockMG_Update(t *testing.T) {
	var maintenance_group_node_default models.NodeBlockMG
	var maintenance_group_node_updated models.NodeBlockMG
	resourceName := "aci_maintenance_group_node.test"
	rName := makeTestVariable(acctest.RandString(5))

	maintMaintGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockMGDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccNodeBlockMGConfig(maintMaintGrpName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockMGExists(resourceName, &maintenance_group_node_default),
				),
			},
			{
				Config: CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "to_", "16000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockMGExists(resourceName, &maintenance_group_node_updated),
					resource.TestCheckResourceAttr(resourceName, "to_", "16000"),
					testAccCheckAciNodeBlockMGIdEqual(&maintenance_group_node_default, &maintenance_group_node_updated),
				),
			},
			{
				Config: CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "from_", "16000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockMGExists(resourceName, &maintenance_group_node_updated),
					resource.TestCheckResourceAttr(resourceName, "from_", "16000"),
					testAccCheckAciNodeBlockMGIdEqual(&maintenance_group_node_default, &maintenance_group_node_updated),
				),
			},
			{
				Config: CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "from_", "7999"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockMGExists(resourceName, &maintenance_group_node_updated),
					resource.TestCheckResourceAttr(resourceName, "from_", "7999"),
					testAccCheckAciNodeBlockMGIdEqual(&maintenance_group_node_default, &maintenance_group_node_updated),
				),
			},
			{
				Config: CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "to_", "7999"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciNodeBlockMGExists(resourceName, &maintenance_group_node_updated),
					resource.TestCheckResourceAttr(resourceName, "to_", "7999"),
					testAccCheckAciNodeBlockMGIdEqual(&maintenance_group_node_default, &maintenance_group_node_updated),
				),
			},

			{
				Config: CreateAccNodeBlockMGConfig(maintMaintGrpName, rName),
			},
		},
	})
}

func TestAccAciNodeBlockMG_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	maintMaintGrpName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockMGDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccNodeBlockMGConfig(maintMaintGrpName, rName),
			},
			{
				Config:      CreateAccNodeBlockMGWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "from_", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "from_", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "from_", "16001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "to_", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "to_", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, "to_", "16001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccNodeBlockMGConfig(maintMaintGrpName, rName),
			},
		},
	})
}

func TestAccAciNodeBlockMG_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	maintMaintGrpName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockMGDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccNodeBlockMGConfigMultiple(maintMaintGrpName, rName),
			},
		},
	})
}

func testAccCheckAciNodeBlockMGExists(name string, maintenance_group_node *models.NodeBlockMG) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Node Block MG %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Node Block MG dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		maintenance_group_nodeFound := models.NodeBlockFromContainerMG(cont)
		if maintenance_group_nodeFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Node Block MG %s not found", rs.Primary.ID)
		}
		*maintenance_group_node = *maintenance_group_nodeFound
		return nil
	}
}

func testAccCheckAciNodeBlockMGDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing maintenance_group_node destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_maintenance_group_node" {
			cont, err := client.Get(rs.Primary.ID)
			maintenance_group_node := models.NodeBlockFromContainerMG(cont)
			if err == nil {
				return fmt.Errorf("Node Block MG %s Still exists", maintenance_group_node.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciNodeBlockMGIdEqual(m1, m2 *models.NodeBlockMG) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("maintenance_group_node DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciNodeBlockMGIdNotEqual(m1, m2 *models.NodeBlockMG) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("maintenance_group_node DNs are equal")
		}
		return nil
	}
}

func CreateNodeBlockMGWithoutRequired(maintMaintGrpName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing maintenance_group_node creation without ", attrName)
	rBlock := `
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "pod_maintenance_group_dn":
		rBlock += `
	resource "aci_maintenance_group_node" "test" {
	#	pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, maintMaintGrpName, rName)
}

func CreateAccNodeBlockMGConfigWithRequiredParams(prName, rName string) string {
	fmt.Printf("=== STEP  testing maintenance_group_node creation with parent resource name %s and resource name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "%s"
	}
	`, prName, rName)
	return resource
}
func CreateAccNodeBlockMGConfigUpdatedName(maintMaintGrpName, rName string) string {
	fmt.Println("=== STEP  testing maintenance_group_node creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "%s"
	}
	`, maintMaintGrpName, rName)
	return resource
}

func CreateAccNodeBlockMGConfig(maintMaintGrpName, rName string) string {
	fmt.Println("=== STEP  testing maintenance_group_node creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "%s"
	}
	`, maintMaintGrpName, rName)
	return resource
}

func CreateAccNodeBlockMGConfigMultiple(maintMaintGrpName, rName string) string {
	fmt.Println("=== STEP  testing multiple maintenance_group_node creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "%s_${count.index}"
		from_ = (count.index+1)*10 
		to_ = (count.index+1)*10+5
		count = 5
	}
	`, maintMaintGrpName, rName)
	return resource
}

func CreateAccNodeBlockMGWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing maintenance_group_node creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccNodeBlockMGConfigWithOptionalValues(maintMaintGrpName, rName string) string {
	fmt.Println("=== STEP  Basic: testing maintenance_group_node creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = "${aci_pod_maintenance_group.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_maintenance_group_node"
		from_ = "2"
		to_ = "2"
		
	}
	`, maintMaintGrpName, rName)

	return resource
}

func CreateAccNodeBlockMGRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing maintenance_group_node updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_maintenance_group_node" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_maintenance_group_node"
		from_ = "2"
		to_ = "2"
		
	}
	`)

	return resource
}

func CreateAccNodeBlockMGUpdatedAttr(maintMaintGrpName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing maintenance_group_node attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "%s"
		%s = "%s"
	}
	`, maintMaintGrpName, rName, attribute, value)
	return resource
}

func CreateAccNodeBlockMGUpdatedAttrList(maintMaintGrpName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing maintenance_group_node attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_pod_maintenance_group" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_maintenance_group_node" "test" {
		pod_maintenance_group_dn  = aci_pod_maintenance_group.test.id
		name  = "%s"
		%s = %s
	}
	`, maintMaintGrpName, rName, attribute, value)
	return resource
}
