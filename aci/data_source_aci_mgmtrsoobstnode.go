package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAciMgmtStaticNode() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciAciMgmtStaticNodeRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"management_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"t_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in_band",
					"out_of_band",
				}, false),
			},

			"addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"gw": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"v6_addr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"v6_gw": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciAciMgmtStaticNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	tDn := d.Get("t_dn").(string)

	bandType := d.Get("type").(string)

	if bandType == "in_band" {
		rn := fmt.Sprintf("rsinBStNode-[%s]", tDn)
		managementEPgDn := d.Get("management_epg_dn").(string)

		dn := fmt.Sprintf("%s/%s", managementEPgDn, rn)

		mgmtRsInBStNode, err := getRemoteInbandStaticNode(aciClient, dn)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(dn)
		_, err = setMgmtStaticNodeAttributes(nil, mgmtRsInBStNode, "in_band", d)
		if err != nil {
			return diag.FromErr(err)
		}

	} else {
		rn := fmt.Sprintf("rsooBStNode-[%s]", tDn)
		managementEPgDn := d.Get("management_epg_dn").(string)

		dn := fmt.Sprintf("%s/%s", managementEPgDn, rn)

		mgmtRsOoBStNode, err := getRemoteOutofbandStaticNode(aciClient, dn)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(dn)
		_, err = setMgmtStaticNodeAttributes(mgmtRsOoBStNode, nil, "out_of_band", d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}
