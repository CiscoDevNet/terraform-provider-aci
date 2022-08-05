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

func TestAccAciFEXProfile_Basic(t *testing.T) {
	var fex_profile_default models.FEXProfile
	var fex_profile_updated models.FEXProfile
	resourceName := "aci_fex_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFEXProfileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFEXProfileWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFEXProfileConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFEXProfileExists(resourceName, &fex_profile_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccFEXProfileConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFEXProfileExists(resourceName, &fex_profile_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_fex_profile"),

					testAccCheckAciFEXProfileIdEqual(&fex_profile_default, &fex_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFEXProfileConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccFEXProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccFEXProfileConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFEXProfileExists(resourceName, &fex_profile_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciFEXProfileIdNotEqual(&fex_profile_default, &fex_profile_updated),
				),
			},
		},
	})
}

func TestAccAciFEXProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFEXProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFEXProfileConfig(rName),
			},

			{
				Config:      CreateAccFEXProfileUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFEXProfileUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFEXProfileUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccFEXProfileUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFEXProfileConfig(rName),
			},
		},
	})
}

func TestAccAciFEXProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFEXProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFEXProfileConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciFEXProfileExists(name string, fex_profile *models.FEXProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fex Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fex Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fex_profileFound := models.FEXProfileFromContainer(cont)
		if fex_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fex Profile %s not found", rs.Primary.ID)
		}
		*fex_profile = *fex_profileFound
		return nil
	}
}

func testAccCheckAciFEXProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing fex_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_fex_profile" {
			cont, err := client.Get(rs.Primary.ID)
			fex_profile := models.FEXProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fex Profile %s Still exists", fex_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFEXProfileIdEqual(m1, m2 *models.FEXProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("fex_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFEXProfileIdNotEqual(m1, m2 *models.FEXProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("fex_profile DNs are equal")
		}
		return nil
	}
}

func CreateFEXProfileWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fex_profile creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_fex_profile" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFEXProfileConfigWithRequiredParams(rName string) string {
	fmt.Println("=== STEP  testing fex_profile creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccFEXProfileConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing fex_profile creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFEXProfileConfig(rName string) string {
	fmt.Println("=== STEP  testing fex_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccFEXProfileConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple fex_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccFEXProfileConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing fex_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fex_profile"
		
	}
	`, rName)

	return resource
}

func CreateAccFEXProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing fex_profile updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_fex_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_fex_profile"
		
	}
	`)

	return resource
}

func CreateAccFEXProfileUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing fex_profile attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_fex_profile" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
