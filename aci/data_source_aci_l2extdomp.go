package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciL2Domain() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciL2DomainRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAciL2DomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("l2dom-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)

	l2extDomP, err := getRemoteL2Domain(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setL2DomainAttributes(l2extDomP, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
