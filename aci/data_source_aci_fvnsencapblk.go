package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciRanges() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciRangesRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"vlan_pool_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"from": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"to": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"alloc_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciRangesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	from := d.Get("from").(string)

	to := d.Get("to").(string)

	rn := fmt.Sprintf("from-[%s]-to-[%s]", from, to)
	VLANPoolDn := d.Get("vlan_pool_dn").(string)

	dn := fmt.Sprintf("%s/%s", VLANPoolDn, rn)

	fvnsEncapBlk, err := getRemoteRanges(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)
	_, err = setRangesAttributes(fvnsEncapBlk, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
