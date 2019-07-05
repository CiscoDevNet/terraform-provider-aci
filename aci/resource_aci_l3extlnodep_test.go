package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciLogicalNodeProfile_Basic(t *testing.T) {
	var logical_node_profile models.LogicalNodeProfile
	description := "logical_node_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalNodeProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLogicalNodeProfileConfig_basic(description, "black"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalNodeProfileExists("aci_logical_node_profile.foological_node_profile", &logical_node_profile),
					testAccCheckAciLogicalNodeProfileAttributes(description, "black", &logical_node_profile),
				),
			},
			{
				ResourceName:      "aci_logical_node_profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciLogicalNodeProfile_update(t *testing.T) {
	var logical_node_profile models.LogicalNodeProfile
	description := "logical_node_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalNodeProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLogicalNodeProfileConfig_basic(description, "black"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalNodeProfileExists("aci_logical_node_profile.foological_node_profile", &logical_node_profile),
					testAccCheckAciLogicalNodeProfileAttributes(description, "black", &logical_node_profile),
				),
			},
			{
				Config: testAccCheckAciLogicalNodeProfileConfig_basic(description, "white"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalNodeProfileExists("aci_logical_node_profile.foological_node_profile", &logical_node_profile),
					testAccCheckAciLogicalNodeProfileAttributes(description, "white", &logical_node_profile),
				),
			},
		},
	})
}

func testAccCheckAciLogicalNodeProfileConfig_basic(description, tag string) string {
	return fmt.Sprintf(`

	resource "aci_logical_node_profile" "foological_node_profile" {
		l3_outside_dn = "${aci_l3_outside.example.id}"
		description   = "%s"
		name          = "demo_node"
		annotation    = "tag_node"
		config_issues = "none"
		name_alias    = "alias_node"
		tag           = "%s"
		target_dscp   = "unspecified"
	  }	  
	`, description, tag)
}

func testAccCheckAciLogicalNodeProfileExists(name string, logical_node_profile *models.LogicalNodeProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Logical Node Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Logical Node Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		logical_node_profileFound := models.LogicalNodeProfileFromContainer(cont)
		if logical_node_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Logical Node Profile %s not found", rs.Primary.ID)
		}
		*logical_node_profile = *logical_node_profileFound
		return nil
	}
}

func testAccCheckAciLogicalNodeProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_logical_node_profile" {
			cont, err := client.Get(rs.Primary.ID)
			logical_node_profile := models.LogicalNodeProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Logical Node Profile %s Still exists", logical_node_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLogicalNodeProfileAttributes(description, tag string, logical_node_profile *models.LogicalNodeProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != logical_node_profile.Description {
			return fmt.Errorf("Bad logical_node_profile Description %s", logical_node_profile.Description)
		}

		if "demo_node" != logical_node_profile.Name {
			return fmt.Errorf("Bad logical_node_profile name %s", logical_node_profile.Name)
		}

		if "tag_node" != logical_node_profile.Annotation {
			return fmt.Errorf("Bad logical_node_profile annotation %s", logical_node_profile.Annotation)
		}

		if "none" != logical_node_profile.ConfigIssues {
			return fmt.Errorf("Bad logical_node_profile config_issues %s", logical_node_profile.ConfigIssues)
		}

		if "alias_node" != logical_node_profile.NameAlias {
			return fmt.Errorf("Bad logical_node_profile name_alias %s", logical_node_profile.NameAlias)
		}

		if tag != logical_node_profile.Tag {
			return fmt.Errorf("Bad logical_node_profile tag %s", logical_node_profile.Tag)
		}

		if "unspecified" != logical_node_profile.TargetDscp {
			return fmt.Errorf("Bad logical_node_profile target_dscp %s", logical_node_profile.TargetDscp)
		}

		return nil
	}
}
