package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciFabricNode_Basic(t *testing.T) {
	var fabric_node models.FabricNode
	description := "fabric_node"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFabricNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFabricNodeConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists("aci_logical_node_to_fabric_node.foofabric_node", &fabric_node),
					testAccCheckAciFabricNodeAttributes(description, &fabric_node),
				),
			},
		},
	})
}

func TestAccAciFabricNode_update(t *testing.T) {
	var fabric_node models.FabricNode
	description := "fabric_node"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFabricNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFabricNodeConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists("aci_logical_node_to_fabric_node.foofabric_node", &fabric_node),
					testAccCheckAciFabricNodeAttributes(description, &fabric_node),
				),
			},
			{
				Config: testAccCheckAciFabricNodeConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists("aci_logical_node_to_fabric_node.foofabric_node", &fabric_node),
					testAccCheckAciFabricNodeAttributes(description, &fabric_node),
				),
			},
		},
	})
}

func testAccCheckAciFabricNodeConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_logical_node_to_fabric_node" "foofabric_node" {
		logical_node_profile_dn = "uni/tn-demo_dev_tenant/out-demo_l3out/lnodep-demo_node"
		#logical_node_profile_dn = aci_logical_node_profile.example.id
		tdn              = "topology/pod-1/node-201"
		annotation       = "%s"
		config_issues    = "none"
		rtr_id           = "10.0.1.1"
		rtr_id_loop_back = "no"
	}
	`, description)
}

func testAccCheckAciFabricNodeExists(name string, fabric_node *models.FabricNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric Node %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric Node dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabric_nodeFound := models.FabricNodeFromContainer(cont)
		if fabric_nodeFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric Node %s not found", rs.Primary.ID)
		}
		*fabric_node = *fabric_nodeFound
		return nil
	}
}

func testAccCheckAciFabricNodeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_logical_node_to_fabric_node" {
			cont, err := client.Get(rs.Primary.ID)
			fabric_node := models.FabricNodeFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fabric Node %s Still exists", fabric_node.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciFabricNodeAttributes(description string, fabric_node *models.FabricNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if "topology/pod-1/node-201" != fabric_node.TDn {
			return fmt.Errorf("Bad fabric_node t_dn %s", fabric_node.TDn)
		}

		if description != fabric_node.Annotation {
			return fmt.Errorf("Bad fabric_node annotation %s", fabric_node.Annotation)
		}

		if "10.0.1.1" != fabric_node.RtrId {
			return fmt.Errorf("Bad fabric_node rtr_id %s", fabric_node.RtrId)
		}

		if "no" != fabric_node.RtrIdLoopBack {
			return fmt.Errorf("Bad fabric_node rtr_id_loop_back %s", fabric_node.RtrIdLoopBack)
		}

		return nil
	}
}
