package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciVLANPool() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciVLANPoolRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"alloc_mode": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciVLANPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	allocmode := d.Get("alloc_mode").(string)

	rn := fmt.Sprintf("infra/vlanns-[%s]-%s", name, allocmode)

	dn := fmt.Sprintf("uni/%s", rn)

	fvnsVlanInstP, err := getRemoteVLANPool(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	setVLANPoolAttributes(fvnsVlanInstP, d)
	return nil
}
