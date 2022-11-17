package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciHSRPGroupProfile_Basic(t *testing.T) {
	var hsrp_group_profile models.HSRPGroupProfile
	description := "hsrp_group_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciHSRPGroupProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciHSRPGroupProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupProfileExists("aci_l3out_hsrp_interface_group.test", &hsrp_group_profile),
					testAccCheckAciHSRPGroupProfileAttributes(description, &hsrp_group_profile),
				),
			},
		},
	})
}

func TestAccAciHSRPGroupProfile_update(t *testing.T) {
	var hsrp_group_profile models.HSRPGroupProfile
	description := "hsrp_group_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciHSRPGroupProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciHSRPGroupProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupProfileExists("aci_l3out_hsrp_interface_group.test", &hsrp_group_profile),
					testAccCheckAciHSRPGroupProfileAttributes(description, &hsrp_group_profile),
				),
			},
			{
				Config: testAccCheckAciHSRPGroupProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupProfileExists("aci_l3out_hsrp_interface_group.test", &hsrp_group_profile),
					testAccCheckAciHSRPGroupProfileAttributes(description, &hsrp_group_profile),
				),
			},
		},
	})
}

func testAccCheckAciHSRPGroupProfileConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_l3out_hsrp_interface_group" "test" {
		l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.example.id
		name                            = "one"
		annotation                      = "example"
		description                     = "%s"
		config_issues                   = "Secondary-vip-conflicts-if-ip"
		group_af                        = "ipv4"
		group_id                        = "20"
		group_name                      = "test"
		ip                              = "10.22.30.40"
		ip_obtain_mode                  = "admin"
		mac                             = "02:10:45:00:00:56"
		name_alias                      = "example"
	}  
	`, description)
}

func testAccCheckAciHSRPGroupProfileExists(name string, hsrp_group_profile *models.HSRPGroupProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("HSRP Group Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No HSRP Group Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		hsrp_group_profileFound := models.HSRPGroupProfileFromContainer(cont)
		if hsrp_group_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("HSRP Group Profile %s not found", rs.Primary.ID)
		}
		*hsrp_group_profile = *hsrp_group_profileFound
		return nil
	}
}

func testAccCheckAciHSRPGroupProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l3out_hsrp_interface_group" {
			cont, err := client.Get(rs.Primary.ID)
			hsrp_group_profile := models.HSRPGroupProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("HSRP Group Profile %s Still exists", hsrp_group_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciHSRPGroupProfileAttributes(description string, hsrp_group_profile *models.HSRPGroupProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != hsrp_group_profile.Description {
			return fmt.Errorf("Bad hsrp_group_profile Description %s", hsrp_group_profile.Description)
		}

		if "one" != hsrp_group_profile.Name {
			return fmt.Errorf("Bad hsrp_group_profile name %s", hsrp_group_profile.Name)
		}

		if "example" != hsrp_group_profile.Annotation {
			return fmt.Errorf("Bad hsrp_group_profile annotation %s", hsrp_group_profile.Annotation)
		}

		if "Secondary-vip-conflicts-if-ip" != hsrp_group_profile.ConfigIssues {
			return fmt.Errorf("Bad hsrp_group_profile config_issues %s", hsrp_group_profile.ConfigIssues)
		}

		if "ipv4" != hsrp_group_profile.GroupAf {
			return fmt.Errorf("Bad hsrp_group_profile group_af %s", hsrp_group_profile.GroupAf)
		}

		if "20" != hsrp_group_profile.GroupId {
			return fmt.Errorf("Bad hsrp_group_profile group_id %s", hsrp_group_profile.GroupId)
		}

		if "test" != hsrp_group_profile.GroupName {
			return fmt.Errorf("Bad hsrp_group_profile group_name %s", hsrp_group_profile.GroupName)
		}

		if "10.22.30.40" != hsrp_group_profile.Ip {
			return fmt.Errorf("Bad hsrp_group_profile ip %s", hsrp_group_profile.Ip)
		}

		if "admin" != hsrp_group_profile.IpObtainMode {
			return fmt.Errorf("Bad hsrp_group_profile ip_obtain_mode %s", hsrp_group_profile.IpObtainMode)
		}

		if "02:10:45:00:00:56" != hsrp_group_profile.Mac {
			return fmt.Errorf("Bad hsrp_group_profile mac %s", hsrp_group_profile.Mac)
		}

		if "example" != hsrp_group_profile.NameAlias {
			return fmt.Errorf("Bad hsrp_group_profile name_alias %s", hsrp_group_profile.NameAlias)
		}

		return nil
	}
}
