package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciEndpointSecurityGroupSelector() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceAciEndpointSecurityGroupSelectorRead,
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
			"match_expression": {
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

func dataSourceAciEndpointSecurityGroupSelectorRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)
	matchExpression := d.Get("matchExpression").(string)
	EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)
	rn := fmt.Sprintf("epselector-[%s]", matchExpression)
	dn := fmt.Sprintf("%s/%s", EndpointSecurityGroupDn, rn)
	fvEPSelector, err := getRemoteEndpointSecurityGroupSelector(aciClient, dn)
	if err != nil {
		return err
	}
	d.SetId(dn)
	setEndpointSecurityGroupSelectorAttributes(fvEPSelector, d)
	return nil
}
