package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciSPANSourcedestinationGroupMatchLabel() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciSPANSourcedestinationGroupMatchLabelRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"span_source_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciSPANSourcedestinationGroupMatchLabelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("spanlbl-%s", name)
	SPANSourceGroupDn := d.Get("span_source_group_dn").(string)

	dn := fmt.Sprintf("%s/%s", SPANSourceGroupDn, rn)

	spanSpanLbl, err := getRemoteSPANSourcedestinationGroupMatchLabel(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	setSPANSourcedestinationGroupMatchLabelAttributes(spanSpanLbl, d)
	return nil
}
