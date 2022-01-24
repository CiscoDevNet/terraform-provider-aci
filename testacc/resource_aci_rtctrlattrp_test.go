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

func TestAccAciActionRuleProfile_Basic(t *testing.T) {
	var action_rule_profile_default models.ActionRuleProfile
	var action_rule_profile_updated models.ActionRuleProfile
	resourceName := "aci_action_rule_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciActionRuleProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateActionRuleProfileWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateActionRuleProfileWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccActionRuleProfileConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists(resourceName, &action_rule_profile_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccActionRuleProfileConfigWithOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists(resourceName, &action_rule_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_action_rule_profile"),

					testAccCheckAciActionRuleProfileIdEqual(&action_rule_profile_default, &action_rule_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccActionRuleProfileConfigUpdatedName(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccActionRuleProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccActionRuleProfileConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists(resourceName, &action_rule_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciActionRuleProfileIdNotEqual(&action_rule_profile_default, &action_rule_profile_updated),
				),
			},
			{
				Config: CreateAccActionRuleProfileConfig(fvTenantName, rName),
			},
			{
				Config: CreateAccActionRuleProfileConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists(resourceName, &action_rule_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciActionRuleProfileIdNotEqual(&action_rule_profile_default, &action_rule_profile_updated),
				),
			},
		},
	})
}

func TestAccAciActionRuleProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciActionRuleProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccActionRuleProfileConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccActionRuleProfileWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccActionRuleProfileUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccActionRuleProfileUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccActionRuleProfileUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccActionRuleProfileUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccActionRuleProfileConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciActionRuleProfile_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciActionRuleProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccActionRuleProfileConfigMultiple(fvTenantName, rName),
			},
		},
	})
}

func testAccCheckAciActionRuleProfileExists(name string, action_rule_profile *models.ActionRuleProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Action Rule Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Action Rule Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		action_rule_profileFound := models.ActionRuleProfileFromContainer(cont)
		if action_rule_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Action Rule Profile %s not found", rs.Primary.ID)
		}
		*action_rule_profile = *action_rule_profileFound
		return nil
	}
}

func testAccCheckAciActionRuleProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing action_rule_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_action_rule_profile" {
			cont, err := client.Get(rs.Primary.ID)
			action_rule_profile := models.ActionRuleProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Action Rule Profile %s Still exists", action_rule_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciActionRuleProfileIdEqual(m1, m2 *models.ActionRuleProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("action_rule_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciActionRuleProfileIdNotEqual(m1, m2 *models.ActionRuleProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("action_rule_profile DNs are equal")
		}
		return nil
	}
}

func CreateActionRuleProfileWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing action_rule_profile creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_action_rule_profile" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccActionRuleProfileConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing action_rule_profile creation with parent resource name %s and resource name %s", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}
func CreateAccActionRuleProfileConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing action_rule_profile creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccActionRuleProfileConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing action_rule_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccActionRuleProfileConfigMultiple(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing multiple action_rule_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccActionRuleProfileWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing action_rule_profile creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		tenant_dn  = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccActionRuleProfileConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing action_rule_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_action_rule_profile"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccActionRuleProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing action_rule_profile updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_action_rule_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_action_rule_profile"
		
	}
	`)

	return resource
}

func CreateAccActionRuleProfileUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing action_rule_profile attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_action_rule_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
