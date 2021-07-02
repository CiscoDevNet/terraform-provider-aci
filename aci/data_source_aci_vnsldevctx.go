package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLogicalDeviceContext() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciLogicalDeviceContextRead,

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ctrct_name_or_lbl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"graph_name_or_lbl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"node_name_or_lbl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"context": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func dataSourceAciLogicalDeviceContextRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	ctrctNameOrLbl := d.Get("ctrct_name_or_lbl").(string)

	graphNameOrLbl := d.Get("graph_name_or_lbl").(string)

	nodeNameOrLbl := d.Get("node_name_or_lbl").(string)

	rn := fmt.Sprintf("ldevCtx-c-%s-g-%s-n-%s", ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
	TenantDn := d.Get("tenant_dn").(string)

	dn := fmt.Sprintf("%s/%s", TenantDn, rn)

	vnsLDevCtx, err := getRemoteLogicalDeviceContext(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setLogicalDeviceContextAttributes(vnsLDevCtx, d)
	return nil
}
