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

func TestAccAciBGPTimersPolicy_Basic(t *testing.T) {
	var bgp_timers_policy_default models.BGPTimersPolicy
	var bgp_timers_policy_updated models.BGPTimersPolicy
	resourceName := "aci_bgp_timers.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPTimersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBGPTimersPolicyWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBGPTimersPolicyWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBGPTimersPolicyConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPTimersPolicyExists(resourceName, &bgp_timers_policy_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "gr_ctrl", "helper"),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "180"),
					resource.TestCheckResourceAttr(resourceName, "ka_intvl", "60"),
					resource.TestCheckResourceAttr(resourceName, "max_as_limit", "0"),
				),
			},
			{
				Config: CreateAccBGPTimersPolicyConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPTimersPolicyExists(resourceName, &bgp_timers_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_bgp_timers_policy"),
					resource.TestCheckResourceAttr(resourceName, "gr_ctrl", "helper"),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "1800"),
					resource.TestCheckResourceAttr(resourceName, "ka_intvl", "1500"),
					resource.TestCheckResourceAttr(resourceName, "max_as_limit", "1"),
					resource.TestCheckResourceAttr(resourceName, "stale_intvl", "1"),
					testAccCheckAciBGPTimersPolicyIdEqual(&bgp_timers_policy_default, &bgp_timers_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccBGPTimersPolicyConfigWithRequiredParams(rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccBGPTimersPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBGPTimersPolicyConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPTimersPolicyExists(resourceName, &bgp_timers_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciBGPTimersPolicyIdNotEqual(&bgp_timers_policy_default, &bgp_timers_policy_updated),
				),
			},
			{
				Config: CreateAccBGPTimersPolicyConfig(rName, rName),
			},
			{
				Config: CreateAccBGPTimersPolicyConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPTimersPolicyExists(resourceName, &bgp_timers_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciBGPTimersPolicyIdNotEqual(&bgp_timers_policy_default, &bgp_timers_policy_updated),
				),
			},
		},
	})
}

func TestAccAciBGPTimersPolicy_Update(t *testing.T) {
	var bgp_timers_policy_default models.BGPTimersPolicy
	var bgp_timers_policy_updated models.BGPTimersPolicy
	resourceName := "aci_bgp_timers.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPTimersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBGPTimersPolicyConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPTimersPolicyExists(resourceName, &bgp_timers_policy_default),
				),
			},
			{
				Config: CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "max_as_limit", "2000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPTimersPolicyExists(resourceName, &bgp_timers_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "max_as_limit", "2000"),
				),
			},
			{
				Config: CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "max_as_limit", "1000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPTimersPolicyExists(resourceName, &bgp_timers_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "max_as_limit", "1000"),
				),
			},
			{
				Config: CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "stale_intvl", "1800"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPTimersPolicyExists(resourceName, &bgp_timers_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "stale_intvl", "1800"),
				),
			},
			{
				Config: CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "stale_intvl", "3600"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPTimersPolicyExists(resourceName, &bgp_timers_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "stale_intvl", "3600"),
				),
			},
		},
	})
}

func TestAccAciBGPTimersPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPTimersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBGPTimersPolicyConfig(rName, rName),
			},
			{
				Config:      CreateAccBGPTimersPolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "gr_ctrl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "hold_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "ka_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "max_as_limit", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "max_as_limit", "2001"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "stale_intvl", "0"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, "stale_intvl", "3601"),
				ExpectError: regexp.MustCompile(` Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccBGPTimersPolicyUpdatedAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccBGPTimersPolicyConfig(rName, rName),
			},
		},
	})
}

func TestAccAciBGPTimersPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBGPTimersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBGPTimersPolicyConfigs(rName),
			},
		},
	})
}

func testAccCheckAciBGPTimersPolicyExists(name string, bgp_timers_policy *models.BGPTimersPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("BGP Timers Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BGP Timers Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bgp_timers_policyFound := models.BGPTimersPolicyFromContainer(cont)
		if bgp_timers_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("BGP Timers Policy %s not found", rs.Primary.ID)
		}
		*bgp_timers_policy = *bgp_timers_policyFound
		return nil
	}
}

func testAccCheckAciBGPTimersPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing bgp_timers_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bgp_timers" {
			cont, err := client.Get(rs.Primary.ID)
			bgp_timers_policy := models.BGPTimersPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("BGP Timers Policy %s Still exists", bgp_timers_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBGPTimersPolicyIdEqual(m1, m2 *models.BGPTimersPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("bgp_timers_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciBGPTimersPolicyIdNotEqual(m1, m2 *models.BGPTimersPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("bgp_timers_policy DNs are equal")
		}
		return nil
	}
}

func CreateBGPTimersPolicyWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_timers_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_bgp_timers" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBGPTimersPolicyConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing bgp_timers_policy creation with parent resource name %s and resource name %s\n", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBGPTimersPolicyConfigs(rName string) string {
	fmt.Println("=== STEP  testing multiple bgp_timers_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_bgp_timers" "test1" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_bgp_timers" "test2" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, rName, rName, rName+"1", rName+"2")
	return resource
}

func CreateAccBGPTimersPolicyConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_timers_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBGPTimersPolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing bgp_timers_policy creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		name = "%s"
		tenant_dn  = aci_tenant.test.id
	}
	resource "aci_bgp_timers" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccBGPTimersPolicyConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_timers_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_timers" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_bgp_timers_policy"
		gr_ctrl = "helper"
		hold_intvl = "1800"
		ka_intvl = "1500"
		max_as_limit = "1"
		stale_intvl = "1"
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccBGPTimersPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing bgp_timers_policy update without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_bgp_timers" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_bgp_timers_policy"
		gr_ctrl = "helper"
		hold_intvl = "1"
		ka_intvl = "1"
		max_as_limit = "1"
		stale_intvl = "2"
	}
	`)

	return resource
}

func CreateAccBGPTimersPolicyUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing bgp_timers_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
