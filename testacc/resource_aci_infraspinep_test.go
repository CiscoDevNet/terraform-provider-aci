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

func TestAccAciSpineProfile_Basic(t *testing.T) {
	var spine_profile_default models.SpineProfile
	var spine_profile_updated models.SpineProfile
	resourceName := "aci_spine_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpineProfileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateSpineProfileWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpineProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineProfileExists(resourceName, &spine_profile_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccSpineProfileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineProfileExists(resourceName, &spine_profile_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_spine_profile"),

					testAccCheckAciSpineProfileIdEqual(&spine_profile_default, &spine_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSpineProfileConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccSpineProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccSpineProfileConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineProfileExists(resourceName, &spine_profile_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciSpineProfileIdNotEqual(&spine_profile_default, &spine_profile_updated),
				),
			},
		},
	})
}

func TestAccAciSpineProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpineProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpineProfileConfig(rName),
			},

			{
				Config:      CreateAccSpineProfileUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSpineProfileUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSpineProfileUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccSpineProfileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSpineProfileConfig(rName),
			},
		},
	})
}

func TestAccAciSpineProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpineProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSpineProfileConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciSpineProfileExists(name string, spine_profile *models.SpineProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spine Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spine Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spine_profileFound := models.SpineProfileFromContainer(cont)
		if spine_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spine Profile %s not found", rs.Primary.ID)
		}
		*spine_profile = *spine_profileFound
		return nil
	}
}

func testAccCheckAciSpineProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing spine_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_spine_profile" {
			cont, err := client.Get(rs.Primary.ID)
			spine_profile := models.SpineProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spine Profile %s Still exists", spine_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSpineProfileIdEqual(m1, m2 *models.SpineProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("spine_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSpineProfileIdNotEqual(m1, m2 *models.SpineProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("spine_profile DNs are equal")
		}
		return nil
	}
}

func CreateSpineProfileWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spine_profile creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_spine_profile" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSpineProfileConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing spine_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccSpineProfileConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing spine_profile creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSpineProfileConfig(rName string) string {
	fmt.Println("=== STEP  testing spine_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSpineProfileConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple spine_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccSpineProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing spine_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_spine_profile"
		
	}
	`, rName)

	return resource
}

func CreateAccSpineProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing spine_profile updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_spine_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_spine_profile"
		
	}
	`)

	return resource
}

func CreateAccSpineProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing spine_profile attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
