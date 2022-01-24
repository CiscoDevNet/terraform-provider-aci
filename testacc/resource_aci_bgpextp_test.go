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

func TestAccAciL3outBGPExternalPolicy_Basic(t *testing.T) {
	var l3out_bgp_external_policy_default models.L3outBgpExternalPolicy
	var l3out_bgp_external_policy_updated models.L3outBgpExternalPolicy
	resourceName := "aci_l3out_bgp_external_policy.test"
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	l3extOutName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outBGPExternalPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outBGPExternalPolicyWithoutRequired(fvTenantName, l3extOutName, "l3_outside_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outBGPExternalPolicyConfig(fvTenantName, l3extOutName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPExternalPolicyExists(resourceName, &l3out_bgp_external_policy_default),
					resource.TestCheckResourceAttr(resourceName, "l3_outside_dn", fmt.Sprintf("uni/tn-%s/out-%s", fvTenantName, l3extOutName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccL3outBGPExternalPolicyConfigWithOptionalValues(fvTenantName, l3extOutName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPExternalPolicyExists(resourceName, &l3out_bgp_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "l3_outside_dn", fmt.Sprintf("uni/tn-%s/out-%s", fvTenantName, l3extOutName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_bgp_external_policy"),
					testAccCheckAciL3outBGPExternalPolicyIdEqual(&l3out_bgp_external_policy_default, &l3out_bgp_external_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccL3outBGPExternalPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outBGPExternalPolicyConfigWithRequiredParams(fvTenantName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPExternalPolicyExists(resourceName, &l3out_bgp_external_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "l3_outside_dn", fmt.Sprintf("uni/tn-%s/out-%s", fvTenantName, rNameUpdated)),
					testAccCheckAciL3outBGPExternalPolicyIdNotEqual(&l3out_bgp_external_policy_default, &l3out_bgp_external_policy_updated),
				),
			},
			{
				Config: CreateAccL3outBGPExternalPolicyConfig(fvTenantName, l3extOutName),
			},
		},
	})
}

func TestAccAciL3outBGPExternalPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	l3extOutName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outBGPExternalPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outBGPExternalPolicyConfig(fvTenantName, l3extOutName),
			},
			{
				Config:      CreateAccL3outBGPExternalPolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outBGPExternalPolicyUpdatedAttr(fvTenantName, l3extOutName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outBGPExternalPolicyUpdatedAttr(fvTenantName, l3extOutName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outBGPExternalPolicyUpdatedAttr(fvTenantName, l3extOutName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outBGPExternalPolicyUpdatedAttr(fvTenantName, l3extOutName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3outBGPExternalPolicyConfig(fvTenantName, l3extOutName),
			},
		},
	})
}

func testAccCheckAciL3outBGPExternalPolicyExists(name string, l3out_bgp_external_policy *models.L3outBgpExternalPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out BGP External Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out BGP External Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_bgp_external_policyFound := models.L3outBgpExternalPolicyFromContainer(cont)
		if l3out_bgp_external_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out BGP External Policy %s not found", rs.Primary.ID)
		}
		*l3out_bgp_external_policy = *l3out_bgp_external_policyFound
		return nil
	}
}

func testAccCheckAciL3outBGPExternalPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_bgp_external_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_bgp_external_policy" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_bgp_external_policy := models.L3outBgpExternalPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out BGP External Policy %s Still exists", l3out_bgp_external_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outBGPExternalPolicyIdEqual(m1, m2 *models.L3outBgpExternalPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_bgp_external_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outBGPExternalPolicyIdNotEqual(m1, m2 *models.L3outBgpExternalPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_bgp_external_policy DNs are equal")
		}
		return nil
	}
}

func CreateL3outBGPExternalPolicyWithoutRequired(fvTenantName, l3extOutName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_bgp_external_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	`
	switch attrName {
	case "l3_outside_dn":
		rBlock += `
	resource "aci_l3out_bgp_external_policy" "test" {
	#	l3_outside_dn  = aci_l3_outside.test.id
	}
	`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName)
}

func CreateAccL3outBGPExternalPolicyConfigWithRequiredParams(fvTenantName, l3extOutName string) string {
	fmt.Println("=== STEP  testing l3out_bgp_external_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}
	`, fvTenantName, l3extOutName)
	return resource
}

func CreateAccL3outBGPExternalPolicyConfig(fvTenantName, l3extOutName string) string {
	fmt.Println("=== STEP  testing l3out_bgp_external_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
	}
	`, fvTenantName, l3extOutName)
	return resource
}

func CreateAccL3outBGPExternalPolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing l3out_bgp_external_policy creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = aci_tenant.test.id
	}
	`, rName)
	return resource
}

func CreateAccL3outBGPExternalPolicyConfigWithOptionalValues(fvTenantName, l3extOutName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_bgp_external_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = "${aci_l3_outside.test.id}"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_bgp_external_policy"
	}
	`, fvTenantName, l3extOutName)

	return resource
}

func CreateAccL3outBGPExternalPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_bgp_external_policy creation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_l3out_bgp_external_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_bgp_external_policy"
	}
	`)
	return resource
}

func CreateAccL3outBGPExternalPolicyUpdatedAttr(fvTenantName, l3extOutName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_bgp_external_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_l3out_bgp_external_policy" "test" {
		l3_outside_dn  = aci_l3_outside.test.id
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, attribute, value)
	return resource
}
