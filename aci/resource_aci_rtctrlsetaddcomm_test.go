package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciRtctrlSetAddComm_Basic(t *testing.T) {
	var rtctrl_set_add_comm models.RtctrlSetAddComm
	fv_tenant_name := acctest.RandString(5)
	rtctrl_attr_p_name := acctest.RandString(5)
	rtctrl_set_add_comm_name := acctest.RandString(5)
	description := "rtctrl_set_add_comm created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRtctrlSetAddCommDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRtctrlSetAddCommConfig_basic(fv_tenant_name, rtctrl_attr_p_name, rtctrl_set_add_comm_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRtctrlSetAddCommExists("aci_action_rule_profile_additional_communities.foo_rtctrl_set_add_comm", &rtctrl_set_add_comm),
					testAccCheckAciRtctrlSetAddCommAttributes(fv_tenant_name, rtctrl_attr_p_name, rtctrl_set_add_comm_name, description, &rtctrl_set_add_comm),
				),
			},
		},
	})
}

func TestAccAciRtctrlSetAddComm_Update(t *testing.T) {
	var rtctrl_set_add_comm models.RtctrlSetAddComm
	fv_tenant_name := acctest.RandString(5)
	rtctrl_attr_p_name := acctest.RandString(5)
	rtctrl_set_add_comm_name := acctest.RandString(5)
	description := "rtctrl_set_add_comm created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRtctrlSetAddCommDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRtctrlSetAddCommConfig_basic(fv_tenant_name, rtctrl_attr_p_name, rtctrl_set_add_comm_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRtctrlSetAddCommExists("aci_action_rule_profile_additional_communities.foo_rtctrl_set_add_comm", &rtctrl_set_add_comm),
					testAccCheckAciRtctrlSetAddCommAttributes(fv_tenant_name, rtctrl_attr_p_name, rtctrl_set_add_comm_name, description, &rtctrl_set_add_comm),
				),
			},
			{
				Config: testAccCheckAciRtctrlSetAddCommConfig_basic(fv_tenant_name, rtctrl_attr_p_name, rtctrl_set_add_comm_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRtctrlSetAddCommExists("aci_action_rule_profile_additional_communities.foo_rtctrl_set_add_comm", &rtctrl_set_add_comm),
					testAccCheckAciRtctrlSetAddCommAttributes(fv_tenant_name, rtctrl_attr_p_name, rtctrl_set_add_comm_name, description, &rtctrl_set_add_comm),
				),
			},
		},
	})
}

func testAccCheckAciRtctrlSetAddCommConfig_basic(fv_tenant_name, rtctrl_attr_p_name, rtctrl_set_add_comm_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name        = "%s"
		description = "tenant created while acceptance testing"
	}

	resource "aci_action_rule_profile" "foo_action_rule_profile" {
		name        = "%s"
		description = "action_rule_profile created while acceptance testing"
		tenant_dn   = aci_tenant.foo_tenant.id
	}

	resource "aci_action_rule_profile_additional_communities" "foo_rtctrl_set_add_comm" {
		name                   = "%s"
		community              = "no-advertise"
		description            = "rtctrl_set_add_comm created while acceptance testing"
		action_rule_profile_dn = aci_action_rule_profile.foo_action_rule_profile.id
	}

	`, fv_tenant_name, rtctrl_attr_p_name, rtctrl_set_add_comm_name)
}

func testAccCheckAciRtctrlSetAddCommExists(name string, rtctrl_set_add_comm *models.RtctrlSetAddComm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("RtctrlSetAddComm %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No RtctrlSetAddComm dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		rtctrl_set_add_commFound := models.RtctrlSetAddCommFromContainer(cont)
		if rtctrl_set_add_commFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("RtctrlSetAddComm %s not found", rs.Primary.ID)
		}
		*rtctrl_set_add_comm = *rtctrl_set_add_commFound
		return nil
	}
}

func testAccCheckAciRtctrlSetAddCommDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_action_rule_profile_additional_communities" {
			cont, err := client.Get(rs.Primary.ID)
			rtctrl_set_add_comm := models.RtctrlSetAddCommFromContainer(cont)
			if err == nil {
				return fmt.Errorf("RtctrlSetAddComm %s Still exists", rtctrl_set_add_comm.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRtctrlSetAddCommAttributes(fv_tenant_name, rtctrl_attr_p_name, rtctrl_set_add_comm_name, description string, rtctrl_set_add_comm *models.RtctrlSetAddComm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if rtctrl_set_add_comm_name != GetMOName(rtctrl_set_add_comm.DistinguishedName) {
			return fmt.Errorf("Bad rtctrl_set_add_comm %s", GetMOName(rtctrl_set_add_comm.DistinguishedName))
		}

		if rtctrl_attr_p_name != GetMOName(GetParentDn(rtctrl_set_add_comm.DistinguishedName, rtctrl_set_add_comm.Rn)) {
			return fmt.Errorf(" Bad rtctrl_attr_p %s", GetMOName(GetParentDn(rtctrl_set_add_comm.DistinguishedName, rtctrl_set_add_comm.Rn)))
		}
		if description != rtctrl_set_add_comm.Description {
			return fmt.Errorf("Bad rtctrl_set_add_comm Description %s", rtctrl_set_add_comm.Description)
		}
		return nil
	}
}
