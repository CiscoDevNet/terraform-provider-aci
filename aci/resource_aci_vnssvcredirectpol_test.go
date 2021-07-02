package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciServiceRedirectPolicy_Basic(t *testing.T) {
	var service_redirect_policy models.ServiceRedirectPolicy
	description := "service_redirect_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciServiceRedirectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciServiceRedirectPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciServiceRedirectPolicyExists("aci_service_redirect_policy.example", &service_redirect_policy),
					testAccCheckAciServiceRedirectPolicyAttributes(description, &service_redirect_policy),
				),
			},
		},
	})
}

func TestAccAciServiceRedirectPolicy_update(t *testing.T) {
	var service_redirect_policy models.ServiceRedirectPolicy
	description := "service_redirect_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciServiceRedirectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciServiceRedirectPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciServiceRedirectPolicyExists("aci_service_redirect_policy.example", &service_redirect_policy),
					testAccCheckAciServiceRedirectPolicyAttributes(description, &service_redirect_policy),
				),
			},
			{
				Config: testAccCheckAciServiceRedirectPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciServiceRedirectPolicyExists("aci_service_redirect_policy.example", &service_redirect_policy),
					testAccCheckAciServiceRedirectPolicyAttributes(description, &service_redirect_policy),
				),
			},
		},
	})
}

func testAccCheckAciServiceRedirectPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_service_redirect_policy" "example" {
		tenant_dn  = "${aci_tenant.tenentcheck.id}"
		name  = "check"
		dest_type = "L3"
		max_threshold_percent = "50"
		hashing_algorithm = "sip"
		description = "%s"
		anycast_enabled = "yes"
		resilient_hash_enabled = "no"
		threshold_enable = "yes"
	}
	`, description)
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

func testAccCheckAciServiceRedirectPolicyAttributes(description string, service_redirect_policy *models.ServiceRedirectPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != service_redirect_policy.Description {
			return fmt.Errorf("Bad service_redirect_policy Description %s", service_redirect_policy.Description)
		}

		if "check" != service_redirect_policy.Name {
			return fmt.Errorf("Bad service_redirect_policy name %s", service_redirect_policy.Name)
		}

		if "yes" != service_redirect_policy.AnycastEnabled {
			return fmt.Errorf("Bad service_redirect_policy anycast_enabled %s", service_redirect_policy.AnycastEnabled)
		}

		if "L3" != service_redirect_policy.DestType {
			return fmt.Errorf("Bad service_redirect_policy dest_type %s", service_redirect_policy.DestType)
		}

		if "50" != service_redirect_policy.MaxThresholdPercent {
			return fmt.Errorf("Bad service_redirect_policy max_threshold_percent %s", service_redirect_policy.MaxThresholdPercent)
		}

		if "no" != service_redirect_policy.ResilientHashEnabled {
			return fmt.Errorf("Bad service_redirect_policy resilient_hash_enabled %s", service_redirect_policy.ResilientHashEnabled)
		}

		if "yes" != service_redirect_policy.ThresholdEnable {
			return fmt.Errorf("Bad service_redirect_policy threshold_enable %s", service_redirect_policy.ThresholdEnable)
		}

		return nil
	}
}
