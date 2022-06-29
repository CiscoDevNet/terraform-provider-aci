package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciAbstractionofIPAddressBlock_Basic(t *testing.T) {
	var abstractionof_ipaddress_block models.AbstractionofIPAddressBlock
	fvns_mcast_addr_inst_p_name := acctest.RandString(5)
	fvns_mcast_addr_blk_name := acctest.RandString(5)
	description := "abstractionof_ipaddress_block created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAbstractionofIPAddressBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAbstractionofIPAddressBlockConfig_basic(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAbstractionofIPAddressBlockExists("aci_abstractionof_ipaddress_block.foo_abstractionof_ipaddress_block", &abstractionof_ipaddress_block),
					testAccCheckAciAbstractionofIPAddressBlockAttributes(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name, description, &abstractionof_ipaddress_block),
				),
			},
		},
	})
}

func TestAccAciAbstractionofIPAddressBlock_Update(t *testing.T) {
	var abstractionof_ipaddress_block models.AbstractionofIPAddressBlock
	fvns_mcast_addr_inst_p_name := acctest.RandString(5)
	fvns_mcast_addr_blk_name := acctest.RandString(5)
	description := "abstractionof_ipaddress_block created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAbstractionofIPAddressBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAbstractionofIPAddressBlockConfig_basic(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAbstractionofIPAddressBlockExists("aci_abstractionof_ipaddress_block.foo_abstractionof_ipaddress_block", &abstractionof_ipaddress_block),
					testAccCheckAciAbstractionofIPAddressBlockAttributes(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name, description, &abstractionof_ipaddress_block),
				),
			},
			{
				Config: testAccCheckAciAbstractionofIPAddressBlockConfig_basic(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAbstractionofIPAddressBlockExists("aci_abstractionof_ipaddress_block.foo_abstractionof_ipaddress_block", &abstractionof_ipaddress_block),
					testAccCheckAciAbstractionofIPAddressBlockAttributes(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name, description, &abstractionof_ipaddress_block),
				),
			},
		},
	})
}

func testAccCheckAciAbstractionofIPAddressBlockConfig_basic(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name string) string {
	return fmt.Sprintf(`

	resource "aci_multicast_address_pool" "foo_multicast_address_pool" {
		name 		= "%s"
		description = "multicast_address_pool created while acceptance testing"

	}

	resource "aci_abstractionof_ipaddress_block" "foo_abstractionof_ipaddress_block" {
		name 		= "%s"
		description = "abstractionof_ipaddress_block created while acceptance testing"
		multicast_pool_dn = aci_multicast_address_pool.foo_multicast_address_pool.id
	}

	`, fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name)
}

func testAccCheckAciAbstractionofIPAddressBlockExists(name string, abstractionof_ipaddress_block *models.AbstractionofIPAddressBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Abstraction of IP Address Block %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Abstraction of IP Address Block dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		abstractionof_ipaddress_blockFound := models.AbstractionofIPAddressBlockFromContainer(cont)
		if abstractionof_ipaddress_blockFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Abstraction of IP Address Block %s not found", rs.Primary.ID)
		}
		*abstractionof_ipaddress_block = *abstractionof_ipaddress_blockFound
		return nil
	}
}

func testAccCheckAciAbstractionofIPAddressBlockDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_abstractionof_ipaddress_block" {
			cont, err := client.Get(rs.Primary.ID)
			abstractionof_ipaddress_block := models.AbstractionofIPAddressBlockFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Abstraction of IP Address Block %s Still exists", abstractionof_ipaddress_block.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciAbstractionofIPAddressBlockAttributes(fvns_mcast_addr_inst_p_name, fvns_mcast_addr_blk_name, description string, abstractionof_ipaddress_block *models.AbstractionofIPAddressBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fvns_mcast_addr_blk_name != GetMOName(abstractionof_ipaddress_block.DistinguishedName) {
			return fmt.Errorf("Bad fvnsmcast_addr_blk %s", GetMOName(abstractionof_ipaddress_block.DistinguishedName))
		}

		if fvns_mcast_addr_inst_p_name != GetMOName(GetParentDn(abstractionof_ipaddress_block.DistinguishedName, abstractionof_ipaddress_block.Rn)) {
			return fmt.Errorf(" Bad fvnsmcast_addr_inst_p %s", GetMOName(GetParentDn(abstractionof_ipaddress_block.DistinguishedName, abstractionof_ipaddress_block.Rn)))
		}
		if description != abstractionof_ipaddress_block.Description {
			return fmt.Errorf("Bad abstractionof_ipaddress_block Description %s", abstractionof_ipaddress_block.Description)
		}
		return nil
	}
}
