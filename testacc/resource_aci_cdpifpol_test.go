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

func TestAccAciCDPInterfacePolicy_Basic(t *testing.T) {
	var cdp_interface_policy_default models.CDPInterfacePolicy
	var cdp_interface_policy_updated models.CDPInterfacePolicy
	resourceName := "aci_cdp_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCDPInterfacePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateCDPInterfacePolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCDPInterfacePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCDPInterfacePolicyExists(resourceName, &cdp_interface_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "enabled"),
				),
			},
			{
				Config: CreateAccCDPInterfacePolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCDPInterfacePolicyExists(resourceName, &cdp_interface_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_cdp_interface_policy"),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "disabled"),

					testAccCheckAciCDPInterfacePolicyIdEqual(&cdp_interface_policy_default, &cdp_interface_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccCDPInterfacePolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccCDPInterfacePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccCDPInterfacePolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCDPInterfacePolicyExists(resourceName, &cdp_interface_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciCDPInterfacePolicyIdNotEqual(&cdp_interface_policy_default, &cdp_interface_policy_updated),
				),
			},
		},
	})
}

func TestAccAciCDPInterfacePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCDPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCDPInterfacePolicyConfigs(rName),
			},
		},
	})
}

func TestAccAciCDPInterfacePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCDPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCDPInterfacePolicyConfig(rName),
			},

			{
				Config:      CreateAccCDPInterfacePolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCDPInterfacePolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCDPInterfacePolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCDPInterfacePolicyUpdatedAttr(rName, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccCDPInterfacePolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCDPInterfacePolicyConfig(rName),
			},
		},
	})
}

func testAccCheckAciCDPInterfacePolicyExists(name string, cdp_interface_policy *models.CDPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("CDP Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CDP Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cdp_interface_policyFound := models.CDPInterfacePolicyFromContainer(cont)
		if cdp_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("CDP Interface Policy %s not found", rs.Primary.ID)
		}
		*cdp_interface_policy = *cdp_interface_policyFound
		return nil
	}
}

func testAccCheckAciCDPInterfacePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing cdp_interface_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_cdp_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			cdp_interface_policy := models.CDPInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("CDP Interface Policy %s Still exists", cdp_interface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciCDPInterfacePolicyIdEqual(m1, m2 *models.CDPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("cdp_interface_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciCDPInterfacePolicyIdNotEqual(m1, m2 *models.CDPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("cdp_interface_policy DNs are equal")
		}
		return nil
	}
}

func CreateCDPInterfacePolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cdp_interface_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_cdp_interface_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccCDPInterfacePolicyConfigWithRequiredParams(rName string) string {
	fmt.Printf("=== STEP  testing cdp_interface_policy creation with name = %s\n", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_cdp_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccCDPInterfacePolicyConfigs(rName string) string {
	fmt.Println("=== STEP  testing multiple cdp_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_cdp_interface_policy" "test" {
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccCDPInterfacePolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing cdp_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_cdp_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccCDPInterfacePolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing cdp_interface_policy creation with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_cdp_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccCDPInterfacePolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing cdp_interface_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_cdp_interface_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cdp_interface_policy"
		admin_st = "disabled"
	}
	`, rName)

	return resource
}

func CreateAccCDPInterfacePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing cdp_interface_policy update without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_cdp_interface_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_cdp_interface_policy"
		admin_st = "disabled"
	}
	`)

	return resource
}

func CreateAccCDPInterfacePolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing cdp_interface_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_cdp_interface_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
