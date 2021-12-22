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

func TestAccAciVRF_Basic(t *testing.T) {
	var vrf_default models.VRF
	var vrf_updated models.VRF
	resourceName := "aci_vrf.test"
	rName := makeTestVariable(acctest.RandString(5))
	rOther := makeTestVariable(acctest.RandString(5))
	prOther := makeTestVariable(acctest.RandString(5))
	longrName := acctest.RandString(65)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccVRFWithoutTenant(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccVRFWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVRFConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFExists(resourceName, &vrf_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "bd_enforced_enable", "no"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "ip_data_plane_learning", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "knw_mcast_act", "permit"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "pc_enf_dir", "ingress"),
					resource.TestCheckResourceAttr(resourceName, "pc_enf_pref", "enforced"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_bgp_ctx_af_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_eigrp_ctx_af_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_ospf_ctx_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
				),
			},
			{
				Config: CreateAccVRFConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFExists(resourceName, &vrf_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "test_annotation"),
					resource.TestCheckResourceAttr(resourceName, "bd_enforced_enable", "yes"),
					resource.TestCheckResourceAttr(resourceName, "description", "test_desc"),
					resource.TestCheckResourceAttr(resourceName, "ip_data_plane_learning", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "knw_mcast_act", "deny"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_alias"),
					resource.TestCheckResourceAttr(resourceName, "pc_enf_dir", "egress"),
					resource.TestCheckResourceAttr(resourceName, "pc_enf_pref", "unenforced"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_mcast_to.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_bgp_ctx_af_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_eigrp_ctx_af_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_ospf_ctx_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAciVRFIdEqual(&vrf_default, &vrf_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccVRFRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccVRFConfigWithParentAndName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of ctx-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccVRFConfigWithParentAndName(rName, rOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFExists(resourceName, &vrf_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rOther),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					testAccCheckAcVRFIdNotEqual(&vrf_default, &vrf_updated),
				),
			},
			{
				Config: CreateAccVRFConfig(rName),
			},
			{
				Config: CreateAccVRFConfigWithParentAndName(prOther, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFExists(resourceName, &vrf_updated),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", prOther)),
					testAccCheckAcVRFIdNotEqual(&vrf_default, &vrf_updated),
				),
			},
		},
	})
}

func TestAccAciVRF_NegativeCases(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVRFConfig(rName),
			},
			{
				Config:      CreateAccVRFWithInValidTenantDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvCtx (.)+`),
			},
			{
				Config:      CreateAccVRFUpdatedAttr(rName, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVRFUpdatedAttr(rName, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVRFUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVRFUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccVRFUpdatedAttr(rName, "bd_enforced_enable", randomValue),
				ExpectError: regexp.MustCompile(`expected bd_enforced_enable to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccVRFUpdatedAttr(rName, "ip_data_plane_learning", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name ipDataPlaneLearning, class fvCtx (.)+`),
			},
			{
				Config:      CreateAccVRFUpdatedAttr(rName, "knw_mcast_act", randomValue),
				ExpectError: regexp.MustCompile(`expected knw_mcast_act to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccVRFUpdatedAttr(rName, "pc_enf_dir", randomValue),
				ExpectError: regexp.MustCompile(`expected pc_enf_dir to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccVRFUpdatedAttr(rName, "pc_enf_pref", randomValue),
				ExpectError: regexp.MustCompile(`expected pc_enf_pref to be one of (.)+, got (.)+`),
			},
			{
				Config: CreateAccVRFConfig(rName),
			},
		},
	})
}

func TestAccAciVRF_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVRFsConfig(rName),
			},
		},
	})
}

func TestAccAciVRF_RelationParameters(t *testing.T) {
	var vrf_default models.VRF
	var vrf_rel1 models.VRF
	var vrf_rel2 models.VRF
	resourceName := "aci_vrf.test"
	rName := makeTestVariable(acctest.RandString(5))
	relRes1 := makeTestVariable(acctest.RandString(5))
	relRes2 := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVRFConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFExists(resourceName, &vrf_default),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_bgp_ctx_af_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_eigrp_ctx_af_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_ospf_ctx_pol.#", "0"),
				),
			},
			{
				Config: CreateAccVRFRelationsIntial(rName, relRes1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFExists(resourceName, &vrf_rel1),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bgp_ctx_pol", fmt.Sprintf("uni/tn-%s/bgpCtxP-%s", rName, relRes1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_mcast_to.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rName, relRes1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_bgp_ctx_af_pol.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "relation_fv_rs_ctx_to_bgp_ctx_af_pol.*", map[string]string{
						"af":                     "ipv4-ucast",
						"tn_bgp_ctx_af_pol_name": fmt.Sprintf("uni/tn-%s/bgpCtxAfP-%s", rName, relRes1),
					}),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_eigrp_ctx_af_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_ep_ret", fmt.Sprintf("uni/tn-%s/epRPol-%s", rName, relRes1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_ext_route_tag_pol", fmt.Sprintf("uni/tn-%s/rttag-%s", rName, relRes1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_ospf_ctx_pol.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "relation_fv_rs_ctx_to_ospf_ctx_pol.*", map[string]string{
						"af":                   "ipv6-ucast",
						"tn_ospf_ctx_pol_name": fmt.Sprintf("uni/tn-%s/ospfCtxP-%s", rName, relRes1),
					}),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ospf_ctx_pol", fmt.Sprintf("uni/tn-%s/ospfCtxP-%s", rName, relRes1)),
					testAccCheckAciVRFIdEqual(&vrf_default, &vrf_rel1),
				),
			},
			{
				Config: CreateAccVRFRelationsFinal(rName, relRes1, relRes2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFExists(resourceName, &vrf_rel2),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bgp_ctx_pol", fmt.Sprintf("uni/tn-%s/bgpCtxP-%s", rName, relRes2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_mcast_to.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rName, relRes2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_bgp_ctx_af_pol.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "relation_fv_rs_ctx_to_bgp_ctx_af_pol.*", map[string]string{
						"af":                     "ipv4-ucast",
						"tn_bgp_ctx_af_pol_name": fmt.Sprintf("uni/tn-%s/bgpCtxAfP-%s", rName, relRes1),
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "relation_fv_rs_ctx_to_bgp_ctx_af_pol.*", map[string]string{
						"af":                     "ipv6-ucast",
						"tn_bgp_ctx_af_pol_name": fmt.Sprintf("uni/tn-%s/bgpCtxAfP-%s", rName, relRes2),
					}),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_eigrp_ctx_af_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_ep_ret", fmt.Sprintf("uni/tn-%s/epRPol-%s", rName, relRes2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_ext_route_tag_pol", fmt.Sprintf("uni/tn-%s/rttag-%s", rName, relRes2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_ospf_ctx_pol.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "relation_fv_rs_ctx_to_ospf_ctx_pol.*", map[string]string{
						"af":                   "ipv6-ucast",
						"tn_ospf_ctx_pol_name": fmt.Sprintf("uni/tn-%s/ospfCtxP-%s", rName, relRes1),
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "relation_fv_rs_ctx_to_ospf_ctx_pol.*", map[string]string{
						"af":                   "ipv4-ucast",
						"tn_ospf_ctx_pol_name": fmt.Sprintf("uni/tn-%s/ospfCtxP-%s", rName, relRes2),
					}),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ospf_ctx_pol", fmt.Sprintf("uni/tn-%s/ospfCtxP-%s", rName, relRes2)),
					testAccCheckAciVRFIdEqual(&vrf_default, &vrf_rel2),
				),
			},
			{
				Config: CreateAccVRFConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFExists(resourceName, &vrf_rel2),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bgp_ctx_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_mcast_to.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_bgp_ctx_af_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_eigrp_ctx_af_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_ep_ret", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_ext_route_tag_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ctx_to_ospf_ctx_pol.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_ospf_ctx_pol", ""),
				),
			},
		},
	})
}

func CreateAccVRFRelationsFinal(rName, relName1, relName2 string) string {
	fmt.Println("=== STEP  testing vrf creation with final relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_ospf_timers" "test1"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_ospf_timers" "test2"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_end_point_retention_policy" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s" 
	}

	resource "aci_bgp_timers" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s" 
	}

	resource "aci_monitoring_policy" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s" 
	}

	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s" 
	}

	resource "aci_bgp_address_family_context" "test1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_bgp_address_family_context" "test2" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		relation_fv_rs_ospf_ctx_pol = aci_ospf_timers.test2.id
		relation_fv_rs_ctx_to_ospf_ctx_pol{
			tn_ospf_ctx_pol_name = aci_ospf_timers.test1.id
      		af = "ipv6-ucast" 
		}
		relation_fv_rs_ctx_to_ospf_ctx_pol{
			tn_ospf_ctx_pol_name = aci_ospf_timers.test2.id
      		af = "ipv4-ucast" 
		}
		relation_fv_rs_ctx_to_ep_ret = aci_end_point_retention_policy.test.id
		relation_fv_rs_bgp_ctx_pol = aci_bgp_timers.test.id
		relation_fv_rs_ctx_mon_pol = aci_monitoring_policy.test.id
		relation_fv_rs_ctx_to_ext_route_tag_pol = aci_l3out_route_tag_policy.test.id
		relation_fv_rs_ctx_to_bgp_ctx_af_pol{
			tn_bgp_ctx_af_pol_name = aci_bgp_address_family_context.test1.id
			af = "ipv4-ucast"
		}
		relation_fv_rs_ctx_to_bgp_ctx_af_pol{
			tn_bgp_ctx_af_pol_name = aci_bgp_address_family_context.test2.id
			af = "ipv6-ucast"
		}
	}
	`, rName, relName1, relName2, relName2, relName2, relName2, relName2, relName1, relName2, rName)
	return resource
}

func CreateAccVRFRelationsIntial(rName, relName string) string {
	fmt.Println("=== STEP  testing vrf creation with initial relational parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_ospf_timers" "test1"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_end_point_retention_policy" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s" 
	}

	resource "aci_bgp_timers" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s" 
	}

	resource "aci_monitoring_policy" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s" 
	}

	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s" 
	}

	resource "aci_bgp_address_family_context" "test1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		relation_fv_rs_ospf_ctx_pol = aci_ospf_timers.test1.id
		relation_fv_rs_ctx_to_ospf_ctx_pol{
			tn_ospf_ctx_pol_name = aci_ospf_timers.test1.id
      		af = "ipv6-ucast" 
		}
		relation_fv_rs_ctx_to_ep_ret = aci_end_point_retention_policy.test.id
		relation_fv_rs_bgp_ctx_pol = aci_bgp_timers.test.id
		relation_fv_rs_ctx_mon_pol = aci_monitoring_policy.test.id
		relation_fv_rs_ctx_to_ext_route_tag_pol = aci_l3out_route_tag_policy.test.id
		relation_fv_rs_ctx_to_bgp_ctx_af_pol{
			tn_bgp_ctx_af_pol_name = aci_bgp_address_family_context.test1.id
			af = "ipv4-ucast"
		}
	}
	`, rName, relName, relName, relName, relName, relName, relName, rName)
	return resource
}

func CreateAccVRFsConfig(rName string) string {
	fmt.Println("=== STEP  creating multiple vrf")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_vrf" "test1"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_vrf" "test2"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_vrf" "test3"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	`, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccVRFUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vrf attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccVRFWithInValidTenantDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing vrf creation with invalid tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_vrf" "test"{
		tenant_dn = aci_application_profile.test.id
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func testAccCheckAcVRFIdNotEqual(vrf1, vrf2 *models.VRF) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vrf1.DistinguishedName == vrf2.DistinguishedName {
			return fmt.Errorf("VRF DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciVRFIdEqual(vrf1, vrf2 *models.VRF) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vrf1.DistinguishedName != vrf2.DistinguishedName {
			return fmt.Errorf("VRF DNs are not equal")
		}
		return nil
	}
}

func CreateAccVRFConfigWithParentAndName(prName, rName string) string {
	fmt.Printf("=== STEP  Basic: testing vrf creation with tenant name %s name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, prName, rName)
	return resource
}

func CreateAccVRFRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vrf updation without required fields")
	resource := fmt.Sprintln(`
	resource "aci_vrf" "test" {
		annotation = "tag"
		bd_enforced_enable = "yes"
		description = "test_desc"
		ip_data_plane_learning = "disabled"
		knw_mcast_act = "deny"
		name_alias = "test_alias"
		pc_enf_dir = "egress"
		pc_enf_pref = "unenforced"
	}
	`)
	return resource
}

func CreateAccVRFConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  testing vrf creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
		annotation = "test_annotation"
		bd_enforced_enable = "yes"
		description = "test_desc"
		ip_data_plane_learning = "disabled"
		knw_mcast_act = "deny"
		name_alias = "test_alias"
		pc_enf_dir = "egress"
		pc_enf_pref = "unenforced"
	}
	`, rName, rName)
	return resource
}

func CreateAccVRFConfig(rName string) string {
	fmt.Println("=== STEP  testing vrf creation with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccVRFWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing vrf creation without name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
	}
	`, rName)
	return resource
}

func CreateAccVRFWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing vrf creation without creating tenant")
	resource := fmt.Sprintf(`
	resource "aci_vrf" "test" {
		name = "%s"
	}
	`, rName)
	return resource
}

func testAccCheckAciVRFExists(name string, vrf *models.VRF) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VRF %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VRF dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vrfFound := models.VRFFromContainer(cont)
		if vrfFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VRF %s not found", rs.Primary.ID)
		}
		*vrf = *vrfFound
		return nil
	}
}

func testAccCheckAciVRFDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vrf destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_vrf" {
			cont, err := client.Get(rs.Primary.ID)
			vrf := models.VRFFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VRF %s Still exists", vrf.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}
