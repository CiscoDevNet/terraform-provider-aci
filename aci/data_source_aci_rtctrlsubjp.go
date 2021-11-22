package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciMatchRule() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciMatchRuleRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": {
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
		})),
	}
}

func dataSourceAciMatchRuleRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)
	rn := fmt.Sprintf("subj-%s", name)
	dn := fmt.Sprintf("%s/%s", TenantDn, rn)
	rtctrlSubjP, err := getRemoteMatchRule(aciClient, dn)
	if err != nil {
		return err
	}
	d.SetId(dn)
	setMatchRuleAttributes(rtctrlSubjP, d)
	return nil
}
