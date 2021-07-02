package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciFabricNodeMember_Basic(t *testing.T) {
	var fabric_node_member models.FabricNodeMember
	description := "fabric_node_member created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFabricNodeMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFabricNodeMemberConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists("aci_fabric_node_member.foofabric_node_member", &fabric_node_member),
					testAccCheckAciFabricNodeMemberAttributes(description, &fabric_node_member),
				),
			},
		},
	})
}

func TestAccAciFabricNodeMember_update(t *testing.T) {
	var fabric_node_member models.FabricNodeMember
	description := "fabric_node_member created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFabricNodeMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFabricNodeMemberConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists("aci_fabric_node_member.foofabric_node_member", &fabric_node_member),
					testAccCheckAciFabricNodeMemberAttributes(description, &fabric_node_member),
				),
			},
			{
				Config: testAccCheckAciFabricNodeMemberConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeMemberExists("aci_fabric_node_member.foofabric_node_member", &fabric_node_member),
					testAccCheckAciFabricNodeMemberAttributes(description, &fabric_node_member),
				),
			},
		},
	})
}

func testAccCheckAciFabricNodeMemberConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_fabric_node_member" "foofabric_node_member" {
		description = "%s"
		name = "test"
		serial  = "example"
		  annotation  = "example"
		  ext_pool_id  = "example"
		  fabric_id  = "example"
		  name_alias  = "example"
		  node_id  = "example"
		  node_type  = "unspecified"
		  pod_id  = "example"
		  role  = "unspecified"
		}
	`, description)
}

func testAccCheckAciFabricNodeMemberExists(name string, fabric_node_member *models.FabricNodeMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric Node Member %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric Node Member dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabric_node_memberFound := models.FabricNodeMemberFromContainer(cont)
		if fabric_node_memberFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric Node Member %s not found", rs.Primary.ID)
		}
		*fabric_node_member = *fabric_node_memberFound
		return nil
	}
}

func testAccCheckAciFabricNodeMemberDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_fabric_node_member" {
			cont, err := client.Get(rs.Primary.ID)
			fabric_node_member := models.FabricNodeMemberFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fabric Node Member %s Still exists", fabric_node_member.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciFabricNodeMemberAttributes(description string, fabric_node_member *models.FabricNodeMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != fabric_node_member.Description {
			return fmt.Errorf("Bad fabric_node_member Description %s", fabric_node_member.Description)
		}

		if "example" != fabric_node_member.Serial {
			return fmt.Errorf("Bad fabric_node_member serial %s", fabric_node_member.Serial)
		}

		if "example" != fabric_node_member.Annotation {
			return fmt.Errorf("Bad fabric_node_member annotation %s", fabric_node_member.Annotation)
		}

		if "example" != fabric_node_member.ExtPoolId {
			return fmt.Errorf("Bad fabric_node_member ext_pool_id %s", fabric_node_member.ExtPoolId)
		}

		if "example" != fabric_node_member.FabricId {
			return fmt.Errorf("Bad fabric_node_member fabric_id %s", fabric_node_member.FabricId)
		}

		if "example" != fabric_node_member.NameAlias {
			return fmt.Errorf("Bad fabric_node_member name_alias %s", fabric_node_member.NameAlias)
		}

		if "example" != fabric_node_member.NodeId {
			return fmt.Errorf("Bad fabric_node_member node_id %s", fabric_node_member.NodeId)
		}

		if "unspecified" != fabric_node_member.NodeType {
			return fmt.Errorf("Bad fabric_node_member node_type %s", fabric_node_member.NodeType)
		}

		if "example" != fabric_node_member.PodId {
			return fmt.Errorf("Bad fabric_node_member pod_id %s", fabric_node_member.PodId)
		}

		if "unspecified" != fabric_node_member.Role {
			return fmt.Errorf("Bad fabric_node_member role %s", fabric_node_member.Role)
		}

		return nil
	}
}
