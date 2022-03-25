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

func TestAccAciActionRuleProfile_Basic(t *testing.T) {
	var action_rule_profile models.ActionRuleProfile
	fv_tenant_name := acctest.RandString(5)
	rtctrl_attr_p_name := acctest.RandString(5)
	description := "action_rule_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciActionRuleProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciActionRuleProfileConfig_basic(fv_tenant_name, rtctrl_attr_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists("aci_action_rule_profile.fooaction_rule_profile", &action_rule_profile),
					testAccCheckAciActionRuleProfileAttributes(fv_tenant_name, rtctrl_attr_p_name, description, &action_rule_profile),
				),
			},
		},
	})
}

func TestAccAciActionRuleProfile_Update(t *testing.T) {
	var action_rule_profile models.ActionRuleProfile
	fv_tenant_name := acctest.RandString(5)
	rtctrl_attr_p_name := acctest.RandString(5)
	description := "action_rule_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciActionRuleProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciActionRuleProfileConfig_basic(fv_tenant_name, rtctrl_attr_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists("aci_action_rule_profile.fooaction_rule_profile", &action_rule_profile),
					testAccCheckAciActionRuleProfileAttributes(fv_tenant_name, rtctrl_attr_p_name, description, &action_rule_profile),
				),
			},
			{
				Config: testAccCheckAciActionRuleProfileConfig_basic(fv_tenant_name, rtctrl_attr_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists("aci_action_rule_profile.fooaction_rule_profile", &action_rule_profile),
					testAccCheckAciActionRuleProfileAttributes(fv_tenant_name, rtctrl_attr_p_name, description, &action_rule_profile),
				),
			},
		},
	})
}

func testAccCheckAciActionRuleProfileConfig_basic(fv_tenant_name, rtctrl_attr_p_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_action_rule_profile" "fooaction_rule_profile" {
		name 		= "%s"
		description = "action_rule_profile created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	`, fv_tenant_name, rtctrl_attr_p_name)
}

func testAccCheckAciActionRuleProfileExists(name string, action_rule_profile *models.ActionRuleProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Action Rule Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Action Rule Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		action_rule_profileFound := models.ActionRuleProfileFromContainer(cont)
		if action_rule_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Action Rule Profile %s not found", rs.Primary.ID)
		}
		*action_rule_profile = *action_rule_profileFound
		return nil
	}
}

func testAccCheckAciActionRuleProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_action_rule_profile" {
			cont, err := client.Get(rs.Primary.ID)
			action_rule_profile := models.ActionRuleProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Action Rule Profile %s Still exists", action_rule_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciActionRuleProfileAttributes(fv_tenant_name, rtctrl_attr_p_name, description string, action_rule_profile *models.ActionRuleProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if rtctrl_attr_p_name != GetMOName(action_rule_profile.DistinguishedName) {
			return fmt.Errorf("Bad rtctrl_attr_p %s", GetMOName(action_rule_profile.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(action_rule_profile.DistinguishedName, action_rule_profile.Rn)) {
			return fmt.Errorf(" Bad fv_tenant %s", GetMOName(GetParentDn(action_rule_profile.DistinguishedName, action_rule_profile.Rn)))
		}
		if description != action_rule_profile.Description {
			return fmt.Errorf("Bad action_rule_profile Description %s", action_rule_profile.Description)
		}
		return nil
	}
}
