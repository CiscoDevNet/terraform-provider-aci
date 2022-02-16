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

func TestAccAciDHCPOptionPolicy_Basic(t *testing.T) {
	var dhcp_option_policy_default models.DHCPOptionPolicy
	var dhcp_option_policy_updated models.DHCPOptionPolicy
	resourceName := "aci_dhcp_option_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDHCPOptionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateDHCPOptionPolicyWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateDHCPOptionPolicyWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDHCPOptionPolicyConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPOptionPolicyExists(resourceName, &dhcp_option_policy_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccDHCPOptionPolicyConfigWithOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPOptionPolicyExists(resourceName, &dhcp_option_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_dhcp_option_policy"),

					testAccCheckAciDHCPOptionPolicyIdEqual(&dhcp_option_policy_default, &dhcp_option_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccDHCPOptionPolicyConfigWithAllOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPOptionPolicyExists(resourceName, &dhcp_option_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_dhcp_option_policy"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option.0.annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option.0.name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option.0.data", ""),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option.0.dhcp_option_id", "0"),
					testAccCheckAciDHCPOptionPolicyIdEqual(&dhcp_option_policy_default, &dhcp_option_policy_updated),
				),
			},
			{
				Config: CreateAccDHCPOptionPolicyConfigWithAllOptionalValuesofOption(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPOptionPolicyExists(resourceName, &dhcp_option_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_dhcp_option_policy"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option.0.annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option.0.name_alias", "test_dhcp_option"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option.0.data", "test_data"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_option.0.dhcp_option_id", "1"),
					testAccCheckAciDHCPOptionPolicyIdEqual(&dhcp_option_policy_default, &dhcp_option_policy_updated),
				),
			},
			{
				Config:      CreateAccDHCPOptionPolicyConfigUpdatedName(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},
			{
				Config:      CreateAccDHCPOptionPolicyConfigWithInvalidOptionName(fvTenantName, rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccDHCPOptionPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccDHCPOptionPolicyOptionWithoutName(fvTenantName, rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDHCPOptionPolicyConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPOptionPolicyExists(resourceName, &dhcp_option_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciDHCPOptionPolicyIdNotEqual(&dhcp_option_policy_default, &dhcp_option_policy_updated),
				),
			},
			{
				Config: CreateAccDHCPOptionPolicyConfig(fvTenantName, rName),
			},
			{
				Config: CreateAccDHCPOptionPolicyConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPOptionPolicyExists(resourceName, &dhcp_option_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciDHCPOptionPolicyIdNotEqual(&dhcp_option_policy_default, &dhcp_option_policy_updated),
				),
			},
		},
	})
}

func TestAccAciDHCPOptionPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDHCPOptionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDHCPOptionPolicyConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccDHCPOptionPolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccDHCPOptionPolicyUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDHCPOptionPolicyUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDHCPOptionPolicyUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccDHCPOptionPolicyUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccDHCPOptionPolicyOptionUpdatedAttr(fvTenantName, rName, "data", acctest.RandString(257)),
				ExpectError: regexp.MustCompile(`failed validation for value`),
			},
			{
				Config:      CreateAccDHCPOptionPolicyOptionUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDHCPOptionPolicyOptionUpdatedAttr(fvTenantName, rName, "dhcp_option_id", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccDHCPOptionPolicyOptionUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDHCPOptionPolicyOptionUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},
			{
				Config: CreateAccDHCPOptionPolicyConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciDHCPOptionPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDHCPOptionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDHCPOptionPolicyConfigMultiple(fvTenantName, rName),
			},
		},
	})
}

func testAccCheckAciDHCPOptionPolicyExists(name string, dhcp_option_policy *models.DHCPOptionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("DHCP Option Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DHCP Option Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		dhcp_option_policyFound := models.DHCPOptionPolicyFromContainer(cont)
		if dhcp_option_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("DHCP Option Policy %s not found", rs.Primary.ID)
		}
		*dhcp_option_policy = *dhcp_option_policyFound
		return nil
	}
}

func testAccCheckAciDHCPOptionPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing dhcp_option_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_dhcp_option_policy" {
			cont, err := client.Get(rs.Primary.ID)
			dhcp_option_policy := models.DHCPOptionPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("DHCP Option Policy %s Still exists", dhcp_option_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciDHCPOptionPolicyIdEqual(m1, m2 *models.DHCPOptionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("dhcp_option_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciDHCPOptionPolicyIdNotEqual(m1, m2 *models.DHCPOptionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("dhcp_option_policy DNs are equal")
		}
		return nil
	}
}

func CreateDHCPOptionPolicyWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing dhcp_option_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_dhcp_option_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccDHCPOptionPolicyConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing dhcp_option_policy creation with parent resource name %s and resource name %s\n", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccDHCPOptionPolicyConfigWithInvalidOptionName(fvTenantName, rName, longName string) string {
	fmt.Println("=== STEP  testing dhcp_option creation with invalid name = ", longName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		dhcp_option {
			name = "%s"
		}
	}
	`, fvTenantName, rName, longName)
	return resource
}
func CreateAccDHCPOptionPolicyConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_option_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccDHCPOptionPolicyOptionWithoutName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_option creation without name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		dhcp_option {}
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccDHCPOptionPolicyConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_option_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccDHCPOptionPolicyConfigMultiple(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing multiple dhcp_option_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccDHCPOptionPolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing dhcp_option_policy creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccDHCPOptionPolicyConfigWithAllOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing dhcp_option_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_dhcp_option_policy"
		dhcp_option {
			name  = "%s"
		}
	}
	`, fvTenantName, rName, rName)

	return resource
}

func CreateAccDHCPOptionPolicyConfigWithAllOptionalValuesofOption(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing dhcp_option_policy creation with optional parameters of dhcp_option_policy and dhcp_option")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_dhcp_option_policy"
		dhcp_option {
			name  = "%s"
			annotation = "orchestrator:terraform_testacc"
			name_alias = "test_dhcp_option"
			data = "test_data"
			dhcp_option_id = "1"
		}
	}
	`, fvTenantName, rName, rName)

	return resource
}

func CreateAccDHCPOptionPolicyConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing dhcp_option_policy creation with optional parameters of dhcp_option_policy")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_dhcp_option_policy"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccDHCPOptionPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing dhcp_option_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_dhcp_option_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_dhcp_option_policy"
		
	}
	`)

	return resource
}

func CreateAccDHCPOptionPolicyOptionUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing dhcp_option attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		dhcp_option {
			name = "%s"
			%s = "%s"
		}
	}
	`, fvTenantName, rName, rName, attribute, value)
	return resource
}

func CreateAccDHCPOptionPolicyUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing dhcp_option_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}

func CreateAccDHCPOptionPolicyUpdatedAttrList(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing dhcp_option_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_option_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = %s
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
