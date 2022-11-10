package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBridgeDomain_Basic(t *testing.T) {
	var bridge_domain models.BridgeDomain
	description := "bridge_domain created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBridgeDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBridgeDomainConfig_basic(description, "yes"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists("aci_bridge_domain.foobridge_domain", &bridge_domain),
					testAccCheckAciBridgeDomainAttributes(description, "yes", &bridge_domain),
				),
			},
		},
	})
}

func TestAccAciBridgeDomain_update(t *testing.T) {
	var bridge_domain models.BridgeDomain
	description := "bridge_domain created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBridgeDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBridgeDomainConfig_basic(description, "yes"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists("aci_bridge_domain.foobridge_domain", &bridge_domain),
					testAccCheckAciBridgeDomainAttributes(description, "yes", &bridge_domain),
				),
			},
			{
				Config: testAccCheckAciBridgeDomainConfig_basic(description, "no"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists("aci_bridge_domain.foobridge_domain", &bridge_domain),
					testAccCheckAciBridgeDomainAttributes(description, "no", &bridge_domain),
				),
			},
		},
	})
}

func testAccCheckAciBridgeDomainConfig_basic(description, ipv6_mcast_allow string) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "tenant_for_bd" {
		name        = "tenant_for_bd"
		description = "This tenant is created by terraform ACI provider"
	}
	resource "aci_bridge_domain" "foobridge_domain" {
		tenant_dn                   = aci_tenant.tenant_for_bd.id
		description                 = "%s"
		name                        = "demo_bd"
		optimize_wan_bandwidth      = "no"
		annotation                  = "tag_bd"
		arp_flood                   = "no"
		ep_clear                    = "no"
		ep_move_detect_mode         = "garp"
		host_based_routing          = "no"
		intersite_bum_traffic_allow = "yes"
		intersite_l2_stretch        = "yes"
		ip_learning                 = "yes"
		ipv6_mcast_allow            = "%s"
		limit_ip_learn_to_subnets   = "yes"
		mac                         = "00:22:BD:F8:19:FF"
		mcast_allow                 = "yes"
		multi_dst_pkt_act           = "bd-flood"
		name_alias                  = "alias_bd"
		bridge_domain_type          = "regular"
		unicast_route               = "no"
		unk_mac_ucast_act           = "flood"
		unk_mcast_act               = "flood"
		vmac                        = "not-applicable"
	}
	`, description, ipv6_mcast_allow)
}

func testAccCheckAciBridgeDomainExists(name string, bridge_domain *models.BridgeDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Bridge Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Bridge Domain dn was set")
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

func testAccCheckAciBridgeDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_bridge_domain" {
			cont, err := client.Get(rs.Primary.ID)
			bridge_domain := models.BridgeDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Bridge Domain %s Still exists", bridge_domain.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciBridgeDomainAttributes(description, ipv6_mcast_allow string, bridge_domain *models.BridgeDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != bridge_domain.Description {
			return fmt.Errorf("Bad bridge_domain Description %s", bridge_domain.Description)
		}

		if "demo_bd" != bridge_domain.Name {
			return fmt.Errorf("Bad bridge_domain name %s", bridge_domain.Name)
		}

		if "no" != bridge_domain.OptimizeWanBandwidth {
			return fmt.Errorf("Bad bridge_domain optimize_wan_bandwidth %s", bridge_domain.OptimizeWanBandwidth)
		}

		if "tag_bd" != bridge_domain.Annotation {
			return fmt.Errorf("Bad bridge_domain annotation %s", bridge_domain.Annotation)
		}

		if "no" != bridge_domain.ArpFlood {
			return fmt.Errorf("Bad bridge_domain arp_flood %s", bridge_domain.ArpFlood)
		}

		if "no" != bridge_domain.EpClear {
			return fmt.Errorf("Bad bridge_domain ep_clear %s", bridge_domain.EpClear)
		}

		if "garp" != bridge_domain.EpMoveDetectMode {
			return fmt.Errorf("Bad bridge_domain ep_move_detect_mode %s", bridge_domain.EpMoveDetectMode)
		}

		if "no" != bridge_domain.HostBasedRouting {
			return fmt.Errorf("Bad bridge_domain host_based_routing %s", bridge_domain.HostBasedRouting)
		}

		if "yes" != bridge_domain.IntersiteBumTrafficAllow {
			return fmt.Errorf("Bad bridge_domain intersite_bum_traffic_allow %s", bridge_domain.IntersiteBumTrafficAllow)
		}

		if "yes" != bridge_domain.IntersiteL2Stretch {
			return fmt.Errorf("Bad bridge_domain intersite_l2_stretch %s", bridge_domain.IntersiteL2Stretch)
		}

		if "yes" != bridge_domain.IpLearning {
			return fmt.Errorf("Bad bridge_domain ip_learning %s", bridge_domain.IpLearning)
		}

		if ipv6_mcast_allow != bridge_domain.Ipv6McastAllow {
			return fmt.Errorf("Bad bridge_domain ipv6_mcast_allow %s", bridge_domain.Ipv6McastAllow)
		}

		if "yes" != bridge_domain.LimitIpLearnToSubnets {
			return fmt.Errorf("Bad bridge_domain limit_ip_learn_to_subnets %s", bridge_domain.LimitIpLearnToSubnets)
		}

		if "00:22:BD:F8:19:FF" != bridge_domain.Mac {
			return fmt.Errorf("Bad bridge_domain mac %s", bridge_domain.Mac)
		}

		if "yes" != bridge_domain.McastAllow {
			return fmt.Errorf("Bad bridge_domain mcast_allow %s", bridge_domain.McastAllow)
		}

		if "bd-flood" != bridge_domain.MultiDstPktAct {
			return fmt.Errorf("Bad bridge_domain multi_dst_pkt_act %s", bridge_domain.MultiDstPktAct)
		}

		if "alias_bd" != bridge_domain.NameAlias {
			return fmt.Errorf("Bad bridge_domain name_alias %s", bridge_domain.NameAlias)
		}

		if "regular" != bridge_domain.BridgeDomain_type {
			return fmt.Errorf("Bad bridge_domain bridge_domain_type %s", bridge_domain.BridgeDomain_type)
		}

		if "no" != bridge_domain.UnicastRoute {
			return fmt.Errorf("Bad bridge_domain unicast_route %s", bridge_domain.UnicastRoute)
		}

		if "flood" != bridge_domain.UnkMacUcastAct {
			return fmt.Errorf("Bad bridge_domain unk_mac_ucast_act %s", bridge_domain.UnkMacUcastAct)
		}

		if "flood" != bridge_domain.UnkMcastAct {
			return fmt.Errorf("Bad bridge_domain unk_mcast_act %s", bridge_domain.UnkMcastAct)
		}

		if "not-applicable" != bridge_domain.Vmac {
			return fmt.Errorf("Bad bridge_domain vmac %s", bridge_domain.Vmac)
		}

		return nil
	}
}
