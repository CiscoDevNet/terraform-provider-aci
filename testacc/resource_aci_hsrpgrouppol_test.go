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

func TestAccAciHSRPGroupPolicy_Basic(t *testing.T) {
	var hsrp_group_policy_default models.HSRPGroupPolicy
	var hsrp_group_policy_updated models.HSRPGroupPolicy
	resourceName := "aci_hsrp_group_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciHSRPGroupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateHSRPGroupPolicyWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateHSRPGroupPolicyWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccHSRPGroupPolicyConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists(resourceName, &hsrp_group_policy_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "0"),
					resource.TestCheckResourceAttr(resourceName, "hello_intvl", "3000"),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "10000"),
					resource.TestCheckResourceAttr(resourceName, "preempt_delay_min", "0"),
					resource.TestCheckResourceAttr(resourceName, "preempt_delay_reload", "0"),
					resource.TestCheckResourceAttr(resourceName, "preempt_delay_sync", "0"),
					resource.TestCheckResourceAttr(resourceName, "prio", "100"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "0"),
					resource.TestCheckResourceAttr(resourceName, "hsrp_group_policy_type", "simple"),
				),
			},
			{
				Config: CreateAccHSRPGroupPolicyConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists(resourceName, &hsrp_group_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_hsrp_group_policy"),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "preempt"),
					resource.TestCheckResourceAttr(resourceName, "hello_intvl", "6000"),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "18000"),
					resource.TestCheckResourceAttr(resourceName, "key", "cisco"),
					resource.TestCheckResourceAttr(resourceName, "preempt_delay_min", "3600"),
					resource.TestCheckResourceAttr(resourceName, "preempt_delay_reload", "3600"),
					resource.TestCheckResourceAttr(resourceName, "preempt_delay_sync", "3600"),
					resource.TestCheckResourceAttr(resourceName, "prio", "255"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "32767"),
					resource.TestCheckResourceAttr(resourceName, "hsrp_group_policy_type", "md5"),
					testAccCheckAciHSRPGroupPolicyIdEqual(&hsrp_group_policy_default, &hsrp_group_policy_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key"},
			},
			{
				Config:      CreateAccHSRPGroupPolicyConfigWithRequiredParams(rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccHSRPGroupPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccHSRPGroupPolicyConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists(resourceName, &hsrp_group_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciHSRPGroupPolicyIdNotEqual(&hsrp_group_policy_default, &hsrp_group_policy_updated),
				),
			},
			{
				Config: CreateAccHSRPGroupPolicyConfig(rName, rName),
			},
			{
				Config: CreateAccHSRPGroupPolicyConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists(resourceName, &hsrp_group_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciHSRPGroupPolicyIdNotEqual(&hsrp_group_policy_default, &hsrp_group_policy_updated),
				),
			},
		},
	})
}

func TestAccAciHSRPGroupPolicy_Update(t *testing.T) {
	var hsrp_group_policy_default models.HSRPGroupPolicy
	var hsrp_group_policy_updated models.HSRPGroupPolicy
	resourceName := "aci_hsrp_group_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciHSRPGroupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccHSRPGroupPolicyConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists(resourceName, &hsrp_group_policy_default),
				),
			},
			{
				Config: CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "preempt_delay_min", "1800"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists(resourceName, &hsrp_group_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "preempt_delay_min", "1800"),
					testAccCheckAciHSRPGroupPolicyIdEqual(&hsrp_group_policy_default, &hsrp_group_policy_updated),
				),
			},
			{
				Config: CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "preempt_delay_reload", "1800"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists(resourceName, &hsrp_group_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "preempt_delay_reload", "1800"),
					testAccCheckAciHSRPGroupPolicyIdEqual(&hsrp_group_policy_default, &hsrp_group_policy_updated),
				),
			},
			{
				Config: CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "preempt_delay_sync", "1800"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists(resourceName, &hsrp_group_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "preempt_delay_sync", "1800"),
					testAccCheckAciHSRPGroupPolicyIdEqual(&hsrp_group_policy_default, &hsrp_group_policy_updated),
				),
			},
			{
				Config: CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "prio", "0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists(resourceName, &hsrp_group_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "0"),
					testAccCheckAciHSRPGroupPolicyIdEqual(&hsrp_group_policy_default, &hsrp_group_policy_updated),
				),
			},
			{
				Config: CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "prio", "125"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists(resourceName, &hsrp_group_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "125"),
					testAccCheckAciHSRPGroupPolicyIdEqual(&hsrp_group_policy_default, &hsrp_group_policy_updated),
				),
			},
			{
				Config: CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "timeout", "16000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists(resourceName, &hsrp_group_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "timeout", "16000"),
					testAccCheckAciHSRPGroupPolicyIdEqual(&hsrp_group_policy_default, &hsrp_group_policy_updated),
				),
			},
		},
	})
}

func TestAccAciHSRPGroupPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciHSRPGroupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccHSRPGroupPolicyConfig(rName, rName),
			},
			{
				Config:      CreateAccHSRPGroupPolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "ctrl", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "hello_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "hold_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "preempt_delay_min", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "preempt_delay_reload", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "preempt_delay_sync", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "preempt_delay_min", "3601"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "preempt_delay_reload", "3601"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "preempt_delay_sync", "3601"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "prio", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "prio", "256"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "timeout", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "timeout", "32768"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "hsrp_group_policy_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "hello_intvl", "750"),
				ExpectError: regexp.MustCompile(`Invalid Configuration HSRP Hold and Hello timers not allowed in fractions of secs`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "hold_intvl", "1750"),
				ExpectError: regexp.MustCompile(`Invalid Configuration HSRP Hold and Hello timers not allowed in fractions of secs`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, "hold_intvl", "1000"),
				ExpectError: regexp.MustCompile(`Invalid Configuration HSRP Hold timer should be 3X the HSRP Hello timer.`),
			},
			{
				Config:      CreateAccHSRPGroupPolicyUpdatedAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccHSRPGroupPolicyConfig(rName, rName),
			},
		},
	})
}

func TestAccAciHSRPGroupPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciHSRPGroupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccHSRPGroupPolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciHSRPGroupPolicyExists(name string, hsrp_group_policy *models.HSRPGroupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("HSRP Group Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No HSRP Group Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		hsrp_group_policyFound := models.HSRPGroupPolicyFromContainer(cont)
		if hsrp_group_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("HSRP Group Policy %s not found", rs.Primary.ID)
		}
		*hsrp_group_policy = *hsrp_group_policyFound
		return nil
	}
}

func testAccCheckAciHSRPGroupPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing hsrp_group_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_hsrp_group_policy" {
			cont, err := client.Get(rs.Primary.ID)
			hsrp_group_policy := models.HSRPGroupPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("HSRP Group Policy %s Still exists", hsrp_group_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciHSRPGroupPolicyIdEqual(m1, m2 *models.HSRPGroupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("hsrp_group_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciHSRPGroupPolicyIdNotEqual(m1, m2 *models.HSRPGroupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("hsrp_group_policy DNs are equal")
		}
		return nil
	}
}

func CreateHSRPGroupPolicyWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing hsrp_group_policy creation without ", attrName)
	rBlock := `

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_hsrp_group_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccHSRPGroupPolicyConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing hsrp_group_policy creation with tenant name %s and hsrp_group_policy name %s\n", fvTenantName, rName)
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccHSRPGroupPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple hsrp_group_policy creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName, rName)
	return resource
}

func CreateAccHSRPGroupPolicyConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing hsrp_group_policy creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccHSRPGroupPolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing hsrp_group_policy creation with invalid tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccHSRPGroupPolicyConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing hsrp_group_policy creation with optional parameters")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_hsrp_group_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_hsrp_group_policy"
		ctrl = "preempt"
		hello_intvl = "6000"
		hold_intvl = "18000"
		key = "cisco"
		preempt_delay_min = "3600"
		preempt_delay_reload = "3600"
		preempt_delay_sync = "3600"
		prio = "255"
		timeout = "32767"
		hsrp_group_policy_type = "md5"
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccHSRPGroupPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing hsrp_group_policy update without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_hsrp_group_policy" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_hsrp_group_policy"
		ctrl = "preempt"
		hello_intvl = "251"
		hold_intvl = "751"
		key = ""
		preempt_delay_min = "1"
		preempt_delay_reload = "1"
		preempt_delay_sync = "1"
		prio = "1"
		secure_auth_key = ""
		timeout = "32767"
		hsrp_group_policy_type = "md5"
	}
	`)

	return resource
}

func CreateAccHSRPGroupPolicyUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing hsrp_group_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_hsrp_group_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
