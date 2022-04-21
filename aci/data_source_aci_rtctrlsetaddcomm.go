package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciRtctrlSetAddComm() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciRtctrlSetAddCommRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"action_rule_profile_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"community": {
				Type:     schema.TypeString,
				Required: true,
			},
			"set_criteria": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciRtctrlSetAddCommRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	community := d.Get("community").(string)
	ActionRuleProfileDn := d.Get("action_rule_profile_dn").(string)
	rn := fmt.Sprintf(models.RnrtctrlSetAddComm, community)
	dn := fmt.Sprintf("%s/%s", ActionRuleProfileDn, rn)

	rtctrlSetAddComm, err := getRemoteRtctrlSetAddComm(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setRtctrlSetAddCommAttributes(rtctrlSetAddComm, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
