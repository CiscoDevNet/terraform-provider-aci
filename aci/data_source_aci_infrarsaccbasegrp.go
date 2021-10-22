package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciAccessGroup() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciAccessAccessGroupRead,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"access_port_selector_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"fex_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAciAccessAccessGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	rn := fmt.Sprintf("rsaccBaseGrp")
	AccessPortSelectorDn := d.Get("access_port_selector_dn").(string)

	dn := fmt.Sprintf("%s/%s", AccessPortSelectorDn, rn)

	infraRsAccBaseGrp, err := getRemoteAccessAccessGroup(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setAccessAccessGroupAttributes(infraRsAccBaseGrp, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
