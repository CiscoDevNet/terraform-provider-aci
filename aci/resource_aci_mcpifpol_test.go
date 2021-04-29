package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciMiscablingProtocolInterfacePolicy_Basic(t *testing.T) {
	var miscabling_protocol_interface_policy models.MiscablingProtocolInterfacePolicy
	description := "mis-cabling_protocol_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMiscablingProtocolInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMiscablingProtocolInterfacePolicyConfig_basic(description, "enabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInterfacePolicyExists("aci_mis-cabling_protocol_interface_policy.foomis-cabling_protocol_interface_policy", &miscabling_protocol_interface_policy),
					testAccCheckAciMiscablingProtocolInterfacePolicyAttributes(description, "enabled", &miscabling_protocol_interface_policy),
				),
			},
			{
				ResourceName:      "aci_miscabling_protocol_interface_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciMiscablingProtocolInterfacePolicy_update(t *testing.T) {
	var miscabling_protocol_interface_policy models.MiscablingProtocolInterfacePolicy
	description := "mis-cabling_protocol_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMiscablingProtocolInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMiscablingProtocolInterfacePolicyConfig_basic(description, "enabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInterfacePolicyExists("aci_mis-cabling_protocol_interface_policy.foomis-cabling_protocol_interface_policy", &miscabling_protocol_interface_policy),
					testAccCheckAciMiscablingProtocolInterfacePolicyAttributes(description, "enabled", &miscabling_protocol_interface_policy),
				),
			},
			{
				Config: testAccCheckAciMiscablingProtocolInterfacePolicyConfig_basic(description, "disabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMiscablingProtocolInterfacePolicyExists("aci_mis-cabling_protocol_interface_policy.foomis-cabling_protocol_interface_policy", &miscabling_protocol_interface_policy),
					testAccCheckAciMiscablingProtocolInterfacePolicyAttributes(description, "disabled", &miscabling_protocol_interface_policy),
				),
			},
		},
	})
}

func testAccCheckAciMiscablingProtocolInterfacePolicyConfig_basic(description, admin_st string) string {
	return fmt.Sprintf(`

	resource "aci_miscabling_protocol_interface_policy" "foomiscabling_protocol_interface_policy" {
		description = "%s"
		name        = "demo_mcpol"
		admin_st    = "%s"
		annotation  = "tag_mcpol"
		name_alias  = "alias_mcpol"
	}  
	`, description, admin_st)
}

func testAccCheckAciMiscablingProtocolInterfacePolicyExists(name string, miscabling_protocol_interface_policy *models.MiscablingProtocolInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Mis-cabling Protocol Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Mis-cabling Protocol Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		miscabling_protocol_interface_policyFound := models.MiscablingProtocolInterfacePolicyFromContainer(cont)
		if miscabling_protocol_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Mis-cabling Protocol Interface Policy %s not found", rs.Primary.ID)
		}
		*miscabling_protocol_interface_policy = *miscabling_protocol_interface_policyFound
		return nil
	}
}

func testAccCheckAciMiscablingProtocolInterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_mis-cabling_protocol_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			miscabling_protocol_interface_policy := models.MiscablingProtocolInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Mis-cabling Protocol Interface Policy %s Still exists", miscabling_protocol_interface_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciMiscablingProtocolInterfacePolicyAttributes(description, admin_st string, miscabling_protocol_interface_policy *models.MiscablingProtocolInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != miscabling_protocol_interface_policy.Description {
			return fmt.Errorf("Bad miscabling_protocol_interface_policy Description %s", miscabling_protocol_interface_policy.Description)
		}

		if "demo_mcpol" != miscabling_protocol_interface_policy.Name {
			return fmt.Errorf("Bad miscabling_protocol_interface_policy name %s", miscabling_protocol_interface_policy.Name)
		}

		if admin_st != miscabling_protocol_interface_policy.AdminSt {
			return fmt.Errorf("Bad miscabling_protocol_interface_policy admin_st %s", miscabling_protocol_interface_policy.AdminSt)
		}

		if "tag_mcpol" != miscabling_protocol_interface_policy.Annotation {
			return fmt.Errorf("Bad miscabling_protocol_interface_policy annotation %s", miscabling_protocol_interface_policy.Annotation)
		}

		if "alias_mcpol" != miscabling_protocol_interface_policy.NameAlias {
			return fmt.Errorf("Bad miscabling_protocol_interface_policy name_alias %s", miscabling_protocol_interface_policy.NameAlias)
		}

		return nil
	}
}
