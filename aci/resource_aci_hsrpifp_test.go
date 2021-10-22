package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL3outHSRPInterfaceProfile_Basic(t *testing.T) {
	var l3out_hsrp_interface_profile models.L3outHSRPInterfaceProfile
	description := "interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHSRPInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outHSRPInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPInterfaceProfileExists("aci_l3out_hsrp_interface_profile.fool3out_hsrp_interface_profile", &l3out_hsrp_interface_profile),
					testAccCheckAciL3outHSRPInterfaceProfileAttributes(description, &l3out_hsrp_interface_profile),
				),
			},
		},
	})
}

func TestAccAciL3outHSRPInterfaceProfile_update(t *testing.T) {
	var l3out_hsrp_interface_profile models.L3outHSRPInterfaceProfile
	description := "interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHSRPInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outHSRPInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPInterfaceProfileExists("aci_l3out_hsrp_interface_profile.fool3out_hsrp_interface_profile", &l3out_hsrp_interface_profile),
					testAccCheckAciL3outHSRPInterfaceProfileAttributes(description, &l3out_hsrp_interface_profile),
				),
			},
			{
				Config: testAccCheckAciL3outHSRPInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPInterfaceProfileExists("aci_l3out_hsrp_interface_profile.fool3out_hsrp_interface_profile", &l3out_hsrp_interface_profile),
					testAccCheckAciL3outHSRPInterfaceProfileAttributes(description, &l3out_hsrp_interface_profile),
				),
			},
		},
	})
}

func testAccCheckAciL3outHSRPInterfaceProfileConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_l3out_hsrp_interface_profile" "fool3out_hsrp_interface_profile" {
		logical_interface_profile_dn = aci_logical_interface_profile.example.id
		description = "%s"
		annotation  = "example"
		name_alias  = "example"
		version     = "v1"
	}
	`, description)
}

func testAccCheckAciL3outHSRPInterfaceProfileExists(name string, l3out_hsrp_interface_profile *models.L3outHSRPInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out HSRP Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out HSRP Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_hsrp_interface_profileFound := models.L3outHSRPInterfaceProfileFromContainer(cont)
		if l3out_hsrp_interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out HSRP Interface Profile %s not found", rs.Primary.ID)
		}
		*l3out_hsrp_interface_profile = *l3out_hsrp_interface_profileFound
		return nil
	}
}

func testAccCheckAciL3outHSRPInterfaceProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_hsrp_interface_profile := models.L3outHSRPInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out HSRP Interface Profile %s Still exists", l3out_hsrp_interface_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3outHSRPInterfaceProfileAttributes(description string, l3out_hsrp_interface_profile *models.L3outHSRPInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3out_hsrp_interface_profile.Description {
			return fmt.Errorf("Bad l3out_hsrp_interface_profile Description %s", l3out_hsrp_interface_profile.Description)
		}

		if "example" != l3out_hsrp_interface_profile.Annotation {
			return fmt.Errorf("Bad l3out_hsrp_interface_profile annotation %s", l3out_hsrp_interface_profile.Annotation)
		}

		if "example" != l3out_hsrp_interface_profile.NameAlias {
			return fmt.Errorf("Bad l3out_hsrp_interface_profile name_alias %s", l3out_hsrp_interface_profile.NameAlias)
		}

		if "v1" != l3out_hsrp_interface_profile.Version {
			return fmt.Errorf("Bad l3out_hsrp_interface_profile version %s", l3out_hsrp_interface_profile.Version)
		}

		return nil
	}
}
