package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBDDHCPLabel_Basic(t *testing.T) {
	var bd_dhcp_label models.BDDHCPLabel
	description := "bd_dhcp_label created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBDDHCPLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBDDHCPLabelConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBDDHCPLabelExists("aci_bd_dhcp_label.foobd_dhcp_label", &bd_dhcp_label),
					testAccCheckAciBDDHCPLabelAttributes(description, &bd_dhcp_label),
				),
			},
		},
	})
}

func TestAccAciBDDHCPLabel_update(t *testing.T) {
	var bd_dhcp_label models.BDDHCPLabel
	description := "bd_dhcp_label created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBDDHCPLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBDDHCPLabelConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBDDHCPLabelExists("aci_bd_dhcp_label.foobd_dhcp_label", &bd_dhcp_label),
					testAccCheckAciBDDHCPLabelAttributes(description, &bd_dhcp_label),
				),
			},
			{
				Config: testAccCheckAciBDDHCPLabelConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBDDHCPLabelExists("aci_bd_dhcp_label.foobd_dhcp_label", &bd_dhcp_label),
					testAccCheckAciBDDHCPLabelAttributes(description, &bd_dhcp_label),
				),
			},
		},
	})
}

func testAccCheckAciBDDHCPLabelConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_bd_dhcp_label" "foobd_dhcp_label" {
		#bridge_domain_dn = aci_bridge_domain.example.id
		bridge_domain_dn  = "uni/tn-check_tenantnk/BD-demo_bd"
		description = "%s"
		name        = "example"
		annotation  = "example"
		name_alias  = "example"
		owner       = "tenant"
		tag         = "blanched-almond"
	}
	`, description)
}

func testAccCheckAciBDDHCPLabelExists(name string, bd_dhcp_label *models.BDDHCPLabel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("BD DHCP Label %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BD DHCP Label dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bd_dhcp_labelFound := models.BDDHCPLabelFromContainer(cont)
		if bd_dhcp_labelFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("BD DHCP Label %s not found", rs.Primary.ID)
		}
		*bd_dhcp_label = *bd_dhcp_labelFound
		return nil
	}
}

func testAccCheckAciBDDHCPLabelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_bd_dhcp_label" {
			cont, err := client.Get(rs.Primary.ID)
			bd_dhcp_label := models.BDDHCPLabelFromContainer(cont)
			if err == nil {
				return fmt.Errorf("BD DHCP Label %s Still exists", bd_dhcp_label.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciBDDHCPLabelAttributes(description string, bd_dhcp_label *models.BDDHCPLabel) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != bd_dhcp_label.Description {
			return fmt.Errorf("Bad bd_dhcp_label Description %s", bd_dhcp_label.Description)
		}

		if "example" != bd_dhcp_label.Name {
			return fmt.Errorf("Bad bd_dhcp_label name %s", bd_dhcp_label.Name)
		}

		if "example" != bd_dhcp_label.Annotation {
			return fmt.Errorf("Bad bd_dhcp_label annotation %s", bd_dhcp_label.Annotation)
		}

		if "example" != bd_dhcp_label.NameAlias {
			return fmt.Errorf("Bad bd_dhcp_label name_alias %s", bd_dhcp_label.NameAlias)
		}

		if "tenant" != bd_dhcp_label.Owner {
			return fmt.Errorf("Bad bd_dhcp_label owner %s", bd_dhcp_label.Owner)
		}

		if "blanched-almond" != bd_dhcp_label.Tag {
			return fmt.Errorf("Bad bd_dhcp_label tag %s", bd_dhcp_label.Tag)
		}

		return nil
	}
}
