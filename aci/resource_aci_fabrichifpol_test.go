package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLinkLevelPolicy_Basic(t *testing.T) {
	var link_level_policy models.LinkLevelPolicy
	description := "Fabric interface policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLinkLevelPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLinkLevelPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLinkLevelPolicyExists("aci_fabric_if_pol.foolink_level_policy", &link_level_policy),
					testAccCheckAciLinkLevelPolicyAttributes(description, &link_level_policy),
				),
			},
		},
	})
}

func TestAccAciLinkLevelPolicy_update(t *testing.T) {
	var link_level_policy models.LinkLevelPolicy
	description := "Fabric interface policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLinkLevelPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLinkLevelPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLinkLevelPolicyExists("aci_fabric_if_pol.foolink_level_policy", &link_level_policy),
					testAccCheckAciLinkLevelPolicyAttributes(description, &link_level_policy),
				),
			},
			{
				Config: testAccCheckAciLinkLevelPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLinkLevelPolicyExists("aci_fabric_if_pol.foolink_level_policy", &link_level_policy),
					testAccCheckAciLinkLevelPolicyAttributes(description, &link_level_policy),
				),
			},
		},
	})
}

func testAccCheckAciLinkLevelPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_fabric_if_pol" "foolink_level_policy" {
		description = "%s"
		name  = "fabric_if_pol_1"
		annotation  = "example"
		auto_neg  = "on"
		fec_mode  = "inherit"
		link_debounce  = "100"
		name_alias  = "example"
		speed  = "inherit"
	}
	`, description)
}

func testAccCheckAciLinkLevelPolicyExists(name string, link_level_policy *models.LinkLevelPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		link_level_policyFound := models.LinkLevelPolicyFromContainer(cont)
		if link_level_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric interface Policy %s not found", rs.Primary.ID)
		}
		*link_level_policy = *link_level_policyFound
		return nil
	}
}

func testAccCheckAciLinkLevelPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_fabric_if_pol" {
			cont, err := client.Get(rs.Primary.ID)
			link_level_policy := models.LinkLevelPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fabric interface Policy %s Still exists", link_level_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLinkLevelPolicyAttributes(description string, link_level_policy *models.LinkLevelPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != link_level_policy.Description {
			return fmt.Errorf("Bad Fabric interface policy Description %s", link_level_policy.Description)
		}

		if "fabric_if_pol_1" != link_level_policy.Name {
			return fmt.Errorf("Bad Fabric interface policy name %s", link_level_policy.Name)
		}

		if "example" != link_level_policy.Annotation {
			return fmt.Errorf("Bad Fabric interface policy annotation %s", link_level_policy.Annotation)
		}

		if "on" != link_level_policy.AutoNeg {
			return fmt.Errorf("Bad Fabric interface policy auto_neg %s", link_level_policy.AutoNeg)
		}

		if "inherit" != link_level_policy.FecMode {
			return fmt.Errorf("Bad Fabric interface policy fec_mode %s", link_level_policy.FecMode)
		}

		if "100" != link_level_policy.LinkDebounce {
			return fmt.Errorf("Bad Fabric interface policy link_debounce %s", link_level_policy.LinkDebounce)
		}

		if "example" != link_level_policy.NameAlias {
			return fmt.Errorf("Bad Fabric interface policy name_alias %s", link_level_policy.NameAlias)
		}

		if "inherit" != link_level_policy.Speed {
			return fmt.Errorf("Bad Fabric interface policy speed %s", link_level_policy.Speed)
		}

		return nil
	}
}
