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

func TestAccAciLeafInterfaceProfile_Basic(t *testing.T) {
	var leaf_interface_profile_default models.LeafInterfaceProfile
	var leaf_interface_profile_updated models.LeafInterfaceProfile
	resourceName := "aci_leaf_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafInterfaceProfileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLeafInterfaceProfileWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLeafInterfaceProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafInterfaceProfileExists(resourceName, &leaf_interface_profile_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccLeafInterfaceProfileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafInterfaceProfileExists(resourceName, &leaf_interface_profile_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_leaf_interface_profile"),

					testAccCheckAciLeafInterfaceProfileIdEqual(&leaf_interface_profile_default, &leaf_interface_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccLeafInterfaceProfileConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccLeafInterfaceProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccLeafInterfaceProfileConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeafInterfaceProfileExists(resourceName, &leaf_interface_profile_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciLeafInterfaceProfileIdNotEqual(&leaf_interface_profile_default, &leaf_interface_profile_updated),
				),
			},
		},
	})
}

func TestAccAciLeafInterfaceProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafInterfaceProfileConfig(rName),
			},

			{
				Config:      CreateAccLeafInterfaceProfileUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafInterfaceProfileUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLeafInterfaceProfileUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLeafInterfaceProfileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLeafInterfaceProfileConfig(rName),
			},
		},
	})
}

func TestAccAciLeafInterfaceProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLeafInterfaceProfileConfigs(rName),
			},
		},
	})
}

func testAccCheckAciLeafInterfaceProfileExists(name string, leaf_interface_profile *models.LeafInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Leaf Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Leaf Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		leaf_interface_profileFound := models.LeafInterfaceProfileFromContainer(cont)
		if leaf_interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Leaf Interface Profile %s not found", rs.Primary.ID)
		}
		*leaf_interface_profile = *leaf_interface_profileFound
		return nil
	}
}

func testAccCheckAciLeafInterfaceProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing leaf_interface_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_leaf_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			leaf_interface_profile := models.LeafInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Leaf Interface Profile %s Still exists", leaf_interface_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLeafInterfaceProfileIdEqual(m1, m2 *models.LeafInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("leaf_interface_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLeafInterfaceProfileIdNotEqual(m1, m2 *models.LeafInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("leaf_interface_profile DNs are equal")
		}
		return nil
	}
}

func CreateLeafInterfaceProfileWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_interface_profile creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_leaf_interface_profile" "test" {
	
	#	name  = "%s"
		description = "created while acceptance testing"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccLeafInterfaceProfileConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing leaf_interface_profile creation with updated name")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafInterfaceProfileConfigs(rName string) string {
	fmt.Println("=== STEP  testing multiple leaf_interface_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_leaf_interface_profile" "test" {
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccLeafInterfaceProfileConfig(rName string) string {
	fmt.Println("=== STEP  testing leaf_interface_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafInterfaceProfileConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing leaf_interface_profile creation with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafInterfaceProfileUpdated(rName string) string {
	fmt.Println("=== STEP  testing leaf_interface_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccLeafInterfaceProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_interface_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_interface_profile"
		
	}
	`, rName)

	return resource
}

func CreateAccLeafInterfaceProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing leaf_interface_profile update without required parameters")
	resource := `
	resource "aci_leaf_interface_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_leaf_interface_profile"
		
	}
	`

	return resource
}

func CreateAccLeafInterfaceProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing leaf_interface_profile attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
