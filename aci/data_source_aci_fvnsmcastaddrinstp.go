package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciMulticastAddressPool() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciMulticastAddressPoolRead,
		SchemaVersion: 1,
		Schema: AppendAttrSchemas(map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"multicast_address_block": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: AppendAttrSchemas(map[string]*schema.Schema{
						"dn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"from": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"to": {
							Type:     schema.TypeString,
							Required: true,
						},
					}, GetBaseAttrSchema(), GetNameAliasAttrSchema()),
				},
			},
		}, GetBaseAttrSchema(), GetNameAliasAttrSchema()),
	}
}

func dataSourceAciMulticastAddressPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)

	rn := fmt.Sprintf(models.RnfvnsMcastAddrInstP, name)
	dn := fmt.Sprintf("%s/%s", models.ParentDnfvnsMcastAddrInstP, rn)

	fvnsMcastAddrInstP, err := getRemoteMulticastAddressPool(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	_, err = setMulticastAddressPoolAttributes(fvnsMcastAddrInstP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	multicastAddressBlocks := getMulticastAddressBlocks("Read", fvnsMcastAddrInstP.Name, aciClient, d)
	if multicastAddressBlocks != nil {
		setMulticastAddressBlocks("Read", fvnsMcastAddrInstP.Name, multicastAddressBlocks, aciClient, d)
	}

	return nil
}
