package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciBGPAddressFamilyContextPolicy_Basic(t *testing.T) {
	var bgp_address_family_context_policy models.BGPAddressFamilyContextPolicy
	description := "bgp_address_family_context_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBGPAddressFamilyContextPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBGPAddressFamilyContextPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists("aci_bgp_address_family_context.test", &bgp_address_family_context_policy),
					testAccCheckAciBGPAddressFamilyContextPolicyAttributes(description, &bgp_address_family_context_policy),
				),
			},
		},
	})
}

func TestAccAciBGPAddressFamilyContextPolicy_update(t *testing.T) {
	var bgp_address_family_context_policy models.BGPAddressFamilyContextPolicy
	description := "bgp_address_family_context_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBGPAddressFamilyContextPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBGPAddressFamilyContextPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists("aci_bgp_address_family_context.test", &bgp_address_family_context_policy),
					testAccCheckAciBGPAddressFamilyContextPolicyAttributes(description, &bgp_address_family_context_policy),
				),
			},
			{
				Config: testAccCheckAciBGPAddressFamilyContextPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists("aci_bgp_address_family_context.test", &bgp_address_family_context_policy),
					testAccCheckAciBGPAddressFamilyContextPolicyAttributes(description, &bgp_address_family_context_policy),
				),
			},
		},
	})
}

func testAccCheckAciBGPAddressFamilyContextPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_bgp_address_family_context" "test" {
		tenant_dn     = "uni/tn-test"
		name          = "one"
		description   = "%s"
		annotation    = "example"
		ctrl          = "host-rt-leak"
		e_dist        = "25"
		i_dist        = "198"
		local_dist    = "100"
		max_ecmp      = "18"
		max_ecmp_ibgp = "25"
		name_alias    = "example"
	}
	`, description)
}

func testAccCheckAciBGPAddressFamilyContextPolicyExists(name string, bgp_address_family_context_policy *models.BGPAddressFamilyContextPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("BGP Address Family Context Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BGP Address Family Context Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bgp_address_family_context_policyFound := models.BGPAddressFamilyContextPolicyFromContainer(cont)
		if bgp_address_family_context_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("BGP Address Family Context Policy %s not found", rs.Primary.ID)
		}
		*bgp_address_family_context_policy = *bgp_address_family_context_policyFound
		return nil
	}
}

func testAccCheckAciBGPAddressFamilyContextPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_bgp_address_family_context" {
			cont, err := client.Get(rs.Primary.ID)
			bgp_address_family_context_policy := models.BGPAddressFamilyContextPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("BGP Address Family Context Policy %s Still exists", bgp_address_family_context_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciBGPAddressFamilyContextPolicyAttributes(description string, bgp_address_family_context_policy *models.BGPAddressFamilyContextPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != bgp_address_family_context_policy.Description {
			return fmt.Errorf("Bad bgp_address_family_context_policy Description %s", bgp_address_family_context_policy.Description)
		}

		if "one" != bgp_address_family_context_policy.Name {
			return fmt.Errorf("Bad bgp_address_family_context_policy name %s", bgp_address_family_context_policy.Name)
		}

		if "example" != bgp_address_family_context_policy.Annotation {
			return fmt.Errorf("Bad bgp_address_family_context_policy annotation %s", bgp_address_family_context_policy.Annotation)
		}

		if "host-rt-leak" != bgp_address_family_context_policy.Ctrl {
			return fmt.Errorf("Bad bgp_address_family_context_policy ctrl %s", bgp_address_family_context_policy.Ctrl)
		}

		if "25" != bgp_address_family_context_policy.EDist {
			return fmt.Errorf("Bad bgp_address_family_context_policy e_dist %s", bgp_address_family_context_policy.EDist)
		}

		if "198" != bgp_address_family_context_policy.IDist {
			return fmt.Errorf("Bad bgp_address_family_context_policy i_dist %s", bgp_address_family_context_policy.IDist)
		}

		if "100" != bgp_address_family_context_policy.LocalDist {
			return fmt.Errorf("Bad bgp_address_family_context_policy local_dist %s", bgp_address_family_context_policy.LocalDist)
		}

		if "18" != bgp_address_family_context_policy.MaxEcmp {
			return fmt.Errorf("Bad bgp_address_family_context_policy max_ecmp %s", bgp_address_family_context_policy.MaxEcmp)
		}

		if "25" != bgp_address_family_context_policy.MaxEcmpIbgp {
			return fmt.Errorf("Bad bgp_address_family_context_policy max_ecmp_ibgp %s", bgp_address_family_context_policy.MaxEcmpIbgp)
		}

		if "example" != bgp_address_family_context_policy.NameAlias {
			return fmt.Errorf("Bad bgp_address_family_context_policy name_alias %s", bgp_address_family_context_policy.NameAlias)
		}

		return nil
	}
}
