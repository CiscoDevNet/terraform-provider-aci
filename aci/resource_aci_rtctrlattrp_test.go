package aci

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciActionRuleProfile_Basic(t *testing.T) {
	var action_rule_profile models.ActionRuleProfile
	var set_route_tag models.RtctrlSetTag
	var set_preference models.RtctrlSetPref
	var set_weight models.RtctrlSetWeight
	var set_metric models.RtctrlSetRtMetric
	var set_metric_type models.RtctrlSetRtMetricType
	var set_next_hop models.RtctrlSetNh
	var set_communities models.RtctrlSetComm
	var set_dampening models.RtctrlSetDamp
	var set_as_path_prepend_last_as models.SetASPath
	var set_as_path_prepend_as models.ASNumber

	fv_tenant_name := acctest.RandString(5)
	rtctrl_attr_p_name := acctest.RandString(5)
	set_route_tag_value := acctest.RandIntRange(0, 2147483647)
	set_preference_value := acctest.RandIntRange(0, 2147483647)
	set_weight_value := acctest.RandIntRange(0, 65535)
	set_as_path_prepend_last_as_value := acctest.RandIntRange(1, 10)
	set_as_path_prepend_as_order_value := acctest.RandIntRange(0, 31)
	set_as_path_prepend_as_asn_value := acctest.RandIntRange(1, 2147483647)
	set_metric_value := acctest.RandIntRange(0, 2147483647)
	description := "action_rule_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciActionRuleProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciActionRuleProfileConfig_basic(
					fv_tenant_name,
					rtctrl_attr_p_name,
					set_route_tag_value,
					set_preference_value,
					set_weight_value,
					set_as_path_prepend_last_as_value,
					set_as_path_prepend_as_order_value,
					set_as_path_prepend_as_asn_value,
					set_metric_value,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists(
						"aci_action_rule_profile.fooaction_rule_profile",
						&action_rule_profile,
						&set_route_tag,
						&set_preference,
						&set_weight,
						&set_metric,
						&set_metric_type,
						&set_next_hop,
						&set_communities,
						&set_dampening,
						&set_as_path_prepend_last_as,
						&set_as_path_prepend_as,
					),
					testAccCheckAciActionRuleProfileAttributes(
						fv_tenant_name,
						rtctrl_attr_p_name,
						description,
						set_route_tag_value,
						set_preference_value,
						set_weight_value,
						set_as_path_prepend_last_as_value,
						set_as_path_prepend_as_order_value,
						set_as_path_prepend_as_asn_value,
						set_metric_value,
						&action_rule_profile,
						&set_route_tag,
						&set_preference,
						&set_weight,
						&set_metric,
						&set_metric_type,
						&set_next_hop,
						&set_communities,
						&set_dampening,
						&set_as_path_prepend_last_as,
						&set_as_path_prepend_as,
					),
				),
			},
		},
	})
}

func TestAccAciActionRuleProfile_Update(t *testing.T) {
	var action_rule_profile models.ActionRuleProfile
	var set_route_tag models.RtctrlSetTag
	var set_preference models.RtctrlSetPref
	var set_weight models.RtctrlSetWeight
	var set_metric models.RtctrlSetRtMetric
	var set_metric_type models.RtctrlSetRtMetricType
	var set_next_hop models.RtctrlSetNh
	var set_communities models.RtctrlSetComm
	var set_dampening models.RtctrlSetDamp
	var set_as_path_prepend_last_as models.SetASPath
	var set_as_path_prepend_as models.ASNumber

	fv_tenant_name := acctest.RandString(5)
	rtctrl_attr_p_name := acctest.RandString(5)
	description := "action_rule_profile created while acceptance testing"
	set_route_tag_value := acctest.RandIntRange(0, 2147483647)
	set_preference_value := acctest.RandIntRange(0, 2147483647)
	set_weight_value := acctest.RandIntRange(0, 65535)
	set_as_path_prepend_last_as_value := acctest.RandIntRange(1, 10)
	set_as_path_prepend_as_order_value := acctest.RandIntRange(0, 31)
	set_as_path_prepend_as_asn_value := acctest.RandIntRange(1, 2147483647)
	set_metric_value := acctest.RandIntRange(0, 2147483647)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciActionRuleProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciActionRuleProfileConfig_basic(
					fv_tenant_name,
					rtctrl_attr_p_name,
					set_route_tag_value,
					set_preference_value,
					set_weight_value,
					set_as_path_prepend_last_as_value,
					set_as_path_prepend_as_order_value,
					set_as_path_prepend_as_asn_value,
					set_metric_value,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists(
						"aci_action_rule_profile.fooaction_rule_profile",
						&action_rule_profile,
						&set_route_tag,
						&set_preference,
						&set_weight,
						&set_metric,
						&set_metric_type,
						&set_next_hop,
						&set_communities,
						&set_dampening,
						&set_as_path_prepend_last_as,
						&set_as_path_prepend_as,
					),
					testAccCheckAciActionRuleProfileAttributes(
						fv_tenant_name,
						rtctrl_attr_p_name,
						description,
						set_route_tag_value,
						set_preference_value,
						set_weight_value,
						set_as_path_prepend_last_as_value,
						set_as_path_prepend_as_order_value,
						set_as_path_prepend_as_asn_value,
						set_metric_value,
						&action_rule_profile,
						&set_route_tag,
						&set_preference,
						&set_weight,
						&set_metric,
						&set_metric_type,
						&set_next_hop,
						&set_communities,
						&set_dampening,
						&set_as_path_prepend_last_as,
						&set_as_path_prepend_as,
					),
				),
			},
			{
				Config: testAccCheckAciActionRuleProfileConfig_basic(
					fv_tenant_name,
					rtctrl_attr_p_name,
					set_route_tag_value,
					set_preference_value,
					set_weight_value,
					set_as_path_prepend_last_as_value,
					set_as_path_prepend_as_order_value,
					set_as_path_prepend_as_asn_value,
					set_metric_value,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileExists(
						"aci_action_rule_profile.fooaction_rule_profile",
						&action_rule_profile,
						&set_route_tag,
						&set_preference,
						&set_weight,
						&set_metric,
						&set_metric_type,
						&set_next_hop,
						&set_communities,
						&set_dampening,
						&set_as_path_prepend_last_as,
						&set_as_path_prepend_as,
					),
					testAccCheckAciActionRuleProfileAttributes(
						fv_tenant_name,
						rtctrl_attr_p_name,
						description,
						set_route_tag_value,
						set_preference_value,
						set_weight_value,
						set_as_path_prepend_last_as_value,
						set_as_path_prepend_as_order_value,
						set_as_path_prepend_as_asn_value,
						set_metric_value,
						&action_rule_profile,
						&set_route_tag,
						&set_preference,
						&set_weight,
						&set_metric,
						&set_metric_type,
						&set_next_hop,
						&set_communities,
						&set_dampening,
						&set_as_path_prepend_last_as,
						&set_as_path_prepend_as,
					),
				),
			},
		},
	})
}

func testAccCheckAciActionRuleProfileConfig_basic(
	fv_tenant_name,
	rtctrl_attr_p_name string,
	set_route_tag_value,
	set_preference_value,
	set_weight_value,
	set_as_path_prepend_last_as_value,
	set_as_path_prepend_as_order_value,
	set_as_path_prepend_as_asn_value,
	set_metric_value int,
) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_action_rule_profile" "fooaction_rule_profile" {
		name 		= "%s"
		description = "action_rule_profile created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
		set_route_tag = %d
		set_preference = %d
		set_weight      = %d
		set_as_path_prepend_last_as = %d
		set_as_path_prepend_as {
			order = %d
			asn   = %d
		  }
		set_metric      = %d
		set_metric_type = "ospf-type1"
		set_next_hop    = "1.1.1.1"
		set_communities = {
		  community = "no-advertise"
		  criteria  = "replace"
		}
		set_dampening = {
			half_life         = 10
			reuse             = 1
			suppress          = 10
			max_suppress_time = 100
		  }

	}
	`, fv_tenant_name, rtctrl_attr_p_name, set_route_tag_value, set_preference_value, set_weight_value, set_as_path_prepend_last_as_value, set_as_path_prepend_as_order_value, set_as_path_prepend_as_asn_value, set_metric_value)

}

func testAccCheckAciActionRuleProfileExists(
	name string,
	action_rule_profile *models.ActionRuleProfile,
	set_route_tag *models.RtctrlSetTag,
	set_preference *models.RtctrlSetPref,
	set_weight *models.RtctrlSetWeight,
	set_metric *models.RtctrlSetRtMetric,
	set_metric_type *models.RtctrlSetRtMetricType,
	set_next_hop *models.RtctrlSetNh,
	set_communities *models.RtctrlSetComm,
	set_dampening *models.RtctrlSetDamp,
	set_as_path_prepend_last_as *models.SetASPath,
	set_as_path_prepend_as *models.ASNumber,
) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Action Rule Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Action Rule Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		action_rule_profileFound := models.ActionRuleProfileFromContainer(cont)
		if action_rule_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Action Rule Profile %s not found", rs.Primary.ID)
		}
		*action_rule_profile = *action_rule_profileFound

		// rtctrlSetTag - Beginning of Import
		set_route_tag_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnrtctrlSetTag)
		set_route_tag_cont, err := client.Get(set_route_tag_dn)
		if err != nil {
			return err
		}

		set_route_tag_found := models.RtctrlSetTagFromContainer(set_route_tag_cont)
		if set_route_tag_found.DistinguishedName != set_route_tag_dn {
			return fmt.Errorf("Set Route Tag %s not found", set_route_tag_dn)
		}
		*set_route_tag = *set_route_tag_found
		// rtctrlSetTag - Import finished successfully

		// rtctrlSetPref - Beginning of Import
		set_preference_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnrtctrlSetPref)
		set_preference_cont, err := client.Get(set_preference_dn)
		if err != nil {
			return err
		}

		set_preference_found := models.RtctrlSetPrefFromContainer(set_preference_cont)
		if set_preference_found.DistinguishedName != set_preference_dn {
			return fmt.Errorf("Set Preference %s not found", set_preference_dn)
		}
		*set_preference = *set_preference_found
		// rtctrlSetPref - Import finished successfully

		// rtctrlSetWeight - Beginning Import
		set_weight_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnrtctrlSetWeight)
		set_weight_cont, err := client.Get(set_weight_dn)
		if err != nil {
			return err
		}

		set_weight_found := models.RtctrlSetWeightFromContainer(set_weight_cont)
		if set_weight_found.DistinguishedName != set_weight_dn {
			return fmt.Errorf("Set Weight %s not found", set_weight_dn)
		}
		*set_weight = *set_weight_found
		// rtctrlSetWeight - Import finished successfully

		// rtctrlSetRtMetric - Beginning of Import
		set_metric_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnrtctrlSetRtMetric)
		set_metric_cont, err := client.Get(set_metric_dn)
		if err != nil {
			return err
		}

		set_metric_found := models.RtctrlSetRtMetricFromContainer(set_metric_cont)
		if set_metric_found.DistinguishedName != set_metric_dn {
			return fmt.Errorf("Set Metric %s not found", set_metric_dn)
		}
		*set_metric = *set_metric_found
		// rtctrlSetRtMetric - Import finished successfully

		// rtctrlSetRtMetricType - Beginning of Import
		set_metric_type_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnrtctrlSetRtMetricType)
		set_metric_type_cont, err := client.Get(set_metric_type_dn)
		if err != nil {
			return err
		}

		set_metric_type_found := models.RtctrlSetRtMetricTypeFromContainer(set_metric_type_cont)
		if set_metric_type_found.DistinguishedName != set_metric_type_dn {
			return fmt.Errorf("Set Metric Type %s not found", set_metric_type_dn)
		}
		*set_metric_type = *set_metric_type_found
		// rtctrlSetRtMetricType - Import finished successfully

		// rtctrlSetNh - Beginning of Import
		set_next_hop_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnrtctrlSetNh)
		set_next_hop_cont, err := client.Get(set_next_hop_dn)
		if err != nil {
			return err
		}

		set_next_hop_found := models.RtctrlSetNhFromContainer(set_next_hop_cont)
		if set_next_hop_found.DistinguishedName != set_next_hop_dn {
			return fmt.Errorf("Set Next Hop %s not found", set_next_hop_dn)
		}
		*set_next_hop = *set_next_hop_found
		// rtctrlSetNh - Import finished successfully

		// rtctrlSetComm - Beginning of Import
		set_communities_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnrtctrlSetComm)
		set_communities_cont, err := client.Get(set_communities_dn)
		if err != nil {
			return err
		}

		set_communities_found := models.RtctrlSetCommFromContainer(set_communities_cont)
		if set_communities_found.DistinguishedName != set_communities_dn {
			return fmt.Errorf("Set Communities %s not found", set_communities_dn)
		}
		*set_communities = *set_communities_found
		// rtctrlSetComm - Import finished successfully

		// rtctrlSetDamp - Beginning of Import
		set_dampening_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnrtctrlSetDamp)
		set_dampening_cont, err := client.Get(set_dampening_dn)
		if err != nil {
			return err
		}

		set_dampening_found := models.RtctrlSetDampFromContainer(set_dampening_cont)
		if set_dampening_found.DistinguishedName != set_dampening_dn {
			return fmt.Errorf("Set Dampening %s not found", set_dampening_dn)
		}
		*set_dampening = *set_dampening_found
		// rtctrlSetDamp - Import finished successfully

		// rtctrlSetASPath - Beginning of Import
		set_as_path_prepend_last_as_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnrtctrlSetASPath, "prepend-last-as")
		set_as_path_prepend_last_as_cont, err := client.Get(set_as_path_prepend_last_as_dn)
		if err != nil {
			return err
		}

		set_as_path_prepend_last_as_found := models.SetASPathFromContainer(set_as_path_prepend_last_as_cont)
		if set_as_path_prepend_last_as_found.DistinguishedName != set_as_path_prepend_last_as_dn {
			return fmt.Errorf("Set As Path Prepend Last As %s not found", set_as_path_prepend_last_as_dn)
		}
		*set_as_path_prepend_last_as = *set_as_path_prepend_last_as_found
		// rtctrlSetASPath - Import finished successfully

		// rtctrlSetASPathASN - Beginning of Import
		set_as_path_prepend_as_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnrtctrlSetASPath, "prepend")
		ReadRelationASNumberData, err := client.ListSetAsPathASNs(set_as_path_prepend_as_dn)
		if err == nil {
			for _, record := range ReadRelationASNumberData {
				*set_as_path_prepend_as = *record
			}
		} else {
			return fmt.Errorf("Set As Path Prepend Last As %s not found", set_as_path_prepend_as_dn)
		}
		// rtctrlSetASPathASN - Import finished successfully

		return nil
	}
}

func testAccCheckAciActionRuleProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_action_rule_profile" {
			cont, err := client.Get(rs.Primary.ID)
			action_rule_profile := models.ActionRuleProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Action Rule Profile %s Still exists", action_rule_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciActionRuleProfileAttributes(
	fv_tenant_name,
	rtctrl_attr_p_name,
	description string,
	set_route_tag_value,
	set_preference_value,
	set_weight_value,
	set_as_path_prepend_last_as_value,
	set_as_path_prepend_as_order_value,
	set_as_path_prepend_as_asn_value,
	set_metric_value int,
	action_rule_profile *models.ActionRuleProfile,
	set_route_tag *models.RtctrlSetTag,
	set_preference *models.RtctrlSetPref,
	set_weight *models.RtctrlSetWeight,
	set_metric *models.RtctrlSetRtMetric,
	set_metric_type *models.RtctrlSetRtMetricType,
	set_next_hop *models.RtctrlSetNh,
	set_communities *models.RtctrlSetComm,
	set_dampening *models.RtctrlSetDamp,
	set_as_path_prepend_last_as *models.SetASPath,
	set_as_path_prepend_as *models.ASNumber,

) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if rtctrl_attr_p_name != GetMOName(action_rule_profile.DistinguishedName) {
			return fmt.Errorf("Bad rtctrl_attr_p %s", GetMOName(action_rule_profile.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(action_rule_profile.DistinguishedName, fmt.Sprintf("/attr-%s", rtctrl_attr_p_name))) {
			return fmt.Errorf(" Bad fv_tenant %s", GetMOName(GetParentDn(action_rule_profile.DistinguishedName, fmt.Sprintf("/attr-%s", rtctrl_attr_p_name))))
		}
		if description != action_rule_profile.Description {
			return fmt.Errorf("Bad action_rule_profile Description %s", action_rule_profile.Description)
		}
		if strconv.Itoa(set_route_tag_value) != set_route_tag.Tag {
			return fmt.Errorf("Bad set_route_tag Tag value %s", set_route_tag.Tag)
		}

		if strconv.Itoa(set_preference_value) != set_preference.LocalPref {
			return fmt.Errorf("Bad set_preference LocalPref value %s", set_preference.LocalPref)
		}

		if strconv.Itoa(set_weight_value) != set_weight.Weight {
			return fmt.Errorf("Bad set_weight Weight value %s", set_weight.Weight)
		}

		if strconv.Itoa(set_metric_value) != set_metric.Metric {
			return fmt.Errorf("Bad set_metric Metric value %s", set_metric.Metric)
		}

		if "ospf-type1" != set_metric_type.MetricType {
			return fmt.Errorf("Bad set_metric_type MetricType value %s", set_metric_type.MetricType)
		}

		if "1.1.1.1" != set_next_hop.Addr {
			return fmt.Errorf("Bad set_next_hop Addr value %s", set_next_hop.Addr)
		}

		if "replace" != set_communities.SetCriteria {
			return fmt.Errorf("Bad set_communities SetCriteria value %s", set_communities.SetCriteria)
		}

		if "no-advertise" != set_communities.Community {
			return fmt.Errorf("Bad set_communities Community value %s", set_communities.Community)
		}

		if "10" != set_dampening.HalfLife {
			return fmt.Errorf("Bad set_dampening HalfLife value %s", set_dampening.HalfLife)
		}

		if "1" != set_dampening.Reuse {
			return fmt.Errorf("Bad set_dampening Reuse value %s", set_dampening.Reuse)
		}

		if "10" != set_dampening.Suppress {
			return fmt.Errorf("Bad set_dampening Suppress value %s", set_dampening.Suppress)
		}

		if "100" != set_dampening.MaxSuppressTime {
			return fmt.Errorf("Bad set_dampening MaxSuppressTime value %s", set_dampening.MaxSuppressTime)
		}

		if strconv.Itoa(set_as_path_prepend_last_as_value) != set_as_path_prepend_last_as.Lastnum {
			return fmt.Errorf("Bad set_as_path_prepend_last_as Lastnum value %s", set_as_path_prepend_last_as.Lastnum)
		}

		if strconv.Itoa(set_as_path_prepend_as_order_value) != set_as_path_prepend_as.Order {
			return fmt.Errorf("Bad set_as_path_prepend_as Order value %s", set_as_path_prepend_as.Order)
		}

		if strconv.Itoa(set_as_path_prepend_as_asn_value) != set_as_path_prepend_as.Asn {
			return fmt.Errorf("Bad set_as_path_prepend_as Asn value %s", set_as_path_prepend_as.Asn)
		}

		return nil
	}
}

func TestAccAciActionRuleProfileMultipath_Basic(t *testing.T) {
	var action_rule_profile models.ActionRuleProfile
	var next_hop_propagation models.NexthopUnchangedAction
	var multipath models.RedistributeMultipathAction

	fv_tenant_name := acctest.RandString(5)
	rtctrl_attr_p_name := acctest.RandString(5)
	description := "action_rule_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciActionRuleProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciActionRuleProfileMultipathConfig_basic(fv_tenant_name, rtctrl_attr_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileMultipathExists("aci_action_rule_profile.fooaction_rule_profile",
						&action_rule_profile,
						&next_hop_propagation,
						&multipath,
					),
					testAccCheckAciActionRuleProfileMultipathAttributes(fv_tenant_name, rtctrl_attr_p_name, description,
						&action_rule_profile,
						&next_hop_propagation,
						&multipath,
					),
				),
			},
		},
	})
}

func TestAccAciActionRuleProfileMultipath_Update(t *testing.T) {
	var action_rule_profile models.ActionRuleProfile
	var next_hop_propagation models.NexthopUnchangedAction
	var multipath models.RedistributeMultipathAction

	fv_tenant_name := acctest.RandString(5)
	rtctrl_attr_p_name := acctest.RandString(5)
	description := "action_rule_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciActionRuleProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciActionRuleProfileMultipathConfig_basic(fv_tenant_name, rtctrl_attr_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileMultipathExists("aci_action_rule_profile.fooaction_rule_profile",
						&action_rule_profile,
						&next_hop_propagation,
						&multipath,
					),
					testAccCheckAciActionRuleProfileMultipathAttributes(fv_tenant_name, rtctrl_attr_p_name, description,
						&action_rule_profile,
						&next_hop_propagation,
						&multipath,
					),
				),
			},
			{
				Config: testAccCheckAciActionRuleProfileMultipathConfig_basic(fv_tenant_name, rtctrl_attr_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciActionRuleProfileMultipathExists("aci_action_rule_profile.fooaction_rule_profile",
						&action_rule_profile,
						&next_hop_propagation,
						&multipath,
					),
					testAccCheckAciActionRuleProfileMultipathAttributes(fv_tenant_name, rtctrl_attr_p_name, description,
						&action_rule_profile,
						&next_hop_propagation,
						&multipath,
					),
				),
			},
		},
	})
}

func testAccCheckAciActionRuleProfileMultipathConfig_basic(fv_tenant_name, rtctrl_attr_p_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_action_rule_profile" "fooaction_rule_profile" {
		name 		= "%s"
		description = "action_rule_profile created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
		next_hop_propagation        = "yes" # Can not be configured along with set_route_tag
		multipath                   = "yes" # Can not be configured along with set_route_tag
	}
	`, fv_tenant_name, rtctrl_attr_p_name)
}

func testAccCheckAciActionRuleProfileMultipathExists(name string,
	action_rule_profile *models.ActionRuleProfile,
	next_hop_propagation *models.NexthopUnchangedAction,
	multipath *models.RedistributeMultipathAction,

) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Action Rule Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Action Rule Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		action_rule_profileFound := models.ActionRuleProfileFromContainer(cont)
		if action_rule_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Action Rule Profile %s not found", rs.Primary.ID)
		}
		*action_rule_profile = *action_rule_profileFound

		// rtctrlSetNhUnchanged - Beginning of Import
		next_hop_propagation_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnrtctrlSetNhUnchanged)
		next_hop_propagation_cont, err := client.Get(next_hop_propagation_dn)
		if err != nil {
			return err
		}

		next_hop_propagation_found := models.NexthopUnchangedActionFromContainer(next_hop_propagation_cont)
		if next_hop_propagation_found.DistinguishedName != next_hop_propagation_dn {
			return fmt.Errorf("Next Hop Propagation %s not found", next_hop_propagation_dn)
		}
		*next_hop_propagation = *next_hop_propagation_found
		// rtctrlSetNhUnchanged - Import finished successfully

		// rtctrlSetRedistMultipath - Beginning of Import
		multipath_dn := rs.Primary.ID + fmt.Sprintf("/"+models.RnrtctrlSetRedistMultipath)
		multipath_cont, err := client.Get(multipath_dn)
		if err != nil {
			return err
		}

		multipath_found := models.RedistributeMultipathActionFromContainer(multipath_cont)
		if multipath_found.DistinguishedName != multipath_dn {
			return fmt.Errorf("Multipath %s not found", multipath_dn)
		}
		*multipath = *multipath_found
		// rtctrlSetRedistMultipath - Import finished successfully

		return nil
	}
}

func testAccCheckAciActionRuleProfileMultipathDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_action_rule_profile" {
			cont, err := client.Get(rs.Primary.ID)
			action_rule_profile := models.ActionRuleProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Action Rule Profile %s Still exists", action_rule_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciActionRuleProfileMultipathAttributes(fv_tenant_name, rtctrl_attr_p_name, description string,
	action_rule_profile *models.ActionRuleProfile,
	next_hop_propagation *models.NexthopUnchangedAction,
	multipath *models.RedistributeMultipathAction,

) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if rtctrl_attr_p_name != GetMOName(action_rule_profile.DistinguishedName) {
			return fmt.Errorf("Bad rtctrl_attr_p %s", GetMOName(action_rule_profile.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(action_rule_profile.DistinguishedName, fmt.Sprintf("/attr-%s", rtctrl_attr_p_name))) {
			return fmt.Errorf(" Bad fv_tenant %s", GetMOName(GetParentDn(action_rule_profile.DistinguishedName, fmt.Sprintf("/attr-%s", rtctrl_attr_p_name))))
		}
		if description != action_rule_profile.Description {
			return fmt.Errorf("Bad action_rule_profile Description %s", action_rule_profile.Description)
		}
		if "nh-unchanged" != next_hop_propagation.Type {
			return fmt.Errorf("Bad Next Hop Propagation Type value %s", next_hop_propagation.Type)
		}

		if "redist-multipath" != multipath.Type {
			return fmt.Errorf("Bad Multipath Type value %s", multipath.Type)
		}
		return nil
	}
}
