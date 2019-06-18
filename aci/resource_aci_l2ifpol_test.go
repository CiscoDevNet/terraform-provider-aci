package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciL2InterfacePolicy_Basic(t *testing.T) {
	var l2_interface_policy models.L2InterfacePolicy
	description := "l2_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL2InterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL2InterfacePolicyConfig_basic(description, "global"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2InterfacePolicyExists("aci_l2_interface_policy.fool2_interface_policy", &l2_interface_policy),
					testAccCheckAciL2InterfacePolicyAttributes(description, "global", &l2_interface_policy),
				),
			},
			{
				ResourceName:      "aci_l2_interface_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciL2InterfacePolicy_update(t *testing.T) {
	var l2_interface_policy models.L2InterfacePolicy
	description := "l2_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL2InterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL2InterfacePolicyConfig_basic(description, "global"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2InterfacePolicyExists("aci_l2_interface_policy.fool2_interface_policy", &l2_interface_policy),
					testAccCheckAciL2InterfacePolicyAttributes(description, "global", &l2_interface_policy),
				),
			},
			{
				Config: testAccCheckAciL2InterfacePolicyConfig_basic(description, "portlocal"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2InterfacePolicyExists("aci_l2_interface_policy.fool2_interface_policy", &l2_interface_policy),
					testAccCheckAciL2InterfacePolicyAttributes(description, "portlocal", &l2_interface_policy),
				),
			},
		},
	})
}

func testAccCheckAciL2InterfacePolicyConfig_basic(description, vlan_scope string) string {
	return fmt.Sprintf(`

	resource "aci_l2_interface_policy" "fool2_interface_policy" {
		description = "%s"
		name        = "demo_l2_pol"
		annotation  = "tag_l2_pol"
		name_alias  = "alias_l2_pol"
		qinq        = "disabled"
		vepa        = "disabled"
		vlan_scope  = "%s"
	}
	  
	`, description, vlan_scope)
}

func testAccCheckAciL2InterfacePolicyExists(name string, l2_interface_policy *models.L2InterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L2 Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L2 Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l2_interface_policyFound := models.L2InterfacePolicyFromContainer(cont)
		if l2_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L2 Interface Policy %s not found", rs.Primary.ID)
		}
		*l2_interface_policy = *l2_interface_policyFound
		return nil
	}
}

func testAccCheckAciL2InterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l2_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			l2_interface_policy := models.L2InterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L2 Interface Policy %s Still exists", l2_interface_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL2InterfacePolicyAttributes(description, vlan_scope string, l2_interface_policy *models.L2InterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l2_interface_policy.Description {
			return fmt.Errorf("Bad l2_interface_policy Description %s", l2_interface_policy.Description)
		}

		if "demo_l2_pol" != l2_interface_policy.Name {
			return fmt.Errorf("Bad l2_interface_policy name %s", l2_interface_policy.Name)
		}

		if "tag_l2_pol" != l2_interface_policy.Annotation {
			return fmt.Errorf("Bad l2_interface_policy annotation %s", l2_interface_policy.Annotation)
		}

		if "alias_l2_pol" != l2_interface_policy.NameAlias {
			return fmt.Errorf("Bad l2_interface_policy name_alias %s", l2_interface_policy.NameAlias)
		}

		if "disabled" != l2_interface_policy.Qinq {
			return fmt.Errorf("Bad l2_interface_policy qinq %s", l2_interface_policy.Qinq)
		}

		if "disabled" != l2_interface_policy.Vepa {
			return fmt.Errorf("Bad l2_interface_policy vepa %s", l2_interface_policy.Vepa)
		}

		if vlan_scope != l2_interface_policy.VlanScope {
			return fmt.Errorf("Bad l2_interface_policy vlan_scope %s", l2_interface_policy.VlanScope)
		}

		return nil
	}
}
