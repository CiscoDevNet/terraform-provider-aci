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
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciCoopPolicy_Basic(t *testing.T) {
	var coop_policy_default models.COOPGroupPolicy
	var coop_policy_updated models.COOPGroupPolicy
	resourceName := "aci_coop_policy.test"
	coopPolicy, err := aci.GetRemoteCOOPGroupPolicy(sharedAciClient(), "uni/fabric/pol-default")
	if err != nil {
		t.Errorf("reading initial config of coopPolicy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCoopPolicyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCoopPolicyExists(resourceName, &coop_policy_default),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
			{
				Config: CreateAccCoopPolicyConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCoopPolicyExists(resourceName, &coop_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_coop_policy"),
					resource.TestCheckResourceAttr(resourceName, "type", "compatible"),
					testAccCheckAciCoopPolicyIdEqual(&coop_policy_default, &coop_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccCoopPolicyConfigWithUpdatedType(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCoopPolicyExists(resourceName, &coop_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "type", "strict"),
					testAccCheckAciCoopPolicyIdEqual(&coop_policy_default, &coop_policy_updated),
				),
			},
			{
				Config: CreateAccCoopPolicyInitialConfig(coopPolicy),
			},
		},
	})
}

func TestAccAciCoopPolicy_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	coopPolicy, err := aci.GetRemoteCOOPGroupPolicy(sharedAciClient(), "uni/fabric/pol-default")
	if err != nil {
		t.Errorf("reading initial config of coopPolicy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCoopPolicyConfig(),
			},
			{
				Config:      CreateAccCoopPolicyUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCoopPolicyUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCoopPolicyUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccCoopPolicyUpdatedAttr("type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccCoopPolicyUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCoopPolicyInitialConfig(coopPolicy),
			},
		},
	})
}

func testAccCheckAciCoopPolicyExists(name string, coop_policy *models.COOPGroupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Coop Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Coop Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		coop_policyFound := models.COOPGroupPolicyFromContainer(cont)
		if coop_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Coop Policy %s not found", rs.Primary.ID)
		}
		*coop_policy = *coop_policyFound
		return nil
	}
}

func testAccCheckAciCoopPolicyIdEqual(m1, m2 *models.COOPGroupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("coop_policy DNs are not equal")
		}
		return nil
	}
}

func CreateAccCoopPolicyConfig() string {
	fmt.Println("=== STEP  testing coop_policy creation")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	}
	`)
	return resource
}

func CreateAccCoopPolicyConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing coop_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_coop_policy"
		type = "compatible"
		
	}
	`)

	return resource
}

func CreateAccCoopPolicyConfigWithUpdatedType() string {
	fmt.Println("=== STEP  Basic: testing coop_policy creation with Updated Type")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
		type = "strict"	
	}
	`)

	return resource
}

func CreateAccCoopPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing coop_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_coop_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_coop_policy"
		type = "strict"
	}
	`)

	return resource
}

func CreateAccCoopPolicyUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing coop_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
		%s = "%s"
	}
	`, attribute, value)
	return resource
}

func CreateAccCoopPolicyInitialConfig(coopPolicy *models.COOPGroupPolicy) string {
	fmt.Println("=== STEP  testing coop_policy with initial config")
	resource := fmt.Sprintf(`

	resource "aci_coop_policy" "test"{
		annotation = "%s"
		name_alias = "%s"
		description = "%s"
		type = "%s"

	}
	`, coopPolicy.Annotation, coopPolicy.NameAlias, coopPolicy.Description, coopPolicy.Type)
	return resource
}
