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

func TestAccAciSpineSwitchPolicyGroup_Basic(t *testing.T) {
	var spine_switch_policy_group_default models.SpineSwitchPolicyGroup
	var spine_switch_policy_group_updated models.SpineSwitchPolicyGroup
	resourceName := "aci_spine_switch_policy_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpineSwitchPolicyGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateSpineSwitchPolicyGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpineSwitchPolicyGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineSwitchPolicyGroupExists(resourceName, &spine_switch_policy_group_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccSpineSwitchPolicyGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineSwitchPolicyGroupExists(resourceName, &spine_switch_policy_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_spine_switch_policy_group"),

					testAccCheckAciSpineSwitchPolicyGroupIdEqual(&spine_switch_policy_group_default, &spine_switch_policy_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSpineSwitchPolicyGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccSpineSwitchPolicyGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccSpineSwitchPolicyGroupConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineSwitchPolicyGroupExists(resourceName, &spine_switch_policy_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciSpineSwitchPolicyGroupIdNotEqual(&spine_switch_policy_group_default, &spine_switch_policy_group_updated),
				),
			},
		},
	})
}

func TestAccAciSpineSwitchPolicyGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpineSwitchPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpineSwitchPolicyGroupConfig(rName),
			},

			{
				Config:      CreateAccSpineSwitchPolicyGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSpineSwitchPolicyGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSpineSwitchPolicyGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccSpineSwitchPolicyGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSpineSwitchPolicyGroupConfig(rName),
			},
		},
	})
}

func TestAccAciSpineSwitchPolicyGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpineSwitchPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpineSwitchPolicyGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciSpineSwitchPolicyGroupExists(name string, spine_switch_policy_group *models.SpineSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spine Switch Policy Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spine Switch Policy Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spine_switch_policy_groupFound := models.SpineSwitchPolicyGroupFromContainer(cont)
		if spine_switch_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spine Switch Policy Group %s not found", rs.Primary.ID)
		}
		*spine_switch_policy_group = *spine_switch_policy_groupFound
		return nil
	}
}

func testAccCheckAciSpineSwitchPolicyGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing spine_switch_policy_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_spine_switch_policy_group" {
			cont, err := client.Get(rs.Primary.ID)
			spine_switch_policy_group := models.SpineSwitchPolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spine Switch Policy Group %s Still exists", spine_switch_policy_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSpineSwitchPolicyGroupIdEqual(m1, m2 *models.SpineSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("spine_switch_policy_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSpineSwitchPolicyGroupIdNotEqual(m1, m2 *models.SpineSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("spine_switch_policy_group DNs are equal")
		}
		return nil
	}
}

func CreateSpineSwitchPolicyGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spine_switch_policy_group creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_spine_switch_policy_group" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSpineSwitchPolicyGroupConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing spine_switch_policy_group creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_switch_policy_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccSpineSwitchPolicyGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing spine_switch_policy_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_spine_switch_policy_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSpineSwitchPolicyGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing spine_switch_policy_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_switch_policy_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSpineSwitchPolicyGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple spine_switch_policy_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_switch_policy_group" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccSpineSwitchPolicyGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing spine_switch_policy_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_switch_policy_group" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_spine_switch_policy_group"
		
	}
	`, rName)

	return resource
}

func CreateAccSpineSwitchPolicyGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing spine_switch_policy_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_spine_switch_policy_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_spine_switch_policy_group"
		
	}
	`)

	return resource
}

func CreateAccSpineSwitchPolicyGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing spine_switch_policy_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_spine_switch_policy_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
