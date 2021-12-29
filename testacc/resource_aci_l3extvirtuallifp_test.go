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

func TestAccAciL3outFloatingSVI_Basic(t *testing.T) {
	var l3out_floating_svi_default models.VirtualLogicalInterfaceProfile
	var l3out_floating_svi_updated models.VirtualLogicalInterfaceProfile
	resourceName := "aci_l3out_floating_svi.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	nodeDn := "topology/pod-1/node-111"
	nodeDnUpdated := "topology/pod-1/node-101"
	encap := "vlan-20"
	encapUpdated := "vlan-60"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outFloatingSVIDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outFloatingSVIWithoutRequired(rName, rName, rName, rName, nodeDn, encap, "logical_interface_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outFloatingSVIWithoutRequired(rName, rName, rName, rName, nodeDn, encap, "encap"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config:      CreateL3outFloatingSVIWithoutRequired(rName, rName, rName, rName, nodeDn, encap, "nodeDn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outFloatingSVIConfig(rName, rName, rName, rName, nodeDn, encap),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_default),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "node_dn", nodeDn),
					resource.TestCheckResourceAttr(resourceName, "encap", encap),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "addr", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "autostate", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "encap_scope", "local"),
					resource.TestCheckResourceAttr(resourceName, "if_inst_t", "ext-svi"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_dad", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "ll_addr", "::"),
					resource.TestCheckResourceAttr(resourceName, "mac", "00:22:BD:F8:19:FF"),
					resource.TestCheckResourceAttr(resourceName, "mode", "regular"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "inherit"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "unspecified"),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIConfigWithOptionalValues(rName, rName, rName, rName, nodeDn, encap),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "addr", "1.2.1.1/16"),
					resource.TestCheckResourceAttr(resourceName, "autostate", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "encap_scope", "ctx"),
					resource.TestCheckResourceAttr(resourceName, "if_inst_t", "ext-svi"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_dad", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "ll_addr", "fe80::1"),
					resource.TestCheckResourceAttr(resourceName, "mac", "00:22:BD:F8:19:F0"),
					resource.TestCheckResourceAttr(resourceName, "mode", "native"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "577"),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS0"),
					resource.TestCheckResourceAttr(resourceName, "relation_l3ext_rs_dyn_path_att.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccL3outFloatingSVIConfigWithUpdatedRequiredParams(rNameUpdated, nodeDn, encap),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rNameUpdated, rNameUpdated, rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "node_dn", nodeDn),
					resource.TestCheckResourceAttr(resourceName, "encap", encap),
					testAccCheckAciL3outFloatingSVIIdNotEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIConfig(rName, rName, rName, rName, nodeDn, encap),
			},
			{
				Config: CreateAccL3outFloatingSVIConfigWithUpdatedRequiredParams(rName, nodeDnUpdated, encap),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "node_dn", nodeDnUpdated),
					testAccCheckAciL3outFloatingSVIIdNotEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIConfig(rName, rName, rName, rName, nodeDn, encap),
			},
			{
				Config: CreateAccL3outFloatingSVIConfigWithUpdatedRequiredParams(rName, nodeDn, encapUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_interface_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/lifp-%s", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "node_dn", nodeDn),
					resource.TestCheckResourceAttr(resourceName, "encap", encapUpdated),
					testAccCheckAciL3outFloatingSVIIdNotEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config:      CreateAccL3outFloatingSVIConfigUpdateWithoutRequiredParameters(rName, rName, rName, rName, "description", "test_coverage"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outFloatingSVIConfig(rName, rName, rName, rName, nodeDn, encap),
			},
		},
	})
}

func TestAccAciL3outFloatingSVI_Update(t *testing.T) {
	var l3out_floating_svi_default models.VirtualLogicalInterfaceProfile
	var l3out_floating_svi_updated models.VirtualLogicalInterfaceProfile
	resourceName := "aci_l3out_floating_svi.test"
	rName := makeTestVariable(acctest.RandString(5))
	nodeDn := "topology/pod-1/node-111"
	encap := "vlan-20"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outFloatingSVIDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outFloatingSVIConfig(rName, rName, rName, rName, nodeDn, encap),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_default),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "mode", "untagged"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "mode", "untagged"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "CS1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS1"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "CS2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS2"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "CS3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS3"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "CS4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS4"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "CS5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS5"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "CS6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS6"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "CS7"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "CS7"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "EF"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "EF"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "VA"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "VA"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "AF11"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF11"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "AF12"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF12"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "AF13"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF13"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "AF21"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF21"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "AF22"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF22"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "AF23"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF23"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "AF31"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF31"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "AF32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF32"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "AF33"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF33"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "AF41"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF41"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "AF42"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF42"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", "AF43"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outFloatingSVIExists(resourceName, &l3out_floating_svi_updated),
					resource.TestCheckResourceAttr(resourceName, "target_dscp", "AF43"),
					testAccCheckAciL3outFloatingSVIIdEqual(&l3out_floating_svi_default, &l3out_floating_svi_updated),
				),
			},
			{
				Config: CreateAccL3outFloatingSVIConfig(rName, rName, rName, rName, nodeDn, encap),
			},
		},
	})
}

func TestAccAciL3outFloatingSVI_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	nodeDn := "topology/pod-1/node-111"
	encap := "vlan-20"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outFloatingSVIDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outFloatingSVIConfig(rName, rName, rName, rName, nodeDn, encap),
			},
			{
				Config:      CreateAccL3outFloatingSVIWithInValidParentDn(rName, rName, rName, rName, nodeDn, encap),
				ExpectError: regexp.MustCompile(`configured object (.)+ not found (.)+,`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "addr", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "addr", "1.1.1.1/32"),
				ExpectError: regexp.MustCompile(`Invalid External Interface Address`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "autostate", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "encap_scope", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "ipv6_dad", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "ll_addr", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "mtu", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "mtu", "100"),
				ExpectError: regexp.MustCompile(`Property mtu of (.)+ is out of range`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, "target_dscp", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)*to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccL3outFloatingSVIUpdatedAttr(rName, rName, rName, rName, nodeDn, encap, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)*is not expected here.`),
			},
			{
				Config: CreateAccL3outFloatingSVIConfig(rName, rName, rName, rName, nodeDn, encap),
			},
		},
	})
}

func TestAccAciL3outFloatingSVI_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	node_dn := "topology/pod-1/node-111"
	encap := "vlan-"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outFloatingSVIDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outFloatingSVIsConfig(rName, rName, rName, rName, node_dn, encap),
			},
		},
	})
}

func testAccCheckAciL3outFloatingSVIExists(name string, l3out_floating_svi *models.VirtualLogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Floating SVI %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Floating SVI dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_floating_sviFound := models.VirtualLogicalInterfaceProfileFromContainer(cont)
		if l3out_floating_sviFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Floating SVI %s not found", rs.Primary.ID)
		}
		*l3out_floating_svi = *l3out_floating_sviFound
		return nil
	}
}

func testAccCheckAciL3outFloatingSVIDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_floating_svi destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_floating_svi" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_floating_svi := models.VirtualLogicalInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Floating SVI %s Still exists", l3out_floating_svi.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outFloatingSVIIdEqual(m1, m2 *models.VirtualLogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_floating_svi DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outFloatingSVIIdNotEqual(m1, m2 *models.VirtualLogicalInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_floating_svi DNs are equal")
		}
		return nil
	}
}

func CreateL3outFloatingSVIWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, nodeDn, encap, attrName string) string {
	fmt.Printf("=== STEP  Basic: testing l3out_floating_svi creation without required parameter %s", attrName)
	rBlock := `
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
		
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}

	`
	switch attrName {
	case "logical_interface_profile_dn":
		rBlock += `
	resource "aci_l3out_floating_svi" "test" {
	#	logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"	
		encap  = "%s"
		description = "created while acceptance testing"
		if_inst_t = "ext-svi"
	}
	`
	case "nodeDn":
		rBlock += `
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	#	node_dn  = "%s"
		encap  = "%s"
		description = "created while acceptance testing"
		if_inst_t = "ext-svi"
	}
		`
	case "encap":
		rBlock += `
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
	#	encap  = "%s"
		description = "created while acceptance testing"
		if_inst_t = "ext-svi"
	}
		`
	}
	return rBlock
}

func CreateAccL3outFloatingSVIConfigWithUpdatedRequiredParams(rName, nodeDn, encap string) string {
	fmt.Println("=== STEP  testing l3out_floating_svi updation of required arguements")
	resource := fmt.Sprintf(
		`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}
	`, rName, rName, rName, rName, nodeDn, encap)
	return resource
}

func CreateAccL3outFloatingSVIConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, nodeDn, encap string) string {
	fmt.Println("=== STEP  testing l3out_floating_svi creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, nodeDn, encap)
	return resource
}

func CreateAccL3outFloatingSVIsConfig(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, nodeDn, encap string) string {
	fmt.Println("=== STEP  testing Multiple l3out_floating_svi creation with required arguements only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}

	resource "aci_l3out_floating_svi" "test1" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}

	resource "aci_l3out_floating_svi" "test2" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}

	resource "aci_l3out_floating_svi" "test3" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, nodeDn, encap+"10", nodeDn, encap+"20", nodeDn, encap+"30", nodeDn, encap+"40")
	return resource
}

func CreateAccL3outFloatingSVIWithInValidParentDn(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, nodeDn, encap string) string {
	fmt.Println("=== STEP  Negative Case: testing l3out_floating_svi creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = "${aci_logical_interface_profile.test.id}invalid"
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, nodeDn, encap)
	return resource
}

func CreateAccL3outFloatingSVIConfigWithOptionalValues(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, nodeDn, encap string) string {
	fmt.Println("=== STEP  Basic: testing l3out_floating_svi creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = "${aci_logical_interface_profile.test.id}"
		node_dn  = "%s"
		encap  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		addr = "1.2.1.1/16"
		autostate = "enabled"
		encap_scope = "ctx"
		if_inst_t = "ext-svi"
		ipv6_dad = "disabled"
		ll_addr = "fe80::1"
		mac = "00:22:BD:F8:19:F0"
		mode = "native"
		mtu = "577"
		target_dscp = "CS0"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, nodeDn, encap)
	return resource
}

func CreateAccL3outFloatingSVIUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, nodeDn, encap, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_floating_svi updation with attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		node_dn  = "%s"
		encap  = "%s"
		if_inst_t = "ext-svi"
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, nodeDn, encap, attribute, value)
	return resource
}

func CreateAccL3outFloatingSVIConfigUpdateWithoutRequiredParameters(fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_floating_svi updation with attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_floating_svi" "test" {
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLIfPName, attribute, value)
	return resource
}
