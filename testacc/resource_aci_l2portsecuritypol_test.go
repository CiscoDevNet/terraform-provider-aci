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

func TestAccAciPortSecurityPolicy_Basic(t *testing.T) {
	var port_security_policy_default models.PortSecurityPolicy
	var port_security_policy_updated models.PortSecurityPolicy
	resourceName := "aci_port_security_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPortSecurityPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreatePortSecurityPolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccPortSecurityPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortSecurityPolicyExists(resourceName, &port_security_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "maximum", "0"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "60"),
					resource.TestCheckResourceAttr(resourceName, "violation", "protect"),
				),
			},
			{
				Config: CreateAccPortSecurityPolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortSecurityPolicyExists(resourceName, &port_security_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_port_security_policy"),
					resource.TestCheckResourceAttr(resourceName, "maximum", "1"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "61"),

					resource.TestCheckResourceAttr(resourceName, "violation", "protect"),

					testAccCheckAciPortSecurityPolicyIdEqual(&port_security_policy_default, &port_security_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccPortSecurityPolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccPortSecurityPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccPortSecurityPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortSecurityPolicyExists(resourceName, &port_security_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciPortSecurityPolicyIdNotEqual(&port_security_policy_default, &port_security_policy_updated),
				),
			},
		},
	})
}

func TestAccAciPortSecurityPolicy_Update(t *testing.T) {
	var port_security_policy_default models.PortSecurityPolicy
	var port_security_policy_updated models.PortSecurityPolicy
	resourceName := "aci_port_security_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPortSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPortSecurityPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortSecurityPolicyExists(resourceName, &port_security_policy_default),
				),
			},
			{
				Config: CreateAccPortSecurityPolicyUpdatedAttr(rName, "maximum", "12000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortSecurityPolicyExists(resourceName, &port_security_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "maximum", "12000"),
					testAccCheckAciPortSecurityPolicyIdEqual(&port_security_policy_default, &port_security_policy_updated),
				),
			},
			{
				Config: CreateAccPortSecurityPolicyUpdatedAttr(rName, "maximum", "6000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortSecurityPolicyExists(resourceName, &port_security_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "maximum", "6000"),
					testAccCheckAciPortSecurityPolicyIdEqual(&port_security_policy_default, &port_security_policy_updated),
				),
			},
			{
				Config: CreateAccPortSecurityPolicyUpdatedAttr(rName, "timeout", "3600"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortSecurityPolicyExists(resourceName, &port_security_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "3600"),
					testAccCheckAciPortSecurityPolicyIdEqual(&port_security_policy_default, &port_security_policy_updated),
				),
			},
			{
				Config: CreateAccPortSecurityPolicyUpdatedAttr(rName, "timeout", "1770"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortSecurityPolicyExists(resourceName, &port_security_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "1770"),
					testAccCheckAciPortSecurityPolicyIdEqual(&port_security_policy_default, &port_security_policy_updated),
				),
			},

			{
				Config: CreateAccPortSecurityPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciPortSecurityPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPortSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPortSecurityPolicyConfig(rName),
			},

			{
				Config:      CreateAccPortSecurityPolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPortSecurityPolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPortSecurityPolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccPortSecurityPolicyUpdatedAttr(rName, "maximum", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccPortSecurityPolicyUpdatedAttr(rName, "maximum", "12001"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccPortSecurityPolicyUpdatedAttr(rName, "timeout", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccPortSecurityPolicyUpdatedAttr(rName, "timeout", "59"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccPortSecurityPolicyUpdatedAttr(rName, "timeout", "3601"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccPortSecurityPolicyUpdatedAttr(rName, "violation", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccPortSecurityPolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccPortSecurityPolicyConfig(rName),
			},
		},
	})
}

func TestAccAciPortSecurityPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPortSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPortSecurityPolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciPortSecurityPolicyExists(name string, port_security_policy *models.PortSecurityPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Port Security Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Port Security Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		port_security_policyFound := models.PortSecurityPolicyFromContainer(cont)
		if port_security_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Port Security Policy %s not found", rs.Primary.ID)
		}
		*port_security_policy = *port_security_policyFound
		return nil
	}
}

func testAccCheckAciPortSecurityPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing port_security_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_port_security_policy" {
			cont, err := client.Get(rs.Primary.ID)
			port_security_policy := models.PortSecurityPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Port Security Policy %s Still exists", port_security_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciPortSecurityPolicyIdEqual(m1, m2 *models.PortSecurityPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("port_security_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciPortSecurityPolicyIdNotEqual(m1, m2 *models.PortSecurityPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("port_security_policy DNs are equal")
		}
		return nil
	}
}

func CreatePortSecurityPolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing port_security_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_port_security_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccPortSecurityPolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing port_security_policy creation with updated name")
	resource := fmt.Sprintf(`
	
	resource "aci_port_security_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccPortSecurityPolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing port_security_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_port_security_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccPortSecurityPolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing port_security_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_port_security_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccPortSecurityPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple port_security_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_port_security_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccPortSecurityPolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing port_security_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_port_security_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_port_security_policy"
		maximum = "1"
		timeout = "61"
		violation = "protect"
		
	}
	`, rName)

	return resource
}

func CreateAccPortSecurityPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing port_security_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_port_security_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_port_security_policy"
		maximum = "1"
		timeout = "61"
		violation = "protect"
		
	}
	`)

	return resource
}

func CreateAccPortSecurityPolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing port_security_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_port_security_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
