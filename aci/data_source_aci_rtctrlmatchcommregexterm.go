package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciMatchRuleBasedonCommunityRegularExpression() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciMatchRuleBasedonCommunityRegularExpressionRead,
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
			"community_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"regex": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciMatchRuleBasedonCommunityRegularExpressionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	commType := d.Get("community_type").(string)
	MatchRuleDn := d.Get("match_rule_dn").(string)
	rn := fmt.Sprintf(models.RnrtctrlMatchCommRegexTerm, commType)
	dn := fmt.Sprintf("%s/%s", MatchRuleDn, rn)

	rtctrlMatchCommRegexTerm, err := getRemoteMatchRuleBasedonCommunityRegularExpression(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setMatchRuleBasedonCommunityRegularExpressionAttributes(rtctrlMatchCommRegexTerm, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
