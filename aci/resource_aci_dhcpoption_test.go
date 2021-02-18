package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciDHCPOption_Basic(t *testing.T) {
	var dhcp_option models.DHCPOption
	//description := "dhcp_option created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDHCPOptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDHCPOptionConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPOptionExists("aci_dhcp_option.foodhcp_option", &dhcp_option),
					testAccCheckAciDHCPOptionAttributes(&dhcp_option),
				),
			},
		},
	})
}

func TestAccAciDHCPOption_update(t *testing.T) {
	var dhcp_option models.DHCPOption
	//description := "dhcp_option created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciDHCPOptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciDHCPOptionConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPOptionExists("aci_dhcp_option.foodhcp_option", &dhcp_option),
					testAccCheckAciDHCPOptionAttributes(&dhcp_option),
				),
			},
			{
				Config: testAccCheckAciDHCPOptionConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciDHCPOptionExists("aci_dhcp_option.foodhcp_option", &dhcp_option),
					testAccCheckAciDHCPOptionAttributes(&dhcp_option),
				),
			},
		},
	})
}

func testAccCheckAciDHCPOptionConfig_basic() string {
	return fmt.Sprintf(`

	resource "aci_dhcp_option" "foodhcp_option" {
		#dhcp_option_policy_dn  = "${aci_dhcp_option_policy.example.id}"
		dhcp_option_policy_dn  = "uni/tn-check_context_tenant/dhcpoptpol-foodhcp_option_policy"
		name  = "example"
		annotation  = "example"
		data  = "example"
		dhcp_option_id  = "33"
		name_alias  = "example"
		}
	`)
}

func testAccCheckAciDHCPOptionExists(name string, dhcp_option *models.DHCPOption) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("DHCP Option %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DHCP Option dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		dhcp_optionFound := models.DHCPOptionFromContainer(cont)
		if dhcp_optionFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("DHCP Option %s not found", rs.Primary.ID)
		}
		*dhcp_option = *dhcp_optionFound
		return nil
	}
}

func testAccCheckAciDHCPOptionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_dhcp_option" {
			cont, err := client.Get(rs.Primary.ID)
			dhcp_option := models.DHCPOptionFromContainer(cont)
			if err == nil {
				return fmt.Errorf("DHCP Option %s Still exists", dhcp_option.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciDHCPOptionAttributes(dhcp_option *models.DHCPOption) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		// if description != dhcp_option.Description {
		// 	return fmt.Errorf("Bad dhcp_option Description %s", dhcp_option.Description)
		// }

		if "example" != dhcp_option.Name {
			return fmt.Errorf("Bad dhcp_option name %s", dhcp_option.Name)
		}

		if "example" != dhcp_option.Annotation {
			return fmt.Errorf("Bad dhcp_option annotation %s", dhcp_option.Annotation)
		}

		if "example" != dhcp_option.Data {
			return fmt.Errorf("Bad dhcp_option data %s", dhcp_option.Data)
		}

		if "33" != dhcp_option.DHCPOption_id {
			return fmt.Errorf("Bad dhcp_option dhcp_option_id %s", dhcp_option.DHCPOption_id)
		}

		if "example" != dhcp_option.NameAlias {
			return fmt.Errorf("Bad dhcp_option name_alias %s", dhcp_option.NameAlias)
		}

		return nil
	}
}
