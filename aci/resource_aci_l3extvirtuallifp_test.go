package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciVirtualLogicalInterfaceProfile_Basic(t *testing.T) {
	var logical_interface_profile models.VirtualLogicalInterfaceProfile
	description := "logical_interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVirtualLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVirtualLogicalInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVirtualLogicalInterfaceProfileExists("aci_l3out_floating_svi.test", &logical_interface_profile),
					testAccCheckAciVirtualLogicalInterfaceProfileAttributes(description, &logical_interface_profile),
				),
			},
		},
	})
}

func TestAccAciVirtualLogicalInterfaceProfile_update(t *testing.T) {
	var logical_interface_profile models.VirtualLogicalInterfaceProfile
	description := "logical_interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVirtualLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVirtualLogicalInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVirtualLogicalInterfaceProfileExists("aci_l3out_floating_svi.test", &logical_interface_profile),
					testAccCheckAciVirtualLogicalInterfaceProfileAttributes(description, &logical_interface_profile),
				),
			},
			{
				Config: testAccCheckAciVirtualLogicalInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVirtualLogicalInterfaceProfileExists("aci_l3out_floating_svi.test", &logical_interface_profile),
					testAccCheckAciVirtualLogicalInterfaceProfileAttributes(description, &logical_interface_profile),
				),
			},
		},
	})
}

func testAccCheckAciVirtualLogicalInterfaceProfileConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn = "uni/tn-aaaaa/out-demo_l3out/lnodep-demo_node/lifp-demo_int_prof"
		node_dn                      = "topology/pod-1/node-201"
		encap                        = "vlan-20"
		addr                         = "10.20.30.40/16"
		annotation                   = "example"
		description                  = "%s"
		autostate                    = "enabled"
		encap_scope                  = "ctx"
		if_inst_t                    = "ext-svi"
		ipv6_dad                     = "disabled"
		ll_addr                      = "::"
		mac                          = "12:23:34:45:56:67"
		mode                         = "untagged"
		mtu                          = "580"
		target_dscp                  = "CS1"
	}
	`, description)
}

func testAccCheckAciVirtualLogicalInterfaceProfileExists(name string, logical_interface_profile *models.VirtualLogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Logical Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Logical Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		logical_interface_profileFound := models.VirtualLogicalInterfaceProfileFromContainer(cont)
		if logical_interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Logical Interface Profile %s not found", rs.Primary.ID)
		}
		*logical_interface_profile = *logical_interface_profileFound
		return nil
	}
}

func testAccCheckAciVirtualLogicalInterfaceProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l3out_floating_svi" {
			cont, err := client.Get(rs.Primary.ID)
			logical_interface_profile := models.LogicalInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Logical Interface Profile %s Still exists", logical_interface_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciVirtualLogicalInterfaceProfileAttributes(description string, logical_interface_profile *models.VirtualLogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != logical_interface_profile.Description {
			return fmt.Errorf("Bad logical_interface_profile Description %s", logical_interface_profile.Description)
		}

		if "topology/pod-1/node-201" != logical_interface_profile.NodeDn {
			return fmt.Errorf("Bad logical_interface_profile node_dn %s", logical_interface_profile.NodeDn)
		}

		if "vlan-20" != logical_interface_profile.Encap {
			return fmt.Errorf("Bad logical_interface_profile encap %s", logical_interface_profile.Encap)
		}

		if "10.20.30.40/16" != logical_interface_profile.Addr {
			return fmt.Errorf("Bad logical_interface_profile addr %s", logical_interface_profile.Addr)
		}

		if "example" != logical_interface_profile.Annotation {
			return fmt.Errorf("Bad logical_interface_profile annotation %s", logical_interface_profile.Annotation)
		}

		if "enabled" != logical_interface_profile.Autostate {
			return fmt.Errorf("Bad logical_interface_profile autostate %s", logical_interface_profile.Autostate)
		}

		if "ctx" != logical_interface_profile.EncapScope {
			return fmt.Errorf("Bad logical_interface_profile encap_scope %s", logical_interface_profile.EncapScope)
		}

		if "ext-svi" != logical_interface_profile.IfInstT {
			return fmt.Errorf("Bad logical_interface_profile if_inst_t %s", logical_interface_profile.IfInstT)
		}

		if "disabled" != logical_interface_profile.Ipv6Dad {
			return fmt.Errorf("Bad logical_interface_profile ipv6_dad %s", logical_interface_profile.Ipv6Dad)
		}

		if "::" != logical_interface_profile.LlAddr {
			return fmt.Errorf("Bad logical_interface_profile ll_addr %s", logical_interface_profile.LlAddr)
		}

		if "12:23:34:45:56:67" != logical_interface_profile.Mac {
			return fmt.Errorf("Bad logical_interface_profile mac %s", logical_interface_profile.Mac)
		}

		if "untagged" != logical_interface_profile.Mode {
			return fmt.Errorf("Bad logical_interface_profile mode %s", logical_interface_profile.Mode)
		}

		if "580" != logical_interface_profile.Mtu {
			return fmt.Errorf("Bad logical_interface_profile mtu %s", logical_interface_profile.Mtu)
		}

		if "CS1" != logical_interface_profile.TargetDscp {
			return fmt.Errorf("Bad logical_interface_profile target_dscp %s", logical_interface_profile.TargetDscp)
		}

		return nil
	}
}
