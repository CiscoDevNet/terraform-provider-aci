package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLoopBackInterfaceProfile_Basic(t *testing.T) {
	var loop_back_interface_profile models.LoopBackInterfaceProfile
	description := "L3out Loopback interface profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLoopBackInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLoopBackInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoopBackInterfaceProfileExists("aci_l3out_loopback_interface_profile.test", &loop_back_interface_profile),
					testAccCheckAciLoopBackInterfaceProfileAttributes(description, &loop_back_interface_profile),
				),
			},
		},
	})
}

func TestAccAciLoopBackInterfaceProfile_update(t *testing.T) {
	var loop_back_interface_profile models.LoopBackInterfaceProfile
	description := "L3out Loopback interface profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLoopBackInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLoopBackInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoopBackInterfaceProfileExists("aci_l3out_loopback_interface_profile.test", &loop_back_interface_profile),
					testAccCheckAciLoopBackInterfaceProfileAttributes(description, &loop_back_interface_profile),
				),
			},
			{
				Config: testAccCheckAciLoopBackInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoopBackInterfaceProfileExists("aci_l3out_loopback_interface_profile.test", &loop_back_interface_profile),
					testAccCheckAciLoopBackInterfaceProfileAttributes(description, &loop_back_interface_profile),
				),
			},
		},
	})
}

func testAccCheckAciLoopBackInterfaceProfileConfig_basic(description string) string {
	return fmt.Sprintf(`	

	resource "aci_l3out_loopback_interface_profile" "test" {
		fabric_node_dn = aci_logical_node_to_fabric_node.example.id
		addr           = "1.2.3.5"
		description    = "%s"
		annotation     = "example"
		name_alias     = "example"
	}
	`, description)
}

func testAccCheckAciLoopBackInterfaceProfileExists(name string, loop_back_interface_profile *models.LoopBackInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Loopback Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Loopback Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		loop_back_interface_profileFound := models.LoopBackInterfaceProfileFromContainer(cont)
		if loop_back_interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Loopback Interface Profile %s not found", rs.Primary.ID)
		}
		*loop_back_interface_profile = *loop_back_interface_profileFound
		return nil
	}
}

func testAccCheckAciLoopBackInterfaceProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l3out_loopback_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			loop_back_interface_profile := models.LoopBackInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Loopback Interface Profile %s Still exists", loop_back_interface_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLoopBackInterfaceProfileAttributes(description string, loop_back_interface_profile *models.LoopBackInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != loop_back_interface_profile.Description {
			return fmt.Errorf("Bad loop_back_interface_profile Description %s", loop_back_interface_profile.Description)
		}

		if "1.2.3.5" != loop_back_interface_profile.Addr {
			return fmt.Errorf("Bad loop_back_interface_profile addr %s", loop_back_interface_profile.Addr)
		}

		if "example" != loop_back_interface_profile.Annotation {
			return fmt.Errorf("Bad loop_back_interface_profile annotation %s", loop_back_interface_profile.Annotation)
		}

		if "example" != loop_back_interface_profile.NameAlias {
			return fmt.Errorf("Bad loop_back_interface_profile name_alias %s", loop_back_interface_profile.NameAlias)
		}

		return nil
	}
}
