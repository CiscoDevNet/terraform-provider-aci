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

func TestAccAciFabricNodeControl_Basic(t *testing.T) {
	var fabric_node_control_default models.FabricNodeControl
	var fabric_node_control_updated models.FabricNodeControl
	resourceName := "aci_fabric_node_control.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeControlDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFabricNodeControlWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricNodeControlConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeControlExists(resourceName, &fabric_node_control_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "control", "None"),
					resource.TestCheckResourceAttr(resourceName, "feature_sel", "telemetry"),
				),
			},
			{
				Config: CreateAccFabricNodeControlConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeControlExists(resourceName, &fabric_node_control_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_fabric_node_control"),
					resource.TestCheckResourceAttr(resourceName, "control", "Dom"),
					resource.TestCheckResourceAttr(resourceName, "feature_sel", "analytics"),
					testAccCheckAciFabricNodeControlIdEqual(&fabric_node_control_default, &fabric_node_control_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFabricNodeControlConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccFabricNodeControlRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricNodeControlConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeControlExists(resourceName, &fabric_node_control_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciFabricNodeControlIdNotEqual(&fabric_node_control_default, &fabric_node_control_updated),
				),
			},
		},
	})
}

func TestAccAciFabricNodeControl_Update(t *testing.T) {
	var fabric_node_control_default models.FabricNodeControl
	var fabric_node_control_updated models.FabricNodeControl
	resourceName := "aci_fabric_node_control.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeControlDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricNodeControlConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeControlExists(resourceName, &fabric_node_control_default),
				),
			},
			{
				Config: CreateAccFabricNodeControlUpdatedAttr(rName, "feature_sel", "netflow"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeControlExists(resourceName, &fabric_node_control_updated),
					resource.TestCheckResourceAttr(resourceName, "feature_sel", "netflow"),
					testAccCheckAciFabricNodeControlIdEqual(&fabric_node_control_default, &fabric_node_control_updated),
				),
			},
			{
				Config: CreateAccFabricNodeControlConfig(rName),
			},
		},
	})
}

func TestAccAciFabricNodeControl_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeControlDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricNodeControlConfig(rName),
			},

			{
				Config:      CreateAccFabricNodeControlUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricNodeControlUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricNodeControlUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricNodeControlUpdatedAttr(rName, "control", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccFabricNodeControlUpdatedAttr(rName, "feature_sel", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccFabricNodeControlUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFabricNodeControlConfig(rName),
			},
		},
	})
}

func TestAccAciFabricNodeControl_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeControlDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricNodeControlConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciFabricNodeControlExists(name string, fabric_node_control *models.FabricNodeControl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric Node Control %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric Node Control dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabric_node_controlFound := models.FabricNodeControlFromContainer(cont)
		if fabric_node_controlFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric Node Control %s not found", rs.Primary.ID)
		}
		*fabric_node_control = *fabric_node_controlFound
		return nil
	}
}

func testAccCheckAciFabricNodeControlDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing fabric_node_control destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_fabric_node_control" {
			cont, err := client.Get(rs.Primary.ID)
			fabric_node_control := models.FabricNodeControlFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fabric Node Control %s Still exists", fabric_node_control.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFabricNodeControlIdEqual(m1, m2 *models.FabricNodeControl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("fabric_node_control DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFabricNodeControlIdNotEqual(m1, m2 *models.FabricNodeControl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("fabric_node_control DNs are equal")
		}
		return nil
	}
}

func CreateFabricNodeControlWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_node_control creation without ", attrName)
	rBlock := `
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_fabric_node_control" "test" {
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFabricNodeControlConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing fabric_node_control creation with updated Name")
	resource := fmt.Sprintf(`
	resource "aci_fabric_node_control" "test" {	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccFabricNodeControlConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing fabric_node_control creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	resource "aci_fabric_node_control" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFabricNodeControlConfig(rName string) string {
	fmt.Println("=== STEP  testing fabric_node_control creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_fabric_node_control" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFabricNodeControlConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple fabric_node_control creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_fabric_node_control" "test" {
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccFabricNodeControlConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_node_control creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_node_control" "test" {
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_node_control"
		control = "Dom"
		feature_sel = "analytics"
	}
	`, rName)

	return resource
}

func CreateAccFabricNodeControlRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing fabric_node_control updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_fabric_node_control" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_node_control"
		control = "Dom"
		feature_sel = "analytics"
	}
	`)

	return resource
}

func CreateAccFabricNodeControlUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing fabric_node_control attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_fabric_node_control" "test" {
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
