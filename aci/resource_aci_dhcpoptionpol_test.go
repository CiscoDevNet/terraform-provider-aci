package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciDHCPOptionPolicy_Basic(t *testing.T) {
	var dhcp_option_policy models.DHCPOptionPolicy
	description := "dhcp_option_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDHCPOptionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDHCPOptionPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPOptionPolicyExists("aci_dhcp_option_policy.foodhcp_option_policy", &dhcp_option_policy),
					testAccCheckAciDHCPOptionPolicyAttributes(description, &dhcp_option_policy),
				),
			},
		},
	})
}

func TestAccAciDHCPOptionPolicy_update(t *testing.T) {
	var dhcp_option_policy models.DHCPOptionPolicy
	description := "dhcp_option_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDHCPOptionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDHCPOptionPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPOptionPolicyExists("aci_dhcp_option_policy.foodhcp_option_policy", &dhcp_option_policy),
					testAccCheckAciDHCPOptionPolicyAttributes(description, &dhcp_option_policy),
				),
			},
			{
				Config: testAccCheckAciDHCPOptionPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPOptionPolicyExists("aci_dhcp_option_policy.foodhcp_option_policy", &dhcp_option_policy),
					testAccCheckAciDHCPOptionPolicyAttributes(description, &dhcp_option_policy),
				),
			},
		},
	})
}

func testAccCheckAciDHCPOptionPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_dhcp_option_policy" "foodhcp_option_policy" {
		#tenant_dn  = "${aci_tenant.example.id}"
		tenant_dn  = "uni/tn-check_context_tenant"
		description = "%s"
		name  = "example"
		annotation  = "example"
		name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciDHCPOptionPolicyExists(name string, dhcp_option_policy *models.DHCPOptionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("DHCP Option Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DHCP Option Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		dhcp_option_policyFound := models.DHCPOptionPolicyFromContainer(cont)
		if dhcp_option_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("DHCP Option Policy %s not found", rs.Primary.ID)
		}
		*dhcp_option_policy = *dhcp_option_policyFound
		return nil
	}
}

func testAccCheckAciDHCPOptionPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_dhcp_option_policy" {
			cont, err := client.Get(rs.Primary.ID)
			dhcp_option_policy := models.DHCPOptionPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("DHCP Option Policy %s Still exists", dhcp_option_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciDHCPOptionPolicyAttributes(description string, dhcp_option_policy *models.DHCPOptionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != dhcp_option_policy.Description {
			return fmt.Errorf("Bad dhcp_option_policy Description %s", dhcp_option_policy.Description)
		}

		if "example" != dhcp_option_policy.Name {
			return fmt.Errorf("Bad dhcp_option_policy name %s", dhcp_option_policy.Name)
		}

		if "example" != dhcp_option_policy.Annotation {
			return fmt.Errorf("Bad dhcp_option_policy annotation %s", dhcp_option_policy.Annotation)
		}

		if "example" != dhcp_option_policy.NameAlias {
			return fmt.Errorf("Bad dhcp_option_policy name_alias %s", dhcp_option_policy.NameAlias)
		}

		return nil
	}
}
