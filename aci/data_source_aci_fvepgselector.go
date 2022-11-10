package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciEndpointSecurityGroupEPgSelector() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciEndpointSecurityGroupEPgSelectorRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"endpoint_security_group_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"match_epg_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciEndpointSecurityGroupEPgSelectorRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	matchEpgDn := d.Get("match_epg_dn").(string)
	EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)
	rn := fmt.Sprintf("epgselector-[%s]", matchEpgDn)
	dn := fmt.Sprintf("%s/%s", EndpointSecurityGroupDn, rn)
	fvEPgSelector, err := getRemoteEndpointSecurityGroupEPgSelector(aciClient, dn)
	if err != nil {
		return err
	}
	d.SetId(dn)
	setEndpointSecurityGroupEPgSelectorAttributes(fvEPgSelector, d)
	return nil
}
