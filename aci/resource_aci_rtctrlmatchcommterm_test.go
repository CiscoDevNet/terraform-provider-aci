package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciMatchCommunityTerm_Basic(t *testing.T) {
	var match_community_term models.MatchCommunityTerm
	fv_tenant_name := acctest.RandString(5)
	rtctrl_subj_p_name := acctest.RandString(5)
	rtctrl_match_comm_term_name := acctest.RandString(5)
	description := "match_community_term created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMatchCommunityTermDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMatchCommunityTermConfig_basic(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_term_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchCommunityTermExists("aci_match_community_term.foomatch_community_term", &match_community_term),
					testAccCheckAciMatchCommunityTermAttributes(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_term_name, description, &match_community_term),
				),
			},
		},
	})
}

func TestAccAciMatchCommunityTerm_Update(t *testing.T) {
	var match_community_term models.MatchCommunityTerm
	fv_tenant_name := acctest.RandString(5)
	rtctrl_subj_p_name := acctest.RandString(5)
	rtctrl_match_comm_term_name := acctest.RandString(5)
	description := "match_community_term created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMatchCommunityTermDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMatchCommunityTermConfig_basic(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_term_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchCommunityTermExists("aci_match_community_term.foomatch_community_term", &match_community_term),
					testAccCheckAciMatchCommunityTermAttributes(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_term_name, description, &match_community_term),
				),
			},
			{
				Config: testAccCheckAciMatchCommunityTermConfig_basic(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_term_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMatchCommunityTermExists("aci_match_community_term.foomatch_community_term", &match_community_term),
					testAccCheckAciMatchCommunityTermAttributes(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_term_name, description, &match_community_term),
				),
			},
		},
	})
}

func testAccCheckAciMatchCommunityTermConfig_basic(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_term_name string) string {
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

	resource "aci_match_community_term" "foomatch_community_term" {
		name 		= "%s"
		description = "match_community_term created while acceptance testing"
		match_rule_dn = aci_match_rule.foomatch_rule.id
	}

	`, fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_term_name)
}

func testAccCheckAciMatchCommunityTermExists(name string, match_community_term *models.MatchCommunityTerm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Match Community Term %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Match Community Term dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		match_community_termFound := models.MatchCommunityTermFromContainer(cont)
		if match_community_termFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Match Community Term %s not found", rs.Primary.ID)
		}
		*match_community_term = *match_community_termFound
		return nil
	}
}

func testAccCheckAciMatchCommunityTermDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_match_community_term" {
			cont, err := client.Get(rs.Primary.ID)
			match_community_term := models.MatchCommunityTermFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Match Community Term %s Still exists", match_community_term.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciMatchCommunityTermAttributes(fv_tenant_name, rtctrl_subj_p_name, rtctrl_match_comm_term_name, description string, match_community_term *models.MatchCommunityTerm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if rtctrl_match_comm_term_name != GetMOName(match_community_term.DistinguishedName) {
			return fmt.Errorf("Bad rtctrl_match_comm_term %s", GetMOName(match_community_term.DistinguishedName))
		}

		if rtctrl_subj_p_name != GetMOName(GetParentDn(match_community_term.DistinguishedName)) {
			return fmt.Errorf(" Bad rtctrl_subj_p %s", GetMOName(GetParentDn(match_community_term.DistinguishedName)))
		}
		if description != match_community_term.Description {
			return fmt.Errorf("Bad match_community_term Description %s", match_community_term.Description)
		}
		return nil
	}
}
