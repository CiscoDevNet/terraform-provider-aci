package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciPCVPCInterfacePolicyGroup_Basic(t *testing.T) {
	var pcvpc_interface_policy_group models.PCVPCInterfacePolicyGroup
	description := "pc/vpc_interface_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPCVPCInterfacePolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPCVPCInterfacePolicyGroupConfig_basic(description, "link"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPCVPCInterfacePolicyGroupExists("aci_pc/vpc_interface_policy_group.foopc/vpc_interface_policy_group", &pcvpc_interface_policy_group),
					testAccCheckAciPCVPCInterfacePolicyGroupAttributes(description, "link", &pcvpc_interface_policy_group),
				),
			},
			{
				ResourceName:      "aci_pcvpc_interface_policy_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciPCVPCInterfacePolicyGroup_update(t *testing.T) {
	var pcvpc_interface_policy_group models.PCVPCInterfacePolicyGroup
	description := "pc/vpc_interface_policy_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPCVPCInterfacePolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPCVPCInterfacePolicyGroupConfig_basic(description, "link"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPCVPCInterfacePolicyGroupExists("aci_pc/vpc_interface_policy_group.foopc/vpc_interface_policy_group", &pcvpc_interface_policy_group),
					testAccCheckAciPCVPCInterfacePolicyGroupAttributes(description, "link", &pcvpc_interface_policy_group),
				),
			},
			{
				Config: testAccCheckAciPCVPCInterfacePolicyGroupConfig_basic(description, "node"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPCVPCInterfacePolicyGroupExists("aci_pc/vpc_interface_policy_group.foopc/vpc_interface_policy_group", &pcvpc_interface_policy_group),
					testAccCheckAciPCVPCInterfacePolicyGroupAttributes(description, "node", &pcvpc_interface_policy_group),
				),
			},
		},
	})
}

func testAccCheckAciPCVPCInterfacePolicyGroupConfig_basic(description, lag_t string) string {
	return fmt.Sprintf(`

	resource "aci_pcvpc_interface_policy_group" "foopcvpc_interface_policy_group" {
		description = "%s"
		name        = "demo_if_pol_grp"
		annotation  = "tag_if_pol"
		lag_t       = "%s"
		name_alias  = "alias_if_pol"
	}  
	`, description, lag_t)
}

func testAccCheckAciPCVPCInterfacePolicyGroupExists(name string, pcvpc_interface_policy_group *models.PCVPCInterfacePolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("PC/VPC Interface Policy Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PC/VPC Interface Policy Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		pcvpc_interface_policy_groupFound := models.PCVPCInterfacePolicyGroupFromContainer(cont)
		if pcvpc_interface_policy_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("PC/VPC Interface Policy Group %s not found", rs.Primary.ID)
		}
		*pcvpc_interface_policy_group = *pcvpc_interface_policy_groupFound
		return nil
	}
}

func testAccCheckAciPCVPCInterfacePolicyGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_pc/vpc_interface_policy_group" {
			cont, err := client.Get(rs.Primary.ID)
			pcvpc_interface_policy_group := models.PCVPCInterfacePolicyGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("PC/VPC Interface Policy Group %s Still exists", pcvpc_interface_policy_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciPCVPCInterfacePolicyGroupAttributes(description, lag_t string, pcvpc_interface_policy_group *models.PCVPCInterfacePolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != pcvpc_interface_policy_group.Description {
			return fmt.Errorf("Bad pcvpc_interface_policy_group Description %s", pcvpc_interface_policy_group.Description)
		}

		if "demo_if_pol_grp" != pcvpc_interface_policy_group.Name {
			return fmt.Errorf("Bad pcvpc_interface_policy_group name %s", pcvpc_interface_policy_group.Name)
		}

		if "tag_if_pol" != pcvpc_interface_policy_group.Annotation {
			return fmt.Errorf("Bad pcvpc_interface_policy_group annotation %s", pcvpc_interface_policy_group.Annotation)
		}

		if lag_t != pcvpc_interface_policy_group.LagT {
			return fmt.Errorf("Bad pcvpc_interface_policy_group lag_t %s", pcvpc_interface_policy_group.LagT)
		}

		if "alias_if_pol" != pcvpc_interface_policy_group.NameAlias {
			return fmt.Errorf("Bad pcvpc_interface_policy_group name_alias %s", pcvpc_interface_policy_group.NameAlias)
		}

		return nil
	}
}
