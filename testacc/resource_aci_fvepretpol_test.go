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

func TestAccAciEndPointRetentionPolicy_Basic(t *testing.T) {
	var end_point_retention_policy_default models.EndPointRetentionPolicy
	var end_point_retention_policy_updated models.EndPointRetentionPolicy
	resourceName := "aci_end_point_retention_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEndPointRetentionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateEndPointRetentionPolicyWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateEndPointRetentionPolicyWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEndPointRetentionPolicyConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndPointRetentionPolicyExists(resourceName, &end_point_retention_policy_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "bounce_age_intvl", "630"),
					resource.TestCheckResourceAttr(resourceName, "bounce_trig", "protocol"),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "300"),
					resource.TestCheckResourceAttr(resourceName, "local_ep_age_intvl", "900"),
					resource.TestCheckResourceAttr(resourceName, "move_freq", "256"),
					resource.TestCheckResourceAttr(resourceName, "remote_ep_age_intvl", "300"),
				),
			},
			{
				Config: CreateAccEndPointRetentionPolicyConfigWithOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndPointRetentionPolicyExists(resourceName, &end_point_retention_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_end_point_retention_policy"),
					resource.TestCheckResourceAttr(resourceName, "bounce_age_intvl", "600"),
					resource.TestCheckResourceAttr(resourceName, "bounce_trig", "rarp-flood"),
					resource.TestCheckResourceAttr(resourceName, "hold_intvl", "6"),
					resource.TestCheckResourceAttr(resourceName, "local_ep_age_intvl", "600"),
					resource.TestCheckResourceAttr(resourceName, "move_freq", "1"),
					resource.TestCheckResourceAttr(resourceName, "remote_ep_age_intvl", "600"),

					testAccCheckAciEndPointRetentionPolicyIdEqual(&end_point_retention_policy_default, &end_point_retention_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccEndPointRetentionPolicyConfigUpdatedName(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccEndPointRetentionPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccEndPointRetentionPolicyConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndPointRetentionPolicyExists(resourceName, &end_point_retention_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciEndPointRetentionPolicyIdNotEqual(&end_point_retention_policy_default, &end_point_retention_policy_updated),
				),
			},
			{
				Config: CreateAccEndPointRetentionPolicyConfig(fvTenantName, rName),
			},
			{
				Config: CreateAccEndPointRetentionPolicyConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndPointRetentionPolicyExists(resourceName, &end_point_retention_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciEndPointRetentionPolicyIdNotEqual(&end_point_retention_policy_default, &end_point_retention_policy_updated),
				),
			},
		},
	})
}

func TestAccAciEndPointRetentionPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEndPointRetentionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndPointRetentionPolicyConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "bounce_age_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "bounce_age_intvl", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "bounce_age_intvl", "1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "bounce_trig", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "hold_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "hold_intvl", "4"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "hold_intvl", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "local_ep_age_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "local_ep_age_intvl", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "local_ep_age_intvl", "1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "move_freq", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "move_freq", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "move_freq", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "remote_ep_age_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "remote_ep_age_intvl", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, "remote_ep_age_intvl", "1"),
				ExpectError: regexp.MustCompile(`out of range`),
			},

			{
				Config:      CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEndPointRetentionPolicyConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciEndPointRetentionPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciEndPointRetentionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndPointRetentionPolicyConfigMultiple(fvTenantName, rName),
			},
		},
	})
}

func testAccCheckAciEndPointRetentionPolicyExists(name string, end_point_retention_policy *models.EndPointRetentionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("End Point Retention Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No End Point Retention Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		end_point_retention_policyFound := models.EndPointRetentionPolicyFromContainer(cont)
		if end_point_retention_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("End Point Retention Policy %s not found", rs.Primary.ID)
		}
		*end_point_retention_policy = *end_point_retention_policyFound
		return nil
	}
}

func testAccCheckAciEndPointRetentionPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing end_point_retention_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_end_point_retention_policy" {
			cont, err := client.Get(rs.Primary.ID)
			end_point_retention_policy := models.EndPointRetentionPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("End Point Retention Policy %s Still exists", end_point_retention_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciEndPointRetentionPolicyIdEqual(m1, m2 *models.EndPointRetentionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("end_point_retention_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciEndPointRetentionPolicyIdNotEqual(m1, m2 *models.EndPointRetentionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("end_point_retention_policy DNs are equal")
		}
		return nil
	}
}

func CreateEndPointRetentionPolicyWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing end_point_retention_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_end_point_retention_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccEndPointRetentionPolicyConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing end_point_retention_policy creation with parent resource name %s and resource name %s\n", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}
func CreateAccEndPointRetentionPolicyConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing end_point_retention_policy creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccEndPointRetentionPolicyConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing end_point_retention_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccEndPointRetentionPolicyConfigMultiple(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing multiple end_point_retention_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccEndPointRetentionPolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing end_point_retention_policy creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		name = "acctest_ap"
		tenant_dn  = aci_tenant.test.id
	}
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccEndPointRetentionPolicyConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing end_point_retention_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_end_point_retention_policy"
		bounce_age_intvl = "600"
		bounce_trig = "rarp-flood"
		hold_intvl = "6"
		local_ep_age_intvl = "600"
		move_freq = "1"
		remote_ep_age_intvl = "600"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccEndPointRetentionPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing end_point_retention_policy updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_end_point_retention_policy" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_end_point_retention_policy"
		bounce_age_intvl = "1"
		bounce_trig = "rarp-flood"
		hold_intvl = "6"
		local_ep_age_intvl = "1"
		move_freq = "1"
		remote_ep_age_intvl = "1"
		
	}
	`)

	return resource
}

func CreateAccEndPointRetentionPolicyUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing end_point_retention_policy attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_end_point_retention_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
