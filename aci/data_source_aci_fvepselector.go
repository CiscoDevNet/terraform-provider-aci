package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciEndpointSecurityGroupSelector() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciEndpointSecurityGroupSelectorRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"endpoint_security_group_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"match_expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func dataSourceAciEndpointSecurityGroupSelectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	matchExpression := d.Get("match_expression").(string)
	EndpointSecurityGroupDn := d.Get("endpoint_security_group_dn").(string)
	rn := fmt.Sprintf("epselector-[%s]", matchExpression)
	dn := fmt.Sprintf("%s/%s", EndpointSecurityGroupDn, rn)
	fvEPSelector, err := getRemoteEndpointSecurityGroupSelector(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setEndpointSecurityGroupSelectorAttributes(fvEPSelector, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
