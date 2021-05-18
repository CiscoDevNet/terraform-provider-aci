package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAciLogicalInterfaceContext() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceAciLogicalInterfaceContextRead,

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

func dataSourceAciLogicalInterfaceContextRead(d *schema.ResourceData, m interface{}) error {
	aciClient := m.(*client.Client)

	connNameOrLbl := d.Get("conn_name_or_lbl").(string)

	rn := fmt.Sprintf("lIfCtx-c-%s", connNameOrLbl)
	LogicalDeviceContextDn := d.Get("logical_device_context_dn").(string)

	dn := fmt.Sprintf("%s/%s", LogicalDeviceContextDn, rn)

	vnsLIfCtx, err := getRemoteLogicalInterfaceContext(aciClient, dn)

	if err != nil {
		return err
	}
	d.SetId(dn)
	setLogicalInterfaceContextAttributes(vnsLIfCtx, d)
	return nil
}
