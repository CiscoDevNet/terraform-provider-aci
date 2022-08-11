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

func TestAccAciOverridePolicyGroup_Basic(t *testing.T) {
	var override_policy_group models.OverridePolicyGroup
	infra_acc_bndl_grp_name := acctest.RandString(5)
	infra_acc_bndl_subgrp_name := acctest.RandString(5)
	description := "override_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOverridePolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOverridePolicyGroupConfig_basic(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOverridePolicyGroupExists("aci_leaf_access_bundle_policy_sub_group.test_policy_sub_group", &override_policy_group),
					testAccCheckAciOverridePolicyGroupAttributes(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name, description, &override_policy_group),
				),
			},
		},
	})
}

func TestAccAciOverridePolicyGroup_Update(t *testing.T) {
	var override_policy_group models.OverridePolicyGroup
	infra_acc_bndl_grp_name := acctest.RandString(5)
	infra_acc_bndl_subgrp_name := acctest.RandString(5)
	description := "override_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOverridePolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOverridePolicyGroupConfig_basic(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOverridePolicyGroupExists("aci_leaf_access_bundle_policy_sub_group.test_policy_sub_group", &override_policy_group),
					testAccCheckAciOverridePolicyGroupAttributes(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name, description, &override_policy_group),
				),
			},
			{
				Config: testAccCheckAciOverridePolicyGroupConfig_basic(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOverridePolicyGroupExists("aci_leaf_access_bundle_policy_sub_group.test_policy_sub_group", &override_policy_group),
					testAccCheckAciOverridePolicyGroupAttributes(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name, description, &override_policy_group),
				),
			},
		},
	})
}

func testAccCheckAciOverridePolicyGroupConfig_basic(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name string) string {
	return fmt.Sprintf(`

	resource "aci_leaf_access_bundle_policy_group" "test_policy_group" {
		name 		= "%s"
		description = "aci_leaf_access_bundle_policy_group created while acceptance testing"

	}

	resource "aci_leaf_access_bundle_policy_sub_group" "test_policy_sub_group" {
		name 		= "%s"
		description = "override_policy_group created while acceptance testing"
		leaf_access_bundle_policy_group_dn = aci_leaf_access_bundle_policy_group.test_policy_group.id
	}

	`, infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name)
}

func testAccCheckAciOverridePolicyGroupExists(name string, override_policy_group *models.OverridePolicyGroup) resource.TestCheckFunc {
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

		override_policy_groupFound := models.OverridePolicyGroupFromContainer(cont)
		if override_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Override Policy Group %s not found", rs.Primary.ID)
		}
		*override_policy_group = *override_policy_groupFound
		return nil
	}
}

func testAccCheckAciOverridePolicyGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_leaf_access_bundle_policy_sub_group" {
			cont, err := client.Get(rs.Primary.ID)
			override_policy_group := models.OverridePolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Override Policy Group %s Still exists", override_policy_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciOverridePolicyGroupAttributes(infra_acc_bndl_grp_name, infra_acc_bndl_subgrp_name, description string, override_policy_group *models.OverridePolicyGroup) resource.TestCheckFunc {
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
