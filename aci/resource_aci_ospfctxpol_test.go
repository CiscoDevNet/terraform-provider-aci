package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciOSPFTimersPolicy_Basic(t *testing.T) {
	var ospf_timers_policy models.OSPFTimersPolicy
	description := "ospf_timers_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOSPFTimersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOSPFTimersPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersPolicyExists("aci_ospf_timers.test", &ospf_timers_policy),
					testAccCheckAciOSPFTimersPolicyAttributes(description, &ospf_timers_policy),
				),
			},
		},
	})
}

func TestAccAciOSPFTimersPolicy_update(t *testing.T) {
	var ospf_timers_policy models.OSPFTimersPolicy
	description := "ospf_timers_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOSPFTimersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOSPFTimersPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersPolicyExists("aci_ospf_timers.test", &ospf_timers_policy),
					testAccCheckAciOSPFTimersPolicyAttributes(description, &ospf_timers_policy),
				),
			},
			{
				Config: testAccCheckAciOSPFTimersPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersPolicyExists("aci_ospf_timers.test", &ospf_timers_policy),
					testAccCheckAciOSPFTimersPolicyAttributes(description, &ospf_timers_policy),
				),
			},
		},
	})
}

func testAccCheckAciOSPFTimersPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_ospf_timers" "test" {
		tenant_dn           = "uni/tn-aaaaa"
		name                = "ospf_timers_1"
		annotation          = "example"
		description 		= "%s"
		bw_ref              = "30000"
		ctrl                = "name-lookup"
		dist                = "200"
		gr_ctrl             = "helper"
		lsa_arrival_intvl   = "2000"
		lsa_gp_pacing_intvl = "50"
		lsa_hold_intvl      = "1000"
		lsa_max_intvl       = "1000"
		lsa_start_intvl     = "5"
		max_ecmp            = "10"
		max_lsa_action      = "restart"
		max_lsa_num         = "56"
		max_lsa_reset_intvl = "10"
		max_lsa_sleep_cnt   = "10"
		max_lsa_sleep_intvl = "10"
		max_lsa_thresh      = "50"
		name_alias          = "example"
		spf_hold_intvl      = "100"
		spf_init_intvl      = "500"
		spf_max_intvl       = "10"
	  }
	`, description)
}

func testAccCheckAciOSPFTimersPolicyExists(name string, ospf_timers_policy *models.OSPFTimersPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("OSPF Timers Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OSPF Timers Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ospf_timers_policyFound := models.OSPFTimersPolicyFromContainer(cont)
		if ospf_timers_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("OSPF Timers Policy %s not found", rs.Primary.ID)
		}
		*ospf_timers_policy = *ospf_timers_policyFound
		return nil
	}
}

func testAccCheckAciOSPFTimersPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_ospf_timers" {
			cont, err := client.Get(rs.Primary.ID)
			ospf_timers_policy := models.OSPFTimersPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("OSPF Timers Policy %s Still exists", ospf_timers_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciOSPFTimersPolicyAttributes(description string, ospf_timers_policy *models.OSPFTimersPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != ospf_timers_policy.Description {
			return fmt.Errorf("Bad ospf_timers_policy Description %s", ospf_timers_policy.Description)
		}

		if "ospf_timers_1" != ospf_timers_policy.Name {
			return fmt.Errorf("Bad ospf_timers_policy name %s", ospf_timers_policy.Name)
		}

		if "example" != ospf_timers_policy.Annotation {
			return fmt.Errorf("Bad ospf_timers_policy annotation %s", ospf_timers_policy.Annotation)
		}

		if "30000" != ospf_timers_policy.BwRef {
			return fmt.Errorf("Bad ospf_timers_policy bw_ref %s", ospf_timers_policy.BwRef)
		}

		if "name-lookup" != ospf_timers_policy.Ctrl {
			return fmt.Errorf("Bad ospf_timers_policy ctrl %s", ospf_timers_policy.Ctrl)
		}

		if "200" != ospf_timers_policy.Dist {
			return fmt.Errorf("Bad ospf_timers_policy dist %s", ospf_timers_policy.Dist)
		}

		if "helper" != ospf_timers_policy.GrCtrl {
			return fmt.Errorf("Bad ospf_timers_policy gr_ctrl %s", ospf_timers_policy.GrCtrl)
		}

		if "2000" != ospf_timers_policy.LsaArrivalIntvl {
			return fmt.Errorf("Bad ospf_timers_policy lsa_arrival_intvl %s", ospf_timers_policy.LsaArrivalIntvl)
		}

		if "50" != ospf_timers_policy.LsaGpPacingIntvl {
			return fmt.Errorf("Bad ospf_timers_policy lsa_gp_pacing_intvl %s", ospf_timers_policy.LsaGpPacingIntvl)
		}

		if "1000" != ospf_timers_policy.LsaHoldIntvl {
			return fmt.Errorf("Bad ospf_timers_policy lsa_hold_intvl %s", ospf_timers_policy.LsaHoldIntvl)
		}

		if "1000" != ospf_timers_policy.LsaMaxIntvl {
			return fmt.Errorf("Bad ospf_timers_policy lsa_max_intvl %s", ospf_timers_policy.LsaMaxIntvl)
		}

		if "5" != ospf_timers_policy.LsaStartIntvl {
			return fmt.Errorf("Bad ospf_timers_policy lsa_start_intvl %s", ospf_timers_policy.LsaStartIntvl)
		}

		if "10" != ospf_timers_policy.MaxEcmp {
			return fmt.Errorf("Bad ospf_timers_policy max_ecmp %s", ospf_timers_policy.MaxEcmp)
		}

		if "restart" != ospf_timers_policy.MaxLsaAction {
			return fmt.Errorf("Bad ospf_timers_policy max_lsa_action %s", ospf_timers_policy.MaxLsaAction)
		}

		if "56" != ospf_timers_policy.MaxLsaNum {
			return fmt.Errorf("Bad ospf_timers_policy max_lsa_num %s", ospf_timers_policy.MaxLsaNum)
		}

		if "10" != ospf_timers_policy.MaxLsaResetIntvl {
			return fmt.Errorf("Bad ospf_timers_policy max_lsa_reset_intvl %s", ospf_timers_policy.MaxLsaResetIntvl)
		}

		if "10" != ospf_timers_policy.MaxLsaSleepCnt {
			return fmt.Errorf("Bad ospf_timers_policy max_lsa_sleep_cnt %s", ospf_timers_policy.MaxLsaSleepCnt)
		}

		if "10" != ospf_timers_policy.MaxLsaSleepIntvl {
			return fmt.Errorf("Bad ospf_timers_policy max_lsa_sleep_intvl %s", ospf_timers_policy.MaxLsaSleepIntvl)
		}

		if "50" != ospf_timers_policy.MaxLsaThresh {
			return fmt.Errorf("Bad ospf_timers_policy max_lsa_thresh %s", ospf_timers_policy.MaxLsaThresh)
		}

		if "example" != ospf_timers_policy.NameAlias {
			return fmt.Errorf("Bad ospf_timers_policy name_alias %s", ospf_timers_policy.NameAlias)
		}

		if "100" != ospf_timers_policy.SpfHoldIntvl {
			return fmt.Errorf("Bad ospf_timers_policy spf_hold_intvl %s", ospf_timers_policy.SpfHoldIntvl)
		}

		if "500" != ospf_timers_policy.SpfInitIntvl {
			return fmt.Errorf("Bad ospf_timers_policy spf_init_intvl %s", ospf_timers_policy.SpfInitIntvl)
		}

		if "10" != ospf_timers_policy.SpfMaxIntvl {
			return fmt.Errorf("Bad ospf_timers_policy spf_max_intvl %s", ospf_timers_policy.SpfMaxIntvl)
		}

		return nil
	}
}
