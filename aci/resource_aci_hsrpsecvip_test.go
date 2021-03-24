package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciL3outHSRPSecondaryVIP_Basic(t *testing.T) {
	var l3out_hsrp_secondary_vip models.L3outHSRPSecondaryVIP
	description := "secondary_virtual_ip_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHSRPSecondaryVIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outHSRPSecondaryVIPConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists("aci_l3out_hsrp_secondary_vip.fool3out_hsrp_secondary_vip", &l3out_hsrp_secondary_vip),
					testAccCheckAciL3outHSRPSecondaryVIPAttributes(description, &l3out_hsrp_secondary_vip),
				),
			},
		},
	})
}

func TestAccAciL3outHSRPSecondaryVIP_update(t *testing.T) {
	var l3out_hsrp_secondary_vip models.L3outHSRPSecondaryVIP
	description := "secondary_virtual_ip_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outHSRPSecondaryVIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outHSRPSecondaryVIPConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists("aci_l3out_hsrp_secondary_vip.fool3out_hsrp_secondary_vip", &l3out_hsrp_secondary_vip),
					testAccCheckAciL3outHSRPSecondaryVIPAttributes(description, &l3out_hsrp_secondary_vip),
				),
			},
			{
				Config: testAccCheckAciL3outHSRPSecondaryVIPConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outHSRPSecondaryVIPExists("aci_l3out_hsrp_secondary_vip.fool3out_hsrp_secondary_vip", &l3out_hsrp_secondary_vip),
					testAccCheckAciL3outHSRPSecondaryVIPAttributes(description, &l3out_hsrp_secondary_vip),
				),
			},
		},
	})
}

func testAccCheckAciL3outHSRPSecondaryVIPConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3out_hsrp_secondary_vip" "fool3out_hsrp_secondary_vip" {
		hsrp_group_profile_dn  = "${aci_hsrp_group_profile.example.id}"
		description = "%s"
		ip  = "example"
  		annotation  = "example"
  		config_issues = "GroupMac-Conflicts-Other-Group"
  		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciL3outHSRPSecondaryVIPExists(name string, l3out_hsrp_secondary_vip *models.L3outHSRPSecondaryVIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out HSRP Secondary VIP %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out HSRP Secondary VIP dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_hsrp_secondary_vipFound := models.L3outHSRPSecondaryVIPFromContainer(cont)
		if l3out_hsrp_secondary_vipFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out HSRP Secondary VIP %s not found", rs.Primary.ID)
		}
		*l3out_hsrp_secondary_vip = *l3out_hsrp_secondary_vipFound
		return nil
	}
}

func testAccCheckAciL3outHSRPSecondaryVIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_secondary_virtual_ip_policy" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_hsrp_secondary_vip := models.L3outHSRPSecondaryVIPFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out HSRP Secondary VIP %s Still exists", l3out_hsrp_secondary_vip.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3outHSRPSecondaryVIPAttributes(description string, l3out_hsrp_secondary_vip *models.L3outHSRPSecondaryVIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3out_hsrp_secondary_vip.Description {
			return fmt.Errorf("Bad l3out_hsrp_secondary_vip Description %s", l3out_hsrp_secondary_vip.Description)
		}

		if "example" != l3out_hsrp_secondary_vip.Ip {
			return fmt.Errorf("Bad l3out_hsrp_secondary_vip ip %s", l3out_hsrp_secondary_vip.Ip)
		}

		if "example" != l3out_hsrp_secondary_vip.Annotation {
			return fmt.Errorf("Bad l3out_hsrp_secondary_vip annotation %s", l3out_hsrp_secondary_vip.Annotation)
		}

		if "GroupMac-Conflicts-Other-Group" != l3out_hsrp_secondary_vip.ConfigIssues {
			return fmt.Errorf("Bad l3out_hsrp_secondary_vip config_issues %s", l3out_hsrp_secondary_vip.ConfigIssues)
		}

		if "example" != l3out_hsrp_secondary_vip.NameAlias {
			return fmt.Errorf("Bad l3out_hsrp_secondary_vip name_alias %s", l3out_hsrp_secondary_vip.NameAlias)
		}

		return nil
	}
}
