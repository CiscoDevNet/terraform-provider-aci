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

func TestAccAciOutofServiceFabricPath_Basic(t *testing.T) {
	var outof_service_fabric_path models.OutofServiceFabricPath
	fabric_rs_oos_path_name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOutofServiceFabricPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOutofServiceFabricPathConfig_basic(fabric_rs_oos_path_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutofServiceFabricPathExists("aci_interface_blacklist.foointerface_blacklist", &outof_service_fabric_path),
				),
			},
		},
	})
}

func TestAccAciOutofServiceFabricPath_Update(t *testing.T) {
	var outof_service_fabric_path models.OutofServiceFabricPath
	fabric_rs_oos_path_name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOutofServiceFabricPathDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOutofServiceFabricPathConfig_basic(fabric_rs_oos_path_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutofServiceFabricPathExists("interface_blacklist.foointerface_blacklist", &outof_service_fabric_path),
					testAccCheckAciOutofServiceFabricPathAttributes(fabric_rs_oos_path_name, &outof_service_fabric_path),
				),
			},
			{
				Config: testAccCheckAciOutofServiceFabricPathConfig_basic(fabric_rs_oos_path_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOutofServiceFabricPathExists("interface_blacklist.foointerface_blacklist", &outof_service_fabric_path),
					testAccCheckAciOutofServiceFabricPathAttributes(fabric_rs_oos_path_name, &outof_service_fabric_path),
				),
			},
		},
	})
}

func testAccCheckAciOutofServiceFabricPathConfig_basic(fabric_rs_oos_path_name string) string {
	return fmt.Sprintf(`

	resource "interface_blacklist" "foointerface_blacklist" {
		pod_id    = 1
		node_id   = 101
		interface = "%s"
	}

	`, fabric_rs_oos_path_name)
}

func testAccCheckAciOutofServiceFabricPathExists(name string, outof_service_fabric_path *models.OutofServiceFabricPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Out of Service Fabric Path %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Out of Service Fabric Path dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		outof_service_fabric_pathFound := models.OutofServiceFabricPathFromContainer(cont)
		if outof_service_fabric_pathFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Out of Service Fabric Path %s not found", rs.Primary.ID)
		}
		*outof_service_fabric_path = *outof_service_fabric_pathFound
		return nil
	}
}

func testAccCheckAciOutofServiceFabricPathDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_interface_blacklist" {
			cont, err := client.Get(rs.Primary.ID)
			outof_service_fabric_path := models.OutofServiceFabricPathFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Out of Service Fabric Path %s Still exists", outof_service_fabric_path.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciOutofServiceFabricPathAttributes(planner_tenant_tmpl_name string, outof_service_fabric_path *models.OutofServiceFabricPath) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if planner_tenant_tmpl_name != GetMOName(outof_service_fabric_path.DistinguishedName) {
			return fmt.Errorf("Bad planner_tenant_tmpl_name %s", GetMOName(outof_service_fabric_path.DistinguishedName))
		}

		return nil
	}
}
