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

func TestAccAciVSwitchPolicy_Basic(t *testing.T) {
	var vswitch_policy_default models.VSwitchPolicyGroup
	var vswitch_policy_updated models.VSwitchPolicyGroup
	resourceName := "aci_vswitch_policy.test"
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	vmmDomPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVSwitchPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVSwitchPolicyWithoutRequired(vmmDomPName, "vmm_domain_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVSwitchPolicyConfig(vmmDomPName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVSwitchPolicyExists(resourceName, &vswitch_policy_default),
					resource.TestCheckResourceAttr(resourceName, "vmm_domain_dn", fmt.Sprintf("%v/dom-%s", providerProfileDn, vmmDomPName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccVSwitchPolicyConfigWithOptionalValues(vmmDomPName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVSwitchPolicyExists(resourceName, &vswitch_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "vmm_domain_dn", fmt.Sprintf("%v/dom-%s", providerProfileDn, vmmDomPName)), resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_vswitch_policy"),
					testAccCheckAciVSwitchPolicyIdEqual(&vswitch_policy_default, &vswitch_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccVSwitchPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVSwitchPolicyConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVSwitchPolicyExists(resourceName, &vswitch_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "vmm_domain_dn", fmt.Sprintf("%v/dom-%s", providerProfileDn, rNameUpdated)),
					testAccCheckAciVSwitchPolicyIdNotEqual(&vswitch_policy_default, &vswitch_policy_updated),
				),
			},
			{
				Config: CreateAccVSwitchPolicyConfig(vmmDomPName),
			},
		},
	})
}

func TestAccAciVSwitchPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	vmmDomPName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVSwitchPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVSwitchPolicyConfig(vmmDomPName),
			},
			{
				Config:      CreateAccVSwitchPolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVSwitchPolicyUpdatedAttr(vmmDomPName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVSwitchPolicyUpdatedAttr(vmmDomPName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVSwitchPolicyUpdatedAttr(vmmDomPName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVSwitchPolicyUpdatedAttr(vmmDomPName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccVSwitchPolicyConfig(vmmDomPName),
			},
		},
	})
}

func testAccCheckAciVSwitchPolicyExists(name string, vswitch_policy *models.VSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VSwitch Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VSwitch Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vswitch_policyFound := models.VSwitchPolicyGroupFromContainer(cont)
		if vswitch_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VSwitch Policy %s not found", rs.Primary.ID)
		}
		*vswitch_policy = *vswitch_policyFound
		return nil
	}
}

func testAccCheckAciVSwitchPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vswitch_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vswitch_policy" {
			cont, err := client.Get(rs.Primary.ID)
			vswitch_policy := models.VSwitchPolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VSwitch Policy %s Still exists", vswitch_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVSwitchPolicyIdEqual(m1, m2 *models.VSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vswitch_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVSwitchPolicyIdNotEqual(m1, m2 *models.VSwitchPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vswitch_policy DNs are equal")
		}
		return nil
	}
}

func CreateVSwitchPolicyWithoutRequired(vmmDomPName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vswitch_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	`
	switch attrName {
	case "vmm_domain_dn":
		rBlock += `
	resource "aci_vswitch_policy" "test" {
	#	vmm_domain_dn  = aci_vmm_domain.test.id
	}
		`

	}
	return fmt.Sprintf(rBlock, vmmDomPName, providerProfileDn)
}

func CreateAccVSwitchPolicyConfigWithRequiredParams(vmmDomPName string) string {
	fmt.Println("=== STEP  testing vswitch_policy creation with Updated required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vswitch_policy" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
	}
	`, vmmDomPName, providerProfileDn)
	return resource
}

func CreateAccVSwitchPolicyConfig(vmmDomPName string) string {
	fmt.Println("=== STEP  testing vswitch_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	resource "aci_vswitch_policy" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
	}
	`, vmmDomPName, providerProfileDn)
	return resource
}

func CreateAccVSwitchPolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing vswitch_policy creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_vswitch_policy" "test" {
		vmm_domain_dn  = aci_tenant.test.id	
	}
	`, rName)
	return resource
}

func CreateAccVSwitchPolicyConfigWithOptionalValues(vmmDomPName string) string {
	fmt.Println("=== STEP  Basic: testing vswitch_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vswitch_policy" "test" {
		vmm_domain_dn  = "${aci_vmm_domain.test.id}"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vswitch_policy"
		
	}
	`, vmmDomPName, providerProfileDn)

	return resource
}

func CreateAccVSwitchPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vswitch_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_vswitch_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vswitch_policy"	
	}
	`)

	return resource
}

func CreateAccVSwitchPolicyUpdatedAttr(vmmDomPName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vswitch_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vswitch_policy" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		%s = "%s"
	}
	`, vmmDomPName, providerProfileDn, attribute, value)
	return resource
}
