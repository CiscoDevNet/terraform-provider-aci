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

func TestAccAciLACPPolicy_Basic(t *testing.T) {
	var lacp_policy_default models.LACPPolicy
	var lacp_policy_updated models.LACPPolicy
	resourceName := "aci_lacp_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLACPPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLACPPolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLACPPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "fast-sel-hot-stdby"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "graceful-conv"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.2", "susp-individual"),
					resource.TestCheckResourceAttr(resourceName, "max_links", "16"),
					resource.TestCheckResourceAttr(resourceName, "min_links", "1"),
					resource.TestCheckResourceAttr(resourceName, "mode", "off"),
				),
			},
			{
				Config: CreateAccLACPPolicyConfigWithOptionalValues(rName), // configuration to update optional filelds
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_lacp_policy"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "load-defer"),
					resource.TestCheckResourceAttr(resourceName, "max_links", "2"),
					resource.TestCheckResourceAttr(resourceName, "min_links", "2"),
					resource.TestCheckResourceAttr(resourceName, "mode", "active"),

					testAccCheckAciLACPPolicyIdEqual(&lacp_policy_default, &lacp_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccLACPPolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccLACPPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccLACPPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLACPPolicyIdNotEqual(&lacp_policy_default, &lacp_policy_updated),
				),
			},
		},
	})
}

func TestAccAciLACPPolicy_Update(t *testing.T) {
	var lacp_policy_default models.LACPPolicy
	var lacp_policy_updated models.LACPPolicy
	resourceName := "aci_lacp_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLACPPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLACPPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_default),
				),
			},

			{

				Config: CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{"fast-sel-hot-stdby"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "fast-sel-hot-stdby"),
				),
			},
			{

				Config: CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{"fast-sel-hot-stdby", "graceful-conv"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "fast-sel-hot-stdby"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "graceful-conv"),
				),
			},
			{

				Config: CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{"graceful-conv"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "graceful-conv"),
				),
			},
			{

				Config: CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{"fast-sel-hot-stdby", "graceful-conv", "load-defer"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "fast-sel-hot-stdby"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "graceful-conv"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.2", "load-defer"),
				),
			},
			{

				Config: CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{"graceful-conv", "load-defer"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "graceful-conv"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "load-defer"),
				),
			},
			{

				Config: CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{"load-defer"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "load-defer"),
				),
			},
			{

				Config: CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{"fast-sel-hot-stdby", "graceful-conv", "load-defer", "susp-individual"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "fast-sel-hot-stdby"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "graceful-conv"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.2", "load-defer"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.3", "susp-individual"),
				),
			},
			{

				Config: CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{"graceful-conv", "load-defer", "susp-individual"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "graceful-conv"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "load-defer"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.2", "susp-individual"),
				),
			},
			{

				Config: CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{"load-defer", "graceful-conv", "susp-individual"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "load-defer"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "graceful-conv"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.2", "susp-individual"),
				),
			},
			{
				Config: CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{"load-defer", "graceful-conv", "susp-individual", "symmetric-hash", "fast-sel-hot-stdby"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "5"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "load-defer"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "graceful-conv"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.2", "susp-individual"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.3", "symmetric-hash"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.4", "fast-sel-hot-stdby"),
				),
			},
			{
				Config: CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{"susp-individual", "graceful-conv", "fast-sel-hot-stdby"})),
			},
			{
				Config: CreateAccLACPPolicyUpdatedAttr(rName, "mode", "explicit-failover"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "explicit-failover"),
					testAccCheckAciLACPPolicyIdEqual(&lacp_policy_default, &lacp_policy_updated),
				),
			},
			{
				Config: CreateAccLACPPolicyUpdatedAttr(rName, "mode", "mac-pin"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "mac-pin"),
					testAccCheckAciLACPPolicyIdEqual(&lacp_policy_default, &lacp_policy_updated),
				),
			},
			{
				Config: CreateAccLACPPolicyUpdatedAttr(rName, "mode", "mac-pin-nicload"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "mac-pin-nicload"),
					testAccCheckAciLACPPolicyIdEqual(&lacp_policy_default, &lacp_policy_updated),
				),
			},
			{
				Config: CreateAccLACPPolicyUpdatedAttr(rName, "mode", "passive"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists(resourceName, &lacp_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "passive"),
					testAccCheckAciLACPPolicyIdEqual(&lacp_policy_default, &lacp_policy_updated),
				),
			},
			{
				Config: CreateAccLACPPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciLACPPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLACPPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLACPPoliciesConfig(rName),
			},
		},
	})
}
func TestAccAciLACPPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLACPPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLACPPolicyConfig(rName),
			},

			{
				Config:      CreateAccLACPPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value`),
			},
			{
				Config:      CreateAccLACPPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value`),
			},
			{
				Config:      CreateAccLACPPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value`),
			},
			{
				Config:      CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccLACPPolicyUpdatedListAttr(rName, "ctrl", StringListtoString([]string{"fast-sel-hot-stdby", "fast-sel-hot-stdby"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccLACPPolicyUpdatedAttr(rName, "max_links", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLACPPolicyUpdatedAttr(rName, "min_links", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLACPPolicyUpdatedAttr(rName, "mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLACPPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLACPPolicyConfig(rName),
			},
		},
	})
}

func testAccCheckAciLACPPolicyExists(name string, lacp_policy *models.LACPPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LACP Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LACP Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		lacp_policyFound := models.LACPPolicyFromContainer(cont)
		if lacp_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LACP Policy %s not found", rs.Primary.ID)
		}
		*lacp_policy = *lacp_policyFound
		return nil
	}
}

func testAccCheckAciLACPPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing lacp_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_lacp_policy" {
			cont, err := client.Get(rs.Primary.ID)
			lacp_policy := models.LACPPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LACP Policy %s Still exists", lacp_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLACPPolicyIdEqual(m1, m2 *models.LACPPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("lacp_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLACPPolicyIdNotEqual(m1, m2 *models.LACPPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("lacp_policy DNs are equal")
		}
		return nil
	}
}

func CreateLACPPolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing lacp_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_lacp_policy" "test" {
	
	#	name  = "%s"
		description = "created while acceptance testing"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLACPPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing lacp_policy creation with updated required arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_lacp_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLACPPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing lacp_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_lacp_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLACPPoliciesConfig(rName string) string {
	fmt.Println("=== STEP  testing Multiple lacp_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_lacp_policy" "test" {
		name  = "%s"
	}

	resource "aci_lacp_policy" "test1" {
		name  = "%s"
	}

	resource "aci_lacp_policy" "test2" {
		name  = "%s"
	}

	resource "aci_lacp_policy" "test3" {
		name  = "%s"
	}
	`, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccLACPPolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing lacp_policy creation with Invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_lacp_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLACPPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing lacp_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_lacp_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_lacp_policy"
		ctrl = ["load-defer"]
		max_links = "2"
		min_links = "2"
		mode = "active"
	}
	`, rName)

	return resource
}

func CreateAccLACPPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing lacp_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_lacp_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_lacp_policy"
		ctrl = ["load-defer"]
		max_links = "2"
		min_links = "2"
		mode = "active"
	}
	`)

	return resource
}

func CreateAccLACPPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing lacp_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_lacp_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccLACPPolicyUpdatedListAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing lacp_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_lacp_policy" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
