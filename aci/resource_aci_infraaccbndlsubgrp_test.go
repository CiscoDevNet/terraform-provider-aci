package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciOverridePCVPCPolicyGroup_Basic(t *testing.T) {
	var override_policy_group models.OverridePCVPCPolicyGroup
	infra_acc_bndl_grp_name := acctest.RandString(5)
	infra_acc_bndl_subgrp_name := acctest.RandString(5)
	description := "override_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOverridePCVPCPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOverridePCVPCPolicyGroupConfig_basic(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOverridePCVPCPolicyGroupExists("aci_leaf_access_bundle_policy_sub_group.test_policy_sub_group", &override_policy_group),
					testAccCheckAciOverridePCVPCPolicyGroupAttributes(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name, description, &override_policy_group),
				),
			},
		},
	})
}

func TestAccAciOverridePCVPCPolicyGroup_Update(t *testing.T) {
	var override_policy_group models.OverridePCVPCPolicyGroup
	infra_acc_bndl_grp_name := acctest.RandString(5)
	infra_acc_bndl_subgrp_name := acctest.RandString(5)
	description := "override_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOverridePCVPCPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOverridePCVPCPolicyGroupConfig_basic(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOverridePCVPCPolicyGroupExists("aci_leaf_access_bundle_policy_sub_group.test_policy_sub_group", &override_policy_group),
					testAccCheckAciOverridePCVPCPolicyGroupAttributes(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name, description, &override_policy_group),
				),
			},
			{
				Config: testAccCheckAciOverridePCVPCPolicyGroupConfig_basic(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOverridePCVPCPolicyGroupExists("aci_leaf_access_bundle_policy_sub_group.test_policy_sub_group", &override_policy_group),
					testAccCheckAciOverridePCVPCPolicyGroupAttributes(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name, description, &override_policy_group),
				),
			},
		},
	})
}

func testAccCheckAciOverridePCVPCPolicyGroupConfig_basic(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name string) string {
	return fmt.Sprintf(`

	resource "aci_lacp_member_policy" "test_lacp_member_policy" {
	    name = "test-policy"
	    description = "This policy member is created by terraform"
    }

	resource "aci_leaf_access_bundle_policy_group" "test_policy_group" {
		name 		= "%s"
		description = "aci_leaf_access_bundle_policy_group created while acceptance testing"
	}

	resource "aci_leaf_access_bundle_policy_sub_group" "test_policy_sub_group" {
	    leaf_access_bundle_policy_group_dn = aci_leaf_access_bundle_policy_group.test_policy_group.id
		name 		= "%s"
		description = "override_policy_group created while acceptance testing"
		port_channel_member = aci_lacp_member_policy.test_lacp_member_policy.id
	}

	`, infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name)
}

func testAccCheckAciOverridePCVPCPolicyGroupExists(name string, override_policy_group *models.OverridePCVPCPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Override Policy Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Override Policy Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		override_policy_groupFound := models.OverridePCVPCPolicyGroupFromContainer(cont)
		if override_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Override Policy Group %s not found", rs.Primary.ID)
		}
		*override_policy_group = *override_policy_groupFound
		return nil
	}
}

func testAccCheckAciOverridePCVPCPolicyGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_leaf_access_bundle_policy_sub_group" {
			cont, err := client.Get(rs.Primary.ID)
			override_policy_group := models.OverridePCVPCPolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Override Policy Group %s Still exists", override_policy_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciOverridePCVPCPolicyGroupAttributes(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name, description string, override_policy_group *models.OverridePCVPCPolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if infra_acc_bndl_subgrp_name != GetMOName(override_policy_group.DistinguishedName) {
			return fmt.Errorf("Bad infraacc_bndl_subgrp %s", GetMOName(override_policy_group.DistinguishedName))
		}

		if infra_acc_bndl_grp_name != GetMOName(GetParentDn(override_policy_group.DistinguishedName, fmt.Sprintf("/"+models.RninfraAccBndlSubgrp, infra_acc_bndl_subgrp_name))) {
			return fmt.Errorf(" Bad infraacc_bndl_grp %s ", GetMOName(GetParentDn(override_policy_group.DistinguishedName, fmt.Sprintf("/"+models.RninfraAccBndlSubgrp, infra_acc_bndl_subgrp_name))))
		}
		if description != override_policy_group.Description {
			return fmt.Errorf("Bad override_policy_group Description %s", override_policy_group.Description)
		}
		return nil
	}
}
