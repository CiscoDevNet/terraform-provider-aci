package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciFabricNodeControl_Basic(t *testing.T) {
	var fabric_node_control models.FabricNodeControl
	description := "fabric_node_control created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFabricNodeControlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFabricNodeControlConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeControlExists("aci_fabric_node_control.foofabric_node_control", &fabric_node_control),
					testAccCheckAciFabricNodeControlAttributes(description, &fabric_node_control),
				),
			},
		},
	})
}

func TestAccAciFabricNodeControl_Update(t *testing.T) {
	var fabric_node_control models.FabricNodeControl
	description := "fabric_node_control created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFabricNodeControlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFabricNodeControlConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeControlExists("aci_fabric_node_control.foofabric_node_control", &fabric_node_control),
					testAccCheckAciFabricNodeControlAttributes(description, &fabric_node_control),
				),
			},
			{
				Config: testAccCheckAciFabricNodeControlConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeControlExists("aci_fabric_node_control.foofabric_node_control", &fabric_node_control),
					testAccCheckAciFabricNodeControlAttributes(description, &fabric_node_control),
				),
			},
		},
	})
}

func testAccCheckAciFabricNodeControlConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_fabric_node_control" "foofabric_node_control" {
		name = "test"
		description = "%s"
		control = "Dom"
		feature_sel = "telemetry"
		name_alias = "test_alias"
		annotation = "test_annotation"
	}

	`, description)
}

func testAccCheckAciFabricNodeControlExists(name string, fabric_node_control *models.FabricNodeControl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric Node Control %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric Node Control dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabric_node_controlFound := models.FabricNodeControlFromContainer(cont)
		if fabric_node_controlFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric Node Control %s not found", rs.Primary.ID)
		}
		*fabric_node_control = *fabric_node_controlFound
		return nil
	}
}

func testAccCheckAciFabricNodeControlDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_fabric_node_control" {
			cont, err := client.Get(rs.Primary.ID)
			fabric_node_control := models.FabricNodeControlFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fabric Node Control %s Still exists", fabric_node_control.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFabricNodeControlAttributes(description string, fabric_node_control *models.FabricNodeControl) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != GetMOName(fabric_node_control.DistinguishedName) {
			return fmt.Errorf("Bad fabric_node_control %s", GetMOName(fabric_node_control.DistinguishedName))
		}

		if description != fabric_node_control.Description {
			return fmt.Errorf("Bad fabric_node_control Description %s", fabric_node_control.Description)
		}

		if "Dom" != fabric_node_control.Control {
			return fmt.Errorf("Bad fabric_node_control Control %s", fabric_node_control.Control)
		}

		if "telemetry" != fabric_node_control.FeatureSel {
			return fmt.Errorf("Bad fabric_node_control FeatureSel %s", fabric_node_control.FeatureSel)
		}

		if "test_alias" != fabric_node_control.NameAlias {
			return fmt.Errorf("Bad fabric_node_control NameAlias %s", fabric_node_control.NameAlias)
		}

		if "test_annotation" != fabric_node_control.Annotation {
			return fmt.Errorf("Bad fabric_node_control Annotation %s", fabric_node_control.Annotation)
		}
		return nil
	}
}
