package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciDHCPRelayPolicy_Basic(t *testing.T) {
	var dhcp_relay_policy models.DHCPRelayPolicy
	description := "dhcp_relay_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDHCPRelayPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDHCPRelayPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPRelayPolicyExists("aci_dhcp_relay_policy.foodhcp_relay_policy", &dhcp_relay_policy),
					testAccCheckAciDHCPRelayPolicyAttributes(description, &dhcp_relay_policy),
				),
			},
		},
	})
}

func TestAccAciDHCPRelayPolicy_update(t *testing.T) {
	var dhcp_relay_policy models.DHCPRelayPolicy
	description := "dhcp_relay_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDHCPRelayPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDHCPRelayPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPRelayPolicyExists("aci_dhcp_relay_policy.foodhcp_relay_policy", &dhcp_relay_policy),
					testAccCheckAciDHCPRelayPolicyAttributes(description, &dhcp_relay_policy),
				),
			},
			{
				Config: testAccCheckAciDHCPRelayPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPRelayPolicyExists("aci_dhcp_relay_policy.foodhcp_relay_policy", &dhcp_relay_policy),
					testAccCheckAciDHCPRelayPolicyAttributes(description, &dhcp_relay_policy),
				),
			},
		},
	})
}

func testAccCheckAciDHCPRelayPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_dhcp_relay_policy" "foodhcp_relay_policy" {
		tenant_dn  = "uni/tn-crest_test_dhrumil_tenant" 
		description = "%s"
		
		name  = "name_example"
		  annotation  = "annotation_example"
		  mode  = "visible"
		  name_alias  = "alias_example"
		  owner  = "infra"
		}
	`, description)
}

func testAccCheckAciDHCPRelayPolicyExists(name string, dhcp_relay_policy *models.DHCPRelayPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("DHCP Relay Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DHCP Relay Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		dhcp_relay_policyFound := models.DHCPRelayPolicyFromContainer(cont)
		if dhcp_relay_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("DHCP Relay Policy %s not found", rs.Primary.ID)
		}
		*dhcp_relay_policy = *dhcp_relay_policyFound
		return nil
	}
}

func testAccCheckAciDHCPRelayPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_dhcp_relay_policy" {
			cont, err := client.Get(rs.Primary.ID)
			dhcp_relay_policy := models.DHCPRelayPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("DHCP Relay Policy %s Still exists", dhcp_relay_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciDHCPRelayPolicyAttributes(description string, dhcp_relay_policy *models.DHCPRelayPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != dhcp_relay_policy.Description {
			return fmt.Errorf("Bad dhcp_relay_policy Description %s", dhcp_relay_policy.Description)
		}

		if "name_example" != dhcp_relay_policy.Name {
			return fmt.Errorf("Bad dhcp_relay_policy name %s", dhcp_relay_policy.Name)
		}

		if "annotation_example" != dhcp_relay_policy.Annotation {
			return fmt.Errorf("Bad dhcp_relay_policy annotation %s", dhcp_relay_policy.Annotation)
		}

		if "visible" != dhcp_relay_policy.Mode {
			return fmt.Errorf("Bad dhcp_relay_policy mode %s", dhcp_relay_policy.Mode)
		}

		if "alias_example" != dhcp_relay_policy.NameAlias {
			return fmt.Errorf("Bad dhcp_relay_policy name_alias %s", dhcp_relay_policy.NameAlias)
		}

		if "infra" != dhcp_relay_policy.Owner {
			return fmt.Errorf("Bad dhcp_relay_policy owner %s", dhcp_relay_policy.Owner)
		}

		return nil
	}
}
