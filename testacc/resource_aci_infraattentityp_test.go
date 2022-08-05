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

func TestAccAciAttachableAccessEntityProfile_Basic(t *testing.T) {
	var attachable_access_entity_profile_default models.AttachableAccessEntityProfile
	var attachable_access_entity_profile_updated models.AttachableAccessEntityProfile
	resourceName := "aci_attachable_access_entity_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAttachableAccessEntityProfileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateAttachableAccessEntityProfileWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAttachableAccessEntityProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAttachableAccessEntityProfileExists(resourceName, &attachable_access_entity_profile_default),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccAttachableAccessEntityProfileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAttachableAccessEntityProfileExists(resourceName, &attachable_access_entity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_attachable_access_entity_profile"),
					testAccCheckAciAttachableAccessEntityProfileIdEqual(&attachable_access_entity_profile_default, &attachable_access_entity_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccAttachableAccessEntityProfileConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccAttachableAccessEntityProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAttachableAccessEntityProfileConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAttachableAccessEntityProfileExists(resourceName, &attachable_access_entity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciAttachableAccessEntityProfileIdNotEqual(&attachable_access_entity_profile_default, &attachable_access_entity_profile_updated),
				),
			},
		},
	})
}

func TestAccAciAttachableAccessEntityProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAttachableAccessEntityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAttachableAccessEntityProfileConfig(rName),
			},

			{
				Config:      CreateAccAttachableAccessEntityProfileUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAttachableAccessEntityProfileUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAttachableAccessEntityProfileUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAttachableAccessEntityProfileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccAttachableAccessEntityProfileConfig(rName),
			},
		},
	})
}

func TestAccAciAttachableAccessEntityProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAttachableAccessEntityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAttachableAccessEntityProfileConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciAttachableAccessEntityProfileExists(name string, attachable_access_entity_profile *models.AttachableAccessEntityProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Attachable Access Entity Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Attachable Access Entity Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		attachable_access_entity_profileFound := models.AttachableAccessEntityProfileFromContainer(cont)
		if attachable_access_entity_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Attachable Access Entity Profile %s not found", rs.Primary.ID)
		}
		*attachable_access_entity_profile = *attachable_access_entity_profileFound
		return nil
	}
}

func testAccCheckAciAttachableAccessEntityProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing attachable_access_entity_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_attachable_access_entity_profile" {
			cont, err := client.Get(rs.Primary.ID)
			attachable_access_entity_profile := models.AttachableAccessEntityProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Attachable Access Entity Profile %s Still exists", attachable_access_entity_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAttachableAccessEntityProfileIdEqual(m1, m2 *models.AttachableAccessEntityProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("attachable_access_entity_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciAttachableAccessEntityProfileIdNotEqual(m1, m2 *models.AttachableAccessEntityProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("attachable_access_entity_profile DNs are equal")
		}
		return nil
	}
}

func CreateAttachableAccessEntityProfileWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing attachable_access_entity_profile creation without ", attrName)
	rBlock := `
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_attachable_access_entity_profile" "test" {
	#	name  = "%s"
	}
	`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccAttachableAccessEntityProfileConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing attachable_access_entity_profile creation with Updated required arguments only ", rName)
	resource := fmt.Sprintf(`
	resource "aci_attachable_access_entity_profile" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccAttachableAccessEntityProfileConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing attachable_access_entity_profile creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	resource "aci_attachable_access_entity_profile" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccAttachableAccessEntityProfileConfig(rName string) string {
	fmt.Println("=== STEP  testing attachable_access_entity_profile creation with required arguments only ", rName)
	resource := fmt.Sprintf(`
	resource "aci_attachable_access_entity_profile" "test" {
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccAttachableAccessEntityProfileConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple attachable_access_entity_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_attachable_access_entity_profile" "test" {
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccAttachableAccessEntityProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing attachable_access_entity_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_attachable_access_entity_profile" "test" {
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_attachable_access_entity_profile"
	}
	`, rName)

	return resource
}

func CreateAccAttachableAccessEntityProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing attachable_access_entity_profile updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_attachable_access_entity_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_attachable_access_entity_profile"
	}
	`)

	return resource
}

func CreateAccAttachableAccessEntityProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing attachable_access_entity_profile attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_attachable_access_entity_profile" "test" {	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
