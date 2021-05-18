package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciOSPFInterfaceProfile_Basic(t *testing.T) {
	var L3outOSPF models.OSPFInterfaceProfile
	description := "L3out OSPF Interface Profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOSPFInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOSPFInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfaceProfileExists("aci_l3out_ospf_interface_profile.test", &L3outOSPF),
					testAccCheckAciOSPFInterfaceProfileAttributes(description, &L3outOSPF),
				),
			},
		},
	})
}

func TestAccAciOSPFInterfaceProfile_update(t *testing.T) {
	var L3outOSPF models.OSPFInterfaceProfile
	description := "L3out OSPF Interface Profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOSPFInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOSPFInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfaceProfileExists("aci_l3out_ospf_interface_profile.test", &L3outOSPF),
					testAccCheckAciOSPFInterfaceProfileAttributes(description, &L3outOSPF),
				),
			},
			{
				Config: testAccCheckAciOSPFInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfaceProfileExists("aci_l3out_ospf_interface_profile.test", &L3outOSPF),
					testAccCheckAciOSPFInterfaceProfileAttributes(description, &L3outOSPF),
				),
			},
		},
	})
}

func testAccCheckAciOSPFInterfaceProfileConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn = "${aci_logical_interface_profile.example.id}"
		description                  = "%s"
		annotation                   = "example"
		auth_key                     = "example"
		auth_key_id                  = "255"
		auth_type                    = "simple"
		name_alias                   = "example"
	}
	`, description)
}

func testAccCheckAciOSPFInterfaceProfileExists(name string, interface_profile *models.OSPFInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out OSPF Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out OSPF Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		interface_profileFound := models.OSPFInterfaceProfileFromContainer(cont)
		if interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out OSPF Interface Profile %s not found", rs.Primary.ID)
		}
		*interface_profile = *interface_profileFound
		return nil
	}
}

func testAccCheckAciOSPFInterfaceProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l3out_ospf_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			interface_profile := models.OSPFInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out OSPF Interface Profile %s Still exists", interface_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciOSPFInterfaceProfileAttributes(description string, interface_profile *models.OSPFInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != interface_profile.Description {
			return fmt.Errorf("Bad L3out OSPF Interface Profile Description %s", interface_profile.Description)
		}

		if "example" != interface_profile.Annotation {
			return fmt.Errorf("Bad L3out OSPF Interface Profile annotation %s", interface_profile.Annotation)
		}

		if "255" != interface_profile.AuthKeyId {
			return fmt.Errorf("Bad L3out OSPF Interface Profile auth_key_id %s", interface_profile.AuthKeyId)
		}

		if "simple" != interface_profile.AuthType {
			return fmt.Errorf("Bad L3out OSPF Interface Profile auth_type %s", interface_profile.AuthType)
		}

		if "example" != interface_profile.NameAlias {
			return fmt.Errorf("Bad L3out OSPF Interface Profile name_alias %s", interface_profile.NameAlias)
		}

		return nil
	}
}
