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

func TestAccAciMulticastAddressBlock_Basic(t *testing.T) {
	var multicast_address_block models.MulticastAddressBlock
	fvns_mcast_addr_inst_p_name := acctest.RandString(5)
	fvns_mcast_addr_blk_name := acctest.RandString(5)
	description := "multicast_address_block created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMulticastAddressBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMulticastAddressBlockConfig_basic(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMulticastAddressBlockExists("aci_multicast_pool_block.foo_multicast_pool_block", &multicast_address_block),
					testAccCheckAciMulticastAddressBlockAttributes(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name, description, &multicast_address_block),
				),
			},
		},
	})
}

func TestAccAciMulticastAddressBlock_Update(t *testing.T) {
	var multicast_address_block models.MulticastAddressBlock
	fvns_mcast_addr_inst_p_name := acctest.RandString(5)
	fvns_mcast_addr_blk_name := acctest.RandString(5)
	description := "multicast_address_block created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMulticastAddressBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMulticastAddressBlockConfig_basic(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMulticastAddressBlockExists("aci_multicast_pool_block.foo_multicast_pool_block", &multicast_address_block),
					testAccCheckAciMulticastAddressBlockAttributes(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name, description, &multicast_address_block),
				),
			},
			{
				Config: testAccCheckAciMulticastAddressBlockConfig_basic(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMulticastAddressBlockExists("aci_multicast_pool_block.foo_multicast_pool_block", &multicast_address_block),
					testAccCheckAciMulticastAddressBlockAttributes(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name, description, &multicast_address_block),
				),
			},
		},
	})
}

func testAccCheckAciMulticastAddressBlockConfig_basic(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name string) string {
	return fmt.Sprintf(`

	resource "aci_multicast_pool" "foo_multicast_pool" {
		name 		= "%s"
		description = "multicast_address_pool created while acceptance testing"
	}

	resource "aci_multicast_pool_block" "foo_multicast_pool_block" {
		name 		= "%s"
		description = "multicast_address_block created while acceptance testing"
		multicast_pool_dn = aci_multicast_pool.foo_multicast_pool.id
		from = "224.0.0.0"
		to = "224.0.0.10"
	}

	`, fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name)
}

func testAccCheckAciMulticastAddressBlockExists(name string, multicast_address_block *models.MulticastAddressBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Multicast Address Block %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Multicast Address Block dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		multicast_address_blockFound := models.MulticastAddressBlockFromContainer(cont)
		if multicast_address_blockFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Multicast Address Block %s not found", rs.Primary.ID)
		}
		*multicast_address_block = *multicast_address_blockFound
		return nil
	}
}

func testAccCheckAciMulticastAddressBlockDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_multicast_pool_block" {
			cont, err := client.Get(rs.Primary.ID)
			multicast_address_block := models.MulticastAddressBlockFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Multicast Address Block %s Still exists", multicast_address_block.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciMulticastAddressBlockAttributes(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name, description string, multicast_address_block *models.MulticastAddressBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fvns_mcast_addr_blk_name != multicast_address_block.Name {
			return fmt.Errorf("Bad fvnsmcast_addr_blk %s", multicast_address_block.Name)
		}
		if fvns_mcast_addr_inst_p_name != GetMOName(GetParentDn(multicast_address_block.DistinguishedName, fmt.Sprintf("/"+models.RnfvnsMcastAddrBlk, multicast_address_block.From, multicast_address_block.To))) {
			return fmt.Errorf(" Bad fvnsmcast_addr_inst_p %s %s", fvns_mcast_addr_inst_p_name, GetMOName(GetParentDn(multicast_address_block.DistinguishedName, fmt.Sprintf("/"+models.RnfvnsMcastAddrBlk, multicast_address_block.From, multicast_address_block.To))))
		}
		if description != multicast_address_block.Description {
			return fmt.Errorf("Bad multicast_address_block Description %s", multicast_address_block.Description)
		}
		return nil
	}
}
