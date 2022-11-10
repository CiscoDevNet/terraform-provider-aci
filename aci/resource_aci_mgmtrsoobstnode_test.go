package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciOutofbandStaticNode_Basic(t *testing.T) {
	var mgmtStNode models.OutofbandStaticNode
	description := "manaegement static node created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOutofbandStaticNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOutofbandStaticNodeConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutofbandStaticNodeExists("aci_static_node_mgmt_address.test", &mgmtStNode),
					testAccCheckAciOutofbandStaticNodeAttributes(description, &mgmtStNode),
				),
			},
		},
	})
}

func TestAccAciOutofbandStaticNode_update(t *testing.T) {
	var mgmtStNode models.OutofbandStaticNode
	description := "manaegement static node created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOutofbandStaticNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOutofbandStaticNodeConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutofbandStaticNodeExists("aci_static_node_mgmt_address.test", &mgmtStNode),
					testAccCheckAciOutofbandStaticNodeAttributes(description, &mgmtStNode),
				),
			},
			{
				Config: testAccCheckAciOutofbandStaticNodeConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutofbandStaticNodeExists("aci_static_node_mgmt_address.test", &mgmtStNode),
					testAccCheckAciOutofbandStaticNodeAttributes(description, &mgmtStNode),
				),
			},
		},
	})
}

func testAccCheckAciOutofbandStaticNodeConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_static_node_mgmt_address" "test" {
		management_epg_dn = aci_node_mgmt_epg.foo_aci_node_mgmt_epg.id
		description       = "%s"
		t_dn              = "topology/pod-1/node-1"
		type              = "out_of_band"
		addr              = "10.20.30.40/20"
		annotation        = "example"
		gw                = "10.20.30.41"
		v6_addr           = "1::40/64"
		v6_gw             = "1::21"
	}
	`, description)
}

func testAccCheckAciOutofbandStaticNodeExists(name string, mgmtStNode *models.OutofbandStaticNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Out-of-band Static Node %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Out-of-band Static Node dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		mgmtStNodeFound := models.OutofbandStaticNodeFromContainer(cont)
		if mgmtStNodeFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Out-of-band Static Node %s not found", rs.Primary.ID)
		}
		*mgmtStNode = *mgmtStNodeFound
		return nil
	}
}

func testAccCheckAciOutofbandStaticNodeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_static_node_mgmt_address" {
			cont, err := client.Get(rs.Primary.ID)
			mgmtStNode := models.OutofbandStaticNodeFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Out-of-band Static Node %s Still exists", mgmtStNode.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciOutofbandStaticNodeAttributes(description string, mgmtStNode *models.OutofbandStaticNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != mgmtStNode.Description {
			return fmt.Errorf("Bad mgmtStNode Description %s", mgmtStNode.Description)
		}

		if "topology/pod-1/node-1" != mgmtStNode.TDn {
			return fmt.Errorf("Bad mgmtStNode t_dn %s", mgmtStNode.TDn)
		}

		if "10.20.30.40/20" != mgmtStNode.Addr {
			return fmt.Errorf("Bad mgmtStNode addr %s", mgmtStNode.Addr)
		}

		if "example" != mgmtStNode.Annotation {
			return fmt.Errorf("Bad mgmtStNode annotation %s", mgmtStNode.Annotation)
		}

		if "10.20.30.41" != mgmtStNode.Gw {
			return fmt.Errorf("Bad mgmtStNode gw %s", mgmtStNode.Gw)
		}

		if "1::40/64" != mgmtStNode.V6Addr {
			return fmt.Errorf("Bad mgmtStNode v6_addr %s", mgmtStNode.V6Addr)
		}

		if "1::21" != mgmtStNode.V6Gw {
			return fmt.Errorf("Bad mgmtStNode v6_gw %s", mgmtStNode.V6Gw)
		}

		return nil
	}
}
