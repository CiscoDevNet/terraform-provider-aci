package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciInterfaceFCPolicy_Basic(t *testing.T) {
	var interface_fc_policy models.InterfaceFCPolicy
	description := "interface_fc_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInterfaceFCPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInterfaceFCPolicyConfig_basic(description, "64"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists("aci_interface_fc_policy.foointerface_fc_policy", &interface_fc_policy),
					testAccCheckAciInterfaceFCPolicyAttributes(description, "64", &interface_fc_policy),
				),
			},
			{
				ResourceName:      "aci_interface_fc_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciInterfaceFCPolicy_update(t *testing.T) {
	var interface_fc_policy models.InterfaceFCPolicy
	description := "interface_fc_policy created while acceptance testing"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInterfaceFCPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInterfaceFCPolicyConfig_basic(description, "64"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists("aci_interface_fc_policy.foointerface_fc_policy", &interface_fc_policy),
					testAccCheckAciInterfaceFCPolicyAttributes(description, "64", &interface_fc_policy),
				),
			},
			{
				Config: testAccCheckAciInterfaceFCPolicyConfig_basic(description, "70"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceFCPolicyExists("aci_interface_fc_policy.foointerface_fc_policy", &interface_fc_policy),
					testAccCheckAciInterfaceFCPolicyAttributes(description, "70", &interface_fc_policy),
				),
			},
		},
	})
}

func testAccCheckAciInterfaceFCPolicyConfig_basic(description, rx_bb_credit string) string {
	return fmt.Sprintf(`

	resource "aci_interface_fc_policy" "foointerface_fc_policy" {	
		name          = "demo_policy"
		description   = "%s"
		annotation    = "tag_if_policy"
		fill_pattern  = "default"
		name_alias    = "demo_alias"
		port_mode     = "f"
		rx_bb_credit  = "%s"
		speed         = "auto"
		trunk_mode    = "auto"
	}
	`, description, rx_bb_credit)
}

func testAccCheckAciInterfaceFCPolicyExists(name string, interface_fc_policy *models.InterfaceFCPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Interface FC Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Interface FC Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		interface_fc_policyFound := models.InterfaceFCPolicyFromContainer(cont)
		if interface_fc_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Interface FC Policy %s not found", rs.Primary.ID)
		}
		*interface_fc_policy = *interface_fc_policyFound
		return nil
	}
}

func testAccCheckAciInterfaceFCPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_interface_fc_policy" {
			cont, err := client.Get(rs.Primary.ID)
			interface_fc_policy := models.InterfaceFCPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Interface FC Policy %s Still exists", interface_fc_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciInterfaceFCPolicyAttributes(description, rx_bb_credit string, interface_fc_policy *models.InterfaceFCPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != interface_fc_policy.Description {
			return fmt.Errorf("Bad interface_fc_policy Description %s", interface_fc_policy.Description)
		}

		if "demo_policy" != interface_fc_policy.Name {
			return fmt.Errorf("Bad interface_fc_policy name %s", interface_fc_policy.Name)
		}

		if "tag_if_policy" != interface_fc_policy.Annotation {
			return fmt.Errorf("Bad interface_fc_policy annotation %s", interface_fc_policy.Annotation)
		}

		if "default" != interface_fc_policy.FillPattern {
			return fmt.Errorf("Bad interface_fc_policy fill_pattern %s", interface_fc_policy.FillPattern)
		}

		if "demo_alias" != interface_fc_policy.NameAlias {
			return fmt.Errorf("Bad interface_fc_policy name_alias %s", interface_fc_policy.NameAlias)
		}

		if "f" != interface_fc_policy.PortMode {
			return fmt.Errorf("Bad interface_fc_policy port_mode %s", interface_fc_policy.PortMode)
		}

		if rx_bb_credit != interface_fc_policy.RxBBCredit {
			return fmt.Errorf("Bad interface_fc_policy rx_bb_credit %s", interface_fc_policy.RxBBCredit)
		}

		if "auto" != interface_fc_policy.Speed {
			return fmt.Errorf("Bad interface_fc_policy speed %s", interface_fc_policy.Speed)
		}

		if "auto" != interface_fc_policy.TrunkMode {
			return fmt.Errorf("Bad interface_fc_policy trunk_mode %s", interface_fc_policy.TrunkMode)
		}
		return nil
	}
}
