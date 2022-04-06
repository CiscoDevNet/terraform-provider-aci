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
			"match_community_factors": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "Create Community Factors",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"transitive",
								"non-transitive",
							}, false),
						},
						"community": {
							Required: true,
							Type:     schema.TypeString,
						},
						"description": {
							Optional: true,
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
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

	rtctrlMatchCommFactors, err := aciClient.ListMatchCommFactorsFromCommunityTerm(dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)

	matchCommFactors := make([]map[string]string, 0, 1)

	for _, factor := range rtctrlMatchCommFactors {
		factorSet, err := setMatchCommunityFactorAttributes(factor, make(map[string]string))
		if err != nil {
			return diag.FromErr(err)
		}
		matchCommFactors = append(matchCommFactors, factorSet)
	}
	d.Set("match_community_factors", matchCommFactors)

	return nil
}
