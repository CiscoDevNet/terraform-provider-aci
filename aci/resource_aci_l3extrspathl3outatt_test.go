package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL3outPathAttachment_Basic(t *testing.T) {
	var l3out_path_attachment models.L3outPathAttachment
	description := "leaf_port created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outPathAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outPathAttachmentConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists("aci_l3out_path_attachment.fool3out_path_attachment", &l3out_path_attachment),
					testAccCheckAciL3outPathAttachmentAttributes(description, &l3out_path_attachment),
				),
			},
		},
	})
}

func TestAccAciL3outPathAttachment_update(t *testing.T) {
	var l3out_path_attachment models.L3outPathAttachment
	description := "leaf_port created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outPathAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outPathAttachmentConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists("aci_l3out_path_attachment.fool3out_path_attachment", &l3out_path_attachment),
					testAccCheckAciL3outPathAttachmentAttributes(description, &l3out_path_attachment),
				),
			},
			{
				Config: testAccCheckAciL3outPathAttachmentConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outPathAttachmentExists("aci_l3out_path_attachment.fool3out_path_attachment", &l3out_path_attachment),
					testAccCheckAciL3outPathAttachmentAttributes(description, &l3out_path_attachment),
				),
			},
		},
	})
}

func testAccCheckAciL3outPathAttachmentConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		description = "sample aci_tenant from terraform"
		name        = "tenant_1"
		annotation  = "tenant_1_tag"
		name_alias  = "alias_tenant"
	  }
			
	  resource "aci_l3_outside" "fool3_outside" {
		tenant_dn      = aci_tenant.footenant.id
		description    = "sample aci_l3_outside"
		name           = "l3_outside_1"
		annotation     = "l3_outside_1_tag"
		enforce_rtctrl = "export"
		name_alias     = "alias_out"
		target_dscp    = "unspecified"
	  }
	  
	  resource "aci_logical_node_profile" "foological_node_profile" {
		l3_outside_dn = aci_l3_outside.fool3_outside.id
		description   = "sample logical node profile"
		name          = "logical_node_profile_1"
		annotation    = "logical_node_profile_1_tag"
		config_issues = "none"
		name_alias    = "alias_node"
		tag           = "black"
		target_dscp   = "unspecified"
	  }
		  
	  resource "aci_logical_interface_profile" "foological_interface_profile" {
		logical_node_profile_dn = aci_logical_node_profile.foological_node_profile.id
		description             = "aci_logical_interface_profile from terraform"
		name                    = "logical_interface_profile_1"
		annotation              = "logical_interface_profile_1_tag"
		name_alias              = "alias_prof"
		prio                    = "unspecified"
		tag                     = "black"
	  }	  
	  
	resource "aci_l3out_path_attachment" "fool3out_path_attachment" {
		logical_interface_profile_dn  = aci_logical_interface_profile.foological_interface_profile.id
		description = "%s"
		target_dn  = "topology/pod-1/paths-101/pathep-[eth1/1]"
  		addr  = "0.0.0.0"
  		annotation  = "example"
  		autostate = "disabled"
  		encap  = "vlan-1"
  		encap_scope = "ctx"
  		if_inst_t = "ext-svi"
  		ipv6_dad = "disabled"
  		ll_addr  = "::"
  		mac  = "00:22:BD:F8:19:FF"
  		mode = "native"
  		mtu = "inherit"
  		target_dscp = "AF11"
	}
	`, description)
}

func testAccCheckAciL3outPathAttachmentExists(name string, l3out_path_attachment *models.L3outPathAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3outPathAttachment %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3outPathAttachment dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_path_attachmentFound := models.L3outPathAttachmentFromContainer(cont)
		if l3out_path_attachmentFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3outPathAttachment %s not found", rs.Primary.ID)
		}
		*l3out_path_attachment = *l3out_path_attachmentFound
		return nil
	}
}

func testAccCheckAciL3outPathAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_leaf_port" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_path_attachment := models.L3outPathAttachmentFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3outPathAttachment %s Still exists", l3out_path_attachment.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3outPathAttachmentAttributes(description string, l3out_path_attachment *models.L3outPathAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3out_path_attachment.Description {
			return fmt.Errorf("Bad l3out_path_attachment Description %s", l3out_path_attachment.Description)
		}

		if "0.0.0.0" != l3out_path_attachment.Addr {
			return fmt.Errorf("Bad l3out_path_attachment addr %s", l3out_path_attachment.Addr)
		}

		if "example" != l3out_path_attachment.Annotation {
			return fmt.Errorf("Bad l3out_path_attachment annotation %s", l3out_path_attachment.Annotation)
		}

		if "disabled" != l3out_path_attachment.Autostate {
			return fmt.Errorf("Bad l3out_path_attachment autostate %s", l3out_path_attachment.Autostate)
		}

		if "vlan-1" != l3out_path_attachment.Encap {
			return fmt.Errorf("Bad l3out_path_attachment encap %s", l3out_path_attachment.Encap)
		}

		if "ctx" != l3out_path_attachment.EncapScope {
			return fmt.Errorf("Bad l3out_path_attachment encap_scope %s", l3out_path_attachment.EncapScope)
		}

		if "ext-svi" != l3out_path_attachment.IfInstT {
			return fmt.Errorf("Bad l3out_path_attachment if_inst_t %s", l3out_path_attachment.IfInstT)
		}

		if "disabled" != l3out_path_attachment.Ipv6Dad {
			return fmt.Errorf("Bad l3out_path_attachment ipv6_dad %s", l3out_path_attachment.Ipv6Dad)
		}

		if "::" != l3out_path_attachment.LlAddr {
			return fmt.Errorf("Bad l3out_path_attachment ll_addr %s", l3out_path_attachment.LlAddr)
		}

		if "00:22:BD:F8:19:FF" != l3out_path_attachment.Mac {
			return fmt.Errorf("Bad l3out_path_attachment mac %s", l3out_path_attachment.Mac)
		}

		if "native" != l3out_path_attachment.Mode {
			return fmt.Errorf("Bad l3out_path_attachment mode %s", l3out_path_attachment.Mode)
		}

		if "inherit" != l3out_path_attachment.Mtu {
			return fmt.Errorf("Bad l3out_path_attachment mtu %s", l3out_path_attachment.Mtu)
		}

		if "AF11" != l3out_path_attachment.TargetDscp {
			return fmt.Errorf("Bad l3out_path_attachment target_dscp %s", l3out_path_attachment.TargetDscp)
		}

		return nil
	}
}
