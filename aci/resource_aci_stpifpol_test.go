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

func TestAccAciSpanningTreeInterfacePolicy_Basic(t *testing.T) {
	var spanning_tree_interface_policy models.SpanningTreeInterfacePolicy
	stp_if_pol_name := acctest.RandString(5)
	description := "spanning_tree_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpanningTreeInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpanningTreeInterfacePolicyConfig_basic(stp_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpanningTreeInterfacePolicyExists("aci_spanning_tree_interface_policy.foospanning_tree_interface_policy", &spanning_tree_interface_policy),
					testAccCheckAciSpanningTreeInterfacePolicyAttributes(stp_if_pol_name, description, &spanning_tree_interface_policy),
				),
			},
		},
	})
}

func testAccCheckAciSpanningTreeInterfacePolicyConfig_basic(stp_if_pol_name string) string {
	return fmt.Sprintf(`

	resource "aci_spanning_tree_interface_policy" "foospanning_tree_interface_policy" {
		name 		= "%s"
		description = "spanning_tree_interface_policy created while acceptance testing"

	}

	`, stp_if_pol_name)
}

func testAccCheckAciSpanningTreeInterfacePolicyExists(name string, spanning_tree_interface_policy *models.SpanningTreeInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spanning Tree Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spanning Tree Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spanning_tree_interface_policyFound := models.SpanningTreeInterfacePolicyFromContainer(cont)
		if spanning_tree_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spanning Tree Interface Policy %s not found", rs.Primary.ID)
		}
		*spanning_tree_interface_policy = *spanning_tree_interface_policyFound
		return nil
	}
}

func testAccCheckAciSpanningTreeInterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_spanning_tree_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			spanning_tree_interface_policy := models.SpanningTreeInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spanning Tree Interface Policy %s Still exists", spanning_tree_interface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSpanningTreeInterfacePolicyAttributes(stp_if_pol_name, description string, spanning_tree_interface_policy *models.SpanningTreeInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if stp_if_pol_name != GetMOName(spanning_tree_interface_policy.DistinguishedName) {
			return fmt.Errorf("Bad stp_if_pol %s", GetMOName(spanning_tree_interface_policy.DistinguishedName))
		}

		if description != spanning_tree_interface_policy.Description {
			return fmt.Errorf("Bad spanning_tree_interface_policy Description %s", spanning_tree_interface_policy.Description)
		}
		return nil
	}
}
