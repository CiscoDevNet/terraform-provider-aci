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

func TestAccAciL4ToL7Devices_Basic(t *testing.T) {
	var l4Tol7_devices models.L4ToL7Devices
	fv_tenant_name := acctest.RandString(5)
	vns_l_dev_vip_name := acctest.RandString(5)
	description := "l4Tol7_device created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL4ToL7DevicesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL4ToL7DevicesConfig_basic(fv_tenant_name, vns_l_dev_vip_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4ToL7DevicesExists("aci_l4Tol7_devices.fool4Tol7_devices", &l4Tol7_devices),
					testAccCheckAciL4ToL7DevicesAttributes(fv_tenant_name, vns_l_dev_vip_name, description, &l4Tol7_devices),
				),
			},
		},
	})
}

func TestAccAciL4ToL7Devices_Update(t *testing.T) {
	var l4Tol7_devices models.L4ToL7Devices
	fv_tenant_name := acctest.RandString(5)
	vns_l_dev_vip_name := acctest.RandString(5)
	description := "l4Tol7_devices created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL4ToL7DevicesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL4ToL7DevicesConfig_basic(fv_tenant_name, vns_l_dev_vip_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4ToL7DevicesExists("aci_l4Tol7_devices.fool4Tol7_devices", &l4Tol7_devices),
					testAccCheckAciL4ToL7DevicesAttributes(fv_tenant_name, vns_l_dev_vip_name, description, &l4Tol7_devices),
				),
			},
			{
				Config: testAccCheckAciL4ToL7DevicesConfig_basic(fv_tenant_name, vns_l_dev_vip_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4ToL7DevicesExists("aci_l4Tol7_devices.fool4Tol7_devices", &l4Tol7_devices),
					testAccCheckAciL4ToL7DevicesAttributes(fv_tenant_name, vns_l_dev_vip_name, description, &l4Tol7_devices),
				),
			},
		},
	})
}

func testAccCheckAciL4ToL7DevicesConfig_basic(fv_tenant_name, vns_l_dev_vip_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_l4Tol7_device" "fool4Tol7_device" {
		name 		= "%s"
		description = "l4Tol7_device created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	`, fv_tenant_name, vns_l_dev_vip_name)
}

func testAccCheckAciL4ToL7DevicesExists(name string, l4Tol7_devices *models.L4ToL7Devices) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L4ToL7 Device %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L4ToL7 Device dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l4Tol7_devicesFound := models.L4ToL7DevicesFromContainer(cont)
		if l4Tol7_devicesFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L4ToL7 Device %s not found", rs.Primary.ID)
		}
		*l4Tol7_devices = *l4Tol7_devicesFound
		return nil
	}
}

func testAccCheckAciL4ToL7DevicesDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l4Tol7_devices" {
			cont, err := client.Get(rs.Primary.ID)
			l4Tol7_devices := models.L4ToL7DevicesFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L4ToL7 Device %s Still exists", l4Tol7_devices.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL4ToL7DevicesAttributes(fv_tenant_name, vns_l_dev_vip_name, description string, l4Tol7_devices *models.L4ToL7Devices) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vns_l_dev_vip_name != GetMOName(l4Tol7_devices.DistinguishedName) {
			return fmt.Errorf("Bad vnsl_dev_vip %s", GetMOName(l4Tol7_devices.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(l4Tol7_devices.DistinguishedName, l4Tol7_devices.Rn)) {
			return fmt.Errorf(" Bad fv_tenant %s", GetMOName(GetParentDn(l4Tol7_devices.DistinguishedName, l4Tol7_devices.Rn)))
		}
		if description != l4Tol7_devices.Description {
			return fmt.Errorf("Bad l4Tol7_device Description %s", l4Tol7_devices.Description)
		}
		return nil
	}
}
