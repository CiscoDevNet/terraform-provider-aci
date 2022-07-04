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

func TestAccAciLeakInternalSubnet_Basic(t *testing.T) {
	var leak_internal_subnet models.LeakInternalSubnet
	fv_tenant_name := acctest.RandString(5)
	fv_ctx_name := acctest.RandString(5)
	leak_internal_subnet_name := acctest.RandString(5)
	description := "leak_internal_subnet created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLeakInternalSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLeakInternalSubnetConfig_basic(fv_tenant_name, fv_ctx_name, leak_internal_subnet_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeakInternalSubnetExists("aci_leak_internal_subnet.foo_leak_internal_subnet", &leak_internal_subnet),
					testAccCheckAciLeakInternalSubnetAttributes(fv_tenant_name, fv_ctx_name, leak_internal_subnet_name, description, &leak_internal_subnet),
				),
			},
		},
	})
}

func TestAccAciLeakInternalSubnet_Update(t *testing.T) {
	var leak_internal_subnet models.LeakInternalSubnet
	fv_tenant_name := acctest.RandString(5)
	fv_ctx_name := acctest.RandString(5)
	leak_internal_subnet_name := acctest.RandString(5)
	description := "leak_internal_subnet created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLeakInternalSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLeakInternalSubnetConfig_basic(fv_tenant_name, fv_ctx_name, leak_internal_subnet_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeakInternalSubnetExists("aci_leak_internal_subnet.foo_leak_internal_subnet", &leak_internal_subnet),
					testAccCheckAciLeakInternalSubnetAttributes(fv_tenant_name, fv_ctx_name, leak_internal_subnet_name, description, &leak_internal_subnet),
				),
			},
			{
				Config: testAccCheckAciLeakInternalSubnetConfig_basic(fv_tenant_name, fv_ctx_name, leak_internal_subnet_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeakInternalSubnetExists("aci_leak_internal_subnet.foo_leak_internal_subnet", &leak_internal_subnet),
					testAccCheckAciLeakInternalSubnetAttributes(fv_tenant_name, fv_ctx_name, leak_internal_subnet_name, description, &leak_internal_subnet),
				),
			},
		},
	})
}

func testAccCheckAciLeakInternalSubnetConfig_basic(fv_tenant_name, fv_ctx_name, leak_internal_subnet_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_vrf" "foo_vrf" {
		name 		= "%s"
		description = "vrf created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	resource "aci_leak_internal_subnet" "foo_leak_internal_subnet" {
		name 		= "%s"
		description = "leak_internal_subnet created while acceptance testing"
		vrf_dn = aci_vrf.foo_vrf.id
	}

	`, fv_tenant_name, fv_ctx_name, leak_internal_subnet_name)
}

func testAccCheckAciLeakInternalSubnetExists(name string, leak_internal_subnet *models.LeakInternalSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Inter-VRF Leaked EPG/BD Subnet %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Inter-VRF Leaked EPG/BD Subnet dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		leak_internal_subnetFound := models.LeakInternalSubnetFromContainer(cont)
		if leak_internal_subnetFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Inter-VRF Leaked EPG/BD Subnet %s not found", rs.Primary.ID)
		}
		*leak_internal_subnet = *leak_internal_subnetFound
		return nil
	}
}

func testAccCheckAciLeakInternalSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_leak_internal_subnet" {
			cont, err := client.Get(rs.Primary.ID)
			leak_internal_subnet := models.LeakInternalSubnetFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Inter-VRF Leaked EPG/BD Subnet %s Still exists", leak_internal_subnet.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLeakInternalSubnetAttributes(fv_tenant_name, fv_ctx_name, leak_internal_subnet_name, description string, leak_internal_subnet *models.LeakInternalSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if leak_internal_subnet_name != GetMOName(leak_internal_subnet.DistinguishedName) {
			return fmt.Errorf("Bad leak_internal_subnet %s", GetMOName(leak_internal_subnet.DistinguishedName))
		}

		// if fv_ctx_name != GetMOName(GetParentDn(leak_internal_subnet.DistinguishedName)) {
		// 	return fmt.Errorf(" Bad fv_ctx %s", GetMOName(GetParentDn(leak_internal_subnet.DistinguishedName)))
		// }
		if description != leak_internal_subnet.Description {
			return fmt.Errorf("Bad leak_internal_subnet Description %s", leak_internal_subnet.Description)
		}
		return nil
	}
}
