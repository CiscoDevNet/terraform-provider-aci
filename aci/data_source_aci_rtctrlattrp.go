package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			},
			"set_next_hop": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"set_communities": {
				Optional: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"next_hop_propagation": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"multipath": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"set_as_path_prepend_last_as": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"set_as_path_prepend_as": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"order": {
							Optional: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"set_dampening": {
				Optional: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	setWeightDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetWeight)
	rtctrlSetWeight, err := getRemoteRtctrlSetWeight(aciClient, setWeightDn)
	if err == nil {
		_, err = setRtctrlSetWeightAttributes(rtctrlSetWeight, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetWeight - Read finished successfully

	// rtctrlSetRtMetric - Beginning of Read
	setRtMetricDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRtMetric)
	rtctrlSetRtMetric, err := getRemoteRtctrlSetRtMetric(aciClient, setRtMetricDn)
	if err == nil {
		_, err = setRtctrlSetRtMetricAttributes(rtctrlSetRtMetric, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetRtMetric - Read finished successfully

	// rtctrlSetRtMetricType - Beginning of Read
	setRtMetricTypeDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRtMetricType)
	rtctrlSetRtMetricType, err := getRemoteRtctrlSetRtMetricType(aciClient, setRtMetricTypeDn)
	if err == nil {
		_, err = setRtctrlSetRtMetricTypeAttributes(rtctrlSetRtMetricType, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetRtMetricType - Read finished successfully

	// rtctrlSetNh - Beginning of Read
	setNhDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetNh)
	rtctrlSetNh, err := getRemoteRtctrlSetNh(aciClient, setNhDn)
	if err == nil {
		_, err = setRtctrlSetNhAttributes(rtctrlSetNh, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetNh - Read finished successfully

	// rtctrlSetComm - Beginning of Read
	setCommDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetComm)
	rtctrlSetComm, err := getRemoteRtctrlSetComm(aciClient, setCommDn)
	if err == nil {
		_, err = setRtctrlSetCommAttributes(rtctrlSetComm, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetComm - Read finished successfully

	// rtctrlSetNhUnchanged - Beginning of Read
	setNhUnchangedDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetNhUnchanged)
	rtctrlSetNhUnchanged, err := getRemoteNexthopUnchangedAction(aciClient, setNhUnchangedDn)
	if err == nil {
		_, err = setNexthopUnchangedActionAttributes(rtctrlSetNhUnchanged, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetNhUnchanged - Read finished successfully

	// rtctrlSetRedistMultipath - Beginning of Read
	setRedistMultipathDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetRedistMultipath)
	rtctrlSetRedistMultipath, err := getRemoteRtctrlSetRedistMultipath(aciClient, setRedistMultipathDn)
	if err == nil {
		_, err = setRtctrlSetRedistMultipathAttributes(rtctrlSetRedistMultipath, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetRedistMultipath - Read finished successfully

	// rtctrlSetASPath - Beginning of Read

	setASPathPrependLastASDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetASPath, "prepend-last-as")

	rtctrlSetASPathLastAS, err := getRemoteRtctrlSetASPath(aciClient, setASPathPrependLastASDn)
	if err == nil {
		_, err = setRtctrlSetASPathAttributes(rtctrlSetASPathLastAS, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetASPath - Read finished successfully

	// rtctrlSetASPathASN - Beginning of Read
	setASNumberDn := rtctrlAttrP.DistinguishedName + "/" + fmt.Sprintf(models.RnrtctrlSetASPath, "prepend")
	_, err = getAndSetRemoteSetASPathASNs(aciClient, setASNumberDn, d)
	if err != nil {
		return diag.FromErr(err)
	}
	// rtctrlSetASPathASN - Read finished successfully

	// rtctrlSetDamp - Beginning of Read
	setDampeningDn := rtctrlAttrP.DistinguishedName + fmt.Sprintf("/"+models.RnrtctrlSetDamp)
	rtctrlSetDamp, err := getRemoteRtctrlSetDamp(aciClient, setDampeningDn)
	if err == nil {
		_, err = setRtctrlSetDampAttributes(rtctrlSetDamp, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// rtctrlSetDamp - Read finished successfully

	return nil
}
