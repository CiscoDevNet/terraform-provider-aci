package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciL3outPathAttachmentSecondaryIp_Basic(t *testing.T) {
	var l3out_path_attachment_secondary_ip models.L3outPathAttachmentSecondaryIp
	description := "secondary_i_paddress created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outPathAttachmentSecondaryIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outPathAttachmentSecondaryIpConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentSecondaryIpExists("aci_l3out_path_attachment_secondary_ip.fool3out_path_attachment_secondary_ip", &l3out_path_attachment_secondary_ip),
					testAccCheckAciL3outPathAttachmentSecondaryIpAttributes(description, &l3out_path_attachment_secondary_ip),
				),
			},
		},
	})
}

func TestAccAciL3outPathAttachmentSecondaryIp_update(t *testing.T) {
	var l3out_path_attachment_secondary_ip models.L3outPathAttachmentSecondaryIp
	description := "secondary_i_paddress created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outPathAttachmentSecondaryIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outPathAttachmentSecondaryIpConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentSecondaryIpExists("aci_l3out_path_attachment_secondary_ip.fool3out_path_attachment_secondary_ip", &l3out_path_attachment_secondary_ip),
					testAccCheckAciL3outPathAttachmentSecondaryIpAttributes(description, &l3out_path_attachment_secondary_ip),
				),
			},
			{
				Config: testAccCheckAciL3outPathAttachmentSecondaryIpConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentSecondaryIpExists("aci_l3out_path_attachment_secondary_ip.fool3out_path_attachment_secondary_ip", &l3out_path_attachment_secondary_ip),
					testAccCheckAciL3outPathAttachmentSecondaryIpAttributes(description, &l3out_path_attachment_secondary_ip),
				),
			},
		},
	})
}

func testAccCheckAciL3outPathAttachmentSecondaryIpConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3out_path_attachment_secondary_ip" "fool3out_path_attachment_secondary_ip" {
		#l3out_path_attachment_dn  = "${aci_leaf_port.example.id}"
		l3out_path_attachment_dn  = "uni/tn-check_tenantnk/out-crest_test_rutvik_l3out/lnodep-crest_test_rutvik_node/lifp-crest_test_rutvik_int_prof/rspathL3OutAtt-[topology/pod-1/paths-101/pathep-[eth1/26]]"
		description = "%s"
		addr  = "10.0.0.1/24"
  		annotation  = "example"
  		ipv6_dad = "disabled"
  		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciL3outPathAttachmentSecondaryIpExists(name string, l3out_path_attachment_secondary_ip *models.L3outPathAttachmentSecondaryIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3outPathAttachmentSecondaryIp %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3outPathAttachmentSecondaryIp dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_path_attachment_secondary_ipFound := models.L3outPathAttachmentSecondaryIpFromContainer(cont)
		if l3out_path_attachment_secondary_ipFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3outPathAttachmentSecondaryIp %s not found", rs.Primary.ID)
		}
		*l3out_path_attachment_secondary_ip = *l3out_path_attachment_secondary_ipFound
		return nil
	}
}

func testAccCheckAciL3outPathAttachmentSecondaryIpDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_secondary_i_paddress" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_path_attachment_secondary_ip := models.L3outPathAttachmentSecondaryIpFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3outPathAttachmentSecondaryIp %s Still exists", l3out_path_attachment_secondary_ip.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3outPathAttachmentSecondaryIpAttributes(description string, l3out_path_attachment_secondary_ip *models.L3outPathAttachmentSecondaryIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3out_path_attachment_secondary_ip.Description {
			return fmt.Errorf("Bad l3out_path_attachment_secondary_ip Description %s", l3out_path_attachment_secondary_ip.Description)
		}

		if "10.0.0.1/24" != l3out_path_attachment_secondary_ip.Addr {
			return fmt.Errorf("Bad l3out_path_attachment_secondary_ip addr %s", l3out_path_attachment_secondary_ip.Addr)
		}

		if "example" != l3out_path_attachment_secondary_ip.Annotation {
			return fmt.Errorf("Bad l3out_path_attachment_secondary_ip annotation %s", l3out_path_attachment_secondary_ip.Annotation)
		}

		if "disabled" != l3out_path_attachment_secondary_ip.Ipv6Dad {
			return fmt.Errorf("Bad l3out_path_attachment_secondary_ip ipv6_dad %s", l3out_path_attachment_secondary_ip.Ipv6Dad)
		}

		if "example" != l3out_path_attachment_secondary_ip.NameAlias {
			return fmt.Errorf("Bad l3out_path_attachment_secondary_ip name_alias %s", l3out_path_attachment_secondary_ip.NameAlias)
		}

		return nil
	}
}
