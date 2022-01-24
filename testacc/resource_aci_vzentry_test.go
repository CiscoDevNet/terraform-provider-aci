package testacc

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

func TestAccAciFilterEntry_Basic(t *testing.T) {
	var filter_entry_default models.FilterEntry
	var filter_entry_updated models.FilterEntry
	resourceName := "aci_filter_entry.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOtherName := makeTestVariable(acctest.RandString(5))
	parentOtherName := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFilterEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccFilterEntrynWithoutFilter(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccFilterEntryWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFilterEntryConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_default),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "apply_to_frag", "no"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "arp_opc", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "d_from_port", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "ether_t", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "icmpv4_t", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "icmpv6_t", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "prot", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "s_from_port", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "s_to_port", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "stateful", "no"),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.0", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "filter_dn", fmt.Sprintf("uni/tn-%s/flt-%s", rName, rName)),
				),
			},
			{
				Config: CreateAccFilterEntryConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "test_annotation"),
					resource.TestCheckResourceAttr(resourceName, "arp_opc", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "apply_to_frag", "no"),
					resource.TestCheckResourceAttr(resourceName, "d_from_port", "https"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "https"),
					resource.TestCheckResourceAttr(resourceName, "ether_t", "ip"),
					resource.TestCheckResourceAttr(resourceName, "icmpv4_t", "echo-rep"),
					resource.TestCheckResourceAttr(resourceName, "icmpv6_t", "dst-unreach"),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "CS0"),
					resource.TestCheckResourceAttr(resourceName, "prot", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "s_from_port", "https"),
					resource.TestCheckResourceAttr(resourceName, "s_to_port", "https"),
					resource.TestCheckResourceAttr(resourceName, "stateful", "yes"),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.0", "est"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_name_alias"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "filter_dn", fmt.Sprintf("uni/tn-%s/flt-%s", rName, rName)),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccFilterEntryConfigUpdatedName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of e-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccFilterEntryConfigWithParentAndName(rName, rOtherName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "filter_dn", fmt.Sprintf("uni/tn-%s/flt-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rOtherName),
					testAccCheckAciFilterEntryIdNotEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryConfig(rName),
			},
			{
				Config: CreateAccFilterEntryConfigWithParentAndName(parentOtherName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "filter_dn", fmt.Sprintf("uni/tn-%s/flt-%s", parentOtherName, parentOtherName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciFilterEntryIdNotEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
		},
	})
}

func TestAccAciFilterEntry_Update(t *testing.T) {
	var filter_entry_default models.FilterEntry
	var filter_entry_updated models.FilterEntry
	resourceName := "aci_filter_entry.test"
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFilterEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFilterEntryConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_default),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "ether_t", "ip"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "ether_t", "ip"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "apply_to_frag", "yes"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "apply_to_frag", "yes"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "apply_to_frag", "no"),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "ether_t", "arp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "ether_t", "arp"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "arp_opc", "req"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "arp_opc", "req"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "arp_opc", "reply"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "arp_opc", "reply"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "arp_opc", "unspecified"),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "ether_t", "ipv4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "ether_t", "ipv4"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "ether_t", "trill"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "ether_t", "trill"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "ether_t", "mpls_ucast"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "ether_t", "mpls_ucast"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "ether_t", "mac_security"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "ether_t", "mac_security"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "ether_t", "fcoe"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "ether_t", "fcoe"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "ether_t", "ipv6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "ether_t", "ipv6"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "icmpv4_t", "dst-unreach"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "icmpv4_t", "dst-unreach"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "icmpv4_t", "echo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "icmpv4_t", "echo"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "icmpv4_t", "time-exceeded"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "icmpv4_t", "time-exceeded"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "icmpv4_t", "src-quench"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "icmpv4_t", "src-quench"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "icmpv6_t", "time-exceeded"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "icmpv6_t", "time-exceeded"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "icmpv6_t", "echo-req"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "icmpv6_t", "echo-req"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "icmpv6_t", "nbr-solicit"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "icmpv6_t", "nbr-solicit"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "icmpv6_t", "nbr-advert"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "icmpv6_t", "nbr-advert"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "icmpv6_t", "redirect"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "icmpv6_t", "redirect"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "CS1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "CS1"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "AF11"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "AF11"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "AF12"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "AF12"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "AF13"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "AF13"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "CS2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "CS2"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "AF21"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "AF21"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "AF22"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "AF22"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "AF23"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "AF23"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "CS3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "CS3"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "AF31"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "AF31"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "AF32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "AF32"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "AF33"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "AF33"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "CS4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "CS4"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "AF42"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "AF42"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "AF43"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "AF43"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "VA"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "VA"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", "EF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "match_dscp", "EF"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "prot", "icmp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "prot", "icmp"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "prot", "igmp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "prot", "igmp"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "prot", "egp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "prot", "egp"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "prot", "igp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "prot", "igp"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "prot", "icmpv6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "prot", "icmpv6"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "prot", "eigrp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "prot", "eigrp"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "prot", "ospfigp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "prot", "ospfigp"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "prot", "pim"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "prot", "pim"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "prot", "l2tp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "prot", "l2tp"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "prot", "udp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "prot", "udp"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttr(rName, "prot", "tcp"),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttrList(rName, "tcp_rules", StringListtoString([]string{"syn"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.0", "syn"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttrList(rName, "tcp_rules", StringListtoString([]string{"ack"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.0", "ack"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttrList(rName, "tcp_rules", StringListtoString([]string{"fin"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.0", "fin"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttrList(rName, "tcp_rules", StringListtoString([]string{"rst"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.0", "rst"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttrList(rName, "tcp_rules", StringListtoString([]string{"syn"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.0", "syn"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedAttrList(rName, "tcp_rules", StringListtoString([]string{"rst", "fin"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.0", "rst"),
					resource.TestCheckResourceAttr(resourceName, "tcp_rules.1", "fin"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedPortAttr(rName, "ftpData"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "d_from_port", "ftpData"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "ftpData"),
					resource.TestCheckResourceAttr(resourceName, "s_from_port", "ftpData"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "ftpData"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedPortAttr(rName, "smtp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "d_from_port", "smtp"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "smtp"),
					resource.TestCheckResourceAttr(resourceName, "s_from_port", "smtp"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "smtp"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedPortAttr(rName, "dns"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "d_from_port", "dns"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "dns"),
					resource.TestCheckResourceAttr(resourceName, "s_from_port", "dns"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "dns"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedPortAttr(rName, "http"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "d_from_port", "http"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "http"),
					resource.TestCheckResourceAttr(resourceName, "s_from_port", "http"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "http"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedPortAttr(rName, "pop3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "d_from_port", "pop3"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "pop3"),
					resource.TestCheckResourceAttr(resourceName, "s_from_port", "pop3"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "pop3"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
			{
				Config: CreateAccFilterEntryUpdatedPortAttr(rName, "rtsp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterEntryExists(resourceName, &filter_entry_updated),
					resource.TestCheckResourceAttr(resourceName, "d_from_port", "rtsp"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "rtsp"),
					resource.TestCheckResourceAttr(resourceName, "s_from_port", "rtsp"),
					resource.TestCheckResourceAttr(resourceName, "d_to_port", "rtsp"),
					testAccCheckAciFilterEntryIdEqual(&filter_entry_default, &filter_entry_updated),
				),
			},
		},
	})
}

func TestAccAciFilterEntry_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longAnnotationDesc := acctest.RandString(129)
	longNameAlias := acctest.RandString(65)
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFilterEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFilterEntryConfig(rName),
			},
			{
				Config:      CreateAccFilterEntryWithInvalidFilter(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "description", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "annotation", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "apply_to_frag", randomValue),
				ExpectError: regexp.MustCompile(`expected apply_to_frag to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "arp_opc", randomValue),
				ExpectError: regexp.MustCompile(`expected arp_opc to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "d_from_port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dFromPort, class vzEntry (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "d_to_port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dToPort, class vzEntry (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "ether_t", randomValue),
				ExpectError: regexp.MustCompile(`expected ether_t to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "icmpv4_t", randomValue),
				ExpectError: regexp.MustCompile(`expected icmpv4_t to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "icmpv6_t", randomValue),
				ExpectError: regexp.MustCompile(`expected icmpv6_t to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "match_dscp", randomValue),
				ExpectError: regexp.MustCompile(`expected match_dscp to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "prot", randomValue),
				ExpectError: regexp.MustCompile(`expected prot to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "s_from_port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name sFromPort, class vzEntry (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "s_to_port", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name sToPort, class vzEntry (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, "stateful", randomValue),
				ExpectError: regexp.MustCompile(`expected stateful to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttrList(rName, "tcp_rules", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)* to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttrList(rName, "tcp_rules", StringListtoString([]string{"ack", "ack"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttrList(rName, "tcp_rules", StringListtoString([]string{"unspecified", "ack"})),
				ExpectError: regexp.MustCompile(`unspecified should not be used along with other values`),
			},
			{
				Config:      CreateAccFilterEntryUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFilterEntryConfig(rName),
			},
		},
	})
}

func TestAccAciFilterEntry_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFilterEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFilterEntriesConfig(rName),
			},
		},
	})
}

func CreateAccFilterEntrynWithoutFilter(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter entry without creating filter")
	resource := fmt.Sprintf(`
	resource "aci_filter_entry" "test" {
		name = "%s"
	}
	`, rName)
	return resource

}

func CreateAccFilterEntryWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter entry without passing name attribute")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
	}
	`, rName, rName)
	return resource
}

func CreateAccFilterEntryConfig(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter entry creation with required paramters only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccFilterEntryConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter entry creation with optional paramters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
		description = "test_description"
		annotation = "test_annotation"
		name_alias = "test_name_alias"
		ether_t = "ip"
		stateful = "yes"
		apply_to_frag = "no"
		arp_opc = "unspecified"
		d_from_port = "https"
		d_to_port = "https"
		icmpv4_t = "echo-rep"
		icmpv6_t = "dst-unreach"
		match_dscp = "CS0"
		prot = "tcp"
		s_from_port = "https"
		s_to_port = "https"
		tcp_rules = ["est"]
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccFilterEntryConfigUpdatedName(rName, longerName string) string {
	fmt.Println("=== STEP  Basic: testing filter entry creation with invalid name with long lenght")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
	}
	`, rName, rName, longerName)
	return resource
}

func CreateAccFilterEntryConfigWithParentAndName(parentName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing filter entry creation with filter name %s and filter entry name %s\n", parentName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
	}
	`, parentName, parentName, rName)
	return resource
}

func CreateAccFilterEntryUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing filter entry %s = %s\n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
		%s = "%s"
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccFilterEntryUpdatedAttrList(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing filter entry %s = %s\n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
		%s = %s
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccFilterEntryUpdatedPortAttr(rName, value string) string {
	fmt.Printf("=== STEP  Basic: testing filter entry ports %s\n", value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name = "%s"
		prot = "tcp"
		d_from_port = "%s"
		d_to_port = "%s"
		s_from_port = "%s"
		s_to_port = "%s"
	}
	`, rName, rName, rName, value, value, value, value)
	return resource
}

func CreateAccFilterEntryWithInvalidFilter(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter entry updation with Invalid filter_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccFilterEntriesConfig(rName string) string {
	fmt.Println("=== STEP  creating multiple filter entries")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_filter_entry" "test"{
		filter_dn = aci_filter.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName, rName, rName)
	return resource

}

func testAccCheckAciFilterEntryDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing Filter Entry destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_filter_entry" {
			cont, err := client.Get(rs.Primary.ID)
			filter_entry := models.FilterEntryFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Filter Entry %s still exists", filter_entry.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFilterEntryExists(name string, filter_entry *models.FilterEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Filter Entry %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Filter Entry Dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		filter_entryFound := models.FilterEntryFromContainer(cont)
		if filter_entryFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Filter Entry %s not found", rs.Primary.ID)
		}
		*filter_entry = *filter_entryFound
		return nil
	}
}

func testAccCheckAciFilterEntryIdNotEqual(fe1, fe2 *models.FilterEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fe1.DistinguishedName == fe2.DistinguishedName {
			return fmt.Errorf("Filter Entry DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciFilterEntryIdEqual(fe1, fe2 *models.FilterEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fe1.DistinguishedName != fe2.DistinguishedName {
			return fmt.Errorf("Filter Entry DNs are no equal")
		}
		return nil
	}
}
