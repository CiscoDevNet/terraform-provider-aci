package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
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

	return nil
}
