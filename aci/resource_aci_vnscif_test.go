package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciConcreteInterface_Basic(t *testing.T) {
	var concrete_interface models.ConcreteInterface
	fv_tenant_name := acctest.RandString(5)
	vns_l_dev_vip_name := acctest.RandString(5)
	vns_c_dev_name := acctest.RandString(5)
	vns_c_if_name := acctest.RandString(5)
	description := "concrete_interface created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciConcreteInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciConcreteInterfaceConfig_basic(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, vns_c_if_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConcreteInterfaceExists("aci_concrete_interface.foo_concrete_interface", &concrete_interface),
					testAccCheckAciConcreteInterfaceAttributes(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, vns_c_if_name, description, &concrete_interface),
				),
			},
		},
	})
}

func TestAccAciConcreteInterface_Update(t *testing.T) {
	var concrete_interface models.ConcreteInterface
	fv_tenant_name := acctest.RandString(5)
	vns_l_dev_vip_name := acctest.RandString(5)
	vns_c_dev_name := acctest.RandString(5)
	vns_c_if_name := acctest.RandString(5)
	description := "concrete_interface created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciConcreteInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciConcreteInterfaceConfig_basic(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, vns_c_if_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConcreteInterfaceExists("aci_concrete_interface.foo_concrete_interface", &concrete_interface),
					testAccCheckAciConcreteInterfaceAttributes(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, vns_c_if_name, description, &concrete_interface),
				),
			},
			{
				Config: testAccCheckAciConcreteInterfaceConfig_basic(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, vns_c_if_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConcreteInterfaceExists("aci_concrete_interface.foo_concrete_interface", &concrete_interface),
					testAccCheckAciConcreteInterfaceAttributes(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, vns_c_if_name, description, &concrete_interface),
				),
			},
		},
	})
}

func testAccCheckAciConcreteInterfaceConfig_basic(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, vns_c_if_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_l4-l7_devices" "foo_l4-l7_devices" {
		name 		= "%s"
		description = "l4-l7_devices created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	resource "aci_concrete_device" "foo_concrete_device" {
		name 		= "%s"
		description = "concrete_device created while acceptance testing"
		l4-l7_devices_dn = aci_l4-l7_devices.foo_l4-l7_devices.id
	}

	resource "aci_concrete_interface" "foo_concrete_interface" {
		name 		= "%s"
		description = "concrete_interface created while acceptance testing"
		concrete_device_dn = aci_concrete_device.foo_concrete_device.id
	}

	`, fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, vns_c_if_name)
}

func testAccCheckAciConcreteInterfaceExists(name string, concrete_interface *models.ConcreteInterface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Concrete Interface %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Concrete Interface dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		concrete_interfaceFound := models.ConcreteInterfaceFromContainer(cont)
		if concrete_interfaceFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Concrete Interface %s not found", rs.Primary.ID)
		}
		*concrete_interface = *concrete_interfaceFound
		return nil
	}
}

func testAccCheckAciConcreteInterfaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_concrete_interface" {
			cont, err := client.Get(rs.Primary.ID)
			concrete_interface := models.ConcreteInterfaceFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Concrete Interface %s Still exists", concrete_interface.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciConcreteInterfaceAttributes(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, vns_c_if_name, description string, concrete_interface *models.ConcreteInterface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vns_c_if_name != GetMOName(concrete_interface.DistinguishedName) {
			return fmt.Errorf("Bad vnsc_if %s", GetMOName(concrete_interface.DistinguishedName))
		}

		if vns_c_dev_name != GetMOName(GetParentDn(concrete_interface.DistinguishedName, concrete_interface.Rn)) {
			return fmt.Errorf(" Bad vnsc_dev %s", GetMOName(GetParentDn(concrete_interface.DistinguishedName, concrete_interface.Rn)))
		}
		if description != concrete_interface.Description {
			return fmt.Errorf("Bad concrete_interface Description %s", concrete_interface.Description)
		}
		return nil
	}
}
