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

func TestAccAciL3InterfacePolicy_Basic(t *testing.T) {
	var l3_interface_policy_default models.L3InterfacePolicy
	var l3_interface_policy_updated models.L3InterfacePolicy
	resourceName := "aci_l3_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3InterfacePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateL3InterfacePolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3InterfacePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3InterfacePolicyExists(resourceName, &l3_interface_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "bfd_isis", "disabled"),
				),
			},
			{
				Config: CreateAccL3InterfacePolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3InterfacePolicyExists(resourceName, &l3_interface_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3_interface_policy"),

					resource.TestCheckResourceAttr(resourceName, "bfd_isis", "enabled"),

					testAccCheckAciL3InterfacePolicyIdEqual(&l3_interface_policy_default, &l3_interface_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL3InterfacePolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccL3InterfacePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccL3InterfacePolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3InterfacePolicyExists(resourceName, &l3_interface_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciL3InterfacePolicyIdNotEqual(&l3_interface_policy_default, &l3_interface_policy_updated),
				),
			},
		},
	})
}

func TestAccAciL3InterfacePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3InterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3InterfacePolicyConfig(rName),
			},

			{
				Config:      CreateAccL3InterfacePolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3InterfacePolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3InterfacePolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccL3InterfacePolicyUpdatedAttr(rName, "bfd_isis", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL3InterfacePolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3InterfacePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciL3InterfacePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3InterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3InterfacePolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciL3InterfacePolicyExists(name string, l3_interface_policy *models.L3InterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3 Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3 Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3_interface_policyFound := models.L3InterfacePolicyFromContainer(cont)
		if l3_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3 Interface Policy %s not found", rs.Primary.ID)
		}
		*l3_interface_policy = *l3_interface_policyFound
		return nil
	}
}

func testAccCheckAciL3InterfacePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3_interface_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			l3_interface_policy := models.L3InterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3 Interface Policy %s Still exists", l3_interface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3InterfacePolicyIdEqual(m1, m2 *models.L3InterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3_interface_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3InterfacePolicyIdNotEqual(m1, m2 *models.L3InterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3_interface_policy DNs are equal")
		}
		return nil
	}
}

func CreateL3InterfacePolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3_interface_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_l3_interface_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccL3InterfacePolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing l3_interface_policy creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_l3_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccL3InterfacePolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing l3_interface_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_l3_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccL3InterfacePolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing l3_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_l3_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccL3InterfacePolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple l3_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_l3_interface_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccL3InterfacePolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing l3_interface_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_l3_interface_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3_interface_policy"
		bfd_isis = "enabled"
		
	}
	`, rName)

	return resource
}

func CreateAccL3InterfacePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3_interface_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_l3_interface_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3_interface_policy"
		bfd_isis = "enabled"
		
	}
	`)

	return resource
}

func CreateAccL3InterfacePolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3_interface_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_l3_interface_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
