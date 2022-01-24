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

func TestAccAciSpinePortPolicyGroup_Basic(t *testing.T) {
	var spine_port_policy_group_default models.SpineAccessPortPolicyGroup
	var spine_port_policy_group_updated models.SpineAccessPortPolicyGroup
	resourceName := "aci_spine_port_policy_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpinePortPolicyGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateSpinePortPolicyGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpinePortPolicyGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpinePortPolicyGroupExists(resourceName, &spine_port_policy_group_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccSpinePortPolicyGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpinePortPolicyGroupExists(resourceName, &spine_port_policy_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_spine_port_policy_group"),

					testAccCheckAciSpinePortPolicyGroupIdEqual(&spine_port_policy_group_default, &spine_port_policy_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSpinePortPolicyGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccSpinePortPolicyGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccSpinePortPolicyGroupConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpinePortPolicyGroupExists(resourceName, &spine_port_policy_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciSpinePortPolicyGroupIdNotEqual(&spine_port_policy_group_default, &spine_port_policy_group_updated),
				),
			},
		},
	})
}

func TestAccAciSpinePortPolicyGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpinePortPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpinePortPolicyGroupConfig(rName),
			},

			{
				Config:      CreateAccSpinePortPolicyGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSpinePortPolicyGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSpinePortPolicyGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccSpinePortPolicyGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSpinePortPolicyGroupConfig(rName),
			},
		},
	})
}

func TestAccAciSpinePortPolicyGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpinePortPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpinePortPolicyGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciSpinePortPolicyGroupExists(name string, spine_port_policy_group *models.SpineAccessPortPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spine Port Policy Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spine Port Policy Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spine_port_policy_groupFound := models.SpineAccessPortPolicyGroupFromContainer(cont)
		if spine_port_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spine Port Policy Group %s not found", rs.Primary.ID)
		}
		*spine_port_policy_group = *spine_port_policy_groupFound
		return nil
	}
}

func testAccCheckAciSpinePortPolicyGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing spine_port_policy_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_spine_port_policy_group" {
			cont, err := client.Get(rs.Primary.ID)
			spine_port_policy_group := models.SpineAccessPortPolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spine Port Policy Group %s Still exists", spine_port_policy_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSpinePortPolicyGroupIdEqual(m1, m2 *models.SpineAccessPortPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("spine_port_policy_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSpinePortPolicyGroupIdNotEqual(m1, m2 *models.SpineAccessPortPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("spine_port_policy_group DNs are equal")
		}
		return nil
	}
}

func CreateSpinePortPolicyGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spine_port_policy_group creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_spine_port_policy_group" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSpinePortPolicyGroupConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing spine_port_policy_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_port_policy_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccSpinePortPolicyGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing spine_port_policy_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_spine_port_policy_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSpinePortPolicyGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing spine_port_policy_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_port_policy_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSpinePortPolicyGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple spine_port_policy_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_port_policy_group" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccSpinePortPolicyGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing spine_port_policy_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_port_policy_group" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_spine_port_policy_group"
		
	}
	`, rName)

	return resource
}

func CreateAccSpinePortPolicyGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing spine_port_policy_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_spine_port_policy_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_spine_port_policy_group"
		
	}
	`)

	return resource
}

func CreateAccSpinePortPolicyGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing spine_port_policy_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_spine_port_policy_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
