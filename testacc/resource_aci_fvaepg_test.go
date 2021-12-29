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

func TestAccAciApplicationEPG_Basic(t *testing.T) {
	var application_epg_default models.ApplicationEPG
	var application_epg_updated models.ApplicationEPG
	resourceName := "aci_application_epg.test"
	rName := acctest.RandString(5)
	rOtherName := acctest.RandString(5)
	parentOtherName := acctest.RandString(5)
	longrName := acctest.RandString(65)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciApplicationEPGDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccApplicationEPGWithoutApplicationProfile(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccApplicationEPGWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccApplicationEPGConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_default),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", ""),
					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "fwd_ctrl", "none"),
					resource.TestCheckResourceAttr(resourceName, "has_mcast_source", "no"),
					resource.TestCheckResourceAttr(resourceName, "is_attr_based_epg", "no"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtleastOne"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "pc_enf_pref", "unenforced"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "exclude"),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "shutdown", "no"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "application_profile_dn", fmt.Sprintf("uni/tn-%s/ap-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd", ""),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cust_qos_pol", "uni/tn-common/qoscustom-default"),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_fc_path_att", ""),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov", ""),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_graph_def", ""),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons_if", ""),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_sec_inherited", ""),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_node_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_dpp_pol", ""),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons", ""),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov_def", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_trust_ctrl", ""),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prot_by", ""),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_path_att", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_aepg_mon_pol", ""),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_intra_epg", ""),
				),
			},
			{
				Config: CreateAccApplicationEPGConfigWithOPtionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "description", "From Terraform"),
					resource.TestCheckResourceAttr(resourceName, "annotation", "tag"),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", "0"),
					resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "fwd_ctrl", "none"),
					resource.TestCheckResourceAttr(resourceName, "has_mcast_source", "no"),
					resource.TestCheckResourceAttr(resourceName, "is_attr_based_epg", "no"),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtleastOne"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "application_epg_alias"),
					resource.TestCheckResourceAttr(resourceName, "pc_enf_pref", "unenforced"),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "exclude"),
					resource.TestCheckResourceAttr(resourceName, "prio", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "shutdown", "no"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "application_profile_dn", fmt.Sprintf("uni/tn-%s/ap-%s", rName, rName)),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd", ""),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cust_qos_pol", "uni/tn-common/qoscustom-default"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_fc_path_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_graph_def.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons_if.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_sec_inherited.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_node_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_dpp_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov_def.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_trust_ctrl", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prot_by.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_path_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_aepg_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_intra_epg.#", "0"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccApplicationEPGConfigUpdateWithoutRequiredParameters(rName, "description", "test_coverage"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccApplicationEPGConfigUpdatedName(rName, longrName),
				ExpectError: regexp.MustCompile(fmt.Sprintf("property name of epg-%s failed validation for value '%s'", longrName, longrName)),
			},
			{
				Config: CreateAccApplicationEPGConfigWithParentAndName(rName, rOtherName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "application_profile_dn", fmt.Sprintf("uni/tn-%s/ap-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rOtherName),
					testAccCheckAciApplicationEPGIdNotEqual(&application_epg_default, &application_epg_updated),
				),
			},
			{
				Config: CreateAccApplicationEPGConfigWithParentAndName(parentOtherName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "application_profile_dn", fmt.Sprintf("uni/tn-%s/ap-%s", parentOtherName, parentOtherName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciApplicationEPGIdNotEqual(&application_epg_default, &application_epg_updated),
				),
			},
		},
	})
}

func TestAccAciApplicationEPG_Update(t *testing.T) {
	var application_epg_default models.ApplicationEPG
	var application_epg_updated models.ApplicationEPG
	resourceName := "aci_application_epg.test"
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciApplicationEPGDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccApplicationEPGConfig(rName),
				Check:  resource.ComposeTestCheckFunc(testAccCheckAciApplicationEPGExists(resourceName, &application_epg_default)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "prio", "level1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level1"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "prio", "level2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level2"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "prio", "level3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level3"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "prio", "level4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level4"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "prio", "level5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level5"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "prio", "level6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "prio", "level6"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "exception_tag", "20"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "exception_tag", "20"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			// {
			// 	Config: CreateAccApplicationEPGUpdatedAttr(rName, "flood_on_encap", "enabled"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "flood_on_encap", "enabled"),
			// 		testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			// },
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "fwd_ctrl", "proxy-arp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "fwd_ctrl", "proxy-arp"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "has_mcast_source", "yes"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "has_mcast_source", "yes"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			// {
			// 	Config: CreateAccApplicationEPGUpdatedAttr(rName, "is_attr_based_epg", "yes"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
			// 		resource.TestCheckResourceAttr(resourceName, "is_attr_based_epg", "yes"),
			// 		testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			// },
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "match_t", "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "All"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "match_t", "AtmostOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "AtmostOne"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "match_t", "None"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "match_t", "None"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "pc_enf_pref", "enforced"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "pc_enf_pref", "enforced"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "pref_gr_memb", "include"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "pref_gr_memb", "include"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGUpdatedAttr(rName, "shutdown", "yes"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_updated),
					resource.TestCheckResourceAttr(resourceName, "shutdown", "yes"),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_updated)),
			},
			{
				Config: CreateAccApplicationEPGConfigWithOPtionalValues(rName),
			},
		},
	})
}

func TestAccAciApplicationEPG_NegativeCases(t *testing.T) {
	rName := acctest.RandString(5)
	longAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(65)
	longDescription := acctest.RandString(129)
	// longExceptionTag := "513"
	longFloodOnEncap := acctest.RandString(10)
	longFwdCtrl := acctest.RandString(10)
	longHasMcastSource := acctest.RandString(5)
	longIsAttrBasedEpg := acctest.RandString(5)
	longMatchT := acctest.RandString(10)
	longPcEnfPref := acctest.RandString(10)
	longPrefGrMemb := acctest.RandString(10)
	longPrio := acctest.RandString(15)
	longShutdown := acctest.RandString(5)
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(12)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciApplicationEPGDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccApplicationEPGConfig(rName),
			},
			{
				Config:      CreateAccApplicationEPGWithInvalidApplicationProfile(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvAEPg (.)+`),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, "description", longDescription),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, "annotation", longAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			// {
			// 	Config:	CreateAccApplicationEPGUpdatedAttr(rName,"exception_tag",longExceptionTag),
			// 	ExpectError: regexp.MustCompile(`property  of (.)+ failed validation for value '(.)+'`),
			// },
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, "flood_on_encap", longFloodOnEncap),
				ExpectError: regexp.MustCompile(`expected flood_on_encap to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, "fwd_ctrl", longFwdCtrl),
				ExpectError: regexp.MustCompile(`expected fwd_ctrl to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, "is_attr_based_epg", longIsAttrBasedEpg),
				ExpectError: regexp.MustCompile(`expected is_attr_based_epg to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, "match_t", longMatchT),
				ExpectError: regexp.MustCompile(`expected match_t to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, "has_mcast_source", longHasMcastSource),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name hasMcastSource, class fvAEPg (.)+`),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, "pc_enf_pref", longPcEnfPref),
				ExpectError: regexp.MustCompile(`expected pc_enf_pref to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, "pref_gr_memb", longPrefGrMemb),
				ExpectError: regexp.MustCompile(`expected pref_gr_memb to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, "prio", longPrio),
				ExpectError: regexp.MustCompile(`expected prio to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, "shutdown", longShutdown),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name shutdown, class fvAEPg (.)+`),
			},
			{
				Config:      CreateAccApplicationEPGUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccApplicationEPGConfig(rName),
			},
		},
	})
}

func TestAccAciApplicationEPG_RelationParameters(t *testing.T) {
	var application_epg_default models.ApplicationEPG
	var application_epg_rel1 models.ApplicationEPG
	var application_epg_rel2 models.ApplicationEPG
	resourceName := "aci_application_epg.test"
	rName := acctest.RandString(5)
	randomName1 := acctest.RandString(5)
	randomName2 := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciApplicationEPGDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccApplicationEPGConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_default),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd", ""),
				),
			},
			{
				Config: CreateAccApplicationEPGConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_default),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_aepg_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons_if.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_fc_path_att.#", "0"),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cust_qos_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_dpp_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_sec_inherited.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_graph_def.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_intra_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_node_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_path_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prot_by.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov_def.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_trust_ctrl", ""),
				),
			},
			{
				Config: CreateAccApplicationEPGRelConfig(rName, randomName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_rel1),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd", fmt.Sprintf("uni/tn-%s/BD-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_aepg_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_cons.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons_if.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_cons_if.*", fmt.Sprintf("uni/tn-%s/cif-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_fc_path_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cust_qos_pol", "uni/tn-common/qoscustom-default"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_dpp_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_sec_inherited.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_sec_inherited.*", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_graph_def.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_intra_epg.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_intra_epg.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_node_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_path_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prot_by.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_prot_by.*", fmt.Sprintf("uni/tn-%s/taboo-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_prov.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov_def.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_trust_ctrl", ""),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_rel1),
				),
			},
			{
				Config: CreateAccApplicationEPGRelFinalConfig(rName, randomName1, randomName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_rel2),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd", fmt.Sprintf("uni/tn-%s/BD-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_aepg_mon_pol", fmt.Sprintf("uni/tn-%s/monepg-%s", rName, randomName1)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_cons.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, randomName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_cons.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons_if.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_cons_if.*", fmt.Sprintf("uni/tn-%s/cif-%s", rName, randomName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_cons_if.*", fmt.Sprintf("uni/tn-%s/cif-%s", rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_fc_path_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cust_qos_pol", "uni/tn-common/qoscustom-default"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_dpp_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_sec_inherited.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_sec_inherited.*", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, randomName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_sec_inherited.*", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_graph_def.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_intra_epg.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_intra_epg.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, randomName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_intra_epg.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_node_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_path_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prot_by.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_prot_by.*", fmt.Sprintf("uni/tn-%s/taboo-%s", rName, randomName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_prot_by.*", fmt.Sprintf("uni/tn-%s/taboo-%s", rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_prov.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, randomName1)),
					resource.TestCheckTypeSetElemAttr(resourceName, "relation_fv_rs_prov.*", fmt.Sprintf("uni/tn-%s/brc-%s", rName, randomName2)),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov_def.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_trust_ctrl", ""),
					testAccCheckAciApplicationEPGIdEqual(&application_epg_default, &application_epg_rel2),
				),
			},
			{
				Config: CreateAccApplicationEPGConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciApplicationEPGExists(resourceName, &application_epg_default),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_bd", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_aepg_mon_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cons_if.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_fc_path_att.#", "0"),
					//resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_cust_qos_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_dpp_pol", ""),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_sec_inherited.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_graph_def.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_intra_epg.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_node_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_path_att.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prot_by.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_prov_def.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "relation_fv_rs_trust_ctrl", "")),
			},
		},
	})
}

func TestAccApplicationEpg_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciApplicationEPGDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAcEPGsConfig(rName),
			},
		},
	})
}

func CreateAcEPGsConfig(rName string) string {
	fmt.Println("=== STEP  creating multiple application epgs")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test1"{
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	resource "aci_application_epg" "test2"{
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	resource "aci_application_epg" "test3"{
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}
	`, rName, rName, rName+"1", rName+"2", rName+"3")
	return resource

}

func testAccCheckAciApplicationEPGIdNotEqual(epg1, epg2 *models.ApplicationEPG) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if epg1.DistinguishedName == epg2.DistinguishedName {
			return fmt.Errorf("Application EPG DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciApplicationEPGExists(name string, application_epg *models.ApplicationEPG) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Application EPG %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Application EPG Dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		application_epgFound := models.ApplicationEPGFromContainer(cont)
		if application_epgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Application Profile %s not found", rs.Primary.ID)
		}
		*application_epg = *application_epgFound
		return nil
	}
}

func testAccCheckAciApplicationEPGDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing application EPG destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_application_epg" {
			cont, err := client.Get(rs.Primary.ID)
			application_epg := models.ApplicationEPGFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Application EPG %s still exists", application_epg.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciApplicationEPGIdEqual(epg1, epg2 *models.ApplicationEPG) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if epg1.DistinguishedName != epg2.DistinguishedName {
			return fmt.Errorf("Application EPG DNs are not equal")
		}
		return nil
	}
}

func CreateAccApplicationEPGWithoutApplicationProfile(rName string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg without creating application_profile")
	resource := fmt.Sprintf(`
	resource "aci_application_epg" "test" {
		name = "%s"
	}
	`, rName)
	return resource
}

func CreateAccApplicationEPGWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg without passing name attribute")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	
	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	} 

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
	}

	`, rName, rName)
	return resource
}

func CreateAccApplicationEPGConfigWithParentAndName(parentName, rName string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg with same parent name and different resource name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}
	`, parentName, parentName, rName)
	return resource
}

func CreateAccApplicationEPGConfig(rName string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg with required paramters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccApplicationEPGConfigWithOPtionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
		annotation = "tag"
		description = "From Terraform"
		name_alias = "application_epg_alias"
		prio = "unspecified"
		exception_tag = "0"
		flood_on_encap = "disabled"
		fwd_ctrl = "none"
		has_mcast_source = "no"
		is_attr_based_epg = "no"
		match_t = "AtleastOne"
		pc_enf_pref = "unenforced"
		pref_gr_memb = "exclude"
		shutdown = "no"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccApplicationEPGConfigUpdatedName(rName, longrName string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg creation with invalid name with long length")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}
	`, rName, rName, longrName)
	return resource
}

func CreateAccApplicationEPGConfigUpdateWithoutRequiredParameters(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_application_epg" "test"{
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccApplicationEPGUpdatedAttr(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg updation with attribute values updation")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
		%s = "%s"
	}
	`, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccApplicationEPGWithInvalidApplicationProfile(rName string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg updation with attribute values updation")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_application_epg" "test"{
		application_profile_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccApplicationEPGRelConfig(rName, relName string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg with relations")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract" "name" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_imported_contract" "name" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_monitoring_policy" "name" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_taboo_contract" "name" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "name1" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
		relation_fv_rs_bd = aci_bridge_domain.test.id
		relation_fv_rs_cust_qos_pol = "uni/tn-common/qoscustom-default"
		relation_fv_rs_prov = [aci_contract.name.id] 
		relation_fv_rs_cons_if = [aci_imported_contract.name.id]
		relation_fv_rs_aepg_mon_pol = aci_monitoring_policy.name.id
		relation_fv_rs_prot_by = [aci_taboo_contract.name.id]
		relation_fv_rs_cons = [aci_contract.name.id]
		relation_fv_rs_intra_epg = [aci_contract.name.id]
		relation_fv_rs_sec_inherited = [aci_application_epg.name1.id]
	}
	`, rName, rName, relName, relName, relName, relName, relName, relName, rName)
	return resource
}

func CreateAccApplicationEPGRelFinalConfig(rName, relName1, relName2 string) string {
	fmt.Println("=== STEP  Basic: Testing application_epg with relations")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_bridge_domain" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_contract" "name" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_contract" "name1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_imported_contract" "name" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_imported_contract" "name1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_monitoring_policy" "name" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
	resource "aci_taboo_contract" "name" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_taboo_contract" "name1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_application_epg" "name1" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_application_epg" "name2" {
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
	}

	resource "aci_application_epg" "test"{
		application_profile_dn = aci_application_profile.test.id
		name = "%s"
		relation_fv_rs_bd = aci_bridge_domain.test.id
		relation_fv_rs_cust_qos_pol = "uni/tn-common/qoscustom-default"
		relation_fv_rs_prov = [aci_contract.name.id,aci_contract.name1.id] 
		relation_fv_rs_cons_if = [aci_imported_contract.name.id,aci_imported_contract.name1.id]
		relation_fv_rs_aepg_mon_pol = aci_monitoring_policy.name.id
		relation_fv_rs_prot_by = [aci_taboo_contract.name.id,aci_taboo_contract.name1.id]
		relation_fv_rs_cons = [aci_contract.name.id,aci_contract.name1.id]
		relation_fv_rs_intra_epg = [aci_contract.name.id,aci_contract.name1.id]
		relation_fv_rs_sec_inherited = [aci_application_epg.name1.id,aci_application_epg.name2.id]
	}
	`, rName, rName, relName1, relName1, relName2, relName1, relName2, relName1, relName1, relName2, relName1, relName2, rName)
	return resource
}
