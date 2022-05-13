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

func TestAccAciConcreteDevice_Basic(t *testing.T) {
	var concrete_device models.ConcreteDevice
	fv_tenant_name := acctest.RandString(5)
	vns_l_dev_vip_name := acctest.RandString(5)
	vns_c_dev_name := acctest.RandString(5)
	description := "concrete_device created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciConcreteDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciConcreteDeviceConfig_basic(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConcreteDeviceExists("aci_concrete_device.foo_concrete_device", &concrete_device),
					testAccCheckAciConcreteDeviceAttributes(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, description, &concrete_device),
				),
			},
		},
	})
}

func TestAccAciConcreteDevice_Update(t *testing.T) {
	var concrete_device models.ConcreteDevice
	fv_tenant_name := acctest.RandString(5)
	vns_l_dev_vip_name := acctest.RandString(5)
	vns_c_dev_name := acctest.RandString(5)
	description := "concrete_device created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciConcreteDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciConcreteDeviceConfig_basic(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConcreteDeviceExists("aci_concrete_device.foo_concrete_device", &concrete_device),
					testAccCheckAciConcreteDeviceAttributes(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, description, &concrete_device),
				),
			},
			{
				Config: testAccCheckAciConcreteDeviceConfig_basic(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConcreteDeviceExists("aci_concrete_device.foo_concrete_device", &concrete_device),
					testAccCheckAciConcreteDeviceAttributes(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, description, &concrete_device),
				),
			},
		},
	})
}

func testAccCheckAciConcreteDeviceConfig_basic(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_l4_l7_devices" "foo_l4_l7_devices" {
		name 		= "%s"
		description = "l4_l7_devices created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	resource "aci_concrete_device" "foo_concrete_device" {
		name 		= "%s"
		description = "concrete_device created while acceptance testing"
		l4_l7_devices_dn = aci_l4_l7_devices.foo_l4_l7_devices.id
	}

	`, fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name)
}

func testAccCheckAciConcreteDeviceExists(name string, concrete_device *models.ConcreteDevice) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Concrete Device %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Concrete Device dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		concrete_deviceFound := models.ConcreteDeviceFromContainer(cont)
		if concrete_deviceFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Concrete Device %s not found", rs.Primary.ID)
		}
		*concrete_device = *concrete_deviceFound
		return nil
	}
}

func testAccCheckAciConcreteDeviceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_concrete_device" {
			cont, err := client.Get(rs.Primary.ID)
			concrete_device := models.ConcreteDeviceFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Concrete Device %s Still exists", concrete_device.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciConcreteDeviceAttributes(fv_tenant_name, vns_l_dev_vip_name, vns_c_dev_name, description string, concrete_device *models.ConcreteDevice) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vns_c_dev_name != GetMOName(concrete_device.DistinguishedName) {
			return fmt.Errorf("Bad vnsc_dev %s", GetMOName(concrete_device.DistinguishedName))
		}

		if vns_l_dev_vip_name != GetMOName(GetParentDn(concrete_device.DistinguishedName, concrete_device.Rn)) {
			return fmt.Errorf(" Bad vnsl_dev_vip %s", GetMOName(GetParentDn(concrete_device.DistinguishedName, concrete_device.Rn)))
		}
		if description != concrete_device.Description {
			return fmt.Errorf("Bad concrete_device Description %s", concrete_device.Description)
		}
		return nil
	}
}
