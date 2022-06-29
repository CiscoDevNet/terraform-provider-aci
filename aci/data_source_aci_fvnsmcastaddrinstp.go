package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciMulticastAddressPool() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciMulticastAddressPoolRead,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func dataSourceAciMulticastAddressPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/maddrns-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)

	fvnsMcastAddrInstP, err := getRemoteMulticastAddressPool(aciClient, dn)
	if err != nil {
		return nil
	}

	d.SetId(dn)

	_, err = setMulticastAddressPoolAttributes(fvnsMcastAddrInstP, d)
	if err != nil {
		return nil
	}

	return nil
}
