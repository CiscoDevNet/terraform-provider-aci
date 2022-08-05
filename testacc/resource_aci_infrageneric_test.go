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

func TestAccAciAccessGeneric_Basic(t *testing.T) {
	var access_generic_default models.AccessGeneric
	var access_generic_updated models.AccessGeneric
	resourceName := "aci_access_generic.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	infraAttEntityPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessGenericDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccessGenericWithoutRequired(infraAttEntityPName, "default", "attachable_access_entity_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccessGenericWithoutRequired(infraAttEntityPName, "default", "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAccessGenericConfig(infraAttEntityPName, "default"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessGenericExists(resourceName, &access_generic_default),
					resource.TestCheckResourceAttr(resourceName, "attachable_access_entity_profile_dn", fmt.Sprintf("uni/infra/attentp-%s", infraAttEntityPName)),
					resource.TestCheckResourceAttr(resourceName, "name", "default"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccAccessGenericConfigWithOptionalValues(infraAttEntityPName, "default"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessGenericExists(resourceName, &access_generic_updated),
					resource.TestCheckResourceAttr(resourceName, "attachable_access_entity_profile_dn", fmt.Sprintf("uni/infra/attentp-%s", infraAttEntityPName)),
					resource.TestCheckResourceAttr(resourceName, "name", "default"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_access_generic"),

					testAccCheckAciAccessGenericIdEqual(&access_generic_default, &access_generic_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccAccessGenericConfigUpdatedName(infraAttEntityPName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccAccessGenericRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAccessGenericConfigWithRequiredParams(rNameUpdated, "default"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessGenericExists(resourceName, &access_generic_updated),
					resource.TestCheckResourceAttr(resourceName, "attachable_access_entity_profile_dn", fmt.Sprintf("uni/infra/attentp-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", "default"),
					testAccCheckAciAccessGenericIdNotEqual(&access_generic_default, &access_generic_updated),
				),
			},
			{
				Config: CreateAccAccessGenericConfig(infraAttEntityPName, "default"),
			},
			{
				Config:      CreateAccAccessGenericConfigWithRequiredParams(rName, rNameUpdated),
				ExpectError: regexp.MustCompile(`default allowed`),
			},
		},
	})
}

func TestAccAciAccessGeneric_Negative(t *testing.T) {
	rName := "default"
	infraAttEntityPName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessGenericDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAccessGenericConfig(infraAttEntityPName, rName),
			},
			{
				Config:      CreateAccAccessGenericWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccAccessGenericUpdatedAttr(infraAttEntityPName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAccessGenericUpdatedAttr(infraAttEntityPName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAccessGenericUpdatedAttr(infraAttEntityPName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccAccessGenericUpdatedAttr(infraAttEntityPName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccAccessGenericConfig(infraAttEntityPName, rName),
			},
		},
	})
}

func testAccCheckAciAccessGenericExists(name string, access_generic *models.AccessGeneric) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Access Generic %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access Generic dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		access_genericFound := models.AccessGenericFromContainer(cont)
		if access_genericFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Access Generic %s not found", rs.Primary.ID)
		}
		*access_generic = *access_genericFound
		return nil
	}
}

func testAccCheckAciAccessGenericDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing access_generic destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_access_generic" {
			cont, err := client.Get(rs.Primary.ID)
			access_generic := models.AccessGenericFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Access Generic %s Still exists", access_generic.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAccessGenericIdEqual(m1, m2 *models.AccessGeneric) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("access_generic DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciAccessGenericIdNotEqual(m1, m2 *models.AccessGeneric) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("access_generic DNs are equal")
		}
		return nil
	}
}

func CreateAccessGenericWithoutRequired(infraAttEntityPName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing access_generic creation without ", attrName)
	rBlock := `
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "attachable_access_entity_profile_dn":
		rBlock += `
	resource "aci_access_generic" "test" {
	#	attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, infraAttEntityPName, rName)
}

func CreateAccAccessGenericConfigWithRequiredParams(infraAttEntityPName, rName string) string {
	fmt.Printf("=== STEP  testing access_generic creation with parent resource name %s and resource name %s\n", infraAttEntityPName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = "%s"
	}
	`, infraAttEntityPName, rName)
	return resource
}
func CreateAccAccessGenericConfigUpdatedName(infraAttEntityPName, rName string) string {
	fmt.Println("=== STEP  testing access_generic creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = "%s"
	}
	`, infraAttEntityPName, rName)
	return resource
}

func CreateAccAccessGenericConfig(infraAttEntityPName, rName string) string {
	fmt.Println("=== STEP  testing access_generic creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = "%s"
	}
	`, infraAttEntityPName, rName)
	return resource
}

func CreateAccAccessGenericWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing access_generic creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccAccessGenericConfigWithOptionalValues(infraAttEntityPName, rName string) string {
	fmt.Println("=== STEP  Basic: testing access_generic creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = "${aci_attachable_access_entity_profile.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_access_generic"
		
	}
	`, infraAttEntityPName, rName)

	return resource
}

func CreateAccAccessGenericRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing access_generic updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_access_generic" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_access_generic"
		
	}
	`)

	return resource
}

func CreateAccAccessGenericUpdatedAttr(infraAttEntityPName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing access_generic attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = "%s"
		%s = "%s"
	}
	`, infraAttEntityPName, rName, attribute, value)
	return resource
}
