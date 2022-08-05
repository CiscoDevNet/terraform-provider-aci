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

func TestAccAciAccessSwitchPolicyGroup_Basic(t *testing.T) {
	var access_switch_policy_group_default models.AccessSwitchPolicyGroup
	var access_switch_policy_group_updated models.AccessSwitchPolicyGroup
	resourceName := "aci_access_switch_policy_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessSwitchPolicyGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateAccessSwitchPolicyGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAccessSwitchPolicyGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessSwitchPolicyGroupExists(resourceName, &access_switch_policy_group_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccAccessSwitchPolicyGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessSwitchPolicyGroupExists(resourceName, &access_switch_policy_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_access_switch_policy_group"),

					testAccCheckAciAccessSwitchPolicyGroupIdEqual(&access_switch_policy_group_default, &access_switch_policy_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"relation_infra_rs_netflow_node_pol",
				},
			},
			{
				Config:      CreateAccAccessSwitchPolicyGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccAccessSwitchPolicyGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccAccessSwitchPolicyGroupConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessSwitchPolicyGroupExists(resourceName, &access_switch_policy_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciAccessSwitchPolicyGroupIdNotEqual(&access_switch_policy_group_default, &access_switch_policy_group_updated),
				),
			},
		},
	})
}

func TestAccAciAccessSwitchPolicyGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessSwitchPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAccessSwitchPolicyGroupConfig(rName),
			},

			{
				Config:      CreateAccAccessSwitchPolicyGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAccessSwitchPolicyGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAccessSwitchPolicyGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccAccessSwitchPolicyGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccAccessSwitchPolicyGroupConfig(rName),
			},
		},
	})
}

func TestAccAciAccessSwitchPolicyGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessSwitchPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAccessSwitchPolicyGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciAccessSwitchPolicyGroupExists(name string, access_switch_policy_group *models.AccessSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Access Switch Policy Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access Switch Policy Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		access_switch_policy_groupFound := models.AccessSwitchPolicyGroupFromContainer(cont)
		if access_switch_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Access Switch Policy Group %s not found", rs.Primary.ID)
		}
		*access_switch_policy_group = *access_switch_policy_groupFound
		return nil
	}
}

func testAccCheckAciAccessSwitchPolicyGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing access_switch_policy_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_access_switch_policy_group" {
			cont, err := client.Get(rs.Primary.ID)
			access_switch_policy_group := models.AccessSwitchPolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Access Switch Policy Group %s Still exists", access_switch_policy_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAccessSwitchPolicyGroupIdEqual(m1, m2 *models.AccessSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("access_switch_policy_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciAccessSwitchPolicyGroupIdNotEqual(m1, m2 *models.AccessSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("access_switch_policy_group DNs are equal")
		}
		return nil
	}
}

func CreateAccessSwitchPolicyGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing access_switch_policy_group creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_access_switch_policy_group" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccAccessSwitchPolicyGroupConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing access_switch_policy_group creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_access_switch_policy_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccAccessSwitchPolicyGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing access_switch_policy_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_access_switch_policy_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccAccessSwitchPolicyGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing access_switch_policy_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_access_switch_policy_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccAccessSwitchPolicyGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple access_switch_policy_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_access_switch_policy_group" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccAccessSwitchPolicyGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing access_switch_policy_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_access_switch_policy_group" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_access_switch_policy_group"
		
	}
	`, rName)

	return resource
}

func CreateAccAccessSwitchPolicyGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing access_switch_policy_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_access_switch_policy_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_access_switch_policy_group"
		
	}
	`)

	return resource
}

func CreateAccAccessSwitchPolicyGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing access_switch_policy_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_access_switch_policy_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
