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

func TestAccAciMatchRuleBasedonCommunityRegularExpression_Basic(t *testing.T) {
	var match_rule_basedon_community_regular_expression models.MatchRuleBasedonCommunityRegularExpression
	fv_tenant_name := acctest.RandString(5)
	rtctrl_subj_p_name := acctest.RandString(5)
	rtctrl_match_comm_regex_term_name := acctest.RandString(5)
	description := "match_rule_basedon_community_regular_expression created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMatchRuleBasedonCommunityRegularExpressionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMatchRuleBasedonCommunityRegularExpressionConfig_basic(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_regex_term_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRuleBasedonCommunityRegularExpressionExists("aci_match_rule_basedon_community_regular_expression.foomatch_rule_basedon_community_regular_expression", &match_rule_basedon_community_regular_expression),
					testAccCheckAciMatchRuleBasedonCommunityRegularExpressionAttributes(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_regex_term_name, description, &match_rule_basedon_community_regular_expression),
				),
			},
		},
	})
}

func TestAccAciMatchRuleBasedonCommunityRegularExpression_Update(t *testing.T) {
	var match_rule_basedon_community_regular_expression models.MatchRuleBasedonCommunityRegularExpression
	fv_tenant_name := acctest.RandString(5)
	rtctrl_subj_p_name := acctest.RandString(5)
	rtctrl_match_comm_regex_term_name := acctest.RandString(5)
	description := "match_rule_basedon_community_regular_expression created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMatchRuleBasedonCommunityRegularExpressionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMatchRuleBasedonCommunityRegularExpressionConfig_basic(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_regex_term_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRuleBasedonCommunityRegularExpressionExists("aci_match_rule_basedon_community_regular_expression.foomatch_rule_basedon_community_regular_expression", &match_rule_basedon_community_regular_expression),
					testAccCheckAciMatchRuleBasedonCommunityRegularExpressionAttributes(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_regex_term_name, description, &match_rule_basedon_community_regular_expression),
				),
			},
			{
				Config: testAccCheckAciMatchRuleBasedonCommunityRegularExpressionConfig_basic(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_regex_term_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRuleBasedonCommunityRegularExpressionExists("aci_match_rule_basedon_community_regular_expression.foomatch_rule_basedon_community_regular_expression", &match_rule_basedon_community_regular_expression),
					testAccCheckAciMatchRuleBasedonCommunityRegularExpressionAttributes(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_regex_term_name, description, &match_rule_basedon_community_regular_expression),
				),
			},
		},
	})
}

func testAccCheckAciMatchRuleBasedonCommunityRegularExpressionConfig_basic(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_regex_term_name string) string {
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

	resource "aci_match_rule_basedon_community_regular_expression" "foomatch_rule_basedon_community_regular_expression" {
		name 		= "%s"
		description = "match_rule_basedon_community_regular_expression created while acceptance testing"
		match_rule_dn = aci_match_rule.foomatch_rule.id
	}

	`, fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_regex_term_name)
}

func testAccCheckAciMatchRuleBasedonCommunityRegularExpressionExists(name string, match_rule_basedon_community_regular_expression *models.MatchRuleBasedonCommunityRegularExpression) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Match Rule Based on Community Regular Expression %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Match Rule Based on Community Regular Expression dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		match_rule_basedon_community_regular_expressionFound := models.MatchRuleBasedonCommunityRegularExpressionFromContainer(cont)
		if match_rule_basedon_community_regular_expressionFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Match Rule Based on Community Regular Expression %s not found", rs.Primary.ID)
		}
		*match_rule_basedon_community_regular_expression = *match_rule_basedon_community_regular_expressionFound
		return nil
	}
}

func testAccCheckAciMatchRuleBasedonCommunityRegularExpressionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_match_rule_basedon_community_regular_expression" {
			cont, err := client.Get(rs.Primary.ID)
			match_rule_basedon_community_regular_expression := models.MatchRuleBasedonCommunityRegularExpressionFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Match Rule Based on Community Regular Expression %s Still exists", match_rule_basedon_community_regular_expression.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciMatchRuleBasedonCommunityRegularExpressionAttributes(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_regex_term_name, description string, match_rule_basedon_community_regular_expression *models.MatchRuleBasedonCommunityRegularExpression) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if rtctrl_match_comm_regex_term_name != GetMOName(match_rule_basedon_community_regular_expression.DistinguishedName) {
			return fmt.Errorf("Bad rtctrl_match_comm_regex_term %s", GetMOName(match_rule_basedon_community_regular_expression.DistinguishedName))
		}

		if rtctrl_subj_p_name != GetMOName(GetParentDn(match_rule_basedon_community_regular_expression.DistinguishedName)) {
			return fmt.Errorf(" Bad rtctrl_subj_p %s", GetMOName(GetParentDn(match_rule_basedon_community_regular_expression.DistinguishedName)))
		}
		if description != match_rule_basedon_community_regular_expression.Description {
			return fmt.Errorf("Bad match_rule_basedon_community_regular_expression Description %s", match_rule_basedon_community_regular_expression.Description)
		}
		return nil
	}
}
