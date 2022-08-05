package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSystemDataSource_Basic(t *testing.T) {
	dataSourceName := "data.aci_system.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      CreateSystemDSWithoutRequired(systemPodId, systemNodeId, "pod_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSystemDSWithoutRequired(systemPodId, systemNodeId, "system_id"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSystemConfigDataSource(systemPodId, systemNodeId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "boot_strap_state"),
					resource.TestCheckResourceAttrSet(dataSourceName, "control_plane_mtu"),
					resource.TestCheckResourceAttrSet(dataSourceName, "current_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "enforce_subnet_check"),
					resource.TestCheckResourceAttrSet(dataSourceName, "etep_addr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "fabric_domain"),
					resource.TestCheckResourceAttrSet(dataSourceName, "fabric_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "fabric_mac"),
					resource.TestCheckResourceAttrSet(dataSourceName, "inb_mgmt_addr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "inb_mgmt_addr6"),
					resource.TestCheckResourceAttrSet(dataSourceName, "inb_mgmt_addr6_mask"),
					resource.TestCheckResourceAttrSet(dataSourceName, "inb_mgmt_addr_mask"),
					resource.TestCheckResourceAttrSet(dataSourceName, "inb_mgmt_gateway"),
					resource.TestCheckResourceAttrSet(dataSourceName, "inb_mgmt_gateway6"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_reboot_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_reset_reason"),
					resource.TestCheckResourceAttrSet(dataSourceName, "lc_own"),
					resource.TestCheckResourceAttrSet(dataSourceName, "mod_ts"),
					resource.TestCheckResourceAttrSet(dataSourceName, "mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "mon_pol_dn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "node_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "oob_mgmt_addr"),
					resource.TestCheckResourceAttrSet(dataSourceName, "oob_mgmt_addr6"),
					resource.TestCheckResourceAttrSet(dataSourceName, "oob_mgmt_addr6_mask"),
					resource.TestCheckResourceAttrSet(dataSourceName, "oob_mgmt_addr_mask"),
					resource.TestCheckResourceAttrSet(dataSourceName, "oob_mgmt_gateway"),
					resource.TestCheckResourceAttrSet(dataSourceName, "oob_mgmt_gateway6"),
					resource.TestCheckResourceAttr(dataSourceName, "pod_id", systemPodId),
					resource.TestCheckResourceAttrSet(dataSourceName, "remote_network_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "remote_node"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rl_oper_pod_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rl_routable_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rldirect_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "role"),
					resource.TestCheckResourceAttrSet(dataSourceName, "serial"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "site_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "state"),
					resource.TestCheckResourceAttr(dataSourceName, "system_id", systemNodeId),
					resource.TestCheckResourceAttrSet(dataSourceName, "system_uptime"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tep_pool"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unicast_xr_ep_learn_disable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_mode"),
				),
			},
			{
				Config:      CreateAccSystemDataSourceUpdate(systemPodId, systemNodeId, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccSystemDSWithInvalidSystemId(systemPodId, randomValue),
				ExpectError: regexp.MustCompile(`Invalid RN`),
			},
			{
				Config: CreateAccSystemConfigDataSource(systemPodId, systemNodeId),
			},
		},
	})
}

func CreateAccSystemConfigDataSource(pid, sid string) string {
	fmt.Println("=== STEP  testing system Data Source with required arguments only")
	resource := fmt.Sprintf(`
	data "aci_system" "test" {
		pod_id  = "%s"
		system_id = "%s"
	}
	`, pid, sid)
	return resource
}

func CreateSystemDSWithoutRequired(pid, sid, attrName string) string {
	fmt.Println("=== STEP  Basic: testing system Data Source without ", attrName)
	rBlock := `
	`
	switch attrName {
	case "pod_id":
		rBlock += `
		data "aci_system" "test" {
		#	pod_id = "%s"
			system_id = "%s"
		}
		`
	case "system_id":
		rBlock += `
		data "aci_system" "test" {
			pod_id = "%s"
		#	system_id = "%s"
		}
		`
	}
	return fmt.Sprintf(rBlock, pid, sid)
}

func CreateAccSystemDSWithInvalidSystemId(pid, sid string) string {
	fmt.Println("=== STEP  testing system Data Source with invalid system id")
	resource := fmt.Sprintf(`
	data "aci_system" "test" {
		pod_id = "%s"
		system_id = "%s"
	}
	`, pid, sid)
	return resource
}

func CreateAccSystemDataSourceUpdate(pid, sid, key, value string) string {
	fmt.Println("=== STEP  testing system Data Source with random attribute")
	resource := fmt.Sprintf(`
	data "aci_system" "test" {
		pod_id = "%s"
		system_id = "%s"
		%s = "%s"
	  }
	`, pid, sid, key, value)
	return resource
}
