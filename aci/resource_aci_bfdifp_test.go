package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBFDInterfaceProfile_Basic(t *testing.T) {
	var interface_profile models.BFDInterfaceProfile
	description := "interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBFDInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBFDInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBFDInterfaceProfileExists("aci_l3out_bfd_interface_profile.test", &interface_profile),
					testAccCheckAciBFDInterfaceProfileAttributes(description, &interface_profile),
				),
			},
		},
	})
}

func TestAccAciBFDInterfaceProfile_update(t *testing.T) {
	var interface_profile models.BFDInterfaceProfile
	description := "interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBFDInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBFDInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBFDInterfaceProfileExists("aci_l3out_bfd_interface_profile.test", &interface_profile),
					testAccCheckAciBFDInterfaceProfileAttributes(description, &interface_profile),
				),
			},
			{
				Config: testAccCheckAciBFDInterfaceProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBFDInterfaceProfileExists("aci_l3out_bfd_interface_profile.test", &interface_profile),
					testAccCheckAciBFDInterfaceProfileAttributes(description, &interface_profile),
				),
			},
		},
	})
}

func testAccCheckAciBFDInterfaceProfileConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn = aci_logical_interface_profile.example.id
		annotation                   = "example"
		description                  = "%s"
		key                          = "example"
		key_id                       = "25"
		interface_profile_type       = "sha1"
	}
	`, description)
}

func testAccCheckAciBFDInterfaceProfileExists(name string, interface_profile *models.BFDInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		interface_profileFound := models.BFDInterfaceProfileFromContainer(cont)
		if interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Interface Profile %s not found", rs.Primary.ID)
		}
		*interface_profile = *interface_profileFound
		return nil
	}
}

func testAccCheckAciBFDInterfaceProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l3out_bfd_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			interface_profile := models.BFDInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Interface Profile %s Still exists", interface_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciBFDInterfaceProfileAttributes(description string, interface_profile *models.BFDInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != interface_profile.Description {
			return fmt.Errorf("Bad interface_profile Description %s", interface_profile.Description)
		}

		if "example" != interface_profile.Annotation {
			return fmt.Errorf("Bad interface_profile annotation %s", interface_profile.Annotation)
		}

		if "25" != interface_profile.KeyId {
			return fmt.Errorf("Bad interface_profile key_id %s", interface_profile.KeyId)
		}

		if "sha1" != interface_profile.InterfaceProfileType {
			return fmt.Errorf("Bad interface_profile interface_profile_type %s", interface_profile.InterfaceProfileType)
		}

		return nil
	}
}
