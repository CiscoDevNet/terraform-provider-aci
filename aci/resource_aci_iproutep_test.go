package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciL3outStaticRoute_Basic(t *testing.T) {
	var l3out_static_route models.L3outStaticRoute
	description := "static_route created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outStaticRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outStaticRouteConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteExists("aci_l3out_static_route.fool3out_static_route", &l3out_static_route),
					testAccCheckAciL3outStaticRouteAttributes(description, &l3out_static_route),
				),
			},
		},
	})
}

func TestAccAciL3outStaticRoute_update(t *testing.T) {
	var l3out_static_route models.L3outStaticRoute
	description := "static_route created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outStaticRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outStaticRouteConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteExists("aci_l3out_static_route.fool3out_static_route", &l3out_static_route),
					testAccCheckAciL3outStaticRouteAttributes(description, &l3out_static_route),
				),
			},
			{
				Config: testAccCheckAciL3outStaticRouteConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteExists("aci_l3out_static_route.fool3out_static_route", &l3out_static_route),
					testAccCheckAciL3outStaticRouteAttributes(description, &l3out_static_route),
				),
			},
		},
	})
}

func testAccCheckAciL3outStaticRouteConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3out_static_route" "fool3out_static_route" {
		fabric_node_dn  = "${aci_fabric_node.example.id}"
		description = "%s"
		ip  = "example"
  		aggregate = "no"
  		annotation  = "example"
  		name_alias  = "example"
  		pref  = "example"
  		rt_ctrl = "bfd"
	}
	`, description)
}

func testAccCheckAciL3outStaticRouteExists(name string, l3out_static_route *models.L3outStaticRoute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Static Route %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Static Route dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_static_routeFound := models.L3outStaticRouteFromContainer(cont)
		if l3out_static_routeFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Static Route %s not found", rs.Primary.ID)
		}
		*l3out_static_route = *l3out_static_routeFound
		return nil
	}
}

func testAccCheckAciL3outStaticRouteDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_static_route" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_static_route := models.L3outStaticRouteFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Static Route %s Still exists", l3out_static_route.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3outStaticRouteAttributes(description string, l3out_static_route *models.L3outStaticRoute) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3out_static_route.Description {
			return fmt.Errorf("Bad l3out_static_route Description %s", l3out_static_route.Description)
		}

		if "example" != l3out_static_route.Ip {
			return fmt.Errorf("Bad l3out_static_route ip %s", l3out_static_route.Ip)
		}

		if "no" != l3out_static_route.Aggregate {
			return fmt.Errorf("Bad l3out_static_route aggregate %s", l3out_static_route.Aggregate)
		}

		if "example" != l3out_static_route.Annotation {
			return fmt.Errorf("Bad l3out_static_route annotation %s", l3out_static_route.Annotation)
		}

		if "example" != l3out_static_route.NameAlias {
			return fmt.Errorf("Bad l3out_static_route name_alias %s", l3out_static_route.NameAlias)
		}

		if "example" != l3out_static_route.Pref {
			return fmt.Errorf("Bad l3out_static_route pref %s", l3out_static_route.Pref)
		}

		if "bfd" != l3out_static_route.RtCtrl {
			return fmt.Errorf("Bad l3out_static_route rt_ctrl %s", l3out_static_route.RtCtrl)
		}

		return nil
	}
}
