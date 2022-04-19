package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAciActionRuleProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciActionRuleProfileRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"set_route_tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"set_preference": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"set_weight": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"set_metric": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"set_metric_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ospf-type1",
					"ospf-type2",
				}, false),
			},
		})),
	}
}

func dataSourceAciActionRuleProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("attr-%s", name)
	TenantDn := d.Get("tenant_dn").(string)

	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	rtctrlAttrP, err := getRemoteActionRuleProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setActionRuleProfileAttributes(rtctrlAttrP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// rtctrlSetTag - Beginning of Read
	setRouteTagDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetTag)
	rtctrlSetTag, err := getRemoteRtctrlSetTag(aciClient, setRouteTagDn)
	if err == nil {
		_, err = setRtctrlSetTagAttributes(rtctrlSetTag, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetTag - Read finished successfully

	// rtctrlSetPref - Beginning of Read
	setPrefDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetPref)
	rtctrlSetPref, err := getRemoteRtctrlSetPref(aciClient, setPrefDn)
	if err == nil {
		_, err = setRtctrlSetPrefAttributes(rtctrlSetPref, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetPref - Read finished successfully

	// rtctrlSetWeight - Beginning of Read
	setWeightCheckDns := make([]string, 0, 1)

	setWeightDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetWeight)

	setWeightCheckDns = append(setWeightCheckDns, setWeightDn)

	err = checkTDn(aciClient, setWeightCheckDns)
	if err == nil {

		rtctrlSetWeight, err := getRemoteRtctrlSetWeight(aciClient, setWeightDn)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = setRtctrlSetWeightAttributes(rtctrlSetWeight, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetWeight - Read finished successfully

	// rtctrlSetRtMetric - Beginning of Read
	setRtMetricCheckDns := make([]string, 0, 1)

	setRtMetricDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRtMetric)

	setRtMetricCheckDns = append(setRtMetricCheckDns, setRtMetricDn)

	err = checkTDn(aciClient, setRtMetricCheckDns)
	if err == nil {

		rtctrlSetRtMetric, err := getRemoteRtctrlSetRtMetric(aciClient, setRtMetricDn)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = setRtctrlSetRtMetricAttributes(rtctrlSetRtMetric, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetRtMetric - Read finished successfully

	// rtctrlSetRtMetricType - Beginning of Read
	setRtMetricTypeCheckDns := make([]string, 0, 1)

	setRtMetricTypeDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRtMetricType)

	setRtMetricTypeCheckDns = append(setRtMetricTypeCheckDns, setRtMetricTypeDn)

	err = checkTDn(aciClient, setRtMetricTypeCheckDns)
	if err == nil {

		rtctrlSetRtMetricType, err := getRemoteRtctrlSetRtMetricType(aciClient, setRtMetricTypeDn)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = setRtctrlSetRtMetricTypeAttributes(rtctrlSetRtMetricType, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetRtMetricType - Read finished successfully

	return nil
}
