package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciManagedNodeConnectivityGroup_Basic(t *testing.T) {
	var managed_node_connectivity_group models.ManagedNodeConnectivityGroup
	annotation := "testing annoration"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciManagedNodeConnectivityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciManagedNodeConnectivityGroupConfig_basic(annotation),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciManagedNodeConnectivityGroupExists("aci_managed_node_connectivity_group.test", &managed_node_connectivity_group),
					testAccCheckAciManagedNodeConnectivityGroupAttributes(annotation, &managed_node_connectivity_group),
				),
			},
		},
	})
}

func TestAccAciManagedNodeConnectivityGroup_update(t *testing.T) {
	var managed_node_connectivity_group models.ManagedNodeConnectivityGroup
	annotation := "testing annoration"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciManagedNodeConnectivityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciManagedNodeConnectivityGroupConfig_basic(annotation),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciManagedNodeConnectivityGroupExists("aci_managed_node_connectivity_group.test", &managed_node_connectivity_group),
					testAccCheckAciManagedNodeConnectivityGroupAttributes(annotation, &managed_node_connectivity_group),
				),
			},
			{
				Config: testAccCheckAciManagedNodeConnectivityGroupConfig_basic(annotation),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciManagedNodeConnectivityGroupExists("aci_managed_node_connectivity_group.test", &managed_node_connectivity_group),
					testAccCheckAciManagedNodeConnectivityGroupAttributes(annotation, &managed_node_connectivity_group),
				),
			},
		},
	})
}

func testAccCheckAciManagedNodeConnectivityGroupConfig_basic(annotation string) string {
	return fmt.Sprintf(`

	resource "aci_managed_node_connectivity_group" "test" {
		name  = "example"
  		annotation  = "%s"
	}
	`, annotation)
}

func testAccCheckAciManagedNodeConnectivityGroupExists(name string, managed_node_connectivity_group *models.ManagedNodeConnectivityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Managed Node Connectivity Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Managed Node Connectivity Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		managed_node_connectivity_groupFound := models.ManagedNodeConnectivityGroupFromContainer(cont)
		if managed_node_connectivity_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Managed Node Connectivity Group %s not found", rs.Primary.ID)
		}
		*managed_node_connectivity_group = *managed_node_connectivity_groupFound
		return nil
	}
}

func testAccCheckAciManagedNodeConnectivityGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_managed_node_connectivity_group" {
			cont, err := client.Get(rs.Primary.ID)
			managed_node_connectivity_group := models.ManagedNodeConnectivityGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Managed Node Connectivity Group %s Still exists", managed_node_connectivity_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciManagedNodeConnectivityGroupAttributes(annotation string, managed_node_connectivity_group *models.ManagedNodeConnectivityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if "example" != managed_node_connectivity_group.Name {
			return fmt.Errorf("Bad managed_node_connectivity_group name %s", managed_node_connectivity_group.Name)
		}

		if annotation != managed_node_connectivity_group.Annotation {
			return fmt.Errorf("Bad managed_node_connectivity_group annotation %s", managed_node_connectivity_group.Annotation)
		}

		return nil
	}
}
