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

func TestAccAciSpanningTreeInterfacePolicy_Basic(t *testing.T) {
	var spanning_tree_interface_policy_default models.SpanningTreeInterfacePolicy
	var spanning_tree_interface_policy_updated models.SpanningTreeInterfacePolicy
	resourceName := "aci_spanning_tree_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpanningTreeInterfacePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateSpanningTreeInterfacePolicyWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpanningTreeInterfacePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpanningTreeInterfacePolicyExists(resourceName, &spanning_tree_interface_policy_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "unspecified"),
				),
			},
			{
				Config: CreateAccSpanningTreeInterfacePolicyConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpanningTreeInterfacePolicyExists(resourceName, &spanning_tree_interface_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_spanning_tree_interface_policy"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "bpdu-filter"),

					testAccCheckAciSpanningTreeInterfacePolicyIdEqual(&spanning_tree_interface_policy_default, &spanning_tree_interface_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSpanningTreeInterfacePolicyConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccSpanningTreeInterfacePolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccSpanningTreeInterfacePolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpanningTreeInterfacePolicyExists(resourceName, &spanning_tree_interface_policy_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciSpanningTreeInterfacePolicyIdNotEqual(&spanning_tree_interface_policy_default, &spanning_tree_interface_policy_updated),
				),
			},
		},
	})
}

func TestAccAciSpanningTreeInterfacePolicy_Update(t *testing.T) {
	var spanning_tree_interface_policy_default models.SpanningTreeInterfacePolicy
	var spanning_tree_interface_policy_updated models.SpanningTreeInterfacePolicy
	resourceName := "aci_spanning_tree_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpanningTreeInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpanningTreeInterfacePolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpanningTreeInterfacePolicyExists(resourceName, &spanning_tree_interface_policy_default),
				),
			},

			{
				Config: CreateAccSpanningTreeInterfacePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"bpdu-filter"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpanningTreeInterfacePolicyExists(resourceName, &spanning_tree_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "bpdu-filter"),
				),
			},
			{

				Config: CreateAccSpanningTreeInterfacePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"bpdu-filter", "bpdu-guard"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpanningTreeInterfacePolicyExists(resourceName, &spanning_tree_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "bpdu-filter"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "bpdu-guard"),
				),
			},
			{

				Config: CreateAccSpanningTreeInterfacePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"bpdu-guard"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpanningTreeInterfacePolicyExists(resourceName, &spanning_tree_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "bpdu-guard"),
				),
			},
			{
				Config: CreateAccSpanningTreeInterfacePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"bpdu-guard", "bpdu-filter"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpanningTreeInterfacePolicyExists(resourceName, &spanning_tree_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "bpdu-guard"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "bpdu-filter"),
				),
			},
			{
				Config: CreateAccSpanningTreeInterfacePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciSpanningTreeInterfacePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpanningTreeInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpanningTreeInterfacePolicyConfig(rName),
			},

			{
				Config:      CreateAccSpanningTreeInterfacePolicyUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSpanningTreeInterfacePolicyUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSpanningTreeInterfacePolicyUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSpanningTreeInterfacePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccSpanningTreeInterfacePolicyUpdatedAttrList(rName, "ctrl", StringListtoString([]string{"bpdu-filter", "bpdu-filter"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},

			{
				Config:      CreateAccSpanningTreeInterfacePolicyUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSpanningTreeInterfacePolicyConfig(rName),
			},
		},
	})
}

func TestAccAciSpanningTreeInterfacePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpanningTreeInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpanningTreeInterfacePolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciSpanningTreeInterfacePolicyExists(name string, spanning_tree_interface_policy *models.SpanningTreeInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spanning Tree Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spanning Tree Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spanning_tree_interface_policyFound := models.SpanningTreeInterfacePolicyFromContainer(cont)
		if spanning_tree_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spanning Tree Interface Policy %s not found", rs.Primary.ID)
		}
		*spanning_tree_interface_policy = *spanning_tree_interface_policyFound
		return nil
	}
}

func testAccCheckAciSpanningTreeInterfacePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing spanning_tree_interface_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_spanning_tree_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			spanning_tree_interface_policy := models.SpanningTreeInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spanning Tree Interface Policy %s Still exists", spanning_tree_interface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSpanningTreeInterfacePolicyIdEqual(m1, m2 *models.SpanningTreeInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("spanning_tree_interface_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSpanningTreeInterfacePolicyIdNotEqual(m1, m2 *models.SpanningTreeInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("spanning_tree_interface_policy DNs are equal")
		}
		return nil
	}
}

func CreateSpanningTreeInterfacePolicyWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spanning_tree_interface_policy creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_spanning_tree_interface_policy" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSpanningTreeInterfacePolicyConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing spanning_tree_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spanning_tree_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccSpanningTreeInterfacePolicyConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing spanning_tree_interface_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_spanning_tree_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSpanningTreeInterfacePolicyConfig(rName string) string {
	fmt.Println("=== STEP  testing spanning_tree_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spanning_tree_interface_policy" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSpanningTreeInterfacePolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple spanning_tree_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spanning_tree_interface_policy" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccSpanningTreeInterfacePolicyConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing spanning_tree_interface_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_spanning_tree_interface_policy" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_spanning_tree_interface_policy"
		ctrl = ["bpdu-filter"]
		
	}
	`, rName)

	return resource
}

func CreateAccSpanningTreeInterfacePolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing spanning_tree_interface_policy updation with required parameters")
	resource := fmt.Sprintf(`
	resource "aci_spanning_tree_interface_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_spanning_tree_interface_policy"
		ctrl = ["bpdu-filter"]
		
	}
	`)

	return resource
}

func CreateAccSpanningTreeInterfacePolicyUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing spanning_tree_interface_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_spanning_tree_interface_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccSpanningTreeInterfacePolicyUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing spanning_tree_interface_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_spanning_tree_interface_policy" "test" {
	
		name  = "%s"
		%s = %s
	}
	`, rName, attribute, value)
	return resource
}
