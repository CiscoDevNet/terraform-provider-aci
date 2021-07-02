package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciAutonomousSystemProfile_Basic(t *testing.T) {
	var autonomous_system_profile models.AutonomousSystemProfile
	description := "autonomous_system_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAutonomousSystemProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAutonomousSystemProfileConfig_basic(description, "121"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAutonomousSystemProfileExists("aci_autonomous_system_profile.fooautonomous_system_profile", &autonomous_system_profile),
					testAccCheckAciAutonomousSystemProfileAttributes(description, "121", &autonomous_system_profile),
				),
			},
			{
				ResourceName:      "aci_autonomous_system_profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciAutonomousSystemProfile_update(t *testing.T) {
	var autonomous_system_profile models.AutonomousSystemProfile
	description := "autonomous_system_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAutonomousSystemProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAutonomousSystemProfileConfig_basic(description, "121"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAutonomousSystemProfileExists("aci_autonomous_system_profile.fooautonomous_system_profile", &autonomous_system_profile),
					testAccCheckAciAutonomousSystemProfileAttributes(description, "121", &autonomous_system_profile),
				),
			},
			{
				Config: testAccCheckAciAutonomousSystemProfileConfig_basic(description, "131"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAutonomousSystemProfileExists("aci_autonomous_system_profile.fooautonomous_system_profile", &autonomous_system_profile),
					testAccCheckAciAutonomousSystemProfileAttributes(description, "131", &autonomous_system_profile),
				),
			},
		},
	})
}

func testAccCheckAciAutonomousSystemProfileConfig_basic(description, asn string) string {
	return fmt.Sprintf(`

	resource "aci_autonomous_system_profile" "fooautonomous_system_profile" {
		description = "%s"
		annotation  = "tag_system"
		asn         = "%s"
		name_alias  = "alias_sys_prof"
	}   
	`, description, asn)
}

func testAccCheckAciAutonomousSystemProfileExists(name string, autonomous_system_profile *models.AutonomousSystemProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Autonomous System Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Autonomous System Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		autonomous_system_profileFound := models.AutonomousSystemProfileFromContainer(cont)
		if autonomous_system_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Autonomous System Profile %s not found", rs.Primary.ID)
		}
		*autonomous_system_profile = *autonomous_system_profileFound
		return nil
	}
}

func testAccCheckAciAutonomousSystemProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_autonomous_system_profile" {
			cont, err := client.Get(rs.Primary.ID)
			autonomous_system_profile := models.AutonomousSystemProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Autonomous System Profile %s Still exists", autonomous_system_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciAutonomousSystemProfileAttributes(description, asn string, autonomous_system_profile *models.AutonomousSystemProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != autonomous_system_profile.Description {
			return fmt.Errorf("Bad autonomous_system_profile Description %s", autonomous_system_profile.Description)
		}

		if "tag_system" != autonomous_system_profile.Annotation {
			return fmt.Errorf("Bad autonomous_system_profile annotation %s", autonomous_system_profile.Annotation)
		}

		if asn != autonomous_system_profile.Asn {
			return fmt.Errorf("Bad autonomous_system_profile asn %s", autonomous_system_profile.Asn)
		}

		if "alias_sys_prof" != autonomous_system_profile.NameAlias {
			return fmt.Errorf("Bad autonomous_system_profile name_alias %s", autonomous_system_profile.NameAlias)
		}

		return nil
	}
}
