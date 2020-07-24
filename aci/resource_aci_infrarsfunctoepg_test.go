package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciEPGsUsingFunction_Basic(t *testing.T) {
	var ep_gs_using_function models.EPGsUsingFunction

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEPGsUsingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciEPGsUsingFunctionConfig_basic("vlan-5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists("aci_epgs_using_function.fooep_gs_using_function", &ep_gs_using_function),
					testAccCheckAciEPGsUsingFunctionAttributes("vlan-5", &ep_gs_using_function),
				),
			},
		},
	})
}

func TestAccAciEPGsUsingFunction_update(t *testing.T) {
	var ep_gs_using_function models.EPGsUsingFunction

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEPGsUsingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciEPGsUsingFunctionConfig_basic("vlan-5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists("aci_epgs_using_function.fooep_gs_using_function", &ep_gs_using_function),
					testAccCheckAciEPGsUsingFunctionAttributes("vlan-5", &ep_gs_using_function),
				),
			},
			{
				Config: testAccCheckAciEPGsUsingFunctionConfig_basic("vlan-10"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEPGsUsingFunctionExists("aci_epgs_using_function.fooep_gs_using_function", &ep_gs_using_function),
					testAccCheckAciEPGsUsingFunctionAttributes("vlan-10", &ep_gs_using_function),
				),
			},
		},
	})
}

func testAccCheckAciEPGsUsingFunctionConfig_basic(encap string) string {
	return fmt.Sprintf(`

	resource "aci_epgs_using_function" "fooep_gs_using_function" {
		  access_generic_dn  = "${aci_access_generic.fooaccess_generic.id}"
		  annotation  = "example"
		  t_dn  = "${aci_application_epg.epg1.id}"
		  encap  = "%s"
		  mode  = "regular"
		  primary_encap  = "vlan-7"
		}
	`, encap)
}

func testAccCheckAciEPGsUsingFunctionExists(name string, ep_gs_using_function *models.EPGsUsingFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("EPGs Using Function %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EPGs Using Function dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ep_gs_using_functionFound := models.EPGsUsingFunctionFromContainer(cont)
		if ep_gs_using_functionFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("EPGs Using Function %s not found", rs.Primary.ID)
		}
		*ep_gs_using_function = *ep_gs_using_functionFound
		return nil
	}
}

func testAccCheckAciEPGsUsingFunctionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_epgs_using_function" {
			cont, err := client.Get(rs.Primary.ID)
			ep_gs_using_function := models.EPGsUsingFunctionFromContainer(cont)
			if err == nil {
				return fmt.Errorf("EPGs Using Function %s Still exists", ep_gs_using_function.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciEPGsUsingFunctionAttributes(encap string, ep_gs_using_function *models.EPGsUsingFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if encap != ep_gs_using_function.Encap {
			return fmt.Errorf("Bad epgs_using_function Description %s", ep_gs_using_function.Description)
		}

		if "example" != ep_gs_using_function.Annotation {
			return fmt.Errorf("Bad epgs_using_function annotation %s", ep_gs_using_function.Annotation)
		}

		if "regular" != ep_gs_using_function.Mode {
			return fmt.Errorf("Bad epgs_using_function mode %s", ep_gs_using_function.Mode)
		}

		if "vlan-7" != ep_gs_using_function.PrimaryEncap {
			return fmt.Errorf("Bad epgs_using_function primary_encap %s", ep_gs_using_function.PrimaryEncap)
		}

		return nil
	}
}
