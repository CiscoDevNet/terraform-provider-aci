package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLACPMemberPolicy_Basic(t *testing.T) {
	var lacpmember_policy models.LACPMemberPolicy
	lacp_if_pol_name := acctest.RandString(5)
	description := "aci_lacp_member_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLACPMemberPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLACPMemberPolicyConfig_basic(lacp_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPMemberPolicyExists("aci_lacp_member_policy.foo_lacpmember_policy", &lacpmember_policy),
					testAccCheckAciLACPMemberPolicyAttributes(lacp_if_pol_name, description, &lacpmember_policy),
				),
			},
		},
	})
}

func TestAccAciLACPMemberPolicy_Update(t *testing.T) {
	var lacpmember_policy models.LACPMemberPolicy
	lacp_if_pol_name := acctest.RandString(5)
	description := "aci_lacp_member_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLACPMemberPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLACPMemberPolicyConfig_basic(lacp_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPMemberPolicyExists("aci_lacp_member_policy.foo_lacpmember_policy", &lacpmember_policy),
					testAccCheckAciLACPMemberPolicyAttributes(lacp_if_pol_name, description, &lacpmember_policy),
				),
			},
			{
				Config: testAccCheckAciLACPMemberPolicyConfig_basic(lacp_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPMemberPolicyExists("aci_lacp_member_policy.foo_lacpmember_policy", &lacpmember_policy),
					testAccCheckAciLACPMemberPolicyAttributes(lacp_if_pol_name, description, &lacpmember_policy),
				),
			},
		},
	})
}

func testAccCheckAciLACPMemberPolicyConfig_basic(lacp_if_pol_name string) string {
	return fmt.Sprintf(`

	resource "aci_lacp_member_policy" "foo_lacpmember_policy" {
		name 		= "%s"
		description = "aci_lacp_member_policy created while acceptance testing"
	}
	`, lacp_if_pol_name)
}

func testAccCheckAciLACPMemberPolicyExists(name string, lacpmember_policy *models.LACPMemberPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LACP Member Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LACP Member Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		lacpmember_policyFound := models.LACPMemberPolicyFromContainer(cont)
		if lacpmember_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LACP Member Policy %s not found", rs.Primary.ID)
		}
		*lacpmember_policy = *lacpmember_policyFound
		return nil
	}
}

func testAccCheckAciLACPMemberPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_lacp_member_policy" {
			cont, err := client.Get(rs.Primary.ID)
			lacpmember_policy := models.LACPMemberPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LACP Member Policy %s Still exists", lacpmember_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLACPMemberPolicyAttributes(lacp_if_pol_name, description string, lacpmember_policy *models.LACPMemberPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if lacp_if_pol_name != GetMOName(lacpmember_policy.DistinguishedName) {
			return fmt.Errorf("Bad lacpif_pol %s", GetMOName(lacpmember_policy.DistinguishedName))
		}

		if description != lacpmember_policy.Description {
			return fmt.Errorf("Bad lacpmember_policy Description %s", lacpmember_policy.Description)
		}
		return nil
	}
}
