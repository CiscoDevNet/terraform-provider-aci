package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBGPPeerPrefixPolicy_Basic(t *testing.T) {
	var bgp_peer_prefix_policy models.BGPPeerPrefixPolicy
	description := "bgp_peer_prefix_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBGPPeerPrefixPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBGPPeerPrefixPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPPeerPrefixPolicyExists("aci_bgp_peer_prefix.test", &bgp_peer_prefix_policy),
					testAccCheckAciBGPPeerPrefixPolicyAttributes(description, &bgp_peer_prefix_policy),
				),
			},
		},
	})
}

func TestAccAciBGPPeerPrefixPolicy_update(t *testing.T) {
	var bgp_peer_prefix_policy models.BGPPeerPrefixPolicy
	description := "bgp_peer_prefix_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBGPPeerPrefixPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBGPPeerPrefixPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPPeerPrefixPolicyExists("aci_bgp_peer_prefix.test", &bgp_peer_prefix_policy),
					testAccCheckAciBGPPeerPrefixPolicyAttributes(description, &bgp_peer_prefix_policy),
				),
			},
			{
				Config: testAccCheckAciBGPPeerPrefixPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPPeerPrefixPolicyExists("aci_bgp_peer_prefix.test", &bgp_peer_prefix_policy),
					testAccCheckAciBGPPeerPrefixPolicyAttributes(description, &bgp_peer_prefix_policy),
				),
			},
		},
	})
}

func testAccCheckAciBGPPeerPrefixPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_bgp_peer_prefix" "test" {
		tenant_dn    = aci_tenant.demo_dev_tenant_test.id
		name         = "one"
		description  = "%s"
		action       = "shut"
		annotation   = "example"
		max_pfx      = "200"
		name_alias   = "example"
		restart_time = "200"
		thresh       = "85"
	}
	`, description)
}

func testAccCheckAciBGPPeerPrefixPolicyExists(name string, bgp_peer_prefix_policy *models.BGPPeerPrefixPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("BGP Peer Prefix Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BGP Peer Prefix Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bgp_peer_prefix_policyFound := models.BGPPeerPrefixPolicyFromContainer(cont)
		if bgp_peer_prefix_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("BGP Peer Prefix Policy %s not found", rs.Primary.ID)
		}
		*bgp_peer_prefix_policy = *bgp_peer_prefix_policyFound
		return nil
	}
}

func testAccCheckAciBGPPeerPrefixPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_bgp_peer_prefix" {
			cont, err := client.Get(rs.Primary.ID)
			bgp_peer_prefix_policy := models.BGPPeerPrefixPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("BGP Peer Prefix Policy %s Still exists", bgp_peer_prefix_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciBGPPeerPrefixPolicyAttributes(description string, bgp_peer_prefix_policy *models.BGPPeerPrefixPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != bgp_peer_prefix_policy.Description {
			return fmt.Errorf("Bad bgp_peer_prefix_policy Description %s", bgp_peer_prefix_policy.Description)
		}

		if "one" != bgp_peer_prefix_policy.Name {
			return fmt.Errorf("Bad bgp_peer_prefix_policy name %s", bgp_peer_prefix_policy.Name)
		}

		if "shut" != bgp_peer_prefix_policy.Action {
			return fmt.Errorf("Bad bgp_peer_prefix_policy action %s", bgp_peer_prefix_policy.Action)
		}

		if "example" != bgp_peer_prefix_policy.Annotation {
			return fmt.Errorf("Bad bgp_peer_prefix_policy annotation %s", bgp_peer_prefix_policy.Annotation)
		}

		if "200" != bgp_peer_prefix_policy.MaxPfx {
			return fmt.Errorf("Bad bgp_peer_prefix_policy max_pfx %s", bgp_peer_prefix_policy.MaxPfx)
		}

		if "example" != bgp_peer_prefix_policy.NameAlias {
			return fmt.Errorf("Bad bgp_peer_prefix_policy name_alias %s", bgp_peer_prefix_policy.NameAlias)
		}

		if "200" != bgp_peer_prefix_policy.RestartTime {
			return fmt.Errorf("Bad bgp_peer_prefix_policy restart_time %s", bgp_peer_prefix_policy.RestartTime)
		}

		if "85" != bgp_peer_prefix_policy.Thresh {
			return fmt.Errorf("Bad bgp_peer_prefix_policy thresh %s", bgp_peer_prefix_policy.Thresh)
		}

		return nil
	}
}
