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

func TestAccAciSAMLProviderGroup_Basic(t *testing.T) {
	var saml_provider_group_default models.SAMLProviderGroup
	var saml_provider_group_updated models.SAMLProviderGroup
	resourceName := "aci_saml_provider_group.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSAMLProviderGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateSAMLProviderGroupWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSAMLProviderGroupConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderGroupExists(resourceName, &saml_provider_group_default),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccSAMLProviderGroupConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderGroupExists(resourceName, &saml_provider_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_saml_provider_group"),

					testAccCheckAciSAMLProviderGroupIdEqual(&saml_provider_group_default, &saml_provider_group_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccSAMLProviderGroupConfigUpdatedName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccSAMLProviderGroupRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccSAMLProviderGroupConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderGroupExists(resourceName, &saml_provider_group_updated),

					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciSAMLProviderGroupIdNotEqual(&saml_provider_group_default, &saml_provider_group_updated),
				),
			},
		},
	})
}

func TestAccAciSAMLProviderGroup_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSAMLProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSAMLProviderGroupConfig(rName),
			},

			{
				Config:      CreateAccSAMLProviderGroupUpdatedAttr(rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSAMLProviderGroupUpdatedAttr(rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccSAMLProviderGroupUpdatedAttr(rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccSAMLProviderGroupUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccSAMLProviderGroupConfig(rName),
			},
		},
	})
}

func TestAccAciSAMLProviderGroup_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSAMLProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccSAMLProviderGroupConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciSAMLProviderGroupExists(name string, saml_provider_group *models.SAMLProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("SAML Provider Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SAML Provider Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		saml_provider_groupFound := models.SAMLProviderGroupFromContainer(cont)
		if saml_provider_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("SAML Provider Group %s not found", rs.Primary.ID)
		}
		*saml_provider_group = *saml_provider_groupFound
		return nil
	}
}

func testAccCheckAciSAMLProviderGroupDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing saml_provider_group destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_saml_provider_group" {
			cont, err := client.Get(rs.Primary.ID)
			saml_provider_group := models.SAMLProviderGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("SAML Provider Group %s Still exists", saml_provider_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSAMLProviderGroupIdEqual(m1, m2 *models.SAMLProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("saml_provider_group DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciSAMLProviderGroupIdNotEqual(m1, m2 *models.SAMLProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("saml_provider_group DNs are equal")
		}
		return nil
	}
}

func CreateSAMLProviderGroupWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing saml_provider_group creation without ", attrName)
	rBlock := `
	
	`
	switch attrName {
	case "name":
		rBlock += `
	resource "aci_saml_provider_group" "test" {
	
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSAMLProviderGroupConfigWithRequiredParams(rName string) string {
	fmt.Printf("=== STEP  testing saml_provider_group creation with name %s\n", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}
func CreateAccSAMLProviderGroupConfigUpdatedName(rName string) string {
	fmt.Println("=== STEP  testing saml_provider_group creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSAMLProviderGroupConfig(rName string) string {
	fmt.Println("=== STEP  testing saml_provider_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider_group" "test" {
	
		name  = "%s"
	}
	`, rName)
	return resource
}

func CreateAccSAMLProviderGroupConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple saml_provider_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider_group" "test" {
	
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName)
	return resource
}

func CreateAccSAMLProviderGroupConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing saml_provider_group creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider_group" "test" {
	
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_saml_provider_group"
		
	}
	`, rName)

	return resource
}

func CreateAccSAMLProviderGroupRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing saml_provider_group updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_saml_provider_group" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_saml_provider_group"
		
	}
	`)

	return resource
}

func CreateAccSAMLProviderGroupUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing saml_provider_group attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_saml_provider_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}
