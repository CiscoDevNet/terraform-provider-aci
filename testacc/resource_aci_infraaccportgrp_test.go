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

func TestAccAciLeafAccessPortPolicyGroup_Basic(t *testing.T) {
	var leaf_access_port_policy_group_default models.LeafAccessPortPolicyGroup
	var leaf_access_port_policy_group_updated models.LeafAccessPortPolicyGroup
	resourceName := "aci_leaf_access_port_policy_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafAccessPortPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateLeafAccessPortPolicyGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLeafAccessPortPolicyGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafAccessPortPolicyGroupExists(resourceName, &leaf_access_port_policy_group_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccLeafAccessPortPolicyGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafAccessPortPolicyGroupExists(resourceName, &leaf_access_port_policy_group_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_leaf_access_port_policy_group"),
					testAccCheckAciLeafAccessPortPolicyGroupIdEqual(&leaf_access_port_policy_group_default, &leaf_access_port_policy_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccLeafAccessPortPolicyGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccLeafAccessPortPolicyGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLeafAccessPortPolicyGroupConfigWithUpdatedRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafAccessPortPolicyGroupExists(resourceName, &leaf_access_port_policy_group_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLeafAccessPortPolicyGroupIdNotEqual(&leaf_access_port_policy_group_default, &leaf_access_port_policy_group_updated),
				),
			},
		},
	})
}

func TestAccAciLeafAccessPortPolicyGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafAccessPortPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafAccessPortPolicyGroupConfig(rName),
			},
			{
				Config:      CreateAccLeafAccessPortPolicyGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafAccessPortPolicyGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafAccessPortPolicyGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafAccessPortPolicyGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLeafAccessPortPolicyGroupConfig(rName),
			},
		},
	})
}

func TestAccAciLeafAccessPortPolicyGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafAccessPortPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafAccessPortPolicyGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciLeafAccessPortPolicyGroupExists(name string, leaf_access_port_policy_group *models.LeafAccessPortPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Leaf Access Port Policy Group %s not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Leaf Access Port Policy Group dn was set")
		}
		client := testAccProvider.Meta().(*client.Client)
		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		leaf_access_port_policy_groupFound := models.LeafAccessPortPolicyGroupFromContainer(cont)
		if leaf_access_port_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Leaf Access Port Policy Group %s not found", rs.Primary.ID)
		}
		*leaf_access_port_policy_group = *leaf_access_port_policy_groupFound
		return nil
	}
}

func testAccCheckAciLeafAccessPortPolicyGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing leaf_access_port_policy_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_leaf_access_port_policy_group" {
			cont, err := client.Get(rs.Primary.ID)
			leaf_access_port_policy_group := models.LeafAccessPortPolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Leaf Access Port Policy Group %s Still exists", leaf_access_port_policy_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLeafAccessPortPolicyGroupIdEqual(m1, m2 *models.LeafAccessPortPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("leaf_access_port_policy_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLeafAccessPortPolicyGroupIdNotEqual(m1, m2 *models.LeafAccessPortPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("leaf_access_port_policy_group DNs are equal")
		}
		return nil
	}
}

func CreateLeafAccessPortPolicyGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_access_port_policy_group creation without ", attrName)
	rBlock := `
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_leaf_access_port_policy_group" "test" {
	#	name  = "%s"
	}
	`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLeafAccessPortPolicyGroupConfigWithUpdatedRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing leaf_access_port_policy_group creation with updated required arguments ", rName)
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_port_policy_group" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccLeafAccessPortPolicyGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing leaf_access_port_policy_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_port_policy_group" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafAccessPortPolicyGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing leaf_access_port_policy_group creation with required arguments", rName)
	resource := fmt.Sprintf(`	
	resource "aci_leaf_access_port_policy_group" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafAccessPortPolicyGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple leaf_access_port_policy_group creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_port_policy_group" "test" {
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccLeafAccessPortPolicyGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_access_port_policy_group creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_port_policy_group" "test" {
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_access_port_policy_group"
	}
	`, rName)
	return resource
}

func CreateAccLeafAccessPortPolicyGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing leaf_access_port_policy_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_port_policy_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_access_port_policy_group"
	}
	`)
	return resource
}

func CreateAccLeafAccessPortPolicyGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing leaf_access_port_policy_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_port_policy_group" "test" {
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
