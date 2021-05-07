package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL3outVPCMember_Basic(t *testing.T) {
	var l3out_vpc_member models.L3outVPCMember
	description := "member_node_configuration created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outVPCMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outVPCMemberConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outVPCMemberExists("aci_l3out_vpc_member.fool3out_vpc_member", &l3out_vpc_member),
					testAccCheckAciL3outVPCMemberAttributes(description, &l3out_vpc_member),
				),
			},
		},
	})
}

func TestAccAciL3outVPCMember_update(t *testing.T) {
	var l3out_vpc_member models.L3outVPCMember
	description := "member_node_configuration created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outVPCMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outVPCMemberConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outVPCMemberExists("aci_l3out_vpc_member.fool3out_vpc_member", &l3out_vpc_member),
					testAccCheckAciL3outVPCMemberAttributes(description, &l3out_vpc_member),
				),
			},
			{
				Config: testAccCheckAciL3outVPCMemberConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outVPCMemberExists("aci_l3out_vpc_member.fool3out_vpc_member", &l3out_vpc_member),
					testAccCheckAciL3outVPCMemberAttributes(description, &l3out_vpc_member),
				),
			},
		},
	})
}

func testAccCheckAciL3outVPCMemberConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3out_vpc_member" "fool3out_vpc_member" {
		leaf_port_dn  = "${aci_l3out_path_attachment.example.id}"
		description = "%s"
		side  = "A"
  		addr  = "10.0.0.1/16"
  		annotation  = "example"
  		ipv6_dad = "disabled"
  		ll_addr  = "::"
  		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciL3outVPCMemberExists(name string, l3out_vpc_member *models.L3outVPCMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out VPC Member %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out VPC Member dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_vpc_memberFound := models.L3outVPCMemberFromContainer(cont)
		if l3out_vpc_memberFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out VPC Member %s not found", rs.Primary.ID)
		}
		*l3out_vpc_member = *l3out_vpc_memberFound
		return nil
	}
}

func testAccCheckAciL3outVPCMemberDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_member_node_configuration" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_vpc_member := models.L3outVPCMemberFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out VPC Member %s Still exists", l3out_vpc_member.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3outVPCMemberAttributes(description string, l3out_vpc_member *models.L3outVPCMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3out_vpc_member.Description {
			return fmt.Errorf("Bad l3out_vpc_member Description %s", l3out_vpc_member.Description)
		}

		if "A" != l3out_vpc_member.Side {
			return fmt.Errorf("Bad l3out_vpc_member side %s", l3out_vpc_member.Side)
		}

		if "10.0.0.1/16" != l3out_vpc_member.Addr {
			return fmt.Errorf("Bad l3out_vpc_member addr %s", l3out_vpc_member.Addr)
		}

		if "example" != l3out_vpc_member.Annotation {
			return fmt.Errorf("Bad l3out_vpc_member annotation %s", l3out_vpc_member.Annotation)
		}

		if "disabled" != l3out_vpc_member.Ipv6Dad {
			return fmt.Errorf("Bad l3out_vpc_member ipv6_dad %s", l3out_vpc_member.Ipv6Dad)
		}

		if "::" != l3out_vpc_member.LlAddr {
			return fmt.Errorf("Bad l3out_vpc_member ll_addr %s", l3out_vpc_member.LlAddr)
		}

		if "example" != l3out_vpc_member.NameAlias {
			return fmt.Errorf("Bad l3out_vpc_member name_alias %s", l3out_vpc_member.NameAlias)
		}

		return nil
	}
}
