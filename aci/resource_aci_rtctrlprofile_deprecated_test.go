package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBgpRouteControlProfile_Basic(t *testing.T) {
	var route_control_profile models.RouteControlProfile
	description := "route_control_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBgpRouteControlProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBgpRouteControlProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteControlProfileExists("aci_bgp_route_control_profile.test", &route_control_profile),
					testAccCheckAciBgpRouteControlProfileAttributes(description, &route_control_profile),
				),
			},
		},
	})
}

func TestAccAciBgpRouteControlProfile_update(t *testing.T) {
	var route_control_profile models.RouteControlProfile
	description := "route_control_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBgpRouteControlProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBgpRouteControlProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteControlProfileExists("aci_bgp_route_control_profile.test", &route_control_profile),
					testAccCheckAciBgpRouteControlProfileAttributes(description, &route_control_profile),
				),
			},
			{
				Config: testAccCheckAciBgpRouteControlProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteControlProfileExists("aci_bgp_route_control_profile.test", &route_control_profile),
					testAccCheckAciBgpRouteControlProfileAttributes(description, &route_control_profile),
				),
			},
		},
	})
}

func testAccCheckAciBgpRouteControlProfileConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name        = "tenant_1"
		description = "This tenant is created by terraform ACI provider"
	}
	  
	resource "aci_l3_outside" "fool3_outside" {
		tenant_dn      = aci_tenant.foo_tenant.id
		name           = "l3_outside_1"
		annotation     = "l3_outside_tag"
		name_alias     = "alias_out"
		target_dscp    = "unspecified"
	}	

	resource "aci_bgp_route_control_profile" "test" {
		parent_dn                  = aci_l3_outside.fool3_outside.id
		name                       = "one"
		annotation                 = "example"
		description                = "%s"
		name_alias                 = "example"
		route_control_profile_type = "global"
	  }
	`, description)
}

func testAccCheckAciBgpRouteControlProfileExists(name string, route_control_profile *models.RouteControlProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Route Control Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route Control Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		route_control_profileFound := models.RouteControlProfileFromContainer(cont)
		if route_control_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Route Control Profile %s not found", rs.Primary.ID)
		}
		*route_control_profile = *route_control_profileFound
		return nil
	}
}

func testAccCheckAciBgpRouteControlProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_bgp_route_control_profile" {
			cont, err := client.Get(rs.Primary.ID)
			route_control_profile := models.RouteControlProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Route Control Profile %s Still exists", route_control_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciBgpRouteControlProfileAttributes(description string, route_control_profile *models.RouteControlProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != route_control_profile.Description {
			return fmt.Errorf("Bad route_control_profile Description %s", route_control_profile.Description)
		}

		if "one" != route_control_profile.Name {
			return fmt.Errorf("Bad route_control_profile name %s", route_control_profile.Name)
		}

		if "example" != route_control_profile.Annotation {
			return fmt.Errorf("Bad route_control_profile annotation %s", route_control_profile.Annotation)
		}

		if "example" != route_control_profile.NameAlias {
			return fmt.Errorf("Bad route_control_profile name_alias %s", route_control_profile.NameAlias)
		}

		if "global" != route_control_profile.RouteControlProfileType {
			return fmt.Errorf("Bad route_control_profile route_control_profile_type %s", route_control_profile.RouteControlProfileType)
		}

		return nil
	}
}
