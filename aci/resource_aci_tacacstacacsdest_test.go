package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciTACACSDestination_Basic(t *testing.T) {
	var tacacs_destination models.TACACSDestination
	description := "tacacs_destination created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTACACSDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTACACSDestinationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSDestinationExists("aci_tacacs_accounting_destination.footacacs_destination", &tacacs_destination),
					testAccCheckAciTACACSDestinationAttributes(description, &tacacs_destination),
				),
			},
		},
	})
}

func TestAccAciTACACSDestination_Update(t *testing.T) {
	var tacacs_destination models.TACACSDestination
	description := "tacacs_destination created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTACACSDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTACACSDestinationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSDestinationExists("aci_tacacs_accounting_destination.footacacs_destination", &tacacs_destination),
					testAccCheckAciTACACSDestinationAttributes(description, &tacacs_destination),
				),
			},
			{
				Config: testAccCheckAciTACACSDestinationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSDestinationExists("aci_tacacs_accounting_destination.footacacs_destination", &tacacs_destination),
					testAccCheckAciTACACSDestinationAttributes(description, &tacacs_destination),
				),
			},
		},
	})
}

func testAccCheckAciTACACSDestinationConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_tacacs_accounting_destination" "footacacs_destination" {
		description = "%s"
		tacacs_accounting_dn = aci_tacacs_accounting.test.id
		host = "test"
		port = "49"
		name = "test_name_example"
		auth_protocol = "pap"
	}

	`, description)
}

func testAccCheckAciTACACSDestinationExists(name string, tacacs_destination *models.TACACSDestination) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("TACACS Destination %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TACACS Destination dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tacacs_destinationFound := models.TACACSDestinationFromContainer(cont)
		if tacacs_destinationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("TACACS Destination %s not found", rs.Primary.ID)
		}
		*tacacs_destination = *tacacs_destinationFound
		return nil
	}
}

func testAccCheckAciTACACSDestinationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tacacs_accounting_destination" {
			cont, err := client.Get(rs.Primary.ID)
			tacacs_destination := models.TACACSDestinationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("TACACS Destination %s Still exists", tacacs_destination.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTACACSDestinationAttributes(description string, tacacs_destination *models.TACACSDestination) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if description != tacacs_destination.Description {
			return fmt.Errorf("Bad tacacs_accounting_destination Description %s", tacacs_destination.Description)
		}

		if "uni/fabric/tacacsgroup-any" != GetParentDn(tacacs_destination.DistinguishedName, fmt.Sprintf("/tacacsdest-%s-port-%s", tacacs_destination.Host, tacacs_destination.Port)) {
			return fmt.Errorf("Bad tacacs_accounting_destination ParentDn %s", fmt.Sprintf("/tacacsdest-%s-port-%s", tacacs_destination.Host, tacacs_destination.Port))
		}

		if "49" != tacacs_destination.Port {
			return fmt.Errorf("Bad tacacs_accounting_destination Port %s", tacacs_destination.Port)
		}

		if "test" != tacacs_destination.Host {
			return fmt.Errorf("Bad tacacs_accounting_destination Host %s", tacacs_destination.Host)
		}

		if "test_name_example" != tacacs_destination.Name {
			return fmt.Errorf("Bad tacacs_accounting_destination Name %s", tacacs_destination.Name)
		}

		if "pap" != tacacs_destination.AuthProtocol {
			return fmt.Errorf("Bad tacacs_accounting_destination AuthProtocol %s", tacacs_destination.AuthProtocol)
		}

		return nil
	}
}
