package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLogicalInterfaceContext() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciLogicalInterfaceContextRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_device_context_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"conn_name_or_lbl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"l3_dest": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"permit_log": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciLogicalInterfaceContextRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	connNameOrLbl := d.Get("conn_name_or_lbl").(string)

	rn := fmt.Sprintf("lIfCtx-c-%s", connNameOrLbl)
	LogicalDeviceContextDn := d.Get("logical_device_context_dn").(string)

	dn := fmt.Sprintf("%s/%s", LogicalDeviceContextDn, rn)

	vnsLIfCtx, err := getRemoteLogicalInterfaceContext(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setLogicalInterfaceContextAttributes(vnsLIfCtx, d)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
