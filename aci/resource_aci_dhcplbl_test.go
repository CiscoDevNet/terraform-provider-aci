package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciDHCPRelayLabel_Basic(t *testing.T) {
	var dhcp_relay_label models.DHCPRelayLabel
	description := "dhcp_relay_label created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDHCPRelayLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDHCPRelayLabelConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPRelayLabelExists("aci_dhcp_relay_label.foodhcp_relay_label", &dhcp_relay_label),
					testAccCheckAciDHCPRelayLabelAttributes(description, &dhcp_relay_label),
				),
			},
		},
	})
}

func TestAccAciDHCPRelayLabel_update(t *testing.T) {
	var dhcp_relay_label models.DHCPRelayLabel
	description := "dhcp_relay_label created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDHCPRelayLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDHCPRelayLabelConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPRelayLabelExists("aci_dhcp_relay_label.foodhcp_relay_label", &dhcp_relay_label),
					testAccCheckAciDHCPRelayLabelAttributes(description, &dhcp_relay_label),
				),
			},
			{
				Config: testAccCheckAciDHCPRelayLabelConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPRelayLabelExists("aci_dhcp_relay_label.foodhcp_relay_label", &dhcp_relay_label),
					testAccCheckAciDHCPRelayLabelAttributes(description, &dhcp_relay_label),
				),
			},
		},
	})
}

func testAccCheckAciDHCPRelayLabelConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_dhcp_relay_label" "foodhcp_relay_label" {
		#bridge_domain_dn  = "${aci_bridge_domain.example.id}"
		bridge_domain_dn  = "uni/tn-check_tenantnk/BD-demo_bd"
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  name_alias  = "example"
		  owner  = "tenant"
		  tag  = "blanched-almond"
		}
	`, description)
}

func testAccCheckAciDHCPRelayLabelExists(name string, dhcp_relay_label *models.DHCPRelayLabel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("DHCP Relay Label %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DHCP Relay Label dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		dhcp_relay_labelFound := models.DHCPRelayLabelFromContainer(cont)
		if dhcp_relay_labelFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("DHCP Relay Label %s not found", rs.Primary.ID)
		}
		*dhcp_relay_label = *dhcp_relay_labelFound
		return nil
	}
}

func testAccCheckAciDHCPRelayLabelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_dhcp_relay_label" {
			cont, err := client.Get(rs.Primary.ID)
			dhcp_relay_label := models.DHCPRelayLabelFromContainer(cont)
			if err == nil {
				return fmt.Errorf("DHCP Relay Label %s Still exists", dhcp_relay_label.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciDHCPRelayLabelAttributes(description string, dhcp_relay_label *models.DHCPRelayLabel) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != dhcp_relay_label.Description {
			return fmt.Errorf("Bad dhcp_relay_label Description %s", dhcp_relay_label.Description)
		}

		if "example" != dhcp_relay_label.Name {
			return fmt.Errorf("Bad dhcp_relay_label name %s", dhcp_relay_label.Name)
		}

		if "example" != dhcp_relay_label.Annotation {
			return fmt.Errorf("Bad dhcp_relay_label annotation %s", dhcp_relay_label.Annotation)
		}

		if "example" != dhcp_relay_label.NameAlias {
			return fmt.Errorf("Bad dhcp_relay_label name_alias %s", dhcp_relay_label.NameAlias)
		}

		if "tenant" != dhcp_relay_label.Owner {
			return fmt.Errorf("Bad dhcp_relay_label owner %s", dhcp_relay_label.Owner)
		}

		if "blanched-almond" != dhcp_relay_label.Tag {
			return fmt.Errorf("Bad dhcp_relay_label tag %s", dhcp_relay_label.Tag)
		}

		return nil
	}
}
