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

func TestAccAciSpineInterfaceProfile_Basic(t *testing.T) {
	var spine_interface_profile_default models.SpineInterfaceProfile
	var spine_interface_profile_updated models.SpineInterfaceProfile
	resourceName := "aci_spine_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpineInterfaceProfileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateSpineInterfaceProfileWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpineInterfaceProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineInterfaceProfileExists(resourceName, &spine_interface_profile_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccSpineInterfaceProfileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineInterfaceProfileExists(resourceName, &spine_interface_profile_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_spine_interface_profile"),

					testAccCheckAciSpineInterfaceProfileIdEqual(&spine_interface_profile_default, &spine_interface_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSpineInterfaceProfileConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccSpineInterfaceProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccSpineInterfaceProfileConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineInterfaceProfileExists(resourceName, &spine_interface_profile_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciSpineInterfaceProfileIdNotEqual(&spine_interface_profile_default, &spine_interface_profile_updated),
				),
			},
		},
	})
}

func TestAccAciSpineInterfaceProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpineInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpineInterfaceProfileConfig(rName),
			},

			{
				Config:      CreateAccSpineInterfaceProfileUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSpineInterfaceProfileUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSpineInterfaceProfileUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccSpineInterfaceProfileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSpineInterfaceProfileConfig(rName),
			},
		},
	})
}

func TestAccAciSpineInterfaceProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpineInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpineInterfaceProfileConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciSpineInterfaceProfileExists(name string, spine_interface_profile *models.SpineInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spine Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spine Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spine_interface_profileFound := models.SpineInterfaceProfileFromContainer(cont)
		if spine_interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spine Interface Profile %s not found", rs.Primary.ID)
		}
		*spine_interface_profile = *spine_interface_profileFound
		return nil
	}
}

func testAccCheckAciSpineInterfaceProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing spine_interface_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_spine_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			spine_interface_profile := models.SpineInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spine Interface Profile %s Still exists", spine_interface_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSpineInterfaceProfileIdEqual(m1, m2 *models.SpineInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("spine_interface_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSpineInterfaceProfileIdNotEqual(m1, m2 *models.SpineInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("spine_interface_profile DNs are equal")
		}
		return nil
	}
}

func CreateSpineInterfaceProfileWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spine_interface_profile creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_spine_interface_profile" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSpineInterfaceProfileConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing spine_interface_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_interface_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccSpineInterfaceProfileConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing spine_interface_profile creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_spine_interface_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSpineInterfaceProfileConfig(rName string) string {
	fmt.Println("=== STEP  testing spine_interface_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_interface_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSpineInterfaceProfileConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple spine_interface_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_interface_profile" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccSpineInterfaceProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing spine_interface_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_interface_profile" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_spine_interface_profile"
		
	}
	`, rName)

	return resource
}

func CreateAccSpineInterfaceProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing spine_interface_profile updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_spine_interface_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_spine_interface_profile"
		
	}
	`)

	return resource
}

func CreateAccSpineInterfaceProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing spine_interface_profile attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_spine_interface_profile" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
