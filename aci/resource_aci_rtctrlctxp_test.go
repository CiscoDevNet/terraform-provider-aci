package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciRouteControlContext_Basic(t *testing.T) {
	var route_control_context models.RouteControlContext
	fv_tenant_name := acctest.RandString(5)
	rtctrl_profile_name := acctest.RandString(5)
	rtctrl_ctx_p_name := acctest.RandString(5)
	description := "route_control_context created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRouteControlContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRouteControlContextConfig_basic(fv_tenant_name, rtctrl_profile_name, rtctrl_ctx_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRouteControlContextExists("aci_route_control_context.fooroute_control_context", &route_control_context),
					testAccCheckAciRouteControlContextAttributes(fv_tenant_name, rtctrl_profile_name, rtctrl_ctx_p_name, description, &route_control_context),
				),
			},
		},
	})
}

func TestAccAciRouteControlContext_Update(t *testing.T) {
	var route_control_context models.RouteControlContext
	fv_tenant_name := acctest.RandString(5)
	rtctrl_profile_name := acctest.RandString(5)
	rtctrl_ctx_p_name := acctest.RandString(5)
	description := "route_control_context created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRouteControlContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRouteControlContextConfig_basic(fv_tenant_name, rtctrl_profile_name, rtctrl_ctx_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRouteControlContextExists("aci_route_control_context.fooroute_control_context", &route_control_context),
					testAccCheckAciRouteControlContextAttributes(fv_tenant_name, rtctrl_profile_name, rtctrl_ctx_p_name, description, &route_control_context),
				),
			},
			{
				Config: testAccCheckAciRouteControlContextConfig_basic(fv_tenant_name, rtctrl_profile_name, rtctrl_ctx_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRouteControlContextExists("aci_route_control_context.fooroute_control_context", &route_control_context),
					testAccCheckAciRouteControlContextAttributes(fv_tenant_name, rtctrl_profile_name, rtctrl_ctx_p_name, description, &route_control_context),
				),
			},
		},
	})
}

func testAccCheckAciRouteControlContextConfig_basic(fv_tenant_name, rtctrl_profile_name, rtctrl_ctx_p_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_route_control_profile" "fooroute_control_profile" {
		name 		= "%s"
		description = "route_control_profile created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	resource "aci_route_control_context" "fooroute_control_context" {
		name 		= "%s"
		description = "route_control_context created while acceptance testing"
		route_control_profile_dn = aci_route_control_profile.fooroute_control_profile.id
	}

	`, fv_tenant_name, rtctrl_profile_name, rtctrl_ctx_p_name)
}

func testAccCheckAciRouteControlContextExists(name string, route_control_context *models.RouteControlContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Route Control Context %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route Control Context dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		route_control_contextFound := models.RouteControlContextFromContainer(cont)
		if route_control_contextFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Route Control Context %s not found", rs.Primary.ID)
		}
		*route_control_context = *route_control_contextFound
		return nil
	}
}

func testAccCheckAciRouteControlContextDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_route_control_context" {
			cont, err := client.Get(rs.Primary.ID)
			route_control_context := models.RouteControlContextFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Route Control Context %s Still exists", route_control_context.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRouteControlContextAttributes(fv_tenant_name, rtctrl_profile_name, rtctrl_ctx_p_name, description string, route_control_context *models.RouteControlContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if rtctrl_ctx_p_name != GetMOName(route_control_context.DistinguishedName) {
			return fmt.Errorf("Bad rtctrl_ctx_p %s", GetMOName(route_control_context.DistinguishedName))
		}

		if rtctrl_profile_name != GetMOName(GetParentDn(route_control_context.DistinguishedName)) {
			return fmt.Errorf(" Bad rtctrl_profile %s", GetMOName(GetParentDn(route_control_context.DistinguishedName)))
		}
		if description != route_control_context.Description {
			return fmt.Errorf("Bad route_control_context Description %s", route_control_context.Description)
		}
		return nil
	}
}
