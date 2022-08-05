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

func TestAccAciLLDPInterfacePolicy_Basic(t *testing.T) {
	var lldp_interface_policy_default models.LLDPInterfacePolicy
	var lldp_interface_policy_updated models.LLDPInterfacePolicy
	resourceName := "aci_lldp_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLLDPInterfacePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLLDPInterfacePolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLLDPInterfacePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLLDPInterfacePolicyExists(resourceName, &lldp_interface_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "admin_rx_st", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "admin_tx_st", "enabled"),
				),
			},
			{
				Config: CreateAccLLDPInterfacePolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLLDPInterfacePolicyExists(resourceName, &lldp_interface_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_lldp_interface_policy"),
					resource.TestCheckResourceAttr(resourceName, "admin_rx_st", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "admin_tx_st", "disabled"),

					testAccCheckAciLLDPInterfacePolicyIdEqual(&lldp_interface_policy_default, &lldp_interface_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccLLDPInterfacePolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccLLDPInterfacePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLLDPInterfacePolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLLDPInterfacePolicyExists(resourceName, &lldp_interface_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLLDPInterfacePolicyIdNotEqual(&lldp_interface_policy_default, &lldp_interface_policy_updated),
				),
			},
		},
	})
}

func TestAccAciLLDPInterfacePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLLDPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLLDPInterfacePolicyConfigs(rName),
			},
		},
	})
}

func TestAccAciLLDPInterfacePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLLDPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLLDPInterfacePolicyConfig(rName),
			},

			{
				Config:      CreateAccLLDPInterfacePolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLLDPInterfacePolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLLDPInterfacePolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLLDPInterfacePolicyUpdatedAttr(rName, "admin_rx_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccLLDPInterfacePolicyUpdatedAttr(rName, "admin_tx_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLLDPInterfacePolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLLDPInterfacePolicyConfig(rName),
			},
		},
	})
}

func testAccCheckAciLLDPInterfacePolicyExists(name string, lldp_interface_policy *models.LLDPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LLDP Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LLDP Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		lldp_interface_policyFound := models.LLDPInterfacePolicyFromContainer(cont)
		if lldp_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LLDP Interface Policy %s not found", rs.Primary.ID)
		}
		*lldp_interface_policy = *lldp_interface_policyFound
		return nil
	}
}

func testAccCheckAciLLDPInterfacePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing lldp_interface_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_lldp_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			lldp_interface_policy := models.LLDPInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LLDP Interface Policy %s Still exists", lldp_interface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLLDPInterfacePolicyIdEqual(m1, m2 *models.LLDPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("lldp_interface_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLLDPInterfacePolicyIdNotEqual(m1, m2 *models.LLDPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("lldp_interface_policy DNs are equal")
		}
		return nil
	}
}

func CreateLLDPInterfacePolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing lldp_interface_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_lldp_interface_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLLDPInterfacePolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing lldp_interface_policy creation with updated name")
	resource := fmt.Sprintf(`
	
	resource "aci_lldp_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLLDPInterfacePolicyConfigs(rName string) string {
	fmt.Println("=== STEP  testing multiple lldp_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_lldp_interface_policy" "test" {
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccLLDPInterfacePolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing lldp_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_lldp_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLLDPInterfacePolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing lldp_interface_policy creation with longer name")
	resource := fmt.Sprintf(`
	
	resource "aci_lldp_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLLDPInterfacePolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing lldp_interface_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_lldp_interface_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_lldp_interface_policy"
		admin_rx_st = "disabled"
		admin_tx_st = "disabled"
	}
	`, rName)

	return resource
}

func CreateAccLLDPInterfacePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing lldp_interface_policy update without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_lldp_interface_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_lldp_interface_policy"
		admin_rx_st = "disabled"
		admin_tx_st = "disabled"
	}
	`)

	return resource
}

func CreateAccLLDPInterfacePolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing lldp_interface_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_lldp_interface_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
