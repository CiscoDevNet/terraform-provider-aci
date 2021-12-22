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

func TestAccAciMatchRule_Basic(t *testing.T) {
	var match_rule models.MatchRule
	fv_tenant_name := acctest.RandString(5)
	rtctrl_subj_p_name := acctest.RandString(5)
	description := "match_rule created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMatchRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMatchRuleConfig_basic(fv_tenant_name, rtctrl_subj_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRuleExists("aci_match_rule.foomatch_rule", &match_rule),
					testAccCheckAciMatchRuleAttributes(fv_tenant_name, rtctrl_subj_p_name, description, &match_rule),
				),
			},
		},
	})
}

func TestAccAciMatchRule_Update(t *testing.T) {
	var match_rule models.MatchRule
	fv_tenant_name := acctest.RandString(5)
	rtctrl_subj_p_name := acctest.RandString(5)
	description := "match_rule created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMatchRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMatchRuleConfig_basic(fv_tenant_name, rtctrl_subj_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRuleExists("aci_match_rule.foomatch_rule", &match_rule),
					testAccCheckAciMatchRuleAttributes(fv_tenant_name, rtctrl_subj_p_name, description, &match_rule),
				),
			},
			{
				Config: testAccCheckAciMatchRuleConfig_basic(fv_tenant_name, rtctrl_subj_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchRuleExists("aci_match_rule.foomatch_rule", &match_rule),
					testAccCheckAciMatchRuleAttributes(fv_tenant_name, rtctrl_subj_p_name, description, &match_rule),
				),
			},
		},
	})
}

func testAccCheckAciMatchRuleConfig_basic(fv_tenant_name, rtctrl_subj_p_name string) string {
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

	`, fv_tenant_name, rtctrl_subj_p_name)
}

func testAccCheckAciMatchRuleExists(name string, match_rule *models.MatchRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Match Rule %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Match Rule dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		match_ruleFound := models.MatchRuleFromContainer(cont)
		if match_ruleFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Match Rule %s not found", rs.Primary.ID)
		}
		*match_rule = *match_ruleFound
		return nil
	}
}

func testAccCheckAciMatchRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_match_rule" {
			cont, err := client.Get(rs.Primary.ID)
			match_rule := models.MatchRuleFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Match Rule %s Still exists", match_rule.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciMatchRuleAttributes(fv_tenant_name, rtctrl_subj_p_name, description string, match_rule *models.MatchRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if rtctrl_subj_p_name != GetMOName(match_rule.DistinguishedName) {
			return fmt.Errorf("Bad rtctrl_subj_p %s", GetMOName(match_rule.DistinguishedName))
		}

		if description != match_rule.Description {
			return fmt.Errorf("Bad match_rule Description %s", match_rule.Description)
		}
		return nil
	}
}
