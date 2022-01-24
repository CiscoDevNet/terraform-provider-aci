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

func TestAccAciInterfaceFcPolicyDataSource_Basic(t *testing.T) {
	var interface_fc_policy_default models.InterfaceFCPolicy
	var interface_fc_policy_updated models.InterfaceFCPolicy
	resourceName := "aci_interface_fc_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciInterfaceFcPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateInterfaceFcPolicyWithoutName(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccInterfaceFcPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "automaxspeed", "32G"),
					resource.TestCheckResourceAttr(resourceName, "fill_pattern", "IDLE"),
					resource.TestCheckResourceAttr(resourceName, "port_mode", "f"),
					resource.TestCheckResourceAttr(resourceName, "rx_bb_credit", "64"),
					resource.TestCheckResourceAttr(resourceName, "speed", "auto"),
					resource.TestCheckResourceAttr(resourceName, "trunk_mode", "trunk-off"),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_interface_fc_policy"),
					resource.TestCheckResourceAttr(resourceName, "automaxspeed", "16G"),
					resource.TestCheckResourceAttr(resourceName, "fill_pattern", "ARBFF"),
					resource.TestCheckResourceAttr(resourceName, "port_mode", "np"),
					resource.TestCheckResourceAttr(resourceName, "rx_bb_credit", "17"),
					resource.TestCheckResourceAttr(resourceName, "speed", "unknown"),
					resource.TestCheckResourceAttr(resourceName, "trunk_mode", "auto"),
					testAccCheckAciInterfaceFcPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccInterfaceFcPolicyConfig(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)* failed validation`),
			},
			{
				Config: CreateAccInterfaceFcPolicyConfigWithUpdatedRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciInterfaceFcPolicyIdNotEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config:      CreateAccInterfaceFcPolicyConfigUpdateWithoutName("description", randomValue),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccInterfaceFcPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciInterfaceFcPolicy_Update(t *testing.T) {
	var interface_fc_policy_default models.InterfaceFCPolicy
	var interface_fc_policy_updated models.InterfaceFCPolicy
	resourceName := "aci_interface_fc_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciInterfaceFcPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccInterfaceFcPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_default),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyUpdatedAttr(rName, "automaxspeed", "2G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "automaxspeed", "2G"),
					testAccCheckAciInterfaceFcPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyUpdatedAttr(rName, "automaxspeed", "4G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "automaxspeed", "4G"),
					testAccCheckAciInterfaceFcPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyUpdatedAttr(rName, "automaxspeed", "8G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "automaxspeed", "8G"),
					testAccCheckAciInterfaceFcPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyUpdatedAttr(rName, "speed", "32G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "32G"),
					testAccCheckAciInterfaceFcPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyUpdatedAttr(rName, "speed", "4G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "4G"),
					testAccCheckAciInterfaceFcPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyUpdatedAttr(rName, "speed", "8G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "8G"),
					testAccCheckAciInterfaceFcPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyUpdatedAttr(rName, "speed", "16G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "speed", "16G"),
					testAccCheckAciInterfaceFcPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyUpdatedAttr(rName, "trunk_mode", "trunk-on"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "trunk_mode", "trunk-on"),
					testAccCheckAciInterfaceFcPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyUpdatedAttr(rName, "trunk_mode", "un-init"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "trunk_mode", "un-init"),
					testAccCheckAciInterfaceFcPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyUpdatedAttr(rName, "rx_bb_credit", "16"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "rx_bb_credit", "16"),
					testAccCheckAciInterfaceFcPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyUpdatedAttr(rName, "rx_bb_credit", "64"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFcPolicyExists(resourceName, &interface_fc_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "rx_bb_credit", "64"),
					testAccCheckAciInterfaceFcPolicyIdEqual(&interface_fc_policy_default, &interface_fc_policy_updated),
				),
			},
			{
				Config: CreateAccInterfaceFcPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciInterfaceFcPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciInterfaceFcPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccInterfaceFcPolicyConfig(rName),
			},
			{
				Config:      CreateAccInterfaceFcPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccInterfaceFcPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccInterfaceFcPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccInterfaceFcPolicyUpdatedAttr(rName, "automaxspeed", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccInterfaceFcPolicyUpdatedAttr(rName, "fill_pattern", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccInterfaceFcPolicyUpdatedAttr(rName, "port_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccInterfaceFcPolicyUpdatedAttr(rName, "rx_bb_credit", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccInterfaceFcPolicyUpdatedAttr(rName, "rx_bb_credit", "65"),
				ExpectError: regexp.MustCompile(`Property rxBBCredit of (.)+ is out of range`),
			},
			{
				Config:      CreateAccInterfaceFcPolicyUpdatedAttr(rName, "speed", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccInterfaceFcPolicyUpdatedAttr(rName, "trunk_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},

			{
				Config:      CreateAccInterfaceFcPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)*is not expected here.`),
			},
			{
				Config: CreateAccInterfaceFcPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciInterfaceFcPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciInterfaceFcPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccInterfaceFcPoliciesConfig(rName),
			},
		},
	})
}
func testAccCheckAciInterfaceFcPolicyExists(name string, interface_fc_policy *models.InterfaceFCPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Interface Fc Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Interface Fc Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		interface_fc_policyFound := models.InterfaceFCPolicyFromContainer(cont)
		if interface_fc_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Interface Fc Policy %s not found", rs.Primary.ID)
		}
		*interface_fc_policy = *interface_fc_policyFound
		return nil
	}
}

func testAccCheckAciInterfaceFcPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing interface_fc_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_interface_fc_policy" {
			cont, err := client.Get(rs.Primary.ID)
			interface_fc_policy := models.InterfaceFCPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Interface Fc Policy %s Still exists", interface_fc_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciInterfaceFcPolicyIdEqual(m1, m2 *models.InterfaceFCPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("interface_fc_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciInterfaceFcPolicyIdNotEqual(m1, m2 *models.InterfaceFCPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("interface_fc_policy DNs are equal")
		}
		return nil
	}
}

func CreateInterfaceFcPolicyWithoutName(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing interface_fc_policy creation without Name ", attrName)
	rBlock := `
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_interface_fc_policy" "test" {	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccInterfaceFcPolicyConfigWithUpdatedRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing interface_fc_policy updation of required arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccInterfaceFcPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing interface_fc_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccInterfaceFcPoliciesConfig(rName string) string {
	fmt.Println("=== STEP  testing Multiple interface_fc_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_interface_fc_policy" "test" {
		name  = "%s"
	}

	resource "aci_interface_fc_policy" "test1" {
		name  = "%s"
	}

	resource "aci_interface_fc_policy" "test2" {
		name  = "%s"
	}

	resource "aci_interface_fc_policy" "test3" {
		name  = "%s"
	}
	`, rName, rName+"1", rName+"2", rName+"3")
	return resource
}
func CreateAccInterfaceFcPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing interface_fc_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_interface_fc_policy"
		automaxspeed = "16G"
		fill_pattern = "ARBFF"
		port_mode = "np"
		rx_bb_credit = "17"
		speed = "unknown"
		trunk_mode = "auto"
	}
	`, rName)

	return resource
}

func CreateAccInterfaceFcPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing interface_fc_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccInterfaceFcPolicyConfigUpdateWithoutName(attribute, value string) string {
	fmt.Printf("=== STEP  testing interface_fc_policy attribute: %s=%s without name \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_interface_fc_policy" "test" {
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
