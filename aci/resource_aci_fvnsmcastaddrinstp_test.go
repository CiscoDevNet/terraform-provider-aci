package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciMulticastAddressPool_Basic(t *testing.T) {
	var multicast_address_pool models.MulticastAddressPool
	fvns_mcast_addr_inst_p_name := acctest.RandString(5)
	description := "multicast_address_pool created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMulticastAddressPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMulticastAddressPoolConfig_basic(fvns_mcast_addr_inst_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMulticastAddressPoolExists("aci_multicast_address_pool.foo_multicast_address_pool", &multicast_address_pool),
					testAccCheckAciMulticastAddressPoolAttributes(fvns_mcast_addr_inst_p_name, description, &multicast_address_pool),
				),
			},
		},
	})
}

func TestAccAciMulticastAddressPool_Update(t *testing.T) {
	var multicast_address_pool models.MulticastAddressPool
	fvns_mcast_addr_inst_p_name := acctest.RandString(5)
	description := "multicast_address_pool created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMulticastAddressPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMulticastAddressPoolConfig_basic(fvns_mcast_addr_inst_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMulticastAddressPoolExists("aci_multicast_address_pool.foo_multicast_address_pool", &multicast_address_pool),
					testAccCheckAciMulticastAddressPoolAttributes(fvns_mcast_addr_inst_p_name, description, &multicast_address_pool),
				),
			},
			{
				Config: testAccCheckAciMulticastAddressPoolConfig_basic(fvns_mcast_addr_inst_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMulticastAddressPoolExists("aci_multicast_address_pool.foo_multicast_address_pool", &multicast_address_pool),
					testAccCheckAciMulticastAddressPoolAttributes(fvns_mcast_addr_inst_p_name, description, &multicast_address_pool),
				),
			},
		},
	})
}

func testAccCheckAciMulticastAddressPoolConfig_basic(fvns_mcast_addr_inst_p_name string) string {
	return fmt.Sprintf(`

	resource "aci_multicast_address_pool" "foo_multicast_address_pool" {
		name 		= "%s"
		description = "multicast_address_pool created while acceptance testing"

	}

	`, fvns_mcast_addr_inst_p_name)
}

func testAccCheckAciMulticastAddressPoolExists(name string, multicast_address_pool *models.MulticastAddressPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Multicast Address Pool %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Multicast Address Pool dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		multicast_address_poolFound := models.MulticastAddressPoolFromContainer(cont)
		if multicast_address_poolFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Multicast Address Pool %s not found", rs.Primary.ID)
		}
		*multicast_address_pool = *multicast_address_poolFound
		return nil
	}
}

func testAccCheckAciMulticastAddressPoolDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_multicast_address_pool" {
			cont, err := client.Get(rs.Primary.ID)
			multicast_address_pool := models.MulticastAddressPoolFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Multicast Address Pool %s Still exists", multicast_address_pool.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciMulticastAddressPoolAttributes(fvns_mcast_addr_inst_p_name, description string, multicast_address_pool *models.MulticastAddressPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fvns_mcast_addr_inst_p_name != GetMOName(multicast_address_pool.DistinguishedName) {
			return fmt.Errorf("Bad fvnsmcast_addr_inst_p %s", GetMOName(multicast_address_pool.DistinguishedName))
		}

		if description != multicast_address_pool.Description {
			return fmt.Errorf("Bad multicast_address_pool Description %s", multicast_address_pool.Description)
		}
		return nil
	}
}
