package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLLDPInterfacePolicy_Basic(t *testing.T) {
	var lldp_interface_policy models.LLDPInterfacePolicy
	description := "lldp_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLLDPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLLDPInterfacePolicyConfig_basic(description, "enabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLLDPInterfacePolicyExists("aci_lldp_interface_policy.foolldp_interface_policy", &lldp_interface_policy),
					testAccCheckAciLLDPInterfacePolicyAttributes(description, "enabled", &lldp_interface_policy),
				),
			},
		},
	})
}

func TestAccAciLLDPInterfacePolicy_update(t *testing.T) {
	var lldp_interface_policy models.LLDPInterfacePolicy
	description := "lldp_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLLDPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLLDPInterfacePolicyConfig_basic(description, "enabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLLDPInterfacePolicyExists("aci_lldp_interface_policy.foolldp_interface_policy", &lldp_interface_policy),
					testAccCheckAciLLDPInterfacePolicyAttributes(description, "enabled", &lldp_interface_policy),
				),
			},
			{
				Config: testAccCheckAciLLDPInterfacePolicyConfig_basic(description, "disabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLLDPInterfacePolicyExists("aci_lldp_interface_policy.foolldp_interface_policy", &lldp_interface_policy),
					testAccCheckAciLLDPInterfacePolicyAttributes(description, "disabled", &lldp_interface_policy),
				),
			},
		},
	})
}

func testAccCheckAciLLDPInterfacePolicyConfig_basic(description, admin_rx_st string) string {
	return fmt.Sprintf(`

	resource "aci_lldp_interface_policy" "foolldp_interface_policy" {
		description = "%s"
		name        = "demo_lldp_pol"
		admin_rx_st = "%s"
		admin_tx_st = "enabled"
		annotation  = "tag_lldp"
		name_alias  = "alias_lldp"
	}  
	`, description, admin_rx_st)
}

func testAccCheckAciLLDPInterfacePolicyExists(name string, lldp_interface_policy *models.LLDPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LLDP Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LLDP Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		lldp_interface_policyFound := models.LLDPInterfacePolicyFromContainer(cont)
		if lldp_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LLDP Interface Policy %s not found", rs.Primary.ID)
		}
		*lldp_interface_policy = *lldp_interface_policyFound
		return nil
	}
}

func testAccCheckAciLLDPInterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_lldp_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			lldp_interface_policy := models.LLDPInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LLDP Interface Policy %s Still exists", lldp_interface_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLLDPInterfacePolicyAttributes(description, admin_rx_st string, lldp_interface_policy *models.LLDPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != lldp_interface_policy.Description {
			return fmt.Errorf("Bad lldp_interface_policy Description %s", lldp_interface_policy.Description)
		}

		if "demo_lldp_pol" != lldp_interface_policy.Name {
			return fmt.Errorf("Bad lldp_interface_policy name %s", lldp_interface_policy.Name)
		}

		if admin_rx_st != lldp_interface_policy.AdminRxSt {
			return fmt.Errorf("Bad lldp_interface_policy admin_rx_st %s", lldp_interface_policy.AdminRxSt)
		}

		if "enabled" != lldp_interface_policy.AdminTxSt {
			return fmt.Errorf("Bad lldp_interface_policy admin_tx_st %s", lldp_interface_policy.AdminTxSt)
		}

		if "tag_lldp" != lldp_interface_policy.Annotation {
			return fmt.Errorf("Bad lldp_interface_policy annotation %s", lldp_interface_policy.Annotation)
		}

		if "alias_lldp" != lldp_interface_policy.NameAlias {
			return fmt.Errorf("Bad lldp_interface_policy name_alias %s", lldp_interface_policy.NameAlias)
		}

		return nil
	}
}
