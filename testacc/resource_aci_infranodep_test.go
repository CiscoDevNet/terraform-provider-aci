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

func TestAccAciLeafProfile_Basic(t *testing.T) {
	var leaf_profile_default models.LeafProfile
	var leaf_profile_updated models.LeafProfile
	resourceName := "aci_leaf_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafProfileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLeafProfileWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLeafProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafProfileExists(resourceName, &leaf_profile_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccLeafProfileConfigWithOptionalValuesWithoutSelectorNodeBlock(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafProfileExists(resourceName, &leaf_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_leaf_profile"),
					testAccCheckAciLeafProfileIdEqual(&leaf_profile_default, &leaf_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccLeafProfileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafProfileExists(resourceName, &leaf_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_leaf_profile"),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.switch_association_type", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.node_block.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.node_block.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.node_block.0.from_", "1"),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.node_block.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.node_block.0.to_", "1"),
					testAccCheckAciLeafProfileIdEqual(&leaf_profile_default, &leaf_profile_updated),
				),
			},
			{
				Config:      CreateAccLeafProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLeafProfileWithSelectorWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLeafProfileWithSelectorWithoutSwitchAssocType(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLeafProfileWithNodeblockWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLeafProfileConfigWithOptionalValuesSelectorNodeBlock(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafProfileExists(resourceName, &leaf_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_leaf_profile"),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.description", "leaf_selector description"),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.switch_association_type", "range"),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.node_block.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.node_block.0.description", "node_block description"),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.node_block.0.from_", "1600"),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.node_block.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.node_block.0.to_", "1600"),
					testAccCheckAciLeafProfileIdEqual(&leaf_profile_default, &leaf_profile_updated),
				),
			},
			{
				Config:      CreateAccLeafProfileConfigWithRequiredParams(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config: CreateAccLeafProfileConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafProfileExists(resourceName, &leaf_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLeafProfileIdNotEqual(&leaf_profile_default, &leaf_profile_updated),
				),
			},
		},
	})
}

func TestAccAciLeafProfile_Update(t *testing.T) {
	var leaf_profile_default models.LeafProfile
	var leaf_profile_updated models.LeafProfile
	resourceName := "aci_leaf_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafProfileExists(resourceName, &leaf_profile_default),
				),
			},
			{
				Config: CreateAccLeafProfileUpdatedAttrSetSWAssocType(rName, "ALL_IN_POD"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafProfileExists(resourceName, &leaf_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.switch_association_type", "ALL_IN_POD"),
					testAccCheckAciLeafProfileIdEqual(&leaf_profile_default, &leaf_profile_updated),
				),
			},
			{
				Config: CreateAccLeafProfileUpdatedAttrNode(rName, "to_", "8000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafProfileExists(resourceName, &leaf_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.node_block.0.to_", "8000"),
					testAccCheckAciLeafProfileIdEqual(&leaf_profile_default, &leaf_profile_updated),
				),
			},
			{
				Config: CreateAccLeafProfileUpdatedAttrNode(rName, "from_", "8000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafProfileExists(resourceName, &leaf_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "leaf_selector.0.node_block.0.from_", "8000"),
					testAccCheckAciLeafProfileIdEqual(&leaf_profile_default, &leaf_profile_updated),
				),
			},
		},
	})
}

func TestAccAciLeafProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafProfileConfig(rName),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttrLeafSelector(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttrSetSWAssocType(rName, randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+`),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttrLeafSelector(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttrNode(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttrNode(rName, "from_", "0"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttrNode(rName, "to_", "0"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttrNode(rName, "from_", "16001"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttrNode(rName, "to_", "16001"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttrNode(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},
			{
				Config: CreateAccLeafProfileUpdatedAttrNode(rName, "to_", "100"),
			},
			{
				Config:      CreateAccLeafProfileUpdatedAttrNode(rName, "from_", "200"),
				ExpectError: regexp.MustCompile(`to_ cannot be less than from_`),
			},
		},
	})
}

func TestAccAciLeafProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafProfileConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciLeafProfileExists(name string, leaf_profile *models.LeafProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Leaf Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Leaf Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		leaf_profileFound := models.LeafProfileFromContainer(cont)
		if leaf_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Leaf Profile %s not found", rs.Primary.ID)
		}
		*leaf_profile = *leaf_profileFound
		return nil
	}
}

func testAccCheckAciLeafProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing leaf_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_leaf_profile" {
			cont, err := client.Get(rs.Primary.ID)
			leaf_profile := models.LeafProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Leaf Profile %s Still exists", leaf_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLeafProfileIdEqual(m1, m2 *models.LeafProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("leaf_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLeafProfileIdNotEqual(m1, m2 *models.LeafProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("leaf_profile DNs are equal")
		}
		return nil
	}
}

func CreateAccLeafProfileConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple leaf_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
	
		name  = "%s_${count.index}"
		
		count = 5
	}
	`, rName)
	return resource
}

func CreateLeafProfileWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_profile creation without ", attrName)
	rBlock := `

	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_leaf_profile" "test" {

	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLeafProfileConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing leaf_profile creation with name =", rName)
	resource := fmt.Sprintf(`

	resource "aci_leaf_profile" "test" {

		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafProfileConfig(rName string) string {
	fmt.Println("=== STEP  testing leaf_profile creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_leaf_profile" "test" {

		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafProfileConfigWithOptionalValuesWithoutSelectorNodeBlock(rName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_profile creation with optional parameters of leaf_profile")
	resource := fmt.Sprintf(`

	resource "aci_leaf_profile" "test" {
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_profile"
	}
	`, rName)

	return resource
}

func CreateAccLeafProfileConfigWithOptionalValuesSelectorNodeBlock(rName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_profile creation with optional parameters of leaf_selector and node block")
	resource := fmt.Sprintf(`

	resource "aci_leaf_profile" "test" {
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_profile"
		leaf_selector {
			name = "%s"
			switch_association_type = "range"
			description = "leaf_selector description"
			node_block {
				name = "%s"
				from_ = "1600"
				to_ = "1600"
				description = "node_block description"
			}
		}
	}
	`, rName, rName, rName)

	return resource
}

func CreateAccLeafProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_profile creation with optional parameters")
	resource := fmt.Sprintf(`

	resource "aci_leaf_profile" "test" {
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_profile"
		leaf_selector {
			name = "%s"
			switch_association_type = "ALL"
			node_block {
				name = "%s"
			}
		}
	}
	`, rName, rName, rName)

	return resource
}

func CreateAccLeafProfileWithNodeblockWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_profile without leaf_selector's switch_association_type")
	resource := fmt.Sprintf(`
	resource "aci_leaf_profile" "test" {
		name = "%s"
		leaf_selector {
			name = "%s"
			switch_association_type = "ALL"
			node_block {
			}
		}
	}
	`, rName, rName)

	return resource
}

func CreateAccLeafProfileWithSelectorWithoutSwitchAssocType(rName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_profile without leaf_selector's switch_association_type")
	resource := fmt.Sprintf(`
	resource "aci_leaf_profile" "test" {
		name = "%s"
		leaf_selector {
			name = "%s"
		}
	}
	`, rName, rName)

	return resource
}

func CreateAccLeafProfileWithSelectorWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_profile without leaf_selector's name")
	resource := fmt.Sprintf(`
	resource "aci_leaf_profile" "test" {
		name = "%s"
		leaf_selector {
			switch_association_type = "ALL"
		}
	}
	`, rName)

	return resource
}

func CreateAccLeafProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing leaf_profile update without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_leaf_profile" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_leaf_profile"

	}
	`)

	return resource
}

func CreateAccLeafProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing leaf_profile attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_leaf_profile" "test" {
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccLeafProfileUpdatedAttrLeafSelector(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing leaf_profile's leaf_selector attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_leaf_profile" "test" {
		name  = "%s"
		leaf_selector {
			name = "%s"
			switch_association_type = "range"
			%s = "%s"
		}
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccLeafProfileUpdatedAttrSetSWAssocType(rName, swtype string) string {
	fmt.Printf("=== STEP  testing leaf_profile's leaf_selector for switch_association_type = %s\n", swtype)
	resource := fmt.Sprintf(`
	resource "aci_leaf_profile" "test" {
		name  = "%s"
		leaf_selector {
			name = "%s"
			switch_association_type = "%s"
		}
	}
	`, rName, rName, swtype)
	return resource
}

func CreateAccLeafProfileUpdatedAttrNode(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing leaf_profile's node_block attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_leaf_profile" "test" {
		name  = "%s"
		leaf_selector {
			name = "%s"
			switch_association_type = "range"
			node_block {
				name = "%s"
				%s = "%s"
			}
		}
	}
	`, rName, rName, rName, attribute, value)
	return resource
}
