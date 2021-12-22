package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciMatchRouteDestinationRule_Basic(t *testing.T) {
	var match_route_destination_rule models.MatchRouteDestinationRule
	fv_tenant_name := acctest.RandString(5)
	rtctrl_subj_p_name := acctest.RandString(5)
	rtctrl_match_rt_dest_name := acctest.RandString(5)
	description := "match_route_destination_rule created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMatchRouteDestinationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMatchRouteDestinationRuleConfig_basic(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_rt_dest_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRouteDestinationRuleExists("aci_match_route_destination_rule.foomatch_route_destination_rule", &match_route_destination_rule),
					testAccCheckAciMatchRouteDestinationRuleAttributes(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_rt_dest_name, description, &match_route_destination_rule),
				),
			},
		},
	})
}

func TestAccAciMatchRouteDestinationRule_Update(t *testing.T) {
	var match_route_destination_rule models.MatchRouteDestinationRule
	fv_tenant_name := acctest.RandString(5)
	rtctrl_subj_p_name := acctest.RandString(5)
	rtctrl_match_rt_dest_name := acctest.RandString(5)
	description := "match_route_destination_rule created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMatchRouteDestinationRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMatchRouteDestinationRuleConfig_basic(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_rt_dest_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRouteDestinationRuleExists("aci_match_route_destination_rule.foomatch_route_destination_rule", &match_route_destination_rule),
					testAccCheckAciMatchRouteDestinationRuleAttributes(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_rt_dest_name, description, &match_route_destination_rule),
				),
			},
			{
				Config: testAccCheckAciMatchRouteDestinationRuleConfig_basic(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_rt_dest_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRouteDestinationRuleExists("aci_match_route_destination_rule.foomatch_route_destination_rule", &match_route_destination_rule),
					testAccCheckAciMatchRouteDestinationRuleAttributes(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_rt_dest_name, description, &match_route_destination_rule),
				),
			},
		},
	})
}

func testAccCheckAciMatchRouteDestinationRuleConfig_basic(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_rt_dest_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_match_rule" "foomatch_rule" {
		name 		= "%s"
		description = "match_rule created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	resource "aci_match_route_destination_rule" "foomatch_route_destination_rule" {
		name 		= "%s"
		description = "match_route_destination_rule created while acceptance testing"
		match_rule_dn = aci_match_rule.foomatch_rule.id
	}

	`, fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_rt_dest_name)
}

func testAccCheckAciMatchRouteDestinationRuleExists(name string, match_route_destination_rule *models.MatchRouteDestinationRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Match Route Destination Rule %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Match Route Destination Rule dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		match_route_destination_ruleFound := models.MatchRouteDestinationRuleFromContainer(cont)
		if match_route_destination_ruleFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Match Route Destination Rule %s not found", rs.Primary.ID)
		}
		*match_route_destination_rule = *match_route_destination_ruleFound
		return nil
	}
}

func testAccCheckAciMatchRouteDestinationRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_match_route_destination_rule" {
			cont, err := client.Get(rs.Primary.ID)
			match_route_destination_rule := models.MatchRouteDestinationRuleFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Match Route Destination Rule %s Still exists", match_route_destination_rule.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciMatchRouteDestinationRuleAttributes(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_rt_dest_name, description string, match_route_destination_rule *models.MatchRouteDestinationRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if rtctrl_match_rt_dest_name != GetMOName(match_route_destination_rule.DistinguishedName) {
			return fmt.Errorf("Bad rtctrl_match_rt_dest %s", GetMOName(match_route_destination_rule.DistinguishedName))
		}

		if description != match_route_destination_rule.Description {
			return fmt.Errorf("Bad match_route_destination_rule Description %s", match_route_destination_rule.Description)
		}
		return nil
	}
}
