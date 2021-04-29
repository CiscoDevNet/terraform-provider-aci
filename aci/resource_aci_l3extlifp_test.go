package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLogicalInterfaceProfile_Basic(t *testing.T) {
	var logical_interface_profile models.LogicalInterfaceProfile
	description := "logical_interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLogicalInterfaceProfileConfig_basic(description, "black"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceProfileExists("aci_logical_interface_profile.foological_interface_profile", &logical_interface_profile),
					testAccCheckAciLogicalInterfaceProfileAttributes(description, "black", &logical_interface_profile),
				),
			},
			{
				ResourceName:      "aci_logical_interface_profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciLogicalInterfaceProfile_update(t *testing.T) {
	var logical_interface_profile models.LogicalInterfaceProfile
	description := "logical_interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLogicalInterfaceProfileConfig_basic(description, "black"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceProfileExists("aci_logical_interface_profile.foological_interface_profile", &logical_interface_profile),
					testAccCheckAciLogicalInterfaceProfileAttributes(description, "black", &logical_interface_profile),
				),
			},
			{
				Config: testAccCheckAciLogicalInterfaceProfileConfig_basic(description, "white"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceProfileExists("aci_logical_interface_profile.foological_interface_profile", &logical_interface_profile),
					testAccCheckAciLogicalInterfaceProfileAttributes(description, "white", &logical_interface_profile),
				),
			},
		},
	})
}

func testAccCheckAciLogicalInterfaceProfileConfig_basic(description, tag string) string {
	return fmt.Sprintf(`

	resource "aci_logical_interface_profile" "foological_interface_profile" {
		logical_node_profile_dn = "${aci_logical_node_profile.example.id}"
		description             = "%s"
		name                    = "demo_int_prof"
		annotation              = "tag_prof"
		name_alias              = "alias_prof"
		prio                    = "unspecified"
		tag                     = "%s"
	  }	  
	`, description, tag)
}

func testAccCheckAciLogicalInterfaceProfileExists(name string, logical_interface_profile *models.LogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Logical Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Logical Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		logical_interface_profileFound := models.LogicalInterfaceProfileFromContainer(cont)
		if logical_interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Logical Interface Profile %s not found", rs.Primary.ID)
		}
		*logical_interface_profile = *logical_interface_profileFound
		return nil
	}
}

func testAccCheckAciLogicalInterfaceProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_logical_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			logical_interface_profile := models.LogicalInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Logical Interface Profile %s Still exists", logical_interface_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLogicalInterfaceProfileAttributes(description, tag string, logical_interface_profile *models.LogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != logical_interface_profile.Description {
			return fmt.Errorf("Bad logical_interface_profile Description %s", logical_interface_profile.Description)
		}

		if "demo_int_prof" != logical_interface_profile.Name {
			return fmt.Errorf("Bad logical_interface_profile name %s", logical_interface_profile.Name)
		}

		if "tag_prof" != logical_interface_profile.Annotation {
			return fmt.Errorf("Bad logical_interface_profile annotation %s", logical_interface_profile.Annotation)
		}

		if "alias_prof" != logical_interface_profile.NameAlias {
			return fmt.Errorf("Bad logical_interface_profile name_alias %s", logical_interface_profile.NameAlias)
		}

		if "unspecified" != logical_interface_profile.Prio {
			return fmt.Errorf("Bad logical_interface_profile prio %s", logical_interface_profile.Prio)
		}

		if tag != logical_interface_profile.Tag {
			return fmt.Errorf("Bad logical_interface_profile tag %s", logical_interface_profile.Tag)
		}

		return nil
	}
}
