package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciMatchCommunityTerm() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciMatchCommunityTermRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"match_rule_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"match_community_factors": {
				Type:     schema.TypeString,
				Optional: true,
			},
		})),
	}
}

func dataSourceAciMatchCommunityTermRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	matchRuleDn := d.Get("match_rule_dn").(string)
	rn := fmt.Sprintf(models.RnrtctrlMatchCommTerm, name)
	dn := fmt.Sprintf("%s/%s", matchRuleDn, rn)

	rtctrlMatchCommTerm, err := getRemoteMatchCommunityTerm(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setMatchCommunityTermAttributes(rtctrlMatchCommTerm, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
