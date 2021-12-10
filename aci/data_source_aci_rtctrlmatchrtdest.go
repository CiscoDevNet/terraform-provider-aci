package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciMatchRouteDestinationRule() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciMatchRouteDestinationRuleRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"match_rule_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aggregate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"greater_than_mask": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"less_than_mask": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciMatchRouteDestinationRuleRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	ip := d.Get("ip").(string)
	MatchRuleDn := d.Get("match_rule_dn").(string)
	rn := fmt.Sprintf("dest-[%s]", ip)
	dn := fmt.Sprintf("%s/%s", MatchRuleDn, rn)
	rtctrlMatchRtDest, err := getRemoteMatchRouteDestinationRule(aciClient, dn)
	if err != nil {
		return err
	}
	d.SetId(dn)
	setMatchRouteDestinationRuleAttributes(rtctrlMatchRtDest, d)
	return nil
}
