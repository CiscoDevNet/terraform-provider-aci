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

func TestAccAciL3outRouteTagPolicy_Basic(t *testing.T) {
	var l3out_route_tag_policy_default models.L3outRouteTagPolicy
	var l3out_route_tag_policy_updated models.L3outRouteTagPolicy
	resourceName := "aci_l3out_route_tag_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outRouteTagPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outRouteTagPolicyWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outRouteTagPolicyWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outRouteTagPolicyConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outRouteTagPolicyExists(resourceName, &l3out_route_tag_policy_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "tag", "4294967295"),
				),
			},
			{

				Config: CreateAccL3outRouteTagPolicyConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outRouteTagPolicyExists(resourceName, &l3out_route_tag_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_route_tag_policy"),
					resource.TestCheckResourceAttr(resourceName, "tag", "6546738"),
					testAccCheckAciL3outRouteTagPolicyIdEqual(&l3out_route_tag_policy_default, &l3out_route_tag_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccL3outRouteTagPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccL3outRouteTagPolicyConfigUpdateWithInvalidName(rName, longrName),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config: CreateAccL3outRouteTagPolicyConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outRouteTagPolicyExists(resourceName, &l3out_route_tag_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciL3outRouteTagPolicyIdNotEqual(&l3out_route_tag_policy_default, &l3out_route_tag_policy_updated),
				),
			},

			{
				Config: CreateAccL3outRouteTagPolicyConfig(rName, rName),
			},
			{
				Config: CreateAccL3outRouteTagPolicyConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outRouteTagPolicyExists(resourceName, &l3out_route_tag_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciL3outRouteTagPolicyIdNotEqual(&l3out_route_tag_policy_default, &l3out_route_tag_policy_updated),
				),
			},
		},
	})
}

func TestAccAciL3outRouteTagPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outRouteTagPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outRouteTagPolicyConfig(rName, rName),
			},
			{
				Config:      CreateAccL3outRouteTagPolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class l3extRouteTagPol (.)+`),
			},
			{
				Config:      CreateAccL3outRouteTagPolicyUpdatedAttr(rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outRouteTagPolicyUpdatedAttr(rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outRouteTagPolicyUpdatedAttr(rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outRouteTagPolicyUpdatedAttr(rName, rName, "tag", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccL3outRouteTagPolicyUpdatedAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3outRouteTagPolicyConfig(rName, rName),
			},
		},
	})
}
func TestAccAciL3outRouteTagPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outRouteTagPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outRouteTagPolicyConfigMultiple(rName),
			},
		},
	})
}

func CreateAccL3outRouteTagPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  creating multiple L3outRouteTagPolicy")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3out_route_tag_policy" "test" {
			  tenant_dn      = aci_tenant.test.id
			  name           = "%s"
	  }
	  resource "aci_l3out_route_tag_policy" "test1" {
		tenant_dn      = aci_tenant.test.id
		name           = "%s"
	}
	resource "aci_l3out_route_tag_policy" "test2" {
		tenant_dn      = aci_tenant.test.id
		name           = "%s"
	}
	`, rName, rName+"1", rName+"2", rName+"3")
	return resource
}
func CreateAccL3outRouteTagPolicyConfigUpdateWithInvalidName(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing L3outRouteTagPolicy creation with parent resource name %s and name %s\n", parentName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name  = "%s"
	  }
	  resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, parentName, rName)
	return resource
}
func testAccCheckAciL3outRouteTagPolicyExists(name string, l3out_route_tag_policy *models.L3outRouteTagPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Route Tag Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Route Tag Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_route_tag_policyFound := models.L3outRouteTagPolicyFromContainer(cont)
		if l3out_route_tag_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Route Tag Policy %s not found", rs.Primary.ID)
		}
		*l3out_route_tag_policy = *l3out_route_tag_policyFound
		return nil
	}
}

func testAccCheckAciL3outRouteTagPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_route_tag_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_route_tag_policy" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_route_tag_policy := models.L3outRouteTagPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Route Tag Policy %s Still exists", l3out_route_tag_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outRouteTagPolicyIdEqual(m1, m2 *models.L3outRouteTagPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_route_tag_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outRouteTagPolicyIdNotEqual(m1, m2 *models.L3outRouteTagPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_route_tag_policy DNs are equal")
		}
		return nil
	}
}

func CreateL3outRouteTagPolicyWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_route_tag_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_l3out_route_tag_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}	`
	}

	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccL3outRouteTagPolicyConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing L3outRouteTagPolicy creation with parent resource name %s and name %s\n", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccL3outRouteTagPolicyConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing l3out_route_tag_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccL3outRouteTagPolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing l3out_route_tag_policy creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = "aci_application_profile.test.id"
		name  = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccL3outRouteTagPolicyConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_route_tag_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_route_tag_policy"
		tag = "6546738"
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccL3outRouteTagPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_route_tag_policy updation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_l3out_route_tag_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_route_tag_policy"
		tag = "686789"
	}
	`)

	return resource
}

func CreateAccL3outRouteTagPolicyUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_route_tag_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
