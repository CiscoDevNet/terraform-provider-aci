package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBGPAddressFamilyContextPolicy_Basic(t *testing.T) {
	var bgp_address_family_context_policy models.BGPAddressFamilyContextPolicy
	fv_tenant_name := acctest.RandString(5)
	fv_ctx_name := acctest.RandString(5)
	fv_rs_ctx_to_bgp_ctx_af_pol_name := acctest.RandString(5)
	description := "bgp_address_family_context_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBGPAddressFamilyContextPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBGPAddressFamilyContextPolicyConfig_basic(fv_tenant_name, fv_ctx_name, fv_rs_ctx_to_bgp_ctx_af_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists("aci_bgp_address_family_context_policy.foobgp_address_family_context_policy", &bgp_address_family_context_policy),
					testAccCheckAciBGPAddressFamilyContextPolicyAttributes(fv_tenant_name, fv_ctx_name, fv_rs_ctx_to_bgp_ctx_af_pol_name, description, &bgp_address_family_context_policy),
				),
			},
		},
	})
}

func TestAccAciBGPAddressFamilyContextPolicy_Update(t *testing.T) {
	var bgp_address_family_context_policy models.BGPAddressFamilyContextPolicy
	fv_tenant_name := acctest.RandString(5)
	fv_ctx_name := acctest.RandString(5)
	fv_rs_ctx_to_bgp_ctx_af_pol_name := acctest.RandString(5)
	description := "bgp_address_family_context_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBGPAddressFamilyContextPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBGPAddressFamilyContextPolicyConfig_basic(fv_tenant_name, fv_ctx_name, fv_rs_ctx_to_bgp_ctx_af_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists("aci_bgp_address_family_context_policy.foobgp_address_family_context_policy", &bgp_address_family_context_policy),
					testAccCheckAciBGPAddressFamilyContextPolicyAttributes(fv_tenant_name, fv_ctx_name, fv_rs_ctx_to_bgp_ctx_af_pol_name, description, &bgp_address_family_context_policy),
				),
			},
			{
				Config: testAccCheckAciBGPAddressFamilyContextPolicyConfig_basic(fv_tenant_name, fv_ctx_name, fv_rs_ctx_to_bgp_ctx_af_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBGPAddressFamilyContextPolicyExists("aci_bgp_address_family_context_policy.foobgp_address_family_context_policy", &bgp_address_family_context_policy),
					testAccCheckAciBGPAddressFamilyContextPolicyAttributes(fv_tenant_name, fv_ctx_name, fv_rs_ctx_to_bgp_ctx_af_pol_name, description, &bgp_address_family_context_policy),
				),
			},
		},
	})
}

func testAccCheckAciBGPAddressFamilyContextPolicyConfig_basic(fv_tenant_name, fv_ctx_name, fv_rs_ctx_to_bgp_ctx_af_pol_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_vrf" "foovrf" {
		name 		= "%s"
		description = "vrf created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	resource "aci_bgp_address_family_context_policy" "foobgp_address_family_context_policy" {
		name 		= "%s"
		description = "bgp_address_family_context_policy created while acceptance testing"
		vrf_dn = aci_vrf.foovrf.id
	}

	`, fv_tenant_name, fv_ctx_name, fv_rs_ctx_to_bgp_ctx_af_pol_name)
}

func testAccCheckAciBGPAddressFamilyContextPolicyExists(name string, bgp_address_family_context_policy *models.BGPAddressFamilyContextPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("BGP  Address Family Context Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BGP  Address Family Context Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bgp_address_family_context_policyFound := models.BGPAddressFamilyContextPolicyFromContainer(cont)
		if bgp_address_family_context_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("BGP  Address Family Context Policy %s not found", rs.Primary.ID)
		}
		*bgp_address_family_context_policy = *bgp_address_family_context_policyFound
		return nil
	}
}

func testAccCheckAciBGPAddressFamilyContextPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bgp_address_family_context_policy" {
			cont, err := client.Get(rs.Primary.ID)
			bgp_address_family_context_policy := models.BGPAddressFamilyContextPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("BGP  Address Family Context Policy %s Still exists", bgp_address_family_context_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBGPAddressFamilyContextPolicyAttributes(fv_tenant_name, fv_ctx_name, fv_rs_ctx_to_bgp_ctx_af_pol_name, description string, bgp_address_family_context_policy *models.BGPAddressFamilyContextPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fv_rs_ctx_to_bgp_ctx_af_pol_name != GetMOName(bgp_address_family_context_policy.DistinguishedName) {
			return fmt.Errorf("Bad fv_rs_ctx_to_bgp_ctx_af_pol %s", GetMOName(bgp_address_family_context_policy.DistinguishedName))
		}

		if fv_ctx_name != GetMOName(GetParentDn(bgp_address_family_context_policy.DistinguishedName)) {
			return fmt.Errorf(" Bad fv_ctx %s", GetMOName(GetParentDn(bgp_address_family_context_policy.DistinguishedName)))
		}
		if description != bgp_address_family_context_policy.Description {
			return fmt.Errorf("Bad bgp_address_family_context_policy Description %s", bgp_address_family_context_policy.Description)
		}
		return nil
	}
}
