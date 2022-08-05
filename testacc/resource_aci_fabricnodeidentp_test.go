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

func TestAccAciFabricNodeMember_Basic(t *testing.T) {
	var fabric_node_member_default models.FabricNodeMember
	var fabric_node_member_updated models.FabricNodeMember
	resourceName := "aci_fabric_node_member.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	serial1 := makeTestVariable(acctest.RandString(5))
	serial2 := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeMemberDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFabricNodeMemberWithoutRequired(rName, serial1, "serial"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateFabricNodeMemberWithoutRequired(rName, serial1, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricNodeMemberConfig(rName, serial1, fabricNodeMemNodeId1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_default),
					resource.TestCheckResourceAttr(resourceName, "serial", serial1),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ext_pool_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "fabric_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "node_id", fabricNodeMemNodeId1),
					resource.TestCheckResourceAttr(resourceName, "node_type", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "pod_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "role", "unspecified"),
				),
			},
			{
				Config: CreateAccFabricNodeMemberConfigWithOptionalValues(rNameUpdated, serial1, fabricNodeMemNodeId1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "serial", serial1),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_fabric_node_member"),
					resource.TestCheckResourceAttr(resourceName, "node_id", fabricNodeMemNodeId1),
					resource.TestCheckResourceAttr(resourceName, "node_type", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "pod_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "role", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "ext_pool_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "fabric_id", "1"),
					testAccCheckAciFabricNodeMemberIdEqual(&fabric_node_member_default, &fabric_node_member_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccFabricNodeMemberRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccFabricNodeMemberConfigWithRequiredParams(acctest.RandString(65), serial2, fabricNodeMemNodeId4),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config: CreateAccFabricNodeMemberConfigWithRequiredParams(rNameUpdated, serial2, fabricNodeMemNodeId1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "serial", serial2),
					testAccCheckAciFabricNodeMemberIdNotEqual(&fabric_node_member_default, &fabric_node_member_updated),
				),
			},
		},
	})
}

func TestAccAciFabricNodeMember_Update(t *testing.T) {
	var fabric_node_member_default models.FabricNodeMember
	var fabric_node_member_updated models.FabricNodeMember
	resourceName := "aci_fabric_node_member.test"
	serial1 := makeTestVariable(acctest.RandString(5))
	serial2 := makeTestVariable(acctest.RandString(5))
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricNodeMemberConfig(rName, serial1, fabricNodeMemNodeId2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_default),
				),
			},
			{
				Config: CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId2, "ext_pool_id", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "ext_pool_id", "100"),
				),
			},
			{
				Config: CreateAccFabricNodeMemberUpdatedAttr(rName, serial1, fabricNodeMemNodeId2, "fabric_id", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "fabric_id", "100"),
				),
			},
			{
				Config: CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId2, "node_type", "remote-leaf-wan"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "node_type", "remote-leaf-wan"),
				),
			},
			{
				Config: CreateAccFabricNodeMemberUpdatedAttr(rName, serial1, fabricNodeMemNodeId2, "pod_id", "126"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "pod_id", "126"),
				),
			},
			{
				Config: CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId2, "pod_id", "254"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "pod_id", "254"),
				),
			},
			{
				Config: CreateAccFabricNodeMemberUpdatedAttr(rName, serial1, fabricNodeMemNodeId2, "role", "leaf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "role", "leaf"),
				),
			},
			{
				Config: CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId2, "role", "spine"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists(resourceName, &fabric_node_member_updated),
					resource.TestCheckResourceAttr(resourceName, "role", "spine"),
				),
			},
			{
				Config: CreateAccFabricNodeMemberConfig(rName, serial1, fabricNodeMemNodeId2),
			},
		},
	})
}

func TestAccAciFabricNodeMember_Negative(t *testing.T) {
	serial1 := makeTestVariable(acctest.RandString(5))
	serial2 := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricNodeMemberConfig(rName, serial1, fabricNodeMemNodeId3),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId4, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId4, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId4, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId4, "ext_pool_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId4, "fabric_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId4, "node_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId4, "pod_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId4, "pod_id", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId4, "pod_id", "255"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId4, "role", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttrNode(rName, serial2, "100"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttrNode(rName, serial2, "4001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttrNode(rName, serial2, randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccFabricNodeMemberUpdatedAttr(rName, serial2, fabricNodeMemNodeId4, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFabricNodeMemberConfig(rName, serial2, fabricNodeMemNodeId4),
			},
		},
	})
}

func testAccCheckAciFabricNodeMemberExists(name string, fabric_node_member *models.FabricNodeMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric Node Member %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric Node Member dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabric_node_memberFound := models.FabricNodeMemberFromContainer(cont)
		if fabric_node_memberFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric Node Member %s not found", rs.Primary.ID)
		}
		*fabric_node_member = *fabric_node_memberFound
		return nil
	}
}

func testAccCheckAciFabricNodeMemberDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing fabric_node_member destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_fabric_node_member" {
			cont, err := client.Get(rs.Primary.ID)
			fabric_node_member := models.FabricNodeMemberFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fabric Node Member %s Still exists", fabric_node_member.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFabricNodeMemberIdEqual(m1, m2 *models.FabricNodeMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("fabric_node_member DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFabricNodeMemberIdNotEqual(m1, m2 *models.FabricNodeMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("fabric_node_member DNs are equal")
		}
		return nil
	}
}

func CreateFabricNodeMemberWithoutRequired(name, serial, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_node_member creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "serial":
		rBlock += `
	resource "aci_fabric_node_member" "test" {
		name = "%s"
	#	serial  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_fabric_node_member" "test" {
	#	name = "%s"
		serial  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, name, serial)
}

func CreateAccFabricNodeMemberConfigWithRequiredParams(name, serial, node string) string {
	fmt.Printf("=== STEP  testing fabric_node_member creation with serial %s and name %s\n", serial, name)
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_member" "test" {
		serial  = "%s"
		node_id = "%s"
		name = "%s"
	}
	`, serial, node, name)
	return resource
}

func CreateAccFabricNodeMemberConfig(name, serial, node string) string {
	fmt.Println("=== STEP  testing fabric_node_member creation with required arguments and node_id")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_member" "test" {
		serial  = "%s"
		name = "%s"
		node_id = "%s"
	}
	`, serial, name, node)
	return resource
}

func CreateAccFabricNodeMemberConfigWithOptionalValues(name, serial, node string) string {
	fmt.Println("=== STEP  Basic: testing fabric_node_member creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_member" "test" {
		name = "%s"
		serial  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_node_member"
		node_id = "%s"
		node_type = "unspecified"
		pod_id = "1"
		role = "unspecified"
		ext_pool_id = "0"
		fabric_id = "1"
	}
	`, name, serial, node)

	return resource
}

func CreateAccFabricNodeMemberRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing fabric_node_member updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_fabric_node_member" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_node_member"
		node_id = "102"
		node_type = "remote-leaf-wan"
		pod_id = "2"
		role = "leaf"
		
	}
	`)

	return resource
}

func CreateAccFabricNodeMemberUpdatedAttrNode(name, serial, node string) string {
	fmt.Printf("=== STEP  testing fabric_node_member attribute: node = %s \n", node)
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_member" "test" {
		name = "%s"
		serial  = "%s"
		node_id = "%s"
	}
	`, name, serial, node)
	return resource
}

func CreateAccFabricNodeMemberUpdatedAttr(name, serial, node, attribute, value string) string {
	fmt.Printf("=== STEP  testing fabric_node_member attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_member" "test" {
		name = "%s"
		serial  = "%s"
		node_id = "%s"
		%s = "%s"
	}
	`, name, serial, node, attribute, value)
	return resource
}
