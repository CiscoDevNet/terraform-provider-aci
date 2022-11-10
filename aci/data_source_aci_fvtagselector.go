package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciEndpointSecurityGroupTagSelector() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciEndpointSecurityGroupTagSelectorRead,
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
			"match_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"match_value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value_operator": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciEndpointSecurityGroupTagSelectorRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	matchKey := d.Get("match_key").(string)
	matchValue := d.Get("match_value").(string)
	EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)
	rn := fmt.Sprintf("tagselectorkey-[%s]-value-[%s]", matchKey, matchValue)
	dn := fmt.Sprintf("%s/%s", EndpointSecurityGroupDn, rn)
	fvTagSelector, err := getRemoteEndpointSecurityGroupTagSelector(aciClient, dn)
	if err != nil {
		return err
	}
	d.SetId(dn)
	setEndpointSecurityGroupTagSelectorAttributes(fvTagSelector, d)
	return nil
}
