package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciStaticPath() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciStaticPathRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"application_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instr_imedcy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"primary_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciStaticPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	tDn := d.Get("tdn").(string)

	rn := fmt.Sprintf("rspathAtt-[%s]", tDn)
	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	dn := fmt.Sprintf("%s/%s", ApplicationEPGDn, rn)

	fvRsPathAtt, err := getRemoteStaticPath(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setStaticPathAttributes(fvRsPathAtt, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
