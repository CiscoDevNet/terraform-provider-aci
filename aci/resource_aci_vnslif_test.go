package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLogicalInterface_Basic(t *testing.T) {
	var logical_interface models.LogicalInterface
	fv_tenant_name := acctest.RandString(5)
	vns_l_dev_vip_name := acctest.RandString(5)
	vns_l_if_name := acctest.RandString(5)
	description := "logical_interface created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLogicalInterfaceConfig_basic(fv_tenant_name, vns_l_dev_vip_name, vns_l_if_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceExists("aci_logical_interface.foo_logical_interface", &logical_interface),
					testAccCheckAciLogicalInterfaceAttributes(fv_tenant_name, vns_l_dev_vip_name, vns_l_if_name, description, &logical_interface),
				),
			},
		},
	})
}

func TestAccAciLogicalInterface_Update(t *testing.T) {
	var logical_interface models.LogicalInterface
	fv_tenant_name := acctest.RandString(5)
	vns_l_dev_vip_name := acctest.RandString(5)
	vns_l_if_name := acctest.RandString(5)
	description := "logical_interface created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLogicalInterfaceConfig_basic(fv_tenant_name, vns_l_dev_vip_name, vns_l_if_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceExists("aci_logical_interface.foo_logical_interface", &logical_interface),
					testAccCheckAciLogicalInterfaceAttributes(fv_tenant_name, vns_l_dev_vip_name, vns_l_if_name, description, &logical_interface),
				),
			},
			{
				Config: testAccCheckAciLogicalInterfaceConfig_basic(fv_tenant_name, vns_l_dev_vip_name, vns_l_if_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceExists("aci_logical_interface.foo_logical_interface", &logical_interface),
					testAccCheckAciLogicalInterfaceAttributes(fv_tenant_name, vns_l_dev_vip_name, vns_l_if_name, description, &logical_interface),
				),
			},
		},
	})
}

func testAccCheckAciLogicalInterfaceConfig_basic(fv_tenant_name, vns_l_dev_vip_name, vns_l_if_name string) string {
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

	resource "aci_logical_interface" "foo_logical_interface" {
		name 		= "%s"
		description = "logical_interface created while acceptance testing"
		l4-l7_devices_dn = aci_l4-l7_devices.foo_l4-l7_devices.id
	}

	`, fv_tenant_name, vns_l_dev_vip_name, vns_l_if_name)
}

func testAccCheckAciLogicalInterfaceExists(name string, logical_interface *models.LogicalInterface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Logical Interface %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Logical Interface dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		logical_interfaceFound := models.LogicalInterfaceFromContainer(cont)
		if logical_interfaceFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Logical Interface %s not found", rs.Primary.ID)
		}
		*logical_interface = *logical_interfaceFound
		return nil
	}
}

func testAccCheckAciLogicalInterfaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_logical_interface" {
			cont, err := client.Get(rs.Primary.ID)
			logical_interface := models.LogicalInterfaceFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Logical Interface %s Still exists", logical_interface.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLogicalInterfaceAttributes(fv_tenant_name, vns_l_dev_vip_name, vns_l_if_name, description string, logical_interface *models.LogicalInterface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vns_l_if_name != GetMOName(logical_interface.DistinguishedName) {
			return fmt.Errorf("Bad vnsl_if %s", GetMOName(logical_interface.DistinguishedName))
		}

		if vns_l_dev_vip_name != GetMOName(GetParentDn(logical_interface.DistinguishedName, logical_interface.Rn)) {
			return fmt.Errorf(" Bad vnsl_dev_vip %s", GetMOName(GetParentDn(logical_interface.DistinguishedName, logical_interface.Rn)))
		}
		if description != logical_interface.Description {
			return fmt.Errorf("Bad logical_interface Description %s", logical_interface.Description)
		}
		return nil
	}
}
