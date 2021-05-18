package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciL3outStaticRouteNextHop_Basic(t *testing.T) {
	var l3out_static_route_next_hop models.L3outStaticRouteNextHop
	description := "nexthop_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outStaticRouteNextHopDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outStaticRouteNextHopConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteNextHopExists("aci_l3out_static_route_next_hop.fool3out_static_route_next_hop", &l3out_static_route_next_hop),
					testAccCheckAciL3outStaticRouteNextHopAttributes(description, &l3out_static_route_next_hop),
				),
			},
		},
	})
}

func TestAccAciL3outStaticRouteNextHop_update(t *testing.T) {
	var l3out_static_route_next_hop models.L3outStaticRouteNextHop
	description := "nexthop_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outStaticRouteNextHopDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outStaticRouteNextHopConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteNextHopExists("aci_l3out_static_route_next_hop.fool3out_static_route_next_hop", &l3out_static_route_next_hop),
					testAccCheckAciL3outStaticRouteNextHopAttributes(description, &l3out_static_route_next_hop),
				),
			},
			{
				Config: testAccCheckAciL3outStaticRouteNextHopConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteNextHopExists("aci_l3out_static_route_next_hop.fool3out_static_route_next_hop", &l3out_static_route_next_hop),
					testAccCheckAciL3outStaticRouteNextHopAttributes(description, &l3out_static_route_next_hop),
				),
			},
		},
	})
}

func testAccCheckAciL3outStaticRouteNextHopConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3out_static_route_next_hop" "fool3out_static_route_next_hop" {
		static_route_dn  = "${aci_static_route.example.id}"
		nh_addr  = "example"
		description = "%s"
  		annotation  = "example"
  		name_alias  = "example"
  		pref = "unspecified"
  		nexthop_profile_type = "none"
	}
	`, description)
}

func testAccCheckAciL3outStaticRouteNextHopExists(name string, l3out_static_route_next_hop *models.L3outStaticRouteNextHop) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Static Route Next Hop %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Static Route Next Hop dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_static_route_next_hopFound := models.L3outStaticRouteNextHopFromContainer(cont)
		if l3out_static_route_next_hopFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Static Route Next Hop %s not found", rs.Primary.ID)
		}
		*l3out_static_route_next_hop = *l3out_static_route_next_hopFound
		return nil
	}
}

func testAccCheckAciL3outStaticRouteNextHopDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_nexthop_profile" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_static_route_next_hop := models.L3outStaticRouteNextHopFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Static Route Next Hop %s Still exists", l3out_static_route_next_hop.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3outStaticRouteNextHopAttributes(description string, l3out_static_route_next_hop *models.L3outStaticRouteNextHop) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3out_static_route_next_hop.Description {
			return fmt.Errorf("Bad l3out_static_route_next_hop Description %s", l3out_static_route_next_hop.Description)
		}

		if "example" != l3out_static_route_next_hop.NhAddr {
			return fmt.Errorf("Bad l3out_static_route_next_hop nh_addr %s", l3out_static_route_next_hop.NhAddr)
		}

		if "example" != l3out_static_route_next_hop.Annotation {
			return fmt.Errorf("Bad l3out_static_route_next_hop annotation %s", l3out_static_route_next_hop.Annotation)
		}

		if "example" != l3out_static_route_next_hop.NameAlias {
			return fmt.Errorf("Bad l3out_static_route_next_hop name_alias %s", l3out_static_route_next_hop.NameAlias)
		}

		if "example" != l3out_static_route_next_hop.NhAddr {
			return fmt.Errorf("Bad l3out_static_route_next_hop nh_addr %s", l3out_static_route_next_hop.NhAddr)
		}

		if "unspecified" != l3out_static_route_next_hop.Pref {
			return fmt.Errorf("Bad l3out_static_route_next_hop pref %s", l3out_static_route_next_hop.Pref)
		}

		if "none" != l3out_static_route_next_hop.NexthopProfile_type {
			return fmt.Errorf("Bad l3out_static_route_next_hop nexthop_profile_type %s", l3out_static_route_next_hop.NexthopProfile_type)
		}

		return nil
	}
}
