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

func TestAccAciOSPFTimers_Basic(t *testing.T) {
	var ospf_timers_default models.OSPFTimersPolicy
	var ospf_timers_updated models.OSPFTimersPolicy
	resourceName := "aci_ospf_timers.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOSPFTimersDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateOSPFTimersWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateOSPFTimersWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccOSPFTimersConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "bw_ref", "40000"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "dist", "110"),
					resource.TestCheckResourceAttr(resourceName, "gr_ctrl", ""),
					resource.TestCheckResourceAttr(resourceName, "lsa_arrival_intvl", "1000"),
					resource.TestCheckResourceAttr(resourceName, "lsa_gp_pacing_intvl", "10"),
					resource.TestCheckResourceAttr(resourceName, "lsa_hold_intvl", "5000"),
					resource.TestCheckResourceAttr(resourceName, "lsa_max_intvl", "5000"),
					resource.TestCheckResourceAttr(resourceName, "lsa_start_intvl", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_ecmp", "8"),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_action", "reject"),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_num", "20000"),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_reset_intvl", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_sleep_cnt", "5"),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_sleep_intvl", "5"),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_thresh", "75"),
					resource.TestCheckResourceAttr(resourceName, "spf_hold_intvl", "1000"),
					resource.TestCheckResourceAttr(resourceName, "spf_init_intvl", "200"),
					resource.TestCheckResourceAttr(resourceName, "spf_max_intvl", "5000"),
				),
			},
			{
				Config: CreateAccOSPFTimersConfigWithOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ospf_timers"),
					resource.TestCheckResourceAttr(resourceName, "bw_ref", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "name-lookup"),
					resource.TestCheckResourceAttr(resourceName, "dist", "1"),
					resource.TestCheckResourceAttr(resourceName, "gr_ctrl", "helper"),
					resource.TestCheckResourceAttr(resourceName, "lsa_arrival_intvl", "10"),
					resource.TestCheckResourceAttr(resourceName, "lsa_gp_pacing_intvl", "1"),
					resource.TestCheckResourceAttr(resourceName, "lsa_hold_intvl", "50"),
					resource.TestCheckResourceAttr(resourceName, "lsa_max_intvl", "50"),
					resource.TestCheckResourceAttr(resourceName, "lsa_start_intvl", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_ecmp", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_action", "log"),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_reset_intvl", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_sleep_cnt", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_sleep_intvl", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_thresh", "1"),
					resource.TestCheckResourceAttr(resourceName, "spf_hold_intvl", "1"),
					resource.TestCheckResourceAttr(resourceName, "spf_init_intvl", "1"),
					resource.TestCheckResourceAttr(resourceName, "spf_max_intvl", "1"),
					testAccCheckAciOSPFTimersIdEqual(&ospf_timers_default, &ospf_timers_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccOSPFTimersConfig(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccOSPFTimersRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccOSPFTimersConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciOSPFTimersIdNotEqual(&ospf_timers_default, &ospf_timers_updated),
				),
			},
			{
				Config: CreateAccOSPFTimersConfig(fvTenantName, rName),
			},
			{
				Config: CreateAccOSPFTimersConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciOSPFTimersIdNotEqual(&ospf_timers_default, &ospf_timers_updated),
				),
			},
		},
	})
}

func TestAccAciOSPFTimers_Update(t *testing.T) {
	var ospf_timers_default models.OSPFTimersPolicy
	var ospf_timers_updated models.OSPFTimersPolicy
	resourceName := "aci_ospf_timers.test"
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOSPFTimersDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOSPFTimersConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_default),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"name-lookup", "pfx-suppress"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "name-lookup"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "pfx-suppress"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"pfx-suppress", "name-lookup"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "pfx-suppress"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.1", "name-lookup"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"pfx-suppress"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ctrl.0", "pfx-suppress"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "bw_ref", "4000000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "bw_ref", "4000000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "bw_ref", "100000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "bw_ref", "100000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "dist", "225"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "dist", "225"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "dist", "170"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "dist", "170"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_arrival_intvl", "600000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "lsa_arrival_intvl", "600000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_arrival_intvl", "10000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "lsa_arrival_intvl", "10000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_gp_pacing_intvl", "1800"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "lsa_gp_pacing_intvl", "1800"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_gp_pacing_intvl", "1000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "lsa_gp_pacing_intvl", "1000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_hold_intvl", "30000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "lsa_hold_intvl", "30000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_hold_intvl", "2000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "lsa_hold_intvl", "2000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_max_intvl", "30000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "lsa_max_intvl", "30000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_max_intvl", "2000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "lsa_max_intvl", "2000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_start_intvl", "5000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "lsa_start_intvl", "5000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_start_intvl", "2500"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "lsa_start_intvl", "2500"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_ecmp", "64"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_ecmp", "64"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_ecmp", "30"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_ecmp", "30"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_action", "restart"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_action", "restart"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_num", "4294967295"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_num", "4294967295"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_num", "10000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_num", "10000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_reset_intvl", "1440"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_reset_intvl", "1440"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_reset_intvl", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_reset_intvl", "100"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_sleep_cnt", "4294967295"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_sleep_cnt", "4294967295"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_sleep_cnt", "10000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_sleep_cnt", "10000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_sleep_intvl", "1440"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_sleep_intvl", "1440"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_sleep_intvl", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_sleep_intvl", "100"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_thresh", "100"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_thresh", "100"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_thresh", "50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "max_lsa_thresh", "50"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "spf_hold_intvl", "60000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "spf_hold_intvl", "60000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "spf_hold_intvl", "5000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "spf_hold_intvl", "5000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "spf_init_intvl", "60000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "spf_init_intvl", "60000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "spf_init_intvl", "5000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "spf_init_intvl", "5000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "spf_max_intvl", "60000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "spf_max_intvl", "60000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "spf_max_intvl", "1000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "spf_max_intvl", "1000"),
				),
			},
			{
				Config: CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "gr_ctrl", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFTimersExists(resourceName, &ospf_timers_updated),
					resource.TestCheckResourceAttr(resourceName, "gr_ctrl", ""),
				),
			},
			{
				Config: CreateAccOSPFTimersConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciOSPFTimers_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOSPFTimersDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOSPFTimersConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccOSPFTimersWithInValidParentDn(fvTenantName, rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "bw_ref", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "bw_ref", "0"),
				ExpectError: regexp.MustCompile(`Property bwRef of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttrList(fvTenantName, rName, "ctrl", StringListtoString([]string{"name-lookup", "name-lookup"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "dist", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "dist", "0"),
				ExpectError: regexp.MustCompile(`Property dist of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "gr_ctrl", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_arrival_intvl", "9"),
				ExpectError: regexp.MustCompile(`Property lsaArrivalIntvl of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_arrival_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_gp_pacing_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_gp_pacing_intvl", "0"),
				ExpectError: regexp.MustCompile(`Property lsaGpPacingIntvl of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_hold_intvl", "40"),
				ExpectError: regexp.MustCompile(`Property lsaHoldIntvl of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_hold_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_max_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_max_intvl", "40"),
				ExpectError: regexp.MustCompile(`Property lsaMaxIntvl of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_start_intvl", "5001"),
				ExpectError: regexp.MustCompile(`Property lsaStartIntvl of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "lsa_start_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_ecmp", "65"),
				ExpectError: regexp.MustCompile(`Property maxEcmp of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_ecmp", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_action", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_num", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_num", "0"),
				ExpectError: regexp.MustCompile(`Property maxLsaNum of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_reset_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_reset_intvl", "0"),
				ExpectError: regexp.MustCompile(`Property maxLsaResetIntvl of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_sleep_cnt", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_sleep_cnt", "0"),
				ExpectError: regexp.MustCompile(`Property maxLsaSleepCnt of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_sleep_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_sleep_intvl", "0"),
				ExpectError: regexp.MustCompile(`Property maxLsaSleepIntvl of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_thresh", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "max_lsa_thresh", "0"),
				ExpectError: regexp.MustCompile(`Property maxLsaThresh of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "spf_hold_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "spf_hold_intvl", "0"),
				ExpectError: regexp.MustCompile(`Property spfHoldIntvl of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "spf_init_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "spf_init_intvl", "0"),
				ExpectError: regexp.MustCompile(`Property spfInitIntvl of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "spf_max_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, "spf_max_intvl", "0"),
				ExpectError: regexp.MustCompile(`Property spfMaxIntvl of (.)* is out of range`),
			},
			{
				Config:      CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccOSPFTimersConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciOSPFTimers_MultipleCreateDestroy(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOSPFTimersDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOSPFTimersMultipleConfig(rName, rName),
			},
		},
	})
}
func testAccCheckAciOSPFTimersExists(name string, ospf_timers *models.OSPFTimersPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("OSPF Timers %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OSPF Timers dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ospf_timersFound := models.OSPFTimersPolicyFromContainer(cont)
		if ospf_timersFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("OSPF Timers %s not found", rs.Primary.ID)
		}
		*ospf_timers = *ospf_timersFound
		return nil
	}
}

func testAccCheckAciOSPFTimersDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing ospf_timers destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ospf_timers" {
			cont, err := client.Get(rs.Primary.ID)
			ospf_timers := models.OSPFTimersPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("OSPF Timers %s Still exists", ospf_timers.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciOSPFTimersIdEqual(m1, m2 *models.OSPFTimersPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("ospf_timers DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciOSPFTimersIdNotEqual(m1, m2 *models.OSPFTimersPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("ospf_timers DNs are equal")
		}
		return nil
	}
}

func CreateOSPFTimersWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ospf_timers creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_ospf_timers" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccOSPFTimersConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing ospf_timers creation with Tenant Name %s and OSPF Timers Name %s required arguments\n", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccOSPFTimersConfig(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing ospf_timers creation with Tenant Name %s and OSPF Timers Name %s required arguments\n", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccOSPFTimersMultipleConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing Multiple ospf_timers creation with required arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_ospf_timers" "test1" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_ospf_timers" "test2" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_ospf_timers" "test3" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccOSPFTimersWithInValidParentDn(prName, rName string) string {
	fmt.Println("=== STEP  Negative Case: testing ospf_timers creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_aaa_domain" "test"{
		name = "%s"
	}
	resource "aci_ospf_timers" "test" {
		tenant_dn  = aci_aaa_domain.test.id
		name  = "%s"
	}
	`, prName, rName)
	return resource
}

func CreateAccOSPFTimersConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing ospf_timers creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_timers" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ospf_timers"
		bw_ref = "1"
		ctrl = ["name-lookup"]
		dist = "1"
		gr_ctrl = "helper"
		lsa_arrival_intvl = "10"
		lsa_gp_pacing_intvl = "1"
		lsa_hold_intvl = "50"
		lsa_max_intvl = "50"
		lsa_start_intvl = "0"
		max_ecmp = "1"
		max_lsa_action = "log"
		max_lsa_num = "1"
		max_lsa_reset_intvl = "1"
		max_lsa_sleep_cnt = "1"
		max_lsa_sleep_intvl = "1"
		max_lsa_thresh = "1"
		spf_hold_intvl = "1"
		spf_init_intvl = "1"
		spf_max_intvl = "1"
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccOSPFTimersRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing ospf_timers updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_ospf_timers" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ospf_timers"
		bw_ref = "2"
		ctrl = ["name-lookup"]
		dist = "2"
		gr_ctrl = "helper"
		lsa_arrival_intvl = "11"
		lsa_gp_pacing_intvl = "2"
		lsa_hold_intvl = "51"
		lsa_max_intvl = "51"
		lsa_start_intvl = "1"
		max_ecmp = "2"
		max_lsa_action = "log"
		max_lsa_num = "2"
		max_lsa_reset_intvl = "2"
		max_lsa_sleep_cnt = "2"
		max_lsa_sleep_intvl = "2"
		max_lsa_thresh = "2"
		spf_hold_intvl = "2"
		spf_init_intvl = "2"
		spf_max_intvl = "2"
	}
	`)

	return resource
}

func CreateAccOSPFTimersUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing ospf_timers attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}

func CreateAccOSPFTimersUpdatedAttrList(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing ospf_timers attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_timers" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = %s
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
