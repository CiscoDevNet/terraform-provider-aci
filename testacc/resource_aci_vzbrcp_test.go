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

func TestAccAciContract_Basic(t *testing.T) {
	var contract_default models.Contract
	var contract_updated models.Contract
	resourceName := "aci_contract.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	prOther := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciContractDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateAccContractWithoutTenant(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccContractWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccContractConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_graph_att", ""),
					resource.TestCheckResourceAttr(resourceName, "scope", "context"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
				),
			},
			{
				Config: CreateAccContractConfigOptionalWithoutFilterParameter(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "test_annotation"),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "filter.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_alias"),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_graph_att", ""),
					resource.TestCheckResourceAttr(resourceName, "scope", "tenant"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS0"),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccContractConfigOptional(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "test_annotation"),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "filter.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.annotation", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_name", rName),
					resource.TestCheckResourceAttr(resourceName, "filter.0.name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.apply_to_frag", "no"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.arp_opc", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_from_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_to_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.entry_annotation", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.entry_description", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.entry_name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.ether_t", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.filter_entry_name", rName),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv4_t", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv6_t", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.match_dscp", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_from_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_to_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.stateful", "no"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.tcp_rules", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_alias"),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					resource.TestCheckResourceAttr(resourceName, "relation_vz_rs_graph_att", ""),
					resource.TestCheckResourceAttr(resourceName, "scope", "tenant"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS0"),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config:      CreateAccContractRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccContractWithoutRequiredFieldFilter(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccContractWithoutRequiredFieldFilterEntry(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccContractConfigWithFilterResourcesOptional(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.annotation", "filter_annotation"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.description", "filter_description"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_name", rName),
					resource.TestCheckResourceAttr(resourceName, "filter.0.name_alias", "filter_name_alias"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.apply_to_frag", "no"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.arp_opc", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_from_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_to_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.entry_annotation", "entry_annotation"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.entry_description", "entry_description"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.entry_name_alias", "entry_name_alias"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.ether_t", "ip"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.filter_entry_name", rName),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv4_t", "echo-rep"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv6_t", "dst-unreach"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.match_dscp", "CS0"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_from_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_to_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.stateful", "yes"),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.tcp_rules", "est"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config:      CreateAccContractConfigWithParentAndName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of brc-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccContractConfigWithParentAndName(rName, rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rOther),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciContrctIdNotEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractConfig(rName),
			},
			{
				Config: CreateAccContractConfigWithParentAndName(prOther, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", prOther)),
					testAccCheckAciContrctIdNotEqual(&contract_default, &contract_updated),
				),
			},
		},
	})
}

func TestAccAciContract_Update(t *testing.T) {
	var contract_default models.Contract
	var contract_updated models.Contract
	resourceName := "aci_contract.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccContractConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_default),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "scope", "global"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "scope", "global"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "scope", "application-profile"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "scope", "application-profile"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "CS1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS1"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF11"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF11"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "CS2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS2"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "AF21"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF21"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "VA"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "VA"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttr(rName, "target_dscp", "EF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "EF"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "ether_t", "ip"),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "apply_to_frag", "yes"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.apply_to_frag", "yes"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "apply_to_frag", "no"),
			},
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntry(rName, "ether_t", "arp"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.ether_t", "arp"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntry(rName, "arp_opc", "req"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.arp_opc", "req"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntry(rName, "arp_opc", "reply"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.arp_opc", "reply"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "arp_opc", "unspecified"),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "ether_t", "ipv4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.ether_t", "ipv4"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntry(rName, "ether_t", "trill"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.ether_t", "trill"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntry(rName, "ether_t", "mpls_ucast"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.ether_t", "mpls_ucast"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntry(rName, "ether_t", "mac_security"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.ether_t", "mac_security"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntry(rName, "ether_t", "fcoe"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.ether_t", "fcoe"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "ether_t", "ipv6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.ether_t", "ipv6"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "icmpv4_t", "dst-unreach"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv4_t", "dst-unreach"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "icmpv4_t", "echo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv4_t", "echo"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "icmpv4_t", "time-exceeded"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv4_t", "time-exceeded"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "icmpv4_t", "src-quench"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv4_t", "src-quench"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "icmpv6_t", "time-exceeded"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv6_t", "time-exceeded"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "icmpv6_t", "echo-req"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv6_t", "echo-req"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "icmpv6_t", "echo-rep"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv6_t", "echo-rep"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "icmpv6_t", "nbr-solicit"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv6_t", "nbr-solicit"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "icmpv6_t", "nbr-advert"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv6_t", "nbr-advert"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "icmpv6_t", "redirect"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.icmpv6_t", "redirect"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "match_dscp", "CS1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.match_dscp", "CS1"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "match_dscp", "AF11"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.match_dscp", "AF11"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "match_dscp", "CS2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.match_dscp", "CS2"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "match_dscp", "AF21"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.match_dscp", "AF21"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "match_dscp", "VA"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.match_dscp", "VA"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "match_dscp", "EF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.match_dscp", "EF"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "icmp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "icmp"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "igmp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "igmp"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "egp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "egp"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "igp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "igp"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "icmpv6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "icmpv6"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "eigrp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "eigrp"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "ospfigp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "ospfigp"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "pim"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "pim"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "l2tp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "l2tp"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "udp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.prot", "udp"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "tcp"),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "tcp_rules", "syn"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.tcp_rules", "syn"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "tcp_rules", "ack"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.tcp_rules", "ack"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "tcp_rules", "fin"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.tcp_rules", "fin"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "tcp_rules", "rst"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciContractExists(resourceName, &contract_updated),
					resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.tcp_rules", "rst"),
					testAccCheckAciContractdEqual(&contract_default, &contract_updated),
				),
			},
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntryForPortAttr(rName, "ftpData"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_from_port", "20"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_to_port", "20"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_from_port", "20"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_to_port", "20"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntryForPortAttr(rName, "smtp"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_from_port", "25"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_to_port", "25"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_from_port", "25"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_to_port", "25"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntryForPortAttr(rName, "dns"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_from_port", "53"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_to_port", "53"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_from_port", "53"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_to_port", "53"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntryForPortAttr(rName, "http"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_from_port", "80"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_to_port", "80"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_from_port", "80"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_to_port", "80"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntryForPortAttr(rName, "pop3"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_from_port", "110"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_to_port", "110"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_from_port", "110"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_to_port", "110"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
			// {
			// 	Config: CreateAccContractUpdatedAttrFilterEntryForPortAttr(rName, "rtsp"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciContractExists(resourceName, &contract_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_from_port", "554"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.d_to_port", "554"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_from_port", "554"),
			// 		resource.TestCheckResourceAttr(resourceName, "filter.0.filter_entry.0.s_to_port", "554"),
			// 		testAccCheckAciContractdEqual(&contract_default, &contract_updated),
			// 	),
			// },
		},
	})
}

func TestAccAciContract_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longAnnotationDesc := acctest.RandString(129)
	longNameAlias := acctest.RandString(65)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	longrName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccContractConfig(rName),
			},
			{
				Config:      CreateAccContractWithInValidTenantDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class vzBrCP (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, "description", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, "annotation", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, "prio", randomValue),
				ExpectError: regexp.MustCompile(`expected prio to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, "target_dscp", randomValue),
				ExpectError: regexp.MustCompile(`expected target_dscp to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, "scope", randomValue),
				ExpectError: regexp.MustCompile(`expected scope to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccContractUpdatedFilterAttr(rName, "description", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedFilterAttr(rName, "annotation", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedFilterAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedFilterAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},
			{
				Config:      CreateAccContractFilterWithInvalidName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of flt-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config:      CreateAccContractFilterEntryWithInvalidName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of e-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "entry_annotation", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "entry_description", longAnnotationDesc),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "entry_name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "apply_to_frag", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "arp_opc", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "ether_t", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "icmpv4_t", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "icmpv6_t", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "match_dscp", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "prot", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "stateful", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "tcp_rules", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`Unsupported argument`),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "ether_t", "ip"),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "tcp"),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "d_from_port", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "d_to_port", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "s_from_port", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "s_to_port", "65536"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "arp_opc", "req"),
				ExpectError: regexp.MustCompile(`ArpOpcode cannot be combined with non arp Ethertype`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "arp_opc", "reply"),
				ExpectError: regexp.MustCompile(`ArpOpcode cannot be combined with non arp Ethertype`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "apply_to_frag", "yes"),
				ExpectError: regexp.MustCompile(`non-IP Ethertype cannot be combined with other l4 properties`),
			},
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "prot", "unspecified"),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "d_from_port", "http"),
				ExpectError: regexp.MustCompile(`non-IP Ethertype cannot be combined with other l4 properties`),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "d_to_port", "http"),
				ExpectError: regexp.MustCompile(`non-IP Ethertype cannot be combined with other l4 properties`),
			},
			// {
			// 	Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "s_from_port", "http"),
			// 	ExpectError: regexp.MustCompile(`non-IP Ethertype cannot be combined with other l4 properties`),
			// },
			// {
			// 	Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "s_to_port", "http"),
			// 	ExpectError: regexp.MustCompile(`non-IP Ethertype cannot be combined with other l4 properties`),
			// },
			{
				Config: CreateAccContractUpdatedAttrFilterEntry(rName, "ether_t", "unspecified"),
			},
			{
				Config:      CreateAccContractUpdatedAttrFilterEntry(rName, "tcp_rules", "est"),
				ExpectError: regexp.MustCompile(`non-IP Ethertype cannot be combined with other l4 properties`),
			},
			{
				Config: CreateAccContractConfig(rName),
			},
		},
	})
}

func TestAccAciContract_MultipleCreateDestroy(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccContractConfigMultiple(rName),
			},
		},
	})
}

func CreateAccContractFilterEntryWithInvalidName(rName, longrName string) string {
	fmt.Printf("=== STEP  testing contract's filter_entry creation with name = %s\n", longrName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		filter{
			filter_name = "%s"
			filter_entry{
				filter_entry_name = "%s"
			}
		}
	}

	`, rName, rName, rName, longrName)
	return resource
}

func CreateAccContractFilterWithInvalidName(rName, longrName string) string {
	fmt.Printf("=== STEP  testing contract's filter creation with name = %s\n", longrName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		filter{
			filter_name = "%s"
		}
	}

	`, rName, rName, longrName)
	return resource
}

func CreateAccContractUpdatedFilterAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing contract's filter updation with %s = %s\n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		filter{
			filter_name = "%s"
			%s = "%s"
		}
	}

	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccContractWithInValidTenantDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing contract creation with invalid tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_application_profile.test.id
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccContractUpdatedAttrFilterEntryForPortAttr(rName, value string) string {
	fmt.Printf("=== STEP  testing contract by updating filter_entry ports %s\n", value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		filter {
			filter_name = "%s"
			filter_entry {
				filter_entry_name = "%s"
				prot = "tcp"
				d_from_port = "%s"
				d_to_port = "%s"
				s_from_port = "%s"
				s_to_port = "%s"
			}
		}
	}

	`, rName, rName, rName, rName, value, value, value, value)
	return resource
}

func CreateAccContractUpdatedAttrFilterEntry(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing contract by updating filter_entry's attribute with %s = %s\n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		filter {
			filter_name = "%s"
			filter_entry {
				filter_entry_name = "%s"
				%s = "%s"
			}
		}
	}

	`, rName, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccContractUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing contract updation with %s = %s\n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}

	`, rName, rName, attribute, value)
	return resource
}

func CreateAccContractConfigWithFilterResourcesOptional(rName string) string {
	fmt.Println("=== STEP  testing contract creation with optional parameters of filter resources")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		filter {
		  filter_name = "%s"
		  annotation = "filter_annotation"
		  description = "filter_description"
		  name_alias = "filter_name_alias"
		  filter_entry {
			filter_entry_name = "%s"
			apply_to_frag = "no"
			arp_opc = "unspecified"
			// d_from_port = "https"
			// d_to_port = "https"
			entry_annotation = "entry_annotation"
			entry_description = "entry_description"
			entry_name_alias = "entry_name_alias"
			ether_t = "ip"
			icmpv4_t = "echo-rep"
			icmpv6_t = "dst-unreach"
			match_dscp = "CS0"
			prot = "tcp"
			// s_from_port = "https"
			// s_to_port = "https"
			stateful = "yes"
			tcp_rules = "est"
		  }
		}
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccContractConfigWithParentAndName(prName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing contract creation with tenant name %s name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, prName, rName)
	return resource
}

func CreateAccContractRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing contract updation without required fields")
	resource := fmt.Sprintln(`
	resource "aci_contract" "test" {
		annotation = "tag"
		description = "test_description"
		name_alias = "test_alias"
		prio = "level1"
		scope = "tenant"
		target_dscp = "CS0"
	}
	`)
	return resource
}

func CreateAccContractConfigWithFilterResources(rName string) string {
	fmt.Println("=== STEP  testing contract creation with filter resources")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		filter {
			filter_name = "%s"
			filter_entry {
			  filter_entry_name = "%s"
			}
		}
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccContractConfigOptionalWithoutFilterParameter(rName string) string {
	fmt.Println("=== STEP  testing contract creation with optional parameters without filter parameter")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		annotation = "test_annotation"
		description = "test_description"
		name_alias = "test_alias"
		prio = "level1"
		scope = "tenant"
		target_dscp = "CS0"
	}
	`, rName, rName)
	return resource
}

func CreateAccContractConfigOptional(rName string) string {
	fmt.Println("=== STEP  testing contract creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		annotation = "test_annotation"
		description = "test_description"
		name_alias = "test_alias"
		prio = "level1"
		scope = "tenant"
		target_dscp = "CS0"
		filter {
			filter_name = "%s"
			filter_entry {
			  filter_entry_name = "%s"
			}
		}
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccContractConfig(rName string) string {
	fmt.Println("=== STEP  testing contract creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccContractConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple contract creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract" "test2" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract" "test3" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccContractWithoutRequiredFieldFilterEntry(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract's filter_entry without required parameter")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		filter{
			filter_name = "%s"
			filter_entry{

			}
		}
	}
	`, rName, rName)
	return resource
}

func CreateAccContractWithoutRequiredFieldFilter(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract's filter without required parameter")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
		filter{}
	}
	`, rName)
	return resource
}

func CreateAccContractWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract creation without name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_contract" "test" {
		tenant_dn = aci_tenant.test.id
	}
	`, rName)
	return resource
}

func CreateAccContractWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract creation without creating tenant")
	resource := fmt.Sprintf(`
	resource "aci_contract" "test" {
		name = "%s"
	}
	`, rName)
	return resource
}

func testAccCheckAciContractExists(name string, contract *models.Contract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Contract %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Contract dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		contractFound := models.ContractFromContainer(cont)
		if contractFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Contract %s not found", rs.Primary.ID)
		}
		*contract = *contractFound
		return nil
	}
}

func testAccCheckAciContractDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing contract destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_contract" {
			cont, err := client.Get(rs.Primary.ID)
			contract := models.ContractFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Contract %s Still exists", contract.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciContrctIdNotEqual(c1, c2 *models.Contract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if c1.DistinguishedName == c2.DistinguishedName {
			return fmt.Errorf("Contract DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciContractdEqual(c1, c2 *models.Contract) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if c1.DistinguishedName != c2.DistinguishedName {
			return fmt.Errorf("Contract DNs are not equal")
		}
		return nil
	}
}
