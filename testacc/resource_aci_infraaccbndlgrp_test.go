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

func TestAccAciLeafAccessBundlePolicyGroup_Basic(t *testing.T) {
	var leaf_access_bundle_policy_group_default models.PCVPCInterfacePolicyGroup
	var leaf_access_bundle_policy_group_updated models.PCVPCInterfacePolicyGroup
	resourceName := "aci_leaf_access_bundle_policy_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafAccessBundlePolicyGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLeafAccessBundlePolicyGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLeafAccessBundlePolicyGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafAccessBundlePolicyGroupExists(resourceName, &leaf_access_bundle_policy_group_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "lag_t", "link"),
				),
			},
			{
				Config: CreateAccLeafAccessBundlePolicyGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafAccessBundlePolicyGroupExists(resourceName, &leaf_access_bundle_policy_group_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_leaf_access_bundle_policy_group"),
					resource.TestCheckResourceAttr(resourceName, "lag_t", "node"),
					testAccCheckAciLeafAccessBundlePolicyGroupIdEqual(&leaf_access_bundle_policy_group_default, &leaf_access_bundle_policy_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccLeafAccessBundlePolicyGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccLeafAccessBundlePolicyGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLeafAccessBundlePolicyGroupConfigWithUpdatedRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafAccessBundlePolicyGroupExists(resourceName, &leaf_access_bundle_policy_group_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLeafAccessBundlePolicyGroupIdNotEqual(&leaf_access_bundle_policy_group_default, &leaf_access_bundle_policy_group_updated),
				),
			},
		},
	})
}

func TestAccAciLeafAccessBundlePolicyGroup_Update(t *testing.T) {
	var leaf_access_bundle_policy_group_default models.PCVPCInterfacePolicyGroup
	var leaf_access_bundle_policy_group_updated models.PCVPCInterfacePolicyGroup
	resourceName := "aci_leaf_access_bundle_policy_group.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafAccessBundlePolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafAccessBundlePolicyGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafAccessBundlePolicyGroupExists(resourceName, &leaf_access_bundle_policy_group_default),
				),
			},
			{
				Config: CreateAccLeafAccessBundlePolicyGroupUpdatedAttr(rName, "lag_t", "not-aggregated"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafAccessBundlePolicyGroupExists(resourceName, &leaf_access_bundle_policy_group_updated),
					resource.TestCheckResourceAttr(resourceName, "lag_t", "not-aggregated"),
					testAccCheckAciLeafAccessBundlePolicyGroupIdEqual(&leaf_access_bundle_policy_group_default, &leaf_access_bundle_policy_group_updated),
				),
			},
			{
				Config: CreateAccLeafAccessBundlePolicyGroupConfig(rName),
			},
		},
	})
}

func TestAccAciLeafAccessBundlePolicyGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafAccessBundlePolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafAccessBundlePolicyGroupConfig(rName),
			},
			{
				Config:      CreateAccLeafAccessBundlePolicyGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafAccessBundlePolicyGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafAccessBundlePolicyGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafAccessBundlePolicyGroupUpdatedAttr(rName, "lag_t", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccLeafAccessBundlePolicyGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLeafAccessBundlePolicyGroupConfig(rName),
			},
		},
	})
}

func TestAccAciLeafAccessBundlePolicyGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafAccessBundlePolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafAccessBundlePolicyGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciLeafAccessBundlePolicyGroupExists(name string, leaf_access_bundle_policy_group *models.PCVPCInterfacePolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Leaf Access Bundle Policy Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Leaf Access Bundle Policy Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		leaf_access_bundle_policy_groupFound := models.PCVPCInterfacePolicyGroupFromContainer(cont)
		if leaf_access_bundle_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Leaf Access Bundle Policy Group %s not found", rs.Primary.ID)
		}
		*leaf_access_bundle_policy_group = *leaf_access_bundle_policy_groupFound
		return nil
	}
}

func testAccCheckAciLeafAccessBundlePolicyGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing leaf_access_bundle_policy_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_leaf_access_bundle_policy_group" {
			cont, err := client.Get(rs.Primary.ID)
			leaf_access_bundle_policy_group := models.PCVPCInterfacePolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Leaf Access Bundle Policy Group %s Still exists", leaf_access_bundle_policy_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLeafAccessBundlePolicyGroupIdEqual(m1, m2 *models.PCVPCInterfacePolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("leaf_access_bundle_policy_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLeafAccessBundlePolicyGroupIdNotEqual(m1, m2 *models.PCVPCInterfacePolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("leaf_access_bundle_policy_group DNs are equal")
		}
		return nil
	}
}

func CreateLeafAccessBundlePolicyGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_access_bundle_policy_group creation without ", attrName)
	rBlock := `
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_leaf_access_bundle_policy_group" "test" {
	#	name  = "%s"
	}
	`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLeafAccessBundlePolicyGroupConfigWithUpdatedRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing leaf_access_bundle_policy_group creation with Updated required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_bundle_policy_group" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafAccessBundlePolicyGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing leaf_access_bundle_policy_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_bundle_policy_group" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafAccessBundlePolicyGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing leaf_access_bundle_policy_group creation with required arguments ", rName)
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_bundle_policy_group" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafAccessBundlePolicyGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple leaf_access_bundle_policy_group creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_bundle_policy_group" "test" {
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccLeafAccessBundlePolicyGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_access_bundle_policy_group creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_bundle_policy_group" "test" {
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_access_bundle_policy_group"
		lag_t = "node"
	}
	`, rName)

	return resource
}

func CreateAccLeafAccessBundlePolicyGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing leaf_access_bundle_policy_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_bundle_policy_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_access_bundle_policy_group"
		lag_t = "node"
	}
	`)
	return resource
}

func CreateAccLeafAccessBundlePolicyGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing leaf_access_bundle_policy_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_leaf_access_bundle_policy_group" "test" {
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
