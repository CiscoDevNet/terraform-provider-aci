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

func TestAccAciServiceRedirectPolicy_Basic(t *testing.T) {
	var service_redirect_policy_default models.ServiceRedirectPolicy
	var service_redirect_policy_updated models.ServiceRedirectPolicy
	resourceName := "aci_service_redirect_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciServiceRedirectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateServiceRedirectPolicyWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateServiceRedirectPolicyWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccServiceRedirectPolicyConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciServiceRedirectPolicyExists(resourceName, &service_redirect_policy_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "anycast_enabled", "no"),
					resource.TestCheckResourceAttr(resourceName, "dest_type", "L3"),
					resource.TestCheckResourceAttr(resourceName, "hashing_algorithm", "sip-dip-prototype"),
					resource.TestCheckResourceAttr(resourceName, "max_threshold_percent", "0"),
					resource.TestCheckResourceAttr(resourceName, "min_threshold_percent", "0"),
					resource.TestCheckResourceAttr(resourceName, "program_local_pod_only", "no"),
					resource.TestCheckResourceAttr(resourceName, "resilient_hash_enabled", "no"),
					resource.TestCheckResourceAttr(resourceName, "threshold_down_action", "permit"),
					resource.TestCheckResourceAttr(resourceName, "threshold_enable", "no"),
				),
			},
			{
				Config: CreateAccServiceRedirectPolicyConfigWithOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciServiceRedirectPolicyExists(resourceName, &service_redirect_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_service_redirect_policy"),

					resource.TestCheckResourceAttr(resourceName, "anycast_enabled", "no"),

					resource.TestCheckResourceAttr(resourceName, "dest_type", "L1"),

					resource.TestCheckResourceAttr(resourceName, "hashing_algorithm", "dip"),
					resource.TestCheckResourceAttr(resourceName, "max_threshold_percent", "10"),
					resource.TestCheckResourceAttr(resourceName, "min_threshold_percent", "1"),

					resource.TestCheckResourceAttr(resourceName, "program_local_pod_only", "yes"),

					resource.TestCheckResourceAttr(resourceName, "resilient_hash_enabled", "yes"),

					resource.TestCheckResourceAttr(resourceName, "threshold_down_action", "bypass"),

					resource.TestCheckResourceAttr(resourceName, "threshold_enable", "yes"),

					testAccCheckAciServiceRedirectPolicyIdEqual(&service_redirect_policy_default, &service_redirect_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccServiceRedirectPolicyConfigUpdatedName(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccServiceRedirectPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccServiceRedirectPolicyConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciServiceRedirectPolicyExists(resourceName, &service_redirect_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciServiceRedirectPolicyIdNotEqual(&service_redirect_policy_default, &service_redirect_policy_updated),
				),
			},
			{
				Config: CreateAccServiceRedirectPolicyConfig(fvTenantName, rName),
			},
			{
				Config: CreateAccServiceRedirectPolicyConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciServiceRedirectPolicyExists(resourceName, &service_redirect_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciServiceRedirectPolicyIdNotEqual(&service_redirect_policy_default, &service_redirect_policy_updated),
				),
			},
		},
	})
}

func TestAccAciServiceRedirectPolicy_Update(t *testing.T) {
	var service_redirect_policy_default models.ServiceRedirectPolicy
	var service_redirect_policy_updated models.ServiceRedirectPolicy
	resourceName := "aci_service_redirect_policy.test"
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciServiceRedirectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccServiceRedirectPolicyConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciServiceRedirectPolicyExists(resourceName, &service_redirect_policy_default),
				),
			},
			{
				Config: CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "max_threshold_percent", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciServiceRedirectPolicyExists(resourceName, &service_redirect_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "max_threshold_percent", "100"),
					testAccCheckAciServiceRedirectPolicyIdEqual(&service_redirect_policy_default, &service_redirect_policy_updated),
				),
			},
			{
				Config: CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "max_threshold_percent", "50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciServiceRedirectPolicyExists(resourceName, &service_redirect_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "max_threshold_percent", "50"),
					testAccCheckAciServiceRedirectPolicyIdEqual(&service_redirect_policy_default, &service_redirect_policy_updated),
				),
			},
			{
				Config: CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "max_threshold_percent", "100"),
			},
			{
				Config: CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "min_threshold_percent", "99"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciServiceRedirectPolicyExists(resourceName, &service_redirect_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "min_threshold_percent", "99"),
					testAccCheckAciServiceRedirectPolicyIdEqual(&service_redirect_policy_default, &service_redirect_policy_updated),
				),
			},
			{
				Config: CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "min_threshold_percent", "50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciServiceRedirectPolicyExists(resourceName, &service_redirect_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "min_threshold_percent", "50"),
					testAccCheckAciServiceRedirectPolicyIdEqual(&service_redirect_policy_default, &service_redirect_policy_updated),
				),
			},

			{
				Config: CreateAccServiceRedirectPolicyConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciServiceRedirectPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciServiceRedirectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccServiceRedirectPolicyConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccServiceRedirectPolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "anycast_enabled", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "dest_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "hashing_algorithm", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "max_threshold_percent", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "max_threshold_percent", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "max_threshold_percent", "101"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "min_threshold_percent", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "min_threshold_percent", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "min_threshold_percent", "101"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "program_local_pod_only", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "resilient_hash_enabled", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "threshold_down_action", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, "threshold_enable", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccServiceRedirectPolicyConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciServiceRedirectPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciServiceRedirectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccServiceRedirectPolicyConfigMultiple(fvTenantName, rName),
			},
		},
	})
}

func testAccCheckAciServiceRedirectPolicyExists(name string, service_redirect_policy *models.ServiceRedirectPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Service Redirect Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service Redirect Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		service_redirect_policyFound := models.ServiceRedirectPolicyFromContainer(cont)
		if service_redirect_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Service Redirect Policy %s not found", rs.Primary.ID)
		}
		*service_redirect_policy = *service_redirect_policyFound
		return nil
	}
}

func testAccCheckAciServiceRedirectPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing service_redirect_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_service_redirect_policy" {
			cont, err := client.Get(rs.Primary.ID)
			service_redirect_policy := models.ServiceRedirectPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Service Redirect Policy %s Still exists", service_redirect_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciServiceRedirectPolicyIdEqual(m1, m2 *models.ServiceRedirectPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("service_redirect_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciServiceRedirectPolicyIdNotEqual(m1, m2 *models.ServiceRedirectPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("service_redirect_policy DNs are equal")
		}
		return nil
	}
}

func CreateServiceRedirectPolicyWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing service_redirect_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_service_redirect_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccServiceRedirectPolicyConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing service_redirect_policy creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}
func CreateAccServiceRedirectPolicyConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing service_redirect_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccServiceRedirectPolicyConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing service_redirect_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccServiceRedirectPolicyConfigMultiple(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing multiple service_redirect_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccServiceRedirectPolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing service_redirect_policy creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		tenant_dn  = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccServiceRedirectPolicyConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing service_redirect_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_service_redirect_policy"
		anycast_enabled = "no"
		dest_type = "L1"
		hashing_algorithm = "dip"
		max_threshold_percent = "10"
		min_threshold_percent = "1"
		program_local_pod_only = "yes"
		resilient_hash_enabled = "yes"
		threshold_down_action = "bypass"
		threshold_enable = "yes"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccServiceRedirectPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing service_redirect_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_service_redirect_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_service_redirect_policy"
		anycast_enabled = "yes"
		dest_type = "L1"
		hashing_algorithm = "dip"
		max_threshold_percent = "1"
		min_threshold_percent = "1"
		program_local_pod_only = "yes"
		resilient_hash_enabled = "yes"
		threshold_down_action = "bypass"
		threshold_enable = "yes"
		
	}
	`)

	return resource
}

func CreateAccServiceRedirectPolicyUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing service_redirect_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
