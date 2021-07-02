package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAny() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciAnyRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"vrf_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"match_t": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pref_gr_memb": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciAnyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("any")
	VRFDn := d.Get("vrf_dn").(string)

	dn := fmt.Sprintf("%s/%s", VRFDn, rn)

	vzAny, err := getRemoteAny(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setAnyAttributes(vzAny, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
