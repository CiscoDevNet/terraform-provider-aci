package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciInterfaceProfile_Basic(t *testing.T) {
	var interface_profile models.InterfaceProfile
	annotation := "port_selector"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInterfaceProfileConfig_basic(annotation),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceProfileExists("aci_spine_interface_profile_selector.foointerface_profile", &interface_profile),
					testAccCheckAciInterfaceProfileAttributes(annotation, &interface_profile),
				),
			},
		},
	})
}

func TestAccAciInterfaceProfile_update(t *testing.T) {
	var interface_profile models.InterfaceProfile
	annotation := "interface_profile_selector"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInterfaceProfileConfig_basic(annotation),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceProfileExists("aci_spine_interface_profile_selector.foointerface_profile", &interface_profile),
					testAccCheckAciInterfaceProfileAttributes(annotation, &interface_profile),
				),
			},
			{
				Config: testAccCheckAciInterfaceProfileConfig_basic(annotation),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceProfileExists("aci_spine_interface_profile_selector.foointerface_profile", &interface_profile),
					testAccCheckAciInterfaceProfileAttributes(annotation, &interface_profile),
				),
			},
		},
	})
}

func testAccCheckAciInterfaceProfileConfig_basic(annotation string) string {
	return fmt.Sprintf(`
	resource "aci_spine_interface_profile_selector" "foointerface_profile" {
		spine_profile_dn = aci_spine_profile.foospine_profile.id
		tdn              = aci_spine_interface_profile.foospine_interface_profile.id
		annotation       = "%s"
	}
	`, annotation)
}

func testAccCheckAciInterfaceProfileExists(name string, interface_profile *models.InterfaceProfile) resource.TestCheckFunc {
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

		interface_profileFound := models.InterfaceProfileFromContainer(cont)
		if interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Interface Profile %s not found", rs.Primary.ID)
		}
		*interface_profile = *interface_profileFound
		return nil
	}
}

func testAccCheckAciInterfaceProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_spine_interface_profile_selector" {
			cont, err := client.Get(rs.Primary.ID)
			interface_profile := models.InterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Interface Profile %s Still exists", interface_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciInterfaceProfileAttributes(annotation string, interface_profile *models.InterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if annotation != interface_profile.Annotation {
			return fmt.Errorf("Bad port_selector annotation %s", interface_profile.Annotation)
		}

		return nil
	}
}
