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

func TestAccAciMiscablingProtocolInterfacePolicy_Basic(t *testing.T) {
	var miscabling_protocol_interface_policy_default models.MiscablingProtocolInterfacePolicy
	var miscabling_protocol_interface_policy_updated models.MiscablingProtocolInterfacePolicy
	resourceName := "aci_miscabling_protocol_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMiscablingProtocolInterfacePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateMiscablingProtocolInterfacePolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMiscablingProtocolInterfacePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInterfacePolicyExists(resourceName, &miscabling_protocol_interface_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "admin_st", "enabled"),
				),
			},
			{
				Config: CreateAccMiscablingProtocolInterfacePolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInterfacePolicyExists(resourceName, &miscabling_protocol_interface_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_miscabling_protocol_interface_policy"),

					resource.TestCheckResourceAttr(resourceName, "admin_st", "disabled"),

					testAccCheckAciMiscablingProtocolInterfacePolicyIdEqual(&miscabling_protocol_interface_policy_default, &miscabling_protocol_interface_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccMiscablingProtocolInterfacePolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccMiscablingProtocolInterfacePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccMiscablingProtocolInterfacePolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInterfacePolicyExists(resourceName, &miscabling_protocol_interface_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciMiscablingProtocolInterfacePolicyIdNotEqual(&miscabling_protocol_interface_policy_default, &miscabling_protocol_interface_policy_updated),
				),
			},
		},
	})
}

func TestAccAciMiscablingProtocolInterfacePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMiscablingProtocolInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMiscablingProtocolInterfacePolicyConfig(rName),
			},

			{
				Config:      CreateAccMiscablingProtocolInterfacePolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMiscablingProtocolInterfacePolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccMiscablingProtocolInterfacePolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccMiscablingProtocolInterfacePolicyUpdatedAttr(rName, "admin_st", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccMiscablingProtocolInterfacePolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccMiscablingProtocolInterfacePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciMiscablingProtocolInterfacePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMiscablingProtocolInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMiscablingProtocolInterfacePolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciMiscablingProtocolInterfacePolicyExists(name string, miscabling_protocol_interface_policy *models.MiscablingProtocolInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Miscabling Protocol Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Miscabling Protocol Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		miscabling_protocol_interface_policyFound := models.MiscablingProtocolInterfacePolicyFromContainer(cont)
		if miscabling_protocol_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Miscabling Protocol Interface Policy %s not found", rs.Primary.ID)
		}
		*miscabling_protocol_interface_policy = *miscabling_protocol_interface_policyFound
		return nil
	}
}

func testAccCheckAciMiscablingProtocolInterfacePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing miscabling_protocol_interface_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_miscabling_protocol_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			miscabling_protocol_interface_policy := models.MiscablingProtocolInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Miscabling Protocol Interface Policy %s Still exists", miscabling_protocol_interface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciMiscablingProtocolInterfacePolicyIdEqual(m1, m2 *models.MiscablingProtocolInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("miscabling_protocol_interface_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciMiscablingProtocolInterfacePolicyIdNotEqual(m1, m2 *models.MiscablingProtocolInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("miscabling_protocol_interface_policy DNs are equal")
		}
		return nil
	}
}

func CreateMiscablingProtocolInterfacePolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing miscabling_protocol_interface_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_miscabling_protocol_interface_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccMiscablingProtocolInterfacePolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing miscabling_protocol_interface_policy creation with updated name", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccMiscablingProtocolInterfacePolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing miscabling_protocol_interface_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccMiscablingProtocolInterfacePolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing miscabling_protocol_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccMiscablingProtocolInterfacePolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple miscabling_protocol_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccMiscablingProtocolInterfacePolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing miscabling_protocol_interface_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_miscabling_protocol_interface_policy"
		admin_st = "disabled"
		
	}
	`, rName)

	return resource
}

func CreateAccMiscablingProtocolInterfacePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing miscabling_protocol_interface_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_miscabling_protocol_interface_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_miscabling_protocol_interface_policy"
		admin_st = "disabled"
		
	}
	`)

	return resource
}

func CreateAccMiscablingProtocolInterfacePolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing miscabling_protocol_interface_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
