package acctest

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBridgeDomain_Basic(t *testing.T) {
	var bridge_domain_default models.BridgeDomain
	var bridge_domain_updated models.BridgeDomain
	resourceName := "aci_bridge_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOtherName := makeTestVariable(acctest.RandString(5))
	parentOtherName := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBridgeDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccBridgeDomainWithoutTenant(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccBridgeDomainWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBridgeDomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists(resourceName, &bridge_domain_default),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "optimize_wan_bandwidth", "no"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "arp_flood", "no"),
					resource.TestCheckResourceAttr(resourceName, "ep_clear", "no"),
					//resource.TestCheckResourceAttr(resourceName, "ep_move_detect_mode", ""), no need to check these parameter for now will look into it in regression testing
					resource.TestCheckResourceAttr(resourceName, "host_based_routing", "no"),
					resource.TestCheckResourceAttr(resourceName, "intersite_bum_traffic_allow", "no"),
					resource.TestCheckResourceAttr(resourceName, "intersite_l2_stretch", "no"),
					resource.TestCheckResourceAttr(resourceName, "ip_learning", "yes"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_mcast_allow", "no"),
					resource.TestCheckResourceAttr(resourceName, "limit_ip_learn_to_subnets", "yes"),
					resource.TestCheckResourceAttr(resourceName, "ll_addr", "::"),
					// resource.TestCheckResourceAttr(resourceName, "mac", ""), no need to check these parameter as its value is APIC server dependent
					resource.TestCheckResourceAttr(resourceName, "mcast_allow", "no"),
					resource.TestCheckResourceAttr(resourceName, "multi_dst_pkt_act", "bd-flood"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "bridge_domain_type", "regular"),
					resource.TestCheckResourceAttr(resourceName, "unicast_route", "yes"),
					resource.TestCheckResourceAttr(resourceName, "unk_mac_ucast_act", "proxy"),
					resource.TestCheckResourceAttr(resourceName, "unk_mcast_act", "flood"),
					resource.TestCheckResourceAttr(resourceName, "v6unk_mcast_act", "flood"),
					resource.TestCheckResourceAttr(resourceName, "vmac", "not-applicable"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_abd_pol_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_fhs", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_netflow_monitor_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_profile", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_relay_p", ""),
				),
			},
			{
				Config: CreateAccBridgeDomainConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists(resourceName, &bridge_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "description", "test_desc"),
					resource.TestCheckResourceAttr(resourceName, "optimize_wan_bandwidth", "yes"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "test_annotation"),
					resource.TestCheckResourceAttr(resourceName, "arp_flood", "yes"),
					resource.TestCheckResourceAttr(resourceName, "ep_clear", "yes"),
					//resource.TestCheckResourceAttr(resourceName, "ep_move_detect_mode", ""),
					resource.TestCheckResourceAttr(resourceName, "host_based_routing", "yes"),
					resource.TestCheckResourceAttr(resourceName, "intersite_bum_traffic_allow", "yes"),
					resource.TestCheckResourceAttr(resourceName, "intersite_l2_stretch", "yes"),
					resource.TestCheckResourceAttr(resourceName, "ip_learning", "no"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_mcast_allow", "yes"),
					resource.TestCheckResourceAttr(resourceName, "limit_ip_learn_to_subnets", "no"),
					resource.TestCheckResourceAttr(resourceName, "ll_addr", "fe80::1"),
					//resource.TestCheckResourceAttr(resourceName, "mac", "00:22:BD:F8:19:FF"),
					resource.TestCheckResourceAttr(resourceName, "mcast_allow", "yes"),
					resource.TestCheckResourceAttr(resourceName, "multi_dst_pkt_act", "encap-flood"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_alias"),
					resource.TestCheckResourceAttr(resourceName, "bridge_domain_type", "regular"),
					resource.TestCheckResourceAttr(resourceName, "unicast_route", "no"),
					resource.TestCheckResourceAttr(resourceName, "unk_mac_ucast_act", "flood"),
					resource.TestCheckResourceAttr(resourceName, "unk_mcast_act", "opt-flood"),
					resource.TestCheckResourceAttr(resourceName, "v6unk_mcast_act", "opt-flood"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vmac", "not-applicable"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_abd_pol_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_fhs", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_netflow_monitor_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_profile", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_relay_p", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_out.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_flood_to.#", "0"),
					testAccCheckAciBridgeDomainIdEqual(&bridge_domain_default, &bridge_domain_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccBridgeDomainConfigUpdatedName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of BD-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccBridgeDomainConfigWithParentAndName(rName, rOtherName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists(resourceName, &bridge_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rOtherName),
					testAccCheckAciBridgeDomainIdNotEqual(&bridge_domain_default, &bridge_domain_updated),
				),
			},
			{
				Config: CreateAccBridgeDomainConfig(rName),
			},
			{
				Config: CreateAccBridgeDomainConfigWithParentAndName(parentOtherName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists(resourceName, &bridge_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", parentOtherName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciBridgeDomainIdNotEqual(&bridge_domain_default, &bridge_domain_updated),
				),
			},
		},
	})
}

func TestAccAciBridgeDomain_Update(t *testing.T) {
	var bridge_domain_default models.BridgeDomain
	var bridge_domain_updated models.BridgeDomain
	resourceName := "aci_bridge_domain.test"
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBridgeDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBridgeDomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists(resourceName, &bridge_domain_default),
				),
			},
			{
				Config: CreateAccBridgeDomainUpdatedAttr(rName, "multi_dst_pkt_act", "drop"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists(resourceName, &bridge_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "multi_dst_pkt_act", "drop"),
					testAccCheckAciBridgeDomainIdEqual(&bridge_domain_default, &bridge_domain_updated),
				),
			},
			{
				Config: CreateAccBridgeDomainUpdatedAttr(rName, "unicast_route", "no"),
			},
			{
				Config: CreateAccBridgeDomainUpdatedAttr(rName, "bridge_domain_type", "fc"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists(resourceName, &bridge_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "bridge_domain_type", "fc"),
					testAccCheckAciBridgeDomainIdEqual(&bridge_domain_default, &bridge_domain_updated),
				),
			},
		},
	})
}

func TestAccAciBridgeDomain_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longAnnotationDesc := acctest.RandString(129)
	longNameAlias := acctest.RandString(65)
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBridgeDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBridgeDomainConfig(rName),
			},
			{
				Config:      CreateAccBridgeDomainWithInvalidTenant(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvBD (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "description", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "annotation", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "arp_flood", randomValue),
				ExpectError: regexp.MustCompile(`expected arp_flood to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "optimize_wan_bandwidth", randomValue),
				ExpectError: regexp.MustCompile(`expected optimize_wan_bandwidth to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "ep_clear", randomValue),
				ExpectError: regexp.MustCompile(`expected ep_clear to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "host_based_routing", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name hostBasedRouting, class fvBD (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "intersite_bum_traffic_allow", randomValue),
				ExpectError: regexp.MustCompile(`expected intersite_bum_traffic_allow to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "intersite_l2_stretch", randomValue),
				ExpectError: regexp.MustCompile(`expected intersite_l2_stretch to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "ip_learning", randomValue),
				ExpectError: regexp.MustCompile(`expected ip_learning to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "ipv6_mcast_allow", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name ipv6McastAllow, class fvBD (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "limit_ip_learn_to_subnets", randomValue),
				ExpectError: regexp.MustCompile(`expected limit_ip_learn_to_subnets to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "ll_addr", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name llAddr, class fvBD (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "mcast_allow", randomValue),
				ExpectError: regexp.MustCompile(`expected mcast_allow to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "multi_dst_pkt_act", randomValue),
				ExpectError: regexp.MustCompile(`expected multi_dst_pkt_act to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "bridge_domain_type", randomValue),
				ExpectError: regexp.MustCompile(`expected bridge_domain_type to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "unicast_route", randomValue),
				ExpectError: regexp.MustCompile(`expected unicast_route to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "unk_mac_ucast_act", randomValue),
				ExpectError: regexp.MustCompile(`expected unk_mac_ucast_act to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "unk_mcast_act", randomValue),
				ExpectError: regexp.MustCompile(`expected unk_mcast_act to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "v6unk_mcast_act", randomValue),
				ExpectError: regexp.MustCompile(`expected v6unk_mcast_act to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, "vmac", randomValue),
				ExpectError: regexp.MustCompile(`expected vmac to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccBridgeDomainUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccBridgeDomainForUnicastAndType(rName),
				ExpectError: regexp.MustCompile(`Invalid Configuration : Unicast Routing is not allowed for FC/FCOE BD`),
			},
			{
				Config: CreateAccBridgeDomainConfig(rName),
			},
		},
	})
}

func TestAccAciBridgeDomain_RelationParameters(t *testing.T) {
	var bridge_domain_default models.BridgeDomain
	var bridge_domain_rel1 models.BridgeDomain
	var bridge_domain_rel2 models.BridgeDomain
	resourceName := "aci_bridge_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomName1 := makeTestVariable(acctest.RandString(5))
	randomName2 := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBridgeDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBridgeDomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists(resourceName, &bridge_domain_default),
				),
			},
			{
				Config: CreateAccBridgeDomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists(resourceName, &bridge_domain_default),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_abd_pol_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_fhs", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_netflow_monitor_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_profile", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_relay_p", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_out.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_flood_to.#", "0"),
				),
			},
			{
				Config: CreateAccBridgeDomainRelConfigInitial(rName, randomName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists(resourceName, &bridge_domain_rel1),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_netflow_monitor_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_fhs", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_abd_pol_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_flood_to.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_bd_flood_to.*", fmt.Sprintf("uni/tn-%s/flt-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_ep_ret", fmt.Sprintf("uni/tn-%s/epRPol-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_out.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_bd_to_out.*", fmt.Sprintf("uni/tn-%s/out-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_profile", fmt.Sprintf("uni/tn-%s/prof-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_relay_p", fmt.Sprintf("uni/tn-%s/relayp-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx", fmt.Sprintf("uni/tn-%s/ctx-%s", rName, randomName1)),
					testAccCheckAciBridgeDomainIdEqual(&bridge_domain_default, &bridge_domain_rel1),
				),
			},
			{
				Config: CreateAccBridgeDomainRelConfigFinal(rName, randomName1, randomName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists(resourceName, &bridge_domain_rel2),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_netflow_monitor_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_fhs", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_abd_pol_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_flood_to.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_bd_flood_to.*", fmt.Sprintf("uni/tn-%s/flt-%s", rName, randomName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_bd_flood_to.*", fmt.Sprintf("uni/tn-%s/flt-%s", rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_ep_ret", fmt.Sprintf("uni/tn-%s/epRPol-%s", rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_out.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_bd_to_out.*", fmt.Sprintf("uni/tn-%s/out-%s", rName, randomName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_bd_to_out.*", fmt.Sprintf("uni/tn-%s/out-%s", rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_profile", fmt.Sprintf("uni/tn-%s/prof-%s", rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_relay_p", fmt.Sprintf("uni/tn-%s/relayp-%s", rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx", fmt.Sprintf("uni/tn-%s/ctx-%s", rName, randomName2)),
					testAccCheckAciBridgeDomainIdEqual(&bridge_domain_default, &bridge_domain_rel2),
				),
			},
			{
				Config: CreateAccBridgeDomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists(resourceName, &bridge_domain_default),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_ep_ret", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_abd_pol_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_fhs", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_netflow_monitor_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_profile", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_relay_p", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_to_out.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd_flood_to.#", "0"),
				),
			},
		},
	})
}

func TestAccBridgeDomain_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBridgeDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBDsConfig(rName),
			},
		},
	})
}

func CreateAccBDsConfig(rName string) string {
	fmt.Println("=== STEP  creating multiple bridge domains")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_bridge_domain" "test1"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_bridge_domain" "test2"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_bridge_domain" "test3"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	`, rName, rName+"1", rName+"2", rName+"3")
	return resource

}

func CreateAccBridgeDomainForUnicastAndType(rName string) string {
	fmt.Println("=== STEP  testing bridge domain with unicast_route=yes and bridge_domain_type=fc")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
		unicast_route= "yes"
		bridge_domain_type = "fc"
	}
	`, rName, rName)
	return resource

}

func CreateAccBridgeDomainRelConfigFinal(rName, relName1, relName2 string) string {
	fmt.Println("=== STEP  testing bridge domain with final relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bgp_route_control_profile" "test"{
		parent_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_monitoring_policy" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter" "test1"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter" "test2"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_dhcp_relay_policy" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_vrf" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_end_point_retention_policy" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_l3_outside" "test1"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_l3_outside" "test2"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
		relation_fv_rs_bd_to_profile = aci_bgp_route_control_profile.test.id
		relation_fv_rs_abd_pol_mon_pol = aci_monitoring_policy.test.id
		relation_fv_rs_bd_flood_to = [aci_filter.test1.id, aci_filter.test2.id]
		relation_fv_rs_bd_to_relay_p = aci_dhcp_relay_policy.test.id
		relation_fv_rs_ctx = aci_vrf.test.id
		relation_fv_rs_bd_to_ep_ret = aci_end_point_retention_policy.test.id
		relation_fv_rs_bd_to_out = [aci_l3_outside.test1.id, aci_l3_outside.test2.id]
	}
	`, rName, relName2, relName2, relName1, relName2, relName2, relName2, relName2, relName1, relName2, rName)
	return resource

}

func CreateAccBridgeDomainRelConfigInitial(rName, relName string) string {
	fmt.Println("=== STEP  testing bridge domain with final relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bgp_route_control_profile" "test"{
		parent_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_monitoring_policy" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter" "test1"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_dhcp_relay_policy" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_vrf" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_end_point_retention_policy" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_l3_outside" "test1"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
		relation_fv_rs_bd_to_profile = aci_bgp_route_control_profile.test.id
		relation_fv_rs_abd_pol_mon_pol = aci_monitoring_policy.test.id
		relation_fv_rs_bd_flood_to = [aci_filter.test1.id]
		relation_fv_rs_bd_to_relay_p = aci_dhcp_relay_policy.test.id
		relation_fv_rs_ctx = aci_vrf.test.id
		relation_fv_rs_bd_to_ep_ret = aci_end_point_retention_policy.test.id
		relation_fv_rs_bd_to_out = [aci_l3_outside.test1.id]
	}
	`, rName, relName, relName, relName, relName, relName, relName, relName, rName)
	return resource

}

func CreateAccBridgeDomainWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing bridge domain without creating tenant")
	resource := fmt.Sprintf(`
	resource "aci_bridge_domain" "test" {
		name = "%s"
	}
	`, rName)
	return resource

}

func CreateAccBridgeDomainWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing bridge domain without passing name attribute")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		tenant_dn = aci_tenant.test.id
	}
	`, rName)
	return resource
}

func CreateAccBridgeDomainConfigWithParentAndName(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing bridge domain creation with tenant name %s and bridge domain name %s\n", parentName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, parentName, rName)
	return resource
}

func CreateAccBridgeDomainConfig(rName string) string {
	fmt.Println("=== STEP  testing bridge domain creation with required parameters only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccBridgeDomainConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing bridge domain creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		annotation= "test_annotation"
		arp_flood= "yes"
		description= "test_desc"
		ep_clear="yes"
		host_based_routing="yes"
		intersite_bum_traffic_allow="yes"
		intersite_l2_stretch="yes"
		ip_learning="no"
		ipv6_mcast_allow="yes"
		limit_ip_learn_to_subnets="no"
		ll_addr="fe80::1"
		mcast_allow="yes"
		multi_dst_pkt_act="encap-flood"
		name_alias="test_alias"
		optimize_wan_bandwidth="yes"
		unicast_route="no"
		unk_mac_ucast_act="flood"
		unk_mcast_act="opt-flood"
		v6unk_mcast_act="opt-flood"
	  }
	`, rName, rName)
	return resource
}

func CreateAccBridgeDomainConfigUpdatedName(rName, longrName string) string {
	fmt.Println("=== STEP  Basic: testing bridge domain creation with invalid name with long length")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	`, rName, longrName)
	return resource
}

func CreateAccBridgeDomainUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing bridge domain attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_bridge_domain" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}

	`, rName, rName, attribute, value)
	return resource
}

func CreateAccBridgeDomainWithInvalidTenant(rName string) string {
	fmt.Println("=== STEP  testing bridge_domain updation with invalid tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_bridge_domain" "test"{
		tenant_dn = aci_application_profile.test.id
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func testAccCheckAciBridgeDomainDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing bridge domain destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_bridge_domain" {
			cont, err := client.Get(rs.Primary.ID)
			bridge_domain := models.BridgeDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Bridge Domain %s still exists", bridge_domain.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBridgeDomainIdNotEqual(bd1, bd2 *models.BridgeDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if bd1.DistinguishedName == bd2.DistinguishedName {
			return fmt.Errorf("Bridge Domain DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciBridgeDomainIdEqual(bd1, bd2 *models.BridgeDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if bd1.DistinguishedName != bd2.DistinguishedName {
			return fmt.Errorf("Bridge Domain DNs are no equal")
		}
		return nil
	}
}

func testAccCheckAciBridgeDomainExists(name string, bridge_domain *models.BridgeDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Bridge Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Bridge Domain Dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bridge_domainFound := models.BridgeDomainFromContainer(cont)
		if bridge_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Bridge Domain %s not found", rs.Primary.ID)
		}
		*bridge_domain = *bridge_domainFound
		return nil
	}
}
