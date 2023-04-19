package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLeakInternalSubnet_Basic(t *testing.T) {
	var leak_internal_subnet models.LeakInternalSubnet
	var leak_to_1 models.TenantandVRFdestinationforInterVRFLeakedRoutes
	var leak_to_2 models.TenantandVRFdestinationforInterVRFLeakedRoutes

	fv_tenant_name := acctest.RandString(5)
	fv_dest_vrf_name := acctest.RandString(5)
	fv_src_vrf_name := acctest.RandString(5)
	leak_internal_subnet_ip := "1.1.1.10/24"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLeakInternalSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLeakInternalSubnetConfig_basic(fv_tenant_name, fv_src_vrf_name, fv_dest_vrf_name, leak_internal_subnet_ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeakInternalSubnetExists("aci_vrf_leak_epg_bd_subnet.foo_leak_internal_subnet", &leak_internal_subnet, &leak_to_1, &leak_to_2),
					testAccCheckAciLeakInternalSubnetAttributes(fv_tenant_name, fv_src_vrf_name, fv_dest_vrf_name, leak_internal_subnet_ip, &leak_internal_subnet, &leak_to_1, &leak_to_2),
				),
			},
		},
	})
}

func TestAccAciLeakInternalSubnet_Update(t *testing.T) {
	var leak_internal_subnet models.LeakInternalSubnet
	var leak_to_1 models.TenantandVRFdestinationforInterVRFLeakedRoutes
	var leak_to_2 models.TenantandVRFdestinationforInterVRFLeakedRoutes

	fv_tenant_name := acctest.RandString(5)
	fv_dest_vrf_name := acctest.RandString(5)
	fv_src_vrf_name := acctest.RandString(5)
	leak_internal_subnet_ip := "1.1.1.20/24"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLeakInternalSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLeakInternalSubnetConfig_basic(fv_tenant_name, fv_src_vrf_name, fv_dest_vrf_name, leak_internal_subnet_ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeakInternalSubnetExists("aci_vrf_leak_epg_bd_subnet.foo_leak_internal_subnet", &leak_internal_subnet, &leak_to_1, &leak_to_2),
					testAccCheckAciLeakInternalSubnetAttributes(fv_tenant_name, fv_src_vrf_name, fv_dest_vrf_name, leak_internal_subnet_ip, &leak_internal_subnet, &leak_to_1, &leak_to_2),
				),
			},
			{
				Config: testAccCheckAciLeakInternalSubnetConfig_basic(fv_tenant_name, fv_src_vrf_name, fv_dest_vrf_name, leak_internal_subnet_ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLeakInternalSubnetExists("aci_vrf_leak_epg_bd_subnet.foo_leak_internal_subnet", &leak_internal_subnet, &leak_to_1, &leak_to_2),
					testAccCheckAciLeakInternalSubnetAttributes(fv_tenant_name, fv_src_vrf_name, fv_dest_vrf_name, leak_internal_subnet_ip, &leak_internal_subnet, &leak_to_1, &leak_to_2),
				),
			},
		},
	})
}

func testAccCheckAciLeakInternalSubnetConfig_basic(fv_tenant_name, fv_src_vrf_name, fv_dest_vrf_name, leak_internal_subnet_ip string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_vrf" "foo_src_vrf" {
		name 		= "%s"
		description = "vrf created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	resource "aci_vrf" "foo_dest_vrf" {
		name 		= "%s"
		description = "vrf created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	resource "aci_vrf_leak_epg_bd_subnet" "foo_leak_internal_subnet" {
		vrf_dn    = aci_vrf.foo_src_vrf.id
		ip        = "%s"
		leak_to {
		  vrf_dn    = "uni/tn-common/ctx-default"
		}
		leak_to {
		  vrf_dn    = aci_vrf.foo_dest_vrf.id
		  allow_l3out_advertisement   = false
		}
	}
	`, fv_tenant_name, fv_src_vrf_name, fv_dest_vrf_name, leak_internal_subnet_ip)
}

func testAccCheckAciLeakInternalSubnetExists(name string, leak_internal_subnet *models.LeakInternalSubnet, leak_to_1 *models.TenantandVRFdestinationforInterVRFLeakedRoutes, leak_to_2 *models.TenantandVRFdestinationforInterVRFLeakedRoutes) resource.TestCheckFunc {
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

		// leakTo - Beginning Read
		leakToObjects, _ := client.ListTenantandVRFdestinationforInterVRFLeakedRoutes(rs.Primary.ID)
		*leak_to_1 = *leakToObjects[0]
		*leak_to_2 = *leakToObjects[1]
		// leakTo - Read finished successfully

		return nil
	}
}

func testAccCheckAciLeakInternalSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vrf_leak_epg_bd_subnet" {
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

func testAccCheckAciLeakInternalSubnetAttributes(fv_tenant_name, fv_src_vrf_name, fv_dest_vrf_name, leak_internal_subnet_ip string, leak_internal_subnet *models.LeakInternalSubnet, leak_to_1 *models.TenantandVRFdestinationforInterVRFLeakedRoutes, leak_to_2 *models.TenantandVRFdestinationforInterVRFLeakedRoutes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if leak_internal_subnet_ip != leak_internal_subnet.Ip {
			return fmt.Errorf("Bad leak_internal_subnet IP %s", leak_internal_subnet.Ip)
		}

		if leak_internal_subnet.Scope != "private" {
			return fmt.Errorf("Bad leak_internal_subnet Scope %s", leak_internal_subnet.Scope)
		}

		if leak_to_1.DestinationCtxName != "default" && leak_to_1.DestinationCtxName != fv_dest_vrf_name {
			return fmt.Errorf("Bad leakTo - 1 Destination VRF Name %s", leak_to_1.DestinationCtxName)
		}

		if leak_to_1.DestinationTenantName != "common" && leak_to_1.DestinationTenantName != fv_tenant_name {
			return fmt.Errorf("Bad leakTo - 1 Destination Tenant Name %s", leak_to_1.DestinationTenantName)
		}

		if leak_to_2.DestinationCtxName != "default" && leak_to_2.DestinationCtxName != fv_dest_vrf_name {
			return fmt.Errorf("Bad leakTo - 2 Destination VRF Name %s", leak_to_2.DestinationCtxName)
		}

		if leak_to_2.DestinationTenantName != "common" && leak_to_2.DestinationTenantName != fv_tenant_name {
			return fmt.Errorf("Bad leakTo - 2 Destination Tenant Name %s", leak_to_2.DestinationTenantName)
		}

		return nil
	}
}
