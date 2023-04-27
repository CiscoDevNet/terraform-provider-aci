package aci

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Access Interface Configuration
func TestAccAciInterfaceConfiguration_Basic(t *testing.T) {
	var interfaceConfiguration models.InfraPortConfiguration
	node := acctest.RandIntRange(101, 4000)
	interfaceVal := "1/1/0" // If the card is absent, the test will fail.
	description := "Create Interface Configuration description"
	adminState := "up"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInterfaceConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInterfaceConfigurationConfig_basic(adminState, description, interfaceVal, node),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceConfigurationExists("aci_interface_config.foo_interface_configuration", &interfaceConfiguration),
					testAccCheckAciInterfaceConfigurationAttributes(adminState, description, interfaceVal, node, &interfaceConfiguration),
				),
			},
		},
	})
}

func TestAccAciInterfaceConfiguration_Update(t *testing.T) {
	var interfaceConfiguration models.InfraPortConfiguration
	node := acctest.RandIntRange(101, 4000)
	interfaceVal := "1/1/0" // If the card is absent, the test will fail.
	description := "Update Interface Configuration description"
	adminState := "down"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInterfaceConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInterfaceConfigurationConfig_basic(adminState, description, interfaceVal, node),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceConfigurationExists("aci_interface_config.foo_interface_configuration", &interfaceConfiguration),
					testAccCheckAciInterfaceConfigurationAttributes(adminState, description, interfaceVal, node, &interfaceConfiguration),
				),
			},
			{
				Config: testAccCheckAciInterfaceConfigurationConfig_basic(adminState, description, interfaceVal, node),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceConfigurationExists("aci_interface_config.foo_interface_configuration", &interfaceConfiguration),
					testAccCheckAciInterfaceConfigurationAttributes(adminState, description, interfaceVal, node, &interfaceConfiguration),
				),
			},
		},
	})
}

func testAccCheckAciInterfaceConfigurationConfig_basic(adminState, description, interfaceVal string, node int) string {
	return fmt.Sprintf(`
	resource "aci_interface_config" "foo_interface_configuration" {
		node        = %d
		interface   = "%s"
		port_type   = "access"
		description = "%s"
		admin_state = "%s"
	}`, node, interfaceVal, description, adminState)
}

func testAccCheckAciInterfaceConfigurationExists(name string, interfaceConfiguration *models.InfraPortConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Interface Configuration %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Interface Configuration dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		interfaceConfigurationFound := models.InfraPortConfigurationFromContainer(cont)
		if interfaceConfigurationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Interface Configuration %s not found", rs.Primary.ID)
		}
		*interfaceConfiguration = *interfaceConfigurationFound
		return nil
	}
}

func testAccCheckAciInterfaceConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_interface_config" {
			cont, err := client.Get(rs.Primary.ID)
			interfaceConfiguration := models.InfraPortConfigurationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Interface Configuration %s Still exists", interfaceConfiguration.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciInterfaceConfigurationAttributes(adminState, description, interfaceVal string, node int, interfaceConfiguration *models.InfraPortConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if description != interfaceConfiguration.Descr {
			return fmt.Errorf("Bad Interface Configuration - Description %s", interfaceConfiguration.Descr)
		}
		nodeId := strconv.Itoa(node)
		if nodeId != interfaceConfiguration.Node {
			return fmt.Errorf("Bad Interface Configuration - Node value %s", interfaceConfiguration.Node)
		}
		objectIntfVal := fmt.Sprintf("%s/%s/%s", interfaceConfiguration.Card, interfaceConfiguration.Port, interfaceConfiguration.SubPort)
		if interfaceVal != objectIntfVal {
			return fmt.Errorf("Bad Interface Configuration - Interface Parts %s", objectIntfVal)
		}
		if interfaceConfiguration.Shutdown != getAdminState(adminState) {
			return fmt.Errorf("Bad Interface Configuration - Admin State %s", interfaceConfiguration.Shutdown)
		}
		if interfaceConfiguration.BrkoutMap != "none" {
			return fmt.Errorf("Bad Interface Configuration - BreakoutMap %s", interfaceConfiguration.BrkoutMap)
		}

		return nil
	}
}

// Fabric Interface Configuration
func TestAccAciFabricInterfaceConfiguration_Basic(t *testing.T) {
	var fabricInterfaceConfiguration models.FabricPortConfiguration
	node := acctest.RandIntRange(101, 4000)
	interfaceVal := "1/1/0" // If the card is absent, the test will fail.
	description := "Create Interface Configuration description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFabricInterfaceConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFabricInterfaceConfigurationConfig_basic(description, interfaceVal, node),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricInterfaceConfigurationExists("aci_interface_config.foo_interface_configuration", &fabricInterfaceConfiguration),
					testAccCheckAciFabricInterfaceConfigurationAttributes(description, interfaceVal, node, &fabricInterfaceConfiguration),
				),
			},
		},
	})
}

func TestAccAciFabricInterfaceConfiguration_Update(t *testing.T) {
	var fabricInterfaceConfiguration models.FabricPortConfiguration
	node := acctest.RandIntRange(101, 4000)
	interfaceVal := "1/1/0" // If the card is absent, the test will fail.
	description := "Update Interface Configuration description"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFabricInterfaceConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFabricInterfaceConfigurationConfig_basic(description, interfaceVal, node),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricInterfaceConfigurationExists("aci_interface_config.foo_interface_configuration", &fabricInterfaceConfiguration),
					testAccCheckAciFabricInterfaceConfigurationAttributes(description, interfaceVal, node, &fabricInterfaceConfiguration),
				),
			},
			{
				Config: testAccCheckAciFabricInterfaceConfigurationConfig_basic(description, interfaceVal, node),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricInterfaceConfigurationExists("aci_interface_config.foo_interface_configuration", &fabricInterfaceConfiguration),
					testAccCheckAciFabricInterfaceConfigurationAttributes(description, interfaceVal, node, &fabricInterfaceConfiguration),
				),
			},
		},
	})
}

func testAccCheckAciFabricInterfaceConfigurationConfig_basic(description, interfaceVal string, node int) string {
	return fmt.Sprintf(`
	resource "aci_interface_config" "foo_interface_configuration" {
		node        = %d
		interface   = "%s"
		port_type   = "fabric"
		description = "%s"
	}`, node, interfaceVal, description)
}

func testAccCheckAciFabricInterfaceConfigurationExists(name string, fabricInterfaceConfiguration *models.FabricPortConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Interface Configuration %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Interface Configuration dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabricInterfaceConfigurationFound := models.FabricPortConfigurationFromContainer(cont)
		if fabricInterfaceConfigurationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Interface Configuration %s not found", rs.Primary.ID)
		}
		*fabricInterfaceConfiguration = *fabricInterfaceConfigurationFound
		return nil
	}
}

func testAccCheckAciFabricInterfaceConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_interface_config" {
			cont, err := client.Get(rs.Primary.ID)
			fabricInterfaceConfiguration := models.FabricPortConfigurationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Interface Configuration %s Still exists", fabricInterfaceConfiguration.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFabricInterfaceConfigurationAttributes(description, interfaceVal string, node int, fabricInterfaceConfiguration *models.FabricPortConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if description != fabricInterfaceConfiguration.Descr {
			return fmt.Errorf("Bad Interface Configuration Description %s", fabricInterfaceConfiguration.Descr)
		}
		nodeId := strconv.Itoa(node)
		if nodeId != fabricInterfaceConfiguration.Node {
			return fmt.Errorf("Bad Interface Configuration Node value %s", fabricInterfaceConfiguration.Node)
		}
		objectIntfVal := fmt.Sprintf("%s/%s/%s", fabricInterfaceConfiguration.Card, fabricInterfaceConfiguration.Port, fabricInterfaceConfiguration.SubPort)
		if interfaceVal != objectIntfVal {
			return fmt.Errorf("Bad Interface Configuration - Interface Parts %s", objectIntfVal)
		}
		return nil
	}
}

// Breakout Map Configuration
func TestAccAciBreakOutInterfaceConfiguration_Basic(t *testing.T) {
	var interfaceConfiguration models.InfraPortConfiguration
	node := acctest.RandIntRange(101, 4000)
	interfaceVal := "1/1/0" // If the card is absent, the test will fail.
	description := "Create Interface Configuration description"
	breakout := "100g-4x"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBreakOutInterfaceConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBreakOutInterfaceConfigurationConfig_basic(description, interfaceVal, breakout, node),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBreakOutInterfaceConfigurationExists("aci_interface_config.foo_interface_configuration", &interfaceConfiguration),
					testAccCheckAciBreakOutInterfaceConfigurationAttributes(description, interfaceVal, breakout, node, &interfaceConfiguration),
				),
			},
		},
	})
}

func TestAccAciBreakOutInterfaceConfiguration_Update(t *testing.T) {
	var interfaceConfiguration models.InfraPortConfiguration
	node := acctest.RandIntRange(101, 4000)
	interfaceVal := "1/1/0" // If the card is absent, the test will fail.
	description := "Update Interface Configuration description"
	breakout := "25g-4x"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBreakOutInterfaceConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBreakOutInterfaceConfigurationConfig_basic(description, interfaceVal, breakout, node),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBreakOutInterfaceConfigurationExists("aci_interface_config.foo_interface_configuration", &interfaceConfiguration),
					testAccCheckAciBreakOutInterfaceConfigurationAttributes(description, interfaceVal, breakout, node, &interfaceConfiguration),
				),
			},
			{
				Config: testAccCheckAciBreakOutInterfaceConfigurationConfig_basic(description, interfaceVal, breakout, node),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBreakOutInterfaceConfigurationExists("aci_interface_config.foo_interface_configuration", &interfaceConfiguration),
					testAccCheckAciBreakOutInterfaceConfigurationAttributes(description, interfaceVal, breakout, node, &interfaceConfiguration),
				),
			},
		},
	})
}

func testAccCheckAciBreakOutInterfaceConfigurationConfig_basic(description, interfaceVal, breakout string, node int) string {
	return fmt.Sprintf(`
	resource "aci_interface_config" "foo_interface_configuration" {
		node        = %d
		interface   = "%s"
		port_type   = "access"
		description = "%s"
		breakout    = "%s"
	}`, node, interfaceVal, description, breakout)
}

func testAccCheckAciBreakOutInterfaceConfigurationExists(name string, interfaceConfiguration *models.InfraPortConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Interface Configuration %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Interface Configuration dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		interfaceConfigurationFound := models.InfraPortConfigurationFromContainer(cont)
		if interfaceConfigurationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Interface Configuration %s not found", rs.Primary.ID)
		}
		*interfaceConfiguration = *interfaceConfigurationFound
		return nil
	}
}

func testAccCheckAciBreakOutInterfaceConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_interface_config" {
			cont, err := client.Get(rs.Primary.ID)
			interfaceConfiguration := models.InfraPortConfigurationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Interface Configuration %s Still exists", interfaceConfiguration.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBreakOutInterfaceConfigurationAttributes(description, interfaceVal, breakout string, node int, interfaceConfiguration *models.InfraPortConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if description != interfaceConfiguration.Descr {
			return fmt.Errorf("Bad Interface Configuration Description %s", interfaceConfiguration.Descr)
		}
		nodeId := strconv.Itoa(node)
		if nodeId != interfaceConfiguration.Node {
			return fmt.Errorf("Bad Interface Configuration Node value %s", interfaceConfiguration.Node)
		}
		objectIntfVal := fmt.Sprintf("%s/%s/%s", interfaceConfiguration.Card, interfaceConfiguration.Port, interfaceConfiguration.SubPort)
		if interfaceVal != objectIntfVal {
			return fmt.Errorf("Bad Interface Configuration - Interface Parts %s", objectIntfVal)
		}
		if breakout != interfaceConfiguration.BrkoutMap {
			return fmt.Errorf("Bad Interface Configuration - Breakout Map %s", interfaceConfiguration.BrkoutMap)
		}
		return nil
	}
}
