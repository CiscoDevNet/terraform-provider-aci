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

func TestAccAciL2InterfacePolicy_Basic(t *testing.T) {
	var l2_interface_policy_default models.L2InterfacePolicy
	var l2_interface_policy_updated models.L2InterfacePolicy
	resourceName := "aci_l2_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2InterfacePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateL2InterfacePolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL2InterfacePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2InterfacePolicyExists(resourceName, &l2_interface_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "qinq", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "vepa", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "vlan_scope", "global"),
				),
			},
			{
				Config: CreateAccL2InterfacePolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2InterfacePolicyExists(resourceName, &l2_interface_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l2_interface_policy"),

					resource.TestCheckResourceAttr(resourceName, "qinq", "corePort"),

					resource.TestCheckResourceAttr(resourceName, "vepa", "enabled"),

					resource.TestCheckResourceAttr(resourceName, "vlan_scope", "portlocal"),

					testAccCheckAciL2InterfacePolicyIdEqual(&l2_interface_policy_default, &l2_interface_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL2InterfacePolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccL2InterfacePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccL2InterfacePolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2InterfacePolicyExists(resourceName, &l2_interface_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciL2InterfacePolicyIdNotEqual(&l2_interface_policy_default, &l2_interface_policy_updated),
				),
			},
		},
	})
}

func TestAccAciL2InterfacePolicy_Update(t *testing.T) {
	var l2_interface_policy_default models.L2InterfacePolicy
	var l2_interface_policy_updated models.L2InterfacePolicy
	resourceName := "aci_l2_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2InterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL2InterfacePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2InterfacePolicyExists(resourceName, &l2_interface_policy_default),
				),
			},

			{
				Config: CreateAccL2InterfacePolicyUpdatedAttr(rName, "qinq", "doubleQtagPort"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2InterfacePolicyExists(resourceName, &l2_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "qinq", "doubleQtagPort"),
					testAccCheckAciL2InterfacePolicyIdEqual(&l2_interface_policy_default, &l2_interface_policy_updated),
				),
			},
			{
				Config: CreateAccL2InterfacePolicyUpdatedAttr(rName, "qinq", "edgePort"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2InterfacePolicyExists(resourceName, &l2_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "qinq", "edgePort"),
					testAccCheckAciL2InterfacePolicyIdEqual(&l2_interface_policy_default, &l2_interface_policy_updated),
				),
			},
			{
				Config: CreateAccL2InterfacePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciL2InterfacePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2InterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL2InterfacePolicyConfig(rName),
			},

			{
				Config:      CreateAccL2InterfacePolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL2InterfacePolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL2InterfacePolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccL2InterfacePolicyUpdatedAttr(rName, "qinq", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL2InterfacePolicyUpdatedAttr(rName, "vepa", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL2InterfacePolicyUpdatedAttr(rName, "vlan_scope", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL2InterfacePolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL2InterfacePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciL2InterfacePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2InterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL2InterfacePolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciL2InterfacePolicyExists(name string, l2_interface_policy *models.L2InterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L2 Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L2 Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l2_interface_policyFound := models.L2InterfacePolicyFromContainer(cont)
		if l2_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L2 Interface Policy %s not found", rs.Primary.ID)
		}
		*l2_interface_policy = *l2_interface_policyFound
		return nil
	}
}

func testAccCheckAciL2InterfacePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l2_interface_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l2_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			l2_interface_policy := models.L2InterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L2 Interface Policy %s Still exists", l2_interface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL2InterfacePolicyIdEqual(m1, m2 *models.L2InterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l2_interface_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL2InterfacePolicyIdNotEqual(m1, m2 *models.L2InterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l2_interface_policy DNs are equal")
		}
		return nil
	}
}

func CreateL2InterfacePolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l2_interface_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_l2_interface_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccL2InterfacePolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing l2_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccL2InterfacePolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing l2_interface_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_l2_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccL2InterfacePolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing l2_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccL2InterfacePolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple l2_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_interface_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccL2InterfacePolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing l2_interface_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_interface_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l2_interface_policy"
		qinq = "corePort"
		vepa = "enabled"
		vlan_scope = "portlocal"
		
	}
	`, rName)

	return resource
}

func CreateAccL2InterfacePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l2_interface_policy updation with required parameters")
	resource := fmt.Sprintf(`
	resource "aci_l2_interface_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l2_interface_policy"
		qinq = "corePort"
		vepa = "enabled"
		vlan_scope = "portlocal"
		
	}
	`)

	return resource
}

func CreateAccL2InterfacePolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l2_interface_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_l2_interface_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
