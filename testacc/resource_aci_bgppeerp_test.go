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

func TestAccAciPeerConnectivityProfile_Basic(t *testing.T) {
	var peer_connectivity_profile_default models.BgpPeerConnectivityProfile
	var peer_connectivity_profile_updated models.BgpPeerConnectivityProfile
	resourceName := "aci_bgp_peer_connectivity_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	addr, _ := acctest.RandIpAddress("10.0.0.0/16")
	addr = fmt.Sprintf("%s/16", addr)
	addrUpdated, _ := acctest.RandIpAddress("10.0.0.0/17")
	addrUpdated = fmt.Sprintf("%s/16", addrUpdated)
	fvTenantName := makeTestVariable(acctest.RandString(5))
	l3extOutName := makeTestVariable(acctest.RandString(5))
	l3extLNodePName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPeerConnectivityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreatePeerConnectivityProfileWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, addr, "parent_dn"),
				ExpectError: regexp.MustCompile(`parent_dn is required`),
			},
			{
				Config:      CreatePeerConnectivityProfileWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, addr, "addr"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccPeerConnectivityProfileConfig(fvTenantName, l3extOutName, l3extLNodePName, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_default),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", fvTenantName, l3extOutName, l3extLNodePName)),
					resource.TestCheckResourceAttr(resourceName, "addr", addr),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "addr_t_ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "addr_t_ctrl.0", "af-ucast"),
					resource.TestCheckResourceAttr(resourceName, "admin_state", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "allowed_self_as_cnt", "3"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "1"),
					resource.TestCheckResourceAttr(resourceName, "weight", "0"),
					resource.TestCheckResourceAttr(resourceName, "as_number", ""),
					resource.TestCheckResourceAttr(resourceName, "local_asn", ""),
					resource.TestCheckResourceAttr(resourceName, "local_asn_propagate", ""),
				),
			},
			{
				Config: CreateAccPeerConnectivityProfileConfigWithOptionalValues(fvTenantName, l3extOutName, l3extLNodePName, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", fvTenantName, l3extOutName, l3extLNodePName)),
					resource.TestCheckResourceAttr(resourceName, "addr", addr),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_peer_connectivity_profile"),
					resource.TestCheckResourceAttr(resourceName, "addr_t_ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "addr_t_ctrl.0", "af-mcast"),
					resource.TestCheckResourceAttr(resourceName, "admin_state", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "allowed_self_as_cnt", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "allow-self-as"),
					resource.TestCheckResourceAttr(resourceName, "password", "cisco"),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.0", "bfd"),
					resource.TestCheckResourceAttr(resourceName, "private_a_sctrl.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "private_a_sctrl.0", "remove-all"),
					resource.TestCheckResourceAttr(resourceName, "private_a_sctrl.1", "remove-exclusive"),
					resource.TestCheckResourceAttr(resourceName, "private_a_sctrl.2", "replace-as"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "2"),
					resource.TestCheckResourceAttr(resourceName, "weight", "1"),
					resource.TestCheckResourceAttr(resourceName, "as_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "local_asn", "2"),
					resource.TestCheckResourceAttr(resourceName, "local_asn_propagate", "dual-as"),

					testAccCheckAciPeerConnectivityProfileIdEqual(&peer_connectivity_profile_default, &peer_connectivity_profile_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config:      CreateAccPeerConnectivityProfileWithInvalidIP(rName, rName, rName, addr),
				ExpectError: regexp.MustCompile(`unknown property value (.)+`),
			},

			{
				Config:      CreateAccPeerConnectivityProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccPeerConnectivityProfileConfigWithRequiredParams(rName, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "addr", addr),
					testAccCheckAciPeerConnectivityProfileIdNotEqual(&peer_connectivity_profile_default, &peer_connectivity_profile_updated),
				),
			},
			{
				Config: CreateAccPeerConnectivityProfileConfig(fvTenantName, l3extOutName, l3extLNodePName, addr),
			},
			{
				Config: CreateAccPeerConnectivityProfileConfigWithRequiredParams(rNameUpdated, addrUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "parent_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", rNameUpdated, rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "addr", addrUpdated),
					testAccCheckAciPeerConnectivityProfileIdNotEqual(&peer_connectivity_profile_default, &peer_connectivity_profile_updated),
				),
			},
		},
	})
}

func TestAccAciPeerConnectivityProfile_Update(t *testing.T) {
	var peer_connectivity_profile_default models.BgpPeerConnectivityProfile
	var peer_connectivity_profile_updated models.BgpPeerConnectivityProfile
	resourceName := "aci_bgp_peer_connectivity_profile.test"
	// rName := makeTestVariable(acctest.RandString(5))

	addr, _ := acctest.RandIpAddress("10.0.0.0/16")
	addr = fmt.Sprintf("%s/16", addr)
	addrUpdated, _ := acctest.RandIpAddress("10.0.0.0/16")
	addrUpdated = fmt.Sprintf("%s/16", addrUpdated)
	fvTenantName := makeTestVariable(acctest.RandString(5))
	l3extOutName := makeTestVariable(acctest.RandString(5))
	l3extLNodePName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPeerConnectivityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPeerConnectivityProfileConfig(fvTenantName, l3extOutName, l3extLNodePName, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_default),
				),
			},
			{

				Config: CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "addr_t_ctrl", StringListtoString([]string{"af-mcast"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "addr_t_ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "addr_t_ctrl.0", "af-mcast"),
				),
			},
			{

				Config: CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "addr_t_ctrl", StringListtoString([]string{"af-mcast", "af-ucast"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "addr_t_ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "addr_t_ctrl.0", "af-mcast"),
					resource.TestCheckResourceAttr(resourceName, "addr_t_ctrl.1", "af-ucast"),
				),
			},
			{

				Config: CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "addr_t_ctrl", StringListtoString([]string{"af-ucast"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "addr_t_ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "addr_t_ctrl.0", "af-ucast"),
				),
			},
			{

				Config: CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "ctrl", StringListtoString([]string{"allow-self-as"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "allow-self-as"),
				),
			},
			{
				Config: CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "ctrl", StringListtoString([]string{"send-ext-com", "send-com", "nh-self", "dis-peer-as-check", "as-override", "allow-self-as"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "send-ext-com"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "send-com"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.2", "nh-self"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.3", "dis-peer-as-check"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.4", "as-override"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.5", "allow-self-as"),
				),
			},
			{

				Config: CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "peer_ctrl", StringListtoString([]string{"bfd"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.0", "bfd"),
				),
			},
			{

				Config: CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "peer_ctrl", StringListtoString([]string{"bfd", "dis-conn-check"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.0", "bfd"),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.1", "dis-conn-check"),
				),
			},
			{

				Config: CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "peer_ctrl", StringListtoString([]string{"dis-conn-check"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.0", "dis-conn-check"),
				),
			},
			{
				Config: CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "peer_ctrl", StringListtoString([]string{"dis-conn-check", "bfd"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.0", "dis-conn-check"),
					resource.TestCheckResourceAttr(resourceName, "peer_ctrl.1", "bfd"),
				),
			},
			{

				Config: CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "private_a_sctrl", StringListtoString([]string{"remove-all", "remove-exclusive", "replace-as"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "private_a_sctrl.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "private_a_sctrl.0", "remove-all"),
					resource.TestCheckResourceAttr(resourceName, "private_a_sctrl.1", "remove-exclusive"),
					resource.TestCheckResourceAttr(resourceName, "private_a_sctrl.2", "replace-as"),
				),
			},
			{
				Config: CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "private_a_sctrl", StringListtoString([]string{"replace-as", "remove-exclusive", "remove-all"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "private_a_sctrl.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "private_a_sctrl.0", "replace-as"),
					resource.TestCheckResourceAttr(resourceName, "private_a_sctrl.1", "remove-exclusive"),
					resource.TestCheckResourceAttr(resourceName, "private_a_sctrl.2", "remove-all"),
				),
			},
			{
				Config: CreateAccPeerConnectivityProfileUpdatedLocalAsnAttrs(fvTenantName, l3extOutName, l3extLNodePName, addr, "local_asn_propagate", "no-prepend"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "local_asn_propagate", "no-prepend"),
				),
			},
			{
				Config: CreateAccPeerConnectivityProfileUpdatedLocalAsnAttrs(fvTenantName, l3extOutName, l3extLNodePName, addr, "local_asn_propagate", "replace-as"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPeerConnectivityProfileExists(resourceName, &peer_connectivity_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "local_asn_propagate", "replace-as"),
				),
			},
			{
				Config: CreateAccPeerConnectivityProfileConfig(fvTenantName, l3extOutName, l3extLNodePName, addr),
			},
		},
	})
}

func TestAccAciPeerConnectivityProfile_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	addr, _ := acctest.RandIpAddress("10.0.0.0/16")
	addr = fmt.Sprintf("%s/16", addr)
	addrUpdated, _ := acctest.RandIpAddress("10.0.0.0/16")
	addrUpdated = fmt.Sprintf("%s/16", addrUpdated)
	fvTenantName := makeTestVariable(acctest.RandString(5))
	l3extOutName := makeTestVariable(acctest.RandString(5))
	l3extLNodePName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPeerConnectivityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPeerConnectivityProfileConfig(fvTenantName, l3extOutName, l3extLNodePName, addr),
			},
			{
				Config:      CreateAccPeerConnectivityProfileWithInValidParentDn(rName, addr),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, addr, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, addr, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, addr, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "addr_t_ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "addr_t_ctrl", StringListtoString([]string{"af-mcast", "af-mcast"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, addr, "admin_state", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, addr, "allowed_self_as_cnt", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "ctrl", StringListtoString([]string{"allow-self-as", "allow-self-as"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "peer_ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "peer_ctrl", StringListtoString([]string{"bfd", "bfd"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "private_a_sctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, "private_a_sctrl", StringListtoString([]string{"remove-all", "remove-all"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedLocalAsnAttrs(fvTenantName, l3extOutName, l3extLNodePName, addr, "local_asn_propagate", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, addr, "ttl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, addr, "weight", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccPeerConnectivityProfileUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, addr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccPeerConnectivityProfileConfig(fvTenantName, l3extOutName, l3extLNodePName, addr),
			},
		},
	})
}

func TestAccAciPeerConnectivityProfile_MultipleCreateDelete(t *testing.T) {

	addr, _ := acctest.RandIpAddress("10.0.0.0/16")
	fvTenantName := makeTestVariable(acctest.RandString(5))
	l3extOutName := makeTestVariable(acctest.RandString(5))
	l3extLNodePName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPeerConnectivityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPeerConnectivityProfileConfigMultiple(fvTenantName, l3extOutName, l3extLNodePName, addr[:(len(addr)-1)]),
			},
		},
	})
}

func testAccCheckAciPeerConnectivityProfileExists(name string, peer_connectivity_profile *models.BgpPeerConnectivityProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Peer Connectivity Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Peer Connectivity Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		peer_connectivity_profileFound := models.BgpPeerConnectivityProfileFromContainer(cont)
		if peer_connectivity_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Peer Connectivity Profile %s not found", rs.Primary.ID)
		}
		*peer_connectivity_profile = *peer_connectivity_profileFound
		return nil
	}
}

func testAccCheckAciPeerConnectivityProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing peer_connectivity_profile destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bgp_peer_connectivity_profile" {
			cont, err := client.Get(rs.Primary.ID)
			peer_connectivity_profile := models.BgpPeerConnectivityProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Peer Connectivity Profile %s Still exists", peer_connectivity_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciPeerConnectivityProfileIdEqual(m1, m2 *models.BgpPeerConnectivityProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("peer_connectivity_profile DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciPeerConnectivityProfileIdNotEqual(m1, m2 *models.BgpPeerConnectivityProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("peer_connectivity_profile DNs are equal")
		}
		return nil
	}
}

func CreateAccPeerConnectivityProfileConfigMultiple(fvTenantName, l3extOutName, l3extLNodePName, addr string) string {
	fmt.Println("=== STEP  testing multiple peer_connectivity_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_bgp_peer_connectivity_profile" "test" {
		parent_dn  = aci_logical_node_profile.test.id
		addr  = "%s${count.index}/16"
		count = 5
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, addr)
	return resource
}

func CreatePeerConnectivityProfileWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, addr, attrName string) string {
	fmt.Println("=== STEP  Basic: testing peer_connectivity_profile creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	`
	switch attrName {
	case "parent_dn":
		rBlock += `
	resource "aci_bgp_peer_connectivity_profile" "test" {
	#	parent_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
	}
		`
	case "addr":
		rBlock += `
	resource "aci_bgp_peer_connectivity_profile" "test" {
		parent_dn  = aci_logical_node_profile.test.id
	#	addr  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName, addr)
}

func CreateAccPeerConnectivityProfileConfigWithRequiredParams(prName, addr string) string {
	fmt.Printf("=== STEP  testing peer_connectivity_profile creation with parent resource name %s and address %s\n", prName, addr)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_bgp_peer_connectivity_profile" "test" {
		parent_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
	}
	`, prName, prName, prName, addr)
	return resource
}

func CreateAccPeerConnectivityProfileConfig(fvTenantName, l3extOutName, l3extLNodePName, addr string) string {
	fmt.Println("=== STEP  testing peer_connectivity_profile creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_bgp_peer_connectivity_profile" "test" {
		parent_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, addr)
	return resource
}

func CreateAccPeerConnectivityProfileWithInvalidIP(fvTenantName, l3extOutName, l3extLNodePName, addr string) string {
	fmt.Println("=== STEP  testing peer_connectivity_profile creation with invalid IP")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_bgp_peer_connectivity_profile" "test" {
		parent_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, l3extLNodePName)
	return resource
}

func CreateAccPeerConnectivityProfileWithInValidParentDn(rName, addr string) string {
	fmt.Println("=== STEP  Negative Case: testing peer_connectivity_profile creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_bgp_peer_connectivity_profile" "test" {
		parent_dn  = aci_tenant.test.id
		addr  = "%s"	
	}
	`, rName, addr)
	return resource
}

func CreateAccPeerConnectivityProfileConfigWithOptionalValues(fvTenantName, l3extOutName, l3extLNodePName, addr string) string {
	fmt.Println("=== STEP  Basic: testing peer_connectivity_profile creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_bgp_peer_connectivity_profile" "test" {
		parent_dn  = "${aci_logical_node_profile.test.id}"
		addr  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_peer_connectivity_profile"
		addr_t_ctrl = ["af-mcast"]
		admin_state = "disabled"
		allowed_self_as_cnt = "2"
		ctrl = ["allow-self-as"]
		password = "cisco"
		peer_ctrl = ["bfd"]
		private_a_sctrl = ["remove-all", "remove-exclusive", "replace-as"]
		ttl = "2"
		as_number = "1"
		local_asn = "2"
		local_asn_propagate = "dual-as"
		weight = "1"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, addr)

	return resource
}

func CreateAccPeerConnectivityProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing peer_connectivity_profile creation with optional parameters")
	resource := `
	resource "aci_bgp_peer_connectivity_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_peer_connectivity_profile"
		addr_t_ctrl = ["af-label-ucast"]
		admin_state = "disabled"
		allowed_self_as_cnt = "2"
		ctrl = ["allow-self-as"]
		password = "Cisco@123"
		peer_ctrl = ["bfd"]
		private_a_sctrl = ["remove-all", "remove-exclusive", "replace-as"]
		ttl = "2"
		as_number = "1"
		local_asn = "2"
		local_asn_propagate = "dual-as"
		weight = "1"
	}
	`

	return resource
}

func CreateAccPeerConnectivityProfileUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, addr, attribute, value string) string {
	fmt.Printf("=== STEP  testing peer_connectivity_profile attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_bgp_peer_connectivity_profile" "test" {
		parent_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, addr, attribute, value)
	return resource
}

func CreateAccPeerConnectivityProfileUpdatedLocalAsnAttrs(fvTenantName, l3extOutName, l3extLNodePName, addr, attribute, value string) string {
	fmt.Printf("=== STEP  testing peer_connectivity_profile attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_bgp_peer_connectivity_profile" "test" {
		parent_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
		local_asn = 3
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, addr, attribute, value)
	return resource
}

func CreateAccPeerConnectivityProfileUpdatedAttrList(fvTenantName, l3extOutName, l3extLNodePName, addr, attribute, value string) string {
	fmt.Printf("=== STEP  testing peer_connectivity_profile attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_bgp_peer_connectivity_profile" "test" {
		parent_dn  = aci_logical_node_profile.test.id
		addr  = "%s"
		%s = %s
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, addr, attribute, value)
	return resource
}
