package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBGPTimersPolicy_Basic(t *testing.T) {
	var bgp_timers_policy models.BGPTimersPolicy
	description := "bgp_timers_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBGPTimersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBGPTimersPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPTimersPolicyExists("aci_bgp_timers.test", &bgp_timers_policy),
					testAccCheckAciBGPTimersPolicyAttributes(description, &bgp_timers_policy),
				),
			},
		},
	})
}

func TestAccAciBGPTimersPolicy_update(t *testing.T) {
	var bgp_timers_policy models.BGPTimersPolicy
	description := "bgp_timers_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBGPTimersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBGPTimersPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPTimersPolicyExists("aci_bgp_timers.test", &bgp_timers_policy),
					testAccCheckAciBGPTimersPolicyAttributes(description, &bgp_timers_policy),
				),
			},
			{
				Config: testAccCheckAciBGPTimersPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPTimersPolicyExists("aci_bgp_timers.test", &bgp_timers_policy),
					testAccCheckAciBGPTimersPolicyAttributes(description, &bgp_timers_policy),
				),
			},
		},
	})
}

func testAccCheckAciBGPTimersPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_bgp_timers" "test" {
		tenant_dn    = aci_tenant.tenentcheck.id
		name         = "one"
		description  = "%s"
		annotation   = "example"
		gr_ctrl      = "helper"
		hold_intvl   = "189"
		ka_intvl     = "65"
		max_as_limit = "70"
		name_alias   = "aliasing"
		stale_intvl  = "15"
	}
	`, description)
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

func testAccCheckAciBGPTimersPolicyAttributes(description string, bgp_timers_policy *models.BGPTimersPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != bgp_timers_policy.Description {
			return fmt.Errorf("Bad bgp_timers_policy Description %s", bgp_timers_policy.Description)
		}

		if "one" != bgp_timers_policy.Name {
			return fmt.Errorf("Bad bgp_timers_policy name %s", bgp_timers_policy.Name)
		}

		if "example" != bgp_timers_policy.Annotation {
			return fmt.Errorf("Bad bgp_timers_policy annotation %s", bgp_timers_policy.Annotation)
		}

		if "helper" != bgp_timers_policy.GrCtrl {
			return fmt.Errorf("Bad bgp_timers_policy gr_ctrl %s", bgp_timers_policy.GrCtrl)
		}

		if "189" != bgp_timers_policy.HoldIntvl {
			return fmt.Errorf("Bad bgp_timers_policy hold_intvl %s", bgp_timers_policy.HoldIntvl)
		}

		if "65" != bgp_timers_policy.KaIntvl {
			return fmt.Errorf("Bad bgp_timers_policy ka_intvl %s", bgp_timers_policy.KaIntvl)
		}

		if "70" != bgp_timers_policy.MaxAsLimit {
			return fmt.Errorf("Bad bgp_timers_policy max_as_limit %s", bgp_timers_policy.MaxAsLimit)
		}

		if "aliasing" != bgp_timers_policy.NameAlias {
			return fmt.Errorf("Bad bgp_timers_policy name_alias %s", bgp_timers_policy.NameAlias)
		}

		if "15" != bgp_timers_policy.StaleIntvl {
			return fmt.Errorf("Bad bgp_timers_policy stale_intvl %s", bgp_timers_policy.StaleIntvl)
		}

		return nil
	}
}
