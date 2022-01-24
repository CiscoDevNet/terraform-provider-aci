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

func TestAccAciDHCPRelayPolicy_Basic(t *testing.T) {
	var dhcp_relay_policy_default models.DHCPRelayPolicy
	var dhcp_relay_policy_updated models.DHCPRelayPolicy
	resourceName := "aci_dhcp_relay_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDHCPRelayPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateDHCPRelayPolicyWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateDHCPRelayPolicyWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDHCPRelayPolicyConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPRelayPolicyExists(resourceName, &dhcp_relay_policy_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "mode", "visible"),
					resource.TestCheckResourceAttr(resourceName, "owner", "infra"),
				),
			},
			{
				Config: CreateAccDHCPRelayPolicyConfigWithOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPRelayPolicyExists(resourceName, &dhcp_relay_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_dhcp_relay_policy"),

					resource.TestCheckResourceAttr(resourceName, "mode", "visible"),

					resource.TestCheckResourceAttr(resourceName, "owner", "tenant"),

					testAccCheckAciDHCPRelayPolicyIdEqual(&dhcp_relay_policy_default, &dhcp_relay_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccDHCPRelayPolicyConfigUpdatedName(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccDHCPRelayPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccDHCPRelayPolicyConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPRelayPolicyExists(resourceName, &dhcp_relay_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciDHCPRelayPolicyIdNotEqual(&dhcp_relay_policy_default, &dhcp_relay_policy_updated),
				),
			},
			{
				Config: CreateAccDHCPRelayPolicyConfig(fvTenantName, rName),
			},
			{
				Config: CreateAccDHCPRelayPolicyConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPRelayPolicyExists(resourceName, &dhcp_relay_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciDHCPRelayPolicyIdNotEqual(&dhcp_relay_policy_default, &dhcp_relay_policy_updated),
				),
			},
		},
	})
}

func TestAccAciDHCPRelayPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDHCPRelayPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDHCPRelayPolicyConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccDHCPRelayPolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccDHCPRelayPolicyUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDHCPRelayPolicyUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccDHCPRelayPolicyUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccDHCPRelayPolicyUpdatedAttr(fvTenantName, rName, "mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccDHCPRelayPolicyUpdatedAttr(fvTenantName, rName, "owner", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccDHCPRelayPolicyUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccDHCPRelayPolicyConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciDHCPRelayPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciDHCPRelayPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccDHCPRelayPolicyConfigMultiple(fvTenantName, rName),
			},
		},
	})
}

func testAccCheckAciDHCPRelayPolicyExists(name string, dhcp_relay_policy *models.DHCPRelayPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("DHCP Relay Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DHCP Relay Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		dhcp_relay_policyFound := models.DHCPRelayPolicyFromContainer(cont)
		if dhcp_relay_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("DHCP Relay Policy %s not found", rs.Primary.ID)
		}
		*dhcp_relay_policy = *dhcp_relay_policyFound
		return nil
	}
}

func testAccCheckAciDHCPRelayPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing dhcp_relay_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_dhcp_relay_policy" {
			cont, err := client.Get(rs.Primary.ID)
			dhcp_relay_policy := models.DHCPRelayPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("DHCP Relay Policy %s Still exists", dhcp_relay_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciDHCPRelayPolicyIdEqual(m1, m2 *models.DHCPRelayPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("dhcp_relay_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciDHCPRelayPolicyIdNotEqual(m1, m2 *models.DHCPRelayPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("dhcp_relay_policy DNs are equal")
		}
		return nil
	}
}

func CreateDHCPRelayPolicyWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing dhcp_relay_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_dhcp_relay_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccDHCPRelayPolicyConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_relay_policy creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}
func CreateAccDHCPRelayPolicyConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_relay_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccDHCPRelayPolicyConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing dhcp_relay_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccDHCPRelayPolicyConfigMultiple(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing multiple dhcp_relay_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccDHCPRelayPolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing dhcp_relay_policy creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		tenant_dn  = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccDHCPRelayPolicyConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing dhcp_relay_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_dhcp_relay_policy"
		mode = "visible"
		owner = "tenant"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccDHCPRelayPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing dhcp_relay_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_dhcp_relay_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_dhcp_relay_policy"
		mode = "not-visible"
		owner = "tenant"
		
	}
	`)

	return resource
}

func CreateAccDHCPRelayPolicyUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing dhcp_relay_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_dhcp_relay_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
