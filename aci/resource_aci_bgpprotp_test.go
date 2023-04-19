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

func TestAccAciL3outBGPProtocolProfile_Basic(t *testing.T) {
	var l3out_bgp_protocol_profile models.L3outBGPProtocolProfile
	fv_tenant_name := acctest.RandString(5)
	l3ext_out_name := acctest.RandString(5)
	l3ext_l_node_p_name := acctest.RandString(5)
	bgp_prot_p_name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBGPProtocolProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outBGPProtocolProfileConfig_basic(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, bgp_prot_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPProtocolProfileExists("aci_l3out_bgp_protocol_profile.fool3out_bgp_protocol_profile", &l3out_bgp_protocol_profile),
					testAccCheckAciL3outBGPProtocolProfileAttributes(&l3out_bgp_protocol_profile),
				),
			},
		},
	})
}

func TestAccAciL3outBGPProtocolProfile_update(t *testing.T) {
	var l3out_bgp_protocol_profile models.L3outBGPProtocolProfile
	fv_tenant_name := acctest.RandString(5)
	l3ext_out_name := acctest.RandString(5)
	l3ext_l_node_p_name := acctest.RandString(5)
	bgp_prot_p_name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBGPProtocolProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outBGPProtocolProfileConfig_basic(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, bgp_prot_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPProtocolProfileExists("aci_l3out_bgp_protocol_profile.fool3out_bgp_protocol_profile", &l3out_bgp_protocol_profile),
					testAccCheckAciL3outBGPProtocolProfileAttributes(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, bgp_prot_p_name, &l3out_bgp_protocol_profile),
				),
			},
			{
				Config: testAccCheckAciL3outBGPProtocolProfileConfig_basic(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, bgp_prot_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPProtocolProfileExists("aci_l3out_bgp_protocol_profile.fool3out_bgp_protocol_profile", &l3out_bgp_protocol_profile),
					testAccCheckAciL3outBGPProtocolProfileAttributes(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, bgp_prot_p_name, &l3out_bgp_protocol_profile),
				),
			},
		},
	})
}

func testAccCheckAciL3outBGPProtocolProfileConfig_basic(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, bgp_prot_p_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_l3_outside" "fool3_outside" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	resource "aci_logical_node_profile" "foological_node_profile" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.fool3_outside.id
	}

	resource "aci_l3out_bgp_protocol_profile" "fool3out_bgp_protocol_profile" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.foological_node_profile.id
	}
	`, fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, bgp_prot_p_name)
}

func testAccCheckAciL3outBGPProtocolProfileExists(name string, l3out_bgp_protocol_profile *models.L3outBGPProtocolProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out BGP Protocol Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out BGP Protocol Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_bgp_protocol_profileFound := models.L3outBGPProtocolProfileFromContainer(cont)
		if l3out_bgp_protocol_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out BGP Protocol Profile %s not found", rs.Primary.ID)
		}
		*l3out_bgp_protocol_profile = *l3out_bgp_protocol_profileFound
		return nil
	}
}

func testAccCheckAciL3outBGPProtocolProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_protocol_profile" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_bgp_protocol_profile := models.L3outBGPProtocolProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out BGP Protocol Profile %s Still exists", l3out_bgp_protocol_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outBGPProtocolProfileAttributes(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, bgp_prot_p_name string, l3out_bgp_protocol_profile *models.L3outBGPProtocolProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "bgp_prot_p_name" != GetMOName(l3out_bgp_protocol_profile.DistinguishedName) {
			return fmt.Errorf("Bad l3out_bgp_protocol_profile %s", GetMOName(l3out_bgp_protocol_profile.DistinguishedName))
		}

		if "l3ext_l_node_p_name" != GetMOName(GetParentDn(l3out_bgp_protocol_profile.DistinguishedName)) {
			return fmt.Errorf("Bad l3extl_node_p %s", GetMOName(GetParentDn(l3out_bgp_protocol_profile.DistinguishedName)))
		}
		return nil
	}
}
