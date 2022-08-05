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

func TestAccAciFabricIfPolicy_Basic(t *testing.T) {
	var fabric_if_policy_default models.LinkLevelPolicy
	var fabric_if_policy_updated models.LinkLevelPolicy
	resourceName := "aci_fabric_if_pol.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricIfPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFabricIfPolicyWithoutName(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricIfPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "auto_neg", "on"),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "inherit"),
					resource.TestCheckResourceAttr(resourceName, "link_debounce", "100"),
					resource.TestCheckResourceAttr(resourceName, "speed", "inherit"),
				),
			},
			{
				Config: CreateAccFabricIfPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_fabric_if_pol"),
					resource.TestCheckResourceAttr(resourceName, "auto_neg", "off"),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "cl91-rs-fec"),
					resource.TestCheckResourceAttr(resourceName, "link_debounce", "0"),
					resource.TestCheckResourceAttr(resourceName, "speed", "unknown"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFabricIfPolicyConfigWithUpdatedRequiredParams(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)* failed validation`),
			},
			{
				Config: CreateAccFabricIfPolicyConfigWithUpdatedRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciFabricIfPolicyIdNotEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config:      CreateAccFabricIfPolicyConfigUpdateWithoutName("description", randomValue),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricIfPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciFabricIfPolicy_Update(t *testing.T) {
	var fabric_if_policy_default models.LinkLevelPolicy
	var fabric_if_policy_updated models.LinkLevelPolicy
	resourceName := "aci_fabric_if_pol.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricIfPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricIfPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_default),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "fec_mode", "cl74-fc-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "cl74-fc-fec"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "fec_mode", "cons16-rs-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "cons16-rs-fec"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "fec_mode", "disable-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "disable-fec"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "fec_mode", "ieee-rs-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "ieee-rs-fec"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "fec_mode", "kp-fec"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "fec_mode", "kp-fec"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "speed", "100M"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "100M"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "speed", "1G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "1G"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "speed", "10G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "10G"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "speed", "25G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "25G"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "speed", "40G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "40G"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "speed", "50G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "50G"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "speed", "100G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "100G"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "speed", "200G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "200G"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "speed", "400G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "400G"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "link_debounce", "5000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "link_debounce", "5000"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyUpdatedAttr(rName, "link_debounce", "2999"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricIfPolicyExists(resourceName, &fabric_if_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "link_debounce", "2999"),
					testAccCheckAciFabricIfPolicyIdEqual(&fabric_if_policy_default, &fabric_if_policy_updated),
				),
			},
			{
				Config: CreateAccFabricIfPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciFabricIfPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricIfPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricIfPolicyConfig(rName),
			},
			{
				Config:      CreateAccFabricIfPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricIfPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricIfPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricIfPolicyUpdatedAttr(rName, "auto_neg", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccFabricIfPolicyUpdatedAttr(rName, "fec_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccFabricIfPolicyUpdatedAttr(rName, "link_debounce", "10000"),
				ExpectError: regexp.MustCompile(`Property linkDebounce of (.)* is out of range`),
			},
			{
				Config:      CreateAccFabricIfPolicyUpdatedAttr(rName, "link_debounce", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccFabricIfPolicyUpdatedAttr(rName, "speed", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccFabricIfPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)*is not expected here.`),
			},
			{
				Config: CreateAccFabricIfPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciFabricIfPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricIfPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricIfPoliciesConfig(rName),
			},
		},
	})
}

func testAccCheckAciFabricIfPolicyExists(name string, fabric_if_pol *models.LinkLevelPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric If Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric If Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabric_if_polFound := models.LinkLevelPolicyFromContainer(cont)
		if fabric_if_polFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric If Policy %s not found", rs.Primary.ID)
		}
		*fabric_if_pol = *fabric_if_polFound
		return nil
	}
}

func testAccCheckAciFabricIfPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing fabric_if_pol destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_fabric_if_pol" {
			cont, err := client.Get(rs.Primary.ID)
			fabric_if_pol := models.LinkLevelPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fabric If Policy %s Still exists", fabric_if_pol.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFabricIfPolicyIdEqual(m1, m2 *models.LinkLevelPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("fabric_if_pol DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFabricIfPolicyIdNotEqual(m1, m2 *models.LinkLevelPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("fabric_if_pol DNs are equal")
		}
		return nil
	}
}

func CreateFabricIfPolicyWithoutName(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_if_pol creation without", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_fabric_if_pol" "test" {
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFabricIfPolicyConfigWithUpdatedRequiredParams(rName string) string {
	fmt.Printf("=== STEP  testing fabric_if_pol creation with name = %s\n", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFabricIfPoliciesConfig(rName string) string {
	fmt.Println("=== STEP  testing fabric_if_pol multiple creation")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
		name  = "%s"
	}

	resource "aci_fabric_if_pol" "test1" {
		name  = "%s"
	}

	resource "aci_fabric_if_pol" "test2" {
		name  = "%s"
	}

	resource "aci_fabric_if_pol" "test3" {
		name  = "%s"
	}
	`, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccFabricIfPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing fabric_if_pol creation with required arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccFabricIfPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing fabric_if_pol creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fabric_if_pol"
		auto_neg = "off"
		fec_mode = "cl91-rs-fec"
		link_debounce = "0"
		speed = "unknown"
	}
	`, rName)

	return resource
}

func CreateAccFabricIfPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing fabric_if_pol attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccFabricIfPolicyConfigUpdateWithoutName(attribute, value string) string {
	fmt.Printf("=== STEP  testing fabric_if_pol attribute: %s=%s without name\n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fabric_if_pol" "test" {
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
