package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciActionRuleProfile_Basic(t *testing.T) {
	var action_rule_profile models.ActionRuleProfile
	description := "action_rule_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciActionRuleProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciActionRuleProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists("aci_action_rule_profile.fooaction_rule_profile", &action_rule_profile),
					testAccCheckAciActionRuleProfileAttributes(description, &action_rule_profile),
				),
			},
		},
	})
}

func TestAccAciActionRuleProfile_update(t *testing.T) {
	var action_rule_profile models.ActionRuleProfile
	description := "action_rule_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciActionRuleProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciActionRuleProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists("aci_action_rule_profile.fooaction_rule_profile", &action_rule_profile),
					testAccCheckAciActionRuleProfileAttributes(description, &action_rule_profile),
				),
			},
			{
				Config: testAccCheckAciActionRuleProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists("aci_action_rule_profile.fooaction_rule_profile", &action_rule_profile),
					testAccCheckAciActionRuleProfileAttributes(description, &action_rule_profile),
				),
			},
		},
	})
}

func testAccCheckAciActionRuleProfileConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_action_rule_profile" "fooaction_rule_profile" {
		  	tenant_dn  = aci_tenant.example.id
			description = "%s"
			name  = "example"
		  	annotation  = "example"
		  	name_alias  = "example"
		}
	`, description)
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

func testAccCheckAciActionRuleProfileAttributes(description string, action_rule_profile *models.ActionRuleProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != action_rule_profile.Description {
			return fmt.Errorf("Bad action_rule_profile Description %s", action_rule_profile.Description)
		}

		if "example" != action_rule_profile.Name {
			return fmt.Errorf("Bad action_rule_profile name %s", action_rule_profile.Name)
		}

		if "example" != action_rule_profile.Annotation {
			return fmt.Errorf("Bad action_rule_profile annotation %s", action_rule_profile.Annotation)
		}

		if "example" != action_rule_profile.NameAlias {
			return fmt.Errorf("Bad action_rule_profile name_alias %s", action_rule_profile.NameAlias)
		}

		return nil
	}
}
