package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciL3outBGPProtocolProfile_Basic(t *testing.T) {
	var l3out_bgp_protocol_profile models.L3outBGPProtocolProfile
	description := "protocol_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBGPProtocolProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outBGPProtocolProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPProtocolProfileExists("aci_l3out_bgp_protocol_profile.fool3out_bgp_protocol_profile", &l3out_bgp_protocol_profile),
					testAccCheckAciL3outBGPProtocolProfileAttributes(description, &l3out_bgp_protocol_profile),
				),
			},
		},
	})
}

func TestAccAciL3outBGPProtocolProfile_update(t *testing.T) {
	var l3out_bgp_protocol_profile models.L3outBGPProtocolProfile
	description := "protocol_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBGPProtocolProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outBGPProtocolProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPProtocolProfileExists("aci_l3out_bgp_protocol_profile.fool3out_bgp_protocol_profile", &l3out_bgp_protocol_profile),
					testAccCheckAciL3outBGPProtocolProfileAttributes(description, &l3out_bgp_protocol_profile),
				),
			},
			{
				Config: testAccCheckAciL3outBGPProtocolProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBGPProtocolProfileExists("aci_l3out_bgp_protocol_profile.fool3out_bgp_protocol_profile", &l3out_bgp_protocol_profile),
					testAccCheckAciL3outBGPProtocolProfileAttributes(description, &l3out_bgp_protocol_profile),
				),
			},
		},
	})
}

func testAccCheckAciL3outBGPProtocolProfileConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3out_bgp_protocol_profile" "fool3out_bgp_protocol_profile" {
		logical_node_profile_dn = "${aci_logical_node_profile.example.id}"
  		annotation  = "example"
  		name_alias  = "example"
	}
	`)
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

func testAccCheckAciL3outBGPProtocolProfileAttributes(description string, l3out_bgp_protocol_profile *models.L3outBGPProtocolProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if "example" != l3out_bgp_protocol_profile.Annotation {
			return fmt.Errorf("Bad l3out_bgp_protocol_profile annotation %s", l3out_bgp_protocol_profile.Annotation)
		}

		if "example" != l3out_bgp_protocol_profile.NameAlias {
			return fmt.Errorf("Bad l3out_bgp_protocol_profile name_alias %s", l3out_bgp_protocol_profile.NameAlias)
		}

		return nil
	}
}
