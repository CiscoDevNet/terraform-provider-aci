package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciBridgeDomainDataSource_Basic(t *testing.T) {
	resourceName := "aci_bridge_domain.test"
	dataSourceName := "data.aci_bridge_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBridgeDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccBridgeDomainDSWithoutTenant(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccBridgeDomainDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBridgeDomainConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "optimize_wan_bandwidth", resourceName, "optimize_wan_bandwidth"),
					resource.TestCheckResourceAttrPair(dataSourceName, "arp_flood", resourceName, "arp_flood"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ep_clear", resourceName, "ep_clear"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ep_move_detect_mode", resourceName, "ep_move_detect_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "host_based_routing", resourceName, "host_based_routing"),
					resource.TestCheckResourceAttrPair(dataSourceName, "intersite_bum_traffic_allow", resourceName, "intersite_bum_traffic_allow"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "intersite_l2_stretch", resourceName, "intersite_l2_stretch"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ip_learning", resourceName, "ip_learning"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ipv6_mcast_allow", resourceName, "ipv6_mcast_allow"),
					resource.TestCheckResourceAttrPair(dataSourceName, "limit_ip_learn_to_subnets", resourceName, "limit_ip_learn_to_subnets"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ll_addr", resourceName, "ll_addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mac", resourceName, "mac"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mcast_allow", resourceName, "mcast_allow"),
					resource.TestCheckResourceAttrPair(dataSourceName, "multi_dst_pkt_act", resourceName, "multi_dst_pkt_act"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bridge_domain_type", resourceName, "bridge_domain_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "unicast_route", resourceName, "unicast_route"),
					resource.TestCheckResourceAttrPair(dataSourceName, "unk_mac_ucast_act", resourceName, "unk_mac_ucast_act"),
					resource.TestCheckResourceAttrPair(dataSourceName, "unk_mcast_act", resourceName, "unk_mcast_act"),
					resource.TestCheckResourceAttrPair(dataSourceName, "v6unk_mcast_act", resourceName, "v6unk_mcast_act"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vmac", resourceName, "vmac"),
				),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedConfigDataSourceRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccBridgeDomainUpdatedConfigDataSource(rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
			{
				Config:      CreateAccBridgeDomainDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
		},
	})
}

func CreateAccBridgeDomainUpdatedConfigDataSourceRandomAttr(rName, attribute, value string) string {
	fmt.Println("=== STEP  testing bridge domain data source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = aci_bridge_domain.test.name
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccBridgeDomainUpdatedConfigDataSource(rName, attribute, value string) string {
	fmt.Println("=== STEP  testing bridge domain data source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}

	data "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = aci_bridge_domain.test.name
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccBridgeDomainConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing bridge domain data source")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		description = "test_description"
	}

	data "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = aci_bridge_domain.test.name
	}
	`, rName, rName)
	return resource
}

func CreateAccBridgeDomainDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing bridge domain reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "${aci_tenant.test.name}abc"
	}
	`, rName, rName)
	return resource
}

func CreateAccBridgeDomainDSWithoutTenant(rName string) string {
	fmt.Println("=== STEP  testing bridge domain reading without giving tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_bridge_domain" "test" {
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccBridgeDomainDSWithoutName(rName string) string {
	fmt.Println("=== STEP  testing bridge domain reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
	}
	`, rName, rName)
	return resource
}
