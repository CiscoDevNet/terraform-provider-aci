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

func TestAccAciBGPAddressFamilyContextPolicy_Basic(t *testing.T) {
	var bgp_address_family_context_default models.BGPAddressFamilyContextPolicy
	var bgp_address_family_context_updated models.BGPAddressFamilyContextPolicy
	resourceName := "aci_bgp_address_family_context.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPAddressFamilyContextPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBGPAddressFamilyContextPolicyWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBGPAddressFamilyContextPolicyWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "ctrl", ""),
					resource.TestCheckResourceAttr(resourceName, "e_dist", "20"),
					resource.TestCheckResourceAttr(resourceName, "i_dist", "200"),
					resource.TestCheckResourceAttr(resourceName, "local_dist", "220"),
					resource.TestCheckResourceAttr(resourceName, "max_ecmp", "16"),
					resource.TestCheckResourceAttr(resourceName, "max_ecmp_ibgp", "16"),
				),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_bgp_address_family_context"),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "host-rt-leak"),
					resource.TestCheckResourceAttr(resourceName, "e_dist", "1"),
					resource.TestCheckResourceAttr(resourceName, "i_dist", "1"),
					resource.TestCheckResourceAttr(resourceName, "local_dist", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_ecmp", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_ecmp_ibgp", "1"),
					testAccCheckAciBGPAddressFamilyContextPolicyIdEqual(&bgp_address_family_context_default, &bgp_address_family_context_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyConfigWithRequiredParams(rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccBGPAddressFamilyContextPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciBGPAddressFamilyContextPolicyIdNotEqual(&bgp_address_family_context_default, &bgp_address_family_context_updated),
				),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyConfig(rName, rName),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciBGPAddressFamilyContextPolicyIdNotEqual(&bgp_address_family_context_default, &bgp_address_family_context_updated),
				),
			},
		},
	})
}

func TestAccAciBGPAddressFamilyContextPolicy_Update(t *testing.T) {
	var bgp_address_family_context_default models.BGPAddressFamilyContextPolicy
	var bgp_address_family_context_updated models.BGPAddressFamilyContextPolicy
	resourceName := "aci_bgp_address_family_context.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPAddressFamilyContextPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBGPAddressFamilyContextPolicyConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_default),
				),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyUpdated(rName, rName, "e_dist", "255"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "e_dist", "255"),
				),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyUpdated(rName, rName, "e_dist", "125"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "e_dist", "125"),
				),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyUpdated(rName, rName, "i_dist", "255"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "i_dist", "255"),
				),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyUpdated(rName, rName, "i_dist", "125"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "i_dist", "125"),
				),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyUpdated(rName, rName, "local_dist", "255"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "local_dist", "255"),
				),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyUpdated(rName, rName, "local_dist", "125"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "local_dist", "125"),
				),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyUpdated(rName, rName, "max_ecmp", "64"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "max_ecmp", "64"),
				),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyUpdated(rName, rName, "max_ecmp", "32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "max_ecmp", "32"),
				),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyUpdated(rName, rName, "max_ecmp_ibgp", "64"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "max_ecmp_ibgp", "64"),
				),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyUpdated(rName, rName, "max_ecmp_ibgp", "32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists(resourceName, &bgp_address_family_context_updated),
					resource.TestCheckResourceAttr(resourceName, "max_ecmp_ibgp", "32"),
				),
			},
		},
	})
}

func TestAccAciBGPAddressFamilyContextPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPAddressFamilyContextPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBGPAddressFamilyContextPolicyConfig(rName, rName),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "ctrl", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "e_dist", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "e_dist", "256"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "i_dist", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "i_dist", "256"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "local_dist", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "local_dist", "256"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "max_ecmp", "0"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "max_ecmp", "65"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "max_ecmp_ibgp", "0"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, "max_ecmp_ibgp", "65"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccBGPAddressFamilyContextPolicyConfig(rName, rName),
			},
		},
	})
}

func TestAccAciBGPAddressFamilyContextPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPAddressFamilyContextPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBGPAddressFamilyContextPolicyConfigs(rName),
			},
		},
	})
}

func testAccCheckAciBGPAddressFamilyContextPolicyExists(name string, bgp_address_family_context *models.BGPAddressFamilyContextPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("BGP Address Family Context %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BGP Address Family Context dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bgp_address_family_contextFound := models.BGPAddressFamilyContextPolicyFromContainer(cont)
		if bgp_address_family_contextFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("BGP Address Family Context %s not found", rs.Primary.ID)
		}
		*bgp_address_family_context = *bgp_address_family_contextFound
		return nil
	}
}

func testAccCheckAciBGPAddressFamilyContextPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing bgp_address_family_context destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bgp_address_family_context" {
			cont, err := client.Get(rs.Primary.ID)
			bgp_address_family_context := models.BGPAddressFamilyContextPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("BGP Address Family Context %s Still exists", bgp_address_family_context.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBGPAddressFamilyContextPolicyIdEqual(m1, m2 *models.BGPAddressFamilyContextPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("bgp_address_family_context DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciBGPAddressFamilyContextPolicyIdNotEqual(m1, m2 *models.BGPAddressFamilyContextPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("bgp_address_family_context DNs are equal")
		}
		return nil
	}
}

func CreateBGPAddressFamilyContextPolicyWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_address_family_context creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_bgp_address_family_context" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBGPAddressFamilyContextPolicyConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing bgp_address_family_context creation with parent resource name %s and name %s\n", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBGPAddressFamilyContextPolicyConfigs(rName string) string {
	fmt.Println("=== STEP  testing multiple bgp_address_family_context creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_bgp_address_family_context" "test1" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_bgp_address_family_context" "test2" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, rName, rName, rName+"1", rName+"2")
	return resource
}

func CreateAccBGPAddressFamilyContextPolicyConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_address_family_context creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBGPAddressFamilyContextPolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing bgp_address_family_context creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		name = "%s"
		tenant_dn  = aci_tenant.test.id
	}

	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccBGPAddressFamilyContextPolicyConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_address_family_context creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_bgp_address_family_context"
		ctrl = "host-rt-leak"
		e_dist = "1"
		i_dist = "1"
		local_dist = "1"
		max_ecmp = "1"
		max_ecmp_ibgp = "1"
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccBGPAddressFamilyContextPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing bgp_address_family_context updadation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_bgp_address_family_context" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_bgp_address_family_context"
		ctrl = "host-rt-leak"
		e_dist = "2"
		i_dist = "2"
		local_dist = "2"
		max_ecmp = "2"
		max_ecmp_ibgp = "2"
	}
	`)

	return resource
}

func CreateAccBGPAddressFamilyContextPolicyUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing bgp_address_family_context attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}

func CreateAccBGPAddressFamilyContextPolicyUpdated(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing bgp_address_family_context attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_address_family_context" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
